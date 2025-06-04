package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/jaygaha/roadmap-go-projects/todo-api/internal/auth"
	"github.com/jaygaha/roadmap-go-projects/todo-api/internal/db"
	"github.com/jaygaha/roadmap-go-projects/todo-api/internal/models"
	"github.com/jaygaha/roadmap-go-projects/todo-api/internal/utils"
	"golang.org/x/crypto/bcrypt"
)

// RegisterRequest represents a request to register a user.
type RegisterRequest struct {
	Name     string `json:"name"`
	Username string `json:"username"`
	Password string `json:"password"`
}

// LoginRequest represents a request to login a user.
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// RegisterHandler handles requests to register a user.
func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondJSON(w, http.StatusBadRequest, "Invalid request", nil)
		return
	}

	req.Name = strings.TrimSpace(req.Name)
	req.Username = strings.TrimSpace(req.Username)
	req.Password = strings.TrimSpace(req.Password)
	if req.Name == "" || req.Username == "" || req.Password == "" {
		utils.RespondJSON(w, http.StatusUnprocessableEntity, "Invalid request", nil)
		return
	}

	// Check if the username is already taken
	var count int
	err := db.DB.QueryRow("SELECT COUNT(*) FROM users WHERE username = $1", req.Username).Scan(&count)
	if err != nil {
		utils.RespondJSON(w, http.StatusInternalServerError, "Username is already taken", nil)
		return
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		utils.RespondJSON(w, http.StatusInternalServerError, "Failed to register user", nil)
		return
	}

	// Insert the user into the database
	userStatement := `
		INSERT INTO users (name, username, password)
		VALUES ($1, $2, $3)
		RETURNING id
	`

	row, err := db.DB.Exec(userStatement, req.Name, req.Username, string(hashedPassword))
	if err != nil {
		utils.RespondJSON(w, http.StatusInternalServerError, "Failed to register user", nil)
		return
	}

	// Check if the user was inserted successfully
	userId, err := row.LastInsertId()
	if err != nil {
		utils.RespondJSON(w, http.StatusInternalServerError, "Failed to register user", nil)
		return
	}

	// If successfull return jwt token
	token, err := auth.GenerateToken(userId, req.Username)
	if err != nil {
		utils.RespondJSON(w, http.StatusInternalServerError, "Failed to generate token", nil)
		return
	}
	utils.RespondJSON(w, http.StatusCreated, "User registered successfully", map[string]string{"token": token})
}

// LoginHandler handles requests to login a user.
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondJSON(w, http.StatusBadRequest, "Invalid request", nil)
		return
	}

	req.Username = strings.TrimSpace(req.Username)
	req.Password = strings.TrimSpace(req.Password)
	if req.Username == "" || req.Password == "" {
		utils.RespondJSON(w, http.StatusUnprocessableEntity, "Invalid request", nil)
		return
	}

	// Get the user from the database
	var user models.User
	row := db.DB.QueryRow("SELECT id, name, username, password FROM users WHERE username = $1", req.Username)
	err := row.Scan(&user.ID, &user.Name, &user.Username, &user.Password)
	if err == sql.ErrNoRows {
		utils.RespondJSON(w, http.StatusUnauthorized, "Invalid username or password", nil)
		return
	} else if err != nil {
		utils.RespondJSON(w, http.StatusInternalServerError, "Failed to login user", nil)
		return
	}

	// Compare the password hash
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		utils.RespondJSON(w, http.StatusUnauthorized, "Invalid username or password", nil)
		return
	}

	// If successfull return jwt token
	token, err := auth.GenerateToken(user.ID, user.Username)
	if err != nil {
		utils.RespondJSON(w, http.StatusInternalServerError, "Failed to generate token", nil)
		return
	}

	utils.RespondJSON(w, http.StatusOK, "User logged in successfully", map[string]string{"token": token})
}

// Helper function to get the user ID from the request context
func GetUserIDFromContext(r *http.Request) int64 {
	userID, ok := r.Context().Value("userID").(int64)
	if !ok {
		return 0
	}

	return userID
}
