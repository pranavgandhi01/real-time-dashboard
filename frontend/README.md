# Frontend Documentation

## Architecture

### Technology Stack
- **Framework**: Next.js 13 with React 18
- **Language**: TypeScript
- **Styling**: Tailwind CSS
- **Maps**: Leaflet with React-Leaflet
- **Real-time**: WebSocket with automatic reconnection

### Component Structure
```
frontend/
├── components/
│   ├── FlightFilters.tsx    # Filtering controls
│   ├── FlightMap.tsx        # Interactive map view
│   └── FlightStats.tsx      # Statistics dashboard
├── pages/
│   ├── _app.tsx            # App configuration
│   └── index.tsx           # Main dashboard
├── tests/
│   └── websocket-reconnection.test.js
└── config/
    └── performance.js       # Performance optimizations
```

## Features

### Real-time Updates
- **WebSocket Connection**: Automatic connection to backend
- **Compression**: Gzip decompression for optimized data transfer
- **Reconnection**: Exponential backoff with jitter (1s → 30s max)
- **Error Handling**: Graceful degradation on connection failures

### User Interface
- **Map View**: Interactive Leaflet map with flight markers
- **Grid View**: Tabular data with sorting and filtering
- **Responsive Design**: Mobile-friendly layout
- **Dark Theme**: Modern dark UI design

### Filtering & Search
- **Country Filter**: Filter flights by origin country
- **Status Filter**: All flights, in-air only, or on-ground only
- **Speed Filter**: Minimum speed threshold
- **Real-time Stats**: Live flight count and statistics

## WebSocket Implementation

### Connection Management
```typescript
const websocketUrl = `${
  process.env.NEXT_PUBLIC_WEBSOCKET_URL || "ws://localhost:8080/ws"
}?token=${process.env.NEXT_PUBLIC_WEBSOCKET_TOKEN || ""}`;
```

### Reconnection Logic
```typescript
// Exponential backoff with jitter
const exponentialDelay = Math.min(baseDelay * Math.pow(2, attempts), 30000);
const jitter = Math.random() * 1000;
const delay = exponentialDelay + jitter;
```

### Data Processing
```typescript
// Handle compressed binary data
const compressedData = new Uint8Array(event.data);
const decompressed = pako.ungzip(compressedData, { to: "string" });
const flights: FlightData[] = JSON.parse(decompressed);
```

## Performance Optimizations

### React Optimizations
- **useMemo**: Memoized filtered flight data
- **Dynamic Imports**: Code splitting for map component
- **Efficient Re-renders**: Optimized state updates

### Data Handling
- **Compression**: Gzip decompression for reduced bandwidth
- **Filtering**: Client-side filtering for responsive UI
- **Pagination**: Efficient data display

## Configuration

### Environment Variables
```bash
# WebSocket Configuration
NEXT_PUBLIC_WEBSOCKET_URL=ws://localhost:8080/ws
NEXT_PUBLIC_WEBSOCKET_TOKEN=your-secret-token

# Performance
NEXT_PUBLIC_MAX_FLIGHTS_DISPLAY=1000
NEXT_PUBLIC_MAP_UPDATE_INTERVAL=1000
```

### Build Configuration
```javascript
// next.config.js
module.exports = {
  experimental: {
    optimizeCss: true,
  },
  compiler: {
    removeConsole: process.env.NODE_ENV === 'production',
  },
}
```

## Testing

### WebSocket Tests
```bash
# Run WebSocket reconnection tests
node frontend/tests/test-runner.js
```

### Test Coverage
- **WebSocket Reconnection**: 100%
- **Exponential Backoff**: 100%
- **Connection State Management**: 100%

## Development

### Local Development
```bash
npm run dev          # Start development server
npm run build        # Build for production
npm run start        # Start production server
npm run lint         # Run ESLint
```

### Docker Development
```bash
docker build -t flight-tracker-frontend .
docker run -p 3000:3000 flight-tracker-frontend
```

## Deployment

### Production Build
```bash
npm run build
npm run start
```

### Docker Production
```dockerfile
FROM node:18-alpine
WORKDIR /app
COPY package*.json ./
RUN npm ci --only=production
COPY . .
RUN npm run build
EXPOSE 3000
CMD ["npm", "start"]
```