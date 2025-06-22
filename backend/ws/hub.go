// backend/ws/hub.go
package ws

import (
	"bytes"
	"compress/gzip"
	"context"
	"net/http"
	"os"
	"strings"
	flightlog "real-time-dashboard/log"
	"real-time-dashboard/schema"
	"real-time-dashboard/config"
	"real-time-dashboard/memory"
	"real-time-dashboard/scaling"
	"time"

	"github.com/gorilla/websocket"
	"github.com/segmentio/kafka-go"
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
        allowedOrigins := os.Getenv("ALLOWED_ORIGINS")
        if allowedOrigins == "" {
            allowedOrigins = "http://localhost:3000,http://127.0.0.1:3000"
            flightlog.LogWarn("ALLOWED_ORIGINS not set, using defaults: %s", allowedOrigins)
        }
        
        origin := r.Header.Get("Origin")
        if origin == "" {
            // Allow requests without Origin header (like curl)
            return true
        }
        
        // Check if origin is in allowed list
        for _, allowed := range strings.Split(allowedOrigins, ",") {
            if strings.TrimSpace(allowed) == origin {
                return true
            }
        }
        
        flightlog.LogWarn("WebSocket connection denied for origin: %s", origin)
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

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client

	// Kafka Reader for consuming flight data
	kafkaReader *kafka.Reader

	// Configuration
	config *config.Config

	// Connection pool metrics
	activeConnections int

	// Memory management
	memoryWindow *memory.SlidingWindow

	// Auto-scaling
	autoScaler *scaling.AutoScaler
}

// NewHub creates a new Hub instance.
func NewHub(cfg *config.Config) *Hub {
	// Create Kafka reader with retry logic
	r := createKafkaReader(cfg)

	return &Hub{
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
		kafkaReader: r,
		config: cfg,
		activeConnections: 0,
		memoryWindow: memory.NewSlidingWindow(
			cfg.Memory.WindowMinutes,
			cfg.Memory.MaxSize,
			cfg.Memory.CleanupInterval,
		),
		autoScaler: scaling.NewAutoScaler(
			cfg.WebSocket.MaxConnections,
			cfg.Scaling.ScaleUpThreshold,
			cfg.Scaling.ScaleDownThreshold,
			cfg.Scaling.CooldownMinutes,
			cfg.Scaling.MonitorInterval,
		),
	}
}

// createKafkaReader creates Kafka reader with connection retry
func createKafkaReader(cfg *config.Config) *kafka.Reader {
	for attempt := 1; attempt <= cfg.Kafka.MaxRetries; attempt++ {
		flightlog.LogInfo("Attempting Kafka connection (%d/%d) to %s", attempt, cfg.Kafka.MaxRetries, cfg.Kafka.BrokerAddress)
		
		r := kafka.NewReader(kafka.ReaderConfig{
			Brokers: []string{cfg.Kafka.BrokerAddress},
			Topic:   cfg.Kafka.Topic,
			GroupID: cfg.Kafka.GroupID,
			MinBytes: 10e3,
			MaxBytes: 10e6,
			MaxWait: 1 * time.Second,
			StartOffset: kafka.LastOffset,
		})
		
		// Test connection
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		_, err := r.ReadMessage(ctx)
		cancel()
		
		if err == nil || !isConnectionError(err) {
			flightlog.LogInfo("Kafka connection established")
			return r
		}
		
		flightlog.LogWarn("Kafka connection failed (attempt %d/%d): %v", attempt, cfg.Kafka.MaxRetries, err)
		r.Close()
		
		if attempt < cfg.Kafka.MaxRetries {
			time.Sleep(time.Duration(cfg.Kafka.RetryInterval) * time.Second)
		}
	}
	
	if cfg.Kafka.FailFast {
		flightlog.LogFatal("Failed to connect to Kafka after %d attempts, exiting", cfg.Kafka.MaxRetries)
	}
	
	flightlog.LogWarn("Kafka unavailable, continuing without Kafka integration")
	return nil
}

// isConnectionError checks if error is connection-related
func isConnectionError(err error) bool {
	errorStr := err.Error()
	return strings.Contains(errorStr, "connection refused") ||
		strings.Contains(errorStr, "no such host") ||
		strings.Contains(errorStr, "timeout")
}

// Run starts the hub's operations, including Kafka message consumption.
func (h *Hub) Run() {
	// Goroutine to consume messages from Kafka
	go h.readKafkaMessages()

	for {
		select {
		case client := <-h.register:
			if h.activeConnections >= h.config.WebSocket.MaxConnections {
				flightlog.LogWarn("Connection pool full, rejecting client")
				client.conn.Close()
				continue
			}
			h.clients[client] = true
			h.activeConnections++
			h.autoScaler.UpdateConnections(h.activeConnections)
			flightlog.LogDebug("Client registered. Active: %d/%d", h.activeConnections, h.config.WebSocket.MaxConnections)
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
				h.activeConnections--
				flightlog.LogDebug("Client unregistered. Active: %d/%d", h.activeConnections, h.config.WebSocket.MaxConnections)
			}
		// The broadcast channel is no longer directly used for fetching data from main.
		// It's effectively replaced by the Kafka consumer.
		}
	}
}

// readKafkaMessages consumes messages from the Kafka topic and broadcasts them to clients.
func (h *Hub) readKafkaMessages() {
    if h.kafkaReader == nil {
        flightlog.LogWarn("Kafka reader not available, skipping message consumption")
        return
    }
    
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

            // Store in sliding window for memory management
            h.memoryWindow.Add(m.Value)

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

            // Track message queue size
            queueSize := 0
            for client := range h.clients {
                queueSize += len(client.send)
                select {
                case client.send <- compressedData:
                default:
                    close(client.send)
                    delete(h.clients, client)
                    h.activeConnections--
                    flightlog.LogWarn("Client send buffer full, client disconnected.")
                }
            }
            
            // Update auto-scaler metrics
            h.autoScaler.UpdateQueueDepth(queueSize)
            h.autoScaler.UpdateConnections(h.activeConnections)
            
            flightlog.LogDebug("Message queue size: %d, Active clients: %d", queueSize, len(h.clients))
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
    expectedToken := hub.config.WebSocket.Token
    
    // Only validate token if one is configured
    if expectedToken != "" && token != expectedToken {
        http.Error(w, "Unauthorized", http.StatusUnauthorized)
        flightlog.LogWarn("Unauthorized WebSocket connection attempt from %s", r.RemoteAddr)
        return
    }
    
    flightlog.LogInfo("WebSocket connection authorized from %s", r.RemoteAddr)
    conn, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        flightlog.LogError("WebSocket upgrade failed: %v", err)
        return
    }
    client := &Client{hub: hub, conn: conn, send: make(chan []byte, 256)}
    client.hub.register <- client
    flightlog.LogInfo("WebSocket connection established from %s", r.RemoteAddr)
    go client.writePump()
    go client.readPump()
}