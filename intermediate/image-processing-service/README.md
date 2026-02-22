# Image Processing Service

A robust, scalable image processing service built with Go, featuring upload, processing, and storage capabilities with MongoDB and S3-compatible storage.

## Features

- **Image Upload**: Support for multiple image formats (JPEG, PNG, GIF, WebP, BMP, TIFF)
- **Image Processing**: Various operations including resize, crop, rotate, flip, grayscale, blur, sharpen, brightness, contrast, and compression
- **User Management**: User registration and authentication
- **Storage**: S3-compatible storage with LocalStack for development
- **Database**: MongoDB for metadata storage
- **API**: RESTful API with comprehensive error handling
- **Docker**: Containerized deployment with Docker Compose
- **Monitoring**: Health checks and structured logging

## Architecture

The service follows a clean architecture pattern with the following layers:

- **Handlers**: HTTP request handlers
- **Services**: Business logic layer
- **Models**: Data models and DTOs
- **Database**: MongoDB integration
- **Storage**: S3-compatible file storage
- **Middleware**: Cross-cutting concerns (CORS, logging, error handling)
- **Utils**: Utility functions for validation and responses

## Prerequisites

- Go 1.24+
- Docker and Docker Compose
- Make (optional, for using Makefile commands)

## Quick Start

### 1. Clone and Setup

```bash
cd /path/to/image-processing-service
make dev-setup
```

### 2. Start Services

```bash
# Start all services (MongoDB, LocalStack, Application)
make up

# Check service status
make status

# View logs
make logs
```

### 3. Create S3 Bucket

```bash
make create-bucket
```

### 4. Test the Service

```bash
# Health check
make health

# Or manually
curl http://localhost:8800/health
```

## API Endpoints
This section documents the HTTP API. All responses use a standard envelope:

```json
{
  "success": true,
  "message": "Text message",
  "data": { },
  "error": "optional error text"
}
```

- Authentication
  - Use Bearer JWT in `Authorization` header for all `/api/v1/images/*` routes.
  - Example: `Authorization: Bearer <token>`
- Rate Limiting
  - `POST /api/v1/images/:id/process` allows 10 requests per user per minute. Exceeding returns `429 Too Many Requests`.
- Async Processing
  - Processing requests return `202 Accepted` and run in background workers. Poll `GET /api/v1/images/:id` to check status and processed variants.

### Health Check
- `GET /`  
  - 200 OK
- `GET /health`  
  - 200 OK

### Authentication
- `POST /api/v1/auth/register`  
  - Body:
    ```json
    { "name": "John Doe", "email": "john@example.com", "password": "SecurePass123!" }
    ```
  - Responses:
    - 201 Created on success, returns user object with token (implementation may vary)
    - 400 on validation errors, 409 if email exists
- `POST /api/v1/auth/login`  
  - Body:
    ```json
    { "email": "john@example.com", "password": "SecurePass123!" }
    ```
  - Responses:
    - 200 OK with user object and token
    - 401 on invalid credentials

### Images
- `POST /api/v1/images/upload`  
  - Auth: required  
  - Content-Type: `multipart/form-data` with field `image`  
  - Valid types: JPEG, PNG, GIF, WebP, BMP, TIFF  
  - Responses:
    - 201 Created with image metadata
    - 400 on invalid file type/size
- `GET /api/v1/images/`  
  - Auth: required  
  - Query: `page` (default 1), `limit` (default 10, max 100)  
  - 200 OK with list and pagination info
- `GET /api/v1/images/:id`  
  - Auth: required  
  - 200 OK with image metadata (including status and processed versions) or 404 if not found
- `POST /api/v1/images/:id/process`  
  - Auth: required  
  - Body:
    ```json
    {
      "operation": "resize",
      "parameters": {
        "width": "800",
        "height": "0",
        "format": "jpeg",
        "quality": "90"
      }
    }
    ```
  - Responses:
    - 202 Accepted with current image metadata (processing scheduled)
    - 400 on invalid operation/parameters
    - 404 if image not found
    - 429 if rate limit exceeded
- `GET /api/v1/images/:id/download`  
  - Auth: required  
  - Query: `processed` (optional processed key). If omitted, returns original  
  - 200 OK with image bytes or 404 if not found
- `DELETE /api/v1/images/:id`  
  - Auth: required  
  - 200 OK on success or 404 if not found

### Supported Operations and Parameters
- `resize`
  - parameters: `width` (int), `height` (int). At least one non-zero.
- `crop`
  - parameters: `width` (int > 0), `height` (int > 0). Center crop.
- `rotate`
  - parameters: `angle` in `"90" | "180" | "270"`.
- `flip`
  - parameters: `mode` in `"horizontal" | "vertical"` (default horizontal).
- `grayscale`
  - parameters: none.
- `blur`
  - parameters: `sigma` (float, default 1.5).
- `sharpen`
  - parameters: `sigma` (float, default 1.0).
- `brightness`
  - parameters: `value` (float, typically -100..100).
- `contrast`
  - parameters: `value` (float, typically -100..100).
- `compress`
  - parameters: none (optionally combine with `format` and `quality` for JPEG).

Format conversion
- Any operation may include `format` to set the output format: `jpeg|jpg|png|gif|bmp|tiff`.
- For JPEG, you can provide `quality` `1..100` to control compression.

## API Usage Examples

### Register a User

```bash
curl -X POST http://localhost:8800/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "name": "John Doe",
    "email": "john@example.com",
    "password": "SecurePass123!"
  }'
```

### Upload an Image

```bash
curl -X POST http://localhost:8800/api/v1/images/upload \
  -F "image=@/path/to/your/image.jpg"
```

### Process an Image

```bash
curl -X POST http://localhost:8800/api/v1/images/{image_id}/process \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{
    "operation": "resize",
    "parameters": {
      "width": "800",
      "height": "0",
      "format": "jpeg",
      "quality": "85"
    }
  }'
```

Response: `202 Accepted` and image metadata. Poll image:

```bash
curl -H "Authorization: Bearer <token>" http://localhost:8800/api/v1/images/{image_id}
```

Download processed:

```bash
curl -L -H "Authorization: Bearer <token>" \
  "http://localhost:8800/api/v1/images/{image_id}/download?processed=<processed_key>" \
  -o processed.jpg
```
```

## Development

### Local Development

```bash
# Install dependencies
make deps

# Run tests
make test

# Run with coverage
make test-coverage

# Format code
make fmt

# Run linter (requires golangci-lint)
make lint

# Run security check (requires gosec)
make sec

# Run locally (without Docker)
make run-local
```

### Docker Development

```bash
# Build Docker image
make build-docker

# Rebuild and restart
make rebuild

# View specific service logs
make logs-app
make logs-mongo
make logs-localstack

# Restart specific service
make restart-app
```

### Database Operations

```bash
# Connect to MongoDB shell
make db-shell

# In MongoDB shell:
use image_processing
db.users.find()
db.images.find()
```

### S3 Operations

```bash
# List buckets
make list-buckets

# Create bucket
make create-bucket
```

## Configuration

The service uses environment variables for configuration. Copy `.env.example` to `.env` and modify as needed:

```bash
cp .env.example .env
```

### Key Configuration Variables

- `APP_PORT`: Application port (default: 8800)
- `MONGODB_URI`: MongoDB connection string
- `AWS_ACCESS_KEY_ID`: AWS/LocalStack access key
- `AWS_SECRET_ACCESS_KEY`: AWS/LocalStack secret key
- `AWS_BUCKET`: S3 bucket name
- `S3_ENDPOINT_URL`: S3 endpoint URL (for LocalStack)

## Project Structure

```
.
├── cmd/
│   └── api/
│       └── main.go              # Application entry point
├── internal/
│   ├── config/
│   │   └── app.go               # Configuration management
│   ├── database/
│   │   └── mongodb.go           # MongoDB connection
│   ├── errors/
│   │   └── errors.go            # Custom error types
│   ├── handlers/
│   │   ├── auth_handler.go      # Authentication handlers
│   │   ├── base_handler.go      # Base handler
│   │   └── image_handler.go     # Image processing handlers
│   ├── logger/
│   │   └── logger.go            # Structured logging
│   ├── middleware/
│   │   └── error_handler.go     # Middleware functions
│   ├── models/
│   │   ├── image.go             # Image models
│   │   └── user.go              # User models
│   ├── routes/
│   │   └── api.go               # Route definitions
│   ├── services/
│   │   ├── file_service.go      # File operations
│   │   ├── image_service.go     # Image business logic
│   │   ├── s3_service.go        # S3 operations
│   │   └── user_service.go      # User business logic
│   └── utils/
│       ├── response.go          # API response utilities
│       └── validation.go        # Input validation
├── infra/
│   ├── builds/
│   │   └── Dockerfile           # Application Dockerfile
│   ├── docker-compose.yml       # Docker Compose configuration
│   └── init-scripts/            # Initialization scripts
├── go.mod                       # Go module definition
├── go.sum                       # Go module checksums
├── Makefile                     # Development commands
├── README.md                    # This file
└── .env.example                 # Environment variables template
```

## Standardization

### Code Organization
- **Clean Architecture**: Separated concerns into distinct layers
- **Error Handling**: Custom error types with proper HTTP status codes
- **Validation**: Comprehensive input validation for all endpoints
- **Response Utilities**: Standardized API response format
- **Middleware**: Centralized cross-cutting concerns

### Features Added
- **Image Processing**: Complete image processing pipeline
- **File Management**: Robust file upload and storage
- **User Management**: User registration with validation
- **Health Checks**: Service health monitoring
- **Pagination**: Efficient data retrieval

### Infrastructure
- **Docker Improvements**: Health checks and proper service dependencies
- **Configuration**: Environment-based configuration management
- **Development Tools**: Comprehensive Makefile with development commands
- **Logging**: Structured logging with Zap

### Security
- **Input Validation**: Comprehensive validation for all inputs
- **Password Security**: Strong password requirements and hashing
- **File Type Validation**: Restricted file types for security
- **Size Limits**: File size restrictions

## Troubleshooting

### Docker Issues

```bash
# Check Docker daemon
docker info

# Restart Docker services
make down
make up

# Remove volumes and restart
make down-volumes
make up
```

### Service Issues

```bash
# Check service status
make status

# View logs
make logs

# Health check
make health
```

### Database Issues

```bash
# Connect to MongoDB
make db-shell

# Check MongoDB logs
make logs-mongo
```

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Run tests: `make test`
5. Format code: `make fmt`
6. Run linter: `make lint`
7. Submit a pull request

## License

This project is licensed under the MIT License.
