#!/bin/bash
#
# A simple CLI wrapper for managing the EagleChat environment.
# This script provides a clean interface to the functions in common.sh.
#
# Usage: ./mgr.sh <command> [arguments]
#
# Commands:
#   initialize [--build]      Initializes the environment, building images only if --build is specified.
#   create-manager <name>     Creates and starts a new ID manager container.
#   start-client <name>       Creates and starts a new client container interactively.
#   delete <name>             Stops and removes a specific container.
#   ls                        Lists the short names of all created containers.
#   cleanup                   Stops and removes all containers, the network, and credentials.
#

set -e

# Navigate to the script's directory to ensure correct relative paths
cd "$(dirname "$0")"
# Source common functions and variables
source ./scripts/common.sh

# --- Helper Functions ---

# Provides usage instructions.
usage() {
  echo "Usage: $0 <command> [arguments]"
  echo "Commands:"
  echo "   initialize [--build]      Initializes the environment, building images only if --build is specified."
  echo "   create-manager <name>     Creates and starts a new ID manager container."
  echo "   start-client <name>       Creates and starts a new client container interactively."
  echo "   delete <name>             Stops and removes a specific container."
  echo "   ls                        Lists the short names of all created containers."
  echo "   cleanup                   Stops and removes all containers, the network, and credentials."
  exit 1
}

# Combines all cleanup steps into a single command.
cleanup() {
  echo "====> Cleaning up the environment..."
  stop_and_remove_containers
  remove_network
  cleanup_credentials
  echo -e "\nEnvironment has been cleaned up."
}

# --- Main Command Dispatcher ---

# Check if a command was provided.
if [ -z "$1" ]; then
  usage
fi

case "$1" in
initialize)
  if [ -z "$2" ]; then
    initialize
  elif [ "$2" == "--build" ]; then
    initialize --build
  else
    echo "Error: Unknown argument '$2' for initialize."
    usage
  fi
  ;;
create-manager)
  if [ -z "$2" ]; then
    echo "Error: Missing <name> for create-manager."
    usage
  fi
  create_manager "$2"
  ;;
start-client)
  if [ -z "$2" ]; then
    echo "Error: Missing <name> for start-client."
    usage
  fi
  start_client "$2"
  ;;
delete)
  if [ -z "$2" ]; then
    echo "Error: Missing <name> for delete."
    usage
  fi
  remove_container "$2"
  ;;
ls)
  list_containers
  ;;
cleanup)
  cleanup
  ;;
*)
  echo "Error: Unknown command '$1'"
  usage
  ;;
esac

echo -e "\nOperation '$1' completed successfully."
