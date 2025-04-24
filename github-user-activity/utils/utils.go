package utils

import (
	"fmt"

	"github.com/jaygaha/roadmap-go-projects/github-user-activity/api"
)

func FormatAndPrintEvents(activities []api.GitHubActivityStrct, filter string) {

	if filter != "" {
		activities = FilterEvents(activities, filter)
	}

	for _, activity := range activities {
		var message string

		switch activity.Type {
		case "PushEvent":
			totalCommits := len(activity.Payload.Commits)
			message = fmt.Sprintf("Pushed %d commits to %s", totalCommits, activity.Repo.Name)
		case "CreateEvent":
			message = fmt.Sprintf("Created %s %s in %s", activity.Payload.RefType, activity.Payload.Ref, activity.Repo.Name)
		case "DeleteEvent":
			message = fmt.Sprintf("Deleted %s %s from %s", activity.Payload.RefType, activity.Payload.Ref, activity.Repo.Name)
		case "PullRequestEvent":
			message = fmt.Sprintf("Pull request %s in %s", activity.Payload.Action, activity.Repo.Name)
		case "ForkEvent":
			message = fmt.Sprintf("Forked %s", activity.Repo.Name)
		case "IssuesEvent":
			message = fmt.Sprintf("Issues %s in %s", activity.Payload.Action, activity.Repo.Name)
		case "WatchEvent":
			message = fmt.Sprintf("Starred %s", activity.Repo.Name)
		default:
			message = fmt.Sprintf("[TODO] %s in %s", activity.Type, activity.Repo.Name)
		}

		fmt.Printf("- %s\n", message)
	}
}

func FilterEvents(events []api.GitHubActivityStrct, filter string) []api.GitHubActivityStrct {
	var filteredEvents []api.GitHubActivityStrct

	for _, event := range events {
		if event.Type == filter {
			filteredEvents = append(filteredEvents, event)
		}
	}

	return filteredEvents
}
