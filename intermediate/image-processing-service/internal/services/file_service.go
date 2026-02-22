package services

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"path/filepath"
	"time"

	"github.com/jaygaha/roadmap-go-projects/intermediate/image-processing-service/internal/logger"
	"go.uber.org/zap"
)

type FileService struct {
	s3Service *S3Service
	logger    *zap.Logger
}

func NewFileService(s3Service *S3Service, logger *zap.Logger) *FileService {
	return &FileService{
		s3Service: s3Service,
		logger:    logger,
	}
}

// UploadFile uploads a file to S3 and returns the key
func (s *FileService) UploadFile(ctx context.Context, file *multipart.FileHeader) (string, error) {
	src, err := file.Open()
	if err != nil {
		s.logger.Error("Failed to open uploaded file", logger.Error(err))
		return "", fmt.Errorf("failed to open file: %w", err)
	}
	defer src.Close()

	// Read file data
	data, err := io.ReadAll(src)
	if err != nil {
		s.logger.Error("Failed to read file data", logger.Error(err))
		return "", fmt.Errorf("failed to read file: %w", err)
	}

	// Generate unique key
	ext := filepath.Ext(file.Filename)
	key := fmt.Sprintf("uploads/%d%s", time.Now().UnixNano(), ext)

	// Upload to S3
	err = s.s3Service.UploadFile(ctx, key, data, file.Header.Get("Content-Type"))
	if err != nil {
		s.logger.Error("Failed to upload file to S3", logger.Error(err))
		return "", fmt.Errorf("failed to upload to S3: %w", err)
	}

	s.logger.Info("File uploaded successfully", logger.String("key", key))
	return key, nil
}

// GetFile retrieves a file from S3 by key
func (s *FileService) GetFile(ctx context.Context, key string) ([]byte, error) {
	output, err := s.s3Service.GetFile(ctx, key)
	if err != nil {
		s.logger.Error("Failed to get file from S3", logger.Error(err))
		return nil, fmt.Errorf("failed to get file from S3: %w", err)
	}

	return output, nil
}

// DeleteFile deletes a file from S3 by key
func (s *FileService) DeleteFile(ctx context.Context, key string) error {
	err := s.s3Service.DeleteFile(ctx, key)
	if err != nil {
		s.logger.Error("Failed to delete file from S3", logger.Error(err))
		return fmt.Errorf("failed to delete file from S3: %w", err)
	}

	s.logger.Info("File deleted successfully", logger.String("key", key))
	return nil
}
