package models

import "time"

// Post denotes the blog post data structure
type Post struct {
	ID        int       `json:"id"` // primary key with auto-increment id
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	Category  string    `json:"category"`
	Tags      []string  `json:"tags"` // array
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
