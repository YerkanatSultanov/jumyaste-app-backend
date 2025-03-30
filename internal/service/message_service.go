package service

import (
	"context"
	"github.com/lib/pq"
	"jumyste-app-backend/internal/entity"
	"jumyste-app-backend/internal/repository"
	"jumyste-app-backend/pkg/logger"
	"log/slog"
)

type MessageService struct {
	MessageRepo *repository.MessageRepository
}

func NewMessageService(messageRepo *repository.MessageRepository) *MessageService {
	return &MessageService{MessageRepo: messageRepo}
}

// SendMessage - Creates a new message
func (s *MessageService) SendMessage(chatID, senderID int, msgType entity.MessageType, content *string, fileURL *string) (*entity.Message, error) {
	message := &entity.Message{
		ChatID:   chatID,
		SenderID: senderID,
		Type:     msgType,
		Content:  content,
		FileURL:  fileURL,
		ReadBy:   pq.Int64Array{},
	}

	messageID, err := s.MessageRepo.CreateMessage(message)
	if err != nil {
		logger.Log.Error("Failed to create message", slog.Int("chat_id", chatID), slog.Int("sender_id", senderID), slog.String("error", err.Error()))
		return nil, err
	}

	message.ID = messageID
	logger.Log.Info("Message sent", slog.Int("message_id", messageID), slog.Int("chat_id", chatID), slog.Int("sender_id", senderID))
	return message, nil
}

// GetMessagesByChatID - Fetch all messages from a chat
func (s *MessageService) GetMessagesByChatID(chatID int) ([]entity.Message, error) {
	return s.MessageRepo.GetMessagesByChatID(chatID)
}

// GetMessageByID - Fetch a message by ID
func (s *MessageService) GetMessageByID(messageID int) (*entity.Message, error) {
	return s.MessageRepo.GetMessageByID(messageID)
}

func (s *MessageService) MarkMessageAsRead(ctx context.Context, messageID, userID int) error {
	err := s.MessageRepo.MarkMessageAsRead(messageID, userID)
	if err != nil {
		logger.Log.Error("Failed to update read status", slog.Int("message_id", messageID), slog.Int("user_id", userID), slog.String("error", err.Error()))
		return err
	}

	logger.Log.Info("Message marked as read", slog.Int("message_id", messageID), slog.Int("user_id", userID))
	return nil
}
