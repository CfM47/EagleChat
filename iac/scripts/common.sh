#!/bin/bash

# This script contains common functions and variables for setting up the EagleChat environment.

set -e

# --- Common Variables ---
# The script using this common file is expected to be in a subdirectory of iac/
# The root of the iac directory
IAC_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
# Project root
REPO_ROOT="$(cd "$IAC_DIR/.." && pwd)"

CA_SK="$REPO_ROOT/ca.key"

export TMP_CERT_DIR="$IAC_DIR/tmp_certs"
export COMMON_CA_DIR="$TMP_CERT_DIR/common_ca"
export CA_PK="$COMMON_CA_DIR/ca_public_key.pem"
export NETWORK_NAME="eaglechat-net"
export MANAGER_ALIAS="id_manager.eaglechat.local"
export CT_NAME_PREFIX="eaglechat."

prefix_ct_name() {
  local name=$1

  echo "$CT_NAME_PREFIX$name"
}

# Generates a new Certificate Authority (CA) key pair if one doesn't exist.
generate_ca() {
  mkdir -p "$COMMON_CA_DIR"
  if [ -f "$CA_SK" ] && [ -f "$CA_PK" ]; then
    echo "====> CA key pair already exists. Skipping generation."
    return
  fi
  echo "====> Generating new CA key pair..."
  openssl genrsa -out "$CA_SK" 4096
  openssl rsa -in "$CA_SK" -pubout -out "$CA_PK"
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
  openssl dgst -sha256 -sigopt rsa_padding_mode:pss -sign "$CA_SK" -out "$signature_path" "$pub_key_path"

  # 3. Copy the CA's public key for the container to use
  cp "$CA_PK" "$output_dir/ca_public_key.pem"

  echo "Credentials for $manager_cn created in $output_dir"
}

create_manager() {
  local name=$1
  name=$(prefix_ct_name "$name")

  local credentials_path="$TMP_CERT_DIR/$name"

  generate_manager_credentials "$name.eaglechat.local" "$credentials_path"

  echo "Starting manager: ${name}"
  docker run -d --rm \
    --name "${name}" \
    --network "${NETWORK_NAME}" \
    --network-alias "${MANAGER_ALIAS}" \
    -v "$credentials_path":/etc/eaglechat/certs:ro \
    -e ID_MANAGER_PRIV_KEY_PATH=/etc/eaglechat/certs/private_key.pem \
    -e ID_MANAGER_SIGNATURE_PATH=/etc/eaglechat/certs/id_manager_signature.pem \
    -e CA_PUBLIC_KEY_PATH=/etc/eaglechat/certs/ca_public_key.pem \
    -e COMMON_NAME="$name.eaglechat.local" \
    eaglechat-id-manager
}

start_client() {
  local name=$1
  name=$(prefix_ct_name "$name")

  echo "Starting client: ${name}"
  docker run -it \
    --name "${name}" \
    --network "${NETWORK_NAME}" \
    -v "$CA_PK":/etc/eaglechat/certs/ca_public_key.pem:ro \
    -e CA_PUBLIC_KEY_PATH=/etc/eaglechat/certs/ca_public_key.pem \
    -e TZ="${TZ:-America/Havana}" \
    eaglechat-client
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

# Stops and removes all containers with the 'eaglechat.' prefix.
stop_and_remove_containers() {
  echo "====> Stopping and removing all containers with 'eaglechat.' prefix..."
  container_ids=$(docker ps -a --filter "name=$CT_NAME_PREFIX" -q)

  if [ -z "$container_ids" ]; then
    echo "No containers with prefix 'eaglechat.' found."
    return
  fi

  echo "Stopping and removing containers..."
  docker stop $container_ids
  docker rm $container_ids
  echo "Container cleanup complete."
}

# Stops and removes a single container by its short name.
# $1: The short name of the container (e.g., "id-manager-1").
remove_container() {
  local name=$1
  if [ -z "$name" ]; then
    echo "Error: container name not provided."
    return 1
  fi
  local prefixed_name
  prefixed_name=$(prefix_ct_name "$name")

  echo "====> Attempting to stop and remove container: $prefixed_name"
  if docker ps -a --format '{{.Names}}' | grep -q "^${prefixed_name}$"; then
    docker stop "$prefixed_name" >/dev/null
    docker rm "$prefixed_name" >/dev/null
    echo "Container '$prefixed_name' removed successfully."
  else
    echo "Container '$prefixed_name' not found, skipping."
  fi
}

# Lists the short names of all containers with the 'eaglechat.' prefix.
list_containers() {
  echo "====> Listing all containers with prefix '$CT_NAME_PREFIX'..."
  docker ps -a --filter "name=$CT_NAME_PREFIX" --format '{{.Names}}' | sed "s/^$CT_NAME_PREFIX//"
}

remove_network() {
  echo -e "\n====> Removing network '${NETWORK_NAME}'..."

  if docker network ls --format '{{.Name}}' | grep -q "^${NETWORK_NAME}$"; then
    docker network rm "${NETWORK_NAME}" || echo "Warning: Could not remove network '${NETWORK_NAME}'. It may still be in use."
  else
    echo "Network '${NETWORK_NAME}' not found, skipping removal."
  fi
}

# Removes temporary credentials.
cleanup_credentials() {
  echo -e "\n====> Cleaning up temporary credentials..."
  rm -f "$CA_SK" "$CA_PK"
  rm -rf "$TMP_CERT_DIR"
  echo "Temporary credentials removed."
}

initialize() {
  local build_arg=$1
  generate_ca
  if [ "$build_arg" == "--build" ]; then
    build_docker_images
  else
    echo "====> Skipping image builds. Use --build to force a build."
  fi
  create_network

  echo "Initialization complete"
}

cleanup() {
  stop_and_remove_containers
  remove_network
  cleanup_credentials

  echo "Cleanup complete"
}

echo "common.sh sourced successfully"
