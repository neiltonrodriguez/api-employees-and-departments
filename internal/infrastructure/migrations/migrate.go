package migrations

import (
	"api-employees-and-departments/internal/domain/department"
	"api-employees-and-departments/internal/domain/employee"
	"fmt"

	"gorm.io/gorm"
)

// Run executes all database migrations
func Run(db *gorm.DB) error {
	// Enable UUID extension for PostgreSQL
	if err := db.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\"").Error; err != nil {
		return fmt.Errorf("failed to create uuid extension: %w", err)
	}

	// Run AutoMigrate for all models
	if err := db.AutoMigrate(
		&department.Department{},
		&employee.Employee{},
	); err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	return nil
}
