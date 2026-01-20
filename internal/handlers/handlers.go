package handlers

import (
	"chat-api/internal/service"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type Handler struct {
	service *service.ChatService
}

func NewHandler(service *service.ChatService) *Handler {
	return &Handler{service: service}
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
	ID        int64       `json:"id"`
	Title     string      `json:"title"`
	CreatedAt string      `json:"created_at"`
	Messages  interface{} `json:"messages"`
}

func respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if data != nil {
		_ = json.NewEncoder(w).Encode(data)
	}
}

func respondError(w http.ResponseWriter, status int, msg string) {
	respondJSON(w, status, ErrorResponse{Error: msg})
}

func parseChatID(r *http.Request) (int64, error) {
	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(parts) < 2 {
		return 0, errors.New("invalid path")
	}
	return strconv.ParseInt(parts[1], 10, 64)
}

// POST /chats
func (h *Handler) CreateChat(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		respondError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	var req CreateChatRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	chat, err := h.service.CreateChat(req.Title)
	if err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	respondJSON(w, http.StatusCreated, chat)
}

// POST /chats/{id}/messages
func (h *Handler) CreateMessage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		respondError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	chatID, err := parseChatID(r)
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid chat id")
		return
	}

	var req CreateMessageRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	message, err := h.service.CreateMessage(chatID, req.Text)
	if err != nil {
		if errors.Is(err, service.ErrChatNotFound) {
			respondError(w, http.StatusNotFound, "chat not found")
			return
		}
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	respondJSON(w, http.StatusCreated, message)
}

// GET /chats/{id}?limit=N
func (h *Handler) GetChat(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		respondError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	chatID, err := parseChatID(r)
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid chat id")
		return
	}

	limit := 20
	if v := r.URL.Query().Get("limit"); v != "" {
		limit, err = strconv.Atoi(v)
		if err != nil {
			respondError(w, http.StatusBadRequest, "invalid limit")
			return
		}
	}

	chat, messages, err := h.service.GetChatWithMessages(chatID, limit)
	if err != nil {
		if errors.Is(err, service.ErrChatNotFound) {
			respondError(w, http.StatusNotFound, "chat not found")
			return
		}
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	resp := GetChatResponse{
		ID:        chat.ID,
		Title:     chat.Title,
		CreatedAt: chat.CreatedAt.Format(time.RFC3339),
		Messages:  messages,
	}

	respondJSON(w, http.StatusOK, resp)
}

// DELETE /chats/{id}
func (h *Handler) DeleteChat(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		respondError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	chatID, err := parseChatID(r)
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid chat id")
		return
	}

	if err := h.service.DeleteChat(chatID); err != nil {
		if errors.Is(err, service.ErrChatNotFound) {
			respondError(w, http.StatusNotFound, "chat not found")
			return
		}
		respondError(w, http.StatusInternalServerError, "failed to delete chat")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
