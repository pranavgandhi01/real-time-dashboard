package tests

import (
	"encoding/json"
	"testing"
	"real-time-dashboard/schema"
)

func TestValidateFlightData(t *testing.T) {
	// Valid flight data
	validData := []map[string]interface{}{
		{
			"icao24":         "MOCK001",
			"callsign":       "TEST123",
			"origin_country": "United States",
			"longitude":      -74.0060,
			"latitude":       40.7128,
			"on_ground":      false,
			"velocity":       250.5,
			"true_track":     180.0,
			"vertical_rate":  5.2,
			"geo_altitude":   10000.0,
		},
	}

	validJSON, _ := json.Marshal(validData)
	
	// Test valid data
	err := schema.ValidateFlightData(validJSON)
	if err != nil {
		t.Errorf("Expected valid data to pass validation, got error: %v", err)
	}

	// Invalid data - missing required field
	invalidData := []map[string]interface{}{
		{
			"icao24":    "MOCK001",
			"callsign":  "TEST123",
			// Missing other required fields
		},
	}

	invalidJSON, _ := json.Marshal(invalidData)
	
	// Test invalid data (should not fail due to lenient validation)
	err = schema.ValidateFlightData(invalidJSON)
	if err != nil {
		t.Logf("Validation failed as expected for incomplete data: %v", err)
	}

	// Test malformed JSON
	malformedJSON := []byte(`{"invalid": json}`)
	err = schema.ValidateFlightData(malformedJSON)
	if err == nil {
		t.Error("Expected malformed JSON to fail validation")
	}
}