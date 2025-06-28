package server

import (
	"log"

	"github.com/jaygaha/roadmap-go-projects/intermediate/broadcast-server/pkg/message"
)

// Hub maintains the set of active clients and broadcasts messages to them
type Hub struct {
	// Registered client
	clients map[*Client]bool

	// Inbound messages from the client
	broadcast chan []byte

	// Register requests from the client
	register chan *Client

	// Unregister requests from the client
	unregister chan *Client
}

// NewHub creates a new hub
func NewHub() *Hub {
	return &Hub{
		clients:    make(map[*Client]bool),
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

// Run starts the hub and handles client registration/unregistration and message broadcasting
func (h *Hub) Run() {
	for {
		// The select statement allows the method to wait on multiple channel operations. It will execute the case that is ready first.
		select {
		// When a client connects, it is received and assigned to the variable client.
		case client := <-h.register:
			h.clients[client] = true // Add the client to the map
			log.Printf("Client connected: %s (Total: %d)", client.GetUsername(), len(h.clients))

			// Send join message to all clients
			joinMsg := message.NewMessage(message.TypeJoin, client.GetUsername(), "")
			if msgBytes, err := joinMsg.ToJson(); err == nil {
				h.broadcastToAll(msgBytes) // Send the join message to all clients
			}

			// Send user count update
			h.sendUserCount()

			// When a client disconnects, it is received and assigned to the variable client.
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client) // Remove the client from the map
				close(client.send)        // Close the send channel to signal the client to stop listening
				log.Printf("Client disconnected: %s (Total: %d)", client.GetUsername(), len(h.clients))

				// Send leave message to all clients
				leaveMsg := message.NewMessage(message.TypeLeave, client.GetUsername(), "")
				if msgBytes, err := leaveMsg.ToJson(); err == nil {
					h.broadcastToAll(msgBytes)
				}

				// Send user count update
				h.sendUserCount()
			}

			// When a message is received, it is broadcasted to all connected clients using the broadcastToAll method.
		case message := <-h.broadcast:
			h.broadcastToAll(message)
		}
	}
}

// broadcastToAll sends a message to all connected clients
func (h *Hub) broadcastToAll(message []byte) {
	// Iterate over all clients and send the message to their send channel
	for client := range h.clients {
		select {
		// Attempt to send the message to the client's send channel
		case client.send <- message:
		default:
			close(client.send)
			delete(h.clients, client)
		}
	}
}

// sendUserCount sends the current user count to all clients
func (h *Hub) sendUserCount() {
	// Create a user count message
	countMsg := &message.Message{
		Type:      message.TypeUserCount,
		UserCount: len(h.clients),
	}

	if msgBytes, err := countMsg.ToJson(); err == nil {
		h.broadcastToAll(msgBytes)
	}
}

// GetClientCount returns the number of connected clients
func (h *Hub) GetClientCount() int {
	return len(h.clients)
}
