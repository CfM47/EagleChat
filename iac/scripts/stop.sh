#!/bin/bash
# This script stops and removes all containers and the network created by the run-*.sh scripts.

set -e

# --- Navigate to script directory ---
cd "$(dirname "$0")" # Stay in iac/scripts
# shellcheck source=./common.sh
source ./common.sh

# --- Configuration ---
MANAGER_NAMES=("id-manager-1" "id-manager-2")
CLIENT_NAMES=("client-1" "client-2" "client-3")
ALL_CONTAINER_NAMES=("${MANAGER_NAMES[@]}" "${CLIENT_NAMES[@]}")

# --- Main Execution ---
stop_and_remove_containers "${ALL_CONTAINER_NAMES[@]}"
remove_network "$NETWORK_NAME"
cleanup_credentials

echo -e "\n====> Environment stopped and cleaned up."
