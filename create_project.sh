#!/bin/bash

# Create directory structure
mkdir -p real-time-dashboard/{backend/{fetcher,ws},frontend/{pages,components,lib}}

# backend/main.go
cat > real-time-dashboard/backend/main.go << 'EOF'
package main

import (
	"log"
	"net/http"
	"time"
	"real-time-dashboard/fetcher"
	"real-time-dashboard/ws"
)

func main() {
	// Create a new hub for managing WebSocket connections.
	hub := ws.NewHub()
	// Run the hub in a separate goroutine.
	go hub.Run()

	// Ticker to fetch data every 15 seconds.
	// OpenSky asks not to poll more than every 10 seconds.
	ticker := time.NewTicker(15 * time.Second)
	go func() {
		for t := range ticker.C {
			log.Println("Fetching flight data at", t)
			// Fetch the flight data.
			flights, err := fetcher.FetchFlights()
			if err != nil {
				log.Printf("Error fetching flight data: %v", err)
				continue
			}
			// Broadcast the fetched data to all connected clients.
			hub.Broadcast <- flights
		}
	}()

	// Configure the HTTP server to handle WebSocket connections.
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		ws.HandleConnections(hub, w, r)
	})

	log.Println("HTTP and WebSocket server started on :8080")
	// Start the server.
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
EOF

# backend/ws/hub.go
cat > real-time-dashboard/backend/ws/hub.go << 'EOF'
package ws

import (
	"encoding/json"
	"log"
	"net/http"
	"real-time-dashboard/fetcher"

	"github.com/gorilla/websocket"
)

// The Upgrader is used to upgrade an HTTP connection to a WebSocket connection.
var upgrader = websocket.Upgrader{
	// CheckOrigin allows connections from any origin.
	// In a production environment, you should restrict this to your frontend's domain.
	CheckOrigin: func(r *http.Request) bool {
		return true
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

// Run starts the hub's event loop.
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
			c.conn.WriteMessage(websocket.CloseMessage, []byte{})
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

	// Allow collection of memory referenced by the caller by doing all work in
	// new goroutines.
	go client.writePump()
}
EOF

# backend/fetcher/flight.go
cat > real-time-dashboard/backend/fetcher/flight.go << 'EOF'
package fetcher

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// FlightData struct to hold the parsed flight information.
type FlightData struct {
	ICAO24        string  `json:"icao24"`
	Callsign      string  `json:"callsign"`
	OriginCountry string  `json:"origin_country"`
	Longitude     float64 `json:"longitude"`
	Latitude      float64 `json:"latitude"`
	OnGround      bool    `json:"on_ground"`
	Velocity      float64 `json:"velocity"`     // meters/second
	TrueTrack     float64 `json:"true_track"`   // degrees (0-360)
	VerticalRate  float64 `json:"vertical_rate"` // meters/second
	GeoAltitude   float64 `json:"geo_altitude"` // meters
}

const openSkyURL = "https://opensky-network.org/api/states/all"

// FetchFlights fetches flight data from the OpenSky Network API.
func FetchFlights() ([]FlightData, error) {
	resp, err := http.Get(openSkyURL)
	if err != nil {
		return nil, fmt.Errorf("failed to get data from OpenSky: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("opensky API returned non-200 status: %d", resp.StatusCode)
	}
	
	body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return nil, fmt.Errorf("failed to read response body: %w", err)
    }

	// The API returns a JSON object with a "states" field,
	// which is an array of state vectors.
	var response struct {
		Time   int           `json:"time"`
		States [][]interface{} `json:"states"`
	}

	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to decode opensky response: %w", err)
	}

	var flights []FlightData
	for _, state := range response.States {
		// Basic validation to ensure we have enough fields.
		if len(state) < 13 {
			continue
		}

		// Type assertions with checks to prevent panics.
		longitude, lonOK := state[5].(float64)
		latitude, latOK := state[6].(float64)
		
		// We only care about flights that have coordinate data.
		if !lonOK || !latOK {
			continue
		}

		flight := FlightData{
			ICAO24:        getString(state[0]),
			Callsign:      getString(state[1]), // Can be null, trim spaces
			OriginCountry: getString(state[2]),
			Longitude:     longitude,
			Latitude:      latitude,
			OnGround:      getBool(state[8]),
			Velocity:      getFloat(state[9]),
			TrueTrack:     getFloat(state[10]),
			VerticalRate:  getFloat(state[11]),
			GeoAltitude:   getFloat(state[13]),
		}
		flights = append(flights, flight)
	}

	return flights, nil
}

// Helper functions to safely parse interface{} types from the state vector.
func getString(v interface{}) string {
	if s, ok := v.(string); ok {
		return s
	}
	return ""
}

func getFloat(v interface{}) float64 {
	if f, ok := v.(float64); ok {
		return f
	}
	return 0.0
}

func getBool(v interface{}) bool {
	if b, ok := v.(bool); ok {
		return b
	}
	return false
}
EOF

# backend/go.mod
cat > real-time-dashboard/backend/go.mod << 'EOF'
module real-time-dashboard

go 1.18

require github.com/gorilla/websocket v1.5.0
EOF

# backend/Dockerfile
cat > real-time-dashboard/backend/Dockerfile << 'EOF'
# Stage 1: Build the Go application
FROM golang:1.18-alpine AS builder

WORKDIR /app

# Copy go.mod and go.sum files to download dependencies
COPY go.mod ./
RUN go mod download

# Copy the source code
COPY . .

# Build the application
# CGO_ENABLED=0 is important for creating a static binary for alpine
# -o /app/server creates the binary named 'server' in the /app directory
RUN CGO_ENABLED=0 go build -o /app/server .

# Stage 2: Create a small production image
FROM alpine:latest

WORKDIR /root/

# Copy the pre-built binary from the builder stage
COPY --from=builder /app/server .

# Expose port 8080 to the outside world
EXPOSE 8080

# Command to run the executable
CMD ["./server"]
EOF

# Create empty frontend files (you can populate these later)
touch real-time-dashboard/frontend/pages/index.tsx
touch real-time-dashboard/frontend/components/FlightMap.tsx
touch real-time-dashboard/frontend/lib/socket.ts
touch real-time-dashboard/docker-compose.yml
touch real-time-dashboard/README.md

echo "Project structure and files created successfully!"