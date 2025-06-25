package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// URL represents a shortened URL entry.
type URL struct {
	Id         primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Url        string             `bson:"url" json:"url"`
	ShortCode  string             `bson:"short_code" json:"short_code"`
	ClickCount int64              `bson:"click_count" json:"click_count"`
	CreatedAt  time.Time          `bson:"created_at" json:"created_at"`
}

// CreateURLRequest represents the request structure for creating a new URL.
type CreateURLRequest struct {
	URL string `json:"url" binding:"required"`
}

// URLResponse represents the response structure for a shortened URL.
type URLResponse struct {
	Id         string    `json:"id"`
	Url        string    `json:"url"`
	ShortCode  string    `json:"short_code"`
	ShortURL   string    `json:"short_url"`
	ClickCount int64     `json:"click_count"`
	CreatedAt  time.Time `json:"created_at"`
}
