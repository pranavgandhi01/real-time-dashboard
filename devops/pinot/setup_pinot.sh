#!/bin/bash
set -e

CLUSTER_NAME="flight-tracker"
PINOT_NAMESPACE="pinot"

print_status() { echo -e "\033[0;32m[INFO]\033[0m $1"; }

print_status "Setting up Pinot on Kind cluster"

# Ensure kubectl context
kubectl config use-context kind-$CLUSTER_NAME

# Create namespace
print_status "Creating Pinot namespace"
kubectl apply -f k8s/namespace.yaml

# Deploy Pinot cluster
print_status "Deploying Pinot cluster"
kubectl apply -f k8s/pinot-cluster.yaml

# Wait for Pinot components
print_status "Waiting for Pinot components to be ready"
kubectl wait --for=condition=ready pod -l app=pinot-zookeeper -n $PINOT_NAMESPACE --timeout=300s
kubectl wait --for=condition=ready pod -l app=pinot-controller -n $PINOT_NAMESPACE --timeout=300s
kubectl wait --for=condition=ready pod -l app=pinot-broker -n $PINOT_NAMESPACE --timeout=300s
kubectl wait --for=condition=ready pod -l app=pinot-server -n $PINOT_NAMESPACE --timeout=300s

# Wait a bit for services to be fully ready
sleep 30

# Create schema and table
print_status "Creating flight events schema"
kubectl exec -n $PINOT_NAMESPACE deployment/pinot-controller -- /opt/pinot/bin/pinot-admin.sh AddSchema \
  -schemaFile /dev/stdin \
  -exec < schemas/flight-events-schema.json || echo "Schema may already exist"

print_status "Creating flight events table"
kubectl exec -n $PINOT_NAMESPACE deployment/pinot-controller -- /opt/pinot/bin/pinot-admin.sh AddTable \
  -tableConfigFile /dev/stdin \
  -exec < config/flight-events-table.json || echo "Table may already exist"

print_status "Pinot setup completed!"
echo "Pinot Controller UI: http://localhost:30900"
echo "Pinot Query Console: http://localhost:30099/query"