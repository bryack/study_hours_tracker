package database

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/bryack/study_hours_tracker/domain"
	_ "github.com/jackc/pgx/v5/stdlib"
)

const (
	createTableQuery = `CREATE TABLE IF NOT EXISTS subjects (
	id SERIAL PRIMARY KEY,
	subject TEXT NOT NULL UNIQUE,
	hours INTEGER NOT NULL DEFAULT 0
	);`
	selectHoursQuery = "SELECT hours FROM subjects WHERE subject = $1"
	insertHoursQuery = `INSERT INTO subjects (subject, hours) 
	VALUES ($1, $2)
	ON CONFLICT (subject)
	DO UPDATE SET hours = subjects.hours + EXCLUDED.hours`
	selectReportQuery = "SELECT subject, hours FROM subjects ORDER BY hours DESC"
	driverName        = "pgx"
)

type PostgresSubjectStore struct {
	db *sql.DB
}

func NewPostgresSubjectStore(connStr string) (*PostgresSubjectStore, error) {
	store := &PostgresSubjectStore{}
	if err := store.initDatabase(connStr); err != nil {
		return nil, fmt.Errorf("failed to init DB: %w", err)
	}
	if err := store.createTable(); err != nil {
		store.db.Close()
		return nil, fmt.Errorf("failed to create DB: %w", err)
	}
	return store, nil
}

func (ps *PostgresSubjectStore) initDatabase(connStr string) error {
	db, err := sql.Open(driverName, connStr)
	if err != nil {
		return fmt.Errorf("failed to open DB with %q: %w", connStr, err)
	}
	ps.db = db

	if err := db.Ping(); err != nil {
		return fmt.Errorf("failed to ping db: %w", err)
	}

	return nil
}

func (ps *PostgresSubjectStore) createTable() error {
	if _, err := ps.db.Exec(createTableQuery); err != nil {
		return fmt.Errorf("failed to create table: %w", err)
	}
	return nil
}

func (ps *PostgresSubjectStore) GetHours(subject string) (int, error) {
	var hours int
	err := ps.db.QueryRow(selectHoursQuery, subject).Scan(&hours)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, domain.ErrSubjectNotFound
		}
		return 0, fmt.Errorf("failed to make DB query for subject %s: %w", subject, err)
	}
	return hours, nil
}

func (ps *PostgresSubjectStore) RecordHour(subject string, numHours int) error {
	if _, err := ps.db.Exec(insertHoursQuery, subject, numHours); err != nil {
		return fmt.Errorf("failed to insert %s: %w", subject, err)
	}
	return nil
}

func (ps *PostgresSubjectStore) GetReport() (domain.Report, error) {
	rows, err := ps.db.Query(selectReportQuery)
	if err != nil {
		return nil, fmt.Errorf("failed to make query from subjects: %w", err)
	}
	defer rows.Close()

	report := make(domain.Report, 0)
	for rows.Next() {
		var sa domain.StudyActivity
		err = rows.Scan(&sa.Subject, &sa.Hours)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		report = append(report, sa)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate rows: %w", err)
	}

	return report, nil
}
