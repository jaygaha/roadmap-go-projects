package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/jaygaha/roadmap-go-projects/intermediate/e-commerce-api-service/services"
)

// contextKey is a type to ensure unique keys in context
type contextKey string

const (
	UserIdKey contextKey = "user_id"
	RoleKey   contextKey = "role"
)

// Auth returns middleware that validates the JWT from the Authorization header
// and injects user_id and role into the request context.
func Auth(authSvc *services.AuthService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			header := r.Header.Get("Authorization")
			if header == "" || !strings.HasPrefix(header, "Bearer ") {
				http.Error(w, `{"error":"missing or invalid authorization header"}`, http.StatusUnauthorized)
				return
			}

			tokenStr := strings.TrimPrefix(header, "Bearer ")
			claims, err := authSvc.ValidateToken(tokenStr)
			if err != nil {
				http.Error(w, `{"error":"invalid or expired token"}`, http.StatusUnauthorized)
				return
			}

			// Inject claims into context for downstream handlers.
			userId := int64(claims["user_id"].(float64))
			role := claims["role"].(string)

			ctx := context.WithValue(r.Context(), UserIdKey, userId)
			ctx = context.WithValue(ctx, RoleKey, role)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
