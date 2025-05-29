package models

import (
	"time"

	"gorm.io/gorm"
)

// ExpenseCategory represents a category for expenses
type ExpenseCategory struct {
	ID            uint           `gorm:"primaryKey" json:"id"`
	Name          string         `gorm:"unique;not null" json:"name"`
	CreatedUserID uint           `gorm:"not null" json:"created_user_id"`  // user/admin ID who created the category
	UpdatedUserID uint           `gorm:"not null"  json:"updated_user_id"` // user/admin ID who last updated the category
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
}
