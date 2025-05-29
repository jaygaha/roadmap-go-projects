package respository

import (
	"time"

	"github.com/jaygaha/roadmap-go-projects/expense-tracker-api/config"
	"github.com/jaygaha/roadmap-go-projects/expense-tracker-api/internal/database"
	"github.com/jaygaha/roadmap-go-projects/expense-tracker-api/internal/models"
	"gorm.io/gorm"
)

// CreateExpense creates a new expense
func CreateExpense(expense *models.Expense) error {
	return database.DB.Create(expense).Error
}

// GetExpensesByUser
func GetExpensesByUser(userID uint, cfg *config.Config) ([]models.Expense, error) {
	var expenses []models.Expense

	err := database.DB.Preload("ExpenseCategory").Where("user_id = ?", userID).Order("created_at desc").Find(&expenses).Error

	return expenses, err
}

// GetExpensesByCategory retrieves expenses for a specific category
func GetExpensesByUserWithFilter(userID uint, startDate, endDate *time.Time, cfg *config.Config) ([]models.Expense, error) {
	var expenses []models.Expense
	query := database.DB.Preload("ExpenseCategory").Where("user_id = ?", userID)

	if startDate != nil && endDate != nil {
		query = query.Where("created_at BETWEEN ? AND ?", *startDate, *endDate)
	}

	err := query.Order("created_at desc").Find(&expenses).Error

	return expenses, err
}

// GetExpenseByID retrieves an expense by its ID
func GetExpenseByID(expenseID uint) (*models.Expense, error) {
	var expense models.Expense
	err := database.DB.Preload("ExpenseCategory").Where("id =?", expenseID).First(&expense).Error
	if err != nil {
		return nil, err
	}
	return &expense, nil
}

// UpdateExpense updates an existing expense
func UpdateExpense(expense *models.Expense) error {
	return database.DB.Save(expense).Error
}

// DeleteExpense deletes an expense by its ID
// won't delete the data from the database but update the deleted_at field
func DeleteExpense(expense *models.Expense) error {
	expense.DeletedAt = gorm.DeletedAt{Time: time.Now(), Valid: true}

	return database.DB.Save(expense).Error
}

// HasExpensesForCategory checks if a category has any associated expenses
func HasExpensesForCategory(categoryID uint) (bool, error) {
	var count int64

	err := database.DB.Model(&models.Expense{}).Where("expense_category_id = ?", categoryID).Count(&count).Error
	if err != nil {
		return false, err
	}

	return count > 0, nil
}
