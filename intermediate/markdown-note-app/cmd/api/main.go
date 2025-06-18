package main

import (
	"log"
	"net/http"

	"github.com/jaygaha/roadmap-go-projects/intermediate/markdown-note-app/internal/handlers"
)

func main() {
	// Make sure upload directory exists
	if err := handlers.CreateUploadDir(); err != nil {
		log.Fatalf("Failed to create upload directory: %v", err)
	}

	// Frontend
	http.HandleFunc("/", handlers.HomeHandler)

	// APIs
	http.HandleFunc("/api/notes/check-grammers", handlers.GrammerCheckHander)
	http.HandleFunc("/api/notes/save", handlers.SaveNoteHandler)
	http.HandleFunc("/api/notes", handlers.ListNotesHandler)
	http.HandleFunc("/api/notes/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handlers.RenderNoteHandler(w, r)
		case http.MethodDelete:
			handlers.DeleteNoteHandler(w, r)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
			w.Write([]byte("Method not allowed"))
		}
	})

	// Static file server for serving uploaded files
	http.Handle("/public/", http.StripPrefix("/public/", http.FileServer(http.Dir("./web/static/uploads"))))
	// for js and css
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./web/static"))))

	log.Println("Server started on port 8800")
	log.Fatalf("Server failed to start: %v", http.ListenAndServe(":8800", nil))
}
