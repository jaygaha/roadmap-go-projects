package handlers

import (
	"html/template"
	"net/http"

	"github.com/jaygaha/roadmap-go-projects/personal-blog-app/providers"
)

// Define global variables for templates
type TemplateData struct {
	Title string
	Data  any
}

// global template function accessible to all templates
var funcMap = template.FuncMap{
	"add": func(a, b int) int { return a + b },
	"sub": func(a, b int) int { return a - b },
}

func WelcomeHandler(w http.ResponseWriter, r *http.Request) {
	posts, pagination := paginatedBlogs(w, r)

	temp := template.Must(template.New("base.html").Funcs(funcMap).ParseFiles("templates/base.html", "templates/pages/blog_list.html"))

	isLoggedIn := providers.IsAuthenticated(r)

	data := TemplateData{
		Title: "Latest",
		Data: map[string]any{
			"IsLoggedIn": isLoggedIn,
			"Title":      "Latest Blogs",
			"Posts":      posts,
			"Pagination": pagination,
		},
	}

	err := temp.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
