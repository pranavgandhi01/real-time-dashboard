# Shared Cluster Setup

Memory-efficient shared Kubernetes cluster for local development using Helm charts.

## Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                 flight-tracker namespace                    │
├─────────────────────────────────────────────────────────────┤
│  Kafka (512Mi)  │  Flink (512Mi)  │  Pinot (256Mi)         │
│  ZooKeeper      │  TaskManager    │  Controller/Broker     │
│  (256Mi)        │  (512Mi)        │  Server (512Mi)        │
└─────────────────────────────────────────────────────────────┘
```

## Resource Allocation

**Total Memory**: ~3Gi (vs 8Gi+ for separate clusters)
**Total CPU**: ~2 cores (vs 4+ cores for separate clusters)

### Per Service:
- **Kafka**: 512Mi memory, 250m CPU
- **Zookeeper**: 256Mi memory, 125m CPU  
- **Flink JobManager**: 512Mi memory, 250m CPU
- **Flink TaskManager**: 512Mi memory, 250m CPU
- **Pinot Controller**: 512Mi memory, 250m CPU
- **Pinot Broker**: 256Mi memory, 125m CPU
- **Pinot Server**: 512Mi memory, 250m CPU

## Quick Start

```bash
# Setup shared cluster
./setup.sh

# Check status
kubectl get all -n flight-tracker

# Cleanup
./cleanup.sh
```

## Helm Charts Used

- **Kafka**: Strimzi Kafka Operator (production-ready)
- **Flink**: Apache Flink Kubernetes Operator
- **Pinot**: Apache Pinot Helm chart

## External Access

- **Kafka**: localhost:32092
- **Flink UI**: localhost:30081 (when enabled)
- **Pinot Controller**: localhost:30900 (when enabled)

## Environment Variables

Configure resources via `devops/.env`:
- Memory/CPU limits per service
- Storage sizes
- External port mappings