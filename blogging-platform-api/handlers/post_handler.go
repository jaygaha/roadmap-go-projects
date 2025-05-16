package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/jaygaha/roadmap-go-projects/blogging-platform-api/database"
	"github.com/jaygaha/roadmap-go-projects/blogging-platform-api/models"
)

// Create response struct for consistent response format across all endpoints
type response struct {
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

// PostWOIdHandler handles different HTTP methods for a given post like /posts
// Post Without ID
func PostWOIdHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getAllPostsHandler(w, r)
	case http.MethodPost:
		createPostHandler(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// PostHandlerSelector handles different HTTP methods for a given post ID like /posts/{id}
func PostHandlerSelector(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getPostHandler(w, r)
	case http.MethodPut:
		updatePostHandler(w, r)
	case http.MethodDelete:
		deletePostHandler(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// getAllPostsHandler handles the retrieval of all blog posts using filter
func getAllPostsHandler(w http.ResponseWriter, r *http.Request) {
	// only allow get method
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// get optional query params
	term := r.URL.Query().Get("term")

	// build base query
	query := `SELECT id, title, content, category, tags, created_at, updated_at FROM posts`
	var args []any

	// add filter to query if term is provided
	if term != "" {
		query += ` WHERE title LIKE ? OR content LIKE ? OR category LIKE ?`
		likeTerm := "%" + term + "%"
		args = append(args, likeTerm, likeTerm, likeTerm)
	}

	// order by
	query += ` ORDER BY created_at ASC`

	rows, err := database.DB.Query(query, args...)

	if err != nil {
		log.Println("Database query error: ", err)
		http.Error(w, "Failed to fetch posts", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	// iterate over the rows and store in a slice
	var posts []models.Post
	for rows.Next() { // iterate over rows
		var post models.Post
		var tagsJson string // store tags as JSON string

		err := rows.Scan(
			&post.ID,
			&post.Title,
			&post.Content,
			&post.Category,
			&tagsJson,
			&post.CreatedAt,
			&post.UpdatedAt,
		)
		if err != nil {
			log.Println("Database scan error: ", err)
			http.Error(w, "Failed to fetch posts", http.StatusInternalServerError)
			return
		}

		// convert tagsJson to slice of strings
		if err := json.Unmarshal([]byte(tagsJson), &post.Tags); err != nil {
			log.Println("Failed to parse tags: ", err)
			http.Error(w, "Failed to fetch posts", http.StatusInternalServerError)
			return
		}
		posts = append(posts, post)
	}

	// return the posts
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response{
		Message: "Posts retrieved successfully",
		Data:    posts,
	})
}

// createPostHandler handles the creation of a new blog post
func createPostHandler(w http.ResponseWriter, r *http.Request) {
	// only allow post method
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// read request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}

	// parse JSON into Post struct
	var post models.Post
	err = json.Unmarshal(body, &post)
	if err != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	// validate post data
	if strings.TrimSpace(post.Title) == "" ||
		strings.TrimSpace(post.Content) == "" ||
		strings.TrimSpace(post.Category) == "" ||
		len(post.Tags) == 0 {
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}

	// convert tags array to JSON strings for storage
	tagsJson, err := json.Marshal(post.Tags)
	if err != nil {
		http.Error(w, "Failed to encode tags", http.StatusInternalServerError)
	}

	// Insert the post into the database
	// Set the time to now with ISO 8601 format
	post.CreatedAt = time.Now().UTC()
	post.UpdatedAt = time.Now().UTC()

	query := `
		INSERT INTO posts (title, content, category, tags, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5,$6)
	`

	row, err := database.DB.Exec(query, post.Title, post.Content, post.Category, tagsJson, post.CreatedAt, post.UpdatedAt)

	if err != nil {
		log.Println("Database insert error: ", err)
		http.Error(w, "Failed to create post", http.StatusInternalServerError)
		return
	}

	id, err := row.LastInsertId()
	if err != nil {
		log.Println("Failed to retrieve last insert ID: ", err)
		http.Error(w, "Failed to create post", http.StatusInternalServerError)
		return
	}
	post.ID = int(id)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response{
		Message: "Post created successfully",
		Data:    post,
	})
}

// getPostHandler handles the retrieval of a blog post using existing ID
func getPostHandler(w http.ResponseWriter, r *http.Request) {
	// only allow get method
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	// get id from url
	id, err := getId(r.URL.Path)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// validate id exists
	post, err := getPostById(id)
	if err != nil {
		http.Error(w, "Post not found", http.StatusNotFound)
		return
	}
	// return the post
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response{
		Message: "Post retrieved successfully",
		Data:    post,
	})
}

// updatePostHandler handles the update of a blog post using existing ID
func updatePostHandler(w http.ResponseWriter, r *http.Request) {
	// only allow put method
	if r.Method != http.MethodPut {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// get id from url
	id, err := getId(r.URL.Path)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// read request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}

	// validate id exists
	_, err = getPostById(id)
	if err != nil {
		http.Error(w, "Post not found", http.StatusNotFound)
		return
	}

	// parse JSON into Post struct
	var post models.Post
	err = json.Unmarshal(body, &post)
	if err != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	// validate post data
	if strings.TrimSpace(post.Title) == "" ||
		strings.TrimSpace(post.Content) == "" ||
		strings.TrimSpace(post.Category) == "" ||
		len(post.Tags) == 0 {
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}

	// convert tags array to JSON strings for storage
	tagsJson, err := json.Marshal(post.Tags)
	if err != nil {
		http.Error(w, "Failed to encode tags", http.StatusInternalServerError)
		return
	}

	// update the post in the database
	post.UpdatedAt = time.Now().UTC()
	query := `
		UPDATE posts
		SET title = $1, content = $2, category = $3, tags = $4, updated_at = $5
		WHERE id = $6
	`
	_, err = database.DB.Exec(query, post.Title, post.Content, post.Category, tagsJson, post.UpdatedAt, id)
	if err != nil {
		log.Println("Database update error: ", err)
		http.Error(w, "Failed to update post", http.StatusInternalServerError)
		return
	}

	// Get the updated post
	updatedPost, err := getPostById(id)
	if err != nil {
		log.Println("Failed to retrieve updated post: ", err)
		http.Error(w, "Failed to retrieve updated post", http.StatusInternalServerError)
		return
	}

	// return the updated post
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response{
		Message: "Post updated successfully",
		Data:    updatedPost,
	})
}

// deletePostHandler handles the deletion of a blog post using existing ID
func deletePostHandler(w http.ResponseWriter, r *http.Request) {
	// only allow delete method
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// extract id from url
	id, err := getId(r.URL.Path)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// validate id exists
	_, err = getPostById(id)
	if err != nil {
		http.Error(w, "Post not found", http.StatusNotFound)
		return
	}
	// delete the post from the database
	query := `
		DELETE FROM posts
		WHERE id = $1
	`
	_, err = database.DB.Exec(query, id)
	if err != nil {
		log.Println("Database delete error: ", err)
		http.Error(w, "Failed to delete post", http.StatusInternalServerError)
		return
	}
	// return success response
	w.WriteHeader(http.StatusNoContent)
}

func getId(urlPath string) (int, error) {
	// extract the id from the url path
	parts := strings.Split(urlPath, "/")

	if len(parts) != 3 {
		return 0, fmt.Errorf("Invalid URL path")
	}

	id, err := strconv.Atoi(parts[2])
	if err != nil {
		return 0, fmt.Errorf("Invalid post ID")
	}

	return id, nil
}

func getPostById(id int) (*models.Post, error) {
	const query = `
		SELECT id, title, content, category, tags, created_at, updated_at
		FROM posts 
		WHERE id = $1`

	post := &models.Post{}
	var tagsJson string

	if err := database.DB.QueryRow(query, id).Scan(
		&post.ID,
		&post.Title,
		&post.Content,
		&post.Category,
		&tagsJson,
		&post.CreatedAt,
		&post.UpdatedAt,
	); err != nil {
		return nil, fmt.Errorf("failed to fetch post: %w", err)
	}

	if err := json.Unmarshal([]byte(tagsJson), &post.Tags); err != nil {
		return nil, fmt.Errorf("failed to parse tags: %w", err)
	}

	return post, nil
}
