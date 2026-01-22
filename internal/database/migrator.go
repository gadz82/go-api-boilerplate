package database

import (
	"database/sql"
	"embed"
	"errors"
	"fmt"

	"github.com/pressly/goose/v3"
)

//go:embed migrations/*.sql
var embedMigrations embed.FS

// Migrator handles database migrations using goose
type Migrator struct {
	db      *sql.DB
	dialect string
}

// NewMigrator creates a new Migrator instance with the specified dialect
// dialect should be "mysql" or "sqlite3"
func NewMigrator(db *sql.DB, dialect string) *Migrator {
	return &Migrator{db: db, dialect: dialect}
}

// Up runs all available migrations
func (m *Migrator) Up() error {
	goose.SetBaseFS(embedMigrations)

	if err := goose.SetDialect(m.dialect); err != nil {
		return fmt.Errorf("failed to set dialect %s: %w", m.dialect, err)
	}

	if err := goose.Up(m.db, "migrations"); err != nil {
		// ErrNoNextVersion means the database is already up to date - not an error
		if errors.Is(err, goose.ErrNoNextVersion) {
			return nil
		}
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	return nil
}

// Down rolls back the last migration
func (m *Migrator) Down() error {
	goose.SetBaseFS(embedMigrations)

	if err := goose.SetDialect(m.dialect); err != nil {
		return fmt.Errorf("failed to set dialect %s: %w", m.dialect, err)
	}

	if err := goose.Down(m.db, "migrations"); err != nil {
		return fmt.Errorf("failed to rollback migration: %w", err)
	}

	return nil
}

// Status prints the status of all migrations
func (m *Migrator) Status() error {
	goose.SetBaseFS(embedMigrations)

	if err := goose.SetDialect(m.dialect); err != nil {
		return fmt.Errorf("failed to set dialect %s: %w", m.dialect, err)
	}

	if err := goose.Status(m.db, "migrations"); err != nil {
		return fmt.Errorf("failed to get migration status: %w", err)
	}

	return nil
}
