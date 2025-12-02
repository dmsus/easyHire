package database

import (
	"github.com/easyhire/backend/internal/models"
	"github.com/easyhire/backend/internal/pkg/logger"
	"gorm.io/gorm"
)

func AutoMigrate(db *gorm.DB) error {
	log := logger.Global().With().Str("component", "migrations").Logger()
	log.Info().Msg("Running GORM auto-migrations...")

	modelsToMigrate := []interface{}{
		&models.User{},
		&models.Competency{},
		&models.LevelWeight{},
		&models.LevelThreshold{},
		&models.Assessment{},
		&models.AssessmentCompetency{},
		&models.AssessmentQuestion{},
		&models.Candidate{},
		&models.Question{},
		&models.Result{},
		&models.Invitation{},
		&models.CodeExecution{},
	}

	for _, model := range modelsToMigrate {
		log.Debug().Msgf("Migrating: %T", model)
		if err := db.AutoMigrate(model); err != nil {
			return err
		}
	}

	log.Info().Msg("✅ GORM auto-migrations completed successfully")
	return nil
}
