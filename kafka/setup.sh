#!/bin/bash
set -e

echo "🚀 Setting up Strimzi Kafka on Kind..."

# Create Kind cluster
echo "Creating Kind cluster..."
kind create cluster --config=kind-cluster.yaml

# Install Strimzi operator
echo "Installing Strimzi operator..."
kubectl create namespace kafka
kubectl apply -f 'https://strimzi.io/install/latest?namespace=kafka' -n kafka

# Wait for operator
echo "Waiting for Strimzi operator..."
kubectl wait --for=condition=Ready pod -l name=strimzi-cluster-operator -n kafka --timeout=300s

# Create Kafka cluster
echo "Creating Kafka cluster..."
kubectl apply -f kafka-cluster.yaml

# Wait for Kafka
echo "Waiting for Kafka cluster..."
kubectl wait kafka/flight-tracker-kafka --for=condition=Ready --timeout=300s -n kafka

# Create topic
echo "Creating flights topic..."
kubectl apply -f flights-topic.yaml

echo "✅ Strimzi Kafka setup complete!"
echo "Kafka bootstrap: localhost:9092"
echo "Schema Registry: http://localhost:8081"