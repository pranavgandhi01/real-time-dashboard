// backend/main.go
package main

import (
	"log"
	"net/http"
	"os"   // Import the os package
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

	// Get server port from environment variable, or default to 8080
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Default port
	}

	serverAddr := ":" + port
	log.Printf("HTTP and WebSocket server started on %s", serverAddr)
	// Start the server.
	err := http.ListenAndServe(serverAddr, nil)
	if err != nil {
		log.Fatalf("ListenAndServe: %v", err)
	}
}