package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"time"
	"github.com/gin-gonic/gin"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"mock-data-service/pkg/config"
	"mock-data-service/pkg/log"
	"mock-data-service/pkg/health"
)

type MockDataService struct {
	producer *kafka.Producer
	config   *config.Config
	baseFlights [][]interface{}
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

	// Load base flight data
	data, err := ioutil.ReadFile("../../docs/openskyapi_response.json")
	if err != nil {
		log.LogFatal("Failed to read OpenSky data: %v", err)
	}

	var openSkyData OpenSkyResponse
	if err := json.Unmarshal(data, &openSkyData); err != nil {
		log.LogFatal("Failed to parse OpenSky data: %v", err)
	}

	mds := &MockDataService{
		producer: producer,
		config:   cfg,
		baseFlights: openSkyData.States,
	}
	
	go mds.startMockBroadcast()
	return mds
}

func (mds *MockDataService) startMockBroadcast() {
	ticker := time.NewTicker(mds.config.FetchInterval)
	defer ticker.Stop()

	for range ticker.C {
		// Generate large dataset: 500-2000 flights per batch
		flightCount := rand.Intn(1500) + 500
		baseCount := len(mds.baseFlights)
		
		for i := 0; i < flightCount; i++ {
			// Use base flight data with variations
			baseIndex := i % baseCount
			baseFlight := mds.baseFlights[baseIndex]
			
			// Create variations of base flights
			avroFlight := mds.createFlightVariation(baseFlight, i)

			data, _ := json.Marshal(avroFlight)
			mds.producer.Produce(&kafka.Message{
				TopicPartition: kafka.TopicPartition{
					Topic:     &mds.config.KafkaTopic,
					Partition: kafka.PartitionAny,
				},
				Key:   []byte(avroFlight.ICAO24),
				Value: data,
			}, nil)
		}
		
		log.LogInfo("Published %d mock flights to Kafka topic %s", flightCount, mds.config.KafkaTopic)
	}
}

func (mds *MockDataService) createFlightVariation(baseFlight []interface{}, variation int) FlightAvro {
	// Extract base data with safe type assertions
	icao24 := getString(baseFlight[0])
	callsign := getString(baseFlight[1])
	baseLat := getFloat(baseFlight[6])
	baseLon := getFloat(baseFlight[5])
	baseAlt := getFloat(baseFlight[7])
	baseSpeed := getFloat(baseFlight[9])
	baseHeading := getFloat(baseFlight[10])
	baseVertical := getFloat(baseFlight[11])
	
	// Create variations
	variationFactor := float64(variation + 1)
	latVariation := (rand.Float64()-0.5) * 0.1 * variationFactor
	lonVariation := (rand.Float64()-0.5) * 0.1 * variationFactor
	altVariation := (rand.Float64()-0.5) * 1000
	speedVariation := (rand.Float64()-0.5) * 50
	
	return FlightAvro{
		Timestamp: time.Now().Unix(),
		ICAO24:    fmt.Sprintf("%s%03d", icao24, variation%1000),
		Callsign:  fmt.Sprintf("%s%03d", callsign[:3], (variation%9000)+1000),
		Position: Position{
			Latitude:  baseLat + latVariation,
			Longitude: baseLon + lonVariation,
			Altitude:  baseAlt + altVariation,
		},
		Velocity: Velocity{
			Speed:        (baseSpeed * 3.6) + speedVariation, // m/s to km/h
			Heading:      baseHeading + (rand.Float64()-0.5)*30,
			VerticalRate: baseVertical + (rand.Float64()-0.5)*10,
		},
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
		"interval": mds.config.FetchInterval.String(),
		"base_flights": len(mds.baseFlights),
		"flights_per_batch": "500-2000",
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