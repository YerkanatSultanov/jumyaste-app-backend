package handler

import (
	"github.com/gin-gonic/gin"
	"jumyste-app-backend/internal/dto"
	"jumyste-app-backend/internal/entity"
	"jumyste-app-backend/internal/service"
	"jumyste-app-backend/pkg/logger"
	"net/http"
	"strconv"
)

type DepartmentsHandler struct {
	DepartmentService *service.DepartmentsService
}

func NewDepartmentsHandler(departmentService *service.DepartmentsService) *DepartmentsHandler {
	return &DepartmentsHandler{DepartmentService: departmentService}
}

// GetMyDepartments godoc
//
// @Summary Get all Departments
// @Description Retrieves a list of all departments for the authenticated user's company
// @Tags Departments
// @Produce json
// @Security BearerAuth
// @Success 200 {array} entity.Department "List of departments"
// @Failure 401 {object} dto.ErrorResponse "Unauthorized"
// @Failure 500 {object} dto.ErrorResponse "Failed to fetch departments"
// @Router /departments/all [get]
func (h *DepartmentsHandler) GetMyDepartments(c *gin.Context) {
	companyID, exists := c.Get("company_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	departments, err := h.DepartmentService.GetDepartmentsByCompany(companyID.(int))
	if err != nil {
		logger.Log.Error("Failed to fetch departments", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch departments"})
		return
	}

	c.JSON(http.StatusOK, departments)
}

// CreateDepartment godoc
//
// @Summary Create a new department
// @Description Creates a new department within the authenticated user's company
// @Tags Departments
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.CreateDepartmentRequest true "Department data"
// @Success 200 {object} entity.Department "Created department"
// @Failure 400 {object} dto.ErrorResponse "Invalid request"
// @Failure 401 {object} dto.ErrorResponse "Unauthorized"
// @Failure 500 {object} dto.ErrorResponse "Failed to create department"
// @Router /departments [post]
func (h *DepartmentsHandler) CreateDepartment(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	companyID, exists := c.Get("company_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var req dto.CreateDepartmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Log.Warn("Invalid department creation request", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	dep := &entity.Department{
		Name:      req.Name,
		Color:     req.Color,
		CompanyId: companyID.(int),
		HrCount:   0,
	}

	err := h.DepartmentService.CreateDepartment(userID.(int), dep)
	if err != nil {
		logger.Log.Error("Failed to create department", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create department"})
		return
	}

	logger.Log.Info("Department created", "id", dep.ID, "name", dep.Name, "company_id", dep.CompanyId)
	c.JSON(http.StatusOK, dep)
}

// GetDepartmentByID godoc
//
// @Summary Get a department by ID
// @Description Get details of a specific department
// @Tags Departments
// @Produce json
// @Security BearerAuth
// @Param id path int true "Department ID"
// @Success 200 {object} entity.Department
// @Failure 404 {object} dto.ErrorResponse "Department not found"
// @Failure 500 {object} dto.ErrorResponse "Internal server error"
// @Router /departments/{id} [get]
func (h *DepartmentsHandler) GetDepartmentByID(c *gin.Context) {
	idStr := c.Param("id")
	depID, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid department ID"})
		return
	}

	department, err := h.DepartmentService.GetDepartmentByID(depID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch department"})
		return
	}
	if department == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Department not found"})
		return
	}

	c.JSON(http.StatusOK, department)
}
