package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jaygaha/roadmap-go-projects/intermediate/image-processing-service/internal/config"
	"github.com/jaygaha/roadmap-go-projects/intermediate/image-processing-service/internal/services"
	"github.com/jaygaha/roadmap-go-projects/intermediate/image-processing-service/internal/utils"
	"go.uber.org/zap"
)

type Handler struct {
	userService  *services.UserService
	fileService  *services.FileService
	imageService *services.ImageService
	logger       *zap.Logger
	cfg          *config.Config
}

func New(
	userService *services.UserService,
	fileService *services.FileService,
	imageService *services.ImageService,
	logger *zap.Logger,
	cfg *config.Config,
) *Handler {
	return &Handler{
		userService:  userService,
		fileService:  fileService,
		imageService: imageService,
		logger:       logger,
		cfg:          cfg,
	}
}

// HealthCheck endpoint provides service health status
func (h *Handler) HealthCheck(c *gin.Context) {
	healthData := gin.H{
		"service":   "image-processing-service",
		"version":   "1.0.0",
		"status":    "healthy",
		"timestamp": time.Now().Format(time.RFC3339),
	}

	utils.SuccessResponse(c, http.StatusOK, "Service is running", healthData)
}
