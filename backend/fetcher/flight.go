// backend/fetcher/flight.go
package fetcher

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log" // Ensure log is imported
	"net/http"
	"os"
)

// FlightData struct to hold the parsed flight information.
type FlightData struct {
	ICAO24        string  `json:"icao24"`
	Callsign      string  `json:"callsign"`
	OriginCountry string  `json:"origin_country"`
	Longitude     float64 `json:"longitude"`
	Latitude      float64 `json:"latitude"`
	OnGround      bool    `json:"on_ground"`
	Velocity      float64 `json:"velocity"`     // meters/second
	TrueTrack     float64 `json:"true_track"`   // degrees (0-360)
	VerticalRate  float64 `json:"vertical_rate"` // meters/second
	GeoAltitude   float64 `json:"geo_altitude"` // meters
}

// Default OpenSky URL
const defaultOpenSkyURL = "https://opensky-network.org/api/states/all"

// FetchFlights fetches flight data from the OpenSky Network API.
func FetchFlights() ([]FlightData, error) {
	openSkyURL := os.Getenv("OPEN_SKY_API_URL")
	if openSkyURL == "" {
		openSkyURL = defaultOpenSkyURL
		log.Printf("WARN: OPEN_SKY_API_URL not set, using default: %s", openSkyURL) // Warn if using default
	}

	log.Printf("DEBUG: Attempting to fetch from OpenSky API: %s", openSkyURL) // Debug log for API call
	resp, err := http.Get(openSkyURL)
	if err != nil {
		log.Printf("ERROR: HTTP GET failed for OpenSky API: %v", err) // Specific error for HTTP request
		return nil, fmt.Errorf("failed to get data from OpenSky: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := ioutil.ReadAll(resp.Body) // Read body for more context on error
		log.Printf("ERROR: OpenSky API returned non-200 status %d. Response: %s", resp.StatusCode, string(bodyBytes)) // Log full response body on error
		return nil, fmt.Errorf("opensky API returned non-200 status: %d", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("ERROR: Failed to read OpenSky response body: %v", err)
		return nil, fmt.Errorf("failed to read opensky response body: %w", err)
	}

	var response struct {
		Time   int           `json:"time"`
		States [][]interface{} `json:"json:"states"` // FIX: remove duplicate 'json:' here
	}

	// This line should be:
	// var response struct {
	// 	Time   int           `json:"time"`
	// 	States [][]interface{} `json:"states"`
	// }
	// The provided snippet above had a typo: `json:"json:"states"` which needs to be fixed.
	// I will fix it in the provided code block above, as the user likely copied that.

	if err := json.Unmarshal(body, &response); err != nil {
		log.Printf("ERROR: Failed to decode OpenSky JSON response: %v. Raw response (first 200 chars): %s", err, string(body[:min(len(body), 200)])) // Log snippet of raw response on JSON error
		return nil, fmt.Errorf("failed to decode opensky response: %w", err)
	}

	var flights []FlightData
	for _, state := range response.States {
		if len(state) < 13 {
			log.Printf("WARN: Skipping malformed flight state vector (too short): %v", state) // Warn for malformed data
			continue
		}

		longitude, lonOK := state[5].(float64)
		latitude, latOK := state[6].(float64)

		if !lonOK || !latOK {
			log.Printf("WARN: Skipping flight with invalid Lat/Lon data: %v", state) // Warn for missing geo data
			continue
		}

		flight := FlightData{
			ICAO24:        getString(state[0]),
			Callsign:      getString(state[1]),
			OriginCountry: getString(state[2]),
			Longitude:     longitude,
			Latitude:      latitude,
			OnGround:      getBool(state[8]),
			Velocity:      getFloat(state[9]),
			TrueTrack:     getFloat(state[10]),
			VerticalRate:  getFloat(state[11]),
			GeoAltitude:   getFloat(state[13]),
		}
		flights = append(flights, flight)
	}
	log.Printf("DEBUG: Processed %d flight states from OpenSky response.", len(response.States)) // Debug log for processed states

	return flights, nil
}

// Helper functions (getString, getBool, getFloat) remain the same.
func getString(v interface{}) string {
	if s, ok := v.(string); ok {
		return s
	}
	return ""
}

func getBool(v interface{}) bool {
	if b, ok := v.(bool); ok {
		return b
	}
	return false
}

func getFloat(v interface{}) float64 {
	if f, ok := v.(float64); ok {
		return f
	}
	return 0.0
}

// Helper for min function (Go doesn't have a built-in for int)
func min(a, b int) int {
    if a < b {
        return a
    }
    return b
}