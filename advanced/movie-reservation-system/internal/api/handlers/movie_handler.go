package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jaygaha/roadmap-go-projects/advanced/movie-reservation-system/internal/model"
	"github.com/jaygaha/roadmap-go-projects/advanced/movie-reservation-system/internal/service"
)

// MovieHandler handles movie-related HTTP requests
type MovieHandler struct {
	service *service.MovieService
}

// NewMovieHandler creates a new movie handler
func NewMovieHandler(service *service.MovieService) *MovieHandler {
	return &MovieHandler{service: service}
}

// GetByID handles GET /api/movies/:id
func (h *MovieHandler) GetByID(c *gin.Context) {
	id := parseInt(c.Param("id"))
	if id == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid movie ID"})
		return
	}

	movie, err := h.service.GetByID(c.Request.Context(), int64(id))
	if err != nil {
		handleAPIError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"movie": movie})
}

// GetAll handles GET /api/movies
func (h *MovieHandler) GetAll(c *gin.Context) {
	page, limit := parsePagination(c)
	movies, total, err := h.service.GetAll(c.Request.Context(), page, limit)
	if err != nil {
		handleAPIError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"movies": movies,
		"total":  total,
	})
}

// Create handles POST /api/movies
func (h *MovieHandler) Create(c *gin.Context) {
	var req model.MovieCreateRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	movie, err := h.service.Create(c.Request.Context(), &req)
	if err != nil {
		handleAPIError(c, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{"movie": movie})
}

// Update handles PUT /api/movies/:id
func (h *MovieHandler) Update(c *gin.Context) {
	id := parseInt(c.Param("id"))
	if id == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid movie ID"})
		return
	}

	var req model.MovieUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	// Get existing movie
	movie, err := h.service.GetByID(c.Request.Context(), int64(id))
	if err != nil {
		handleAPIError(c, err)
		return
	}

	// Apply updates
	if req.Title != "" {
		movie.Title = req.Title
	}
	if req.Description != "" {
		movie.Description = req.Description
	}
	if req.PosterURL != "" {
		movie.PosterURL = req.PosterURL
	}
	if req.GenreID != nil {
		movie.GenreID = *req.GenreID
	}

	if err := h.service.Update(c.Request.Context(), movie); err != nil {
		handleAPIError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"movie": movie})
}

// Delete handles DELETE /api/movies/:id
func (h *MovieHandler) Delete(c *gin.Context) {
	id := parseInt(c.Param("id"))
	if id == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid movie ID"})
		return
	}

	if err := h.service.Delete(c.Request.Context(), int64(id)); err != nil {
		handleAPIError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "movie deleted successfully"})
}

// FindByGenre handles GET /api/movies/genre/:genreId
func (h *MovieHandler) FindByGenre(c *gin.Context) {
	genreID := parseInt(c.Param("genreId"))
	if genreID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid genre ID"})
		return
	}

	page, limit := parsePagination(c)
	movies, total, err := h.service.FindByGenre(c.Request.Context(), int64(genreID), page, limit)
	if err != nil {
		handleAPIError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"movies": movies,
		"total":  total,
	})
}

// GetAllGenres handles GET /api/genres
func (h *MovieHandler) GetAllGenres(c *gin.Context) {
	genres, err := h.service.GetAllGenres(c.Request.Context())
	if err != nil {
		handleAPIError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"genres": genres,
		"total":  len(genres),
	})
}
