#!/bin/bash
# This script builds the Docker images for all applications.

set -e

# Navigate to the script's directory to ensure paths are correct
cd "$(dirname "$0")"

echo $COMPOSE_BAKE

echo "Building eaglechat-id-manager image..."
docker build \
  --tag eaglechat-id-manager \
  --file ./id_manager/Dockerfile \
  ..

echo "\nBuilding eaglechat-client image..."
docker build \
  --tag eaglechat-client \
  --file ./client/Dockerfile \
  ..

echo
echo "Images built successfully."
