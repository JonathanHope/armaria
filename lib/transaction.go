package lib

import (
	"database/sql"
	"embed"
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
func queryWithTransaction[T any](dbFile NullString, connectFn connectDBFn, queryFn queryTxFn[T]) (val T, err error) {
	db, err := connectFn(dbFile)
	if err != nil {
		return
	}
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		err = fmt.Errorf("%w: %w", ErrUnexpected, err)
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
				err = fmt.Errorf("%w: %w", ErrUnexpected, err)
			}
		}
	}()

	val, err = queryFn(tx)
	return val, err
}

// QueryTxFn is a function that operates on a transaction which doesn't return results.
type execTxFn func(transaction) error

// execWithTransaction creates a scope for functions that operate on a transaction which doesn't return results.
// Handles connecting to the DB, creating the transaction, committing/rolling back, and closing the connection.
// Will also handle creating the DB if it doesn't exist and applying missing migrations to it.
func execWithTransaction(dbFile NullString, connectFn connectDBFn, execFn execTxFn) (err error) {
	db, err := connectFn(dbFile)
	if err != nil {
		return
	}
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		err = fmt.Errorf("%w: %w", ErrUnexpected, err)
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
				err = fmt.Errorf("%w: %w", ErrUnexpected, err)
			}
		}
	}()

	err = execFn(tx)
	return err
}

// queryWithDB creates a scope for for functions that operate on a database connection which return results.
// Handles connecting to the DB and closing the connection.
// Will also handle creating the DB if it doesn't exist and applying missing migrations to it.
func queryWithDB[T any](dbFile NullString, connectFn connectDBFn, queryFn queryTxFn[T]) (T, error) {
	var val T

	db, err := connectFn(dbFile)
	if err != nil {
		return val, err
	}
	defer db.Close()

	return queryFn(db)
}

// connectDBFn is a function that connects to SQLite database.
type connectDBFn func(dbFile NullString) (*sql.DB, error)

// connectDB returns a connection to the bookmarks database.
func connectDB(dbFile NullString) (*sql.DB, error) {
	dbLocation, err := getDBLocation(dbFile, runtime.GOOS, os.MkdirAll, os.UserHomeDir, filepath.Join)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrUnexpected, err)
	}

	db, err := sql.Open("sqlite3", dbLocation)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrUnexpected, err)
	}

	if err := configureDB(db); err != nil {
		return nil, fmt.Errorf("%w: %w", ErrUnexpected, err)
	}

	if err := migrateDB(db); err != nil {
		return nil, fmt.Errorf("%w: %w", ErrUnexpected, err)
	}

	return db, nil
}

// mkDirAllFn creates a directory if it doesn't alread exist.
type mkDirAllFn func(path string, perm os.FileMode) error

// userHomeFn returns the home directory of the current user.
type userHomeFn func() (string, error)

// joinFn joins path segments together.
type joinFn func(elem ...string) string

// getDBLocation returns the default location for the bookmarks database.
func getDBLocation(dbFile NullString, goos string, mkDirAll mkDirAllFn, userHome userHomeFn, join joinFn) (string, error) {
	if !dbFile.Valid || !dbFile.Dirty {
		home, err := userHome()
		if err != nil {
			return "", err
		}

		var folder string
		if goos == "linux" {
			folder = join(home, ".armaria")
		} else if goos == "windows" {
			folder = join(home, "AppData", "Local", "Armaria")
		} else if goos == "darwin" {
			folder = join(home, "Library", "Application Support", "Armaria")
		} else {
			panic("Unsupported operating system")
		}

		err = mkDirAll(folder, os.ModePerm)
		if err != nil {
			return "", err
		}

		return filepath.Join(folder, "bookmarks.db"), nil
	} else {
		return dbFile.String, nil
	}
}

// configureDB uses PRAGMA to configure SQLite:
// - Use a WAL for journal mode
// - Enforce foreign keys.
func configureDB(db *sql.DB) error {
	if _, err := db.Exec("PRAGMA journal_mode=WAL"); err != nil {
		return err
	}

	if _, err := db.Exec("PRAGMA foreign_keys=on"); err != nil {
		return err
	}

	return nil
}

// migrateDB brings the database up to the latest migration.
func migrateDB(db *sql.DB) error {
	goose.SetBaseFS(embedMigrations)
	goose.SetLogger(goose.NopLogger())

	if err := goose.SetDialect("sqlite3"); err != nil {
		return err
	}

	if err := goose.Up(db, "migrations"); err != nil {
		return err
	}

	return nil
}

// query runs SQL that returns results.
func query[T any](tx transaction, q *bqb.Query) ([]T, error) {
	var results []T

	sql, args, err := q.ToSql()
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrUnexpected, err)
	}

	rows, err := tx.Query(sql, args...)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrUnexpected, err)
	}

	defer rows.Close()

	err = scan.RowsStrict(&results, rows)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrUnexpected, err)
	}

	return results, nil
}

// exec runs SQL that does't return results.
func exec(tx transaction, q *bqb.Query) error {
	sql, args, err := q.ToSql()
	if err != nil {
		return fmt.Errorf("%w: %w", ErrUnexpected, err)
	}

	if _, err = tx.Exec(sql, args...); err != nil {
		return fmt.Errorf("%w: %w", ErrUnexpected, err)
	}

	return nil
}

// count runs SQL that returns a count.
func count(tx transaction, q *bqb.Query) (int, error) {
	sql, args, err := q.ToSql()
	if err != nil {
		return 0, fmt.Errorf("%w: %w", ErrUnexpected, err)
	}

	var count = 0
	if err = tx.QueryRow(sql, args...).Scan(&count); err != nil {
		return 0, fmt.Errorf("%w: %w", ErrUnexpected, err)
	}

	return count, nil
}
