#!/bin/bash
set -eo pipefail

# This file is the startup script for the full-nodes. Copies the correct binaries into
# the full-node home directories, and starts the node using `cosmovisor` to run `dydxprotocold`.
# Any arguments passed into this script is forwarded to `cosmovisor`.
# Example usage: ./start.sh run start --home chain/.full-node-1

source "./vars.sh"

# Set up CosmosVisor for full-nodes.
for i in $(seq 0 $LAST_FULL_NODE_INDEX); do
	FULL_NODE_HOME_DIR="$HOME/chain/.full-node-$i"
	# Copy binaries for `cosmovisor` from the docker image into the home directory.
	# Work-around to ensure docker volume contains the same binaries as the git repo.
	cp -r "$HOME/cosmovisor" "$FULL_NODE_HOME_DIR/"
done

cosmovisor "$@"
