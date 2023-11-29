package db

import (
	"database/sql"
	"errors"
	"testing"

	"github.com/jonathanhope/armaria/internal/null"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
)

// QueryWithTransaction

func TestQueryWithTransactionCommitsTransactionIfNoError(t *testing.T) {
	want := "value"

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	mock.ExpectBegin()
	mock.ExpectCommit()
	mock.ExpectClose()

	got, err := queryWithTransactionInternal[string](
		null.NullStringFrom("bookmarks.db"),
		"",
		func(inputPath null.NullString, confiPath string) (*sql.DB, error) {
			return db, nil
		},
		func(tx Transaction) (string, error) {
			return want, nil
		},
	)

	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled expections: %s", err)
	}

	if got != want {
		t.Errorf("got %+v; want %+v", got, want)
	}
}

func TestQueryWithTransactionRollsBackTransactionIfError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	mock.ExpectBegin()
	mock.ExpectRollback()
	mock.ExpectClose()

	_, err = queryWithTransactionInternal[string](
		null.NullStringFrom("bookmarks.db"),
		"",
		func(inputPath null.NullString, configPath string) (*sql.DB, error) {
			return db, nil
		},
		func(tx Transaction) (string, error) {
			return "", errors.New("test")
		},
	)

	if err == nil {
		t.Fatalf("unexpected success")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled expections: %s", err)
	}
}

func TestQueryWithTransactionRollsBackTransactionIfPanic(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	mock.ExpectBegin()
	mock.ExpectRollback()
	mock.ExpectClose()

	defer func() {
		_ = recover()

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("unfulfilled expections: %s", err)
		}
	}()

	_, _ = queryWithTransactionInternal[string](
		null.NullStringFrom("bookmarks.db"),
		"",
		func(inputPath null.NullString, configPath string) (*sql.DB, error) {
			return db, nil
		},
		func(tx Transaction) (string, error) {
			panic("")
		},
	)

	t.Errorf("did not panic")
}

// ExecWithTransaction

func TestExecWithTransactionCommitsTransactionIfNoError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	mock.ExpectBegin()
	mock.ExpectCommit()
	mock.ExpectClose()

	err = execWithTransactionInternal(
		null.NullStringFrom("bookmarks.db"),
		"",
		func(inputPath null.NullString, configPath string) (*sql.DB, error) {
			return db, nil
		},
		func(tx Transaction) error {
			return nil
		},
	)

	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled expections: %s", err)
	}
}

func TestExecWithTransactionRollsBackTransactionIfError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	mock.ExpectBegin()
	mock.ExpectRollback()
	mock.ExpectClose()

	err = execWithTransactionInternal(
		null.NullStringFrom("bookmarks.db"),
		"",
		func(inputPath null.NullString, configPath string) (*sql.DB, error) {
			return db, nil
		},
		func(tx Transaction) error {
			return errors.New("test")
		},
	)

	if err == nil {
		t.Fatalf("unexpected success")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled expections: %s", err)
	}
}

func TestExecWithTransactionRollsBackTransactionIfPanic(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	mock.ExpectBegin()
	mock.ExpectRollback()
	mock.ExpectClose()

	defer func() {
		_ = recover()

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("unfulfilled expections: %s", err)
		}
	}()

	_ = execWithTransactionInternal(
		null.NullStringFrom("bookmarks.db"),
		"",
		func(inputPath null.NullString, configPath string) (*sql.DB, error) {
			return db, nil
		},
		func(tx Transaction) error {
			panic("")
		},
	)

	t.Errorf("did not panic")
}

// queryWithDB

func TestQueryWithDB(t *testing.T) {
	want := "value"

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	mock.ExpectClose()

	got, err := queryWithDBInternal[string](
		null.NullStringFrom("bookmarks.db"),
		"",
		func(inputPath null.NullString, configPath string) (*sql.DB, error) {
			return db, nil
		},
		func(tx Transaction) (string, error) {
			return want, nil
		},
	)

	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled expections: %s", err)
	}

	if got != want {
		t.Errorf("got %+v; want %+v", got, want)
	}
}
