interface FlightFiltersProps {
  countries: string[];
  selectedCountry: string;
  onCountryChange: (country: string) => void;
  statusFilter: "all" | "air" | "ground";
  onStatusChange: (status: "all" | "air" | "ground") => void;
  minSpeed: number;
  onMinSpeedChange: (speed: number) => void;
}

export default function FlightFilters({
  countries,
  selectedCountry,
  onCountryChange,
  statusFilter,
  onStatusChange,
  minSpeed,
  onMinSpeedChange,
}: FlightFiltersProps) {
  return (
    <div className="bg-gray-800 border border-gray-700 p-6 rounded-xl mb-8 shadow-lg">
      <div className="flex items-center gap-2 mb-6">
        <h3 className="text-xl font-bold text-white">
          <label className="block text-sm font-semibold text-gray-300 mb-3">
            🔍 Filters
          </label>
        </h3>
      </div>
      <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
        {/* Country Filter */}
        <div>
          <select
            value={selectedCountry}
            onChange={(e) => onCountryChange(e.target.value)}
            className="w-full bg-gray-700 border border-gray-600 text-white rounded-lg px-4 py-3 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500 transition-all"
          >
            <option value="">All Countries ({countries.length})</option>
            {countries.map((country) => (
              <option key={country} value={country}>
                {country}
              </option>
            ))}
          </select>
        </div>

        {/* Status Filter */}
        <div>
          <label className="block text-sm font-semibold text-gray-300 mb-3">
            ✈️ Flight Status
          </label>
          <select
            value={statusFilter}
            onChange={(e) =>
              onStatusChange(e.target.value as "all" | "air" | "ground")
            }
            className="w-full bg-gray-700 border border-gray-600 text-white rounded-lg px-4 py-3 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500 transition-all"
          >
            <option value="all">🌍 All Flights</option>
            <option value="air">🚁 In Air Only</option>
            <option value="ground">🛩️ On Ground Only</option>
          </select>
        </div>

        {/* Speed Filter */}
        <div>
          <div className="flex items-center gap-2 mb-2">
            <span className="text-lg">📊</span>
            <span className="text-sm font-semibold text-gray-300">
              Min Speed:
            </span>
            <input
              type="range"
              min="0"
              max="300"
              value={minSpeed}
              onChange={(e) => onMinSpeedChange(Number(e.target.value))}
              className="flex-1 h-2 bg-gray-700 rounded-lg appearance-none cursor-pointer slider focus:outline-none focus:ring-2 focus:ring-blue-500"
            />
          </div>
          <div className="flex justify-between text-xs text-gray-500 px-2">
            <span>0 km/h</span>
            <span className="font-medium text-gray-300">
              {(minSpeed * 3.6).toFixed(0)} km/h
            </span>
            <span>1080 km/h</span>
          </div>
        </div>
      </div>
    </div>
  );
}
