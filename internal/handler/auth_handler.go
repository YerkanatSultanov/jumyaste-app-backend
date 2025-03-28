package handler

import (
	"errors"
	"github.com/gin-gonic/gin"
	"jumyste-app-backend/internal/dto"
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

// Login godoc
// @Summary User login
// @Description Authenticates a user and returns a JWT token.
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body dto.LoginRequest true "Login credentials"
// @Success 200 {object} dto.LoginResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Router /auth/login [post]
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

// RequestPasswordReset godoc
// @Summary Request password reset
// @Description Sends a password reset code to the user's email.
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body dto.RequestPasswordResetRequest true "User email"
// @Success 200 {object} dto.RequestPasswordResetResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /auth/forgot-password [post]
func (h *AuthHandler) RequestPasswordReset(c *gin.Context) {
	var req dto.RequestPasswordResetRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "Invalid request"})
		return
	}

	err := h.AuthService.RequestPasswordReset(req.Email)
	if err != nil {
		logger.Log.Error("Failed to request password reset", slog.String("email", req.Email), slog.String("error", err.Error()))
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "Failed to request password reset"})
		return
	}

	c.JSON(http.StatusOK, dto.RequestPasswordResetResponse{Message: "Reset code sent to your email"})
}

// ResetPassword godoc
// @Summary Reset user password
// @Description Resets the password using a reset code.
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body dto.ResetPasswordRequest true "Reset password request"
// @Success 200 {object} dto.ResetPasswordResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /auth/reset-password [post]
func (h *AuthHandler) ResetPassword(c *gin.Context) {
	var req dto.ResetPasswordRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "Invalid request"})
		return
	}

	err := h.AuthService.ResetPassword(req.Email, req.ResetCode, req.NewPassword, req.ConfirmPassword)
	if err != nil {
		if errors.Is(err, service.ErrInvalidResetCode) {
			c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "Invalid or expired reset code"})
			return
		}
		if errors.Is(err, service.ErrUserNotFound) {
			c.JSON(http.StatusNotFound, dto.ErrorResponse{Error: "User not found"})
			return
		}
		logger.Log.Error("Failed to reset password", slog.String("email", req.Email), slog.String("error", err.Error()))
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "Failed to reset password"})
		return
	}

	c.JSON(http.StatusOK, dto.ResetPasswordResponse{Message: "Password reset successful"})
}

// Register godoc
// @Summary Register a new user
// @Description Creates a new user account with the provided details.
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body dto.RegisterUserRequest true "User registration data"
// @Success 200 {object} dto.SuccessResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var request dto.RegisterUserRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	if request.RoleID == 0 {
		request.RoleID = 3
	}

	user := &entity.User{
		Email:     request.Email,
		Password:  request.Password,
		FirstName: request.FirstName,
		LastName:  request.LastName,
		RoleID:    request.RoleID,
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
