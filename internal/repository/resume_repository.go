package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/lib/pq"
	"jumyste-app-backend/internal/dto"
	"jumyste-app-backend/internal/entity"
	"jumyste-app-backend/pkg/logger"
)

type ResumeRepository struct {
	DB *sql.DB
}

func NewResumeRepository(db *sql.DB) *ResumeRepository {
	return &ResumeRepository{DB: db}
}

func (r *ResumeRepository) CreateResume(ctx context.Context, resume *entity.Resume) error {
	skillsArray := pq.Array(resume.Skills)

	parsedDataJSON, err := json.Marshal(resume.ParsedData)
	if err != nil {
		return err
	}

	query := `
		INSERT INTO resume (user_id, full_name, desired_position, skills, city, about, parsed_data)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id
	`

	err = r.DB.QueryRowContext(ctx, query,
		resume.UserID,
		resume.FullName,
		resume.DesiredPosition,
		skillsArray,
		resume.City,
		resume.About,
		parsedDataJSON,
	).Scan(&resume.ID)

	return err
}

func (r *ResumeRepository) GetResumeByUserID(ctx context.Context, userID int) (*entity.Resume, *entity.User, error) {
	var resume entity.Resume
	var user entity.User
	var parsedData []byte

	query := `
		SELECT 
			r.id, r.user_id, r.full_name, r.desired_position, r.skills, r.city, r.about, r.parsed_data, r.created_at,
			u.id, u.email, u.password, u.first_name, u.last_name, u.profile_picture, u.role_id
		FROM resume r
		JOIN users u ON r.user_id = u.id
		WHERE r.user_id = $1
		LIMIT 1
	`

	err := r.DB.QueryRowContext(ctx, query, userID).Scan(
		&resume.ID,
		&resume.UserID,
		&resume.FullName,
		&resume.DesiredPosition,
		pq.Array(&resume.Skills),
		&resume.City,
		&resume.About,
		&parsedData,
		&resume.CreatedAt,

		&user.ID,
		&user.Email,
		&user.Password,
		&user.FirstName,
		&user.LastName,
		&user.ProfilePicture,
		&user.RoleId,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil, nil
		}
		logger.Log.Error("Failed to get resume and user by user_id", "error", err)
		return nil, nil, err
	}

	if len(parsedData) > 0 {
		if err := json.Unmarshal(parsedData, &resume.ParsedData); err != nil {
			logger.Log.Error("Failed to unmarshal parsed_data", "error", err)
			return nil, nil, err
		}
	}

	return &resume, &user, nil
}

func (r *ResumeRepository) DeleteResumeByUserID(ctx context.Context, userID int) error {
	query := `DELETE FROM resume WHERE user_id = $1`
	_, err := r.DB.ExecContext(ctx, query, userID)
	if err != nil {
		logger.Log.Error("Failed to delete resume by user_id", "error", err)
		return err
	}
	return nil
}

func (r *ResumeRepository) GetByUserID(ctx context.Context, userId int) (*entity.Resume, error) {
	var resume entity.Resume
	var parsedData []byte

	query := `
		SELECT id, user_id, full_name, desired_position, skills, city, about, parsed_data, created_at
		FROM resume
		WHERE user_id = $1
	`

	err := r.DB.QueryRowContext(ctx, query, userId).Scan(
		&resume.ID,
		&resume.UserID,
		&resume.FullName,
		&resume.DesiredPosition,
		pq.Array(&resume.Skills),
		&resume.City,
		&resume.About,
		&parsedData,
		&resume.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		logger.Log.Error("Failed to get resume by ID", "error", err)
		return nil, err
	}

	if len(parsedData) > 0 {
		if err := json.Unmarshal(parsedData, &resume.ParsedData); err != nil {
			logger.Log.Error("Failed to unmarshal parsed_data in GetByID", "error", err)
			return nil, err
		}
	}

	return &resume, nil
}
func (r *ResumeRepository) GetByID(ctx context.Context, resumeId int) (*entity.Resume, error) {
	var resume entity.Resume
	var parsedData []byte

	query := `
		SELECT id, user_id, full_name, desired_position, skills, city, about, parsed_data, created_at
		FROM resume
		WHERE id = $1
	`

	err := r.DB.QueryRowContext(ctx, query, resumeId).Scan(
		&resume.ID,
		&resume.UserID,
		&resume.FullName,
		&resume.DesiredPosition,
		pq.Array(&resume.Skills),
		&resume.City,
		&resume.About,
		&parsedData,
		&resume.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		logger.Log.Error("Failed to get resume by ID", "error", err)
		return nil, err
	}

	if len(parsedData) > 0 {
		if err := json.Unmarshal(parsedData, &resume.ParsedData); err != nil {
			logger.Log.Error("Failed to unmarshal parsed_data in GetByID", "error", err)
			return nil, err
		}
	}

	return &resume, nil
}
func (r *ResumeRepository) FilterCandidates(ctx context.Context, filter dto.CandidateFilter) ([]entity.JobApplicationWithResume, error) {
	query := `
		SELECT
			ja.id, ja.user_id, ja.vacancy_id, ja.first_name, ja.last_name, ja.email, ja.status, ja.applied_at, ja.resume_id, ja.ai_matching_score,
			r.id, r.full_name, r.desired_position, r.skills, r.city, r.about, r.parsed_data, r.created_at,
			u.id, u.email, u.first_name, u.last_name, u.profile_picture, u.role_id
		FROM job_applications ja
		JOIN resume r ON ja.resume_id = r.id
		JOIN users u ON r.user_id = u.id
		WHERE 1=1
	`

	args := []interface{}{}
	argID := 1

	if filter.AIMatchMin > 0 {
		query += fmt.Sprintf(" AND ja.ai_matching_score >= $%d", argID)
		args = append(args, filter.AIMatchMin)
		argID++
	}

	if len(filter.Skills) > 0 {
		query += fmt.Sprintf(" AND r.skills && $%d", argID) // PostgreSQL array intersection
		args = append(args, pq.Array(filter.Skills))
		argID++
	}

	if filter.City != "" {
		query += fmt.Sprintf(" AND r.city ILIKE $%d", argID)
		args = append(args, "%"+filter.City+"%")
		argID++
	}

	if filter.Position != "" {
		query += fmt.Sprintf(" AND r.desired_position ILIKE $%d", argID)
		args = append(args, "%"+filter.Position+"%")
		argID++
	}

	rows, err := r.DB.QueryContext(ctx, query, args...)
	if err != nil {
		logger.Log.Error("Failed to filter candidates", "error", err)
		return nil, err
	}
	defer rows.Close()

	var results []entity.JobApplicationWithResume

	for rows.Next() {
		var app entity.JobApplication
		var resume entity.Resume
		var user entity.User

		var skills []string
		var parsedData []byte

		err := rows.Scan(
			&app.ID, &app.UserID, &app.VacancyID, &app.FirstName, &app.LastName, &app.Email, &app.Status, &app.AppliedAt, &app.ResumeID, &app.AIMatchingScore,
			&resume.ID, &resume.FullName, &resume.DesiredPosition, pq.Array(&skills), &resume.City, &resume.About, &parsedData, &resume.CreatedAt,
			&user.ID, &user.Email, &user.FirstName, &user.LastName, &user.ProfilePicture, &user.RoleId,
		)
		if err != nil {
			logger.Log.Error("Failed to scan row", "error", err)
			continue
		}

		resume.Skills = skills
		if len(parsedData) > 0 {
			if err := json.Unmarshal(parsedData, &resume.ParsedData); err != nil {
				logger.Log.Error("Failed to unmarshal parsed_data in GetByID", "error", err)
				return nil, err
			}
		}

		appWithResume := entity.JobApplicationWithResume{
			JobApplication: app,
			Resume:         resume,
			User:           user,
		}

		results = append(results, appWithResume)
	}

	return results, nil
}
