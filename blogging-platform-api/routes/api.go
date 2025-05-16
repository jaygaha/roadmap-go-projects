package routes

import (
	"fmt"
	"net/http"

	"github.com/jaygaha/roadmap-go-projects/blogging-platform-api/handlers"
)

// RegisterRoutes defines and registers the routes for the API.
func RegisterRoutes() {
	// Health check
	http.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "API service is running...")
	})

	// Blog routes
	http.HandleFunc("/posts", handlers.PostWOIdHandler)      // Multiple route handlers for /posts
	http.HandleFunc("/posts/", handlers.PostHandlerSelector) // Single route handler for /posts/{id}
}
