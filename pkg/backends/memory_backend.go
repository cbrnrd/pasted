package backends

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
)

// MemoryBackend is a backend that stores files in memory
type MemoryBackend struct {
	// mapping is a map from keys to file contents
	mapping map[string][]byte
}

var _ Backend = (*MemoryBackend)(nil)

func NewMemoryBackend() *MemoryBackend {
	return &MemoryBackend{mapping: make(map[string][]byte)}
}

// Put stores the contents of r in memory and returns the key
// The key is the first 4 bytes of the SHA256 hash of the contents
func (m *MemoryBackend) Put(r io.Reader) (string, error) {
	fmt.Println("writing contents")
	contents, err := io.ReadAll(r)
	if err != nil {
		return "", err
	}

	hash := sha256.Sum256(contents)
	path := hex.EncodeToString(hash[:4])
	m.mapping[path] = contents
	fmt.Println(path)
	return path, nil

}

// Get writes the contents of the file at key to w
func (m *MemoryBackend) Get(key string, w io.Writer) error {
	_, err := w.Write(m.mapping[key])
	return err
}
