// backend/ws/hub.go
package ws

import (
	"encoding/json"
	"net/http"
	"os"
	"strings"
	"real-time-dashboard/fetcher"
	flightlog "real-time-dashboard/log" // Import the new log package

	"github.com/gorilla/websocket"
)

// The Upgrader is used to upgrade an HTTP connection to a WebSocket connection.
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		allowedOriginsStr := os.Getenv("ALLOWED_ORIGINS")
		if allowedOriginsStr == "" {
			flightlog.LogWarn("ALLOWED_ORIGINS environment variable not set. Allowing all WebSocket origins (NOT recommended for production).") // Use flightlog.LogWarn
			return true
		}

		origin := r.Header.Get("Origin")
		if origin == "" {
			flightlog.LogWarn("WebSocket connection denied - no Origin header provided from %s", r.RemoteAddr) // Use flightlog.LogWarn
			return false
		}

		allowedOrigins := strings.Split(allowedOriginsStr, ",")
		for _, o := range allowedOrigins {
			trimmedOrigin := strings.TrimSpace(o)
			if origin == trimmedOrigin {
				flightlog.LogInfo("WebSocket connection allowed for origin '%s' from %s", origin, r.RemoteAddr) // Use flightlog.LogInfo
				return true
			}
		}
		flightlog.LogWarn("WebSocket connection denied - origin '%s' not in allowed list from %s", origin, r.RemoteAddr) // Use flightlog.LogWarn
		return false
	},
}

// Client is a middleman between the websocket connection and the hub.
type Client struct {
	hub *Hub

	// The websocket connection.
	conn *websocket.Conn

	// Buffered channel of outbound messages.
	send chan []byte
}

// Hub maintains the set of active clients and broadcasts messages to the
// clients.
type Hub struct {
	// Registered clients.
	clients map[*Client]bool

	// Inbound messages from the clients.
	Broadcast chan []fetcher.FlightData

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client
}

// NewHub creates a new Hub instance.
func NewHub() *Hub {
	return &Hub{
		Broadcast:  make(chan []fetcher.FlightData),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
	}
}

// Run starts the hub's event loop. It handles client registration,
// unregistration, and message broadcasting.
func (h *Hub) Run() {
	flightlog.LogInfo("WebSocket Hub started.") // Use flightlog.LogInfo
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
			flightlog.LogInfo("Client registered. Total clients: %d", len(h.clients)) // Use flightlog.LogInfo
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
				flightlog.LogInfo("Client unregistered. Total clients: %d", len(h.clients)) // Use flightlog.LogInfo
			}
		case flightData := <-h.Broadcast:
			message, err := json.Marshal(flightData)
			if err != nil {
				flightlog.LogError("Failed to marshal flight data for broadcast: %v", err) // Use flightlog.LogError
				continue
			}
			flightlog.LogDebug("Broadcasting %d bytes of flight data to %d clients.", len(message), len(h.clients)) // Use flightlog.LogDebug
			for client := range h.clients {
				select {
				case client.send <- message:
					// Message sent successfully
				default:
					close(client.send)
					delete(h.clients, client)
					flightlog.LogWarn("Client send buffer full, disconnecting client. Total clients: %d", len(h.clients)) // Use flightlog.LogWarn
				}
			}
		}
	}
}

// writePump pumps messages from the hub to the websocket connection.
func (c *Client) writePump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
		flightlog.LogInfo("writePump for client %s exiting.", c.conn.RemoteAddr().String()) // Use flightlog.LogInfo
	}()
	for {
		message, ok := <-c.send
		if !ok {
			flightlog.LogInfo("Hub closed send channel for client %s.", c.conn.RemoteAddr().String()) // Use flightlog.LogInfo
			c.conn.WriteMessage(websocket.CloseMessage, []byte{})
			return
		}

		w, err := c.conn.NextWriter(websocket.TextMessage)
		if err != nil {
			flightlog.LogError("Failed to get next writer for client %s: %v", c.conn.RemoteAddr().String(), err) // Use flightlog.LogError
			return
		}
		w.Write(message)

		if err := w.Close(); err != nil {
			flightlog.LogError("Failed to close writer for client %s: %v", c.conn.RemoteAddr().String(), err) // Use flightlog.LogError
			return
		}
		flightlog.LogDebug("Sent %d bytes to client %s.", len(message), c.conn.RemoteAddr().String()) // Use flightlog.LogDebug
	}
}

// HandleConnections handles websocket requests from the peer.
func HandleConnections(hub *Hub, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		flightlog.LogError("Failed to upgrade HTTP to WebSocket for %s: %v", r.RemoteAddr, err) // Use flightlog.LogError
		return
	}
	client := &Client{hub: hub, conn: conn, send: make(chan []byte, 256)}
	client.hub.register <- client
	flightlog.LogInfo("WebSocket connection established for client %s", conn.RemoteAddr().String()) // Use flightlog.LogInfo

	go client.writePump()

	select {}
}