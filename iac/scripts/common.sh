#!/bin/bash

# This script contains common functions and variables for setting up the EagleChat environment.

set -e

# --- Common Variables ---
# The script using this common file is expected to be in a subdirectory of iac/
# The root of the iac directory
IAC_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
# Project root
REPO_ROOT="$(cd "$IAC_DIR/.." && pwd)"

NETWORK_NAME="eaglechat-net"
CA_KEY="$REPO_ROOT/ca.key"
CA_PUB_KEY="$REPO_ROOT/ca_public_key.pem"
TMP_CERT_DIR="$IAC_DIR/tmp_certs"

# --- Helper Functions ---

# Generates a new Certificate Authority (CA) key pair if one doesn't exist.
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

# Generates credentials (private/public key, signature) for an ID Manager.
# $1: The Common Name (CN) for the manager (e.g., "id-manager-1.eaglechat.local").
# $2: The output directory for the generated credentials.
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

# Builds the Docker images for the ID Manager and Client.
build_docker_images() {
    echo "====> Building Docker images..."
    docker build -t eaglechat-id-manager -f "$IAC_DIR/id_manager/Dockerfile" "$REPO_ROOT"
    docker build -t eaglechat-client -f "$IAC_DIR/client/Dockerfile" "$REPO_ROOT"
}

# Creates the Docker network if it doesn't already exist.
create_network() {
    echo -e "\n====> Creating network '${NETWORK_NAME}' if it doesn't exist..."
    if ! docker network ls --format '{{.Name}}' | grep -q "^${NETWORK_NAME}$"; then
      docker network create --driver overlay --attachable "$NETWORK_NAME"
    else
      echo "Network '${NETWORK_NAME}' already exists."
    fi
}

# Stops and removes a list of containers.
# $@: A list of container names to stop and remove.
stop_and_remove_containers() {
    local containers=("$@")
    echo "====> Stopping and removing containers..."
    for container in "${containers[@]}"; do
        if docker ps -a --format '{{.Names}}' | grep -q "^${container}$"; then
            echo "Stopping and removing container: ${container}"
            docker stop "${container}" > /dev/null
            docker rm "${container}" > /dev/null
        else
            echo "Container '${container}' not found, skipping."
        fi
    done
}

# Removes the Docker network.
# $1: The name of the network to remove.
remove_network() {
    local network_to_remove="$1"
    echo -e "\n====> Removing network '${network_to_remove}'..."
    if docker network ls --format '{{.Name}}' | grep -q "^${network_to_remove}$"; then
        docker network rm "${network_to_remove}" || echo "Warning: Could not remove network '${network_to_remove}'. It may still be in use."
    else
        echo "Network '${network_to_remove}' not found, skipping removal."
    fi
}

# Removes temporary credentials.
cleanup_credentials() {
    echo -e "\n====> Cleaning up temporary credentials..."
    rm -rf "$TMP_CERT_DIR"
    echo "Temporary credentials removed."
}

echo "common.sh sourced successfully"
