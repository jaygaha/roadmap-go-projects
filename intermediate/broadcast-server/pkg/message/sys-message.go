package message

import (
	"encoding/json"
	"fmt"
	"time"
)

// MessageType represents different types of system messages
type MessageType string

// constants of message types
const (
	TypeMessage   MessageType = "message"
	TypeJoin      MessageType = "join"
	TypeLeave     MessageType = "leave"
	TypeUserCount MessageType = "user_count"
)

// Message holds a message sent between client and server
type Message struct {
	Type      MessageType `json:"type"`
	Username  string      `json:"username"`
	Content   string      `json:"content"`
	UserCount int         `json:"user_count,omitempty"`
	CreatedAt time.Time   `json:"created_at"`
}

// NewMessage creates a new message
func NewMessage(msgType MessageType, username, content string) *Message {
	return &Message{
		Type:      msgType,
		Username:  username,
		Content:   content,
		CreatedAt: time.Now(),
	}
}

// ToJson converts message to Json bytes
func (m *Message) ToJson() ([]byte, error) {
	return json.Marshal(m)
}

// FromJson converts Json bytes to message
func FromJson(data []byte) (*Message, error) {
	var msg Message
	// log.Println("FromJson:", string(data))
	err := json.Unmarshal(data, &msg)

	return &msg, err
}

// String returns a formatted string representation of the message
func (m *Message) String() string {
	switch m.Type {
	case TypeJoin:
		return fmt.Sprintf("[%s] %s joined the chat", m.CreatedAt.Format(time.RFC3339), m.Username)
	case TypeLeave:
		return fmt.Sprintf("[%s] %s left the chat", m.CreatedAt.Format(time.RFC3339), m.Username)
	case TypeMessage:
		return fmt.Sprintf("[%s] %s: %s", m.CreatedAt.Format(time.RFC3339), m.Username, m.Content)
	case TypeUserCount:
		return fmt.Sprintf("Users online: %d", m.UserCount)
	default:
		return fmt.Sprintf("[%s] %s: %s", m.CreatedAt.Format(time.RFC3339), m.Username, m.Content)
	}
}

// SplitJSONMessages splits a byte array that might contain multiple JSON objects
// It returns a slice of byte arrays, each containing a valid JSON object
func SplitJSONMessages(data []byte) [][]byte {
	var result [][]byte
	var depth, start int
	inString := false
	escaped := false

	for i, b := range data {
		switch b {
		case '{':
			if !inString {
				if depth == 0 {
					start = i
				}
				depth++
			}
		case '}':
			if !inString {
				depth--
				if depth == 0 {
					result = append(result, data[start:i+1])
				}
			}
		case '"':
			if !escaped {
				inString = !inString
			}
		case '\\':
			escaped = !escaped
			continue
		}
		escaped = false
	}

	// If no valid JSON objects were found, return the original data
	if len(result) == 0 {
		return [][]byte{data}
	}

	return result
}
