package utils

import (
	"encoding/json"
	"net/http"
)

// writeJSONResponse writes a JSON response with the specified status code and data
func writeJSONResponse(w http.ResponseWriter, statusCode int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}

func WriteErrorResponse(w http.ResponseWriter, statusCode int, err error) {
	writeJSONResponse(w, statusCode, map[string]string{"error": err.Error()})
}

func WriteSuccessResponse(w http.ResponseWriter, statusCode int, message string, data any) {
	writeJSONResponse(w, statusCode, map[string]any{"message": message, "data": data})
}
