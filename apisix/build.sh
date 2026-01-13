#!/bin/bash
# Build script for apisix example

set -e

echo "Building apisix server image..."

# Build from repository root
cd "$(dirname "$0")/.."

# Build server image
echo "Building server image..."
docker build -f apisix/go-server/Dockerfile -t dubbo-go-apisix-server:latest .

echo "Build completed successfully!"
echo ""
echo "To run the example:"
echo "  1. Create docker network: docker network create default_network"
echo "  2. Start etcd: docker-compose -f apisix/deploy/etcd-compose/docker-compose.yml up -d"
echo "  3. Start APISIX: docker-compose -f apisix/deploy/apisix-compose/docker-compose.yml up -d"
echo "  4. Start Nacos: docker-compose -f apisix/deploy/nacos2.0.3-compose/docker-compose.yml up -d"
echo "  5. Start the server: docker run -d --name dubbo-go-apisix-server --network default_network -e NACOS_ADDR=<nacos_ip>:8848 dubbo-go-apisix-server:latest"
