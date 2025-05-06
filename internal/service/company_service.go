package service

import (
	"errors"
	"jumyste-app-backend/internal/entity"
	"jumyste-app-backend/internal/repository"
	"jumyste-app-backend/pkg/logger"
	"log/slog"
)

type CompanyService struct {
	repo *repository.CompanyRepository
}

func NewCompanyService(repo *repository.CompanyRepository) *CompanyService {
	return &CompanyService{repo: repo}
}

func (s *CompanyService) CreateCompany(company *entity.Company) error {
	logger.Log.Info("Creating company", slog.String("name", company.Name))

	err := s.repo.Create(company)
	if err != nil {
		logger.Log.Error("Failed to create company", slog.String("error", err.Error()))
		return err
	}

	logger.Log.Info("Company created successfully", slog.Int("company_id", company.ID))
	return nil
}

func (s *CompanyService) GetCompanyByID(id int) (*entity.Company, error) {
	logger.Log.Info("Getting company by ID", slog.Int("company_id", id))

	company, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if company == nil {
		return nil, errors.New("company not found")
	}

	return company, nil
}

func (s *CompanyService) UpdateCompany(company *entity.Company) error {
	logger.Log.Info("Updating company", slog.Int("company_id", company.ID))

	err := s.repo.Update(company)
	if err != nil {
		logger.Log.Error("Failed to update company", slog.String("error", err.Error()))
		return err
	}

	logger.Log.Info("Company updated successfully", slog.Int("company_id", company.ID))
	return nil
}

func (s *CompanyService) DeleteCompany(id int) error {
	logger.Log.Info("Deleting company", slog.Int("company_id", id))

	err := s.repo.Delete(id)
	if err != nil {
		logger.Log.Error("Failed to delete company", slog.String("error", err.Error()))
		return err
	}

	logger.Log.Info("Company deleted successfully", slog.Int("company_id", id))
	return nil
}
