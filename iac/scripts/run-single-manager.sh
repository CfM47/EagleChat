#!/bin/bash
# This script sets up a single-manager, three-client environment using raw docker commands.

set -e

# --- Configuration ---
NETWORK_NAME="eaglechat-net"
MANAGER_NAME="id-manager-1"
CLIENT_NAMES=("client-1" "client-2" "client-3")
ID_MANAGER_ENDPOINT="http://${MANAGER_NAME}:8080" # Clients connect to the manager by its container name

CA_KEY="../../ca.key"
CA_PUB_KEY="../../ca_public_key.pem"
TMP_CERT_DIR="./tmp_certs"

MANAGER_1_CERT_DIR="$TMP_CERT_DIR/$MANAGER_NAME"

# --- Helper Functions (Copied from run-compose.sh, adjusted paths) ---
generate_ca() {
  if [ -f "$CA_KEY" ] && [ -f "$CA_PUB_KEY" ]; then
    echo "====> CA key pair already exists. Skipping generation."
    return
  fi
  echo "====> Generating new CA key pair..."
  openssl genrsa -out "$CA_KEY" 4096
  openssl rsa -in "$CA_KEY" -pubout -out "$CA_PUB_KEY"
  echo "CA key pair generated."
}

generate_manager_credentials() {
  local manager_cn="$1"
  local output_dir="$2"
  
  echo "====> Generating credentials for $manager_cn..."
  mkdir -p "$output_dir"

  local priv_key_path="$output_dir/private_key.pem"
  local pub_key_path="$output_dir/public_key.pem"
  local signature_path="$output_dir/id_manager_signature.pem"

  # 1. Generate ID Manager private and public key
  openssl genrsa -out "$priv_key_path" 4096
  openssl rsa -in "$priv_key_path" -pubout -out "$pub_key_path"

  # 2. Sign the ID Manager's public key with the CA private key using PSS padding
  openssl dgst -sha256 -sigopt rsa_padding_mode:pss -sign "$CA_KEY" -out "$signature_path" "$pub_key_path"

  # 3. Copy the CA's public key for the container to use
  cp "$CA_PUB_KEY" "$output_dir/ca_public_key.pem"
  
  echo "Credentials for $manager_cn created in $output_dir"
}


# --- Navigate to script directory ---
cd "$(dirname "$0")" # Stay in iac/scripts

# --- Main Execution ---

# 1. Clean up old temporary certs and generate new ones
echo "====> Preparing credentials..."
rm -rf "$TMP_CERT_DIR"
generate_ca
generate_manager_credentials "$MANAGER_NAME.eaglechat.local" "$MANAGER_1_CERT_DIR"


echo "====> Building Docker images..."
docker build -t eaglechat-id-manager -f ../id_manager/Dockerfile ../
docker build -t eaglechat-client -f ../client/Dockerfile ../

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
  -v "$(pwd)/$MANAGER_1_CERT_DIR":/etc/eaglechat/certs:ro \
  -e ID_MANAGER_PRIV_KEY_PATH=/etc/eaglechat/certs/private_key.pem \
  -e ID_MANAGER_SIGNATURE_PATH=/etc/eaglechat/certs/id_manager_signature.pem \
  -e CA_PUBLIC_KEY_PATH=/etc/eaglechat/certs/ca_public_key.pem \
  -e COMMON_NAME="$MANAGER_NAME.eaglechat.local" \
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
echo "====> To stop and clean up, run 'iac/scripts/stop.sh'"
