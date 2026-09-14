package database

import (
	"context"
)

type Listener struct {
	Name     string
	Protocol string
	Config   string
}

func (db *Database) ListenerCreate(listener Listener) error {
	query := `INSERT INTO listeners (protocol, name, config) VALUES (?, ?, ?)`
	_, err := db.ExecContext(context.Background(), query, listener.Protocol, listener.Name, listener.Config)
	if err != nil {
		return err
	}
	return nil
}

func (db *Database) ListenerGetAll() ([]Listener, error) {
	query := `SELECT protocol, name, config FROM listeners`
	rows, err := db.QueryContext(context.Background(), query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var listeners []Listener
	for rows.Next() {
		var listener Listener
		err := rows.Scan(&listener.Protocol, &listener.Name, &listener.Config)
		if err != nil {
			return nil, err
		}
		listeners = append(listeners, listener)
	}
	return listeners, rows.Err()
}
