// backend/fetcher/flight.go
package fetcher

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	flightlog "real-time-dashboard/log"
	"time"
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

func sanitizeError(err error) string {
    if err == nil {
        return ""
    }
    // Remove newlines and control characters
    msg := strings.ReplaceAll(err.Error(), "\n", " ")
    msg = strings.ReplaceAll(msg, "\r", " ")
    // Remove other control characters
    return strings.Map(func(r rune) rune {
        if r < 32 && r != 9 { // Keep tab but remove other control chars
            return -1
        }
        return r
    }, msg)
}

// FetchFlights fetches flight data from the OpenSky Network API with retries,
// or generates mock data if USE_MOCK_DATA environment variable is "true".
func FetchFlights() ([]FlightData, error) {
	// The `GenerateMockFlights` is now in backend/fetcher/mock.go,
	// but since both files are in the same `fetcher` package,
	// `GenerateMockFlights` is directly accessible without needing to import
	// `fetcher/mock` explicitly.
	if os.Getenv("USE_MOCK_DATA") == "true" {
		return GenerateMockFlights(), nil
	}

	openSkyURL := os.Getenv("OPEN_SKY_API_URL")
	if openSkyURL == "" {
		openSkyURL = defaultOpenSkyURL
		flightlog.LogWarn("OPEN_SKY_API_URL not set, using default: %s", openSkyURL)
	}

	const maxRetries = 3
	const retryDelay = 2 * time.Second

	for i := 0; i < maxRetries; i++ {
		flightlog.LogDebug("Attempting to fetch from OpenSky API (%d/%d): %s", i+1, maxRetries, openSkyURL)
		resp, err := http.Get(openSkyURL)
		if err != nil {
			flightlog.LogError("HTTP GET failed for OpenSky API (attempt %d/%d): %v", i+1, maxRetries, err)
			if i < maxRetries-1 {
				time.Sleep(retryDelay)
				continue
			}
			return nil, fmt.Errorf("failed to get data from OpenSky after %d retries: %w", maxRetries, err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			bodyBytes, _ := ioutil.ReadAll(resp.Body)
			flightlog.LogError("OpenSky API returned non-200 status %d (attempt %d/%d). Response: %s", resp.StatusCode, i+1, maxRetries, string(bodyBytes))
			if i < maxRetries-1 {
				time.Sleep(retryDelay)
				continue
			}
			return nil, fmt.Errorf("opensky API returned non-200 status %d after %d retries", resp.StatusCode, maxRetries)
		}

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			flightlog.LogError("Failed to read OpenSky response body (attempt %d/%d): %v", i+1, maxRetries, err)
			if i < maxRetries-1 {
				time.Sleep(retryDelay)
				continue
			}
			return nil, fmt.Errorf("failed to read opensky response body after %d retries: %v", maxRetries, sanitizeError(err))
		}

		var response struct {
			Time   int           `json:"time"`
			States [][]interface{} `json:"states"`
		}

		if err := json.Unmarshal(body, &response); err != nil {
			flightlog.LogError("Failed to decode OpenSky JSON response (attempt %d/%d): %v. Raw response (first 200 chars): %s", i+1, maxRetries, err, string(body[:min(len(body), 200)]))
			if i < maxRetries-1 {
				time.Sleep(retryDelay)
				continue
			}
			return nil, fmt.Errorf("failed to decode opensky response after %d retries: %v", maxRetries, sanitizeError(err))
		}

		var flights []FlightData
		for _, state := range response.States {
			if len(state) < 13 {
				flightlog.LogWarn("Skipping malformed flight state vector (too short): %v", state)
				continue
			}

			longitude, lonOK := state[5].(float64)
			latitude, latOK := state[6].(float64)

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
		return flights, nil // Successfully fetched data, exit loop
	}
	return nil, fmt.Errorf("unexpected error: should have returned or errored out before this point") // Should not be reached
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