// backend/main.go
package main

import (
	"log"
	"net/http"
	"os"
	"time"
	"real-time-dashboard/fetcher"
	"real-time-dashboard/ws"
)

func main() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile) // Add file and line number to logs

	// Create a new hub for managing WebSocket connections.
	hub := ws.NewHub()
	// Run the hub in a separate goroutine.
	go hub.Run()

	// Ticker to fetch data every 15 seconds.
	ticker := time.NewTicker(15 * time.Second)
	go func() {
		for t := range ticker.C {
			log.Printf("INFO: Starting flight data fetch cycle at %v", t) // More descriptive log
			// Fetch the flight data.
			flights, err := fetcher.FetchFlights()
			if err != nil {
				log.Printf("ERROR: Failed to fetch flight data: %v", err) // More descriptive error log
				continue
			}
			log.Printf("INFO: Successfully fetched %d flights. Broadcasting to clients.", len(flights)) // Success log
			// Broadcast the fetched data to all connected clients.
			hub.Broadcast <- flights
		}
	}()

	// Configure the HTTP server to handle WebSocket connections.
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("INFO: Incoming WebSocket connection request from %s", r.RemoteAddr) // Log new connection attempts
		ws.HandleConnections(hub, w, r)
	})

	// Get server port from environment variable, or default to 8080
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Default port
	}

	serverAddr := ":" + port
	log.Printf("INFO: HTTP and WebSocket server starting on %s", serverAddr) // Informative startup log
	// Start the server.
	err := http.ListenAndServe(serverAddr, nil)
	if err != nil {
		log.Fatalf("FATAL: Server failed to start: %v", err) // Use Fatal for unrecoverable errors
	}
}