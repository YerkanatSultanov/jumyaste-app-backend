package handler

import (
	"github.com/gin-gonic/gin"
	_ "jumyste-app-backend/internal/dto"
	_ "jumyste-app-backend/internal/entity"
	"jumyste-app-backend/internal/service"
	"net/http"
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
// @Description Retrieves a list of all vacancies
// @Tags Departments
// @Produce json
// @Security BearerAuth
// @Success 200 {array} entity.Department "List of departments"
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch departments"})
		return
	}

	c.JSON(http.StatusOK, departments)
}
