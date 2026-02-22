# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
BINARY_NAME=image-processing-service
BINARY_UNIX=$(BINARY_NAME)_unix

# Docker parameters
DOCKER_COMPOSE=docker-compose
COMPOSE_FILE=./infra/docker-compose.yml
DOCKER_IMAGE_NAME=image-processing-service

.PHONY: all build clean test deps up down logs restart build-docker help dev-setup

# Default target
all: deps test build

# Build the application
build:
	$(GOBUILD) -o bin/$(BINARY_NAME) -v ./cmd/api

# Build for Linux
build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o bin/$(BINARY_UNIX) -v ./cmd/api

# Clean build artifacts
clean:
	$(GOCLEAN)
	rm -rf bin/
	rm -f coverage.out coverage.html

# Run tests
test:
	$(GOTEST) -v ./...

# Run tests with coverage
test-coverage:
	$(GOTEST) -v -coverprofile=coverage.out ./...
	$(GOCMD) tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Install dependencies
deps:
	$(GOMOD) download
	$(GOMOD) tidy

# Development setup
dev-setup: deps
	@echo "Setting up development environment..."
	@if [ ! -f .env ]; then cp .env.example .env; echo "Created .env file from .env.example"; fi
	@echo "Development setup complete!"

# Start all services
up:
	@echo "Starting services..."
	$(DOCKER_COMPOSE) -f $(COMPOSE_FILE) up -d
	@echo "Services started. Use 'make logs' to view logs."

# Stop all services
down:
	@echo "Stopping services..."
	$(DOCKER_COMPOSE) -f $(COMPOSE_FILE) down

# Stop all services and remove volumes
down-volumes:
	@echo "Stopping services and removing volumes..."
	$(DOCKER_COMPOSE) -f $(COMPOSE_FILE) down -v

# View logs
logs:
	$(DOCKER_COMPOSE) -f $(COMPOSE_FILE) logs -f

# View logs for specific service
logs-app:
	$(DOCKER_COMPOSE) -f $(COMPOSE_FILE) logs -f app

logs-mongo:
	$(DOCKER_COMPOSE) -f $(COMPOSE_FILE) logs -f mongodb

logs-localstack:
	$(DOCKER_COMPOSE) -f $(COMPOSE_FILE) logs -f localstack

# Restart all services
restart: down up

# Restart specific service
restart-app:
	$(DOCKER_COMPOSE) -f $(COMPOSE_FILE) restart app

# Build Docker image
build-docker:
	$(DOCKER_COMPOSE) -f $(COMPOSE_FILE) build app

# Rebuild and restart services
rebuild: down build-docker up

# Check service status
status:
	$(DOCKER_COMPOSE) -f $(COMPOSE_FILE) ps

# Run the application locally (without Docker)
run-local: build
	@echo "Starting application locally..."
	@if [ ! -f .env ]; then echo "Warning: .env file not found. Using default values."; fi
	./bin/$(BINARY_NAME)

# Format Go code
fmt:
	$(GOCMD) fmt ./...

# Run Go linter
lint:
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run; \
	else \
		echo "golangci-lint not installed. Install it with: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"; \
	fi

# Run security check
sec:
	@if command -v gosec >/dev/null 2>&1; then \
		gosec ./...; \
	else \
		echo "gosec not installed. Install it with: go install github.com/securecodewarrior/gosec/v2/cmd/gosec@latest"; \
	fi

# Database operations
db-shell:
	$(DOCKER_COMPOSE) -f $(COMPOSE_FILE) exec mongodb mongosh -u root -p root@123 --authenticationDatabase admin

# Create S3 bucket in LocalStack
create-bucket:
	@echo "Creating S3 bucket in LocalStack..."
	docker exec image-processing-localstack awslocal s3 mb s3://img-bucket
	@echo "Bucket created successfully!"

# List S3 buckets
list-buckets:
	docker exec image-processing-localstack awslocal s3 ls

# Health check
health:
	@echo "Checking service health..."
	@curl -f http://localhost:8800/health || echo "Service is not responding"

# Install development tools
install-tools:
	@echo "Installing development tools..."
	$(GOCMD) install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	$(GOCMD) install github.com/securecodewarrior/gosec/v2/cmd/gosec@latest
	@echo "Development tools installed!"

# Help
help:
	@echo "Available commands:"
	@echo "  build         - Build the application"
	@echo "  build-linux   - Build for Linux"
	@echo "  clean         - Clean build artifacts"
	@echo "  test          - Run tests"
	@echo "  test-coverage - Run tests with coverage"
	@echo "  deps          - Install dependencies"
	@echo "  dev-setup     - Setup development environment"
	@echo "  up            - Start all services with Docker"
	@echo "  down          - Stop all services"
	@echo "  down-volumes  - Stop services and remove volumes"
	@echo "  logs          - View all service logs"
	@echo "  logs-app      - View application logs"
	@echo "  logs-mongo    - View MongoDB logs"
	@echo "  logs-localstack - View LocalStack logs"
	@echo "  restart       - Restart all services"
	@echo "  restart-app   - Restart application service"
	@echo "  build-docker  - Build Docker image"
	@echo "  rebuild       - Rebuild and restart services"
	@echo "  status        - Check service status"
	@echo "  run-local     - Run application locally"
	@echo "  fmt           - Format Go code"
	@echo "  lint          - Run Go linter"
	@echo "  sec           - Run security check"
	@echo "  db-shell      - Connect to MongoDB shell"
	@echo "  create-bucket - Create S3 bucket in LocalStack"
	@echo "  list-buckets  - List S3 buckets"
	@echo "  health        - Check service health"
	@echo "  install-tools - Install development tools"
	@echo "  help          - Show this help message"