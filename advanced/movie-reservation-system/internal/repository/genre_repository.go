package repository

import (
	"context"
	"github.com/jaygaha/roadmap-go-projects/advanced/movie-reservation-system/internal/model"
	"github.com/jaygaha/roadmap-go-projects/advanced/movie-reservation-system/pkg"
)

type genreRepository struct {
	db *DB
}

// NewGenreRepository creates a new genre repository
func NewGenreRepository(db *DB) GenreRepository {
	return &genreRepository{db: db}
}

func (r *genreRepository) GetByID(ctx context.Context, id int64) (*model.Genre, error) {
	query := `SELECT id, name FROM genres WHERE id = $1`

	row := r.db.QueryRowContext(ctx, query, id)

	genre := &model.Genre{}
	err := row.Scan(&genre.ID, &genre.Name)

	if err != nil {
		if err == pkg.ErrRecordNotFound {
			return nil, pkg.ErrRecordNotFound
		}
		return nil, pkg.ErrDatabase
	}

	return genre, nil
}

func (r *genreRepository) GetAll(ctx context.Context) ([]model.Genre, error) {
	query := `SELECT id, name FROM genres ORDER BY name`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, pkg.ErrDatabase
	}
	defer rows.Close()

	var genres []model.Genre
	for rows.Next() {
		var g model.Genre
		err := rows.Scan(&g.ID, &g.Name)
		if err != nil {
			return nil, pkg.ErrDatabase
		}
		genres = append(genres, g)
	}

	return genres, nil
}

func (r *genreRepository) Create(ctx context.Context, genre *model.Genre) error {
	query := `INSERT INTO genres (name) VALUES ($1) RETURNING id`

	err := r.db.QueryRowContext(ctx, query, genre.Name).Scan(&genre.ID)

	if err != nil {
		return pkg.ErrDatabase
	}

	return nil
}

func (r *genreRepository) Update(ctx context.Context, genre *model.Genre) error {
	query := `UPDATE genres SET name = $1 WHERE id = $2`

	result, err := r.db.ExecContext(ctx, query, genre.Name, genre.ID)

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

func (r *genreRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM genres WHERE id = $1`

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
