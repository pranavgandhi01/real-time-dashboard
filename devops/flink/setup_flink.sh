#!/bin/bash
set -e

CLUSTER_NAME="flight-tracker"
FLINK_NAMESPACE="flink"

print_status() { echo -e "\033[0;32m[INFO]\033[0m $1"; }

print_status "Setting up Flink on Kind cluster"

# Ensure kubectl context
kubectl config use-context kind-$CLUSTER_NAME

# Create namespace
print_status "Creating Flink namespace"
kubectl apply -f k8s/namespace.yaml

# Install Flink Kubernetes Operator
print_status "Installing Flink Kubernetes Operator"
kubectl create -f https://github.com/apache/flink-kubernetes-operator/releases/download/release-1.7.0/flink-kubernetes-operator-1.7.0.yaml || echo "Operator already exists"

# Wait for operator
print_status "Waiting for Flink operator to be ready"
kubectl wait --for=condition=ready pod -l app.kubernetes.io/name=flink-kubernetes-operator -n flink-kubernetes-operator --timeout=300s

# Apply RBAC
print_status "Setting up RBAC for Flink"
kubectl apply -f k8s/rbac.yaml

# Deploy Flink cluster
print_status "Deploying Flink cluster"
kubectl apply -f k8s/flink-cluster.yaml

print_status "Waiting for Flink cluster to be ready"
kubectl wait --for=condition=ready pod -l app=flight-tracker-flink -n $FLINK_NAMESPACE --timeout=300s

print_status "Flink setup completed!"
echo "Flink UI: http://localhost:30081"
echo "Flink Metrics: http://localhost:30249"