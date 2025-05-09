package providers

import (
	"encoding/json"
	"errors"
	"os"
	"strconv"
	"time"

	"github.com/jaygaha/roadmap-go-projects/personal-blog-app/models"
)

const storageFile = "posts.json"

// Load paginated lists of posts from the JSON file
func LoadPaginatedPosts(page, pageSize int) ([]models.Blog, int, error) {
	posts, err := loadPosts()
	if err != nil {
		return nil, 0, err
	}

	totalPosts := len(posts)
	startIndex := (page - 1) * pageSize
	endIndex := startIndex + pageSize
	if startIndex >= totalPosts {
		return nil, 0, nil // No posts on this page
	}
	if endIndex > totalPosts {
		endIndex = totalPosts
	}
	paginatedPosts := posts[startIndex:endIndex]

	return paginatedPosts, totalPosts, nil
}

// reads the posts from the JSON file
func loadPosts() ([]models.Blog, error) {
	if _, err := os.Stat(storageFile); errors.Is(err, os.ErrNotExist) {
		err := os.WriteFile(storageFile, []byte("[]"), 0644) // create the file if it doesn't exist

		if err != nil {
			return nil, err
		}
	}

	data, err := os.ReadFile(storageFile)
	if err != nil {
		return nil, err
	}

	var posts []models.Blog
	if err := json.Unmarshal(data, &posts); err != nil {
		return nil, err
	}

	return posts, nil
}

func StringToInt(s string) int {
	intValue, err := strconv.Atoi(s)
	if err == nil {
		return intValue
	}

	// Handle the error or return a default value
	return 0
}

// save all posts to the JSON file
func savePosts(posts []models.Blog) error {
	data, err := json.MarshalIndent(posts, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(storageFile, data, 0644)
}

// create a new post
func CreatePost(title, description string, createdAt time.Time) (models.Blog, error) {
	posts, err := loadPosts()
	if err != nil {
		return models.Blog{}, err
	}

	newPost := models.Blog{
		ID:          len(posts) + 1,
		Title:       title,
		Description: description,
		CreatedAt:   createdAt,
	}

	posts = append(posts, newPost)

	err = savePosts(posts)

	return newPost, err
}

func UpdatePost(id int, title, description string) (models.Blog, error) {
	posts, err := loadPosts()
	if err != nil {
		return models.Blog{}, err
	}

	updated := false
	var updatedPost models.Blog

	for i, post := range posts {
		if post.ID == id {
			posts[i].Title = title
			posts[i].Description = description
			now := time.Now()
			posts[i].UpdatedAt = &now

			updatedPost = posts[i]
			updated = true
			break
		}
	}

	if !updated {
		return models.Blog{}, errors.New("post not found")
	}

	err = savePosts(posts)
	return updatedPost, err
}

// show a single post
func GetPostByID(id int) (models.Blog, error) {
	posts, err := loadPosts()
	if err != nil {
		return models.Blog{}, err
	}
	for _, post := range posts {
		if post.ID == id {
			return post, nil
		}
	}
	return models.Blog{}, errors.New("post not found")
}

// delete a post
func DeletePost(id int) error {
	posts, err := loadPosts()
	if err != nil {
		return err
	}

	updatedPosts := make([]models.Blog, 0)
	for _, post := range posts {
		if post.ID != id {
			updatedPosts = append(updatedPosts, post)
		}
	}

	return savePosts(updatedPosts)
}
