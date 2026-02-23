package models

import (
	"testing"
	"time"
)

func TestCreateWorkoutRequest_Validate(t *testing.T) {
	valid := CreateWorkoutRequest{
		Name:         "Test",
		Description:  "Desc",
		ScheduledFor: time.Now(),
		Exercises: []WorkoutExerciseRequest{
			{ExerciseID: 1, Sets: 3, Reps: 10, Weight: 50, Notes: ""},
		},
	}
	if err := valid.Validate(); err != nil {
		t.Fatalf("expected valid request, got error: %v", err)
	}

	noName := valid
	noName.Name = ""
	if err := noName.Validate(); err == nil {
		t.Fatalf("expected error for missing name")
	}

	noExercises := valid
	noExercises.Exercises = []WorkoutExerciseRequest{}
	if err := noExercises.Validate(); err == nil {
		t.Fatalf("expected error for missing exercises")
	}

	badNumbers := valid
	badNumbers.Exercises[0] = WorkoutExerciseRequest{ExerciseID: 1, Sets: 0, Reps: -1, Weight: -5}
	if err := badNumbers.Validate(); err == nil {
		t.Fatalf("expected error for invalid sets/reps/weight")
	}
}
