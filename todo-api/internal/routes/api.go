package routes

import (
	"net/http"

	"github.com/jaygaha/roadmap-go-projects/todo-api/internal/handlers"
	"github.com/jaygaha/roadmap-go-projects/todo-api/internal/middleware"
)

// RegisterRoutes registers all the routes for the application.
func RegisterRoutes() *http.ServeMux {
	mux := http.NewServeMux()

	// Public routes
	mux.HandleFunc("/health", handlers.HealthCheckHandler)
	mux.HandleFunc("/register", handlers.RegisterHandler)
	mux.HandleFunc("/login", handlers.LoginHandler)

	// Protected routes — wrap with AuthMiddleware
	todoHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handlers.ListTodosHandler(w, r)
		case http.MethodPost:
			handlers.CreateTodoHandler(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
	mux.Handle("/todos", middleware.AuthMiddleware(todoHandler))

	todoIdHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handlers.FetchTodoHandler(w, r)
		case http.MethodPut:
			handlers.UpdateTodoHandler(w, r)
		case http.MethodDelete:
			handlers.DeleteTodoHandler(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	mux.Handle("/todos/", middleware.AuthMiddleware(todoIdHandler))

	// mux.HandleFunc("/todos/", func(w http.ResponseWriter, r *http.Request) {
	// 	switch r.Method {
	// 	case http.MethodGet:
	// 		handlers.FetchTodoHandler(w, r)
	// 	case http.MethodPut:
	// 		handlers.UpdateTodoHandler(w, r)
	// 	case http.MethodDelete:
	// 		handlers.DeleteTodoHandler(w, r)
	// 	default:
	// 		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	// 	}
	// })

	return mux
}
