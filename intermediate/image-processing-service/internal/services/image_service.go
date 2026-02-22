package services

import (
	"bytes"
	"context"
	"fmt"
	"image"
	"strconv"
	"strings"
	"time"

	"github.com/disintegration/imaging"
	"github.com/jaygaha/roadmap-go-projects/intermediate/image-processing-service/internal/database"
	"github.com/jaygaha/roadmap-go-projects/intermediate/image-processing-service/internal/logger"
	"github.com/jaygaha/roadmap-go-projects/intermediate/image-processing-service/internal/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.uber.org/zap"
)

type imageProcessingJob struct {
	ImageID    string
	Operation  string
	Parameters map[string]string
}

type ImageService struct {
	db      *database.MongoDB
	storage *S3Service
	logger  *zap.Logger
	jobs    chan imageProcessingJob
}

func NewImageService(db *database.MongoDB, storage *S3Service, logger *zap.Logger) *ImageService {
	s := &ImageService{
		db:      db,
		storage: storage,
		logger:  logger,
		jobs:    make(chan imageProcessingJob, 100),
	}
	go s.startWorkers(2)
	return s
}

func (s *ImageService) startWorkers(count int) {
	for i := 0; i < count; i++ {
		go s.worker()
	}
}

func (s *ImageService) worker() {
	for job := range s.jobs {
		img, err := s.GetImage(job.ImageID)
		if err != nil {
			s.logger.Error("Failed to load image for processing", logger.Error(err))
			continue
		}
		if img == nil {
			s.logger.Warn("Image not found for processing", logger.String("image_id", job.ImageID))
			continue
		}
		ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
		_, err = s.ProcessImage(ctx, img, job.Operation, job.Parameters)
		cancel()
		if err != nil {
			s.logger.Error("Failed to process image asynchronously", logger.Error(err))
		}
	}
}

func (s *ImageService) EnqueueProcessing(imageID, operation string, parameters map[string]string) error {
	if s.jobs == nil {
		return fmt.Errorf("image processing queue not initialized")
	}
	job := imageProcessingJob{
		ImageID:    imageID,
		Operation:  operation,
		Parameters: parameters,
	}
	select {
	case s.jobs <- job:
		return nil
	default:
		return fmt.Errorf("image processing queue is full")
	}
}

// CreateImage creates a new image record in the database
func (s *ImageService) CreateImage(image *models.Image) error {
	image.ID = primitive.NewObjectID()
	image.CreatedAt = time.Now()
	image.UpdatedAt = time.Now()

	collection := s.db.Collection("images")
	_, err := collection.InsertOne(context.Background(), image)
	if err != nil {
		s.logger.Error("Failed to insert image", logger.Error(err))
		return fmt.Errorf("failed to insert image: %w", err)
	}

	s.logger.Info("Image created successfully", logger.String("id", image.ID.Hex()))
	return nil
}

// GetImage retrieves an image by ID
func (s *ImageService) GetImage(id string) (*models.Image, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		s.logger.Error("Invalid image ID", logger.Error(err))
		return nil, fmt.Errorf("invalid image ID: %w", err)
	}

	collection := s.db.Collection("images")
	var image models.Image
	err = collection.FindOne(context.Background(), bson.M{"_id": objectID}).Decode(&image)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		s.logger.Error("Failed to get image", logger.Error(err))
		return nil, err
	}

	return &image, nil
}

// GetUserImages retrieves all images for a specific user with pagination
func (s *ImageService) GetUserImages(userID string, page, limit int) ([]*models.Image, int64, error) {
	userObjectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		s.logger.Error("Invalid user ID", logger.Error(err))
		return nil, 0, fmt.Errorf("invalid user ID: %w", err)
	}

	collection := s.db.Collection("images")
	filter := bson.M{"user_id": userObjectID}

	// Count total documents
	total, err := collection.CountDocuments(context.Background(), filter)
	if err != nil {
		s.logger.Error("Failed to count images", logger.Error(err))
		return nil, 0, err
	}

	// Calculate skip
	skip := (page - 1) * limit

	// Find with pagination
	opts := options.Find().
		SetSkip(int64(skip)).
		SetLimit(int64(limit)).
		SetSort(bson.D{{"created_at", -1}})

	cursor, err := collection.Find(context.Background(), filter, opts)
	if err != nil {
		s.logger.Error("Failed to find images", logger.Error(err))
		return nil, 0, err
	}
	defer cursor.Close(context.Background())

	var images []*models.Image
	for cursor.Next(context.Background()) {
		var image models.Image
		if err := cursor.Decode(&image); err != nil {
			s.logger.Error("Failed to decode image", logger.Error(err))
			continue
		}
		images = append(images, &image)
	}

	if err := cursor.Err(); err != nil {
		s.logger.Error("Cursor error", logger.Error(err))
		return nil, 0, err
	}

	return images, total, nil
}

// UpdateImageStatus updates the status of an image
func (s *ImageService) UpdateImageStatus(id string, status models.ImageStatus) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		s.logger.Error("Invalid image ID", logger.Error(err))
		return fmt.Errorf("invalid image ID: %w", err)
	}

	collection := s.db.Collection("images")
	update := bson.M{
		"$set": bson.M{
			"status":     status,
			"updated_at": time.Now(),
		},
	}

	_, err = collection.UpdateOne(context.Background(), bson.M{"_id": objectID}, update)
	if err != nil {
		s.logger.Error("Failed to update image status", logger.Error(err))
		return err
	}

	s.logger.Info("Image status updated", logger.String("image_id", id), logger.String("status", string(status)))
	return nil
}

// AddProcessedImage adds a processed version to an image
func (s *ImageService) AddProcessedImage(id string, processed models.ProcessedImage) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		s.logger.Error("Invalid image ID", logger.Error(err))
		return fmt.Errorf("invalid image ID: %w", err)
	}

	processed.CreatedAt = time.Now()

	collection := s.db.Collection("images")
	update := bson.M{
		"$push": bson.M{"processed": processed},
		"$set":  bson.M{"updated_at": time.Now()},
	}

	_, err = collection.UpdateOne(context.Background(), bson.M{"_id": objectID}, update)
	if err != nil {
		s.logger.Error("Failed to add processed image", logger.Error(err))
		return err
	}

	s.logger.Info("Processed image added", logger.String("image_id", id), logger.String("operation", processed.Operation))
	return nil
}

// ProcessImage processes an image with the specified operation
func (s *ImageService) ProcessImage(ctx context.Context, image *models.Image, operation string, parameters map[string]string) (*models.ProcessedImage, error) {
	if err := s.UpdateImageStatus(image.ID.Hex(), models.ImageStatusProcessing); err != nil {
		return nil, fmt.Errorf("failed to update status: %w", err)
	}

	processedKey := fmt.Sprintf("processed/%s_%s_%d", image.ID.Hex(), operation, time.Now().UnixNano())

	s.logger.Info("Processing image",
		logger.String("image_id", image.ID.Hex()),
		logger.String("operation", operation),
		logger.Any("parameters", parameters),
	)

	originalData, err := s.storage.GetFile(ctx, image.OriginalKey)
	if err != nil {
		_ = s.UpdateImageStatus(image.ID.Hex(), models.ImageStatusFailed)
		return nil, fmt.Errorf("failed to get original file: %w", err)
	}

	src, err := imaging.Decode(bytes.NewReader(originalData))
	if err != nil {
		_ = s.UpdateImageStatus(image.ID.Hex(), models.ImageStatusFailed)
		return nil, fmt.Errorf("failed to decode image: %w", err)
	}

	dst, contentType, format, err := s.applyOperation(src, image.ContentType, operation, parameters)
	if err != nil {
		_ = s.UpdateImageStatus(image.ID.Hex(), models.ImageStatusFailed)
		return nil, fmt.Errorf("failed to apply operation: %w", err)
	}

	var buf bytes.Buffer
	encodeOpts := []imaging.EncodeOption{}

	if format == imaging.JPEG {
		q := 90
		if v, ok := parameters["quality"]; ok {
			if parsed, err := strconv.Atoi(v); err == nil && parsed >= 1 && parsed <= 100 {
				q = parsed
			}
		}
		encodeOpts = append(encodeOpts, imaging.JPEGQuality(q))
	}

	if err := imaging.Encode(&buf, dst, format, encodeOpts...); err != nil {
		_ = s.UpdateImageStatus(image.ID.Hex(), models.ImageStatusFailed)
		return nil, fmt.Errorf("failed to encode processed image: %w", err)
	}

	data := buf.Bytes()

	if err := s.storage.UploadFile(ctx, processedKey, data, contentType); err != nil {
		_ = s.UpdateImageStatus(image.ID.Hex(), models.ImageStatusFailed)
		return nil, fmt.Errorf("failed to upload processed image: %w", err)
	}

	size := int64(len(data))
	bounds := dst.Bounds()

	processed := models.ProcessedImage{
		Key:         processedKey,
		Operation:   operation,
		Parameters:  parameters,
		ContentType: contentType,
		Size:        size,
		Width:       bounds.Dx(),
		Height:      bounds.Dy(),
	}

	if err := s.AddProcessedImage(image.ID.Hex(), processed); err != nil {
		return nil, fmt.Errorf("failed to add processed image: %w", err)
	}

	if err := s.UpdateImageStatus(image.ID.Hex(), models.ImageStatusCompleted); err != nil {
		s.logger.Error("Failed to update status to completed", logger.Error(err))
	}

	return &processed, nil
}

func (s *ImageService) applyOperation(src image.Image, originalContentType string, operation string, parameters map[string]string) (image.Image, string, imaging.Format, error) {
	targetFormat, contentType := s.resolveFormat(originalContentType, parameters)

	switch operation {
	case string(models.OperationResize):
		width, height, err := parseDimensions(parameters)
		if err != nil {
			return nil, "", 0, err
		}
		if width == 0 && height == 0 {
			return nil, "", 0, fmt.Errorf("width or height must be provided for resize")
		}
		dst := imaging.Resize(src, width, height, imaging.Lanczos)
		return dst, contentType, targetFormat, nil
	case string(models.OperationCrop):
		width, height, err := parseDimensions(parameters)
		if err != nil {
			return nil, "", 0, err
		}
		if width <= 0 || height <= 0 {
			return nil, "", 0, fmt.Errorf("width and height must be greater than zero for crop")
		}
		dst := imaging.CropCenter(src, width, height)
		return dst, contentType, targetFormat, nil
	case string(models.OperationRotate):
		angleStr := parameters["angle"]
		switch angleStr {
		case "90":
			return imaging.Rotate90(src), contentType, targetFormat, nil
		case "180":
			return imaging.Rotate180(src), contentType, targetFormat, nil
		case "270":
			return imaging.Rotate270(src), contentType, targetFormat, nil
		default:
			return nil, "", 0, fmt.Errorf("unsupported rotate angle: %s", angleStr)
		}
	case string(models.OperationFlip):
		mode := strings.ToLower(parameters["mode"])
		if mode == "vertical" {
			return imaging.FlipV(src), contentType, targetFormat, nil
		}
		return imaging.FlipH(src), contentType, targetFormat, nil
	case string(models.OperationGrayscale):
		return imaging.Grayscale(src), contentType, targetFormat, nil
	case string(models.OperationBlur):
		sigma := parseFloat(parameters["sigma"], 1.5)
		return imaging.Blur(src, sigma), contentType, targetFormat, nil
	case string(models.OperationSharpen):
		sigma := parseFloat(parameters["sigma"], 1.0)
		return imaging.Sharpen(src, sigma), contentType, targetFormat, nil
	case string(models.OperationBrightness):
		percent := parseFloat(parameters["value"], 0)
		return imaging.AdjustBrightness(src, percent), contentType, targetFormat, nil
	case string(models.OperationContrast):
		percent := parseFloat(parameters["value"], 0)
		return imaging.AdjustContrast(src, percent), contentType, targetFormat, nil
	case string(models.OperationCompress):
		return src, contentType, targetFormat, nil
	default:
		return nil, "", 0, fmt.Errorf("unsupported operation: %s", operation)
	}
}

func (s *ImageService) resolveFormat(originalContentType string, parameters map[string]string) (imaging.Format, string) {
	if parameters == nil {
		parameters = map[string]string{}
	}

	if formatStr, ok := parameters["format"]; ok {
		switch strings.ToLower(formatStr) {
		case "jpeg", "jpg":
			return imaging.JPEG, "image/jpeg"
		case "png":
			return imaging.PNG, "image/png"
		case "gif":
			return imaging.GIF, "image/gif"
		case "bmp":
			return imaging.BMP, "image/bmp"
		case "tiff", "tif":
			return imaging.TIFF, "image/tiff"
		}
	}

	switch strings.ToLower(originalContentType) {
	case "image/jpeg", "image/jpg":
		return imaging.JPEG, "image/jpeg"
	case "image/png":
		return imaging.PNG, "image/png"
	case "image/gif":
		return imaging.GIF, "image/gif"
	case "image/bmp":
		return imaging.BMP, "image/bmp"
	case "image/tiff":
		return imaging.TIFF, "image/tiff"
	default:
		return imaging.PNG, "image/png"
	}
}

func parseDimensions(parameters map[string]string) (int, int, error) {
	width := 0
	height := 0

	if v, ok := parameters["width"]; ok && v != "" {
		parsed, err := strconv.Atoi(v)
		if err != nil {
			return 0, 0, fmt.Errorf("invalid width: %s", v)
		}
		width = parsed
	}

	if v, ok := parameters["height"]; ok && v != "" {
		parsed, err := strconv.Atoi(v)
		if err != nil {
			return 0, 0, fmt.Errorf("invalid height: %s", v)
		}
		height = parsed
	}

	return width, height, nil
}

func parseFloat(value string, defaultValue float64) float64 {
	if value == "" {
		return defaultValue
	}
	parsed, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return defaultValue
	}
	return parsed
}

// DeleteImage deletes an image and all its processed versions
func (s *ImageService) DeleteImage(ctx context.Context, id string) error {
	// Get image first to get all keys for deletion
	image, err := s.GetImage(id)
	if err != nil {
		return fmt.Errorf("failed to get image: %w", err)
	}

	if image == nil {
		return fmt.Errorf("image not found")
	}

	// TODO: Delete files from S3 storage
	// Delete original file
	s.logger.Info("Would delete original file", logger.String("key", image.OriginalKey))

	// Delete processed files
	for _, processed := range image.Processed {
		s.logger.Info("Would delete processed file", logger.String("key", processed.Key))
	}

	// Delete from database
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		s.logger.Error("Invalid image ID", logger.Error(err))
		return fmt.Errorf("invalid image ID: %w", err)
	}

	collection := s.db.Collection("images")
	_, err = collection.DeleteOne(context.Background(), bson.M{"_id": objectID})
	if err != nil {
		s.logger.Error("Failed to delete image from database", logger.Error(err))
		return err
	}

	s.logger.Info("Image deleted successfully", logger.String("image_id", id))
	return nil
}
