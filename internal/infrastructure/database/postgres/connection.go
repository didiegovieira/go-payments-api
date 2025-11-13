package postgres

import (
	"database/sql"
	"fmt"
	"go-payments-api/internal/settings"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	_ "github.com/lib/pq"
)

type DB struct {
	conn *sql.DB
}

func NewConnection() (*DB, error) {
	connStr := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		settings.Settings.Database.Host,
		settings.Settings.Database.Port,
		settings.Settings.Database.User,
		settings.Settings.Database.Password,
		settings.Settings.Database.Name,
	)

	conn, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	conn.SetMaxOpenConns(25)
	conn.SetMaxIdleConns(5)
	conn.SetConnMaxLifetime(time.Hour)
	conn.SetConnMaxIdleTime(time.Minute * 30)

	if err := conn.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	log.Println("✅ Database connection established")

	// Run migrations automatically
	if err := runMigrationsFromFiles(conn); err != nil {
		return nil, fmt.Errorf("failed to run migrations: %w", err)
	}

	return &DB{conn: conn}, nil
}

func runMigrationsFromFiles(conn *sql.DB) error {
	log.Println("🔄 Running database migrations...")

	// Create schema_migrations table
	_, err := conn.Exec(`
        CREATE TABLE IF NOT EXISTS schema_migrations (
            version INTEGER PRIMARY KEY,
            filename VARCHAR(255) NOT NULL,
            applied_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
        )
    `)
	if err != nil {
		return fmt.Errorf("failed to create schema_migrations table: %w", err)
	}

	// Get migration files from scripts/migrations directory
	migrationsPath := "scripts/migrations"
	files, err := os.ReadDir(migrationsPath)
	if err != nil {
		return fmt.Errorf("failed to read migrations directory: %w", err)
	}

	// Filter and sort SQL files
	var sqlFiles []os.DirEntry
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".sql") {
			sqlFiles = append(sqlFiles, file)
		}
	}

	if len(sqlFiles) == 0 {
		log.Println("⚠️  No migration files found in scripts/migrations")
		return nil
	}

	// Sort files by name (assuming they start with version number like 001_, 002_, etc.)
	sort.Slice(sqlFiles, func(i, j int) bool {
		return sqlFiles[i].Name() < sqlFiles[j].Name()
	})

	// Apply each migration
	for idx, file := range sqlFiles {
		version := idx + 1
		filename := file.Name()

		// Check if migration already applied
		var exists bool
		err := conn.QueryRow(
			"SELECT EXISTS(SELECT 1 FROM schema_migrations WHERE version = $1)",
			version,
		).Scan(&exists)
		if err != nil {
			return fmt.Errorf("failed to check migration %d: %w", version, err)
		}

		if exists {
			log.Printf("⏭️  Skipping migration %d: %s (already applied)", version, filename)
			continue
		}

		log.Printf("⚙️  Applying migration %d: %s", version, filename)

		// Read SQL file
		sqlContent, err := os.ReadFile(filepath.Join(migrationsPath, filename))
		if err != nil {
			return fmt.Errorf("failed to read migration file %s: %w", filename, err)
		}

		// Begin transaction
		tx, err := conn.Begin()
		if err != nil {
			return fmt.Errorf("failed to begin transaction for migration %d: %w", version, err)
		}

		// Execute migration SQL
		if _, err := tx.Exec(string(sqlContent)); err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to execute migration %d (%s): %w", version, filename, err)
		}

		// Record migration
		if _, err := tx.Exec(
			"INSERT INTO schema_migrations (version, filename) VALUES ($1, $2)",
			version,
			filename,
		); err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to record migration %d: %w", version, err)
		}

		// Commit transaction
		if err := tx.Commit(); err != nil {
			return fmt.Errorf("failed to commit migration %d: %w", version, err)
		}

		log.Printf("✅ Applied migration %d: %s", version, filename)
	}

	log.Println("✅ All migrations completed successfully")
	return nil
}

func (db *DB) GetConnection() *sql.DB {
	return db.conn
}

func (db *DB) Close() error {
	if db.conn != nil {
		log.Printf("Closing PostgreSQL connection")
		return db.conn.Close()
	}
	return nil
}
