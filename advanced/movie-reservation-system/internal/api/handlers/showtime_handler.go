package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jaygaha/roadmap-go-projects/advanced/movie-reservation-system/internal/model"
	"github.com/jaygaha/roadmap-go-projects/advanced/movie-reservation-system/internal/service"
)

// ShowtimeHandler handles showtime-related HTTP requests
type ShowtimeHandler struct {
	service *service.ShowtimeService
}

// NewShowtimeHandler creates a new showtime handler
func NewShowtimeHandler(service *service.ShowtimeService) *ShowtimeHandler {
	return &ShowtimeHandler{service: service}
}

// GetByID handles GET /api/showtimes/:id
func (h *ShowtimeHandler) GetByID(c *gin.Context) {
	id := parseInt(c.Param("id"))
	if id == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid showtime ID"})
		return
	}

	showtime, err := h.service.GetByID(c.Request.Context(), int64(id))
	if err != nil {
		handleAPIError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"showtime": showtime})
}

// GetAll handles GET /api/showtimes
func (h *ShowtimeHandler) GetAll(c *gin.Context) {
	page, limit := parsePagination(c)
	showtimes, total, err := h.service.GetAll(c.Request.Context(), page, limit)
	if err != nil {
		handleAPIError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"showtimes": showtimes,
		"total":     total,
	})
}

// GetByDate handles GET /api/showtimes/by-date
func (h *ShowtimeHandler) GetByDate(c *gin.Context) {
	date := c.Query("date")
	if date == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "date parameter is required"})
		return
	}

	// Validate date format
	if _, err := time.Parse("2006-01-02", date); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid date format, use YYYY-MM-DD"})
		return
	}

	showtimes, err := h.service.GetByDate(c.Request.Context(), date)
	if err != nil {
		handleAPIError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"date":      date,
		"showtimes": showtimes,
	})
}

// Create handles POST /api/showtimes
func (h *ShowtimeHandler) Create(c *gin.Context) {
	var req model.ShowtimeCreateRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	showtime, err := h.service.Create(c.Request.Context(), &req)
	if err != nil {
		handleAPIError(c, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{"showtime": showtime})
}

// Update handles PUT /api/showtimes/:id
func (h *ShowtimeHandler) Update(c *gin.Context) {
	id := parseInt(c.Param("id"))
	if id == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid showtime ID"})
		return
	}

	var req map[string]interface{}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	// Get existing showtime
	showtime, err := h.service.GetByID(c.Request.Context(), int64(id))
	if err != nil {
		handleAPIError(c, err)
		return
	}

	// Apply updates
	if movieID, ok := req["movie_id"]; ok {
		if mid, ok := movieID.(float64); ok {
			showtime.MovieID = int64(mid)
		}
	}
	if startTime, ok := req["start_time"]; ok {
		if st, ok := startTime.(string); ok {
			if t, err := time.Parse(time.RFC3339, st); err == nil {
				showtime.StartTime = t
			}
		}
	}
	if endTime, ok := req["end_time"]; ok {
		if et, ok := endTime.(string); ok {
			if t, err := time.Parse(time.RFC3339, et); err == nil {
				showtime.EndTime = t
			}
		}
	}
	showtime.UpdatedAt = time.Now()

	if err := h.service.Update(c.Request.Context(), showtime); err != nil {
		handleAPIError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"showtime": showtime})
}

// Delete handles DELETE /api/showtimes/:id
func (h *ShowtimeHandler) Delete(c *gin.Context) {
	id := parseInt(c.Param("id"))
	if id == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid showtime ID"})
		return
	}

	if err := h.service.Delete(c.Request.Context(), int64(id)); err != nil {
		handleAPIError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "showtime deleted successfully"})
}

// GetByMovieID handles GET /api/showtimes/movie/:movieId
func (h *ShowtimeHandler) GetByMovieID(c *gin.Context) {
	movieID := parseInt(c.Param("movieId"))
	if movieID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid movie ID"})
		return
	}

	showtimes, err := h.service.GetByMovieID(c.Request.Context(), int64(movieID))
	if err != nil {
		handleAPIError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"showtimes": showtimes})
}
