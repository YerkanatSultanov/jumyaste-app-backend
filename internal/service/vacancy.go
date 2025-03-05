package service

import (
	"errors"
	"jumyste-app-backend/internal/entity"
	"jumyste-app-backend/internal/repository"
	"jumyste-app-backend/pkg/logger"
	"log/slog"
)

type VacancyService struct {
	repo *repository.VacancyRepository
}

func NewVacancyService(repo *repository.VacancyRepository) *VacancyService {
	return &VacancyService{repo: repo}
}

func (s *VacancyService) CreateVacancy(v *entity.Vacancy) error {
	logger.Log.Info("Creating new vacancy", slog.String("title", v.Title), slog.Int("created_by", v.CreatedBy))

	err := s.repo.CreateVacancy(v)
	if err != nil {
		logger.Log.Error("Failed to create vacancy", slog.String("title", v.Title), slog.Int("created_by", v.CreatedBy), slog.String("error", err.Error()))
		return err
	}

	logger.Log.Info("Vacancy created successfully", slog.String("title", v.Title), slog.Int("created_by", v.CreatedBy))
	return nil
}

func (s *VacancyService) UpdateVacancy(v *entity.Vacancy) error {
	logger.Log.Info("Updating vacancy in service", slog.Int("vacancy_id", v.ID), slog.Int("user_id", v.CreatedBy))

	err := s.repo.UpdateVacancy(v)
	if err != nil {
		logger.Log.Error("Failed to update vacancy in repository", slog.String("error", err.Error()))
		return err
	}

	logger.Log.Info("Vacancy updated successfully in service", slog.Int("vacancy_id", v.ID))
	return nil
}

func (s *VacancyService) GetVacancyById(id int) (*entity.Vacancy, error) {
	logger.Log.Info("Getting vacancy by id", slog.Int("id", id))

	vac, err := s.repo.GetVacancyById(id)
	if err != nil {
		logger.Log.Error("Failed to find vacancy", slog.String("error", err.Error()))
		return nil, err
	}
	logger.Log.Info("Vacancy fetched successfully in service", slog.Int("vacancy_id", vac.ID))
	return vac, nil
}

func (s *VacancyService) DeleteVacancy(id, userID int) error {
	vacancy, err := s.repo.GetVacancyById(id)
	if err != nil {
		logger.Log.Error("Vacancy not found", slog.Int("vacancy_id", id), slog.String("error", err.Error()))
		return errors.New("vacancy not found")
	}

	if vacancy.CreatedBy != userID {
		logger.Log.Warn("Unauthorized vacancy delete attempt",
			slog.Int("vacancy_id", id), slog.Int("user_id", userID))
		return errors.New("you can only delete your own vacancies")
	}

	return s.repo.DeleteVacancy(id)
}

func (s *VacancyService) GetAllVacancies() ([]*entity.Vacancy, error) {
	logger.Log.Info("Fetching all vacancies")

	vacancies, err := s.repo.GetAllVacancies()
	if err != nil {
		logger.Log.Error("Failed to retrieve vacancies", slog.String("error", err.Error()))
		return nil, err
	}

	logger.Log.Info("Vacancies retrieved successfully", slog.Int("count", len(vacancies)))
	return vacancies, nil
}
