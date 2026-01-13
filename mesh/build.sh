#!/bin/bash
# Build script for mesh example

set -e

echo "Building mesh images for Kubernetes..."

# Build from repository root
cd "$(dirname "$0")/.."

# Build server image
echo "Building server image..."
docker build -f mesh/go-server/Dockerfile -t dubbo-go-mesh-provider:latest .

# Build client image
echo "Building client image..."
docker build -f mesh/go-client/Dockerfile -t dubbo-go-mesh-consumer:latest .

# If using minikube, load images into minikube
if command -v minikube &> /dev/null; then
    echo "Loading images into minikube..."
    minikube image load dubbo-go-mesh-provider:latest
    minikube image load dubbo-go-mesh-consumer:latest
fi

echo "Build completed successfully!"
echo ""
echo "To deploy:"
echo "  kubectl apply -f mesh/deploy/Namespace.yml"
echo "  kubectl apply -f mesh/deploy/provider/"
echo "  kubectl apply -f mesh/deploy/consumer/"


