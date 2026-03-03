package handlers

import (
	"net/http"

	"github.com/jaygaha/roadmap-go-projects/intermediate/e-commerce-api-service/models"
	"github.com/jaygaha/roadmap-go-projects/intermediate/e-commerce-api-service/services"
)

// AuthHandler handles authentication requests.
type AuthHandler struct {
	svc *services.AuthService
}

// NewAuthHandler creates a new AuthHandler.
func NewAuthHandler(svc *services.AuthService) *AuthHandler {
	return &AuthHandler{svc: svc}
}

// Register handles user registration requests.
// @Summary      Register a new user
// @Description  Create a user account and return a JWT
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        payload  body      models.UserRegisterRequest  true  "User credentials"
// @Success      201      {object}  models.AuthResponse
// @Failure      400      {object}  map[string]string
// @Router       /auth/register [post]
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req models.UserRegisterRequest
	if err := readJSON(r, &req); err != nil {
		handleError(w, models.ErrBadRequest)
		return
	}

	resp, err := h.svc.Register(req)
	if err != nil {
		handleError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, resp)
}

// @Summary      Login
// @Description  Authenticate a user and return a JWT
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        payload  body      models.UserLoginRequest  true  "User credentials"
// @Success      200      {object}  models.AuthResponse
// @Failure      400      {object}  map[string]string
// @Failure      401      {object}  map[string]string
// @Router       /auth/login [post]
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req models.UserLoginRequest
	if err := readJSON(r, &req); err != nil {
		handleError(w, models.ErrBadRequest)
		return
	}

	resp, err := h.svc.Login(req)
	if err != nil {
		handleError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, resp)
}
