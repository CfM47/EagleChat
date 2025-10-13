#!/bin/bash
# Finds and deletes a container by its IP address on the simulated network.
# Usage: ./delete.sh <ip>

set -e

IP=$1
: "${SIMULATED_NETWORK_NAME:=eaglechat_lan}"

if [ -z "$IP" ]; then
  echo "Error: IP address must be provided."
  exit 1
fi

# Find the container name/ID associated with the IP address
# We use a simple grep/sed combo for portability instead of requiring jq
CONTAINER_ID=$(docker network inspect "$SIMULATED_NETWORK_NAME" | \
  grep -B 4 "\"IPv4Address\": \"${IP}/" | \
  grep "Name" | \
  sed 's/.*\"\(.*\)\".*/\1/')

if [ -z "$CONTAINER_ID" ]; then
  echo "Error: No container found at IP ${IP} on network '${SIMULATED_NETWORK_NAME}'."
  exit 1
fi

echo "Found container '${CONTAINER_ID}' at ${IP}. Stopping and removing it..."
docker stop "$CONTAINER_ID"
echo "Container removed."
