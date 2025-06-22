#!/bin/bash
set -e

# Check if kind is installed
if ! command -v kind &> /dev/null; then
    echo "kind is not installed. Please install it first."
    exit 1
fi

# Create Kind cluster
echo "Creating Kind cluster..."
kind create cluster --config kind/config.yaml

# Install Strimzi operator
echo "Installing Strimzi operator..."
kubectl create namespace kafka
kubectl create -f 'https://strimzi.io/install/latest?namespace=kafka' -n kafka

# Wait for Strimzi operator
echo "Waiting for Strimzi operator to be ready..."
kubectl wait deployment/strimzi-cluster-operator --for=condition=Available=True --timeout=300s -n kafka

# Create Kafka cluster
echo "Creating Kafka cluster..."
kubectl apply -f strimzi/kafka-cluster.yaml -n kafka

# Wait for Kafka cluster
echo "Waiting for Kafka cluster to be ready..."
kubectl wait kafka/flight-cluster --for=condition=Ready --timeout=300s -n kafka

# Create flights topic
echo "Creating flights topic..."
kubectl apply -f strimzi/kafka-topic.yaml -n kafka

echo "Setup complete!"
echo "Kafka bootstrap server: localhost:9094"
echo "To check status: kubectl get kafka,kafkatopic -n kafka"
