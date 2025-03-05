package handler

import (
	"github.com/gin-gonic/gin"
	"jumyste-app-backend/internal/entity"
	"jumyste-app-backend/internal/service"
	"jumyste-app-backend/pkg/logger"
	"log/slog"
	"net/http"
	"strconv"
)

type VacancyHandler struct {
	VacancyService *service.VacancyService
}

func NewVacancyHandler(vacancyService *service.VacancyService) *VacancyHandler {
	return &VacancyHandler{VacancyService: vacancyService}
}

func (h *VacancyHandler) CreateVacancy(c *gin.Context) {
	userID := c.GetInt("user_id")

	var vacancy entity.Vacancy
	if err := c.ShouldBindJSON(&vacancy); err != nil {
		logger.Log.Error("Invalid vacancy input", slog.String("error", err.Error()))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	vacancy.CreatedBy = userID

	logger.Log.Info("Creating vacancy", slog.Int("created_by", vacancy.CreatedBy), slog.String("title", vacancy.Title))

	if err := h.VacancyService.CreateVacancy(&vacancy); err != nil {
		logger.Log.Error("Failed to create vacancy", slog.String("error", err.Error()))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create vacancy"})
		return
	}

	logger.Log.Info("Vacancy created successfully", slog.String("title", vacancy.Title))
	c.JSON(http.StatusCreated, vacancy)
}

func (h *VacancyHandler) UpdateVacancy(c *gin.Context) {
	vacancyID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		logger.Log.Error("Invalid vacancy ID", slog.String("error", err.Error()))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid vacancy ID"})
		return
	}

	userID := c.GetInt("user_id")
	vacancy, err := h.VacancyService.GetVacancyById(vacancyID)
	if err != nil {
		logger.Log.Error("Vacancy not found", slog.Int("vacancy_id", vacancyID), slog.String("error", err.Error()))
		c.JSON(http.StatusNotFound, gin.H{"error": "Vacancy not found"})
		return
	}

	if vacancy.CreatedBy != userID {
		logger.Log.Warn("Unauthorized vacancy update attempt",
			slog.Int("vacancy_id", vacancyID), slog.Int("user_id", userID))
		c.JSON(http.StatusForbidden, gin.H{"error": "You can only update your own vacancies"})
		return
	}

	if err := c.ShouldBindJSON(&vacancy); err != nil {
		logger.Log.Error("Invalid vacancy update input", slog.String("error", err.Error()))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	if err := h.VacancyService.UpdateVacancy(vacancy); err != nil {
		logger.Log.Error("Failed to update vacancy", slog.String("error", err.Error()))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update vacancy"})
		return
	}

	logger.Log.Info("Vacancy updated successfully", slog.Int("vacancy_id", vacancyID))
	c.JSON(http.StatusOK, gin.H{"message": "Vacancy updated successfully"})
}

func (h *VacancyHandler) DeleteVacancy(c *gin.Context) {
	vacancyID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		logger.Log.Error("Invalid vacancy ID", slog.String("error", err.Error()))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid vacancy ID"})
		return
	}

	userID := c.GetInt("user_id")

	logger.Log.Info("Attempting to delete vacancy",
		slog.Int("vacancy_id", vacancyID), slog.Int("user_id", userID))

	err = h.VacancyService.DeleteVacancy(vacancyID, userID)
	if err != nil {
		if err.Error() == "vacancy not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "Vacancy not found"})
			return
		}
		if err.Error() == "you can only delete your own vacancies" {
			c.JSON(http.StatusForbidden, gin.H{"error": "You can only delete your own vacancies"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete vacancy"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Vacancy deleted successfully"})
}

func (h *VacancyHandler) GetAllVacancies(c *gin.Context) {
	vacancies, err := h.VacancyService.GetAllVacancies()
	if err != nil {
		logger.Log.Error("Failed to fetch vacancies", slog.String("error", err.Error()))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch vacancies"})
		return
	}

	logger.Log.Info("Vacancies retrieved", slog.Int("count", len(vacancies)))
	c.JSON(http.StatusOK, vacancies)
}
