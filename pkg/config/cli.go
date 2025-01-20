package config

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/go-redis/redis/v8"
)

type CLIConfig struct {
	// Backend is the backend to use for storing files
	Backend string `yaml:"backend"`

	FileBackendRoot string `yaml:"file_backend_root"`

	SizeLimitBytes int64 `yaml:"size_limit_bytes"`

	// ListenAddr is the address to listen on for incoming connections
	ListenAddr string `yaml:"listen_addr"`

	// HTTPListenAddr is the address to listen on for incoming HTTP connections
	HttpListenAddr string `yaml:"http_listen_addr"`

	// Domain is the domain to use for generating URLs
	Domain string `yaml:"domain"`

	// TLS is whether or not to use TLS
	TLS TLSConfig `yaml:"tls"`

	Transformers []string `yaml:"transformers"`

	AESTransform struct {
		Key string `yaml:"key"`
	} `yaml:"aes_transform"`

	PgxConfig struct {
		ConnString   string `yaml:"conn_string"`
		CreateTables bool   `yaml:"create_tables"`
	} `yaml:"postgres"`

	SQLiteConfig struct {
		Path         string `yaml:"path"`
		CreateTables bool   `yaml:"create_tables"`
	} `yaml:"sqlite"`

	S3Config struct {
		aws.Config `yaml:"config"`
		Bucket string `yaml:"bucket"` 
	} `yaml:"s3"`

	RedisConfig redis.Options `yaml:"redis"`
}

type TLSConfig struct {
	// CertFile is the path to the certificate file
	CertFile string `yaml:"cert_file"`

	// KeyFile is the path to the key file
	KeyFile string `yaml:"key_file"`
}
