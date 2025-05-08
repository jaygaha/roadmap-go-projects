package main

import (
	"log"
	"net/http"

	"github.com/jaygaha/roadmap-go-projects/unit-converter-app/handlers"
)

func main() {
	// serve static files
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	// handle routes
	http.HandleFunc("/", handlers.WelcomeHandler)

	// API routes
	http.HandleFunc("/api/length", handlers.LengthHandler)
	http.HandleFunc("/api/weight", handlers.WeightHandler)
	http.HandleFunc("/api/temperature", handlers.TemperatureHandler)

	// start server
	log.Println("Server started on port 8800")
	if err := http.ListenAndServe(":8800", nil); err != nil {
		log.Println("Server failed to start due to: ", err)
	}
}
