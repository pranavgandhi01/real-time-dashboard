#!/bin/bash
set -e

CLUSTER_NAME="flight-tracker"
NAMESPACE="flight-tracker"

print_status() { echo -e "\033[0;32m[INFO]\033[0m $1"; }

print_status "Cleaning up shared Flight Tracker cluster"

# Check cleanup mode
if [[ "$1" == "--full" ]]; then
    print_status "Full cleanup mode - removing Helm releases"
    
    # Uninstall Helm releases
    if helm list -n $NAMESPACE 2>/dev/null | grep -q strimzi-kafka-operator; then
        print_status "Uninstalling Strimzi Kafka Operator"
        helm uninstall strimzi-kafka-operator -n $NAMESPACE
    fi
    
    # Delete namespace
    print_status "Deleting namespace"
    kubectl delete namespace $NAMESPACE --ignore-not-found=true
else
    print_status "Quick cleanup mode - keeping cluster, removing resources"
    kubectl delete kafka,kafkatopic --all -n $NAMESPACE --ignore-not-found=true
fi

# Delete Kind cluster
if kind get clusters | grep -q "^$CLUSTER_NAME$"; then
    print_status "Deleting Kind cluster: $CLUSTER_NAME"
    kind delete cluster --name $CLUSTER_NAME
fi

print_status "Cleanup completed!"