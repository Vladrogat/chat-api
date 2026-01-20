package tests

import (
	"chat-api/internal/service"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateChat(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		wantTitle string
		wantErr   error
	}{
		{
			name:      "Valid title",
			input:     "Test Chat",
			wantTitle: "Test Chat",
			wantErr:   nil,
		},
		{
			name:      "Title with spaces",
			input:     "  Test Chat  ",
			wantTitle: "Test Chat",
			wantErr:   nil,
		},
		{
			name:      "Empty title",
			input:     "",
			wantTitle: "",
			wantErr:   service.ErrEmptyTitle,
		},
		{
			name:      "Title with only spaces",
			input:     "   ",
			wantTitle: "",
			wantErr:   service.ErrEmptyTitle,
		},
		{
			name:      "Title too long",
			input:     strings.Repeat("a", 201),
			wantTitle: "",
			wantErr:   service.ErrTitleTooLong,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotTitle, err := service.ValidateChat(tt.input)
			assert.Equal(t, tt.wantTitle, gotTitle)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}

func TestValidateMessage(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		wantText string
		wantErr  error
	}{
		{
			name:     "Valid text",
			input:    "Hello, World!",
			wantText: "Hello, World!",
			wantErr:  nil,
		},
		{
			name:     "Text with spaces",
			input:    "  Hello, World!  ",
			wantText: "Hello, World!",
			wantErr:  nil,
		},
		{
			name:     "Empty text",
			input:    "",
			wantText: "",
			wantErr:  service.ErrEmptyText,
		},
		{
			name:     "Text with only spaces",
			input:    "   ",
			wantText: "",
			wantErr:  service.ErrEmptyText,
		},
		{
			name:     "Text too long",
			input:    strings.Repeat("a", 5001),
			wantText: "",
			wantErr:  service.ErrTextTooLong,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotText, err := service.ValidateMessage(tt.input)
			assert.Equal(t, tt.wantText, gotText)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}

func TestValidateLimit(t *testing.T) {
	tests := []struct {
		name    string
		input   int
		wantErr error
	}{
		{
			name:    "Valid limit",
			input:   20,
			wantErr: nil,
		},
		{
			name:    "Min limit",
			input:   1,
			wantErr: nil,
		},
		{
			name:    "Max limit",
			input:   100,
			wantErr: nil,
		},
		{
			name:    "Limit too low",
			input:   0,
			wantErr: service.ErrInvalidLimit,
		},
		{
			name:    "Limit too high",
			input:   101,
			wantErr: service.ErrInvalidLimit,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := service.ValidateLimit(tt.input)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}
