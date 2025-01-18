package config

import (
	"fmt"

	"github.com/cbrnrd/pasted/pkg/backends"
	"github.com/cbrnrd/pasted/pkg/util"
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
	default:
		return nil, fmt.Errorf("unknown backend %s", config.Backend)
	}
}
