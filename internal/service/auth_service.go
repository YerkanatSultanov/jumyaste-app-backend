package service

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"jumyste-app-backend/config"
	"jumyste-app-backend/internal/entity"
	"jumyste-app-backend/internal/repository"
	"jumyste-app-backend/pkg/logger"
	"jumyste-app-backend/pkg/mail"
	"jumyste-app-backend/utils"
	"log/slog"
	"time"
)

type AuthService struct {
	repo      *repository.AuthRepository
	JWTSecret string
}

func NewAuthService(repo *repository.AuthRepository) *AuthService {
	return &AuthService{
		repo:      repo,
		JWTSecret: config.AppConfig.JWT.Secret,
	}
}

func (s *AuthService) RegisterUser(email, password, firstName, lastName string) error {
	logger.Log.Info("Starting user registration",
		slog.String("email", email),
		slog.String("Name", firstName))

	exists, err := s.repo.UserExistsByEmail(email)
	if err != nil {
		logger.Log.Error("Error checking user existence",
			slog.String("email", email),
			slog.String("error", err.Error()))
		return err
	}
	if exists {
		logger.Log.Warn("Registration attempt for existing user",
			slog.String("email", email))
		return errors.New("user with this email already exists")
	}

	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		logger.Log.Error("Failed to hash password", slog.String("error", err.Error()))
		return err
	}

	user := &entity.User{
		Email:     email,
		FirstName: firstName,
		LastName:  lastName,
		Password:  hashedPassword,
	}

	err = s.repo.CreateUser(user)
	if err != nil {
		logger.Log.Error("Failed to create user",
			slog.String("email", email),
			slog.String("error", err.Error()))
		return err
	}

	logger.Log.Info("User registered successfully",
		slog.String("email", email),
		slog.String("name", firstName))

	return nil
}

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

func (s *AuthService) VerifyToken(tokenString string) (int, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(s.JWTSecret), nil
	})

	if err != nil || !token.Valid {
		return 0, errors.New("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0, errors.New("invalid token claims")
	}

	exp, ok := claims["exp"].(float64)
	if !ok || time.Now().Unix() > int64(exp) {
		return 0, errors.New("token expired")
	}

	userID, ok := claims["user_id"].(float64)
	if !ok {
		return 0, errors.New("missing user_id")
	}

	return int(userID), nil
}

func generateToken() string {
	b := make([]byte, 32)
	rand.Read(b)
	return base64.URLEncoding.EncodeToString(b)
}

// RequestPasswordReset генерирует токен и отправляет письмо
func (s *AuthService) RequestPasswordReset(email string) error {
	user, err := s.repo.GetUserByEmail(email)
	if err != nil {
		return errors.New("пользователь не найден")
	}

	token := generateToken()
	expiration := time.Now().Add(1 * time.Hour)

	if err := s.repo.SavePasswordResetToken(user.ID, token, expiration); err != nil {
		return err
	}

	resetLink := fmt.Sprintf("http://localhost:3000/reset-password?token=%s", token)

	emailBody := fmt.Sprintf("Для сброса пароля перейдите по ссылке: %s", resetLink)
	if err := mail.SendEmail(user.Email, "Сброс пароля", emailBody); err != nil {
		return err
	}

	return nil
}

func (s *AuthService) ResetPassword(token, newPassword string) error {
	user, err := s.repo.GetUserByResetToken(token)
	if err != nil {
		return errors.New("неверный или истекший токен")
	}

	hashedPassword, err := utils.HashPassword(newPassword)
	if err != nil {
		return errors.New("ошибка хеширования пароля")
	}

	if err := s.repo.UpdateUserPassword(user.ID, hashedPassword); err != nil {
		return err
	}

	if err := s.repo.DeletePasswordResetToken(token); err != nil {
		return err
	}

	return nil
}
