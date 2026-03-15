package repository

import (
	"context"
	"database/sql"

	"github.com/jaygaha/roadmap-go-projects/advanced/movie-reservation-system/internal/model"
	"github.com/jaygaha/roadmap-go-projects/advanced/movie-reservation-system/pkg"
)

type reservationRepository struct {
	db *DB
}

// NewReservationRepository creates a new reservation repository
func NewReservationRepository(db *DB) ReservationRepository {
	return &reservationRepository{db: db}
}

func (r *reservationRepository) GetByID(ctx context.Context, id int64) (*model.Reservation, error) {
	query := `SELECT res.id, res.user_id, res.showtime_id, res.seat_number, res.status, res.created_at, res.updated_at,
					 u.email, u.name, u.role,
					 s.start_time, s.end_time, s.theater_capacity, s.available_seats,
					 m.title, m.description, m.poster_url, m.genre_id
			  FROM reservations res
			  LEFT JOIN users u ON res.user_id = u.id
			  LEFT JOIN showtimes s ON res.showtime_id = s.id
			  LEFT JOIN movies m ON s.movie_id = m.id
			  WHERE res.id = $1`

	row := r.db.QueryRowContext(ctx, query, id)

	reservation := &model.Reservation{}
	user := &model.User{}
	showtime := &model.Showtime{}
	movie := &model.Movie{}
	err := row.Scan(
		&reservation.ID,
		&reservation.UserID,
		&reservation.ShowtimeID,
		&reservation.SeatNumber,
		&reservation.Status,
		&reservation.CreatedAt,
		&reservation.UpdatedAt,
		&user.Email,
		&user.Name,
		&user.Role,
		&showtime.StartTime,
		&showtime.EndTime,
		&showtime.TheaterCapacity,
		&showtime.AvailableSeats,
		&movie.Title,
		&movie.Description,
		&movie.PosterURL,
		&movie.GenreID,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, pkg.ErrRecordNotFound
		}
		return nil, pkg.ErrDatabase
	}

	reservation.User = user
	reservation.Showtime = showtime
	showtime.Movie = movie

	return reservation, nil
}

func (r *reservationRepository) GetAllByUserID(ctx context.Context, userID int64, page, limit int) ([]model.Reservation, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	offset := (page - 1) * limit

	// Get total count
	var total int64
	err := r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM reservations WHERE user_id = $1`, userID).Scan(&total)
	if err != nil {
		return nil, 0, pkg.ErrDatabase
	}

	// Get reservations
	query := `SELECT res.id, res.user_id, res.showtime_id, res.seat_number, res.status, res.created_at, res.updated_at,
					 s.start_time, s.end_time, s.theater_capacity, s.available_seats,
					 m.title
			  FROM reservations res
			  LEFT JOIN showtimes s ON res.showtime_id = s.id
			  LEFT JOIN movies m ON s.movie_id = m.id
			  WHERE res.user_id = $1 ORDER BY res.created_at DESC LIMIT $2 OFFSET $3`

	rows, err := r.db.QueryContext(ctx, query, userID, limit, offset)
	if err != nil {
		return nil, 0, pkg.ErrDatabase
	}
	defer rows.Close()

	var reservations []model.Reservation
	for rows.Next() {
		var res model.Reservation
		showtime := &model.Showtime{}
		movie := &model.Movie{}
		err := rows.Scan(
			&res.ID,
			&res.UserID,
			&res.ShowtimeID,
			&res.SeatNumber,
			&res.Status,
			&res.CreatedAt,
			&res.UpdatedAt,
			&showtime.StartTime,
			&showtime.EndTime,
			&showtime.TheaterCapacity,
			&showtime.AvailableSeats,
			&movie.Title,
		)
		if err != nil {
			return nil, 0, pkg.ErrDatabase
		}
		res.Showtime = showtime
		showtime.Movie = movie
		reservations = append(reservations, res)
	}

	return reservations, total, nil
}

func (r *reservationRepository) GetAll(ctx context.Context, page, limit int) ([]model.Reservation, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	offset := (page - 1) * limit

	// Get total count
	var total int64
	err := r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM reservations`).Scan(&total)
	if err != nil {
		return nil, 0, pkg.ErrDatabase
	}

	// Get reservations
	query := `SELECT res.id, res.user_id, res.showtime_id, res.seat_number, res.status, res.created_at, res.updated_at,
					 u.email, u.name,
					 s.start_time, s.end_time, s.theater_capacity, s.available_seats,
					 m.title
			  FROM reservations res
			  LEFT JOIN users u ON res.user_id = u.id
			  LEFT JOIN showtimes s ON res.showtime_id = s.id
			  LEFT JOIN movies m ON s.movie_id = m.id
			  ORDER BY res.created_at DESC LIMIT $1 OFFSET $2`

	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, 0, pkg.ErrDatabase
	}
	defer rows.Close()

	var reservations []model.Reservation
	for rows.Next() {
		var res model.Reservation
		user := &model.User{}
		showtime := &model.Showtime{}
		movie := &model.Movie{}
		err := rows.Scan(
			&res.ID,
			&res.UserID,
			&res.ShowtimeID,
			&res.SeatNumber,
			&res.Status,
			&res.CreatedAt,
			&res.UpdatedAt,
			&user.Email,
			&user.Name,
			&showtime.StartTime,
			&showtime.EndTime,
			&showtime.TheaterCapacity,
			&showtime.AvailableSeats,
			&movie.Title,
		)
		if err != nil {
			return nil, 0, pkg.ErrDatabase
		}
		res.User = user
		res.Showtime = showtime
		showtime.Movie = movie
		reservations = append(reservations, res)
	}

	return reservations, total, nil
}

func (r *reservationRepository) GetByShowtimeID(ctx context.Context, showtimeID int64) ([]model.Reservation, error) {
	query := `SELECT id, user_id, showtime_id, seat_number, status, created_at, updated_at
			  FROM reservations WHERE showtime_id = $1 ORDER BY seat_number`

	rows, err := r.db.QueryContext(ctx, query, showtimeID)
	if err != nil {
		return nil, pkg.ErrDatabase
	}
	defer rows.Close()

	var reservations []model.Reservation
	for rows.Next() {
		var res model.Reservation
		err := rows.Scan(
			&res.ID,
			&res.UserID,
			&res.ShowtimeID,
			&res.SeatNumber,
			&res.Status,
			&res.CreatedAt,
			&res.UpdatedAt,
		)
		if err != nil {
			return nil, pkg.ErrDatabase
		}
		reservations = append(reservations, res)
	}

	return reservations, nil
}

func (r *reservationRepository) Create(ctx context.Context, reservation *model.Reservation) error {
	query := `INSERT INTO reservations (user_id, showtime_id, seat_number, status)
			  VALUES ($1, $2, $3, $4) RETURNING id`

	err := r.db.QueryRowContext(ctx, query,
		reservation.UserID,
		reservation.ShowtimeID,
		reservation.SeatNumber,
		reservation.Status,
	).Scan(&reservation.ID)

	if err != nil {
		return pkg.ErrDatabase
	}

	return nil
}

func (r *reservationRepository) Update(ctx context.Context, reservation *model.Reservation) error {
	query := `UPDATE reservations SET showtime_id = $1, seat_number = $2, status = $3, updated_at = CURRENT_TIMESTAMP
			  WHERE id = $4`

	result, err := r.db.ExecContext(ctx, query,
		reservation.ShowtimeID,
		reservation.SeatNumber,
		reservation.Status,
		reservation.ID,
	)

	if err != nil {
		return pkg.ErrDatabase
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return pkg.ErrDatabase
	}

	if rowsAffected == 0 {
		return pkg.ErrRecordNotFound
	}

	return nil
}

func (r *reservationRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM reservations WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)

	if err != nil {
		return pkg.ErrDatabase
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return pkg.ErrDatabase
	}

	if rowsAffected == 0 {
		return pkg.ErrRecordNotFound
	}

	return nil
}

func (r *reservationRepository) Cancel(ctx context.Context, id int64) error {
	const StatusCancelled = "cancelled"

	// First get the reservation to check its status and get the showtime
	var status string
	err := r.db.QueryRowContext(ctx, `SELECT status FROM reservations WHERE id = $1`, id).Scan(&status)
	if err != nil {
		return pkg.ErrRecordNotFound
	}

	if status == StatusCancelled {
		return nil // Already cancelled
	}

	query := `UPDATE reservations SET status = $1, updated_at = CURRENT_TIMESTAMP WHERE id = $2`

	result, err := r.db.ExecContext(ctx, query, StatusCancelled, id)
	if err != nil {
		return pkg.ErrDatabase
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return pkg.ErrDatabase
	}

	if rowsAffected == 0 {
		return pkg.ErrRecordNotFound
	}

	return nil
}

func (r *reservationRepository) GetByUserIDAndShowtimeID(ctx context.Context, userID, showtimeID int64) (*model.Reservation, error) {
	query := `SELECT id, user_id, showtime_id, seat_number, status, created_at, updated_at
			  FROM reservations WHERE user_id = $1 AND showtime_id = $2`

	row := r.db.QueryRowContext(ctx, query, userID, showtimeID)

	reservation := &model.Reservation{}
	err := row.Scan(
		&reservation.ID,
		&reservation.UserID,
		&reservation.ShowtimeID,
		&reservation.SeatNumber,
		&reservation.Status,
		&reservation.CreatedAt,
		&reservation.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, pkg.ErrRecordNotFound
		}
		return nil, pkg.ErrDatabase
	}

	return reservation, nil
}

func (r *reservationRepository) GetStats(ctx context.Context) (*model.ReservationStats, error) {
	const (
		StatusActive    = "active"
		StatusCancelled = "cancelled"
	)
	stats := &model.ReservationStats{}

	// Total reservations
	err := r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM reservations`).Scan(&stats.TotalReservations)
	if err != nil {
		return nil, pkg.ErrDatabase
	}

	// Active reservations
	err = r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM reservations WHERE status = $1`, StatusActive).Scan(&stats.ActiveReservations)
	if err != nil {
		return nil, pkg.ErrDatabase
	}

	// Cancelled reservations
	err = r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM reservations WHERE status = $1`, StatusCancelled).Scan(&stats.CancelledReservations)
	if err != nil {
		return nil, pkg.ErrDatabase
	}

	// Total capacity used
	var totalCapacity int64
	err = r.db.QueryRowContext(ctx, `SELECT COALESCE(SUM(s.theater_capacity), 0) FROM reservations res
			LEFT JOIN showtimes s ON res.showtime_id = s.id
			WHERE res.status = $1`, StatusActive).Scan(&totalCapacity)
	if err != nil {
		return nil, pkg.ErrDatabase
	}

	if stats.TotalReservations > 0 {
		stats.AverageCapacity = float64(totalCapacity) / float64(stats.TotalReservations)
	}

	return stats, nil
}
