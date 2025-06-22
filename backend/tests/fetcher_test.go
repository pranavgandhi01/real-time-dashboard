package tests

import (
	"testing"
	"real-time-dashboard/fetcher"
)

func TestFetchFlights(t *testing.T) {
	flights, err := fetcher.FetchFlights()
	if err != nil {
		t.Errorf("FetchFlights failed: %v", err)
	}
	
	if len(flights) == 0 {
		t.Error("Expected at least some flight data")
	}
	
	// Validate flight data structure
	for i, flight := range flights {
		if flight.ICAO24 == "" {
			t.Errorf("Flight %d missing ICAO24", i)
		}
		
		if flight.OriginCountry == "" {
			t.Errorf("Flight %d missing OriginCountry", i)
		}
		
		// Validate coordinate ranges
		if flight.Latitude < -90 || flight.Latitude > 90 {
			t.Errorf("Flight %d invalid latitude: %f", i, flight.Latitude)
		}
		
		if flight.Longitude < -180 || flight.Longitude > 180 {
			t.Errorf("Flight %d invalid longitude: %f", i, flight.Longitude)
		}
		
		// Validate velocity is non-negative
		if flight.Velocity < 0 {
			t.Errorf("Flight %d negative velocity: %f", i, flight.Velocity)
		}
	}
}

func TestGenerateMockFlights(t *testing.T) {
	mockFlights := fetcher.GenerateMockFlights()
	
	if len(mockFlights) == 0 {
		t.Error("Expected mock flights to be generated")
	}
	
	// Test that mock data is realistic
	for i, flight := range mockFlights {
		// Test ICAO24 format
		if len(flight.ICAO24) != 6 {
			t.Errorf("Mock flight %d invalid ICAO24 length: %s", i, flight.ICAO24)
		}
		
		// Test callsign format
		if flight.Callsign == "" {
			t.Errorf("Mock flight %d missing callsign", i)
		}
		
		// Test coordinate validity
		if flight.Latitude < -85 || flight.Latitude > 85 {
			t.Errorf("Mock flight %d invalid latitude: %f", i, flight.Latitude)
		}
		
		// Test altitude ranges
		if flight.OnGround && flight.GeoAltitude > 500 {
			t.Errorf("Mock flight %d on ground but high altitude: %f", i, flight.GeoAltitude)
		}
		
		if !flight.OnGround && flight.GeoAltitude < 1000 {
			t.Errorf("Mock flight %d in air but low altitude: %f", i, flight.GeoAltitude)
		}
	}
}

func TestFlightDataConsistency(t *testing.T) {
	// Test multiple calls return consistent data structure
	flights1, _ := fetcher.FetchFlights()
	flights2, _ := fetcher.FetchFlights()
	
	if len(flights1) != len(flights2) {
		t.Logf("Flight count varies between calls: %d vs %d (expected for mock data)", len(flights1), len(flights2))
	}
	
	// Both should have same structure
	if len(flights1) > 0 && len(flights2) > 0 {
		f1, f2 := flights1[0], flights2[0]
		
		// Check that all fields are populated
		if (f1.ICAO24 == "") != (f2.ICAO24 == "") {
			t.Error("Inconsistent ICAO24 field population")
		}
	}
}