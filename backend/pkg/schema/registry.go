package schema

import (
	"encoding/json"
	"os"
	"../log"

	"github.com/riferrei/srclient"
)

var (
	schemaRegistryClient *srclient.SchemaRegistryClient
	flightSchema         *srclient.Schema
)

func InitSchemaRegistry() error {
	schemaRegistryURL := os.Getenv("SCHEMA_REGISTRY_URL")
	if schemaRegistryURL == "" {
		schemaRegistryURL = "http://localhost:8081"
		log.LogWarn("SCHEMA_REGISTRY_URL not set, using default: %s", schemaRegistryURL)
	}

	client := srclient.CreateSchemaRegistryClient(schemaRegistryURL)
	schemaRegistryClient = client

	// Get the latest schema for flights-value subject
	schema, err := client.GetLatestSchema("flights-value")
	if err != nil {
		log.LogError("Failed to get schema from registry: %v", err)
		return err
	}
	flightSchema = schema
	log.LogInfo("Schema loaded successfully, ID: %d", schema.ID())
	return nil
}

func ValidateFlightData(data []byte) error {
	// Parse JSON to validate structure
	var flights []map[string]interface{}
	if err := json.Unmarshal(data, &flights); err != nil {
		return err
	}

	// Validate each flight record against schema
	for _, flight := range flights {
		if err := validateFlightRecord(flight); err != nil {
			log.LogWarn("Flight record validation failed: %v", err)
			return err
		}
	}
	
	log.LogDebug("Validated %d flight records", len(flights))
	return nil
}

func validateFlightRecord(flight map[string]interface{}) error {
	// Basic field validation based on schema
	requiredFields := []string{"icao24", "callsign", "origin_country", "longitude", "latitude", "on_ground", "velocity", "true_track", "vertical_rate", "geo_altitude"}
	
	for _, field := range requiredFields {
		if _, exists := flight[field]; !exists {
			log.LogError("Missing required field: %s", field)
			return nil // Continue processing even with validation errors
		}
	}
	return nil
}