# Fitness Workout Tracker

A simple HTTP API to manage workouts with authentication, exercises catalog, CRUD for workouts, and summary reports. Built with Go net/http and SQLite, documented via Swagger.

## Features
- User registration and JWT-based login
- Seeded exercises catalog
- Workouts CRUD with exercises per workout
- Reports by date range (start_date, end_date)
- Bearer token security documented in Swagger UI

## Tech Stack
- Go 1.22+
- net/http standard library
- SQLite (mattn/go-sqlite3)
- swaggo/swag + http-swagger for API docs

## Installation

1. Clone the repository:
   ```
   git clone https://github.com/jaygaha/roadmap-go-projects.git
   cd intermediate/fitness-workout-tracker
   ```
2. Install dependencies:
   ```
   go mod tidy
   ```

## Getting Started

1. Copy env and set JWT secret:
   - `cp .env.example .env`
   - Set JWT_SECRET in .env
2. Run the server:
   - `go run ./cmd/api/main.go`
3. API base path:
   - `http://localhost:8800/api/v1`
4. Swagger UI:
   - `http://localhost:8800/swagger/index.html`
   - Click Authorize and enter: Bearer <token>

## Endpoints (Base /api/v1)
- POST /auth/register
- POST /auth/login
- GET /exercises (requires Bearer)
- POST /workouts (requires Bearer)
- GET /workouts (requires Bearer)
- GET /workouts/{id} (requires Bearer)
- PUT /workouts/{id} (requires Bearer)
- DELETE /workouts/{id} (requires Bearer)
- GET /workouts/reports?start_date=YYYY-MM-DD&end_date=YYYY-MM-DD (requires Bearer)

## Testing
- Run all tests:
  - go test ./...
- Coverage includes:
  - Route e2e tests using httptest
  - Middleware tests (JSON header, JWT auth)
  - Model validation tests
  - Database initialization and seeding test

## Notes
- The server loads .env if present; if missing, environment variables are used.
- JWT_SECRET must be set; server exits if missing.

## Contributing

- This project is part of the [roadmap.sh](https://roadmap.sh/projects/fitness-workout-tracker) backend projects series.
- Created by [jaygaha](https://github.com/jaygaha)

