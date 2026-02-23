package models

import (
	"fmt"
	"time"
)

// CreateWorkoutRequest represents the incoming JSON for a new workout
type CreateWorkoutRequest struct {
	Name         string                   `json:"name"`
	Description  string                   `json:"description"`
	ScheduledFor time.Time                `json:"scheduled_for"`
	Exercises    []WorkoutExerciseRequest `json:"exercises"` // The collection
}

// WorkoutExerciseRequest represents the specific details for one exercise in the workout
type WorkoutExerciseRequest struct {
	ExerciseID int     `json:"exercise_id"`
	Sets       int     `json:"sets"`
	Reps       int     `json:"reps"`
	Weight     float64 `json:"weight"`
	Notes      string  `json:"notes"`
}

type WorkoutExerciseResponse struct {
	ExerciseID int      `json:"exercise_id"`
	Sets       int      `json:"sets"`
	Reps       int      `json:"reps"`
	Weight     float64  `json:"weight"`
	Notes      string   `json:"notes"`
	Exercise   Exercise `json:"exercise"`
}

// WorkoutResponse represents the outgoing JSON for a workout
type WorkoutResponse struct {
	ID           int                       `json:"id"`
	Name         string                    `json:"name"`
	Description  string                    `json:"description"`
	ScheduledFor time.Time                 `json:"scheduled_for"`
	Exercises    []WorkoutExerciseResponse `json:"exercises"`
}

// Validate checks if the CreateWorkoutRequest is valid
func (r CreateWorkoutRequest) Validate() error {
	if r.Name == "" {
		return fmt.Errorf("workout name is required")
	}
	if len(r.Exercises) == 0 {
		return fmt.Errorf("at least one exercise is required")
	}
	for _, ex := range r.Exercises {
		if ex.Sets <= 0 || ex.Reps <= 0 || ex.Weight < 0 {
			return fmt.Errorf("invalid sets, reps, or weight values")
		}
	}
	return nil
}

type WorkoutReportItem struct {
	ID           int              `json:"id"`
	Name         string           `json:"name"`
	ScheduledFor time.Time        `json:"scheduled_for"`
	Exercises    []ReportExercise `json:"exercises"`
}

type ReportExercise struct {
	ExerciseID int     `json:"exercise_id"`
	Sets       int     `json:"sets"`
	Reps       int     `json:"reps"`
	Weight     float64 `json:"weight"`
}
