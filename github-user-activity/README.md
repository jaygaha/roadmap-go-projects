# GitHub User Activity

A simple command line interface (CLI) to fetch the recent activity of a GitHub user and display it in the terminal.

This project is built using the Go programming language and is a solution to the challenge provided by [roadmap.sh](https://roadmap.sh/projects/github-user-activity).

## Features

- Fetches recent public activity for any GitHub user
- Supports filtering by event type (e.g., PushEvent, CreateEvent, etc.)
- Displays activity in a human-readable format in the terminal

## Installation

1. **Clone the repository:**
   ```bash
   git clone https://github.com/jaygaha/roadmap-go-projects.git
   cd roadmap-go-projects/github-user-activity
   ```
2. **Build the project:**
   ```bash
   go build -o github-activity
   ```

## Usage

Run the CLI with a GitHub username to fetch their recent activity:

```bash
./github-activity <username>
```

Filter by event type using the -filter flag:

```bash
./github-activity -filter=PushEvent <github-username>
```

### Example

```bash
./github-activity jaygaha
./github-activity -filter=PushEvent jaygaha
```

### Supported filters

Please refer to the [`GitHub event types documentation`](https://docs.github.com/en/rest/using-the-rest-api/github-event-types?apiVersion=2022-11-28) for a list of supported event types.

## Project Structure

- **main.go** : Entry point of the application
- **cmd/command.go** : Command-line parsing and execution logic
- **api/github.go** : Functions to interact with the GitHub API
- **api/types.go** : Data structures for GitHub API responses
- **utils/utils.go** : Utility functions for formatting and filtering events

## Contributing

Contributions are welcome! Please open an issue or submit a pull request for any improvements or bug fixes.

## Project Link

- [Challenge Link](https://roadmap.sh/projects/github-user-activity)

Happy coding! 🚀