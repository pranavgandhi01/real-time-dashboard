// backend/ws/hub.go
package ws

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"
	"real-time-dashboard/fetcher"

	"github.com/gorilla/websocket"
)

// The Upgrader is used to upgrade an HTTP connection to a WebSocket connection.
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		allowedOriginsStr := os.Getenv("ALLOWED_ORIGINS")
		if allowedOriginsStr == "" {
			log.Println("WARN: ALLOWED_ORIGINS environment variable not set. Allowing all WebSocket origins (NOT recommended for production).")
			return true
		}

		origin := r.Header.Get("Origin")
		if origin == "" {
			log.Printf("WARN: WebSocket connection denied - no Origin header provided from %s", r.RemoteAddr)
			return false
		}

		allowedOrigins := strings.Split(allowedOriginsStr, ",")
		for _, o := range allowedOrigins {
			trimmedOrigin := strings.TrimSpace(o)
			if origin == trimmedOrigin {
				log.Printf("INFO: WebSocket connection allowed for origin '%s' from %s", origin, r.RemoteAddr)
				return true
			}
		}
		log.Printf("WARN: WebSocket connection denied - origin '%s' not in allowed list from %s", origin, r.RemoteAddr)
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
	log.Println("INFO: WebSocket Hub started.")
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
			log.Printf("INFO: Client registered. Total clients: %d", len(h.clients))
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
				log.Printf("INFO: Client unregistered. Total clients: %d", len(h.clients))
			}
		case flightData := <-h.Broadcast:
			message, err := json.Marshal(flightData)
			if err != nil {
				log.Printf("ERROR: Failed to marshal flight data for broadcast: %v", err)
				continue
			}
			log.Printf("DEBUG: Broadcasting %d bytes of flight data to %d clients.", len(message), len(h.clients))
			for client := range h.clients {
				select {
				case client.send <- message:
					// Message sent successfully
				default:
					close(client.send)
					delete(h.clients, client)
					log.Printf("WARN: Client send buffer full, disconnecting client. Total clients: %d", len(h.clients))
				}
			}
		}
	}
}

// writePump pumps messages from the hub to the websocket connection.
func (c *Client) writePump() {
	defer func() {
		c.hub.unregister <- c // Ensure client is unregistered on exit
		c.conn.Close()
		log.Printf("INFO: writePump for client %s exiting.", c.conn.RemoteAddr().String())
	}()
	for {
		message, ok := <-c.send
		if !ok {
			log.Printf("INFO: Hub closed send channel for client %s.", c.conn.RemoteAddr().String())
			c.conn.WriteMessage(websocket.CloseMessage, []byte{})
			return
		}

		w, err := c.conn.NextWriter(websocket.TextMessage)
		if err != nil {
			log.Printf("ERROR: Failed to get next writer for client %s: %v", c.conn.RemoteAddr().String(), err)
			return
		}
		w.Write(message)

		if err := w.Close(); err != nil {
			log.Printf("ERROR: Failed to close writer for client %s: %v", c.conn.RemoteAddr().String(), err)
			return
		}
		log.Printf("DEBUG: Sent %d bytes to client %s.", len(message), c.conn.RemoteAddr().String())
	}
}

// HandleConnections handles websocket requests from the peer.
func HandleConnections(hub *Hub, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("ERROR: Failed to upgrade HTTP to WebSocket for %s: %v", r.RemoteAddr, err)
		return
	}
	client := &Client{hub: hub, conn: conn, send: make(chan []byte, 256)}
	client.hub.register <- client
	log.Printf("INFO: WebSocket connection established for client %s", conn.RemoteAddr().String())

	// Start the writePump in a goroutine
	go client.writePump()

	// Keep the goroutine alive to handle potential read messages (though not used now)
	// or simply to ensure the client remains registered until the connection closes.
	// This select{} statement effectively blocks the goroutine until the connection is closed
	// or something else causes it to exit.
	select {}
}