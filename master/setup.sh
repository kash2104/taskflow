#!/bin/bash

set -e

echo "Creating namespace"
kubectl apply -f k8s/namespace.yml

echo "Secrets and ConfigMaps"
kubectl apply -f k8s/configmap.yml
kubectl apply -f k8s/secret.yml

echo "Deploying the pods"
kubectl apply -f k8s/deployment.yml

echo "creating http service"
kubectl apply -f k8s/http-service.yml

echo "creating grpc service"
kubectl apply -f k8s/grpc-service.yml
