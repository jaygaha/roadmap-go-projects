package models

import (
	"time"

	"gorm.io/gorm"
)

// Expense represents an expense table in the database
type Expense struct {
	ID                uint            `gorm:"primaryKey" json:"id"`
	Title             string          `gorm:"size:255;not null" json:"title"`
	Amount            float64         `gorm:"not null" json:"amount"`
	UserID            uint            `gorm:"not null" json:"user_id"`             // user ID who created the expense
	ExpenseCategoryID uint            `gorm:"not null" json:"expense_category_id"` // category ID for the expense
	ExpenseCategory   ExpenseCategory `json:"category"`                            // category for the expense
	CreatedAt         time.Time       `json:"created_at"`
	UpdatedAt         time.Time       `json:"updated_at"`
	DeletedAt         gorm.DeletedAt  `gorm:"index" json:"-"`
}
