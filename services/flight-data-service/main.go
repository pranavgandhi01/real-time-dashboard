package main

import (
	"net/http"
	"sync"
	"time"
	"github.com/gin-gonic/gin"
	"../../../pkg/config"
	"../../../pkg/types"
	"../../../pkg/client"
	"../../../pkg/health"
	"../../../pkg/log"
)

type FlightService struct {
	flights map[string]types.Flight
	mu      sync.RWMutex
	fetcher *client.FlightFetcher
	config  *config.Config
}

func NewFlightService(cfg *config.Config) *FlightService {
	fs := &FlightService{
		flights: make(map[string]types.Flight),
		fetcher: client.NewFlightFetcher(),
		config:  cfg,
	}
	go fs.startFetching()
	return fs
}

func (fs *FlightService) GetAllFlights(c *gin.Context) {
	fs.mu.RLock()
	defer fs.mu.RUnlock()
	
	flights := make([]types.Flight, 0, len(fs.flights))
	for _, flight := range fs.flights {
		flights = append(flights, flight)
	}
	c.JSON(200, flights)
}

func (fs *FlightService) GetStats(c *gin.Context) {
	fs.mu.RLock()
	defer fs.mu.RUnlock()
	
	totalFlights := len(fs.flights)
	inAir := 0
	onGround := 0
	
	for _, flight := range fs.flights {
		if flight.OnGround {
			onGround++
		} else {
			inAir++
		}
	}
	
	stats := types.FlightStats{
		TotalFlights: totalFlights,
		InAir:        inAir,
		OnGround:     onGround,
		LastUpdated:  time.Now(),
	}
	c.JSON(200, stats)
}

func (fs *FlightService) startFetching() {
	ticker := time.NewTicker(fs.config.FetchInterval)
	defer ticker.Stop()
	
	for range ticker.C {
		flights, err := fs.fetcher.FetchFlights()
		if err != nil {
			log.LogError("Failed to fetch flights: %v", err)
			continue
		}
		
		fs.mu.Lock()
		for _, flight := range flights {
			fs.flights[flight.ICAO24] = flight
		}
		fs.mu.Unlock()
		
		log.LogInfo("Updated %d flights", len(flights))
	}
}

func main() {
	cfg := config.Load()
	flightService := NewFlightService(cfg)
	r := gin.Default()
	
	r.GET("/health", gin.WrapF(health.HealthHandler))

	r.GET("/flights", flightService.GetAllFlights)
	r.GET("/stats", flightService.GetStats)

	log.LogInfo("Flight Data Service starting on port %s", cfg.Port)
	log.LogFatal("Server failed: %v", http.ListenAndServe(":"+cfg.Port, r))
}