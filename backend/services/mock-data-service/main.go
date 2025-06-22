package main

import (
	"context"
	"encoding/json"
	"math/rand"
	"net/http"
	"strings"
	"time"
	"github.com/gin-gonic/gin"
	"github.com/segmentio/kafka-go"
	"github.com/real-time-dashboard/backend/pkg/config"
	"github.com/real-time-dashboard/backend/pkg/health"
	"github.com/real-time-dashboard/backend/pkg/log"
	"github.com/real-time-dashboard/backend/pkg/types"
)

func generateMockFlights() []types.Flight {
	flights := make([]types.Flight, 50)
	for i := 0; i < 50; i++ {
		flights[i] = types.Flight{
			ICAO24:        generateRandomICAO(),
			Callsign:      generateRandomCallsign(),
			OriginCountry: "Mock Country",
			Longitude:     -180 + rand.Float64()*360,
			Latitude:      -90 + rand.Float64()*180,
			OnGround:      rand.Float64() < 0.3,
			Velocity:      rand.Float64() * 500,
			LastUpdated:   time.Now(),
		}
	}
	return flights
}

func generateRandomICAO() string {
	chars := "ABCDEF0123456789"
	result := make([]byte, 6)
	for i := range result {
		result[i] = chars[rand.Intn(len(chars))]
	}
	return string(result)
}

func generateRandomCallsign() string {
	airlines := []string{"UAL", "DAL", "AAL", "SWA", "JBU"}
	return airlines[rand.Intn(len(airlines))] + "123"
}

type KafkaProducer struct {
	writer *kafka.Writer
}

func NewKafkaProducer(brokers, topic string) *KafkaProducer {
	return &KafkaProducer{
		writer: &kafka.Writer{
			Addr:         kafka.TCP(strings.Split(brokers, ",")...),
			Topic:        topic,
			Balancer:     &kafka.LeastBytes{},
			RequiredAcks: kafka.RequireOne,
			Compression:  kafka.Snappy,
		},
	}
}

func (p *KafkaProducer) PublishFlights(flights []types.Flight) error {
	messages := make([]kafka.Message, 0, len(flights))
	
	for _, flight := range flights {
		data, _ := json.Marshal(flight)
		messages = append(messages, kafka.Message{
			Key:   []byte(flight.ICAO24),
			Value: data,
			Time:  time.Now(),
		})
	}
	
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	return p.writer.WriteMessages(ctx, messages...)
}

func main() {
	cfg := config.Load()
	producer := NewKafkaProducer(cfg.KafkaBroker, cfg.KafkaTopic)
	
	// Start periodic publishing
	go func() {
		ticker := time.NewTicker(10 * time.Second)
		for range ticker.C {
			flights := generateMockFlights()
			if err := producer.PublishFlights(flights); err != nil {
				log.LogError("Failed to publish to Kafka: %v", err)
			} else {
				log.LogInfo("Published %d mock flights to Kafka", len(flights))
			}
		}
	}()
	
	r := gin.Default()
	r.GET("/health", gin.WrapF(health.HealthHandler))
	r.GET("/flights", func(c *gin.Context) {
		c.JSON(200, generateMockFlights())
	})
	r.GET("/stats", func(c *gin.Context) {
		c.JSON(200, gin.H{"total": 50, "mock": true})
	})

	log.LogInfo("Mock Data Service starting on port %s", cfg.Port)
	log.LogFatal("Server failed: %v", http.ListenAndServe(":"+cfg.Port, r))
}