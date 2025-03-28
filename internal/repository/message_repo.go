package repository

import (
	"database/sql"
	"jumyste-app-backend/internal/entity"
	"time"
)

type MessageRepository struct {
	DB *sql.DB
}

func NewMessageRepository(db *sql.DB) *MessageRepository {
	return &MessageRepository{DB: db}
}

func (r *MessageRepository) CreateMessage(message *entity.Message) (int, error) {
	query := `INSERT INTO messages (chat_id, sender_id, type, content, file_url, created_at)
	          VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`

	var messageID int

	var content sql.NullString
	if message.Content != nil {
		content = sql.NullString{String: *message.Content, Valid: true}
	}

	var fileURL sql.NullString
	if message.FileURL != nil {
		fileURL = sql.NullString{String: *message.FileURL, Valid: true}
	}

	err := r.DB.QueryRow(query, message.ChatID, message.SenderID, message.Type, content, fileURL, time.Now()).Scan(&messageID)
	if err != nil {
		return 0, err
	}

	message.ID = messageID
	return messageID, nil
}

// GetMessagesByChatID - Fetches all messages in a chat
func (r *MessageRepository) GetMessagesByChatID(chatID int) ([]entity.Message, error) {
	query := `SELECT id, chat_id, sender_id, type, content, file_url, created_at 
	          FROM messages WHERE chat_id = $1 ORDER BY created_at`
	rows, err := r.DB.Query(query, chatID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []entity.Message
	for rows.Next() {
		var message entity.Message
		err := rows.Scan(&message.ID, &message.ChatID, &message.SenderID, &message.Type, &message.Content, &message.FileURL, &message.CreatedAt)
		if err != nil {
			return nil, err
		}
		messages = append(messages, message)
	}

	return messages, nil
}

// GetMessageByID - Retrieves a single message by ID
func (r *MessageRepository) GetMessageByID(messageID int) (*entity.Message, error) {
	query := `SELECT id, chat_id, sender_id, type, content, file_url, created_at 
	          FROM messages WHERE id = $1`
	row := r.DB.QueryRow(query, messageID)

	var message entity.Message
	err := row.Scan(&message.ID, &message.ChatID, &message.SenderID, &message.Type, &message.Content, &message.FileURL, &message.CreatedAt)
	if err != nil {
		return nil, err
	}

	return &message, nil
}
