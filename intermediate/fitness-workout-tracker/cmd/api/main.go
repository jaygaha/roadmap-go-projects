package main

import (
	"log"
	"net/http"

	"github.com/jaygaha/roadmap-go-projects/intermediate/fitness-workout-tracker/internal/database"
	"github.com/jaygaha/roadmap-go-projects/intermediate/fitness-workout-tracker/internal/handlers"
)

// @title           Fitness Workout Tracker API
// @version         1.0
// @description     This is a sample server fitness workout tracker server.
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    https://jaygaha.com.np
// @contact.email  jaygaha@gmail.com

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8800
// @BasePath  /api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Provide JWT token as: Bearer <token>

func main() {
	// Initialize the database (creates tracker.db if it doesn't exist)
	db, err := database.InitDB("tracker.db")
	if err != nil {
		log.Fatalf("Error initializing database: %v\n", err)
	}
	defer db.Close()

	// Call our new router setup function, passing the database connection
	mux := handlers.SetupRoutes(db)

	log.Println("Server starting on :8800...")
	if err := http.ListenAndServe(":8800", mux); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
