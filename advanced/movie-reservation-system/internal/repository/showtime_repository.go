package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/jaygaha/roadmap-go-projects/advanced/movie-reservation-system/internal/model"
	"github.com/jaygaha/roadmap-go-projects/advanced/movie-reservation-system/pkg"
)

type showtimeRepository struct {
	db *DB
}

// NewShowtimeRepository creates a new showtime repository
func NewShowtimeRepository(db *DB) ShowtimeRepository {
	return &showtimeRepository{db: db}
}

func (r *showtimeRepository) GetByID(ctx context.Context, id int64) (*model.Showtime, error) {
	query := `SELECT s.id, s.movie_id, s.start_time, s.end_time, s.theater_capacity, s.available_seats, s.created_at, s.updated_at,
					 m.title, m.description, m.poster_url, m.genre_id, m.created_at as movie_created_at, m.updated_at as movie_updated_at
			  FROM showtimes s
			  LEFT JOIN movies m ON s.movie_id = m.id
			  WHERE s.id = $1`

	row := r.db.QueryRowContext(ctx, query, id)

	showtime := &model.Showtime{}
	movie := &model.Movie{}
	err := row.Scan(
		&showtime.ID,
		&showtime.MovieID,
		&showtime.StartTime,
		&showtime.EndTime,
		&showtime.TheaterCapacity,
		&showtime.AvailableSeats,
		&showtime.CreatedAt,
		&showtime.UpdatedAt,
		&movie.Title,
		&movie.Description,
		&movie.PosterURL,
		&movie.GenreID,
		&movie.CreatedAt,
		&movie.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, pkg.ErrRecordNotFound
		}
		return nil, pkg.ErrDatabase
	}

	showtime.Movie = movie
	return showtime, nil
}

func (r *showtimeRepository) GetAll(ctx context.Context, page, limit int) ([]model.Showtime, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	offset := (page - 1) * limit

	// Get total count
	var total int64
	err := r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM showtimes`).Scan(&total)
	if err != nil {
		return nil, 0, pkg.ErrDatabase
	}

	// Get showtimes
	query := `SELECT s.id, s.movie_id, s.start_time, s.end_time, s.theater_capacity, s.available_seats, s.created_at, s.updated_at
			  FROM showtimes s ORDER BY s.start_time LIMIT $1 OFFSET $2`

	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, 0, pkg.ErrDatabase
	}
	defer rows.Close()

	var showtimes []model.Showtime
	for rows.Next() {
		var s model.Showtime
		err := rows.Scan(
			&s.ID,
			&s.MovieID,
			&s.StartTime,
			&s.EndTime,
			&s.TheaterCapacity,
			&s.AvailableSeats,
			&s.CreatedAt,
			&s.UpdatedAt,
		)
		if err != nil {
			return nil, 0, pkg.ErrDatabase
		}
		showtimes = append(showtimes, s)
	}

	return showtimes, total, nil
}

func (r *showtimeRepository) GetByMovieID(ctx context.Context, movieID int64) ([]model.Showtime, error) {
	query := `SELECT id, movie_id, start_time, end_time, theater_capacity, available_seats, created_at, updated_at
			  FROM showtimes WHERE movie_id = $1 ORDER BY start_time`

	rows, err := r.db.QueryContext(ctx, query, movieID)
	if err != nil {
		return nil, pkg.ErrDatabase
	}
	defer rows.Close()

	var showtimes []model.Showtime
	for rows.Next() {
		var s model.Showtime
		err := rows.Scan(
			&s.ID,
			&s.MovieID,
			&s.StartTime,
			&s.EndTime,
			&s.TheaterCapacity,
			&s.AvailableSeats,
			&s.CreatedAt,
			&s.UpdatedAt,
		)
		if err != nil {
			return nil, pkg.ErrDatabase
		}
		showtimes = append(showtimes, s)
	}

	return showtimes, nil
}

func (r *showtimeRepository) GetByDate(ctx context.Context, date string) ([]model.Showtime, error) {
	// Parse date and get start and end of day
	startOfDay, err := time.Parse("2006-01-02", date)
	if err != nil {
		return nil, pkg.ErrInvalidDate
	}
	endOfDay := startOfDay.Add(24 * time.Hour)

	query := `SELECT id, movie_id, start_time, end_time, theater_capacity, available_seats, created_at, updated_at
			  FROM showtimes
			  WHERE start_time >= $1 AND start_time < $2
			  ORDER BY start_time`

	rows, err := r.db.QueryContext(ctx, query, startOfDay, endOfDay)
	if err != nil {
		return nil, pkg.ErrDatabase
	}
	defer rows.Close()

	var showtimes []model.Showtime
	for rows.Next() {
		var s model.Showtime
		err := rows.Scan(
			&s.ID,
			&s.MovieID,
			&s.StartTime,
			&s.EndTime,
			&s.TheaterCapacity,
			&s.AvailableSeats,
			&s.CreatedAt,
			&s.UpdatedAt,
		)
		if err != nil {
			return nil, pkg.ErrDatabase
		}
		showtimes = append(showtimes, s)
	}

	return showtimes, nil
}

func (r *showtimeRepository) Create(ctx context.Context, showtime *model.Showtime) error {
	query := `INSERT INTO showtimes (movie_id, start_time, end_time, theater_capacity, available_seats)
			  VALUES ($1, $2, $3, $4, $5) RETURNING id`

	err := r.db.QueryRowContext(ctx, query,
		showtime.MovieID,
		showtime.StartTime,
		showtime.EndTime,
		showtime.TheaterCapacity,
		showtime.AvailableSeats,
	).Scan(&showtime.ID)

	if err != nil {
		return pkg.ErrDatabase
	}

	return nil
}

func (r *showtimeRepository) Update(ctx context.Context, showtime *model.Showtime) error {
	query := `UPDATE showtimes SET movie_id = $1, start_time = $2, end_time = $3, theater_capacity = $4, available_seats = $5, updated_at = CURRENT_TIMESTAMP
			  WHERE id = $6`

	result, err := r.db.ExecContext(ctx, query,
		showtime.MovieID,
		showtime.StartTime,
		showtime.EndTime,
		showtime.TheaterCapacity,
		showtime.AvailableSeats,
		showtime.ID,
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

func (r *showtimeRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM showtimes WHERE id = $1`

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

func (r *showtimeRepository) DecrementAvailableSeats(ctx context.Context, showtimeID int64) error {
	query := `UPDATE showtimes SET available_seats = available_seats - 1 WHERE id = $1 AND available_seats > 0`

	result, err := r.db.ExecContext(ctx, query, showtimeID)
	if err != nil {
		return pkg.ErrDatabase
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return pkg.ErrDatabase
	}

	if rowsAffected == 0 {
		return pkg.ErrSeatNotAvailable
	}

	return nil
}

func (r *showtimeRepository) IncrementAvailableSeats(ctx context.Context, showtimeID int64) error {
	query := `UPDATE showtimes SET available_seats = available_seats + 1 WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, showtimeID)
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
