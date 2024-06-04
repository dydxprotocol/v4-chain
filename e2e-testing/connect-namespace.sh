#!/bin/bash

# Ensure the script exits on any command failure
set -e

BRIDGE=$1

if [ -z "$BRIDGE" ]; then
  echo "Error: BRIDGE parameter is missing"
  exit 1
fi

echo "Creating veth pair..."
sudo ip link add veth0 type veth peer name veth1

echo "Retrieving container ID..."
CONTAINER_ID=$(docker inspect -f '{{.State.Pid}}' interchain-security-instance)
if [ -z "$CONTAINER_ID" ]; then
  echo "Error: Failed to retrieve container ID"
  exit 1
fi

echo "Container ID: $CONTAINER_ID"

echo "Setting veth1 to container's network namespace..."
sudo ip link set veth1 netns $CONTAINER_ID

echo "Setting veth0 master to bridge $BRIDGE..."
sudo ip link set veth0 master $BRIDGE

echo "Bringing up veth0..."
sudo ip link set veth0 up

echo "Script completed successfully."