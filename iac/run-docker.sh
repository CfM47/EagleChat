#!/bin/bash

echo "cd ../ to go to project root"
cd ../

# Exit immediately if a command exits with a non-zero status.
set -e

# --- Configuration ---
# Network details from your 'ip route' command
SUBNET="10.37.129.0/24"
GATEWAY="10.37.129.84"
PARENT_INTERFACE="wlan0"

# IPs for the host bridge and containers. These are now exported to be used by docker-compose.
export HOST_BRIDGE_IP="10.37.129.214"
export ID_MANAGER_IP="10.37.129.215"
export CLIENT1_IP="10.37.129.216"
export CLIENT2_IP="10.37.129.217"

# --- Cleanup ---
echo "--- Cleaning up previous environment ---"
# Use 'docker compose down' to stop services and remove containers, networks
docker compose -f iac/docker-compose.local.yml down --remove-orphans || true
docker network rm eaglechat_net || true
sudo ip link del mynet-host || true
echo "Cleanup complete."
echo

# --- Build Images ---
echo "--- Building Docker images ---"
docker build -t eaglechat-client:latest -f iac/client/Dockerfile .
docker build -t eaglechat-id_manager:latest -f iac/id_manager/Dockerfile .
echo "Image build complete."
echo

# --- Network Setup ---
echo "--- Setting up macvlan network ---"
echo "CRITICAL WARNING: macvlan over Wi-Fi/hotspots is unreliable. If containers"
echo "cannot communicate, this is the likely cause. A wired connection is strongly recommended."
echo

# Create the macvlan network for containers
docker network create -d macvlan \
  --subnet="$SUBNET" \
  --gateway="$GATEWAY" \
  -o parent="$PARENT_INTERFACE" \
  eaglechat_net

# Create the macvlan bridge for host-to-container communication
echo "--- Creating host network bridge (requires sudo) ---"
sudo ip link add mynet-host link "$PARENT_INTERFACE" type macvlan mode bridge
sudo ip addr add "${HOST_BRIDGE_IP}/24" dev mynet-host
sudo ip link set mynet-host up
echo "Network setup complete."
echo

# --- Deploy with Docker Compose ---
echo "--- Deploying EagleChat services with docker compose ---"
docker compose -f iac/docker-compose.local.yml up -d
echo "Deployment complete."
echo

echo "Waiting 10s for services to start..."
sleep 10

# --- Verification ---
echo "--- Verifying id_manager status (via http://${ID_MANAGER_IP}:8080/status) ---"
# Curl the service using its dedicated IP. This also verifies host-to-container communication.
curl --fail http://${ID_MANAGER_IP}:8080/status || echo "Verification failed. This may be due to Wi-Fi client isolation."
echo
echo

# --- Setup Complete ---
echo "--- Setup Complete ---"
echo "Services are running. You can get a shell in a client using:"
echo "docker compose -f iac/docker-compose.local.yml exec client1 /bin/bash"

echo
echo "You can attach to the TUI of a client using:"
echo "docker attach $(docker compose -f iac/docker-compose.local.yml ps -q client1)"
echo "docker attach $(docker compose -f iac/docker-compose.local.yml ps -q client2)"
