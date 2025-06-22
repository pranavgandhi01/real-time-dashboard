package types

import (
	"testing"
	"time"
)

func TestFlightStruct(t *testing.T) {
	now := time.Now()
	flight := Flight{
		ICAO24:        "abc123",
		Callsign:      "TEST123",
		OriginCountry: "United States",
		Longitude:     -122.4194,
		Latitude:      37.7749,
		OnGround:      false,
		Velocity:      250.5,
		LastUpdated:   now,
	}
	
	if flight.ICAO24 != "abc123" {
		t.Errorf("Expected ICAO24 abc123, got %s", flight.ICAO24)
	}
	
	if flight.OnGround != false {
		t.Errorf("Expected OnGround false, got %t", flight.OnGround)
	}
	
	if flight.Velocity != 250.5 {
		t.Errorf("Expected Velocity 250.5, got %f", flight.Velocity)
	}
}

func TestFlightStats(t *testing.T) {
	now := time.Now()
	stats := FlightStats{
		TotalFlights: 100,
		InAir:        75,
		OnGround:     25,
		LastUpdated:  now,
	}
	
	if stats.TotalFlights != 100 {
		t.Errorf("Expected TotalFlights 100, got %d", stats.TotalFlights)
	}
	
	if stats.InAir+stats.OnGround != stats.TotalFlights {
		t.Errorf("InAir + OnGround should equal TotalFlights")
	}
}