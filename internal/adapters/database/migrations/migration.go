package migrations

import (
	"github.com/bonarizki-dat/boilerplate-gin-dat/internal/adapters/database"
	"github.com/bonarizki-dat/boilerplate-gin-dat/internal/domain/models"
	"github.com/bonarizki-dat/boilerplate-gin-dat/pkg/logger"
)

// Migrate runs database migrations for all models.
//
// Currently uses GORM AutoMigrate for development convenience.
// For production, consider using golang-migrate or similar versioned migration tool.
//
// Models are migrated in order to handle foreign key dependencies.
func Migrate() {
	var migrationModels = []interface{}{
		&models.User{},    // Users table (for authentication)
		&models.Example{}, // Example table
	}

	logger.Infof("Starting database migrations...")

	err := database.DB.AutoMigrate(migrationModels...)
	if err != nil {
		logger.Errorf("migration failed: %v", err)
		return
	}

	logger.Infof("Database migrations completed successfully")
}
