package models

import "time"

type Exercise struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Category    string    `json:"category"`    // e.g., Strength, Cardio, Flexibility
	MuscleGroup string    `json:"muscleGroup"` // e.g., Chest, Legs, Back
	Description string    `json:"description"`
	UpdatedAt   time.Time `json:"updated_at" bson:"updated_at"`
	CreatedAt   time.Time `json:"created_at" bson:"created_at"`
}
