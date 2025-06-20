// backend/main.go
package main

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"real-time-dashboard/cache"
	"real-time-dashboard/fetcher"
	"real-time-dashboard/health"
	flightlog "real-time-dashboard/log"
	"real-time-dashboard/schema"
	"real-time-dashboard/ws"
	"github.com/prometheus/client_golang/prometheus"
    "github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/segmentio/kafka-go"
)

var (
    connectedClients = prometheus.NewGauge(prometheus.GaugeOpts{
        Name: "websocket_connected_clients",
        Help: "Number of active WebSocket clients",
    })
    fetchLatency = prometheus.NewHistogram(prometheus.HistogramOpts{
        Name: "flight_fetch_latency_seconds",
        Help: "Latency of fetching flight data",
        Buckets: prometheus.DefBuckets,
    })
)

func init() {
    prometheus.MustRegister(connectedClients, fetchLatency)
}

const (
	defaultKafkaBroker = "localhost:9092" // Default Kafka broker address
	defaultKafkaTopic  = "flights"       // Default Kafka topic name
)

func main() {
    // Initialize Schema Registry
    if err := schema.InitSchemaRegistry(); err != nil {
        flightlog.LogWarn("Schema Registry initialization failed: %v", err)
    }
    
    // Initialize Redis
    if err := cache.InitRedis(); err != nil {
        flightlog.LogWarn("Redis initialization failed: %v", err)
    }

    hub := ws.NewHub()
    go hub.Run()

    // Note: Client count tracking removed due to unexported field

    // Kafka producer setup
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
    w := &kafka.Writer{
        Addr:     kafka.TCP(kafkaBroker),
        Topic:    kafkaTopic,
        Balancer: &kafka.LeastBytes{},
    }
    defer w.Close()

    // Fetch flights with latency tracking
    ticker := time.NewTicker(15 * time.Second)
    go func() {
        for t := range ticker.C {
            flightlog.LogDebug("Fetching flight data at %v", t)
            start := time.Now()
            flights, err := fetcher.FetchFlights()
            fetchLatency.Observe(time.Since(start).Seconds())
            if err != nil {
                flightlog.LogError("Error fetching flight data: %v", err)
                continue
            }
            message, err := json.Marshal(flights)
            if err != nil {
                flightlog.LogError("Error marshalling flight data for Kafka: %v", err)
                continue
            }
            
            // Validate against schema
            if err := schema.ValidateFlightData(message); err != nil {
                flightlog.LogError("Schema validation failed: %v", err)
                continue
            }
            err = w.WriteMessages(context.Background(), kafka.Message{Value: message})
            if err != nil {
                flightlog.LogError("Failed to write message to Kafka: %v", err)
            } else {
                flightlog.LogDebug("Successfully published %d flights to Kafka topic '%s'", len(flights), kafkaTopic)
            }
        }
    }()

    // Setup HTTP routes
    http.Handle("/metrics", promhttp.Handler())
    http.HandleFunc("/health", health.HealthHandler)
    http.HandleFunc("/ready", health.ReadinessHandler)
    http.HandleFunc("/docs", func(w http.ResponseWriter, r *http.Request) {
        http.ServeFile(w, r, "docs/swagger-ui.html")
    })
    http.HandleFunc("/api-docs", func(w http.ResponseWriter, r *http.Request) {
        http.ServeFile(w, r, "docs/api-swagger.yaml")
    })
    http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
        ws.HandleConnections(hub, w, r)
    })

    port := os.Getenv("PORT")
    if port == "" {
        port = "8080"
    }
    
    server := &http.Server{
        Addr: ":" + port,
        TLSConfig: &tls.Config{
            MinVersion: tls.VersionTLS12,
        },
    }
    
    // Setup graceful shutdown
    go func() {
        sigChan := make(chan os.Signal, 1)
        signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
        <-sigChan
        
        flightlog.LogInfo("Shutting down server gracefully...")
        
        ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
        defer cancel()
        
        if err := server.Shutdown(ctx); err != nil {
            flightlog.LogError("Server shutdown error: %v", err)
        }
        
        // Close Kafka writer
        if err := w.Close(); err != nil {
            flightlog.LogError("Kafka writer close error: %v", err)
        }
        
        // Close Redis
        if err := cache.Close(); err != nil {
            flightlog.LogError("Redis close error: %v", err)
        }
        
        flightlog.LogInfo("Server shutdown complete")
    }()
    
    certPath := os.Getenv("TLS_CERT_PATH")
    keyPath := os.Getenv("TLS_KEY_PATH")
    
    if certPath != "" && keyPath != "" {
        flightlog.LogInfo("Starting HTTPS server on port %s", port)
        err := server.ListenAndServeTLS(certPath, keyPath)
        if err != nil && err != http.ErrServerClosed {
            flightlog.LogFatal("Failed to start HTTPS server: %v", err)
        }
    } else {
        flightlog.LogInfo("Starting HTTP server on port %s (TLS disabled)", port)
        err := server.ListenAndServe()
        if err != nil && err != http.ErrServerClosed {
            flightlog.LogFatal("Failed to start HTTP server: %v", err)
        }
    }
}