package backends

import (
	"context"
	"fmt"
	"io"

	"github.com/jackc/pgx/v5/pgxpool"
)

// PgxBackend is a backend that stores pastes in a PostgreSQL database.
// It uses pgx as the driver.
type PgxBackend struct {
	ctx         context.Context
	pool        *pgxpool.Pool
	pathGenFunc PathGenFunc
}

var _ Backend = (*PgxBackend)(nil)

// NewPostgresBackend creates a new PostgresBackend.
// If createTables is true, the necessary tables will be created if they do not exist.
func NewPostgresBackend(ctx context.Context, connString string, pgf PathGenFunc, createTables bool) (*PgxBackend, error) {
	pool, err := pgxpool.New(ctx, connString)
	if err != nil {
		return nil, err
	}

	if createTables {
		_, err = pool.Exec(ctx, `CREATE TABLE IF NOT EXISTS pastes (
			id TEXT PRIMARY KEY,
			data BYTEA
		)`)
		if err != nil {
			return nil, err
		}
	}

	return &PgxBackend{ctx: ctx, pool: pool, pathGenFunc: pgf}, nil
}

func (b *PgxBackend) Put(r io.Reader) (string, error) {
	data, err := io.ReadAll(r)
	if err != nil {
		return "", err
	}

	key := b.pathGenFunc()
	fmt.Println(key)
	fmt.Println(data)
	// spew.Dump(b)

	_, err = b.pool.Exec(b.ctx, "INSERT INTO pastes (id, data) VALUES ($1, $2) RETURNING id", key, data)
	if err != nil {
		return "", err
	}

	return key, nil
}

func (b *PgxBackend) Get(key string, w io.Writer) error {
	var data []byte
	err := b.pool.QueryRow(b.ctx, "SELECT data FROM pastes WHERE id=$1", key).Scan(&data)
	if err != nil {
		return err
	}

	_, err = w.Write(data)
	return err
}
