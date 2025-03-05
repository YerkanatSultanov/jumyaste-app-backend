package middleware

import (
	"github.com/gin-gonic/gin"
	"jumyste-app-backend/pkg/logger"
	"jumyste-app-backend/utils"
	"log/slog"
	"net/http"
	"strings"
)

func RequireRole(requiredRole int) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			logger.Log.Warn("Authorization header is missing")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization token required"})
			c.Abort()
			return
		}

		parts := strings.Split(tokenString, " ")
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			logger.Log.Warn("Invalid Authorization header format")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token format"})
			c.Abort()
			return
		}

		tokenString = parts[1]

		claims, err := utils.ParseJWT(tokenString)
		if err != nil {
			logger.Log.Warn("Invalid token", slog.String("error", err.Error()))
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		if claims.RoleID != requiredRole {
			logger.Log.Warn("Unauthorized access attempt", slog.Int("user_id", claims.UserID), slog.Int("role_id", claims.RoleID))
			c.JSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions"})
			c.Abort()
			return
		}
		logger.Log.Info("Authorized access")
		c.Set("user_id", claims.UserID)
		c.Set("role_id", claims.RoleID)

		c.Next()
	}
}
