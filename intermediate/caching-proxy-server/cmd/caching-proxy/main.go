package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/jaygaha/roadmap-go-projects/intermediate/caching-proxy-server/internal/cache"
	"github.com/jaygaha/roadmap-go-projects/intermediate/caching-proxy-server/internal/config"
	"github.com/jaygaha/roadmap-go-projects/intermediate/caching-proxy-server/internal/handler"
)

func main() {
	// configure logger
	// it helps to print the file name and line number where the log is printed
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// parse command line flags
	config.ParseFlags()

	// initialize cache with the given TTL
	cache.Init(config.CacheTTL)

	// Handle clear flag
	if config.ClearCache {
		cache.Clear()
		fmt.Println("Cache cleared successfully.")
		return
	}

	// Setup routes
	http.HandleFunc("/clear-cache", handler.HandleClearCache)
	http.HandleFunc("/", handler.HandleProxyRequest)

	// Start the server
	fmt.Println("Server started on port:", config.ServerPort)
	fmt.Println("Upstream server:", config.OriginServerURL.String())
	fmt.Println("Cache TTL:", config.CacheTTL)

	if err := http.ListenAndServe(":"+config.ServerPort, nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}

}
