package handler

import (
	"github.com/gin-gonic/gin"
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

type sendInvitationRequest struct {
	Email     string `json:"email"`
	CompanyID int    `json:"company_id"`
	DepID     int    `json:"dep_id"`
}

func (h *InvitationHandler) SendInvitationHandler(c *gin.Context) {
	var req sendInvitationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Log.Error("Invalid request body", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	if req.Email == "" || req.CompanyID == 0 || req.DepID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing required fields"})
		return
	}

	err := h.service.SendInvitation(req.Email, req.CompanyID, req.DepID)
	if err != nil {
		logger.Log.Error("Failed to send invitation", "email", req.Email, "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send invitation"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Invitation sent successfully"})
}
