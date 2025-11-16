#!/bin/bash
# This script sets up a multi-manager, three-client environment using raw docker commands.

set -e

# --- Configuration ---
NETWORK_NAME="eaglechat-net"
MANAGER_NAMES=("id-manager-1" "id-manager-2")
CLIENT_NAMES=("client-1" "client-2" "client-3")
SHARED_ALIAS="id_manager.eaglechat.local"
ID_MANAGER_ENDPOINT="${SHARED_ALIAS}" # Clients connect to the shared alias

# --- Navigate to script directory ---
cd "$(dirname "$0")/.." # Go up to iac directory

echo "====> Building Docker images..."
docker build -t eaglechat-id-manager -f id_manager/Dockerfile ..
docker build -t eaglechat-client -f client/Dockerfile ..

echo "\n====> Creating network '${NETWORK_NAME}' if it doesn't exist..."
if ! docker network ls --format '{{.Name}}' | grep -q "^${NETWORK_NAME}$\"; then
  docker network create --driver overlay --attachable "$NETWORK_NAME"
else
  echo "Network '${NETWORK_NAME}' already exists."
fi

echo "\n====> Starting ID Managers (2 instances)..."
for MANAGER_NAME in "${MANAGER_NAMES[@]}"; do
  echo "Starting manager: ${MANAGER_NAME}"
  docker run -d --rm \
    --name "${MANAGER_NAME}" \
    --network "${NETWORK_NAME}" \
    --network-alias "${SHARED_ALIAS}" \
    eaglechat-id-manager
done

echo "\n====> Starting Clients (3 instances)..."
for CLIENT_NAME in "${CLIENT_NAMES[@]}"; do
  echo "Starting client: ${CLIENT_NAME}"
  docker run -d --rm \
    --name "${CLIENT_NAME}" \
    --network "${NETWORK_NAME}" \
    -e ID_MANAGER_ENDPOINTS="${ID_MANAGER_ENDPOINT}" \
    eaglechat-client
done

echo "\n====> Multi-manager environment is now RUNNING."
echo "To stop: run 'iac/scripts/stop.sh'"
