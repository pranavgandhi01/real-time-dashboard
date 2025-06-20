import { useEffect } from 'react';
import { MapContainer, TileLayer, Marker, Popup, useMap } from 'react-leaflet';
import L from 'leaflet';

// Fix Leaflet default icons
delete (L.Icon.Default.prototype as any)._getIconUrl;
L.Icon.Default.mergeOptions({
  iconRetinaUrl: 'https://cdnjs.cloudflare.com/ajax/libs/leaflet/1.7.1/images/marker-icon-2x.png',
  iconUrl: 'https://cdnjs.cloudflare.com/ajax/libs/leaflet/1.7.1/images/marker-icon.png',
  shadowUrl: 'https://cdnjs.cloudflare.com/ajax/libs/leaflet/1.7.1/images/marker-shadow.png',
});

interface FlightData {
  icao24: string;
  callsign: string;
  origin_country: string;
  longitude: number;
  latitude: number;
  on_ground: boolean;
  velocity: number;
  true_track: number;
  vertical_rate: number;
  geo_altitude: number;
}

interface FlightMapProps {
  flights: FlightData[];
}

// Custom flight icon
const createFlightIcon = (onGround: boolean, heading: number) => {
  const color = onGround ? '#fbbf24' : '#3b82f6';
  return L.divIcon({
    html: `<div style="transform: rotate(${heading}deg); color: ${color}; font-size: 16px; text-align: center; line-height: 20px;">✈️</div>`,
    className: 'flight-icon',
    iconSize: [20, 20],
    iconAnchor: [10, 10],
    popupAnchor: [0, -10],
  });
};

// Component to update map bounds
function MapUpdater({ flights }: { flights: FlightData[] }) {
  const map = useMap();
  
  useEffect(() => {
    if (flights.length > 0) {
      const bounds = L.latLngBounds(
        flights.map(flight => [flight.latitude, flight.longitude])
      );
      map.fitBounds(bounds, { padding: [20, 20] });
    }
  }, [flights, map]);
  
  return null;
}

export default function FlightMap({ flights }: FlightMapProps) {
  return (
    <div className="h-96 w-full rounded-lg overflow-hidden shadow-lg">
      <MapContainer
        center={[40.7128, -74.0060]} // NYC default
        zoom={2}
        style={{ height: '100%', width: '100%' }}
        className="z-0"
      >
        <TileLayer
          attribution='&copy; <a href="https://www.openstreetmap.org/copyright">OpenStreetMap</a>'
          url="https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png"
        />
        <MapUpdater flights={flights} />
        {flights.map((flight) => (
          <Marker
            key={flight.icao24}
            position={[flight.latitude, flight.longitude]}
            icon={createFlightIcon(flight.on_ground, flight.true_track)}
          >
            <Popup>
              <div className="text-sm">
                <h3 className="font-bold text-blue-600">{flight.callsign || 'N/A'}</h3>
                <p><strong>From:</strong> {flight.origin_country}</p>
                <p><strong>Altitude:</strong> {flight.geo_altitude.toFixed(0)} m</p>
                <p><strong>Speed:</strong> {(flight.velocity * 3.6).toFixed(1)} km/h</p>
                <p><strong>Status:</strong> {flight.on_ground ? 'On Ground' : 'In Air'}</p>
              </div>
            </Popup>
          </Marker>
        ))}
      </MapContainer>
    </div>
  );
}