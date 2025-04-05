package repository

import (
	"database/sql"
	"errors"
	"jumyste-app-backend/pkg/logger"
)

type InvitationRepository struct {
	db *sql.DB
}

func NewInvitationRepository(db *sql.DB) *InvitationRepository {
	return &InvitationRepository{db: db}
}

func (r *InvitationRepository) CheckInvitationExists(email string) (bool, error) {
	var exists bool
	query := `SELECT EXISTS (SELECT 1 FROM invitations WHERE email = $1)`
	err := r.db.QueryRow(query, email).Scan(&exists)
	if err != nil {
		logger.Log.Error("Error checking invitation existence", "email", email, "error", err)
		return false, err
	}
	return exists, nil
}

func (r *InvitationRepository) DeleteInvitation(email string) error {
	query := `DELETE FROM invitations WHERE email = $1`
	_, err := r.db.Exec(query, email)
	if err != nil {
		logger.Log.Error("Error deleting invitation", "email", email, "error", err)
	}
	return err
}

func (r *InvitationRepository) CreateInvitation(email string, companyId, departmentId int) error {
	var id int
	query := `INSERT INTO invitations(email, company_id, dep_id) VALUES ($1, $2, $3) RETURNING id`
	err := r.db.QueryRow(query, email, companyId, departmentId).Scan(&id)
	if err != nil {
		logger.Log.Error("Error creating invitation", "email", email, "company_id", companyId, "dep_id", departmentId, "error", err)
		return err
	}
	logger.Log.Info("Invitation created successfully", "id", id, "email", email, "company_id", companyId, "dep_id", departmentId)
	return nil
}

func (r *InvitationRepository) GetInvitationByEmail(email string) (*struct {
	CompanyID int
	DepID     int
}, error) {
	var invitation struct {
		CompanyID int
		DepID     int
	}
	query := `SELECT company_id, dep_id FROM invitations WHERE email = $1`
	err := r.db.QueryRow(query, email).Scan(&invitation.CompanyID, &invitation.DepID)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		logger.Log.Error("Error fetching invitation", "email", email, "error", err)
		return nil, err
	}
	return &invitation, nil
}
