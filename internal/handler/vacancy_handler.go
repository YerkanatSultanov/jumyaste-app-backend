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
	vacancy.Status = "open"

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
	vacancy, err := h.VacancyService.GetVacancyById(vacancyID, true)
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
// @Tags All can use
// @Produce json
// @Security BearerAuth
// @Success 200 {array} entity.Vacancy "List of vacancies"
// @Failure 500 {object} dto.ErrorResponse "Failed to fetch vacancies"
// @Router /users/vacancy [get]
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
// @Description Allows searching for vacancies based on various filters, including an optional status filter.
//
//	If 'status' is set to 'all', it will return vacancies regardless of their status (open or closed).
//
// @Tags All can use
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param query query string false "Search query"
// @Param employment_type query []string false "Employment type filter" collectionFormat(multi)
// @Param work_format query []string false "Work format filter" collectionFormat(multi)
// @Param skills query []string false "Skills filter" collectionFormat(multi)
// @Param status query string false "Filter vacancies by status (open, closed, or 'all' for all vacancies)" Enum(open,closed,all) default(all)
// @Success 200 {array} entity.Vacancy "List of matching vacancies"
// @Failure 400 {object} dto.ErrorResponse "Invalid search parameters"
// @Failure 500 {object} dto.ErrorResponse "Failed to search vacancies"
// @Router /users/vacancy/search [get]
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

	if status := c.DefaultQuery("status", "all"); status != "all" {
		filter.Status = status
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

// GetVacancyByCompanyID godoc
// @Summary      Get vacancies by company ID
// @Description  Retrieve all vacancies for the company of the current HR
// @Tags         Vacancies
// @Security     BearerAuth
// @Produce      json
// @Success      200  {array}   entity.Vacancy
// @Failure      500  {object}  dto.ErrorResponse  "Failed to retrieve vacancies"
// @Router       /vacancies/company [get]
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

// GetVacancyByIDForUser godoc
// @Summary      Get vacancy by ID
// @Description  Retrieve a specific vacancy by its ID
// @Security     BearerAuth
// @Tags         All can use
// @Produce      json
// @Param        id   path      int  true  "Vacancy ID"
// @Success      200  {object}  entity.Vacancy
// @Failure      400  {object}  dto.ErrorResponse  "Invalid vacancy ID"
// @Failure      404  {object}  dto.ErrorResponse  "Vacancy not found"
// @Router       /users/vacancy/{id} [get]
func (h *VacancyHandler) GetVacancyByIDForUser(c *gin.Context) {
	vacancyID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		logger.Log.Error("Invalid vacancy ID", slog.String("error", err.Error()))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid vacancy ID"})
		return
	}

	vacancy, err := h.VacancyService.GetVacancyById(vacancyID, false)
	if err != nil {
		logger.Log.Error("Vacancy not found", slog.Int("vacancy_id", vacancyID), slog.String("error", err.Error()))
		c.JSON(http.StatusNotFound, gin.H{"error": "Vacancy not found"})
		return
	}

	logger.Log.Info("Vacancy retrieved successfully", slog.Int("vacancy_id", vacancy.ID))
	c.JSON(http.StatusOK, vacancy)
}

// GetVacancyByIDForHr godoc
// @Summary      Get vacancy by ID
// @Description  Retrieve a specific vacancy by its ID
// @Security     BearerAuth
// @Tags         Vacancies
// @Produce      json
// @Param        id   path      int  true  "Vacancy ID"
// @Success      200  {object}  entity.Vacancy
// @Failure      400  {object}  dto.ErrorResponse  "Invalid vacancy ID"
// @Failure      404  {object}  dto.ErrorResponse  "Vacancy not found"
// @Router       /vacancies/hr/{id} [get]
func (h *VacancyHandler) GetVacancyByIDForHr(c *gin.Context) {
	vacancyID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		logger.Log.Error("Invalid vacancy ID", slog.String("error", err.Error()))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid vacancy ID"})
		return
	}

	vacancy, err := h.VacancyService.GetVacancyById(vacancyID, true)
	if err != nil {
		logger.Log.Error("Vacancy not found", slog.Int("vacancy_id", vacancyID), slog.String("error", err.Error()))
		c.JSON(http.StatusNotFound, gin.H{"error": "Vacancy not found"})
		return
	}

	logger.Log.Info("Vacancy retrieved successfully", slog.Int("vacancy_id", vacancy.ID))
	c.JSON(http.StatusOK, vacancy)
}

// UpdateVacancyStatusHandler godoc
// @Summary      Update vacancy status
// @Description  Update the status of a vacancy by its ID
// @Security     BearerAuth
// @Tags         Vacancies
// @Accept       json
// @Produce      json
// @Param        id     path     int  true  "Vacancy ID"
// @Param        request body dto.UpdateVacancyStatusRequest true  "New status"
// @Success      200 {object} dto.SuccessResponse
// @Failure      400 {object} dto.ErrorResponse
// @Failure      500 {object} dto.ErrorResponse
// @Router       /vacancies/status/{id} [put]
func (h *VacancyHandler) UpdateVacancyStatusHandler(c *gin.Context) {
	vacancyID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		logger.Log.Error("Invalid vacancy ID", slog.String("error", err.Error()))
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "Invalid vacancy ID"})
		return
	}
	var request dto.UpdateVacancyStatusRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		logger.Log.Error("Invalid request body", slog.String("error", err.Error()))
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "Invalid request body"})
		return
	}

	logger.Log.Info("Attempting to update vacancy status", slog.Int("vacancy_id", vacancyID), slog.String("new_status", request.Status))

	err = h.VacancyService.UpdateVacancyStatus(vacancyID, request.Status)
	if err != nil {
		logger.Log.Error("Failed to update vacancy status", slog.Int("vacancy_id", vacancyID), slog.String("error", err.Error()))
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "Failed to update vacancy status"})
		return
	}

	logger.Log.Info("Successfully updated vacancy status", slog.Int("vacancy_id", vacancyID))

	c.JSON(http.StatusOK, dto.SuccessResponse{Message: "Successfully updated vacancy status"})
}

// GetFeedData godoc
//
// @Summary Get feed data for user
// @Description Returns count of new vacancies since last feed view
// @Tags Vacancies
// @Security BearerAuth
// @Success 200 {object} dto.FeedDataResponse "Feed data"
// @Failure 401 {object} dto.ErrorResponse "Unauthorized"
// @Failure 500 {object} dto.ErrorResponse "Internal error"
// @Router /vacancies/feed/data [get]
func (h *VacancyHandler) GetFeedData(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: "Unauthorized"})
		return
	}

	data, err := h.VacancyService.GetFeedData(userID.(int))
	if err != nil {
		logger.Log.Error("Failed to get feed data", "error", err)
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "Failed to get feed data"})
		return
	}

	c.JSON(http.StatusOK, data)
}

// GenerateVacancyDescription godoc
// @Summary Generate vacancy description
// @Description Generate a detailed HTML description for the given vacancy details.
// @Tags Vacancies
// @Security BearerAuth
// @Accept  json
// @Produce  json
// @Param vacancyInput body dto.VacancyInput true "Vacancy Input"
// @Success 200 {object} dto.DescriptionResponse "Description generated successfully"
// @Failure 400 {object} dto.ErrorResponse "Invalid input"
// @Failure 500 {object} dto.ErrorResponse "Failed to generate description"
// @Router /vacancies/generate-description [post]
func (h *VacancyHandler) GenerateVacancyDescription(c *gin.Context) {
	var input dto.VacancyInput
	if err := c.ShouldBindJSON(&input); err != nil {
		logger.Log.Error("Invalid vacancy input", "error", err)
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "Invalid input"})
		return
	}

	description, err := h.VacancyService.GenerateDescription(input)
	if err != nil {
		logger.Log.Error("Failed to generate vacancy description", "error", err)
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "Failed to generate description"})
		return
	}

	c.JSON(http.StatusOK, dto.DescriptionResponse{Description: description})
}
