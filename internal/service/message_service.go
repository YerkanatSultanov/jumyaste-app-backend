package service

import (
	"jumyste-app-backend/internal/entity"
	"jumyste-app-backend/internal/repository"
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
	}

	messageID, err := s.MessageRepo.CreateMessage(message)
	if err != nil {
		return nil, err
	}

	message.ID = messageID
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
