package backends

import (
	"errors"
	"io"
)

type Backend interface {
	Put(r io.Reader) (string, error)
	Get(key string, w io.Writer) error
}

type PathGenFunc func() string

var (
	ErrFileTooLarge = errors.New("file too large")
)
