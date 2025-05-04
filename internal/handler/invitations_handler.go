package handler

import (
	"github.com/gin-gonic/gin"
	_ "jumyste-app-backend/internal/dto"
	"jumyste-app-backend/internal/service"
	"jumyste-app-backend/pkg/logger"
	"net/http"
)

type InvitationHandler struct {
	service *service.InvitationService
}

func NewInvitationHandler(service *service.InvitationService) *InvitationHandler {
	return &InvitationHandler{service: service}
}

type SendInvitationRequest struct {
	Email string `json:"email"`
	DepID int    `json:"dep_id"`
}

// SendInvitationHandler godoc
// @Summary      Send invitation to register
// @Description  Send an invitation email to a user with company and department information
// @Tags         Invitations
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request  body      SendInvitationRequest  true  "Invitation request body"
// @Success      200      {object}  dto.SuccessResponse
// @Failure      400      {object}  dto.ErrorResponse  "Invalid request body or missing required fields"
// @Failure      500      {object}  dto.ErrorResponse  "Failed to send invitation"
// @Router       /invitations [post]
func (h *InvitationHandler) SendInvitationHandler(c *gin.Context) {
	var req SendInvitationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Log.Error("Invalid request body", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	companyID, exists := c.Get("company_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	err := h.service.SendInvitation(req.Email, companyID.(int), req.DepID)
	if err != nil {
		logger.Log.Error("Failed to send invitation", "email", req.Email, "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send invitation"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Invitation sent successfully"})
}
