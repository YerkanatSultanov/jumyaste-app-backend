package repository

import (
	"context"
	"database/sql"
	"jumyste-app-backend/internal/entity"
	"jumyste-app-backend/pkg/logger"
)

type JobApplicationRepository struct {
	DB *sql.DB
}

func NewJobApplicationRepository(db *sql.DB) *JobApplicationRepository {
	return &JobApplicationRepository{DB: db}
}

func (r *JobApplicationRepository) CreateJobApplication(ctx context.Context, application *entity.JobApplication) error {
	query := `
        INSERT INTO job_applications (user_id, vacancy_id, first_name, last_name, email, status, resume_id, ai_matching_score)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
        RETURNING id, applied_at
    `
	err := r.DB.QueryRowContext(ctx, query,
		application.UserID,
		application.VacancyID,
		application.FirstName,
		application.LastName,
		application.Email,
		application.Status,
		application.ResumeID,
		application.AIMatchingScore,
	).Scan(&application.ID, &application.AppliedAt)
	return err
}

func (r *JobApplicationRepository) GetJobApplicationsByVacancyID(ctx context.Context, vacancyID int) ([]entity.JobApplication, error) {
	query := `
        SELECT id, user_id, vacancy_id, first_name, last_name, email, status, applied_at, resume_id, ai_matching_score
        FROM job_applications
        WHERE vacancy_id = $1
    `
	rows, err := r.DB.QueryContext(ctx, query, vacancyID)
	if err != nil {
		logger.Log.Error("Failed to get job applications", "error", err)
		return nil, err
	}
	defer rows.Close()

	var applications []entity.JobApplication
	for rows.Next() {
		var application entity.JobApplication
		if err := rows.Scan(&application.ID, &application.UserID, &application.VacancyID, &application.FirstName, &application.LastName, &application.Email, &application.Status, &application.AppliedAt, &application.ResumeID, &application.AIMatchingScore); err != nil {
			logger.Log.Error("Failed to scan job application", "error", err)
			return nil, err
		}
		applications = append(applications, application)
	}
	return applications, nil
}

func (r *JobApplicationRepository) UpdateJobApplicationStatus(ctx context.Context, applicationID int, status string) error {
	query := `
        UPDATE job_applications
        SET status = $1
        WHERE id = $2
    `
	_, err := r.DB.ExecContext(ctx, query, status, applicationID)
	if err != nil {
		logger.Log.Error("Failed to update job application status", "error", err)
		return err
	}
	return nil
}

func (r *JobApplicationRepository) DeleteJobApplication(ctx context.Context, applicationID int) error {
	query := `DELETE FROM job_applications WHERE id = $1`
	_, err := r.DB.ExecContext(ctx, query, applicationID)
	if err != nil {
		logger.Log.Error("Failed to delete job application", "error", err)
		return err
	}
	return nil
}

func (r *JobApplicationRepository) GetJobAppAnalytics(ctx context.Context, hrID int) ([]entity.ApplicationStatusStat, error) {
	query := `
		SELECT ja.status, COUNT(*) as count
		FROM job_applications ja
		JOIN vacancies v ON ja.vacancy_id = v.id
		WHERE v.created_by = $1
		GROUP BY ja.status;
	`

	rows, err := r.DB.QueryContext(ctx, query, hrID)
	if err != nil {
		logger.Log.Error("Failed to get job application status stats", "error", err)
		return nil, err
	}
	defer rows.Close()

	var stats []entity.ApplicationStatusStat
	for rows.Next() {
		var stat entity.ApplicationStatusStat
		if err := rows.Scan(&stat.Status, &stat.Count); err != nil {
			return nil, err
		}
		stats = append(stats, stat)
	}

	if len(stats) == 0 {
		return nil, nil
	}

	return stats, nil
}
