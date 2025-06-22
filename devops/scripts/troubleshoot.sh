#!/bin/bash

echo "🔍 Kafka Cluster Troubleshooting"
echo "================================"

# Check current context
echo "📋 Current Kubernetes Context:"
kubectl config current-context

echo ""
echo "🏗️ Namespace Status:"
kubectl get ns flight-tracker -o wide

echo ""
echo "📦 Strimzi Operator Status:"
kubectl get pods -n flight-tracker -l name=strimzi-cluster-operator

echo ""
echo "☕ Kafka Cluster Resources:"
kubectl get kafka -n flight-tracker -o wide

echo ""
echo "📊 All Pods in Namespace:"
kubectl get pods -n flight-tracker -o wide

echo ""
echo "📝 Kafka Cluster Events:"
kubectl get events -n flight-tracker --sort-by='.lastTimestamp' | tail -20

echo ""
echo "🔍 Kafka Resource Description:"
kubectl describe kafka kafka -n flight-tracker

echo ""
echo "📋 Strimzi Operator Logs (last 50 lines):"
kubectl logs -n flight-tracker -l name=strimzi-cluster-operator --tail=50

echo ""
echo "💾 Resource Usage:"
kubectl top pods -n flight-tracker 2>/dev/null || echo "Metrics server not available"

echo ""
echo "🔧 Cluster Node Status:"
kubectl get nodes -o wide

echo ""
echo "📊 Persistent Volumes:"
kubectl get pv | grep flight-tracker || echo "No PVs found for flight-tracker"

echo ""
echo "💿 Persistent Volume Claims:"
kubectl get pvc -n flight-tracker

echo ""
echo "⚙️ ConfigMaps and Secrets:"
kubectl get cm,secrets -n flight-tracker

echo ""
echo "🚪 Services:"
kubectl get svc -n flight-tracker