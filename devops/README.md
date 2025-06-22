# DevOps Infrastructure

Production-ready observability, messaging, and infrastructure services for the Real-Time Flight Dashboard.

## 📚 Documentation

- **[Architecture Overview](ARCHITECTURE.md)** - System design and component relationships
- **[Deployment Guide](DEPLOYMENT.md)** - Step-by-step deployment instructions
- **[Monitoring Guide](MONITORING.md)** - Observability, metrics, and alerting
- **[Security Guide](SECURITY.md)** - Security best practices and configurations
- **[Troubleshooting Guide](TROUBLESHOOTING.md)** - Common issues and solutions
- **[Port Allocation](PORTS.md)** - Service port mappings and conflicts

## 📁 Structure

```
devops/
├── observability/
│   ├── docker-compose.yml           # Monitoring & logging stack
│   ├── prometheus.yml               # Prometheus configuration
│   ├── filebeat.yml                 # Log collection config
│   └── grafana/                     # Dashboard configurations
├── infrastructure/
│   └── docker-compose.yml           # Core infrastructure services
├── kafka/                           # Kafka setup with Strimzi
│   ├── strimzi/                     # Kubernetes Kafka cluster
│   └── kind/                        # Local Kubernetes setup
├── scripts/                         # Management scripts
└── docs/                            # Documentation files
```

## Quick Start

### Start Observability Stack
```bash
cd devops/observability
docker-compose up -d
```

### Start Infrastructure Services
```bash
cd devops/infrastructure
docker-compose up -d
```

### Start Kafka (Kubernetes)
```bash
cd devops/kafka
./setup_strimzi.sh setup
```

### Start Everything
```bash
cd devops
./scripts/start-all.sh
```

### Check Status
```bash
cd devops
./scripts/status.sh
```

### Stop Everything
```bash
cd devops
./scripts/stop-all.sh
```

## 🚀 Services

### Observability Stack
- **Jaeger** (16686) - Distributed tracing and request flow analysis
- **Prometheus** (9090) - Metrics collection and time-series storage
- **Grafana** (3000) - Metrics visualization and alerting dashboards
- **Elasticsearch** (9200) - Centralized log storage and indexing
- **Kibana** (5601) - Log analysis and visualization
- **Filebeat** - Docker container log collection

### Infrastructure Services
- **Redis** (6379) - High-performance caching and session storage

### Event Streaming (Kubernetes)
- **Kafka** (9094/32092) - Distributed event streaming with Strimzi operator
- **Zookeeper** - Kafka cluster coordination
- **Strimzi Operator** - Kubernetes-native Kafka management

## 🔗 Access Points

| Service | URL | Credentials | Purpose |
|---------|-----|-------------|----------|
| Jaeger | http://localhost:16686 | - | Distributed tracing UI |
| Prometheus | http://localhost:9090 | - | Metrics and targets |
| Grafana | http://localhost:3000 | admin/admin | Dashboards and alerts |
| Kibana | http://localhost:5601 | - | Log analysis and search |
| Elasticsearch | http://localhost:9200 | - | REST API for logs |
| Redis | localhost:6379 | - | Cache and sessions |
| Kafka | localhost:32092 | - | Event streaming (external) |

## 🛠️ Management

### Quick Commands
```bash
# Start all infrastructure
./scripts/start-all.sh

# Check status of all services
./scripts/status.sh

# Stop all infrastructure
./scripts/stop-all.sh

# Kafka-specific status
./scripts/kafka-status.sh
```

### Service Management
```bash
# View logs
cd observability && docker-compose logs -f [service-name]

# Restart specific service
docker-compose restart [service-name]

# Scale services
docker-compose up -d --scale [service-name]=3
```

### Kafka Management
```bash
# Check Kafka cluster
kubectl get kafka,kafkatopic -n kafka

# View Kafka logs
kubectl logs -f deployment/strimzi-cluster-operator -n kafka

# Port forward for direct access
kubectl port-forward svc/flight-tracker-kafka-kafka-bootstrap 9092:9092 -n kafka
```

## 🔧 Prerequisites

- **Docker**: 20.10+ with Docker Compose v2
- **System**: 8GB RAM, 10GB storage
- **Optional**: Kind + kubectl for Kafka

## 📊 Monitoring Features

- **Real-time Metrics**: Application and infrastructure monitoring
- **Distributed Tracing**: Request flow across microservices
- **Centralized Logging**: Structured log aggregation and analysis
- **Custom Dashboards**: Pre-configured Grafana dashboards
- **Alerting**: Configurable alerts for critical metrics
- **Health Checks**: Automated service health monitoring

## 🔒 Security Features

- **Network Isolation**: Separate Docker networks for service tiers
- **Resource Limits**: CPU and memory constraints
- **Non-root Containers**: Security-hardened container execution
- **Persistent Storage**: Data retention with volume management
- **Access Control**: Authentication for sensitive services

For detailed information, see the [Security Guide](SECURITY.md).