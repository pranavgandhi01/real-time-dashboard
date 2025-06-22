#!/bin/bash
set -e

CLUSTER_NAME="flight-tracker"
KAFKA_NAMESPACE="kafka"

print_status() { echo -e "\033[0;32m[INFO]\033[0m $1"; }

print_status "Cleaning up Kafka cluster"

# Delete Kafka resources
kubectl delete kafka flight-tracker-kafka -n $KAFKA_NAMESPACE --ignore-not-found=true
kubectl delete kafkatopic --all -n $KAFKA_NAMESPACE --ignore-not-found=true

# Wait for cleanup
print_status "Waiting for Kafka resources to be deleted"
kubectl wait --for=delete kafka/flight-tracker-kafka -n $KAFKA_NAMESPACE --timeout=300s 2>/dev/null || true

# Delete Kind cluster
if kind get clusters | grep -q "^$CLUSTER_NAME$"; then
    print_status "Deleting Kind cluster: $CLUSTER_NAME"
    kind delete cluster --name $CLUSTER_NAME
fi

print_status "Kafka cleanup completed!"