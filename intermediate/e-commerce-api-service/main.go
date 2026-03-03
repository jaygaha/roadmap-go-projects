// E-Commerce API Service
//
// @title           E-Commerce API Service
// @version         1.0
// @description     Layered REST API for auth, products, cart, checkout, and orders.
// @BasePath        /api/v1
//
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer {token}" to authenticate.
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
	"github.com/jaygaha/roadmap-go-projects/intermediate/e-commerce-api-service/docs"
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

	// Page Handler (HTML templates)
	pageH := handlers.NewPageHandler("templates", productSvc, cartSvc, orderSvc, authSvc)

	// 4. Seed admin user if not exists
	database.SeedAdmin(userRepo, cfg.AdminEmail, cfg.AdminPassword)

	// 5. Build router and start server
	// Configure swagger metadata
	docs.SwaggerInfo.Title = "E-Commerce API Service"
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.BasePath = "/api/v1"
	docs.SwaggerInfo.Host = "localhost:" + cfg.Port

	r := router.New(authSvc, authH, productH, cartH, orderH, pageH)

	log.Printf("[SERVER] Starting on :%s", cfg.Port)
	log.Printf("[SERVER] Open http://localhost:%s in your browser", cfg.Port)

	if err := http.ListenAndServe(":"+cfg.Port, r); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
