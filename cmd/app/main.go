package main

import (
	"chat-api/internal/config"
	"chat-api/internal/database"
	"chat-api/internal/handlers"
	"chat-api/internal/repository"
	"chat-api/internal/service"
	"log"
	"net/http"
)

func main() {
	cfg := config.Load()
	db, err := database.Connect(&cfg.Database)

	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	repo := repository.NewRepository(db)
	srv := service.NewChatService(repo)
	handler := handlers.NewHandler(srv)

	mux := http.NewServeMux()

	mux.Handle("/chats/", handlers.LoggingMiddleware(handler.HandleChats))
	mux.Handle("/chats", handlers.LoggingMiddleware(handler.HandleChats))

	addr := ":" + cfg.Server.Port
	log.Printf("Starting server on %s", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
