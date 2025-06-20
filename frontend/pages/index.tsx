// frontend/pages/index.tsx
import { useEffect, useState } from "react";
import pako from "pako";

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
  const [error, setError] = useState<string | null>(null); // New state for error messages

  useEffect(() => {
    const websocketUrl = `${
      process.env.NEXT_PUBLIC_WEBSOCKET_URL || "ws://localhost:8080/ws"
    }?token=${process.env.NEXT_PUBLIC_WEBSOCKET_TOKEN}`;
    let socket: WebSocket;
    let reconnectAttempts = 0;
    const maxReconnectAttempts = 5;
    const baseReconnectDelay = 2000; // 2 seconds

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
        setStatus("Disconnected");
        if (process.env.NODE_ENV === 'development') {
          console.info(`[WebSocket] Connection closed. Code: ${event.code}, Clean: ${event.wasClean}`);
        }
        if (!event.wasClean && reconnectAttempts < maxReconnectAttempts) {
          const delay = baseReconnectDelay * Math.pow(2, reconnectAttempts);
          setError(
            `Reconnecting in ${delay / 1000}s... (Attempt ${
              reconnectAttempts + 1
            }/${maxReconnectAttempts})`
          );
          setTimeout(() => {
            reconnectAttempts++;
            connect();
          }, delay);
        } else {
          setError(
            `WebSocket disconnected. Code: ${event.code}, Reason: ${event.reason}`
          );
        }
      };

      socket.onerror = (event) => {
        setStatus("Error");
        setError("WebSocket connection error. Retrying...");
        if (process.env.NODE_ENV === 'development') {
          console.error('[WebSocket] Connection error:', event);
        }
      };
    }

    connect();

    return () => {
      socket.close();
    };
  }, []);

  return (
    <div className="min-h-screen bg-gray-900 text-white p-6">
      <header className="flex justify-between items-center mb-6">
        <h1 className="text-4xl font-extrabold text-teal-400">
          Live Flight Tracker
        </h1>
        <div className="text-lg">
          Status:{" "}
          <span
            className={`font-semibold ${
              status === "Connected"
                ? "text-green-500"
                : status === "Disconnected"
                ? "text-yellow-500"
                : "text-red-500"
            }`}
          >
            {status}
          </span>
        </div>
      </header>

      {error && (
        <div className="bg-red-700 p-4 rounded-lg shadow-md mb-6 text-center text-lg">
          Error: {error}
        </div>
      )}

      <div className="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-4">
        {flights.length > 0 ? (
          flights.map((flight) => (
            <div
              key={flight.icao24}
              className="bg-gray-800 rounded-lg p-4 shadow-lg transform hover:scale-105 transition-transform duration-300"
            >
              <div className="flex justify-between items-center mb-2">
                <h2 className="text-xl font-bold text-teal-300">
                  {flight.callsign || "N/A"}
                </h2>
                <span
                  className={`px-2 py-1 text-xs font-semibold rounded-full ${
                    flight.on_ground
                      ? "bg-yellow-600 text-yellow-100"
                      : "bg-blue-600 text-blue-100"
                  }`}
                >
                  {flight.on_ground ? "On Ground" : "In Air"}
                </span>
              </div>
              <p className="text-gray-400 text-sm mb-4">
                From: {flight.origin_country}
              </p>
              <div className="text-sm space-y-2">
                <p>
                  <span className="font-semibold text-gray-300">Altitude:</span>{" "}
                  {flight.geo_altitude.toFixed(0)} m
                </p>
                <p>
                  <span className="font-semibold text-gray-300">Speed:</span>{" "}
                  {(flight.velocity * 3.6).toFixed(2)} km/h
                </p>
                <p>
                  <span className="font-semibold text-gray-300">Heading:</span>{" "}
                  {flight.true_track.toFixed(2)}°
                </p>
                <p>
                  <span className="font-semibold text-gray-300">
                    Vertical Rate:
                  </span>{" "}
                  {flight.vertical_rate.toFixed(2)} m/s
                </p>
                <p>
                  <span className="font-semibold text-gray-300">Lat/Lon:</span>{" "}
                  {flight.latitude.toFixed(4)}, {flight.longitude.toFixed(4)}
                </p>
              </div>
            </div>
          ))
        ) : (
          <div className="col-span-full text-center text-gray-500 text-xl py-10">
            {status === "Connected"
              ? "Waiting for flight data..."
              : "No flight data available."}
          </div>
        )}
      </div>
    </div>
  );
}
