package routes

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jaygaha/roadmap-go-projects/intermediate/image-processing-service/internal/config"
	"github.com/jaygaha/roadmap-go-projects/intermediate/image-processing-service/internal/handlers"
	"github.com/jaygaha/roadmap-go-projects/intermediate/image-processing-service/internal/middleware"
	"go.uber.org/zap"
)

func SetupRouter(h *handlers.Handler, logger *zap.Logger, cfg *config.Config) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()

	// Apply middleware
	r.Use(middleware.CORS())
	r.Use(middleware.RequestLogger(logger))
	r.Use(middleware.RecoveryHandler(logger))
	r.Use(middleware.ErrorHandler(logger))

	// Health check
	r.GET("/", h.HealthCheck)
	r.GET("/health", h.HealthCheck)

	// API routes
	api := r.Group("/api/v1")
	{
		auth := api.Group("/auth")
		{
			auth.POST("/register", h.SignUpHandler)
			auth.POST("/login", h.LoginHandler)
		}

		images := api.Group("/images")
		{
			images.Use(middleware.AuthMiddleware(cfg))

			images.POST("/upload", h.UploadImageHandler)
			images.GET("/", h.GetUserImagesHandler)
			images.GET("/:id", h.GetImageHandler)
			images.POST("/:id/process", middleware.RateLimitMiddleware(10, time.Minute), h.ProcessImageHandler)
			images.GET("/:id/download", h.DownloadImageHandler)
			images.DELETE("/:id", h.DeleteImageHandler)
		}
	}

	return r
}
