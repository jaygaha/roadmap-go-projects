package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jaygaha/roadmap-go-projects/intermediate/image-processing-service/internal/errors"
	"github.com/jaygaha/roadmap-go-projects/intermediate/image-processing-service/internal/logger"
	"github.com/jaygaha/roadmap-go-projects/intermediate/image-processing-service/internal/models"
	"github.com/jaygaha/roadmap-go-projects/intermediate/image-processing-service/internal/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// UploadImageHandler handles image upload requests
func (h *Handler) UploadImageHandler(c *gin.Context) {
	// Get user ID from context (would be set by auth middleware)
	userIDStr := c.GetString("user_id")
	if userIDStr == "" {
		c.Error(errors.NewUnauthorizedError("User not authenticated"))
		return
	}

	userID, err := primitive.ObjectIDFromHex(userIDStr)
	if err != nil {
		c.Error(errors.NewValidationError("Invalid user ID"))
		return
	}

	// Get file from form
	file, err := c.FormFile("image")
	if err != nil {
		utils.ValidationErrorResponse(c, err)
		return
	}

	// Validate file size
	if err := utils.ValidateFileSize(file.Size); err != nil {
		utils.ValidationErrorResponse(c, err)
		return
	}

	// Validate file type
	contentType := file.Header.Get("Content-Type")
	if err := utils.ValidateFileType(contentType); err != nil {
		utils.ValidationErrorResponse(c, err)
		return
	}

	// Upload file to storage
	key, err := h.fileService.UploadFile(c.Request.Context(), file)
	if err != nil {
		h.logger.Error("Failed to upload file", logger.Error(err))
		c.Error(errors.NewFileError("upload", err))
		return
	}

	// Create image record
	image := &models.Image{
		UserID:      userID,
		OriginalKey: key,
		Filename:    file.Filename,
		ContentType: contentType,
		Size:        file.Size,
		Status:      models.ImageStatusUploaded,
	}

	if err := h.imageService.CreateImage(image); err != nil {
		h.logger.Error("Failed to create image record", logger.Error(err))
		c.Error(errors.NewDatabaseError("image creation", err))
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "Image uploaded successfully", image.ToResponse())
}

// GetImageHandler retrieves an image by ID
func (h *Handler) GetImageHandler(c *gin.Context) {
	imageID := c.Param("id")
	if imageID == "" {
		c.Error(errors.NewValidationError("Image ID is required"))
		return
	}

	image, err := h.imageService.GetImage(imageID)
	if err != nil {
		h.logger.Error("Failed to get image", logger.Error(err))
		c.Error(errors.NewDatabaseError("image retrieval", err))
		return
	}

	if image == nil {
		c.Error(errors.NewNotFoundError("Image"))
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Image retrieved successfully", image.ToResponse())
}

// GetUserImagesHandler retrieves all images for a user
func (h *Handler) GetUserImagesHandler(c *gin.Context) {
	userIDStr := c.GetString("user_id")
	if userIDStr == "" {
		c.Error(errors.NewUnauthorizedError("User not authenticated"))
		return
	}

	// Parse pagination parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	images, total, err := h.imageService.GetUserImages(userIDStr, page, limit)
	if err != nil {
		h.logger.Error("Failed to get user images", logger.Error(err))
		c.Error(errors.NewDatabaseError("images retrieval", err))
		return
	}

	// Convert to response format
	responses := make([]*models.ImageResponse, len(images))
	for i, img := range images {
		responses[i] = img.ToResponse()
	}

	result := gin.H{
		"images": responses,
		"pagination": gin.H{
			"page":  page,
			"limit": limit,
			"total": total,
		},
	}

	utils.SuccessResponse(c, http.StatusOK, "Images retrieved successfully", result)
}

// ProcessImageHandler handles image processing requests
func (h *Handler) ProcessImageHandler(c *gin.Context) {
	imageID := c.Param("id")
	if imageID == "" {
		c.Error(errors.NewValidationError("Image ID is required"))
		return
	}

	var req models.ProcessImageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, err)
		return
	}

	// Validate operation
	if !models.IsValidOperation(req.Operation) {
		c.Error(errors.NewValidationError("Invalid operation"))
		return
	}

	image, err := h.imageService.GetImage(imageID)
	if err != nil {
		h.logger.Error("Failed to get image", logger.Error(err))
		c.Error(errors.NewDatabaseError("image retrieval", err))
		return
	}

	if image == nil {
		c.Error(errors.NewNotFoundError("Image"))
		return
	}

	if err := h.imageService.EnqueueProcessing(imageID, req.Operation, req.Parameters); err != nil {
		h.logger.Error("Failed to enqueue image processing", logger.Error(err))
		c.Error(errors.NewInternalError("Image processing enqueue failed", err))
		return
	}

	utils.SuccessResponse(c, http.StatusAccepted, "Image processing started", image.ToResponse())
}

// DownloadImageHandler handles image download requests
func (h *Handler) DownloadImageHandler(c *gin.Context) {
	imageID := c.Param("id")
	processedKey := c.Query("processed")

	if imageID == "" {
		c.Error(errors.NewValidationError("Image ID is required"))
		return
	}

	// Get image
	image, err := h.imageService.GetImage(imageID)
	if err != nil {
		h.logger.Error("Failed to get image", logger.Error(err))
		c.Error(errors.NewDatabaseError("image retrieval", err))
		return
	}

	if image == nil {
		c.Error(errors.NewNotFoundError("Image"))
		return
	}

	// Determine which version to download
	var key string
	var contentType string
	var filename string

	if processedKey != "" {
		// Find processed version
		found := false
		for _, processed := range image.Processed {
			if processed.Key == processedKey {
				key = processed.Key
				contentType = processed.ContentType
				filename = "processed_" + image.Filename
				found = true
				break
			}
		}
		if !found {
			c.Error(errors.NewNotFoundError("Processed image"))
			return
		}
	} else {
		// Download original
		key = image.OriginalKey
		contentType = image.ContentType
		filename = image.Filename
	}

	// Get file data
	data, err := h.fileService.GetFile(c.Request.Context(), key)
	if err != nil {
		h.logger.Error("Failed to get file", logger.Error(err))
		c.Error(errors.NewFileError("download", err))
		return
	}

	// Set headers for file download
	c.Header("Content-Type", contentType)
	c.Header("Content-Disposition", "attachment; filename=\""+filename+"\"")
	c.Header("Content-Length", strconv.Itoa(len(data)))

	c.Data(http.StatusOK, contentType, data)
}

// DeleteImageHandler handles image deletion requests
func (h *Handler) DeleteImageHandler(c *gin.Context) {
	imageID := c.Param("id")
	if imageID == "" {
		c.Error(errors.NewValidationError("Image ID is required"))
		return
	}

	userIDStr := c.GetString("user_id")
	if userIDStr == "" {
		c.Error(errors.NewUnauthorizedError("User not authenticated"))
		return
	}

	// Get image to verify ownership
	image, err := h.imageService.GetImage(imageID)
	if err != nil {
		h.logger.Error("Failed to get image", logger.Error(err))
		c.Error(errors.NewDatabaseError("image retrieval", err))
		return
	}

	if image == nil {
		c.Error(errors.NewNotFoundError("Image"))
		return
	}

	// Verify ownership
	if image.UserID.Hex() != userIDStr {
		c.Error(errors.NewUnauthorizedError("Not authorized to delete this image"))
		return
	}

	// Delete image and all processed versions
	if err := h.imageService.DeleteImage(c.Request.Context(), imageID); err != nil {
		h.logger.Error("Failed to delete image", logger.Error(err))
		c.Error(errors.NewDatabaseError("image deletion", err))
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Image deleted successfully", nil)
}
