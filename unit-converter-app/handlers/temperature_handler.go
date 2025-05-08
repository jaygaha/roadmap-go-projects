package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/jaygaha/roadmap-go-projects/unit-converter-app/converter"
)

// parse request body
var requestTempData struct {
	Temp     string `json:"converter_temp"`
	FromUnit string `json:"converter_unit_from"`
	ToUnit   string `json:"converter_unit_to"`
}

func TemperatureHandler(w http.ResponseWriter, r *http.Request) {
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
	if err := decoder.Decode(&requestTempData); err != nil {
		APIResponse(w, "", fmt.Errorf("Invalid JSON payload"), http.StatusBadRequest)
		return
	}

	tempStr := requestTempData.Temp
	fromUnit := requestTempData.FromUnit
	toUnit := requestTempData.ToUnit

	if tempStr == "" || fromUnit == "" || toUnit == "" {
		APIResponse(w, "", fmt.Errorf("All fields are required"), http.StatusUnprocessableEntity)
		return
	}

	// convert str to float64
	temp, err := strconv.ParseFloat(tempStr, 64)
	if err != nil {
		APIResponse(w, "", fmt.Errorf("Invalid length value"), http.StatusUnprocessableEntity)
		return
	}

	// convert temperature
	convertedTemp, err := converter.ConvertTemperature(temp, fromUnit, toUnit)
	if err != nil {
		APIResponse(w, "", fmt.Errorf("Invalid unit"), http.StatusUnprocessableEntity)
		return
	}
	// convert result to string with 2 decimal places
	resultStr := strconv.FormatFloat(convertedTemp, 'f', 2, 64)
	responseStr := fmt.Sprintf("%s%s = %s%s", tempStr, fromUnit, resultStr, toUnit)

	APIResponse(w, responseStr, nil, http.StatusOK)
}
