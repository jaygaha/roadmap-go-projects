package handlers

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/jaygaha/roadmap-go-projects/expense-tracker-api/config"
	"github.com/jaygaha/roadmap-go-projects/expense-tracker-api/internal/middleware"
	"github.com/jaygaha/roadmap-go-projects/expense-tracker-api/internal/request"
	"github.com/jaygaha/roadmap-go-projects/expense-tracker-api/internal/response"
	"github.com/jaygaha/roadmap-go-projects/expense-tracker-api/internal/services"
)

// CreateCategoryRequest represents the request body for the create category endpoint
type CreateCategoryRequest struct {
	Name string `json:"name" validate:"required"`
}

// CreateCategoryHandler handles POST requests to create a new expense category
func CreateCategoryHandler(cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req CreateCategoryRequest

		if err := response.JSON(r, &req); err != nil {
			response.Error(w, http.StatusBadRequest, "Invalid request body")
			return
		}

		// Validate the request
		if err := request.ValidateStruct(req); err != nil {
			response.Error(w, http.StatusUnprocessableEntity, err.Error())
			return
		}

		// check if category already exists
		_, err := services.GetCategoryByName(req.Name, cfg)
		if err == nil {
			response.Error(w, http.StatusUnprocessableEntity, "Category already exists")
			return
		}

		userID := r.Context().Value(middleware.UserIDKey).(uint)

		// Create the category
		category, err := services.CreateExpenseCategory(req.Name, userID, cfg)
		if err != nil {
			response.Error(w, http.StatusInternalServerError, "Failed to create category")
			return
		}

		// Return the created category in the response
		response.Success(w, http.StatusCreated, category)
	}
}

// ListCategoriesHandler handles GET requests to retrieve a list of expense categories
func ListCategoriesHandler(cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Retrieve the list of categories from the repository
		categories, err := services.ListExpenseCategories(cfg)
		if err != nil {
			response.Error(w, http.StatusInternalServerError, "Failed to retrieve categories")
			return
		}

		// Return the list of categories in the response
		response.Success(w, http.StatusOK, categories)
	}
}

// GetCategoryByIDHandler handles GET requests to retrieve a category by ID
func GetCategoryByIDHandler(cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idParam := mux.Vars(r)["id"] // GET /api/v1/expenses/categories/{id}
		// Convert the ID parameter to an integer
		categoryID, err := strconv.Atoi(idParam)
		if err != nil {
			response.Error(w, http.StatusBadRequest, "Invalid category ID")
			return
		}

		// Retrieve the category by ID from the repository
		category, err := services.GetCategoryByID(uint(categoryID), cfg)
		if err != nil {
			response.Error(w, http.StatusNotFound, "Category not found")
			return
		}

		// Return the category in the response
		response.Success(w, http.StatusOK, category)
	}
}

// UpdateCategoryHandler handles PUT requests to update an existing category
func UpdateCategoryHandler(cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idParam := mux.Vars(r)["id"] // PUT /api/v1/expenses/categories/{id}
		// Convert the ID parameter to an integer
		categoryID, err := strconv.Atoi(idParam)
		if err != nil {
			response.Error(w, http.StatusBadRequest, "Invalid category ID")
			return
		}

		uCategoryID := uint(categoryID)

		// Retrieve the category by ID from the repository
		_, err = services.GetCategoryByID(uCategoryID, cfg)
		if err != nil {
			response.Error(w, http.StatusNotFound, "Category not found")
			return
		}

		// Parse the request body
		var req CreateCategoryRequest
		if err = response.JSON(r, &req); err != nil {
			response.Error(w, http.StatusBadRequest, "Invalid request body")
			return
		}

		// Validate the request
		if err = request.ValidateStruct(req); err != nil {
			response.Error(w, http.StatusUnprocessableEntity, err.Error())
			return
		}
		// check if category already exists
		_, err = services.GetCategoryByName(req.Name, cfg)
		if err == nil {
			response.Error(w, http.StatusUnprocessableEntity, "Category name already exists")
			return
		}

		// Simulate creator — later from JWT claims
		userID := uint(1)

		// Save the updated category
		updatedCategory, err := services.UpdateExpenseCategory(uCategoryID, req.Name, userID, cfg)
		if err != nil {
			response.Error(w, http.StatusInternalServerError, "Failed to update category")
			return
		}

		// Return the updated category in the response
		response.Success(w, http.StatusOK, updatedCategory)
	}
}

// DeleteCategoryHandler handles DELETE requests to delete a category
func DeleteCategoryHandler(cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idParam := mux.Vars(r)["id"] // DELETE /api/v1/expenses/categories/{id}
		// Convert the ID parameter to an integer
		categoryID, err := strconv.Atoi(idParam)
		if err != nil {
			response.Error(w, http.StatusBadRequest, "Invalid category ID")
			return
		}
		uCategoryID := uint(categoryID)
		// Retrieve the category by ID from the repository
		_, err = services.GetCategoryByID(uCategoryID, cfg)
		if err != nil {
			response.Error(w, http.StatusNotFound, "Category not found")
			return
		}

		// Check if the category has any associated expenses
		hasExpenses, err := services.HasExpensesForCategory(uCategoryID, cfg)
		if err != nil {
			response.Error(w, http.StatusInternalServerError, "Failed to check expenses")
			return
		}

		if hasExpenses {
			response.Error(w, http.StatusUnprocessableEntity, "Category has associated expenses")
			return
		}

		userID := r.Context().Value(middleware.UserIDKey).(uint)

		// Delete the category
		err = services.DeleteExpenseCategory(uCategoryID, userID, cfg)
		if err != nil {
			response.Error(w, http.StatusInternalServerError, "Failed to delete category")
			return
		}

		// Return a success response
		response.Success(w, http.StatusOK, nil)
	}
}
