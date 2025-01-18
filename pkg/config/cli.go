package config

type CLIConfig struct {
	// Backend is the backend to use for storing files
	Backend string `yaml:"backend"`

	FileBackendRoot string `yaml:"file_backend_root"`

	SizeLimitBytes int64 `yaml:"size_limit_bytes"`

	// ListenAddr is the address to listen on for incoming connections
	ListenAddr string `yaml:"listen_addr"`

	// Domain is the domain to use for generating URLs
	Domain string `yaml:"domain"`

	// TLS is whether or not to use TLS
	TLS TLSConfig `yaml:"tls"`

	Transformers []string `yaml:"transformers"`

	AESTransform struct {
		Key string `yaml:"key"`
	} `yaml:"aes_transform"`
}

type TLSConfig struct {
	// CertFile is the path to the certificate file
	CertFile string `yaml:"cert_file"`

	// KeyFile is the path to the key file
	KeyFile string `yaml:"key_file"`
}
