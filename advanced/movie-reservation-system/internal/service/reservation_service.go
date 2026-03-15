package service

import (
	"context"
	"time"

	"github.com/jaygaha/roadmap-go-projects/advanced/movie-reservation-system/internal/model"
	"github.com/jaygaha/roadmap-go-projects/advanced/movie-reservation-system/internal/repository"
	pkg "github.com/jaygaha/roadmap-go-projects/advanced/movie-reservation-system/pkg"
)

// ReservationService handles reservation business logic
type ReservationService struct {
	reservationRepo repository.ReservationRepository
	showtimeRepo    repository.ShowtimeRepository
	userRepo        repository.UserRepository
	seatRepo        repository.SeatRepository
}

// NewReservationService creates a new reservation service
func NewReservationService(reservationRepo repository.ReservationRepository, showtimeRepo repository.ShowtimeRepository, userRepo repository.UserRepository, seatRepo repository.SeatRepository) *ReservationService {
	return &ReservationService{
		reservationRepo: reservationRepo,
		showtimeRepo:    showtimeRepo,
		userRepo:        userRepo,
		seatRepo:        seatRepo,
	}
}

// GetByID retrieves a reservation by ID
func (s *ReservationService) GetByID(ctx context.Context, id int64) (*model.Reservation, error) {
	return s.reservationRepo.GetByID(ctx, id)
}

// GetAllByUserID retrieves reservations for a specific user
func (s *ReservationService) GetAllByUserID(ctx context.Context, userID int64, page, limit int) ([]model.Reservation, int64, error) {
	return s.reservationRepo.GetAllByUserID(ctx, userID, page, limit)
}

// GetAll retrieves all reservations (admin)
func (s *ReservationService) GetAll(ctx context.Context, page, limit int) ([]model.Reservation, int64, error) {
	return s.reservationRepo.GetAll(ctx, page, limit)
}

// GetStats retrieves reservation statistics (admin)
func (s *ReservationService) GetStats(ctx context.Context) (*model.ReservationStats, error) {
	return s.reservationRepo.GetStats(ctx)
}

// Create creates a new reservation
func (s *ReservationService) Create(ctx context.Context, req *model.CreateReservationRequest, userID int64) (*model.Reservation, error) {
	// Check if showtime exists
	_, err := s.showtimeRepo.GetByID(ctx, req.ShowtimeID)
	if err != nil {
		if err == pkg.ErrRecordNotFound {
			return nil, pkg.ErrShowtimeNotFound
		}
		return nil, pkg.ErrDatabase
	}

	// Check if seat is available
	seat, err := s.seatRepo.GetByShowtimeIDAndNumber(ctx, req.ShowtimeID, req.SeatNumber)
	if err != nil {
		if err == pkg.ErrRecordNotFound {
			return nil, pkg.ErrSeatNotAvailable
		}
		return nil, pkg.ErrDatabase
	}

	if !seat.IsAvailable {
		return nil, pkg.ErrSeatNotAvailable
	}

	// Check if user already has a reservation for this seat
	existing, err := s.reservationRepo.GetByUserIDAndShowtimeID(ctx, userID, req.ShowtimeID)
	if err == nil && existing.SeatNumber == req.SeatNumber {
		return nil, pkg.ErrAlreadyReserved
	}

	// Create reservation
	reservation := &model.Reservation{
		UserID:     userID,
		ShowtimeID: req.ShowtimeID,
		SeatNumber: req.SeatNumber,
		Status:     model.StatusActive,
	}

	if err := s.reservationRepo.Create(ctx, reservation); err != nil {
		return nil, pkg.ErrDatabase
	}

	// Mark seat as unavailable
	if err := s.seatRepo.UpdateAvailability(ctx, seat.ID, false); err != nil {
		return nil, pkg.ErrDatabase
	}

	// Decrement available seats in showtime
	if err := s.showtimeRepo.DecrementAvailableSeats(ctx, req.ShowtimeID); err != nil {
		return nil, pkg.ErrDatabase
	}

	return reservation, nil
}

// Cancel cancels a reservation
func (s *ReservationService) Cancel(ctx context.Context, id int64, userID int64) error {
	// Get reservation
	reservation, err := s.reservationRepo.GetByID(ctx, id)
	if err != nil {
		if err == pkg.ErrRecordNotFound {
			return pkg.ErrReservationNotFound
		}
		return pkg.ErrDatabase
	}

	// Check ownership
	if reservation.UserID != userID {
		return pkg.ErrForbidden
	}

	// Check if showtime is in the past
	now := contextTime(ctx)
	if now.After(reservation.Showtime.StartTime) {
		return pkg.ErrReservationPast
	}

	// Mark reservation as cancelled
	if err := s.reservationRepo.Cancel(ctx, id); err != nil {
		return err
	}

	// Mark seat as available
	seat, err := s.seatRepo.GetByShowtimeIDAndNumber(ctx, reservation.ShowtimeID, reservation.SeatNumber)
	if err != nil {
		return err
	}
	if err := s.seatRepo.UpdateAvailability(ctx, seat.ID, true); err != nil {
		return err
	}

	// Increment available seats in showtime
	if err := s.showtimeRepo.IncrementAvailableSeats(ctx, reservation.ShowtimeID); err != nil {
		return err
	}

	return nil
}

// contextTime gets the current time from context or returns now()
func contextTime(_ context.Context) time.Time {
	// In real implementation, you could use context deadline
	return time.Now()
}
