# Blogging Platform API

A RESTful API for a blogging platform built with Go. This project provides endpoints to create, read, update, and delete blog posts with filtering capabilities.

## Features

- RESTful API architecture
- CRUD operations for blog posts
- SQLite database for data persistence
- Search and filter posts by terms
- JSON response format
- Structured error handling

## Project Structure

- **main.go**: Entry point of the application
- **routes/api.go**: API route definitions
- **handlers/post_handler.go**: HTTP request handlers for blog posts
- **models/post.go**: Data structures for blog posts
- **database/**:
  - **db.go**: Database connection and migration setup
  - **blog.db**: SQLite database file

## Technologies Used

- Go (Golang)
- SQLite (via modernc.org/sqlite)
- Standard library HTTP server
- JSON for data interchange

## Installation

### Prerequisites

- Go 1.18 or higher
- Git

### Steps

1. Clone the repository:
   ```bash
   git clone https://github.com/yourusername/blogging-platform-api.git
   cd blogging-platform-api
   ```

2. Install dependencies:
   ```bash
   go mod download
   ```

3. Run the application:
   ```bash
   go run main.go
   ```

4. The server will start on port 8800. You can access it at http://localhost:8800

## API Endpoints

### Health Check
- `GET /ping`: Check if the API is running

### Blog Posts
- `GET /posts`: Get all blog posts
  - Query parameters:
    - `term`: Filter posts by title, content, or category
- `POST /posts`: Create a new blog post
- `GET /posts/{id}`: Get a specific blog post by ID
- `PUT /posts/{id}`: Update a specific blog post
- `DELETE /posts/{id}`: Delete a specific blog post

## API Request/Response Examples

### Create a Post

**Request:**
```http
POST /posts
Content-Type: application/json

{
  "title": "Getting Started with Go",
  "content": "Go is a statically typed, compiled programming language...",
  "category": "Programming",
  "tags": ["golang", "programming", "tutorial"]
}
```

**Response:**
```json
{
  "message": "Post created successfully",
  "data": {
    "id": 1,
    "title": "Getting Started with Go",
    "content": "Go is a statically typed, compiled programming language...",
    "category": "Programming",
    "tags": ["golang", "programming", "tutorial"],
    "created_at": "2025-05-15T03:49:46.547885Z",
    "updated_at": "2025-05-15T03:49:46.547885Z"
  }
}
```

### Get All Posts

**Request:**
```http
GET /posts
```

**Response:**
```json
{
  "message": "Posts retrieved successfully",
  "data": [
    {
      "id": 1,
      "title": "Getting Started with Go",
      "content": "Go is a statically typed, compiled programming language...",
      "category": "Programming",
      "tags": ["golang", "programming", "tutorial"],
      "created_at": "2025-05-15T03:49:46.547885Z",
      "updated_at": "2025-05-15T03:49:46.547885Z"
    }
    // ...
  ]
}
```

## Error Handling

The API returns appropriate HTTP status codes and error messages in case of failures:

- `400 Bad Request`: Invalid input data
- `404 Not Found`: Resource not found
- `405 Method Not Allowed`: HTTP method not supported for the endpoint
- `500 Internal Server Error`: Server-side error

## Development

### Adding New Features

1. Define new routes in `routes/api.go`
2. Create handlers in `handlers/` directory
3. Add models in `models/` directory if needed
4. Update database migrations in `database/db.go` if needed

## Sample Collection
You can import the provided Postman collection for testing the API endpoints:
- [Blogging Platform API.postman_collection.json](./collection_postman_go-bloging-platform.json)

## Acknowledgments

- [roadmap.sh](https://roadmap.sh/projects/blogging-platform-api) for the project inspiration
- Created by [jaygaha](https://github.com/jaygaha)