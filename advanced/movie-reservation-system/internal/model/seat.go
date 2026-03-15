package model

import "time"

// Seat represents a seat in a showtime
type Seat struct {
	ID            int64     `json:"id" db:"id"`
	ShowtimeID    int64     `json:"showtime_id" db:"showtime_id"`
	SeatNumber    string    `json:"seat_number" db:"seat_number"`
	IsAvailable   bool      `json:"is_available" db:"is_available"`
	CreatedAt     time.Time `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time `json:"updated_at" db:"updated_at"`
}

// SeatListResponse represents the response for seat list
type SeatListResponse struct {
	Seats      []Seat `json:"seats"`
	Total      int    `json:"total"`
	Available  int    `json:"available"`
	Occupied   int    `json:"occupied"`
}

// SeatAvailabilityRequest represents the request for seat availability
type SeatAvailabilityRequest struct {
	ShowtimeID int64 `json:"showtime_id" validate:"required"`
}

// SeatAvailabilityResponse represents the response for seat availability
type SeatAvailabilityResponse struct {
	ShowtimeID      int64   `json:"showtime_id"`
	TotalSeats      int     `json:"total_seats"`
	AvailableSeats  int     `json:"available_seats"`
	OccupiedSeats   int     `json:"occupied_seats"`
	Seats           []Seat  `json:"seats"`
}
