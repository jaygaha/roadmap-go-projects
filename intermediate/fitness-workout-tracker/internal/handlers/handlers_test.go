package handlers

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jaygaha/roadmap-go-projects/intermediate/fitness-workout-tracker/internal/database"
	"github.com/jaygaha/roadmap-go-projects/intermediate/fitness-workout-tracker/internal/models"
)

func setupTestServer(t *testing.T) (*sql.DB, http.Handler, string) {
	t.Helper()
	dbPath := filepath.Join(t.TempDir(), "test.db")
	db, err := database.InitDB(dbPath)
	if err != nil {
		t.Fatalf("init db: %v", err)
	}

	jwtSecret := "testsecret"

	// Build a mux like SetupRoutes, but without env lookups
	root := http.NewServeMux()
	root.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"message":"Fitness Tracker API is up and running!"}`))
	})

	v1 := http.NewServeMux()
	v1.HandleFunc("POST /auth/register", RegisterUser(db))
	v1.HandleFunc("POST /auth/login", LoginUser(db, jwtSecret))
	auth := AuthMiddleware(jwtSecret)
	v1.Handle("GET /exercises", auth(GetExercises(db)))
	v1.Handle("POST /workouts", auth(CreateWorkout(db)))
	v1.Handle("GET /workouts", auth(ListWorkouts(db)))
	v1.Handle("PUT /workouts/{id}", auth(UpdateWorkout(db)))
	v1.Handle("GET /workouts/{id}", auth(GetWorkout(db)))
	v1.Handle("DELETE /workouts/{id}", auth(DeleteWorkout(db)))
	v1.Handle("GET /workouts/reports", auth(GetWorkoutReport(db)))
	root.Handle("/api/v1/", JSONMiddleware(http.StripPrefix("/api/v1", v1)))

	return db, root, jwtSecret
}

func doRequest(t *testing.T, handler http.Handler, method, path string, body any, headers map[string]string) *httptest.ResponseRecorder {
	t.Helper()
	var buf bytes.Buffer
	if body != nil {
		if err := json.NewEncoder(&buf).Encode(body); err != nil {
			t.Fatalf("encode body: %v", err)
		}
	}
	req := httptest.NewRequest(method, path, &buf)
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	return rr
}

func registerAndLogin(t *testing.T, h http.Handler) string {
	t.Helper()
	reg := map[string]string{
		"username": "alice",
		"email":    "alice@example.com",
		"password": "Secret123!",
	}
	rr := doRequest(t, h, http.MethodPost, "/api/v1/auth/register", reg, nil)
	if rr.Code != http.StatusCreated {
		t.Fatalf("register status = %d, body=%s", rr.Code, rr.Body.String())
	}
	login := map[string]string{
		"email":    "alice@example.com",
		"password": "Secret123!",
	}
	rr = doRequest(t, h, http.MethodPost, "/api/v1/auth/login", login, nil)
	if rr.Code != http.StatusOK {
		t.Fatalf("login status = %d, body=%s", rr.Code, rr.Body.String())
	}
	var resp map[string]string
	_ = json.Unmarshal(rr.Body.Bytes(), &resp)
	return resp["token"]
}

func TestRootHealth(t *testing.T) {
	_, h, _ := setupTestServer(t)
	rr := doRequest(t, h, http.MethodGet, "/", nil, nil)
	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}
}

func TestAuthRegisterAndLogin(t *testing.T) {
	_, h, _ := setupTestServer(t)
	token := registerAndLogin(t, h)
	if token == "" {
		t.Fatalf("expected token, got empty")
	}
}

func TestExercisesProtected(t *testing.T) {
	_, h, _ := setupTestServer(t)
	// Unauthorized
	rr := doRequest(t, h, http.MethodGet, "/api/v1/exercises", nil, nil)
	if rr.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", rr.Code)
	}
	// Authorized
	token := registerAndLogin(t, h)
	headers := map[string]string{"Authorization": "Bearer " + token}
	rr = doRequest(t, h, http.MethodGet, "/api/v1/exercises", nil, headers)
	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}
}

func TestWorkoutCRUDAndReport(t *testing.T) {
	db, h, _ := setupTestServer(t)
	token := registerAndLogin(t, h)
	headers := map[string]string{"Authorization": "Bearer " + token}

	// Create workout
	create := models.CreateWorkoutRequest{
		Name:         "Leg Day",
		Description:  "Heavy squats",
		ScheduledFor: time.Now().Add(24 * time.Hour).UTC(),
		Exercises: []models.WorkoutExerciseRequest{
			{ExerciseID: 1, Sets: 5, Reps: 5, Weight: 100, Notes: "Warmup included"},
		},
	}
	rr := doRequest(t, h, http.MethodPost, "/api/v1/workouts", create, headers)
	if rr.Code != http.StatusCreated {
		t.Fatalf("create workout expected 201, got %d, body=%s", rr.Code, rr.Body.String())
	}

	// Query last workout id (for tests)
	var wid int
	if err := db.QueryRow("SELECT id FROM workouts ORDER BY id DESC LIMIT 1").Scan(&wid); err != nil {
		t.Fatalf("get workout id: %v", err)
	}

	// Get workout
	rr = doRequest(t, h, http.MethodGet, "/api/v1/workouts/"+intToPath(wid), nil, headers)
	if rr.Code != http.StatusOK {
		t.Fatalf("get workout expected 200, got %d", rr.Code)
	}

	// List workouts
	rr = doRequest(t, h, http.MethodGet, "/api/v1/workouts", nil, headers)
	if rr.Code != http.StatusOK {
		t.Fatalf("list workouts expected 200, got %d", rr.Code)
	}

	// Update workout
	update := models.CreateWorkoutRequest{
		Name:         "Leg Day Updated",
		Description:  "Add accessory",
		ScheduledFor: time.Now().Add(48 * time.Hour).UTC(),
		Exercises: []models.WorkoutExerciseRequest{
			{ExerciseID: 1, Sets: 4, Reps: 8, Weight: 90, Notes: "Backoff sets"},
		},
	}
	rr = doRequest(t, h, http.MethodPut, "/api/v1/workouts/"+intToPath(wid), update, headers)
	if rr.Code != http.StatusOK {
		t.Fatalf("update workout expected 200, got %d, body=%s", rr.Code, rr.Body.String())
	}

	// Report
	start := time.Now().Add(-24 * time.Hour).Format("2006-01-02")
	end := time.Now().Add(72 * time.Hour).Format("2006-01-02")
	rr = doRequest(t, h, http.MethodGet, "/api/v1/workouts/reports?start_date="+start+"&end_date="+end, nil, headers)
	if rr.Code != http.StatusOK {
		t.Fatalf("report expected 200, got %d, body=%s", rr.Code, rr.Body.String())
	}

	// Delete workout
	rr = doRequest(t, h, http.MethodDelete, "/api/v1/workouts/"+intToPath(wid), nil, headers)
	if rr.Code != http.StatusNoContent {
		t.Fatalf("delete workout expected 204, got %d, body=%s", rr.Code, rr.Body.String())
	}
}

func intToPath(id int) string {
	// simple helper to avoid fmt import
	digits := []byte{}
	if id == 0 {
		return "0"
	}
	for id > 0 {
		d := byte('0' + id%10)
		digits = append([]byte{d}, digits...)
		id /= 10
	}
	return string(digits)
}

func TestJWTMiddlewareRejectsInvalidToken(t *testing.T) {
	_, h, _ := setupTestServer(t)
	// token signed with wrong secret
	claims := &models.Claims{
		UserID:   1,
		Username: "evil",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Hour)),
		},
	}
	bad, _ := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte("wrong"))
	headers := map[string]string{"Authorization": "Bearer " + bad}
	rr := doRequest(t, h, http.MethodGet, "/api/v1/exercises", nil, headers)
	if rr.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401 for invalid token, got %d", rr.Code)
	}
}
