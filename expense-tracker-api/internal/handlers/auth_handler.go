package handlers

import (
	"net/http"
	"strings"
	"time"

	"github.com/jaygaha/roadmap-go-projects/expense-tracker-api/config"
	"github.com/jaygaha/roadmap-go-projects/expense-tracker-api/internal/request"
	"github.com/jaygaha/roadmap-go-projects/expense-tracker-api/internal/response"
	"github.com/jaygaha/roadmap-go-projects/expense-tracker-api/internal/services"
	"github.com/jaygaha/roadmap-go-projects/expense-tracker-api/pkg/utils"
)

// SignupRequest represents the request body for the signup endpoint
type SignupRequest struct {
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

// LoginRequest represents the request body for the login endpoint
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

// SignUpHandler handles the signup request
func SignUpHandler(cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req SignupRequest
		if err := response.JSON(r, &req); err != nil {
			response.Error(w, http.StatusBadRequest, "Invalid request body")
			return
		}

		// Validate the request
		if err := request.ValidateStruct(req); err != nil {
			response.Error(w, http.StatusUnprocessableEntity, err.Error())
			return
		}

		// Check if the email is already registered
		_, err := services.GetUserByEmail(req.Email, cfg)
		if err != nil {
			response.Error(w, http.StatusUnprocessableEntity, "Email already registered")
			return
		}

		// Create a new user
		token, err := services.Signup(req.Name, req.Email, req.Password, cfg)
		if err != nil {
			response.Error(w, http.StatusBadRequest, err.Error())
			return
		}

		// Return the token in the response
		response.Success(w, http.StatusCreated, map[string]string{"token": token})
	}
}

// LoginHandler handles the login request
func LoginHandler(cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req LoginRequest
		if err := response.JSON(r, &req); err != nil {
			response.Error(w, http.StatusBadRequest, "Invalid request body")
			return
		}

		// Validate the request
		if err := request.ValidateStruct(req); err != nil {
			response.Error(w, http.StatusUnprocessableEntity, err.Error())
			return
		}

		// Authenticate the user
		token, err := services.Login(req.Email, req.Password, cfg)
		if err != nil {
			response.Error(w, http.StatusUnauthorized, err.Error())
			return
		}

		// Return the token in the response
		response.Success(w, http.StatusOK, map[string]string{"token": token})
	}
}

// LogoutHandler handles the logout request
func LogoutHandler(cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

		claims, err := utils.ParseToken(tokenStr, cfg.JWTSecret)
		if err != nil {
			response.Error(w, http.StatusUnauthorized, "Invalid token")
			return
		}

		expUnix := int64(claims["exp"].(float64))
		utils.AddTokenToBlocklist(tokenStr, time.Unix(expUnix, 0))

		response.Success(w, http.StatusOK, map[string]string{
			"message": "Logged out successfully. Token invalidated.",
		})
	}
}
