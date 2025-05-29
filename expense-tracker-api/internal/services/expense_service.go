package services

import (
	"errors"
	"time"

	"github.com/jaygaha/roadmap-go-projects/expense-tracker-api/config"
	"github.com/jaygaha/roadmap-go-projects/expense-tracker-api/internal/models"
	respository "github.com/jaygaha/roadmap-go-projects/expense-tracker-api/internal/repositories"
)

// CreateExpense creates a new expense
func CreateExpense(expense *models.Expense) error {
	// Create the expense in the database
	err := respository.CreateExpense(expense)
	if err != nil {
		return errors.New("Failed to create expense")
	}

	return nil
}

// ListExpenses retrieves a list of expenses
func ListExpensesWithFilter(userID uint, startDate, endDate *time.Time, cfg *config.Config) ([]models.Expense, error) {
	expenses, err := respository.GetExpensesByUserWithFilter(userID, startDate, endDate, cfg)
	if err != nil {
		return nil, errors.New("Failed to retrieve expenses")
	}

	return expenses, nil
}

// GetExpenseByID retrieves an expense by its ID
// Only the owner of the expense can view it
func GetExpenseByID(expenseID uint, userID uint, cfg *config.Config) (*models.Expense, error) {
	expense, err := respository.GetExpenseByID(expenseID)
	if err != nil {
		return nil, errors.New("Expense not found")
	}

	// Check if the user is the owner of the expense
	if expense.UserID != userID {
		return nil, errors.New("Unauthorized")
	}

	return expense, nil
}

// UpdateExpense updates an expense
func UpdateExpense(expense *models.Expense, cfg *config.Config) error {
	// Update the expense in the database
	err := respository.UpdateExpense(expense)
	if err != nil {
		return errors.New("Failed to update expense")
	}

	return nil
}

// DeleteExpense deletes an expense
func DeleteExpense(expense *models.Expense, cfg *config.Config) error {
	// Delete the expense from the database
	err := respository.DeleteExpense(expense)
	if err != nil {
		return errors.New("Failed to delete expense")
	}

	return nil
}

// HasExpensesForCategory checks if a category has associated expenses
func HasExpensesForCategory(categoryID uint, cfg *config.Config) (bool, error) {
	hasExpenses, err := respository.HasExpensesForCategory(categoryID)
	if err != nil {
		return false, errors.New("Failed to check expenses")
	}
	return hasExpenses, nil
}
