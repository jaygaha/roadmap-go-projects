package services

import (
	"fmt"
	"log"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/jaygaha/roadmap-go-projects/intermediate/e-commerce-api-service/models"
	"github.com/jaygaha/roadmap-go-projects/intermediate/e-commerce-api-service/repository"
	"golang.org/x/crypto/bcrypt"
)

// AuthService is the service for authentication
type AuthService struct {
	userRepo  *repository.UserRepo
	jwtSecret []byte
}

// NewAuthService creates a new instance of AuthService
func NewAuthService(repo *repository.UserRepo, secret string) *AuthService {
	return &AuthService{userRepo: repo, jwtSecret: []byte(secret)}
}

// Register registers a new user
func (s *AuthService) Register(req models.UserRegisterRequest) (*models.AuthResponse, error) {
	if req.Email == "" || req.Password == "" {
		return nil, fmt.Errorf("%w: email and password required", models.ErrBadRequest)
	}
	if len(req.Password) < 8 {
		return nil, fmt.Errorf("%w: password must be at least 8 characters", models.ErrBadRequest)
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("hashing password: %w", err)
	}

	user, err := s.userRepo.Create(req.Email, string(hashed), "customer")
	if err != nil {
		return nil, err
	}

	token, err := s.generateToken(user)
	if err != nil {
		return nil, err
	}

	return &models.AuthResponse{Token: token, User: *user}, nil
}

// Login logs in a user
func (s *AuthService) Login(req models.UserLoginRequest) (*models.AuthResponse, error) {
	user, err := s.userRepo.FindByEmail(req.Email)
	log.Printf("[SVC] Login user: %v", user)
	if err != nil {
		return nil, fmt.Errorf("%w: invalid credentials", models.ErrUnauthorized)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, fmt.Errorf("%w: invalid credentials", models.ErrUnauthorized)
	}

	token, err := s.generateToken(user)
	if err != nil {
		return nil, err
	}

	return &models.AuthResponse{Token: token, User: *user}, nil
}

// generateToken generates a JWT token for the user
func (s *AuthService) generateToken(user *models.User) (string, error) {
	claims := jwt.MapClaims{
		"user_id": user.ID,
		"email":   user.Email,
		"role":    user.Role,
		"exp":     time.Now().Add(72 * time.Hour).Unix(),
		"iat":     time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(s.jwtSecret)
}

// ValidateToken parses and validates a JWT, returning claims.
func (s *AuthService) ValidateToken(tokenStr string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return s.jwtSecret, nil
	})
	if err != nil || !token.Valid {
		return nil, models.ErrUnauthorized
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, models.ErrUnauthorized
	}
	return claims, nil
}
