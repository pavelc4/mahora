package db

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "modernc.org/sqlite"

	dbgen "github.com/pavelc4/mahora/internal/db/gen"
)

type DB struct {
	*dbgen.Queries
	conn *sql.DB
}

func New(ctx context.Context, dsn string) (*DB, error) {
	conn, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, fmt.Errorf("db.New open: %w", err)
	}

	conn.SetMaxOpenConns(1)
	if _, err = conn.ExecContext(ctx, `PRAGMA journal_mode=WAL; PRAGMA foreign_keys=ON;`); err != nil {
		return nil, fmt.Errorf("db.New pragma: %w", err)
	}

	if err = runMigrations(conn); err != nil {
		return nil, fmt.Errorf("db.New migrate: %w", err)
	}

	return &DB{
		Queries: dbgen.New(conn),
		conn:    conn,
	}, nil
}

func (d *DB) Close() error {
	if err := d.conn.Close(); err != nil {
		return fmt.Errorf("db.Close: %w", err)
	}
	return nil
}

func runMigrations(conn *sql.DB) error {
	driver, err := sqlite.WithInstance(conn, &sqlite.Config{})
	if err != nil {
		return fmt.Errorf("runMigrations driver: %w", err)
	}

	m, err := migrate.NewWithDatabaseInstance("file://internal/db/migrations", "sqlite", driver)
	if err != nil {
		return fmt.Errorf("runMigrations init: %w", err)
	}

	if err = m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("runMigrations up: %w", err)
	}
	return nil
}
