package server

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/jaygaha/roadmap-go-projects/intermediate/broadcast-server/internal/config"
)

// Server represents the broadcast server
type Server struct {
	config *config.Config
	hub    *Hub
}

// NewServer creates a new server instance
func NewServer(cfg *config.Config) *Server {
	return &Server{
		config: cfg,
		hub:    NewHub(),
	}
}

// Start starts the server
func (s *Server) Start() error {
	// Start the hub
	go s.hub.Run()

	// Setup HTTP routes
	http.HandleFunc("/ws", s.handleWebSocket)

	// Setup graceful shutdown
	go s.setupGracefulShutdown()

	// Start the server
	addr := ":" + s.config.ServerPort
	log.Printf("WebSocket endpoint: ws://localhost%s/ws", addr)

	return http.ListenAndServe(addr, nil)
}

// handleWebSocket handles WebSocket connections
func (s *Server) handleWebSocket(w http.ResponseWriter, r *http.Request) {
	// Get username from query parameter
	username := r.URL.Query().Get("username")
	if username == "" {
		http.Error(w, "Username is required", http.StatusBadRequest)
		return
	}

	// Upgrade HTTP connection to WebSocket
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket upgrade error: %v", err)
		return
	}

	// Create new client
	client := NewClient(s.hub, conn, username)
	s.hub.register <- client

	// Start client goroutines
	go client.writePump()
	go client.readPump()
}

// setupGracefulShutdown handles graceful shutdown
func (s *Server) setupGracefulShutdown() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	<-c
	log.Println("\nReceived shutdown signal. Gracefully shutting down...")

	// Here you could add cleanup logic like:
	// - Closing database connections
	// - Saving state
	// - Notifying clients of shutdown

	os.Exit(0)
}
