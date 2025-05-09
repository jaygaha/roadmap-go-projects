package models

import "time"

type Blog struct {
	ID          int        `json:"id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   *time.Time `json:"updated_at"` // nullable
}

type Pagination struct {
	CurrentPage int
	TotalPages  int
	HasNext     bool
	HasPrev     bool
}
