#!/bin/bash

# Run backend tests

echo "Running backend unit tests..."
cd backend
go test ./tests/... -v

echo ""
echo "Running test coverage..."
go test ./tests/... -cover

echo ""
echo "Testing health endpoints..."
echo "Make sure the server is running on localhost:8080"
echo ""

# Test health endpoint
echo "Testing /health endpoint:"
curl -s http://localhost:8080/health | jq . || echo "Health endpoint test failed"

echo ""
echo "Testing /ready endpoint:"
curl -s http://localhost:8080/ready | jq . || echo "Ready endpoint test failed"

echo ""
echo "Testing /metrics endpoint:"
curl -s http://localhost:8080/metrics | head -5 || echo "Metrics endpoint test failed"

echo ""
echo "Tests completed!"