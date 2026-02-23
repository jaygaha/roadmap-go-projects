package models

import "github.com/golang-jwt/jwt/v5"

// RegisterRequest represents the expected JSON payload
type RegisterRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// LoginRequest represents the login payload
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Claims defines the payload of our JWT
type Claims struct {
	UserID   int    `json:"user_id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}
