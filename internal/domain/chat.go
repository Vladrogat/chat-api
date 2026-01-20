package domain

import "time"

type Chat struct {
	ID        int64     `gorm:"primaryKey;column:id"`
	Title     string    `gorm:"column:title;size:200;not null"`
	CreatedAt time.Time `gorm:"column:created_at"`

	Messages []Message `gorm:"foreignKey:ChatID"`
}
