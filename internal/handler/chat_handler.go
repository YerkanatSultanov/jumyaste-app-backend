package handler

import (
	"github.com/gin-gonic/gin"
	"jumyste-app-backend/internal/dto"
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
		SecondUserId int `json:"second_user_id"`
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: "Unauthorized"})
		return
	}

	id, ok := userID.(int)
	if !ok {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "Invalid user ID type"})
		return
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	chat, err := h.ChatService.CreateChat(req.SecondUserId, id)
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

func (h *ChatHandler) GetChatsByUserIDHandler(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	uid, ok := userID.(int)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID"})
		return
	}

	chats, err := h.ChatService.GetChatsByUserID(uid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch chats"})
		return
	}

	c.JSON(http.StatusOK, chats)
}
