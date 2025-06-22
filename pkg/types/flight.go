package types

import "time"

type Flight struct {
	ICAO24        string    `json:"icao24"`
	Callsign      string    `json:"callsign"`
	OriginCountry string    `json:"origin_country"`
	Longitude     float64   `json:"longitude"`
	Latitude      float64   `json:"latitude"`
	OnGround      bool      `json:"on_ground"`
	Velocity      float64   `json:"velocity"`
	TrueTrack     float64   `json:"true_track"`
	VerticalRate  float64   `json:"vertical_rate"`
	GeoAltitude   float64   `json:"geo_altitude"`
	LastUpdated   time.Time `json:"last_updated"`
}

type FlightStats struct {
	TotalFlights int       `json:"total_flights"`
	InAir        int       `json:"in_air"`
	OnGround     int       `json:"on_ground"`
	LastUpdated  time.Time `json:"last_updated"`
}