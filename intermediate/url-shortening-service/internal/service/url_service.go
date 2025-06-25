package service

import (
	"context"
	"fmt"
	"log"

	"github.com/jaygaha/roadmap-go-projects/intermediate/url-shortening-service/internal/models"
	"github.com/jaygaha/roadmap-go-projects/intermediate/url-shortening-service/internal/repository"
	"github.com/jaygaha/roadmap-go-projects/intermediate/url-shortening-service/internal/utils"
)

// urlService implements URLService interface
type urlService struct {
	repo    repository.URLRepository
	baseURL string
}

// NewURLService
func NewURLService(repo repository.URLRepository, baseURL string) URLService {
	return &urlService{
		repo:    repo,
		baseURL: baseURL,
	}
}

// CreateShortURL
func (s *urlService) CreateShortURL(ctx context.Context, Url string) (*models.URLResponse, error) {
	shortCode, err := utils.GenerateShortCode()
	if err != nil {
		return nil, fmt.Errorf("failed to generate short code: %w", err)
	}

	urlModel := &models.URL{
		Url:       Url,
		ShortCode: shortCode,
	}
	if err := s.repo.Create(ctx, urlModel); err != nil {
		return nil, fmt.Errorf("failed to create short URL: %w", err)
	}

	return s.buildURLResponse(urlModel, true), nil
}

// GetAllURLs returns all URLs
func (s *urlService) GetAllURLs(ctx context.Context, page, limit int) ([]*models.URLResponse, int, error) {
	urls, total, err := s.repo.GetAll(ctx, page, limit)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get all URLs: %w", err)
	}

	var urlResponses []*models.URLResponse
	for _, url := range urls {
		urlResponses = append(urlResponses, s.buildURLResponse(url, true))
	}

	return urlResponses, total, nil
}

// GetOriginalURL returns the original URL which needs to be redirected
func (s *urlService) GetOriginalURL(ctx context.Context, shortCode string) (string, error) {
	url, err := s.repo.GetByShortCode(ctx, shortCode)
	if err != nil {
		return "", fmt.Errorf("failed to get original URL: %w", err)
	}

	// Increment click count asynchronously
	go func() {
		if err := s.repo.IncrementClickCount(context.Background(), shortCode); err != nil {
			log.Printf("Failed to increment click count: %v", err)
		}
	}()

	return url.Url, nil
}

// GetOriginalURLStats returns the original URL with stats
func (s *urlService) GetOriginalURLStats(ctx context.Context, shortCode string) (*models.URLResponse, error) {
	url, err := s.repo.GetByShortCode(ctx, shortCode)
	if err != nil {
		return nil, fmt.Errorf("failed to get original URL: %w", err)
	}

	return s.buildURLResponse(url, false), nil
}

// UpdateURL updates a URL
func (s *urlService) UpdateURL(ctx context.Context, shortCode string, longURL string) (*models.URLResponse, error) {
	url, err := s.repo.GetByShortCode(ctx, shortCode)
	if err != nil {
		return nil, fmt.Errorf("failed to get URL: %w", err)
	}

	url.Url = longURL
	if err := s.repo.UpdateByShortCode(ctx, shortCode, url); err != nil {
		return nil, fmt.Errorf("failed to update URL: %w", err)
	}

	return s.buildURLResponse(url, true), nil
}

// DeleteURL deletes a URL
func (s *urlService) DeleteURL(ctx context.Context, shortCode string) error {
	_, err := s.repo.GetByShortCode(ctx, shortCode)
	if err != nil {
		return fmt.Errorf("failed to get URL: %w", err)
	}

	if err := s.repo.DeleteByShortCode(ctx, shortCode); err != nil {
		return fmt.Errorf("failed to delete URL: %w", err)
	}

	return nil
}

// buildURLResponse builds a URLResponse from a URL
func (s *urlService) buildURLResponse(url *models.URL, isHideStats bool) *models.URLResponse {
	res := &models.URLResponse{
		Id:        url.Id.Hex(),
		Url:       url.Url,
		ShortCode: url.ShortCode,
		ShortURL:  fmt.Sprintf("%s/s/%s", s.baseURL, url.ShortCode), // e.g. http://localhost:8080/s/123456
		CreatedAt: url.CreatedAt,
	}

	if !isHideStats {
		res.ClickCount = url.ClickCount
	}

	return res
}
