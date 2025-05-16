package main

import (
	"log"
	"net/http"

	"github.com/jaygaha/roadmap-go-projects/blogging-platform-api/database"
	"github.com/jaygaha/roadmap-go-projects/blogging-platform-api/routes"
)

func main() {
	// Initialize the database
	database.ConnectDB()

	defer database.DB.Close()

	// Register the API routes
	routes.RegisterRoutes()

	// start the HTTP server
	log.Println("Server is running on port 8800...")
	log.Fatal(http.ListenAndServe(":8800", nil))
}
