package service

import (
	"context"
	"fmt"
	"time"

	"github.com/jaygaha/roadmap-go-projects/advanced/movie-reservation-system/internal/model"
	"github.com/jaygaha/roadmap-go-projects/advanced/movie-reservation-system/internal/repository"
	pkg "github.com/jaygaha/roadmap-go-projects/advanced/movie-reservation-system/pkg"
)

// ShowtimeService handles showtime business logic
type ShowtimeService struct {
	showtimeRepo repository.ShowtimeRepository
	movieRepo    repository.MovieRepository
	seatRepo     repository.SeatRepository
}

// NewShowtimeService creates a new showtime service
func NewShowtimeService(showtimeRepo repository.ShowtimeRepository, movieRepo repository.MovieRepository, seatRepo repository.SeatRepository) *ShowtimeService {
	return &ShowtimeService{
		showtimeRepo: showtimeRepo,
		movieRepo:    movieRepo,
		seatRepo:     seatRepo,
	}
}

// GetByID retrieves a showtime by ID
func (s *ShowtimeService) GetByID(ctx context.Context, id int64) (*model.Showtime, error) {
	return s.showtimeRepo.GetByID(ctx, id)
}

// GetAll retrieves paginated showtimes
func (s *ShowtimeService) GetAll(ctx context.Context, page, limit int) ([]model.Showtime, int64, error) {
	return s.showtimeRepo.GetAll(ctx, page, limit)
}

// GetByMovieID retrieves showtimes for a specific movie
func (s *ShowtimeService) GetByMovieID(ctx context.Context, movieID int64) ([]model.Showtime, error) {
	// Check if movie exists
	_, err := s.movieRepo.GetByID(ctx, movieID)
	if err != nil {
		if err == pkg.ErrRecordNotFound {
			return nil, pkg.ErrMovieNotFound
		}
		return nil, pkg.ErrDatabase
	}
	return s.showtimeRepo.GetByMovieID(ctx, movieID)
}

// GetByDate retrieves showtimes for a specific date
func (s *ShowtimeService) GetByDate(ctx context.Context, date string) ([]model.Showtime, error) {
	return s.showtimeRepo.GetByDate(ctx, date)
}

// Create creates a new showtime
func (s *ShowtimeService) Create(ctx context.Context, req *model.ShowtimeCreateRequest) (*model.Showtime, error) {
	// Check if movie exists
	_, err := s.movieRepo.GetByID(ctx, req.MovieID)
	if err != nil {
		if err == pkg.ErrRecordNotFound {
			return nil, pkg.ErrMovieNotFound
		}
		return nil, pkg.ErrDatabase
	}

	// Parse times
	startTime, err := time.Parse(time.RFC3339, req.StartTime)
	if err != nil {
		return nil, pkg.ErrValidation
	}

	endTime, err := time.Parse(time.RFC3339, req.EndTime)
	if err != nil {
		return nil, pkg.ErrValidation
	}

	showtime := &model.Showtime{
		MovieID:         req.MovieID,
		StartTime:       startTime,
		EndTime:         endTime,
		TheaterCapacity: req.TheaterCapacity,
		AvailableSeats:  req.TheaterCapacity,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	if err := s.showtimeRepo.Create(ctx, showtime); err != nil {
		return nil, pkg.ErrDatabase
	}

	// Create seats
	seatNumbers := make([]string, req.TheaterCapacity)
	for i := 0; i < req.TheaterCapacity; i++ {
		seatNumbers[i] = np(i+1, 3, "0")
	}

	if err := s.seatRepo.CreateBatch(ctx, showtime.ID, seatNumbers); err != nil {
		return nil, pkg.ErrDatabase
	}

	return showtime, nil
}

// Update updates a showtime
func (s *ShowtimeService) Update(ctx context.Context, showtime *model.Showtime) error {
	// Check if movie exists
	if showtime.MovieID != 0 {
		_, err := s.movieRepo.GetByID(ctx, showtime.MovieID)
		if err != nil {
			return pkg.ErrMovieNotFound
		}
	}
	return s.showtimeRepo.Update(ctx, showtime)
}

// Delete deletes a showtime
func (s *ShowtimeService) Delete(ctx context.Context, id int64) error {
	return s.showtimeRepo.Delete(ctx, id)
}

// np formats a seat number as a zero-padded string of the given width
func np(n int, place int, _ string) string {
	return fmt.Sprintf("%0*d", place, n)
}
