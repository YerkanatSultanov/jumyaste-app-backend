package service

import (
	"fmt"
	"jumyste-app-backend/internal/repository"
	"jumyste-app-backend/pkg/logger"
	"jumyste-app-backend/pkg/mail"
)

type InvitationService struct {
	repo *repository.InvitationRepository
}

func NewInvitationService(repo *repository.InvitationRepository) *InvitationService {
	return &InvitationService{repo: repo}
}

func (s *InvitationService) SendInvitation(email string, companyID, depID int) error {
	logger.Log.Info("Attempting to send an invitation", "email", email, "company_id", companyID, "dep_id", depID)

	exists, err := s.repo.CheckInvitationExists(email)
	if err != nil {
		logger.Log.Error("Error checking invitation existence", "error", err)
		return err
	}
	if exists {
		logger.Log.Warn("Invitation already exists", "email", email)
		return fmt.Errorf("invitation has already been sent")
	}

	err = s.repo.CreateInvitation(email, companyID, depID)
	if err != nil {
		logger.Log.Error("Error creating invitation", "error", err)
		return err
	}

	subject := "Company Invitation"
	body := fmt.Sprintf("You have been invited to a company. Register at: https://jumyste.click")

	err = mail.SendEmail(email, subject, body)
	if err != nil {
		logger.Log.Error("Error sending email", "error", err)
		return err
	}

	logger.Log.Info("Invitation successfully sent", "email", email)
	return nil
}
