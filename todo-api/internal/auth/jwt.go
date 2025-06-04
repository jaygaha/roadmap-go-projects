package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// SecretKey is the secret key used to sign JWT tokens.
// IN PRODUCTION, THIS SHOULD BE STORED IN AN ENVIRONMENT VARIABLE.
var SecretKey = []byte("74dd5f08d68cbcc3d4b85af09a1c48866e7b02fd051f840079595e6002c12a82")

// Claims is the struct that represents the claims in a JWT token.
type Claims struct {
	UserID   int64  `json:"user_id"`
	UserName string `json:"user_name"`
	jwt.RegisteredClaims
}

// GenerateToken generates a JWT token for the given user.
func GenerateToken(userID int64, userName string) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)

	claims := &Claims{
		UserID:   userID,
		UserName: userName,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(SecretKey)
}

// ValidateToken validates the given JWT token.
func ValidateToken(tokenString string) (*Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (any, error) {
		return SecretKey, nil
	})

	if err != nil || !token.Valid {
		return nil, err
	}

	return claims, nil
}
