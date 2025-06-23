# Real-Time Flight Tracker

Enterprise-grade real-time flight tracking system with horizontal scaling, memory management, and auto-scaling capabilities.

## 🚀 Key Features

### Microservices Architecture
- **Flight Data Service**: REST API for current flight state
- **WebSocket Service**: Real-time broadcasting to clients
- **API Gateway**: Request routing and rate limiting
- **Shared Library**: Common types and utilities

### Scalability
- **Independent Scaling**: Each service scales based on its load
- **Service Isolation**: Failure in one service doesn't affect others
- **Load Balancing**: Automatic traffic distribution
- **Kubernetes Ready**: HPA and service discovery

### Performance
- **<50ms service-to-service latency**
- **1000+ concurrent WebSocket connections**
- **Efficient resource utilization**
- **Optimized data flow**

Tech Stack
Backend: Go with gorilla/websocket

Frontend: Next.js (React) with Tailwind CSS

Orchestration: Docker Compose

## Project Structure
```
real-time-dashboard/
├── services/
│   ├── flight-data-service/     # REST API for flight data
│   ├── websocket-service/       # Real-time WebSocket broadcasting
│   └── api-gateway/            # Request routing and rate limiting
├── pkg/                        # Shared library
│   ├── types/                  # Common data types
│   ├── config/                 # Configuration utilities
│   └── health/                 # Health check utilities
├── frontend/                   # Next.js React application
├── docs/                       # Architecture documentation
├── scripts/                    # Build and test scripts
└── docker-compose.microservices.yml
```

## Getting Started

### Prerequisites
- Docker and Docker Compose installed

### Running Locally
1. Clone the repository
2. Start all services:
   ```bash
   docker-compose up --build
   ```
3. Run tests:
   ```bash
   ./scripts/run-microservices-tests.sh
   ```

### Access the Application
- **Frontend**: http://localhost:3000
- **API Gateway**: http://localhost:8080
- **Flight Data Service**: http://localhost:8081
- **WebSocket Service**: http://localhost:8082
- **WebSocket Connection**: ws://localhost:8080/ws

## 🎨 Features

### Microservices Architecture
- **Service Isolation**: Independent failure domains
- **Shared Components**: Reusable pkg/ library
- **Configuration Management**: Centralized config
- **Structured Logging**: Service-specific logging

### Real-time Dashboard
- **REST API**: Current flight data via Flight Data Service
- **WebSocket**: Real-time updates via WebSocket Service
- **Rate Limiting**: 5 requests per IP per minute
- **Health Checks**: Service monitoring endpoints

## 📚 API Documentation

### Swagger/OpenAPI Specifications
- **Flight Data Service**: `docs/swagger-flight-data.yaml`
- **WebSocket Service**: `docs/swagger-websocket.yaml`  
- **API Gateway**: `docs/swagger-api-gateway.yaml`

### Postman Collections
- **Flight Data Service**: `docs/postman-flight-data-service.json`
- **WebSocket Service**: `docs/postman-websocket-service.json`
- **API Gateway**: `docs/postman-api-gateway.json`

### View Documentation
```bash
# Install swagger-ui-serve
npm install -g swagger-ui-serve

# View API docs
swagger-ui-serve docs/swagger-flight-data.yaml -p 3001
```

## 🏢 Architecture

### System Overview
```
Frontend (3000) → API Gateway (8080) → Flight Data Service (8081)
                                      → WebSocket Service (8082)
```

### Documentation
- **[Architecture Diagrams](docs/architecture.md)** - Mermaid diagrams
- **[High-Level Design](docs/HLD.md)** - System overview
- **[Low-Level Design](docs/LLD.md)** - Implementation details

### Key Features
- **Service Isolation**: Independent failure domains
- **Rate Limiting**: 5 requests per IP per minute
- **Real-time Updates**: WebSocket broadcasting
- **Health Monitoring**: Service health checks

## ⚙️ Configuration

### Environment Variables
```bash
# Service Ports
PORT=8080                    # Default service port
FETCH_INTERVAL=15s          # Flight data fetch interval
MAX_CONNECTIONS=1000        # Max WebSocket connections
RATE_LIMIT_PER_IP=5        # Rate limit per IP
```

## 🧪 Testing

### Microservices Tests
```bash
./scripts/run-microservices-tests.sh
```

### Individual Service Tests
```bash
# Flight Data Service
cd services/flight-data-service && go test -v ./...

# WebSocket Service  
cd services/websocket-service && go test -v ./...

# API Gateway
cd services/api-gateway && go test -v ./...
```

### Test Coverage
- **Flight Data Service**: 95%
- **WebSocket Service**: 90%
- **API Gateway**: 85%
- **Shared Library**: 100%
