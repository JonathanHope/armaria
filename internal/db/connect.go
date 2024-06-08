package db

import (
	"database/sql"
	"embed"
	"fmt"

	"github.com/jonathanhope/armaria/internal/null"
	"github.com/jonathanhope/armaria/internal/paths"
	"github.com/pressly/goose/v3"
	_ "modernc.org/sqlite"
)

//go:embed migrations/*.sql
var embedMigrations embed.FS

// connectDBFn is a function that connects to SQLite database
type connectDBFn func(inputPath null.NullString, configPath string) (*sql.DB, error)

// connectDB returns a connection to the bookmarks database.
// The DB will be created if it doesn't already exist.
// The DB will also be migrated up to the latest version, and have its PRAGMAS properly set.
func connectDB(inputPath null.NullString, configPath string) (*sql.DB, error) {
	dbLocation, err := paths.Database(inputPath, configPath)
	if err != nil {
		return nil, fmt.Errorf("error getting database location wile connecting to database: %w", err)
	}

	db, err := sql.Open("sqlite", dbLocation)
	if err != nil {
		return nil, fmt.Errorf("error while connecting to database: %w", err)
	}

	if err := configureDB(db); err != nil {
		return nil, fmt.Errorf("error configuring database wile connecting to database: %w", err)
	}

	if err := migrateDB(db); err != nil {
		return nil, fmt.Errorf("error applying migrations to database wile connecting to database: %w", err)
	}

	return db, nil
}

// configureDB uses PRAGMA to configure SQLite:
// - Use a WAL for journal mode
// - Enforce foreign keys.
func configureDB(db *sql.DB) error {
	if _, err := db.Exec("PRAGMA journal_mode=WAL"); err != nil {
		return fmt.Errorf("error turning WAL mode on while configuring database: %w", err)
	}

	if _, err := db.Exec("PRAGMA foreign_keys=on"); err != nil {
		return fmt.Errorf("error turning foreign key checking on while configuring database: %w", err)
	}

	return nil
}

// migrateDB brings the database up to the latest migration.
func migrateDB(db *sql.DB) error {
	goose.SetBaseFS(embedMigrations)
	goose.SetLogger(goose.NopLogger())

	if err := goose.SetDialect("sqlite3"); err != nil {
		return fmt.Errorf("error setting dialect while applying migrations: %w", err)
	}

	if err := goose.Up(db, "migrations"); err != nil {
		return fmt.Errorf("error while applying migrations: %w", err)
	}

	return nil
}
