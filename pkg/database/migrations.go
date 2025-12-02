package database

import (
	"embed"
	"fmt"
	"io/fs"
	"sort"
	"strings"

	"github.com/easyhire/backend/internal/pkg/logger"
	"gorm.io/gorm"
)

//go:embed migrations/*.sql
var migrationsFS embed.FS

type Migration struct {
	Version int64
	Name    string
	Content string
}

func RunMigrations(db *gorm.DB) error {
	if db == nil {
		logger.Global().Warn().Msg("Database not configured, skipping migrations")
		return nil
	}

	log := logger.Global().With().Str("component", "migrations").Logger()
	log.Info().Msg("Starting database migrations...")

	// Get SQL DB connection
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("failed to get sql.DB: %w", err)
	}

	// Read migration files
	migrations, err := loadMigrations()
	if err != nil {
		return fmt.Errorf("failed to load migrations: %w", err)
	}

	// Sort migrations by version
	sort.Slice(migrations, func(i, j int) bool {
		return migrations[i].Version < migrations[j].Version
	})

	// Get applied migrations
	appliedVersions, err := getAppliedVersions(db)
	if err != nil {
		return fmt.Errorf("failed to get applied migrations: %w", err)
	}

	// Apply new migrations
	for _, migration := range migrations {
		if _, applied := appliedVersions[migration.Version]; applied {
			log.Debug().Int64("version", migration.Version).Str("name", migration.Name).Msg("Migration already applied")
			continue
		}

		log.Info().Int64("version", migration.Version).Str("name", migration.Name).Msg("Applying migration")

		// Start transaction
		tx := db.Begin()
		if tx.Error != nil {
			return fmt.Errorf("failed to start transaction: %w", tx.Error)
		}

		// Execute migration
		if err := tx.Exec(migration.Content).Error; err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to execute migration %d: %w", migration.Version, err)
		}

		// Record migration
		recordSQL := `INSERT INTO schema_migrations (version, name) VALUES (?, ?)`
		if err := tx.Exec(recordSQL, migration.Version, migration.Name).Error; err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to record migration %d: %w", migration.Version, err)
		}

		// Commit transaction
		if err := tx.Commit().Error; err != nil {
			return fmt.Errorf("failed to commit migration %d: %w", migration.Version, err)
		}

		log.Info().Int64("version", migration.Version).Str("name", migration.Name).Msg("Migration applied successfully")
	}

	log.Info().Msg("All migrations completed successfully")
	return nil
}

func loadMigrations() ([]Migration, error) {
	var migrations []Migration

	// Read migration files
	files, err := fs.ReadDir(migrationsFS, "migrations")
	if err != nil {
		// No migration files yet, that's OK
		return migrations, nil
	}

	for _, file := range files {
		if !strings.HasSuffix(file.Name(), ".sql") {
			continue
		}

		// Parse version from filename (e.g., 001_init_schema.sql -> 1)
		var version int64
		var name string
		if _, err := fmt.Sscanf(file.Name(), "%d_%s", &version, &name); err != nil {
			return nil, fmt.Errorf("invalid migration filename: %s", file.Name())
		}

		// Remove .sql extension from name
		name = strings.TrimSuffix(name, ".sql")

		// Read migration content
		content, err := fs.ReadFile(migrationsFS, "migrations/"+file.Name())
		if err != nil {
			return nil, fmt.Errorf("failed to read migration file %s: %w", file.Name(), err)
		}

		migrations = append(migrations, Migration{
			Version: version,
			Name:    name,
			Content: string(content),
		})
	}

	return migrations, nil
}

func getAppliedVersions(db *gorm.DB) (map[int64]bool, error) {
	applied := make(map[int64]bool)

	// Check if schema_migrations table exists
	var tableExists bool
	checkSQL := `SELECT EXISTS (
		SELECT FROM information_schema.tables 
		WHERE table_schema = 'public' 
		AND table_name = 'schema_migrations'
	)`

	if err := db.Raw(checkSQL).Scan(&tableExists).Error; err != nil {
		return nil, fmt.Errorf("failed to check schema_migrations table: %w", err)
	}

	if !tableExists {
		return applied, nil
	}

	// Get applied migrations
	rows, err := db.Raw("SELECT version FROM schema_migrations ORDER BY version").Rows()
	if err != nil {
		return nil, fmt.Errorf("failed to query applied migrations: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var version int64
		if err := rows.Scan(&version); err != nil {
			return nil, fmt.Errorf("failed to scan migration version: %w", err)
		}
		applied[version] = true
	}

	return applied, nil
}
