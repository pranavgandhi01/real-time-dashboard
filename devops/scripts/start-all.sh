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

# Start shared cluster (automatic if Kind/kubectl available)
echo ""
if command -v kind &> /dev/null && command -v kubectl &> /dev/null; then
    echo "📨 Starting shared cluster with Kafka, Flink, and Pinot..."
    cd ../shared-cluster
    ./setup.sh
    SHARED_CLUSTER_STARTED=true
else
    echo "⚠️  Kind/kubectl not found - installing..."
    echo "⚠️  Kind/kubectl not found - installing..."
    curl -Lo ./kind https://kind.sigs.k8s.io/dl/v0.20.0/kind-linux-amd64
    chmod +x ./kind && sudo mv ./kind /usr/local/bin/kind
    echo ""
    echo "📨 Starting shared cluster..."
    cd ../../shared-cluster
    ./setup.sh
    SHARED_CLUSTER_STARTED=true
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
if [[ "$SHARED_CLUSTER_STARTED" == "true" ]]; then
    echo "   Kafka (Streaming):    localhost:32092"
    echo "   Shared Cluster:       kubectl get all -n flight-tracker"
fi
echo ""
echo "📊 View running services: docker ps"
echo "📋 View logs: cd devops/observability && docker-compose logs -f"