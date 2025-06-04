package models

import "time"

// User represents a user in the system.
type User struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	Username  string    `json:"username"`
	Password  string    `json:"-"` // Password is hidden
	CreatedAt time.Time `json:"created_at"`
}
