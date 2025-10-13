#!/bin/bash
# Creates the simulated network if it doesn't already exist.
# Usage: ./start_network.sh <subnet>
# Example: ./start_network.sh 172.20.0.0/16

set -e

SUBNET=$1
: "${SIMULATED_NETWORK_NAME:=eaglechat_lan}"

if [ -z "$SUBNET" ]; then
  echo "Error: Subnet must be provided."
  echo "Usage: $0 <subnet>"
  exit 1
fi

if docker network ls --format '{{.Name}}' | grep -q "^${SIMULATED_NETWORK_NAME}$"; then
  echo "Network '${SIMULATED_NETWORK_NAME}' already exists."
else
  echo "Creating network '${SIMULATED_NETWORK_NAME}' with subnet ${SUBNET}..."
  docker network create \
    --driver bridge \
    --subnet="$SUBNET" \
    "$SIMULATED_NETWORK_NAME"
fi
