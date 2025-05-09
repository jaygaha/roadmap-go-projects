package handlers

import (
	"html/template"
	"math"
	"net/http"

	"github.com/jaygaha/roadmap-go-projects/personal-blog-app/models"
	"github.com/jaygaha/roadmap-go-projects/personal-blog-app/providers"
)

func AdminDashboardHandler(w http.ResponseWriter, r *http.Request) {
	posts, pagination := paginatedBlogs(w, r)

	data := TemplateData{
		Title: "Dashboard",
		Data: map[string]any{
			"IsLoggedIn": true,
			"Title":      "Manage Blog Posts",
			"Posts":      posts,
			"Pagination": pagination,
		},
	}

	temp := template.Must(template.New("base.html").Funcs(funcMap).ParseFiles("templates/base.html", "templates/pages/admin_blog_list.html"))
	err := temp.Execute(w, data)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// common blogs list with pagination
func paginatedBlogs(w http.ResponseWriter, r *http.Request) ([]models.Blog, models.Pagination) {
	// get query params
	page := r.URL.Query().Get("page")
	pageSize := r.URL.Query().Get("pageSize")
	// convert page and pageSize to integer
	pageInt := providers.StringToInt(page)
	if pageInt == 0 {
		pageInt = 1
	}
	pageSizeInt := providers.StringToInt(pageSize)
	if pageSizeInt == 0 {
		pageSizeInt = 10
	}

	posts, totalPosts, err := providers.LoadPaginatedPosts(pageInt, pageSizeInt)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return nil, models.Pagination{}
	}

	// calculate total pages
	totalPages := int(math.Ceil(float64(totalPosts) / float64(pageSizeInt)))

	// create pagination metadata
	pagination := models.Pagination{
		CurrentPage: pageInt,
		TotalPages:  totalPages,
		HasPrev:     pageInt > 1,
		HasNext:     pageInt < totalPages,
	}

	return posts, pagination
}
