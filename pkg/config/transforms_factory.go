package config

import (
	"crypto/sha256"
	"fmt"

	"github.com/cbrnrd/pasted/pkg/transforms"
)

func (config *CLIConfig) GetTransforms() ([]transforms.Transformer, error) {
	var t []transforms.Transformer
	for _, tc := range config.Transformers {
		switch tc {
		case "aes":
			hash := sha256.Sum256([]byte(config.AESTransform.Key))
			aesTransformer, err := transforms.NewAESTransformer(hash[:])
			if err != nil {
				return nil, err
			}
			t = append(t, aesTransformer)
		case "gzip":
			t = append(t, &transforms.GZipTransformer{})
		case "base64":
			t = append(t, &transforms.Base64Transformer{})
		default:
			return nil, fmt.Errorf("unknown transform %s", tc)
		}
	}
	return t, nil
}
