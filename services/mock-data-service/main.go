package main

import (
	"encoding/json"
	"io/ioutil"
	"math/rand"
	"net/http"
	"time"
	"github.com/gin-gonic/gin"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"../../../pkg/config"
	"../../../pkg/log"
	"../../../pkg/health"
)

type MockDataService struct {
	producer *kafka.Producer
	config   *config.Config
	mockData OpenSkyResponse
}

type OpenSkyResponse struct {
	Time   int64           `json:"time"`
	States [][]interface{} `json:"states"`
}

type FlightAvro struct {
	Timestamp int64    `json:"timestamp"`
	ICAO24    string   `json:"icao24"`
	Callsign  string   `json:"callsign"`
	Position  Position `json:"position"`
	Velocity  Velocity `json:"velocity"`
}

type Position struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Altitude  float64 `json:"altitude"`
}

type Velocity struct {
	Speed        float64 `json:"speed"`
	Heading      float64 `json:"heading"`
	VerticalRate float64 `json:"verticalRate"`
}

func NewMockDataService(cfg *config.Config) *MockDataService {
	producer, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": cfg.KafkaBroker,
	})
	if err != nil {
		log.LogFatal("Failed to create producer: %v", err)
	}

	// Load mock data
	data, err := ioutil.ReadFile("../../docs/openskyapi_response.json")
	if err != nil {
		log.LogFatal("Failed to read mock data: %v", err)
	}

	var mockData OpenSkyResponse
	if err := json.Unmarshal(data, &mockData); err != nil {
		log.LogFatal("Failed to parse mock data: %v", err)
	}

	mds := &MockDataService{
		producer: producer,
		config:   cfg,
		mockData: mockData,
	}
	
	go mds.startMockBroadcast()
	return mds
}

func (mds *MockDataService) startMockBroadcast() {
	ticker := time.NewTicker(mds.config.FetchInterval)
	defer ticker.Stop()

	for range ticker.C {
		for _, state := range mds.mockData.States {
			if len(state) < 13 {
				continue
			}

			// Add random variations to simulate movement
			lat := getFloat(state[6]) + (rand.Float64()-0.5)*0.01
			lon := getFloat(state[5]) + (rand.Float64()-0.5)*0.01
			alt := getFloat(state[7]) + (rand.Float64()-0.5)*100

			avroFlight := FlightAvro{
				Timestamp: time.Now().Unix(),
				ICAO24:    getString(state[0]),
				Callsign:  getString(state[1]),
				Position: Position{
					Latitude:  lat,
					Longitude: lon,
					Altitude:  alt,
				},
				Velocity: Velocity{
					Speed:        getFloat(state[9]) * 3.6, // m/s to km/h
					Heading:      getFloat(state[10]),
					VerticalRate: getFloat(state[11]),
				},
			}

			data, _ := json.Marshal(avroFlight)
			// Use ICAO24 as partition key for consistent routing
			mds.producer.Produce(&kafka.Message{
				TopicPartition: kafka.TopicPartition{
					Topic:     &mds.config.KafkaTopic,
					Partition: kafka.PartitionAny,
				},
				Key:   []byte(avroFlight.ICAO24),
				Value: data,
			}, nil)
		}
		
		log.LogInfo("Published %d mock flights to Kafka", len(mds.mockData.States))
	}
}

func getString(v interface{}) string {
	if s, ok := v.(string); ok {
		return s
	}
	return ""
}

func getFloat(v interface{}) float64 {
	if f, ok := v.(float64); ok {
		return f
	}
	return 0.0
}

func (mds *MockDataService) GetStats(c *gin.Context) {
	c.JSON(200, gin.H{
		"service": "mock-data",
		"status":  "active",
		"topic":   mds.config.KafkaTopic,
		"flights": len(mds.mockData.States),
	})
}

func main() {
	cfg := config.Load()
	mockService := NewMockDataService(cfg)
	defer mockService.producer.Close()

	r := gin.Default()
	r.GET("/health", gin.WrapF(health.HealthHandler))
	r.GET("/stats", mockService.GetStats)

	log.LogInfo("Mock Data Service starting on port %s", cfg.Port)
	log.LogFatal("Server failed: %v", http.ListenAndServe(":"+cfg.Port, r))
}