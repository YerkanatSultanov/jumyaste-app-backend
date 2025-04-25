package service

import (
	"context"
	"fmt"
	"jumyste-app-backend/internal/ai"
	"jumyste-app-backend/internal/dto"
	"jumyste-app-backend/internal/entity"
	"jumyste-app-backend/internal/repository"
	"jumyste-app-backend/pkg/logger"
	"strings"
)

type JobApplicationService struct {
	JobApplicationRepo *repository.JobApplicationRepository
	ResumeRepo         *repository.ResumeRepository
	VacancyRepo        *repository.VacancyRepository
	AIClient           *ai.OpenAIClient
}

func NewJobApplicationService(repo *repository.JobApplicationRepository,
	resumeRepo *repository.ResumeRepository,
	vacancyRepo *repository.VacancyRepository,
	aiClient *ai.OpenAIClient) *JobApplicationService {
	return &JobApplicationService{JobApplicationRepo: repo,
		ResumeRepo:  resumeRepo,
		VacancyRepo: vacancyRepo,
		AIClient:    aiClient}
}

func (s *JobApplicationService) ApplyForJob(
	ctx context.Context,
	userID, vacancyID int,
	firstName, lastName, email string,
	resumeID int,
) (*entity.JobApplication, error) {
	resume, err := s.ResumeRepo.GetByID(ctx, userID)
	if err != nil {
		logger.Log.Error("Failed to get resume", "error", err)
		return nil, err
	}

	vacancy, err := s.VacancyRepo.GetVacancyById(vacancyID)
	if err != nil {
		logger.Log.Error("Failed to get vacancy", "error", err)
		return nil, err
	}
	resumeText := fmt.Sprintf("Имя: %s\nДолжность: %s\nНавыки: %s\nГород: %s\nО себе: %s",
		resume.FullName,
		resume.DesiredPosition,
		strings.Join(resume.Skills, ", "),
		resume.City,
		resume.About,
	)

	vacancyText := fmt.Sprintf("Название: %s\nТип занятости: %s\nФормат работы: %s\nНавыки: %s\nЛокация: %s\nОпыт: %s\nЗарплата: от %d до %d",
		vacancy.Title,
		vacancy.EmploymentType,
		vacancy.WorkFormat,
		strings.Join(vacancy.Skills, ", "),
		vacancy.Location,
		vacancy.Experience,
		vacancy.SalaryMin,
		vacancy.SalaryMax,
	)

	matchScore, err := s.AIClient.GetMatchingScore(resumeText, vacancyText)
	if err != nil {
		logger.Log.Error("AI matching failed", "error", err)
		return nil, err
	}

	application := &entity.JobApplication{
		UserID:          userID,
		VacancyID:       vacancyID,
		ResumeID:        resumeID,
		FirstName:       firstName,
		LastName:        lastName,
		Email:           email,
		Status:          "new",
		AIMatchingScore: matchScore,
	}

	err = s.JobApplicationRepo.CreateJobApplication(ctx, application)
	if err != nil {
		logger.Log.Error("Failed to save application", "error", err)
		return nil, err
	}

	return application, nil
}

func (s *JobApplicationService) GetJobApplicationsByVacancyID(ctx context.Context, vacancyID int) ([]dto.JobApplicationWithResumeResponse, error) {
	applications, err := s.JobApplicationRepo.GetJobApplicationsByVacancyID(ctx, vacancyID)
	if err != nil {
		logger.Log.Error("Failed to get job applications", "error", err)
		return nil, err
	}

	var response []dto.JobApplicationWithResumeResponse
	for _, app := range applications {
		resume, user, err := s.ResumeRepo.GetResumeByUserID(ctx, app.ResumeID)
		if err != nil {
			logger.Log.Error("Failed to get resume for application", "error", err)
			return nil, err
		}

		appResponse := dto.JobApplicationWithResumeResponse{
			ID:              app.ID,
			UserID:          app.UserID,
			VacancyID:       app.VacancyID,
			FirstName:       user.FirstName,
			LastName:        user.LastName,
			Email:           user.Email,
			Status:          app.Status,
			AppliedAt:       app.AppliedAt.Format("2006-01-02 15:04:05"),
			AIMatchingScore: app.AIMatchingScore,
			Resume: dto.ResumeResponse{
				FullName:        resume.FullName,
				DesiredPosition: resume.DesiredPosition,
				Skills:          resume.Skills,
				City:            resume.City,
				About:           resume.About,
				ParsedData:      resume.ParsedData,
				User: dto.UserResponse{
					ID:             user.ID,
					Email:          user.Email,
					FirstName:      user.FirstName,
					LastName:       user.LastName,
					ProfilePicture: user.ProfilePicture,
					RoleID:         user.RoleId,
				},
			},
		}

		response = append(response, appResponse)
	}

	return response, nil
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

func (s *JobApplicationService) GetJobAppAnalytics(ctx context.Context, hrID int) ([]entity.ApplicationStatusStat, error) {
	logger.Log.Info("Fetching job application analytics", "hr_id", hrID)

	stats, err := s.JobApplicationRepo.GetJobAppAnalytics(ctx, hrID)
	if err != nil {
		logger.Log.Error("Failed to get job application analytics", "hr_id", hrID, "error", err)
		return nil, err
	}

	totalCount := 0
	for _, stat := range stats {
		totalCount += stat.Count
	}

	if totalCount == 0 {
		logger.Log.Info("No job applications found for analytics", "hr_id", hrID)
		return stats, nil
	}

	for i := range stats {
		stats[i].Percentage = float64(stats[i].Count) / float64(totalCount) * 100
	}

	logger.Log.Info("Successfully fetched job application analytics", "hr_id", hrID, "total_count", totalCount)
	return stats, nil
}
