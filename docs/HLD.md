# High-Level Design: Real-Time Flight Dashboard

## System Overview
Microservices-based real-time flight tracking system with independent scaling and service isolation.

## Architecture Components

### Core Services
- **Flight Data Service** (Port 8081): Flight state management and REST API
- **WebSocket Service** (Port 8082): Real-time client connections
- **API Gateway** (Port 8080): Request routing and rate limiting
- **Frontend**: Next.js React application

### Shared Components
- **pkg/config**: Centralized configuration management
- **pkg/log**: Structured logging with configurable levels
- **pkg/middleware**: Rate limiting and common middleware
- **pkg/types**: Shared data structures
- **pkg/health**: Health check endpoints
- **pkg/client**: OpenSky API client

### Scalability Features
- **Service Isolation**: Independent failure domains
- **Rate Limiting**: 5 requests per IP per minute
- **Health Monitoring**: Service health checks
- **Configuration**: Environment-based settings

## Data Flow
```
OpenSky API → Flight Data Service → WebSocket Service → Frontend
                     ↓                       ↑
                REST API ← API Gateway ← Frontend
```

## Key Metrics
- **Throughput**: 1000+ concurrent WebSocket connections
- **Latency**: <50ms service-to-service communication
- **Availability**: 99.9% uptime with service isolation
- **Scaling**: Independent service scaling based on load

## Technology Stack
- **Flight Data Service**: Go 1.22, Gin, OpenSky API client
- **WebSocket Service**: Go 1.22, Gorilla WebSocket
- **API Gateway**: Go 1.22, Gin, Rate limiting middleware
- **Frontend**: Next.js 13, React 18, TypeScript
- **Shared Library**: Configuration, logging, middleware
- **Deployment**: Docker Compose