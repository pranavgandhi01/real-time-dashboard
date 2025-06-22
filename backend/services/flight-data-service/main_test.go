package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
	"github.com/gin-gonic/gin"
	"../../../pkg/config"
	"../../../pkg/types"
)

func TestFlightService_GetAllFlights(t *testing.T) {
	gin.SetMode(gin.TestMode)
	
	cfg := &config.Config{Port: "8081"}
	fs := NewFlightService(cfg)
	fs.flights["test123"] = types.Flight{
		ICAO24:        "test123",
		Callsign:      "TEST123",
		OriginCountry: "Test Country",
		Longitude:     -122.4194,
		Latitude:      37.7749,
		OnGround:      false,
		Velocity:      250.5,
		LastUpdated:   time.Now(),
	}
	
	r := gin.New()
	r.GET("/flights", fs.GetAllFlights)
	
	req, _ := http.NewRequest("GET", "/flights", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	
	if w.Code != 200 {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
	
	var flights []types.Flight
	json.Unmarshal(w.Body.Bytes(), &flights)
	
	if len(flights) != 1 {
		t.Errorf("Expected 1 flight, got %d", len(flights))
	}
	
	if flights[0].ICAO24 != "test123" {
		t.Errorf("Expected ICAO24 test123, got %s", flights[0].ICAO24)
	}
}

func TestFlightService_GetStats(t *testing.T) {
	gin.SetMode(gin.TestMode)
	
	cfg := &config.Config{Port: "8081"}
	fs := NewFlightService(cfg)
	fs.flights["air1"] = types.Flight{OnGround: false}
	fs.flights["ground1"] = types.Flight{OnGround: true}
	
	r := gin.New()
	r.GET("/stats", fs.GetStats)
	
	req, _ := http.NewRequest("GET", "/stats", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	
	if w.Code != 200 {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
	
	var stats types.FlightStats
	json.Unmarshal(w.Body.Bytes(), &stats)
	
	if stats.TotalFlights != 2 {
		t.Errorf("Expected 2 total flights, got %d", stats.TotalFlights)
	}
	
	if stats.InAir != 1 {
		t.Errorf("Expected 1 in air, got %d", stats.InAir)
	}
	
	if stats.OnGround != 1 {
		t.Errorf("Expected 1 on ground, got %d", stats.OnGround)
	}
}