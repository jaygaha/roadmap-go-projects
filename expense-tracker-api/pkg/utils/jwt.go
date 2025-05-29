package utils

import (
	"errors"
	"sync"
	"time"

	"github.com/golang-jwt/jwt"
)

// JWT blocklist (in-memory for demo, replace with Redis for prod)
var (
	// JWT token blocklist: token string → expiry time
	blocklist = make(map[string]time.Time)
	mu        sync.Mutex
)

// GenerateToken generates a JWT token with userID and expiry duration
func GenerateToken(userID uint, secret string, duration time.Duration) (string, error) {
	claims := jwt.MapClaims{
		"sub": userID,
		"exp": time.Now().Add(duration).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

// ParseToken parses a JWT token and returns its claims
func ParseToken(tokenStr, secret string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil || !token.Valid {
		return nil, errors.New("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("invalid claims format")
	}

	return claims, nil
}

// AddTokenToBlocklist adds a token string to blocklist with its expiry time
func AddTokenToBlocklist(token string, expiry time.Time) {
	mu.Lock()
	defer mu.Unlock()
	blocklist[token] = expiry
}

// IsTokenBlocked checks if a token is in the blocklist
func IsTokenBlocked(token string) bool {
	mu.Lock()
	defer mu.Unlock()

	expiry, found := blocklist[token]
	if !found {
		return false
	}

	if time.Now().After(expiry) {
		delete(blocklist, token)
		return false
	}

	return true
}
