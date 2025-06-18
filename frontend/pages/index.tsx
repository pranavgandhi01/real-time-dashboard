// frontend/pages/index.tsx
import { useEffect, useState } from "react";

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

  useEffect(() => {
    // Get WebSocket URL from environment variable, or use default for local development
    const wsUrl =
      process.env.NEXT_PUBLIC_WEBSOCKET_URL || "ws://localhost:8080/ws";

    // Use the WebSocket API available in the browser
    const socket = new WebSocket(wsUrl);

    socket.onopen = () => {
      setStatus("Connected");
      console.log("WebSocket connection established");
    };

    socket.onmessage = (event) => {
      try {
        const data: FlightData[] = JSON.parse(event.data);
        // We might get null if the fetch failed, so we check for it.
        if (data) {
          setFlights(data);
        }
      } catch (error) {
        console.error("Failed to parse flight data:", error);
      }
    };

    socket.onclose = () => {
      setStatus("Disconnected");
      console.log("WebSocket connection closed");
    };

    socket.onerror = (error) => {
      setStatus("Error");
      console.error("WebSocket error:", error);
    };

    // Cleanup function to close the WebSocket connection when the component unmounts
    return () => {
      socket.close();
    };
  }, []); // Empty dependency array means this effect runs once on mount and cleans up on unmount

  return (
    <div className="min-h-screen bg-gray-900 text-white p-6">
      <header className="flex justify-between items-center mb-8">
        <h1 className="text-4xl font-extrabold text-transparent bg-clip-text bg-gradient-to-r from-teal-400 to-blue-500">
          Real-Time Flight Dashboard
        </h1>
        <span
          className={`text-lg font-semibold px-4 py-2 rounded-full ${
            status === "Connected"
              ? "bg-green-600"
              : status === "Disconnected"
              ? "bg-red-600"
              : "bg-yellow-600"
          }`}
        >
          Status: {status}
        </span>
      </header>

      <main>
        {flights.length === 0 && status === "Connected" && (
          <p className="text-center text-gray-500 text-xl">
            Waiting for flight data...
          </p>
        )}
        {status === "Error" && (
          <p className="text-center text-red-500 text-xl">
            Error connecting to WebSocket. Please check the backend server.
          </p>
        )}

        <div className="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-4">
          {flights.length > 0
            ? flights.map((flight) => (
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
                      <span className="font-semibold text-gray-300">
                        Altitude:
                      </span>{" "}
                      {flight.geo_altitude.toFixed(0)} m
                    </p>
                    <p>
                      <span className="font-semibold text-gray-300">
                        Speed:
                      </span>{" "}
                      {(flight.velocity * 3.6).toFixed(2)} km/h
                    </p>
                    <p>
                      <span className="font-semibold text-gray-300">
                        Heading:
                      </span>{" "}
                      {flight.true_track.toFixed(0)}°
                    </p>
                    <p>
                      <span className="font-semibold text-gray-300">
                        Vertical Rate:
                      </span>{" "}
                      {flight.vertical_rate.toFixed(1)} m/s
                    </p>
                    <p>
                      <span className="font-semibold text-gray-300">
                        Lat/Lon:
                      </span>{" "}
                      {flight.latitude.toFixed(4)},{" "}
                      {flight.longitude.toFixed(4)}
                    </p>
                  </div>
                </div>
              ))
            : status === "Connected" && (
                <p className="col-span-full text-center text-gray-400">
                  No flights currently in view, or data is being fetched...
                </p>
              )}
        </div>
      </main>
    </div>
  );
}
