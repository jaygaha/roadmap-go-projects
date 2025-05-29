package services

import (
	"errors"

	"github.com/jaygaha/roadmap-go-projects/expense-tracker-api/config"
	"github.com/jaygaha/roadmap-go-projects/expense-tracker-api/internal/models"
	userRespository "github.com/jaygaha/roadmap-go-projects/expense-tracker-api/internal/repositories"
)

// GetUserByEmail retrieves a user by their email address
func GetUserByEmail(email string, cfg *config.Config) (*models.User, error) {
	// Check if user already exists
	user, err := userRespository.GetUserByEmail(email)
	if err == nil {
		return &models.User{}, errors.New("user already exists")
	}

	return user, nil
}
