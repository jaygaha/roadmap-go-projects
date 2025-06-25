package routes

import (
	"errors"
	"net/http"

	"github.com/jaygaha/roadmap-go-projects/intermediate/url-shortening-service/internal/handler"
	"github.com/jaygaha/roadmap-go-projects/intermediate/url-shortening-service/internal/utils"
)

// SetupRouter sets up the router
func SetupRouter(urlHandler *handler.URLHandler) *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", healthHandler)

	// Create and GetAll
	crUrlHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			urlHandler.CreateURL(w, r)
		case http.MethodGet:
			urlHandler.GetAllURLs(w, r)
		default:
			utils.WriteErrorResponse(w, http.StatusMethodNotAllowed, errors.New("Method Not Allowed"))
		}
	})

	mux.Handle("/api/shorten", crUrlHandler)

	// Update, get and delete
	rudUrlHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			urlHandler.RedirectOriginalURL(w, r)
		case http.MethodPatch:
			urlHandler.UpdateOriginalURL(w, r)
		case http.MethodDelete:
			urlHandler.DeleteOriginalURL(w, r)
		default:
			utils.WriteErrorResponse(w, http.StatusMethodNotAllowed, errors.New("Method Not Allowed"))
		}
	})

	mux.Handle("/api/shorten/{shortCode}", rudUrlHandler)

	mux.HandleFunc("/api/shorten/{shortCode}/stats", urlHandler.GetOriginalURLStats)

	return mux
}

// healthHandler
func healthHandler(w http.ResponseWriter, r *http.Request) {
	utils.WriteSuccessResponse(w, http.StatusOK, "OK", nil)
}
