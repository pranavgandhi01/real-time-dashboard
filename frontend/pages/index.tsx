// frontend/pages/index.tsx
import { useEffect, useState, useMemo } from "react";
import dynamic from "next/dynamic";
import pako from "pako";
import FlightStats from "../components/FlightStats";
import FlightFilters from "../components/FlightFilters";

// Dynamic import to avoid SSR issues with Leaflet
const FlightMap = dynamic(() => import("../components/FlightMap"), {
  ssr: false,
  loading: () => <div className="h-96 bg-gray-800 rounded-lg flex items-center justify-center">Loading map...</div>
});

// Define the structure of the flight data we expect from the backend
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

export default function Home() {
  const [flights, setFlights] = useState<FlightData[]>([]);
  const [status, setStatus] = useState("Connecting...");
  const [error, setError] = useState<string | null>(null);
  
  // Filter states
  const [selectedCountry, setSelectedCountry] = useState("");
  const [statusFilter, setStatusFilter] = useState<'all' | 'air' | 'ground'>('all');
  const [minSpeed, setMinSpeed] = useState(0);
  const [viewMode, setViewMode] = useState<'map' | 'grid'>('grid');
  
  // Memoized filtered flights
  const filteredFlights = useMemo(() => {
    return flights.filter(flight => {
      if (selectedCountry && flight.origin_country !== selectedCountry) return false;
      if (statusFilter === 'air' && flight.on_ground) return false;
      if (statusFilter === 'ground' && !flight.on_ground) return false;
      if (flight.velocity < minSpeed) return false;
      return true;
    });
  }, [flights, selectedCountry, statusFilter, minSpeed]);
  
  // Get unique countries for filter
  const countries = useMemo(() => {
    return Array.from(new Set(flights.map(f => f.origin_country))).sort();
  }, [flights]);

  useEffect(() => {
    const websocketUrl = `${
      process.env.NEXT_PUBLIC_WEBSOCKET_URL || "ws://localhost:8080/ws"
    }?token=${process.env.NEXT_PUBLIC_WEBSOCKET_TOKEN || ""}`;
    let socket: WebSocket;
    let reconnectAttempts = 0;
    const maxReconnectAttempts = 10;
    const baseReconnectDelay = 1000;
    
    console.log('[WebSocket] Connecting to:', websocketUrl);

    function connect() {
      socket = new WebSocket(websocketUrl);
      socket.binaryType = 'arraybuffer'; // Ensure binary data is received as ArrayBuffer
      setStatus("Connecting...");

      socket.onopen = () => {
        setStatus("Connected");
        setError(null);
        reconnectAttempts = 0; // Reset attempts on successful connection
        if (process.env.NODE_ENV === 'development') {
          console.info('[WebSocket] Connection established');
        }
      };

      socket.onmessage = (event) => {
        try {
          // Handle binary compressed data
          if (process.env.NODE_ENV === 'development') {
            console.log('[WebSocket] Received data type:', typeof event.data, 'Size:', event.data.byteLength || event.data.length);
          }
          const compressedData = new Uint8Array(event.data);
          const decompressed = pako.ungzip(compressedData, { to: "string" });
          const data: FlightData[] = JSON.parse(decompressed);
          if (data && Array.isArray(data)) {
            setFlights(data);
            setError(null);
          } else {
            setError("Received empty or invalid flight data from server.");
          }
        } catch (error) {
          if (error instanceof SyntaxError) {
            setError(`Failed to parse flight data: Invalid JSON format`);
          } else if (error instanceof Error) {
            setError(`Error processing flight data: ${error.message}`);
          } else {
            setError(`Unknown error occurred while processing flight data`);
          }
          if (process.env.NODE_ENV === 'development') {
            console.error('[WebSocket] Message processing error:', error);
          }
        }
      };

      socket.onclose = (event) => {
        console.log('[WebSocket] Connection closed:', event.code, event.reason);
        setStatus("Disconnected");
        
        if (event.code !== 1000 && reconnectAttempts < maxReconnectAttempts) {
          // Exponential backoff with jitter
          const exponentialDelay = Math.min(baseReconnectDelay * Math.pow(2, reconnectAttempts), 30000);
          const jitter = Math.random() * 1000;
          const delay = exponentialDelay + jitter;
          
          setError(`Reconnecting in ${Math.round(delay/1000)}s... (${reconnectAttempts + 1}/${maxReconnectAttempts})`);
          
          setTimeout(() => {
            reconnectAttempts++;
            connect();
          }, delay);
        } else {
          setError('WebSocket connection failed. Please refresh the page.');
        }
      };

      socket.onerror = (event) => {
        console.error('[WebSocket] Connection error:', event);
        setStatus("Error");
        setError('WebSocket connection error. Check if backend is running.');
      };
    }

    connect();

    return () => {
      socket.close();
    };
  }, []);

  return (
    <div className="min-h-screen bg-gray-900 text-white p-6">
      <header className="mb-8">
        <div className="flex flex-col lg:flex-row justify-between items-start lg:items-end gap-4">
          <div>
            <h1 className="text-4xl font-extrabold text-teal-400 mb-1">
              Live Flight Tracker
            </h1>
            <p className="text-gray-400">Real-time flight monitoring dashboard</p>
          </div>
          <div className="flex items-center gap-3">
            <div className="flex bg-gray-800 rounded-lg p-1">
              <button
                onClick={() => setViewMode('map')}
                className={`px-3 py-2 rounded text-sm font-medium transition-all ${
                  viewMode === 'map' ? 'bg-blue-600 text-white' : 'text-gray-400 hover:text-white'
                }`}
              >
                🗺️ Map View
              </button>
              <button
                onClick={() => setViewMode('grid')}
                className={`px-3 py-2 rounded text-sm font-medium transition-all ${
                  viewMode === 'grid' ? 'bg-blue-600 text-white' : 'text-gray-400 hover:text-white'
                }`}
              >
                📋 Grid View
              </button>
            </div>
            <div className="flex items-center gap-2 bg-gray-800 px-3 py-2 rounded-lg">
              <div className={`w-2 h-2 rounded-full ${
                status === "Connected" ? "bg-green-500" : "bg-red-500"
              }`}></div>
              <span className="text-sm text-gray-300">Status:</span>
              <span className={`text-sm font-medium ${
                status === "Connected" ? "text-green-400" : "text-red-400"
              }`}>
                {status}
              </span>
            </div>
          </div>
        </div>
      </header>

      {error && (
        <div className="bg-red-700 p-4 rounded-lg shadow-md mb-6 text-center text-lg">
          Error: {error}
        </div>
      )}
      
      <FlightStats flights={filteredFlights} />
      
      <FlightFilters
        countries={countries}
        selectedCountry={selectedCountry}
        onCountryChange={setSelectedCountry}
        statusFilter={statusFilter}
        onStatusChange={setStatusFilter}
        minSpeed={minSpeed}
        onMinSpeedChange={setMinSpeed}
      />

      {viewMode === 'map' ? (
        <FlightMap flights={filteredFlights} />
      ) : (
        <div className="bg-gray-800 rounded-lg overflow-hidden">
          {filteredFlights.length > 0 ? (
            <table className="w-full">
              <thead className="bg-gray-700">
                <tr>
                  <th className="px-4 py-3 text-left text-sm font-semibold text-gray-300">Flight</th>
                  <th className="px-4 py-3 text-left text-sm font-semibold text-gray-300">Status</th>
                  <th className="px-4 py-3 text-left text-sm font-semibold text-gray-300">Country</th>
                  <th className="px-4 py-3 text-left text-sm font-semibold text-gray-300">Altitude</th>
                  <th className="px-4 py-3 text-left text-sm font-semibold text-gray-300">Speed</th>
                  <th className="px-4 py-3 text-left text-sm font-semibold text-gray-300">Heading</th>
                  <th className="px-4 py-3 text-left text-sm font-semibold text-gray-300">Position</th>
                </tr>
              </thead>
              <tbody>
                {filteredFlights.map((flight, index) => (
                  <tr key={flight.icao24} className={`border-t border-gray-700 hover:bg-gray-750 ${index % 2 === 0 ? 'bg-gray-800' : 'bg-gray-850'}`}>
                    <td className="px-4 py-3">
                      <span className="font-bold text-teal-300">{flight.callsign || "N/A"}</span>
                    </td>
                    <td className="px-4 py-3">
                      <span className={`px-2 py-1 text-xs rounded-full ${
                        flight.on_ground ? "bg-yellow-600 text-yellow-100" : "bg-blue-600 text-blue-100"
                      }`}>
                        {flight.on_ground ? "Ground" : "Air"}
                      </span>
                    </td>
                    <td className="px-4 py-3 text-sm text-gray-300">{flight.origin_country}</td>
                    <td className="px-4 py-3 text-sm text-gray-300">{flight.geo_altitude.toFixed(0)} m</td>
                    <td className="px-4 py-3 text-sm text-gray-300">{(flight.velocity * 3.6).toFixed(1)} km/h</td>
                    <td className="px-4 py-3 text-sm text-gray-300">{flight.true_track.toFixed(0)}°</td>
                    <td className="px-4 py-3 text-sm text-gray-300">{flight.latitude.toFixed(2)}, {flight.longitude.toFixed(2)}</td>
                  </tr>
                ))}
              </tbody>
            </table>
          ) : (
            <div className="text-center text-gray-500 py-10">
              {status === "Connected"
                ? filteredFlights.length === 0 && flights.length > 0
                  ? "No flights match current filters"
                  : "Waiting for flight data..."
                : "No flight data available."}
            </div>
          )}
        </div>
      )}
    </div>
  );
}
