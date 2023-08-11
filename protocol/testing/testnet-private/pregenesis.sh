#!/bin/bash
set -eo pipefail

# This script creates the pregenesis file for the private testnet.
# The pregenesis file includes initial genesis module states for the testnet, as well as two internal validators (dydx-1 and dydx-2) run by the dYdX team.
# External validator accounts and their gentx's will need to be added to get the finalized genesis file.
#
# The script must be run from the root of the repo.
# dasel v1 (https://daseldocs.tomwright.me/v/v1/) is required to run this script. Note that dasel v2, the default version installed with `brew`, uses a different syntax
# (https://daseldocs.tomwright.me/#v1-to-v2-breaking-changes) from currently used in `genesis.sh`.
# To install dasel v1.27.3 (same version used in test docker containers) on macOS:
#   curl -sSLf "$(curl -sSLf https://api.github.com/repos/tomwright/dasel/releases | grep browser_download_url | grep 1.27.3 | grep -v .gz | grep darwin_amd64 | cut -d\" -f 4)" -L -o dasel && chmod +x dasel
#   sudo mv ./dasel /usr/local/bin/dasel

source "./testing/genesis.sh"
CHAIN_ID="dydxprotocol-testnet"
FAUCET_ACCOUNTS=(
	"dydx1nzuttarf5k2j0nug5yzhr6p74t9avehn9hlh8m" # main faucet
	"dydx10du0qegtt73ynv5ctenh565qha27ptzr6dz8c3" # backup #1
	"dydx1axstmx84qtv0avhjwek46v6tcmyc8agu03nafv" # backup #2
)
TMP_GENTX_DIR="/tmp/gentx"
TMP_CHAIN_DIR="/tmp/chain"

# Define monikers for each validator. These are made up strings and can be anything.
# This also controls in which directory the validator's home will be located. i.e. `/tmp/chain/.dydx-1`
MONIKERS=(
	"dydx-1"
	"dydx-2"
)

# Public IPs for each validator. Taken from `Private Test-net Internal Infra Playbook`
IPS=(
	"3.131.247.148"
	"13.112.168.94"
)

# Define mnemonics for internal validators.
MNEMONICS=(
	# dydx-1
	# Consensus Address: dydxvalcons1zf9csp5ygq95cqyxh48w3qkuckmpealrw2ug4d
	"merge panther lobster crazy road hollow amused security before critic about cliff exhibit cause coyote talent happy where lion river tobacco option coconut small"

	# dydx-2
	# Consensus Address: dydxvalcons1s7wykslt83kayxuaktep9fw8qxe5n73ucftkh4
	"color habit donor nurse dinosaur stable wonder process post perfect raven gold census inside worth inquiry mammal panic olive toss shadow strong name drum"
)

# Define node keys for internal validators.
NODE_KEYS=(
	# Node ID: 17e5e45691f0d01449c84fd4ae87279578cdd7ec
	"8EGQBxfGMcRfH0C45UTedEG5Xi3XAcukuInLUqFPpskjp1Ny0c5XvwlKevAwtVvkwoeYYQSe0geQG/cF3GAcUA=="

	# Node ID: b69182310be02559483e42c77b7b104352713166
	"3OZf5HenMmeTncJY40VJrNYKIKcXoILU5bkYTLzTJvewowU2/iV2+8wSlGOs9LoKdl0ODfj8UutpMhLn5cORlw=="
)

VALIDATOR_ACCOUNTS=(
	"dydx199tqg4wdlnu4qjlxchpd7seg454937hjrknju4" # dydx-1
	"dydx10fx7sy6ywd5senxae9dwytf8jxek3t2gcen2vs" # dydx-2
)

cleanup_tmp_dir() {
	if [ -d "$TMP_GENTX_DIR" ]; then
		rm -r "$TMP_GENTX_DIR"
	fi
	if [ -d "$TMP_CHAIN_DIR" ]; then
		rm -r "$TMP_CHAIN_DIR"
	fi
}

create_pregenesis_file() {
	VAL_HOME_DIR="$TMP_CHAIN_DIR/.dydxprotocol"
	VAL_CONFIG_DIR="$VAL_HOME_DIR/config"
	dydxprotocold init "test-moniker" -o --chain-id=$CHAIN_ID --home "$VAL_HOME_DIR"
	# Using "*" as a subscript results in a single arg: "dydx1... dydx1... dydx1..."
	# Using "@" as a subscript results in separate args: "dydx1..." "dydx1..." "dydx1..."
	# Note: `edit_genesis` must be called before `add-genesis-account`.
	edit_genesis "$VAL_CONFIG_DIR" "${TEST_ACCOUNTS[*]}" "${FAUCET_ACCOUNTS[*]}"

	for acct in "${FAUCET_ACCOUNTS[@]}"; do
		dydxprotocold add-genesis-account "$acct" 900000000000000000usdc,900000000000stake --home "$VAL_HOME_DIR"
	done

	# Create temporary directory for all gentx files.
	mkdir "$TMP_GENTX_DIR"

	# Iterate over internal validators and set up their home directories, as well as generate `gentx` transaction for each.
	for i in "${!MONIKERS[@]}"; do
		INDIVIDUAL_VAL_HOME_DIR=""$TMP_CHAIN_DIR"/.${MONIKERS[$i]}"
		INDIVIDUAL_VAL_CONFIG_DIR="$INDIVIDUAL_VAL_HOME_DIR/config"

		# Initialize the chain and validator files.
		dydxprotocold init "${MONIKERS[$i]}" -o --chain-id=$CHAIN_ID --home "$INDIVIDUAL_VAL_HOME_DIR"

		# Overwrite the randomly generated `priv_validator_key.json` with a key generated deterministically from the mnemonic.
		dydxprotocold tendermint gen-priv-key --home "$INDIVIDUAL_VAL_HOME_DIR" --mnemonic "${MNEMONICS[$i]}"

		# Note: `dydxprotocold init` non-deterministically creates `node_id.json` for each validator.
		# This is inconvenient for persistent peering during testing in Terraform configuration as the `node_id`
		# would change with every build of this container.
		#
		# For that reason we overwrite the non-deterministically generated one with a deterministic key defined in this file here.
		new_file=$(jq ".priv_key.value = \"${NODE_KEYS[$i]}\"" "$INDIVIDUAL_VAL_CONFIG_DIR"/node_key.json)
		cat <<<"$new_file" >"$INDIVIDUAL_VAL_CONFIG_DIR"/node_key.json

		echo "${MNEMONICS[$i]}" | dydxprotocold keys add "${MONIKERS[$i]}" --recover --keyring-backend=test --home "$INDIVIDUAL_VAL_HOME_DIR"

		dydxprotocold add-genesis-account "${VALIDATOR_ACCOUNTS[$i]}" 100000000000000000usdc,100000000000stake --home "$INDIVIDUAL_VAL_HOME_DIR"

		dydxprotocold add-genesis-account "${VALIDATOR_ACCOUNTS[$i]}" 100000000000000000usdc,100000000000stake --home "$VAL_HOME_DIR"

		dydxprotocold gentx "${MONIKERS[$i]}" 50000000000stake --moniker="${MONIKERS[$i]}" --keyring-backend=test --chain-id=$CHAIN_ID --home "$INDIVIDUAL_VAL_HOME_DIR" --ip="${IPS[$i]}"

		# Copy the gentx to a shared directory.
		cp -a "$INDIVIDUAL_VAL_CONFIG_DIR/gentx/." "$TMP_GENTX_DIR"
	done

	cp -r "$TMP_GENTX_DIR" "$VAL_CONFIG_DIR"

	dydxprotocold collect-gentxs --home "$VAL_HOME_DIR"
}

cleanup_tmp_dir
create_pregenesis_file
echo "Wrote pregenesis file to $VAL_CONFIG_DIR/genesis.json"
