package database

import (
	"context"
	"database/sql"
	"fmt"

	_ "modernc.org/sqlite"
)

type Database struct {
	*sql.DB
}

func New(path string) (*Database, error) {
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, fmt.Errorf("opening database: %w", err)
	}

	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("pinging database: %w", err)
	}

	return &Database{db}, nil
}

func (db *Database) Init(ctx context.Context) error {
	query := `
	CREATE TABLE IF NOT EXISTS listeners (
		protocol TEXXT NOT NULL,
		config TEXT NOT NULL,
		name TEXT UNIQUE NOT NULL
	);

	CREATE TABLE IF NOT EXISTS implants (
		id STRING UNIQUE NOT NULL,
		meta TEXT NOT NULL,
		channels TEXT NOT NULL
	);`
	_, err := db.ExecContext(ctx, query)
	return err
}

func (db *Database) Close() error {
	return db.DB.Close()
}
