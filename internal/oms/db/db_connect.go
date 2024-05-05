// Package db holds code related to making requests to the postgres db
package db

import (
	"database/sql"
	"os"

	_ "github.com/lib/pq"
	"github.com/pkg/errors"
)

var (
	errDatabaseURLEnvNotSet = errors.New("DATABASE_URL is not set")
	errConnectingToDatabase = errors.New("error connecting to database")
)

// NewDBContext creates a new db context.
func NewDBContext() (*Queries, error) {
	dbValue, err := setupDatabase()
	if err != nil {
		return nil, err
	}

	return New(dbValue), nil
}

func setupDatabase() (*sql.DB, error) {
	// Example DSN: "host=localhost user=devuser dbname=devdb password=devpassword sslmode=disable"
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		return nil, errDatabaseURLEnvNotSet
	}

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, errors.Wrap(err, "Error opening database")
	}

	err = db.Ping()
	if err != nil {
		return nil, errConnectingToDatabase
	}

	return db, nil
}
