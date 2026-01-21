package service

import (
	"errors"
	"strings"
)

const (
	MinChatTitleLength = 1
	MaxChatTitleLength = 200

	MinMessageLength = 1
	MaxMessageLength = 5000

	DefaultMessagesLimit = 20
	MinMessagesLimit     = 1
	MaxMessagesLimit     = 100
)

var (
	ErrChatNotFound = errors.New("chat not found")
	ErrEmptyTitle   = errors.New("title cannot be empty")
	ErrTitleTooLong = errors.New("title must be between 1 and 200 characters")
	ErrEmptyText    = errors.New("text cannot be empty")
	ErrTextTooLong  = errors.New("text must be between 1 and 5000 characters")
	ErrInvalidLimit = errors.New("limit must be between 1 and 100")
)

// Validates chat data
func ValidateChat(title string) (string, error) {
	title = strings.TrimSpace(title)

	if title == "" {
		return "", ErrEmptyTitle
	}

	if len(title) < MinChatTitleLength || len(title) > MaxChatTitleLength {
		return "", ErrTitleTooLong
	}
	return title, nil
}

// Validates message data
func ValidateMessage(text string) (string, error) {
	text = strings.TrimSpace(text)

	if text == "" {
		return "", ErrEmptyText
	}

	if len(text) < MinMessageLength || len(text) > MaxMessageLength {
		return "", ErrTextTooLong
	}

	return text, nil
}

// Validates limit parameter
func ValidateLimit(limit int) error {
	if limit < MinMessagesLimit || limit > MaxMessagesLimit {
		return ErrInvalidLimit
	}

	return nil
}
