package backends

import (
	"io"
	"os"
	"path/filepath"

	"github.com/cbrnrd/pasted/pkg/util"
)

type FileBackend struct {

	// root is the root directory where files are stored
	Root string `yaml:"root"`

	// maxSize is the maximum size of a submission
	MaxSize int64 `yaml:"max_size"`

	// pathGen is a function that generates a path for a file
	pathGen PathGenFunc `yaml:"-"`
}

var _ Backend = (*FileBackend)(nil)

func NewFileBackend(root string, maxSize int64, pathGen PathGenFunc) (*FileBackend, error) {
	if err := os.MkdirAll(root, os.ModePerm); err != nil {
		return nil, err
	}

	return &FileBackend{
		Root:    root,
		MaxSize: maxSize,
		pathGen: pathGen,
	}, nil
}

// Returns a file backend with sensible defaults
func DefaultFileBackend(root string) (*FileBackend, error) {
	// 50kB
	return NewFileBackend(root, 50*1024, func() string {
		r, _ := util.GenerateRandomString(6)
		return r
	})
}

// Put stores the contents of r in a file and returns the generated path to the file
func (f *FileBackend) Put(r io.Reader) (string, error) {
	path := f.pathGen()

	outFile, err := os.Create(filepath.Join(f.Root, path))
	if err != nil {
		return "", err
	}
	defer outFile.Close()

	n, err := io.CopyN(outFile, r, f.MaxSize)
	if err != io.EOF {
		return "", err
	}
	if n == f.MaxSize {
		return "", ErrFileTooLarge
	}

	return path, nil
}

// Get writes the contents of the file at key to w
func (f *FileBackend) Get(path string, w io.Writer) error {
	c := filepath.Clean(path)

	inFile, err := os.Open(filepath.Join(f.Root, c))
	if err != nil {
		return err
	}

	_, err = io.Copy(w, inFile)
	if err != nil {
		return err
	}

	return nil
}
