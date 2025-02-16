package repository

import (
	"database/sql"
	"jumyste-app-backend/internal/entity"
	"jumyste-app-backend/pkg/logger"
	"log/slog"
	"time"
)

type AuthRepository struct {
	db *sql.DB
}

func NewAuthRepository(db *sql.DB) *AuthRepository {
	return &AuthRepository{db: db}
}

func (r *AuthRepository) UserExistsByEmail(email string) (bool, error) {
	logger.Log.Info("Checking if user exists", slog.String("email", email))

	var exists bool
	query := "SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)"
	err := r.db.QueryRow(query, email).Scan(&exists)
	if err != nil {
		logger.Log.Error("Failed to check user existence",
			slog.String("email", email),
			slog.String("error", err.Error()))
		return false, err
	}

	if exists {
		logger.Log.Warn("User already exists", slog.String("email", email))
	} else {
		logger.Log.Info("User does not exist", slog.String("email", email))
	}

	return exists, nil
}

func (r *AuthRepository) CreateUser(user *entity.User) error {
	logger.Log.Info("Creating new user", slog.String("email", user.Email))

	query := "INSERT INTO users(email,password, first_name, last_name, profile_picture) VALUES ($1, $2, $3, $4, $5) RETURNING id"
	err := r.db.QueryRow(query, user.Email, user.Password, user.FirstName, user.LastName, user.ProfilePicture).Scan(&user.ID)
	if err != nil {
		logger.Log.Error("Failed to insert user",
			slog.String("email", user.Email),
			slog.String("error", err.Error()))
		return err
	}

	logger.Log.Info("User created successfully", slog.String("email", user.Email))
	return nil
}

func (r *AuthRepository) GetUserByEmail(email string) (*entity.User, error) {
	logger.Log.Info("Fetching user by email", slog.String("email", email))

	var user entity.User
	query := "SELECT id, email, password, first_name, last_name, profile_picture  FROM users WHERE email = $1"
	err := r.db.QueryRow(query, email).Scan(&user.ID, &user.Email, &user.Password, &user.FirstName, &user.LastName, &user.ProfilePicture)
	if err != nil {
		logger.Log.Error("User not found",
			slog.String("email", email),
			slog.String("error", err.Error()))
		return nil, err
	}

	logger.Log.Info("User retrieved successfully", slog.String("email", email))
	return &user, nil
}

func (r *AuthRepository) SavePasswordResetToken(userID int, token string, expiration time.Time) error {
	_, err := r.db.Exec(
		"INSERT INTO password_resets(user_id, token, expires_at) VALUES ($1, $2, $3)",
		userID, token, expiration,
	)
	return err
}

func (r *AuthRepository) GetUserByResetToken(token string) (*entity.User, error) {
	var user entity.User
	query := `
		SELECT u.id, u.email, u.password 
		FROM users u
		JOIN password_resets pr ON u.id = pr.user_id
		WHERE pr.token = $1 AND pr.expires_at > NOW()`
	err := r.db.QueryRow(query, token).Scan(&user.ID, &user.Email, &user.Password)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *AuthRepository) UpdateUserPassword(userID int, hashedPassword string) error {
	_, err := r.db.Exec("UPDATE users SET password = $1 WHERE id = $2", hashedPassword, userID)
	return err
}

func (r *AuthRepository) DeletePasswordResetToken(token string) error {
	_, err := r.db.Exec("DELETE FROM password_resets WHERE token = $1", token)
	return err
}
