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

  const stats = [
    { label: 'Total Flights', value: totalFlights, color: 'text-blue-400', icon: '✈️', bg: 'bg-blue-500/10' },
    { label: 'In Air', value: inAir, color: 'text-green-400', icon: '🚁', bg: 'bg-green-500/10' },
    { label: 'On Ground', value: onGround, color: 'text-yellow-400', icon: '🛩️', bg: 'bg-yellow-500/10' },
    { label: 'Avg Speed (km/h)', value: avgSpeed, color: 'text-purple-400', icon: '📊', bg: 'bg-purple-500/10' },
    { label: 'Avg Altitude (m)', value: avgAltitude, color: 'text-indigo-400', icon: '⛰️', bg: 'bg-indigo-500/10' },
    { label: 'Countries', value: countries, color: 'text-pink-400', icon: '🌍', bg: 'bg-pink-500/10' },
  ];

  return (
    <div className="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-6 gap-4 mb-6">
      {stats.map((stat, index) => (
        <div key={index} className={`bg-gray-800 border border-gray-700 p-3 rounded-lg`}>
          <div className="flex items-center gap-2">
            <span className="text-lg">{stat.icon}</span>
            <span className="text-sm text-gray-300">{stat.label}:</span>
            <span className={`text-lg font-bold ${stat.color}`}>{stat.value}</span>
          </div>
        </div>
      ))}
    </div>
  );
}