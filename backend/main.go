// backend/main.go
package main

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"real-time-dashboard/cache"
	"real-time-dashboard/config"
	"real-time-dashboard/fetcher"
	"real-time-dashboard/health"
	flightlog "real-time-dashboard/log"
	"real-time-dashboard/ratelimit"
	"real-time-dashboard/scaling"
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
    flightDataProcessingTime = prometheus.NewHistogram(prometheus.HistogramOpts{
        Name: "flight_data_processing_seconds",
        Help: "Time taken to process flight data",
        Buckets: []float64{0.001, 0.005, 0.01, 0.05, 0.1, 0.5, 1.0},
    })
    websocketMessageQueueSize = prometheus.NewGauge(prometheus.GaugeOpts{
        Name: "websocket_message_queue_size",
        Help: "Current size of WebSocket message queue",
    })
    kafkaConsumerLag = prometheus.NewGauge(prometheus.GaugeOpts{
        Name: "kafka_consumer_lag",
        Help: "Kafka consumer lag in messages",
    })
    redisCacheHits = prometheus.NewCounter(prometheus.CounterOpts{
        Name: "redis_cache_hits_total",
        Help: "Total number of Redis cache hits",
    })
    redisCacheMisses = prometheus.NewCounter(prometheus.CounterOpts{
        Name: "redis_cache_misses_total",
        Help: "Total number of Redis cache misses",
    })
)

func init() {
    prometheus.MustRegister(
        connectedClients, 
        fetchLatency, 
        flightDataProcessingTime,
        websocketMessageQueueSize,
        kafkaConsumerLag,
        redisCacheHits,
        redisCacheMisses,
    )
}

const (
	defaultKafkaBroker = "localhost:9092" // Default Kafka broker address
	defaultKafkaTopic  = "flights"       // Default Kafka topic name
)

func main() {
    // Load configuration
    cfg := config.Load()
    flightlog.LogInfo("Configuration loaded: MaxConnections=%d, Port=%s", cfg.WebSocket.MaxConnections, cfg.WebSocket.Port)
    
    // Initialize Schema Registry
    if err := schema.InitSchemaRegistry(); err != nil {
        flightlog.LogWarn("Schema Registry initialization failed: %v", err)
    }
    
    // Initialize Redis
    if err := cache.InitRedis(); err != nil {
        flightlog.LogWarn("Redis initialization failed: %v", err)
    }

    // Initialize rate limiter from config
    rateLimiter := ratelimit.NewRateLimiter(
        cfg.Scaling.RateLimitPerIP, 
        time.Duration(cfg.Scaling.RateLimitWindow)*time.Minute,
    )
    
    // Initialize scalable consumer for horizontal scaling
    scalableConsumer := scaling.NewScalableConsumer(cfg)
    defer scalableConsumer.Close()
    
    hub := ws.NewHub(cfg)
    go hub.Run()

    // Note: Client count tracking removed due to unexported field

    // Kafka producer setup with connection test
    w := createKafkaWriter(cfg)
    if w != nil {
        defer w.Close()
    }

    // Fetch flights with latency tracking
    ticker := time.NewTicker(15 * time.Second)
    go func() {
        for t := range ticker.C {
            flightlog.LogDebug("Fetching flight data at %v", t)
            fetchStart := time.Now()
            flights, err := fetcher.FetchFlights()
            fetchLatency.Observe(time.Since(fetchStart).Seconds())
            if err != nil {
                flightlog.LogError("Error fetching flight data: %v", err)
                continue
            }
            
            processStart := time.Now()
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
            
            flightDataProcessingTime.Observe(time.Since(processStart).Seconds())
            if w != nil {
                err = w.WriteMessages(context.Background(), kafka.Message{Value: message})
                if err != nil {
                    flightlog.LogError("Failed to write message to Kafka: %v", err)
                } else {
                    flightlog.LogDebug("Successfully published %d flights to Kafka topic '%s'", len(flights), cfg.Kafka.Topic)
                }
            } else {
                flightlog.LogDebug("Kafka writer unavailable, skipping message publish")
            }
        }
    }()

    // Setup HTTP routes
    http.Handle("/metrics", promhttp.Handler())
    http.HandleFunc("/health", health.HealthHandler)
    http.HandleFunc("/ready", health.ReadinessHandler)
    http.HandleFunc("/scale-health", func(w http.ResponseWriter, r *http.Request) {
        // This would be called by orchestration systems for scaling decisions
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusOK)
        w.Write([]byte(`{"status":"healthy","connections":` + fmt.Sprintf("%d", hub.activeConnections) + `}`))
    })
    http.HandleFunc("/docs", func(w http.ResponseWriter, r *http.Request) {
        http.ServeFile(w, r, "docs/swagger-ui.html")
    })
    http.HandleFunc("/api-docs", func(w http.ResponseWriter, r *http.Request) {
        http.ServeFile(w, r, "docs/api-swagger.yaml")
    })
    http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
        if !rateLimiter.Allow(r.RemoteAddr) {
            http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
            flightlog.LogWarn("Rate limit exceeded for %s", r.RemoteAddr)
            return
        }
        ws.HandleConnections(hub, w, r)
    })

    port := cfg.WebSocket.Port
    
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
        if w != nil {
            if err := w.Close(); err != nil {
                flightlog.LogError("Kafka writer close error: %v", err)
            }
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

// createKafkaWriter creates Kafka writer with connection test
func createKafkaWriter(cfg *config.Config) *kafka.Writer {
    for attempt := 1; attempt <= cfg.Kafka.MaxRetries; attempt++ {
        flightlog.LogInfo("Testing Kafka writer connection (%d/%d) to %s", attempt, cfg.Kafka.MaxRetries, cfg.Kafka.BrokerAddress)
        
        w := &kafka.Writer{
            Addr:         kafka.TCP(cfg.Kafka.BrokerAddress),
            Topic:        cfg.Kafka.Topic,
            Balancer:     &kafka.LeastBytes{},
            BatchSize:    cfg.Kafka.BatchSize,
            BatchTimeout: time.Duration(cfg.Kafka.BatchTimeout) * time.Millisecond,
            MaxAttempts:  cfg.Kafka.MaxAttempts,
            WriteTimeout: time.Duration(cfg.Kafka.WriteTimeout) * time.Second,
            ReadTimeout:  time.Duration(cfg.Kafka.ReadTimeout) * time.Second,
        }
        
        // Test connection with a small message
        ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
        err := w.WriteMessages(ctx, kafka.Message{Value: []byte("test")})
        cancel()
        
        if err == nil {
            flightlog.LogInfo("Kafka writer connection established")
            return w
        }
        
        flightlog.LogWarn("Kafka writer connection failed (attempt %d/%d): %v", attempt, cfg.Kafka.MaxRetries, err)
        w.Close()
        
        if attempt < cfg.Kafka.MaxRetries {
            time.Sleep(time.Duration(cfg.Kafka.RetryInterval) * time.Second)
        }
    }
    
    if cfg.Kafka.FailFast {
        flightlog.LogFatal("Failed to connect Kafka writer after %d attempts, exiting", cfg.Kafka.MaxRetries)
    }
    
    flightlog.LogWarn("Kafka writer unavailable, continuing without Kafka publishing")
    return nil
}