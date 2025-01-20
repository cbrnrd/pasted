package config

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/cbrnrd/pasted/pkg/backends"
	"github.com/cbrnrd/pasted/pkg/util"
	"github.com/go-redis/redis/v8"
	_ "github.com/mattn/go-sqlite3"
)

// GetBackend creates a new backend based on the provided configuration.
func (config *CLIConfig) GetBackend() (backends.Backend, error) {
	switch config.Backend {
	case "memory":
		return backends.NewMemoryBackend(), nil
	case "file":
		return backends.NewFileBackend(config.FileBackendRoot, config.SizeLimitBytes, func() string {
			path, err := util.GenerateRandomString(5)
			if err != nil {
				return ""
			}
			return path
		})
	case "sqlite":
		db, err := sql.Open("sqlite3", config.SQLiteConfig.Path)
		if err != nil {
			return nil, err
		}

		return backends.NewSQLiteBackend(db, func() string {
			path, err := util.GenerateRandomString(5)
			if err != nil {
				return ""
			}
			return path
		}, config.SQLiteConfig.CreateTables)
	case "pgx", "postgres":
		return backends.NewPostgresBackend(context.Background(), config.PgxConfig.ConnString, func() string {
			path, err := util.GenerateRandomString(5)
			if err != nil {
				return ""
			}
			return path
		}, config.PgxConfig.CreateTables)
	case "redis":
		redisClient := redis.NewClient(&config.RedisConfig)
		status := redisClient.Ping(context.Background())
		if status.Err() != nil {
			return nil, status.Err()
		}

		return backends.NewRedisBackend(context.Background(), func() string {
			path, err := util.GenerateRandomString(5)
			if err != nil {
				return ""
			}
			return path
		}, redisClient), nil
	case "s3":
		s3Client := s3.NewFromConfig(config.S3Config.Config)
		return backends.NewS3Backend(context.Background(), func() string {
			path, err := util.GenerateRandomString(5)
			if err != nil {
				return ""
			}
			return path
		}, config.S3Config.Bucket, s3Client), nil

	default:
		return nil, fmt.Errorf("unknown backend %s", config.Backend)
	}
}
