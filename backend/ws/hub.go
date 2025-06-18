// backend/ws/hub.go
package ws

import (
	"context" // Import context for Kafka reader operations
	// "encoding/json" // REMOVED: No longer directly used for marshal/unmarshal
	"log"
	"net/http"
	"os" // Import os for environment variables
	flightlog "real-time-dashboard/log"
	"time" // Import time for consumer group rebalance

	"github.com/gorilla/websocket"
	"github.com/segmentio/kafka-go" // Import the Kafka library
)

const (
	defaultKafkaBroker = "localhost:9092" // Default Kafka broker address
	defaultKafkaTopic  = "flights"       // Default Kafka topic name
	defaultKafkaGroupID = "flight-websocket-group" // Default Kafka consumer group ID
	pongWait = 60 * time.Second
	maxMessageSize = 512
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

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client

	// Kafka Reader for consuming flight data
	kafkaReader *kafka.Reader
}

// NewHub creates a new Hub instance.
func NewHub() *Hub {
	kafkaBroker := os.Getenv("KAFKA_BROKER_ADDRESS")
	if kafkaBroker == "" {
		kafkaBroker = defaultKafkaBroker
		flightlog.LogWarn("KAFKA_BROKER_ADDRESS not set, using default for hub: %s", kafkaBroker)
	}

	kafkaTopic := os.Getenv("KAFKA_TOPIC")
	if kafkaTopic == "" {
		kafkaTopic = defaultKafkaTopic
		flightlog.LogWarn("KAFKA_TOPIC not set, using default for hub: %s", kafkaTopic)
	}

	kafkaGroupID := os.Getenv("KAFKA_GROUP_ID")
	if kafkaGroupID == "" {
		kafkaGroupID = defaultKafkaGroupID
		flightlog.LogWarn("KAFKA_GROUP_ID not set, using default for hub: %s", kafkaGroupID)
	}

	// Create a new Kafka consumer (Reader)
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{kafkaBroker},
		Topic:   kafkaTopic,
		GroupID: kafkaGroupID, // Consumer group for distributed consumption
		MinBytes: 10e3, // 10KB
		MaxBytes: 10e6, // 10MB
		MaxWait:    1 * time.Second, // Maximum amount of time to wait for new data to come when fetching messages from kafka.
		StartOffset: kafka.LastOffset, // Start consuming from the latest message
		// If no messages are available after MaxWait, it will return an empty list or an error, enabling retry.
	})

	return &Hub{
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
		kafkaReader: r, // Assign the Kafka reader to the hub
	}
}

// Run starts the hub's operations, including Kafka message consumption.
func (h *Hub) Run() {
	// Goroutine to consume messages from Kafka
	go h.readKafkaMessages()

	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
			flightlog.LogDebug("Client registered. Total clients: %d", len(h.clients))
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
				flightlog.LogDebug("Client unregistered. Total clients: %d", len(h.clients))
			}
		// The broadcast channel is no longer directly used for fetching data from main.
		// It's effectively replaced by the Kafka consumer.
		}
	}
}

// readKafkaMessages consumes messages from the Kafka topic and broadcasts them to clients.
func (h *Hub) readKafkaMessages() {
	defer func() {
		if err := h.kafkaReader.Close(); err != nil {
			flightlog.LogError("failed to close kafka reader: %v", err)
		}
	}()

	flightlog.LogInfo("Starting Kafka consumer for topic '%s' with group '%s' on broker '%s'",
		h.kafkaReader.Config().Topic, h.kafkaReader.Config().GroupID, h.kafkaReader.Config().Brokers[0])

	for {
		m, err := h.kafkaReader.ReadMessage(context.Background())
		if err != nil {
			flightlog.LogError("Error reading message from Kafka: %v", err)
			// Depending on the error, you might want to exit, retry, or log and continue.
			// For now, just log and continue to try reading next message.
			time.Sleep(1 * time.Second) // Small delay before retrying to avoid tight loop on persistent errors
			continue
		}

		flightlog.LogDebug("Received message from Kafka partition %d offset %d: %s", m.Partition, m.Offset, string(m.Value))

		// The message value is the JSON-marshaled flight data.
		// Send the message directly to all connected clients.
		for client := range h.clients {
			select {
			case client.send <- m.Value: // Send the raw byte slice received from Kafka
			default:
				close(client.send)
				delete(h.clients, client)
				flightlog.LogWarn("Client send buffer full, client disconnected.")
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

// readPump pumps messages from the websocket connection to the hub.
// The application reads messages from the websocket connection and dispatches them to the hub.
func (c *Client) readPump() {
    defer func() {
        c.hub.unregister <- c
        c.conn.Close()
    }()
    // Set read limits and handlers for keeping the connection alive
    // These constants should be defined at the top of your hub.go or in a related package.
    // I'm providing common values here. You might need to add these constants:
    // const (
    //  pongWait     = 60 * time.Second
    //  maxMessageSize = 512
    // )
    c.conn.SetReadLimit(maxMessageSize)
    c.conn.SetReadDeadline(time.Now().Add(pongWait))
    c.conn.SetPongHandler(func(string) error {
        c.conn.SetReadDeadline(time.Now().Add(pongWait))
        return nil
    })

    for {
        // This application doesn't expect messages from the client,
        // but we still need to read from the connection to detect
        // client disconnections or other protocol messages.
        _, _, err := c.conn.ReadMessage()
        if err != nil {
            if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
                // Log unexpected close errors
                log.Printf("error: %v", err)
            }
            break // Exit the loop on error (client disconnected)
        }
        // If you were to receive messages from the client, you would process them here.
    }
}

// HandleConnections handles websocket requests from the peer.
func HandleConnections(hub *Hub, w http.ResponseWriter, r *http.Request) {
    conn, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        log.Printf("WebSocket upgrade failed: %v", err)
        return
    }
    client := &Client{hub: hub, conn: conn, send: make(chan []byte, 256)}
    client.hub.register <- client

    go client.writePump() // Start goroutine to send messages to client
    client.readPump()     // This function runs in the current goroutine
}