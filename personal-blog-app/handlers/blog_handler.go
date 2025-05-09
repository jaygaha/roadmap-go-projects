package handlers

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/jaygaha/roadmap-go-projects/personal-blog-app/providers"
)

func BlogNewHandler(w http.ResponseWriter, r *http.Request) {
	data := TemplateData{
		Title: "Create New Blog",
		Data: map[string]any{
			"IsLoggedIn": true,
			"FormTitle":  "Create a New Blog",
			"FormUrl":    "/blogs/submit",
			"FormMethod": "POST",
			"FormData": map[string]any{
				"ID":          "",
				"Title":       "",
				"Description": "",
				"CreatedAt":   time.Now().Format("2006-01-02"),
			},
		},
	}

	temp := template.Must(template.ParseFiles("templates/base.html", "templates/pages/form.html"))
	err := temp.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func BlogPostHandler(w http.ResponseWriter, r *http.Request) {
	// only accept POST OR PUT requests
	if r.Method != http.MethodPost && r.Method != http.MethodPut {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// parse the form data
	err := r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// validate the form data
	title := r.FormValue("title")
	description := r.FormValue("description")
	createdAd := r.FormValue("created_at")
	// check if id is given
	id := r.FormValue("id")

	if title == "" || description == "" || createdAd == "" {
		http.Error(w, "Title, description and date are required", http.StatusUnprocessableEntity)
		return
	}

	// validate the date
	_, err = time.Parse("2006-01-02", createdAd)
	if err != nil {
		http.Error(w, "Invalid date format", http.StatusUnprocessableEntity)
		return
	}

	// date converted to time.Time
	createdAt, err := time.Parse("2006-01-02", createdAd)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// update if id is gq
	if id != "" {
		id, err := strconv.Atoi(id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// get the blog post from the database
		blogPost, err := providers.GetPostByID(id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		// update the blog post in the database
		blogPost, err = providers.UpdatePost(blogPost.ID, title, description)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// redirect to the blog post page
		http.Redirect(w, r, fmt.Sprintf("/article/%d", blogPost.ID), http.StatusSeeOther)
	} else {
		blogPost, err := providers.CreatePost(title, description, createdAt)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// redirect to the blog post page
		http.Redirect(w, r, fmt.Sprintf("/article/%d", blogPost.ID), http.StatusSeeOther)
	}
}

// show blog
func BlogShowHandler(w http.ResponseWriter, r *http.Request) {
	// get the blog post ID from the URL passed as a parameter like article/1
	pathSegments := strings.Split(r.URL.Path, "/")

	if len(pathSegments) >= 2 && pathSegments[1] == "article" {
		if len(pathSegments) > 2 {
			dynamicID := pathSegments[2]
			// Handle trailing slash
			if dynamicID == "" && len(pathSegments) == 3 {
				http.NotFound(w, r)
				return
			} else if dynamicID == "" && len(pathSegments) > 3 {
				http.NotFound(w, r)
				return
			}

			// convert the dynamicID to an integer
			id, err := strconv.Atoi(dynamicID)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			// get the blog post from the database
			post, err := providers.GetPostByID(id)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			data := TemplateData{
				Title: post.Title,
				Data: map[string]any{
					"IsLoggedIn": providers.IsAuthenticated(r),
					"Post":       post,
				},
			}

			temp := template.Must(template.ParseFiles("templates/base.html", "templates/pages/page.html"))
			err = temp.Execute(w, data)

			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}

		} else {
			// redirect to the home page
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}
	} else {
		http.NotFound(w, r)
	}
}

func BlogEditHandler(w http.ResponseWriter, r *http.Request) {
	// get the blog post ID from the URL passed as a parameter like edit/1
	pathSegments := strings.Split(r.URL.Path, "/")

	if len(pathSegments) >= 2 && pathSegments[1] == "edit" {
		if len(pathSegments) > 2 {
			dynamicID := pathSegments[2]
			// Handle trailing slash
			if dynamicID == "" && len(pathSegments) == 3 {
				http.NotFound(w, r)
				return
			} else if dynamicID == "" && len(pathSegments) > 3 {
				http.NotFound(w, r)
				return
			}
			// convert the dynamicID to an integer
			id, err := strconv.Atoi(dynamicID)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			// get the blog post from the database
			post, err := providers.GetPostByID(id)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			data := TemplateData{
				Title: "Edit a blog",
				Data: map[string]any{
					"IsLoggedIn": true,
					"FormTitle":  "Edit a Blog",
					"FormUrl":    "/blogs/submit",
					"FormMethod": "PUT",
					"FormData": map[string]any{
						"ID":          post.ID,
						"Title":       post.Title,
						"Description": post.Description,
						"CreatedAt":   post.CreatedAt.Format("2006-01-02"),
					},
				},
			}

			temp := template.Must(template.ParseFiles("templates/base.html", "templates/pages/form.html"))
			err = temp.Execute(w, data)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		} else {
			// redirect to the home page
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}
	} else {
		http.NotFound(w, r)
	}
}

func BlogDeleteHandler(w http.ResponseWriter, r *http.Request) {
	// only post method is allowed
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	deleteId := r.FormValue("id")

	// validate the id
	if deleteId == "" {
		http.Error(w, "ID is required", http.StatusUnprocessableEntity)
		return
	}
	// convert the id to an integer
	id, err := strconv.Atoi(deleteId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// delete the blog post from the database
	err = providers.DeletePost(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// redirect to the home page
	http.Redirect(w, r, "/admin", http.StatusSeeOther)
}
