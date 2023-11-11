package lib

import (
	"database/sql"
	"errors"
	"os"
	"path"
	"testing"

	"gopkg.in/DATA-DOG/go-sqlmock.v1"
)

/// queryWithTransactionCommitsTransaction

func TestQueryWithTransactionCommitsTransactionIfNoError(t *testing.T) {
	want := "value"

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("unexpected error: %+v", err)
	}

	mock.ExpectBegin()
	mock.ExpectCommit()
	mock.ExpectClose()

	got, err := queryWithTransaction[string](
		NullStringFrom("bookmarks.db"),
		func(dbFile NullString) (*sql.DB, error) {
			return db, nil
		},
		func(tx transaction) (string, error) {
			return want, nil
		},
	)

	if err != nil {
		t.Fatalf("unexpected error: %+v", err)
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
		t.Fatalf("unexpected error: %+v", err)
	}

	mock.ExpectBegin()
	mock.ExpectRollback()
	mock.ExpectClose()

	_, err = queryWithTransaction[string](
		NullStringFrom("bookmarks.db"),
		func(dbFile NullString) (*sql.DB, error) {
			return db, nil
		},
		func(tx transaction) (string, error) {
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
		t.Fatalf("unexpected error: %+v", err)
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

	_, _ = queryWithTransaction[string](
		NullStringFrom("bookmarks.db"),
		func(dbFile NullString) (*sql.DB, error) {
			return db, nil
		},
		func(tx transaction) (string, error) {
			panic("")
		},
	)

	t.Errorf("did not panic")
}

// execWithTransaction

func TestExecWithTransactionCommitsTransactionIfNoError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("unexpected error: %+v", err)
	}

	mock.ExpectBegin()
	mock.ExpectCommit()
	mock.ExpectClose()

	err = execWithTransaction(
		NullStringFrom("bookmarks.db"),
		func(dbFile NullString) (*sql.DB, error) {
			return db, nil
		},
		func(tx transaction) error {
			return nil
		},
	)

	if err != nil {
		t.Fatalf("unexpected error: %+v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled expections: %s", err)
	}
}

func TestExecWithTransactionRollsBackTransactionIfError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("unexpected error: %+v", err)
	}

	mock.ExpectBegin()
	mock.ExpectRollback()
	mock.ExpectClose()

	err = execWithTransaction(
		NullStringFrom("bookmarks.db"),
		func(dbFile NullString) (*sql.DB, error) {
			return db, nil
		},
		func(tx transaction) error {
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
		t.Fatalf("unexpected error: %+v", err)
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

	_ = execWithTransaction(
		NullStringFrom("bookmarks.db"),
		func(dbFile NullString) (*sql.DB, error) {
			return db, nil
		},
		func(tx transaction) error {
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
		t.Fatalf("unexpected error: %+v", err)
	}

	mock.ExpectClose()

	got, err := queryWithDB[string](
		NullStringFrom("bookmarks.db"),
		func(dbFile NullString) (*sql.DB, error) {
			return db, nil
		},
		func(tx transaction) (string, error) {
			return want, nil
		},
	)

	if err != nil {
		t.Fatalf("unexpected error: %+v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled expections: %s", err)
	}

	if got != want {
		t.Errorf("got %+v; want %+v", got, want)
	}
}

// getDefaultDB

func TestGetDefaultDB(t *testing.T) {
	type test struct {
		input         NullString
		goos          string
		folder        string
		db            string
		folderCreated bool
	}

	tests := []test{
		{
			input:         NullStringFromPtr(nil),
			goos:          "windows",
			folder:        "~/AppData/Local/Armaria",
			db:            "~/AppData/Local/Armaria/bookmarks.db",
			folderCreated: true,
		},
		{
			input:         NullStringFromPtr(nil),
			goos:          "linux",
			folder:        "~/.armaria",
			db:            "~/.armaria/bookmarks.db",
			folderCreated: true,
		},
		{
			input:         NullStringFromPtr(nil),
			goos:          "darwin",
			folder:        "~/Library/Application Support/Armaria",
			db:            "~/Library/Application Support/Armaria/bookmarks.db",
			folderCreated: true,
		},
		{
			input:         NullStringFrom("bookmarks.db"),
			db:            "bookmarks.db",
			folderCreated: false,
		},
	}

	userHome := func() (string, error) {
		return "~", nil
	}

	for _, tc := range tests {
		t.Run(tc.goos, func(t *testing.T) {
			folderCreated := false
			mkDirAll := func(path string, perm os.FileMode) error {
				folderCreated = true
				if path != tc.folder {
					t.Errorf("folder: got %+v; want %+v", path, tc.folder)
				}

				return nil
			}

			got, err := getDBLocation(tc.input, tc.goos, mkDirAll, userHome, path.Join)
			if err != nil {
				t.Fatalf("unexpected error: %+v", err)
			}

			if folderCreated != tc.folderCreated {
				t.Fatalf("folder created: got %+v; want %+v", folderCreated, tc.folderCreated)
			}

			if got != tc.db {
				t.Errorf("db: got %+v; want %+v", got, tc.db)
			}
		})
	}
}

// configureDB

func TestConfigureDB(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("unexpected error: %+v", err)
	}

	mock.ExpectExec("PRAGMA journal_mode=WAL").WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectExec("PRAGMA foreign_keys=on").WillReturnResult(sqlmock.NewResult(0, 0))

	if err := configureDB(db); err != nil {
		t.Fatalf("unexpected error: %+v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled expections: %s", err)
	}
}
