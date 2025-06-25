# URL Shortening Service

A RESTful API service that allows users to shorten long URLs, built with `Go` and `MongoDB`. This project implements a clean architecture pattern with handlers, services, and repositories.

## Features

- Create short URLs from long URLs
- Retrieve original URLs using short codes
- Update existing URLs
- Delete URLs
- Track and view URL access statistics
- Pagination support for listing URLs

## Project Structure

```
url-shortening-service/
├── cmd/
│   └── api/
│       └── main.go              # Application entry point
├── internal/
│   ├── config/                  # Application configuration
│   │   └── app.go
│   ├── database/                # Database connection setup
│   │   └── mongodb.go
│   ├── handler/                 # HTTP request handlers
│   │   └── url_handler.go
│   ├── models/                  # Data structures
│   │   └── url.go
│   ├── repository/              # Data access layer
│   │   ├── interfaces.go        # Repository interfaces
│   │   └── mongodb/             # MongoDB implementation
│   │       └── url_repository.go
│   ├── routes/                  # API route definitions
│   │   └── api.go
│   ├── service/                 # Business logic layer
│   │   ├── interfaces.go        # Service interfaces
│   │   └── url_service.go       # URL service implementation
│   └── utils/                   # Utility functions
│       ├── api_response.go      # API response helpers
│       ├── helper.go            # General helpers
│       └── shortcode.go         # Short code generation
├── .air.toml                    # Air configuration for hot reload
├── docker-compose.yml           # Docker Compose configuration
└── go.mod                       # Go module dependencies
```

## Architecture

This project follows a clean architecture pattern with the following layers:

1. **Handlers**: Handle HTTP requests and responses
2. **Services**: Implement business logic
3. **Repositories**: Handle data storage and retrieval

The application uses dependency injection to connect these layers, making the code modular and testable.

## API Endpoints

### Create Short URL
```
POST /api/shorten

Request Body:
{
  "url": "https://www.example.com/some/long/url"
}

Response (201 Created):
{
  "message": "URL created successfully",
  "data": {
    "id": "...",
    "url": "https://www.example.com/some/long/url",
    "short_code": "abc123",
    "short_url": "http://localhost:8800/abc123",
    "click_count": 0,
    "created_at": "2025-06-25T12:00:00Z"
  }
}
```

### Get All URLs
```
GET /api/shorten?page=1&limit=20

Response (200 OK):
{
  "message": "URLs retrieved successfully",
  "data": {
    "list": [...],
    "pagination": {
      "page": 1,
      "limit": 20,
      "total": 100,
      "pages": 5
    }
  }
}
```

### Get Original URL
```
GET /api/shorten/{shortCode}

Response (200 OK):
{
  "message": "URL retrieved successfully",
  "data": "https://www.example.com/some/long/url"
}
```

### Get URL Statistics
```
GET /api/shorten/{shortCode}/stats

Response (200 OK):
{
  "message": "URL retrieved successfully",
  "data": {
    "id": "...",
    "url": "https://www.example.com/some/long/url",
    "short_code": "abc123",
    "short_url": "http://localhost:8800/abc123",
    "click_count": 10,
    "created_at": "2025-06-25T12:00:00Z"
  }
}
```

### Update URL
```
PATCH /api/shorten/{shortCode}

Request Body:
{
  "url": "https://www.example.com/some/updated/url"
}

Response (200 OK):
{
  "message": "URL updated successfully",
  "data": {
    "id": "...",
    "url": "https://www.example.com/some/updated/url",
    "short_code": "abc123",
    "short_url": "http://localhost:8800/abc123",
    "click_count": 10,
    "created_at": "2025-06-25T12:00:00Z"
  }
}
```

### Delete URL
```
DELETE /api/shorten/{shortCode}

Response (200 OK):
{
  "message": "URL deleted successfully",
  "data": null
}
```

## Setup and Running

### Prerequisites

- Go 1.16 or higher
- Docker and Docker Compose
- MongoDB (or use the provided Docker Compose setup)

### Running with Docker

1. Clone the repository
2. Navigate to the project directory
3. Start the application using Docker Compose:

```bash
# For development with hot reload
docker-compose up api-dev

# For production
docker-compose up api-prod
```

The API will be available at:
- Development: http://localhost:8800
- Production: http://localhost:8801

### Running Locally

1. Clone the repository
2. Navigate to the project directory
3. Install dependencies:

```bash
go mod download
```

4. Start MongoDB (using Docker or locally)
5. Run the application:

```bash
go run cmd/api/main.go
```

The API will be available at http://localhost:8800

## Implementation Details

### Short Code Generation

The service generates unique short codes using a cryptographically secure random number generator and base64 URL encoding, ensuring that each short code is unique and URL-safe.

### Click Tracking

The service tracks the number of times each short URL is accessed. When a user accesses a short URL, the click count is incremented asynchronously to avoid impacting the redirect performance.

### Pagination

The API supports pagination for listing URLs, allowing clients to specify the page number and limit per page.

## Testing

The project includes a Postman collection (`collection_postman_shorten-url.json`) that can be imported to test the API endpoints.

## Contributing

- This project is part of the [roadmap.sh](https://roadmap.sh/projects/url-shortening-service) backend projects series.
- Created by [jaygaha](https://github.com/jaygaha)