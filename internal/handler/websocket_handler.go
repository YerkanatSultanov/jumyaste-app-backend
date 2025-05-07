package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	_ "jumyste-app-backend/internal/dto"
	"jumyste-app-backend/internal/manager"
	"jumyste-app-backend/internal/middleware"
	"jumyste-app-backend/pkg/logger"
	"net/http"
	"strconv"
)

type WebSocketHandler struct {
	WSManager      *manager.WebSocketManager
	AuthMiddleware *middleware.AuthMiddleware
}

func NewWebSocketHandler(wsManager *manager.WebSocketManager, authMiddleware *middleware.AuthMiddleware) *WebSocketHandler {
	return &WebSocketHandler{
		WSManager:      wsManager,
		AuthMiddleware: authMiddleware,
	}
}

// WebSocket connection upgrader
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// HandleWebSocket godoc
//
// @Summary WebSocket соединение с чатом
// @Description Устанавливает WebSocket-соединение с авторизованным пользователем и chat_id в query
// @Tags WebSocket
// @Produce plain
// @Security BearerAuth
// @Param chat_id query int true "ID чата"
// @Success 101 {string} string "Switching Protocols – WebSocket connection established"
// @Failure 400 {object} dto.ErrorResponse "Invalid chat ID"
// @Failure 401 {object} dto.ErrorResponse "Unauthorized – отсутствует или неверный токен"
// @Failure 500 {object} dto.ErrorResponse "Ошибка при апгрейде соединения"
// @Router /ws [get]
func (h *WebSocketHandler) HandleWebSocket(c *gin.Context) {
	tokenString := c.Query("access_token")
	if tokenString == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing token"})
		return
	}

	claims, err := h.AuthMiddleware.VerifyTokenWithClaims(tokenString)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}

	userID := claims.UserID

	chatIDStr := c.Query("chat_id")
	chatID, err := strconv.Atoi(chatIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid chat ID"})
		return
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		logger.Log.Error("WebSocket upgrade failed", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upgrade connection"})
		return
	}

	client := &manager.Client{
		Conn:   conn,
		Send:   make(chan []byte),
		UserID: userID,
		ChatID: chatID,
	}

	h.WSManager.Register <- client

	go client.WriteMessages()
	h.WSManager.HandleClient(client)
}
