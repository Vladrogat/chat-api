package database

import (
	"chat-api/internal/domain"
	"errors"
	"slices"

	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

// Return new repository
//
//	@param db *gorm.DB
//	@return *Repository
func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

// Create chat
//
//	@param chat *domain.Chat
//	@return error
func (r *Repository) CreateChat(chat *domain.Chat) error {
	return r.db.Create(chat).Error
}

// Get chat by ID
//
//	@param id uint
//	@return *domain.Chat
//	@return error
func (r *Repository) GetChatByID(id uint) (*domain.Chat, error) {
	var chat domain.Chat
	err := r.db.First(&chat, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &chat, nil
}

// Delete chat
//
//	@param id uint
//	@return error
func (r *Repository) DeleteChat(id uint) error {
	result := r.db.Delete(&domain.Chat{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

// Create message
//
//	@param message *domain.Message
//	@return error
func (r *Repository) CreateMessage(message *domain.Message) error {
	return r.db.Create(message).Error
}

func (r *Repository) GetMessagesByChatID(chatID uint, limit int) ([]domain.Message, error) {
	var messages []domain.Message

	err := r.db.Where("chat_id = ?", chatID).
		Order("created_at DESC").
		Limit(limit).
		Find(&messages).Error
	if err != nil {
		return nil, err
	}

	slices.Reverse(messages)

	return messages, nil
}
