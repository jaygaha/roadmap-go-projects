package main

import (
	"log"
	"net/http"

	"github.com/jaygaha/roadmap-go-projects/intermediate/e-commerce-api-service/config"
	"github.com/jaygaha/roadmap-go-projects/intermediate/e-commerce-api-service/database"
	"github.com/jaygaha/roadmap-go-projects/intermediate/e-commerce-api-service/handlers"
	"github.com/jaygaha/roadmap-go-projects/intermediate/e-commerce-api-service/repository"
	"github.com/jaygaha/roadmap-go-projects/intermediate/e-commerce-api-service/router"
	"github.com/jaygaha/roadmap-go-projects/intermediate/e-commerce-api-service/services"
)

func main() {
	// 1. Load configuration
	cfg := config.Load()

	// 2. Initialize database
	db, err := database.New(cfg.DBPath)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	if err := database.RunMigrations(db); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	// 3. Wire up layers: repos → services → handlers
	userRepo := repository.NewUserRepo(db)
	productRepo := repository.NewProductRepo(db)
	cartRepo := repository.NewCartRepo(db)
	orderRepo := repository.NewOrderRepo(db)

	authSvc := services.NewAuthService(userRepo, cfg.JWTSecret)
	productSvc := services.NewProductService(productRepo)
	cartSvc := services.NewCartService(cartRepo, productRepo)
	paymentSvc := services.NewPaymentService(cfg.StripeKey)
	orderSvc := services.NewOrderService(orderRepo, cartRepo, productRepo, paymentSvc)

	authH := handlers.NewAuthHandler(authSvc)
	productH := handlers.NewProductHandler(productSvc)
	cartH := handlers.NewCartHandler(cartSvc)
	orderH := handlers.NewOrderHandler(orderSvc)

	// 4. Seed admin user if not exists
	database.SeedAdmin(userRepo, cfg.AdminEmail, cfg.AdminPassword)

	// 5. Build router and start server
	r := router.New(authSvc, authH, productH, cartH, orderH)

	log.Printf("[SERVER] Starting on :%s", cfg.Port)
	if err := http.ListenAndServe(":"+cfg.Port, r); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
