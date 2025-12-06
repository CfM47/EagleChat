#!/bin/bash
# This script sets up a single-manager, three-client environment using raw docker commands.

set -e

# --- Navigate to script directory ---
cd "$(dirname "$0")" # Stay in iac/scripts
# shellcheck source=./common.sh
source ./common.sh

# --- Configuration ---
MANAGER_NAME="id-manager-1"
CLIENT_NAMES=("client-1" "client-2" "client-3")
ID_MANAGER_ENDPOINT="http://${MANAGER_NAME}:8080" # Clients connect to the manager by its container name
MANAGER_1_CERT_DIR="$TMP_CERT_DIR/$MANAGER_NAME"


# --- Main Execution ---

# 1. Clean up old temporary certs and generate new ones
echo "====> Preparing credentials..."
rm -rf "$TMP_CERT_DIR"
generate_ca
generate_manager_credentials "$MANAGER_NAME.eaglechat.local" "$MANAGER_1_CERT_DIR"

# 2. Build images and create network
build_docker_images
create_network

echo -e "\n====> Starting ID Manager: ${MANAGER_NAME}"
docker run -d --rm \
  --name "${MANAGER_NAME}" \
  --network "${NETWORK_NAME}" \
  --network-alias id_manager.eaglechat.local \
  -v "$MANAGER_1_CERT_DIR":/etc/eaglechat/certs:ro \
  -e ID_MANAGER_PRIV_KEY_PATH=/etc/eaglechat/certs/private_key.pem \
  -e ID_MANAGER_SIGNATURE_PATH=/etc/eaglechat/certs/id_manager_signature.pem \
  -e CA_PUBLIC_KEY_PATH=/etc/eaglechat/certs/ca_public_key.pem \
  -e COMMON_NAME="$MANAGER_NAME.eaglechat.local" \
  eaglechat-id-manager

echo -e "\n====> Starting Clients (3 instances)..."
for CLIENT_NAME in "${CLIENT_NAMES[@]}"; do
  echo "Starting client: ${CLIENT_NAME}"
  docker run -d --rm \
    --name "${CLIENT_NAME}" \
    --network "${NETWORK_NAME}" \
    -e ID_MANAGER_ENDPOINTS="${ID_MANAGER_ENDPOINT}" \
    eaglechat-client
done

echo -e "\n====> Single-manager environment is now RUNNING."
echo "====> To stop and clean up, run 'iac/scripts/stop.sh'"
