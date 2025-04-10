package handler

import (
	"github.com/gin-gonic/gin"
	"jumyste-app-backend/internal/dto"
	"jumyste-app-backend/internal/entity"
	"jumyste-app-backend/internal/service"
	"jumyste-app-backend/pkg/logger"
	"log/slog"
	"net/http"
)

type UserHandler struct {
	UserService *service.UserService
}

func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{UserService: userService}
}

func (h *UserHandler) CreateUser(c *gin.Context) {
	var user entity.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	if err := h.UserService.CreateUser(&user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	c.JSON(http.StatusCreated, user)
}

// GetUser godoc
// @Summary      Get user information
// @Description  Retrieves the profile information of the authenticated user
// @Tags         Users
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Success      200  {object}  entity.UserResponse  "User profile information"
// @Failure      401  {object}  dto.ErrorResponse  "Unauthorized - Token is missing or invalid"
// @Failure      404  {object}  dto.ErrorResponse  "User not found - No user associated with the given ID"
// @Failure      500  {object}  dto.ErrorResponse  "Internal server error - Invalid user ID type"
// @Router       /users/me [get]
func (h *UserHandler) GetUser(c *gin.Context) {
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

	user, err := h.UserService.GetUserByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, dto.ErrorResponse{Error: "User not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}

// UpdateUser godoc
// @Summary      Update user information
// @Description  Update user details such as name, email, or profile picture
// @Tags         Users
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        updates  body      map[string]interface{}  true  "Fields to update"
// @Success      200      {object}  dto.SuccessResponse
// @Failure      400      {object}  dto.ErrorResponse  "Invalid request body"
// @Failure      401      {object}  dto.ErrorResponse  "Unauthorized"
// @Failure      500      {object}  dto.ErrorResponse  "Failed to update user"
// @Router       /users/me [patch]
func (h *UserHandler) UpdateUser(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	id, ok := userID.(int)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID"})
		return
	}

	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		logger.Log.Error("Invalid request body", slog.String("error", err.Error()))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	logger.Log.Info("Updating user", slog.Int("user_id", id), slog.Any("updates", updates))

	if err := h.UserService.UpdateUser(id, updates); err != nil {
		logger.Log.Error("Failed to update user", slog.String("error", err.Error()))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User updated successfully"})
}
