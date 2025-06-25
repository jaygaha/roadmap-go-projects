package utils

import (
	"errors"
	"net/http"
)

// GetCodeFromRequest returns short code from the request URL.
func GetCodeFromRequest(r *http.Request) (string, error) {
	codeStr := r.URL.Path[len("/api/shorten/"):]
	if codeStr == "" {
		return "", errors.New("short code is required")
	}

	return codeStr, nil
}
