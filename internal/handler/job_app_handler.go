package handler

import (
	"github.com/gin-gonic/gin"
	_ "jumyste-app-backend/internal/dto"
	"jumyste-app-backend/internal/service"
	"jumyste-app-backend/pkg/logger"
	"net/http"
	"strconv"
)

type JobApplicationHandler struct {
	JobApplicationService *service.JobApplicationService
	ResumeService         *service.ResumeService
}

func NewJobApplicationHandler(service *service.JobApplicationService, resumeService *service.ResumeService) *JobApplicationHandler {
	return &JobApplicationHandler{JobApplicationService: service, ResumeService: resumeService}
}

// ApplyForJob godoc
// @Summary Apply for a job
// @Description Apply for a job by providing vacancy ID and user details
// @Tags Job Applications
// @Accept json
// @Produce json
// @Param vacancy_id path int true "Vacancy ID"
// @Security BearerAuth
// @Success 201 {object} dto.JobApplicationResponse
// @Failure 400 {object} dto.ErrorResponse "Invalid vacancy ID"
// @Failure 401 {object} dto.ErrorResponse "Unauthorized"
// @Failure 500 {object} dto.ErrorResponse "Failed to apply for job"
// @Router /jobs/apply/{vacancy_id} [post]
func (h *JobApplicationHandler) ApplyForJob(c *gin.Context) {
	userID := c.GetInt("user_id")
	vacancyIDStr := c.Param("vacancy_id")

	vacancyID, err := strconv.Atoi(vacancyIDStr)
	if err != nil {
		logger.Log.Error("Invalid vacancy ID", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid vacancy ID"})
		return
	}
	resume, user, err := h.ResumeService.GetResumeAndUserByUserID(c.Request.Context(), userID)
	if err != nil {
		logger.Log.Error("Failed to retrieve user information and resume", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve user information and resume"})
		return
	}

	application, err := h.JobApplicationService.ApplyForJob(c.Request.Context(), userID, vacancyID, user.FirstName, user.LastName, user.Email, resume.ID)
	if err != nil {
		logger.Log.Error("Failed to apply for job", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, application)
}

// GetJobApplicationsByVacancyID godoc
// @Summary Get job applications by vacancy ID
// @Description Retrieve all job applications for a specific vacancy
// @Tags Job Applications
// @Accept json
// @Produce json
// @Param vacancy_id path int true "Vacancy ID"
// @Security BearerAuth
// @Success 200 {array} dto.JobApplicationWithResumeResponse
// @Failure 400 {object} dto.ErrorResponse "Invalid vacancy ID"
// @Failure 401 {object} dto.ErrorResponse "Unauthorized"
// @Failure 500 {object} dto.ErrorResponse "Failed to retrieve job applications"
// @Router /jobs/{vacancy_id} [get]
func (h *JobApplicationHandler) GetJobApplicationsByVacancyID(c *gin.Context) {
	vacancyIDStr := c.Param("vacancy_id")

	vacancyID, err := strconv.Atoi(vacancyIDStr)
	if err != nil {
		logger.Log.Error("Invalid vacancy ID", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid vacancy ID"})
		return
	}

	applications, err := h.JobApplicationService.GetJobApplicationsByVacancyID(c.Request.Context(), vacancyID)
	if err != nil {
		logger.Log.Error("Failed to get job applications", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, applications)
}

// UpdateJobApplicationStatus godoc
// @Summary Update the status of a job application
// @Description Update the status of a specific job application by application ID
// @Tags Job Applications
// @Accept json
// @Produce json
// @Param application_id path int true "Application ID"
// @Param status path string true "New Status"
// @Security BearerAuth
// @Success 200 {object} dto.SuccessResponse
// @Failure 400 {object} dto.ErrorResponse "Invalid application ID or status"
// @Failure 401 {object} dto.ErrorResponse "Unauthorized"
// @Failure 500 {object} dto.ErrorResponse "Failed to update application status"
// @Router /jobs/{application_id}/status/{status} [put]
func (h *JobApplicationHandler) UpdateJobApplicationStatus(c *gin.Context) {
	applicationIDStr := c.Param("application_id")
	status := c.Param("status")

	applicationID, err := strconv.Atoi(applicationIDStr)
	if err != nil {
		logger.Log.Error("Invalid application ID", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid application ID"})
		return
	}

	err = h.JobApplicationService.UpdateJobApplicationStatus(c.Request.Context(), applicationID, status)
	if err != nil {
		logger.Log.Error("Failed to update status", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "updated"})
}

// DeleteJobApplication godoc
// @Summary Delete a job application
// @Description Delete a job application by application ID
// @Tags Job Applications
// @Accept json
// @Produce json
// @Param application_id path int true "Application ID"
// @Security BearerAuth
// @Success 200 {object} dto.SuccessResponse
// @Failure 400 {object} dto.ErrorResponse "Invalid application ID"
// @Failure 401 {object} dto.ErrorResponse "Unauthorized"
// @Failure 500 {object} dto.ErrorResponse "Failed to delete job application"
// @Router /jobs/{application_id} [delete]
func (h *JobApplicationHandler) DeleteJobApplication(c *gin.Context) {
	applicationIDStr := c.Param("application_id")

	applicationID, err := strconv.Atoi(applicationIDStr)
	if err != nil {
		logger.Log.Error("Invalid application ID", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid application ID"})
		return
	}

	err = h.JobApplicationService.DeleteJobApplication(c.Request.Context(), applicationID)
	if err != nil {
		logger.Log.Error("Failed to delete job application", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "deleted"})
}
func (h *JobApplicationHandler) GetJobAppAnalytics(c *gin.Context) {
	userID := c.GetInt("user_id")
	logger.Log.Info("Getting HR analytics", "user_id", userID)

	stats, err := h.JobApplicationService.GetJobAppAnalytics(c.Request.Context(), userID)
	if err != nil {
		logger.Log.Error("Failed to get HR analytics", "error", err, "user_id", userID)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get analytics"})
		return
	}

	logger.Log.Info("Successfully retrieved HR analytics", "user_id", userID, "stats_count", len(stats))
	c.JSON(http.StatusOK, stats)
}
