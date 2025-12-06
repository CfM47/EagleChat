#!/bin/bash
# This script sets up a multi-manager, three-client environment using raw docker commands.

set -e

# --- Navigate to script directory ---
cd "$(dirname "$0")" # Stay in iac/scripts
# shellcheck source=./common.sh
source ./common.sh

# --- Configuration ---
MANAGER_NAMES=("id-manager-1" "id-manager-2")
CLIENT_NAMES=("client-1" "client-2" "client-3")
SHARED_ALIAS="id_manager.eaglechat.local"
ID_MANAGER_ENDPOINT="${SHARED_ALIAS}" # Clients connect to the shared alias

# --- Main Execution ---

# 1. Clean up old temporary certs and generate new ones
echo "====> Preparing credentials..."
rm -rf "$TMP_CERT_DIR"
generate_ca
generate_manager_credentials "${MANAGER_NAMES[0]}.eaglechat.local" "$TMP_CERT_DIR/${MANAGER_NAMES[0]}"
generate_manager_credentials "${MANAGER_NAMES[1]}.eaglechat.local" "$TMP_CERT_DIR/${MANAGER_NAMES[1]}"

# 2. Build images and create network
build_docker_images
create_network

echo -e "\n====> Starting ID Managers (2 instances)..."
for MANAGER_NAME in "${MANAGER_NAMES[@]}"; do
  echo "Starting manager: ${MANAGER_NAME}"
  docker run -d --rm \
    --name "${MANAGER_NAME}" \
    --network "${NETWORK_NAME}" \
    --network-alias "${SHARED_ALIAS}" \
    -v "$TMP_CERT_DIR/$MANAGER_NAME":/etc/eaglechat/certs:ro \
    -e ID_MANAGER_PRIV_KEY_PATH=/etc/eaglechat/certs/private_key.pem \
    -e ID_MANAGER_SIGNATURE_PATH=/etc/eaglechat/certs/id_manager_signature.pem \
    -e CA_PUBLIC_KEY_PATH=/etc/eaglechat/certs/ca_public_key.pem \
    -e COMMON_NAME="$MANAGER_NAME.eaglechat.local" \
    eaglechat-id-manager
done

echo -e "\n====> Starting Clients (3 instances)..."
for CLIENT_NAME in "${CLIENT_NAMES[@]}"; do
  echo "Starting client: ${CLIENT_NAME}"
  docker run -d --rm \
    --name "${CLIENT_NAME}" \
    --network "${NETWORK_NAME}" \
    -e ID_MANAGER_ENDPOINTS="${ID_MANAGER_ENDPOINT}" \
    eaglechat-client
done

echo -e "\n====> Multi-manager environment is now RUNNING."
echo "====> To stop and clean up, run 'iac/scripts/stop.sh'"
