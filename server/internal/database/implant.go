package database

import (
	"bytes"
	"context"
	"encoding/gob"
)

type Implant struct {
	ID       string
	Meta     map[string]any
	Channels []Channel
}

type Channel struct {
	ID       uint32
	Name     string
	Fallback bool
	InUse    bool
}

func (db *Database) ImplantCreate(id string, channels []map[string]any, meta any) error {

	channelsSerialized, err := encode(channels)
	if err != nil {
		return err
	}

	metaSerialized, err := encode(meta)
	if err != nil {
		return err
	}

	query := `INSERT INTO implants (id, meta, channels) VALUES (?, ?, ?)`
	_, err = db.ExecContext(context.Background(), query, id, metaSerialized, channelsSerialized)
	if err != nil {
		return err
	}
	return nil
}

func (db *Database) ImplantGetAll() ([]Implant, error) {
	query := `SELECT id, meta, channels FROM implants`
	rows, err := db.QueryContext(context.Background(), query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var (
		implants []Implant
	)
	for rows.Next() {
		var (
			id                 string
			metaSerialized     []byte
			channelsSerialized []byte
		)

		err := rows.Scan(&id, &metaSerialized, &channelsSerialized)
		if err != nil {
			return nil, err
		}

		meta, err := decode[map[string]any](metaSerialized)
		if err != nil {
			continue
		}

		channels, err := decode[[]map[string]any](channelsSerialized)
		if err != nil {
			continue
		}

		implant := Implant{
			ID:   id,
			Meta: meta,
		}

		for _, ch := range channels {
			implant.Channels = append(implant.Channels, Channel{
				Name:     ch["name"].(string),
				Fallback: ch["fallback"].(bool),
				InUse:    ch["in-use"].(bool),
			})
		}

		implants = append(implants, implant)
	}

	return implants, rows.Err()
}

func encode(data any) ([]byte, error) {
	var (
		encoded bytes.Buffer
	)

	err := gob.NewEncoder(&encoded).Encode(data)
	if err != nil {
		return nil, err
	}

	return encoded.Bytes(), nil
}

func decode[T any](data []byte) (T, error) {
	var (
		decoded T
	)

	buf := bytes.NewBuffer(data)
	err := gob.NewDecoder(buf).Decode(&decoded)
	if err != nil {
		return decoded, err
	}

	return decoded, nil
}
