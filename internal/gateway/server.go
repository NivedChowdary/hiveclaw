package gateway

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/nanilabs/hiveclaw/internal/session"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins for now
	},
}

// Message types for WebSocket protocol
type MessageType string

const (
	TypeRequest  MessageType = "req"
	TypeResponse MessageType = "res"
	TypeEvent    MessageType = "event"
)

// WSMessage is the base WebSocket message structure
type WSMessage struct {
	Type    MessageType     `json:"type"`
	ID      string          `json:"id,omitempty"`
	Method  string          `json:"method,omitempty"`
	Event   string          `json:"event,omitempty"`
	Params  json.RawMessage `json:"params,omitempty"`
	Payload json.RawMessage `json:"payload,omitempty"`
	OK      *bool           `json:"ok,omitempty"`
	Error   *WSError        `json:"error,omitempty"`
}

type WSError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// Client represents a connected WebSocket client
type Client struct {
	ID        string
	Conn      *websocket.Conn
	Send      chan []byte
	Gateway   *Gateway
	SessionID string
	Role      string // "operator" or "node"
}

// Gateway is the main WebSocket server
type Gateway struct {
	Port       int
	ConfigPath string
	Clients    map[string]*Client
	Sessions   *session.Manager
	mu         sync.RWMutex
	hub        *Hub
}

// Hub manages all client connections
type Hub struct {
	Clients    map[*Client]bool
	Broadcast  chan []byte
	Register   chan *Client
	Unregister chan *Client
	mu         sync.RWMutex
}

func newHub() *Hub {
	return &Hub{
		Clients:    make(map[*Client]bool),
		Broadcast:  make(chan []byte),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
	}
}

func (h *Hub) run() {
	for {
		select {
		case client := <-h.Register:
			h.mu.Lock()
			h.Clients[client] = true
			h.mu.Unlock()
			log.Printf("Client connected: %s", client.ID)

		case client := <-h.Unregister:
			h.mu.Lock()
			if _, ok := h.Clients[client]; ok {
				delete(h.Clients, client)
				close(client.Send)
			}
			h.mu.Unlock()
			log.Printf("Client disconnected: %s", client.ID)

		case message := <-h.Broadcast:
			h.mu.RLock()
			for client := range h.Clients {
				select {
				case client.Send <- message:
				default:
					close(client.Send)
					delete(h.Clients, client)
				}
			}
			h.mu.RUnlock()
		}
	}
}

// New creates a new Gateway instance
func New(port int, configPath string) *Gateway {
	return &Gateway{
		Port:       port,
		ConfigPath: configPath,
		Clients:    make(map[string]*Client),
		Sessions:   session.NewManager(),
		hub:        newHub(),
	}
}

// Start starts the gateway server
func (g *Gateway) Start() error {
	go g.hub.run()

	// WebSocket endpoint
	http.HandleFunc("/ws", g.handleWebSocket)

	// REST API endpoints
	http.HandleFunc("/api/health", g.handleHealth)
	http.HandleFunc("/api/sessions", g.handleSessions)
	http.HandleFunc("/api/chat", g.handleChat)

	// Serve static files (React frontend)
	fs := http.FileServer(http.Dir("./web/frontend/dist"))
	http.Handle("/", fs)

	addr := fmt.Sprintf(":%d", g.Port)
	log.Printf("🐝 HiveClaw gateway listening on %s", addr)
	log.Printf("   WebSocket: ws://localhost%s/ws", addr)
	log.Printf("   Dashboard: http://localhost%s", addr)

	return http.ListenAndServe(addr, nil)
}

func (g *Gateway) handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket upgrade error: %v", err)
		return
	}

	client := &Client{
		ID:      fmt.Sprintf("client_%d", time.Now().UnixNano()),
		Conn:    conn,
		Send:    make(chan []byte, 256),
		Gateway: g,
		Role:    "operator",
	}

	g.hub.Register <- client

	go client.writePump()
	go client.readPump()

	// Send welcome message
	welcome := WSMessage{
		Type:  TypeEvent,
		Event: "connected",
		Payload: json.RawMessage(fmt.Sprintf(`{"clientId":"%s","version":"0.1.0"}`, client.ID)),
	}
	data, _ := json.Marshal(welcome)
	client.Send <- data
}

func (c *Client) readPump() {
	defer func() {
		c.Gateway.hub.Unregister <- c
		c.Conn.Close()
	}()

	for {
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket error: %v", err)
			}
			break
		}

		var msg WSMessage
		if err := json.Unmarshal(message, &msg); err != nil {
			log.Printf("Invalid message: %v", err)
			continue
		}

		c.handleMessage(msg)
	}
}

func (c *Client) writePump() {
	defer c.Conn.Close()

	for message := range c.Send {
		if err := c.Conn.WriteMessage(websocket.TextMessage, message); err != nil {
			return
		}
	}
}

func (c *Client) handleMessage(msg WSMessage) {
	switch msg.Method {
	case "connect":
		c.handleConnect(msg)
	case "chat.send":
		c.handleChatSend(msg)
	case "session.list":
		c.handleSessionList(msg)
	case "session.create":
		c.handleSessionCreate(msg)
	default:
		c.sendError(msg.ID, "UNKNOWN_METHOD", fmt.Sprintf("Unknown method: %s", msg.Method))
	}
}

func (c *Client) handleConnect(msg WSMessage) {
	ok := true
	response := WSMessage{
		Type: TypeResponse,
		ID:   msg.ID,
		OK:   &ok,
		Payload: json.RawMessage(`{"type":"hello-ok","protocol":1}`),
	}
	data, _ := json.Marshal(response)
	c.Send <- data
}

func (c *Client) handleChatSend(msg WSMessage) {
	var params struct {
		SessionID string `json:"sessionId"`
		Message   string `json:"message"`
	}
	if err := json.Unmarshal(msg.Params, &params); err != nil {
		c.sendError(msg.ID, "INVALID_PARAMS", "Invalid parameters")
		return
	}

	// TODO: Route to LLM and stream response
	log.Printf("Chat message: %s -> %s", params.SessionID, params.Message)

	// For now, echo back
	ok := true
	response := WSMessage{
		Type: TypeResponse,
		ID:   msg.ID,
		OK:   &ok,
		Payload: json.RawMessage(fmt.Sprintf(`{"response":"Echo: %s"}`, params.Message)),
	}
	data, _ := json.Marshal(response)
	c.Send <- data
}

func (c *Client) handleSessionList(msg WSMessage) {
	sessions := c.Gateway.Sessions.List()
	data, _ := json.Marshal(sessions)
	
	ok := true
	response := WSMessage{
		Type:    TypeResponse,
		ID:      msg.ID,
		OK:      &ok,
		Payload: data,
	}
	respData, _ := json.Marshal(response)
	c.Send <- respData
}

func (c *Client) handleSessionCreate(msg WSMessage) {
	var params struct {
		Name string `json:"name"`
	}
	json.Unmarshal(msg.Params, &params)

	sess := c.Gateway.Sessions.Create(params.Name)
	c.SessionID = sess.ID

	data, _ := json.Marshal(sess)
	ok := true
	response := WSMessage{
		Type:    TypeResponse,
		ID:      msg.ID,
		OK:      &ok,
		Payload: data,
	}
	respData, _ := json.Marshal(response)
	c.Send <- respData
}

func (c *Client) sendError(id, code, message string) {
	ok := false
	response := WSMessage{
		Type:  TypeResponse,
		ID:    id,
		OK:    &ok,
		Error: &WSError{Code: code, Message: message},
	}
	data, _ := json.Marshal(response)
	c.Send <- data
}

// REST handlers
func (g *Gateway) handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  "ok",
		"version": "0.1.0",
		"uptime":  time.Now().Unix(),
	})
}

func (g *Gateway) handleSessions(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(g.Sessions.List())
}

func (g *Gateway) handleChat(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		SessionID string `json:"sessionId"`
		Message   string `json:"message"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// TODO: Route to LLM
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"response": fmt.Sprintf("Echo: %s", req.Message),
	})
}
