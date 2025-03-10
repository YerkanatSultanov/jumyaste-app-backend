package repository

import (
	"database/sql"
	"fmt"
	"github.com/lib/pq"
	"jumyste-app-backend/internal/entity"
	"jumyste-app-backend/pkg/logger"
	"log/slog"
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
        (title, employment_type, work_format, experience, salary_min, salary_max, location, category, skills, description, created_by) 
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
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
	).Scan(&v.ID, &v.CreatedAt)
}

func (r *VacancyRepository) UpdateVacancy(v *entity.Vacancy) error {
	query := `
        UPDATE vacancies 
        SET title = $1, employment_type = $2, work_format = $3, experience = $4, 
            salary_min = $5, salary_max = $6, location = $7, category = $8, 
            skills = $9, description = $10
        WHERE id = $11 AND created_by = $12
        RETURNING created_at`

	err := r.db.QueryRow(
		query, v.Title, v.EmploymentType, v.WorkFormat, v.Experience,
		v.SalaryMin, v.SalaryMax, v.Location, v.Category, pq.Array(v.Skills), v.Description,
		v.ID, v.CreatedBy,
	).Scan(&v.CreatedAt)

	return err
}

func (r *VacancyRepository) GetVacancyById(id int) (*entity.Vacancy, error) {
	query := `
	SELECT id, title, employment_type, work_format, experience, salary_min,
       salary_max, location, category, skills, description, created_by, created_at
	FROM vacancies WHERE id = $1`

	var vacancy entity.Vacancy
	err := r.db.QueryRow(query, id).Scan(
		&vacancy.ID, &vacancy.Title, &vacancy.EmploymentType, &vacancy.WorkFormat,
		&vacancy.Experience, &vacancy.SalaryMin, &vacancy.SalaryMax,
		&vacancy.Location, &vacancy.Category, pq.Array(&vacancy.Skills),
		&vacancy.Description, &vacancy.CreatedBy, &vacancy.CreatedAt,
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
	       salary_max, location, category, skills, description, created_by, created_at
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
			&v.SalaryMin, &v.SalaryMax, &v.Location, &v.Category, pq.Array(&v.Skills), &v.Description, &v.CreatedBy, &v.CreatedAt); err != nil {
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
		       salary_max, location, category, skills, description, created_by
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
			&v.SalaryMin, &v.SalaryMax, &v.Location, &v.Category, pq.Array(&v.Skills), &v.Description, &v.CreatedBy); err != nil {
			logger.Log.Error("Error scanning vacancy row", slog.String("error", err.Error()))
			return nil, err
		}
		vacancies = append(vacancies, &v)
	}

	return vacancies, nil
}

func (r *VacancyRepository) SearchVacancies(filter entity.VacancyFilter) ([]*entity.Vacancy, error) {
	query := `SELECT id, title, employment_type, work_format, experience, salary_min, salary_max, location, category, skills, description, created_by FROM vacancies WHERE 1=1`
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

	if filter.Query != "" {
		query += fmt.Sprintf(" AND search_vector @@ plainto_tsquery('russian', $%d)", argIndex)
		args = append(args, filter.Query)
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
		err := rows.Scan(&vacancy.ID, &vacancy.Title, &vacancy.EmploymentType, &vacancy.WorkFormat, &vacancy.Experience, &vacancy.SalaryMin, &vacancy.SalaryMax, &vacancy.Location, &vacancy.Category, pq.Array(&vacancy.Skills), &vacancy.Description, &vacancy.CreatedBy)
		if err != nil {
			return nil, err
		}
		vacancies = append(vacancies, &vacancy)
	}

	return vacancies, nil
}
