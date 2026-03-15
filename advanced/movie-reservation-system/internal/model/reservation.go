package model

import "time"

// ReservationStatus represents the status of a reservation
type ReservationStatus string

const (
	StatusActive    ReservationStatus = "active"
	StatusCancelled ReservationStatus = "cancelled"
	StatusCompleted ReservationStatus = "completed"
)

// Reservation represents a seat reservation
type Reservation struct {
	ID            int64             `json:"id" db:"id"`
	UserID        int64             `json:"user_id" db:"user_id"`
	User          *User             `json:"user,omitempty" db:"-"`
	ShowtimeID    int64             `json:"showtime_id" db:"showtime_id"`
	Showtime      *Showtime         `json:"showtime,omitempty" db:"-"`
	SeatNumber    string            `json:"seat_number" db:"seat_number"`
	Status        ReservationStatus `json:"status" db:"status"`
	CreatedAt     time.Time         `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time         `json:"updated_at" db:"updated_at"`
}

// CreateReservationRequest represents the request body for creating a reservation
type CreateReservationRequest struct {
	ShowtimeID int64  `json:"showtime_id" validate:"required"`
	SeatNumber string `json:"seat_number" validate:"required,min=1,max=10"`
}

// ReservationListResponse represents the response for reservation list
type ReservationListResponse struct {
	Reservations []Reservation `json:"reservations"`
	Total        int64         `json:"total"`
}

// ReservationDetailResponse represents the response for a single reservation
type ReservationDetailResponse struct {
	Reservation *Reservation `json:"reservation"`
}

// AdminReservationListRequest represents the request for admin reservation list
type AdminReservationListRequest struct {
	Page      int                 `json:"page" validate:"min=1"`
	Limit     int                 `json:"limit" validate:"min=1,max=100"`
	Status    *ReservationStatus
	StartDate *string
	EndDate   *string
}

// ReservationStats represents reservation statistics for admin reports
type ReservationStats struct {
	TotalReservations     int64   `json:"total_reservations"`
	TotalRevenue          float64 `json:"total_revenue"`
	AverageCapacity       float64 `json:"average_capacity"`
	ActiveReservations    int64   `json:"active_reservations"`
	CancelledReservations int64   `json:"cancelled_reservations"`
}
