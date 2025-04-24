package cmd

import (
	"errors"
	"flag"
	"fmt"
	"log"

	"github.com/jaygaha/roadmap-go-projects/github-user-activity/api"
	"github.com/jaygaha/roadmap-go-projects/github-user-activity/utils"
)

func ParseCommand() (string, string, error) {
	// Define flags for the command
	filterFlag := flag.String("filter", "", "Filter by event type")
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage: %s [options] <github-username>\n\n", "github-activity")
		fmt.Fprintf(flag.CommandLine.Output(), "Options:\n")
		flag.PrintDefaults()
		fmt.Fprintf(flag.CommandLine.Output(), "\nExample:\n")
		fmt.Fprintf(flag.CommandLine.Output(), "  github-activity jaygaha\n")
		fmt.Fprintf(flag.CommandLine.Output(), "  github-activity -filter=CreateEvent jaygaha\n")
	}

	// Parse the command-line arguments
	flag.Parse()

	// Get the remaining arguments (command and arguments)
	args := flag.Args()

	// Validate arguments
	if len(args) == 0 {
		flag.Usage()
		return "", "", errors.New("GitHub username is required")
	}

	// Return username and filter value
	return args[0], *filterFlag, nil
}

func ExecuteCommand() error {
	// Parse command-line arguments
	username, filter, err := ParseCommand()
	if err != nil {
		return err
	}

	log.Printf("GitHub Username: %s\n", username)
	// check if username is valid
	if err = api.GetUser(username); err != nil {
		return err
	}

	events, err := api.GetGitHubUserActivity(username)
	if err != nil {
		return err
	}

	// format events and print
	utils.FormatAndPrintEvents(events, filter)

	return nil
}
