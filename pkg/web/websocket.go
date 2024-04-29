package web

import (
	"chat-bot/pkg/chat"
	"chat-bot/pkg/database/redis"
	"chat-bot/pkg/session"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true // Allow connections from any origin for testing purposes
		},
	}
	connections = &ConnectionRegistry{
		Connections: make(map[string]*websocket.Conn),
		lock:        sync.RWMutex{},
	}
)

// Define the structure for messages
type initialSessionMessage struct {
	Type      string `json:"type"`
	SessionID string `json:"sessionID"`
}

type ConnectionRegistry struct {
	Connections map[string]*websocket.Conn
	lock        sync.RWMutex
}

func (cr *ConnectionRegistry) Add(connID string, conn *websocket.Conn) {
	cr.lock.Lock()
	defer cr.lock.Unlock()
	cr.Connections[connID] = conn
}

func (cr *ConnectionRegistry) Remove(connID string) {
	cr.lock.Lock()
	defer cr.lock.Unlock()
	if _, exists := cr.Connections[connID]; exists {
		// Cleanup session data from Redis when connection closes
		delete(cr.Connections, connID)
		go cr.cleanupSessionData(connID)
	}
}

func (cr *ConnectionRegistry) cleanupSessionData(sessionID string) {
	ctx := context.Background()
	keysPattern := "session:" + sessionID + "*"
	iter := redis.RDB.Scan(ctx, 0, keysPattern, 0).Iterator()
	for iter.Next(ctx) {
		err := redis.RDB.Del(ctx, iter.Val()).Err()
		if err != nil {
			log.Printf("Failed to delete session data for %s: %v", sessionID, err)
		}
	}
	if err := iter.Err(); err != nil {
		log.Printf("Error retrieving keys for session %s: %v", sessionID, err)
	}
}

func HandleWebSocketConnections(w http.ResponseWriter, r *http.Request) {
	sessionID := r.URL.Query().Get("sessionID")
	ctx := context.Background()

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("WebSocket upgrade error:", err)
		return
	}

	// Session validation
	if sessionID == "" {
		sessionID = session.GenerateSessionID()
		redis.RDB.Set(ctx, "session:"+sessionID, "active", time.Minute*30)
		msg := initialSessionMessage{
			Type:      "session_init",
			SessionID: sessionID,
		}
		message, err := json.Marshal(msg)
		if err != nil {
			log.Println("Failed to encode session init message:", err)
			return
		}
		conn.WriteMessage(websocket.TextMessage, message)
	} else {
		// Check Redis if sessionID exists and is valid
		_, err := redis.RDB.Get(ctx, "session:"+sessionID).Result()
		if err != nil {
			log.Println("Session error or session does not exist:", err)
			return
		}
	}

	connections.Add(sessionID, conn)
	go handleMessages(conn, sessionID)
}

func handleMessages(conn *websocket.Conn, sessionID string) {
	defer conn.Close()
	ctx := context.Background()
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("Read error:", err)
			break
		}

		jsonResponse := chat.ProcessMessage(ctx, sessionID, string(message))
		if err := conn.WriteMessage(websocket.TextMessage, []byte(jsonResponse)); err != nil {
			log.Println("Write error:", err)
			break
		}
	}
}
