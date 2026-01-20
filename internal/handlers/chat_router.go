package handlers

import (
	"net/http"
	"strings"
)

func (h *Handler) HandleChats(w http.ResponseWriter, r *http.Request) {
	path := strings.Trim(r.URL.Path, "/")
	parts := strings.Split(path, "/")

	switch {
	// POST /chats
	case len(parts) == 1 && r.Method == http.MethodPost:
		h.CreateChat(w, r)

	// POST /chats/{id}/messages
	case len(parts) == 3 && parts[2] == "messages" && r.Method == http.MethodPost:
		h.CreateMessage(w, r)

	// GET /chats/{id}
	case len(parts) == 2 && r.Method == http.MethodGet:
		h.GetChat(w, r)

	// DELETE /chats/{id}
	case len(parts) == 2 && r.Method == http.MethodDelete:
		h.DeleteChat(w, r)

	default:
		http.NotFound(w, r)
	}
}
