package service

import (
	"jumyste-app-backend/internal/entity"
	"jumyste-app-backend/internal/repository"
)

type DepartmentsService struct {
	DepartmentRepo *repository.DepartmentsRepo
}

func NewDepartmentsService(departmentRepo *repository.DepartmentsRepo) *DepartmentsService {
	return &DepartmentsService{
		DepartmentRepo: departmentRepo,
	}
}

func (s *DepartmentsService) GetDepartmentsByCompany(companyID int) ([]*entity.Department, error) {
	return s.DepartmentRepo.GetDepartmentsByCompanyID(companyID)
}
