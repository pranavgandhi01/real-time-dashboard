# Integration Tests

End-to-end testing for the real-time dashboard system.

## Quick Start

```bash
cd tests
npm install
npm test
```

## Tests

### WebSocket Integration
Tests WebSocket connection and data flow from backend services.

```bash
# Run with custom URL
WS_URL=ws://localhost:8080/ws npm test
```

## Environment Variables

- `WS_URL`: WebSocket server URL (default: ws://localhost:8080/ws)