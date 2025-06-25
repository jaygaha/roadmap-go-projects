package service

import (
	"context"

	"github.com/jaygaha/roadmap-go-projects/intermediate/url-shortening-service/internal/models"
)

// URLService defines the interface for URL shortening service
type URLService interface {
	CreateShortURL(ctx context.Context, longURL string) (*models.URLResponse, error)
	GetAllURLs(ctx context.Context, page, limit int) ([]*models.URLResponse, int, error)
	GetOriginalURL(ctx context.Context, shortCode string) (string, error)
	GetOriginalURLStats(ctx context.Context, shortCode string) (*models.URLResponse, error)
	UpdateURL(ctx context.Context, shortCode string, longURL string) (*models.URLResponse, error)
	DeleteURL(ctx context.Context, shortCode string) error
}
