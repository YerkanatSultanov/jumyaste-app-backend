package handler

import (
	"jumyste-app-backend/internal/dto"
	_ "jumyste-app-backend/internal/entity"
	"net/http"
	"strconv"

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
// @Summary Upload a resume
// @Description Upload a resume file and process it
// @Tags Resume
// @Accept multipart/form-data
// @Produce json
// @Param resume formData file true "Resume file"
// @Security BearerAuth
// @Success 200 {object} dto.ResumeResponse
// @Failure 400 {object} dto.ErrorResponse "Failed to retrieve resume file"
// @Failure 401 {object} dto.ErrorResponse "Unauthorized"
// @Failure 500 {object} dto.ErrorResponse "Failed to process resume"
// @Router /resume/upload [post]
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

	c.JSON(http.StatusOK, dto.ResumeResponse{
		FullName:        fullName,
		DesiredPosition: desiredPosition,
		Skills:          skills,
		City:            city,
		About:           about,
		ParsedData:      parsedResume.ParsedData,
	})
}

// CreateResume godoc
// @Summary Create a resume from JSON data
// @Description Create and save a resume using the provided JSON data
// @Tags Resume
// @Accept json
// @Produce json
// @Param resume_request body dto.ResumeRequest true "Resume data"
// @Security BearerAuth
// @Success 200 {object} dto.SuccessResponse "Resume saved successfully"
// @Failure 400 {object} dto.ErrorResponse "Invalid input data"
// @Failure 401 {object} dto.ErrorResponse "Unauthorized"
// @Failure 500 {object} dto.ErrorResponse "Failed to save resume"
// @Router /resume/manual [post]
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

	c.JSON(http.StatusOK, dto.SuccessResponse{
		Message: "Resume saved successfully",
	})
}

// GetResumeByUserID godoc
// @Summary Get a resume by user ID
// @Description Retrieve the resume of a user by their user ID
// @Tags Resume
// @Accept json
// @Produce json
// @Param user_id path int true "User ID"
// @Security BearerAuth
// @Success 200 {object} dto.ResumeResponse
// @Failure 400 {object} dto.ErrorResponse "Invalid user ID"
// @Failure 401 {object} dto.ErrorResponse "Unauthorized"
// @Failure 404 {object} dto.ErrorResponse "Resume or user not found"
// @Failure 500 {object} dto.ErrorResponse "Failed to get resume"
// @Router /resume/{user_id} [get]
func (h *ResumeHandler) GetResumeByUserID(c *gin.Context) {
	userID := c.Param("user_id")

	userIDInt, err := strconv.Atoi(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "Invalid user ID"})
		return
	}

	resume, user, err := h.ResumeService.GetResumeAndUserByUserID(c.Request.Context(), userIDInt)
	if err != nil {
		logger.Log.Error("Failed to get resume and user", "error", err)
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: err.Error()})
		return
	}

	if resume == nil || user == nil {
		c.JSON(http.StatusNotFound, dto.ErrorResponse{Error: "Resume or user not found"})
		return
	}

	var workExps []dto.WorkExperienceResponse
	for _, exp := range resume.Experiences {
		workExps = append(workExps, dto.WorkExperienceResponse{
			CompanyName:    exp.CompanyName,
			Position:       exp.Position,
			StartDate:      exp.StartDate,
			EndDate:        exp.EndDate,
			Location:       exp.Location,
			EmploymentType: exp.EmploymentType,
			Description:    exp.Description,
		})
	}

	c.JSON(http.StatusOK, dto.ResumeResponse{
		FullName:        resume.FullName,
		DesiredPosition: resume.DesiredPosition,
		Skills:          resume.Skills,
		City:            resume.City,
		About:           resume.About,
		ParsedData:      resume.ParsedData,
		User: dto.UserResponse{
			ID:             user.ID,
			Email:          user.Email,
			FirstName:      user.FirstName,
			LastName:       user.LastName,
			ProfilePicture: user.ProfilePicture,
			RoleID:         user.RoleId,
		},
		WorkExperiences: workExps,
	})
}

// DeleteResumeByUserID godoc
// @Summary Delete a resume by user ID
// @Description Delete the resume associated with a given user ID
// @Tags Resume
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} dto.SuccessResponse "Successfully deleted resume"
// @Failure 400 {object} dto.ErrorResponse "Invalid user ID"
// @Failure 401 {object} dto.ErrorResponse "Unauthorized"
// @Failure 500 {object} dto.ErrorResponse "Failed to delete resume"
// @Router /resume/ [delete]
func (h *ResumeHandler) DeleteResumeByUserID(c *gin.Context) {
	userID := c.GetInt("user_id")

	err := h.ResumeService.DeleteResumeByUserID(c.Request.Context(), userID)
	if err != nil {
		logger.Log.Error("Failed to delete resume", "error", err)
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse{
		Message: "Resume deleted successfully",
	})
}

// FilterCandidates godoc
// @Summary Filter candidates based on specified criteria
// @Description Filter candidates using multiple query parameters such as AI match score, skills, city, and position.
// @Tags Resume
// @Accept json
// @Produce json
// @Param ai_match query int false "Minimum AI match score"
// @Param skills query string false "Skills (can be passed multiple times)"
// @Param city query string false "City of the candidate"
// @Param position query string false "Desired position of the candidate"
// @Security BearerAuth
// @Success 200 {array} entity.JobApplicationWithResume "List of filtered candidates"
// @Failure 400 {object} dto.ErrorResponse "Invalid query parameters"
// @Failure 401 {object} dto.ErrorResponse "Unauthorized"
// @Failure 500 {object} dto.ErrorResponse "Failed to filter candidates"
// @Router /resume/candidates [get]
func (h *ResumeHandler) FilterCandidates(c *gin.Context) {
	var filter dto.CandidateFilter

	if err := c.ShouldBindQuery(&filter); err != nil {
		logger.Log.Warn("Invalid filter query params", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid query parameters"})
		return
	}

	ctx := c.Request.Context()
	candidates, err := h.ResumeService.FilterCandidates(ctx, filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to filter candidates"})
		return
	}

	c.JSON(http.StatusOK, candidates)
}
