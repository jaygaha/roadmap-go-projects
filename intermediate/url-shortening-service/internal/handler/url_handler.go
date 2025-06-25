package handler

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/jaygaha/roadmap-go-projects/intermediate/url-shortening-service/internal/models"
	"github.com/jaygaha/roadmap-go-projects/intermediate/url-shortening-service/internal/service"
	"github.com/jaygaha/roadmap-go-projects/intermediate/url-shortening-service/internal/utils"
)

// URLHandler handles URL-related operations
type URLHandler struct {
	service service.URLService
}

// NewURLHandler creates a new URLHandler instance
func NewURLHandler(service service.URLService) *URLHandler {
	return &URLHandler{
		service: service,
	}
}

// CreateURL creates a new short URL
func (h *URLHandler) CreateURL(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.WriteErrorResponse(w, http.StatusMethodNotAllowed, errors.New("Method Not Allowed"))
		return
	}

	var req models.CreateURLRequest
	// JSON parse
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteErrorResponse(w, http.StatusMethodNotAllowed, errors.New("Invalid request body"))
		return
	}

	// Validate
	if req.URL == "" {
		utils.WriteErrorResponse(w, http.StatusUnprocessableEntity, errors.New("URL is required"))
		return
	}

	// Create URL
	shortURL, err := h.service.CreateShortURL(r.Context(), req.URL)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusInternalServerError, err)
		return
	}

	// Return response
	utils.WriteSuccessResponse(w, http.StatusCreated, "URL created successfully", shortURL)
}

// GetAllURLs returns all URLs
func (h *URLHandler) GetAllURLs(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.WriteErrorResponse(w, http.StatusMethodNotAllowed, errors.New("Method Not Allowed"))
		return
	}

	pageParam := r.URL.Query().Get("page")
	limitParam := r.URL.Query().Get("limit")

	// default values
	page := 1
	limit := 20

	// Parse query params if provided
	if pageParam != "" {
		p, err := strconv.Atoi(pageParam)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusBadRequest, errors.New("Invalid page parameter"))
			return
		}
		page = p
	}
	if limitParam != "" {
		l, err := strconv.Atoi(limitParam)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusBadRequest, errors.New("Invalid limit parameter"))
			return
		}
		limit = l
	}

	// Get all URLs
	urls, total, err := h.service.GetAllURLs(r.Context(), page, limit)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusInternalServerError, err)
		return
	}

	totalPages := (total + limit - 1) / limit // ceiling division
	response := map[string]any{
		"list": urls,
		"pagination": map[string]any{
			"page":  page,
			"limit": limit,
			"total": total,
			"pages": totalPages,
		},
	}

	// Return response
	utils.WriteSuccessResponse(w, http.StatusOK, "URLs retrieved successfully", response)
}

// RedirectOriginalURL redirects to the original URL
func (h *URLHandler) RedirectOriginalURL(w http.ResponseWriter, r *http.Request) {
	shortCode, err := utils.GetCodeFromRequest(r)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, err)
		return
	}

	originalURL, err := h.service.GetOriginalURL(r.Context(), shortCode)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusInternalServerError, err)
		return
	}

	// Return response
	utils.WriteSuccessResponse(w, http.StatusOK, "URL retrieved successfully", originalURL)
}

// GetOriginalURLStats returns data with stats
func (h *URLHandler) GetOriginalURLStats(w http.ResponseWriter, r *http.Request) {
	shortCode, err := utils.GetCodeFromRequest(r)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, err)
		return
	}

	// xxx/stats remove the last 6 characters (/stats)
	shortCode = shortCode[:len(shortCode)-6]

	originalURL, err := h.service.GetOriginalURLStats(r.Context(), shortCode)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusInternalServerError, err)
		return
	}

	log.Println(originalURL)

	// Return response
	utils.WriteSuccessResponse(w, http.StatusOK, "URL retrieved successfully", originalURL)
}

// UpdateOriginalURL updates the original URL
func (h *URLHandler) UpdateOriginalURL(w http.ResponseWriter, r *http.Request) {
	shortCode, err := utils.GetCodeFromRequest(r)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, err)
		return
	}

	var req models.CreateURLRequest
	// JSON parse
	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteErrorResponse(w, http.StatusMethodNotAllowed, errors.New("Invalid request body"))
		return
	}

	// Validate
	if req.URL == "" {
		utils.WriteErrorResponse(w, http.StatusUnprocessableEntity, errors.New("URL is required"))
		return
	}

	// Update URL
	url, err := h.service.UpdateURL(r.Context(), shortCode, req.URL)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusInternalServerError, err)
		return
	}

	// Return response
	utils.WriteSuccessResponse(w, http.StatusOK, "URL updated successfully", url)
}

// DeleteOriginalURL deletes the original URL
func (h *URLHandler) DeleteOriginalURL(w http.ResponseWriter, r *http.Request) {
	shortCode, err := utils.GetCodeFromRequest(r)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, err)
		return
	}

	// Delete URL
	err = h.service.DeleteURL(r.Context(), shortCode)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusNotFound, err)
		return
	}

	// Return response
	utils.WriteSuccessResponse(w, http.StatusNoContent, "URL deleted successfully", nil)
}
