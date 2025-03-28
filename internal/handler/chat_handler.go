package handler

import (
	"github.com/gin-gonic/gin"
	"jumyste-app-backend/internal/service"
	"net/http"
	"strconv"
)

type ChatHandler struct {
	ChatService *service.ChatService
}

func NewChatHandler(chatService *service.ChatService) *ChatHandler {
	return &ChatHandler{ChatService: chatService}
}

// CreateChatHandler - Creates chat with received ID's
func (h *ChatHandler) CreateChatHandler(c *gin.Context) {
	var req struct {
		Users []uint `json:"users"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	chat, err := h.ChatService.CreateChat(req.Users)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, chat)
}

// GetChatByIDHandler - Retrieves a chat by ID
func (h *ChatHandler) GetChatByIDHandler(c *gin.Context) {
	chatID, err := strconv.ParseUint(c.Param("chatID"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid chat ID"})
		return
	}

	chat, err := h.ChatService.GetChatByID(uint(chatID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Chat not found"})
		return
	}

	c.JSON(http.StatusOK, chat)
}

// GetAllChatsHandler - Retrieves all chats
func (h *ChatHandler) GetAllChatsHandler(c *gin.Context) {
	chats, err := h.ChatService.GetAllChats()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch chats"})
		return
	}

	c.JSON(http.StatusOK, chats)
}
