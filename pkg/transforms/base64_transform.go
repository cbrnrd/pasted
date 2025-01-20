package transforms

import (
	"bytes"
	"encoding/base64"
	"io"
)

type Base64Transformer struct{}

func NewBase64Transformer() *Base64Transformer {
	return &Base64Transformer{}
}

func (t *Base64Transformer) Transform(input io.Reader) (io.Reader, error) {
	// Read the input data into memory
	plaintext, err := io.ReadAll(input)
	if err != nil {
		return nil, err
	}

	// Encode the data as base64
	encoded := make([]byte, base64.StdEncoding.EncodedLen(len(plaintext)))
	base64.StdEncoding.Encode(encoded, plaintext)

	return bytes.NewReader(encoded), nil
}

func (t *Base64Transformer) ReverseTransform(input io.Reader) (io.Reader, error) {
	// Read the input data into memory
	encoded, err := io.ReadAll(input)
	if err != nil {
		return nil, err
	}

	// Decode the data from base64
	decoded := make([]byte, base64.StdEncoding.DecodedLen(len(encoded)))
	n, err := base64.StdEncoding.Decode(decoded, encoded)
	if err != nil {
		return nil, err
	}

	return bytes.NewReader(decoded[:n]), nil
}
