package repository

import (
	"database/sql"
	"jumyste-app-backend/internal/entity"
	"jumyste-app-backend/pkg/logger"
)

type CompanyRepository struct {
	DB *sql.DB
}

func NewCompanyRepository(db *sql.DB) *CompanyRepository {
	return &CompanyRepository{DB: db}
}

func (r *CompanyRepository) Create(company *entity.Company) error {
	query := `INSERT INTO companies (name, owner_id, photo_url, description)
			  VALUES ($1, $2, $3, $4) RETURNING id`

	err := r.DB.QueryRow(query,
		company.Name,
		company.OwnerId,
		company.PhotoUrl,
		company.Description,
	).Scan(&company.ID)

	if err != nil {
		logger.Log.Error("Failed to create company", "error", err)
		return err
	}

	logger.Log.Info("Company created", "company_id", company.ID)
	return nil
}

func (r *CompanyRepository) GetByID(id int) (*entity.Company, error) {
	query := `SELECT id, name, owner_id, created_at, photo_url, description FROM companies WHERE id = $1`
	row := r.DB.QueryRow(query, id)

	var company entity.Company
	err := row.Scan(&company.ID, &company.Name, &company.OwnerId, &company.CreatedAt, &company.PhotoUrl, &company.Description)
	if err != nil {
		if err == sql.ErrNoRows {
			logger.Log.Warn("Company not found", "company_id", id)
			return nil, nil
		}
		logger.Log.Error("Failed to get company by ID", "error", err)
		return nil, err
	}

	logger.Log.Info("Company retrieved", "company_id", company.ID)
	return &company, nil
}

func (r *CompanyRepository) Update(company *entity.Company) error {
	query := `UPDATE companies SET name = $1, photo_url = $2, description = $3 WHERE id = $4`

	_, err := r.DB.Exec(query,
		company.Name,
		company.PhotoUrl,
		company.Description,
		company.ID,
	)

	if err != nil {
		logger.Log.Error("Failed to update company", "company_id", company.ID, "error", err)
		return err
	}

	logger.Log.Info("Company updated", "company_id", company.ID)
	return nil
}

func (r *CompanyRepository) Delete(id int) error {
	query := `DELETE FROM companies WHERE id = $1`

	_, err := r.DB.Exec(query, id)
	if err != nil {
		logger.Log.Error("Failed to delete company", "company_id", id, "error", err)
		return err
	}

	logger.Log.Info("Company deleted", "company_id", id)
	return nil
}
