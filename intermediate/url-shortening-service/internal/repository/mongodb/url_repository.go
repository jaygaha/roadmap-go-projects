package mongodb

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jaygaha/roadmap-go-projects/intermediate/url-shortening-service/internal/models"
	"github.com/jaygaha/roadmap-go-projects/intermediate/url-shortening-service/internal/repository"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// urlRepository implements URLRepository interface
type urlRepository struct {
	collection *mongo.Collection
}

// NewURLRepository creates a new instance of urlRepository
func NewURLRepository(db *mongo.Database) repository.URLRepository {
	return &urlRepository{
		collection: db.Collection("urls"),
	}
}

// Create creates a new URL
func (r *urlRepository) Create(ctx context.Context, url *models.URL) error {
	url.Id = primitive.NewObjectID()
	url.ClickCount = 0
	url.CreatedAt = time.Now()

	_, err := r.collection.InsertOne(ctx, url)
	if err != nil {
		return fmt.Errorf("failed to insert URL: %w", err)
	}

	return nil
}

// GetAll returns all URLs
func (r *urlRepository) GetAll(ctx context.Context, page, limit int) ([]*models.URL, int, error) {
	var urls []*models.URL

	skip := (page - 1) * limit

	cursor, err := r.collection.Find(ctx, bson.M{}, options.Find().SetSkip(int64(skip)).SetLimit(int64(limit)))
	if err != nil {
		return nil, 0, fmt.Errorf("failed to find URLs: %w", err)
	}

	if err := cursor.All(ctx, &urls); err != nil {
		return nil, 0, fmt.Errorf("failed to decode URLs: %w", err)
	}

	// Count total documents
	total, err := r.collection.CountDocuments(ctx, bson.M{})
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count URLs: %w", err)
	}

	return urls, int(total), nil
}

// GetByShortCode retrieves a URL by its short code
func (r *urlRepository) GetByShortCode(ctx context.Context, shortCode string) (*models.URL, error) {
	var url models.URL
	filter := bson.M{"short_code": shortCode}

	err := r.collection.FindOne(ctx, filter).Decode(&url)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, errors.New("URL not found")
		}

		return nil, fmt.Errorf("failed to get URL by short code: %w", err)
	}
	return &url, nil
}

// UpdateByShortCode updates a URL by its short code
func (r *urlRepository) UpdateByShortCode(ctx context.Context, shortCode string, url *models.URL) error {
	filter := bson.M{"short_code": shortCode}
	update := bson.M{"$set": url}

	_, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("failed to update URL by short code: %w", err)
	}

	return nil
}

// DeleteByShortCode deletes a URL by its short code
func (r *urlRepository) DeleteByShortCode(ctx context.Context, shortCode string) error {
	filter := bson.M{"short_code": shortCode}

	_, err := r.collection.DeleteOne(ctx, filter)
	if err != nil {
		return fmt.Errorf("failed to delete URL by short code: %w", err)
	}

	return nil
}

// IncrementClickCount increments the click count of a URL by its short code
func (r *urlRepository) IncrementClickCount(ctx context.Context, shortCode string) error {
	filter := bson.M{"short_code": shortCode}
	update := bson.M{"$inc": bson.M{"click_count": 1}}

	_, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("failed to increment click count: %w", err)
	}

	return nil
}
