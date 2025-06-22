#!/bin/bash
set -e

echo "🧪 Running Microservices Tests"

# Test services
for service in flight-data-service websocket-service api-gateway; do
    echo "Testing $service..."
    cd services/$service
    go test -v ./...
    cd ../..
done

# Test shared library
echo "Testing shared library..."
cd pkg
go test -v ./...
cd ..

echo "✅ All tests passed!"

# Coverage report
echo "📊 Coverage Report:"
for service in flight-data-service websocket-service api-gateway; do
    echo "$service:"
    cd services/$service && go test -cover ./... && cd ../..
done

echo "pkg:"
cd pkg && go test -cover ./... && cd ..