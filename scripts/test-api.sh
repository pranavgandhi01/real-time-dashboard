#!/bin/bash

# API Testing Script for Real-Time Flight Tracker

BASE_URL="http://localhost:8080"
TOKEN="84cdb6c6cdaba1ca7a862f158bc5afb07729b90361c6086f8a5947e3d6aacc3c"

echo "🚀 Testing Real-Time Flight Tracker API"
echo "========================================"
echo ""

# Test Health Endpoint
echo "1. Testing Health Endpoint..."
HEALTH_RESPONSE=$(curl -s -w "HTTPSTATUS:%{http_code}" "$BASE_URL/health")
HTTP_STATUS=$(echo $HEALTH_RESPONSE | tr -d '\n' | sed -e 's/.*HTTPSTATUS://')
RESPONSE_BODY=$(echo $HEALTH_RESPONSE | sed -e 's/HTTPSTATUS:.*//g')

if [ $HTTP_STATUS -eq 200 ]; then
    echo "✅ Health check passed"
    echo "   Response: $(echo $RESPONSE_BODY | jq -r '.status')"
else
    echo "❌ Health check failed (HTTP $HTTP_STATUS)"
fi
echo ""

# Test Readiness Endpoint
echo "2. Testing Readiness Endpoint..."
READY_RESPONSE=$(curl -s -w "HTTPSTATUS:%{http_code}" "$BASE_URL/ready")
HTTP_STATUS=$(echo $READY_RESPONSE | tr -d '\n' | sed -e 's/.*HTTPSTATUS://')
RESPONSE_BODY=$(echo $READY_RESPONSE | sed -e 's/HTTPSTATUS:.*//g')

if [ $HTTP_STATUS -eq 200 ]; then
    echo "✅ Readiness check passed"
    echo "   Response: $(echo $RESPONSE_BODY | jq -r '.status')"
else
    echo "❌ Readiness check failed (HTTP $HTTP_STATUS)"
fi
echo ""

# Test Metrics Endpoint
echo "3. Testing Metrics Endpoint..."
METRICS_RESPONSE=$(curl -s -w "HTTPSTATUS:%{http_code}" "$BASE_URL/metrics")
HTTP_STATUS=$(echo $METRICS_RESPONSE | tr -d '\n' | sed -e 's/.*HTTPSTATUS://')

if [ $HTTP_STATUS -eq 200 ]; then
    echo "✅ Metrics endpoint accessible"
    METRICS_COUNT=$(echo $METRICS_RESPONSE | sed -e 's/HTTPSTATUS:.*//g' | grep -c "^# HELP")
    echo "   Metrics available: $METRICS_COUNT"
else
    echo "❌ Metrics endpoint failed (HTTP $HTTP_STATUS)"
fi
echo ""

# Test API Documentation
echo "4. Testing API Documentation..."
DOCS_RESPONSE=$(curl -s -w "HTTPSTATUS:%{http_code}" "$BASE_URL/docs")
HTTP_STATUS=$(echo $DOCS_RESPONSE | tr -d '\n' | sed -e 's/.*HTTPSTATUS://')

if [ $HTTP_STATUS -eq 200 ]; then
    echo "✅ API documentation accessible"
    echo "   URL: $BASE_URL/docs"
else
    echo "❌ API documentation failed (HTTP $HTTP_STATUS)"
fi
echo ""

# Test OpenAPI Spec
echo "5. Testing OpenAPI Specification..."
SPEC_RESPONSE=$(curl -s -w "HTTPSTATUS:%{http_code}" "$BASE_URL/api-docs")
HTTP_STATUS=$(echo $SPEC_RESPONSE | tr -d '\n' | sed -e 's/.*HTTPSTATUS://')

if [ $HTTP_STATUS -eq 200 ]; then
    echo "✅ OpenAPI specification accessible"
    echo "   URL: $BASE_URL/api-docs"
else
    echo "❌ OpenAPI specification failed (HTTP $HTTP_STATUS)"
fi
echo ""

# Test WebSocket Connection (basic check)
echo "6. Testing WebSocket Endpoint..."
WS_RESPONSE=$(curl -s -w "HTTPSTATUS:%{http_code}" "$BASE_URL/ws?token=$TOKEN" \
    -H "Upgrade: websocket" \
    -H "Connection: Upgrade" \
    -H "Sec-WebSocket-Key: dGhlIHNhbXBsZSBub25jZQ==" \
    -H "Sec-WebSocket-Version: 13")
HTTP_STATUS=$(echo $WS_RESPONSE | tr -d '\n' | sed -e 's/.*HTTPSTATUS://')

if [ $HTTP_STATUS -eq 101 ] || [ $HTTP_STATUS -eq 400 ]; then
    echo "✅ WebSocket endpoint responding"
    echo "   Note: Use a WebSocket client for full testing"
else
    echo "❌ WebSocket endpoint failed (HTTP $HTTP_STATUS)"
fi
echo ""

echo "🏁 API Testing Complete!"
echo ""
echo "📚 Additional Resources:"
echo "   • Interactive API Docs: $BASE_URL/docs"
echo "   • Postman Collection: docs/postman-collection.json"
echo "   • Architecture Diagram: docs/architecture.md"
echo "   • WebSocket Test: Use a WebSocket client with token=$TOKEN"