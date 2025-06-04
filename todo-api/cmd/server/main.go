package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/jaygaha/roadmap-go-projects/todo-api/internal/db"
	"github.com/jaygaha/roadmap-go-projects/todo-api/internal/routes"
)

func main() {
	port := flag.String("port", "8800", "Server port")
	flag.Parse()

	// Connect to the database
	db.InitDB()
	defer db.DB.Close()

	mux := routes.RegisterRoutes()

	err := http.ListenAndServe(fmt.Sprintf(":%s", *port), mux)
	if err != nil {
		log.Fatalf("Could not start server: %s\n", err.Error())
	}

	log.Println("Server is running on port", *port)
}
