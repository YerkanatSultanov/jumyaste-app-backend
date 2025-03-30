package manager

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"log"
	"sync"
)

// Client представляет соединение WebSocket с пользователем
type Client struct {
	Conn   *websocket.Conn
	Send   chan []byte
	UserID int
	ChatID int
}

// WebSocketManager управляет всеми соединениями
type WebSocketManager struct {
	Clients    map[*Client]bool
	Broadcast  chan []byte
	MarkAsRead chan ReadMessage
	Register   chan *Client
	Unregister chan *Client
	mu         sync.Mutex
}

// ReadMessage структура для события "прочтено"
type ReadMessage struct {
	MessageID int `json:"message_id"`
	UserID    int `json:"user_id"`
	ChatID    int `json:"chat_id"`
}

// NewWebSocketManager инициализирует менеджер WebSocket
func NewWebSocketManager() *WebSocketManager {
	return &WebSocketManager{
		Clients:    make(map[*Client]bool),
		Broadcast:  make(chan []byte),
		MarkAsRead: make(chan ReadMessage),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
	}
}

// Run запускает WebSocket менеджер
func (manager *WebSocketManager) Run() {
	for {
		select {
		case client := <-manager.Register:
			manager.mu.Lock()
			manager.Clients[client] = true
			manager.mu.Unlock()

		case client := <-manager.Unregister:
			manager.mu.Lock()
			if _, ok := manager.Clients[client]; ok {
				delete(manager.Clients, client)
				close(client.Send)
			}
			manager.mu.Unlock()

		case message := <-manager.Broadcast:
			manager.mu.Lock()
			for client := range manager.Clients {
				select {
				case client.Send <- message:
				default:
					close(client.Send)
					delete(manager.Clients, client)
				}
			}
			manager.mu.Unlock()

		case readMsg := <-manager.MarkAsRead:
			manager.mu.Lock()
			// Рассылаем событие "прочитано" всем клиентам в этом чате
			readEvent, _ := json.Marshal(map[string]interface{}{
				"type":       "message_read",
				"message_id": readMsg.MessageID,
				"user_id":    readMsg.UserID,
				"chat_id":    readMsg.ChatID,
			})

			for client := range manager.Clients {
				if client.ChatID == readMsg.ChatID {
					client.Send <- readEvent
				}
			}
			manager.mu.Unlock()
		}
	}
}

// HandleClient обрабатывает сообщения от клиента
func (manager *WebSocketManager) HandleClient(client *Client) {
	defer func() {
		manager.Unregister <- client
		client.Conn.Close()
	}()

	for {
		_, message, err := client.Conn.ReadMessage()
		if err != nil {
			log.Println("Ошибка чтения сообщения:", err)
			break
		}

		// Попробуем распарсить сообщение как "прочитано"
		var readMsg ReadMessage
		if err := json.Unmarshal(message, &readMsg); err == nil && readMsg.MessageID > 0 {
			readMsg.UserID = client.UserID
			readMsg.ChatID = client.ChatID
			manager.MarkAsRead <- readMsg
			continue
		}

		// Если не "прочитано", то отправляем в чат
		manager.mu.Lock()
		for c := range manager.Clients {
			if c.ChatID == client.ChatID {
				c.Send <- message
			}
		}
		manager.mu.Unlock()
	}
}

// WriteMessages отправляет сообщения клиенту
func (client *Client) WriteMessages() {
	defer client.Conn.Close()
	for message := range client.Send {
		if err := client.Conn.WriteMessage(websocket.TextMessage, message); err != nil {
			break
		}
	}
}
