package client

import (
	"bufio"
	"fmt"
	"log"
	"net/url"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/gorilla/websocket"
	"github.com/jaygaha/roadmap-go-projects/intermediate/broadcast-server/pkg/message"
)

// Client represents a WebSocket client
type Client struct {
	host     string
	port     string
	username string
	conn     *websocket.Conn
	done     chan struct{}
}

// NewClient creates a new client instance
func NewClient(host, port, username string) *Client {
	return &Client{
		host:     host,
		port:     port,
		username: username,
		done:     make(chan struct{}),
	}
}

// Connect connects to the server and starts the client
func (c *Client) Connect() error {
	// Build WebSocket URL
	u := url.URL{
		Scheme:   "ws",
		Host:     c.host + ":" + c.port,
		Path:     "/ws",
		RawQuery: "username=" + url.QueryEscape(c.username),
	}

	// Connect to WebSocket
	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		return fmt.Errorf("failed to connect to WebSocket: %w", err)
	}
	c.conn = conn

	// Setup graceful shutdown
	c.setupGracefulShutdown()

	// Start goroutines
	go c.readMessages()
	go c.writeMessages()

	// Wait for done signal
	<-c.done
	return nil
}

// readMessages reads messages from the WebSocket connection
func (c *Client) readMessages() {
	defer func() {
		c.conn.Close()
		close(c.done)
	}()

	for {
		_, messageBytes, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket error: %v", err)
			}
			return
		}
		// log.Printf("Received message: %s", messageBytes)

		// Split the message if it contains multiple JSON objects
		messages := message.SplitJSONMessages(messageBytes)
		for _, msgBytes := range messages {
			// Parse and display message
			msg, err := message.FromJson(msgBytes)
			if err != nil {
				log.Printf("Error parsing message: %v", err)
				continue
			}

			// Don't display our own messages
			if msg.Username == c.username && msg.Type == message.TypeMessage {
				continue
			}

			fmt.Printf("\r%s\n> ", msg.String())
		}
	}
}

// writeMessages reads input from stdin and sends messages
func (c *Client) writeMessages() {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("> ")

	for scanner.Scan() {
		text := strings.TrimSpace(scanner.Text())

		if text == "" {
			fmt.Print("> ")
			continue
		}

		if text == "quit" || text == "exit" {
			break
		}

		// Create and send message
		msg := message.NewMessage(message.TypeMessage, c.username, text)
		msgBytes, err := msg.ToJson()
		if err != nil {
			log.Printf("Error creating message: %v", err)
			fmt.Print("> ")
			continue
		}

		err = c.conn.WriteMessage(websocket.TextMessage, msgBytes)
		if err != nil {
			log.Printf("Error sending message: %v", err)
			break
		}

		fmt.Print("> ")
	}

	// Close connection
	c.conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
	time.Sleep(time.Second)
	close(c.done)
}

// setupGracefulShutdown handles graceful shutdown
func (c *Client) setupGracefulShutdown() {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-ch
		fmt.Println("\nDisconnecting...")
		if c.conn != nil {
			c.conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			c.conn.Close()
		}
		close(c.done)
	}()
}
