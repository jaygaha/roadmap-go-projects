package main

import (
	"context"
	"log"

	_ "github.com/lib/pq"

	"github.com/gin-gonic/gin"
	"github.com/jaygaha/roadmap-go-projects/advanced/movie-reservation-system/internal/api"
	"github.com/jaygaha/roadmap-go-projects/advanced/movie-reservation-system/internal/auth"
	"github.com/jaygaha/roadmap-go-projects/advanced/movie-reservation-system/internal/config"
	"github.com/jaygaha/roadmap-go-projects/advanced/movie-reservation-system/internal/repository"
	"github.com/jaygaha/roadmap-go-projects/advanced/movie-reservation-system/internal/service"
)

func main() {
	// Load configuration
	cfg := config.LoadConfig()

	// Connect to database
	dbConfig := &repository.DBConfig{
		Host:     cfg.Database.Host,
		Port:     cfg.Database.Port,
		Name:     cfg.Database.Name,
		User:     cfg.Database.User,
		Password: cfg.Database.Password,
	}

	db, err := repository.New(dbConfig)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Run database migrations and seed on startup
	if err := db.RunMigrations(context.Background()); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	// Initialize auth service
	authService := auth.NewAuthService(cfg.JWT.Secret, cfg.JWT.ExpireTime)

	// Initialize services
	services := service.NewServices(db, authService)

	// Set up router
	router := gin.Default()

	// Setup API routes
	api.SetupRoutes(router, services, authService)

	// Start server
	addr := cfg.Server.Host + ":" + cfg.Server.Port
	log.Printf("Starting server on %s", addr)
	if err := router.Run(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
