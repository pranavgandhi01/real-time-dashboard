#!/bin/bash

# Generate secure WebSocket token

TOKEN=$(python3 -c "import secrets; print(secrets.token_hex(32))")
echo "Generated secure token: $TOKEN"

# Update environment files if they exist
if [ -f ".env" ]; then
    sed -i "s/WEBSOCKET_TOKEN=.*/WEBSOCKET_TOKEN=$TOKEN/" .env
    echo "Updated .env"
fi

if [ -f ".env.frontend" ]; then
    sed -i "s/NEXT_PUBLIC_WEBSOCKET_TOKEN=.*/NEXT_PUBLIC_WEBSOCKET_TOKEN=$TOKEN/" .env.frontend
    echo "Updated .env.frontend"
fi

if [ -f ".env.prod" ]; then
    sed -i "s/WEBSOCKET_TOKEN=.*/WEBSOCKET_TOKEN=$TOKEN/" .env.prod
    echo "Updated .env.prod"
fi

if [ -f ".env.frontend.prod" ]; then
    sed -i "s/NEXT_PUBLIC_WEBSOCKET_TOKEN=.*/NEXT_PUBLIC_WEBSOCKET_TOKEN=$TOKEN/" .env.frontend.prod
    echo "Updated .env.frontend.prod"
fi

echo "✅ Token generation complete!"