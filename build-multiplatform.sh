#!/bin/bash

# Simple multi-platform Docker build and push script
set -e

# Configuration  
REGISTRY="${DOCKER_REGISTRY:-serbenyuk}"
VERSION="${1:-latest}"
PLATFORMS="linux/amd64,linux/arm64"

echo "🐳 Building and pushing multi-platform images..."
echo "Registry: $REGISTRY"
echo "Version: $VERSION"
echo "Platforms: $PLATFORMS"
echo ""

# Build and push go-app
echo "📦 Building go-monitoring-app..."
cd go-app
docker buildx build \
    --platform $PLATFORMS \
    --tag $REGISTRY/go-monitoring-app:$VERSION \
    --tag $REGISTRY/go-monitoring-app:latest \
    --push \
    .
cd ..
echo "✅ go-monitoring-app pushed successfully"

# Build and push go-client  
echo "📦 Building go-monitoring-client..."
cd go-client
docker buildx build \
    --platform $PLATFORMS \
    --tag $REGISTRY/go-monitoring-client:$VERSION \
    --tag $REGISTRY/go-monitoring-client:latest \
    --push \
    .
cd ..
echo "✅ go-monitoring-client pushed successfully"

echo ""
echo "🎉 All images built and pushed successfully!"
echo ""
echo "📋 Available images:"
echo "   - $REGISTRY/go-monitoring-app:$VERSION"
echo "   - $REGISTRY/go-monitoring-client:$VERSION"
echo ""
echo "🚀 Deploy with:"
echo "   docker run -p 8000:8000 -p 8081:8081 $REGISTRY/go-monitoring-app:$VERSION"
echo "   docker run -p 8082:8082 $REGISTRY/go-monitoring-client:$VERSION"