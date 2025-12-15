#!/bin/bash
# This script stops and removes the Docker Compose environment.
#
# Usage:
#   ./stop-compose.sh            (stops the default single-manager scenario)
#   ./stop-compose.sh multi-manager (stops the multi-manager scenario)

set -e

# Navigate to the script's directory to ensure docker-compose is run from the correct context
cd "$(dirname "$0")"
source ./scripts/common.sh

SCENARIO=$1

echo "====> Tearing down Docker Compose environment..."

if [ "$SCENARIO" == "multi-manager" ]; then
  docker compose --profile multi-manager down --remove-orphans 
else
  docker compose down --remove-orphans
fi

echo "====> Environment stopped."

cleanup_credentials

echo "====> Verifying network cleanup..."
COMPOSE_NETWORK_NAME="iac_eaglechat-net" # Compose prefixes the project name (directory name)
remove_network "$COMPOSE_NETWORK_NAME"