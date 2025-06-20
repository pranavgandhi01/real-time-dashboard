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

interface FlightStatsProps {
  flights: FlightData[];
}

export default function FlightStats({ flights }: FlightStatsProps) {
  const totalFlights = flights.length;
  const inAir = flights.filter(f => !f.on_ground).length;
  const onGround = flights.filter(f => f.on_ground).length;
  const avgSpeed = flights.length > 0 
    ? (flights.reduce((sum, f) => sum + f.velocity, 0) / flights.length * 3.6).toFixed(1)
    : '0';
  const avgAltitude = flights.length > 0
    ? (flights.reduce((sum, f) => sum + f.geo_altitude, 0) / flights.length).toFixed(0)
    : '0';

  const countries = Array.from(new Set(flights.map(f => f.origin_country))).length;

  return (
    <div className="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-6 gap-4 mb-6">
      <div className="bg-gray-800 p-4 rounded-lg text-center">
        <div className="text-2xl font-bold text-blue-400">{totalFlights}</div>
        <div className="text-sm text-gray-400">Total Flights</div>
      </div>
      <div className="bg-gray-800 p-4 rounded-lg text-center">
        <div className="text-2xl font-bold text-green-400">{inAir}</div>
        <div className="text-sm text-gray-400">In Air</div>
      </div>
      <div className="bg-gray-800 p-4 rounded-lg text-center">
        <div className="text-2xl font-bold text-yellow-400">{onGround}</div>
        <div className="text-sm text-gray-400">On Ground</div>
      </div>
      <div className="bg-gray-800 p-4 rounded-lg text-center">
        <div className="text-2xl font-bold text-purple-400">{avgSpeed}</div>
        <div className="text-sm text-gray-400">Avg Speed (km/h)</div>
      </div>
      <div className="bg-gray-800 p-4 rounded-lg text-center">
        <div className="text-2xl font-bold text-indigo-400">{avgAltitude}</div>
        <div className="text-sm text-gray-400">Avg Altitude (m)</div>
      </div>
      <div className="bg-gray-800 p-4 rounded-lg text-center">
        <div className="text-2xl font-bold text-pink-400">{countries}</div>
        <div className="text-sm text-gray-400">Countries</div>
      </div>
    </div>
  );
}