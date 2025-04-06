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

func (h *ResumeHandler) UploadResume(c *gin.Context) {
	userID := c.GetInt("user_id")

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

	parsedResume.UserID = userID

	fullName := parsedResume.FullName
	desiredPosition := parsedResume.DesiredPosition
	skills := parsedResume.Skills
	city := parsedResume.City
	about := parsedResume.About

	if fullName == "" || desiredPosition == "" {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "Missing required fields in resume"})
		return
	}

	err = h.ResumeService.SaveResume(c.Request.Context(), parsedResume)
	if err != nil {
		logger.Log.Error("Failed to save resume", "error", err)
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "Failed to save resume"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"full_name":        fullName,
		"desired_position": desiredPosition,
		"skills":           skills,
		"city":             city,
		"about":            about,
		"parsed_data":      parsedResume.ParsedData,
	})
}

func (h *ResumeHandler) CreateResume(c *gin.Context) {
	userID := c.GetInt("user_id")

	var req dto.ResumeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Log.Error("Failed to bind JSON", "error", err)
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "Invalid input data"})
		return
	}

	err := h.ResumeService.CreateResumeFromRequest(c.Request.Context(), userID, req)
	if err != nil {
		logger.Log.Error("Failed to save resume", "error", err)
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "Failed to save resume"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Resume saved successfully",
	})
}

func (h *ResumeHandler) GetResumeByUserID(c *gin.Context) {
	userID := c.GetInt("user_id")

	resume, user, err := h.ResumeService.GetResumeAndUserByUserID(c.Request.Context(), userID)
	if err != nil {
		logger.Log.Error("Failed to get resume and user", "error", err)
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: err.Error()})
		return
	}

	if resume == nil || user == nil {
		c.JSON(http.StatusNotFound, dto.ErrorResponse{Error: "Resume or user not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"full_name":        resume.FullName,
		"desired_position": resume.DesiredPosition,
		"skills":           resume.Skills,
		"city":             resume.City,
		"about":            resume.About,
		"parsed_data":      resume.ParsedData,
		"user": gin.H{
			"id":              user.ID,
			"email":           user.Email,
			"first_name":      user.FirstName,
			"last_name":       user.LastName,
			"profile_picture": user.ProfilePicture,
			"role_id":         user.RoleId,
		},
	})
}

func (h *ResumeHandler) DeleteResumeByUserID(c *gin.Context) {
	userID := c.GetInt("user_id")

	err := h.ResumeService.DeleteResumeByUserID(c.Request.Context(), userID)
	if err != nil {
		logger.Log.Error("Failed to delete resume", "error", err)
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Resume deleted successfully",
	})
}
