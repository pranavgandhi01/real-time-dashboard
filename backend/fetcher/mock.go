// backend/fetcher/mock.go
package fetcher

import (
	"fmt"
	"math/rand"
	flightlog "real-time-dashboard/log" // Assuming your custom log package is here
	"time"
)

func init() {
	// Seed the random number generator once when the package is initialized.
	rand.Seed(time.Now().UnixNano())
}

// GenerateMockFlights creates and returns a slice of mock flight data.
// This function is now responsible for producing mock data.
func GenerateMockFlights() []FlightData {
	flightlog.LogInfo("Generating mock flight data...")
	mockFlights := make([]FlightData, 0, 5) // Generate 5 mock flights

	countries := []string{"United States", "Germany", "France", "United Kingdom", "Canada"}
	callsigns := []string{"MOCK001", "TEST123", "SIMFLT", "DEVJET", "PLANEABC"}

	for i := 0; i < cap(mockFlights); i++ {
		mockFlights = append(mockFlights, FlightData{
			ICAO24:        fmt.Sprintf("MOCK%04d", i),
			Callsign:      callsigns[rand.Intn(len(callsigns))],
			OriginCountry: countries[rand.Intn(len(countries))],
			Longitude:     -180 + rand.Float64()*360,   // Random longitude
			Latitude:      -90 + rand.Float64()*180,    // Random latitude
			OnGround:      rand.Intn(2) == 0,           // Random true/false
			Velocity:      rand.Float64()*900 + 100,    // 100-1000 m/s (approx 360-3600 km/h)
			TrueTrack:     rand.Float64()*360,          // 0-360 degrees
			VerticalRate:  rand.Float64()*20 - 10,      // -10 to 10 m/s
			GeoAltitude:   rand.Float64()*12000 + 100,  // 100-12100 meters
		})
	}
	flightlog.LogInfo("Generated %d mock flights.", len(mockFlights))
	return mockFlights
}