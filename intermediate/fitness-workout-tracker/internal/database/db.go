package database

import (
	"database/sql"
	"log"

	"github.com/jaygaha/roadmap-go-projects/intermediate/fitness-workout-tracker/internal/models"
	_ "github.com/mattn/go-sqlite3" // Import the SQLite driver
)

// InitDB opens the database connection and runs migrations
func InitDB(filepath string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", filepath)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	err = runMigrations(db)
	if err != nil {
		return nil, err
	}

	// Insert seed exercises
	err = runSeedData(db)
	if err != nil {
		return nil, err
	}

	return db, nil
}

// runMigrations applies the necessary database migrations
func runMigrations(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		username TEXT UNIQUE NOT NULL,
		email TEXT UNIQUE NOT NULL,
		password_hash TEXT NOT NULL,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS exercises (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT UNIQUE NOT NULL, -- Name of the exercise
		category TEXT NOT NULL,    -- Category (e.g., Strength)
		muscle_group TEXT NOT NULL, -- Muscle group (e.g., Legs)
		description TEXT,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS workouts (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER NOT NULL,
		name TEXT NOT NULL,
		description TEXT,
		scheduled_for DATETIME NOT NULL,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
	);

	CREATE TABLE IF NOT EXISTS workout_exercises (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		workout_id INTEGER NOT NULL,
		exercise_id INTEGER NOT NULL,
		sets INTEGER NOT NULL,
		reps INTEGER NOT NULL,
		weight REAL NOT NULL,
		notes TEXT,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (workout_id) REFERENCES workouts(id) ON DELETE CASCADE,
		FOREIGN KEY (exercise_id) REFERENCES exercises(id) ON DELETE CASCADE
	);
	`

	_, err := db.Exec(query)
	if err != nil {
		log.Printf("Error running migrations: %v\n", err)
		return err
	}

	log.Println("Database migrated successfully.")
	return nil
}

// runSeedData inserts initial seed data into the database
func runSeedData(db *sql.DB) error {
	// Insert seed exercises
	exercises := []models.Exercise{
		{Name: "Barbell Bench Press", Category: "Strength", MuscleGroup: "Chest", Description: "Compound movement for chest power."},
		{Name: "Deadlift", Category: "Strength", MuscleGroup: "Back/Legs", Description: "Targets the entire posterior chain."},
		{Name: "Running", Category: "Cardio", MuscleGroup: "Full Body", Description: "High-intensity aerobic exercise."},
		{Name: "Pigeon Pose", Category: "Flexibility", MuscleGroup: "Hips", Description: "Deep stretch for hip mobility."},
		{Name: "Lat Pulldown", Category: "Strength", MuscleGroup: "Back", Description: "Pulling exercise to build a wider back profile."},
		{Name: "Squat", Category: "Strength", MuscleGroup: "Legs", Description: "Compound movement for leg power."},
		{Name: "Plank", Category: "Flexibility", MuscleGroup: "Core", Description: "Incline bench press for core strength."},
		{Name: "Jump Rope", Category: "Cardio", MuscleGroup: "Full Body", Description: "Light aerobic exercise for cardiovascular health."},
		{Name: "Crunches", Category: "Flexibility", MuscleGroup: "Core", Description: "Bodyweight exercise for core strength."},
		{Name: "Rowing Machine", Category: "Cardio", MuscleGroup: "Legs", Description: "Cardio exercise for leg power and cardiovascular health."},
	}

	tx, err := db.Begin()
	if err != nil {
		return err
	}

	// Example using PostgreSQL/SQLite syntax
	query := `INSERT INTO exercises (name, category, muscle_group, description) 
              VALUES ($1, $2, $3, $4) ON CONFLICT (name) DO NOTHING`

	stmt, err := tx.Prepare(query)
	if err != nil {
		tx.Rollback()
		return err
	}
	defer stmt.Close()

	for _, ex := range exercises {
		if _, err := stmt.Exec(ex.Name, ex.Category, ex.MuscleGroup, ex.Description); err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}
