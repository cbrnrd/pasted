package transforms

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"io"
)

// AESTransformer encrypts and decrypts data using AES-GCM.
type AESTransformer struct {
	key []byte
}

// NewAESTransformer creates a new AESTransformer with the given key.
func NewAESTransformer(key []byte) (*AESTransformer, error) {
	if len(key) != 16 && len(key) != 24 && len(key) != 32 {
		return nil, fmt.Errorf("invalid key length: must be 16, 24, or 32 bytes")
	}
	return &AESTransformer{key: key}, nil
}

// Transform encrypts the input data using AES-GCM.
func (t *AESTransformer) Transform(input io.Reader) (io.Reader, error) {
	// Read the input data into memory
	plaintext, err := io.ReadAll(input)
	if err != nil {
		return nil, err
	}

	// Generate a new AES cipher block
	block, err := aes.NewCipher(t.key)
	if err != nil {
		return nil, err
	}

	// Create a GCM cipher mode
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	// Generate a random nonce
	nonce := make([]byte, aesGCM.NonceSize())
	if _, err := rand.Read(nonce); err != nil {
		return nil, err
	}

	// Encrypt the data
	ciphertext := aesGCM.Seal(nonce, nonce, plaintext, nil)

	return bytes.NewReader(ciphertext), nil
}

// ReverseTransform decrypts the input data using AES-GCM.
func (t *AESTransformer) ReverseTransform(input io.Reader) (io.Reader, error) {
	// Read the input data into memory
	ciphertext, err := io.ReadAll(input)
	if err != nil {
		return nil, err
	}

	// Generate a new AES cipher block
	block, err := aes.NewCipher(t.key)
	if err != nil {
		return nil, err
	}

	// Create a GCM cipher mode
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	// Ensure the input is at least the size of the nonce
	nonceSize := aesGCM.NonceSize()
	if len(ciphertext) < nonceSize {
		return nil, fmt.Errorf("ciphertext too short")
	}

	// Extract the nonce and the actual ciphertext
	nonce, encryptedData := ciphertext[:nonceSize], ciphertext[nonceSize:]

	// Decrypt the data
	plaintext, err := aesGCM.Open(nil, nonce, encryptedData, nil)
	if err != nil {
		return nil, err
	}

	return bytes.NewReader(plaintext), nil
}

// // NewAesTransform creates a new AesTransform with the given key
// func NewAesTransform(key string) *AesTransform {
// 	hash := sha256.Sum256([]byte(key))
// 	return &AesTransform{Key: hash[:]}
// }

// // Transform encrypts data using AES
// func (a *AesTransform) Transform(r io.Reader, w io.Writer) error {
// 	fmt.Println("applying AES transform")

// 	var err error
// 	a.writer, err = newEncryptWriter(w, a.Key)
// 	if err != nil {
// 		return err
// 	}

// 	_, err = io.Copy(a.writer, r)
// 	return err

// }

// // ReverseTransform decrypts data using AES
// func (a *AesTransform) ReverseTransform(r io.Reader, w io.Writer) error {
// 	fmt.Println("applying reverse AES transform")
// 	var err error
// 	a.reader, err = newEncryptReader(r, a.Key)
// 	if err != nil {
// 		return err
// 	}

// 	_, err = io.Copy(w, a.reader)
// 	return err
// }

// func newEncryptReader(r io.Reader, key []byte) (*cipher.StreamReader, error) {
// 	// read initial value
// 	iv := make([]byte, aes.BlockSize)
// 	n, err := r.Read(iv)
// 	if err != nil || n != len(iv) {
// 		return nil, errors.New("could not read initial value")
// 	}

// 	block, err := newBlock(key)
// 	if err != nil {
// 		return nil, err
// 	}

// 	stream := cipher.NewOFB(block, iv)
// 	return &cipher.StreamReader{S: stream, R: r}, nil
// }

// func newEncryptWriter(w io.Writer, key []byte) (*cipher.StreamWriter, error) {
// 	// generate random initial value
// 	iv := make([]byte, aes.BlockSize)
// 	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
// 		return nil, err
// 	}

// 	// write clear IV to allow for decryption
// 	n, err := w.Write(iv)
// 	if err != nil || n != len(iv) {
// 		return nil, errors.New("could not write initial value")
// 	}

// 	block, err := newBlock(key)
// 	if err != nil {
// 		return nil, err
// 	}

// 	stream := cipher.NewOFB(block, iv)
// 	return &cipher.StreamWriter{S: stream, W: w}, nil
// }

// func newBlock(key []byte) (cipher.Block, error) {
// 	hash := sha256.Sum256(key)
// 	block, err := aes.NewCipher(hash[:])
// 	if err != nil {
// 		return nil, err
// 	}
// 	return block, nil
// }
