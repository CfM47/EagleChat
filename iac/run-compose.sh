#!/bin/bash
# This script runs the Docker Compose environment for local testing.
# It AUTOMATICALLY generates the required CA and ID Manager credentials.

set -e

# Navigate to the script's directory to ensure correct relative paths
cd "$(dirname "$0")"
# Source common functions and variables
source ./scripts/common.sh

# --- Main Execution ---

# 1. Clean up old credentials and generate new ones
echo "====> Preparing credentials..."
rm -rf "$TMP_CERT_DIR"
generate_ca
generate_manager_credentials "id-manager-1.eaglechat.local" "$TMP_CERT_DIR/id-manager-1"
generate_manager_credentials "id-manager-2.eaglechat.local" "$TMP_CERT_DIR/id-manager-2" # For multi-manager scenario

# 2. Determine scenario and run Docker Compose
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
echo "====> To stop and clean up credentials, run 'iac/stop-compose.sh'"
