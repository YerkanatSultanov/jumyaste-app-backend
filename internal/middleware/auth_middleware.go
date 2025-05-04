package middleware

import (
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"jumyste-app-backend/config"
	"jumyste-app-backend/pkg/logger"
	"jumyste-app-backend/utils"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type AuthMiddleware struct {
	secretKey string
}

func NewAuthMiddleware(cfg config.Config) *AuthMiddleware {
	return &AuthMiddleware{secretKey: cfg.JWT.Secret}
}

var ErrInvalidToken = errors.New("invalid token")
var ErrExpiredToken = errors.New("token is expired")

func (m *AuthMiddleware) VerifyTokenMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		claims, err := m.validateToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}

		logger.Log.Info("User authenticated", slog.Int("user_id", claims.UserID), slog.Int("role_id", claims.RoleID))

		c.Set("user_id", claims.UserID)
		c.Set("role_id", claims.RoleID)
		c.Set("company_id", claims.CompanyID)
		c.Set("dep_id", claims.DepartmentID)
		c.Set("claims", claims)
		c.Next()
	}
}

func (m *AuthMiddleware) validateToken(tokenString string) (*utils.Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &utils.Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(m.secretKey), nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrExpiredToken
		}
		return nil, ErrInvalidToken
	}

	claims, ok := token.Claims.(*utils.Claims)
	if !ok || !token.Valid {
		return nil, ErrInvalidToken
	}

	if claims.ExpiresAt.Time.Before(time.Now()) {
		return nil, ErrExpiredToken
	}

	return claims, nil
}

func (m *AuthMiddleware) VerifyTokenWithClaims(tokenString string) (*utils.Claims, error) {
	return m.validateToken(tokenString)
}
