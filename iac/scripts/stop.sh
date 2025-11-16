#!/bin/bash
# This script stops and removes all containers and the network created by the run-*.sh scripts.

set -e

# --- Configuration ---
NETWORK_NAME="eaglechat-net"
MANAGER_NAMES=("id-manager-1" "id-manager-2")
CLIENT_NAMES=("client-1" "client-2" "client-3")

# --- Navigate to script directory ---
cd "$(dirname "$0")/.." # Go up to iac directory

echo "====> Stopping and removing containers..."

# Collect all potential container names
ALL_CONTAINER_NAMES=("${MANAGER_NAMES[@]}" "${CLIENT_NAMES[@]}")

for CONTAINER_NAME in "${ALL_CONTAINER_NAMES[@]}"; do
  if docker ps -a --format '{{.Names}}' | grep -q "^${CONTAINER_NAME}$\"; then
    echo "Stopping and removing container: ${CONTAINER_NAME}"
    docker stop "${CONTAINER_NAME}" > /dev/null # Stop silently
    docker rm "${CONTAINER_NAME}" > /dev/null   # Remove silently
  else
    echo "Container '${CONTAINER_NAME}' not found, skipping."
  fi
done

echo "\n====> Removing network '${NETWORK_NAME}'..."
if docker network ls --format '{{.Name}}' | grep -q "^${NETWORK_NAME}$\"; then
  docker network rm "${NETWORK_NAME}"
else
  echo "Network '${NETWORK_NAME}' not found, skipping removal."
fi

echo "\n====> Environment stopped and cleaned up."
