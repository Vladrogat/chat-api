package repository

import (
	"chat-api/internal/domain"

	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

// Return new repository
func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

// Create chat
func (r *Repository) CreateChat(chat *domain.Chat) error {
	return r.db.Create(chat).Error
}

// Get chat by ID
func (r *Repository) GetChatByID(id int64) (*domain.Chat, error) {
	var chat domain.Chat
	err := r.db.First(&chat, id).Error
	if err != nil {
		return nil, err
	}
	return &chat, nil
}

// Delete chat
func (r *Repository) DeleteChat(id int64) error {
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
func (r *Repository) CreateMessage(message *domain.Message) error {
	return r.db.Create(message).Error
}

func (r *Repository) GetMessagesByChatID(chatID int64, limit int) ([]domain.Message, error) {
	var messages []domain.Message

	err := r.db.Where("chat_id = ?", chatID).
		Order("created_at DESC").
		Limit(limit).
		Find(&messages).Error
	if err != nil {
		return nil, err
	}

	for i, j := 0, len(messages)-1; i < j; i, j = i+1, j-1 {
		messages[i], messages[j] = messages[j], messages[i]
	}

	return messages, nil
}
