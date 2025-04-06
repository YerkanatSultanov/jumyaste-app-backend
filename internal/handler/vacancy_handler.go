package handler

import (
	"github.com/gin-gonic/gin"
	"jumyste-app-backend/internal/dto"
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

// CreateVacancy godoc
//
// @Summary Create a new vacancy
// @Description Allows an employer to create a new vacancy
// @Tags Vacancies
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param vacancy body dto.CreateVacancyRequest true "Vacancy details"
// @Success 201 {object} entity.Vacancy "Vacancy successfully created"
// @Failure 400 {object} dto.ErrorResponse "Invalid input"
// @Failure 500 {object} dto.ErrorResponse "Failed to create vacancy"
// @Router /vacancies [post]
func (h *VacancyHandler) CreateVacancy(c *gin.Context) {
	userID := c.GetInt("user_id")
	companyID := c.GetInt("company_id")

	var vacancy entity.Vacancy
	if err := c.ShouldBindJSON(&vacancy); err != nil {
		logger.Log.Error("Invalid vacancy input", slog.String("error", err.Error()))
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "Invalid input"})
		return
	}

	vacancy.CreatedBy = userID
	vacancy.CompanyId = companyID

	logger.Log.Info("Creating vacancy", slog.Int("created_by", vacancy.CreatedBy), slog.String("title", vacancy.Title))

	if err := h.VacancyService.CreateVacancy(&vacancy); err != nil {
		logger.Log.Error("Failed to create vacancy", slog.String("error", err.Error()))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create vacancy"})
		return
	}
	logger.Log.Info("Created Vacancy", slog.String("title", vacancy.Title), slog.Time("created_at", vacancy.CreatedAt))

	logger.Log.Info("Vacancy created successfully", slog.String("title", vacancy.Title))
	c.JSON(http.StatusCreated, vacancy)
}

// UpdateVacancy godoc
//
// @Summary Update an existing vacancy
// @Description Allows an employer to update their own vacancy
// @Tags Vacancies
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Vacancy ID"
// @Param vacancy body dto.UpdateVacancyRequest true "Updated vacancy details"
// @Success 200 {object} dto.SuccessResponse "Vacancy updated successfully"
// @Failure 400 {object} dto.ErrorResponse "Invalid input or vacancy ID"
// @Failure 403 {object} dto.ErrorResponse "User does not own the vacancy"
// @Failure 404 {object} dto.ErrorResponse "Vacancy not found"
// @Failure 500 {object} dto.ErrorResponse "Failed to update vacancy"
// @Router /vacancies/{id} [put]
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

// DeleteVacancy godoc
//
// @Summary Delete a vacancy
// @Description Allows an employer to delete their own vacancy
// @Tags Vacancies
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Vacancy ID"
// @Success 200 {object} dto.SuccessResponse "Vacancy deleted successfully"
// @Failure 400 {object} dto.ErrorResponse "Invalid vacancy ID"
// @Failure 403 {object} dto.ErrorResponse "User does not own the vacancy"
// @Failure 404 {object} dto.ErrorResponse "Vacancy not found"
// @Failure 500 {object} dto.ErrorResponse "Failed to delete vacancy"
// @Router /vacancies/{id} [delete]
func (h *VacancyHandler) DeleteVacancy(c *gin.Context) {
	vacancyID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		logger.Log.Error("Invalid vacancy ID", slog.String("error", err.Error()))
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "Invalid vacancy ID"})
		return
	}

	userID := c.GetInt("user_id")

	logger.Log.Info("Attempting to delete vacancy",
		slog.Int("vacancy_id", vacancyID), slog.Int("user_id", userID))

	err = h.VacancyService.DeleteVacancy(vacancyID, userID)
	if err != nil {
		switch err.Error() {
		case "vacancy not found":
			c.JSON(http.StatusNotFound, dto.ErrorResponse{Error: "Vacancy not found"})
		case "you can only delete your own vacancies":
			c.JSON(http.StatusForbidden, dto.ErrorResponse{Error: "You can only delete your own vacancies"})
		default:
			c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "Failed to delete vacancy"})
		}
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse{Message: "Vacancy deleted successfully"})
}

// GetAllVacancies godoc
//
// @Summary Get all vacancies
// @Description Retrieves a list of all vacancies
// @Tags Vacancies
// @Produce json
// @Security BearerAuth
// @Success 200 {array} entity.Vacancy "List of vacancies"
// @Failure 500 {object} dto.ErrorResponse "Failed to fetch vacancies"
// @Router /vacancies [get]
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

// GetMyVacancies godoc
//
// @Summary Get vacancies created by the authenticated HR
// @Description Returns a list of vacancies created by the currently authenticated HR user
// @Tags Vacancies
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {array} entity.Vacancy "List of vacancies"
// @Failure 500 {object} dto.ErrorResponse "Failed to retrieve vacancies"
// @Router /vacancies/my [get]
func (h *VacancyHandler) GetMyVacancies(c *gin.Context) {
	userID := c.GetInt("user_id")

	vacancies, err := h.VacancyService.GetMyVacancies(userID)
	if err != nil {
		logger.Log.Error("Failed to retrieve HR vacancies", slog.Int("user_id", userID), slog.String("error", err.Error()))
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "Failed to retrieve vacancies"})
		return
	}

	c.JSON(http.StatusOK, vacancies)
}

// SearchVacancies godoc
//
// @Summary Search for vacancies
// @Description Allows searching for vacancies based on various filters
// @Tags Vacancies
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param query query string false "Search query"
// @Param employment_type query []string false "Employment type filter" collectionFormat(multi)
// @Param work_format query []string false "Work format filter" collectionFormat(multi)
// @Param skills query []string false "Skills filter" collectionFormat(multi)
// @Success 200 {array} entity.Vacancy "List of matching vacancies"
// @Failure 400 {object} dto.ErrorResponse "Invalid search parameters"
// @Failure 500 {object} dto.ErrorResponse "Failed to search vacancies"
// @Router /vacancies/search [get]
func (h *VacancyHandler) SearchVacancies(c *gin.Context) {
	var filter entity.VacancyFilter

	if err := c.ShouldBindQuery(&filter); err != nil {
		logger.Log.Error("Invalid search parameters", "error", err)
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "Invalid search parameters"})
		return
	}

	filter.Query = c.Query("query")

	if employmentType := c.QueryArray("employment_type"); len(employmentType) > 0 {
		filter.EmploymentType = employmentType
	}

	if workFormat := c.QueryArray("work_format"); len(workFormat) > 0 {
		filter.WorkFormat = workFormat
	}

	if skills := c.QueryArray("skills"); len(skills) > 0 {
		filter.Skills = skills
	}

	logger.Log.Info("Handling search vacancies request", "filter", filter)

	vacancies, err := h.VacancyService.SearchVacancies(filter)
	if err != nil {
		logger.Log.Error("Failed to search vacancies", "error", err)
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "Failed to search vacancies"})
		return
	}

	if len(vacancies) == 0 {
		logger.Log.Info("No vacancies found for the given filters", "filter", filter)
		c.JSON(http.StatusOK, gin.H{"message": "No vacancies found"})
		return
	}

	c.JSON(http.StatusOK, vacancies)
}

func (h *VacancyHandler) GetVacancyByCompanyID(c *gin.Context) {
	companyId := c.GetInt("company_id")
	vacancies, err := h.VacancyService.GetVacanciesByCompanyId(companyId)

	if err != nil {
		logger.Log.Error("Failed to retrieve HR vacancies", slog.Int("user_id", companyId), slog.String("error", err.Error()))
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "Failed to retrieve vacancies"})
		return
	}

	c.JSON(http.StatusOK, vacancies)
}

func (h *VacancyHandler) GetVacancyByID(c *gin.Context) {
	vacancyID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		logger.Log.Error("Invalid vacancy ID", slog.String("error", err.Error()))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid vacancy ID"})
		return
	}

	vacancy, err := h.VacancyService.GetVacancyById(vacancyID)
	if err != nil {
		logger.Log.Error("Vacancy not found", slog.Int("vacancy_id", vacancyID), slog.String("error", err.Error()))
		c.JSON(http.StatusNotFound, gin.H{"error": "Vacancy not found"})
		return
	}

	if err := c.ShouldBindJSON(&vacancy); err != nil {
		logger.Log.Error("Invalid vacancy input", slog.String("error", err.Error()))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	logger.Log.Info(" Vacancy received successfully", slog.Int("vacancy_id", vacancyID))
	c.JSON(http.StatusOK, vacancy)
}
