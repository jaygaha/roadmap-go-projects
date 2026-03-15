package auth

import (
	"golang.org/x/crypto/bcrypt"
)

// HashPassword hashes a password using bcrypt
func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

// CheckPassword checks if the provided password matches the hashed password
func CheckPassword(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}

// IsPasswordHash checks if the given string is a bcrypt hash
func IsPasswordHash(password string) bool {
	// bcrypt hashes start with $2a$, $2b$, or $2y$
	if len(password) < 4 {
		return false
	}
	return password[:4] == "$2a$" || password[:4] == "$2b$" || password[:4] == "$2y$"
}
