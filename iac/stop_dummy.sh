#!/bin/bash
# Tears down the environment created by run_dummy.sh.

set -e

# --- Configuration ---
# This MUST match the name used in run_dummy.sh
export SIMULATED_NETWORK_NAME="dummy_net"

# --- Script Body ---
# Navigate to the script's directory to ensure paths are correct
cd "$(dirname "$0")"

echo "====> Tearing down the '${SIMULATED_NETWORK_NAME}' environment..."
./network_commands/destroy_network.sh
