#!/bin/bash
set -e

CLUSTER_NAME="flight-tracker"
KAFKA_NAMESPACE="kafka"

print_status() { echo -e "\033[0;32m[INFO]\033[0m $1"; }

# Create Kind cluster if it doesn't exist
if ! kind get clusters | grep -q "^$CLUSTER_NAME$"; then
    print_status "Creating Kind cluster: $CLUSTER_NAME"
    kind create cluster --name $CLUSTER_NAME --config ../kind/config.yaml
fi

# Ensure we're using the correct kubectl context
print_status "Setting kubectl context to kind-$CLUSTER_NAME"
kubectl config use-context kind-$CLUSTER_NAME

# Create namespace
print_status "Creating Kafka namespace"
kubectl apply -f namespace.yaml

# Install Strimzi operator (skip if already exists)
print_status "Installing Strimzi operator"
kubectl create -f https://strimzi.io/install/latest?namespace=$KAFKA_NAMESPACE -n $KAFKA_NAMESPACE 2>/dev/null || echo "Strimzi operator already exists"

print_status "Waiting for Strimzi operator to be ready..."
kubectl wait --for=condition=ready pod -l name=strimzi-cluster-operator -n kafka --timeout=300s

print_status "Waiting for Kafka CRD to be available"
until kubectl get crd kafkas.kafka.strimzi.io >/dev/null 2>&1; do
    echo "Waiting for Kafka CRD..."
    sleep 5
done

print_status "Applying Kafka cluster configuration"
kubectl apply -f kafka-cluster.yaml

print_status "Applying Kafka topics"
kubectl apply -f kafka-topic.yaml

print_status "Waiting for Kafka cluster to be ready (this may take several minutes)"
echo "Note: Kafka cluster will continue deploying in background"
echo "Check status with: kubectl get kafka,kafkatopic -n kafka"

print_status "Kafka setup completed successfully!"
echo "External access: localhost:32092"