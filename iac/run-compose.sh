#!/bin/bash
# This script runs the Docker Compose environment for local testing.
# It sets the appropriate environment variables based on the selected scenario.
#
# Usage:
#   ./run-compose.sh            (starts the default single-manager scenario)
#   ./run-compose.sh multi-manager (starts the multi-manager scenario)

set -e

# Navigate to the script's directory to ensure docker-compose is run from the correct context
cd "$(dirname "$0")"

SCENARIO=$1

if [ "$SCENARIO" == "multi-manager" ]; then
  echo "====> Starting scenario: multi-manager"
  export ID_MANAGER_ENDPOINTS="id_manager.eaglechat.local"
  docker compose --profile multi-manager up --build -d 
else
  echo "====> Starting scenario: default (single-manager)"
  export ID_MANAGER_ENDPOINTS="http://id-manager-1:8080"
  docker compose up --build -d
fi

echo
echo "====> Environment is now RUNNING. View status with 'docker-compose ps'"
