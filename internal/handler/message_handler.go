package handler

import (
	"github.com/gin-gonic/gin"
	"jumyste-app-backend/internal/entity"
	"jumyste-app-backend/internal/service"
	"net/http"
	"strconv"
)

type MessageHandler struct {
	MessageService *service.MessageService
}

func NewMessageHandler(messageService *service.MessageService) *MessageHandler {
	return &MessageHandler{MessageService: messageService}
}

func (h *MessageHandler) SendMessageHandler(c *gin.Context) {
	var req struct {
		ChatID  int                `form:"chat_id"`
		Type    entity.MessageType `form:"type"`
		Content *string            `form:"content,omitempty"`
		FileURL *string            `form:"file_url,omitempty"`
	}

	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	senderID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	sender, ok := senderID.(int)
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "Invalid user ID"})
		return
	}

	// Handle file upload
	var fileURL *string
	file, err := c.FormFile("file")
	if err == nil {
		filePath := "uploads/" + file.Filename
		if err := c.SaveUploadedFile(file, filePath); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
			return
		}
		fileURL = &filePath
	}

	message, err := h.MessageService.SendMessage(req.ChatID, sender, req.Type, req.Content, fileURL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send message"})
		return
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
