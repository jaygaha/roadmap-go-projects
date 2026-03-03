package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/jaygaha/roadmap-go-projects/intermediate/e-commerce-api-service/middleware"
	"github.com/jaygaha/roadmap-go-projects/intermediate/e-commerce-api-service/models"
)

// --- JSON helpers (DRY response writing) ---
func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

func readJSON(r *http.Request, v any) error {
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("[JSON] read body error: %v", err)
		return err
	}
	log.Printf("[JSON] raw body: %s", string(bodyBytes))
	dec := json.NewDecoder(bytes.NewReader(bodyBytes))
	dec.DisallowUnknownFields()

	if err := dec.Decode(v); err != nil {
		log.Printf("[JSON] decode error: %v", err)
		return err
	}
	return nil
}

// handleError maps domain errors to HTTP status codes.
// This is the single place where error→status translation happens.
func handleError(w http.ResponseWriter, err error) {
	status := http.StatusInternalServerError
	switch {
	case errors.Is(err, models.ErrNotFound):
		status = http.StatusNotFound
	case errors.Is(err, models.ErrConflict):
		status = http.StatusConflict
	case errors.Is(err, models.ErrUnauthorized):
		status = http.StatusUnauthorized
	case errors.Is(err, models.ErrForbidden):
		status = http.StatusForbidden
	case errors.Is(err, models.ErrBadRequest):
		status = http.StatusBadRequest
	case errors.Is(err, models.ErrInsufficientStock):
		status = http.StatusConflict
	case errors.Is(err, models.ErrEmptyCart):
		status = http.StatusBadRequest
	}
	writeJSON(w, status, map[string]string{"error": err.Error()})
}

// urlParamInt64 parses a URL parameter as an int64.
func urlParamInt64(r *http.Request, key string) (int64, error) {
	s := chi.URLParam(r, key)

	return strconv.ParseInt(s, 10, 64)
}

// getUserId retrieves the user ID from the request context.
func getUserId(r *http.Request) int64 {
	return r.Context().Value(middleware.UserIdKey).(int64)
}
