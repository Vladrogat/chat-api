package service

import (
	"chat-api/internal/domain"
	"chat-api/internal/repository"
	"errors"

	"gorm.io/gorm"
)

type ChatService struct {
	repo *repository.Repository
}

func NewChatService(repo *repository.Repository) *ChatService {
	return &ChatService{repo: repo}
}

// CreateChat creates a new chat
func (s *ChatService) CreateChat(title string) (*domain.Chat, error) {
	title, err := ValidateChat(title)
	if err != nil {
		return nil, err
	}

	chat := &domain.Chat{
		Title: title,
	}

	if err := s.repo.CreateChat(chat); err != nil {
		return nil, err
	}

	return chat, nil
}

// CreateMessage creates a message in a chat
func (s *ChatService) CreateMessage(chatID int64, text string) (*domain.Message, error) {
	text, err := ValidateMessage(text)
	if err != nil {
		return nil, err
	}

	if _, err := s.repo.GetChatByID(chatID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrChatNotFound
		}
		return nil, err
	}

	message := &domain.Message{
		ChatID: chatID,
		Text:   text,
	}

	if err := s.repo.CreateMessage(message); err != nil {
		return nil, err
	}

	return message, nil
}

// GetChatWithMessages returns chat with last N messages
func (s *ChatService) GetChatWithMessages(chatID int64, limit int) (*domain.Chat, []domain.Message, error) {
	if limit == 0 {
		limit = DefaultMessagesLimit
	}

	if err := ValidateLimit(limit); err != nil {
		return nil, nil, err
	}

	chat, err := s.repo.GetChatByID(chatID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil, ErrChatNotFound
		}
		return nil, nil, err
	}

	messages, err := s.repo.GetMessagesByChatID(chatID, limit)
	if err != nil {
		return nil, nil, err
	}

	return chat, messages, nil
}

// DeleteChat deletes chat with cascade messages
func (s *ChatService) DeleteChat(chatID int64) error {
	err := s.repo.DeleteChat(chatID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrChatNotFound
		}
		return err
	}
	return nil
}
