package database

import (
	"log"

	"github.com/jaygaha/roadmap-go-projects/expense-tracker-api/config"
	"github.com/jaygaha/roadmap-go-projects/expense-tracker-api/internal/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

// ConnectDB connects to the database
func ConnectDB(cfg *config.Config) error {
	var err error

	// Connect to the database
	DB, err = gorm.Open(sqlite.Open(cfg.DBName), &gorm.Config{})
	if err != nil {
		return err
	}

	// Migrate the database
	// Creates the tables if they don't exist according to the defined models structures
	if err = DB.AutoMigrate(&models.User{},
		&models.ExpenseCategory{},
		&models.Expense{}); err != nil {
		return err
	}

	// Seed the database
	if err = SeedDefaultDB(); err != nil {
		return err
	}

	log.Println("Connected to the database")
	return nil
}

// SeedDB seeds the database with initial data
func SeedDefaultDB() error {
	// Seed users
	err := seedUsers()
	if err != nil {
		return err
	}
	// Seed expense categories
	err = seedCategories()
	if err != nil {
		return err
	}
	return nil
}

// seedUsers seeds the users table
func seedUsers() error {
	// Check if users have already been seeded
	// If not, seed the users table with initial data
	var count int64

	err := DB.Model(&models.User{}).Count(&count).Error
	if err != nil {
		return err
	}

	if count > 0 {
		log.Println("Users have already been seeded, skipping...")
		return nil
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("root@123"), 14)
	if err != nil {
		return err
	}

	users := []models.User{
		{Name: "Jay", Email: "jaygaha@gmail.com", Password: string(hashedPassword)},
	}

	if err := DB.Create(&users).Error; err != nil {
		return err
	}

	log.Println("Users have been seeded successfully")
	return nil
}

// SeedCategories seeds the expense categories table
func seedCategories() error {
	// Check if expense categories have already been seeded
	// If not, seed the expense categories table with initial data
	var count int64
	err := DB.Model(&models.ExpenseCategory{}).Count(&count).Error
	if err != nil {
		return err
	}
	if count > 0 {
		log.Println("Expense categories have already been seeded, skipping...")
		return nil
	}

	// Get the user ID for the created user
	var user models.User
	if err := DB.First(&user).Error; err != nil {
		return err
	}
	// Seed the expense categories table with initial data
	categories := []models.ExpenseCategory{
		{Name: "Food", CreatedUserID: user.ID, UpdatedUserID: user.ID},
		{Name: "Transportation", CreatedUserID: user.ID, UpdatedUserID: user.ID},
		{Name: "Entertainment", CreatedUserID: user.ID, UpdatedUserID: user.ID},
		{Name: "Bills", CreatedUserID: user.ID, UpdatedUserID: user.ID},
		{Name: "Shopping", CreatedUserID: user.ID, UpdatedUserID: user.ID},
		{Name: "Other", CreatedUserID: user.ID, UpdatedUserID: user.ID},
	}
	if err := DB.Create(&categories).Error; err != nil {
		return err
	}

	log.Println("Expense categories have been seeded successfully")
	return nil
}
