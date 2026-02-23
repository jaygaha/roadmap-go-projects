package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/jaygaha/roadmap-go-projects/intermediate/fitness-workout-tracker/internal/models"
)

// GetExercises godoc
//
//	@ID				getExercises
//	@Summary		Get all exercises
//	@Description	Retrieve a list of all available exercises
//	@Security		BearerToken
//	@Tags			exercises
//	@Accept			json
//	@Produce		json
//	@Success		200	{array}		models.Exercise		"List of exercises"
//	@Failure		401	{object}	map[string]string	"Unauthorized"
//	@Failure		500	{object}	map[string]string	"Internal server error"
//	@Security		BearerAuth
//	@Router			/exercises [get]
func GetExercises(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rows, err := db.Query("SELECT id, name, category, muscle_group, description FROM exercises")
		if err != nil {
			http.Error(w, `{"error": "Failed to fetch exercises"}`, http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var exercises []models.Exercise
		for rows.Next() {
			var ex models.Exercise
			if err := rows.Scan(&ex.ID, &ex.Name, &ex.Category, &ex.MuscleGroup, &ex.Description); err != nil {
				http.Error(w, `{"error": "Error scanning exercises"}`, http.StatusInternalServerError)
				return
			}
			exercises = append(exercises, ex)
		}

		json.NewEncoder(w).Encode(exercises)
	}
}
