package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/jaygaha/roadmap-go-projects/expense-tracker-api/config"
	"github.com/jaygaha/roadmap-go-projects/expense-tracker-api/internal/response"
	"github.com/jaygaha/roadmap-go-projects/expense-tracker-api/pkg/utils"
)

type key string

const UserIDKey key = "userID"

// JWTAuthMiddleware verifies JWT tokens
func JWTAuthMiddleware(cfg *config.Config) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				response.Error(w, http.StatusUnauthorized, "Unauthorized")
				return
			}

			tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
			if tokenStr == authHeader {
				response.Error(w, http.StatusUnauthorized, "Invalid token format")
				return
			}

			if utils.IsTokenBlocked(tokenStr) {
				response.Error(w, http.StatusUnauthorized, "Token has been revoked")
				return
			}

			claims, err := utils.ParseToken(tokenStr, cfg.JWTSecret)
			if err != nil {
				response.Error(w, http.StatusUnauthorized, "Invalid token")
				return
			}

			userID := uint(claims["sub"].(float64))

			// Attach userID to request context
			ctx := context.WithValue(r.Context(), UserIDKey, uint(userID))
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
