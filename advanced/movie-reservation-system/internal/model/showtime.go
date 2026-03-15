package model

import "time"

// Showtime represents a scheduled movie showing
type Showtime struct {
	ID              int64     `json:"id" db:"id"`
	MovieID         int64     `json:"movie_id" db:"movie_id"`
	Movie           *Movie    `json:"movie,omitempty" db:"-"`
	StartTime       time.Time `json:"start_time" db:"start_time"`
	EndTime         time.Time `json:"end_time" db:"end_time"`
	TheaterCapacity int       `json:"theater_capacity" db:"theater_capacity"`
	AvailableSeats  int       `json:"available_seats" db:"available_seats"`
	CreatedAt       time.Time `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time `json:"updated_at" db:"updated_at"`
}

// ShowtimeCreateRequest represents the request body for creating a showtime
type ShowtimeCreateRequest struct {
	MovieID         int64  `json:"movie_id" validate:"required"`
	StartTime       string `json:"start_time" validate:"required,datetime=2006-01-02T15:04:05Z07:00"`
	EndTime         string `json:"end_time" validate:"required,datetime=2006-01-02T15:04:05Z07:00"`
	TheaterCapacity int    `json:"theater_capacity" validate:"required,gt=0"`
}

// ShowtimeListResponse represents the response for showtime list
type ShowtimeListResponse struct {
	Showtimes []Showtime `json:"showtimes"`
	Total     int64      `json:"total"`
}

// ShowtimeDetailResponse represents the response for a single showtime
type ShowtimeDetailResponse struct {
	Showtime *Showtime `json:"showtime"`
}

// ShowtimeByDateRequest represents the request for showtimes by date
type ShowtimeByDateRequest struct {
	Date string `json:"date" validate:"required,date=2006-01-02"`
}

// ShowtimeByDateResponse represents the response for showtimes by date
type ShowtimeByDateResponse struct {
	Date      string     `json:"date"`
	Showtimes []Showtime `json:"showtimes"`
}
