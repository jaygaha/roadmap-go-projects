package middleware

import (
	"net/http"

	"github.com/jaygaha/roadmap-go-projects/personal-blog-app/providers"
)

// authenticate against session

func Authenticate(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Check if the user is authenticated
		// If not, redirect to the login page
		// If authenticated, call the next handler
		if providers.IsAuthenticated(r) == false {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		next.ServeHTTP(w, r)
	}
}
