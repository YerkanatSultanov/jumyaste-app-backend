package handler

import (
	"jumyste-app-backend/internal/dto"
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

// UploadResume godoc
//
// @Summary Upload and parse a resume
// @Description Accepts a PDF file, extracts text, and returns structured resume data
// @Tags Resume
// @Accept multipart/form-data
// @Produce json
// @Param resume formData file true "Resume file (PDF only)"
// @Success 200 {object} map[string]interface{} "Parsed resume data"
// @Failure 400 {object} dto.ErrorResponse "Failed to retrieve resume file"
// @Failure 500 {object} dto.ErrorResponse "Failed to process resume"
// @Security BearerAuth
// @Router /resume/upload [post]
func (h *ResumeHandler) UploadResume(c *gin.Context) {
	file, _, err := c.Request.FormFile("resume")
	if err != nil {
		logger.Log.Error("Failed to retrieve resume file", "error", err)
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "Failed to retrieve resume file"})
		return
	}
	defer file.Close()

	parsedResume, err := h.ResumeService.ProcessResume(file)
	if err != nil {
		logger.Log.Error("Failed to process resume", "error", err)
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "Failed to process resume"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"resume": parsedResume})
}
