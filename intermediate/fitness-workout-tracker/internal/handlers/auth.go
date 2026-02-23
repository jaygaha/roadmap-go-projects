package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jaygaha/roadmap-go-projects/intermediate/fitness-workout-tracker/internal/models"
	"golang.org/x/crypto/bcrypt"
)

// RegisterUser handles the POST /auth/register endpoint
// RegisterUser godoc
//
//	@ID				registerUser
//	@Summary		Register a new user
//	@Description	Register a new user with username, email, and password
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			registerRequest	body		models.RegisterRequest	true	"Register Request"
//	@Success		201	{object}	map[string]string	"User registered successfully. Please log in to continue."
//	@Failure		400	{object}	map[string]string	"Invalid request payload"
//	@Failure		422	{object}	map[string]string	"Username, email, and password are required"
//	@Failure		422	{object}	map[string]string	"Username or email already exists"
//	@Failure		500	{object}	map[string]string	"Internal server error"
//	@Router			/auth/register [post]
func RegisterUser(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req models.RegisterRequest

		// 1. Decode the JSON request body here
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, `{"error": "Invalid request payload"}`, http.StatusBadRequest)
			return
		}

		// 2. Validate the request payload
		if req.Username == "" || req.Email == "" || req.Password == "" {
			http.Error(w, `{"error": "Username, email, and password are required"}`, http.StatusUnprocessableEntity)
			return
		}

		// 3. Hash the password
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			http.Error(w, `{"error": "Failed to process password"}`, http.StatusInternalServerError)
			return
		}

		// 4. Insert into the database
		query := `INSERT INTO users (username, email, password_hash) VALUES (?, ?, ?)`
		_, err = db.Exec(query, req.Username, req.Email, string(hashedPassword))
		if err != nil {
			// Catching the unique constraint violation
			http.Error(w, `{"error": "Username or email already exists"}`, http.StatusUnprocessableEntity)
			return
		}

		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(`{"message": "User registered successfully. Please log in to continue."}`))
	}
}

// LoginUser handles the POST /login endpoint
// LoginUser godoc
//
//	@ID				loginUser
//	@Summary		Login a user
//	@Description	Login a user with email and password
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			loginRequest	body		models.LoginRequest	true	"Login Request"
//	@Success		200	{object}	map[string]string	"JWT token"
//	@Failure		400	{object}	map[string]string	"Invalid request"
//	@Failure		422	{object}	map[string]string	"Email and password are required"
//	@Failure		401	{object}	map[string]string	"Invalid credentials"
//	@Failure		500	{object}	map[string]string	"Internal server error"
//	@Router			/auth/login [post]
func LoginUser(db *sql.DB, jwtSecret string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req models.LoginRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, `{"error": "Invalid request"}`, http.StatusBadRequest)
			return
		}

		// Validate the request payload
		if req.Email == "" || req.Password == "" {
			http.Error(w, `{"error": "Email and password are required"}`, http.StatusUnprocessableEntity)
			return
		}

		// Fetch the user from the database
		var id int
		var username, passwordHash string
		query := `SELECT id, username, password_hash FROM users WHERE email = ?`
		err := db.QueryRow(query, req.Email).Scan(&id, &username, &passwordHash)
		if err != nil {
			http.Error(w, `{"error": "Invalid credentials"}`, http.StatusUnauthorized)
			return
		}

		// Compare the provided password with the stored hash
		err = bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(req.Password))
		if err != nil {
			http.Error(w, `{"error": "Invalid credentials"}`, http.StatusUnauthorized)
			return
		}

		// Next: Generate the JWT
		// Create the claims
		expirationTime := time.Now().Add(24 * time.Hour)
		claims := &models.Claims{
			UserID:   id,
			Username: username,
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(expirationTime),
			},
		}

		// Create the token with the HS256 algorithm
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

		// Sign the token with our secret key
		tokenString, err := token.SignedString([]byte(jwtSecret))
		if err != nil {
			http.Error(w, `{"error": "Failed to generate token"}`, http.StatusInternalServerError)
			return
		}

		// Send the token back to the user
		response := map[string]string{"token": tokenString}
		json.NewEncoder(w).Encode(response)
	}
}
