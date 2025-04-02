#!/bin/bash
set -eo pipefail

# This script initializes and configures a single-validator local testnet that runs directly on the host machine
# (without containerization). It sets up one validator and one full node with deterministic keys and accounts.
#
# Prerequisites:
#   - Must be run from the `protocol` directory
#   - Requires `jq` and `dasel` installed (will attempt to install via brew if missing)
#
# Usage:
#   ./testing/testnet-local/local_native.sh <path-to-binary>
# Example:
#   ./testing/testnet-local/local_native.sh dydxprotocold
#
# The script will:
#   1. Create validator and full node configurations in /tmp/chain
#   2. Set up deterministic keys and genesis accounts
#   3. Configure network parameters for local testing

source "./testing/genesis.sh"
BINARY="$1"

CHAIN_ID="localdydxprotocol"

# Define mnemonics for all validators.
MNEMONICS=(
	# alice
	# Consensus Address: dydxvalcons1zf9csp5ygq95cqyxh48w3qkuckmpealrw2ug4d
	"merge panther lobster crazy road hollow amused security before critic about cliff exhibit cause coyote talent happy where lion river tobacco option coconut small"
)

# Define node keys for all full nodes.
FULL_NODE_KEYS=(
	# Node ID: dfa67970296bbecce14daba6cb0da516ed60458a
	"+c9Wyy9G4VJvVmUQ41CogREJPVMDqnBxefcGoika3Qo7U7eJHVIcjPIFuS0HYm224mWMfYgdNlo5KgJ0z1x/0w=="
)

# Define node keys for all validators.
NODE_KEYS=(
	# Node ID: 17e5e45691f0d01449c84fd4ae87279578cdd7ec
	"8EGQBxfGMcRfH0C45UTedEG5Xi3XAcukuInLUqFPpskjp1Ny0c5XvwlKevAwtVvkwoeYYQSe0geQG/cF3GAcUA=="
)

# Define monikers for each validator. These are made up strings and can be anything.
# This also controls in which directory the validator's home will be located. i.e. `/dydxprotocol/chain/.alice`
MONIKERS=(
	"alice"
)

# Define all test accounts for the chain.
TEST_ACCOUNTS=(
	"dydx199tqg4wdlnu4qjlxchpd7seg454937hjrknju4" # alice
)

FAUCET_ACCOUNTS=(
	"dydx1nzuttarf5k2j0nug5yzhr6p74t9avehn9hlh8m" # main faucet
)

# Define dependencies for this script.
# `jq` and `dasel` are used to manipulate json and yaml files respectively.
install_prerequisites() {
	brew install jq
}

TMP_GENTX_DIR="/tmp/gentx"
TMP_CHAIN_DIR="/tmp/chain"
TMP_EXCHANGE_CONFIG_JSON_DIR="/tmp/exchange_config"

cleanup_tmp_dir() {
	if [ -d "/tmp/chain" ]; then
		rm -r "/tmp/chain"
	fi

	if [ -d "$TMP_EXCHANGE_CONFIG_JSON_DIR" ]; then
		rm -r "$TMP_EXCHANGE_CONFIG_JSON_DIR"
	fi
	if [ -d "$TMP_GENTX_DIR" ]; then
		rm -r "$TMP_GENTX_DIR"
	fi
	if [ -d "$TMP_CHAIN_DIR" ]; then
		rm -r "$TMP_CHAIN_DIR"
	fi
}

# Create all validators for the chain including a full-node.
# Initialize their genesis files and home directories.
create_validators() {
	# Create temporary directory for all gentx files.
	mkdir /tmp/gentx

    # Create temporary directory for exchange config jsons.
    echo "Copying exchange config jsons to $TMP_EXCHANGE_CONFIG_JSON_DIR"
    cp -R ./daemons/pricefeed/client/constants/testdata $TMP_EXCHANGE_CONFIG_JSON_DIR

	# Iterate over all validators and set up their home directories, as well as generate `gentx` transaction for each.
	for i in "${!MONIKERS[@]}"; do
		VAL_HOME_DIR="/tmp/chain/.${MONIKERS[$i]}"
		VAL_CONFIG_DIR="$VAL_HOME_DIR/config"

		# Initialize the chain and validator files.
		$BINARY init "${MONIKERS[$i]}" -o --chain-id=$CHAIN_ID --home "$VAL_HOME_DIR"

		# Overwrite the randomly generated `priv_validator_key.json` with a key generated deterministically from the mnemonic.
		$BINARY tendermint gen-priv-key --home "$VAL_HOME_DIR" --mnemonic "${MNEMONICS[$i]}"

		# Note: `dydxprotocold init` non-deterministically creates `node_id.json` for each validator.
		# This is inconvenient for persistent peering during testing in Terraform configuration as the `node_id`
		# would change with every build of this container.
		#
		# For that reason we overwrite the non-deterministically generated one with a deterministic key defined in this file here.
		new_file=$(jq ".priv_key.value = \"${NODE_KEYS[$i]}\"" "$VAL_CONFIG_DIR"/node_key.json)
		cat <<<"$new_file" >"$VAL_CONFIG_DIR"/node_key.json

    	edit_config "$VAL_CONFIG_DIR"
    
		# Using "*" as a subscript results in a single arg: "dydx1... dydx1... dydx1..."
		# Using "@" as a subscript results in separate args: "dydx1..." "dydx1..." "dydx1..."
		# Note: `edit_genesis` must be called before `add-genesis-account`.
		# main
		edit_genesis "$VAL_CONFIG_DIR" "${TEST_ACCOUNTS[*]}" "${FAUCET_ACCOUNTS[*]}" "" "" "$TMP_EXCHANGE_CONFIG_JSON_DIR" "testing/delaymsg_config" "" ""
		# v5.0
		# edit_genesis "$VAL_CONFIG_DIR" "${TEST_ACCOUNTS[*]}" "${FAUCET_ACCOUNTS[*]}" "$TMP_EXCHANGE_CONFIG_JSON_DIR" "testing/delaymsg_config" "" ""

		echo "${MNEMONICS[$i]}" | $BINARY keys add "${MONIKERS[$i]}" --recover --keyring-backend=test --home "$VAL_HOME_DIR"

		for acct in "${TEST_ACCOUNTS[@]}"; do
			$BINARY add-genesis-account "$acct" 100000000000000000$USDC_DENOM,$TESTNET_VALIDATOR_NATIVE_TOKEN_BALANCE$NATIVE_TOKEN --home "$VAL_HOME_DIR"
		done
		for acct in "${FAUCET_ACCOUNTS[@]}"; do
			$BINARY add-genesis-account "$acct" 900000000000000000$USDC_DENOM,$TESTNET_VALIDATOR_NATIVE_TOKEN_BALANCE$NATIVE_TOKEN --home "$VAL_HOME_DIR"
		done

		$BINARY gentx "${MONIKERS[$i]}" $TESTNET_VALIDATOR_SELF_DELEGATE_AMOUNT$NATIVE_TOKEN --moniker="${MONIKERS[$i]}" --keyring-backend=test --chain-id=$CHAIN_ID --home "$VAL_HOME_DIR"

		echo "$BINARY gentx ${MONIKERS[$i]} $TESTNET_VALIDATOR_SELF_DELEGATE_AMOUNT$NATIVE_TOKEN --moniker=${MONIKERS[$i]} --keyring-backend=test --chain-id=$CHAIN_ID --home $VAL_HOME_DIR"

		# Copy the gentx to a shared directory.
		cp -a "$VAL_CONFIG_DIR/gentx/." /tmp/gentx
	done

	# Copy gentxs to the first validator's home directory to build the genesis json file
	FIRST_VAL_HOME_DIR="/tmp/chain/.${MONIKERS[0]}"
	FIRST_VAL_CONFIG_DIR="$FIRST_VAL_HOME_DIR/config"

	rm -rf "$FIRST_VAL_CONFIG_DIR/gentx"
	mkdir "$FIRST_VAL_CONFIG_DIR/gentx"
	cp -r /tmp/gentx "$FIRST_VAL_CONFIG_DIR"

	echo "Collecting gentxs..."
	# Build the final genesis.json file that all validators and the full-nodes will use.
	$BINARY collect-gentxs --home "$FIRST_VAL_HOME_DIR"

	echo "Createing validators..."
	# Copy this genesis file to each of the other validators
	for i in "${!MONIKERS[@]}"; do
		if [[ "$i" == 0 ]]; then
			# Skip first moniker as it already has the correct genesis file.
			continue
		fi

		VAL_HOME_DIR="/tmp/chain/.${MONIKERS[$i]}"
		VAL_CONFIG_DIR="$VAL_HOME_DIR/config"
		rm -rf "$VAL_CONFIG_DIR/genesis.json"
		cp "$FIRST_VAL_CONFIG_DIR/genesis.json" "$VAL_CONFIG_DIR/genesis.json"
	done

	echo "Creating full nodes..."
    # Create directories for full-nodes to use.
    for i in "${!FULL_NODE_KEYS[@]}"; do
        FULL_NODE_HOME_DIR="/tmp/chain/.full-node-$i"
        FULL_NODE_CONFIG_DIR="$FULL_NODE_HOME_DIR/config"
        $BINARY init "full-node" -o --chain-id=$CHAIN_ID --home "$FULL_NODE_HOME_DIR"

        # Note: `dydxprotocold init` non-deterministically creates `node_id.json` for each validator.
        # This is inconvenient for persistent peering during testing in Terraform configuration as the `node_id`
        # would change with every build of this container.
        #
        # For that reason we overwrite the non-deterministically generated one with a deterministic key defined in this file here.
        new_file=$(jq ".priv_key.value = \"${FULL_NODE_KEYS[$i]}\"" "$FULL_NODE_CONFIG_DIR"/node_key.json)
        cat <<<"$new_file" >"$FULL_NODE_CONFIG_DIR"/node_key.json

        edit_config "$FULL_NODE_CONFIG_DIR"

        # Update ports so that full-nodes don't conflict with validators.
        dasel put -t string -f "$FULL_NODE_CONFIG_DIR"/app.toml '.api.address' -v "tcp://0.0.0.0:11317"
		dasel put -t string -f "$FULL_NODE_CONFIG_DIR"/app.toml '.grpc.address' -v "tcp://0.0.0.0:19090"
        dasel put -t string -f "$FULL_NODE_CONFIG_DIR"/client.toml '.node' -v "tcp://localhost:36657"
        dasel put -t string -f "$FULL_NODE_CONFIG_DIR"/config.toml '.rpc.laddr' -v "tcp://0.0.0.0:36657"
        dasel put -t string -f "$FULL_NODE_CONFIG_DIR"/config.toml '.p2p.laddr' -v "tcp://0.0.0.0:36656"



        # Copy genesis file to full-node home directory.
        cp "$FIRST_VAL_CONFIG_DIR/genesis.json" "$FULL_NODE_CONFIG_DIR/genesis.json"
    done
}

# TODO(DEC-1894): remove this function once we migrate off of persistent peers.
# Note: DO NOT add more config modifications in this method. Use `cmd/config.go` to configure
# the default config values.
edit_config() {
	CONFIG_FOLDER=$1

	# Disable pex
	dasel put -t bool -f "$CONFIG_FOLDER"/config.toml '.p2p.pex' -v 'false'

	# Default `timeout_commit` is 999ms. For local testnet, use a larger value to make 
	# block time longer for easier troubleshooting.
	dasel put -t string -f "$CONFIG_FOLDER"/config.toml '.consensus.timeout_commit' -v '5s'
}

cleanup_tmp_dir
# install_prerequisites
create_validators

