#!/bin/bash
# Creates a container of a specific type on the simulated network at a given IP.
# Usage: ./create.sh <type> <ip> [target_ip]
# Example: ./create.sh id_manager 172.20.0.10
# Example: ./create.sh client 172.20.0.20 172.20.0.10

set -e

TYPE=$1
IP=$2
TARGET_IP=$3
: "${SIMULATED_NETWORK_NAME:=eaglechat_lan}"

# --- Validate Input ---
if [ "$TYPE" != "id_manager" ] && [ "$TYPE" != "client" ]; then
  echo "Error: Invalid type '${TYPE}'. Must be 'id_manager' or 'client'."
  exit 1
fi
if [ -z "$IP" ]; then
  echo "Error: IP address must be provided."
  exit 1
fi
if [ "$TYPE" == "client" ] && [ -z "$TARGET_IP" ]; then
    echo "Error: A client requires a target_ip to connect to."
    exit 1
fi

# --- Check for IP Conflict ---
# The grep pattern looks for the exact IP/mask match in the network inspect output
if docker network inspect "$SIMULATED_NETWORK_NAME" | grep -q "\"IPv4Address\": \"${IP}/"; then
    echo "Error: IP address ${IP} is already in use on network '${SIMULATED_NETWORK_NAME}'."
    exit 1
fi

# --- Run Container ---
CONTAINER_NAME="${TYPE}-$(echo "$IP" | tr '.' '-')"
echo "Creating container '${CONTAINER_NAME}' at ${IP}..."

case $TYPE in
  id_manager)
    docker run \
      -d --rm \
      --network "$SIMULATED_NETWORK_NAME" \
      --ip "$IP" \
      --name "$CONTAINER_NAME" \
      eaglechat-id-manager
    ;;  
  client)
    docker run \
      -d --rm \
      --network "$SIMULATED_NETWORK_NAME" \
      --ip "$IP" \
      --name "$CONTAINER_NAME" \
      eaglechat-client /client "$TARGET_IP" 8080
    ;;esac

echo "Successfully started '${CONTAINER_NAME}'."
