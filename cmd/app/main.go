package main

import (
	"chat-api/internal/config"
	"chat-api/internal/database"
	"log"
)

func main() {
	cfg := config.Load()
	db, err := database.Connect(&cfg.Database)

	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
}
