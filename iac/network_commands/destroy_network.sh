#!/bin/bash
# Stops all containers on the simulated network and removes it.
# Usage: ./destroy_network.sh

set -e

: "${SIMULATED_NETWORK_NAME:=eaglechat_lan}"

if ! docker network ls --format '{{.Name}}' | grep -q "^${SIMULATED_NETWORK_NAME}$"; then
  echo "Network '${SIMULATED_NETWORK_NAME}' does not exist. Nothing to do."
  exit 0
fi

CONTAINERS=$(docker ps -a --filter "network=${SIMULATED_NETWORK_NAME}" -q)

if [ -n "$CONTAINERS" ]; then
  echo "Stopping and removing containers on network '${SIMULATED_NETWORK_NAME}'..."
  docker rm --force $CONTAINERS
else
  echo "No containers found on network '${SIMULATED_NETWORK_NAME}'."
fi

echo "Removing network '${SIMULATED_NETWORK_NAME}'..."
docker network rm "$SIMULATED_NETWORK_NAME"

echo "Network destroyed."
