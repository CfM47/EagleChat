#!/bin/bash
# This script sets up the dummy simulation environment and leaves it running.
# 1. Cleans up any previous runs.
# 2. Starts a new simulated network.
# 3. Builds the required Docker images.
# 4. Creates an id_manager and a client.

set -e

# --- Configuration ---
export SIMULATED_NETWORK_NAME="dummy_net"
SUBNET="172.22.0.0/16"
MANAGER_IP="172.22.0.10"
CLIENT_IP="172.22.0.20"

# --- Script Body ---
# Navigate to the script's directory to ensure paths are correct
cd "$(dirname "$0")"

echo "====> [Step 1/5] Tearing down any previous environment..."
./network_commands/destroy_network.sh || true # Allow to fail if network doesn't exist

echo "\n====> [Step 2/5] Starting a new simulated network '${SIMULATED_NETWORK_NAME}'..."
./network_commands/start_network.sh "$SUBNET"

echo "\n====> [Step 3/5] Building Docker images..."
./build_images.sh

echo "\n====> [Step 4/5] Creating id_manager at ${MANAGER_IP}..."
./network_commands/create.sh id_manager "$MANAGER_IP"

echo "\n====> [Step 5/5] Creating client at ${CLIENT_IP} (targeting manager)..."
./network_commands/create.sh client "$CLIENT_IP" "$MANAGER_IP"

echo "\n====> Dummy environment is now RUNNING. ===="
echo "Current status:"
docker ps
echo "\nRun './stop_dummy.sh' to tear down the environment."
