# Expense Tracker API

A RESTful API built with Go that allows users to track their expenses. This project is a solution to the [Expense Tracker API Project](https://roadmap.sh/projects/expense-tracker-api) from the Backend Developer Roadmap.

This project is build according to the Go standard folder structure.

## Features

### Authentication
- User signup with name, email, and password
- JWT-based authentication
- Protected routes requiring valid JWT tokens
- Logout functionality with token invalidation

### Expense Management
- Create, read, update, and delete expenses
- Each expense includes:
  - Title
  - Amount
  - Category
  - Creation timestamp

### Expense Categories
- Predefined expense categories
- CRUD operations for managing categories
- Category validation for expenses

### Expense Filtering
- Filter expenses by date ranges:
  - Past week
  - Past month
  - Last 3 months
  - Custom date range

## Tech Stack

- **Language**: Go 1.24.0
- **Web Framework**: Gorilla Mux
- **Database**: SQLite with GORM
- **Authentication**: JWT (JSON Web Tokens)
- **Validation**: go-playground/validator

## Project Structure

```
├── cmd/
│   └── server/           # Application entry point
├── config/              # Configuration management
├── internal/
│   ├── database/        # Database connection and setup
│   ├── handlers/        # HTTP request handlers
│   ├── middleware/      # Custom middleware (JWT auth)
│   ├── models/          # Data models
│   ├── repositories/    # Database operations
│   ├── request/         # Request validation
│   ├── response/        # Response formatting
│   ├── routes/          # Route definitions
│   └── services/        # Business logic
└── pkg/
    └── utils/           # Shared utilities
└── .env.example         # Environment variables template
└── Makefile             # Makefile for build and running
└── postman-collection.json # Postman collection for API testing
└── README.md            # Project documentation
```

## API Endpoints

### Authentication
- `POST /api/v1/auth/signup` - Register a new user
- `POST /api/v1/auth/login` - Login and receive JWT token
- `POST /api/v1/auth/logout` - Logout and invalidate token

### Expenses
- `POST /api/v1/expenses` - Create a new expense
- `GET /api/v1/expenses` - List all expenses (with optional filters)
- `GET /api/v1/expenses/{id}` - Get expense by ID
- `PUT /api/v1/expenses/{id}` - Update an expense
- `DELETE /api/v1/expenses/{id}` - Delete an expense

### Categories
- `POST /api/v1/categories` - Create a new category
- `GET /api/v1/categories` - List all categories
- `GET /api/v1/categories/{id}` - Get category by ID
- `PUT /api/v1/categories/{id}` - Update a category
- `DELETE /api/v1/categories/{id}` - Delete a category

## Getting Started

1. Clone the repository
2. Copy `.env.example` to `.env` and configure your environment variables
3. Run the application:
   ```bash
   go mod tidy
   # Using Makefile
   make dev
   # Without Makefile
   go run cmd/server/main.go
   ```

## Security Features

- Password hashing using bcrypt
- JWT token validation and expiration
- Token blocklist for logout functionality
- Protected routes with middleware authentication
- Input validation for all requests

## Error Handling

- Consistent error response format
- Proper HTTP status codes
- Validation error messages
- Database error handling

## Future Improvements

- Add unit and integration tests
- Implement rate limiting
- Add request logging
- Support for multiple currencies
- Export expenses to CSV/PDF
- Add user roles and permissions
- Implement data pagination

## Postman Collection

You can find the Postman collection for [this API collection](./postman-collection.json). Import this collection into Postman to easily test the API endpoints.

## Project Link

- [Expense Tracker API](https://roadmap.sh/projects/expense-tracker-api)

## Acknowledgments

- Part of the Go programming language learning roadmap projects
- Created by [jaygaha](https://github.com/jaygaha)