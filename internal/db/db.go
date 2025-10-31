package db

import (
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"api-employees-and-departments/config"
	"api-employees-and-departments/internal/infrastructure/logging"
)

func Connect(cfg *config.Config) (*gorm.DB, error) {
	// Create custom GORM logger that uses Zap for structured logging
	gormLogger := logging.NewGormLogger(time.Second)

	db, err := gorm.Open(postgres.Open(cfg.DSN()), &gorm.Config{
		Logger: gormLogger,
	})
	if err != nil {
		return nil, err
	}
	return db, nil
}
