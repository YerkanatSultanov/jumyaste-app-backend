package service

import (
	"errors"
	"fmt"
	"jumyste-app-backend/internal/dto"
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

func (s *VacancyService) GetVacancyById(id int, allowClosed bool) (*entity.Vacancy, error) {
	logger.Log.Info("Getting vacancy by id", slog.Int("id", id))

	vac, err := s.repo.GetVacancyById(id)
	if err != nil {
		logger.Log.Error("Failed to find vacancy", slog.String("error", err.Error()))
		return nil, err
	}

	count, err := s.repo.CountResponses(id)
	if err != nil {
		logger.Log.Error("Failed to count responses", slog.String("error", err.Error()))
		return nil, err
	}
	vac.CountResponses = count

	if vac.Status == "closed" && !allowClosed {
		logger.Log.Error("Vacancy is closed", slog.String("title", vac.Title))
		return nil, errors.New("vacancy is closed")
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

func (s *VacancyService) GetMyVacancies(userID int) ([]*entity.Vacancy, error) {
	logger.Log.Info("Fetching vacancies for HR", slog.Int("user_id", userID))

	vacancies, err := s.repo.GetVacanciesByRecruiterID(userID)
	if err != nil {
		logger.Log.Error("Failed to retrieve vacancies", slog.String("error", err.Error()))
		return nil, err
	}
	logger.Log.Info("Vacancies retrieved successfully", slog.Int("count", len(vacancies)))

	return vacancies, nil
}

func (s *VacancyService) SearchVacancies(filter entity.VacancyFilter) ([]*entity.Vacancy, error) {
	logger.Log.Info("Searching vacancies with filters", "filter", filter)

	vacancies, err := s.repo.SearchVacancies(filter)
	if err != nil {
		logger.Log.Error("Failed to search vacancies", "error", err)
		return nil, err
	}

	return vacancies, nil
}

func (s *VacancyService) GetVacanciesByCompanyId(companyId int) ([]*entity.Vacancy, error) {
	logger.Log.Info("Fetching vacancies for company", slog.Int("company_id", companyId))

	vacancies, err := s.repo.GetVacanciesByCompany(companyId)
	if err != nil {
		logger.Log.Error("Failed to retrieve vacancies for company", slog.String("error", err.Error()))
		return nil, err
	}

	return vacancies, nil
}

func (s *VacancyService) UpdateVacancyStatus(vacancyId int, status string) error {
	logger.Log.Info("Updating vacancy status", slog.Int("vacancy_id", vacancyId), slog.String("new_status", status))

	validStatuses := []string{"open", "closed"}
	isValid := false
	for _, validStatus := range validStatuses {
		if status == validStatus {
			isValid = true
			break
		}
	}

	if !isValid {
		logger.Log.Error("Invalid status provided", slog.Int("vacancy_id", vacancyId), slog.String("status", status))
		return fmt.Errorf("invalid status: %s", status)
	}

	err := s.repo.UpdateStatus(vacancyId, status)
	if err != nil {
		logger.Log.Error("Failed to update vacancy status", slog.Int("vacancy_id", vacancyId), slog.String("error", err.Error()))
		return err
	}

	logger.Log.Info("Successfully updated vacancy status", slog.Int("vacancy_id", vacancyId), slog.String("new_status", status))

	return nil
}

func (s *VacancyService) GetFeedData(userID int) (dto.FeedDataResponse, error) {
	lastViewedAt, err := s.repo.GetFeedLastViewedAt(userID)
	if err != nil {
		return dto.FeedDataResponse{}, err
	}

	newVacanciesCount, err := s.repo.CountNewVacancies(userID, lastViewedAt)
	if err != nil {
		return dto.FeedDataResponse{}, err
	}

	err = s.repo.UpdateFeedLastViewedAt(userID)
	if err != nil {
		return dto.FeedDataResponse{}, err
	}

	return dto.FeedDataResponse{
		NewVacanciesCount: newVacanciesCount,
	}, nil
}
