package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

func GetUser(username string) error {
	// make a request to the github api
	url := fmt.Sprintf("https://api.github.com/users/%s", username)

	// make a request to the github api
	user, err := http.Get(url)
	if err != nil {
		return errors.New("Error making request to github api")
	}
	// defer the response body
	defer user.Body.Close()

	// check if the response is successful
	if user.StatusCode != http.StatusOK {
		return errors.New("Error getting user, provide a valid username")
	}

	return nil
}

// get github user activity
func GetGitHubUserActivity(username string) ([]GitHubActivityStrct, error) {
	url := fmt.Sprintf("https://api.github.com/users/%s/events", username)

	// add query params
	url += "?per_page=50" // default is 30 and max is 100

	// make a request to the github api
	resp, err := http.Get(url)
	if err != nil {
		return nil, errors.New("Error making event request to github api")
	}

	defer resp.Body.Close()

	// check if the response is successful
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("Error getting user events")
	}

	var activities []GitHubActivityStrct

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.New("Error reading response body")
	}

	err = json.Unmarshal(body, &activities)
	if err != nil {
		return nil, errors.New("Error unmarshalling response body")
	}

	return activities, nil
}
