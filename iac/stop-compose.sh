#!/bin/bash
# This script stops and removes the Docker Compose environment.
#
# Usage:
#   ./stop-compose.sh            (stops the default single-manager scenario)
#   ./stop-compose.sh multi-manager (stops the multi-manager scenario)

set -e

# Navigate to the script's directory to ensure docker-compose is run from the correct context
cd "$(dirname "$0")"

SCENARIO=$1

echo "====> Tearing down Docker Compose environment..."

if [ "$SCENARIO" == "multi-manager" ]; then
  docker compose --profile multi-manager down --remove-orphans 
else
  docker compose down --remove-orphans
fi

echo "====> Environment stopped."

echo "====> Verifying network cleanup..."
NETWORK_NAME="iac_eaglechat-net" # Compose prefixes the project name (directory name)

if docker network ls --format '{{.Name}}' | grep -q "^${NETWORK_NAME}$"; then
  echo "Network '${NETWORK_NAME}' still exists. Attempting to remove it..."
  # This command might fail if a container is still attached, so we add || true to prevent script exit
  docker network rm "${NETWORK_NAME}" || echo "Warning: Could not remove network. It may still be in use by a dangling container."
else
  echo "Network removed successfully."
fi