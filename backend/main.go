// backend/main.go
package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"real-time-dashboard/fetcher"
	flightlog "real-time-dashboard/log" // Assuming you have this custom log package
	"real-time-dashboard/ws"

	"github.com/segmentio/kafka-go" // Import the Kafka library
)

const (
	defaultKafkaBroker = "localhost:9092" // Default Kafka broker address
	defaultKafkaTopic  = "flights"       // Default Kafka topic name
)

func main() {
	// Initialize WebSocket hub (still needed for handling local clients, but will consume from Kafka)
	hub := ws.NewHub()
	go hub.Run() // Start the hub's message processing loop

	// --- Kafka Producer Setup ---
	kafkaBroker := os.Getenv("KAFKA_BROKER_ADDRESS")
	if kafkaBroker == "" {
		kafkaBroker = defaultKafkaBroker
		flightlog.LogWarn("KAFKA_BROKER_ADDRESS not set, using default: %s", kafkaBroker)
	}

	kafkaTopic := os.Getenv("KAFKA_TOPIC")
	if kafkaTopic == "" {
		kafkaTopic = defaultKafkaTopic
		flightlog.LogWarn("KAFKA_TOPIC not set, using default: %s", kafkaTopic)
	}

	// Create a Kafka producer (Writer)
	w := &kafka.Writer{
		Addr:     kafka.TCP(kafkaBroker),
		Topic:    kafkaTopic,
		Balancer: &kafka.LeastBytes{}, // Choose a balancer (e.g., round-robin, least-bytes)
	}
	defer func() {
		if err := w.Close(); err != nil {
			flightlog.LogError("failed to close kafka writer: %v", err)
		}
	}()
	// --- End Kafka Producer Setup ---

	// Ticker to fetch data and publish to Kafka every 15 seconds.
	ticker := time.NewTicker(15 * time.Second)
	go func() {
		for t := range ticker.C {
			flightlog.LogDebug("Fetching flight data at %v", t)
			flights, err := fetcher.FetchFlights()
			if err != nil {
				flightlog.LogError("Error fetching flight data: %v", err)
				continue
			}

			// Marshal flights to JSON
			message, err := json.Marshal(flights)
			if err != nil {
				flightlog.LogError("Error marshalling flight data for Kafka: %v", err)
				continue
			}

			// Publish to Kafka
			err = w.WriteMessages(context.Background(),
				kafka.Message{
					Value: message,
				},
			)
			if err != nil {
				flightlog.LogError("Failed to write message to Kafka: %v", err)
			} else {
				flightlog.LogDebug("Successfully published %d flights to Kafka topic '%s'", len(flights), kafkaTopic)
			}
		}
	}()

	// Configure the HTTP server to handle WebSocket connections.
	// This part will eventually consume from Kafka within the hub.
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