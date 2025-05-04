package repository

import (
	"database/sql"
	"jumyste-app-backend/internal/entity"
)

type DepartmentsRepo struct {
	DB *sql.DB
}

func NewDepartmentsRepo(db *sql.DB) *DepartmentsRepo {
	return &DepartmentsRepo{DB: db}
}

func (r *DepartmentsRepo) GetDepartmentsByCompanyID(companyID int) ([]*entity.Department, error) {
	query := `SELECT id, company_id, name, hr_count FROM departments WHERE company_id = $1`

	rows, err := r.DB.Query(query, companyID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var departments []*entity.Department
	for rows.Next() {
		var dep entity.Department
		if err := rows.Scan(&dep.ID, &dep.CompanyId, &dep.Name, &dep.HrCount); err != nil {
			return nil, err
		}
		departments = append(departments, &dep)
	}

	return departments, nil
}
