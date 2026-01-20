package main

import (
	"chat-api/internal/config"
	"chat-api/internal/database"
	"chat-api/internal/handlers"
	"log"
	"net/http"
	"strings"
)

func main() {
	cfg := config.Load()
	db, err := database.Connect(&cfg.Database)

	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	repo := database.NewRepository(db)
	handler := handlers.NewHandler(repo)

	mux := http.NewServeMux()

	mux.HandleFunc("/chats/", func(w http.ResponseWriter, r *http.Request) {
		parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")

		if len(parts) == 1 && r.Method == http.MethodPost {
			// POST /chats/
			handlers.LoggingMiddleware(handler.CreateChat)(w, r)
		} else if len(parts) == 3 && parts[2] == "messages" && r.Method == http.MethodPost {
			// POST /chats/{id}/messages/
			handlers.LoggingMiddleware(handler.CreateMessage)(w, r)
		} else if len(parts) == 2 && r.Method == http.MethodGet {
			// GET /chats/{id}
			handlers.LoggingMiddleware(handler.GetChat)(w, r)
		} else if len(parts) == 2 && r.Method == http.MethodDelete {
			// DELETE /chats/{id}
			handlers.LoggingMiddleware(handler.DeleteChat)(w, r)
		} else {
			http.NotFound(w, r)
		}
	})

	addr := ":" + cfg.Server.Port
	log.Printf("Starting server on %s", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
