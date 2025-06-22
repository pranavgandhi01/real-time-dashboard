# Deployment Guide

## Prerequisites

### System Requirements
- **OS**: Linux, macOS, or Windows with WSL2
- **RAM**: 8GB minimum, 16GB recommended
- **Storage**: 10GB free space
- **Docker**: Version 20.10+
- **Docker Compose**: Version 2.0+

### Optional (for Kafka)
- **Kind**: Kubernetes in Docker
- **kubectl**: Kubernetes CLI
- **Kubernetes**: 4GB RAM allocation for Kind cluster

## Quick Start

### 1. Start All Services
```bash
cd devops
./scripts/start-all.sh
```

### 2. Verify Deployment
```bash
./scripts/status.sh
```

### 3. Access Services
- **Grafana**: http://localhost:3000 (admin/admin)
- **Jaeger**: http://localhost:16686
- **Prometheus**: http://localhost:9090
- **Kibana**: http://localhost:5601
- **Elasticsearch**: http://localhost:9200

## Step-by-Step Deployment

### Infrastructure Services

```bash
# Start Redis and core infrastructure
cd devops/infrastructure
docker-compose up -d

# Verify Redis
docker exec infrastructure-redis-1 redis-cli ping
```

### Observability Stack

```bash
# Start monitoring and logging
cd devops/observability
docker-compose up -d

# Wait for services to be ready
curl -f http://localhost:9090 # Prometheus
curl -f http://localhost:16686 # Jaeger
curl -f http://localhost:9200 # Elasticsearch
```

### Kafka Cluster (Optional)

```bash
# Install Kind and kubectl (if needed)
cd devops/kafka/kind
./install-kind.sh

# Deploy Kafka cluster
cd ../strimzi
./setup_strimzi.sh
```

## Service Configuration

### Prometheus Configuration
- **Config**: `devops/observability/prometheus.yml`
- **Targets**: API Gateway, Flight Data Service, WebSocket Service
- **Scrape Interval**: 15 seconds

### Grafana Setup
- **Dashboards**: Auto-provisioned from `grafana/dashboards/`
- **Data Sources**: Prometheus auto-configured
- **Default Login**: admin/admin

### Elasticsearch Configuration
- **Index Pattern**: `flight-tracker-logs-*`
- **Retention**: Managed by Filebeat configuration
- **Memory**: 512MB heap size

### Kafka Topics
- **flight-events**: 3 partitions, 1 day retention
- **flight-alerts**: 3 partitions, 7 days retention
- **Compression**: Snappy

## Environment-Specific Deployments

### Development Environment
```bash
# Minimal setup - infrastructure only
cd devops/infrastructure
docker-compose up -d

# Add observability
cd ../observability
docker-compose up -d
```

### Testing Environment
```bash
# Full stack with Kafka
./scripts/start-all.sh
# Answer 'y' to Kafka setup prompt
```

### Production Considerations

#### Resource Limits
```yaml
# Example production limits
resources:
  requests:
    memory: 2Gi
    cpu: 1000m
  limits:
    memory: 4Gi
    cpu: 2000m
```

#### Security Hardening
- Change default Grafana password
- Enable authentication for Elasticsearch
- Configure TLS for Kafka
- Use secrets management

#### High Availability
- Multi-node Kafka cluster
- Elasticsearch cluster
- Load balancer for services
- Persistent volume replication

## Troubleshooting

### Common Issues

#### Port Conflicts
```bash
# Check port usage
netstat -tulpn | grep :3000
lsof -i :3000

# Solution: Stop conflicting services or change ports
```

#### Docker Issues
```bash
# Clean up containers
docker system prune -f

# Reset Docker networks
docker network prune -f

# Check Docker resources
docker system df
```

#### Kafka Connection Issues
```bash
# Verify Kind cluster
kind get clusters

# Check Kafka status
kubectl get kafka -n kafka

# Port forward for debugging
kubectl port-forward svc/flight-tracker-kafka-kafka-bootstrap 9092:9092 -n kafka
```

#### Service Health Issues
```bash
# Check service logs
cd devops/observability
docker-compose logs -f elasticsearch

# Restart specific service
docker-compose restart grafana

# Check service dependencies
docker-compose ps
```

### Log Analysis

#### Application Logs
```bash
# View all logs
cd devops/observability
docker-compose logs -f

# Service-specific logs
docker-compose logs -f grafana
docker-compose logs -f elasticsearch
```

#### Kafka Logs
```bash
# Kafka cluster logs
kubectl logs -f deployment/strimzi-cluster-operator -n kafka

# Topic status
kubectl describe kafkatopic flight-events -n kafka
```

## Monitoring Deployment Health

### Automated Health Checks
```bash
# Run status check
./scripts/status.sh

# Kafka-specific status
./scripts/kafka-status.sh
```

### Manual Verification
```bash
# Test Prometheus metrics
curl http://localhost:9090/api/v1/targets

# Test Elasticsearch
curl http://localhost:9200/_cluster/health

# Test Kafka connectivity
kubectl exec -it flight-tracker-kafka-kafka-0 -n kafka -- bin/kafka-topics.sh --bootstrap-server localhost:9092 --list
```

## Cleanup and Maintenance

### Stop Services
```bash
# Stop all services
./scripts/stop-all.sh

# Stop specific stack
cd devops/observability
docker-compose down
```

### Data Cleanup
```bash
# Remove volumes (data loss!)
cd devops/observability
docker-compose down -v

# Clean Kafka cluster
cd devops/kafka/strimzi
./cleanup_strimzi.sh
```

### Maintenance Tasks
```bash
# Update images
docker-compose pull
docker-compose up -d

# Backup volumes
docker run --rm -v observability_elasticsearch-data:/data -v $(pwd):/backup alpine tar czf /backup/elasticsearch-backup.tar.gz /data

# Monitor disk usage
docker system df
df -h
```