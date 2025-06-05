package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/jaygaha/roadmap-go-projects/weather-api/internal/config"
	"github.com/jaygaha/roadmap-go-projects/weather-api/internal/db"
	"github.com/jaygaha/roadmap-go-projects/weather-api/internal/handlers"
)

func main() {
	// Load config
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
		return
	}

	// Connect to Redis
	redisClient := db.NewRedisClient(cfg.RedisAddr, cfg.RedisPassword)

	srv := &handlers.Server{
		RedisClient: redisClient,
		Config:      cfg,
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Welcome to the weather API!"))
	})
	mux.HandleFunc("/weathers", srv.GetWeatherData)

	log.Printf("Server is running on port %d", cfg.Port)
	err = http.ListenAndServe(fmt.Sprintf(":%d", cfg.Port), mux)
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}

}
