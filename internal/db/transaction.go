package db

import (
	"database/sql"
	"fmt"

	"github.com/blockloop/scan/v2"
	"github.com/jonathanhope/armaria/internal/null"
	"github.com/nullism/bqb"
)

// The abstractions in this file simplify working with transactions in Go.

// Transaction represents a database/sql connection.
type Transaction interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
}

// QueryTxFn is a function that operates on a transaction which returns results.
type QueryTxFn[T any] func(Transaction) (T, error)

// queryWithTransaction creates a scope for functions that operate on a transaction which returns results.
// Handles connecting to the DB, creating the transaction, committing/rolling back, and closing the connection.
// Will also handle creating the DB if it doesn't exist and applying missing migrations to it.
func QueryWithTransaction[T any](inputPath null.NullString, configPath string, queryFn QueryTxFn[T]) (val T, err error) {
	return queryWithTransactionInternal[T](inputPath, configPath, connectDB, queryFn)
}

// queryWithTransactionInternal enables DI for QueryWithTransaction.
func queryWithTransactionInternal[T any](inputPath null.NullString, configPath string, connectFn connectDBFn, queryFn QueryTxFn[T]) (val T, err error) {
	db, err := connectFn(inputPath, configPath)
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

// ExecTxFn is a function that operates on a transaction which doesn't return results.
type ExecTxFn func(Transaction) error

// execWithTransaction creates a scope for functions that operate on a transaction which doesn't return results.
// Handles connecting to the DB, creating the transaction, committing/rolling back, and closing the connection.
// Will also handle creating the DB if it doesn't exist and applying missing migrations to it.
func ExecWithTransaction(inputPath null.NullString, configPath string, execFn ExecTxFn) (err error) {
	return execWithTransactionInternal(inputPath, configPath, connectDB, execFn)
}

// execWithTransactionInternal enables DI For ExecWithTransaction.
func execWithTransactionInternal(inputPath null.NullString, configPath string, connectFn connectDBFn, execFn ExecTxFn) (err error) {
	db, err := connectFn(inputPath, configPath)
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

// QueryWithDB creates a scope for for functions that operate on a database connection which return results.
// Handles connecting to the DB and closing the connection.
// Will also handle creating the DB if it doesn't exist and applying missing migrations to it.
func QueryWithDB[T any](inputPath null.NullString, configPath string, queryFn QueryTxFn[T]) (T, error) {
	return queryWithDBInternal(inputPath, configPath, connectDB, queryFn)
}

// queryWithDBInternal enables DI for QueryWithDB.
func queryWithDBInternal[T any](inputPath null.NullString, configPath string, connectFn connectDBFn, queryFn QueryTxFn[T]) (T, error) {
	var val T

	db, err := connectFn(inputPath, configPath)
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

// query runs SQL that returns results.
func query[T any](tx Transaction, q *bqb.Query) ([]T, error) {
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
func exec(tx Transaction, q *bqb.Query) error {
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
func count(tx Transaction, q *bqb.Query) (int, error) {
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
