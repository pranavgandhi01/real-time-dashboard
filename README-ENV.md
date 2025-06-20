# Environment Configuration Guide

## Development Setup

Use the default `.env` and `.env.frontend` files for local development:

```bash
# Start development environment
docker-compose up --build
```

## Production Setup

1. Copy production environment templates:
```bash
cp .env.production .env.prod
cp .env.frontend.production .env.frontend.prod
```

2. Update production values in `.env.prod` and `.env.frontend.prod`:
   - Change `WEBSOCKET_TOKEN` to a secure random string
   - Update `ALLOWED_ORIGINS` to your domain
   - Set proper TLS certificate paths
   - Configure production Kafka broker address

3. Start production environment:
```bash
docker-compose -f docker-compose.prod.yml up --build
```

## Environment Variables

### Backend (.env)
- `LOG_LEVEL`: Logging level (debug, info, warn, error, fatal)
- `PORT`: Server port (default: 8080)
- `KAFKA_BROKER_ADDRESS`: Kafka broker address
- `KAFKA_TOPIC`: Kafka topic name
- `KAFKA_GROUP_ID`: Kafka consumer group ID
- `OPEN_SKY_API_URL`: OpenSky API endpoint
- `USE_MOCK_DATA`: Use mock data instead of API (true/false)
- `ALLOWED_ORIGINS`: CORS allowed origins
- `WEBSOCKET_TOKEN`: WebSocket authentication token
- `TLS_CERT_PATH`: TLS certificate file path
- `TLS_KEY_PATH`: TLS private key file path
- `PROMETHEUS_ENABLED`: Enable Prometheus metrics (true/false)

### Frontend (.env.frontend)
- `NEXT_PUBLIC_WEBSOCKET_URL`: WebSocket server URL
- `NEXT_PUBLIC_WEBSOCKET_TOKEN`: WebSocket authentication token

## Schema Registry

The project uses Confluent Schema Registry for data validation:

```bash
# Register flight data schema
./scripts/register-schema.sh

# View registered schema
curl http://localhost:8081/subjects/flights-value/versions/latest
```

## Token Generation

Generate secure WebSocket tokens using the provided script:

```bash
# Generate and update all environment files
./scripts/generate-token.sh

# Or generate manually:
python3 -c "import secrets; print(secrets.token_hex(32))"
```

## Security Notes

- Always change default tokens in production
- Use HTTPS/WSS in production
- Restrict CORS origins to your domain
- Use proper TLS certificates