package repository

import (
	"database/sql"
	"jumyste-app-backend/internal/entity"
)

type CompanyRepository struct {
	DB *sql.DB
}

func NewCompanyRepository(db *sql.DB) *CompanyRepository {
	return &CompanyRepository{DB: db}
}

//func (r *CompanyRepository) Create(company *entity.Company) error {
//	query := `INSERT INTO companies (name, owner_id, created_at) VALUES ($1, $2, $3) RETURNING id`
//	return r.DB.QueryRow(query, company.Name, company.Description, time.Now()).Scan(&company.ID)
//}

func (r *CompanyRepository) GetByID(id int) (*entity.Company, error) {
	query := `SELECT id, name, owner_id, created_at, photo_url FROM companies WHERE id = $1`
	row := r.DB.QueryRow(query, id)

	var company entity.Company
	err := row.Scan(&company.ID, &company.Name, &company.OwnerId, &company.CreatedAt, &company.PhotoUrl)
	if err != nil {
		return nil, err
	}

	return &company, nil
}

//
//func (r *CompanyRepository) Update(company *entity.Company) error {
//	query := `UPDATE companies SET name = $1, description = $2 WHERE id = $3`
//	_, err := r.DB.Exec(query, company.Name, company.Description, company.ID)
//	return err
//}
//
//func (r *CompanyRepository) Delete(id int) error {
//	query := `DELETE FROM companies WHERE id = $1`
//	_, err := r.DB.Exec(query, id)
//	return err
//}
