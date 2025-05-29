package routes

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jaygaha/roadmap-go-projects/expense-tracker-api/config"
	"github.com/jaygaha/roadmap-go-projects/expense-tracker-api/internal/handlers"
	"github.com/jaygaha/roadmap-go-projects/expense-tracker-api/internal/middleware"
)

// RegisterRoutes registers all application routes
func RegisterRoutes(cfg *config.Config) *mux.Router {
	router := mux.NewRouter()

	// Health Check
	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write([]byte("The API server is up and running")); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}).Methods(http.MethodGet)

	// API versioning: /api/v1
	apiRouter := router.PathPrefix("/api/v1").Subrouter()

	// Auth routes
	authRoutes := apiRouter.PathPrefix("/auth").Subrouter()
	authRoutes.HandleFunc("/signup", handlers.SignUpHandler(cfg)).Methods(http.MethodPost)
	authRoutes.HandleFunc("/login", handlers.LoginHandler(cfg)).Methods(http.MethodPost)
	authRoutes.Handle("/logout", middleware.JWTAuthMiddleware(cfg)(handlers.LogoutHandler(cfg)))

	// Protected routes
	// Apply JWT middleware to all protected routes
	protectedRouter := apiRouter.PathPrefix("").Subrouter()
	protectedRouter.Use(middleware.JWTAuthMiddleware(cfg))

	// Expense routes
	expenseRoutes := protectedRouter.PathPrefix("/expenses").Subrouter()
	expenseRoutes.HandleFunc("", handlers.CreateExpensehandler(cfg)).Methods(http.MethodPost)
	expenseRoutes.HandleFunc("", handlers.ListExpensesHandler(cfg)).Methods(http.MethodGet)
	expenseRoutes.HandleFunc("/{id}", handlers.GetExpenseByIDHandler(cfg)).Methods(http.MethodGet)
	expenseRoutes.HandleFunc("/{id}", handlers.UpdateExpenseHandler(cfg)).Methods(http.MethodPut)
	expenseRoutes.HandleFunc("/{id}", handlers.DeleteExpenseHandler(cfg)).Methods(http.MethodDelete)

	// Category routes (moved to top-level protected routes)
	categoryRoutes := protectedRouter.PathPrefix("/categories").Subrouter()
	categoryRoutes.HandleFunc("", handlers.CreateCategoryHandler(cfg)).Methods(http.MethodPost)
	categoryRoutes.HandleFunc("", handlers.ListCategoriesHandler(cfg)).Methods(http.MethodGet)
	categoryRoutes.HandleFunc("/{id}", handlers.GetCategoryByIDHandler(cfg)).Methods(http.MethodGet)
	categoryRoutes.HandleFunc("/{id}", handlers.UpdateCategoryHandler(cfg)).Methods(http.MethodPut)
	categoryRoutes.HandleFunc("/{id}", handlers.DeleteCategoryHandler(cfg)).Methods(http.MethodDelete)

	return router
}
