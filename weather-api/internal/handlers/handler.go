package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/jaygaha/roadmap-go-projects/weather-api/internal/config"
	"github.com/jaygaha/roadmap-go-projects/weather-api/internal/db"
)

type Server struct {
	RedisClient *db.RedisClient
	Config      *config.Config
}

// GetWeatherData retrieves weather data for a given city.
func (s *Server) GetWeatherData(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	// get city from query params
	city := r.URL.Query().Get("city")
	if city == "" {
		http.Error(w, "City parameter is missing", http.StatusBadRequest)
		return
	}

	// Convert city to latitude and longitude
	lat, lon, err := s.getLatLon(city)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	apiUrl := fmt.Sprintf("https://api.openweathermap.org/data/2.5/weather?lat=%f&lon=%f&appid=%s", lat, lon, s.Config.WeatherAPIKey)
	redisKey := fmt.Sprintf("weather:%s", city)

	// Check if weather data is in cache
	weatherValue, err := db.GetKey(s.RedisClient, redisKey)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// If weather data is in cache, return it
	if weatherValue != "" {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(weatherValue))
		return
	}

	// If weather data is not in cache, make a request to the API
	resp, err := http.Get(apiUrl)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// Parse the response
	var data map[string]any
	err = json.NewDecoder(resp.Body).Decode(&data)

	// Set weather data in cache
	jsonValue, err := json.Marshal(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = db.SetKey(s.RedisClient, redisKey, jsonValue)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return the weather data
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

// getLatLon converts a city name to latitude and longitude.
func (s *Server) getLatLon(city string) (float64, float64, error) {
	appUrl := fmt.Sprintf("http://api.openweathermap.org/geo/1.0/direct?q=%s&limit=1&appid=%s", city, s.Config.WeatherAPIKey)
	redisKey := fmt.Sprintf("latlon:%s", city)

	// Check if city is in cache
	cityValue, err := db.GetKey(s.RedisClient, redisKey)
	if err != nil {
		return 0, 0, err
	}

	// If city is in cache, return it
	if cityValue != "" {
		// Convert json value to lat and lon
		var data map[string]any
		err = json.Unmarshal([]byte(cityValue), &data)
		if err != nil {
			return 0, 0, err
		}

		lat := data["lat"].(float64)
		lon := data["lon"].(float64)

		return lat, lon, nil
	}

	// If city is not in cache, make a request to the API
	resp, err := http.Get(appUrl)
	if err != nil {
		return 0, 0, err
	}
	defer resp.Body.Close()

	// Parse the response
	var locations []map[string]any
	err = json.NewDecoder(resp.Body).Decode(&locations)
	if err != nil {
		return 0, 0, err
	}

	if len(locations) == 0 {
		return 0, 0, fmt.Errorf("no location found for city: %s", city)
	}

	// Get the latitude and longitude from the first result
	lat := locations[0]["lat"].(float64)
	lon := locations[0]["lon"].(float64)

	// Set lat and lon in cache as json
	jsonValue, err := json.Marshal(locations[0])
	if err != nil {
		return 0, 0, err
	}

	err = db.SetKey(s.RedisClient, redisKey, jsonValue)
	if err != nil {
		return 0, 0, err
	}

	return lat, lon, nil
}
