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

type CompanyHandler struct {
	CompanyService *service.CompanyService
}

func NewCompanyHandler(companyService *service.CompanyService) *CompanyHandler {
	return &CompanyHandler{CompanyService: companyService}
}

// CreateCompany godoc
//
// @Summary Create a new company
// @Tags Companies
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param company body entity.Company true "Company payload"
// @Success 201 {object} entity.Company
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /companies [post]
func (h *CompanyHandler) CreateCompany(c *gin.Context) {
	var company entity.Company
	if err := c.ShouldBindJSON(&company); err != nil {
		logger.Log.Error("Invalid JSON", slog.String("error", err.Error()))
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "Invalid request payload"})
		return
	}

	if err := h.CompanyService.CreateCompany(&company); err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "Failed to create company"})
		return
	}

	c.JSON(http.StatusCreated, company)
}

// GetCompanyByID godoc
//
// @Summary Get company by ID
// @Tags Companies
// @Produce json
// @Security BearerAuth
// @Param id path int true "Company ID"
// @Success 200 {object} entity.Company
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /companies/{id} [get]
func (h *CompanyHandler) GetCompanyByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "Invalid company ID"})
		return
	}

	company, err := h.CompanyService.GetCompanyByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, dto.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, company)
}

// UpdateCompany godoc
//
// @Summary Update a company
// @Tags Companies
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Company ID"
// @Param company body entity.Company true "Updated company"
// @Success 200 {object} dto.SuccessResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /companies/{id} [put]
func (h *CompanyHandler) UpdateCompany(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "Invalid company ID"})
		return
	}

	var company entity.Company
	if err := c.ShouldBindJSON(&company); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "Invalid request payload"})
		return
	}

	company.ID = id

	if err := h.CompanyService.UpdateCompany(&company); err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "Failed to update company"})
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse{Message: "Company updated successfully"})
}

// DeleteCompany godoc
//
// @Summary Delete a company
// @Tags Companies
// @Produce json
// @Security BearerAuth
// @Param id path int true "Company ID"
// @Success 200 {object} dto.SuccessResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /companies/{id} [delete]
func (h *CompanyHandler) DeleteCompany(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "Invalid company ID"})
		return
	}

	if err := h.CompanyService.DeleteCompany(id); err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "Failed to delete company"})
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse{Message: "Company deleted successfully"})
}
