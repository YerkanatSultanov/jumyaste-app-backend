package repository

import (
	"database/sql"
	"github.com/lib/pq"
	"jumyste-app-backend/internal/entity"
	"time"
)

type ChatRepository struct {
	DB *sql.DB
}

func NewChatRepository(db *sql.DB) *ChatRepository {
	return &ChatRepository{DB: db}
}

// CreateChat - Creates a new chat and adds users to chat_users table
func (r *ChatRepository) CreateChat(chat *entity.Chat) (int, error) {
	tx, err := r.DB.Begin()
	if err != nil {
		return 0, err
	}

	query := "INSERT INTO chats (created_at, updated_at) VALUES ($1, $2) RETURNING id"
	var chatID int
	err = tx.QueryRow(query, time.Now(), time.Now()).Scan(&chatID)
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	chat.ID = chatID

	userQuery := "INSERT INTO chat_users (chat_id, user_id) VALUES ($1, $2)"
	for _, user := range chat.Users {
		_, err := tx.Exec(userQuery, chatID, user.ID)
		if err != nil {
			tx.Rollback()
			return 0, err
		}
	}

	err = tx.Commit()
	if err != nil {
		return 0, err
	}

	return chatID, nil
}

// GetChatByID - Fetches a chat by its ID
func (r *ChatRepository) GetChatByID(chatID uint) (*entity.Chat, error) {
	query := "SELECT id, created_at, updated_at FROM chats WHERE id = $1"
	row := r.DB.QueryRow(query, chatID)

	var chat entity.Chat
	err := row.Scan(&chat.ID, &chat.CreatedAt, &chat.UpdatedAt)
	if err != nil {
		return nil, err
	}

	userQuery := `SELECT u.id, u.first_name, u.last_name,u.email FROM users u 
	              JOIN chat_users cu ON u.id = cu.user_id WHERE cu.chat_id = $1`
	rows, err := r.DB.Query(userQuery, chatID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []entity.UserResponse
	for rows.Next() {
		var user entity.UserResponse
		if err := rows.Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email); err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	chat.Users = users

	return &chat, nil
}

// GetAllChats - Retrieves all chats
func (r *ChatRepository) GetAllChats() ([]entity.Chat, error) {
	query := "SELECT id, created_at, updated_at FROM chats"
	rows, err := r.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var chats []entity.Chat
	for rows.Next() {
		var chat entity.Chat
		if err := rows.Scan(&chat.ID, &chat.CreatedAt, &chat.UpdatedAt); err != nil {
			return nil, err
		}
		chats = append(chats, chat)
	}
	return chats, nil
}

func (r *ChatRepository) GetUsersByIDs(userIDs []uint) ([]entity.UserResponse, error) {
	query := "SELECT id, first_name, last_name, email FROM users WHERE id = ANY($1)"
	rows, err := r.DB.Query(query, pq.Array(userIDs))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []entity.UserResponse
	for rows.Next() {
		var user entity.UserResponse
		if err := rows.Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}
