package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
	"jumyste-app-backend/internal/entity"
	"jumyste-app-backend/internal/repository"
	"jumyste-app-backend/pkg/logger"
	"jumyste-app-backend/pkg/mail"
	"jumyste-app-backend/utils"
	"log/slog"
	"strconv"
	"time"
)

type AuthService struct {
	repo           *repository.AuthRepository
	redis          *redis.Client
	invitationRepo *repository.InvitationRepository
	hrRepo         *repository.HrRepository
}

func NewAuthService(repo *repository.AuthRepository, redis *redis.Client, invitationRepo *repository.InvitationRepository, hrRepo *repository.HrRepository) *AuthService {
	return &AuthService{
		repo:           repo,
		redis:          redis,
		invitationRepo: invitationRepo,
		hrRepo:         hrRepo,
	}
}

var (
	ErrInvalidResetCode = errors.New("invalid or expired reset code")
	ErrUserNotFound     = errors.New("user not found")
)

func (s *AuthService) LoginUser(email, password string) (string, string, error) {
	logger.Log.Info("Attempting user login", slog.String("email", email))

	user, err := s.repo.GetUserByEmail(email)
	if err != nil {
		logger.Log.Warn("Login failed: user not found", slog.String("email", email), slog.String("error", err.Error()))
		return "", "", errors.New("invalid credentials")
	}

	if !utils.CheckPassword(password, user.Password) {
		logger.Log.Warn("Login failed: incorrect password", slog.String("email", email))
		return "", "", errors.New("invalid credentials")
	}

	hr, err := s.hrRepo.GetHRByUserID(user.ID)
	if err != nil {
		logger.Log.Warn("HR data not found for user", slog.Int("user_id", user.ID), slog.String("error", err.Error()))
		user.CompanyID = 0
		user.DepID = 0
	} else {
		user.CompanyID = hr.CompanyID
		user.DepID = hr.DepID
	}

	accessToken, err := utils.GenerateJWT(user.ID, user.RoleId, user.CompanyID, user.DepID)
	if err != nil {
		logger.Log.Error("Failed to generate JWT", slog.String("email", email), slog.String("error", err.Error()))
		return "", "", err
	}

	refreshToken, err := utils.GenerateRefreshToken(user.ID, user.RoleId, user.CompanyID, user.DepID)
	if err != nil {
		logger.Log.Error("Failed to generate refresh token", slog.String("email", email), slog.String("error", err.Error()))
		return "", "", err
	}

	err = s.SaveRefreshToken(user.ID, refreshToken)
	if err != nil {
		logger.Log.Error("Failed to save refresh token to Redis", slog.String("email", email), slog.String("error", err.Error()))
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func (s *AuthService) RequestPasswordReset(email string) error {
	logger.Log.Info("Processing password reset request", slog.String("email", email))

	user, err := s.repo.GetUserByEmail(email)
	if err != nil {
		logger.Log.Warn("User not found", slog.String("email", email))
		return ErrUserNotFound
	}

	resetCode := utils.GenerateResetCode()
	expiresAt := time.Now().Add(10 * time.Minute)

	err = s.repo.SavePasswordResetCode(user.ID, resetCode, expiresAt)
	if err != nil {
		logger.Log.Error("Failed to save reset code", slog.String("email", email), slog.String("error", err.Error()))
		return fmt.Errorf("failed to save reset code")
	}

	subject := "Password Reset Code"
	body := fmt.Sprintf("Your password reset code is: %s", resetCode)

	err = mail.SendEmail(email, subject, body)
	if err != nil {
		logger.Log.Error("Failed to send reset email", slog.String("email", email), slog.String("error", err.Error()))
		return fmt.Errorf("failed to send reset email")
	}

	logger.Log.Info("Password reset email sent", slog.String("email", email))
	return nil
}

func (s *AuthService) ResetPassword(email, resetCode, newPassword, confirmPassword string) error {
	logger.Log.Info("Verifying password reset request", slog.String("email", email))

	if newPassword != confirmPassword {
		logger.Log.Warn("Password confirmation does not match", slog.String("email", email))
		return fmt.Errorf("passwords do not match")
	}

	user, err := s.repo.GetUserByEmail(email)
	if err != nil {
		logger.Log.Warn("User not found", slog.String("email", email))
		return fmt.Errorf("user not found")
	}

	savedCode, expiresAt, err := s.repo.GetPasswordResetCode(user.ID)
	if err != nil {
		logger.Log.Warn("Invalid or expired reset code", slog.String("email", email))
		return fmt.Errorf("invalid or expired reset code")
	}

	if resetCode != savedCode {
		logger.Log.Warn("Incorrect reset code", slog.String("email", email))
		return fmt.Errorf("incorrect reset code")
	}

	if time.Now().After(expiresAt) {
		logger.Log.Warn("Reset code expired", slog.String("email", email))
		return fmt.Errorf("reset code expired")
	}

	hashedPassword, err := utils.HashPassword(newPassword)
	if err != nil {
		logger.Log.Error("Failed to hash new password", slog.String("email", email), slog.String("error", err.Error()))
		return fmt.Errorf("failed to hash new password")
	}

	err = s.repo.UpdateUserPassword(user.ID, hashedPassword)
	if err != nil {
		logger.Log.Error("Failed to update password", slog.String("email", email), slog.String("error", err.Error()))
		return fmt.Errorf("failed to update password")
	}

	_ = s.repo.DeletePasswordResetCode(user.ID)

	logger.Log.Info("Password reset successfully", slog.String("email", email))
	return nil
}

func (s *AuthService) RegisterUser(userReq *entity.User) error {
	logger.Log.Info("Registering user", slog.String("email", userReq.Email))

	exists, err := s.repo.UserExistsByEmail(userReq.Email)
	if err != nil {
		return err
	}
	if exists {
		return errors.New("user with this email already exists")
	}

	hashedPassword, err := utils.HashPassword(userReq.Password)
	if err != nil {
		return err
	}

	user := &entity.User{
		Email:     userReq.Email,
		Password:  hashedPassword,
		FirstName: userReq.FirstName,
		LastName:  userReq.LastName,
		RoleId:    1,
	}

	err = s.repo.CreateUser(user)
	if err != nil {
		return err
	}

	logger.Log.Info("User registered successfully", slog.String("email", userReq.Email))
	return nil
}

func (s *AuthService) RegisterHR(userReq *entity.HRRegistration) error {
	logger.Log.Info("Registering HR", slog.String("email", userReq.Email))

	invitation, err := s.invitationRepo.GetInvitationByEmail(userReq.Email)
	if err != nil {
		return err
	}
	if invitation == nil {
		return errors.New("no invitation found for this email")
	}

	exists, err := s.repo.UserExistsByEmail(userReq.Email)
	if err != nil {
		return err
	}
	if exists {
		return errors.New("user with this email already exists")
	}

	hashedPassword, err := utils.HashPassword(userReq.Password)
	if err != nil {
		return err
	}

	user := &entity.User{
		Email:     userReq.Email,
		Password:  hashedPassword,
		FirstName: userReq.FirstName,
		LastName:  userReq.LastName,
		RoleId:    2,
	}

	err = s.repo.CreateUser(user)
	if err != nil {
		return err
	}

	hr := &entity.HR{
		UserID:    user.ID,
		DepID:     invitation.DepID,
		CompanyID: invitation.CompanyID,
	}

	err = s.hrRepo.CreateHR(hr)
	if err != nil {
		return err
	}

	err = s.invitationRepo.DeleteInvitation(userReq.Email)
	if err != nil {
		logger.Log.Warn("Failed to delete invitation", slog.String("email", userReq.Email), slog.String("error", err.Error()))
	}

	logger.Log.Info("HR registered successfully", slog.String("email", userReq.Email))
	return nil
}

func (s *AuthService) SaveRefreshToken(userID int, refreshToken string) error {
	err := s.redis.Set(context.Background(), s.getRefreshTokenKey(userID), refreshToken, 24*time.Hour*7).Err()
	if err != nil {
		return err
	}
	return nil
}

func (s *AuthService) GetRefreshToken(userID int) (string, error) {
	refreshToken, err := s.redis.Get(context.Background(), s.getRefreshTokenKey(userID)).Result()
	if errors.Is(err, redis.Nil) {
		return "", errors.New("refresh token not found")
	}
	if err != nil {
		return "", err
	}
	return refreshToken, nil
}

func (s *AuthService) DeleteRefreshToken(userID int) error {
	err := s.redis.Del(context.Background(), s.getRefreshTokenKey(userID)).Err()
	if err != nil {
		return err
	}
	return nil
}

func (s *AuthService) getRefreshTokenKey(userID int) string {
	return "refresh_token:" + strconv.Itoa(userID)
}
