package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jaygaha/roadmap-go-projects/advanced/movie-reservation-system/internal/model"
	"github.com/jaygaha/roadmap-go-projects/advanced/movie-reservation-system/internal/service"
)

// ReservationHandler handles reservation-related HTTP requests
type ReservationHandler struct {
	service *service.ReservationService
}

// NewReservationHandler creates a new reservation handler
func NewReservationHandler(service *service.ReservationService) *ReservationHandler {
	return &ReservationHandler{service: service}
}

// GetByID handles GET /api/reservations/:id
func (h *ReservationHandler) GetByID(c *gin.Context) {
	id := parseInt(c.Param("id"))
	if id == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid reservation ID"})
		return
	}

	reservation, err := h.service.GetByID(c.Request.Context(), int64(id))
	if err != nil {
		handleAPIError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"reservation": reservation,
	})
}

// GetAllByUserID handles GET /api/users/me/reservations
func (h *ReservationHandler) GetAllByUserID(c *gin.Context) {
	userID := getUserID(c)
	page, limit := parsePagination(c)

	reservations, total, err := h.service.GetAllByUserID(c.Request.Context(), int64(userID), page, limit)
	if err != nil {
		handleAPIError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"reservations": reservations,
		"total":        total,
	})
}

// GetAll handles GET /api/reservations (admin)
func (h *ReservationHandler) GetAll(c *gin.Context) {
	page, limit := parsePagination(c)

	reservations, total, err := h.service.GetAll(c.Request.Context(), page, limit)
	if err != nil {
		handleAPIError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"reservations": reservations,
		"total":        total,
	})
}

// Create handles POST /api/reservations
func (h *ReservationHandler) Create(c *gin.Context) {
	userID := getUserID(c)

	var req model.CreateReservationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	reservation, err := h.service.Create(c.Request.Context(), &req, int64(userID))
	if err != nil {
		handleAPIError(c, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"reservation": reservation,
	})
}

// Cancel handles DELETE /api/reservations/:id
func (h *ReservationHandler) Cancel(c *gin.Context) {
	userID := getUserID(c)
	id := parseInt(c.Param("id"))

	if id == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid reservation ID"})
		return
	}

	if err := h.service.Cancel(c.Request.Context(), int64(id), int64(userID)); err != nil {
		handleAPIError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "reservation cancelled successfully",
	})
}

// GetStats handles GET /api/reservations/stats (admin)
func (h *ReservationHandler) GetStats(c *gin.Context) {
	stats, err := h.service.GetStats(c.Request.Context())
	if err != nil {
		handleAPIError(c, err)
		return
	}

	c.JSON(http.StatusOK, stats)
}

// getUserID extracts user ID from context
func getUserID(c *gin.Context) int {
	userID, exists := c.Get("userID")
	if !exists {
		return 0
	}

	if id, ok := userID.(int64); ok {
		return int(id)
	}

	return 0
}
