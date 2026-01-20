package database

import (
	"chat-api/internal/config"
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Connect to database
//
//	@param cfg *config.DatabaseConfig
//	@return *gorm.DB
//	@return error
func Connect(cfg *config.DatabaseConfig) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(cfg.DSN()), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		return nil, fmt.Errorf("failed to connection to database: %w", err)
	}

	log.Println("Database connection successfully")
	return db, nil
}
