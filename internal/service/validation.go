package service

import (
	"errors"
	"strings"
)

var (
	ErrEmptyTitle   = errors.New("title cannot be empty")
	ErrTitleTooLong = errors.New("title must be between 1 and 200 characters")
	ErrEmptyText    = errors.New("text cannot be empty")
	ErrTextTooLong  = errors.New("text must be between 1 and 5000 characters")
	ErrInvalidLimit = errors.New("limit must be between 1 and 100")
)

// Validates chat data
//
//	@param title string
//	@return string
//	@return error
func ValidateChat(title string) (string, error) {
	title = strings.TrimSpace(title)

	if title == "" {
		return "", ErrEmptyTitle
	}

	if len(title) < 1 || len(title) > 200 {
		return "", ErrTitleTooLong
	}
	return title, nil
}

// Validates message data
//
//	@param text string
//	@return string
//	@return error
func ValidateMessage(text string) (string, error) {
	text = strings.TrimSpace(text)

	if text == "" {
		return "", ErrEmptyText
	}

	if len(text) < 1 || len(text) > 5000 {
		return "", ErrTextTooLong
	}

	return text, nil
}

// Validates limit parameter
//
//	@param limit int
//	@return error
func ValidateLimit(limit int) error {
	if limit < 1 || limit > 100 {
		return ErrInvalidLimit
	}

	return nil
}
