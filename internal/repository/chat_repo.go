package repository

import (
	"database/sql"
	"github.com/lib/pq"
	"jumyste-app-backend/internal/entity"
	"jumyste-app-backend/pkg/logger"
	"log/slog"
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
//func (r *ChatRepository) GetAllChats() ([]entity.Chat, error) {
//	query := "SELECT id, created_at, updated_at FROM chats"
//	rows, err := r.DB.Query(query)
//	if err != nil {
//		return nil, err
//	}
//	defer rows.Close()
//
//	var chats []entity.Chat
//	for rows.Next() {
//		var chat entity.Chat
//		if err := rows.Scan(&chat.ID, &chat.CreatedAt, &chat.UpdatedAt); err != nil {
//			return nil, err
//		}
//		chats = append(chats, chat)
//	}
//	return chats, nil
//}

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

		users, err := r.GetUsersByChatID(chat.ID, -1)
		if err != nil {
			return nil, err
		}
		chat.Users = users

		chats = append(chats, chat)
	}
	return chats, nil
}

func (r *ChatRepository) GetUsersByIDs(userIDs []int) ([]entity.UserResponse, error) {
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

// GetChatsByUserID - Получает все чаты, в которых состоит пользователь
func (r *ChatRepository) GetChatsByUserID(userID int) ([]entity.Chat, error) {
	query := `
		SELECT 
		    c.id, 
		    c.created_at, 
		    c.updated_at,
		    COALESCE(m.content, '') AS last_message, 
		    m.created_at AS last_message_at,
		    COALESCE(($1 = ANY(m.read_by)), false) AS is_read 
		FROM chats c
		JOIN chat_users cu ON c.id = cu.chat_id
		LEFT JOIN (
		    SELECT DISTINCT ON (chat_id) chat_id, content, created_at, read_by
		    FROM messages
		    ORDER BY chat_id, created_at DESC
		) m ON c.id = m.chat_id
		WHERE cu.user_id = $1
		ORDER BY m.created_at DESC NULLS LAST;
	`

	rows, err := r.DB.Query(query, userID)
	if err != nil {
		logger.Log.Error("Failed to fetch chats", slog.Int("user_id", userID), slog.String("error", err.Error()))
		return nil, err
	}
	defer rows.Close()

	var chats []entity.Chat
	for rows.Next() {
		var chat entity.Chat
		var lastMessage string
		var lastMessageAt sql.NullTime
		var isRead sql.NullBool // ✅ Используем sql.NullBool

		if err := rows.Scan(&chat.ID, &chat.CreatedAt, &chat.UpdatedAt, &lastMessage, &lastMessageAt, &isRead); err != nil {
			return nil, err
		}

		chat.LastMessage = lastMessage
		if lastMessageAt.Valid {
			chat.LastMessageAt = lastMessageAt.Time
		}
		chat.IsRead = isRead.Valid && isRead.Bool

		users, err := r.GetUsersByChatID(chat.ID, userID)
		if err != nil {
			return nil, err
		}
		chat.Users = users

		chats = append(chats, chat)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return chats, nil
}

// GetUsersByChatID - Получает собеседников в чате
func (r *ChatRepository) GetUsersByChatID(chatID, userID int) ([]entity.UserResponse, error) {
	query := `
		SELECT u.id, u.email, u.first_name, u.last_name, u.profile_picture 
		FROM users u 
		JOIN chat_users cu ON u.id = cu.user_id 
		WHERE cu.chat_id = $1 AND u.id != $2`

	rows, err := r.DB.Query(query, chatID, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []entity.UserResponse
	for rows.Next() {
		var user entity.UserResponse
		if err := rows.Scan(&user.ID, &user.Email, &user.FirstName, &user.LastName, &user.ProfilePicture); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}
