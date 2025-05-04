package repository

import (
	"database/sql"
	"errors"
	"jumyste-app-backend/internal/entity"
	"jumyste-app-backend/pkg/logger"
)

type DepartmentsRepo struct {
	DB *sql.DB
}

func NewDepartmentsRepo(db *sql.DB) *DepartmentsRepo {
	return &DepartmentsRepo{DB: db}
}

func (r *DepartmentsRepo) GetDepartmentsByCompanyID(companyID int) ([]*entity.Department, error) {
	logger.Log.Info("Fetching departments by company ID", "company_id", companyID)

	query := `SELECT id, color, company_id, name, hr_count FROM departments WHERE company_id = $1`
	rows, err := r.DB.Query(query, companyID)
	if err != nil {
		logger.Log.Error("Error querying departments", "company_id", companyID, "error", err)
		return nil, err
	}
	defer rows.Close()

	var departments []*entity.Department
	for rows.Next() {
		var dep entity.Department
		if err := rows.Scan(&dep.ID, &dep.Color, &dep.CompanyId, &dep.Name, &dep.HrCount); err != nil {
			logger.Log.Error("Error scanning department row", "error", err)
			return nil, err
		}
		departments = append(departments, &dep)
	}

	logger.Log.Info("Departments fetched successfully", "count", len(departments))
	return departments, nil
}

func (r *DepartmentsRepo) CreateDepartment(dep *entity.Department) error {
	logger.Log.Info("Creating new department", "company_id", dep.CompanyId, "name", dep.Name, "hr_count", dep.HrCount)

	query := `INSERT INTO departments (color, company_id, name, hr_count) VALUES ($1, $2, $3, $4) RETURNING id`
	err := r.DB.QueryRow(query, dep.Color, dep.CompanyId, dep.Name, dep.HrCount).Scan(&dep.ID)
	if err != nil {
		logger.Log.Error("Error creating department", "error", err)
		return err
	}

	logger.Log.Info("Department created successfully", "department_id", dep.ID)
	return nil
}

func (r *DepartmentsRepo) IsUserOwnerOfCompany(userID int, companyID int) (bool, error) {
	var exists bool
	query := `SELECT EXISTS (
		SELECT 1 FROM companies WHERE id = $1 AND owner_id = $2
	)`
	err := r.DB.QueryRow(query, companyID, userID).Scan(&exists)
	if err != nil {
		return false, err
	}
	if !exists {
		logger.Log.Warn("User is not owner of the company", "user_id", userID, "company_id", companyID)
	}
	return exists, nil
}

func (r *DepartmentsRepo) GetDepartmentByID(depID int) (*entity.Department, error) {
	query := `SELECT id, color, company_id, name, hr_count FROM departments WHERE id = $1`

	var dep entity.Department
	err := r.DB.QueryRow(query, depID).Scan(&dep.ID, &dep.Color, &dep.CompanyId, &dep.Name, &dep.HrCount)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		logger.Log.Error("Failed to fetch department by ID", "dep_id", depID, "error", err)
		return nil, err
	}

	return &dep, nil
}
