package utils

import (
	"encoding/json"
	"net/http"
)

// JSONResponse is a struct that represents a JSON response.
type JSONResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

// RespondJSON returns a JSON response.
func RespondJSON(w http.ResponseWriter, status int, message string, data any) {
	response := JSONResponse{
		Status:  http.StatusText(status),
		Message: message,
		Data:    data,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(response)
}
