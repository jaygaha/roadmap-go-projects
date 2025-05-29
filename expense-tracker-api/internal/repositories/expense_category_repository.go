package respository

import (
	"time"

	"github.com/jaygaha/roadmap-go-projects/expense-tracker-api/internal/database"
	"github.com/jaygaha/roadmap-go-projects/expense-tracker-api/internal/models"
	"gorm.io/gorm"
)

// CreateExpenseCategory creates a new expense category
func CreateExpenseCategory(category *models.ExpenseCategory) error {
	return database.DB.Create(category).Error
}

// ListExpenseCategories retrieves all expenses for a given user ID
func ListExpenseCategories() ([]models.ExpenseCategory, error) {
	var categories []models.ExpenseCategory
	err := database.DB.Find(&categories).Error

	return categories, err
}

// GetCategoryByName retrieves a category by its name
func GetCategoryByName(name string) (*models.ExpenseCategory, error) {
	var category models.ExpenseCategory

	// check against email and shouldn't be deleted
	result := database.DB.Where("name = ?", name).First(&category)
	if result.Error != nil {
		return nil, result.Error
	}

	return &category, nil
}

// GetCategoryByID retrieves a category by its ID
func GetCategoryByID(id uint) (*models.ExpenseCategory, error) {
	var category models.ExpenseCategory

	// check against email and shouldn't be deleted
	result := database.DB.Where("id = ?", id).First(&category)
	if result.Error != nil {
		return nil, result.Error
	}

	return &category, nil
}

// UpdateExpenseCategory updates an expense category
func UpdateExpenseCategory(category *models.ExpenseCategory) error {
	return database.DB.Save(category).Error
}

// DeleteExpenseCategory deletes an expense category
func DeleteExpenseCategory(category *models.ExpenseCategory, userID uint) error {
	// Update the deleted_at field to mark the category as deleted
	category.UpdatedUserID = userID
	category.DeletedAt = gorm.DeletedAt{Time: time.Now(), Valid: true}

	return database.DB.Save(category).Error
}
