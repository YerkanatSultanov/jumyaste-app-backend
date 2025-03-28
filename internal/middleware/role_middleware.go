package middleware

import (
	"github.com/gin-gonic/gin"
	"jumyste-app-backend/pkg/logger"
	"log/slog"
	"net/http"
)

func RequireRole(requiredRole int) gin.HandlerFunc {
	return func(c *gin.Context) {
		roleID, exists := c.Get("role_id")
		if !exists {
			logger.Log.Warn("Role ID missing in context")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		if roleID.(int) != requiredRole {
			logger.Log.Warn("Unauthorized access attempt", slog.Int("role_id", roleID.(int)))
			c.JSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions"})
			c.Abort()
			return
		}

		c.Next()
	}
}
