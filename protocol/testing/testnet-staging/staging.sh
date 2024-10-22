#!/bin/bash
set -eo pipefail

# This file initializes muliple validators and a full-node for remote testnet purposes.

source "./genesis.sh"

CHAIN_ID="dydxprotocol-testnet"

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

	# emily
	# Consensus Address: dydxvalcons1zpt0ck6ttuhjm97apawa54ffgthst34u3peqc0
	"brave way sting spin fog process matrix glimpse volcano recall day lab raccoon hand path pig rent mixture just way blouse alone upon prefer"

	# fiona
	# Consensus Address: dydxvalcons14wwueldgtrdjrmx23wcuwk83keywe5w0pfn94t
	"suffer claw truly wife simple mean still mammal bind cake truly runway attack burden lazy peanut unusual such shock twice appear gloom priority kind"

	# greg
	# Consensus Address: dydxvalcons15yzv3qacs0z2jgm5ecn4ywjkvwc6dl63a2za0p
	"step vital slight present group gallery flower gap copy sweet travel bitter arena reject evidence deal ankle motion dismiss trim armed slab life future"

	# henry
	# Consensus Address: dydxvalcons1pggt0hc2drw0j9456vwpu9wmav67h90h03p3h9
	"piece choice region bike tragic error drive defense air venture bean solve income upset physical sun link actor task runway match gauge brand march"

	# ian
	# Consensus Address: dydxvalcons167ajkznjs3wfa565n9emqmey0c2h69ympf0tmk
	"burst section toss rotate law thumb shoe wire only decide meadow aunt flight humble story mammal radar scene wrist essay taxi leisure excess milk"

	# jeff
	# Consensus Address: dydxvalcons1ehwtcwwh25ftac3khhz0jn7wd9053xfeqdwes4
	"fashion charge estate devote jaguar fun swift always road lend scrap panic matter core defense high gas athlete permit crane assume pact fitness matrix"
)

# Define node keys for all full nodes.
FULL_NODE_KEYS=(
	# Node ID: dfa67970296bbecce14daba6cb0da516ed60458a
	"+c9Wyy9G4VJvVmUQ41CogREJPVMDqnBxefcGoika3Qo7U7eJHVIcjPIFuS0HYm224mWMfYgdNlo5KgJ0z1x/0w=="

	# Node ID: 25dd504d86d82673b9cf94fe78c00714f236c9f8
	"6pJcb5ezfttShtuAsPfOVv5Ua4h3MdZWzHvBnLC3w3f2vdhHiKelbWsVtkzIt6WF475k+5me4n6ptiz99WxkIw=="

	# Node ID: 5b0bdffc54d3aa942ab8abc636bd9cfd0e835709
	"HU8oEKbQU5SgIUosbIPr/WBcW+LW/39eFUo1mEQRNVb3VwrJNbaG7It7hCR+6+Jc9y9IN8QIx7011zV66NDevw=="

	# Node ID: c026a2f137552f26867fb90e7f6a025d44a7781f
	"+dGignS2RfEyWDp39HFxlKZ6h9mgvRrFYEyo8aRDW+PN0XTMU5KrSD6B+1sE7/uAeDsvcth6th+6maHisRMDRg=="
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

	# Node ID: 6c0745b95432ceba57505f7d02f36c750ef18736
	"Xw4SsROea66n4EshBd4mnMAIO0qxnSrtR8B+ehUyav6gacBWq3kcT/w+SjRKuxhRW/AbnI1rkHlDwpT1aqUCgw=="

	# Node ID: 59af8504630944fbc705828e78cd0849e2b952e0
	"IyeYTocVczq4klGBe1l+3TXvKxdu/Hgsz4om0HzQl7cYGBWyoFnoDAs+/MnOtrS/L1zq0SqH2iIOuxl0ltbh0w=="

	# Node ID: 1a4673f1cbc27412632359488c7e3c7d1eb52a56
	"a275Dhs6Pmy1/3TpXDc9/KQvbqLzCBbY92E0maUqtoufO3AN72DqZ+SsRVXEJaFVz49ql1/Nwaoajnekbe/Plg=="

	# Node ID: df6cb9eaf88dd8d981f2f590bcf35a24c56c1b1d
	"PBZscluMMvd6lNZtZO/wuRDuiEEPqrDSYofSZQBPRh+geV0bhL+2ft1W205aMWOyzWEqltUfavdUsU0gssLl+g=="

	# Node ID: e1a73913517dfb41ef817db8e86a05e9260a43e9
	"xq70EO+vOaMTiUjOkBMKB4t5kwPGfA4nU2w6wFUOwMvOAfOzZoN9OBXnbuREJSxhsyW1vjy2k2+PpmvnK1YnQA=="

	# Node ID: b5aed81f57ec4e2c1556355ff54f28983acd4a53
	"cHZoQ+2o+ZTaGfbJPJ66z70GjbkPezomHkC7WRNLW9V3n7We/0du9couCffD/vNvvOvFD1H9uwztl+f2UjcKzA=="
)

# Define monikers for each validator. These are made up strings and can be anything.
# This also controls in which directory the validator's home will be located. i.e. `/dydxprotocol/chain/.alice`
MONIKERS=(
	"alice"
	"bob"
	"carl"
	"dave"
	"emily"
	"fiona"
	"greg"
	"henry"
	"ian"
	"jeff"
)

# Define all test accounts for the chain.
TEST_ACCOUNTS=(
	"dydx199tqg4wdlnu4qjlxchpd7seg454937hjrknju4" # alice
	"dydx10fx7sy6ywd5senxae9dwytf8jxek3t2gcen2vs" # bob
	"dydx1fjg6zp6vv8t9wvy4lps03r5l4g7tkjw9wvmh70" # carl
	"dydx1wau5mja7j7zdavtfq9lu7ejef05hm6ffenlcsn" # dave
	"dydx1966p9acs4mpgj40g3awctfvgrz0mnx8w5vc4sz" # emily
	"dydx18swhz9sgh8ecjaz3cm0v53llw0qzm5se026rsn" # fiona
	"dydx1df84hz7y0dd3mrqcv3vrhw9wdttelul8edqmvp" # greg
	"dydx16h7p7f4dysrgtzptxx2gtpt5d8t834g9dj830z" # henry
	"dydx15u9tppy5e2pdndvlrvafxqhuurj9mnpdstzj6z" # ian
	"dydx168pjt8rkru35239fsqvz7rzgeclakp49zx3aum" # jeff
)

FAUCET_ACCOUNTS=(
	"dydx1nzuttarf5k2j0nug5yzhr6p74t9avehn9hlh8m" # main faucet
	"dydx10du0qegtt73ynv5ctenh565qha27ptzr6dz8c3" # backup #1
	"dydx1axstmx84qtv0avhjwek46v6tcmyc8agu03nafv" # backup #2
)

# Addresses of vaults.
# Can use ../scripts/vault/get_vault.go to generate a vault's address.
VAULT_ACCOUNTS=(
	"dydx1c0m5x87llaunl5sgv3q5vd7j5uha26d2q2r2q0" # BTC vault
	"dydx14rplxdyycc6wxmgl8fggppgq4774l70zt6phkw" # ETH vault
	"dydx190te44zcctdgk0qmqtenve2m00g3r2dn7ntd72" # LINK vault
	"dydx1a83cjn83vqh5ss2vccg6uuaeky7947xldp9r2e" # POL vault
	"dydx1nkz8xcar6sxedw0yva6jzjplw7hfg6pp6e7h0l" # CRV vault
)
# Number of each vault above, which for CLOB vaults is the ID of the clob pair it quotes on.
VAULT_NUMBERS=(
	0 # BTC clob pair ID
	1 # ETH clob pair ID
	2 # LINK clob pair ID
	3 # POL clob pair ID
	4 # CRV clob pair ID
)

# Define dependencies for this script.
# `jq` and `dasel` are used to manipulate json and yaml files respectively.
install_prerequisites() {
	apk add dasel jq
}

# Create all validators for the chain including a full-node.
# Initialize their genesis files and home directories.
create_validators() {
	# Create directories for full-nodes to use.
	for i in "${!FULL_NODE_KEYS[@]}"; do
		FULL_NODE_HOME_DIR="$HOME/chain/.full-node-$i"
		FULL_NODE_CONFIG_DIR="$FULL_NODE_HOME_DIR/config"
		dydxprotocold init "full-node" -o --chain-id=$CHAIN_ID --home "$FULL_NODE_HOME_DIR"

		# Note: `dydxprotocold init` non-deterministically creates `node_id.json` for each validator.
		# This is inconvenient for persistent peering during testing in Terraform configuration as the `node_id`
		# would change with every build of this container.
		#
		# For that reason we overwrite the non-deterministically generated one with a deterministic key defined in this file here.
		new_file=$(jq ".priv_key.value = \"${FULL_NODE_KEYS[$i]}\"" "$FULL_NODE_CONFIG_DIR"/node_key.json)
		cat <<<"$new_file" >"$FULL_NODE_CONFIG_DIR"/node_key.json

		edit_config "$FULL_NODE_CONFIG_DIR"
	done

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

		echo "${MNEMONICS[$i]}" | dydxprotocold keys add "${MONIKERS[$i]}" --recover --keyring-backend=test --home "$VAL_HOME_DIR"

		# Using "*" as a subscript results in a single arg: "dydx1... dydx1... dydx1..."
		# Using "@" as a subscript results in separate args: "dydx1..." "dydx1..." "dydx1..."
		# Note: `edit_genesis` must be called before `add-genesis-account`.
		edit_genesis "$VAL_CONFIG_DIR" "${TEST_ACCOUNTS[*]}" "${FAUCET_ACCOUNTS[*]}" "${VAULT_ACCOUNTS[*]}" "${VAULT_NUMBERS[*]}" "" "" "" ""
		update_genesis_use_test_volatile_market "$VAL_CONFIG_DIR"
		update_genesis_complete_bridge_delay "$VAL_CONFIG_DIR" "600"

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

	# Copy the genesis file to the full-node directories.
	for i in "${!FULL_NODE_KEYS[@]}"; do
		FULL_NODE_HOME_DIR="$HOME/chain/.full-node-$i"
		FULL_NODE_CONFIG_DIR="$FULL_NODE_HOME_DIR/config"

		cp "$FIRST_VAL_CONFIG_DIR/genesis.json" "$FULL_NODE_CONFIG_DIR/genesis.json"
	done
}

setup_cosmovisor() {
	for i in "${!FULL_NODE_KEYS[@]}"; do
		FULL_NODE_HOME_DIR="$HOME/chain/.full-node-$i"
		export DAEMON_NAME=dydxprotocold
		export DAEMON_HOME="$HOME/chain/.full-node-$i"

		cosmovisor init /bin/dydxprotocold
	done

	for i in "${!MONIKERS[@]}"; do
		VAL_HOME_DIR="$HOME/chain/.${MONIKERS[$i]}"
		export DAEMON_NAME=dydxprotocold
		export DAEMON_HOME="$HOME/chain/.${MONIKERS[$i]}"

		cosmovisor init /bin/dydxprotocold
	done
}

# TODO(DEC-1894): remove this function once we migrate off of persistent peers.
# Note: DO NOT add more config modifications in this method. Use `cmd/config.go` to configure
# the default config values.
edit_config() {
	CONFIG_FOLDER=$1

	# Disable pex
	dasel put -t bool -f "$CONFIG_FOLDER"/config.toml '.p2p.pex' -v 'false'

  # Enable Slinky Prometheus metrics
	dasel put -t bool -f "$CONFIG_FOLDER"/app.toml '.oracle.metrics_enabled' -v 'true'
	dasel put -t string -f "$CONFIG_FOLDER"/app.toml '.oracle.prometheus_server_address' -v 'localhost:8001'
}

install_prerequisites
create_validators
setup_cosmovisor
