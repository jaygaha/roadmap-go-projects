package respository

import (
	"github.com/jaygaha/roadmap-go-projects/expense-tracker-api/internal/database"
	"github.com/jaygaha/roadmap-go-projects/expense-tracker-api/internal/models"
)

// CreateUser creates a new user in the database
func CreateUser(user *models.User) error {
	result := database.DB.Create(user)

	return result.Error
}

// GetUserByEmail retrieves a user by their email from the database
func GetUserByEmail(email string) (*models.User, error) {
	var user models.User

	// check against email and shouldn't be deleted
	result := database.DB.Where("email = ?", email).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}

	return &user, nil
}

// GetUserByID retrieves a user by their ID from the database
func GetUserByID(userID int) (*models.User, error) {
	var user models.User
	// check against email and shouldn't be deleted
	result := database.DB.Where("id = ?", userID).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}

	return &user, nil
}
