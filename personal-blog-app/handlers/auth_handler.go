package handlers

import (
	"html/template"
	"net/http"

	"github.com/jaygaha/roadmap-go-projects/personal-blog-app/providers"
)

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	data := TemplateData{
		Title: "Login",
		Data: map[string]any{
			"IsLoggedIn": false,
		},
	}

	if r.Method == http.MethodPost {
		// Handle login form submission
		email := r.FormValue("email")
		password := r.FormValue("password")

		// redirect to dashboard if email and password are correct
		if email == "jaygaha@gmail.com" && password == "admin@123" {
			// Create a new session
			providers.CreateSession(w)

			http.Redirect(w, r, "/admin", http.StatusSeeOther)
			return
		}

		// append error message to data
		data.Data = map[string]any{
			"Error": "Invalid email or password",
		}
	}

	// redirect is user is logged in
	if providers.IsAuthenticated(r) {
		http.Redirect(w, r, "/admin", http.StatusSeeOther)
		return
	}

	temp := template.Must(template.ParseFiles("templates/base.html", "templates/pages/login.html"))
	err := temp.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	// Delete the session
	providers.DeleteSession(w, r)

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}
