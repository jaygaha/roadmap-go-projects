package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jaygaha/roadmap-go-projects/advanced/movie-reservation-system/internal/model"
	"github.com/jaygaha/roadmap-go-projects/advanced/movie-reservation-system/internal/service"
	pkg "github.com/jaygaha/roadmap-go-projects/advanced/movie-reservation-system/pkg"
)

// UserHandler handles user-related HTTP requests
type UserHandler struct {
	service *service.UserService
}

// NewUserHandler creates a new user handler
func NewUserHandler(service *service.UserService) *UserHandler {
	return &UserHandler{service: service}
}

// Signup handles POST /api/auth/signup
func (h *UserHandler) Signup(c *gin.Context) {
	var req model.UserSignupRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	resp, err := h.service.Signup(c.Request.Context(), &req)
	if err != nil {
		handleAPIError(c, err)
		return
	}

	c.JSON(http.StatusCreated, resp)
}

// Login handles POST /api/auth/login
func (h *UserHandler) Login(c *gin.Context) {
	var req model.UserLoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	resp, err := h.service.Login(c.Request.Context(), &req)
	if err != nil {
		handleAPIError(c, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

// RefreshToken handles POST /api/auth/refresh
func (h *UserHandler) RefreshToken(c *gin.Context) {
	var req model.RefreshTokenRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	// For now, we validate the old token and generate a new one
	// In a production system, you'd validate the refresh token
	// This is a simplified implementation
	c.JSON(http.StatusNotImplemented, gin.H{"error": "refresh token not fully implemented - use login endpoint"})
}

// GetByID handles GET /api/users/:id
func (h *UserHandler) GetByID(c *gin.Context) {
	id := parseInt(c.Param("id"))
	if id == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	user, err := h.service.GetByID(c.Request.Context(), int64(id))
	if err != nil {
		handleAPIError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": user})
}

// List handles GET /api/users
func (h *UserHandler) List(c *gin.Context) {
	page, limit := parsePagination(c)

	users, total, err := h.service.List(c.Request.Context(), page, limit)
	if err != nil {
		handleAPIError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"users": users,
		"total": total,
	})
}

// PromoteToAdmin handles POST /api/users/:id/promote (admin only)
func (h *UserHandler) PromoteToAdmin(c *gin.Context) {
	id := parseInt(c.Param("id"))
	if id == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	user, err := h.service.GetByID(c.Request.Context(), int64(id))
	if err != nil {
		handleAPIError(c, err)
		return
	}

	user.Role = model.RoleAdmin
	if err := h.service.Update(c.Request.Context(), user); err != nil {
		handleAPIError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "user promoted to admin",
		"user":    user,
	})
}

// parsePagination extracts page and limit from query params
func parsePagination(c *gin.Context) (int, int) {
	page := 1
	limit := 10

	if p := c.Query("page"); p != "" {
		page = parseInt(p, 1)
	}

	if l := c.Query("limit"); l != "" {
		limit = parseInt(l, 10)
	}

	return page, limit
}

// parseInt converts string to int with default fallback
func parseInt(s string, defaults ...int) int {
	defaultValue := 0
	if len(defaults) > 0 {
		defaultValue = defaults[0]
	}

	if s == "" {
		return defaultValue
	}

	result := 0
	negative := false

	for _, c := range s {
		if c == '-' {
			negative = true
			continue
		}
		if c < '0' || c > '9' {
			return defaultValue
		}
		result = result*10 + int(c-'0')
	}

	if negative {
		result = -result
	}

	return result
}

// handleAPIError handles error responses
func handleAPIError(c *gin.Context, err error) {
	switch {
	case err == pkg.ErrUserNotFound:
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	case err == pkg.ErrUserExists:
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
	case err == pkg.ErrInvalidPassword:
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
	case err == pkg.ErrDatabase:
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
}
