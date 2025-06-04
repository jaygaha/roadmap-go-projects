package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/jaygaha/roadmap-go-projects/todo-api/internal/auth"
	"github.com/jaygaha/roadmap-go-projects/todo-api/internal/utils"
)

// AuthMiddleware is a middleware that checks if the user is authenticated.
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")

		if authHeader == "" {
			utils.RespondJSON(w, http.StatusUnauthorized, "Unauthorized", nil)
			return
		}

		// Expecting a token format like "Bearer <token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			utils.RespondJSON(w, http.StatusUnauthorized, "Invalid token format", nil)
			return
		}

		tokenString := parts[1]
		claims, err := auth.ValidateToken(tokenString)
		if err != nil {
			utils.RespondJSON(w, http.StatusUnauthorized, "Invalid token", nil)
			return
		}

		// Add user ID to the context
		ctx := context.WithValue(r.Context(), "userID", claims.UserID)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
