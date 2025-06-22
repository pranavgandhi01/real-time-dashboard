#!/bin/bash
# Get Kafka external IP for Docker containers
KAFKA_IP=$(kubectl get nodes -o jsonpath='{.items[0].status.addresses[?(@.type=="InternalIP")].address}' 2>/dev/null || echo "localhost")
echo "KAFKA_BROKER=${KAFKA_IP}:32092"