#!/bin/bash

set -e

echo "Applying ConfigMaps"
kubectl apply -f k8s/configmap.yml

echo "Deploying the pods"
kubectl apply -f k8s/deployment.yml

echo "Creating horizontal pod autoscaler"
kubectl apply -f k8s/hpa.yml
