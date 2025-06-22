#!/bin/bash

echo "🧪 Running Backend Test Suite"
echo "=============================="

cd backend

echo "📋 Running Unit Tests..."
go test ./ratelimit/... -v

echo ""
echo "🔧 Running Integration Tests..."
go test ./tests/... -v

echo ""
echo "⚡ Running Performance Tests..."
go test ./tests/... -bench=. -benchmem

echo ""
echo "📊 Running Tests with Coverage..."
go test ./... -cover -coverprofile=coverage.out

echo ""
echo "📈 Coverage Report:"
go tool cover -func=coverage.out

echo ""
echo "✅ Backend tests completed!"