#!/bin/bash
set -eo pipefail

# This file initializes a non-validating full node for mainnet

source "./vars.sh"

CHAIN_ID="dydx-mainnet-1"

# Define dependencies for this script.
# `jq` and `dasel` are used to manipulate json and yaml files respectively.
install_prerequisites() {
	apk add dasel jq
}

set_cosmovisor_binary_permissions() {
    # Set up upgrade binaries.
    for version in "${!version_to_url[@]}"; do
        echo "Setting up version ${version}..."
        version_dir="$HOME/cosmovisor/upgrades/$version"
        mkdir -p "$version_dir/bin"
        url=${version_to_url[$version]}
        tar_file=$(basename $url)

        echo "Downloading tar file from ${url}..."
        wget -O $tar_file $url
        tar -xzf $tar_file -C "$version_dir"
        rm $tar_file
		binary_file="${tar_file%.tar.gz}"
        mv "$version_dir/build/$binary_file" "$version_dir/bin/dydxprotocold"
        chmod 755 "$version_dir/bin/dydxprotocold"
        echo "Successfully set up $version_dir/bin/dydxprotocold"
    done
    current_version_path="$HOME/cosmovisor/upgrades/$CURRENT_VERSION/bin"
    mkdir -p $current_version_path
    cp /bin/dydxprotocold $current_version_path
}

create_full_nodes() {
	# Create directories for full-nodes to use.
	for i in $(seq 0 $LAST_FULL_NODE_INDEX); do
		FULL_NODE_HOME_DIR="$HOME/chain/.full-node-$i"
		FULL_NODE_CONFIG_DIR="$FULL_NODE_HOME_DIR/config"
		dydxprotocold init "full-node" -o --chain-id=$CHAIN_ID --home "$FULL_NODE_HOME_DIR"
	done

	# Copy the genesis file to the full-node directories.
	for i in $(seq 0 $LAST_FULL_NODE_INDEX); do
		FULL_NODE_HOME_DIR="$HOME/chain/.full-node-$i"
		FULL_NODE_CONFIG_DIR="$FULL_NODE_HOME_DIR/config"

		cp "$HOME/genesis.json" "$FULL_NODE_CONFIG_DIR/genesis.json"
	done

	# Set up CosmosVisor for full-nodes.
	for i in $(seq 0 $LAST_FULL_NODE_INDEX); do
		FULL_NODE_HOME_DIR="$HOME/chain/.full-node-$i"
		# DAEMON_NAME is the name of the binary.
		export DAEMON_NAME=dydxprotocold

		# DAEMON_HOME is the location where the cosmovisor/ directory is kept
		# that contains the genesis binary, the upgrade binaries, and any additional
		# auxiliary files associated with each binary
		export DAEMON_HOME="$HOME/chain/.full-node-$i"

		# Create the folder structure required for using cosmovisor.
		cosmovisor init /bin/dydxprotocold

		# Override cosmovisor's default symlink to point to current verson's binary.
		ln -sf $DAEMON_HOME/cosmovisor/upgrades/$CURRENT_VERSION $DAEMON_HOME/cosmovisor/current

		cp -r "$HOME/cosmovisor" "$FULL_NODE_HOME_DIR/"
	done
}

install_prerequisites
set_cosmovisor_binary_permissions
create_full_nodes