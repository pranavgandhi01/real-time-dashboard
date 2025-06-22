# Kafka with Strimzi on Kind

## Quick Start

```bash
# Setup Kafka cluster
./setup.sh

# Check status
kubectl get kafka -n kafka
kubectl get kafkatopic -n kafka

# Port forward for local access
kubectl port-forward svc/flight-tracker-kafka-kafka-bootstrap 9092:9092 -n kafka &
```

## Cleanup

```bash
kind delete cluster --name flight-tracker
```

## Configuration

- **Cluster**: flight-tracker-kafka
- **Topic**: flights (6 partitions)
- **Partition Key**: ICAO24
- **Retention**: 1 hour
- **Bootstrap**: localhost:9092

# Kafka Development Setup

This directory contains configurations for running Apache Kafka using Strimzi operator on Kind (Kubernetes in Docker).

## Prerequisites

- Docker
- Kind
- kubectl

## Setup

1. Run the setup script:

```bash
./setup-kind.sh
```

2. Verify the setup:

```bash
kubectl get kafka,kafkatopic -n kafka
```

3. Get Kafka connection details:

```bash
echo "Bootstrap server: localhost:9094"
```

## Configuration

- Kafka version: 3.6.0
- Replicas: 1 (development)
- External port: 9094
- Topic: flights
  - Partitions: 6
  - Key: ICAO24 (aircraft identifier)
  - Retention: 10 minutes

## Development Notes

- The cluster uses ephemeral storage - data is lost when restarting
- Auto topic creation is disabled - topics must be created explicitly
- Not suitable for production use - single replica, no persistence

## Monitoring

Access Kafka metrics:

```bash
curl localhost:9090/metrics
```

## Cleanup

Delete the Kind cluster:

```bash
kind delete cluster --name kafka-dev
```
