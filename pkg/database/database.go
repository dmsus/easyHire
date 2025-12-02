package database

import (
	"fmt"
	"time"

	"github.com/easyhire/backend/internal/pkg/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Database struct {
	DB *gorm.DB
}

func NewDatabase(config *config.DatabaseConfig) (*Database, error) {
	if config.Host == "" {
		// Running without database
		return &Database{DB: nil}, nil
	}

	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		config.Host, config.Port, config.User, config.Password, config.DBName, config.SSLMode)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get sql.DB: %w", err)
	}

	sqlDB.SetMaxOpenConns(25)
	sqlDB.SetMaxIdleConns(25)
	sqlDB.SetConnMaxLifetime(5 * time.Minute)

	// Run auto-migrations
	if err := AutoMigrate(db); err != nil {
		return nil, fmt.Errorf("failed to run auto-migrations: %w", err)
	}

	// Run SQL migrations
	if err := RunMigrations(db); err != nil {
		return nil, fmt.Errorf("failed to run SQL migrations: %w", err)
	}

	return &Database{DB: db}, nil
}

func (d *Database) Close() error {
	if d.DB == nil {
		return nil
	}
	
	sqlDB, err := d.DB.DB()
	if err != nil {
		return fmt.Errorf("failed to get sql.DB: %w", err)
	}
	return sqlDB.Close()
}

func (d *Database) HealthCheck() error {
	if d.DB == nil {
		return fmt.Errorf("database not configured")
	}
	
	sqlDB, err := d.DB.DB()
	if err != nil {
		return fmt.Errorf("failed to get sql.DB: %w", err)
	}
	return sqlDB.Ping()
}
