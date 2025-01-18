package util

import (
	"crypto/rand"
	"math/big"
)

func GenerateRandomString(n int) (string, error) {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"
	bytes := make([]byte, n)
	for i := range bytes {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(letters))))
		if err != nil {
			return "", err
		}
		bytes[i] = letters[n.Int64()]
	}
	return string(bytes), nil
}