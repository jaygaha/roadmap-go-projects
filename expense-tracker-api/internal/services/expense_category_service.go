package services

import (
	"errors"

	"github.com/jaygaha/roadmap-go-projects/expense-tracker-api/config"
	"github.com/jaygaha/roadmap-go-projects/expense-tracker-api/internal/models"
	respository "github.com/jaygaha/roadmap-go-projects/expense-tracker-api/internal/repositories"
)

// GetCategoryByName retrieves a category by its name
func GetCategoryByName(name string, cfg *config.Config) (*models.ExpenseCategory, error) {
	category, err := respository.GetCategoryByName(name)
	if err != nil {
		return &models.ExpenseCategory{}, errors.New("Category not found")
	}

	return category, nil
}

// CreateExpenseCategory creates a new expense category
func CreateExpenseCategory(name string, userID uint, cfg *config.Config) (*models.ExpenseCategory, error) {
	newCategory := &models.ExpenseCategory{
		Name:          name,
		CreatedUserID: userID,
		UpdatedUserID: userID,
	}

	if err := respository.CreateExpenseCategory(newCategory); err != nil {
		return &models.ExpenseCategory{}, errors.New("Failed to create category")
	}

	return newCategory, nil
}

// ListExpenseCategories retrieves a list of expense categories
func ListExpenseCategories(cfg *config.Config) ([]models.ExpenseCategory, error) {
	categories, err := respository.ListExpenseCategories()
	if err != nil {
		return []models.ExpenseCategory{}, errors.New("Failed to retrieve categories")
	}

	return categories, nil
}

// GetCategoryByID retrieves a category by its ID
func GetCategoryByID(id uint, cfg *config.Config) (*models.ExpenseCategory, error) {
	category, err := respository.GetCategoryByID(id)
	if err != nil {
		return &models.ExpenseCategory{}, errors.New("Category not found")
	}

	return category, nil
}

// UpdateExpenseCategory updates an existing expense category
func UpdateExpenseCategory(id uint, name string, userID uint, cfg *config.Config) (*models.ExpenseCategory, error) {
	category, err := respository.GetCategoryByID(id)
	if err != nil {
		return &models.ExpenseCategory{}, errors.New("Category not found")
	}

	category.Name = name
	category.UpdatedUserID = userID
	if err := respository.UpdateExpenseCategory(category); err != nil {
		return &models.ExpenseCategory{}, errors.New("Failed to update category")
	}

	return category, nil
}

// DeleteExpenseCategory deletes an expense category
func DeleteExpenseCategory(id, userID uint, cfg *config.Config) error {
	category, err := respository.GetCategoryByID(id)
	if err != nil {
		return errors.New("Category not found")
	}

	if err := respository.DeleteExpenseCategory(category, userID); err != nil {
		return errors.New("Failed to delete category")
	}

	return nil
}
