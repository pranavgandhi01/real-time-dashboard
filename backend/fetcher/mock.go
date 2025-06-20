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
	mockFlights := make([]FlightData, 0, 15) // Generate 15 mock flights for better filtering

	// More realistic data based on actual flight patterns
	countries := []string{"United States", "Germany", "France", "United Kingdom", "Canada", "Japan", "Australia", "Brazil", "India", "China"}
	airlines := []string{"AAL", "UAL", "DAL", "LUV", "BAW", "DLH", "AFR", "KLM", "SWA", "JBU", "ANA", "JAL", "QFA", "TAM"}
	
	// Country-based coordinates
	countryCoords := map[string]struct{ lat, lon, variation float64 }{
		"United States": {39.0, -98.0, 15.0},
		"Germany":       {51.0, 9.0, 3.0},
		"France":        {46.0, 2.0, 4.0},
		"United Kingdom": {54.0, -2.0, 4.0},
		"Canada":        {60.0, -95.0, 20.0},
		"Japan":         {36.0, 138.0, 5.0},
		"Australia":     {-25.0, 133.0, 15.0},
		"Brazil":        {-14.0, -51.0, 12.0},
		"India":         {20.0, 77.0, 8.0},
		"China":         {35.0, 104.0, 12.0},
	}

	for i := 0; i < cap(mockFlights); i++ {
		// Choose random country
		country := countries[rand.Intn(len(countries))]
		coords := countryCoords[country]
		
		// Generate coordinates within country bounds
		lat := coords.lat + (rand.Float64()-0.5)*coords.variation
		lon := coords.lon + (rand.Float64()-0.5)*coords.variation
		
		// Ensure valid ranges
		if lat > 85 { lat = 85 }
		if lat < -85 { lat = -85 }
		if lon > 180 { lon = 180 }
		if lon < -180 { lon = -180 }
		
		// More realistic flight parameters
		onGround := rand.Float64() < 0.2 // 20% on ground
		var velocity, altitude, verticalRate float64
		
		if onGround {
			velocity = rand.Float64()*30 + 5      // 5-35 m/s for ground
			altitude = rand.Float64()*100 + 50    // 50-150m ground level
			verticalRate = 0                      // No vertical movement on ground
		} else {
			velocity = rand.Float64()*150 + 100   // 100-250 m/s (360-900 km/h) for air
			altitude = rand.Float64()*10000 + 2000 // 2000-12000m cruising altitude
			verticalRate = (rand.Float64()-0.5)*10 // ±5 m/s vertical rate
		}

		mockFlights = append(mockFlights, FlightData{
			ICAO24:        fmt.Sprintf("%06X", rand.Intn(16777216)),
			Callsign:      fmt.Sprintf("%s%d", airlines[rand.Intn(len(airlines))], rand.Intn(9999)+1),
			OriginCountry: country,
			Longitude:     lon,
			Latitude:      lat,
			OnGround:      onGround,
			Velocity:      velocity,
			TrueTrack:     rand.Float64()*360,
			VerticalRate:  verticalRate,
			GeoAltitude:   altitude,
		})
	}
	flightlog.LogInfo("Generated %d mock flights.", len(mockFlights))
	return mockFlights
}