package transforms

import (
	"bytes"
	"compress/gzip"
	"io"
)

type GZipTransformer struct{}

func (t *GZipTransformer) Transform(input io.Reader) (io.Reader, error) {
	var buf bytes.Buffer
	gz := gzip.NewWriter(&buf)
	_, err := io.Copy(gz, input)
	if err != nil {
		return nil, err
	}
	if err := gz.Close(); err != nil {
		return nil, err
	}
	return &buf, nil
}

func (t *GZipTransformer) ReverseTransform(input io.Reader) (io.Reader, error) {
	gz, err := gzip.NewReader(input)
	if err != nil {
		return nil, err
	}
	var buf bytes.Buffer
	_, err = io.Copy(&buf, gz)
	if err != nil {
		return nil, err
	}
	return &buf, nil
}
