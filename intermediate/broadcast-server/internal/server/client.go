package server

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/jaygaha/roadmap-go-projects/intermediate/broadcast-server/pkg/message"
)

const (
	// Time allowed to write a message to the peer
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer
	maxMessageSize = 512
)

// websocket.Upgrader is used to upgrade an HTTP connection to a WebSocket connection
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024, // Maximum size of the buffer used for reading messages from the WebSocket connection
	WriteBufferSize: 1024, // Maximum size of the buffer used for writing messages to the WebSocket connection
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins for simplicity
	},
}

// Client represents a connected client
type Client struct {
	hub      *Hub
	conn     *websocket.Conn
	send     chan []byte
	username string
}

// NewClient creates a new client
func NewClient(hub *Hub, conn *websocket.Conn, username string) *Client {
	return &Client{
		hub:      hub,
		conn:     conn,
		send:     make(chan []byte, 256),
		username: username,
	}
}

// readPump pumps messages from the websocket connection to the hub
// The purpose of this method is to handle reading messages from a WebSocket connection.
func (c *Client) readPump() {
	defer func() {
		// defer is used to ensure that the client is unregistered and the connection is closed when the function returns
		c.hub.unregister <- c
		c.conn.Close()
	}()

	c.conn.SetReadLimit(maxMessageSize)              // a limit on the size of messages that can be read from the WebSocket connection.
	c.conn.SetReadDeadline(time.Now().Add(pongWait)) // a deadline for reading messages from the WebSocket connection.
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	}) // a handler for pong messages received from the WebSocket connection and resets the read deadline.

	for {
		// Read message from the WebSocket connection
		_, messageByte, err := c.conn.ReadMessage() // ReadMessage reads a message from the WebSocket connection.
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("unexpected close error: %v", err)
			}
			break
		}

		// Unmarshal the messageByte into a Message struct
		msg, err := message.FromJson(messageByte)
		if err != nil {
			log.Printf("error unmarshalling message: %v", err)
			continue
		}

		// Set the username from the client connection
		msg.Username = c.username
		msg.CreatedAt = time.Now()

		// Broadcast the message to hub
		if msgBytes, err := msg.ToJson(); err == nil {
			c.hub.broadcast <- msgBytes
		}
	}
}

// writePump pumps messages from the hub to the websocket connection
// The purpose of this method is to handle writing messages to the WebSocket connection and managing periodic ping messages to
// keep the connection alive.
func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod) // used to send periodic ping messages to the WebSocket connection to keep it alive and check if the connection is still active.
	defer func() {
		// Stop the ticker and close the connection when the function returns
		ticker.Stop()
		c.conn.Close()
	}()

	// infinite loop that listens for messages from the c.send channel and sends them to the WebSocket connection.
	for {
		select {
		case message, ok := <-c.send: // listens for messages from the c.send channel.
			c.conn.SetWriteDeadline(time.Now().Add(writeWait)) // If the write operation does not complete before this deadline, it will fail.
			if !ok {                                           // If the channel is closed, send a close message to the peer.
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			// retrieve a writer for the WS connection that will be used to send a text message
			//  If there is an error obtaining the writer, the method exits.
			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message) // received message is written to the WS connection using the writer.

			// Add queued chat messages to the current websocket message.
			// This handle any additional messages that may have been queued in the send channel
			n := len(c.send)
			for range n {
				w.Write([]byte{'\n'})
				w.Write(<-c.send)
			}

			// After all messages have been written, the writer is closed.
			if err := w.Close(); err != nil {
				return
			}

			// Sending ping messages to the peer
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait)) // Set the write deadline for the next write operation.
			// If the write operation does not complete before this deadline, it will fail.
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// GetUsername returns the client's username
func (c *Client) GetUsername() string {
	return c.username
}
