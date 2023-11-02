package common

import (
	"crypto/rand"
	"encoding/hex"
)

// GenerateSecretKey generates a random secret key of specified length.
func GenerateSecretKey(length int) (string, error) {
	bytes := make([]byte, length)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(bytes), nil
}
