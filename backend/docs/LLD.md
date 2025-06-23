# Low-Level Design: Real-Time Flight Dashboard

## Microservices Components

### Flight Data Service
```go
type FlightService struct {
    flights map[string]types.Flight
    mu      sync.RWMutex
    fetcher *client.FlightFetcher
    config  *config.Config
}
```

### WebSocket Service
```go
type WSService struct {
    clients map[*websocket.Conn]bool
}
```

### API Gateway
```go
type APIGateway struct {
    flightDataURL string
    websocketURL  string
    rateLimiter   *middleware.RateLimiter
}
```

## Shared Components

### pkg/types
```go
type Flight struct {
    ICAO24        string    `json:"icao24"`
    Callsign      string    `json:"callsign"`
    OriginCountry string    `json:"origin_country"`
    Longitude     float64   `json:"longitude"`
    Latitude      float64   `json:"latitude"`
    OnGround      bool      `json:"on_ground"`
    Velocity      float64   `json:"velocity"`
    LastUpdated   time.Time `json:"last_updated"`
}

type FlightStats struct {
    TotalFlights int       `json:"total_flights"`
    InAir        int       `json:"in_air"`
    OnGround     int       `json:"on_ground"`
    LastUpdated  time.Time `json:"last_updated"`
}
```

### pkg/config
```go
type Config struct {
    Port           string
    FetchInterval  time.Duration
    MaxConnections int
    RateLimitPerIP int
}
```

### pkg/log
```go
const (
    LogLevelDebug LogLevel = iota
    LogLevelInfo
    LogLevelWarn
    LogLevelError
    LogLevelFatal
)

func LogInfo(format string, v ...interface{})
func LogError(format string, v ...interface{})
func LogFatal(format string, v ...interface{})
```

## Frontend Components

### WebSocket Connection
```typescript
interface FlightData {
    icao24: string
    callsign: string
    latitude: number
    longitude: number
    on_ground: boolean
    velocity: number
}
```

### Reconnection Logic
- **Exponential Backoff**: 1s, 2s, 4s, 8s, 16s, 30s (max)
- **Max Attempts**: 5 reconnection attempts
- **Delay**: 3000ms base reconnection delay

## API Endpoints

### Flight Data Service (Port 8081)
- `GET /flights` - All current flights
- `GET /flights/{icao24}` - Specific flight
- `GET /stats` - Flight statistics
- `GET /health` - Service health

### WebSocket Service (Port 8082)
- `WS /ws` - WebSocket connection
- `GET /health` - Service health
- `GET /metrics` - Connection metrics

### API Gateway (Port 8080)
- `GET /flights/*` - Proxy to Flight Data Service
- `GET /stats` - Proxy to Flight Data Service
- `WS /ws` - Proxy to WebSocket Service
- `GET /health` - Gateway health

## Service Configuration

### Environment Variables
```bash
# Logging
LOG_LEVEL=debug

# Service Configuration
PORT=8080
FETCH_INTERVAL=15s
MAX_CONNECTIONS=1000
RATE_LIMIT_PER_IP=5

# External APIs
OPEN_SKY_API_URL=https://opensky-network.org/api/states/all

# CORS
ALLOWED_ORIGINS=http://localhost:3000
```

## Deployment

### Docker Compose
```yaml
services:
  flight-data-service:
    build: ./services/flight-data-service
    ports: ["8081:8081"]
    
  websocket-service:
    build: ./services/websocket-service
    ports: ["8082:8082"]
    
  api-gateway:
    build: ./services/api-gateway
    ports: ["8080:8080"]
```