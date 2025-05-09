package main

import (
	"fmt"
	"net/http"

	"github.com/jaygaha/roadmap-go-projects/personal-blog-app/routes"
)

func main() {
	// app routes
	routes.WebRoutes()

	// start server
	fmt.Println("Server started on port 8800")
	if err := http.ListenAndServe(":8800", nil); err != nil {
		fmt.Printf("Server error: %v", err)
	}
}
