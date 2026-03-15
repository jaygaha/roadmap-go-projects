package repository

import (
	"context"
	"database/sql"
	"github.com/jaygaha/roadmap-go-projects/advanced/movie-reservation-system/internal/model"
	"github.com/jaygaha/roadmap-go-projects/advanced/movie-reservation-system/pkg"
)

type seatRepository struct {
	db *DB
}

// NewSeatRepository creates a new seat repository
func NewSeatRepository(db *DB) SeatRepository {
	return &seatRepository{db: db}
}

func (r *seatRepository) GetAllByShowtimeID(ctx context.Context, showtimeID int64) ([]model.Seat, error) {
	query := `SELECT id, showtime_id, seat_number, is_available, created_at, updated_at
			  FROM seats WHERE showtime_id = $1 ORDER BY seat_number`

	rows, err := r.db.QueryContext(ctx, query, showtimeID)
	if err != nil {
		return nil, pkg.ErrDatabase
	}
	defer rows.Close()

	var seats []model.Seat
	for rows.Next() {
		var s model.Seat
		err := rows.Scan(
			&s.ID,
			&s.ShowtimeID,
			&s.SeatNumber,
			&s.IsAvailable,
			&s.CreatedAt,
			&s.UpdatedAt,
		)
		if err != nil {
			return nil, pkg.ErrDatabase
		}
		seats = append(seats, s)
	}

	return seats, nil
}

func (r *seatRepository) GetAvailableByShowtimeID(ctx context.Context, showtimeID int64) ([]model.Seat, error) {
	query := `SELECT id, showtime_id, seat_number, is_available, created_at, updated_at
			  FROM seats WHERE showtime_id = $1 AND is_available = true ORDER BY seat_number`

	rows, err := r.db.QueryContext(ctx, query, showtimeID)
	if err != nil {
		return nil, pkg.ErrDatabase
	}
	defer rows.Close()

	var seats []model.Seat
	for rows.Next() {
		var s model.Seat
		err := rows.Scan(
			&s.ID,
			&s.ShowtimeID,
			&s.SeatNumber,
			&s.IsAvailable,
			&s.CreatedAt,
			&s.UpdatedAt,
		)
		if err != nil {
			return nil, pkg.ErrDatabase
		}
		seats = append(seats, s)
	}

	return seats, nil
}

func (r *seatRepository) GetByShowtimeIDAndNumber(ctx context.Context, showtimeID int64, seatNumber string) (*model.Seat, error) {
	query := `SELECT id, showtime_id, seat_number, is_available, created_at, updated_at
			  FROM seats WHERE showtime_id = $1 AND seat_number = $2`

	row := r.db.QueryRowContext(ctx, query, showtimeID, seatNumber)

	seat := &model.Seat{}
	err := row.Scan(
		&seat.ID,
		&seat.ShowtimeID,
		&seat.SeatNumber,
		&seat.IsAvailable,
		&seat.CreatedAt,
		&seat.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, pkg.ErrRecordNotFound
		}
		return nil, pkg.ErrDatabase
	}

	return seat, nil
}

func (r *seatRepository) CreateBatch(ctx context.Context, showtimeID int64, seatNumbers []string) error {
	// Use transaction for batch insert
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return pkg.ErrDatabase
	}
	defer tx.Rollback()

	query := `INSERT INTO seats (showtime_id, seat_number, is_available)
			  VALUES ($1, $2, true)`

	stmt, err := tx.PrepareContext(ctx, query)
	if err != nil {
		return pkg.ErrDatabase
	}
	defer stmt.Close()

	for _, seatNumber := range seatNumbers {
		_, err := stmt.ExecContext(ctx, showtimeID, seatNumber)
		if err != nil {
			return pkg.ErrDatabase
		}
	}

	return tx.Commit()
}

func (r *seatRepository) UpdateAvailability(ctx context.Context, seatID int64, available bool) error {
	query := `UPDATE seats SET is_available = $1, updated_at = CURRENT_TIMESTAMP WHERE id = $2`

	result, err := r.db.ExecContext(ctx, query, available, seatID)
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

func (r *seatRepository) ReserveSeat(ctx context.Context, showtimeID int64, seatNumber string) error {
	// Use transaction to ensure atomicity
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return pkg.ErrDatabase
	}
	defer tx.Rollback()

	// Check if seat is available
	var isAvailable bool
	err = tx.QueryRowContext(ctx, `SELECT is_available FROM seats WHERE showtime_id = $1 AND seat_number = $2`, showtimeID, seatNumber).Scan(&isAvailable)
	if err == sql.ErrNoRows {
		return pkg.ErrSeatNotAvailable
	}
	if err != nil {
		return pkg.ErrDatabase
	}

	if !isAvailable {
		return pkg.ErrSeatNotAvailable
	}

	// Update seat availability
	_, err = tx.ExecContext(ctx, `UPDATE seats SET is_available = false, updated_at = CURRENT_TIMESTAMP WHERE showtime_id = $1 AND seat_number = $2`, showtimeID, seatNumber)
	if err != nil {
		return pkg.ErrDatabase
	}

	return tx.Commit()
}

func (r *seatRepository) ReleaseSeat(ctx context.Context, showtimeID int64, seatNumber string) error {
	query := `UPDATE seats SET is_available = true, updated_at = CURRENT_TIMESTAMP WHERE showtime_id = $1 AND seat_number = $2`

	result, err := r.db.ExecContext(ctx, query, showtimeID, seatNumber)
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
