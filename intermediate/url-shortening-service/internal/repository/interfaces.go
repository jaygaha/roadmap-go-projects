package repository

import (
	"context"

	"github.com/jaygaha/roadmap-go-projects/intermediate/url-shortening-service/internal/models"
)

type URLRepository interface {
	Create(ctx context.Context, url *models.URL) error
	GetAll(ctx context.Context, page, limit int) ([]*models.URL, int, error)
	GetByShortCode(ctx context.Context, shortCode string) (*models.URL, error)
	UpdateByShortCode(ctx context.Context, shortCode string, url *models.URL) error
	DeleteByShortCode(ctx context.Context, shortCode string) error
	IncrementClickCount(ctx context.Context, shortCode string) error
}
