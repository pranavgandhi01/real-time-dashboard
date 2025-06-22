# Microservices Documentation

## Service Overview

### Flight Data Service (Port 8081)
**Responsibility**: Maintain current flight state and provide REST API

**Key Features**:
- Fetches flight data from OpenSky API every 15 seconds
- Maintains in-memory flight state map
- Provides REST endpoints for flight data
- Real-time data updates

**Endpoints**:
- `GET /flights` - All current flights
- `GET /flights/{icao24}` - Specific flight
- `GET /stats` - Flight statistics
- `GET /health` - Health check

### WebSocket Service (Port 8082)
**Responsibility**: Real-time broadcasting to connected clients

**Key Features**:
- Manages WebSocket connections
- Broadcasts real-time updates to clients
- Auto-scales based on connection count
- Connection metrics monitoring

**Endpoints**:
- `WS /ws` - WebSocket connection
- `GET /metrics` - Connection metrics
- `GET /health` - Health check

### API Gateway (Port 8080)
**Responsibility**: Request routing and cross-cutting concerns

**Key Features**:
- Routes requests to appropriate services
- Rate limiting (5 requests per IP per minute)
- Authentication and authorization
- Load balancing and failover

**Routes**:
- `/api/flights/*` → Flight Data Service
- `/ws` → WebSocket Service
- `/health` → All services health

## Service Communication

### HTTP/REST
```
Frontend → API Gateway → Flight Data Service
Frontend → API Gateway → WebSocket Service
```

### Service Discovery
Services communicate using Docker network:
- `flight-data-service:8081`
- `websocket-service:8082`

## Deployment

### Development
```bash
docker-compose up --build
```

## Monitoring

### Health Checks
Each service exposes `/health` endpoint for monitoring

### Metrics
- Flight Data Service: API latency, flight count
- WebSocket Service: Connection count, message rate
- API Gateway: Request rate, error rate

### Logging
Structured logging with correlation IDs for distributed tracing