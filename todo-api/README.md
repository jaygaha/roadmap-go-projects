# Todo API

A simple RESTful API for managing todo items built with Go. This project demonstrates how to create a basic CRUD API with user authentication using JWT tokens.

## Features

- User registration and login with JWT authentication
- Create, read, update, and delete todo items
- Filter todos by completion status
- Pagination and sorting options
- SQLite database for data storage

## Project Structure

```
├── cmd/
│   └── server/         # Application entry point
├── internal/
│   ├── auth/           # JWT authentication
│   ├── db/             # Database connection and migrations
│   ├── handlers/       # HTTP request handlers
│   ├── middleware/     # HTTP middleware (auth)
│   ├── models/         # Data models
│   ├── routes/         # API routes
│   └── utils/          # Utility functions
├── go-todo-api.json    # Postman collection for testing
├── go.mod              # Go module dependencies
└── Makefile            # Build and run commands
```

## Prerequisites

- Go 1.24 or higher
- SQLite (included as a dependency)

## Getting Started

### Installation

1. Clone the repository
2. Navigate to the project directory

```bash
cd todo-api
```

3. Install dependencies

```bash
go mod download
```

### Running the Application

Use the Makefile to run the application:

```bash
make run
```

Or specify a custom port:

```bash
make run port=8080
```

The server will start on the specified port (default: 8800).

## API Endpoints

### Authentication

- `POST /register` - Register a new user
- `POST /login` - Login and get JWT token

### Todo Operations

- `GET /todos` - List all todos (with optional filtering)
- `POST /todos` - Create a new todo
- `GET /todos/{id}` - Get a specific todo
- `PUT /todos/{id}` - Update a todo
- `DELETE /todos/{id}` - Delete a todo

## Authentication

All todo endpoints require authentication. Include the JWT token in the Authorization header:

```
Authorization: Bearer <your_token>
```

## Example Usage

1. Register a new user

```bash
curl -X POST http://localhost:8800/register \
  -H "Content-Type: application/json" \
  -d '{"name":"John","username":"john","password":"password123"}'
```

2. Login to get a token

```bash
curl -X POST http://localhost:8800/login \
  -H "Content-Type: application/json" \
  -d '{"username":"john","password":"password123"}'
```

3. Create a new todo (using the token)

```bash
curl -X POST http://localhost:8800/todos \
  -H "Authorization: Bearer <your_token>" \
  -H "Content-Type: application/json" \
  -d '{"title":"Learn Go","description":"Complete the Go tutorial"}'
```

4. List all todos

```bash
curl -X GET http://localhost:8800/todos \
  -H "Authorization: Bearer <your_token>"
```

## Testing with Postman

Import the included `go-todo-api.json` file into Postman to test the API endpoints.

## Database

The application uses SQLite for data storage. The database file (`todo.db`) is created automatically when you run the application for the first time.

## Project Link

- [Todo List API](https://roadmap.sh/projects/todo-list-api)

## Acknowledgments

- Part of the Go programming language learning roadmap projects
- Created by [jaygaha](https://github.com/jaygaha)