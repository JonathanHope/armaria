package db

import (
	"testing"

	"gopkg.in/DATA-DOG/go-sqlmock.v1"
)

func TestConfigureDB(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	mock.ExpectExec("PRAGMA journal_mode=WAL").WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectExec("PRAGMA foreign_keys=on").WillReturnResult(sqlmock.NewResult(0, 0))

	if err := configureDB(db); err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled expections: %s", err)
	}
}
