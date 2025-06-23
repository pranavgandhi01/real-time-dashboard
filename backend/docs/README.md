# Documentation Index

## Architecture Documentation
- **[Architecture Diagrams](architecture.md)** - System architecture, data flow, and service communication
- **[High-Level Design](HLD.md)** - System overview and components
- **[Low-Level Design](LLD.md)** - Detailed service implementation

## API Documentation

### Swagger/OpenAPI Specifications
- **[Flight Data Service API](swagger-flight-data.yaml)** - REST API for flight data
- **[WebSocket Service API](swagger-websocket.yaml)** - WebSocket service for real-time updates
- **[API Gateway](swagger-api-gateway.yaml)** - Central gateway with rate limiting

### Postman Collections
- **[Flight Data Service Collection](postman-flight-data-service.json)** - Test flight data endpoints
- **[WebSocket Service Collection](postman-websocket-service.json)** - Test WebSocket service
- **[API Gateway Collection](postman-api-gateway.json)** - Test gateway and rate limiting

## Service Documentation
- **[Services Overview](../services/README.md)** - Microservices documentation

## Quick Start

### Deploy Services
```bash
docker-compose up --build
```

### View Swagger Documentation
```bash
# Install swagger-ui-serve (if not installed)
npm install -g swagger-ui-serve

# Serve Flight Data Service API docs
swagger-ui-serve docs/swagger-flight-data.yaml -p 3001

# Serve API Gateway docs  
swagger-ui-serve docs/swagger-api-gateway.yaml -p 3002
```

### Import Postman Collections
1. Open Postman
2. Click "Import" 
3. Select the JSON files from `docs/postman-*.json`
4. Collections will be available in your workspace

### Test WebSocket Connection
```javascript
// Browser console or WebSocket client
const ws = new WebSocket('ws://localhost:8080/ws');
ws.onmessage = (event) => {
    console.log('Flight update:', JSON.parse(event.data));
};
```

## Architecture Overview

```
Frontend (3000) → API Gateway (8080) → Flight Data Service (8081)
                                    → WebSocket Service (8082)
```

### Key Features
- **Rate Limiting**: 5 requests per IP per minute
- **Health Checks**: All services expose `/health` endpoints
- **Real-time Updates**: WebSocket broadcasting
- **Service Isolation**: Independent scaling and failure domains