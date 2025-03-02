package service

import (
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
	"jumyste-app-backend/internal/entity"
	"jumyste-app-backend/internal/repository"
	"jumyste-app-backend/pkg/logger"
	"jumyste-app-backend/pkg/mail"
	"jumyste-app-backend/utils"
	"log/slog"
	"time"
)

type AuthService struct {
	repo  *repository.AuthRepository
	redis *redis.Client
}

func NewAuthService(repo *repository.AuthRepository, redisClient *redis.Client) *AuthService {
	return &AuthService{repo: repo, redis: redisClient}
}

var (
	ErrInvalidResetCode = errors.New("invalid or expired reset code")
	ErrUserNotFound     = errors.New("user not found")
)

//func (s *AuthService) RegisterUser(email, password, firstName, lastName string) error {
//	logger.Log.Info("Starting user registration",
//		slog.String("email", email),
//		slog.String("Name", firstName))
//
//	exists, err := s.repo.UserExistsByEmail(email)
//	if err != nil {
//		logger.Log.Error("Error checking user existence",
//			slog.String("email", email),
//			slog.String("error", err.Error()))
//		return err
//	}
//	if exists {
//		logger.Log.Warn("Registration attempt for existing user",
//			slog.String("email", email))
//		return errors.New("user with this email already exists")
//	}
//
//	hashedPassword, err := utils.HashPassword(password)
//	if err != nil {
//		logger.Log.Error("Failed to hash password", slog.String("error", err.Error()))
//		return err
//	}
//
//	user := &entity.User{
//		Email:     email,
//		FirstName: firstName,
//		LastName:  lastName,
//		Password:  hashedPassword,
//	}
//
//	err = s.repo.CreateUser(user)
//	if err != nil {
//		logger.Log.Error("Failed to create user",
//			slog.String("email", email),
//			slog.String("error", err.Error()))
//		return err
//	}
//
//	logger.Log.Info("User registered successfully",
//		slog.String("email", email),
//		slog.String("name", firstName))
//
//	return nil
//}

func (s *AuthService) LoginUser(email, password string) (string, error) {
	logger.Log.Info("Attempting user login", slog.String("email", email))

	user, err := s.repo.GetUserByEmail(email)
	if err != nil {
		logger.Log.Warn("Login failed: user not found", slog.String("email", email))
		return "", errors.New("invalid credentials")
	}
	logger.Log.Info("Comparing passwords",
		slog.String("input_password", password),
		slog.String("hashed_password", user.Password))

	if !utils.CheckPassword(password, user.Password) {
		logger.Log.Warn("Login failed: incorrect password", slog.String("email", email))
		return "", errors.New("invalid credentials")
	}

	token, err := utils.GenerateJWT(user.ID)
	if err != nil {
		logger.Log.Error("Failed to generate JWT", slog.String("error", err.Error()))
		return "", err
	}

	logger.Log.Info("User logged in successfully", slog.String("email", email))

	return token, nil
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

func (s *AuthService) RegisterUser(email, password, firstName, lastName string) error {
	logger.Log.Info("Registering user and sending verification code", slog.String("email", email))

	exists, err := s.repo.UserExistsByEmail(email)
	if err != nil {
		logger.Log.Error("Error checking user existence", slog.String("email", email), slog.String("error", err.Error()))
		return err
	}
	if exists {
		return errors.New("user with this email already exists")
	}

	//verificationCode := utils.GenerateResetCode()

	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		logger.Log.Error("Failed to hash password", slog.String("error", err.Error()))
		return err
	}
	user := &entity.User{
		Email:     email,
		Password:  hashedPassword,
		FirstName: firstName,
		LastName:  lastName,
	}

	err = s.repo.CreateUser(user)
	if err != nil {
		logger.Log.Error("Failed to create user", slog.String("email", email), slog.String("error", err.Error()))
		return err
	}
	logger.Log.Info("User registered successfully", slog.String("email", email))

	//ctx := context.Background()
	//key := fmt.Sprintf("pending_registration:%s", email)
	//userData := map[string]interface{}{
	//	"email":     email,
	//	"password":  hashedPassword,
	//	"firstName": firstName,
	//	"lastName":  lastName,
	//}
	//err = s.redis.HMSet(ctx, key, userData).Err()
	//if err != nil {
	//	logger.Log.Error("Failed to save user data in Redis", slog.String("email", email), slog.String("error", err.Error()))
	//	return errors.New("failed to save user data")
	//}
	//
	//s.redis.Expire(ctx, key, 10*time.Minute)
	//
	//go func() {
	//	subject := "Verification Code"
	//	body := fmt.Sprintf("Your verification code is: %s", verificationCode)
	//
	//	err := mail.SendEmail(email, subject, body)
	//	if err != nil {
	//		logger.Log.Error("Failed to send verification email", slog.String("email", email), slog.String("error", err.Error()))
	//	} else {
	//		logger.Log.Debug("Verification code sent successfully", slog.String("email", email))
	//	}
	//}()
	//
	//logger.Log.Info("User registration process started", slog.String("email", email))

	return nil
}

//func (s *AuthService) VerifyCodeAndRegister(email, code string) error {
//	ctx := context.Background()
//	key := fmt.Sprintf("pending_registration:%s", email)
//
//	userData, err := s.redis.HGetAll(ctx, key).Result()
//	if err != nil {
//		return errors.New("failed to get registration data")
//	}
//	if len(userData) == 0 {
//		return errors.New("registration data expired or not found")
//	}
//
//	if userData["code"] != code {
//		return errors.New("invalid verification code")
//	}
//
//	user := &entity.User{
//		Email:     userData["email"],
//		Password:  userData["password"],
//		FirstName: userData["firstName"],
//		LastName:  userData["lastName"],
//	}
//
//	err = s.repo.CreateUser(user)
//	if err != nil {
//		logger.Log.Error("Failed to create user", slog.String("email", email), slog.String("error", err.Error()))
//		return err
//	}
//
//	s.redis.Del(ctx, key)
//
//	logger.Log.Info("User registered successfully", slog.String("email", email))
//	return nil
//}
