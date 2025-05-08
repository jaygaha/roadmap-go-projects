package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/jaygaha/roadmap-go-projects/unit-converter-app/converter"
)

// Parse JSON request body
var requestWeightData struct {
	Weight   string `json:"converter_weight"`
	FromUnit string `json:"converter_unit_from"`
	ToUnit   string `json:"converter_unit_to"`
}

func WeightHandler(w http.ResponseWriter, r *http.Request) {
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
	if err := decoder.Decode(&requestWeightData); err != nil {
		APIResponse(w, "", fmt.Errorf("Invalid JSON payload"), http.StatusBadRequest)
		return
	}

	weightStr := requestWeightData.Weight
	fromUnit := requestWeightData.FromUnit
	toUnit := requestWeightData.ToUnit

	if weightStr == "" || fromUnit == "" || toUnit == "" {
		APIResponse(w, "", fmt.Errorf("All fields are required"), http.StatusUnprocessableEntity)
		return
	}

	// convert str to float64
	weight, err := strconv.ParseFloat(weightStr, 64)
	if err != nil {
		APIResponse(w, "", fmt.Errorf("Invalid length value"), http.StatusUnprocessableEntity)
		return
	}

	// convert weight
	convertedWeight, err := converter.ConvertWeight(weight, fromUnit, toUnit)
	if err != nil {
		APIResponse(w, "", fmt.Errorf("Invalid unit"), http.StatusUnprocessableEntity)
		return
	}
	// convert result to string with 2 decimal places
	resultStr := strconv.FormatFloat(convertedWeight, 'f', 2, 64)
	responseStr := fmt.Sprintf("%s%s = %s%s", weightStr, fromUnit, resultStr, toUnit)

	APIResponse(w, responseStr, nil, http.StatusOK)
}
