// backend/ws/hub.go
package ws

import (
	"bytes"
	"compress/gzip"
	"context" // Import context for Kafka reader operations
	// "encoding/json" // REMOVED: No longer directly used for marshal/unmarshal
	"net/http"
	"os" // Import os for environment variables
	flightlog "real-time-dashboard/log"
	"real-time-dashboard/schema"
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
    CheckOrigin: func(r *http.Request) bool {
        allowedOrigin := os.Getenv("ALLOWED_ORIGINS")
        if allowedOrigin == "" {
            allowedOrigin = "http://localhost:3000"
            flightlog.LogWarn("ALLOWED_ORIGINS not set, using default: %s", allowedOrigin)
        }
        origin := r.Header.Get("Origin")
        return origin == allowedOrigin
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
    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()
    defer h.kafkaReader.Close()

    flightlog.LogInfo("Starting Kafka consumer for topic '%s' with group '%s' on broker '%s'",
        h.kafkaReader.Config().Topic, h.kafkaReader.Config().GroupID, h.kafkaReader.Config().Brokers[0])

    for {
        select {
        case <-ctx.Done():
            flightlog.LogInfo("Shutting down Kafka consumer")
            return
        default:
            m, err := h.kafkaReader.ReadMessage(ctx)
            if err != nil {
                if err == context.Canceled {
                    return
                }
                flightlog.LogError("Error reading message from Kafka: %v", err)
                time.Sleep(2 * time.Second)
                continue
            }

            flightlog.LogDebug("Received message from Kafka partition %d offset %d: %s", m.Partition, m.Offset, string(m.Value))
            
            // Validate message against schema before broadcasting
            if err := schema.ValidateFlightData(m.Value); err != nil {
                flightlog.LogError("Schema validation failed for received message: %v", err)
                continue
            }

            // Compress the message before sending to clients
            flightlog.LogDebug("Original message size: %d bytes", len(m.Value))
            var buf bytes.Buffer
            gz := gzip.NewWriter(&buf)
            if _, err := gz.Write(m.Value); err != nil {
                flightlog.LogError("Failed to compress message: %v", err)
                continue
            }
            if err := gz.Close(); err != nil {
                flightlog.LogError("Failed to close gzip writer: %v", err)
                continue
            }
            compressedData := buf.Bytes()
            flightlog.LogDebug("Compressed message size: %d bytes", len(compressedData))

            for client := range h.clients {
                select {
                case client.send <- compressedData:
                default:
                    close(client.send)
                    delete(h.clients, client)
                    flightlog.LogWarn("Client send buffer full, client disconnected.")
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
            c.conn.WriteMessage(websocket.CloseMessage, []byte{})
            return
        }

        w, err := c.conn.NextWriter(websocket.BinaryMessage)
        if err != nil {
            flightlog.LogError("Failed to get WebSocket writer: %v", err)
            return
        }
        
        if _, err := w.Write(message); err != nil {
            flightlog.LogError("Failed to write message to WebSocket: %v", err)
            return
        }

        if err := w.Close(); err != nil {
            flightlog.LogError("Failed to close WebSocket writer: %v", err)
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
                flightlog.LogError("WebSocket connection closed unexpectedly: %v", err)
            } else {
                flightlog.LogDebug("WebSocket connection closed normally: %v", err)
            }
            break // Exit the loop on error (client disconnected)
        }
        // If you were to receive messages from the client, you would process them here.
    }
}

// HandleConnections handles websocket requests from the peer.
func HandleConnections(hub *Hub, w http.ResponseWriter, r *http.Request) {
    token := r.URL.Query().Get("token")
    expectedToken := os.Getenv("WEBSOCKET_TOKEN")
    if expectedToken != "" && token != expectedToken {
        http.Error(w, "Unauthorized", http.StatusUnauthorized)
        flightlog.LogWarn("Unauthorized WebSocket connection attempt")
        return
    }
    conn, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        flightlog.LogError("WebSocket upgrade failed: %v", err)
        return
    }
    client := &Client{hub: hub, conn: conn, send: make(chan []byte, 256)}
    client.hub.register <- client
    flightlog.LogInfo("WebSocket connection established from %s", r.RemoteAddr)
    go client.writePump()
    client.readPump()
}