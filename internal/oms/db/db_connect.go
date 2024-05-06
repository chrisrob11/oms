// Package db holds code related to making requests to the postgres db
package db

import (
	"database/sql"
	"log/slog"
	"os"

	_ "github.com/lib/pq"
	"github.com/pkg/errors"
)

var (
	errDatabaseURLEnvNotSet = errors.New("DATABASE_URL is not set")
)

// NewDBContext creates a new db context.
func NewDBContext(logger *slog.Logger) (*Queries, error) {
	dbValue, err := setupDatabase(logger)
	if err != nil {
		return nil, err
	}

	return New(dbValue), nil
}

func setupDatabase(logger *slog.Logger) (*sql.DB, error) {
	// Example DSN: "host=localhost user=devuser dbname=devdb password=devpassword sslmode=disable"
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		return nil, errDatabaseURLEnvNotSet
	}

	logger.Info("Database dsn", slog.Attr{Key: "dsn", Value: slog.StringValue(dsn)})

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, errors.Wrap(err, "Error opening database")
	}

	err = db.Ping()
	if err != nil {
		return nil, errors.Wrap(err, "error connecting to database")
	}

	return db, nil
}
