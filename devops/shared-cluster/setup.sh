#!/bin/bash
set -e

CLUSTER_NAME="flight-tracker"
NAMESPACE="flight-tracker"

print_status() { echo -e "\033[0;32m[INFO]\033[0m $1"; }

print_status "Setting up shared Flight Tracker cluster"

# Create Kind cluster if it doesn't exist
if ! kind get clusters | grep -q "^$CLUSTER_NAME$"; then
    print_status "Creating Kind cluster: $CLUSTER_NAME"
    kind create cluster --name $CLUSTER_NAME --config kind-config.yaml
fi

# Set kubectl context
kubectl config use-context kind-$CLUSTER_NAME

# Add Helm repositories
print_status "Adding Helm repositories"
helm repo add bitnami https://charts.bitnami.com/bitnami
helm repo update

# Create namespace and resource limits
print_status "Creating namespace with resource limits"
kubectl apply -f k8s/namespace.yaml

# Install Strimzi Kafka Operator
print_status "Installing Strimzi Kafka Operator"
helm upgrade --install strimzi-kafka-operator \
  --repo https://strimzi.io/charts/ strimzi-kafka-operator \
  --namespace $NAMESPACE \
  --create-namespace \
  --wait --timeout=300s

# Wait for operator to be ready
print_status "Waiting for Strimzi operator to be ready"
kubectl wait --for=condition=ready pod -l name=strimzi-cluster-operator -n $NAMESPACE --timeout=300s

# Create KafkaNodePool first
print_status "Creating KafkaNodePool"
kubectl apply -f - <<EOF
apiVersion: kafka.strimzi.io/v1beta2
kind: KafkaNodePool
metadata:
  name: broker
  namespace: $NAMESPACE
  labels:
    strimzi.io/cluster: kafka
spec:
  replicas: 1
  roles:
    - broker
    - controller
  storage:
    type: ephemeral
  resources:
    requests:
      memory: 512Mi
      cpu: 250m
    limits:
      memory: 1Gi
      cpu: 500m
EOF

# Apply Kafka cluster configuration with KRaft and KafkaNodePools
print_status "Deploying Kafka cluster with Strimzi (KRaft mode + NodePools)"
kubectl apply -f - <<EOF
apiVersion: kafka.strimzi.io/v1beta2
kind: Kafka
metadata:
  name: kafka
  namespace: $NAMESPACE
  annotations:
    strimzi.io/node-pools: enabled
    strimzi.io/kraft: enabled
spec:
  kafka:
    version: 3.9.0
    metadataVersion: 3.9-IV0
    listeners:
      - name: plain
        port: 9092
        type: internal
        tls: false
      - name: external
        port: 9094
        type: nodeport
        tls: false
        configuration:
          bootstrap:
            nodePort: 32092
    config:
      auto.create.topics.enable: "false"
      offsets.topic.replication.factor: 1
      transaction.state.log.replication.factor: 1
      min.insync.replicas: 1
      log.retention.hours: 24
      # KRaft specific configs
      process.roles: broker,controller
      controller.quorum.voters: 0@kafka-0.kafka-brokers.flight-tracker.svc.cluster.local:9093
      controller.listener.names: CONTROLLER
      inter.broker.listener.name: REPLICATION
EOF

# Wait for Kafka cluster to be ready
print_status "Waiting for Kafka cluster to be ready (KRaft mode)"
kubectl wait --for=condition=ready kafka/kafka -n $NAMESPACE --timeout=600s

# Create Kafka topics using Strimzi CRDs
print_status "Creating Kafka topics"
kubectl apply -f - <<EOF
apiVersion: kafka.strimzi.io/v1beta2
kind: KafkaTopic
metadata:
  name: flight-events
  namespace: $NAMESPACE
  labels:
    strimzi.io/cluster: kafka
spec:
  partitions: 3
  replicas: 1
  config:
    retention.ms: 86400000
    cleanup.policy: delete
    compression.type: snappy
---
apiVersion: kafka.strimzi.io/v1beta2
kind: KafkaTopic
metadata:
  name: flight-alerts
  namespace: $NAMESPACE
  labels:
    strimzi.io/cluster: kafka
spec:
  partitions: 3
  replicas: 1
  config:
    retention.ms: 604800000
    cleanup.policy: delete
    compression.type: snappy
EOF

# Wait and create topics directly in Kafka broker
print_status "Creating topics directly in Kafka broker"
sleep 15
kubectl exec -n $NAMESPACE kafka-broker-0 -- /opt/kafka/bin/kafka-topics.sh --bootstrap-server localhost:9092 --create --topic flight-events --partitions 3 --replication-factor 1 --if-not-exists 2>/dev/null || true
kubectl exec -n $NAMESPACE kafka-broker-0 -- /opt/kafka/bin/kafka-topics.sh --bootstrap-server localhost:9092 --create --topic flight-alerts --partitions 3 --replication-factor 1 --if-not-exists 2>/dev/null || true

print_status "Shared cluster setup completed!"
echo "Kafka external access: localhost:32092"
echo "Namespace: $NAMESPACE"
echo "Check status: kubectl get all -n $NAMESPACE"