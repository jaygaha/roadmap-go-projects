package handlers

import (
	"encoding/json"
	"net/http"
	"net/url"
)

// TypeGrammerCheckRequest is the request body for the grammer check endpoint
type TypeGrammerCheckRequest struct {
	Text     string `json:"text"`
	Language string `json:"language"`
}

// Replacement represents a single replacement suggestion.
type Replacement struct {
	Value string `json:"value"`
}

// TypeGrammerCheckResponse is the response body for the grammer check endpoint
type TypeGrammerCheckResponse struct {
	Matches []struct {
		Message      string        `json:"message"`
		Offset       int           `json:"offset"`
		Length       int           `json:"length"`
		Replacements []Replacement `json:"replacements"`
		Rule         struct {
			Id        string `json:"id"`
			IssueType string `json:"issueType"`
		} `json:"rule"`
	} `json:"matches"`
}

func GrammerCheckHander(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse the JSON request body
	var request TypeGrammerCheckRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid JSON input", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// Transform JSON data into x-www-form-urlencoded format
	formData := url.Values{}
	// Concert text into json
	txtData := map[string]string{"text": request.Text}
	jsonData, _ := json.Marshal(txtData)
	formData.Set("data", string(jsonData))
	formData.Set("language", request.Language)

	// Send a POST request to the LanguageTool API
	languageToolAPI := "https://api.languagetool.org/v2/check"
	resp, err := http.PostForm(languageToolAPI, formData)
	if err != nil {
		http.Error(w, "Error sending request to third party", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// Read and parse the response from LanguageTool API
	var response TypeGrammerCheckResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		http.Error(w, "Error decoding response from third party", http.StatusInternalServerError)
		return
	}

	// Set the Content-Type header to application/json
	w.Header().Set("Content-Type", "application/json")
	// Set the response status code received from the third-party API
	w.WriteHeader(resp.StatusCode)
	// Write the JSON response to the client
	respData := map[string]any{"message": "Grammer check done successfully", "data": response}
	json.NewEncoder(w).Encode(respData)
}
