# Realtime Leaderboard System

A real-time leaderboard system for ranking and scoring users across various games/activities, powered by Redis sorted sets.

## Project Structure

```
realtime-leaderboard-system/
├── cmd/
│   └── server/
│       └── main.go
├── internal/
│   ├── auth/
│   │   ├── auth.go          # JWT & password hashing
│   │   ├── middleware.go     # JWT auth middleware
│   │   └── ratelimit.go     # Rate limiting middleware
│   ├── database/
│   │   ├── sqlite.go        # SQLite initialization
│   │   └── redis.go         # Redis client setup
│   ├── handlers/
│   │   ├── auth.go          # Register, Login, Profile
│   │   ├── game.go          # Create/List games
│   │   ├── scores.go        # Submit scores, history
│   │   └── leaderboard.go   # Global/game/period leaderboards
│   ├── models/
│   │   ├── user.go
│   │   ├── score.go
│   │   └── game.go
│   └── services/
│       ├── leaderboard.go   # Redis leaderboard operations
│       ├── score.go          # Score queries
│       └── user.go           # User queries
├── pkg/
│   ├── config/
│   │   └── config.go        # Environment configuration
│   └── utils/
│       ├── jwt.go
│       └── validation.go
├── go.mod
├── Dockerfile
├── docker-compose.yml
└── Makefile
```

## Tech Stack

- **Gin Framework** — RESTful API server
- **SQLite + GORM** — Persistent storage for users, scores, and games
- **Redis Sorted Sets** — Efficient real-time leaderboards (ZADD, ZREVRANGE, ZREVRANK)
- **JWT** — Token-based authentication
- **Rate Limiting** — In-memory per-IP rate limiter

## API Endpoints

### Auth
| Method | Path               | Auth | Description        |
|--------|-------------------|------|--------------------|
| POST   | /api/auth/register | No   | Register new user  |
| POST   | /api/auth/login    | No   | Login user         |
| GET    | /api/auth/profile  | Yes  | Get user profile   |

### Games
| Method | Path        | Auth | Description  |
|--------|------------|------|--------------|
| GET    | /api/games | No   | List games   |
| POST   | /api/games | Yes  | Create game  |

### Scores
| Method | Path                | Auth | Description         |
|--------|---------------------|------|---------------------|
| POST   | /api/scores/submit  | Yes  | Submit a score      |
| GET    | /api/scores/history | Yes  | Get score history   |

### Leaderboard
| Method | Path                          | Auth | Description                           |
|--------|-------------------------------|------|---------------------------------------|
| GET    | /api/leaderboard/global       | No   | Top 10 global leaderboard             |
| GET    | /api/leaderboard/game/:gameId | No   | Top 10 for a specific game            |
| GET    | /api/leaderboard/rank/:userId | No   | User's global rank                    |
| GET    | /api/leaderboard/top/:gameId  | No   | Top players (daily/weekly/monthly)    |

**Query params for `/top/:gameId`:** `?period=daily|weekly|monthly`

## Running

```bash
# Start Redis
docker-compose up -d redis

# Run the server
make run

# Or with Docker
make docker-build
make docker-run
```

## Environment Variables

| Variable             | Default          | Description              |
|---------------------|------------------|--------------------------|
| JWT_SECRET          | default-secret   | JWT signing secret       |
| REDIS_ADDR          | localhost:6379   | Redis address            |
| SQLITE_PATH         | ./leaderboard.db | SQLite database path     |
| PORT                | 8080             | Server port              |
| RATE_LIMIT_PER_MINUTE | 10            | Max requests per minute  |


## Contributing

- Challenge: [Realtime Leaderboard](https://roadmap.sh/projects/realtime-leaderboard-system)
- This project is part of the [roadmap.sh](https://roadmap.sh/projects) backend projects series.
- Created by [jaygaha](https://github.com/jaygaha)