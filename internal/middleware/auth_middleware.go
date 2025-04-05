package middleware

import (
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"jumyste-app-backend/config"
	"jumyste-app-backend/pkg/logger"
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
		c.Set("dep_id", claims.DepID)
		c.Next()
	}
}

func (m *AuthMiddleware) validateToken(tokenString string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(m.secretKey), nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrExpiredToken
		}
		return nil, ErrInvalidToken
	}

	claims, ok := token.Claims.(*CustomClaims)
	if !ok || !token.Valid {
		return nil, ErrInvalidToken
	}

	if claims.ExpiresAt.Before(time.Now()) {
		return nil, ErrExpiredToken
	}

	return claims, nil
}
func (m *AuthMiddleware) VerifyTokenWithClaims(tokenString string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(m.secretKey), nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrExpiredToken
		}
		return nil, ErrInvalidToken
	}

	claims, ok := token.Claims.(*CustomClaims)
	if !ok || !token.Valid {
		return nil, ErrInvalidToken
	}

	if claims.ExpiresAt.Before(time.Now()) {
		return nil, ErrExpiredToken
	}

	return claims, nil
}

type CustomClaims struct {
	UserID    int `json:"user_id"`
	RoleID    int `json:"role_id"`
	CompanyID int `json:"company_id,omitempty"`
	DepID     int `json:"dep_id,omitempty"`
	jwt.RegisteredClaims
}
