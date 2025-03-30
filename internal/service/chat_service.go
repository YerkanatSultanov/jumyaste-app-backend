package service

import (
	"errors"
	"jumyste-app-backend/internal/entity"
	"jumyste-app-backend/internal/repository"
	"jumyste-app-backend/pkg/logger"
	"log/slog"
)

type ChatService struct {
	ChatRepo *repository.ChatRepository
}

func NewChatService(chatRepo *repository.ChatRepository) *ChatService {
	return &ChatService{ChatRepo: chatRepo}
}

// CreateChat - Creates a new chat
func (s *ChatService) CreateChat(userIDs []uint) (*entity.Chat, error) {
	if len(userIDs) == 0 {
		return nil, errors.New("at least one user must be provided")
	}

	users, err := s.ChatRepo.GetUsersByIDs(userIDs)
	if err != nil {
		return nil, err
	}
	if len(users) != len(userIDs) {
		return nil, errors.New("one or more users do not exist")
	}

	chat := &entity.Chat{
		Users: users,
	}

	chatID, err := s.ChatRepo.CreateChat(chat)
	if err != nil {
		return nil, err
	}

	chat.ID = chatID
	return chat, nil
}

// GetChatByID - Fetch a chat by ID
func (s *ChatService) GetChatByID(chatID uint) (*entity.Chat, error) {
	chat, err := s.ChatRepo.GetChatByID(chatID)
	if err != nil {
		return nil, err
	}
	return chat, nil
}

// GetAllChats - Fetch all chats
func (s *ChatService) GetAllChats() ([]entity.Chat, error) {
	chats, err := s.ChatRepo.GetAllChats()
	if err != nil {
		return nil, err
	}
	return chats, nil
}

// GetChatsByUserID - Получает чаты, в которых участвует пользователь
func (s *ChatService) GetChatsByUserID(userID int) ([]entity.Chat, error) {
	logger.Log.Info("Fetching chats for user", slog.Int("user_id", userID))

	chats, err := s.ChatRepo.GetChatsByUserID(userID)
	if err != nil {
		logger.Log.Error("Failed to fetch user chats", slog.Int("user_id", userID), slog.String("error", err.Error()))
		return nil, err
	}

	logger.Log.Info("Successfully fetched user chats", slog.Int("user_id", userID), slog.Int("chat_count", len(chats)))
	return chats, nil
}
