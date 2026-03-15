package repository

import (
	"context"
	"database/sql"
	"github.com/jaygaha/roadmap-go-projects/advanced/movie-reservation-system/internal/model"
	"github.com/jaygaha/roadmap-go-projects/advanced/movie-reservation-system/pkg"
)

type movieRepository struct {
	db *DB
}

// NewMovieRepository creates a new movie repository
func NewMovieRepository(db *DB) MovieRepository {
	return &movieRepository{db: db}
}

func (r *movieRepository) GetByID(ctx context.Context, id int64) (*model.Movie, error) {
	query := `SELECT m.id, m.title, m.description, m.poster_url, m.genre_id, m.created_at, m.updated_at, g.name as genre_name
			  FROM movies m LEFT JOIN genres g ON m.genre_id = g.id WHERE m.id = $1`

	row := r.db.QueryRowContext(ctx, query, id)

	movie := &model.Movie{}
	genre := &model.Genre{}
	err := row.Scan(
		&movie.ID,
		&movie.Title,
		&movie.Description,
		&movie.PosterURL,
		&movie.GenreID,
		&movie.CreatedAt,
		&movie.UpdatedAt,
		&genre.Name,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, pkg.ErrRecordNotFound
		}
		return nil, pkg.ErrDatabase
	}

	movie.Genre = genre
	return movie, nil
}

func (r *movieRepository) GetAll(ctx context.Context, page, limit int) ([]model.Movie, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	offset := (page - 1) * limit

	// Get total count
	var total int64
	err := r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM movies`).Scan(&total)
	if err != nil {
		return nil, 0, pkg.ErrDatabase
	}

	// Get movies
	query := `SELECT m.id, m.title, m.description, m.poster_url, m.genre_id, m.created_at, m.updated_at, g.name as genre_name
			  FROM movies m LEFT JOIN genres g ON m.genre_id = g.id ORDER BY m.id LIMIT $1 OFFSET $2`

	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, 0, pkg.ErrDatabase
	}
	defer rows.Close()

	var movies []model.Movie
	for rows.Next() {
		var m model.Movie
		genre := &model.Genre{}
		err := rows.Scan(
			&m.ID,
			&m.Title,
			&m.Description,
			&m.PosterURL,
			&m.GenreID,
			&m.CreatedAt,
			&m.UpdatedAt,
			&genre.Name,
		)
		if err != nil {
			return nil, 0, pkg.ErrDatabase
		}
		m.Genre = genre
		movies = append(movies, m)
	}

	return movies, total, nil
}

func (r *movieRepository) Create(ctx context.Context, movie *model.Movie) error {
	query := `INSERT INTO movies (title, description, poster_url, genre_id)
			  VALUES ($1, $2, $3, $4) RETURNING id`

	err := r.db.QueryRowContext(ctx, query,
		movie.Title,
		movie.Description,
		movie.PosterURL,
		movie.GenreID,
	).Scan(&movie.ID)

	if err != nil {
		return pkg.ErrDatabase
	}

	return nil
}

func (r *movieRepository) Update(ctx context.Context, movie *model.Movie) error {
	query := `UPDATE movies SET title = $1, description = $2, poster_url = $3, genre_id = $4, updated_at = CURRENT_TIMESTAMP
			  WHERE id = $5`

	result, err := r.db.ExecContext(ctx, query,
		movie.Title,
		movie.Description,
		movie.PosterURL,
		movie.GenreID,
		movie.ID,
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

func (r *movieRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM movies WHERE id = $1`

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

func (r *movieRepository) FindByGenre(ctx context.Context, genreID int64, page, limit int) ([]model.Movie, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	offset := (page - 1) * limit

	// Get total count
	var total int64
	err := r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM movies WHERE genre_id = $1`, genreID).Scan(&total)
	if err != nil {
		return nil, 0, pkg.ErrDatabase
	}

	// Get movies
	query := `SELECT m.id, m.title, m.description, m.poster_url, m.genre_id, m.created_at, m.updated_at, g.name as genre_name
			  FROM movies m LEFT JOIN genres g ON m.genre_id = g.id WHERE m.genre_id = $1 ORDER BY m.id LIMIT $2 OFFSET $3`

	rows, err := r.db.QueryContext(ctx, query, genreID, limit, offset)
	if err != nil {
		return nil, 0, pkg.ErrDatabase
	}
	defer rows.Close()

	var movies []model.Movie
	for rows.Next() {
		var m model.Movie
		genre := &model.Genre{}
		err := rows.Scan(
			&m.ID,
			&m.Title,
			&m.Description,
			&m.PosterURL,
			&m.GenreID,
			&m.CreatedAt,
			&m.UpdatedAt,
			&genre.Name,
		)
		if err != nil {
			return nil, 0, pkg.ErrDatabase
		}
		m.Genre = genre
		movies = append(movies, m)
	}

	return movies, total, nil
}
