package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/lib/pq"
	"jumyste-app-backend/internal/entity"
	"jumyste-app-backend/pkg/logger"
	"log/slog"
	"time"
)

type VacancyRepository struct {
	db *sql.DB
}

func NewVacancyRepository(db *sql.DB) *VacancyRepository {
	return &VacancyRepository{db: db}
}

func (r *VacancyRepository) CreateVacancy(v *entity.Vacancy) error {
	query := `
        INSERT INTO vacancies 
        (title, employment_type, work_format, experience, salary_min, salary_max, location, category, skills, description, created_by, company_id, status) 
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
        RETURNING id, created_at`

	return r.db.QueryRow(
		query,
		v.Title,
		v.EmploymentType,
		v.WorkFormat,
		v.Experience,
		v.SalaryMin,
		v.SalaryMax,
		v.Location,
		v.Category,
		pq.Array(v.Skills),
		v.Description,
		v.CreatedBy,
		v.CompanyId,
		v.Status,
	).Scan(&v.ID, &v.CreatedAt)
}

func (r *VacancyRepository) UpdateVacancy(v *entity.Vacancy) error {
	query := `
        UPDATE vacancies 
        SET title = $1, employment_type = $2, work_format = $3, experience = $4, 
            salary_min = $5, salary_max = $6, location = $7, category = $8, 
            skills = $9, description = $10, status = $11
        WHERE id = $11 AND created_by = $12
        RETURNING created_at`

	err := r.db.QueryRow(
		query, v.Title, v.EmploymentType, v.WorkFormat, v.Experience,
		v.SalaryMin, v.SalaryMax, v.Location, v.Category, pq.Array(v.Skills), v.Description,
		v.ID, v.CreatedBy, v.Status,
	).Scan(&v.CreatedAt)

	return err
}

func (r *VacancyRepository) GetVacancyById(id int) (*entity.Vacancy, error) {
	query := `
	SELECT id, title, employment_type, work_format, experience, salary_min,
       salary_max, location, category, skills, description, created_by, created_at, company_id, status
	FROM vacancies WHERE id = $1`

	var vacancy entity.Vacancy
	err := r.db.QueryRow(query, id).Scan(
		&vacancy.ID, &vacancy.Title, &vacancy.EmploymentType, &vacancy.WorkFormat,
		&vacancy.Experience, &vacancy.SalaryMin, &vacancy.SalaryMax,
		&vacancy.Location, &vacancy.Category, pq.Array(&vacancy.Skills),
		&vacancy.Description, &vacancy.CreatedBy, &vacancy.CreatedAt, &vacancy.CompanyId, &vacancy.Status,
	)
	if err != nil {
		logger.Log.Error("Vacancy not found",
			slog.Int("id", id),
			slog.String("error", err.Error()))
		return nil, err
	}
	logger.Log.Info("Vacancy retrieved successfully", slog.Int("email", id))
	return &vacancy, nil
}

func (r *VacancyRepository) DeleteVacancy(id int) error {
	query := "DELETE FROM vacancies WHERE id = $1"
	_, err := r.db.Exec(query, id)
	if err != nil {
		logger.Log.Error("Failed to delete vacancy", slog.Int("vacancy_id", id), slog.String("error", err.Error()))
		return err
	}

	logger.Log.Info("Vacancy deleted successfully", slog.Int("vacancy_id", id))
	return nil
}

func (r *VacancyRepository) GetAllVacancies() ([]*entity.Vacancy, error) {
	query := `
	SELECT id, title, employment_type, work_format, experience, salary_min, 
	       salary_max, location, category, skills, description, created_by, created_at, company_id, status
	FROM vacancies`

	rows, err := r.db.Query(query)
	if err != nil {
		logger.Log.Error("Failed to retrieve vacancies", slog.String("error", err.Error()))
		return nil, err
	}
	defer rows.Close()

	var vacancies []*entity.Vacancy
	for rows.Next() {
		v := &entity.Vacancy{}
		if err := rows.Scan(&v.ID, &v.Title, &v.EmploymentType, &v.WorkFormat, &v.Experience,
			&v.SalaryMin, &v.SalaryMax, &v.Location, &v.Category, pq.Array(&v.Skills), &v.Description, &v.CreatedBy, &v.CreatedAt, &v.CompanyId, &v.Status); err != nil {
			logger.Log.Error("Failed to scan vacancy row", slog.String("error", err.Error()))
			return nil, err
		}
		vacancies = append(vacancies, v)
	}

	if err := rows.Err(); err != nil {
		logger.Log.Error("Error iterating through vacancies", slog.String("error", err.Error()))
		return nil, err
	}

	logger.Log.Info("Vacancies retrieved successfully", slog.Int("count", len(vacancies)))
	return vacancies, nil
}

func (r *VacancyRepository) GetVacanciesByRecruiterID(userID int) ([]*entity.Vacancy, error) {
	query := `
		SELECT id, title, employment_type, work_format, experience, salary_min, 
		       salary_max, location, category, skills, description, created_by, created_at, company_id, status
		FROM vacancies 
		WHERE created_by = $1`

	rows, err := r.db.Query(query, userID)
	if err != nil {
		logger.Log.Error("Failed to fetch HR vacancies", slog.Int("user_id", userID), slog.String("error", err.Error()))
		return nil, err
	}
	defer rows.Close()

	var vacancies []*entity.Vacancy
	for rows.Next() {
		var v entity.Vacancy
		if err := rows.Scan(&v.ID, &v.Title, &v.EmploymentType, &v.WorkFormat, &v.Experience,
			&v.SalaryMin, &v.SalaryMax, &v.Location, &v.Category, pq.Array(&v.Skills), &v.Description, &v.CreatedBy, &v.CreatedAt, &v.CompanyId, &v.Status); err != nil {
			logger.Log.Error("Error scanning vacancy row", slog.String("error", err.Error()))
			return nil, err
		}
		vacancies = append(vacancies, &v)
	}

	return vacancies, nil
}

func (r *VacancyRepository) SearchVacancies(filter entity.VacancyFilter) ([]*entity.Vacancy, error) {
	query := `SELECT id, title, employment_type, work_format, experience, salary_min, salary_max, location, category, skills, description, created_by, company_id, status FROM vacancies WHERE 1=1`
	var args []interface{}
	argIndex := 1

	if filter.Title != "" {
		query += fmt.Sprintf(" AND title ILIKE $%d", argIndex)
		args = append(args, "%"+filter.Title+"%")
		argIndex++
	}

	if len(filter.Skills) > 0 {
		query += fmt.Sprintf(" AND skills @> $%d", argIndex)
		args = append(args, pq.Array(filter.Skills))
		argIndex++
	}

	if filter.Experience != "" {
		query += fmt.Sprintf(" AND experience = $%d", argIndex)
		args = append(args, filter.Experience)
		argIndex++
	}

	if len(filter.EmploymentType) > 0 {
		query += fmt.Sprintf(" AND employment_type = ANY($%d)", argIndex)
		args = append(args, pq.Array(filter.EmploymentType))
		argIndex++
	}

	if len(filter.WorkFormat) > 0 {
		query += fmt.Sprintf(" AND work_format = ANY($%d)", argIndex)
		args = append(args, pq.Array(filter.WorkFormat))
		argIndex++
	}

	if filter.Location != "" {
		query += fmt.Sprintf(" AND location ILIKE $%d", argIndex)
		args = append(args, "%"+filter.Location+"%")
		argIndex++
	}

	if filter.CompanyId != 0 {
		query += fmt.Sprintf(" AND company_id = $%d", argIndex)
		args = append(args, filter.CompanyId)
		argIndex++
	}

	if filter.Query != "" {
		query += fmt.Sprintf(" AND search_vector @@ plainto_tsquery('russian', $%d)", argIndex)
		args = append(args, filter.Query)
		argIndex++
	}

	if filter.Status != "" && filter.Status != "all" {
		query += fmt.Sprintf(" AND status = $%d", argIndex)
		args = append(args, filter.Status)
		argIndex++
	}

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var vacancies []*entity.Vacancy
	for rows.Next() {
		var vacancy entity.Vacancy
		err := rows.Scan(
			&vacancy.ID,
			&vacancy.Title,
			&vacancy.EmploymentType,
			&vacancy.WorkFormat,
			&vacancy.Experience,
			&vacancy.SalaryMin,
			&vacancy.SalaryMax,
			&vacancy.Location,
			&vacancy.Category,
			pq.Array(&vacancy.Skills),
			&vacancy.Description,
			&vacancy.CreatedBy,
			&vacancy.CompanyId,
			&vacancy.Status,
		)
		if err != nil {
			return nil, err
		}
		vacancies = append(vacancies, &vacancy)
	}

	return vacancies, nil
}

func (r *VacancyRepository) GetVacanciesByCompany(companyID int) ([]*entity.Vacancy, error) {
	query := `SELECT id, title, employment_type, work_format, experience, salary_min, 
                  salary_max, location, category, skills, description, created_by, created_at , company_id, status
                  FROM vacancies WHERE company_id = $1`

	rows, err := r.db.Query(query, companyID)
	if err != nil {
		logger.Log.Error("Failed to retrieve vacancies by company", slog.Int("company_id", companyID), slog.String("error", err.Error()))
		return nil, err
	}
	defer rows.Close()

	var vacancies []*entity.Vacancy
	for rows.Next() {
		v := &entity.Vacancy{}
		if err := rows.Scan(&v.ID, &v.Title, &v.EmploymentType, &v.WorkFormat, &v.Experience,
			&v.SalaryMin, &v.SalaryMax, &v.Location, &v.Category, pq.Array(&v.Skills),
			&v.Description, &v.CreatedBy, &v.CreatedAt, &v.CompanyId, &v.Status); err != nil {
			logger.Log.Error("Failed to scan vacancy row", slog.String("error", err.Error()))
			return nil, err
		}
		vacancies = append(vacancies, v)
	}

	if err := rows.Err(); err != nil {
		logger.Log.Error("Error iterating through vacancies", slog.String("error", err.Error()))
		return nil, err
	}

	logger.Log.Info("Vacancies retrieved successfully", slog.Int("company_id", companyID), slog.Int("count", len(vacancies)))
	return vacancies, nil
}

func (r *VacancyRepository) UpdateStatus(vacancyID int, status string) error {
	logger.Log.Info("Updating status", slog.Int("vacancy_id", vacancyID))
	query := `UPDATE vacancies SET status = $1 WHERE id = $2`
	_, err := r.db.Exec(query, status, vacancyID)
	return err
}

func (r *VacancyRepository) CountResponses(vacancyId int) (int, error) {
	var count int
	logger.Log.Info("Counting responses", slog.Int("vacancy_id", vacancyId))
	query := `SELECT COUNT(*) FROM job_applications WHERE vacancy_id = $1`
	err := r.db.QueryRow(query, vacancyId).Scan(&count)
	return count, err
}

func (r *VacancyRepository) GetFeedLastViewedAt(userID int) (time.Time, error) {
	var lastViewed time.Time

	err := r.db.QueryRow(`
		SELECT last_viewed_at 
		FROM vacancy_feed_views 
		WHERE user_id = $1
	`, userID).Scan(&lastViewed)

	if errors.Is(err, sql.ErrNoRows) {
		logger.Log.Info("No last viewed record found, inserting new", slog.Int("user_id", userID))

		_, err = r.db.Exec(`
			INSERT INTO vacancy_feed_views (user_id, last_viewed_at)
			VALUES ($1, NOW())
		`, userID)

		if err != nil {
			return time.Time{}, err
		}

		return time.Time{}, nil
	}

	return lastViewed, err
}

func (r *VacancyRepository) UpdateFeedLastViewedAt(userID int) error {
	logger.Log.Info("Updating feed last viewed at", slog.Int("user_id", userID))

	_, err := r.db.Exec(`
		UPDATE vacancy_feed_views
		SET last_viewed_at = NOW()
		WHERE user_id = $1
	`, userID)

	return err
}

func (r *VacancyRepository) CountNewVacancies(userID int, lastViewedAt time.Time) (int, error) {
	var count int

	logger.Log.Info("Counting new vacancies", slog.Int("user_id", userID), slog.Any("last_viewed_at", lastViewedAt))

	query := `
		SELECT COUNT(*) 
		FROM vacancies 
		WHERE created_at > $1 
		  AND status = 'open'
	`

	err := r.db.QueryRow(query, lastViewedAt).Scan(&count)

	return count, err
}
