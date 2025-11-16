#!/bin/bash
# This script sets up a single-manager, three-client environment using raw docker commands.

set -e

# --- Configuration ---
NETWORK_NAME="eaglechat-net"
MANAGER_NAME="id-manager-1"
CLIENT_NAMES=("client-1" "client-2" "client-3")
ID_MANAGER_ENDPOINT="http://${MANAGER_NAME}:8080" # Clients connect to the manager by its container name

# --- Navigate to script directory ---
cd "$(dirname "$0")/.." # Go up to iac directory

echo "====> Building Docker images..."
docker build -t eaglechat-id-manager -f id_manager/Dockerfile ..
docker build -t eaglechat-client -f client/Dockerfile ..

echo "\n====> Creating network '${NETWORK_NAME}' if it doesn't exist..."
if ! docker network ls --format '{{.Name}}' | grep -q "^${NETWORK_NAME}$"; then
  docker network create --driver overlay --attachable "$NETWORK_NAME"
else
  echo "Network '${NETWORK_NAME}' already exists."
fi

echo "\n====> Starting ID Manager: ${MANAGER_NAME}"
docker run -d --rm \
  --name "${MANAGER_NAME}" \
  --network "${NETWORK_NAME}" \
  --network-alias id_manager.eaglechat.local \
  eaglechat-id-manager

echo "\n====> Starting Clients (3 instances)..."
for CLIENT_NAME in "${CLIENT_NAMES[@]}"; do
  echo "Starting client: ${CLIENT_NAME}"
  docker run -d --rm \
    --name "${CLIENT_NAME}" \
    --network "${NETWORK_NAME}" \
    -e ID_MANAGER_ENDPOINTS="${ID_MANAGER_ENDPOINT}" \
    eaglechat-client
done

echo "\n====> Single-manager environment is now RUNNING."
echo "To stop: run 'iac/scripts/stop.sh'"
