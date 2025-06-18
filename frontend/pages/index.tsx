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
    // Use the WebSocket API available in the browser
    // Replace with your actual backend host in production.
    const socket = new WebSocket("ws://localhost:8080/ws");

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

    // Cleanup function to close the socket when the component unmounts
    return () => {
      socket.close();
    };
  }, []); // Empty dependency array ensures this effect runs only once

  return (
    <div className="bg-gray-900 text-white min-h-screen font-sans">
      <header className="bg-gray-800 shadow-md p-4">
        <div className="container mx-auto flex justify-between items-center">
          <h1 className="text-3xl font-bold text-teal-400">
            Live Flight Tracker
          </h1>
          <div className="flex items-center space-x-2">
            <span
              className={`h-3 w-3 rounded-full ${
                status === "Connected"
                  ? "bg-green-500"
                  : status === "Connecting..."
                  ? "bg-yellow-500 animate-pulse"
                  : "bg-red-500"
              }`}
            ></span>
            <span className="text-gray-400">{status}</span>
          </div>
        </div>
      </header>

      <main className="container mx-auto p-4">
        <p className="text-gray-400 mb-6">
          Displaying {flights.length} flights in real-time. Data refreshes
          automatically.
        </p>
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
                    <span className="font-semibold text-gray-300">
                      Altitude:
                    </span>{" "}
                    {flight.geo_altitude.toFixed(0)} m
                  </p>
                  <p>
                    <span className="font-semibold text-gray-300">Speed:</span>{" "}
                    {(flight.velocity * 3.6).toFixed(2)} km/h
                  </p>
                  <p>
                    <span className="font-semibold text-gray-300">
                      Heading:
                    </span>{" "}
                    {flight.true_track.toFixed(0)}°
                  </p>
                </div>
              </div>
            ))
          ) : (
            <p className="text-gray-500 col-span-full text-center">
              {status === "Connected"
                ? "Waiting for initial flight data..."
                : "Attempting to connect to the server..."}
            </p>
          )}
        </div>
      </main>
    </div>
  );
}
