# Movie Reservation System

A backend service for movie ticket booking built with Go — user authentication, movie/showtime management, seat reservations, and admin reporting.

---

## Tech Stack

| Component | Technology |
|-----------|-----------|
| **Language** | Go 1.25 |
| **Framework** | Gin 1.12 |
| **Database** | PostgreSQL 12+ |
| **Authentication** | JWT HS256 (golang-jwt/jwt/v5) |
| **Password Hashing** | bcrypt (golang.org/x/crypto) |
| **Containerization** | Docker & Docker Compose |

---

## Project Structure

```
movie-reservation-system/
├── cmd/
│   └── main.go                      # Entry point
├── internal/
│   ├── api/
│   │   ├── handlers/                # HTTP handlers
│   │   └── routes.go                # Route definitions + middleware
│   ├── service/                     # Business logic
│   ├── repository/
│   │   ├── migrate.go               # Auto-migration runner
│   │   └── ...                      # Data access layer
│   ├── model/                       # Domain models
│   ├── auth/                        # JWT + password hashing
│   └── config/                      # Configuration
├── pkg/
│   ├── errors.go                    # Custom error types
│   └── validation.go
├── migrations/
│   ├── migrations.go                # Embeds SQL files into binary
│   ├── 000_setup.sql
│   ├── 001_create_users.sql
│   ├── 002_create_movies.sql
│   ├── 003_create_showtimes.sql
│   ├── 004_create_reservations.sql
│   ├── 005_create_seats.sql
│   └── 006_seed_data.sql            # Genres + default admin user
├── postman_collection.json          # Importable API collection
├── docker-compose.yml
└── .env.example
```

---

## Installation & Setup

### Prerequisites

- Go 1.25+
- PostgreSQL 12+ (or Docker)

### 1. Clone

```bash
git clone https://github.com/jaygaha/roadmap-go-projects.git
cd advanced/movie-reservation-system
```

### 2. Configure environment

```bash
cp .env.example .env
```

| Variable | Default | Description |
|----------|---------|-------------|
| `HOST` | `localhost` | Server bind address |
| `PORT` | `8080` | Server port |
| `DB_HOST` | `localhost` | PostgreSQL host |
| `DB_PORT` | `5432` | PostgreSQL port |
| `DB_NAME` | `movie_reservation` | Database name |
| `DB_USER` | `postgres` | Database user |
| `DB_PASSWORD` | `postgres` | Database password |
| `JWT_SECRET` | — | **Required.** Use `openssl rand -base64 32` |
| `JWT_EXPIRE_HOURS` | `24` | Token lifetime in hours |

### 3. Start PostgreSQL

```bash
docker compose up -d
```

### 4. Run

```bash
go run ./cmd/main.go
```

Migrations and seed data run **automatically on first startup**. You will see:

```
[migration] applied: 001_create_users
[migration] applied: 002_create_movies
...
[migration] applied: 006_seed_data
Starting server on localhost:8080
```

Subsequent restarts skip already-applied migrations (tracked in `schema_migrations`).

---

## Default Credentials

Seeded by `006_seed_data.sql`:

| Email | Password | Role |
|-------|----------|------|
| `admin@movieapp.com` | `admin123` | admin |

---

## API Reference

Import `postman_collection.json` into Postman to get all endpoints with pre-configured requests and auto-saved tokens.

**Base URL:** `http://localhost:8080/api`

### Endpoints Summary

| Method | Path | Auth | Description |
|--------|------|------|-------------|
| POST | `/auth/signup` | Public | Register new user |
| POST | `/auth/login` | Public | Login, receive JWT |
| GET | `/users` | User | List users |
| GET | `/users/:id` | User | Get user by ID |
| POST | `/users/:id/promote` | Admin | Promote user to admin |
| GET | `/movies` | User | List movies |
| GET | `/movies/:id` | User | Get movie details |
| GET | `/movies/genres` | User | List all genres |
| GET | `/movies/genre/:genreId` | User | Movies by genre |
| POST | `/movies` | Admin | Create movie |
| PUT | `/movies/:id` | Admin | Update movie |
| DELETE | `/movies/:id` | Admin | Delete movie |
| GET | `/showtimes` | User | List showtimes |
| GET | `/showtimes/:id` | User | Get showtime details |
| GET | `/showtimes/by-date?date=YYYY-MM-DD` | User | Showtimes by date |
| GET | `/showtimes/movie/:movieId` | User | Showtimes by movie |
| POST | `/showtimes` | Admin | Create showtime (auto-generates seats) |
| PUT | `/showtimes/:id` | Admin | Update showtime |
| DELETE | `/showtimes/:id` | Admin | Delete showtime |
| POST | `/reservations` | User | Reserve a seat |
| GET | `/reservations/me` | User | My reservations |
| GET | `/reservations/:id` | User | Get reservation |
| DELETE | `/reservations/:id` | User | Cancel reservation |
| GET | `/reservations` | Admin | All reservations |
| GET | `/reservations/stats` | Admin | Reservation statistics |

### Authentication

All protected endpoints require:

```
Authorization: Bearer <token>
```

The token is returned by `/auth/signup` and `/auth/login`.

---

## Architecture

```
HTTP Handlers  →  Service Layer  →  Repository Layer  →  PostgreSQL
   (routes)      (business logic)    (data access)
```

- Layered separation: handlers → services → repositories
- Interface-based repositories for testability
- All DB operations are context-aware
- Critical operations wrapped in transactions

---

## Acknowledgments

- Challenge: [Movie Reservation System](https://roadmap.sh/projects/movie-reservation-system)
- Part of the [roadmap.sh](https://roadmap.sh/projects) Go backend projects series
- Created by [jaygaha](https://github.com/jaygaha)
