# Task Tracker CLI app using `Go`

This project is a solution for the task-tracker challenge provided by [roadmap.sh](https://roadmap.sh/projects/task-tracker). The task-tracker is a command-line application that allows users to manage tasks with various statuses such as todo, in-progress, and done.

This solution is built using `Go` and provides a command-line interface for managing tasks.

## Features

- Add new tasks with descriptions.
- List tasks with optional filtering by status.
- Update task descriptions and statuses.
- Delete tasks by ID.
- Mark tasks with dynamic statuses like `mark-todo`, `mark-in-progress`, and `mark-done`.

## Dependencies
- Go (Golang)

## Installation

To install and run the task-tracker, ensure you have Go installed on your system. Then, clone the repository and build the application:

```bash
git clone https://github.com/jaygaha/roadmap-go-projects.git
cd roadmap-go-projects/task-tracker
go mod tidy
go build
```

## Usage

After building the application, you can run it using the following commands:

- `Add a Task`: Add a new task with a description.
  
  ```bash
  ./task-tracker add "New Task Description"
  ```
- `List Tasks`: List all tasks, optionally filtering by status.
  
  ```bash
  ./task-tracker list
  ./task-tracker list todo
  ./task-tracker list in-progress
  ./task-tracker list done
   ```
- `Update a Task`: Update a task's description or status by ID.
  
  ```bash
  ./task-tracker update <id> "Updated Task Description"
  ```
- `Delete a Task`: Delete a task by ID.
  
  ```bash
  ./task-tracker delete <id>
   ```
- `Mark Task Status`: Mark a task with a specific status.

  ```bash
  ./task-tracker mark-in-progress <id>
  ./task-tracker mark-done <id>
    ```
- `Help`: Get help for available commands.

  ```bash
    ./task-tracker --help
   ```

## Project Link

- https://roadmap.sh/projects/task-tracker

