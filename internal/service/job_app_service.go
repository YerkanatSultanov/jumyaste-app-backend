package service

import (
	"context"
	"jumyste-app-backend/internal/entity"
	"jumyste-app-backend/internal/repository"
	"jumyste-app-backend/pkg/logger"
)

type JobApplicationService struct {
	JobApplicationRepo *repository.JobApplicationRepository
}

func NewJobApplicationService(repo *repository.JobApplicationRepository) *JobApplicationService {
	return &JobApplicationService{JobApplicationRepo: repo}
}

func (s *JobApplicationService) ApplyForJob(ctx context.Context, userID, vacancyID int, firstName, lastName, email string) (*entity.JobApplication, error) {
	application := &entity.JobApplication{
		UserID:    userID,
		VacancyID: vacancyID,
		FirstName: firstName,
		LastName:  lastName,
		Email:     email,
		Status:    "new",
	}

	err := s.JobApplicationRepo.CreateJobApplication(ctx, application)
	if err != nil {
		logger.Log.Error("Failed to apply for job", "error", err)
		return nil, err
	}

	return application, nil
}

func (s *JobApplicationService) GetJobApplicationsByVacancyID(ctx context.Context, vacancyID int) ([]entity.JobApplication, error) {
	applications, err := s.JobApplicationRepo.GetJobApplicationsByVacancyID(ctx, vacancyID)
	if err != nil {
		logger.Log.Error("Failed to get job applications", "error", err)
		return nil, err
	}
	return applications, nil
}

func (s *JobApplicationService) UpdateJobApplicationStatus(ctx context.Context, applicationID int, status string) error {
	err := s.JobApplicationRepo.UpdateJobApplicationStatus(ctx, applicationID, status)
	if err != nil {
		logger.Log.Error("Failed to update job application status", "error", err)
		return err
	}
	return nil
}

func (s *JobApplicationService) DeleteJobApplication(ctx context.Context, applicationID int) error {
	err := s.JobApplicationRepo.DeleteJobApplication(ctx, applicationID)
	if err != nil {
		logger.Log.Error("Failed to delete job application", "error", err)
		return err
	}
	return nil
}
