Real-Time Flight Tracker
This project is a web application that displays real-time flight data on a dashboard. It's built with a Go backend that fetches data from the OpenSky Network API and pushes it to clients over WebSockets. The frontend is a Next.js application that renders the data.

Tech Stack
Backend: Go with gorilla/websocket

Frontend: Next.js (React) with Tailwind CSS

Orchestration: Docker Compose

Project Structure
real-time-dashboard/
├── backend/
│ ├── main.go
│ ├── ws/
│ │ └── hub.go
│ ├── fetcher/
│ │ └── flight.go
│ ├── go.mod
│ └── Dockerfile
├── frontend/
│ ├── pages/
│ │ └── index.tsx
│ ├── package.json
│ └── Dockerfile
├── docker-compose.yml
└── README.md

Getting Started
Prerequisites
- Docker installed on your machine
- Python 3 (for token generation)
- jq (for testing scripts)

Running Locally
1. Clone the repository and create the file structure:
   Make sure you have all the files provided placed in the correct directories as shown in the project structure above.

2. Generate secure tokens:
   ```bash
   ./scripts/generate-token.sh
   ```

3. Start Kafka services:
   ```bash
   cd kafka && docker-compose up -d
   ```

4. Register schema:
   ```bash
   ./scripts/register-schema.sh
   ```

5. Build and run the main services:
   ```bash
   docker-compose up --build
   ```

6. Run tests:
   ```bash
   ./scripts/run-tests.sh
   ```

This command will build the Docker images for both the backend and frontend services and then start them.

Access the application:

- **Frontend Dashboard**: http://localhost:3000
- **Backend WebSocket**: ws://localhost:8080/ws
- **Health Check**: http://localhost:8080/health
- **Readiness Check**: http://localhost:8080/ready
- **Metrics**: http://localhost:8080/metrics
- **Redis**: localhost:6379
- **Schema Registry**: http://localhost:8081
- **API Documentation**: http://localhost:8080/docs
- **OpenAPI Spec**: http://localhost:8080/api-docs

## Features

### Interactive Map View
- Real-time flight tracking on world map
- Flight status indicators (in-air vs on-ground)
- Click flights for detailed information
- Auto-zoom to fit all flights

### Advanced Filtering
- Filter by country of origin
- Filter by flight status (all/air/ground)
- Filter by minimum speed
- Real-time statistics dashboard

### Performance & Reliability
- Redis caching for improved performance
- Graceful shutdown handling
- Health check endpoints
- Schema validation with Avro
- Comprehensive error handling
- Unit test coverage

## API Documentation

### Interactive Documentation
Access the interactive Swagger UI at: http://localhost:8080/docs

### API Endpoints
- `GET /health` - Application health status
- `GET /ready` - Application readiness check
- `GET /metrics` - Prometheus metrics
- `GET /ws?token=<token>` - WebSocket connection

### Postman Collection
Import the Postman collection from `docs/postman-collection.json` for API testing.

### OpenAPI Specification
The complete API specification is available at:
- **YAML**: `docs/api-swagger.yaml`
- **Endpoint**: http://localhost:8080/api-docs

## Solution Architecture

View the complete solution architecture in `docs/architecture.md`.

### Key Components
- **Frontend**: Next.js with interactive Leaflet maps
- **Backend**: Go with WebSocket and REST API
- **Message Broker**: Apache Kafka with Schema Registry
- **Cache**: Redis for performance optimization
- **Monitoring**: Prometheus metrics and health checks

## Testing

Run the test suite:
```bash
./scripts/run-tests.sh
```

This will run:
- Unit tests for schema validation
- Health endpoint tests
- API endpoint verification
- Test coverage reports
