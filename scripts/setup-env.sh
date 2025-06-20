#!/bin/bash

# Environment setup script for Real-Time Dashboard

echo "Setting up environment configuration..."

# Check if production environment files exist
if [ ! -f ".env.prod" ]; then
    echo "Creating production environment file..."
    cp .env.production .env.prod
    echo "⚠️  Please update .env.prod with your production values"
fi

if [ ! -f ".env.frontend.prod" ]; then
    echo "Creating production frontend environment file..."
    cp .env.frontend.production .env.frontend.prod
    echo "⚠️  Please update .env.frontend.prod with your production values"
fi

# Generate random WebSocket token if not set
if grep -q "your-secret-token-change-in-production" .env.prod 2>/dev/null; then
    TOKEN=$(openssl rand -hex 32)
    sed -i "s/your-secret-token-change-in-production/$TOKEN/g" .env.prod
    sed -i "s/your-secret-token-change-in-production/$TOKEN/g" .env.frontend.prod
    echo "✅ Generated secure WebSocket token"
fi

echo "✅ Environment setup complete!"
echo ""
echo "Next steps:"
echo "1. Review and update .env.prod and .env.frontend.prod"
echo "2. For development: docker-compose up --build"
echo "3. For production: docker-compose -f docker-compose.prod.yml up --build"