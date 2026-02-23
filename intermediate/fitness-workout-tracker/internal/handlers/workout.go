package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/jaygaha/roadmap-go-projects/intermediate/fitness-workout-tracker/internal/models"
)

// CreateWorkout handles POST /workouts
// CreateWorkout godoc
//
//	@ID				createWorkout
//	@Summary		Create a new workout
//	@Description	Create a new workout for the authenticated user
//	@Tags			workouts
//	@Accept			json
//	@Produce		json
//	@Param			createWorkoutRequest	body		models.CreateWorkoutRequest	true	"Create Workout Request"
//	@Success		201	{object}	map[string]string	"Workout created successfully"
//	@Failure		401	{object}	map[string]string	"Unauthorized"
//	@Failure		400	{object}	map[string]string	"Invalid request"
//	@Failure		422	{object}	map[string]string	"Validation errors"
//	@Failure		500	{object}	map[string]string	"Internal server error"
//	@Security		BearerAuth
//	@Router			/workouts [post]
func CreateWorkout(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 1. Get UserID from Context
		userID := r.Context().Value(userIDKey).(int)

		// 2. Decode the Request Body
		var req models.CreateWorkoutRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, `{"error": "Invalid request"}`, http.StatusBadRequest)
			return
		}

		// 3. Validate the Request Payload
		if err := req.Validate(); err != nil {
			http.Error(w, `{"error": "Invalid request"}`, http.StatusUnprocessableEntity)
			return
		}

		// 4. Start Transaction
		tx, err := db.Begin()
		if err != nil {
			http.Error(w, `{"error": "Database error"}`, http.StatusInternalServerError)
			return
		}
		defer tx.Rollback() // Ensure cleanup on error

		// 5. Insert the Parent Workout
		res, err := tx.Exec(`INSERT INTO workouts (user_id, name, scheduled_for, description) VALUES (?, ?, ?, ?)`,
			userID, req.Name, req.ScheduledFor, req.Description)
		if err != nil {
			http.Error(w, `{"error": "Failed to create workout"}`, http.StatusInternalServerError)
			return
		}

		workoutID, _ := res.LastInsertId()

		// 6. Insert the Child Exercises
		for _, ex := range req.Exercises {
			_, err := tx.Exec(`INSERT INTO workout_exercises (workout_id, exercise_id, sets, reps, weight, notes) VALUES (?, ?, ?, ?, ?, ?)`,
				workoutID, ex.ExerciseID, ex.Sets, ex.Reps, ex.Weight, ex.Notes)
			if err != nil {
				http.Error(w, `{"error": "Failed to add exercises"}`, http.StatusInternalServerError)
				return
			}
		}

		// 7. Commit the Transaction
		if err := tx.Commit(); err != nil {
			http.Error(w, `{"error": "Failed to save workout"}`, http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]string{"message": "Workout created successfully"})
	}
}

// ListWorkouts handles GET /workouts
// ListWorkouts godoc
//
//	@ID				listWorkouts
//	@Summary		List all workouts
//	@Description	Retrieve a list of all workouts for the authenticated user
//	@Tags			workouts
//	@Accept			json
//	@Produce		json
//	@Success		200	{array}		models.WorkoutResponse		"List of workouts"
//	@Failure		401	{object}	map[string]string			"Unauthorized"
//	@Failure		500	{object}	map[string]string	"Internal server error"
//	@Security		BearerAuth
//	@Router			/workouts [get]
func ListWorkouts(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := r.Context().Value(userIDKey).(int)

		query := `
			SELECT w.id, w.name, w.scheduled_for, w.description,
			       we.exercise_id, we.sets, we.reps, we.weight, we.notes,
			       e.name, e.category, e.muscle_group, e.description
			FROM workouts w
			LEFT JOIN workout_exercises we ON w.id = we.workout_id
			LEFT JOIN exercises e ON we.exercise_id = e.id
			WHERE w.user_id = ?
			ORDER BY w.scheduled_for DESC`

		rows, err := db.Query(query, userID)
		if err != nil {
			log.Printf("Error querying workouts: %v", err)
			http.Error(w, `{"error": "Failed to fetch workouts"}`, http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		// Tracking map and ordered slice
		workoutMap := make(map[int]*models.WorkoutResponse)
		var result []*models.WorkoutResponse

		var eName, eCat, eMuscle, eDesc sql.NullString

		for rows.Next() {
			var workoutID int
			var name string
			var scheduledFor time.Time
			var description *string
			// Using pointers to handle the NULLs we discussed!
			var exID, sets, reps *int
			var weight *float64
			var notes *string

			err := rows.Scan(&workoutID, &name, &scheduledFor, &description, &exID, &sets, &reps, &weight, &notes, &eName, &eCat, &eMuscle, &eDesc)
			if err != nil {
				continue
			}

			// 1. If this is a new workout, initialize it
			if _, exists := workoutMap[workoutID]; !exists {
				// Initialize description with empty string if nil
				desc := ""
				if description != nil {
					desc = *description
				}
				newWorkout := &models.WorkoutResponse{
					ID:           workoutID,
					Name:         name,
					ScheduledFor: scheduledFor,
					Description:  desc, // Initialize with empty string
					Exercises:    []models.WorkoutExerciseResponse{},
				}

				// Store the pointer in the map
				workoutMap[workoutID] = newWorkout

				// Append the SAME pointer to the slice to preserve order and sync
				result = append(result, newWorkout)
			}

			// Now, appending here updates the object in both the map AND the slice
			if exID != nil {
				workoutMap[workoutID].Exercises = append(workoutMap[workoutID].Exercises, models.WorkoutExerciseResponse{
					ExerciseID: *exID,
					Sets:       *sets,
					Reps:       *reps,
					Weight:     *weight,
					Notes:      *notes,
					Exercise: models.Exercise{
						ID:          *exID,
						Name:        eName.String,
						Category:    eCat.String,
						MuscleGroup: eMuscle.String,
						Description: eDesc.String,
					},
				})
			}
		}

		if len(result) == 0 {
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode([]models.WorkoutResponse{})
			return
		}

		json.NewEncoder(w).Encode(result)
	}
}

// UpdateWorkout handles PUT /workouts/{id}
// UpdateWorkout godoc
//
//	@ID				updateWorkout
//	@Summary		Update a workout
//	@Description	Update the details of a specific workout for the authenticated user
//	@Tags			workouts
//	@Accept			json
//	@Produce		json
//	@Param			id				path		int							true	"Workout ID"
//	@Param			updateWorkoutRequest	body		models.CreateWorkoutRequest	true	"Update Workout Request"
//	@Success		200	{object}	map[string]string	"Workout updated successfully"
//	@Failure		401	{object}	map[string]string	"Unauthorized"
//	@Failure		400	{object}	map[string]string	"Invalid request"
//	@Failure		422	{object}	map[string]string	"Validation errors"
//	@Failure		500	{object}	map[string]string	"Internal server error"
//	@Security		BearerAuth
//	@Router			/workouts/{id} [put]
func UpdateWorkout(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 1. Extract Workout ID from URL and User ID from Context
		idStr := r.PathValue("id")
		workoutID, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, `{"error": "Invalid workout ID"}`, http.StatusBadRequest)
			return
		}
		userID := r.Context().Value(userIDKey).(int)

		// 2. Decode and Validate the Request Body
		var req models.CreateWorkoutRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, `{"error": "Invalid request payload"}`, http.StatusBadRequest)
			return
		}
		if err := req.Validate(); err != nil {
			http.Error(w, fmt.Sprintf(`{"error": "%s"}`, err.Error()), http.StatusBadRequest)
			return
		}

		// 3. Start Transaction
		tx, err := db.Begin()
		if err != nil {
			http.Error(w, `{"error": "Database error"}`, http.StatusInternalServerError)
			return
		}
		defer tx.Rollback()

		// 4. Update Workout Metadata (Ownership Check Included)
		res, err := tx.Exec(`
			UPDATE workouts 
			SET name = ?, scheduled_for = ?, description = ?, updated_at = ?
			WHERE id = ? AND user_id = ?`,
			req.Name, req.ScheduledFor, req.Description, time.Now(), workoutID, userID)

		// Debug req
		log.Printf("UpdateWorkout req: %v", req)
		if err != nil {
			http.Error(w, `{"error": "Failed to update workout"}`, http.StatusInternalServerError)
			return
		}

		rowsAffected, _ := res.RowsAffected()
		if rowsAffected == 0 {
			http.Error(w, `{"error": "Workout not found or unauthorized"}`, http.StatusNotFound)
			return
		}

		// 5. Delete Old Exercises (The "Clear" step)
		_, err = tx.Exec(`DELETE FROM workout_exercises WHERE workout_id = ?`, workoutID)
		if err != nil {
			http.Error(w, `{"error": "Failed to clear existing exercises"}`, http.StatusInternalServerError)
			return
		}

		// 6. Insert New Exercises (The "Re-insert" step)
		for _, ex := range req.Exercises {
			_, err := tx.Exec(`
				INSERT INTO workout_exercises (workout_id, exercise_id, sets, reps, weight, notes) 
				VALUES (?, ?, ?, ?, ?, ?)`,
				workoutID, ex.ExerciseID, ex.Sets, ex.Reps, ex.Weight, ex.Notes)
			if err != nil {
				http.Error(w, `{"error": "Failed to add exercises"}`, http.StatusInternalServerError)
				return
			}
		}

		// 7. Commit changes
		if err := tx.Commit(); err != nil {
			http.Error(w, `{"error": "Failed to update workout"}`, http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"message": "Workout updated successfully"})
	}
}

// GetWorkout godoc
//
//	@ID				getWorkout
//	@Summary		Get a workout
//	@Description	Retrieve the details of a specific workout for the authenticated user
//	@Tags			workouts
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int							true	"Workout ID"
//	@Success		200	{object}	models.WorkoutResponse	"Workout details"
//	@Failure		401	{object}	map[string]string	"Unauthorized"
//	@Failure		400	{object}	map[string]string	"Invalid request"
//	@Failure		404	{object}	map[string]string	"Workout not found"
//	@Failure		500	{object}	map[string]string	"Internal server error"
//	@Security		BearerAuth
//	@Router			/workouts/{id} [get]
func GetWorkout(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		workoutID, err := strconv.Atoi(r.PathValue("id"))
		if err != nil {
			http.Error(w, `{"error": "Invalid workout ID"}`, http.StatusBadRequest)
			return
		}

		userID, ok := r.Context().Value(userIDKey).(int)
		if !ok {
			http.Error(w, `{"error": "Unauthorized"}`, http.StatusUnauthorized)
			return
		}

		rows, err := db.Query(`
			SELECT w.id, w.name, w.scheduled_for, w.description, 
			       e.id, e.name, e.category, e.muscle_group, e.description,
			       we.sets, we.reps, we.weight, we.notes
			FROM workouts w
			LEFT JOIN workout_exercises we ON w.id = we.workout_id
			LEFT JOIN exercises e ON we.exercise_id = e.id
			WHERE w.id = ? AND w.user_id = ?`,
			workoutID, userID)
		if err != nil {
			http.Error(w, `{"error": "Database error"}`, http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var workout models.WorkoutResponse
		initialized := false

		for rows.Next() {
			// Use Null types to safely handle the LEFT JOIN results
			var eID, sets, reps sql.NullInt64
			var eName, eCat, eMuscle, eDesc, weNotes sql.NullString
			var weight sql.NullFloat64

			err := rows.Scan(
				&workout.ID, &workout.Name, &workout.ScheduledFor, &workout.Description,
				&eID, &eName, &eCat, &eMuscle, &eDesc,
				&sets, &reps, &weight, &weNotes,
			)
			if err != nil {
				http.Error(w, `{"error": "Scanning error"}`, http.StatusInternalServerError)
				return
			}
			initialized = true

			// Only append an exercise if the joined row actually exists (ID is not NULL)
			if eID.Valid {
				exerciseDetail := models.Exercise{
					ID:          int(eID.Int64),
					Name:        eName.String,
					Category:    eCat.String,
					MuscleGroup: eMuscle.String,
					Description: eDesc.String,
				}

				workoutExercise := models.WorkoutExerciseResponse{
					ExerciseID: int(eID.Int64),
					Sets:       int(sets.Int64),
					Reps:       int(reps.Int64),
					Weight:     weight.Float64,
					Notes:      weNotes.String,
					Exercise:   exerciseDetail,
				}

				workout.Exercises = append(workout.Exercises, workoutExercise)
			}
		}

		if !initialized {
			http.Error(w, `{"error": "Workout not found"}`, http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(workout)
	}
}

// DeleteWorkout godoc
//
//	@ID				deleteWorkout
//	@Summary		Delete a workout
//	@Description	Delete a specific workout for the authenticated user
//	@Tags			workouts
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int							true	"Workout ID"
//	@Success		204	{string}	string						"Workout deleted successfully"
//	@Failure		401	{object}	map[string]string	"Unauthorized"
//	@Failure		400	{object}	map[string]string	"Invalid request"
//	@Failure		404	{object}	map[string]string	"Workout not found"
//	@Failure		500	{object}	map[string]string	"Internal server error"
//	@Security		BearerAuth
//	@Router			/workouts/{id} [delete]
func DeleteWorkout(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 1. Extract Workout ID from URL and User ID from Context
		idStr := r.PathValue("id")
		workoutID, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, `{"error": "Invalid workout ID"}`, http.StatusBadRequest)
			return
		}
		userID := r.Context().Value(userIDKey).(int)

		// 2. Start Transaction
		tx, err := db.Begin()
		if err != nil {
			http.Error(w, `{"error": "Database error"}`, http.StatusInternalServerError)
			return
		}
		defer tx.Rollback()

		// 3. Delete Exercises First (The Children)
		_, err = tx.Exec(`DELETE FROM workout_exercises WHERE workout_id = ?`, workoutID)
		if err != nil {
			http.Error(w, `{"error": "Failed to delete child records"}`, http.StatusInternalServerError)
			return
		}

		// 4. Delete Workout (The Parent - with Ownership Check)
		res, err := tx.Exec(`DELETE FROM workouts WHERE id = ? AND user_id = ?`, workoutID, userID)
		if err != nil {
			http.Error(w, `{"error": "Failed to delete workout"}`, http.StatusInternalServerError)
			return
		}

		// 5. Check if anything was actually deleted
		rowsAffected, _ := res.RowsAffected()
		if rowsAffected == 0 {
			http.Error(w, `{"error": "Workout not found or unauthorized"}`, http.StatusNotFound)
			return
		}

		// 6. Commit
		if err := tx.Commit(); err != nil {
			http.Error(w, `{"error": "Failed to delete"}`, http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent) // 204 No Content is standard for successful deletes
	}
}

// GetWorkoutReport handles GET /workouts/reports
//
//	@ID				getWorkoutReport
//	@Summary		Get workout report
//	@Description	Retrieve a report of workouts for the authenticated user within a specified date range
//	@Tags			workouts
//	@Accept			json
//	@Produce		json
//	@Param			start_date	query		string					true	"Start Date (YYYY-MM-DD)"
//	@Param			end_date	query		string					true	"End Date (YYYY-MM-DD)"
//	@Success		200			{object}	[]models.WorkoutReportItem	"Workout report"
//	@Failure		400			{object}	map[string]string	"Invalid request"
//	@Failure		401			{object}	map[string]string	"Unauthorized"
//	@Failure		500			{object}	map[string]string	"Internal server error"
//	@Security		BearerAuth
//	@Router			/workouts/reports [get]
func GetWorkoutReport(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, ok := r.Context().Value(userIDKey).(int)
		if !ok {
			http.Error(w, `{"error": "Unauthorized"}`, http.StatusUnauthorized)
			return
		}

		// 1. Improved Date Parsing
		startStr := r.URL.Query().Get("start_date")
		endStr := r.URL.Query().Get("end_date")
		if startStr == "" || endStr == "" {
			http.Error(w, `{"error": "start_date and end_date are required (YYYY-MM-DD)"}`, http.StatusBadRequest)
			return
		}

		start, err := time.Parse("2006-01-02", startStr)
		if err != nil {
			http.Error(w, `{"error": "invalid start date"}`, http.StatusBadRequest)
			return
		}
		end, err := time.Parse("2006-01-02", endStr)
		if err != nil {
			http.Error(w, `{"error": "invalid end date"}`, http.StatusBadRequest)
			return
		}

		// 2. Fetch Data
		query := `
			SELECT w.id, w.name, w.scheduled_for, we.exercise_id, we.sets, we.reps, we.weight
			FROM workouts w
			JOIN workout_exercises we ON w.id = we.workout_id
			WHERE w.user_id = ? AND w.scheduled_for BETWEEN ? AND ?
			ORDER BY w.scheduled_for ASC`

		rows, err := db.Query(query, userID, start, end)
		if err != nil {
			http.Error(w, `{"error": "Database error"}`, http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		// 3. Process with Pointers to maintain order and structure
		workoutMap := make(map[int]*models.WorkoutReportItem)
		var reportSlice []*models.WorkoutReportItem
		totalExercises := 0

		for rows.Next() {
			var wID, eID, sets, reps int
			var wName string
			var wDate time.Time
			var weight float64

			if err := rows.Scan(&wID, &wName, &wDate, &eID, &sets, &reps, &weight); err != nil {
				http.Error(w, `{"error": "Scan error"}`, http.StatusInternalServerError)
				return
			}

			workout, exists := workoutMap[wID]
			if !exists {
				workout = &models.WorkoutReportItem{
					ID:           wID,
					Name:         wName,
					ScheduledFor: wDate,
					Exercises:    []models.ReportExercise{},
				}
				workoutMap[wID] = workout
				reportSlice = append(reportSlice, workout) // Preserves SQL sort order
			}

			workout.Exercises = append(workout.Exercises, models.ReportExercise{
				ExerciseID: eID,
				Sets:       sets,
				Reps:       reps,
				Weight:     weight,
			})
			totalExercises++
		}

		// 4. Final Response Construction
		response := struct {
			StartDate      string                      `json:"start_date"`
			EndDate        string                      `json:"end_date"`
			TotalWorkouts  int                         `json:"total_workouts"`
			TotalExercises int                         `json:"total_exercises"`
			Workouts       []*models.WorkoutReportItem `json:"workouts"`
		}{
			StartDate:      startStr,
			EndDate:        endStr,
			TotalWorkouts:  len(reportSlice),
			TotalExercises: totalExercises,
			Workouts:       reportSlice,
		}

		json.NewEncoder(w).Encode(response)
	}
}
