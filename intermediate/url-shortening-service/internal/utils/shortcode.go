package utils

import (
	"crypto/rand"
	"encoding/base64"
)

// GenerateShortCode generates a unique short code
func GenerateShortCode() (string, error) {
	bytes := make([]byte, 6)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}

	encoded := base64.URLEncoding.EncodeToString(bytes)
	return encoded[:8], nil
}
