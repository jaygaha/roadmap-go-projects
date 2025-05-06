package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/jaygaha/roadmap-go-projects/tmdb-cli-tool/models"
)

func FetchMovieList(movieType string, page int) ([]models.Movie, error) {
	tmdbApiUrl := os.Getenv("TMDB_API_URL")
	tmdbApiKey := os.Getenv("TMDB_API_KEY")

	if tmdbApiUrl == "" || tmdbApiKey == "" {
		return nil, fmt.Errorf("TMDB_API_URL or TMDB_API_KEY is not set in the environment variables")
	}

	url := fmt.Sprintf("%s/movie/%s?api_key=%s&page=%d", tmdbApiUrl, movieType, tmdbApiKey, page)

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch movie list: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch movie list. Status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	var movieResponse models.MovieResponse
	err = json.Unmarshal(body, &movieResponse)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal response body: %v", err)
	}

	return movieResponse.Results, nil
}
