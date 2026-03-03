// For server-rendered pages, the JWT lives in an HttpOnly cookie
// instead of an Authorization header. This middleware reads it.

package middleware

import (
	"context"
	"net/http"

	"github.com/jaygaha/roadmap-go-projects/intermediate/e-commerce-api-service/services"
)

// CookieAuth is like Auth but reads the JWT from a cookie.
// It does NOT block unauthenticated users — it just populates
// context if a valid token is present. Page handlers decide
// whether to redirect based on whether User is nil.
func CookieAuth(authSvc *services.AuthService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			cookie, err := r.Cookie("token")
			if err != nil || cookie.Value == "" {
				next.ServeHTTP(w, r)
				return
			}

			claims, err := authSvc.ValidateToken(cookie.Value)
			if err != nil {
				// Token expired or invalid — clear it
				http.SetCookie(w, &http.Cookie{
					Name: "token", Value: "", Path: "/", MaxAge: -1,
				})
				next.ServeHTTP(w, r)
				return
			}

			userID := int64(claims["user_id"].(float64))
			role := claims["role"].(string)
			ctx := context.WithValue(r.Context(), UserIdKey, userID)
			ctx = context.WithValue(ctx, RoleKey, role)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
