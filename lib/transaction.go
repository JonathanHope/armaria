package lib

import (
	"database/sql"
	"embed"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/blockloop/scan/v2"
	_ "github.com/mattn/go-sqlite3"
	"github.com/nullism/bqb"
	"github.com/pressly/goose/v3"
)

// The abstractions in this file simplify working with transactions in Go.

//go:embed migrations/*.sql
var embedMigrations embed.FS

// transaction represents a database/sql connection.
type transaction interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
}

// queryTxFn is a function that operates on a transaction which returns results.
type queryTxFn[T any] func(transaction) (T, error)

// queryWithTransaction creates a scope for functions that operate on a transaction which returns results.
// Handles connecting to the DB, creating the transaction, committing/rolling back, and closing the connection.
// Will also handle creating the DB if it doesn't exist and applying missing migrations to it.
func queryWithTransaction[T any](inputPath NullString, connectFn connectDBFn, queryFn queryTxFn[T]) (val T, err error) {
	db, err := connectFn(inputPath)
	if err != nil {
		err = fmt.Errorf("error connecting to database while querying with a transaction: %w", err)
		return
	}
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		err = fmt.Errorf("error beginning transaction while querying with a transaction: %w", err)
		return
	}

	defer func() {
		if p := recover(); p != nil {
			// a panic occurred, rollback and repanic
			//nolint:errcheck
			tx.Rollback()
			panic(p)
		} else if err != nil {
			// something went wrong, rollback
			//nolint:errcheck
			tx.Rollback()
		} else {
			// all good, commit
			err = tx.Commit()
			if err != nil {
				err = fmt.Errorf("error committing transaction while querying with a transaction: %w", err)
			}
		}
	}()

	val, err = queryFn(tx)
	if err != nil {
		return val, fmt.Errorf("error while querying with a transaction: %w", err)
	}

	return val, nil
}

// QueryTxFn is a function that operates on a transaction which doesn't return results.
type execTxFn func(transaction) error

// execWithTransaction creates a scope for functions that operate on a transaction which doesn't return results.
// Handles connecting to the DB, creating the transaction, committing/rolling back, and closing the connection.
// Will also handle creating the DB if it doesn't exist and applying missing migrations to it.
func execWithTransaction(inputPath NullString, connectFn connectDBFn, execFn execTxFn) (err error) {
	db, err := connectFn(inputPath)
	if err != nil {
		err = fmt.Errorf("error connecting to database while executing with a transaction: %w", err)
		return
	}
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		err = fmt.Errorf("error beginning transaction while executing with a transaction: %w", err)
		return
	}

	defer func() {
		if p := recover(); p != nil {
			// a panic occurred, rollback and repanic
			//nolint:errcheck
			tx.Rollback()
			panic(p)
		} else if err != nil {
			// something went wrong, rollback
			//nolint:errcheck
			tx.Rollback()
		} else {
			// all good, commit
			err = tx.Commit()
			if err != nil {
				err = fmt.Errorf("error committing transaction while executing with a transaction: %w", err)
			}
		}
	}()

	err = execFn(tx)
	if err != nil {
		return fmt.Errorf("error while executing with a transaction: %w", err)
	}

	return nil
}

// queryWithDB creates a scope for for functions that operate on a database connection which return results.
// Handles connecting to the DB and closing the connection.
// Will also handle creating the DB if it doesn't exist and applying missing migrations to it.
func queryWithDB[T any](inputPath NullString, connectFn connectDBFn, queryFn queryTxFn[T]) (T, error) {
	var val T

	db, err := connectFn(inputPath)
	if err != nil {
		return val, fmt.Errorf("error connecting to database while querying with database: %w", err)
	}
	defer db.Close()

	val, err = queryFn(db)
	if err != nil {
		return val, fmt.Errorf("error while querying with database: %w", err)
	}

	return val, nil
}

// connectDBFn is a function that connects to SQLite database
type connectDBFn func(inputPath NullString) (*sql.DB, error)

// connectDB returns a connection to the bookmarks database.
func connectDB(inputPath NullString) (*sql.DB, error) {
	config, err := GetConfig()
	if err != nil && !errors.Is(err, ErrConfigMissing) {
		return nil, fmt.Errorf("error getting config wile connecting to database: %w", err)
	}

	dbLocation, err := getDatabasePath(inputPath, config.DB, runtime.GOOS, os.MkdirAll, os.UserHomeDir, filepath.Join)
	if err != nil {
		return nil, fmt.Errorf("error getting database location wile connecting to database: %w", err)
	}

	db, err := sql.Open("sqlite3", dbLocation)
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

// query runs SQL that returns results.
func query[T any](tx transaction, q *bqb.Query) ([]T, error) {
	var results []T

	sql, args, err := q.ToSql()
	if err != nil {
		return nil, fmt.Errorf("error building query while querying: %w", err)
	}

	rows, err := tx.Query(sql, args...)
	if err != nil {
		return nil, fmt.Errorf("error while querying: %w", err)
	}

	defer rows.Close()

	err = scan.RowsStrict(&results, rows)
	if err != nil {
		return nil, fmt.Errorf("error building results while querying: %w", err)
	}

	return results, nil
}

// exec runs SQL that does't return results.
func exec(tx transaction, q *bqb.Query) error {
	sql, args, err := q.ToSql()
	if err != nil {
		return fmt.Errorf("error building query while executing: %w", err)
	}

	if _, err = tx.Exec(sql, args...); err != nil {
		return fmt.Errorf("error while executing: %w", err)
	}

	return nil
}

// count runs SQL that returns a count.
func count(tx transaction, q *bqb.Query) (int, error) {
	sql, args, err := q.ToSql()
	if err != nil {
		return 0, fmt.Errorf("error building query while counting: %w", err)
	}

	var count = 0
	if err = tx.QueryRow(sql, args...).Scan(&count); err != nil {
		return 0, fmt.Errorf("error while counting: %w", err)
	}

	return count, nil
}
