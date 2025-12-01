package database

import (
	"context"
	"fmt"
	"time"

	"github.com/easyhire/backend/internal/pkg/config"
	"github.com/rs/zerolog/log"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Database struct {
	DB *gorm.DB
}

func NewDatabase(cfg *config.DatabaseConfig) (*Database, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Name, cfg.SSLMode,
	)

	gormConfig := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent), // Custom logger will be set
		PrepareStmt: true,
		NowFunc: func() time.Time {
			return time.Now().UTC()
		},
	}

	db, err := gorm.Open(postgres.Open(dsn), gormConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Get underlying SQL DB
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get database instance: %w", err)
	}

	// Set connection pool settings
	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	sqlDB.SetConnMaxLifetime(cfg.ConnMaxLifetime)

	// Test connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := sqlDB.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("database ping failed: %w", err)
	}

	log.Info().Msg("✅ Database connection established")

	return &Database{DB: db}, nil
}

func (d *Database) Close() error {
	sqlDB, err := d.DB.DB()
	if err != nil {
		return fmt.Errorf("failed to get database instance: %w", err)
	}

	if err := sqlDB.Close(); err != nil {
		return fmt.Errorf("failed to close database connection: %w", err)
	}

	log.Info().Msg("🔌 Database connection closed")
	return nil
}

func (d *Database) Migrate(models ...interface{}) error {
	if err := d.DB.AutoMigrate(models...); err != nil {
		return fmt.Errorf("failed to migrate database: %w", err)
	}

	log.Info().Int("models_count", len(models)).Msg("✅ Database migration completed")
	return nil
}

func (d *Database) WithTransaction(fn func(*gorm.DB) error) error {
	return d.DB.Transaction(fn)
}

func (d *Database) HealthCheck() error {
	sqlDB, err := d.DB.DB()
	if err != nil {
		return fmt.Errorf("failed to get database instance: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := sqlDB.PingContext(ctx); err != nil {
		return fmt.Errorf("database health check failed: %w", err)
	}

	return nil
}

// Custom query builder for common operations
type QueryBuilder struct {
	db *gorm.DB
}

func NewQueryBuilder(db *gorm.DB) *QueryBuilder {
	return &QueryBuilder{db: db}
}

func (q *QueryBuilder) Where(condition interface{}, args ...interface{}) *QueryBuilder {
	q.db = q.db.Where(condition, args...)
	return q
}

func (q *QueryBuilder) Order(order string) *QueryBuilder {
	q.db = q.db.Order(order)
	return q
}

func (q *QueryBuilder) Limit(limit int) *QueryBuilder {
	q.db = q.db.Limit(limit)
	return q
}

func (q *QueryBuilder) Offset(offset int) *QueryBuilder {
	q.db = q.db.Offset(offset)
	return q
}

func (q *QueryBuilder) Preload(query string, args ...interface{}) *QueryBuilder {
	q.db = q.db.Preload(query, args...)
	return q
}

func (q *QueryBuilder) Find(dest interface{}) error {
	return q.db.Find(dest).Error
}

func (q *QueryBuilder) First(dest interface{}) error {
	return q.db.First(dest).Error
}

func (q *QueryBuilder) Count(count *int64) error {
	return q.db.Count(count).Error
}

func (q *QueryBuilder) Delete(model interface{}) error {
	return q.db.Delete(model).Error
}

func (q *QueryBuilder) Update(column string, value interface{}) error {
	return q.db.Update(column, value).Error
}
