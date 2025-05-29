package services

import (
	"errors"
	"time"

	"github.com/jaygaha/roadmap-go-projects/expense-tracker-api/config"
	"github.com/jaygaha/roadmap-go-projects/expense-tracker-api/internal/models"
	userRespository "github.com/jaygaha/roadmap-go-projects/expense-tracker-api/internal/repositories"
	"github.com/jaygaha/roadmap-go-projects/expense-tracker-api/pkg/utils"
	"golang.org/x/crypto/bcrypt"
)

// HashPassword hashes the given password using bcrypt
func HashPassword(password string) (string, error) {
	// Generate a bcrypt hash of the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return "", err
	}

	return string(hashedPassword), nil
}

// CheckPasswordHash checks if the given password matches the hashed password
func CheckPasswordHash(password, hashedPassword string) error {
	// Compare the hashed password with the plaintext password
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

// Singup creates a new user and returns JWT token
func Signup(name, email, password string, cfg *config.Config) (string, error) {
	// Check if user already exists
	_, err := userRespository.GetUserByEmail(email)
	if err == nil {
		return "", errors.New("user already exists")
	}

	// Hash the password
	hashedPassword, err := HashPassword(password)
	if err != nil {
		return "", err
	}

	// Create the user
	newUser := &models.User{
		Name:     name,
		Email:    email,
		Password: hashedPassword,
	}

	err = userRespository.CreateUser(newUser)
	if err != nil {
		return "", err
	}

	// Generate JWT token
	token, err := utils.GenerateToken(uint(newUser.ID), cfg.JWTSecret, time.Hour*24*time.Duration(cfg.JWTEXP)) // Token expires in specified days
	if err != nil {
		return "", err
	}

	return token, nil
}

// Login validates credentials and returns JWT token
func Login(email, password string, cfg *config.Config) (string, error) {
	// Get the user by email
	user, err := userRespository.GetUserByEmail(email)
	if err != nil {
		return "", errors.New("invalid email or password")
	}

	// Check if the password is correct
	err = CheckPasswordHash(password, user.Password)
	if err != nil {
		return "", errors.New("invalid email or password")
	}

	token, err := utils.GenerateToken(uint(user.ID), cfg.JWTSecret, time.Hour*24*time.Duration(cfg.JWTEXP)) // Token expires in specified days
	if err != nil {
		return "", err
	}

	return token, nil
}
