package service

import (
	"jumyste-app-backend/internal/entity"
	"jumyste-app-backend/internal/repository"
	"jumyste-app-backend/pkg/logger"
	"jumyste-app-backend/utils"
	"log/slog"
)

type UserService struct {
	UserRepo    *repository.UserRepository
	CompanyRepo *repository.CompanyRepository
}

func NewUserService(userRepo *repository.UserRepository, companyRepo *repository.CompanyRepository) *UserService {
	return &UserService{UserRepo: userRepo, CompanyRepo: companyRepo}
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

func (s *UserService) GetUserByIDWithClaims(claims *utils.Claims) (*entity.UserResponse, error) {
	logger.Log.Info("Fetching user by ID", slog.Int("user_id", claims.UserID))

	user, err := s.UserRepo.GetUserByID(claims.UserID)
	if err != nil {
		logger.Log.Error("Failed to fetch user", slog.Int("user_id", claims.UserID), slog.String("error", err.Error()))
		return nil, err
	}

	const hrRoleID = 2
	if claims.RoleID == hrRoleID {
		company, err := s.CompanyRepo.GetByID(claims.CompanyID)
		if err == nil {
			user.Company = company
		} else {
			logger.Log.Warn("Failed to fetch company", slog.Int("company_id", claims.CompanyID), slog.String("error", err.Error()))
		}
	}

	logger.Log.Info("User fetched successfully", slog.Int("user_id", claims.UserID))
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
