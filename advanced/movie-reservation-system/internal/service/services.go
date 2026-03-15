package service

import (
	"github.com/jaygaha/roadmap-go-projects/advanced/movie-reservation-system/internal/auth"
	"github.com/jaygaha/roadmap-go-projects/advanced/movie-reservation-system/internal/repository"
)

// Services holds all service dependencies
type Services struct {
	UserService        *UserService
	MovieService       *MovieService
	ShowtimeService    *ShowtimeService
	ReservationService *ReservationService
}

// NewServices creates and wires up all services
func NewServices(db *repository.DB, auth *auth.AuthService) *Services {
	// Initialize repositories
	userRepo := repository.NewUserRepository(db)
	movieRepo := repository.NewMovieRepository(db)
	genreRepo := repository.NewGenreRepository(db)
	showtimeRepo := repository.NewShowtimeRepository(db)
	reservationRepo := repository.NewReservationRepository(db)
	seatRepo := repository.NewSeatRepository(db)

	// Initialize services
	userService := NewUserService(userRepo, auth)
	movieService := NewMovieService(movieRepo, genreRepo)
	showtimeService := NewShowtimeService(showtimeRepo, movieRepo, seatRepo)
	reservationService := NewReservationService(reservationRepo, showtimeRepo, userRepo, seatRepo)

	return &Services{
		UserService:        userService,
		MovieService:       movieService,
		ShowtimeService:    showtimeService,
		ReservationService: reservationService,
	}
}
