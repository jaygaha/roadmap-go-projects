package handler

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/jaygaha/roadmap-go-projects/intermediate/broadcast-server/internal/client"
)

func ConnectClientCommand() {
	var host, port, username string

	// Parse flags for client command
	clientFlgs := flag.NewFlagSet("connect", flag.ExitOnError)
	clientFlgs.StringVar(&host, "host", "localhost", "Server host to connect to")
	clientFlgs.StringVar(&port, "port", "6000", "Server port to connect to")
	clientFlgs.StringVar(&username, "username", "", "Username to chat with")
	clientFlgs.Parse(os.Args[2:])

	// Prompt for username if not provided
	if username == "" {
		fmt.Print("Enter username: ")
		fmt.Scanln(&username)
	}

	// Validate username
	if username == "" {
		fmt.Println("Username is required")
		os.Exit(1)
	}

	log.Printf("Connecting to server %s:%s as %s...\n", host, port, username)
	fmt.Println("Type your messages and press Enter to send. Type 'quit' to exit.")

	c := client.NewClient(host, port, username)
	if err := c.Connect(); err != nil {
		log.Fatalf("Failed to connect to server: %v", err)
	}
}
