package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/jaygaha/roadmap-go-projects/expense-tracker-api/config"
	"github.com/jaygaha/roadmap-go-projects/expense-tracker-api/internal/database"
	"github.com/jaygaha/roadmap-go-projects/expense-tracker-api/internal/routes"
)

func main() {
	// Load configuration
	cfg := config.LoadConfig()

	// Connect to the database
	err := database.ConnectDB(cfg)
	if err != nil {
		log.Fatal("Error connecting to the database:", err)
		return
	}

	// Register routes
	router := routes.RegisterRoutes(cfg)

	// Start the server
	addr := fmt.Sprintf(":%s", cfg.AppPort)
	log.Println("Starting server on", addr)
	log.Fatal(http.ListenAndServe(addr, router))
}
