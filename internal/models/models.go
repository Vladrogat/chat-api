package models

import "time"

type Chat struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Title     string    `gorm:"type:varchar(200); not null" json:"title"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	Message   []Message `gorm:"constraint:onDelete:CASCAD" json:"-"`
}

type Message struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	ChatID    uint      `gorm:"not null;index" json:"chat_id"`
	Text      string    `gorm:"type:text;not null" json:"text"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	Chat      Chat      `gorm:"foreignKey:ChatID" json:"-"`
}
