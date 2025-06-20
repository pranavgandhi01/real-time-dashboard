import { useMemo } from 'react';
import { MapContainer, TileLayer, Marker, Popup } from 'react-leaflet';
import L from 'leaflet';
import 'leaflet/dist/leaflet.css';
import performanceConfig from '../config/performance';

// Fix Leaflet default icon issue
delete (L.Icon.Default.prototype as any)._getIconUrl;
L.Icon.Default.mergeOptions({
  iconRetinaUrl: '/marker-icon-2x.png',
  iconUrl: '/marker-icon.png',
  shadowUrl: '/marker-shadow.png',
});

// Custom flight icon
const createFlightIcon = (onGround: boolean, heading: number) => {
  const color = onGround ? '#f59e0b' : '#3b82f6';
  const size = onGround ? 12 : 16;
  
  return L.divIcon({
    html: `<div style="width:${size}px;height:${size}px;background:${color};border:2px solid white;border-radius:50%;box-shadow:0 2px 4px rgba(0,0,0,0.3);transform:rotate(${heading}deg)"></div>`,
    className: 'flight-icon',
    iconSize: [size + 4, size + 4],
    iconAnchor: [(size + 4) / 2, (size + 4) / 2],
  });
};

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

export default function FlightMap({ flights }: FlightMapProps) {
  const validFlights = useMemo(() => 
    flights.filter(f => 
      f.latitude >= -90 && f.latitude <= 90 && 
      f.longitude >= -180 && f.longitude <= 180
    ), [flights]
  );

  // Lazy loading for large datasets
  const displayFlights = useMemo(() => {
    const maxDisplay = performanceConfig.maxDisplayFlights;
    if (validFlights.length <= maxDisplay) {
      return validFlights;
    }
    console.log(`[FlightMap] Lazy loading: showing ${maxDisplay}/${validFlights.length} flights`);
    return validFlights.slice(0, maxDisplay);
  }, [validFlights]);

  if (validFlights.length === 0) {
    return (
      <div className="h-96 w-full bg-gray-800 border border-gray-700 flex items-center justify-center">
        <div className="text-center">
          <div className="text-4xl mb-2">🗺️</div>
          <div className="text-gray-400">Waiting for flight data...</div>
        </div>
      </div>
    );
  }

  return (
    <div style={{ height: '400px', width: '100%' }}>
      <MapContainer
        center={[40, 0]}
        zoom={2}
        style={{ height: '100%', width: '100%' }}
      >
        <TileLayer url="https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png" />
        {displayFlights.map((flight) => (
          <Marker 
            key={flight.icao24} 
            position={[flight.latitude, flight.longitude]}
            icon={createFlightIcon(flight.on_ground, flight.true_track)}
          >
            <Popup>
              <strong>{flight.callsign}</strong><br/>
              {flight.origin_country}<br/>
              {flight.on_ground ? 'On Ground' : 'In Air'}
            </Popup>
          </Marker>
        ))}
      </MapContainer>
    </div>
  );
}