// backend/main.go
package main

import (
	"net/http"
	"os"
	"time"
	"real-time-dashboard/fetcher"
	"real-time-dashboard/ws"
	flightlog "real-time-dashboard/log" // Import the new log package
)

func main() {
	// The log.SetFlags and currentLogLevel initialization is now in flightlog.init() function.
	// No need to explicitly call log.SetFlags here.

	// Create a new hub for managing WebSocket connections.
	hub := ws.NewHub()
	// Run the hub in a separate goroutine.
	go hub.Run()

	// Ticker to fetch data every 15 seconds.
	ticker := time.NewTicker(15 * time.Second)
	go func() {
		for t := range ticker.C {
			flightlog.LogInfo("Starting flight data fetch cycle at %v", t) // Use flightlog.LogInfo
			// Fetch the flight data.
			flights, err := fetcher.FetchFlights()
			if err != nil {
				flightlog.LogError("Failed to fetch flight data: %v", err) // Use flightlog.LogError
				continue
			}
			flightlog.LogInfo("Successfully fetched %d flights. Broadcasting to clients.", len(flights)) // Use flightlog.LogInfo
			// Broadcast the fetched data to all connected clients.
			hub.Broadcast <- flights
		}
	}()

	// Configure the HTTP server to handle WebSocket connections.
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		flightlog.LogInfo("Incoming WebSocket connection request from %s", r.RemoteAddr) // Use flightlog.LogInfo
		ws.HandleConnections(hub, w, r)
	})

	// Get server port from environment variable, or default to 8080
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Default port
	}

	serverAddr := ":" + port
	flightlog.LogInfo("HTTP and WebSocket server starting on %s", serverAddr) // Use flightlog.LogInfo
	// Start the server.
	err := http.ListenAndServe(serverAddr, nil)
	if err != nil {
		flightlog.LogFatal("Server failed to start: %v", err) // Use flightlog.LogFatal
	}
}