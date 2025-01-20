package backends

import (
	"database/sql"
	"fmt"
	"io"
)

type SQLiteBackend struct {

	// The underlying database connection.
	db *sql.DB

	// The function used to generate paths/keys.
	pathGenFunc PathGenFunc
}

var _ Backend = (*SQLiteBackend)(nil)

func NewSQLiteBackend(db *sql.DB, pgf PathGenFunc, createTables bool) (*SQLiteBackend, error) {

	if db == nil {
		return nil, fmt.Errorf("db is nil")
	}

	if createTables {
		_, err := db.Exec(`CREATE TABLE IF NOT EXISTS pastes (
			id TEXT PRIMARY KEY,
			data BLOB
		)`)
		if err != nil {
			return nil, err
		}
	}
	return &SQLiteBackend{db: db, pathGenFunc: pgf}, nil
}

func (b *SQLiteBackend) Put(r io.Reader) (string, error) {
	data, err := io.ReadAll(r)
	if err != nil {
		return "", err
	}

	key := b.pathGenFunc()
	_, err = b.db.Exec("INSERT INTO pastes (id, data) VALUES (?, ?)", key, data)
	if err != nil {
		return "", err
	}

	return key, nil
}

func (b *SQLiteBackend) Get(key string, w io.Writer) error {
	var data []byte
	err := b.db.QueryRow("SELECT data FROM pastes WHERE id=?", key).Scan(&data)
	if err != nil {
		return err
	}

	_, err = w.Write(data)
	return err
}