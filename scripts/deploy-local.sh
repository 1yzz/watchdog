#!/bin/bash

# Local Docker Deployment Script for Watchdog
set -e

CONTAINER_NAME="watchdog-local"
IMAGE_NAME="watchdog"
GRPC_PORT="50051"
HTTP_PORT="8080"
DATA_DIR="/var/lib/watchdog"

echo "🚀 Deploying Watchdog locally with Docker..."

# Build the Docker image
echo "📦 Building Docker image..."
docker build -t ${IMAGE_NAME}:latest .

# Stop and remove existing container
echo "🛑 Stopping existing container..."
docker stop ${CONTAINER_NAME} 2>/dev/null || true
docker rm ${CONTAINER_NAME} 2>/dev/null || true

# Create data directory
echo "📁 Creating data directory..."
sudo mkdir -p ${DATA_DIR}
sudo chown $(whoami):$(whoami) ${DATA_DIR}

# Run the new container
echo "🐳 Starting new container..."
docker run -d \
    --name ${CONTAINER_NAME} \
    --restart unless-stopped \
    -p ${GRPC_PORT}:${GRPC_PORT} \
    -p ${HTTP_PORT}:${HTTP_PORT} \
    -v ${DATA_DIR}:/data \
    -e DATABASE_URL=sqlite:///data/watchdog.db \
    -e LOG_LEVEL=info \
    -e GRPC_PORT=${GRPC_PORT} \
    -e HTTP_PORT=${HTTP_PORT} \
    ${IMAGE_NAME}:latest

# Wait for service to start
echo "⏳ Waiting for service to start..."
sleep 10

# Health check
echo "🔍 Checking service health..."
for i in {1..30}; do
    if curl -f -s http://localhost:${HTTP_PORT}/health > /dev/null 2>&1; then
        echo "✅ Watchdog is running successfully!"
        echo ""
        echo "🔗 Service endpoints:"
        echo "   gRPC: localhost:${GRPC_PORT}"
        echo "   HTTP: http://localhost:${HTTP_PORT}"
        echo ""
        echo "📊 Container status:"
        docker ps | grep ${CONTAINER_NAME}
        echo ""
        echo "📝 View logs with: docker logs ${CONTAINER_NAME}"
        echo "🛑 Stop with: docker stop ${CONTAINER_NAME}"
        exit 0
    fi
    echo "Waiting for service... (${i}/30)"
    sleep 2
done

echo "❌ Service failed to start. Check logs:"
docker logs ${CONTAINER_NAME}
exit 1
