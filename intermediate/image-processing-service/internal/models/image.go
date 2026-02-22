package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Image represents an image record in the database
type Image struct {
	ID          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	UserID      primitive.ObjectID `json:"user_id" bson:"user_id"`
	OriginalKey string             `json:"original_key" bson:"original_key"`
	Filename    string             `json:"filename" bson:"filename"`
	ContentType string             `json:"content_type" bson:"content_type"`
	Size        int64              `json:"size" bson:"size"`
	Width       int                `json:"width" bson:"width"`
	Height      int                `json:"height" bson:"height"`
	Status      ImageStatus        `json:"status" bson:"status"`
	Processed   []ProcessedImage   `json:"processed" bson:"processed"`
	CreatedAt   time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt   time.Time          `json:"updated_at" bson:"updated_at"`
}

// ProcessedImage represents a processed version of an image
type ProcessedImage struct {
	Key         string            `json:"key" bson:"key"`
	Operation   string            `json:"operation" bson:"operation"`
	Parameters  map[string]string `json:"parameters" bson:"parameters"`
	ContentType string            `json:"content_type" bson:"content_type"`
	Size        int64             `json:"size" bson:"size"`
	Width       int               `json:"width" bson:"width"`
	Height      int               `json:"height" bson:"height"`
	CreatedAt   time.Time         `json:"created_at" bson:"created_at"`
}

// ImageStatus represents the status of an image
type ImageStatus string

const (
	ImageStatusUploaded   ImageStatus = "uploaded"
	ImageStatusProcessing ImageStatus = "processing"
	ImageStatusCompleted  ImageStatus = "completed"
	ImageStatusFailed     ImageStatus = "failed"
)

// UploadImageRequest represents the request for uploading an image
type UploadImageRequest struct {
	Filename string `json:"filename" binding:"required"`
}

// ProcessImageRequest represents the request for processing an image
type ProcessImageRequest struct {
	Operation  string            `json:"operation" binding:"required"`
	Parameters map[string]string `json:"parameters"`
}

// ImageResponse represents the response for image operations
type ImageResponse struct {
	ID          string             `json:"id"`
	Filename    string             `json:"filename"`
	ContentType string             `json:"content_type"`
	Size        int64              `json:"size"`
	Width       int                `json:"width"`
	Height      int                `json:"height"`
	Status      ImageStatus        `json:"status"`
	Processed   []ProcessedImage   `json:"processed"`
	CreatedAt   time.Time          `json:"created_at"`
	UpdatedAt   time.Time          `json:"updated_at"`
}

// ToResponse converts Image model to ImageResponse
func (i *Image) ToResponse() *ImageResponse {
	return &ImageResponse{
		ID:          i.ID.Hex(),
		Filename:    i.Filename,
		ContentType: i.ContentType,
		Size:        i.Size,
		Width:       i.Width,
		Height:      i.Height,
		Status:      i.Status,
		Processed:   i.Processed,
		CreatedAt:   i.CreatedAt,
		UpdatedAt:   i.UpdatedAt,
	}
}

// ImageOperation represents supported image operations
type ImageOperation string

const (
	OperationResize     ImageOperation = "resize"
	OperationCrop       ImageOperation = "crop"
	OperationRotate     ImageOperation = "rotate"
	OperationFlip       ImageOperation = "flip"
	OperationGrayscale  ImageOperation = "grayscale"
	OperationBlur       ImageOperation = "blur"
	OperationSharpen    ImageOperation = "sharpen"
	OperationBrightness ImageOperation = "brightness"
	OperationContrast   ImageOperation = "contrast"
	OperationCompress   ImageOperation = "compress"
)

// IsValidOperation checks if the operation is supported
func IsValidOperation(operation string) bool {
	validOps := map[string]bool{
		string(OperationResize):     true,
		string(OperationCrop):       true,
		string(OperationRotate):     true,
		string(OperationFlip):       true,
		string(OperationGrayscale):  true,
		string(OperationBlur):       true,
		string(OperationSharpen):    true,
		string(OperationBrightness): true,
		string(OperationContrast):   true,
		string(OperationCompress):   true,
	}
	return validOps[operation]
}