#!/bin/bash
set -eo pipefail

# Below is the shared genesis configuration for local development, as well as our testnets.
# Any changes to genesis state in those environments belong here.
# If you are making a change to genesis which is _required_ for the chain to function,
# then that change probably belongs in `DefaultGenesis` for the module, and not here.

NINE_ZEROS="000000000"
EIGHTEEN_ZEROS="$NINE_ZEROS$NINE_ZEROS"

# Address of the `subaccounts` module account.
# Obtained from `authtypes.NewModuleAddress(subaccounttypes.ModuleName)`.
SUBACCOUNTS_MODACC_ADDR="dydx1v88c3xv9xyv3eetdx0tvcmq7ung3dywp5upwc6"
REWARDS_VESTER_ACCOUNT_ADDR="dydx1ltyc6y4skclzafvpznpt2qjwmfwgsndp458rmp"
# Address of the `bridge` module account.
# Obtained from `authtypes.NewModuleAddress(bridgetypes.ModuleName)`.
BRIDGE_MODACC_ADDR="dydx1zlefkpe3g0vvm9a4h0jf9000lmqutlh9jwjnsv"
USDC_DENOM="ibc/8E27BA2D5493AF5636760E354E46004562C46AB7EC0CC4C1CA14E9E20E2545B5"
REWARD_TOKEN="adv4tnt"
NATIVE_TOKEN="adv4tnt" # public testnet token
DEFAULT_SUBACCOUNT_QUOTE_BALANCE=100000000000000000
DEFAULT_SUBACCOUNT_QUOTE_BALANCE_FAUCET=900000000000000000
NATIVE_TOKEN_WHOLE_COIN="dv4tnt"
COIN_NAME="dYdX V4 Testnet Token"
# Each testnet validator has 1 million whole coins of native token.
TESTNET_VALIDATOR_NATIVE_TOKEN_BALANCE=1000000$EIGHTEEN_ZEROS # 1e24 or 1 million native tokens.
# Each testnet validator self-delegates 500k whole coins of native token.
TESTNET_VALIDATOR_SELF_DELEGATE_AMOUNT=500000$EIGHTEEN_ZEROS # 5e23 or 500k native tokens.
FAUCET_NATIVE_TOKEN_BALANCE=50000000$EIGHTEEN_ZEROS # 5e25 or 50 million native tokens. 
ETH_CHAIN_ID=11155111 # sepolia
# https://sepolia.etherscan.io/address/0xf75012c350e4ad55be2048bd67ce6e03b20de82d
ETH_BRIDGE_ADDRESS="0xf75012c350e4ad55be2048bd67ce6e03b20de82d"
TOTAL_NATIVE_TOKEN_SUPPLY=1000000000$EIGHTEEN_ZEROS # 1e27
BRIDGE_GENESIS_ACKNOWLEDGED_NEXT_ID=5
BRIDGE_GENESIS_ACKNOWLEDGED_ETH_BLOCK_HEIGHT=4322136

# Use a fix genesis time from the past.
GENESIS_TIME="2023-01-01T00:00:00Z"

function edit_genesis() {
	GENESIS=$1/genesis.json

	# IFS stands for "Internal Field Separator" and it's a special var that determines how bash splits strings.
	# So IFS=' ' specifies that we want to split on spaces.
	# "read" is a built in bash command that reads from stdin.
	# The -a flag specifies that the input should be treated as an array and assign it to the var specified after.
	# The -r flag tells the command to not treat a Backslash as an escape character.
	IFS=' ' read -ra INPUT_TEST_ACCOUNTS <<<"${2}"
	IFS=' ' read -ra INPUT_FAUCET_ACCOUNTS <<<"${3}"

	EXCHANGE_CONFIG_JSON_DIR="$4"
	if [ -z "$EXCHANGE_CONFIG_JSON_DIR" ]; then
		# Default to using exchange_config folder within the current directory.
		EXCHANGE_CONFIG_JSON_DIR="exchange_config"
	fi

	DELAY_MSG_JSON_DIR="$5"
	if [ -z "$DELAY_MSG_JSON_DIR" ]; then
		# Default to using exchange_config folder within the current directory.
		DELAY_MSG_JSON_DIR="delaymsg_config"
	fi

	INITIAL_CLOB_PAIR_STATUS="$6"
		if [ -z "$INITIAL_CLOB_PAIR_STATUS" ]; then
		# Default to initialie clob pairs as active.
		INITIAL_CLOB_PAIR_STATUS='STATUS_ACTIVE'
	fi

	REWARDS_VESTER_ACCOUNT_BALANCE="$7"
	if [ -z "$REWARDS_VESTER_ACCOUNT_BALANCE" ]; then
		# Default to 200 million full coins.
		REWARDS_VESTER_ACCOUNT_BALANCE="200000000$EIGHTEEN_ZEROS"
	fi
	
	# Genesis time
	dasel put -t string -f "$GENESIS" '.genesis_time' -v "$GENESIS_TIME"

	# Consensus params
	dasel put -t string -f "$GENESIS" '.consensus_params.block.max_bytes' -v '4194304'
	dasel put -t string -f "$GENESIS" '.consensus_params.block.max_gas' -v '-1'

	# Update crisis module.
	dasel put -t string -f "$GENESIS" '.app_state.crisis.constant_fee.denom' -v "$NATIVE_TOKEN"

	# Update gov module.
	dasel put -t string -f "$GENESIS" '.app_state.gov.params.min_deposit.[0].denom' -v "$NATIVE_TOKEN"
	# reduced deposit period
	dasel put -t string -f "$GENESIS" '.app_state.gov.params.max_deposit_period' -v '300s'
	# reduced voting period
	dasel put -t string -f "$GENESIS" '.app_state.gov.params.voting_period' -v '300s'
	# set initial deposit ratio to prevent spamming
	dasel put -t string -f "$GENESIS" '.app_state.gov.params.min_initial_deposit_ratio' -v '0.20000' # 20%

	# Update staking module.
	dasel put -t string -f "$GENESIS" '.app_state.staking.params.unbonding_time' -v '1814400s' # 21 days
	dasel put -t string -f "$GENESIS" '.app_state.staking.params.bond_denom' -v "$NATIVE_TOKEN"

	# Update assets module.
	dasel put -t int -f "$GENESIS" '.app_state.assets.assets.[0].id' -v '0'
	dasel put -t string -f "$GENESIS" '.app_state.assets.assets.[0].symbol' -v 'USDC'
	dasel put -t string -f "$GENESIS" '.app_state.assets.assets.[0].denom' -v "$USDC_DENOM"
	dasel put -t string -f "$GENESIS" '.app_state.assets.assets.[0].denom_exponent' -v '-6'
	dasel put -t bool -f "$GENESIS" '.app_state.assets.assets.[0].has_market' -v 'false'
	dasel put -t int -f "$GENESIS" '.app_state.assets.assets.[0].market_id' -v '0'
	dasel put -t int -f "$GENESIS" '.app_state.assets.assets.[0].atomic_resolution' -v '-6'

	# Update bridge module.
	dasel put -t string -f "$GENESIS" '.app_state.bridge.event_params.denom' -v "$NATIVE_TOKEN"
	dasel put -t int -f "$GENESIS" '.app_state.bridge.event_params.eth_chain_id' -v "$ETH_CHAIN_ID"
	dasel put -t string -f "$GENESIS" '.app_state.bridge.event_params.eth_address' -v "$ETH_BRIDGE_ADDRESS"
	dasel put -t int -f "$GENESIS" '.app_state.bridge.acknowledged_event_info.next_id' -v "$BRIDGE_GENESIS_ACKNOWLEDGED_NEXT_ID"
	dasel put -t int -f "$GENESIS" '.app_state.bridge.acknowledged_event_info.eth_block_height' -v "$BRIDGE_GENESIS_ACKNOWLEDGED_ETH_BLOCK_HEIGHT"

	# Update ibc module.
	dasel put -t bool -f "$GENESIS" '.app_state.ibc.client_genesis.create_localhost' -v "false"
	dasel put -t json -f "$GENESIS" '.app_state.ibc.client_genesis.params.allowed_clients' -v "[]"
	dasel put -t string -f "$GENESIS" '.app_state.ibc.client_genesis.params.allowed_clients.[]' -v "07-tendermint"

	# Update perpetuals module.
	# Liquidity Tiers.
	dasel put -t json -f "$GENESIS" '.app_state.perpetuals.liquidity_tiers' -v "[]"
	# Liquidity Tier: Large-Cap
	dasel put -t json -f "$GENESIS" '.app_state.perpetuals.liquidity_tiers.[]' -v "{}"
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.liquidity_tiers.[0].id' -v '0'
	dasel put -t string -f "$GENESIS" '.app_state.perpetuals.liquidity_tiers.[0].name' -v 'Large-Cap'
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.liquidity_tiers.[0].initial_margin_ppm' -v '50000' # 5%
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.liquidity_tiers.[0].maintenance_fraction_ppm' -v '600000' # 60% of IM
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.liquidity_tiers.[0].base_position_notional' -v '1000000000000' # 1_000_000 USDC
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.liquidity_tiers.[0].impact_notional' -v '10000000000' # 10_000 USDC (500 USDC / 5%)

	# Liquidity Tier: Mid-Cap
	dasel put -t json -f "$GENESIS" '.app_state.perpetuals.liquidity_tiers.[]' -v "{}"
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.liquidity_tiers.[1].id' -v '1'
	dasel put -t string -f "$GENESIS" '.app_state.perpetuals.liquidity_tiers.[1].name' -v 'Mid-Cap'
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.liquidity_tiers.[1].initial_margin_ppm' -v '100000' # 10%
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.liquidity_tiers.[1].maintenance_fraction_ppm' -v '500000' # 50% of IM
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.liquidity_tiers.[1].base_position_notional' -v '250000000000' # 250_000 USDC
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.liquidity_tiers.[1].impact_notional' -v '5000000000' # 5_000 USDC (500 USDC / 10%)

	# Liquidity Tier: Long-Tail
	dasel put -t json -f "$GENESIS" '.app_state.perpetuals.liquidity_tiers.[]' -v "{}"
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.liquidity_tiers.[2].id' -v '2'
	dasel put -t string -f "$GENESIS" '.app_state.perpetuals.liquidity_tiers.[2].name' -v 'Long-Tail'
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.liquidity_tiers.[2].initial_margin_ppm' -v '200000' # 20%
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.liquidity_tiers.[2].maintenance_fraction_ppm' -v '500000' # 50% of IM
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.liquidity_tiers.[2].base_position_notional' -v '100000000000' # 100_000 USDC
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.liquidity_tiers.[2].impact_notional' -v '2500000000' # 2_500 USDC (500 USDC / 20%)

	# Liquidity Tier: Safety
	dasel put -t json -f "$GENESIS" '.app_state.perpetuals.liquidity_tiers.[]' -v "{}"
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.liquidity_tiers.[3].id' -v '3'
	dasel put -t string -f "$GENESIS" '.app_state.perpetuals.liquidity_tiers.[3].name' -v 'Safety'
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.liquidity_tiers.[3].initial_margin_ppm' -v '1000000' # 100%
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.liquidity_tiers.[3].maintenance_fraction_ppm' -v '200000' # 20% of IM
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.liquidity_tiers.[3].base_position_notional' -v '1000000000' # 1_000 USDC
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.liquidity_tiers.[3].impact_notional' -v '2500000000' # 2_500 USDC (2_500 USDC / 100%)

	# Params.
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.params.funding_rate_clamp_factor_ppm' -v '6000000' # 600 % (same as 75% on hourly rate)
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.params.premium_vote_clamp_factor_ppm' -v '60000000' # 6000 % (some multiples of funding rate clamp factor)
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.params.min_num_votes_per_sample' -v '15' # half of expected number of votes

	# Perpetuals.
	dasel put -t json -f "$GENESIS" '.app_state.perpetuals.perpetuals' -v "[]"

	# Perpetual: BTC-USD
	dasel put -t json -f "$GENESIS" '.app_state.perpetuals.perpetuals.[]' -v "{}"
	dasel put -t string -f "$GENESIS" '.app_state.perpetuals.perpetuals.[0].params.ticker' -v 'BTC-USD'
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.perpetuals.[0].params.id' -v '0'
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.perpetuals.[0].params.market_id' -v '0'
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.perpetuals.[0].params.atomic_resolution' -v '-10'
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.perpetuals.[0].params.default_funding_ppm' -v '0'
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.perpetuals.[0].params.liquidity_tier' -v '0'

	# Perpetual: ETH-USD
	dasel put -t json -f "$GENESIS" '.app_state.perpetuals.perpetuals.[]' -v "{}"
	dasel put -t string -f "$GENESIS" '.app_state.perpetuals.perpetuals.[1].params.ticker' -v 'ETH-USD'
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.perpetuals.[1].params.id' -v '1'
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.perpetuals.[1].params.market_id' -v '1'
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.perpetuals.[1].params.atomic_resolution' -v '-9'
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.perpetuals.[1].params.default_funding_ppm' -v '0'
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.perpetuals.[1].params.liquidity_tier' -v '0'

	# Perpetual: LINK-USD
	dasel put -t json -f "$GENESIS" '.app_state.perpetuals.perpetuals.[]' -v "{}"
	dasel put -t string -f "$GENESIS" '.app_state.perpetuals.perpetuals.[2].params.ticker' -v 'LINK-USD'
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.perpetuals.[2].params.id' -v '2'
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.perpetuals.[2].params.market_id' -v '2'
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.perpetuals.[2].params.atomic_resolution' -v '-6'
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.perpetuals.[2].params.default_funding_ppm' -v '0'
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.perpetuals.[2].params.liquidity_tier' -v '1'

	# Perpetual: MATIC-USD
	dasel put -t json -f "$GENESIS" '.app_state.perpetuals.perpetuals.[]' -v "{}"
	dasel put -t string -f "$GENESIS" '.app_state.perpetuals.perpetuals.[3].params.ticker' -v 'MATIC-USD'
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.perpetuals.[3].params.id' -v '3'
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.perpetuals.[3].params.market_id' -v '3'
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.perpetuals.[3].params.atomic_resolution' -v '-5'
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.perpetuals.[3].params.default_funding_ppm' -v '0'
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.perpetuals.[3].params.liquidity_tier' -v '1'

	# Perpetual: CRV-USD
	dasel put -t json -f "$GENESIS" '.app_state.perpetuals.perpetuals.[]' -v "{}"
	dasel put -t string -f "$GENESIS" '.app_state.perpetuals.perpetuals.[4].params.ticker' -v 'CRV-USD'
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.perpetuals.[4].params.id' -v '4'
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.perpetuals.[4].params.market_id' -v '4'
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.perpetuals.[4].params.atomic_resolution' -v '-5'
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.perpetuals.[4].params.default_funding_ppm' -v '0'
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.perpetuals.[4].params.liquidity_tier' -v '1'

	# Perpetual: SOL-USD
	dasel put -t json -f "$GENESIS" '.app_state.perpetuals.perpetuals.[]' -v "{}"
	dasel put -t string -f "$GENESIS" '.app_state.perpetuals.perpetuals.[5].params.ticker' -v 'SOL-USD'
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.perpetuals.[5].params.id' -v '5'
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.perpetuals.[5].params.market_id' -v '5'
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.perpetuals.[5].params.atomic_resolution' -v '-7'
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.perpetuals.[5].params.default_funding_ppm' -v '0'
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.perpetuals.[5].params.liquidity_tier' -v '1'

	# Perpetual: ADA-USD
	dasel put -t json -f "$GENESIS" '.app_state.perpetuals.perpetuals.[]' -v "{}"
	dasel put -t string -f "$GENESIS" '.app_state.perpetuals.perpetuals.[6].params.ticker' -v 'ADA-USD'
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.perpetuals.[6].params.id' -v '6'
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.perpetuals.[6].params.market_id' -v '6'
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.perpetuals.[6].params.atomic_resolution' -v '-5'
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.perpetuals.[6].params.default_funding_ppm' -v '0'
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.perpetuals.[6].params.liquidity_tier' -v '1'

	# Perpetual: AVAX-USD
	dasel put -t json -f "$GENESIS" '.app_state.perpetuals.perpetuals.[]' -v "{}"
	dasel put -t string -f "$GENESIS" '.app_state.perpetuals.perpetuals.[7].params.ticker' -v 'AVAX-USD'
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.perpetuals.[7].params.id' -v '7'
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.perpetuals.[7].params.market_id' -v '7'
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.perpetuals.[7].params.atomic_resolution' -v '-7'
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.perpetuals.[7].params.default_funding_ppm' -v '0'
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.perpetuals.[7].params.liquidity_tier' -v '1'

	# Perpetual: FIL-USD
	dasel put -t json -f "$GENESIS" '.app_state.perpetuals.perpetuals.[]' -v "{}"
	dasel put -t string -f "$GENESIS" '.app_state.perpetuals.perpetuals.[8].params.ticker' -v 'FIL-USD'
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.perpetuals.[8].params.id' -v '8'
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.perpetuals.[8].params.market_id' -v '8'
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.perpetuals.[8].params.atomic_resolution' -v '-6'
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.perpetuals.[8].params.default_funding_ppm' -v '0'
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.perpetuals.[8].params.liquidity_tier' -v '1'

	# Perpetual: LTC-USD
	dasel put -t json -f "$GENESIS" '.app_state.perpetuals.perpetuals.[]' -v "{}"
	dasel put -t string -f "$GENESIS" '.app_state.perpetuals.perpetuals.[9].params.ticker' -v 'LTC-USD'
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.perpetuals.[9].params.id' -v '9'
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.perpetuals.[9].params.market_id' -v '9'
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.perpetuals.[9].params.atomic_resolution' -v '-7'
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.perpetuals.[9].params.default_funding_ppm' -v '0'
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.perpetuals.[9].params.liquidity_tier' -v '1'

	# Perpetual: DOGE-USD
	dasel put -t json -f "$GENESIS" '.app_state.perpetuals.perpetuals.[]' -v "{}"
	dasel put -t string -f "$GENESIS" '.app_state.perpetuals.perpetuals.[10].params.ticker' -v 'DOGE-USD'
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.perpetuals.[10].params.id' -v '10'
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.perpetuals.[10].params.market_id' -v '10'
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.perpetuals.[10].params.atomic_resolution' -v '-4'
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.perpetuals.[10].params.default_funding_ppm' -v '0'
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.perpetuals.[10].params.liquidity_tier' -v '1'

	# Perpetual: ATOM-USD
	dasel put -t json -f "$GENESIS" '.app_state.perpetuals.perpetuals.[]' -v "{}"
	dasel put -t string -f "$GENESIS" '.app_state.perpetuals.perpetuals.[11].params.ticker' -v 'ATOM-USD'
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.perpetuals.[11].params.id' -v '11'
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.perpetuals.[11].params.market_id' -v '11'
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.perpetuals.[11].params.atomic_resolution' -v '-6'
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.perpetuals.[11].params.default_funding_ppm' -v '0'
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.perpetuals.[11].params.liquidity_tier' -v '1'

	# Perpetual: DOT-USD
	dasel put -t json -f "$GENESIS" '.app_state.perpetuals.perpetuals.[]' -v "{}"
	dasel put -t string -f "$GENESIS" '.app_state.perpetuals.perpetuals.[12].params.ticker' -v 'DOT-USD'
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.perpetuals.[12].params.id' -v '12'
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.perpetuals.[12].params.market_id' -v '12'
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.perpetuals.[12].params.atomic_resolution' -v '-6'
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.perpetuals.[12].params.default_funding_ppm' -v '0'
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.perpetuals.[12].params.liquidity_tier' -v '1'

	# Perpetual: UNI-USD
	dasel put -t json -f "$GENESIS" '.app_state.perpetuals.perpetuals.[]' -v "{}"
	dasel put -t string -f "$GENESIS" '.app_state.perpetuals.perpetuals.[13].params.ticker' -v 'UNI-USD'
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.perpetuals.[13].params.id' -v '13'
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.perpetuals.[13].params.market_id' -v '13'
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.perpetuals.[13].params.atomic_resolution' -v '-6'
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.perpetuals.[13].params.default_funding_ppm' -v '0'
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.perpetuals.[13].params.liquidity_tier' -v '1'

	# Perpetual: BCH-USD
	dasel put -t json -f "$GENESIS" '.app_state.perpetuals.perpetuals.[]' -v "{}"
	dasel put -t string -f "$GENESIS" '.app_state.perpetuals.perpetuals.[14].params.ticker' -v 'BCH-USD'
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.perpetuals.[14].params.id' -v '14'
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.perpetuals.[14].params.market_id' -v '14'
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.perpetuals.[14].params.atomic_resolution' -v '-8'
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.perpetuals.[14].params.default_funding_ppm' -v '0'
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.perpetuals.[14].params.liquidity_tier' -v '1'

	# Perpetual: TRX-USD
	dasel put -t json -f "$GENESIS" '.app_state.perpetuals.perpetuals.[]' -v "{}"
	dasel put -t string -f "$GENESIS" '.app_state.perpetuals.perpetuals.[15].params.ticker' -v 'TRX-USD'
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.perpetuals.[15].params.id' -v '15'
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.perpetuals.[15].params.market_id' -v '15'
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.perpetuals.[15].params.atomic_resolution' -v '-4'
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.perpetuals.[15].params.default_funding_ppm' -v '0'
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.perpetuals.[15].params.liquidity_tier' -v '1'

	# Perpetual: NEAR-USD
	dasel put -t json -f "$GENESIS" '.app_state.perpetuals.perpetuals.[]' -v "{}"
	dasel put -t string -f "$GENESIS" '.app_state.perpetuals.perpetuals.[16].params.ticker' -v 'NEAR-USD'
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.perpetuals.[16].params.id' -v '16'
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.perpetuals.[16].params.market_id' -v '16'
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.perpetuals.[16].params.atomic_resolution' -v '-6'
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.perpetuals.[16].params.default_funding_ppm' -v '0'
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.perpetuals.[16].params.liquidity_tier' -v '1'

	# Perpetual: MKR-USD
	dasel put -t json -f "$GENESIS" '.app_state.perpetuals.perpetuals.[]' -v "{}"
	dasel put -t string -f "$GENESIS" '.app_state.perpetuals.perpetuals.[17].params.ticker' -v 'MKR-USD'
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.perpetuals.[17].params.id' -v '17'
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.perpetuals.[17].params.market_id' -v '17'
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.perpetuals.[17].params.atomic_resolution' -v '-9'
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.perpetuals.[17].params.default_funding_ppm' -v '0'
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.perpetuals.[17].params.liquidity_tier' -v '2'

	# Perpetual: XLM-USD
	dasel put -t json -f "$GENESIS" '.app_state.perpetuals.perpetuals.[]' -v "{}"
	dasel put -t string -f "$GENESIS" '.app_state.perpetuals.perpetuals.[18].params.ticker' -v 'XLM-USD'
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.perpetuals.[18].params.id' -v '18'
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.perpetuals.[18].params.market_id' -v '18'
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.perpetuals.[18].params.atomic_resolution' -v '-5'
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.perpetuals.[18].params.default_funding_ppm' -v '0'
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.perpetuals.[18].params.liquidity_tier' -v '1'

	# Perpetual: ETC-USD
	dasel put -t json -f "$GENESIS" '.app_state.perpetuals.perpetuals.[]' -v "{}"
	dasel put -t string -f "$GENESIS" '.app_state.perpetuals.perpetuals.[19].params.ticker' -v 'ETC-USD'
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.perpetuals.[19].params.id' -v '19'
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.perpetuals.[19].params.market_id' -v '19'
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.perpetuals.[19].params.atomic_resolution' -v '-7'
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.perpetuals.[19].params.default_funding_ppm' -v '0'
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.perpetuals.[19].params.liquidity_tier' -v '1'

	# Perpetual: COMP-USD
	dasel put -t json -f "$GENESIS" '.app_state.perpetuals.perpetuals.[]' -v "{}"
	dasel put -t string -f "$GENESIS" '.app_state.perpetuals.perpetuals.[20].params.ticker' -v 'COMP-USD'
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.perpetuals.[20].params.id' -v '20'
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.perpetuals.[20].params.market_id' -v '20'
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.perpetuals.[20].params.atomic_resolution' -v '-7'
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.perpetuals.[20].params.default_funding_ppm' -v '0'
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.perpetuals.[20].params.liquidity_tier' -v '2'

	# Perpetual: WLD-USD
	dasel put -t json -f "$GENESIS" '.app_state.perpetuals.perpetuals.[]' -v "{}"
	dasel put -t string -f "$GENESIS" '.app_state.perpetuals.perpetuals.[21].params.ticker' -v 'WLD-USD'
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.perpetuals.[21].params.id' -v '21'
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.perpetuals.[21].params.market_id' -v '21'
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.perpetuals.[21].params.atomic_resolution' -v '-6'
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.perpetuals.[21].params.default_funding_ppm' -v '0'
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.perpetuals.[21].params.liquidity_tier' -v '1'

	# Perpetual: APE-USD
	dasel put -t json -f "$GENESIS" '.app_state.perpetuals.perpetuals.[]' -v "{}"
	dasel put -t string -f "$GENESIS" '.app_state.perpetuals.perpetuals.[22].params.ticker' -v 'APE-USD'
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.perpetuals.[22].params.id' -v '22'
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.perpetuals.[22].params.market_id' -v '22'
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.perpetuals.[22].params.atomic_resolution' -v '-6'
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.perpetuals.[22].params.default_funding_ppm' -v '0'
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.perpetuals.[22].params.liquidity_tier' -v '2'

	# Perpetual: APT-USD
	dasel put -t json -f "$GENESIS" '.app_state.perpetuals.perpetuals.[]' -v "{}"
	dasel put -t string -f "$GENESIS" '.app_state.perpetuals.perpetuals.[23].params.ticker' -v 'APT-USD'
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.perpetuals.[23].params.id' -v '23'
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.perpetuals.[23].params.market_id' -v '23'
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.perpetuals.[23].params.atomic_resolution' -v '-6'
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.perpetuals.[23].params.default_funding_ppm' -v '0'
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.perpetuals.[23].params.liquidity_tier' -v '1'

	# Perpetual: ARB-USD
	dasel put -t json -f "$GENESIS" '.app_state.perpetuals.perpetuals.[]' -v "{}"
	dasel put -t string -f "$GENESIS" '.app_state.perpetuals.perpetuals.[24].params.ticker' -v 'ARB-USD'
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.perpetuals.[24].params.id' -v '24'
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.perpetuals.[24].params.market_id' -v '24'
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.perpetuals.[24].params.atomic_resolution' -v '-6'
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.perpetuals.[24].params.default_funding_ppm' -v '0'
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.perpetuals.[24].params.liquidity_tier' -v '1'

	# Perpetual: BLUR-USD
	dasel put -t json -f "$GENESIS" '.app_state.perpetuals.perpetuals.[]' -v "{}"
	dasel put -t string -f "$GENESIS" '.app_state.perpetuals.perpetuals.[25].params.ticker' -v 'BLUR-USD'
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.perpetuals.[25].params.id' -v '25'
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.perpetuals.[25].params.market_id' -v '25'
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.perpetuals.[25].params.atomic_resolution' -v '-5'
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.perpetuals.[25].params.default_funding_ppm' -v '0'
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.perpetuals.[25].params.liquidity_tier' -v '2'

	# Perpetual: LDO-USD
	dasel put -t json -f "$GENESIS" '.app_state.perpetuals.perpetuals.[]' -v "{}"
	dasel put -t string -f "$GENESIS" '.app_state.perpetuals.perpetuals.[26].params.ticker' -v 'LDO-USD'
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.perpetuals.[26].params.id' -v '26'
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.perpetuals.[26].params.market_id' -v '26'
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.perpetuals.[26].params.atomic_resolution' -v '-6'
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.perpetuals.[26].params.default_funding_ppm' -v '0'
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.perpetuals.[26].params.liquidity_tier' -v '2'

	# Perpetual: OP-USD
	dasel put -t json -f "$GENESIS" '.app_state.perpetuals.perpetuals.[]' -v "{}"
	dasel put -t string -f "$GENESIS" '.app_state.perpetuals.perpetuals.[27].params.ticker' -v 'OP-USD'
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.perpetuals.[27].params.id' -v '27'
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.perpetuals.[27].params.market_id' -v '27'
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.perpetuals.[27].params.atomic_resolution' -v '-6'
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.perpetuals.[27].params.default_funding_ppm' -v '0'
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.perpetuals.[27].params.liquidity_tier' -v '1'

	# Perpetual: PEPE-USD
	dasel put -t json -f "$GENESIS" '.app_state.perpetuals.perpetuals.[]' -v "{}"
	dasel put -t string -f "$GENESIS" '.app_state.perpetuals.perpetuals.[28].params.ticker' -v 'PEPE-USD'
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.perpetuals.[28].params.id' -v '28'
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.perpetuals.[28].params.market_id' -v '28'
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.perpetuals.[28].params.atomic_resolution' -v '1'
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.perpetuals.[28].params.default_funding_ppm' -v '0'
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.perpetuals.[28].params.liquidity_tier' -v '1'

	# Perpetual: SEI-USD
	dasel put -t json -f "$GENESIS" '.app_state.perpetuals.perpetuals.[]' -v "{}"
	dasel put -t string -f "$GENESIS" '.app_state.perpetuals.perpetuals.[29].params.ticker' -v 'SEI-USD'
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.perpetuals.[29].params.id' -v '29'
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.perpetuals.[29].params.market_id' -v '29'
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.perpetuals.[29].params.atomic_resolution' -v '-5'
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.perpetuals.[29].params.default_funding_ppm' -v '0'
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.perpetuals.[29].params.liquidity_tier' -v '2'

	# Perpetual: SHIB-USD
	dasel put -t json -f "$GENESIS" '.app_state.perpetuals.perpetuals.[]' -v "{}"
	dasel put -t string -f "$GENESIS" '.app_state.perpetuals.perpetuals.[30].params.ticker' -v 'SHIB-USD'
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.perpetuals.[30].params.id' -v '30'
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.perpetuals.[30].params.market_id' -v '30'
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.perpetuals.[30].params.atomic_resolution' -v '0'
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.perpetuals.[30].params.default_funding_ppm' -v '0'
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.perpetuals.[30].params.liquidity_tier' -v '1'

	# Perpetual: SUI-USD
	dasel put -t json -f "$GENESIS" '.app_state.perpetuals.perpetuals.[]' -v "{}"
	dasel put -t string -f "$GENESIS" '.app_state.perpetuals.perpetuals.[31].params.ticker' -v 'SUI-USD'
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.perpetuals.[31].params.id' -v '31'
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.perpetuals.[31].params.market_id' -v '31'
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.perpetuals.[31].params.atomic_resolution' -v '-5'
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.perpetuals.[31].params.default_funding_ppm' -v '0'
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.perpetuals.[31].params.liquidity_tier' -v '1'

	# Perpetual: XRP-USD
	dasel put -t json -f "$GENESIS" '.app_state.perpetuals.perpetuals.[]' -v "{}"
	dasel put -t string -f "$GENESIS" '.app_state.perpetuals.perpetuals.[32].params.ticker' -v 'XRP-USD'
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.perpetuals.[32].params.id' -v '32'
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.perpetuals.[32].params.market_id' -v '32'
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.perpetuals.[32].params.atomic_resolution' -v '-5'
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.perpetuals.[32].params.default_funding_ppm' -v '0'
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.perpetuals.[32].params.liquidity_tier' -v '1'

	# Update prices module.
	# Market: BTC-USD
	dasel put -t json -f "$GENESIS" '.app_state.prices.market_params' -v "[]"
	dasel put -t json -f "$GENESIS" '.app_state.prices.market_prices' -v "[]"

	dasel put -t json -f "$GENESIS" '.app_state.prices.market_params.[]' -v "{}"
	dasel put -t string -f "$GENESIS" '.app_state.prices.market_params.[0].pair' -v 'BTC-USD'
	dasel put -t int -f "$GENESIS" '.app_state.prices.market_params.[0].id' -v '0'
	dasel put -t int -f "$GENESIS" '.app_state.prices.market_params.[0].exponent' -v '-5'
	dasel put -t int -f "$GENESIS" '.app_state.prices.market_params.[0].min_exchanges' -v '3'
	dasel put -t int -f "$GENESIS" '.app_state.prices.market_params.[0].min_price_change_ppm' -v '1000' # 0.1%
	dasel put -t json -f "$GENESIS" '.app_state.prices.market_prices.[]' -v "{}"
	dasel put -t int -f "$GENESIS" '.app_state.prices.market_prices.[0].id' -v '0'
	dasel put -t int -f "$GENESIS" '.app_state.prices.market_prices.[0].exponent' -v '-5'
	dasel put -t int -f "$GENESIS" '.app_state.prices.market_prices.[0].price' -v '2868819524'          # $28,688 = 1 BTC.
	# BTC Exchange Config
	btc_exchange_config_json=$(cat "$EXCHANGE_CONFIG_JSON_DIR/btc_exchange_config.json" | jq -c '.')
	dasel put -t string -f "$GENESIS" '.app_state.prices.market_params.[0].exchange_config_json' -v "$btc_exchange_config_json"

	# Market: ETH-USD
	dasel put -t json -f "$GENESIS" '.app_state.prices.market_params.[]' -v "{}"
	dasel put -t string -f "$GENESIS" '.app_state.prices.market_params.[1].pair' -v 'ETH-USD'
	dasel put -t int -f "$GENESIS" '.app_state.prices.market_params.[1].id' -v '1'
	dasel put -t int -f "$GENESIS" '.app_state.prices.market_params.[1].exponent' -v '-6'
	dasel put -t int -f "$GENESIS" '.app_state.prices.market_params.[1].min_exchanges' -v '3'
	dasel put -t int -f "$GENESIS" '.app_state.prices.market_params.[1].min_price_change_ppm' -v '1000' # 0.1%
	dasel put -t json -f "$GENESIS" '.app_state.prices.market_prices.[]' -v "{}"
	dasel put -t int -f "$GENESIS" '.app_state.prices.market_prices.[1].id' -v '1'
	dasel put -t int -f "$GENESIS" '.app_state.prices.market_prices.[1].exponent' -v '-6'
	dasel put -t int -f "$GENESIS" '.app_state.prices.market_prices.[1].price' -v '1811985252'          # $1,812 = 1 ETH.
	# ETH Exchange Config
	eth_exchange_config_json=$(cat "$EXCHANGE_CONFIG_JSON_DIR/eth_exchange_config.json" | jq -c '.')
	dasel put -t string -f "$GENESIS" '.app_state.prices.market_params.[1].exchange_config_json' -v "$eth_exchange_config_json"

	# Market: LINK-USD
	dasel put -t json -f "$GENESIS" '.app_state.prices.market_params.[]' -v "{}"
	dasel put -t string -f "$GENESIS" '.app_state.prices.market_params.[2].pair' -v 'LINK-USD'
	dasel put -t int -f "$GENESIS" '.app_state.prices.market_params.[2].id' -v '2'
	dasel put -t int -f "$GENESIS" '.app_state.prices.market_params.[2].exponent' -v '-9'
	dasel put -t int -f "$GENESIS" '.app_state.prices.market_params.[2].min_exchanges' -v '3'
	dasel put -t int -f "$GENESIS" '.app_state.prices.market_params.[2].min_price_change_ppm' -v '2500' # 0.25%
	dasel put -t json -f "$GENESIS" '.app_state.prices.market_prices.[]' -v "{}"
	dasel put -t int -f "$GENESIS" '.app_state.prices.market_prices.[2].id' -v '2'
	dasel put -t int -f "$GENESIS" '.app_state.prices.market_prices.[2].exponent' -v '-9'
	dasel put -t int -f "$GENESIS" '.app_state.prices.market_prices.[2].price' -v '7204646989'          # $7.205 = 1 LINK.
	# LINK Exchange Config
	link_exchange_config_json=$(cat "$EXCHANGE_CONFIG_JSON_DIR/link_exchange_config.json" | jq -c '.')
	dasel put -t string -f "$GENESIS" '.app_state.prices.market_params.[2].exchange_config_json' -v "$link_exchange_config_json"

	# Market: MATIC-USD
	dasel put -t json -f "$GENESIS" '.app_state.prices.market_params.[]' -v "{}"
	dasel put -t string -f "$GENESIS" '.app_state.prices.market_params.[3].pair' -v 'MATIC-USD'
	dasel put -t int -f "$GENESIS" '.app_state.prices.market_params.[3].id' -v '3'
	dasel put -t int -f "$GENESIS" '.app_state.prices.market_params.[3].exponent' -v '-10'
	dasel put -t int -f "$GENESIS" '.app_state.prices.market_params.[3].min_exchanges' -v '3'
	dasel put -t int -f "$GENESIS" '.app_state.prices.market_params.[3].min_price_change_ppm' -v '2500' # 0.25%
	dasel put -t json -f "$GENESIS" '.app_state.prices.market_prices.[]' -v "{}"
	dasel put -t int -f "$GENESIS" '.app_state.prices.market_prices.[3].id' -v '3'
	dasel put -t int -f "$GENESIS" '.app_state.prices.market_prices.[3].exponent' -v '-10'
	dasel put -t int -f "$GENESIS" '.app_state.prices.market_prices.[3].price' -v '6665746387'          # $0.6666 = 1 MATIC.
	# MATIC Exchange Config
	matic_exchange_config_json=$(cat "$EXCHANGE_CONFIG_JSON_DIR/matic_exchange_config.json" | jq -c '.')
	dasel put -t string -f "$GENESIS" '.app_state.prices.market_params.[3].exchange_config_json' -v "$matic_exchange_config_json"

	# Market: CRV-USD
	dasel put -t json -f "$GENESIS" '.app_state.prices.market_params.[]' -v "{}"
	dasel put -t string -f "$GENESIS" '.app_state.prices.market_params.[4].pair' -v 'CRV-USD'
	dasel put -t int -f "$GENESIS" '.app_state.prices.market_params.[4].id' -v '4'
	dasel put -t int -f "$GENESIS" '.app_state.prices.market_params.[4].exponent' -v '-10'
	dasel put -t int -f "$GENESIS" '.app_state.prices.market_params.[4].min_exchanges' -v '3'
	dasel put -t int -f "$GENESIS" '.app_state.prices.market_params.[4].min_price_change_ppm' -v '2500' # 0.25%
	dasel put -t json -f "$GENESIS" '.app_state.prices.market_prices.[]' -v "{}"
	dasel put -t int -f "$GENESIS" '.app_state.prices.market_prices.[4].id' -v '4'
	dasel put -t int -f "$GENESIS" '.app_state.prices.market_prices.[4].exponent' -v '-10'
	dasel put -t int -f "$GENESIS" '.app_state.prices.market_prices.[4].price' -v '6029316660'          # $0.6029 = 1 CRV.
	# CRV Exchange Config
	crv_exchange_config_json=$(cat "$EXCHANGE_CONFIG_JSON_DIR/crv_exchange_config.json" | jq -c '.')
	dasel put -t string -f "$GENESIS" '.app_state.prices.market_params.[4].exchange_config_json' -v "$crv_exchange_config_json"

	# Market: SOL-USD
	dasel put -t json -f "$GENESIS" '.app_state.prices.market_params.[]' -v "{}"
	dasel put -t string -f "$GENESIS" '.app_state.prices.market_params.[5].pair' -v 'SOL-USD'
	dasel put -t int -f "$GENESIS" '.app_state.prices.market_params.[5].id' -v '5'
	dasel put -t int -f "$GENESIS" '.app_state.prices.market_params.[5].exponent' -v '-8'
	dasel put -t int -f "$GENESIS" '.app_state.prices.market_params.[5].min_exchanges' -v '3'
	dasel put -t int -f "$GENESIS" '.app_state.prices.market_params.[5].min_price_change_ppm' -v '2500' # 0.25%
	dasel put -t json -f "$GENESIS" '.app_state.prices.market_prices.[]' -v "{}"
	dasel put -t int -f "$GENESIS" '.app_state.prices.market_prices.[5].id' -v '5'
	dasel put -t int -f "$GENESIS" '.app_state.prices.market_prices.[5].exponent' -v '-8'
	dasel put -t int -f "$GENESIS" '.app_state.prices.market_prices.[5].price' -v '2350695125'          # $23.51 = 1 SOL.
	# SOL Exchange Config
	sol_exchange_config_json=$(cat "$EXCHANGE_CONFIG_JSON_DIR/sol_exchange_config.json" | jq -c '.')
	dasel put -t string -f "$GENESIS" '.app_state.prices.market_params.[5].exchange_config_json' -v "$sol_exchange_config_json"

	# Market: ADA-USD
	dasel put -t json -f "$GENESIS" '.app_state.prices.market_params.[]' -v "{}"
	dasel put -t string -f "$GENESIS" '.app_state.prices.market_params.[6].pair' -v 'ADA-USD'
	dasel put -t int -f "$GENESIS" '.app_state.prices.market_params.[6].id' -v '6'
	dasel put -t int -f "$GENESIS" '.app_state.prices.market_params.[6].exponent' -v '-10'
	dasel put -t int -f "$GENESIS" '.app_state.prices.market_params.[6].min_exchanges' -v '3'
	dasel put -t int -f "$GENESIS" '.app_state.prices.market_params.[6].min_price_change_ppm' -v '2500' # 0.25%
	dasel put -t json -f "$GENESIS" '.app_state.prices.market_prices.[]' -v "{}"
	dasel put -t int -f "$GENESIS" '.app_state.prices.market_prices.[6].id' -v '6'
	dasel put -t int -f "$GENESIS" '.app_state.prices.market_prices.[6].exponent' -v '-10'
	dasel put -t int -f "$GENESIS" '.app_state.prices.market_prices.[6].price' -v '2918831290'          # $0.2919 = 1 ADA.
	# ADA Exchange Config
	ada_exchange_config_json=$(cat "$EXCHANGE_CONFIG_JSON_DIR/ada_exchange_config.json" | jq -c '.')
	dasel put -t string -f "$GENESIS" '.app_state.prices.market_params.[6].exchange_config_json' -v "$ada_exchange_config_json"

	# Market: AVAX-USD
	dasel put -t json -f "$GENESIS" '.app_state.prices.market_params.[]' -v "{}"
	dasel put -t string -f "$GENESIS" '.app_state.prices.market_params.[7].pair' -v 'AVAX-USD'
	dasel put -t int -f "$GENESIS" '.app_state.prices.market_params.[7].id' -v '7'
	dasel put -t int -f "$GENESIS" '.app_state.prices.market_params.[7].exponent' -v '-8'
	dasel put -t int -f "$GENESIS" '.app_state.prices.market_params.[7].min_exchanges' -v '3'
	dasel put -t int -f "$GENESIS" '.app_state.prices.market_params.[7].min_price_change_ppm' -v '2500' # 0.25%
	dasel put -t json -f "$GENESIS" '.app_state.prices.market_prices.[]' -v "{}"
	dasel put -t int -f "$GENESIS" '.app_state.prices.market_prices.[7].id' -v '7'
	dasel put -t int -f "$GENESIS" '.app_state.prices.market_prices.[7].exponent' -v '-8'
	dasel put -t int -f "$GENESIS" '.app_state.prices.market_prices.[7].price' -v '1223293720'          # $12.23 = 1 AVAX.
	# AVAX Exchange Config
	avax_exchange_config_json=$(cat "$EXCHANGE_CONFIG_JSON_DIR/avax_exchange_config.json" | jq -c '.')
	dasel put -t string -f "$GENESIS" '.app_state.prices.market_params.[7].exchange_config_json' -v "$avax_exchange_config_json"

	# Market: FIL-USD
	dasel put -t json -f "$GENESIS" '.app_state.prices.market_params.[]' -v "{}"
	dasel put -t string -f "$GENESIS" '.app_state.prices.market_params.[8].pair' -v 'FIL-USD'
	dasel put -t int -f "$GENESIS" '.app_state.prices.market_params.[8].id' -v '8'
	dasel put -t int -f "$GENESIS" '.app_state.prices.market_params.[8].exponent' -v '-9'
	dasel put -t int -f "$GENESIS" '.app_state.prices.market_params.[8].min_exchanges' -v '3'
	dasel put -t int -f "$GENESIS" '.app_state.prices.market_params.[8].min_price_change_ppm' -v '2500' # 0.25%
	dasel put -t json -f "$GENESIS" '.app_state.prices.market_prices.[]' -v "{}"
	dasel put -t int -f "$GENESIS" '.app_state.prices.market_prices.[8].id' -v '8'
	dasel put -t int -f "$GENESIS" '.app_state.prices.market_prices.[8].exponent' -v '-9'
	dasel put -t int -f "$GENESIS" '.app_state.prices.market_prices.[8].price' -v '4050336602'          # $4.050 = 1 FIL.
	# FIL Exchange Config
	fil_exchange_config_json=$(cat "$EXCHANGE_CONFIG_JSON_DIR/fil_exchange_config.json" | jq -c '.')
	dasel put -t string -f "$GENESIS" '.app_state.prices.market_params.[8].exchange_config_json' -v "$fil_exchange_config_json"

	# Market: LTC-USD
	dasel put -t json -f "$GENESIS" '.app_state.prices.market_params.[]' -v "{}"
	dasel put -t string -f "$GENESIS" '.app_state.prices.market_params.[9].pair' -v 'LTC-USD'
	dasel put -t int -f "$GENESIS" '.app_state.prices.market_params.[9].id' -v '9'
	dasel put -t int -f "$GENESIS" '.app_state.prices.market_params.[9].exponent' -v '-8'
	dasel put -t int -f "$GENESIS" '.app_state.prices.market_params.[9].min_exchanges' -v '3'
	dasel put -t int -f "$GENESIS" '.app_state.prices.market_params.[9].min_price_change_ppm' -v '2500' # 0.25%
	dasel put -t json -f "$GENESIS" '.app_state.prices.market_prices.[]' -v "{}"
	dasel put -t int -f "$GENESIS" '.app_state.prices.market_prices.[9].id' -v '9'
	dasel put -t int -f "$GENESIS" '.app_state.prices.market_prices.[9].exponent' -v '-8'
	dasel put -t int -f "$GENESIS" '.app_state.prices.market_prices.[9].price' -v '8193604950'          # $81.93 = 1 LTC.
	# LTC Exchange Config
	ltc_exchange_config_json=$(cat "$EXCHANGE_CONFIG_JSON_DIR/ltc_exchange_config.json" | jq -c '.')
	dasel put -t string -f "$GENESIS" '.app_state.prices.market_params.[9].exchange_config_json' -v "$ltc_exchange_config_json"

	# Market: DOGE-USD
	dasel put -t json -f "$GENESIS" '.app_state.prices.market_params.[]' -v "{}"
	dasel put -t string -f "$GENESIS" '.app_state.prices.market_params.[10].pair' -v 'DOGE-USD'
	dasel put -t int -f "$GENESIS" '.app_state.prices.market_params.[10].id' -v '10'
	dasel put -t int -f "$GENESIS" '.app_state.prices.market_params.[10].exponent' -v '-11'
	dasel put -t int -f "$GENESIS" '.app_state.prices.market_params.[10].min_exchanges' -v '3'
	dasel put -t int -f "$GENESIS" '.app_state.prices.market_params.[10].min_price_change_ppm' -v '2500' # 0.25%
	dasel put -t json -f "$GENESIS" '.app_state.prices.market_prices.[]' -v "{}"
	dasel put -t int -f "$GENESIS" '.app_state.prices.market_prices.[10].id' -v '10'
	dasel put -t int -f "$GENESIS" '.app_state.prices.market_prices.[10].exponent' -v '-11'
	dasel put -t int -f "$GENESIS" '.app_state.prices.market_prices.[10].price' -v '7320836895'          # $0.07321 = 1 DOGE.
	# DOGE Exchange Config
	doge_exchange_config_json=$(cat "$EXCHANGE_CONFIG_JSON_DIR/doge_exchange_config.json" | jq -c '.')
	dasel put -t string -f "$GENESIS" '.app_state.prices.market_params.[10].exchange_config_json' -v "$doge_exchange_config_json"

	# Market: ATOM-USD
	dasel put -t json -f "$GENESIS" '.app_state.prices.market_params.[]' -v "{}"
	dasel put -t string -f "$GENESIS" '.app_state.prices.market_params.[11].pair' -v 'ATOM-USD'
	dasel put -t int -f "$GENESIS" '.app_state.prices.market_params.[11].id' -v '11'
	dasel put -t int -f "$GENESIS" '.app_state.prices.market_params.[11].exponent' -v '-9'
	dasel put -t int -f "$GENESIS" '.app_state.prices.market_params.[11].min_exchanges' -v '3'
	dasel put -t int -f "$GENESIS" '.app_state.prices.market_params.[11].min_price_change_ppm' -v '2500' # 0.25%
	dasel put -t json -f "$GENESIS" '.app_state.prices.market_prices.[]' -v "{}"
	dasel put -t int -f "$GENESIS" '.app_state.prices.market_prices.[11].id' -v '11'
	dasel put -t int -f "$GENESIS" '.app_state.prices.market_prices.[11].exponent' -v '-9'
	dasel put -t int -f "$GENESIS" '.app_state.prices.market_prices.[11].price' -v '8433494428'          # $8.433 = 1 ATOM.
	# ATOM Exchange Config
	atom_exchange_config_json=$(cat "$EXCHANGE_CONFIG_JSON_DIR/atom_exchange_config.json" | jq -c '.')
	dasel put -t string -f "$GENESIS" '.app_state.prices.market_params.[11].exchange_config_json' -v "$atom_exchange_config_json"

	# Market: DOT-USD
	dasel put -t json -f "$GENESIS" '.app_state.prices.market_params.[]' -v "{}"
	dasel put -t string -f "$GENESIS" '.app_state.prices.market_params.[12].pair' -v 'DOT-USD'
	dasel put -t int -f "$GENESIS" '.app_state.prices.market_params.[12].id' -v '12'
	dasel put -t int -f "$GENESIS" '.app_state.prices.market_params.[12].exponent' -v '-9'
	dasel put -t int -f "$GENESIS" '.app_state.prices.market_params.[12].min_exchanges' -v '3'
	dasel put -t int -f "$GENESIS" '.app_state.prices.market_params.[12].min_price_change_ppm' -v '2500' # 0.25%
	dasel put -t json -f "$GENESIS" '.app_state.prices.market_prices.[]' -v "{}"
	dasel put -t int -f "$GENESIS" '.app_state.prices.market_prices.[12].id' -v '12'
	dasel put -t int -f "$GENESIS" '.app_state.prices.market_prices.[12].exponent' -v '-9'
	dasel put -t int -f "$GENESIS" '.app_state.prices.market_prices.[12].price' -v '4937186533'          # $4.937 = 1 DOT.
	# DOT Exchange Config
	dot_exchange_config_json=$(cat "$EXCHANGE_CONFIG_JSON_DIR/dot_exchange_config.json" | jq -c '.')
	dasel put -t string -f "$GENESIS" '.app_state.prices.market_params.[12].exchange_config_json' -v "$dot_exchange_config_json"

	# Market: UNI-USD
	dasel put -t json -f "$GENESIS" '.app_state.prices.market_params.[]' -v "{}"
	dasel put -t string -f "$GENESIS" '.app_state.prices.market_params.[13].pair' -v 'UNI-USD'
	dasel put -t int -f "$GENESIS" '.app_state.prices.market_params.[13].id' -v '13'
	dasel put -t int -f "$GENESIS" '.app_state.prices.market_params.[13].exponent' -v '-9'
	dasel put -t int -f "$GENESIS" '.app_state.prices.market_params.[13].min_exchanges' -v '3'
	dasel put -t int -f "$GENESIS" '.app_state.prices.market_params.[13].min_price_change_ppm' -v '2500' # 0.25%
	dasel put -t json -f "$GENESIS" '.app_state.prices.market_prices.[]' -v "{}"
	dasel put -t int -f "$GENESIS" '.app_state.prices.market_prices.[13].id' -v '13'
	dasel put -t int -f "$GENESIS" '.app_state.prices.market_prices.[13].exponent' -v '-9'
	dasel put -t int -f "$GENESIS" '.app_state.prices.market_prices.[13].price' -v '5852293356'          # $5.852 = 1 UNI.
	# UNI Exchange Config
	uni_exchange_config_json=$(cat "$EXCHANGE_CONFIG_JSON_DIR/uni_exchange_config.json" | jq -c '.')
	dasel put -t string -f "$GENESIS" '.app_state.prices.market_params.[13].exchange_config_json' -v "$uni_exchange_config_json"

	# Market: BCH-USD
	dasel put -t json -f "$GENESIS" '.app_state.prices.market_params.[]' -v "{}"
	dasel put -t string -f "$GENESIS" '.app_state.prices.market_params.[14].pair' -v 'BCH-USD'
	dasel put -t int -f "$GENESIS" '.app_state.prices.market_params.[14].id' -v '14'
	dasel put -t int -f "$GENESIS" '.app_state.prices.market_params.[14].exponent' -v '-7'
	dasel put -t int -f "$GENESIS" '.app_state.prices.market_params.[14].min_exchanges' -v '3'
	dasel put -t int -f "$GENESIS" '.app_state.prices.market_params.[14].min_price_change_ppm' -v '2500' # 0.25%
	dasel put -t json -f "$GENESIS" '.app_state.prices.market_prices.[]' -v "{}"
	dasel put -t int -f "$GENESIS" '.app_state.prices.market_prices.[14].id' -v '14'
	dasel put -t int -f "$GENESIS" '.app_state.prices.market_prices.[14].exponent' -v '-7'
	dasel put -t int -f "$GENESIS" '.app_state.prices.market_prices.[14].price' -v '2255676327'          # $225.6 = 1 BCH.
	# BCH Exchange Config
	bch_exchange_config_json=$(cat "$EXCHANGE_CONFIG_JSON_DIR/bch_exchange_config.json" | jq -c '.')
	dasel put -t string -f "$GENESIS" '.app_state.prices.market_params.[14].exchange_config_json' -v "$bch_exchange_config_json"

	# Market: TRX-USD
	dasel put -t json -f "$GENESIS" '.app_state.prices.market_params.[]' -v "{}"
	dasel put -t string -f "$GENESIS" '.app_state.prices.market_params.[15].pair' -v 'TRX-USD'
	dasel put -t int -f "$GENESIS" '.app_state.prices.market_params.[15].id' -v '15'
	dasel put -t int -f "$GENESIS" '.app_state.prices.market_params.[15].exponent' -v '-11'
	dasel put -t int -f "$GENESIS" '.app_state.prices.market_params.[15].min_exchanges' -v '3'
	dasel put -t int -f "$GENESIS" '.app_state.prices.market_params.[15].min_price_change_ppm' -v '2500' # 0.25%
	dasel put -t json -f "$GENESIS" '.app_state.prices.market_prices.[]' -v "{}"
	dasel put -t int -f "$GENESIS" '.app_state.prices.market_prices.[15].id' -v '15'
	dasel put -t int -f "$GENESIS" '.app_state.prices.market_prices.[15].exponent' -v '-11'
	dasel put -t int -f "$GENESIS" '.app_state.prices.market_prices.[15].price' -v '7795369902'          # $0.07795 = 1 TRX.
	# TRX Exchange Config
	trx_exchange_config_json=$(cat "$EXCHANGE_CONFIG_JSON_DIR/trx_exchange_config.json" | jq -c '.')
	dasel put -t string -f "$GENESIS" '.app_state.prices.market_params.[15].exchange_config_json' -v "$trx_exchange_config_json"

	# Market: NEAR-USD
	dasel put -t json -f "$GENESIS" '.app_state.prices.market_params.[]' -v "{}"
	dasel put -t string -f "$GENESIS" '.app_state.prices.market_params.[16].pair' -v 'NEAR-USD'
	dasel put -t int -f "$GENESIS" '.app_state.prices.market_params.[16].id' -v '16'
	dasel put -t int -f "$GENESIS" '.app_state.prices.market_params.[16].exponent' -v '-9'
	dasel put -t int -f "$GENESIS" '.app_state.prices.market_params.[16].min_exchanges' -v '3'
	dasel put -t int -f "$GENESIS" '.app_state.prices.market_params.[16].min_price_change_ppm' -v '2500' # 0.25%
	dasel put -t json -f "$GENESIS" '.app_state.prices.market_prices.[]' -v "{}"
	dasel put -t int -f "$GENESIS" '.app_state.prices.market_prices.[16].id' -v '16'
	dasel put -t int -f "$GENESIS" '.app_state.prices.market_prices.[16].exponent' -v '-9'
	dasel put -t int -f "$GENESIS" '.app_state.prices.market_prices.[16].price' -v '1312325536'          # $1.312 = 1 NEAR.
	# NEAR Exchange Config
	near_exchange_config_json=$(cat "$EXCHANGE_CONFIG_JSON_DIR/near_exchange_config.json" | jq -c '.')
	dasel put -t string -f "$GENESIS" '.app_state.prices.market_params.[16].exchange_config_json' -v "$near_exchange_config_json"

	# Market: MKR-USD
	dasel put -t json -f "$GENESIS" '.app_state.prices.market_params.[]' -v "{}"
	dasel put -t string -f "$GENESIS" '.app_state.prices.market_params.[17].pair' -v 'MKR-USD'
	dasel put -t int -f "$GENESIS" '.app_state.prices.market_params.[17].id' -v '17'
	dasel put -t int -f "$GENESIS" '.app_state.prices.market_params.[17].exponent' -v '-6'
	dasel put -t int -f "$GENESIS" '.app_state.prices.market_params.[17].min_exchanges' -v '3'
	dasel put -t int -f "$GENESIS" '.app_state.prices.market_params.[17].min_price_change_ppm' -v '4000' # 0.4%
	dasel put -t json -f "$GENESIS" '.app_state.prices.market_prices.[]' -v "{}"
	dasel put -t int -f "$GENESIS" '.app_state.prices.market_prices.[17].id' -v '17'
	dasel put -t int -f "$GENESIS" '.app_state.prices.market_prices.[17].exponent' -v '-6'
	dasel put -t int -f "$GENESIS" '.app_state.prices.market_prices.[17].price' -v '1199517382'          # $1,200 = 1 MKR.
	# MKR Exchange Config
	mkr_exchange_config_json=$(cat "$EXCHANGE_CONFIG_JSON_DIR/mkr_exchange_config.json" | jq -c '.')
	dasel put -t string -f "$GENESIS" '.app_state.prices.market_params.[17].exchange_config_json' -v "$mkr_exchange_config_json"

	# Market: XLM-USD
	dasel put -t json -f "$GENESIS" '.app_state.prices.market_params.[]' -v "{}"
	dasel put -t string -f "$GENESIS" '.app_state.prices.market_params.[18].pair' -v 'XLM-USD'
	dasel put -t int -f "$GENESIS" '.app_state.prices.market_params.[18].id' -v '18'
	dasel put -t int -f "$GENESIS" '.app_state.prices.market_params.[18].exponent' -v '-10'
	dasel put -t int -f "$GENESIS" '.app_state.prices.market_params.[18].min_exchanges' -v '3'
	dasel put -t int -f "$GENESIS" '.app_state.prices.market_params.[18].min_price_change_ppm' -v '2500' # 0.25%
	dasel put -t json -f "$GENESIS" '.app_state.prices.market_prices.[]' -v "{}"
	dasel put -t int -f "$GENESIS" '.app_state.prices.market_prices.[18].id' -v '18'
	dasel put -t int -f "$GENESIS" '.app_state.prices.market_prices.[18].exponent' -v '-10'
	dasel put -t int -f "$GENESIS" '.app_state.prices.market_prices.[18].price' -v '1398578933'          # $0.1399 = 1 XLM.
	# XLM Exchange Config
	xlm_exchange_config_json=$(cat "$EXCHANGE_CONFIG_JSON_DIR/xlm_exchange_config.json" | jq -c '.')
	dasel put -t string -f "$GENESIS" '.app_state.prices.market_params.[18].exchange_config_json' -v "$xlm_exchange_config_json"

	# Market: ETC-USD
	dasel put -t json -f "$GENESIS" '.app_state.prices.market_params.[]' -v "{}"
	dasel put -t string -f "$GENESIS" '.app_state.prices.market_params.[19].pair' -v 'ETC-USD'
	dasel put -t int -f "$GENESIS" '.app_state.prices.market_params.[19].id' -v '19'
	dasel put -t int -f "$GENESIS" '.app_state.prices.market_params.[19].exponent' -v '-8'
	dasel put -t int -f "$GENESIS" '.app_state.prices.market_params.[19].min_exchanges' -v '3'
	dasel put -t int -f "$GENESIS" '.app_state.prices.market_params.[19].min_price_change_ppm' -v '2500' # 0.25%
	dasel put -t json -f "$GENESIS" '.app_state.prices.market_prices.[]' -v "{}"
	dasel put -t int -f "$GENESIS" '.app_state.prices.market_prices.[19].id' -v '19'
	dasel put -t int -f "$GENESIS" '.app_state.prices.market_prices.[19].exponent' -v '-8'
	dasel put -t int -f "$GENESIS" '.app_state.prices.market_prices.[19].price' -v '1741060746'          # $17.41 = 1 ETC.
	# ETC Exchange Config
	etc_exchange_config_json=$(cat "$EXCHANGE_CONFIG_JSON_DIR/etc_exchange_config.json" | jq -c '.')
	dasel put -t string -f "$GENESIS" '.app_state.prices.market_params.[19].exchange_config_json' -v "$etc_exchange_config_json"

	# Market: COMP-USD
	dasel put -t json -f "$GENESIS" '.app_state.prices.market_params.[]' -v "{}"
	dasel put -t string -f "$GENESIS" '.app_state.prices.market_params.[20].pair' -v 'COMP-USD'
	dasel put -t int -f "$GENESIS" '.app_state.prices.market_params.[20].id' -v '20'
	dasel put -t int -f "$GENESIS" '.app_state.prices.market_params.[20].exponent' -v '-8'
	dasel put -t int -f "$GENESIS" '.app_state.prices.market_params.[20].min_exchanges' -v '3'
	dasel put -t int -f "$GENESIS" '.app_state.prices.market_params.[20].min_price_change_ppm' -v '4000' # 0.4%
	dasel put -t json -f "$GENESIS" '.app_state.prices.market_prices.[]' -v "{}"
	dasel put -t int -f "$GENESIS" '.app_state.prices.market_prices.[20].id' -v '20'
	dasel put -t int -f "$GENESIS" '.app_state.prices.market_prices.[20].exponent' -v '-8'
	dasel put -t int -f "$GENESIS" '.app_state.prices.market_prices.[20].price' -v '5717635307'          # $57.18 = 1 COMP.
	# COMP Exchange Config
	comp_exchange_config_json=$(cat "$EXCHANGE_CONFIG_JSON_DIR/comp_exchange_config.json" | jq -c '.')
	dasel put -t string -f "$GENESIS" '.app_state.prices.market_params.[20].exchange_config_json' -v "$comp_exchange_config_json"

	# Market: WLD-USD
	dasel put -t json -f "$GENESIS" '.app_state.prices.market_params.[]' -v "{}"
	dasel put -t string -f "$GENESIS" '.app_state.prices.market_params.[21].pair' -v 'WLD-USD'
	dasel put -t int -f "$GENESIS" '.app_state.prices.market_params.[21].id' -v '21'
	dasel put -t int -f "$GENESIS" '.app_state.prices.market_params.[21].exponent' -v '-9'
	dasel put -t int -f "$GENESIS" '.app_state.prices.market_params.[21].min_exchanges' -v '3'
	dasel put -t int -f "$GENESIS" '.app_state.prices.market_params.[21].min_price_change_ppm' -v '2500' # 0.25%
	dasel put -t json -f "$GENESIS" '.app_state.prices.market_prices.[]' -v "{}"
	dasel put -t int -f "$GENESIS" '.app_state.prices.market_prices.[21].id' -v '21'
	dasel put -t int -f "$GENESIS" '.app_state.prices.market_prices.[21].exponent' -v '-9'
	dasel put -t int -f "$GENESIS" '.app_state.prices.market_prices.[21].price' -v '1943019371'          # $1.943 = 1 WLD.
	# WLD Exchange Config
	wld_exchange_config_json=$(cat "$EXCHANGE_CONFIG_JSON_DIR/wld_exchange_config.json" | jq -c '.')
	dasel put -t string -f "$GENESIS" '.app_state.prices.market_params.[21].exchange_config_json' -v "$wld_exchange_config_json"

	# Market: APE-USD
	dasel put -t json -f "$GENESIS" '.app_state.prices.market_params.[]' -v "{}"
	dasel put -t string -f "$GENESIS" '.app_state.prices.market_params.[22].pair' -v 'APE-USD'
	dasel put -t int -f "$GENESIS" '.app_state.prices.market_params.[22].id' -v '22'
	dasel put -t int -f "$GENESIS" '.app_state.prices.market_params.[22].exponent' -v '-9'
	dasel put -t int -f "$GENESIS" '.app_state.prices.market_params.[22].min_exchanges' -v '3'
	dasel put -t int -f "$GENESIS" '.app_state.prices.market_params.[22].min_price_change_ppm' -v '4000' # 0.4%
	dasel put -t json -f "$GENESIS" '.app_state.prices.market_prices.[]' -v "{}"
	dasel put -t int -f "$GENESIS" '.app_state.prices.market_prices.[22].id' -v '22'
	dasel put -t int -f "$GENESIS" '.app_state.prices.market_prices.[22].exponent' -v '-9'
	dasel put -t int -f "$GENESIS" '.app_state.prices.market_prices.[22].price' -v '1842365656'          # $1.842 = 1 APE.
	# APE Exchange Config
	ape_exchange_config_json=$(cat "$EXCHANGE_CONFIG_JSON_DIR/ape_exchange_config.json" | jq -c '.')
	dasel put -t string -f "$GENESIS" '.app_state.prices.market_params.[22].exchange_config_json' -v "$ape_exchange_config_json"

	# Market: APT-USD
	dasel put -t json -f "$GENESIS" '.app_state.prices.market_params.[]' -v "{}"
	dasel put -t string -f "$GENESIS" '.app_state.prices.market_params.[23].pair' -v 'APT-USD'
	dasel put -t int -f "$GENESIS" '.app_state.prices.market_params.[23].id' -v '23'
	dasel put -t int -f "$GENESIS" '.app_state.prices.market_params.[23].exponent' -v '-9'
	dasel put -t int -f "$GENESIS" '.app_state.prices.market_params.[23].min_exchanges' -v '3'
	dasel put -t int -f "$GENESIS" '.app_state.prices.market_params.[23].min_price_change_ppm' -v '2500' # 0.25%
	dasel put -t json -f "$GENESIS" '.app_state.prices.market_prices.[]' -v "{}"
	dasel put -t int -f "$GENESIS" '.app_state.prices.market_prices.[23].id' -v '23'
	dasel put -t int -f "$GENESIS" '.app_state.prices.market_prices.[23].exponent' -v '-9'
	dasel put -t int -f "$GENESIS" '.app_state.prices.market_prices.[23].price' -v '6787621897'          # $6.788 = 1 APT.
	# APT Exchange Config
	apt_exchange_config_json=$(cat "$EXCHANGE_CONFIG_JSON_DIR/apt_exchange_config.json" | jq -c '.')
	dasel put -t string -f "$GENESIS" '.app_state.prices.market_params.[23].exchange_config_json' -v "$apt_exchange_config_json"

	# Market: ARB-USD
	dasel put -t json -f "$GENESIS" '.app_state.prices.market_params.[]' -v "{}"
	dasel put -t string -f "$GENESIS" '.app_state.prices.market_params.[24].pair' -v 'ARB-USD'
	dasel put -t int -f "$GENESIS" '.app_state.prices.market_params.[24].id' -v '24'
	dasel put -t int -f "$GENESIS" '.app_state.prices.market_params.[24].exponent' -v '-9'
	dasel put -t int -f "$GENESIS" '.app_state.prices.market_params.[24].min_exchanges' -v '3'
	dasel put -t int -f "$GENESIS" '.app_state.prices.market_params.[24].min_price_change_ppm' -v '2500' # 0.25%
	dasel put -t json -f "$GENESIS" '.app_state.prices.market_prices.[]' -v "{}"
	dasel put -t int -f "$GENESIS" '.app_state.prices.market_prices.[24].id' -v '24'
	dasel put -t int -f "$GENESIS" '.app_state.prices.market_prices.[24].exponent' -v '-9'
	dasel put -t int -f "$GENESIS" '.app_state.prices.market_prices.[24].price' -v '1127629325'          # $1.128 = 1 ARB.
	# ARB Exchange Config
	arb_exchange_config_json=$(cat "$EXCHANGE_CONFIG_JSON_DIR/arb_exchange_config.json" | jq -c '.')
	dasel put -t string -f "$GENESIS" '.app_state.prices.market_params.[24].exchange_config_json' -v "$arb_exchange_config_json"

	# Market: BLUR-USD
	dasel put -t json -f "$GENESIS" '.app_state.prices.market_params.[]' -v "{}"
	dasel put -t string -f "$GENESIS" '.app_state.prices.market_params.[25].pair' -v 'BLUR-USD'
	dasel put -t int -f "$GENESIS" '.app_state.prices.market_params.[25].id' -v '25'
	dasel put -t int -f "$GENESIS" '.app_state.prices.market_params.[25].exponent' -v '-10'
	dasel put -t int -f "$GENESIS" '.app_state.prices.market_params.[25].min_exchanges' -v '3'
	dasel put -t int -f "$GENESIS" '.app_state.prices.market_params.[25].min_price_change_ppm' -v '4000' # 0.4%
	dasel put -t json -f "$GENESIS" '.app_state.prices.market_prices.[]' -v "{}"
	dasel put -t int -f "$GENESIS" '.app_state.prices.market_prices.[25].id' -v '25'
	dasel put -t int -f "$GENESIS" '.app_state.prices.market_prices.[25].exponent' -v '-10'
	dasel put -t int -f "$GENESIS" '.app_state.prices.market_prices.[25].price' -v '2779565892'          # $.2780 = 1 BLUR.
	# BLUR Exchange Config
	blur_exchange_config_json=$(cat "$EXCHANGE_CONFIG_JSON_DIR/blur_exchange_config.json" | jq -c '.')
	dasel put -t string -f "$GENESIS" '.app_state.prices.market_params.[25].exchange_config_json' -v "$blur_exchange_config_json"

	# Market: LDO-USD
	dasel put -t json -f "$GENESIS" '.app_state.prices.market_params.[]' -v "{}"
	dasel put -t string -f "$GENESIS" '.app_state.prices.market_params.[26].pair' -v 'LDO-USD'
	dasel put -t int -f "$GENESIS" '.app_state.prices.market_params.[26].id' -v '26'
	dasel put -t int -f "$GENESIS" '.app_state.prices.market_params.[26].exponent' -v '-9'
	dasel put -t int -f "$GENESIS" '.app_state.prices.market_params.[26].min_exchanges' -v '3'
	dasel put -t int -f "$GENESIS" '.app_state.prices.market_params.[26].min_price_change_ppm' -v '4000' # 0.4%
	dasel put -t json -f "$GENESIS" '.app_state.prices.market_prices.[]' -v "{}"
	dasel put -t int -f "$GENESIS" '.app_state.prices.market_prices.[26].id' -v '26'
	dasel put -t int -f "$GENESIS" '.app_state.prices.market_prices.[26].exponent' -v '-9'
	dasel put -t int -f "$GENESIS" '.app_state.prices.market_prices.[26].price' -v '1855061997'          # $1.855 = 1 LDO.
	# LDO Exchange Config
	ldo_exchange_config_json=$(cat "$EXCHANGE_CONFIG_JSON_DIR/ldo_exchange_config.json" | jq -c '.')
	dasel put -t string -f "$GENESIS" '.app_state.prices.market_params.[26].exchange_config_json' -v "$ldo_exchange_config_json"

	# Market: OP-USD
	dasel put -t json -f "$GENESIS" '.app_state.prices.market_params.[]' -v "{}"
	dasel put -t string -f "$GENESIS" '.app_state.prices.market_params.[27].pair' -v 'OP-USD'
	dasel put -t int -f "$GENESIS" '.app_state.prices.market_params.[27].id' -v '27'
	dasel put -t int -f "$GENESIS" '.app_state.prices.market_params.[27].exponent' -v '-9'
	dasel put -t int -f "$GENESIS" '.app_state.prices.market_params.[27].min_exchanges' -v '3'
	dasel put -t int -f "$GENESIS" '.app_state.prices.market_params.[27].min_price_change_ppm' -v '2500' # 0.25%
	dasel put -t json -f "$GENESIS" '.app_state.prices.market_prices.[]' -v "{}"
	dasel put -t int -f "$GENESIS" '.app_state.prices.market_prices.[27].id' -v '27'
	dasel put -t int -f "$GENESIS" '.app_state.prices.market_prices.[27].exponent' -v '-9'
	dasel put -t int -f "$GENESIS" '.app_state.prices.market_prices.[27].price' -v '1562218603'          # $1.562 = 1 OP.
	# OP Exchange Config
	op_exchange_config_json=$(cat "$EXCHANGE_CONFIG_JSON_DIR/op_exchange_config.json" | jq -c '.')
	dasel put -t string -f "$GENESIS" '.app_state.prices.market_params.[27].exchange_config_json' -v "$op_exchange_config_json"

	# Market: PEPE-USD
	dasel put -t json -f "$GENESIS" '.app_state.prices.market_params.[]' -v "{}"
	dasel put -t string -f "$GENESIS" '.app_state.prices.market_params.[28].pair' -v 'PEPE-USD'
	dasel put -t int -f "$GENESIS" '.app_state.prices.market_params.[28].id' -v '28'
	dasel put -t int -f "$GENESIS" '.app_state.prices.market_params.[28].exponent' -v '-16'
	dasel put -t int -f "$GENESIS" '.app_state.prices.market_params.[28].min_exchanges' -v '3'
	dasel put -t int -f "$GENESIS" '.app_state.prices.market_params.[28].min_price_change_ppm' -v '2500' # 0.25%
	dasel put -t json -f "$GENESIS" '.app_state.prices.market_prices.[]' -v "{}"
	dasel put -t int -f "$GENESIS" '.app_state.prices.market_prices.[28].id' -v '28'
	dasel put -t int -f "$GENESIS" '.app_state.prices.market_prices.[28].exponent' -v '-16'
	dasel put -t int -f "$GENESIS" '.app_state.prices.market_prices.[28].price' -v '2481900353'          # $.000000248190035 = 1 PEPE.
	# PEPE Exchange Config
	pepe_exchange_config_json=$(cat "$EXCHANGE_CONFIG_JSON_DIR/pepe_exchange_config.json" | jq -c '.')
	dasel put -t string -f "$GENESIS" '.app_state.prices.market_params.[28].exchange_config_json' -v "$pepe_exchange_config_json"

	# Market: SEI-USD
	dasel put -t json -f "$GENESIS" '.app_state.prices.market_params.[]' -v "{}"
	dasel put -t string -f "$GENESIS" '.app_state.prices.market_params.[29].pair' -v 'SEI-USD'
	dasel put -t int -f "$GENESIS" '.app_state.prices.market_params.[29].id' -v '29'
	dasel put -t int -f "$GENESIS" '.app_state.prices.market_params.[29].exponent' -v '-10'
	dasel put -t int -f "$GENESIS" '.app_state.prices.market_params.[29].min_exchanges' -v '3'
	dasel put -t int -f "$GENESIS" '.app_state.prices.market_params.[29].min_price_change_ppm' -v '4000' # 0.4%
	dasel put -t json -f "$GENESIS" '.app_state.prices.market_prices.[]' -v "{}"
	dasel put -t int -f "$GENESIS" '.app_state.prices.market_prices.[29].id' -v '29'
	dasel put -t int -f "$GENESIS" '.app_state.prices.market_prices.[29].exponent' -v '-10'
	dasel put -t int -f "$GENESIS" '.app_state.prices.market_prices.[29].price' -v '1686998025'          # $.1687 = 1 SEI.
	# SEI Exchange Config
	sei_exchange_config_json=$(cat "$EXCHANGE_CONFIG_JSON_DIR/sei_exchange_config.json" | jq -c '.')
	dasel put -t string -f "$GENESIS" '.app_state.prices.market_params.[29].exchange_config_json' -v "$sei_exchange_config_json"

	# Market: SHIB-USD
	dasel put -t json -f "$GENESIS" '.app_state.prices.market_params.[]' -v "{}"
	dasel put -t string -f "$GENESIS" '.app_state.prices.market_params.[30].pair' -v 'SHIB-USD'
	dasel put -t int -f "$GENESIS" '.app_state.prices.market_params.[30].id' -v '30'
	dasel put -t int -f "$GENESIS" '.app_state.prices.market_params.[30].exponent' -v '-15'
	dasel put -t int -f "$GENESIS" '.app_state.prices.market_params.[30].min_exchanges' -v '3'
	dasel put -t int -f "$GENESIS" '.app_state.prices.market_params.[30].min_price_change_ppm' -v '2500' # 0.25%
	dasel put -t json -f "$GENESIS" '.app_state.prices.market_prices.[]' -v "{}"
	dasel put -t int -f "$GENESIS" '.app_state.prices.market_prices.[30].id' -v '30'
	dasel put -t int -f "$GENESIS" '.app_state.prices.market_prices.[30].exponent' -v '-15'
	dasel put -t int -f "$GENESIS" '.app_state.prices.market_prices.[30].price' -v '8895882688'          # $.000008896 = 1 SHIB.
	# SHIB Exchange Config
	shib_exchange_config_json=$(cat "$EXCHANGE_CONFIG_JSON_DIR/shib_exchange_config.json" | jq -c '.')
	dasel put -t string -f "$GENESIS" '.app_state.prices.market_params.[30].exchange_config_json' -v "$shib_exchange_config_json"

	# Market: SUI-USD
	dasel put -t json -f "$GENESIS" '.app_state.prices.market_params.[]' -v "{}"
	dasel put -t string -f "$GENESIS" '.app_state.prices.market_params.[31].pair' -v 'SUI-USD'
	dasel put -t int -f "$GENESIS" '.app_state.prices.market_params.[31].id' -v '31'
	dasel put -t int -f "$GENESIS" '.app_state.prices.market_params.[31].exponent' -v '-10'
	dasel put -t int -f "$GENESIS" '.app_state.prices.market_params.[31].min_exchanges' -v '3'
	dasel put -t int -f "$GENESIS" '.app_state.prices.market_params.[31].min_price_change_ppm' -v '2500' # 0.25%
	dasel put -t json -f "$GENESIS" '.app_state.prices.market_prices.[]' -v "{}"
	dasel put -t int -f "$GENESIS" '.app_state.prices.market_prices.[31].id' -v '31'
	dasel put -t int -f "$GENESIS" '.app_state.prices.market_prices.[31].exponent' -v '-10'
	dasel put -t int -f "$GENESIS" '.app_state.prices.market_prices.[31].price' -v '5896318772'          # $.5896 = 1 SUI.
	# SUI Exchange Config
	sui_exchange_config_json=$(cat "$EXCHANGE_CONFIG_JSON_DIR/sui_exchange_config.json" | jq -c '.')
	dasel put -t string -f "$GENESIS" '.app_state.prices.market_params.[31].exchange_config_json' -v "$sui_exchange_config_json"

	# Market: XRP-USD
	dasel put -t json -f "$GENESIS" '.app_state.prices.market_params.[]' -v "{}"
	dasel put -t string -f "$GENESIS" '.app_state.prices.market_params.[32].pair' -v 'XRP-USD'
	dasel put -t int -f "$GENESIS" '.app_state.prices.market_params.[32].id' -v '32'
	dasel put -t int -f "$GENESIS" '.app_state.prices.market_params.[32].exponent' -v '-10'
	dasel put -t int -f "$GENESIS" '.app_state.prices.market_params.[32].min_exchanges' -v '3'
	dasel put -t int -f "$GENESIS" '.app_state.prices.market_params.[32].min_price_change_ppm' -v '2500' # 0.25%
	dasel put -t json -f "$GENESIS" '.app_state.prices.market_prices.[]' -v "{}"
	dasel put -t int -f "$GENESIS" '.app_state.prices.market_prices.[32].id' -v '32'
	dasel put -t int -f "$GENESIS" '.app_state.prices.market_prices.[32].exponent' -v '-10'
	dasel put -t int -f "$GENESIS" '.app_state.prices.market_prices.[32].price' -v '6327613800'          # $.6328 = 1 XRP.
	# XRP Exchange Config
	xrp_exchange_config_json=$(cat "$EXCHANGE_CONFIG_JSON_DIR/xrp_exchange_config.json" | jq -c '.')
	dasel put -t string -f "$GENESIS" '.app_state.prices.market_params.[32].exchange_config_json' -v "$xrp_exchange_config_json"

	# Market: USDT-USD
	dasel put -t json -f "$GENESIS" '.app_state.prices.market_params.[]' -v "{}"
	dasel put -t string -f "$GENESIS" '.app_state.prices.market_params.[33].pair' -v 'USDT-USD'
	dasel put -t int -f "$GENESIS" '.app_state.prices.market_params.[33].id' -v '1000000'
	dasel put -t int -f "$GENESIS" '.app_state.prices.market_params.[33].exponent' -v '-9'
	dasel put -t int -f "$GENESIS" '.app_state.prices.market_params.[33].min_exchanges' -v '3'
	dasel put -t int -f "$GENESIS" '.app_state.prices.market_params.[33].min_price_change_ppm' -v '1000'  # 0.100%
	dasel put -t json -f "$GENESIS" '.app_state.prices.market_prices.[]' -v "{}"
	dasel put -t int -f "$GENESIS" '.app_state.prices.market_prices.[33].id' -v '1000000'
	dasel put -t int -f "$GENESIS" '.app_state.prices.market_prices.[33].exponent' -v '-9'
	dasel put -t int -f "$GENESIS" '.app_state.prices.market_prices.[33].price' -v '1000000000'          # $1 = 1 USDT.
	# USDT Exchange Config
	usdt_exchange_config_json=$(cat "$EXCHANGE_CONFIG_JSON_DIR/usdt_exchange_config.json" | jq -c '.')
	dasel put -t string -f "$GENESIS" '.app_state.prices.market_params.[33].exchange_config_json' -v "$usdt_exchange_config_json"

	# Market: DYDX-USD
	dasel put -t json -f "$GENESIS" '.app_state.prices.market_params.[]' -v "{}"
	dasel put -t string -f "$GENESIS" '.app_state.prices.market_params.[34].pair' -v 'DYDX-USD'
	dasel put -t int -f "$GENESIS" '.app_state.prices.market_params.[34].id' -v '1000001'
	dasel put -t int -f "$GENESIS" '.app_state.prices.market_params.[34].exponent' -v '-9'
	dasel put -t int -f "$GENESIS" '.app_state.prices.market_params.[34].min_exchanges' -v '3'
	dasel put -t int -f "$GENESIS" '.app_state.prices.market_params.[34].min_price_change_ppm' -v '2500'  # 0.25%
	dasel put -t json -f "$GENESIS" '.app_state.prices.market_prices.[]' -v "{}"
	dasel put -t int -f "$GENESIS" '.app_state.prices.market_prices.[34].id' -v '1000001'
	dasel put -t int -f "$GENESIS" '.app_state.prices.market_prices.[34].exponent' -v '-9'
	dasel put -t int -f "$GENESIS" '.app_state.prices.market_prices.[34].price' -v '2050000000'          # $2.05 = 1 DYDX.
	# DYDX Exchange Config
	dydx_exchange_config_json=$(cat "$EXCHANGE_CONFIG_JSON_DIR/dydx_exchange_config.json" | jq -c '.')
	dasel put -t string -f "$GENESIS" '.app_state.prices.market_params.[34].exchange_config_json' -v "$dydx_exchange_config_json"

	# Initialize bridge module account balance as total native token supply.
	bridge_module_account_balance=$TOTAL_NATIVE_TOKEN_SUPPLY
	total_accounts_quote_balance=0
	acct_idx=0
	# Update subaccounts module for load testing accounts and update bridge module account balance.
	for acct in "${INPUT_TEST_ACCOUNTS[@]}"; do
		add_subaccount "$GENESIS" "$acct_idx" "$acct" "$DEFAULT_SUBACCOUNT_QUOTE_BALANCE"
		total_accounts_quote_balance=$(($total_accounts_quote_balance + $DEFAULT_SUBACCOUNT_QUOTE_BALANCE))
		bridge_module_account_balance=$(echo "$bridge_module_account_balance - $TESTNET_VALIDATOR_NATIVE_TOKEN_BALANCE" | bc)
		acct_idx=$(($acct_idx + 1))
	done
	# Update subaccounts module for faucet accounts and update bridge module account balance.
	for acct in "${INPUT_FAUCET_ACCOUNTS[@]}"; do
		add_subaccount "$GENESIS" "$acct_idx" "$acct" "$DEFAULT_SUBACCOUNT_QUOTE_BALANCE_FAUCET"
		total_accounts_quote_balance=$(($total_accounts_quote_balance + $DEFAULT_SUBACCOUNT_QUOTE_BALANCE_FAUCET))
		bridge_module_account_balance=$(echo "$bridge_module_account_balance - $TESTNET_VALIDATOR_NATIVE_TOKEN_BALANCE" | bc)
		acct_idx=$(($acct_idx + 1))
	done

	next_bank_idx=0
	if (( total_accounts_quote_balance > 0 )); then
		# Initialize subaccounts module account balance within bank module.
		dasel put -t json -f "$GENESIS" ".app_state.bank.balances.[]" -v "{}"
		dasel put -t string -f "$GENESIS" ".app_state.bank.balances.[0].address" -v "${SUBACCOUNTS_MODACC_ADDR}"
		dasel put -t json -f "$GENESIS" ".app_state.bank.balances.[0].coins.[]" -v "{}"
		dasel put -t string -f "$GENESIS" ".app_state.bank.balances.[0].coins.[0].denom" -v "$USDC_DENOM"
		# TODO(DEC-969): For testnet, ensure subaccounts module balance >= sum of subaccount quote balances.
		dasel put -t string -f "$GENESIS" ".app_state.bank.balances.[0].coins.[0].amount" -v "${total_accounts_quote_balance}"
		next_bank_idx=$(($next_bank_idx+1))
	fi

	if [ $(echo "$REWARDS_VESTER_ACCOUNT_BALANCE > 0" | bc -l) -eq 1 ]; then
		# Initialize bank balance of reward vester account.
		dasel put -t json -f "$GENESIS" ".app_state.bank.balances.[]" -v "{}"
		dasel put -t string -f "$GENESIS" ".app_state.bank.balances.[$next_bank_idx].address" -v "${REWARDS_VESTER_ACCOUNT_ADDR}"
		dasel put -t json -f "$GENESIS" ".app_state.bank.balances.[$next_bank_idx].coins.[]" -v "{}"
		dasel put -t string -f "$GENESIS" ".app_state.bank.balances.[$next_bank_idx].coins.[0].denom" -v "${REWARD_TOKEN}"
		dasel put -t string -f "$GENESIS" ".app_state.bank.balances.[$next_bank_idx].coins.[0].amount" -v "$REWARDS_VESTER_ACCOUNT_BALANCE"
		next_bank_idx=$(($next_bank_idx+1))

		bridge_module_account_balance=$(echo "$bridge_module_account_balance - $REWARDS_VESTER_ACCOUNT_BALANCE" | bc)
	fi

	# Initialize bank balance of bridge module account.
	dasel put -t json -f "$GENESIS" ".app_state.bank.balances.[]" -v "{}"
	dasel put -t string -f "$GENESIS" ".app_state.bank.balances.[$next_bank_idx].address" -v "${BRIDGE_MODACC_ADDR}"
	dasel put -t json -f "$GENESIS" ".app_state.bank.balances.[$next_bank_idx].coins.[]" -v "{}"
	dasel put -t string -f "$GENESIS" ".app_state.bank.balances.[$next_bank_idx].coins.[0].denom" -v "${NATIVE_TOKEN}"
	dasel put -t string -f "$GENESIS" ".app_state.bank.balances.[$next_bank_idx].coins.[0].amount" -v "${bridge_module_account_balance}"

	# Set denom metadata
	set_denom_metadata "$NATIVE_TOKEN" "$NATIVE_TOKEN_WHOLE_COIN" "$COIN_NAME"

	# Use ATOM-USD as test oracle price of the reward token.
	dasel put -t int -f "$GENESIS" '.app_state.rewards.params.market_id' -v '11'

	# Update clob module.
	# Clob: BTC-USD
	dasel put -t json -f "$GENESIS" '.app_state.clob.clob_pairs.[]' -v "{}"
	dasel put -t int -f "$GENESIS" '.app_state.clob.clob_pairs.[0].id' -v '0'
	dasel put -t string -f "$GENESIS" '.app_state.clob.clob_pairs.[0].status' -v "$INITIAL_CLOB_PAIR_STATUS"
	dasel put -t int -f "$GENESIS" '.app_state.clob.clob_pairs.[0].perpetual_clob_metadata.perpetual_id' -v '0'
	dasel put -t int -f "$GENESIS" '.app_state.clob.clob_pairs.[0].step_base_quantums' -v '1000000'
	dasel put -t int -f "$GENESIS" '.app_state.clob.clob_pairs.[0].subticks_per_tick' -v '100000'
	dasel put -t int -f "$GENESIS" '.app_state.clob.clob_pairs.[0].quantum_conversion_exponent' -v '-9'

	# Clob: ETH-USD
	dasel put -t json -f "$GENESIS" '.app_state.clob.clob_pairs.[]' -v "{}"
	dasel put -t int -f "$GENESIS" '.app_state.clob.clob_pairs.[1].id' -v '1'
	dasel put -t string -f "$GENESIS" '.app_state.clob.clob_pairs.[1].status' -v "$INITIAL_CLOB_PAIR_STATUS"
	dasel put -t int -f "$GENESIS" '.app_state.clob.clob_pairs.[1].perpetual_clob_metadata.perpetual_id' -v '1'
	dasel put -t int -f "$GENESIS" '.app_state.clob.clob_pairs.[1].step_base_quantums' -v '1000000'
	dasel put -t int -f "$GENESIS" '.app_state.clob.clob_pairs.[1].subticks_per_tick' -v '100000'
	dasel put -t int -f "$GENESIS" '.app_state.clob.clob_pairs.[1].quantum_conversion_exponent' -v '-9'

	# Clob: LINK-USD
	dasel put -t json -f "$GENESIS" '.app_state.clob.clob_pairs.[]' -v "{}"
	dasel put -t int -f "$GENESIS" '.app_state.clob.clob_pairs.[2].id' -v '2'
	dasel put -t string -f "$GENESIS" '.app_state.clob.clob_pairs.[2].status' -v "$INITIAL_CLOB_PAIR_STATUS"
	dasel put -t int -f "$GENESIS" '.app_state.clob.clob_pairs.[2].perpetual_clob_metadata.perpetual_id' -v '2'
	dasel put -t int -f "$GENESIS" '.app_state.clob.clob_pairs.[2].step_base_quantums' -v '1000000'
	dasel put -t int -f "$GENESIS" '.app_state.clob.clob_pairs.[2].subticks_per_tick' -v '1000000'
	dasel put -t int -f "$GENESIS" '.app_state.clob.clob_pairs.[2].quantum_conversion_exponent' -v '-9'

	# Clob: MATIC-USD
	dasel put -t json -f "$GENESIS" '.app_state.clob.clob_pairs.[]' -v "{}"
	dasel put -t int -f "$GENESIS" '.app_state.clob.clob_pairs.[3].id' -v '3'
	dasel put -t string -f "$GENESIS" '.app_state.clob.clob_pairs.[3].status' -v "$INITIAL_CLOB_PAIR_STATUS"
	dasel put -t int -f "$GENESIS" '.app_state.clob.clob_pairs.[3].perpetual_clob_metadata.perpetual_id' -v '3'
	dasel put -t int -f "$GENESIS" '.app_state.clob.clob_pairs.[3].step_base_quantums' -v '1000000'
	dasel put -t int -f "$GENESIS" '.app_state.clob.clob_pairs.[3].subticks_per_tick' -v '1000000'
	dasel put -t int -f "$GENESIS" '.app_state.clob.clob_pairs.[3].quantum_conversion_exponent' -v '-9'

	# Clob: CRV-USD
	dasel put -t json -f "$GENESIS" '.app_state.clob.clob_pairs.[]' -v "{}"
	dasel put -t int -f "$GENESIS" '.app_state.clob.clob_pairs.[4].id' -v '4'
	dasel put -t string -f "$GENESIS" '.app_state.clob.clob_pairs.[4].status' -v "$INITIAL_CLOB_PAIR_STATUS"
	dasel put -t int -f "$GENESIS" '.app_state.clob.clob_pairs.[4].perpetual_clob_metadata.perpetual_id' -v '4'
	dasel put -t int -f "$GENESIS" '.app_state.clob.clob_pairs.[4].step_base_quantums' -v '1000000'
	dasel put -t int -f "$GENESIS" '.app_state.clob.clob_pairs.[4].subticks_per_tick' -v '1000000'
	dasel put -t int -f "$GENESIS" '.app_state.clob.clob_pairs.[4].quantum_conversion_exponent' -v '-9'

	# Clob: SOL-USD
	dasel put -t json -f "$GENESIS" '.app_state.clob.clob_pairs.[]' -v "{}"
	dasel put -t int -f "$GENESIS" '.app_state.clob.clob_pairs.[5].id' -v '5'
	dasel put -t string -f "$GENESIS" '.app_state.clob.clob_pairs.[5].status' -v "$INITIAL_CLOB_PAIR_STATUS"
	dasel put -t int -f "$GENESIS" '.app_state.clob.clob_pairs.[5].perpetual_clob_metadata.perpetual_id' -v '5'
	dasel put -t int -f "$GENESIS" '.app_state.clob.clob_pairs.[5].step_base_quantums' -v '1000000'
	dasel put -t int -f "$GENESIS" '.app_state.clob.clob_pairs.[5].subticks_per_tick' -v '1000000'
	dasel put -t int -f "$GENESIS" '.app_state.clob.clob_pairs.[5].quantum_conversion_exponent' -v '-9'

	# Clob: ADA-USD
	dasel put -t json -f "$GENESIS" '.app_state.clob.clob_pairs.[]' -v "{}"
	dasel put -t int -f "$GENESIS" '.app_state.clob.clob_pairs.[6].id' -v '6'
	dasel put -t string -f "$GENESIS" '.app_state.clob.clob_pairs.[6].status' -v "$INITIAL_CLOB_PAIR_STATUS"
	dasel put -t int -f "$GENESIS" '.app_state.clob.clob_pairs.[6].perpetual_clob_metadata.perpetual_id' -v '6'
	dasel put -t int -f "$GENESIS" '.app_state.clob.clob_pairs.[6].step_base_quantums' -v '1000000'
	dasel put -t int -f "$GENESIS" '.app_state.clob.clob_pairs.[6].subticks_per_tick' -v '1000000'
	dasel put -t int -f "$GENESIS" '.app_state.clob.clob_pairs.[6].quantum_conversion_exponent' -v '-9'

	# Clob: AVAX-USD
	dasel put -t json -f "$GENESIS" '.app_state.clob.clob_pairs.[]' -v "{}"
	dasel put -t int -f "$GENESIS" '.app_state.clob.clob_pairs.[7].id' -v '7'
	dasel put -t string -f "$GENESIS" '.app_state.clob.clob_pairs.[7].status' -v "$INITIAL_CLOB_PAIR_STATUS"
	dasel put -t int -f "$GENESIS" '.app_state.clob.clob_pairs.[7].perpetual_clob_metadata.perpetual_id' -v '7'
	dasel put -t int -f "$GENESIS" '.app_state.clob.clob_pairs.[7].step_base_quantums' -v '1000000'
	dasel put -t int -f "$GENESIS" '.app_state.clob.clob_pairs.[7].subticks_per_tick' -v '1000000'
	dasel put -t int -f "$GENESIS" '.app_state.clob.clob_pairs.[7].quantum_conversion_exponent' -v '-9'

	# Clob: FIL-USD
	dasel put -t json -f "$GENESIS" '.app_state.clob.clob_pairs.[]' -v "{}"
	dasel put -t int -f "$GENESIS" '.app_state.clob.clob_pairs.[8].id' -v '8'
	dasel put -t string -f "$GENESIS" '.app_state.clob.clob_pairs.[8].status' -v "$INITIAL_CLOB_PAIR_STATUS"
	dasel put -t int -f "$GENESIS" '.app_state.clob.clob_pairs.[8].perpetual_clob_metadata.perpetual_id' -v '8'
	dasel put -t int -f "$GENESIS" '.app_state.clob.clob_pairs.[8].step_base_quantums' -v '1000000'
	dasel put -t int -f "$GENESIS" '.app_state.clob.clob_pairs.[8].subticks_per_tick' -v '1000000'
	dasel put -t int -f "$GENESIS" '.app_state.clob.clob_pairs.[8].quantum_conversion_exponent' -v '-9'

	# Clob: LTC-USD
	dasel put -t json -f "$GENESIS" '.app_state.clob.clob_pairs.[]' -v "{}"
	dasel put -t int -f "$GENESIS" '.app_state.clob.clob_pairs.[9].id' -v '9'
	dasel put -t string -f "$GENESIS" '.app_state.clob.clob_pairs.[9].status' -v "$INITIAL_CLOB_PAIR_STATUS"
	dasel put -t int -f "$GENESIS" '.app_state.clob.clob_pairs.[9].perpetual_clob_metadata.perpetual_id' -v '9'
	dasel put -t int -f "$GENESIS" '.app_state.clob.clob_pairs.[9].step_base_quantums' -v '1000000'
	dasel put -t int -f "$GENESIS" '.app_state.clob.clob_pairs.[9].subticks_per_tick' -v '1000000'
	dasel put -t int -f "$GENESIS" '.app_state.clob.clob_pairs.[9].quantum_conversion_exponent' -v '-9'

	# Clob: DOGE-USD
	dasel put -t json -f "$GENESIS" '.app_state.clob.clob_pairs.[]' -v "{}"
	dasel put -t int -f "$GENESIS" '.app_state.clob.clob_pairs.[10].id' -v '10'
	dasel put -t string -f "$GENESIS" '.app_state.clob.clob_pairs.[10].status' -v "$INITIAL_CLOB_PAIR_STATUS"
	dasel put -t int -f "$GENESIS" '.app_state.clob.clob_pairs.[10].perpetual_clob_metadata.perpetual_id' -v '10'
	dasel put -t int -f "$GENESIS" '.app_state.clob.clob_pairs.[10].step_base_quantums' -v '1000000'
	dasel put -t int -f "$GENESIS" '.app_state.clob.clob_pairs.[10].subticks_per_tick' -v '1000000'
	dasel put -t int -f "$GENESIS" '.app_state.clob.clob_pairs.[10].quantum_conversion_exponent' -v '-9'

	# Clob: ATOM-USD
	dasel put -t json -f "$GENESIS" '.app_state.clob.clob_pairs.[]' -v "{}"
	dasel put -t int -f "$GENESIS" '.app_state.clob.clob_pairs.[11].id' -v '11'
	dasel put -t string -f "$GENESIS" '.app_state.clob.clob_pairs.[11].status' -v "$INITIAL_CLOB_PAIR_STATUS"
	dasel put -t int -f "$GENESIS" '.app_state.clob.clob_pairs.[11].perpetual_clob_metadata.perpetual_id' -v '11'
	dasel put -t int -f "$GENESIS" '.app_state.clob.clob_pairs.[11].step_base_quantums' -v '1000000'
	dasel put -t int -f "$GENESIS" '.app_state.clob.clob_pairs.[11].subticks_per_tick' -v '1000000'
	dasel put -t int -f "$GENESIS" '.app_state.clob.clob_pairs.[11].quantum_conversion_exponent' -v '-9'

	# Clob: DOT-USD
	dasel put -t json -f "$GENESIS" '.app_state.clob.clob_pairs.[]' -v "{}"
	dasel put -t int -f "$GENESIS" '.app_state.clob.clob_pairs.[12].id' -v '12'
	dasel put -t string -f "$GENESIS" '.app_state.clob.clob_pairs.[12].status' -v "$INITIAL_CLOB_PAIR_STATUS"
	dasel put -t int -f "$GENESIS" '.app_state.clob.clob_pairs.[12].perpetual_clob_metadata.perpetual_id' -v '12'
	dasel put -t int -f "$GENESIS" '.app_state.clob.clob_pairs.[12].step_base_quantums' -v '1000000'
	dasel put -t int -f "$GENESIS" '.app_state.clob.clob_pairs.[12].subticks_per_tick' -v '1000000'
	dasel put -t int -f "$GENESIS" '.app_state.clob.clob_pairs.[12].quantum_conversion_exponent' -v '-9'

	# Clob: UNI-USD
	dasel put -t json -f "$GENESIS" '.app_state.clob.clob_pairs.[]' -v "{}"
	dasel put -t int -f "$GENESIS" '.app_state.clob.clob_pairs.[13].id' -v '13'
	dasel put -t string -f "$GENESIS" '.app_state.clob.clob_pairs.[13].status' -v "$INITIAL_CLOB_PAIR_STATUS"
	dasel put -t int -f "$GENESIS" '.app_state.clob.clob_pairs.[13].perpetual_clob_metadata.perpetual_id' -v '13'
	dasel put -t int -f "$GENESIS" '.app_state.clob.clob_pairs.[13].step_base_quantums' -v '1000000'
	dasel put -t int -f "$GENESIS" '.app_state.clob.clob_pairs.[13].subticks_per_tick' -v '1000000'
	dasel put -t int -f "$GENESIS" '.app_state.clob.clob_pairs.[13].quantum_conversion_exponent' -v '-9'

	# Clob: BCH-USD
	dasel put -t json -f "$GENESIS" '.app_state.clob.clob_pairs.[]' -v "{}"
	dasel put -t int -f "$GENESIS" '.app_state.clob.clob_pairs.[14].id' -v '14'
	dasel put -t string -f "$GENESIS" '.app_state.clob.clob_pairs.[14].status' -v "$INITIAL_CLOB_PAIR_STATUS"
	dasel put -t int -f "$GENESIS" '.app_state.clob.clob_pairs.[14].perpetual_clob_metadata.perpetual_id' -v '14'
	dasel put -t int -f "$GENESIS" '.app_state.clob.clob_pairs.[14].step_base_quantums' -v '1000000'
	dasel put -t int -f "$GENESIS" '.app_state.clob.clob_pairs.[14].subticks_per_tick' -v '1000000'
	dasel put -t int -f "$GENESIS" '.app_state.clob.clob_pairs.[14].quantum_conversion_exponent' -v '-9'

	# Clob: TRX-USD
	dasel put -t json -f "$GENESIS" '.app_state.clob.clob_pairs.[]' -v "{}"
	dasel put -t int -f "$GENESIS" '.app_state.clob.clob_pairs.[15].id' -v '15'
	dasel put -t string -f "$GENESIS" '.app_state.clob.clob_pairs.[15].status' -v "$INITIAL_CLOB_PAIR_STATUS"
	dasel put -t int -f "$GENESIS" '.app_state.clob.clob_pairs.[15].perpetual_clob_metadata.perpetual_id' -v '15'
	dasel put -t int -f "$GENESIS" '.app_state.clob.clob_pairs.[15].step_base_quantums' -v '1000000'
	dasel put -t int -f "$GENESIS" '.app_state.clob.clob_pairs.[15].subticks_per_tick' -v '1000000'
	dasel put -t int -f "$GENESIS" '.app_state.clob.clob_pairs.[15].quantum_conversion_exponent' -v '-9'

	# Clob: NEAR-USD
	dasel put -t json -f "$GENESIS" '.app_state.clob.clob_pairs.[]' -v "{}"
	dasel put -t int -f "$GENESIS" '.app_state.clob.clob_pairs.[16].id' -v '16'
	dasel put -t string -f "$GENESIS" '.app_state.clob.clob_pairs.[16].status' -v "$INITIAL_CLOB_PAIR_STATUS"
	dasel put -t int -f "$GENESIS" '.app_state.clob.clob_pairs.[16].perpetual_clob_metadata.perpetual_id' -v '16'
	dasel put -t int -f "$GENESIS" '.app_state.clob.clob_pairs.[16].step_base_quantums' -v '1000000'
	dasel put -t int -f "$GENESIS" '.app_state.clob.clob_pairs.[16].subticks_per_tick' -v '1000000'
	dasel put -t int -f "$GENESIS" '.app_state.clob.clob_pairs.[16].quantum_conversion_exponent' -v '-9'

	# Clob: MKR-USD
	dasel put -t json -f "$GENESIS" '.app_state.clob.clob_pairs.[]' -v "{}"
	dasel put -t int -f "$GENESIS" '.app_state.clob.clob_pairs.[17].id' -v '17'
	dasel put -t string -f "$GENESIS" '.app_state.clob.clob_pairs.[17].status' -v "$INITIAL_CLOB_PAIR_STATUS"
	dasel put -t int -f "$GENESIS" '.app_state.clob.clob_pairs.[17].perpetual_clob_metadata.perpetual_id' -v '17'
	dasel put -t int -f "$GENESIS" '.app_state.clob.clob_pairs.[17].step_base_quantums' -v '1000000'
	dasel put -t int -f "$GENESIS" '.app_state.clob.clob_pairs.[17].subticks_per_tick' -v '1000000'
	dasel put -t int -f "$GENESIS" '.app_state.clob.clob_pairs.[17].quantum_conversion_exponent' -v '-9'

	# Clob: XLM-USD
	dasel put -t json -f "$GENESIS" '.app_state.clob.clob_pairs.[]' -v "{}"
	dasel put -t int -f "$GENESIS" '.app_state.clob.clob_pairs.[18].id' -v '18'
	dasel put -t string -f "$GENESIS" '.app_state.clob.clob_pairs.[18].status' -v "$INITIAL_CLOB_PAIR_STATUS"
	dasel put -t int -f "$GENESIS" '.app_state.clob.clob_pairs.[18].perpetual_clob_metadata.perpetual_id' -v '18'
	dasel put -t int -f "$GENESIS" '.app_state.clob.clob_pairs.[18].step_base_quantums' -v '1000000'
	dasel put -t int -f "$GENESIS" '.app_state.clob.clob_pairs.[18].subticks_per_tick' -v '1000000'
	dasel put -t int -f "$GENESIS" '.app_state.clob.clob_pairs.[18].quantum_conversion_exponent' -v '-9'

	# Clob: ETC-USD
	dasel put -t json -f "$GENESIS" '.app_state.clob.clob_pairs.[]' -v "{}"
	dasel put -t int -f "$GENESIS" '.app_state.clob.clob_pairs.[19].id' -v '19'
	dasel put -t string -f "$GENESIS" '.app_state.clob.clob_pairs.[19].status' -v "$INITIAL_CLOB_PAIR_STATUS"
	dasel put -t int -f "$GENESIS" '.app_state.clob.clob_pairs.[19].perpetual_clob_metadata.perpetual_id' -v '19'
	dasel put -t int -f "$GENESIS" '.app_state.clob.clob_pairs.[19].step_base_quantums' -v '1000000'
	dasel put -t int -f "$GENESIS" '.app_state.clob.clob_pairs.[19].subticks_per_tick' -v '1000000'
	dasel put -t int -f "$GENESIS" '.app_state.clob.clob_pairs.[19].quantum_conversion_exponent' -v '-9'

	# Clob: COMP-USD
	dasel put -t json -f "$GENESIS" '.app_state.clob.clob_pairs.[]' -v "{}"
	dasel put -t int -f "$GENESIS" '.app_state.clob.clob_pairs.[20].id' -v '20'
	dasel put -t string -f "$GENESIS" '.app_state.clob.clob_pairs.[20].status' -v "$INITIAL_CLOB_PAIR_STATUS"
	dasel put -t int -f "$GENESIS" '.app_state.clob.clob_pairs.[20].perpetual_clob_metadata.perpetual_id' -v '20'
	dasel put -t int -f "$GENESIS" '.app_state.clob.clob_pairs.[20].step_base_quantums' -v '1000000'
	dasel put -t int -f "$GENESIS" '.app_state.clob.clob_pairs.[20].subticks_per_tick' -v '1000000'
	dasel put -t int -f "$GENESIS" '.app_state.clob.clob_pairs.[20].quantum_conversion_exponent' -v '-9'

	# Clob: WLD-USD
	dasel put -t json -f "$GENESIS" '.app_state.clob.clob_pairs.[]' -v "{}"
	dasel put -t int -f "$GENESIS" '.app_state.clob.clob_pairs.[21].id' -v '21'
	dasel put -t string -f "$GENESIS" '.app_state.clob.clob_pairs.[21].status' -v "$INITIAL_CLOB_PAIR_STATUS"
	dasel put -t int -f "$GENESIS" '.app_state.clob.clob_pairs.[21].perpetual_clob_metadata.perpetual_id' -v '21'
	dasel put -t int -f "$GENESIS" '.app_state.clob.clob_pairs.[21].step_base_quantums' -v '1000000'
	dasel put -t int -f "$GENESIS" '.app_state.clob.clob_pairs.[21].subticks_per_tick' -v '1000000'
	dasel put -t int -f "$GENESIS" '.app_state.clob.clob_pairs.[21].quantum_conversion_exponent' -v '-9'

	# Clob: APE-USD
	dasel put -t json -f "$GENESIS" '.app_state.clob.clob_pairs.[]' -v "{}"
	dasel put -t int -f "$GENESIS" '.app_state.clob.clob_pairs.[22].id' -v '22'
	dasel put -t string -f "$GENESIS" '.app_state.clob.clob_pairs.[22].status' -v "$INITIAL_CLOB_PAIR_STATUS"
	dasel put -t int -f "$GENESIS" '.app_state.clob.clob_pairs.[22].perpetual_clob_metadata.perpetual_id' -v '22'
	dasel put -t int -f "$GENESIS" '.app_state.clob.clob_pairs.[22].step_base_quantums' -v '1000000'
	dasel put -t int -f "$GENESIS" '.app_state.clob.clob_pairs.[22].subticks_per_tick' -v '1000000'
	dasel put -t int -f "$GENESIS" '.app_state.clob.clob_pairs.[22].quantum_conversion_exponent' -v '-9'

	# Clob: APT-USD
	dasel put -t json -f "$GENESIS" '.app_state.clob.clob_pairs.[]' -v "{}"
	dasel put -t int -f "$GENESIS" '.app_state.clob.clob_pairs.[23].id' -v '23'
	dasel put -t string -f "$GENESIS" '.app_state.clob.clob_pairs.[23].status' -v "$INITIAL_CLOB_PAIR_STATUS"
	dasel put -t int -f "$GENESIS" '.app_state.clob.clob_pairs.[23].perpetual_clob_metadata.perpetual_id' -v '23'
	dasel put -t int -f "$GENESIS" '.app_state.clob.clob_pairs.[23].step_base_quantums' -v '1000000'
	dasel put -t int -f "$GENESIS" '.app_state.clob.clob_pairs.[23].subticks_per_tick' -v '1000000'
	dasel put -t int -f "$GENESIS" '.app_state.clob.clob_pairs.[23].quantum_conversion_exponent' -v '-9'

	# Clob: ARB-USD
	dasel put -t json -f "$GENESIS" '.app_state.clob.clob_pairs.[]' -v "{}"
	dasel put -t int -f "$GENESIS" '.app_state.clob.clob_pairs.[24].id' -v '24'
	dasel put -t string -f "$GENESIS" '.app_state.clob.clob_pairs.[24].status' -v "$INITIAL_CLOB_PAIR_STATUS"
	dasel put -t int -f "$GENESIS" '.app_state.clob.clob_pairs.[24].perpetual_clob_metadata.perpetual_id' -v '24'
	dasel put -t int -f "$GENESIS" '.app_state.clob.clob_pairs.[24].step_base_quantums' -v '1000000'
	dasel put -t int -f "$GENESIS" '.app_state.clob.clob_pairs.[24].subticks_per_tick' -v '1000000'
	dasel put -t int -f "$GENESIS" '.app_state.clob.clob_pairs.[24].quantum_conversion_exponent' -v '-9'

	# Clob: BLUR-USD
	dasel put -t json -f "$GENESIS" '.app_state.clob.clob_pairs.[]' -v "{}"
	dasel put -t int -f "$GENESIS" '.app_state.clob.clob_pairs.[25].id' -v '25'
	dasel put -t string -f "$GENESIS" '.app_state.clob.clob_pairs.[25].status' -v "$INITIAL_CLOB_PAIR_STATUS"
	dasel put -t int -f "$GENESIS" '.app_state.clob.clob_pairs.[25].perpetual_clob_metadata.perpetual_id' -v '25'
	dasel put -t int -f "$GENESIS" '.app_state.clob.clob_pairs.[25].step_base_quantums' -v '1000000'
	dasel put -t int -f "$GENESIS" '.app_state.clob.clob_pairs.[25].subticks_per_tick' -v '1000000'
	dasel put -t int -f "$GENESIS" '.app_state.clob.clob_pairs.[25].quantum_conversion_exponent' -v '-9'

	# Clob: LDO-USD
	dasel put -t json -f "$GENESIS" '.app_state.clob.clob_pairs.[]' -v "{}"
	dasel put -t int -f "$GENESIS" '.app_state.clob.clob_pairs.[26].id' -v '26'
	dasel put -t string -f "$GENESIS" '.app_state.clob.clob_pairs.[26].status' -v "$INITIAL_CLOB_PAIR_STATUS"
	dasel put -t int -f "$GENESIS" '.app_state.clob.clob_pairs.[26].perpetual_clob_metadata.perpetual_id' -v '26'
	dasel put -t int -f "$GENESIS" '.app_state.clob.clob_pairs.[26].step_base_quantums' -v '1000000'
	dasel put -t int -f "$GENESIS" '.app_state.clob.clob_pairs.[26].subticks_per_tick' -v '1000000'
	dasel put -t int -f "$GENESIS" '.app_state.clob.clob_pairs.[26].quantum_conversion_exponent' -v '-9'

	# Clob: OP-USD
	dasel put -t json -f "$GENESIS" '.app_state.clob.clob_pairs.[]' -v "{}"
	dasel put -t int -f "$GENESIS" '.app_state.clob.clob_pairs.[27].id' -v '27'
	dasel put -t string -f "$GENESIS" '.app_state.clob.clob_pairs.[27].status' -v "$INITIAL_CLOB_PAIR_STATUS"
	dasel put -t int -f "$GENESIS" '.app_state.clob.clob_pairs.[27].perpetual_clob_metadata.perpetual_id' -v '27'
	dasel put -t int -f "$GENESIS" '.app_state.clob.clob_pairs.[27].step_base_quantums' -v '1000000'
	dasel put -t int -f "$GENESIS" '.app_state.clob.clob_pairs.[27].subticks_per_tick' -v '1000000'
	dasel put -t int -f "$GENESIS" '.app_state.clob.clob_pairs.[27].quantum_conversion_exponent' -v '-9'

	# Clob: PEPE-USD
	dasel put -t json -f "$GENESIS" '.app_state.clob.clob_pairs.[]' -v "{}"
	dasel put -t int -f "$GENESIS" '.app_state.clob.clob_pairs.[28].id' -v '28'
	dasel put -t string -f "$GENESIS" '.app_state.clob.clob_pairs.[28].status' -v "$INITIAL_CLOB_PAIR_STATUS"
	dasel put -t int -f "$GENESIS" '.app_state.clob.clob_pairs.[28].perpetual_clob_metadata.perpetual_id' -v '28'
	dasel put -t int -f "$GENESIS" '.app_state.clob.clob_pairs.[28].step_base_quantums' -v '1000000'
	dasel put -t int -f "$GENESIS" '.app_state.clob.clob_pairs.[28].subticks_per_tick' -v '1000000'
	dasel put -t int -f "$GENESIS" '.app_state.clob.clob_pairs.[28].quantum_conversion_exponent' -v '-9'

	# Clob: SEI-USD
	dasel put -t json -f "$GENESIS" '.app_state.clob.clob_pairs.[]' -v "{}"
	dasel put -t int -f "$GENESIS" '.app_state.clob.clob_pairs.[29].id' -v '29'
	dasel put -t string -f "$GENESIS" '.app_state.clob.clob_pairs.[29].status' -v "$INITIAL_CLOB_PAIR_STATUS"
	dasel put -t int -f "$GENESIS" '.app_state.clob.clob_pairs.[29].perpetual_clob_metadata.perpetual_id' -v '29'
	dasel put -t int -f "$GENESIS" '.app_state.clob.clob_pairs.[29].step_base_quantums' -v '1000000'
	dasel put -t int -f "$GENESIS" '.app_state.clob.clob_pairs.[29].subticks_per_tick' -v '1000000'
	dasel put -t int -f "$GENESIS" '.app_state.clob.clob_pairs.[29].quantum_conversion_exponent' -v '-9'

	# Clob: SHIB-USD
	dasel put -t json -f "$GENESIS" '.app_state.clob.clob_pairs.[]' -v "{}"
	dasel put -t int -f "$GENESIS" '.app_state.clob.clob_pairs.[30].id' -v '30'
	dasel put -t string -f "$GENESIS" '.app_state.clob.clob_pairs.[30].status' -v "$INITIAL_CLOB_PAIR_STATUS"
	dasel put -t int -f "$GENESIS" '.app_state.clob.clob_pairs.[30].perpetual_clob_metadata.perpetual_id' -v '30'
	dasel put -t int -f "$GENESIS" '.app_state.clob.clob_pairs.[30].step_base_quantums' -v '1000000'
	dasel put -t int -f "$GENESIS" '.app_state.clob.clob_pairs.[30].subticks_per_tick' -v '1000000'
	dasel put -t int -f "$GENESIS" '.app_state.clob.clob_pairs.[30].quantum_conversion_exponent' -v '-9'

	# Clob: SUI-USD
	dasel put -t json -f "$GENESIS" '.app_state.clob.clob_pairs.[]' -v "{}"
	dasel put -t int -f "$GENESIS" '.app_state.clob.clob_pairs.[31].id' -v '31'
	dasel put -t string -f "$GENESIS" '.app_state.clob.clob_pairs.[31].status' -v "$INITIAL_CLOB_PAIR_STATUS"
	dasel put -t int -f "$GENESIS" '.app_state.clob.clob_pairs.[31].perpetual_clob_metadata.perpetual_id' -v '31'
	dasel put -t int -f "$GENESIS" '.app_state.clob.clob_pairs.[31].step_base_quantums' -v '1000000'
	dasel put -t int -f "$GENESIS" '.app_state.clob.clob_pairs.[31].subticks_per_tick' -v '1000000'
	dasel put -t int -f "$GENESIS" '.app_state.clob.clob_pairs.[31].quantum_conversion_exponent' -v '-9'

	# Clob: XRP-USD
	dasel put -t json -f "$GENESIS" '.app_state.clob.clob_pairs.[]' -v "{}"
	dasel put -t int -f "$GENESIS" '.app_state.clob.clob_pairs.[32].id' -v '32'
	dasel put -t string -f "$GENESIS" '.app_state.clob.clob_pairs.[32].status' -v "$INITIAL_CLOB_PAIR_STATUS"
	dasel put -t int -f "$GENESIS" '.app_state.clob.clob_pairs.[32].perpetual_clob_metadata.perpetual_id' -v '32'
	dasel put -t int -f "$GENESIS" '.app_state.clob.clob_pairs.[32].step_base_quantums' -v '1000000'
	dasel put -t int -f "$GENESIS" '.app_state.clob.clob_pairs.[32].subticks_per_tick' -v '1000000'
	dasel put -t int -f "$GENESIS" '.app_state.clob.clob_pairs.[32].quantum_conversion_exponent' -v '-9'

	# Liquidations
	dasel put -t int -f "$GENESIS" '.app_state.clob.liquidations_config.max_liquidation_fee_ppm' -v '15000'  # 1.5%
	dasel put -t int -f "$GENESIS" '.app_state.clob.liquidations_config.position_block_limits.min_position_notional_liquidated' -v '1000000000' # 1_000 USDC
	dasel put -t int -f "$GENESIS" '.app_state.clob.liquidations_config.position_block_limits.max_position_portion_liquidated_ppm' -v '100000'  # 10%
	dasel put -t int -f "$GENESIS" '.app_state.clob.liquidations_config.subaccount_block_limits.max_notional_liquidated' -v '100000000000'  # 100_000 USDC
	dasel put -t int -f "$GENESIS" '.app_state.clob.liquidations_config.subaccount_block_limits.max_quantums_insurance_lost' -v '1000000000000' # 1_000_000 USDC
	dasel put -t int -f "$GENESIS" '.app_state.clob.liquidations_config.fillable_price_config.bankruptcy_adjustment_ppm' -v '1000000'
	dasel put -t int -f "$GENESIS" '.app_state.clob.liquidations_config.fillable_price_config.spread_to_maintenance_margin_ratio_ppm' -v '1500000'  # 150%

	# Block Rate Limit
	# Max 50 short term orders per block
	dasel put -t json -f "$GENESIS" '.app_state.clob.block_rate_limit_config.max_short_term_orders_per_n_blocks.[]' -v "{}"
	dasel put -t int -f "$GENESIS" '.app_state.clob.block_rate_limit_config.max_short_term_orders_per_n_blocks.[0].limit' -v '200'
	dasel put -t int -f "$GENESIS" '.app_state.clob.block_rate_limit_config.max_short_term_orders_per_n_blocks.[0].num_blocks' -v '1'
	# Max 50 short term order cancellations per block
	dasel put -t json -f "$GENESIS" '.app_state.clob.block_rate_limit_config.max_short_term_order_cancellations_per_n_blocks.[]' -v "{}"
	dasel put -t int -f "$GENESIS" '.app_state.clob.block_rate_limit_config.max_short_term_order_cancellations_per_n_blocks.[0].limit' -v '200'
	dasel put -t int -f "$GENESIS" '.app_state.clob.block_rate_limit_config.max_short_term_order_cancellations_per_n_blocks.[0].num_blocks' -v '1'
	# Max 2 stateful orders per block
	dasel put -t json -f "$GENESIS" '.app_state.clob.block_rate_limit_config.max_stateful_orders_per_n_blocks.[]' -v "{}"
	dasel put -t int -f "$GENESIS" '.app_state.clob.block_rate_limit_config.max_stateful_orders_per_n_blocks.[0].limit' -v '2'
	dasel put -t int -f "$GENESIS" '.app_state.clob.block_rate_limit_config.max_stateful_orders_per_n_blocks.[0].num_blocks' -v '1'
	# Max 20 stateful orders per 100 blocks
	dasel put -t json -f "$GENESIS" '.app_state.clob.block_rate_limit_config.max_stateful_orders_per_n_blocks.[]' -v "{}"
	dasel put -t int -f "$GENESIS" '.app_state.clob.block_rate_limit_config.max_stateful_orders_per_n_blocks.[1].limit' -v '20'
	dasel put -t int -f "$GENESIS" '.app_state.clob.block_rate_limit_config.max_stateful_orders_per_n_blocks.[1].num_blocks' -v '100'

	# Equity Tier Limit
	# Max 0 open short term orders for $0 USDC TNC
	dasel put -t json -f "$GENESIS" '.app_state.clob.equity_tier_limit_config.short_term_order_equity_tiers.[]' -v "{}"
	dasel put -t int -f "$GENESIS" '.app_state.clob.equity_tier_limit_config.short_term_order_equity_tiers.[0].limit' -v '0'
	dasel put -t string -f "$GENESIS" '.app_state.clob.equity_tier_limit_config.short_term_order_equity_tiers.[0].usd_tnc_required' -v '0'
	# Max 1 open short term orders for $20 USDC TNC
	dasel put -t json -f "$GENESIS" '.app_state.clob.equity_tier_limit_config.short_term_order_equity_tiers.[]' -v "{}"
	dasel put -t int -f "$GENESIS" '.app_state.clob.equity_tier_limit_config.short_term_order_equity_tiers.[1].limit' -v '1'
	dasel put -t string -f "$GENESIS" '.app_state.clob.equity_tier_limit_config.short_term_order_equity_tiers.[1].usd_tnc_required' -v '20000000'
	# Max 5 open short term orders for $100 USDC TNC
	dasel put -t json -f "$GENESIS" '.app_state.clob.equity_tier_limit_config.short_term_order_equity_tiers.[]' -v "{}"
	dasel put -t int -f "$GENESIS" '.app_state.clob.equity_tier_limit_config.short_term_order_equity_tiers.[2].limit' -v '5'
	dasel put -t string -f "$GENESIS" '.app_state.clob.equity_tier_limit_config.short_term_order_equity_tiers.[2].usd_tnc_required' -v '100000000'
	# Max 10 open short term orders for $1000 USDC TNC
	dasel put -t json -f "$GENESIS" '.app_state.clob.equity_tier_limit_config.short_term_order_equity_tiers.[]' -v "{}"
	dasel put -t int -f "$GENESIS" '.app_state.clob.equity_tier_limit_config.short_term_order_equity_tiers.[3].limit' -v '10'
	dasel put -t string -f "$GENESIS" '.app_state.clob.equity_tier_limit_config.short_term_order_equity_tiers.[3].usd_tnc_required' -v '1000000000'
	# Max 100 open short term orders for $10,000 USDC TNC
	dasel put -t json -f "$GENESIS" '.app_state.clob.equity_tier_limit_config.short_term_order_equity_tiers.[]' -v "{}"
	dasel put -t int -f "$GENESIS" '.app_state.clob.equity_tier_limit_config.short_term_order_equity_tiers.[4].limit' -v '100'
	dasel put -t string -f "$GENESIS" '.app_state.clob.equity_tier_limit_config.short_term_order_equity_tiers.[4].usd_tnc_required' -v '10000000000'
	# Max 200 open short term orders for $100,000 USDC TNC
	dasel put -t json -f "$GENESIS" '.app_state.clob.equity_tier_limit_config.short_term_order_equity_tiers.[]' -v "{}"
	dasel put -t int -f "$GENESIS" '.app_state.clob.equity_tier_limit_config.short_term_order_equity_tiers.[5].limit' -v '1000'
	dasel put -t string -f "$GENESIS" '.app_state.clob.equity_tier_limit_config.short_term_order_equity_tiers.[5].usd_tnc_required' -v '100000000000'
	# Max 0 open stateful orders for $0 USDC TNC
	dasel put -t json -f "$GENESIS" '.app_state.clob.equity_tier_limit_config.stateful_order_equity_tiers.[]' -v "{}"
	dasel put -t int -f "$GENESIS" '.app_state.clob.equity_tier_limit_config.stateful_order_equity_tiers.[0].limit' -v '0'
	dasel put -t string -f "$GENESIS" '.app_state.clob.equity_tier_limit_config.stateful_order_equity_tiers.[0].usd_tnc_required' -v '0'
	# Max 1 open stateful orders for $20 USDC TNC
	dasel put -t json -f "$GENESIS" '.app_state.clob.equity_tier_limit_config.stateful_order_equity_tiers.[]' -v "{}"
	dasel put -t int -f "$GENESIS" '.app_state.clob.equity_tier_limit_config.stateful_order_equity_tiers.[1].limit' -v '1'
	dasel put -t string -f "$GENESIS" '.app_state.clob.equity_tier_limit_config.stateful_order_equity_tiers.[1].usd_tnc_required' -v '20000000'
	# Max 5 open stateful orders for $100 USDC TNC
	dasel put -t json -f "$GENESIS" '.app_state.clob.equity_tier_limit_config.stateful_order_equity_tiers.[]' -v "{}"
	dasel put -t int -f "$GENESIS" '.app_state.clob.equity_tier_limit_config.stateful_order_equity_tiers.[2].limit' -v '5'
	dasel put -t string -f "$GENESIS" '.app_state.clob.equity_tier_limit_config.stateful_order_equity_tiers.[2].usd_tnc_required' -v '100000000'
	# Max 10 open stateful orders for $1000 USDC TNC
	dasel put -t json -f "$GENESIS" '.app_state.clob.equity_tier_limit_config.stateful_order_equity_tiers.[]' -v "{}"
	dasel put -t int -f "$GENESIS" '.app_state.clob.equity_tier_limit_config.stateful_order_equity_tiers.[3].limit' -v '10'
	dasel put -t string -f "$GENESIS" '.app_state.clob.equity_tier_limit_config.stateful_order_equity_tiers.[3].usd_tnc_required' -v '1000000000'
	# Max 100 open stateful orders for $10,000 USDC TNC
	dasel put -t json -f "$GENESIS" '.app_state.clob.equity_tier_limit_config.stateful_order_equity_tiers.[]' -v "{}"
	dasel put -t int -f "$GENESIS" '.app_state.clob.equity_tier_limit_config.stateful_order_equity_tiers.[4].limit' -v '100'
	dasel put -t string -f "$GENESIS" '.app_state.clob.equity_tier_limit_config.stateful_order_equity_tiers.[4].usd_tnc_required' -v '10000000000'
	# Max 200 open stateful orders for $100,000 USDC TNC
	dasel put -t json -f "$GENESIS" '.app_state.clob.equity_tier_limit_config.stateful_order_equity_tiers.[]' -v "{}"
	dasel put -t int -f "$GENESIS" '.app_state.clob.equity_tier_limit_config.stateful_order_equity_tiers.[5].limit' -v '200'
	dasel put -t string -f "$GENESIS" '.app_state.clob.equity_tier_limit_config.stateful_order_equity_tiers.[5].usd_tnc_required' -v '100000000000'


  # Fee Tiers
  # Schedule a delayed message to swap fee tiers to the standard schedule after ~120 days of blocks.
	dasel put -t int -f "$GENESIS" '.app_state.delaymsg.next_delayed_message_id' -v '1'
	dasel put -t json -f "$GENESIS" '.app_state.delaymsg.delayed_messages.[]' -v "{}"
	dasel put -t int -f "$GENESIS" '.app_state.delaymsg.delayed_messages.[0].id' -v '0'

	delaymsg=$(cat "$DELAY_MSG_JSON_DIR/perpetual_fee_params_msg.json" | jq -c '.')
	dasel put -t json -f "$GENESIS" '.app_state.delaymsg.delayed_messages.[0].msg' -v "$delaymsg"
	# Schedule the message to execute in ~7 days (at 1.6s per block.)
	dasel put -t int -f "$GENESIS" '.app_state.delaymsg.delayed_messages.[0].block_height' -v '378000'
	# Uncomment the following to schedule the message to execute in ~120 days (at 1.6s per block.)
	# dasel put -t int -f "$GENESIS" '.app_state.delaymsg.delayed_messages.[0].block_height' -v '6480000'
}

function add_subaccount() {
	GEN_FILE=$1
	IDX=$2
	ACCOUNT=$3
	QUOTE_BALANCE=$4

	dasel put -t json -f "$GEN_FILE" ".app_state.subaccounts.subaccounts.[]" -v "{}"
	dasel put -t json -f "$GEN_FILE" ".app_state.subaccounts.subaccounts.[$IDX].id" -v "{}"
	dasel put -t string -f "$GEN_FILE" ".app_state.subaccounts.subaccounts.[$IDX].id.owner" -v "$ACCOUNT"
	dasel put -t int -f "$GEN_FILE" ".app_state.subaccounts.subaccounts.[$IDX].id.number" -v '0'
	dasel put -t bool -f "$GEN_FILE" ".app_state.subaccounts.subaccounts.[$IDX].margin_enabled" -v 'true'

	dasel put -t json -f "$GEN_FILE" ".app_state.subaccounts.subaccounts.[$IDX].asset_positions.[]" -v '{}'
	dasel put -t int -f "$GEN_FILE" ".app_state.subaccounts.subaccounts.[$IDX].asset_positions.[0].asset_id" -v '0'
	dasel put -t string -f "$GEN_FILE" ".app_state.subaccounts.subaccounts.[$IDX].asset_positions.[0].quantums" -v "${QUOTE_BALANCE}"
	dasel put -t int -f "$GEN_FILE" ".app_state.subaccounts.subaccounts.[$IDX].asset_positions.[0].index" -v '0'
}

# Modify the genesis file to use only the test exchange for computing index prices. The test exchange is configured
# to serve prices for BTC, ETH and LINK. This must be called after edit_genesis to ensure all markets exist.
function update_genesis_use_test_exchange() {
	GENESIS=$1/genesis.json

	# For BTC, ETH and LINK, remove all exchanges except the test exchange.
	btc_exchange_config_json=$(cat <<-EOF
	{
		"exchanges": [
			{
				"exchangeName": "TestExchange",
				"ticker": "BTC-USD"
			}
		]
	}
	EOF
	)
	dasel put -t string -f "$GENESIS" '.app_state.prices.market_params.[0].exchange_config_json' -v "$btc_exchange_config_json"
	dasel put -t int -f "$GENESIS" '.app_state.prices.market_params.[0].min_exchanges' -v '1'

	eth_exchange_config_json=$(cat <<-EOF
	{
		"exchanges": [
			{
				"exchangeName": "TestExchange",
				"ticker": "ETH-USD"
			}
		]
	}
	EOF
	)
	dasel put -t string -f "$GENESIS" '.app_state.prices.market_params.[1].exchange_config_json' -v "$eth_exchange_config_json"
	dasel put -t int -f "$GENESIS" '.app_state.prices.market_params.[1].min_exchanges' -v '1'

	link_exchange_config_json=$(cat <<-EOF
	{
		"exchanges": [
			{
				"exchangeName": "TestExchange",
				"ticker": "LINK-USD"
			}
		]
	}
	EOF
	)
	dasel put -t string -f "$GENESIS" '.app_state.prices.market_params.[2].exchange_config_json' -v "$link_exchange_config_json"
	dasel put -t int -f "$GENESIS" '.app_state.prices.market_params.[2].min_exchanges' -v '1'

	# All remaining markets can just use the LINK ticker so the daemon will start. All markets must have at least 1
	# exchange. With only one exchange configured, there should not be enough prices to meet the minimum exchange
	# count, and these markets will not have index prices.
	for market_idx in {3..34}
	do
		dasel put -t string -f "$GENESIS" ".app_state.prices.market_params.[$market_idx].exchange_config_json" -v "$link_exchange_config_json"
	done
}

# Modify the genesis file to add test volatile market. Market TEST-USD will be added as market 33.
function update_genesis_use_test_volatile_market() {
	GENESIS=$1/genesis.json
	TEST_USD_MARKET_ID=33

	# Market: TEST-USD
	dasel put -t json -f "$GENESIS" '.app_state.prices.market_params.[]' -v "{}"
	dasel put -t string -f "$GENESIS" '.app_state.prices.market_params.last().pair' -v 'TEST-USD'
	dasel put -t int -f "$GENESIS" '.app_state.prices.market_params.last().id' -v "${TEST_USD_MARKET_ID}"
	dasel put -t int -f "$GENESIS" '.app_state.prices.market_params.last().exponent' -v '-5'
	dasel put -t int -f "$GENESIS" '.app_state.prices.market_params.last().min_exchanges' -v '1'
	dasel put -t int -f "$GENESIS" '.app_state.prices.market_params.last().min_price_change_ppm' -v '250' # 0.025%
	dasel put -t json -f "$GENESIS" '.app_state.prices.market_prices.[]' -v "{}"
	dasel put -t int -f "$GENESIS" '.app_state.prices.market_prices.last().id' -v "${TEST_USD_MARKET_ID}"
	dasel put -t int -f "$GENESIS" '.app_state.prices.market_prices.last().exponent' -v '-5'
	dasel put -t int -f "$GENESIS" '.app_state.prices.market_prices.last().price' -v '10000000'          # $100 = 1 TEST.
	# TEST Exchange Config
	test_exchange_config_json=$(cat "$EXCHANGE_CONFIG_JSON_DIR/test_exchange_config.json" | jq -c '.')
	dasel put -t string -f "$GENESIS" '.app_state.prices.market_params.last().exchange_config_json' -v "$test_exchange_config_json"

	# Liquidity Tier: For TEST-USD. 1% leverage and regular 1m nonlinear margin thresholds.
	NUM_LIQUIDITY_TIERS=$(jq -c '.app_state.perpetuals.liquidity_tiers | length' < ${GENESIS})
	dasel put -t json -f "$GENESIS" '.app_state.perpetuals.liquidity_tiers.[]' -v "{}"
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.liquidity_tiers.last().id' -v "${NUM_LIQUIDITY_TIERS}"
	dasel put -t string -f "$GENESIS" '.app_state.perpetuals.liquidity_tiers.last().name' -v 'test-usd-100x-liq-tier-linear'
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.liquidity_tiers.last().initial_margin_ppm' -v '10007' # 1% + a little prime (100x leverage)
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.liquidity_tiers.last().maintenance_fraction_ppm' -v '500009' # 50% of IM + a little prime
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.liquidity_tiers.last().base_position_notional' -v '1000000000039' # 1_000_000 USDC + a little prime
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.liquidity_tiers.last().impact_notional' -v '50000000000' # 50_000 USDC (500 USDC / 1%)

	# Liquidity Tier: For TEST-USD. 1% leverage and 100 nonlinear margin thresholds.
	NUM_LIQUIDITY_TIERS_2=$(jq -c '.app_state.perpetuals.liquidity_tiers | length' < ${GENESIS})
	dasel put -t json -f "$GENESIS" '.app_state.perpetuals.liquidity_tiers.[]' -v "{}"
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.liquidity_tiers.last().id' -v "${NUM_LIQUIDITY_TIERS_2}"
	dasel put -t string -f "$GENESIS" '.app_state.perpetuals.liquidity_tiers.last().name' -v 'test-usd-100x-liq-tier-nonlinear'
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.liquidity_tiers.last().initial_margin_ppm' -v '10007' # 1% + a little prime (100x leverage)
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.liquidity_tiers.last().maintenance_fraction_ppm' -v '500009' # 50% of IM + a little prime
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.liquidity_tiers.last().base_position_notional' -v '100000007' # 100 USDC + a little prime
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.liquidity_tiers.last().impact_notional' -v '50000000000' # 50_000 USDC (500 USDC / 1%)

	# Perpetual: TEST-USD
	NUM_PERPETUALS=$(jq -c '.app_state.perpetuals.perpetuals | length' < ${GENESIS})
	dasel put -t json -f "$GENESIS" '.app_state.perpetuals.perpetuals.[]' -v "{}"
	dasel put -t string -f "$GENESIS" '.app_state.perpetuals.perpetuals.last().params.ticker' -v 'TEST-USD'
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.perpetuals.last().params.id' -v "${NUM_PERPETUALS}"
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.perpetuals.last().params.market_id' -v "${TEST_USD_MARKET_ID}"
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.perpetuals.last().params.atomic_resolution' -v '-10'
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.perpetuals.last().params.default_funding_ppm' -v '0'
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.perpetuals.last().params.liquidity_tier' -v "${NUM_LIQUIDITY_TIERS}"

	# Clob: TEST-USD
	NUM_CLOB_PAIRS=$(jq -c '.app_state.clob.clob_pairs | length' < ${GENESIS})
	dasel put -t json -f "$GENESIS" '.app_state.clob.clob_pairs.[]' -v "{}"
	dasel put -t int -f "$GENESIS" '.app_state.clob.clob_pairs.last().id' -v "${NUM_CLOB_PAIRS}"
	dasel put -t string -f "$GENESIS" '.app_state.clob.clob_pairs.last().status' -v 'STATUS_ACTIVE'
	dasel put -t int -f "$GENESIS" '.app_state.clob.clob_pairs.last().perpetual_clob_metadata.perpetual_id' -v "${NUM_PERPETUALS}"
	dasel put -t int -f "$GENESIS" '.app_state.clob.clob_pairs.last().step_base_quantums' -v '1000000'
	dasel put -t int -f "$GENESIS" '.app_state.clob.clob_pairs.last().subticks_per_tick' -v '100' # $0.01 ticks
	dasel put -t int -f "$GENESIS" '.app_state.clob.clob_pairs.last().quantum_conversion_exponent' -v '-8'
}

# Modify the genesis file to only use fixed price exchange.
function update_all_markets_with_fixed_price_exchange() {
    GENESIS=$1/genesis.json

    # Read the number of markets
    NUM_MARKETS=$(jq -c '.app_state.prices.market_params | length' < "${GENESIS}")

    # Loop through each market and update the parameters
    for ((j = 0; j < NUM_MARKETS; j++)); do
        # Get the current ticker
        TICKER=$(jq -r ".app_state.prices.market_params[$j].pair" < "${GENESIS}")

        # Update the exchange_config_json using the EOF syntax
        exchange_config_json=$(cat <<-EOF
{
    "exchanges": [
        {
            "exchangeName": "TestFixedPriceExchange",
            "ticker": "${TICKER}"
        }
    ]
}
EOF
        )
        dasel put -t string -f "$GENESIS" ".app_state.prices.market_params.[$j].exchange_config_json" -v "$exchange_config_json"

        # Update the min_exchanges
        dasel put -t int -f "$GENESIS" ".app_state.prices.market_params.[$j].min_exchanges" -v "1"
    done
}

# Modify the genesis file with reduced complete bridge delay (for testing in non-prod envs).
update_genesis_complete_bridge_delay() {
	GENESIS=$1/genesis.json

	# Reduce complete bridge delay to 600 blocks.
	dasel put -t int -f "$GENESIS" '.app_state.bridge.safety_params.delay_blocks' -v "$2"
}

# Set the denom metadata, which is for human readability.
set_denom_metadata() {
	local BASE_DENOM=$1
	local WHOLE_COIN_DENOM=$2
	local COIN_NAME=$3
	dasel put -t json -f "$GENESIS" ".app_state.bank.denom_metadata" -v "[]"
	dasel put -t json -f "$GENESIS" ".app_state.bank.denom_metadata.[]" -v "{}"
	dasel put -t string -f "$GENESIS" ".app_state.bank.denom_metadata.[0].description" -v "The native token of the network"
	dasel put -t json -f "$GENESIS" ".app_state.bank.denom_metadata.[0].denom_units" -v "[]"
	dasel put -t json -f "$GENESIS" ".app_state.bank.denom_metadata.[0].denom_units.[]" -v "{}"
	# Base denom is the minimum unit of the a token and the denom used by `x/bank`.
	dasel put -t string -f "$GENESIS" ".app_state.bank.denom_metadata.[0].denom_units.[0].denom" -v "$BASE_DENOM"
	dasel put -t json -f "$GENESIS" ".app_state.bank.denom_metadata.[0].denom_units.[]" -v "{}"
	dasel put -t string -f "$GENESIS" ".app_state.bank.denom_metadata.[0].denom_units.[1].denom" -v "$WHOLE_COIN_DENOM"
	dasel put -t int -f "$GENESIS" ".app_state.bank.denom_metadata.[0].denom_units.[1].exponent" -v 18
	dasel put -t string -f "$GENESIS" ".app_state.bank.denom_metadata.[0].base" -v "$BASE_DENOM"
	dasel put -t string -f "$GENESIS" ".app_state.bank.denom_metadata.[0].name" -v "$COIN_NAME"
	dasel put -t string -f "$GENESIS" ".app_state.bank.denom_metadata.[0].display" -v "$WHOLE_COIN_DENOM"
	dasel put -t string -f "$GENESIS" ".app_state.bank.denom_metadata.[0].symbol" -v "$WHOLE_COIN_DENOM"
}
