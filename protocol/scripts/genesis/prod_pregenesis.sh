#!/bin/bash
set -eo pipefail

# This script creates the pregenesis file for a production network.
#
# The script must be run from the root of the `v4-chain` repo.
#
# example usage:
# $ make build
# $ ./scripts/genesis/prod_pregenesis.sh ./build/dydxprotocold

# Check for missing required arguments
if [ -z "$1" ]; then
  echo "Error: Missing required argument DYDX_BINARY."
  echo "Usage: $0 <DYDX_BINARY> [-s|--SEED_FAUCET_USDC]"
  exit 1
fi

# Capture the required argument
DYDX_BINARY="$1"

source "./testing/genesis.sh"

TMP_CHAIN_DIR="/tmp/prod-chain"
TMP_EXCHANGE_CONFIG_JSON_DIR="/tmp/prod-exchange_config"
BRIDGE_MODACC_BALANCE="1$NINE_ZEROS$EIGHTEEN_ZEROS" # 1e27
BRIDGE_MODACC_ADDR="dydx1zlefkpe3g0vvm9a4h0jf9000lmqutlh9jwjnsv"

# TODO(GENESIS): Update below values before running this script. Sample values are shown.
################## Start of required values to be updated ##################
CHAIN_ID="dydx-sample-1"
# Base denomination of the native token. Usually comes with a prefix "u-", "a-" to indicate unit.
NATIVE_TOKEN="asample"
# Denomination of the native token in whole coins.
NATIVE_TOKEN_WHOLE_COIN="sample" 
# Human readable name of token.
COIN_NAME="Sample Coin Name"
# Market ID in the oracle price list for the rewards token.
REWARDS_TOKEN_MARKET_ID=1
# The numerical chain ID of the Ethereum chain for bridge daemon to query.
ETH_CHAIN_ID=9
# The address of the Ethereum contract for bridge daemon to monitor for logs.
ETH_BRIDGE_ADDRESS="0xsampleaddress" # default value points to a Sepolia contract
# The next event id (the last processed id plus one) of the logs from the Ethereum contract.
BRIDGE_GENESIS_ACKNOWLEDGED_NEXT_ID=99
# The Ethereum block height of the most recently processed bridge event.
BRIDGE_GENESIS_ACKNOWLEDGED_ETH_BLOCK_HEIGHT=99999
# Genesis time of the chain.
GENESIS_TIME="2023-12-31T00:00:00Z"
# Start time of the community vesting schedule.
COMMUNITY_VEST_START_TIME="2001-01-01T00:00:00Z"
# End time of the community vesting schedule.
COMMUNITY_VEST_END_TIME="2050-01-01T00:00:00Z"
# Start time of the rewards vesting schedule.
REWARDS_VEST_START_TIME="2001-01-01T00:00:00Z"
# End time of the rewards vesting schedule.
REWARDS_VEST_END_TIME="2050-01-01T00:00:00Z"

################## End of required values to be updated ##################

cleanup_tmp_dir() {
	if [ -d "$TMP_EXCHANGE_CONFIG_JSON_DIR" ]; then
		rm -r "$TMP_EXCHANGE_CONFIG_JSON_DIR"
	fi
	if [ -d "$TMP_CHAIN_DIR" ]; then
		rm -r "$TMP_CHAIN_DIR"
	fi
}

# Set production default genesis params.
function overwrite_genesis_production() {	
	# Slashing params
	dasel put -t string -f "$GENESIS" '.app_state.slashing.params.signed_blocks_window' -v '8192' # ~3 hr
	dasel put -t string -f "$GENESIS" '.app_state.slashing.params.min_signed_per_window' -v '0.2' # 20%
	dasel put -t string -f "$GENESIS" '.app_state.slashing.params.downtime_jail_duration' -v '7200s'
	dasel put -t string -f "$GENESIS" '.app_state.slashing.params.slash_fraction_double_sign' -v '0.0' # 0%
	dasel put -t string -f "$GENESIS" '.app_state.slashing.params.slash_fraction_downtime' -v '0.0' # 0%

	# Staking params
	dasel put -t string -f "$GENESIS" '.app_state.staking.params.bond_denom' -v "$NATIVE_TOKEN"
	dasel put -t int -f "$GENESIS" '.app_state.staking.params.max_validators' -v '60'
	dasel put -t string -f "$GENESIS" '.app_state.staking.params.min_commission_rate' -v '0.05' # 5%
	dasel put -t string -f "$GENESIS" '.app_state.staking.params.unbonding_time' -v '2592000s' # 30 days
	dasel put -t int -f "$GENESIS" '.app_state.staking.params.max_entries' -v '7'
	dasel put -t int -f "$GENESIS" '.app_state.staking.params.historical_entries' -v '10000'

	# Distribution params
	dasel put -t string -f "$GENESIS" '.app_state.distribution.params.community_tax' -v '0.0' # 0%
	dasel put -t bool -f "$GENESIS" '.app_state.distribution.params.withdraw_addr_enabled' -v 'true'

	# Bank params
	# Initialize bank balance of bridge module account.
	dasel put -t json -f "$GENESIS" ".app_state.bank.balances" -v "[]"
	dasel put -t json -f "$GENESIS" ".app_state.bank.balances.[]" -v "{}"
	dasel put -t string -f "$GENESIS" ".app_state.bank.balances.[0].address" -v "${BRIDGE_MODACC_ADDR}"
	dasel put -t json -f "$GENESIS" ".app_state.bank.balances.[0].coins.[]" -v "{}"
	dasel put -t string -f "$GENESIS" ".app_state.bank.balances.[0].coins.[0].denom" -v "${NATIVE_TOKEN}"
	dasel put -t string -f "$GENESIS" ".app_state.bank.balances.[0].coins.[0].amount" -v "${BRIDGE_MODACC_BALANCE}"
	# Set denom metadata
	set_denom_metadata "$NATIVE_TOKEN" "$NATIVE_TOKEN_WHOLE_COIN" "$COIN_NAME"

	# Governance params
	dasel put -t string -f "$GENESIS" '.app_state.gov.params.min_deposit.[0].amount' -v "10000$EIGHTEEN_ZEROS" # 10k whole coins of native token
	dasel put -t string -f "$GENESIS" '.app_state.gov.params.min_deposit.[0].denom' -v "$NATIVE_TOKEN"
	dasel put -t string -f "$GENESIS" '.app_state.gov.params.max_deposit_period' -v '172800s' # 2 days
	dasel put -t string -f "$GENESIS" '.app_state.gov.params.expedited_voting_period' -v '86400s' # 1 day
	dasel put -t string -f "$GENESIS" '.app_state.gov.params.voting_period' -v '345600s' # 4 days
	dasel put -t string -f "$GENESIS" '.app_state.gov.params.quorum' -v '0.33400' # 33.4%
	dasel put -t string -f "$GENESIS" '.app_state.gov.params.threshold' -v '0.50000' # 50%
	dasel put -t string -f "$GENESIS" '.app_state.gov.params.veto_threshold' -v '0.33400' # 33.4%
	dasel put -t string -f "$GENESIS" '.app_state.gov.params.min_initial_deposit_ratio' -v '0.20000' # 20%
	dasel put -t bool -f "$GENESIS" '.app_state.gov.params.burn_proposal_deposit_prevote' -v 'false' 
	dasel put -t bool -f "$GENESIS" '.app_state.gov.params.burn_vote_quorum' -v 'false' 
	dasel put -t bool -f "$GENESIS" '.app_state.gov.params.burn_vote_veto' -v 'true'

	# Rewards params
	dasel put -t string -f "$GENESIS" '.app_state.rewards.params.denom' -v "$NATIVE_TOKEN"
	dasel put -t int -f "$GENESIS" '.app_state.rewards.params.fee_multiplier_ppm' -v '0'
	dasel put -t int -f "$GENESIS" '.app_state.rewards.params.market_id' -v "$REWARDS_TOKEN_MARKET_ID"

	# Vest params
	# For community treasury
	dasel put -t string -f "$GENESIS" '.app_state.vest.vest_entries.[0].vester_account' -v "community_vester"
	dasel put -t string -f "$GENESIS" '.app_state.vest.vest_entries.[0].treasury_account' -v "community_treasury"
	dasel put -t string -f "$GENESIS" '.app_state.vest.vest_entries.[0].denom' -v "$NATIVE_TOKEN"
	dasel put -t string -f "$GENESIS" '.app_state.vest.vest_entries.[0].start_time' -v "$COMMUNITY_VEST_START_TIME"
	dasel put -t string -f "$GENESIS" '.app_state.vest.vest_entries.[0].end_time' -v "$COMMUNITY_VEST_END_TIME"
	# For rewards treasury
	dasel put -t string -f "$GENESIS" '.app_state.vest.vest_entries.[1].vester_account' -v "rewards_vester"
	dasel put -t string -f "$GENESIS" '.app_state.vest.vest_entries.[1].treasury_account' -v "rewards_treasury"
	dasel put -t string -f "$GENESIS" '.app_state.vest.vest_entries.[1].denom' -v "$NATIVE_TOKEN"
	dasel put -t string -f "$GENESIS" '.app_state.vest.vest_entries.[1].start_time' -v "$REWARDS_VEST_START_TIME"
	dasel put -t string -f "$GENESIS" '.app_state.vest.vest_entries.[1].end_time' -v "$REWARDS_VEST_END_TIME"

	# Delayed message params
	# Schedule a delayed message to swap fee tiers to the standard schedule after ~120 days of blocks.
	dasel put -t int -f "$GENESIS" '.app_state.delaymsg.next_delayed_message_id' -v '1'
	dasel put -t json -f "$GENESIS" '.app_state.delaymsg.delayed_messages' -v "[]"
	dasel put -t json -f "$GENESIS" '.app_state.delaymsg.delayed_messages.[]' -v "{}"
	dasel put -t int -f "$GENESIS" '.app_state.delaymsg.delayed_messages.[0].id' -v '0'
	delaymsg=$(cat "$DELAY_MSG_JSON_DIR/perpetual_fee_params_msg.json" | jq -c '.')
	dasel put -t json -f "$GENESIS" '.app_state.delaymsg.delayed_messages.[0].msg' -v "$delaymsg"
	# Schedule the message to execute in ~120 days (at 1.5s per block)
	dasel put -t int -f "$GENESIS" '.app_state.delaymsg.delayed_messages.[0].block_height' -v '6912000'

	# Bridge module params.
	dasel put -t string -f "$GENESIS" '.app_state.bridge.event_params.denom' -v "$NATIVE_TOKEN"
	dasel put -t int -f "$GENESIS" '.app_state.bridge.event_params.eth_chain_id' -v "$ETH_CHAIN_ID"
	dasel put -t string -f "$GENESIS" '.app_state.bridge.event_params.eth_address' -v "$ETH_BRIDGE_ADDRESS"
	dasel put -t int -f "$GENESIS" '.app_state.bridge.acknowledged_event_info.next_id' -v "$BRIDGE_GENESIS_ACKNOWLEDGED_NEXT_ID"
	dasel put -t int -f "$GENESIS" '.app_state.bridge.acknowledged_event_info.eth_block_height' -v "$BRIDGE_GENESIS_ACKNOWLEDGED_ETH_BLOCK_HEIGHT"

	# Crisis module
	dasel put -t string -f "$GENESIS" '.app_state.crisis.constant_fee.amount' -v "1$EIGHTEEN_ZEROS" # 1 whole coin of native denom
	dasel put -t string -f "$GENESIS" '.app_state.crisis.constant_fee.denom' -v "$NATIVE_TOKEN"

	# Genesis time
	dasel put -t string -f "$GENESIS" '.genesis_time' -v "$GENESIS_TIME"
}

create_pregenesis_file() {
	VAL_HOME_DIR="$TMP_CHAIN_DIR/.dydxprotocol"
	VAL_CONFIG_DIR="$VAL_HOME_DIR/config"
	# This initializes the $VAL_HOME_DIR folder.
	$DYDX_BINARY init "test-moniker" -o --chain-id=$CHAIN_ID --home "$VAL_HOME_DIR"

	# Create temporary directory for exchange config jsons.
	echo "Copying exchange config jsons to $TMP_EXCHANGE_CONFIG_JSON_DIR"
	cp -R ./daemons/pricefeed/client/constants/testdata $TMP_EXCHANGE_CONFIG_JSON_DIR

	echo "Running edit_genesis..."
	edit_genesis "$VAL_CONFIG_DIR" "" "" "" "" "$TMP_EXCHANGE_CONFIG_JSON_DIR" "./testing/delaymsg_config" "STATUS_INITIALIZING" ""
	
	echo "Oerwriting genesis params for production..."
	overwrite_genesis_production
}

sort_genesis_file(){
	jq -S . $VAL_CONFIG_DIR/genesis.json > $VAL_CONFIG_DIR/sorted_genesis.json
}

cleanup_tmp_dir
create_pregenesis_file
sort_genesis_file
echo "Wrote pregenesis file to $VAL_CONFIG_DIR/sorted_genesis.json"
