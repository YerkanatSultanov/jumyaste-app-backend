package handler

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"jumyste-app-backend/internal/entity"
	"jumyste-app-backend/internal/manager"
	"jumyste-app-backend/internal/service"
	"jumyste-app-backend/pkg/logger"
	"log/slog"
	"net/http"
	"strconv"
)

type MessageHandler struct {
	MessageService *service.MessageService
	WSManager      *manager.WebSocketManager
}

func NewMessageHandler(messageService *service.MessageService, wsManager *manager.WebSocketManager) *MessageHandler {
	return &MessageHandler{
		MessageService: messageService,
		WSManager:      wsManager,
	}
}

func (h *MessageHandler) SendMessageHandler(c *gin.Context) {
	chatID, err := strconv.Atoi(c.PostForm("chat_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid chat ID"})
		return
	}

	messageType := c.PostForm("type")
	content := c.PostForm("content")
	fileData := c.PostForm("file_data")

	senderID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	sender, ok := senderID.(int)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID"})
		return
	}

	var fileURL *string
	if fileData != "" {
		fileURL = &fileData
	}

	message, err := h.MessageService.SendMessage(chatID, sender, entity.MessageType(messageType), &content, fileURL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send message"})
		return
	}

	if content != "" {
		h.WSManager.Broadcast <- []byte(content)
	}

	c.JSON(http.StatusCreated, message)
}

// GetMessagesByChatIDHandler - Retrieves all messages in a chat
func (h *MessageHandler) GetMessagesByChatIDHandler(c *gin.Context) {
	chatID, err := strconv.ParseUint(c.Param("chatID"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid chat ID"})
		return
	}

	messages, err := h.MessageService.GetMessagesByChatID(int(chatID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch messages"})
		return
	}

	c.JSON(http.StatusOK, messages)
}

// GetMessageByIDHandler - Retrieves a specific message by ID
func (h *MessageHandler) GetMessageByIDHandler(c *gin.Context) {
	messageID, err := strconv.ParseUint(c.Param("messageID"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid message ID"})
		return
	}

	message, err := h.MessageService.GetMessageByID(int(messageID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Message not found"})
		return
	}

	c.JSON(http.StatusOK, message)
}

func (h *MessageHandler) MarkAsRead(c *gin.Context) {
	userID := c.GetInt("user_id")
	messageID, err := strconv.Atoi(c.Query("message_id"))
	if err != nil {
		logger.Log.Error("Invalid message ID", slog.String("error", err.Error()))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid message ID"})
		return
	}

	// Отмечаем сообщение как прочитанное в БД
	err = h.MessageService.MarkMessageAsRead(c.Request.Context(), messageID, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to mark message as read"})
		return
	}

	// Получаем информацию о сообщении (например, ChatID)
	message, err := h.MessageService.GetMessageByID(messageID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve message info"})
		return
	}

	readEvent := gin.H{
		"type":       "message_read",
		"message_id": messageID,
		"user_id":    userID,
		"chat_id":    message.ChatID,
	}

	h.WSManager.Broadcast <- toJSON(readEvent)

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func toJSON(data interface{}) []byte {
	jsonData, _ := json.Marshal(data)
	return jsonData
}
