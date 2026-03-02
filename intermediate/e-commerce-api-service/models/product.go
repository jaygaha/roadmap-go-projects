package models

import "time"

// Product represents a product in the system
type Product struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Price       int64     `json:"price"` // cents
	Stock       int       `json:"stock"`
	ImageURL    string    `json:"image_url"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// ProductCreateRequest represents the request payload for creating a product
type ProductCreateRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Price       int64  `json:"price"` // cents
	Stock       int    `json:"stock"`
	ImageURL    string `json:"image_url"`
}

// ProductUpdateRequest represents the request payload for updating a product
type ProductUpdateRequest struct {
	Name        *string `json:"name"`
	Description *string `json:"description"`
	Price       *int64  `json:"price"` // cents
	Stock       *int    `json:"stock"`
	ImageURL    *string `json:"image_url"`
}

// ProductQuery represents the query parameters for product retrieval
type ProductQuery struct {
	Name     string
	MinPrice int64
	MaxPrice int64
	Page     int
	Limit    int
}
