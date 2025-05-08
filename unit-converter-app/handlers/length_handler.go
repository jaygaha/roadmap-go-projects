package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/jaygaha/roadmap-go-projects/unit-converter-app/converter"
)

// Parse JSON request body
var requestData struct {
	Length   string `json:"converter_length"`
	FromUnit string `json:"converter_unit_from"`
	ToUnit   string `json:"converter_unit_to"`
}

// response back with jso
func LengthHandler(w http.ResponseWriter, r *http.Request) {
	// only post method is allowed
	if r.Method != http.MethodPost {
		APIResponse(w, "", fmt.Errorf("Method not allowed"), http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseForm(); err != nil {
		APIResponse(w, "", fmt.Errorf("Internal Server Error"), http.StatusInternalServerError)
		return
	}

	// all form values are required
	// values are passed as json
	// parse json
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&requestData); err != nil {
		APIResponse(w, "", fmt.Errorf("Invalid JSON payload"), http.StatusBadRequest)
		return
	}

	lengthStr := requestData.Length
	fromUnit := requestData.FromUnit
	toUnit := requestData.ToUnit

	if lengthStr == "" || fromUnit == "" || toUnit == "" {
		APIResponse(w, "", fmt.Errorf("All fields are required"), http.StatusUnprocessableEntity)
		return
	}

	// convert str to float64
	length, err := strconv.ParseFloat(lengthStr, 64)
	if err != nil {
		APIResponse(w, "", fmt.Errorf("Invalid length value"), http.StatusUnprocessableEntity)
		return
	}

	// convert length
	result, err := converter.ConvertLength(length, fromUnit, toUnit)
	if err != nil {
		APIResponse(w, "", fmt.Errorf("%s", err.Error()), http.StatusUnprocessableEntity)
		return
	}

	// convert result to string with 2 decimal places
	resultStr := strconv.FormatFloat(result, 'f', 2, 64)
	responseStr := fmt.Sprintf("%s%s = %s%s", lengthStr, fromUnit, resultStr, toUnit)

	APIResponse(w, responseStr, nil, http.StatusOK)
}
