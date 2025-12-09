#!/bin/bash
# This script runs the Docker Compose environment for local testing.
# It AUTOMATICALLY generates the required CA and ID Manager credentials.

set -e

# Navigate to the script's directory to ensure correct relative paths
cd "$(dirname "$0")"

# --- Configuration ---
CA_KEY="../ca.key"
CA_PUB_KEY="../ca_public_key.pem"
TMP_CERT_DIR="./tmp_certs"

# --- Helper Functions ---
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
  local manager_name="$1"
  local manager_dir="$TMP_CERT_DIR/$manager_name"
  
  echo "====> Generating credentials for $manager_name..."
  mkdir -p "$manager_dir"

  local priv_key_path="$manager_dir/private_key.pem"
  local pub_key_path="$manager_dir/public_key.pem"
  local signature_path="$manager_dir/id_manager_signature.pem"

  # 1. Generate ID Manager private and public key
  openssl genrsa -out "$priv_key_path" 4096
  openssl rsa -in "$priv_key_path" -pubout -out "$pub_key_path"

  # 2. Sign the ID Manager's public key with the CA private key using PSS padding
  openssl dgst -sha256 -sigopt rsa_padding_mode:pss -sign "$CA_KEY" -out "$signature_path" "$pub_key_path"

  # 3. Copy the CA's public key for the container to use
  cp "$CA_PUB_KEY" "$manager_dir/ca_public_key.pem"
  
  echo "Credentials for $manager_name created in $manager_dir"
}

# --- Main Execution ---

# 1. Clean up old credentials and generate new ones
echo "====> Preparing credentials..."
rm -rf "$TMP_CERT_DIR"
generate_ca
generate_manager_credentials "id-manager-1"
generate_manager_credentials "id-manager-2" # For multi-manager scenario

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
