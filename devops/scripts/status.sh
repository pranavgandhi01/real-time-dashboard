#!/bin/bash

echo "📊 DevOps Infrastructure Status"
echo "================================"
echo ""

# Check Docker containers
echo "🐳 Docker Containers:"
echo "---------------------"
if docker ps --format "table {{.Names}}\t{{.Status}}\t{{.Ports}}" | grep -E "(observability|infrastructure|real-time-dashboard)" 2>/dev/null; then
    echo ""
else
    echo "❌ No DevOps containers running"
    echo ""
fi

# Check Kafka
echo "📨 Kafka Cluster:"
echo "-----------------"
if command -v kubectl &> /dev/null; then
    if kubectl config current-context 2>/dev/null | grep -q "kind-flight-tracker"; then
        if kubectl get kafka -n kafka 2>/dev/null | grep -q flight-tracker-kafka; then
            echo "✅ Kafka cluster found:"
            kubectl get kafka,kafkatopic -n kafka 2>/dev/null || echo "❌ Kafka not accessible"
        else
            echo "❌ No Kafka cluster found"
        fi
    else
        echo "❌ Not connected to kind-flight-tracker context"
        echo "Run: kubectl config use-context kind-flight-tracker"
    fi
else
    echo "❌ kubectl not available"
fi
echo ""

# Check service endpoints
echo "🔗 Service Health:"
echo "------------------"
services=(
    "Jaeger:http://localhost:16686"
    "Prometheus:http://localhost:9090"
    "Grafana:http://localhost:3000"
    "Kibana:http://localhost:5601"
    "Elasticsearch:http://localhost:9200"
)

for service in "${services[@]}"; do
    name=$(echo $service | cut -d: -f1)
    url=$(echo $service | cut -d: -f2-)
    
    if curl -s --max-time 3 $url > /dev/null 2>&1; then
        echo "✅ $name - $url"
    else
        echo "❌ $name - $url (not responding)"
    fi
done

# Redis check - improved
echo -n "🔄 Redis - localhost:6379 "
if docker exec infrastructure-redis-1 redis-cli ping 2>/dev/null | grep -q PONG; then
    echo "✅"
elif redis-cli -h localhost -p 6379 ping 2>/dev/null | grep -q PONG; then
    echo "✅"
else
    echo "❌ (not responding)"
fi

# Kafka external access check
if kubectl get svc -n kafka 2>/dev/null | grep -q "flight-tracker-kafka-kafka-external-bootstrap"; then
    echo "✅ Kafka - localhost:32092 (external access available)"
else
    echo "⚠️  Kafka - localhost:32092 (checking...)"
fi

echo ""
echo "📋 Quick Commands:"
echo "------------------"
echo "Start all:    ./scripts/start-all.sh"
echo "Stop all:     ./scripts/stop-all.sh"
echo "Kafka status: ./scripts/kafka-status.sh"
echo "View logs:    cd ../observability && docker-compose logs -f"
echo "Port info:    cat ../PORTS.md"
echo "Kafka ctx:    kubectl config use-context kind-flight-tracker"