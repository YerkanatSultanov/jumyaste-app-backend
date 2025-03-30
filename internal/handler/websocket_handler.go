package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"jumyste-app-backend/internal/manager"
	"jumyste-app-backend/internal/middleware"
	"jumyste-app-backend/pkg/logger"
	"net/http"
	"strconv"
	"strings"
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

func (h *WebSocketHandler) HandleWebSocket(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing token"})
		return
	}

	if !strings.HasPrefix(authHeader, "Bearer ") {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token format"})
		return
	}

	tokenString := strings.TrimPrefix(authHeader, "Bearer ")

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
