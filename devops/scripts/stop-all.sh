#!/bin/bash
set -e

echo "🛑 Stopping Complete DevOps Infrastructure..."
echo ""

# Confirmation prompt
read -p "⚠️  This will stop all DevOps services. Continue? [y/N]: " -n 1 -r
echo
if [[ ! $REPLY =~ ^[Yy]$ ]]; then
    echo "❌ Cancelled by user"
    exit 0
fi

# Stop main application services first
echo "📦 Stopping application services..."
cd ../../
docker-compose down 2>/dev/null || echo "⚠️  No application services running"

# Stop Kafka (if running)
echo "📨 Checking for Kafka cluster..."
if command -v kubectl &> /dev/null && kind get clusters 2>/dev/null | grep -q flight-tracker; then
    echo "📨 Stopping Kafka cluster..."
    cd devops/kafka/strimzi
    ./cleanup_strimzi.sh
else
    echo "⏭️  No Kafka cluster found"
fi

# Stop observability stack
echo "📊 Stopping observability stack..."
pwd
cd devops/observability
docker-compose down

# Stop infrastructure services
echo "📦 Stopping infrastructure services..."
cd ../infrastructure
docker-compose down

# Optional cleanup
echo ""
read -p "🧹 Remove volumes and networks? [y/N]: " -n 1 -r
echo
if [[ $REPLY =~ ^[Yy]$ ]]; then
    echo "🧹 Cleaning up volumes and networks..."
    cd ../../
    cd devops/observability
    docker-compose down -v
    cd ../infrastructure
    docker-compose down -v
    
    # Clean up any orphaned containers
    docker container prune -f 2>/dev/null || true
    docker network prune -f 2>/dev/null || true
    
    echo "✅ Cleanup completed!"
fi

echo ""
echo "🎉 All DevOps Infrastructure Stopped!"
echo ""
echo "📊 Check remaining containers: docker ps"
echo "🔄 Restart everything: ./scripts/start-all.sh"