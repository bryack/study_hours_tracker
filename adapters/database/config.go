package database

import (
	"fmt"
	"os"
)

func SetupPostgres() (*PostgresSubjectStore, error) {
	connStr := os.Getenv("DATABASE_URL")
	if connStr == "" {
		connStr = "postgres://postgres:pass@localhost:5432/mydb?sslmode=disable"
	}
	store, err := NewPostgresSubjectStore(connStr)
	if err != nil {
		return nil, fmt.Errorf("could not connect to database: %w", err)
	}
	return store, nil
}
