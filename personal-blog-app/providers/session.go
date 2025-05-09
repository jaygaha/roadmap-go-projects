package providers

import (
	"crypto/rand"
	"encoding/base64"
	"net/http"
	"time"
)

// stores active sessions in memory
type Session struct {
	ID        string
	Data      map[string]any
	ExpiresAt time.Time
}

var (
	sessions        = make(map[string]*Session)
	sessionDuration = 30 * time.Minute // Session duration
)

// create a new session and return the session id
func CreateSession(w http.ResponseWriter) *Session {
	sessionId := generateSessionID()

	session := &Session{
		ID:        sessionId,
		Data:      make(map[string]any),
		ExpiresAt: time.Now().Add(sessionDuration),
	}
	sessions[sessionId] = session

	cookie := &http.Cookie{
		Name:    "session_id",
		Value:   sessionId,
		Expires: session.ExpiresAt,
	}

	http.SetCookie(w, cookie)

	return session
}

// get a session by session id
func GetSession(r *http.Request) *Session {
	cookie, err := r.Cookie("session_id")
	if err != nil || cookie.Value == "" {
		return nil
	}

	sessionId := cookie.Value
	session, ok := sessions[sessionId]

	if !ok || session.ExpiresAt.Before(time.Now()) {
		delete(sessions, sessionId)

		return nil
	}
	session.ExpiresAt = time.Now().Add(sessionDuration) // Update the expiration time

	return session
}

// delete a session by session id
func DeleteSession(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_id")
	if err != nil || cookie.Value == "" {
		return
	}

	sessionId := cookie.Value
	delete(sessions, sessionId)

	// Clear the session cookie
	cookie = &http.Cookie{
		Name:    "session_id",
		Value:   "",
		Expires: time.Now(),
	}

	http.SetCookie(w, cookie)
}

// Is Authenticated checks if the user is authenticated
func IsAuthenticated(r *http.Request) bool {
	session := GetSession(r)

	if session == nil {
		return false
	}

	return true
}

// create a random session id string
func generateSessionID() string {
	b := make([]byte, 32)  // 32 bytes for a 256-bit session ID
	_, err := rand.Read(b) // Fill the slice with random bytes
	if err != nil {
		panic(err) // Handle error appropriately in a real application
	}

	return base64.StdEncoding.EncodeToString(b)
}
