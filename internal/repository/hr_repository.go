package repository

import (
	"database/sql"
	"jumyste-app-backend/internal/entity"
	"jumyste-app-backend/pkg/logger"
	"log/slog"
)

type HrRepository struct {
	DB *sql.DB
}

func NewHrRepository(db *sql.DB) *HrRepository {
	return &HrRepository{DB: db}
}

func (r *HrRepository) CreateHR(hr *entity.HR) error {
	query := `INSERT INTO hr(user_id, dep_id, company_id) 
	          VALUES ($1, $2, $3) RETURNING id`
	return r.DB.QueryRow(query, hr.UserID, hr.DepID, hr.CompanyID).Scan(&hr.ID)
}

func (r *HrRepository) GetHRByUserID(userID int) (*entity.HR, error) {
	query := `SELECT id, user_id, dep_id, company_id FROM hr WHERE user_id = $1`
	row := r.DB.QueryRow(query, userID)

	var hr entity.HR
	err := row.Scan(&hr.ID, &hr.UserID, &hr.DepID, &hr.CompanyID)
	if err != nil {
		return nil, err
	}
	return &hr, nil
}

func (r *HrRepository) UpdateHR(hr *entity.HR) error {
	query := `UPDATE hr SET dep_id = $1, company_id = $2 WHERE user_id = $3`
	_, err := r.DB.Exec(query, hr.DepID, hr.CompanyID, hr.UserID)
	if err != nil {
		logger.Log.Error("Failed to update HR data",
			slog.Int("user_id", hr.UserID),
			slog.String("error", err.Error()))
		return err
	}
	return nil
}

func (r *HrRepository) CheckHRByUserID(userID int) (bool, error) {
	query := `SELECT COUNT(*) FROM hr WHERE user_id = $1`
	var count int
	err := r.DB.QueryRow(query, userID).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *HrRepository) DeleteHR(userID int) error {
	query := `DELETE FROM hr WHERE user_id = $1`
	_, err := r.DB.Exec(query, userID)
	if err != nil {
		return err
	}
	return nil
}
