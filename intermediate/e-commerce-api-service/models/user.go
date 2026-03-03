package models

import "time"

// User represents a user in the system
type User struct {
	ID        int64      `json:"id"`
	Email     string     `json:"email"`
	Password  string     `json:"password"`
	Role      string     `json:"role"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
}

// User registration request payload
type UserRegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// User login request payload
type UserLoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// AuthResponse represents the response payload for authentication
type AuthResponse struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}
