# Markdown Note-taking API

A simple RESTful API for managing markdown notes built with Go. This application allows users to upload markdown files, check grammar, save notes, and render them as HTML.

## Features

- **File Upload**: Upload markdown (.md) files to the server
- **Grammar Check**: Check grammar of note content using LanguageTool API
- **Note Management**: Save, list, and delete notes
- **Markdown Rendering**: Convert markdown files to HTML for display
- **Web Interface**: Simple frontend for interacting with the API

## Project Structure

```
markdown-note-app/
├── cmd/api/main.go          # Application entry point
├── internal/handlers/       # HTTP handlers
│   ├── grammer_handler.go   # Grammar checking functionality
│   ├── home_handler.go      # Frontend handler
│   └── note_handler.go      # Note CRUD operations
├── web/                     # Frontend assets
│   ├── static/             # CSS, JS, and uploaded files
│   └── templates/          # HTML templates
├── go.mod                   # Go module dependencies
└── Makefile                # Build automation
```

## API Endpoints

### Notes Management

- **POST** `/api/notes/save` - Upload and save a markdown file
- **GET** `/api/notes` - List all saved notes
- **GET** `/api/notes/{filename}` - Render a specific note as HTML
- **DELETE** `/api/notes/{filename}` - Delete a specific note

### Grammar Check

- **POST** `/api/notes/check-grammers` - Check grammar of text content

### Frontend

- **GET** `/` - Web interface for the application
- **GET** `/static/*` - Serve static assets (CSS, JS)
- **GET** `/public/*` - Serve uploaded files

## Installation

### Prerequisites

- Go 1.24.0 or higher
- Internet connection (for grammar checking feature)

### Setup

1. Clone the repository:
```bash
git clone https://github.com/jaygaha/roadmap-go-projects.git
cd roadmap-go-projects/markdown-note-app
```

2. Install dependencies:
```bash
go mod tidy
```

3. Run the application:
```bash
make dev
# OR
go run cmd/api/main.go
```

The server will start on `http://localhost:8800`

## API Usage

### Upload a Note

```bash
curl -X POST -F "note=@example.md" http://localhost:8800/api/notes/save
```

### List All Notes

```bash
curl http://localhost:8800/api/notes
```

### Get Note as HTML

```bash
curl http://localhost:8800/api/notes/example_20240101120000.md
```

### Check Grammar

```bash
curl -X POST -H "Content-Type: application/json" \
  -d '{"text":"This are a test.","language":"en-US"}' \
  http://localhost:8800/api/notes/check-grammers
```

### Delete a Note

```bash
curl -X DELETE http://localhost:8800/api/notes/example_20240101120000.md
```

## Dependencies

- **goldmark**: Markdown to HTML converter
- **LanguageTool API**: External service for grammar checking

## File Storage

Uploaded markdown files are stored in the `web/static/uploads/` directory with timestamps appended to ensure unique filenames.

## Response Format

All API endpoints return JSON responses in the following format:

```json
{
  "message": "Operation description",
  "data": "Response data (when applicable)"
}
```

## Error Handling

The API returns appropriate HTTP status codes:

- `200` - Success
- `201` - Created (for file uploads)
- `400` - Bad Request (invalid input)
- `404` - Not Found (file doesn't exist)
- `405` - Method Not Allowed
- `500` - Internal Server Error

## Contributing

- This project is part of the [roadmap.sh](https://roadmap.sh/projects/markdown-note-taking-app) backend projects series.
- Created by [jaygaha](https://github.com/jaygaha)