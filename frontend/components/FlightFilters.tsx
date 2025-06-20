interface FlightFiltersProps {
  countries: string[];
  selectedCountry: string;
  onCountryChange: (country: string) => void;
  statusFilter: 'all' | 'air' | 'ground';
  onStatusChange: (status: 'all' | 'air' | 'ground') => void;
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
    <div className="bg-gray-800 p-4 rounded-lg mb-6">
      <h3 className="text-lg font-semibold text-white mb-4">Filters</h3>
      <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
        {/* Country Filter */}
        <div>
          <label className="block text-sm font-medium text-gray-300 mb-2">
            Country
          </label>
          <select
            value={selectedCountry}
            onChange={(e) => onCountryChange(e.target.value)}
            className="w-full bg-gray-700 text-white rounded px-3 py-2 focus:outline-none focus:ring-2 focus:ring-blue-500"
          >
            <option value="">All Countries</option>
            {countries.map((country) => (
              <option key={country} value={country}>
                {country}
              </option>
            ))}
          </select>
        </div>

        {/* Status Filter */}
        <div>
          <label className="block text-sm font-medium text-gray-300 mb-2">
            Status
          </label>
          <select
            value={statusFilter}
            onChange={(e) => onStatusChange(e.target.value as 'all' | 'air' | 'ground')}
            className="w-full bg-gray-700 text-white rounded px-3 py-2 focus:outline-none focus:ring-2 focus:ring-blue-500"
          >
            <option value="all">All Flights</option>
            <option value="air">In Air</option>
            <option value="ground">On Ground</option>
          </select>
        </div>

        {/* Speed Filter */}
        <div>
          <label className="block text-sm font-medium text-gray-300 mb-2">
            Min Speed (km/h): {(minSpeed * 3.6).toFixed(0)}
          </label>
          <input
            type="range"
            min="0"
            max="300"
            value={minSpeed}
            onChange={(e) => onMinSpeedChange(Number(e.target.value))}
            className="w-full h-2 bg-gray-700 rounded-lg appearance-none cursor-pointer slider"
          />
        </div>
      </div>
    </div>
  );
}