package handler

import (
	"github.com/gin-gonic/gin"
	"jumyste-app-backend/internal/service"
	"jumyste-app-backend/pkg/logger"
	"net/http"
	"strconv"
)

type JobApplicationHandler struct {
	JobApplicationService *service.JobApplicationService
	UserService           *service.UserService
}

func NewJobApplicationHandler(service *service.JobApplicationService, userService *service.UserService) *JobApplicationHandler {
	return &JobApplicationHandler{JobApplicationService: service, UserService: userService}
}

func (h *JobApplicationHandler) ApplyForJob(c *gin.Context) {
	userID := c.GetInt("user_id")
	vacancyIDStr := c.Param("vacancy_id")

	vacancyID, err := strconv.Atoi(vacancyIDStr)
	if err != nil {
		logger.Log.Error("Invalid vacancy ID", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid vacancy ID"})
		return
	}

	user, err := h.UserService.GetUserByID(userID)
	if err != nil {
		logger.Log.Error("Failed to retrieve user information", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve user information"})
		return
	}

	application, err := h.JobApplicationService.ApplyForJob(c.Request.Context(), userID, vacancyID, user.FirstName, user.LastName, user.Email)
	if err != nil {
		logger.Log.Error("Failed to apply for job", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, application)
}

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
