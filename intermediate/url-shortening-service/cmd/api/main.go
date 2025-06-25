package main

import (
	"log"
	"net/http"

	"github.com/jaygaha/roadmap-go-projects/intermediate/url-shortening-service/internal/config"
	"github.com/jaygaha/roadmap-go-projects/intermediate/url-shortening-service/internal/database"
	"github.com/jaygaha/roadmap-go-projects/intermediate/url-shortening-service/internal/handler"
	"github.com/jaygaha/roadmap-go-projects/intermediate/url-shortening-service/internal/repository/mongodb"
	"github.com/jaygaha/roadmap-go-projects/intermediate/url-shortening-service/internal/routes"
	"github.com/jaygaha/roadmap-go-projects/intermediate/url-shortening-service/internal/service"
)

func main() {
	// Load configuration
	cfg := config.LoadConfig()

	// Connect to db
	db, err := database.ConnectMongoDB(cfg.MongoDBURI, cfg.DatabaseName)
	if err != nil {
		log.Fatalf("%v", err)
	}

	// Initialize layers
	urlRepo := mongodb.NewURLRepository(db)
	urlService := service.NewURLService(urlRepo, cfg.BaseURL)
	urlHandler := handler.NewURLHandler(urlService)

	// Setup router
	mux := routes.SetupRouter(urlHandler)

	err = http.ListenAndServe(":8800", mux)
	if err != nil {
		log.Fatalf("Could not start server: %s\n", err.Error())
	}

	log.Println("Server started at 8800")
}
