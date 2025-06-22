package client

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
	"../types"
)

type FlightFetcher struct {
	client  *http.Client
	baseURL string
}

func NewFlightFetcher() *FlightFetcher {
	return &FlightFetcher{
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
		baseURL: "https://opensky-network.org/api/states/all",
	}
}

func (f *FlightFetcher) FetchFlights() ([]types.Flight, error) {
	resp, err := f.client.Get(f.baseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch flights: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status %d", resp.StatusCode)
	}

	var response struct {
		States [][]interface{} `json:"states"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	flights := make([]types.Flight, 0, len(response.States))
	for _, state := range response.States {
		if len(state) < 10 {
			continue
		}

		flight := types.Flight{
			ICAO24:        getString(state[0]),
			Callsign:      getString(state[1]),
			OriginCountry: getString(state[2]),
			Longitude:     getFloat64(state[5]),
			Latitude:      getFloat64(state[6]),
			OnGround:      getBool(state[8]),
			Velocity:      getFloat64(state[9]),
			LastUpdated:   time.Now(),
		}

		flights = append(flights, flight)
	}

	return flights, nil
}

func getString(v interface{}) string {
	if s, ok := v.(string); ok {
		return s
	}
	return ""
}

func getFloat64(v interface{}) float64 {
	if f, ok := v.(float64); ok {
		return f
	}
	return 0
}

func getBool(v interface{}) bool {
	if b, ok := v.(bool); ok {
		return b
	}
	return false
}