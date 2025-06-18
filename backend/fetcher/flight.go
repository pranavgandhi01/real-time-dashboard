// backend/fetcher/flight.go
package fetcher

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os" // Import os for environment variables
	flightlog "real-time-dashboard/log" // Import the new log package
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
		flightlog.LogWarn("OPEN_SKY_API_URL not set, using default: %s", openSkyURL)
	}

	flightlog.LogDebug("Attempting to fetch from OpenSky API: %s", openSkyURL)
	resp, err := http.Get(openSkyURL)
	if err != nil {
		flightlog.LogError("HTTP GET failed for OpenSky API: %v", err)
		return nil, fmt.Errorf("failed to get data from OpenSky: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := ioutil.ReadAll(resp.Body)
		flightlog.LogError("OpenSky API returned non-200 status %d. Response: %s", resp.StatusCode, string(bodyBytes))
		return nil, fmt.Errorf("opensky API returned non-200 status: %d", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		flightlog.LogError("Failed to read OpenSky response body: %v", err)
		return nil, fmt.Errorf("failed to read opensky response body: %w", err)
	}

	var response struct {
		Time   int           `json:"time"`
		States [][]interface{} `json:"states"`
	}

	if err := json.Unmarshal(body, &response); err != nil {
		// Log only the first 200 characters of the response body for brevity and to avoid excessively long logs
		flightlog.LogError("Failed to decode OpenSky JSON response: %v. Raw response (first 200 chars): %s", err, string(body[:min(len(body), 200)]))
		return nil, fmt.Errorf("failed to decode opensky response: %w", err)
	}

	var flights []FlightData
	for _, state := range response.States {
		// Basic validation to ensure we have enough fields.
		if len(state) < 13 {
			flightlog.LogWarn("Skipping malformed flight state vector (too short): %v", state)
			continue
		}

		// Type assertions with checks to prevent panics.
		longitude, lonOK := state[5].(float64)
		latitude, latOK := state[6].(float64)

		// We only care about flights that have coordinate data.
		if !lonOK || !latOK {
			flightlog.LogWarn("Skipping flight with invalid Lat/Lon data: %v", state)
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
	flightlog.LogDebug("Processed %d flight states from OpenSky response.", len(response.States))

	return flights, nil
}

// Helper functions to safely parse interface{} types from the state vector.
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

func getBool(v interface{}) bool {
	if b, ok := v.(bool); ok {
		return b
	}
	return false
}

// min returns the smaller of two integers.
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}