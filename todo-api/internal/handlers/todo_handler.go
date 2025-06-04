package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/jaygaha/roadmap-go-projects/todo-api/internal/db"
	"github.com/jaygaha/roadmap-go-projects/todo-api/internal/models"
	"github.com/jaygaha/roadmap-go-projects/todo-api/internal/utils"
)

// CreateTodoRequest represents a request to create a new todo.
type TodoRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	IsCompleted bool   `json:"is_completed"`
}

// CreateTodoHandler creates a new todo. POST /todos/create
func CreateTodoHandler(w http.ResponseWriter, r *http.Request) {
	// Parse the request body
	var req TodoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondJSON(w, http.StatusBadRequest, "Invalid request body", nil)
		return
	}

	req.Title = strings.TrimSpace(req.Title)
	req.Description = strings.TrimSpace(req.Description)

	if req.Title == "" {
		utils.RespondJSON(w, http.StatusUnprocessableEntity, "Title is required", nil)
		return
	}

	// Get the authenticated user ID from the request context
	userId := GetUserIDFromContext(r)
	if userId == 0 {
		utils.RespondJSON(w, http.StatusInternalServerError, "Failed to get user ID", nil)
		return
	}

	timestamp := time.Now()
	cratedAt := timestamp.Format("2006-01-02 15:04:05")
	updatedAt := timestamp.Format("2006-01-02 15:04:05")

	// Create a new todo
	row, err := db.DB.Exec("INSERT INTO todos (title, description, is_completed, user_id, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?)", req.Title, req.Description, req.IsCompleted, userId, cratedAt, updatedAt)
	if err != nil {
		utils.RespondJSON(w, http.StatusInternalServerError, "Failed to create todo", nil)
		return
	}

	todoId, err := row.LastInsertId()
	if err != nil {
		utils.RespondJSON(w, http.StatusInternalServerError, "Failed to get todo ID", nil)
		return
	}

	todo, err := getTodoById(todoId, userId)
	if err != nil {
		utils.RespondJSON(w, http.StatusInternalServerError, "Failed to get todo", nil)
		return
	}

	utils.RespondJSON(w, http.StatusCreated, "Todo created successfully", todo)
}

// ListTodosHandler lists all todos with extra features. GET /todos?limit=10&page=1&is_completed=true&sort=created_at&order=desc
func ListTodosHandler(w http.ResponseWriter, r *http.Request) {
	// Get the authenticated user ID from the request context
	userId := GetUserIDFromContext(r)
	if userId == 0 {
		utils.RespondJSON(w, http.StatusInternalServerError, "Failed to get user ID", nil)
		return
	}

	// Query parameters
	queryParams := r.URL.Query()
	limitStr := queryParams.Get("limit")
	pageStr := queryParams.Get("page")
	IsCompletedStr := queryParams.Get("is_completed")
	sortBy := queryParams.Get("sort") // title, created_at
	order := queryParams.Get("order") // "asc" or "desc"

	// Default values
	limit := 10
	page := 1

	if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
		limit = l
	}
	if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
		page = p
	}

	offset := (page - 1) * limit

	// Filter by user ID
	filterQuery := "user_id = ?"
	filterArgs := []any{userId}

	if IsCompletedStr != "" {
		if IsCompletedStr == "true" || IsCompletedStr == "false" {
			filterQuery += " AND is_completed =?"
			filterArgs = append(filterArgs, IsCompletedStr == "true")
		}
	}

	// Sorting
	validSortFields := map[string]bool{
		"title":      true,
		"created_at": true,
	}
	if !validSortFields[sortBy] {
		sortBy = "created_at"
	}

	if order != "asc" {
		order = "desc"
	}

	// Query todos
	query := "SELECT id, title, description, is_completed, created_at, updated_at FROM todos WHERE " + filterQuery + " ORDER BY " + sortBy + " " + order + " LIMIT ? OFFSET ?"
	filterArgs = append(filterArgs, limit, offset)

	rows, err := db.DB.Query(query, filterArgs...)
	if err != nil {
		utils.RespondJSON(w, http.StatusInternalServerError, "Failed to list todos", map[string]any{"error": err.Error()})
		return
	}
	defer rows.Close()

	var todos []models.Todo
	for rows.Next() {
		var todo models.Todo
		todo.UserID = userId
		err := rows.Scan(&todo.ID, &todo.Title, &todo.Description, &todo.IsCompleted, &todo.CreatedAt, &todo.UpdatedAt)

		if err != nil {
			utils.RespondJSON(w, http.StatusInternalServerError, "Failed to list todos", map[string]any{"error": err.Error()})
			return
		}

		todos = append(todos, todo)
	}

	utils.RespondJSON(w, http.StatusOK, "Todos listed successfully", map[string]any{
		"page":  page,
		"limit": limit,
		"count": len(todos),
		"todos": todos,
	})
}

// FetchTodoHandler fetches a todo by its ID. GET GET /todos/details?id=1
func FetchTodoHandler(w http.ResponseWriter, r *http.Request) {
	// Get the authenticated user ID from the request context
	userId := GetUserIDFromContext(r)
	if userId == 0 {
		utils.RespondJSON(w, http.StatusInternalServerError, "Failed to get user ID", nil)
		return
	}

	todoId, err := getIdFromRequest(r)
	if err != nil {
		utils.RespondJSON(w, http.StatusBadRequest, "Invalid todo ID", nil)
		return
	}

	todo, err := getTodoById(todoId, userId)
	if err != nil {
		utils.RespondJSON(w, http.StatusNotFound, "Failed to get todo", nil)
		return
	}

	utils.RespondJSON(w, http.StatusOK, "Todo fetched successfully", todo)
}

// UpdateTodoHandler updates a todo by its ID. PUT /todos/update?id=1
func UpdateTodoHandler(w http.ResponseWriter, r *http.Request) {
	// Get the authenticated user ID from the request context
	userId := GetUserIDFromContext(r)
	if userId == 0 {
		utils.RespondJSON(w, http.StatusInternalServerError, "Failed to get user ID", nil)
		return
	}

	// Query ID
	todoId, err := getIdFromRequest(r)
	if err != nil {
		utils.RespondJSON(w, http.StatusBadRequest, "Invalid todo ID", nil)
		return
	}

	// Parse the request body
	var req TodoRequest
	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondJSON(w, http.StatusBadRequest, "Invalid request body", nil)
		return
	}

	req.Title = strings.TrimSpace(req.Title)
	req.Description = strings.TrimSpace(req.Description)
	if req.Title == "" {
		utils.RespondJSON(w, http.StatusUnprocessableEntity, "Title is required", nil)
	}

	timestamp := time.Now()
	updatedAt := timestamp.Format("2006-01-02 15:04:05")

	res, err := db.DB.Exec("UPDATE todos SET title=?, description=?, is_completed=?, updated_at=? WHERE id=? AND user_id=?", req.Title, req.Description, req.IsCompleted, updatedAt, todoId, userId)
	if err != nil {
		utils.RespondJSON(w, http.StatusInternalServerError, "Failed to update todo", map[string]any{"error": err.Error()})
		return
	}
	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		utils.RespondJSON(w, http.StatusNotFound, "Nothing to update or Todo not found", nil)
		return
	}
	todo, err := getTodoById(todoId, userId)
	if err != nil {
		utils.RespondJSON(w, http.StatusInternalServerError, "Failed to get todo", nil)
		return
	}

	utils.RespondJSON(w, http.StatusOK, "Todo updated successfully", todo)
}

// DeleteTodoHandler deletes a todo by its ID. DELETE /todos/delete?id=1
func DeleteTodoHandler(w http.ResponseWriter, r *http.Request) {
	// Get the authenticated user ID from the request context
	userId := GetUserIDFromContext(r)
	if userId == 0 {
		utils.RespondJSON(w, http.StatusInternalServerError, "Failed to get user ID", nil)
		return
	}

	// Query ID
	todoId, err := getIdFromRequest(r)
	if err != nil {
		utils.RespondJSON(w, http.StatusBadRequest, "Invalid todo ID", nil)
		return
	}

	res, err := db.DB.Exec("DELETE FROM todos WHERE id=? AND user_id=?", todoId, userId)
	if err != nil {
		utils.RespondJSON(w, http.StatusInternalServerError, "Failed to delete todo", nil)
		return
	}
	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		utils.RespondJSON(w, http.StatusNotFound, "Todo not found", nil)
		return
	}

	utils.RespondJSON(w, http.StatusOK, "Todo deleted successfully", nil)
}

// getTodoById fetches a todo by its ID and owner.
func getTodoById(id, userId int64) (*models.Todo, error) {
	row := db.DB.QueryRow("SELECT id, title, description, is_completed, user_id, created_at, updated_at FROM todos WHERE id = ? AND user_id = ?", id, userId)

	var todo models.Todo
	err := row.Scan(&todo.ID, &todo.Title, &todo.Description, &todo.IsCompleted, &todo.UserID, &todo.CreatedAt, &todo.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return &todo, nil
}

// getIdFromRequest returns ID from the request URL.
func getIdFromRequest(r *http.Request) (int64, error) {
	idStr := r.URL.Path[len("/todos/"):]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return 0, err
	}

	return int64(id), nil
}
