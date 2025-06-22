#!/bin/bash
set -e

echo "🚀 Starting Complete DevOps Infrastructure..."
echo ""

# Function to check if service is ready
check_service() {
    local url=$1
    local name=$2
    local max_attempts=30
    local attempt=1
    
    echo "⏳ Waiting for $name to be ready..."
    while [ $attempt -le $max_attempts ]; do
        if curl -s $url > /dev/null 2>&1; then
            echo "✅ $name is ready!"
            return 0
        fi
        sleep 2
        attempt=$((attempt + 1))
    done
    echo "❌ $name failed to start within timeout"
    return 1
}

# Stop any existing containers first
echo "🧹 Cleaning up existing containers..."
cd ../../
docker-compose down 2>/dev/null || true

# Start infrastructure services
echo "📦 Starting infrastructure services..."
cd devops/infrastructure
docker-compose up -d

# Start observability stack
echo "📊 Starting observability stack..."
cd ../observability
docker-compose up -d

# Wait for core services
echo ""
echo "🔍 Checking service health..."
check_service "http://localhost:9090" "Prometheus"
check_service "http://localhost:16686" "Jaeger"
check_service "http://localhost:9200" "Elasticsearch"

# Start Kafka (optional - requires Kind)
echo ""
if command -v kind &> /dev/null && command -v kubectl &> /dev/null; then
    read -p "🤔 Start Kafka cluster? (requires Kind/kubectl) [y/N]: " -n 1 -r
    echo
    if [[ $REPLY =~ ^[Yy]$ ]]; then
        echo "📨 Starting Kafka cluster..."
        cd ../kafka/strimzi
        ./setup_strimzi.sh
        KAFKA_STARTED=true
    else
        echo "⏭️  Skipping Kafka setup"
        KAFKA_STARTED=false
    fi
else
    echo "⚠️  Kind/kubectl not found for Kafka setup"
    read -p "📦 Install Kind and kubectl automatically? [y/N]: " -n 1 -r
    echo
    if [[ $REPLY =~ ^[Yy]$ ]]; then
        cd ../kafka/kind
        pwd
        ./install-kind.sh
        echo ""
        read -p "📨 Now start Kafka cluster? [y/N]: " -n 1 -r
        echo
        if [[ $REPLY =~ ^[Yy]$ ]]; then
            cd ../strimzi
            ./setup_strimzi.sh
            KAFKA_STARTED=true
        else
            KAFKA_STARTED=false
        fi
    else
        echo "⏭️  Skipping Kafka setup"
        KAFKA_STARTED=false
    fi
fi

echo ""
echo "🎉 DevOps Infrastructure Started Successfully!"
echo ""
echo "🔗 Access Points:"
echo "   Jaeger (Tracing):     http://localhost:16686"
echo "   Prometheus (Metrics): http://localhost:9090"
echo "   Grafana (Dashboards): http://localhost:3000 (admin/admin)"
echo "   Kibana (Logs):        http://localhost:5601"
echo "   Elasticsearch (API):  http://localhost:9200"
echo "   Redis (Cache):        localhost:6379"
if [[ "$KAFKA_STARTED" == "true" ]]; then
    echo "   Kafka (Streaming):    localhost:32092"
fi
echo ""
echo "📊 View running services: docker ps"
echo "📋 View logs: cd devops/observability && docker-compose logs -f"