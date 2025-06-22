# Production-like Kafka with Strimzi

Kubernetes-native Kafka setup using Strimzi operator for development that mirrors production architecture.

## Overview

- **Distribution**: Strimzi Kafka Operator
- **Platform**: Kind (Kubernetes in Docker)
- **Configuration**: Production-like with persistent storage, resource limits
- **Topics**: `flights` (6 partitions, 1h retention)
- **Access**: External NodePort (32092) + Internal Service

## Prerequisites

- Docker, Kind, kubectl installed
- 4GB+ RAM available for Kind cluster

## Setup

```bash
# Deploy Kafka cluster
./setup_strimzi.sh setup

# Verify deployment
kubectl get kafka,kafkatopic -n kafka

# Cleanup when done
./setup_strimzi.sh cleanup
```

## Production-like Features

- **Persistent Storage**: 10Gi Kafka, 5Gi Zookeeper
- **Resource Limits**: CPU/Memory constraints
- **External Access**: NodePort 32092
- **Topic Management**: Kubernetes CRDs
- **Monitoring**: Built-in JMX metrics

### Topics Configuration

- **flights**: 6 partitions, 1h retention, Snappy compression
- Auto-topic creation disabled for production safety

## Service Integration

```bash
# External access (from host)
KAFKA_BOOTSTRAP_SERVERS=localhost:32092

# Internal access (from pods)
KAFKA_BOOTSTRAP_SERVERS=flight-tracker-kafka-kafka-bootstrap.kafka.svc.cluster.local:9092
```

## Management

```bash
# View cluster status
kubectl get kafka -n kafka

# Topic management
kubectl get kafkatopic -n kafka
kubectl describe kafkatopic flights -n kafka

# Port forward for direct access
kubectl port-forward svc/flight-tracker-kafka-kafka-bootstrap 9092:9092 -n kafka
```

## Files

- `setup_strimzi.sh` - Setup/cleanup script
- `strimzi/kafka-cluster.yaml` - Production-like Kafka cluster
- `strimzi/kafka-topic.yaml` - Topic configuration
- `strimzi/namespace.yaml` - Kafka namespace
- `kind/config.yaml` - Kind cluster config
