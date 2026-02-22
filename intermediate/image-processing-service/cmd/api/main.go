package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jaygaha/roadmap-go-projects/intermediate/image-processing-service/internal/config"
	"github.com/jaygaha/roadmap-go-projects/intermediate/image-processing-service/internal/database"
	"github.com/jaygaha/roadmap-go-projects/intermediate/image-processing-service/internal/handlers"
	"github.com/jaygaha/roadmap-go-projects/intermediate/image-processing-service/internal/logger"
	"github.com/jaygaha/roadmap-go-projects/intermediate/image-processing-service/internal/routes"
	"github.com/jaygaha/roadmap-go-projects/intermediate/image-processing-service/internal/services"
)

func main() {
	// Initialize logger
	log := logger.New()
	defer log.Sync()

	// Load configuration
	cfg := config.LoadConfig()

	// Initialize database
	db, err := database.NewMongoDB(cfg.MongoURI)
	if err != nil {
		log.Fatal("Failed to connect to database", logger.Error(err))
	}

	defer db.Disconnect(context.Background())

	s3Service := services.NewS3Service(cfg)
	userService := services.NewUserService(db, log)
	fileService := services.NewFileService(s3Service, log)
	imageService := services.NewImageService(db, s3Service, log)

	h := handlers.New(userService, fileService, imageService, log, cfg)

	router := routes.SetupRouter(h, log, cfg)

	// Create server
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", cfg.Port),
		Handler: router,
	}

	// Start server in a goroutine
	go func() {
		log.Info("Starting server on :" + cfg.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("Failed to start server", logger.Error(err))
		}
	}()

	// Wait for interrupt signal to gracefully shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Info("Shutting down server...")

	// Graceful shutdown with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown", logger.Error(err))
	}

	log.Info("Server exited")
}
