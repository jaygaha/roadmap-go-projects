package service

import (
	"context"
	"time"

	"github.com/jaygaha/roadmap-go-projects/advanced/movie-reservation-system/internal/model"
	"github.com/jaygaha/roadmap-go-projects/advanced/movie-reservation-system/internal/repository"
	pkg "github.com/jaygaha/roadmap-go-projects/advanced/movie-reservation-system/pkg"
)

// MovieService handles movie business logic
type MovieService struct {
	movieRepo repository.MovieRepository
	genreRepo repository.GenreRepository
}

// NewMovieService creates a new movie service
func NewMovieService(movieRepo repository.MovieRepository, genreRepo repository.GenreRepository) *MovieService {
	return &MovieService{
		movieRepo: movieRepo,
		genreRepo: genreRepo,
	}
}

// GetByID retrieves a movie by ID
func (s *MovieService) GetByID(ctx context.Context, id int64) (*model.Movie, error) {
	return s.movieRepo.GetByID(ctx, id)
}

// GetAll retrieves paginated movies
func (s *MovieService) GetAll(ctx context.Context, page, limit int) ([]model.Movie, int64, error) {
	return s.movieRepo.GetAll(ctx, page, limit)
}

// Create creates a new movie
func (s *MovieService) Create(ctx context.Context, req *model.MovieCreateRequest) (*model.Movie, error) {
	// Validate genre exists
	_, err := s.genreRepo.GetByID(ctx, req.GenreID)
	if err != nil {
		return nil, pkg.ErrGenreNotFound
	}

	movie := &model.Movie{
		Title:       req.Title,
		Description: req.Description,
		PosterURL:   req.PosterURL,
		GenreID:     req.GenreID,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := s.movieRepo.Create(ctx, movie); err != nil {
		return nil, pkg.ErrDatabase
	}

	return movie, nil
}

// Update updates a movie
func (s *MovieService) Update(ctx context.Context, movie *model.Movie) error {
	if movie.GenreID != 0 {
		_, err := s.genreRepo.GetByID(ctx, movie.GenreID)
		if err != nil {
			return pkg.ErrGenreNotFound
		}
	}
	return s.movieRepo.Update(ctx, movie)
}

// Delete deletes a movie
func (s *MovieService) Delete(ctx context.Context, id int64) error {
	return s.movieRepo.Delete(ctx, id)
}

// FindByGenre retrieves movies by genre ID
func (s *MovieService) FindByGenre(ctx context.Context, genreID int64, page, limit int) ([]model.Movie, int64, error) {
	return s.movieRepo.FindByGenre(ctx, genreID, page, limit)
}

// GetAllGenres retrieves all genres
func (s *MovieService) GetAllGenres(ctx context.Context) ([]model.Genre, error) {
	return s.genreRepo.GetAll(ctx)
}
