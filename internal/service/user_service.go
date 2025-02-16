package service

import (
	"jumyste-app-backend/internal/entity"
	"jumyste-app-backend/internal/repository"
	"jumyste-app-backend/pkg/logger"
	"log/slog"
)

type UserService struct {
	UserRepo *repository.UserRepository
}

func NewUserService(userRepo *repository.UserRepository) *UserService {
	return &UserService{UserRepo: userRepo}
}

func (s *UserService) CreateUser(user *entity.User) error {
	logger.Log.Info("Creating new user", slog.String("email", user.Email))
	err := s.UserRepo.CreateUser(user)
	if err != nil {
		logger.Log.Error("Failed to create user", slog.String("email", user.Email), slog.String("error", err.Error()))
		return err
	}

	logger.Log.Info("User created successfully", slog.String("email", user.Email))
	return nil
}

func (s *UserService) GetUserByID(id int) (*entity.UserResponse, error) {
	logger.Log.Info("Fetching user by ID", slog.Int("user_id", id))

	user, err := s.UserRepo.GetUserByID(id)
	if err != nil {
		logger.Log.Error("Failed to fetch user", slog.Int("user_id", id), slog.String("error", err.Error()))
		return nil, err
	}

	logger.Log.Info("User fetched successfully", slog.Int("user_id", id))
	return user, nil
}

func (s *UserService) UpdateUser(userID int, updates map[string]interface{}) error {
	logger.Log.Info("Service: Updating user", slog.Int("user_id", userID), slog.Any("updates", updates))

	if err := s.UserRepo.UpdateUser(userID, updates); err != nil {
		logger.Log.Error("Failed to update user in service", slog.String("error", err.Error()))
		return err
	}

	logger.Log.Info("User updated successfully in service", slog.Int("user_id", userID))
	return nil
}
