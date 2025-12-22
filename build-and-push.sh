#!/usr/bin/env bash
set -e

# IMAGE="bojankrlekrstic/krstenica-svc:latest"
IMAGE="bojankrlekrstic/krstenica-svc:version1.1.0"

echo "➡️ Building Docker image: $IMAGE"
docker build -t $IMAGE .

echo "➡️ Pushing image to Docker Hub..."
docker push $IMAGE

echo "✅ Done! Image is pushed to Docker Hub."
