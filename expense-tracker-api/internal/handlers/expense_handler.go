package handlers

import (
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/jaygaha/roadmap-go-projects/expense-tracker-api/config"
	"github.com/jaygaha/roadmap-go-projects/expense-tracker-api/internal/middleware"
	"github.com/jaygaha/roadmap-go-projects/expense-tracker-api/internal/models"
	"github.com/jaygaha/roadmap-go-projects/expense-tracker-api/internal/response"
	"github.com/jaygaha/roadmap-go-projects/expense-tracker-api/internal/services"
)

// ExpenseRequest represents the request body for the expense endpoint
type ExpenseRequest struct {
	Title             string  `json:"title" validate:"required"`
	Amount            float64 `json:"amount" validate:"required"`
	ExpenseCategoryID uint    `json:"expense_category_id" validate:"required"`
}

// ValidateExpenseRequest validates the expense request
func (e *ExpenseRequest) ValidateExpenseRequest() error {
	if strings.TrimSpace(e.Title) == "" {
		return errors.New("title is required")
	}
	if e.Amount <= 0 {
		return errors.New("amount must be greater than zero")
	}
	if e.ExpenseCategoryID == 0 {
		return errors.New("expense category ID is required")
	}

	return nil
}

// CreateExpenseHandler creates a new expense
func CreateExpensehandler(cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req ExpenseRequest

		if err := response.JSON(r, &req); err != nil {
			response.Error(w, http.StatusBadRequest, "Invalid request body")
			return
		}

		if err := req.ValidateExpenseRequest(); err != nil {
			response.Error(w, http.StatusUnprocessableEntity, err.Error())
			return
		}

		// get user ID from context
		userID, ok := r.Context().Value(middleware.UserIDKey).(uint)

		if !ok {
			response.Error(w, http.StatusInternalServerError, "Failed to get user ID from context")
			return
		}

		// create expense
		expense := &models.Expense{
			Title:             req.Title,
			Amount:            req.Amount,
			UserID:            userID,
			ExpenseCategoryID: req.ExpenseCategoryID,
		}

		if err := services.CreateExpense(expense); err != nil {
			response.Error(w, http.StatusInternalServerError, "Failed to create expense")
			return
		}

		response.Success(w, http.StatusCreated, expense)
	}
}

// ListExpensesHandler lists all expenses
func ListExpensesHandler(cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// get user ID from context
		userID, ok := r.Context().Value(middleware.UserIDKey).(uint)
		if !ok {
			response.Error(w, http.StatusInternalServerError, "Failed to get user ID from context")
			return
		}

		// Get filter parameters from query string
		filter := r.URL.Query().Get("filter")
		startDateStr := r.URL.Query().Get("start_date")
		endDateStr := r.URL.Query().Get("end_date")

		var startDate, endDate *time.Time // it can be null
		now := time.Now()

		switch filter {
		case "past_week":
			s := now.AddDate(0, 0, -7)
			startDate = &s
			endDate = &now
		case "past_month":
			s := now.AddDate(0, -1, 0)
			startDate = &s
			endDate = &now
		case "last_3_months":
			s := now.AddDate(0, -3, 0)
			startDate = &s
			endDate = &now
		case "custom":
			if startDateStr == "" || endDateStr == "" {
				response.Error(w, http.StatusBadRequest, "start_date and end_date are required for custom filter")
				return
			}
			dateLayout := "2006-01-02"
			startDateParsed, err1 := time.Parse(dateLayout, startDateStr)
			endDateParsed, err2 := time.Parse(dateLayout, endDateStr)
			if err1 != nil || err2 != nil {
				response.Error(w, http.StatusBadRequest, "Invalid date format. Use YYYY-MM-DD")
				return
			}

			startDate = &startDateParsed
			endDate = &endDateParsed
		}

		// list expenses
		expenses, err := services.ListExpensesWithFilter(userID, startDate, endDate, cfg)
		if err != nil {
			response.Error(w, http.StatusInternalServerError, "Failed to list expenses")
			return
		}

		response.Success(w, http.StatusOK, expenses)
	}
}

// GetExpenseByIDHandler retrieves an expense by its ID
func GetExpenseByIDHandler(cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// get user ID from context
		userID, ok := r.Context().Value(middleware.UserIDKey).(uint)
		if !ok {
			response.Error(w, http.StatusInternalServerError, "Failed to get user ID from context")
			return
		}

		// get expense ID from URL parameter
		idParam := mux.Vars(r)["id"]
		id, err := strconv.Atoi(idParam)
		if err != nil {
			response.Error(w, http.StatusBadRequest, "Invalid expense ID")
			return
		}

		// get expense
		expense, err := services.GetExpenseByID(uint(id), userID, cfg)
		if err != nil {
			response.Error(w, http.StatusNotFound, "Expense not found")
			return
		}

		response.Success(w, http.StatusOK, expense)
	}
}

// UpdateExpenseHandler updates an expense
// Only owner can update the expense
func UpdateExpenseHandler(cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// get user ID from context
		userID, ok := r.Context().Value(middleware.UserIDKey).(uint)
		if !ok {
			response.Error(w, http.StatusInternalServerError, "Failed to get user ID from context")
			return
		}
		// get expense ID from URL parameter
		idParam := mux.Vars(r)["id"]
		id, err := strconv.Atoi(idParam)
		if err != nil {
			response.Error(w, http.StatusBadRequest, "Invalid expense ID")
			return
		}
		// get expense
		expense, err := services.GetExpenseByID(uint(id), userID, cfg)
		if err != nil {
			response.Error(w, http.StatusNotFound, "Expense not found")
			return
		}

		// parse request body
		var req ExpenseRequest
		if err := response.JSON(r, &req); err != nil {
			response.Error(w, http.StatusBadRequest, "Invalid request body")
			return
		}

		// validate request
		if err := req.ValidateExpenseRequest(); err != nil {
			response.Error(w, http.StatusUnprocessableEntity, err.Error())
			return
		}

		// update expense
		expense.Title = req.Title
		expense.Amount = req.Amount
		expense.ExpenseCategoryID = req.ExpenseCategoryID

		if err := services.UpdateExpense(expense, cfg); err != nil {
			response.Error(w, http.StatusInternalServerError, "Failed to update expense")
			return
		}

		response.Success(w, http.StatusOK, expense)
	}
}

// DeleteExpenseHandler deletes an expense
// Only owner can delete the expense
func DeleteExpenseHandler(cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// get user ID from context
		userID, ok := r.Context().Value(middleware.UserIDKey).(uint)
		if !ok {
			response.Error(w, http.StatusInternalServerError, "Failed to get user ID from context")
			return
		}
		// get expense ID from URL parameter
		idParam := mux.Vars(r)["id"]
		id, err := strconv.Atoi(idParam)
		if err != nil {
			response.Error(w, http.StatusBadRequest, "Invalid expense ID")
			return
		}

		// get expense
		expense, err := services.GetExpenseByID(uint(id), userID, cfg)
		if err != nil {
			response.Error(w, http.StatusNotFound, "Expense not found")
			return
		}

		// delete expense
		if err := services.DeleteExpense(expense, cfg); err != nil {
			response.Error(w, http.StatusInternalServerError, "Failed to delete expense")
			return
		}

		response.Success(w, http.StatusOK, map[string]string{"message": "Expense deleted successfully"})
	}
}
