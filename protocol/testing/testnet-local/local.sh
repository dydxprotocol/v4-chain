#!/bin/bash
set -eo pipefail

# This file initializes muliple validators for local and CI testing purposes.
# This file should be run as part of `docker-compose.yml`.

source "./genesis.sh"

CHAIN_ID="localklyraprotocol"

# Define mnemonics for all validators.
MNEMONICS=(
	# alice
	# Consensus Address: klyravalcons1zf9csp5ygq95cqyxh48w3qkuckmpealrhxq0ye
	"merge panther lobster crazy road hollow amused security before critic about cliff exhibit cause coyote talent happy where lion river tobacco option coconut small"

	# bob
	# Consensus Address: klyravalcons1s7wykslt83kayxuaktep9fw8qxe5n73up9h3xp
	"color habit donor nurse dinosaur stable wonder process post perfect raven gold census inside worth inquiry mammal panic olive toss shadow strong name drum"

	# carl
	# Consensus Address: klyravalcons1vy0nrh7l4rtezrsakaadz4mngwlpdmhyretgwy
	"school artefact ghost shop exchange slender letter debris dose window alarm hurt whale tiger find found island what engine ketchup globe obtain glory manage"

	# dave
	# Consensus Address: klyravalcons1stjspktkshgcsv8sneqk2vs2ws0nw2wrnjkt6l
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
# This also controls in which directory the validator's home will be located. i.e. `/klyraprotocol/chain/.alice`
MONIKERS=(
	"alice"
	"bob"
	"carl"
	"dave"
)

# Define all test accounts for the chain.
TEST_ACCOUNTS=(
	"klyra199tqg4wdlnu4qjlxchpd7seg454937hju8xa57" # alice
	"klyra10fx7sy6ywd5senxae9dwytf8jxek3t2g8gx9ym" # bob
	"klyra1fjg6zp6vv8t9wvy4lps03r5l4g7tkjw93awcky" # carl
	"klyra1wau5mja7j7zdavtfq9lu7ejef05hm6ffxz2hcc" # dave
)

FAUCET_ACCOUNTS=(
	"klyra1nzuttarf5k2j0nug5yzhr6p74t9avehn6x2c0s" # main faucet
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
		klyraprotocold init "${MONIKERS[$i]}" -o --chain-id=$CHAIN_ID --home "$VAL_HOME_DIR"

		# Overwrite the randomly generated `priv_validator_key.json` with a key generated deterministically from the mnemonic.
		klyraprotocold tendermint gen-priv-key --home "$VAL_HOME_DIR" --mnemonic "${MNEMONICS[$i]}"

		# Note: `klyraprotocold init` non-deterministically creates `node_id.json` for each validator.
		# This is inconvenient for persistent peering during testing in Terraform configuration as the `node_id`
		# would change with every build of this container.
		#
		# For that reason we overwrite the non-deterministically generated one with a deterministic key defined in this file here.
		new_file=$(jq ".priv_key.value = \"${NODE_KEYS[$i]}\"" "$VAL_CONFIG_DIR"/node_key.json)
		cat <<<"$new_file" >"$VAL_CONFIG_DIR"/node_key.json

		edit_config "$VAL_CONFIG_DIR"

		# Using "*" as a subscript results in a single arg: "klyra1... klyra1... klyra1..."
		# Using "@" as a subscript results in separate args: "klyra1..." "klyra1..." "klyra1..."
		# Note: `edit_genesis` must be called before `add-genesis-account`.
		edit_genesis "$VAL_CONFIG_DIR" "${TEST_ACCOUNTS[*]}" "${FAUCET_ACCOUNTS[*]}" "" "" "" ""
		update_genesis_use_test_volatile_market "$VAL_CONFIG_DIR"

		echo "${MNEMONICS[$i]}" | klyraprotocold keys add "${MONIKERS[$i]}" --recover --keyring-backend=test --home "$VAL_HOME_DIR"

		for acct in "${TEST_ACCOUNTS[@]}"; do
			klyraprotocold add-genesis-account "$acct" 100000000000000000$TDAI_DENOM,$TESTNET_VALIDATOR_NATIVE_TOKEN_BALANCE$NATIVE_TOKEN --home "$VAL_HOME_DIR"
		done
		for acct in "${FAUCET_ACCOUNTS[@]}"; do
			klyraprotocold add-genesis-account "$acct" 900000000000000000$TDAI_DENOM,$TESTNET_VALIDATOR_NATIVE_TOKEN_BALANCE$NATIVE_TOKEN --home "$VAL_HOME_DIR"
		done

		klyraprotocold gentx "${MONIKERS[$i]}" $TESTNET_VALIDATOR_SELF_DELEGATE_AMOUNT$NATIVE_TOKEN --moniker="${MONIKERS[$i]}" --keyring-backend=test --chain-id=$CHAIN_ID --home "$VAL_HOME_DIR"

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
	klyraprotocold collect-gentxs --home "$FIRST_VAL_HOME_DIR"

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
		export DAEMON_NAME=klyraprotocold
		export DAEMON_HOME="$HOME/chain/.${MONIKERS[$i]}"

		cosmovisor init /bin/klyraprotocold
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

install_prerequisites
create_validators
setup_cosmovisor
