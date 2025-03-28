package manager

import (
	"github.com/gorilla/websocket"
	"log"
	"sync"
)

// Client represents a WebSocket connection
type Client struct {
	Conn *websocket.Conn
	Send chan []byte
}

// WebSocketManager handles all active connections
type WebSocketManager struct {
	Clients    map[*Client]bool
	Broadcast  chan []byte
	Register   chan *Client
	Unregister chan *Client
	mu         sync.Mutex
}

// NewWebSocketManager initializes a WebSocket manager
func NewWebSocketManager() *WebSocketManager {
	return &WebSocketManager{
		Clients:    make(map[*Client]bool),
		Broadcast:  make(chan []byte),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
	}
}

// Run starts the WebSocket manager
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
		}
	}
}

// HandleClient manages an individual WebSocket connection
func (manager *WebSocketManager) HandleClient(client *Client) {
	defer func() {
		manager.Unregister <- client
		client.Conn.Close()
	}()

	for {
		_, message, err := client.Conn.ReadMessage()
		if err != nil {
			log.Println("Error reading:", err)
			break
		}
		manager.Broadcast <- message
	}
}

// WriteMessages sends messages to the client
func (client *Client) WriteMessages() {
	defer client.Conn.Close()
	for message := range client.Send {
		if err := client.Conn.WriteMessage(websocket.TextMessage, message); err != nil {
			break
		}
	}
}
