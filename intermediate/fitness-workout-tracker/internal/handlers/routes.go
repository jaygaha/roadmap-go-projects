package handlers

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"

	_ "github.com/jaygaha/roadmap-go-projects/intermediate/fitness-workout-tracker/docs"
	httpSwagger "github.com/swaggo/http-swagger"
)

// Import generated docs

// SetupRoutes configures all the API endpoints and injects dependencies
func SetupRoutes(db *sql.DB) http.Handler {

	_ = godotenv.Load() // Non-fatal if .env is missing; fall back to environment

	// Access the variable
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("JWT_SECRET not found in .env file")
	}

	mux := http.NewServeMux()

	// "/" handles both the empty path and the trailing slash

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Strict check: ensures we don't handle "/exercises" or other sub-paths here
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"message": "Fitness Tracker API is up and running!"}`))
	})

	// swagger routes
	mux.Handle("/swagger/", httpSwagger.WrapHandler)

	// Create a specific mux for v1
	v1 := http.NewServeMux()
	v1.HandleFunc("POST /auth/register", RegisterUser(db))
	v1.HandleFunc("POST /auth/login", LoginUser(db, jwtSecret))

	// attach auth middleware to all v1 routes
	auth := AuthMiddleware(jwtSecret)
	// exercise routes
	v1.Handle("GET /exercises", auth(GetExercises(db)))

	// workout routes
	v1.Handle("POST /workouts", auth(CreateWorkout(db)))
	v1.Handle("GET /workouts", auth(ListWorkouts(db)))
	v1.Handle("PUT /workouts/{id}", auth(UpdateWorkout(db)))
	v1.Handle("GET /workouts/{id}", auth(GetWorkout(db)))
	v1.Handle("DELETE /workouts/{id}", auth(DeleteWorkout(db)))
	// report routes
	v1.Handle("GET /workouts/reports", auth(GetWorkoutReport(db)))

	// Prefix all v1 routes with /v1
	mux.Handle("/api/v1/", http.StripPrefix("/api/v1", v1))

	// Wrap the mux with our JSONMiddleware
	return JSONMiddleware(mux)
}
