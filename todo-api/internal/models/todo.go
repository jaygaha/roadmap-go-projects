package models

import "time"

// Todo represents a user in the system.
type Todo struct {
	ID          int       `json:"id"`
	UserID      int64     `json:"userId"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	IsCompleted bool      `json:"is_completed"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
