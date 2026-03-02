package router

import (
	"github.com/go-chi/chi/v5"
	chiMw "github.com/go-chi/chi/v5/middleware"

	"github.com/jaygaha/roadmap-go-projects/intermediate/e-commerce-api-service/handlers"
	"github.com/jaygaha/roadmap-go-projects/intermediate/e-commerce-api-service/middleware"
	"github.com/jaygaha/roadmap-go-projects/intermediate/e-commerce-api-service/services"
)

func New(
	authSvc *services.AuthService,
	authH *handlers.AuthHandler,
	productH *handlers.ProductHandler,
	cartH *handlers.CartHandler,
	orderH *handlers.OrderHandler,
) *chi.Mux {
	r := chi.NewRouter()

	// Global middleware stack
	r.Use(chiMw.Recoverer)    // panic recovery
	r.Use(chiMw.RequestID)    // X-Request-Id header
	r.Use(middleware.Logging) // structured request logging

	r.Route("/api/v1", func(r chi.Router) {
		// ── Public routes
		r.Post("/auth/register", authH.Register)
		r.Post("/auth/login", authH.Login)

		r.Get("/products", productH.List)
		r.Get("/products/{id}", productH.GetById)

		// ── Authenticated routes
		r.Group(func(r chi.Router) {
			r.Use(middleware.Auth(authSvc))

			// Cart
			r.Get("/cart", cartH.GetCart)
			r.Post("/cart/items", cartH.AddItem)
			r.Put("/cart/items/{productId}", cartH.UpdateItem)
			r.Delete("/cart/items/{productId}", cartH.RemoveItem)

			// Checkout & Orders
			r.Post("/checkout", orderH.Checkout)
			r.Get("/orders", orderH.ListOrders)
			r.Get("/orders/{id}", orderH.GetOrder)

		})

		// ── Admin routes
		r.Group(func(r chi.Router) {
			r.Use(middleware.Auth(authSvc))
			r.Use(middleware.AdminOnly)

			r.Post("/admin/products", productH.Create)
			r.Put("/admin/products/{id}", productH.Update)
			r.Delete("/admin/products/{id}", productH.Delete)
		})
	})

	return r
}
