package handler

import (
	"github.com/gin-gonic/gin"
	"jumyste-app-backend/internal/dto"
	_ "jumyste-app-backend/internal/entity"
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

// CreateChatHandler godoc
// @Summary Create a chat between two users
// @Description Create a new chat by providing the second user's ID
// @Tags Chats
// @Accept json
// @Produce json
// @Param request body dto.CreateChatRequest true "Chat creation payload"
// @Security BearerAuth
// @Success 201 {object} entity.Chat "Chat created successfully"
// @Failure 400 {object} dto.ErrorResponse "Invalid request"
// @Failure 401 {object} dto.ErrorResponse "Unauthorized"
// @Failure 500 {object} dto.ErrorResponse "Failed to create chat"
// @Router /chats [post]
func (h *ChatHandler) CreateChatHandler(c *gin.Context) {
	var req dto.CreateChatRequest

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

	chat, err := h.ChatService.CreateChat(id, req.SecondUserId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, chat)
}

// GetChatByIDHandler godoc
// @Summary Get chat by ID
// @Description Retrieve a chat by its ID
// @Tags Chats
// @Accept json
// @Produce json
// @Param chatID path int true "Chat ID"
// @Security BearerAuth
// @Success 200 {object} entity.Chat "Chat found"
// @Failure 400 {object} dto.ErrorResponse "Invalid chat ID"
// @Failure 404 {object} dto.ErrorResponse "Chat not found"
// @Failure 401 {object} dto.ErrorResponse "Unauthorized"
// @Failure 500 {object} dto.ErrorResponse "Internal server error"
// @Router /chats/{chatID} [get]
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

// GetAllChatsHandler godoc
// @Summary Get all chats
// @Description Retrieve all chats
// @Tags Chats
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {array} entity.Chat "List of all chats"
// @Failure 401 {object} dto.ErrorResponse "Unauthorized"
// @Failure 500 {object} dto.ErrorResponse "Internal server error"
// @Router /chats [get]
func (h *ChatHandler) GetAllChatsHandler(c *gin.Context) {
	chats, err := h.ChatService.GetAllChats()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch chats"})
		return
	}

	c.JSON(http.StatusOK, chats)
}

// GetChatsByUserIDHandler godoc
// @Summary Get chats by user ID
// @Description Retrieve all chats for a specific user by their user ID
// @Tags Chats
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {array} entity.Chat "List of chats for the user"
// @Failure 401 {object} dto.ErrorResponse "Unauthorized"
// @Failure 500 {object} dto.ErrorResponse "Internal server error"
// @Router /chats/user [get]
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
