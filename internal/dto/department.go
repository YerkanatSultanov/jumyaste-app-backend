package dto

type CreateDepartmentRequest struct {
	Name string `json:"name" binding:"required"`
}
