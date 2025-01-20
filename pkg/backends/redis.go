package backends

import (
	"context"
	"io"

	"github.com/go-redis/redis/v8"
)

type RedisBackend struct {
	client *redis.Client
	ctx    context.Context
	pathGenFunc PathGenFunc
}

var _ Backend = (*RedisBackend)(nil)

func NewRedisBackend(ctx context.Context, pgf PathGenFunc, client *redis.Client) *RedisBackend {
	return &RedisBackend{client: client, ctx: ctx, pathGenFunc: pgf}
}

// Get writes the contents of the file at key to w.
// If the key does not exist, Get should return nil.
// 
// Note that the value is read into memory before being written to w.
func (b *RedisBackend) Get(key string, w io.Writer) error {
	val, err := b.client.Get(b.ctx, key).Result()
	if err == redis.Nil {
		return nil
	}

	if err != nil {
		return err
	}

	_, err = w.Write([]byte(val))

	return err
}

// Put stores the contents of r in memory and returns the key.
// Note that r will be read into memory before being stored.
func (b *RedisBackend) Put(r io.Reader) (string, error) {
	value, err := io.ReadAll(r)
	if err != nil {
		return "", err
	}
	path := b.pathGenFunc()
	status := b.client.Set(b.ctx, path, value, 0).Err()
	if status != nil {
		return "", status
	}
	return path, nil
}
