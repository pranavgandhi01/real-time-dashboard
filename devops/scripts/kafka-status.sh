#!/bin/bash

echo "📨 Kafka Cluster Status"
echo "======================"

if ! command -v kubectl &> /dev/null; then
    echo "❌ kubectl not found"
    exit 1
fi

if ! kind get clusters 2>/dev/null | grep -q flight-tracker; then
    echo "❌ Kind cluster 'flight-tracker' not found"
    exit 1
fi

echo "🔍 Checking Kafka resources..."
kubectl get kafka,kafkatopic -n kafka

echo ""
echo "🐳 Checking Kafka pods..."
kubectl get pods -n kafka

echo ""
echo "📊 Kafka Cluster Details:"
if kubectl get kafka flight-tracker-kafka -n kafka -o jsonpath='{.status.conditions[?(@.type=="Ready")].status}' 2>/dev/null | grep -q True; then
    echo "✅ Kafka Cluster: Ready"
else
    echo "⏳ Kafka Cluster: Not Ready (still deploying)"
fi

echo ""
echo "📋 Topics Status:"
kubectl get kafkatopic -n kafka -o custom-columns="NAME:.metadata.name,READY:.status.conditions[?(@.type=='Ready')].status"