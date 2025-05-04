package service

import (
	"fmt"
	"jumyste-app-backend/internal/entity"
	"jumyste-app-backend/internal/repository"
	"jumyste-app-backend/pkg/logger"
)

type DepartmentsService struct {
	DepartmentRepo *repository.DepartmentsRepo
}

func NewDepartmentsService(departmentRepo *repository.DepartmentsRepo) *DepartmentsService {
	return &DepartmentsService{
		DepartmentRepo: departmentRepo,
	}
}

func (s *DepartmentsService) CreateDepartment(userID int, dep *entity.Department) error {
	logger.Log.Info("Service: creating department", "company_id", dep.CompanyId, "name", dep.Name, "user_id", userID)

	isOwner, err := s.DepartmentRepo.IsUserOwnerOfCompany(userID, dep.CompanyId)
	if err != nil {
		logger.Log.Error("Service error: failed to check ownership", "error", err)
		return err
	}
	if !isOwner {
		logger.Log.Warn("Unauthorized: user is not the owner of the company", "user_id", userID, "company_id", dep.CompanyId)
		return fmt.Errorf("user is not authorized to create department")
	}

	err = s.DepartmentRepo.CreateDepartment(dep)
	if err != nil {
		logger.Log.Error("Service error: failed to create department", "error", err)
	}
	return err
}

func (s *DepartmentsService) GetDepartmentsByCompany(companyID int) ([]*entity.Department, error) {
	logger.Log.Info("Service: fetching departments by company ID", "company_id", companyID)
	departments, err := s.DepartmentRepo.GetDepartmentsByCompanyID(companyID)
	if err != nil {
		logger.Log.Error("Service error: failed to get departments", "error", err)
	}
	return departments, err
}

func (s *DepartmentsService) GetDepartmentByID(depID int) (*entity.Department, error) {
	logger.Log.Info("Fetching department by ID", "dep_id", depID)
	return s.DepartmentRepo.GetDepartmentByID(depID)
}
