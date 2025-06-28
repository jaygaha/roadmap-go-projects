package cmd

import (
	"os"

	"github.com/jaygaha/roadmap-go-projects/intermediate/broadcast-server/internal/handler"
)

func Execute() {
	if len(os.Args) < 2 {
		handler.PrintUsage()
		return
	}

	command := os.Args[1]
	switch command {
	case "start":
		// Handle the start command
		handler.StartServerCommand()
		return
	case "connect":
		// Handle the connect command
		handler.ConnectClientCommand()
	default:
		handler.PrintUsage()
		return
	}
}
