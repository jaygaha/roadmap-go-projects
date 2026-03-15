package model

import "time"

// Genre represents a movie genre
type Genre struct {
	ID   int64  `json:"id" db:"id"`
	Name string `json:"name" db:"name"`
}

// Movie represents a movie in the system
type Movie struct {
	ID            int64     `json:"id" db:"id"`
	Title         string    `json:"title" db:"title"`
	Description   string    `json:"description" db:"description"`
	PosterURL     string    `json:"poster_url" db:"poster_url"`
	GenreID       int64     `json:"genre_id" db:"genre_id"`
	Genre         *Genre    `json:"genre,omitempty" db:"-"`
	CreatedAt     time.Time `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time `json:"updated_at" db:"updated_at"`
}

// MovieCreateRequest represents the request body for creating a movie
type MovieCreateRequest struct {
	Title       string `json:"title" validate:"required,min=1,max=255"`
	Description string `json:"description" validate:"required"`
	PosterURL   string `json:"poster_url" validate:"url"`
	GenreID     int64  `json:"genre_id" validate:"required"`
}

// MovieUpdateRequest represents the request body for updating a movie
type MovieUpdateRequest struct {
	Title       string `json:"title" validate:"min=1,max=255"`
	Description string `json:"description"`
	PosterURL   string `json:"poster_url" validate:"url"`
	GenreID     *int64 `json:"genre_id"`
}

// MovieListResponse represents the response for movie list
type MovieListResponse struct {
	Movies []Movie `json:"movies"`
	Total  int64   `json:"total"`
}

// MovieDetailResponse represents the response for a single movie
type MovieDetailResponse struct {
	Movie *Movie `json:"movie"`
}
