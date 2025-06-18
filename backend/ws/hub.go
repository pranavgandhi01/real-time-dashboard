// backend/ws/hub.go
package ws

import (
	"encoding/json"
	"log"
	"net/http"
	"os"     // Import the os package
	"strings" // Import strings for splitting
	"real-time-dashboard/fetcher"

	"github.com/gorilla/websocket"
)

// The Upgrader is used to upgrade an HTTP connection to a WebSocket connection.
var upgrader = websocket.Upgrader{
	// CheckOrigin allows connections from specific origins for security.
	// In a production environment, you should restrict this to your frontend's domain.
	CheckOrigin: func(r *http.Request) bool {
		allowedOriginsStr := os.Getenv("ALLOWED_ORIGINS")
		if allowedOriginsStr == "" {
			// If no specific origins are set, allow all for development (not recommended for production)
			return true
		}

		origin := r.Header.Get("Origin")
		if origin == "" {
			return false // No origin header, deny
		}

		allowedOrigins := strings.Split(allowedOriginsStr, ",")
		for _, o := range allowedOrigins {
			if origin == strings.TrimSpace(o) {
				return true
			}
		}
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
	broadcast chan []fetcher.FlightData

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client
}

// NewHub creates a new Hub instance.
func NewHub() *Hub {
	return &Hub{
		broadcast:  make(chan []fetcher.FlightData),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
	}
}

// Run starts the hub's event loop. It handles client registration,
// unregistration, and message broadcasting.
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
			log.Println("Client registered. Total clients:", len(h.clients))
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
				log.Println("Client unregistered. Total clients:", len(h.clients))
			}
		case flightData := <-h.broadcast:
			// Marshal the flight data to JSON.
			message, err := json.Marshal(flightData)
			if err != nil {
				log.Printf("Error marshalling flight data: %v", err)
				continue
			}
			// Send the message to all connected clients.
			for client := range h.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
		}
	}
}

// writePump pumps messages from the hub to the websocket connection.
func (c *Client) writePump() {
	defer func() {
		c.conn.Close()
	}()
	for {
		message, ok := <-c.send
		if !ok {
			// The hub closed the channel.
			c.conn.WriteMessage(websocket.CloseMessage, []byte{})\
			return
		}

		w, err := c.conn.NextWriter(websocket.TextMessage)
		if err != nil {
			return
		}
		w.Write(message)

		if err := w.Close(); err != nil {
			return
		}
	}
}

// HandleConnections handles websocket requests from the peer.
func HandleConnections(hub *Hub, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	client := &Client{hub: hub, conn: conn, send: make(chan []byte, 256)}
	client.hub.register <- client

	// Allow collection of old clients to go garbage collected when unregister
	defer func() {
		client.hub.unregister <- client
	}()

	go client.writePump()
	// No need for a readPump if clients don't send messages back.
	// If clients were sending messages, a readPump would be needed here.
	// For this app, it's a broadcast-only system.
	select {} // Block forever to keep the goroutine alive until unregister
}