# Environment Configuration

## Development Setup

```bash
# Start all services
docker-compose up --build
```

## Environment Files

### Backend (.env)
```bash
PORT=8080                    # Service port
FETCH_INTERVAL=15s          # Flight data fetch interval
MAX_CONNECTIONS=1000        # Max WebSocket connections
RATE_LIMIT_PER_IP=5        # Rate limit per IP
OPEN_SKY_API_URL=https://opensky-network.org/api/states/all
ALLOWED_ORIGINS=http://localhost:3000
```

### Frontend (.env.frontend)
```bash
NEXT_PUBLIC_API_URL=http://localhost:8080
NEXT_PUBLIC_WEBSOCKET_URL=ws://localhost:8080/ws
NEXT_PUBLIC_MAX_DISPLAY_FLIGHTS=100
NEXT_PUBLIC_CLUSTER_DISTANCE=0.1
NEXT_PUBLIC_RECONNECT_DELAY=3000
NEXT_PUBLIC_MAX_RECONNECT_ATTEMPTS=5
```

## Service URLs

- **Frontend**: http://localhost:3000
- **API Gateway**: http://localhost:8080
- **Flight Data Service**: http://localhost:8081
- **WebSocket Service**: http://localhost:8082
- **WebSocket Connection**: ws://localhost:8080/ws