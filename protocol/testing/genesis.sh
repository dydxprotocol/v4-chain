#!/bin/bash
set -eo pipefail

# Below is the shared genesis configuration for local development, as well as our testnets.
# Any changes to genesis state in those environments belong here.
# If you are making a change to genesis which is _required_ for the chain to function,
# then that change probably belongs in `DefaultGenesis` for the module, and not here.

# Address of the `subaccounts` module account.
# Obtained from `authtypes.NewModuleAddress(subaccounttypes.ModuleName)`.
SUBACCOUNTS_MODACC_ADDR="dydx1v88c3xv9xyv3eetdx0tvcmq7ung3dywp5upwc6"
USDC_DENOM='ibc/xxx'
DEFAULT_SUBACCOUNT_QUOTE_BALANCE=100000000000000000
DEFAULT_SUBACCOUNT_QUOTE_BALANCE_FAUCET=900000000000000000

function edit_genesis() {
	GENESIS=$1/genesis.json

	# IFS stands for "Internal Field Separator" and it's a special var that determins how bash splits strings.
	# So IFS=' ' specifies that we want to split on spaces.
	# "read" is a built in bash command that reads from stdin.
	# The -a flag specifies that the input should be treated as an array and assign it to the var specified after.
	# The -r flag tells the command to not treat a Backslash as an escape character.
	IFS=' ' read -ra TEST_ACCOUNTS <<<"${2}"
	IFS=' ' read -ra FAUCET_ACCOUNTS <<<"${3}"

	# Update staking module.
	# TODO(DEC-1673): restore or update the unbonding time to a reasonable value.
	dasel put string -f "$GENESIS" '.app_state.staking.params.unbonding_time' '7200s'

	# Update assets module.
	dasel put int -f "$GENESIS" '.app_state.assets.assets.[0].id' '0'
	dasel put string -f "$GENESIS" '.app_state.assets.assets.[0].symbol' 'USDC'
	dasel put string -f "$GENESIS" '.app_state.assets.assets.[0].denom' "$USDC_DENOM"
	dasel put string -f "$GENESIS" '.app_state.assets.assets.[0].denom_exponent' -v '-6'
	dasel put bool -f "$GENESIS" '.app_state.assets.assets.[0].has_market' 'false'
	dasel put int -f "$GENESIS" '.app_state.assets.assets.[0].market_id' '0'
	dasel put int -f "$GENESIS" '.app_state.assets.assets.[0].atomic_resolution' -v '-6'
	dasel put int -f "$GENESIS" '.app_state.assets.assets.[0].long_interest' '0'

	# Update perpetuals module.
	# Liquidity Tiers.
	dasel put int -f "$GENESIS" '.app_state.perpetuals.liquidity_tiers.[0].id' '0'
	dasel put string -f "$GENESIS" '.app_state.perpetuals.liquidity_tiers.[0].name' 'Large-Cap'
	dasel put int -f "$GENESIS" '.app_state.perpetuals.liquidity_tiers.[0].initial_margin_ppm' '50000' # 5 %
	dasel put int -f "$GENESIS" '.app_state.perpetuals.liquidity_tiers.[0].maintenance_fraction_ppm' '600000' # 3 % (60% of IM)
	dasel put int -f "$GENESIS" '.app_state.perpetuals.liquidity_tiers.[0].base_position_notional' '1000000000000' # 1_000_000 USDC

	# Params.
	dasel put int -f "$GENESIS" '.app_state.perpetuals.params.funding_rate_clamp_factor_ppm' '6000000' # 600 % (same as 75% on hourly rate)
	dasel put int -f "$GENESIS" '.app_state.perpetuals.params.premium_vote_clamp_factor_ppm' '60000000' # 6000 % (some multiples of funding rate clamp factor)

	# Perpetuals.
	dasel put string -f "$GENESIS" '.app_state.perpetuals.perpetuals.[0].ticker' 'BTC-USD'
	dasel put int -f "$GENESIS" '.app_state.perpetuals.perpetuals.[0].id' '0'
	dasel put int -f "$GENESIS" '.app_state.perpetuals.perpetuals.[0].market_id' '0'
	dasel put int -f "$GENESIS" '.app_state.perpetuals.perpetuals.[0].atomic_resolution' -v '-10'
	dasel put int -f "$GENESIS" '.app_state.perpetuals.perpetuals.[0].default_funding_ppm' '0'
	dasel put int -f "$GENESIS" '.app_state.perpetuals.perpetuals.[0].liquidity_tier' '0'

	dasel put string -f "$GENESIS" '.app_state.perpetuals.perpetuals.[1].ticker' 'ETH-USD'
	dasel put int -f "$GENESIS" '.app_state.perpetuals.perpetuals.[1].id' '1'
	dasel put int -f "$GENESIS" '.app_state.perpetuals.perpetuals.[1].market_id' '1'
	dasel put int -f "$GENESIS" '.app_state.perpetuals.perpetuals.[1].atomic_resolution' -v '-9'
	dasel put int -f "$GENESIS" '.app_state.perpetuals.perpetuals.[1].default_funding_ppm' '0'
	dasel put int -f "$GENESIS" '.app_state.perpetuals.perpetuals.[0].liquidity_tier' '0'

	# Update prices module.
	# ExchangeFeed: Binance
	dasel put int -f "$GENESIS" '.app_state.prices.exchange_feeds.[0].id' '0'
	dasel put string -f "$GENESIS" '.app_state.prices.exchange_feeds.[0].name' 'Binance'
	dasel put string -f "$GENESIS" '.app_state.prices.exchange_feeds.[0].memo' 'Memo for Binance'
	# ExchangeFeed: BinanceUS
	dasel put int -f "$GENESIS" '.app_state.prices.exchange_feeds.[1].id' '1'
	dasel put string -f "$GENESIS" '.app_state.prices.exchange_feeds.[1].name' 'BinanceUS'
	dasel put string -f "$GENESIS" '.app_state.prices.exchange_feeds.[1].memo' 'Memo for Binance US'
	# ExchangeFeed: Bitfinex
	dasel put int -f "$GENESIS" '.app_state.prices.exchange_feeds.[2].id' '2'
	dasel put string -f "$GENESIS" '.app_state.prices.exchange_feeds.[2].name' 'Bitfinex'
	dasel put string -f "$GENESIS" '.app_state.prices.exchange_feeds.[2].memo' 'Memo for Bitfinex'
# ExchangeFeed: Kraken
	dasel put int -f "$GENESIS" '.app_state.prices.exchange_feeds.[3].id' '3'
	dasel put string -f "$GENESIS" '.app_state.prices.exchange_feeds.[3].name' 'Kraken'
	dasel put string -f "$GENESIS" '.app_state.prices.exchange_feeds.[3].memo' 'Memo for Kraken'

	# Market: BTC-USD
	dasel put string -f "$GENESIS" '.app_state.prices.markets.[0].pair' 'BTC-USD'
	dasel put int -f "$GENESIS" '.app_state.prices.markets.[0].id' '0'
	dasel put int -f "$GENESIS" '.app_state.prices.markets.[0].exponent' -v '-5'
	dasel put int -f "$GENESIS" '.app_state.prices.markets.[0].min_exchanges' '2'
	dasel put int -f "$GENESIS" '.app_state.prices.markets.[0].min_price_change_ppm' '1000' # 0.1%
	dasel put int -f "$GENESIS" '.app_state.prices.markets.[0].price' '2000000000'          # $20,000 = 1 BTC.
	dasel put int -f "$GENESIS" '.app_state.prices.markets.[0].exchanges.[0]' '0'
	dasel put int -f "$GENESIS" '.app_state.prices.markets.[0].exchanges.[1]' '1'
	dasel put int -f "$GENESIS" '.app_state.prices.markets.[0].exchanges.[2]' '2'
	dasel put int -f "$GENESIS" '.app_state.prices.markets.[0].exchanges.[3]' '3'
	# Market: ETH-USD
	dasel put string -f "$GENESIS" '.app_state.prices.markets.[1].pair' 'ETH-USD'
	dasel put int -f "$GENESIS" '.app_state.prices.markets.[1].id' '1'
	dasel put int -f "$GENESIS" '.app_state.prices.markets.[1].exponent' -v '-6'
	dasel put int -f "$GENESIS" '.app_state.prices.markets.[1].min_exchanges' '2'
	dasel put int -f "$GENESIS" '.app_state.prices.markets.[1].min_price_change_ppm' '1000' # 0.1%
	dasel put int -f "$GENESIS" '.app_state.prices.markets.[1].price' '1500000000'          # $1,500 = 1 ETH.
	dasel put int -f "$GENESIS" '.app_state.prices.markets.[1].exchanges.[0]' '0'
	dasel put int -f "$GENESIS" '.app_state.prices.markets.[1].exchanges.[1]' '1'
	dasel put int -f "$GENESIS" '.app_state.prices.markets.[1].exchanges.[2]' '2'
	dasel put int -f "$GENESIS" '.app_state.prices.markets.[1].exchanges.[2]' '3'
	# Market: LINK-USD
	dasel put string -f "$GENESIS" '.app_state.prices.markets.[2].pair' 'LINK-USD'
	dasel put int -f "$GENESIS" '.app_state.prices.markets.[2].id' '2'
	dasel put int -f "$GENESIS" '.app_state.prices.markets.[2].exponent' -v '-8'
	dasel put int -f "$GENESIS" '.app_state.prices.markets.[2].min_exchanges' '1'
	dasel put int -f "$GENESIS" '.app_state.prices.markets.[2].min_price_change_ppm' '1000' # 0.1%
	dasel put int -f "$GENESIS" '.app_state.prices.markets.[2].price' '1000000000'          # $10 = 1 LINK.
	dasel put int -f "$GENESIS" '.app_state.prices.markets.[2].exchanges.[0]' '0'
	dasel put int -f "$GENESIS" '.app_state.prices.markets.[2].exchanges.[1]' '1'
	dasel put int -f "$GENESIS" '.app_state.prices.markets.[2].exchanges.[2]' '3'

	total_accounts_quote_balance=0
	acct_idx=0
	# Update subaccounts module for load testing accounts.
	for acct in "${TEST_ACCOUNTS[@]}"; do
		add_subaccount "$GENESIS" "$acct_idx" "$acct" "$DEFAULT_SUBACCOUNT_QUOTE_BALANCE"
		total_accounts_quote_balance=$(($total_accounts_quote_balance + $DEFAULT_SUBACCOUNT_QUOTE_BALANCE))
		acct_idx=$(($acct_idx + 1))
	done
	# Update subaccounts module for faucet accounts.
	for acct in "${FAUCET_ACCOUNTS[@]}"; do
		add_subaccount "$GENESIS" "$acct_idx" "$acct" "$DEFAULT_SUBACCOUNT_QUOTE_BALANCE_FAUCET"
		total_accounts_quote_balance=$(($total_accounts_quote_balance + $DEFAULT_SUBACCOUNT_QUOTE_BALANCE_FAUCET))
		acct_idx=$(($acct_idx + 1))
	done

	# Initialize subaccounts module account balance within bank module.
	dasel put string -f "$GENESIS" ".app_state.bank.balances.[0].address" "${SUBACCOUNTS_MODACC_ADDR}"
	dasel put string -f "$GENESIS" ".app_state.bank.balances.[0].coins.[0].denom" "$USDC_DENOM"
	# TODO(DEC-969): For testnet, ensure subaccounts module balance >= sum of subaccount quote balances.
	dasel put string -f "$GENESIS" ".app_state.bank.balances.[0].coins.[0].amount" "${total_accounts_quote_balance}"

	# Update clob module.
	dasel put int -f "$GENESIS" '.app_state.clob.clob_pairs.[0].id' '0'
	dasel put string -f "$GENESIS" '.app_state.clob.clob_pairs.[0].status' 'STATUS_ACTIVE'
	dasel put int -f "$GENESIS" '.app_state.clob.clob_pairs.[0].perpetual_clob_metadata.perpetual_id' '0'
	dasel put int -f "$GENESIS" '.app_state.clob.clob_pairs.[0].step_base_quantums' '1000000'
	dasel put int -f "$GENESIS" '.app_state.clob.clob_pairs.[0].subticks_per_tick' '10000'
	dasel put int -f "$GENESIS" '.app_state.clob.clob_pairs.[0].min_order_base_quantums' '1000000'
	dasel put int -f "$GENESIS" '.app_state.clob.clob_pairs.[0].quantum_conversion_exponent' -v '-8'
	dasel put int -f "$GENESIS" '.app_state.clob.clob_pairs.[0].maker_fee_ppm' -v '200'
	dasel put int -f "$GENESIS" '.app_state.clob.clob_pairs.[0].taker_fee_ppm' -v '500'

	dasel put int -f "$GENESIS" '.app_state.clob.clob_pairs.[1].id' '1'
	dasel put string -f "$GENESIS" '.app_state.clob.clob_pairs.[1].status' 'STATUS_ACTIVE'
	dasel put int -f "$GENESIS" '.app_state.clob.clob_pairs.[1].perpetual_clob_metadata.perpetual_id' '1'
	dasel put int -f "$GENESIS" '.app_state.clob.clob_pairs.[1].step_base_quantums' '1000000'
	dasel put int -f "$GENESIS" '.app_state.clob.clob_pairs.[1].subticks_per_tick' '100000'
	dasel put int -f "$GENESIS" '.app_state.clob.clob_pairs.[1].min_order_base_quantums' '1000000'
	dasel put int -f "$GENESIS" '.app_state.clob.clob_pairs.[1].quantum_conversion_exponent' -v '-9'
	dasel put int -f "$GENESIS" '.app_state.clob.clob_pairs.[1].maker_fee_ppm' -v '200'
	dasel put int -f "$GENESIS" '.app_state.clob.clob_pairs.[1].taker_fee_ppm' -v '500'

	dasel put int -f "$GENESIS" '.app_state.clob.liquidations_config.max_liquidation_fee_ppm' '5000'
	dasel put int -f "$GENESIS" '.app_state.clob.liquidations_config.fillable_price_config.bankruptcy_adjustment_ppm' '1000000'
	dasel put int -f "$GENESIS" '.app_state.clob.liquidations_config.fillable_price_config.spread_to_maintenance_margin_ratio_ppm' '100000'
	dasel put int -f "$GENESIS" '.app_state.clob.liquidations_config.position_block_limits.min_position_notional_liquidated' '1000'
	dasel put int -f "$GENESIS" '.app_state.clob.liquidations_config.position_block_limits.max_position_portion_liquidated_ppm' '1000000'
	dasel put int -f "$GENESIS" '.app_state.clob.liquidations_config.subaccount_block_limits.max_notional_liquidated' '100000000000000'
	dasel put int -f "$GENESIS" '.app_state.clob.liquidations_config.subaccount_block_limits.max_quantums_insurance_lost' '100000000000000'
}

function add_subaccount() {
	GEN_FILE=$1
	IDX=$2
	ACCOUNT=$3
	QUOTE_BALANCE=$4

	dasel put string -f "$GEN_FILE" ".app_state.subaccounts.subaccounts.[$IDX].id.owner" "$ACCOUNT"
	dasel put int -f "$GEN_FILE" ".app_state.subaccounts.subaccounts.[$IDX].id.number" '0'
	dasel put bool -f "$GEN_FILE" ".app_state.subaccounts.subaccounts.[$IDX].margin_enabled" 'true'

	dasel put int -f "$GEN_FILE" ".app_state.subaccounts.subaccounts.[$IDX].asset_positions.[0].asset_id" '0'
	dasel put string -f "$GEN_FILE" ".app_state.subaccounts.subaccounts.[$IDX].asset_positions.[0].quantums" "${QUOTE_BALANCE}"
	dasel put int -f "$GEN_FILE" ".app_state.subaccounts.subaccounts.[$IDX].asset_positions.[0].index" '0'
}
