# Weather API

A `Go-based` weather API that demonstrates working with 3rd party APIs, caching, and environment variables. This project fetches weather data from `OpenWeatherMap` API and implements Redis caching for improved performance.

## Features

- **3rd Party API Integration**: Integrates with `OpenWeatherMap` API for weather data and geocoding
- **Redis Caching**: Implements caching to reduce API calls and improve response times
- **Environment Variables**: Uses environment variables for configuration management
- **RESTful API**: Provides a clean REST endpoint for weather data retrieval

## Project Structure

```
weather-api/
├── cmd/
│   └── server/
│       └── main.go          # Application entry point
├── internal/
│   ├── config/
│   │   └── app.go           # Configuration management
│   ├── db/
│   │   └── redis.go         # Redis client and operations
│   └── handlers/
│       └── handler.go       # HTTP handlers
├── .env.example             # Environment variables template
├── Makefile                 # Build and run commands
└── go.mod                   # Go module dependencies
```

## Key Learning Concepts

### 1. Working with 3rd Party APIs

This project demonstrates how to:
- Make HTTP requests to external APIs
- Handle API responses and error cases
- Parse JSON data from API responses
- Combine multiple API calls (geocoding + weather data)

**Example API Integration:**
```go
// Geocoding API call
appUrl := fmt.Sprintf("http://api.openweathermap.org/geo/1.0/direct?q=%s&limit=1&appid=%s", city, s.Config.WeatherAPIKey)
resp, err := http.Get(appUrl)

// Weather API call
apiUrl := fmt.Sprintf("https://api.openweathermap.org/data/2.5/weather?lat=%f&lon=%f&appid=%s", lat, lon, s.Config.WeatherAPIKey)
resp, err := http.Get(apiUrl)
```

### 2. Caching with Redis

Implements Redis caching to:
- Reduce API calls to external services
- Improve response times
- Handle cache hits and misses
- Store both geocoding and weather data

**Caching Strategy:**
- City coordinates are cached with key `latlon:{city}`
- Weather data is cached with key `weather:{city}`
- Cache-first approach: check cache before making API calls

**Redis Operations:**
```go
// Set data in cache
func SetKey(rc *RedisClient, key string, value any) error {
    return rc.Client.Set(rc.Client.Context(), key, value, 0).Err()
}

// Get data from cache
func GetKey(rc *RedisClient, key string) (string, error) {
    value, err := rc.Client.Get(rc.Client.Context(), key).Result()
    if err == redis.Nil {
        return "", nil // Cache miss
    }
    return value, err
}
```

### 3. Environment Variables

Demonstrates proper configuration management:
- Sensitive data (API keys) stored in environment variables
- Default values for non-sensitive configuration
- Validation of required environment variables

**Configuration Structure:**
```go
type Config struct {
    Port          int    `env:"PORT" envDefault:"8800"`
    WeatherAPIKey string `env:"WEATHER_API_KEY,required"`
    RedisAddr     string `env:"REDIS_ADDR" envDefault:"localhost:6379"`
    RedisPassword string `env:"REDIS_PASSWORD" envDefault:""`
}
```

## Prerequisites

- Go 1.24 or higher
- Redis server
- OpenWeatherMap API key

## Setup

1. **Clone the repository**
   ```bash
   git clone https://github.com/jaygaha/roadmap-go-projects.git
   cd weather-api
   ```

2. **Install dependencies**
   ```bash
   go mod tidy
   ```

3. **Start Redis server**
   ```bash
   redis-server
   ```

4. **Set up environment variables**
   ```bash
   cp .env.example .env
   ```
   
   Edit `.env` file and add your OpenWeatherMap API key:
   ```
   WEATHER_API_KEY=your_openweathermap_api_key_here
   PORT=8800
   REDIS_ADDR=localhost:6379
   REDIS_PASSWORD=
   ```

5. **Get OpenWeatherMap API Key**
   - Sign up at [OpenWeatherMap](https://openweathermap.org/api)
   - Get your free API key
   - Add it to your `.env` file

## Running the Application

### Using Makefile
```bash
make run
```

### Using Go directly
```bash
go run ./cmd/server/main.go
```

The server will start on port 8800 (or the port specified in your environment variables).

## API Usage

### Get Weather Data

**Endpoint:** `GET /weathers?city={city_name}`

**Example:**
```bash
curl -X GET "http://localhost:8800/weathers?city=tokyo"
```

**Response:**
```json
{
  "coord": {
    "lon": 139.7594,
    "lat": 35.6828
  },
  "weather": [
    {
      "id": 800,
      "main": "Clear",
      "description": "clear sky",
      "icon": "01d"
    }
  ],
  "main": {
    "temp": 298.15,
    "feels_like": 298.74,
    "temp_min": 297.15,
    "temp_max": 299.15,
    "pressure": 1013,
    "humidity": 64
  }
}
```

## How It Works

1. **Request Flow:**
   - Client makes request to `/weathers?city=tokyo`
   - Server checks Redis cache for city coordinates
   - If not cached, calls OpenWeatherMap Geocoding API
   - Caches coordinates for future use
   - Checks Redis cache for weather data
   - If not cached, calls OpenWeatherMap Weather API
   - Caches weather data and returns response

2. **Caching Benefits:**
   - Subsequent requests for the same city are served from cache
   - Reduces API calls to OpenWeatherMap
   - Improves response times
   - Reduces costs (API rate limits)

## Dependencies

- **github.com/go-redis/redis/v8**: Redis client for Go
- **github.com/joho/godotenv**: Load environment variables from .env file
- **github.com/caarlos0/env**: Parse environment variables into structs

## Error Handling

The application handles various error scenarios:
- Missing city parameter
- Invalid city names
- API failures
- Redis connection issues
- JSON parsing errors

## Best Practices Demonstrated

1. **Separation of Concerns**: Clear separation between handlers, database operations, and configuration
2. **Error Handling**: Proper error handling and HTTP status codes
3. **Environment Configuration**: Secure handling of API keys and configuration
4. **Caching Strategy**: Efficient caching to reduce external API calls
5. **Code Organization**: Clean project structure following Go conventions

## Future Enhancements

- Add cache expiration times
- Implement rate limiting
- Add more weather endpoints
- Add unit tests
- Add Docker support
- Add logging middleware

## Project Link

- [Weather API](https://roadmap.sh/projects/weather-api-wrapper-service)

## Acknowledgments

- Part of the Go programming language learning roadmap projects
- Created by [jaygaha](https://github.com/jaygaha)