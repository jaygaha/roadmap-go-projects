package handlers

import (
	"net/http"

	"github.com/jaygaha/roadmap-go-projects/todo-api/internal/utils"
)

// HealthCheckHandler returns a health check response.
func HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	utils.RespondJSON(w, http.StatusOK, "API server is healthy", nil)
}
