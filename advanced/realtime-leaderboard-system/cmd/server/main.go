package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/jaygaha/roadmap-go-projects/advanced/realtime-leaderboard-system/internal/auth"
	"github.com/jaygaha/roadmap-go-projects/advanced/realtime-leaderboard-system/internal/database"
	"github.com/jaygaha/roadmap-go-projects/advanced/realtime-leaderboard-system/internal/handlers"
	"github.com/jaygaha/roadmap-go-projects/advanced/realtime-leaderboard-system/internal/models"
	"github.com/jaygaha/roadmap-go-projects/advanced/realtime-leaderboard-system/pkg/config"
)

func main() {
	// Initialize databases
	database.InitSQLite()
	database.InitRedis()

	// Auto migrate models
	database.DB.AutoMigrate(&models.User{}, &models.Score{}, &models.Game{})

	// Setup Gin router
	r := gin.Default()

	// Global rate limiting
	r.Use(auth.RateLimitMiddleware())

	// Auth routes (public)
	authGroup := r.Group("/api/auth")
	{
		authGroup.POST("/register", handlers.Register)
		authGroup.POST("/login", handlers.Login)
		authGroup.GET("/profile", auth.JWTAuthMiddleware(), handlers.GetProfile)
	}

	// Game routes
	gameGroup := r.Group("/api/games")
	{
		gameGroup.GET("", handlers.ListGames)
		gameGroup.POST("", auth.JWTAuthMiddleware(), handlers.CreateGame)
	}

	// Score routes (protected)
	scoreGroup := r.Group("/api/scores")
	scoreGroup.Use(auth.JWTAuthMiddleware())
	{
		scoreGroup.POST("/submit", handlers.SubmitScore)
		scoreGroup.GET("/history", handlers.GetScoreHistory)
	}

	// Leaderboard routes (public)
	lbGroup := r.Group("/api/leaderboard")
	{
		lbGroup.GET("/global", handlers.GetGlobalLeaderboard)
		lbGroup.GET("/game/:gameId", handlers.GetGameLeaderboard)
		lbGroup.GET("/rank/:userId", handlers.GetUserRank)
		lbGroup.GET("/top/:gameId", handlers.GetTopPlayers)
	}

	log.Printf("Server starting on port %s", config.Port)
	r.Run(":" + config.Port)
}
