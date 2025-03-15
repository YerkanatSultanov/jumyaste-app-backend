package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"jumyste-app-backend/internal/service"
	"jumyste-app-backend/pkg/logger"
)

type ResumeHandler struct {
	ResumeService *service.ResumeService
}

func NewResumeHandler(resumeService *service.ResumeService) *ResumeHandler {
	return &ResumeHandler{ResumeService: resumeService}
}

func (h *ResumeHandler) UploadResume(c *gin.Context) {
	file, _, err := c.Request.FormFile("resume")
	if err != nil {
		logger.Log.Error("Failed to retrieve resume file", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to retrieve resume file"})
		return
	}
	defer file.Close()

	parsedResume, err := h.ResumeService.ProcessResume(file)
	if err != nil {
		logger.Log.Error("Failed to process resume", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process resume"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"resume": parsedResume})
}
