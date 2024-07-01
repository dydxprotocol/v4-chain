#!/bin/bash
set -exo pipefail

# This file initializes muliple validators for local and CI testing purposes.
# This file should be run as part of `docker-compose.yml`.

source "./genesis.sh"

CHAIN_ID="localdydxprotocol"

# Define mnemonics for all validators.
MNEMONICS=(
	# alice
	# Consensus Address: dydxvalcons1zf9csp5ygq95cqyxh48w3qkuckmpealrw2ug4d
	"merge panther lobster crazy road hollow amused security before critic about cliff exhibit cause coyote talent happy where lion river tobacco option coconut small"

	# bob
	# Consensus Address: dydxvalcons1s7wykslt83kayxuaktep9fw8qxe5n73ucftkh4
	"color habit donor nurse dinosaur stable wonder process post perfect raven gold census inside worth inquiry mammal panic olive toss shadow strong name drum"

	# carl
	# Consensus Address: dydxvalcons1vy0nrh7l4rtezrsakaadz4mngwlpdmhy64h0ls
	"school artefact ghost shop exchange slender letter debris dose window alarm hurt whale tiger find found island what engine ketchup globe obtain glory manage"

	# dave
	# Consensus Address: dydxvalcons1stjspktkshgcsv8sneqk2vs2ws0nw2wr272vtt
	"switch boring kiss cash lizard coconut romance hurry sniff bus accident zone chest height merit elevator furnace eagle fetch quit toward steak mystery nest"
)

# Define node keys for all validators.
NODE_KEYS=(
	# Node ID: 17e5e45691f0d01449c84fd4ae87279578cdd7ec
	"8EGQBxfGMcRfH0C45UTedEG5Xi3XAcukuInLUqFPpskjp1Ny0c5XvwlKevAwtVvkwoeYYQSe0geQG/cF3GAcUA=="

	# Node ID: b69182310be02559483e42c77b7b104352713166
	"3OZf5HenMmeTncJY40VJrNYKIKcXoILU5bkYTLzTJvewowU2/iV2+8wSlGOs9LoKdl0ODfj8UutpMhLn5cORlw=="

	# Node ID: 47539956aaa8e624e0f1d926040e54908ad0eb44
	"tWV4uEya9Xvmm/kwcPTnEQIV1ZHqiqUTN/jLPHhIBq7+g/5AEXInokWUGM0shK9+BPaTPTNlzv7vgE8smsFg4w=="

	# Node ID: 5882428984d83b03d0c907c1f0af343534987052
	"++C3kWgFAs7rUfwAHB7Ffrv43muPg0wTD2/UtSPFFkhtobooIqc78UiotmrT8onuT1jg8/wFPbSjhnKRThTRZg=="
)

# Define monikers for each validator. These are made up strings and can be anything.
# This also controls in which directory the validator's home will be located. i.e. `/dydxprotocol/chain/.alice`
MONIKERS=(
	"alice"
	"bob"
	"carl"
	"dave"
)

# Define all test accounts for the chain.
TEST_ACCOUNTS=(
	"dydx199tqg4wdlnu4qjlxchpd7seg454937hjrknju4" # alice
	"dydx10fx7sy6ywd5senxae9dwytf8jxek3t2gcen2vs" # bob
	"dydx1fjg6zp6vv8t9wvy4lps03r5l4g7tkjw9wvmh70" # carl
	"dydx1wau5mja7j7zdavtfq9lu7ejef05hm6ffenlcsn" # dave
)

FAUCET_ACCOUNTS=(
	"dydx1nzuttarf5k2j0nug5yzhr6p74t9avehn9hlh8m" # main faucet
)

# Addresses of vaults.
# Can use ../scripts/vault/get_vault.go to generate a vault's address.
VAULT_ACCOUNTS=(
	"dydx1c0m5x87llaunl5sgv3q5vd7j5uha26d2q2r2q0" # BTC vault
	"dydx14rplxdyycc6wxmgl8fggppgq4774l70zt6phkw" # ETH vault
)
# Number of each vault, which for CLOB vaults is the ID of the clob pair it quotes on.
VAULT_NUMBERS=(
	0 # BTC clob pair ID
	1 # ETH clob pair ID
)

# Define dependencies for this script.
# `jq` and `dasel` are used to manipulate json and yaml files respectively.
install_prerequisites() {
	apk add dasel jq
}

# Create all validators for the chain including a full-node.
# Initialize their genesis files and home directories.
create_validators() {
	# Create temporary directory for all gentx files.
	mkdir /tmp/gentx

	# Iterate over all validators and set up their home directories, as well as generate `gentx` transaction for each.
	for i in "${!MONIKERS[@]}"; do
		VAL_HOME_DIR="$HOME/chain/.${MONIKERS[$i]}"
		VAL_CONFIG_DIR="$VAL_HOME_DIR/config"

		# Initialize the chain and validator files.
		dydxprotocold init "${MONIKERS[$i]}" -o --chain-id=$CHAIN_ID --home "$VAL_HOME_DIR"

		# Overwrite the randomly generated `priv_validator_key.json` with a key generated deterministically from the mnemonic.
		dydxprotocold tendermint gen-priv-key --home "$VAL_HOME_DIR" --mnemonic "${MNEMONICS[$i]}"

		# Note: `dydxprotocold init` non-deterministically creates `node_id.json` for each validator.
		# This is inconvenient for persistent peering during testing in Terraform configuration as the `node_id`
		# would change with every build of this container.
		#
		# For that reason we overwrite the non-deterministically generated one with a deterministic key defined in this file here.
		new_file=$(jq ".priv_key.value = \"${NODE_KEYS[$i]}\"" "$VAL_CONFIG_DIR"/node_key.json)
		cat <<<"$new_file" >"$VAL_CONFIG_DIR"/node_key.json

		edit_config "$VAL_CONFIG_DIR"
		use_slinky "$VAL_CONFIG_DIR"

		# Using "*" as a subscript results in a single arg: "dydx1... dydx1... dydx1..."
		# Using "@" as a subscript results in separate args: "dydx1..." "dydx1..." "dydx1..."
		# Note: `edit_genesis` must be called before `add-genesis-account`.
		# edit_genesis "$VAL_CONFIG_DIR" "${TEST_ACCOUNTS[*]}" "${FAUCET_ACCOUNTS[*]}" "${VAULT_ACCOUNTS[*]}" "${VAULT_NUMBERS[*]}" "" "" "" ""
		edit_genesis "$VAL_CONFIG_DIR" "" "${FAUCET_ACCOUNTS[*]}" "${VAULT_ACCOUNTS[*]}" "${VAULT_NUMBERS[*]}" "" "" "" ""
		update_genesis_use_test_volatile_market "$VAL_CONFIG_DIR"
		update_genesis_complete_bridge_delay "$VAL_CONFIG_DIR" "30"

		echo "${MNEMONICS[$i]}" | dydxprotocold keys add "${MONIKERS[$i]}" --recover --keyring-backend=test --home "$VAL_HOME_DIR"

		for acct in "${TEST_ACCOUNTS[@]}"; do
			dydxprotocold add-genesis-account "$acct" 100000000000000000$USDC_DENOM,$TESTNET_VALIDATOR_NATIVE_TOKEN_BALANCE$NATIVE_TOKEN --home "$VAL_HOME_DIR"
		done
		for acct in "${FAUCET_ACCOUNTS[@]}"; do
			dydxprotocold add-genesis-account "$acct" 900000000000000000$USDC_DENOM,$TESTNET_VALIDATOR_NATIVE_TOKEN_BALANCE$NATIVE_TOKEN --home "$VAL_HOME_DIR"
		done

		dydxprotocold gentx "${MONIKERS[$i]}" $TESTNET_VALIDATOR_SELF_DELEGATE_AMOUNT$NATIVE_TOKEN --moniker="${MONIKERS[$i]}" --keyring-backend=test --chain-id=$CHAIN_ID --home "$VAL_HOME_DIR"

		# Copy the gentx to a shared directory.
		cp -a "$VAL_CONFIG_DIR/gentx/." /tmp/gentx
	done

	# Copy gentxs to the first validator's home directory to build the genesis json file
	FIRST_VAL_HOME_DIR="$HOME/chain/.${MONIKERS[0]}"
	FIRST_VAL_CONFIG_DIR="$FIRST_VAL_HOME_DIR/config"

	rm -rf "$FIRST_VAL_CONFIG_DIR/gentx"
	mkdir "$FIRST_VAL_CONFIG_DIR/gentx"
	cp -r /tmp/gentx "$FIRST_VAL_CONFIG_DIR"

	# Build the final genesis.json file that all validators and the full-nodes will use.
	dydxprotocold collect-gentxs --home "$FIRST_VAL_HOME_DIR"

	# Copy this genesis file to each of the other validators
	for i in "${!MONIKERS[@]}"; do
		if [[ "$i" == 0 ]]; then
			# Skip first moniker as it already has the correct genesis file.
			continue
		fi

		VAL_HOME_DIR="$HOME/chain/.${MONIKERS[$i]}"
		VAL_CONFIG_DIR="$VAL_HOME_DIR/config"
		rm -rf "$VAL_CONFIG_DIR/genesis.json"
		cp "$FIRST_VAL_CONFIG_DIR/genesis.json" "$VAL_CONFIG_DIR/genesis.json"
	done
}

setup_cosmovisor() {
	for i in "${!MONIKERS[@]}"; do
		VAL_HOME_DIR="$HOME/chain/.${MONIKERS[$i]}"
		export DAEMON_NAME=dydxprotocold
		export DAEMON_HOME="$HOME/chain/.${MONIKERS[$i]}"

		cosmovisor init /bin/dydxprotocold
	done
}

use_slinky() {
  CONFIG_FOLDER=$1
  # Enable slinky daemon
  dasel put -t bool -f "$CONFIG_FOLDER"/app.toml 'oracle.enabled' -v true
	dasel put -t string -f "$VAL_CONFIG_DIR"/app.toml 'oracle.oracle_address' -v 'slinky0:8080'
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

  # Enable Slinky Prometheus metrics
	dasel put -t bool -f "$CONFIG_FOLDER"/app.toml '.oracle.metrics_enabled' -v 'true'
	dasel put -t string -f "$CONFIG_FOLDER"/app.toml '.oracle.prometheus_server_address' -v 'localhost:8001'
}

install_prerequisites
create_validators
setup_cosmovisor
