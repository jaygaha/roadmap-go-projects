package handlers

import (
	"encoding/json"
	"html/template"
	"net/http"
	"path/filepath"

	"github.com/jaygaha/roadmap-go-projects/unit-converter-app/converter"
)

func WelcomeHandler(w http.ResponseWriter, r *http.Request) {
	renderGetWelcomePage(w)
}

func renderGetWelcomePage(w http.ResponseWriter) {
	templatesDir := "templates"
	tmplPath := filepath.Join(templatesDir, "layout.tmpl")
	lengthFormPath := filepath.Join(templatesDir, "forms", "length_form.tmpl")
	weightFormPath := filepath.Join(templatesDir, "forms", "weight_form.tmpl")
	temperatureFormPath := filepath.Join(templatesDir, "forms", "temperature_form.tmpl")

	// Parse all templates
	tmpl := template.Must(template.ParseFiles(tmplPath, lengthFormPath, weightFormPath, temperatureFormPath))

	// populate data
	data := map[string]any{
		"Title":        "Unit Converter in Go",
		"Units":        converter.GetLengthUnits(),
		"Weight_Units": converter.GetWeightUnits(),
		"Temp_Units":   converter.GetTemperatureUnits(),
	}

	err := tmpl.ExecuteTemplate(w, "layout.tmpl", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func APIResponse(w http.ResponseWriter, result string, error error, httpStatusCode int) {
	// set content type
	w.Header().Set("Content-Type", "application/json")
	// set http status code
	w.WriteHeader(httpStatusCode)

	// response back with json
	resp := map[string]string{}

	if error != nil {
		resp = map[string]string{"message": "Error while processing data", "error": error.Error()}
	} else {
		resp = map[string]string{"message": "Converter submitted successfully", "result": result}
	}

	json.NewEncoder(w).Encode(resp)
	return
}
