package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"jumyste-app-backend/internal/entity"
	"jumyste-app-backend/pkg/logger"
	"log/slog"
	"strings"
)

type UserRepository struct {
	DB *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{DB: db}
}

func (r *UserRepository) CreateUser(user *entity.User) error {
	query := `INSERT INTO users(email, password, first_name, last_name, profile_picture) 
	          VALUES ($1, $2, $3, $4, $5) RETURNING id`
	return r.DB.QueryRow(query, user.Email, user.Password, user.FirstName, user.LastName, user.ProfilePicture).Scan(&user.ID)
}

func (r *UserRepository) GetUserByID(id int) (*entity.UserResponse, error) {
	query := `SELECT id, email, first_name, last_name, profile_picture FROM users WHERE id = $1`
	row := r.DB.QueryRow(query, id)

	var user entity.UserResponse
	err := row.Scan(&user.ID, &user.Email, &user.FirstName, &user.LastName, &user.ProfilePicture)
	return &user, err
}
func (r *UserRepository) UpdateUser(userID int, updates map[string]interface{}) error {
	if len(updates) == 0 {
		return errors.New("no fields to update")
	}

	var setClauses []string
	params := []interface{}{}
	counter := 1

	for key, value := range updates {
		setClauses = append(setClauses, fmt.Sprintf("%s = $%d", key, counter))
		params = append(params, value)
		counter++
	}

	query := fmt.Sprintf("UPDATE users SET %s WHERE id = $%d", strings.Join(setClauses, ", "), counter)
	params = append(params, userID)

	logger.Log.Info("Executing SQL query", slog.String("query", query), slog.Any("params", params))

	_, err := r.DB.Exec(query, params...)
	if err != nil {
		logger.Log.Error("Failed to update user in repository",
			slog.Int("user_id", userID),
			slog.String("error", err.Error()))
		return err
	}

	logger.Log.Info("User updated successfully in repository",
		slog.Int("user_id", userID))
	return nil
}
