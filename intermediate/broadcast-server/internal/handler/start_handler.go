package handler

import (
	"flag"
	"log"
	"os"

	"github.com/jaygaha/roadmap-go-projects/intermediate/broadcast-server/internal/config"
	"github.com/jaygaha/roadmap-go-projects/intermediate/broadcast-server/internal/server"
)

func StartServerCommand() {
	var port string

	// Parse flags for server command
	serverFlgs := flag.NewFlagSet("start", flag.ExitOnError)
	serverFlgs.StringVar(&port, "port", "6000", "Port to run the server on")
	serverFlgs.Parse(os.Args[2:])

	// Map to config
	config := &config.Config{
		ServerPort: port,
	}

	log.Printf("Server port: %s...\n", config.ServerPort)
	log.Println("Press Ctrl+C to exit")

	srv := server.NewServer(config)
	if err := srv.Start(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
