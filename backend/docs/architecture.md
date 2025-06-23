# Architecture Diagrams

## System Architecture

![System Architecture](system-architecture.mmd)

**Components:**
- **Frontend**: Next.js application (Port 3000)
- **API Gateway**: Request routing and rate limiting (Port 8080)
- **Flight Data Service**: REST API for flight data (Port 8081)
- **WebSocket Service**: Real-time broadcasting (Port 8082)
- **OpenSky API**: External flight data source

## Data Flow Diagram

![Data Flow](data-flow.mmd)

**Flow:**
1. Frontend requests flight data via API Gateway
2. Gateway proxies to Flight Data Service
3. WebSocket connection for real-time updates
4. Continuous data streaming to frontend

## Service Communication

![Service Communication](service-communication.mmd)

**Communication Pattern:**
- Frontend communicates only with API Gateway
- Gateway routes requests to appropriate services
- Services are isolated and independently scalable