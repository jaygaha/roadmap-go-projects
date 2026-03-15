package repository

import (
	"context"
	"github.com/jaygaha/roadmap-go-projects/advanced/movie-reservation-system/internal/model"
)

// UserRepository defines the interface for user database operations
type UserRepository interface {
	GetByID(ctx context.Context, id int64) (*model.User, error)
	GetByEmail(ctx context.Context, email string) (*model.User, error)
	Create(ctx context.Context, user *model.User) error
	Update(ctx context.Context, user *model.User) error
	Delete(ctx context.Context, id int64) error
	List(ctx context.Context, page, limit int) ([]model.User, int64, error)
	ListAdmins(ctx context.Context) ([]model.User, error)
}

// MovieRepository defines the interface for movie database operations
type MovieRepository interface {
	GetByID(ctx context.Context, id int64) (*model.Movie, error)
	GetAll(ctx context.Context, page, limit int) ([]model.Movie, int64, error)
	Create(ctx context.Context, movie *model.Movie) error
	Update(ctx context.Context, movie *model.Movie) error
	Delete(ctx context.Context, id int64) error
	FindByGenre(ctx context.Context, genreID int64, page, limit int) ([]model.Movie, int64, error)
}

// GenreRepository defines the interface for genre database operations
type GenreRepository interface {
	GetByID(ctx context.Context, id int64) (*model.Genre, error)
	GetAll(ctx context.Context) ([]model.Genre, error)
	Create(ctx context.Context, genre *model.Genre) error
	Update(ctx context.Context, genre *model.Genre) error
	Delete(ctx context.Context, id int64) error
}

// ShowtimeRepository defines the interface for showtime database operations
type ShowtimeRepository interface {
	GetByID(ctx context.Context, id int64) (*model.Showtime, error)
	GetAll(ctx context.Context, page, limit int) ([]model.Showtime, int64, error)
	GetByMovieID(ctx context.Context, movieID int64) ([]model.Showtime, error)
	GetByDate(ctx context.Context, date string) ([]model.Showtime, error)
	Create(ctx context.Context, showtime *model.Showtime) error
	Update(ctx context.Context, showtime *model.Showtime) error
	Delete(ctx context.Context, id int64) error
	DecrementAvailableSeats(ctx context.Context, showtimeID int64) error
	IncrementAvailableSeats(ctx context.Context, showtimeID int64) error
}

// ReservationRepository defines the interface for reservation database operations
type ReservationRepository interface {
	GetByID(ctx context.Context, id int64) (*model.Reservation, error)
	GetAllByUserID(ctx context.Context, userID int64, page, limit int) ([]model.Reservation, int64, error)
	GetAll(ctx context.Context, page, limit int) ([]model.Reservation, int64, error)
	GetByShowtimeID(ctx context.Context, showtimeID int64) ([]model.Reservation, error)
	Create(ctx context.Context, reservation *model.Reservation) error
	Update(ctx context.Context, reservation *model.Reservation) error
	Delete(ctx context.Context, id int64) error
	Cancel(ctx context.Context, id int64) error
	GetByUserIDAndShowtimeID(ctx context.Context, userID, showtimeID int64) (*model.Reservation, error)
	GetStats(ctx context.Context) (*model.ReservationStats, error)
}

// SeatRepository defines the interface for seat database operations
type SeatRepository interface {
	GetAllByShowtimeID(ctx context.Context, showtimeID int64) ([]model.Seat, error)
	GetAvailableByShowtimeID(ctx context.Context, showtimeID int64) ([]model.Seat, error)
	GetByShowtimeIDAndNumber(ctx context.Context, showtimeID int64, seatNumber string) (*model.Seat, error)
	CreateBatch(ctx context.Context, showtimeID int64, seatNumbers []string) error
	UpdateAvailability(ctx context.Context, seatID int64, available bool) error
	ReserveSeat(ctx context.Context, showtimeID int64, seatNumber string) error
	ReleaseSeat(ctx context.Context, showtimeID int64, seatNumber string) error
}
