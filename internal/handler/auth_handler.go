package handler

import (
	"errors"
	"github.com/gin-gonic/gin"
	"jumyste-app-backend/internal/entity"
	"jumyste-app-backend/internal/service"
	"jumyste-app-backend/pkg/logger"
	"log/slog"
	"net/http"
)

type AuthHandler struct {
	AuthService *service.AuthService
}

func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{AuthService: authService}
}

//func (h *AuthHandler) Register(c *gin.Context) {
//	var user entity.User
//	if err := c.ShouldBindJSON(&user); err != nil {
//		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
//		return
//	}
//
//	if err := h.AuthService.RegisterUser(user.Email, user.Password, user.FirstName, user.LastName); err != nil {
//		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register"})
//		return
//	}
//
//	c.JSON(http.StatusCreated, gin.H{"message": "User registered successfully"})
//}

func (h *AuthHandler) Login(c *gin.Context) {
	var credentials struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&credentials); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	token, err := h.AuthService.LoginUser(credentials.Email, credentials.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}

func (h *AuthHandler) RequestPasswordReset(c *gin.Context) {
	var req struct {
		Email string `json:"email"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	if req.Email == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email is required"})
		return
	}

	err := h.AuthService.RequestPasswordReset(req.Email)
	if err != nil {
		logger.Log.Error("Failed to request password reset", slog.String("email", req.Email), slog.String("error", err.Error()))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to request password reset"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Reset code sent to your email"})
}

func (h *AuthHandler) ResetPassword(c *gin.Context) {
	var req struct {
		Email           string `json:"email"`
		ResetCode       string `json:"reset_code"`
		NewPassword     string `json:"new_password"`
		ConfirmPassword string `json:"confirm_password"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	if req.Email == "" || req.ResetCode == "" || req.NewPassword == "" || req.ConfirmPassword == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "All fields are required"})
		return
	}

	err := h.AuthService.ResetPassword(req.Email, req.ResetCode, req.NewPassword, req.ConfirmPassword)
	if err != nil {
		if errors.Is(err, service.ErrInvalidResetCode) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid or expired reset code"})
			return
		}
		if errors.Is(err, service.ErrUserNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}
		logger.Log.Error("Failed to reset password", slog.String("email", req.Email), slog.String("error", err.Error()))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to reset password"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Password reset successful"})
}

func (h *AuthHandler) Register(c *gin.Context) {
	var request struct {
		Email     string `json:"email"`
		Password  string `json:"password"`
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		RoleId    int    `json:"role_id"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}
	if request.RoleId == 0 {
		request.RoleId = 3
	}

	user := &entity.User{
		Email:     request.Email,
		Password:  request.Password,
		FirstName: request.FirstName,
		LastName:  request.LastName,
		RoleID:    request.RoleId,
	}

	if err := h.AuthService.RegisterUser(user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User registered successfully"})
}

//func (h *AuthHandler) VerifyCodeAndRegister(c *gin.Context) {
//	var request struct {
//		Email string `json:"email"`
//		Code  string `json:"code"`
//	}
//
//	if err := c.ShouldBindJSON(&request); err != nil {
//		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
//		return
//	}
//
//	if err := h.AuthService.VerifyCodeAndRegister(request.Email, request.Code); err != nil {
//		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
//		return
//	}
//
//	c.JSON(http.StatusCreated, gin.H{"message": "User registered successfully"})
//}
