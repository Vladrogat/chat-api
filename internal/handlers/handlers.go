package handlers

import (
	"chat-api/internal/database"
	"chat-api/internal/domain"
	"chat-api/internal/service"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"strings"

	"gorm.io/gorm"
)

type Handler struct {
	repo *database.Repository
}

func NewHandler(repo *database.Repository) *Handler {
	return &Handler{repo: repo}
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type CreateChatRequest struct {
	Title string `json:"title"`
}

type CreateMessageRequest struct {
	Text string `json:"text"`
}

type GetChatResponse struct {
	ID        int64            `json:"id"`
	Title     string           `json:"title"`
	CreatedAt string           `json:"created_at"`
	Messages  []domain.Message `json:"messages"`
}

// sends a JSON response
//
//	@param w http.ResponseWriter
//	@param status int
//	@param data interface{}
func respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if data != nil {
		json.NewEncoder(w).Encode(data)
	}
}

// sends an error response
//
//	@param w http.ResponseWriter
//	@param status int
//	@param message string
func respondError(w http.ResponseWriter, status int, message string) {
	log.Printf("Error response: %d - %s", status, message)
	respondJSON(w, status, ErrorResponse{Error: message})
}

// handles POST /chats/{id}/messages/
//
//	@param w http.ResponseWriter
//	@param r *http.Request
func (h *Handler) CreateChat(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var req CreateChatRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	title, err := service.ValidateChat(req.Title)

	if err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	chat := &domain.Chat{
		Title: title,
	}

	if err := h.repo.CreateChat(chat); err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to create chat")
		return
	}

	respondJSON(w, http.StatusCreated, chat)
}

// handles POST /chats/{id}/messages/
//
//	@param w http.ResponseWriter
//	@param r *http.Request
func (h *Handler) CreateMessage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	pathParts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(pathParts) < 2 {
		respondError(w, http.StatusBadRequest, "Invalid path")
		return
	}

	chatID, err := strconv.ParseUint(pathParts[1], 10, 32)

	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid chat ID")
		return
	}

	chat, err := h.repo.GetChatByID(uint(chatID))

	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to check chat")
		return
	}

	if chat == nil {
		respondError(w, http.StatusNotFound, "Chat not found")
		return
	}

	var req CreateMessageRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	text, err := service.ValidateMessage(req.Text)
	if err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	message := &domain.Message{
		ChatID: int64(chatID),
		Text:   text,
	}

	if err := h.repo.CreateMessage(message); err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to create message")
		return
	}

	respondJSON(w, http.StatusCreated, message)
}

// handles GET /chats/{id}
//
//	@param w http.ResponseWriter
//	@param r *http.Request
func (h *Handler) GetChat(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	pathParts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(pathParts) < 2 {
		respondError(w, http.StatusBadRequest, "Invalid path")
		return
	}

	chatID, err := strconv.ParseUint(pathParts[1], 10, 32)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid chat ID")
		return
	}

	limit := 20
	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		parsedLimit, err := strconv.Atoi(limitStr)

		if err != nil || service.ValidateLimit(parsedLimit) != nil {
			respondError(w, http.StatusBadRequest, "Limit must be between 1 and 100")
			return
		}
		limit = parsedLimit
	}

	chat, err := h.repo.GetChatByID(uint(chatID))
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to get chat")
		return
	}
	if chat == nil {
		respondError(w, http.StatusNotFound, "Chat not found")
		return
	}

	messages, err := h.repo.GetMessagesByChatID(uint(chatID), limit)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to get messages")
		return
	}

	response := GetChatResponse{
		ID:        chat.ID,
		Title:     chat.Title,
		CreatedAt: chat.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		Messages:  messages,
	}

	respondJSON(w, http.StatusOK, response)
}

// handles DELETE /chats/{id}
//
//	@param w http.ResponseWriter
//	@param r *http.Request
func (h *Handler) DeleteChat(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	pathParts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(pathParts) < 2 {
		respondError(w, http.StatusBadRequest, "Invalid path")
		return
	}

	chatID, err := strconv.ParseUint(pathParts[1], 10, 32)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid chat ID")
		return
	}

	if err = h.repo.DeleteChat(uint(chatID)); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			respondError(w, http.StatusNotFound, "Chat not found")
			return
		}
		respondError(w, http.StatusInternalServerError, "Failed to delete chat")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
