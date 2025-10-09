#!/usr/bin/env bash
set -euo pipefail

# --- Configuration ---
IMAGE_NAME="kalininn/photoview"
IMAGE_TAG="latest"   # or e.g. "$(date +%Y%m%d)" for date-based tags

# --- Build ---
echo "🚀 Building image: $IMAGE_NAME:$IMAGE_TAG"
docker build -t ${IMAGE_NAME}:${IMAGE_TAG} .

# --- Tag ---
echo "🏷  Tagging image..."
docker tag ${IMAGE_NAME}:${IMAGE_TAG} "$IMAGE_NAME:$IMAGE_TAG"

# --- Push ---
echo "📤 Pushing image to registry..."
docker push "$IMAGE_NAME:$IMAGE_TAG"

echo "✅ Done: $IMAGE_NAME:$IMAGE_TAG pushed successfully."
