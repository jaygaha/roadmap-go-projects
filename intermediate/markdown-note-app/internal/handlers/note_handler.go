package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"time"

	"github.com/yuin/goldmark"
)

// upload path
const UPLOAD_PATH = "web/static/uploads"

// CreateUploadDir creates the upload directory if it doesn't exist
func CreateUploadDir() error {
	uploadDir := UPLOAD_PATH
	if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
		if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
			return err
		}
	}

	return nil
}

type Note struct {
	Title    string `json:"title"`
	FileName string `json:"file_name"`
	URL      string `json:"url"`
}

// SaveNoteHandler saves a note
func SaveNoteHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method != http.MethodPost {
		// response in json
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(map[string]string{"message": "Method not allowed"})
		return
	}

	// Parse the form to handle file uploads
	if err := r.ParseMultipartForm(10 << 20); err != nil { // limit up to 10MB
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"message": "Upload size limit exceeds"})
		return
	}

	// Retrieve the file from the form
	file, handler, err := r.FormFile("note")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"message": "Failed to retrieve file"})
		return
	}
	defer file.Close()

	// Validate the file type
	ext := filepath.Ext(handler.Filename)
	if ext != ".md" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"message": "Invalid file type. Only .md files are allowed"})
		return
	}

	// Generate unique file name using timestamp
	ts := time.Now().Format("20060102150405")
	// get filename without extension
	fileNameOnly := handler.Filename[:len(handler.Filename)-len(ext)]
	uniqueFileName := fmt.Sprintf("%s_%s%s", fileNameOnly, ts, ext)

	// Save the file to the server
	savePath := filepath.Join(UPLOAD_PATH, uniqueFileName)
	outFile, err := os.Create(savePath)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"message": "Failed to create file"})
		return
	}
	defer outFile.Close()

	_, err = io.Copy(outFile, file)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"message": "Failed to write file"})
		return
	}

	// Return a success response
	w.WriteHeader(http.StatusCreated)
	// Write the JSON response to the client
	json.NewEncoder(w).Encode(map[string]string{"message": "Note saved successfully"})
}

// ListNotesHandler lists all notes
func ListNotesHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method != http.MethodGet {
		// response in json
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(map[string]string{"message": "Method not allowed"})
		return
	}

	baseURL := "http://localhost:8800/api/notes/render/" // Base URL for accessing files

	// Read all files in the upload directory
	files, err := os.ReadDir(UPLOAD_PATH)
	if err != nil {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(map[string]string{"message": "Failed to read notes"})
		return
	}
	// Collect file names and filter for .md files
	var notes []Note
	for _, file := range files {
		if !file.IsDir() && filepath.Ext(file.Name()) == ".md" {
			name := file.Name()
			title := name[:len(name)-3]
			notes = append(notes, Note{
				Title:    title,
				FileName: name,
				URL:      fmt.Sprintf("%s%s", baseURL, name),
			})
		}
	}

	// Sort files by file name (timestamp)
	sort.Slice(notes, func(i, j int) bool {
		return notes[i].FileName < notes[j].FileName
	})

	// Return a success response with the list of notes
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	respData := map[string]any{"message": "List fetched successfully", "data": notes}
	json.NewEncoder(w).Encode(respData)
}

// RenderNoteHandler renders a note from markdown file to html
func RenderNoteHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method != http.MethodGet {
		// response in json
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(map[string]string{"message": "Method not allowed"})
		return
	}

	// Extract the file name from the URL path
	fileName := r.URL.Path[len("/api/notes/"):]
	if fileName == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"message": "Invalid file name"})
		return
	}

	// Locate the file in the upload directory
	filePath := filepath.Join(UPLOAD_PATH, fileName)
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"message": "File not found"})
		return
	}

	// Read the Markdown file
	markdownContent, err := os.ReadFile(filePath)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"message": "Failed to read file"})
		return
	}

	// Convert Markdown to HTML using a Markdown-to-HTML converter
	// var htmlContent string
	md := goldmark.New()
	var buf bytes.Buffer
	if err := md.Convert(markdownContent, &buf); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"message": "Failed to convert Markdown to HTML"})
		return
	}
	htmlContent := buf.String()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	respData := map[string]any{"message": "Markdown file retrieved successfully", "data": htmlContent}
	json.NewEncoder(w).Encode(respData)
}

// DeleteNoteHandler deletes a note
func DeleteNoteHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method != http.MethodDelete {
		// response in json
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(map[string]string{"message": "Method not allowed"})
		return
	}

	// Extract the file name from the URL path
	fileName := r.URL.Path[len("/api/notes/"):]
	if fileName == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"message": "Invalid file name"})
		return
	}

	// Locate the file in the upload directory
	filePath := filepath.Join(UPLOAD_PATH, fileName)
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"message": "File not found"})
		return
	}

	// Delete the file
	if err := os.Remove(filePath); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"message": "Failed to delete file"})
		return
	}

	// Return a success response
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Note deleted successfully"})
}
