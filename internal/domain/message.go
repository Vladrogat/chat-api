package domain

import "time"

type Message struct {
	ID        int64     `gorm:"primaryKey;column:id"`
	ChatID    int64     `gorm:"column:chat_id;not null;index"`
	Text      string    `gorm:"column:text;size:5000;not null"`
	CreatedAt time.Time `gorm:"column:created_at"`
}
