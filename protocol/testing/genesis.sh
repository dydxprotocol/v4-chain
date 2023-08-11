#!/bin/bash
set -eo pipefail

# Below is the shared genesis configuration for local development, as well as our testnets.
# Any changes to genesis state in those environments belong here.
# If you are making a change to genesis which is _required_ for the chain to function,
# then that change probably belongs in `DefaultGenesis` for the module, and not here.

# Address of the `subaccounts` module account.
# Obtained from `authtypes.NewModuleAddress(subaccounttypes.ModuleName)`.
SUBACCOUNTS_MODACC_ADDR="dydx1v88c3xv9xyv3eetdx0tvcmq7ung3dywp5upwc6"
USDC_DENOM="ibc/usdc-placeholder"
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
	dasel put -t string -f "$GENESIS" '.app_state.staking.params.unbonding_time' -v '7200s'

	# Update assets module.
	dasel put -t int -f "$GENESIS" '.app_state.assets.assets.[0].id' -v '0'
	dasel put -t string -f "$GENESIS" '.app_state.assets.assets.[0].symbol' -v 'USDC'
	dasel put -t string -f "$GENESIS" '.app_state.assets.assets.[0].denom' -v "$USDC_DENOM"
	dasel put -t string -f "$GENESIS" '.app_state.assets.assets.[0].denom_exponent' -v '-6'
	dasel put -t bool -f "$GENESIS" '.app_state.assets.assets.[0].has_market' -v 'false'
	dasel put -t int -f "$GENESIS" '.app_state.assets.assets.[0].market_id' -v '0'
	dasel put -t int -f "$GENESIS" '.app_state.assets.assets.[0].atomic_resolution' -v '-6'
	dasel put -t int -f "$GENESIS" '.app_state.assets.assets.[0].long_interest' -v '0'

	# Update perpetuals module.
	# Liquidity Tiers.
	dasel put -t json -f "$GENESIS" '.app_state.perpetuals.liquidity_tiers' -v "[]"
	dasel put -t json -f "$GENESIS" '.app_state.perpetuals.liquidity_tiers.[]' -v "{}"
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.liquidity_tiers.[0].id' -v '0'
	dasel put -t string -f "$GENESIS" '.app_state.perpetuals.liquidity_tiers.[0].name' -v 'Large-Cap'
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.liquidity_tiers.[0].initial_margin_ppm' -v '50000' # 5 %
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.liquidity_tiers.[0].maintenance_fraction_ppm' -v '600000' # 3 % (60% of IM)
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.liquidity_tiers.[0].base_position_notional' -v '1000000000000' # 1_000_000 USDC

	# Params.
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.params.funding_rate_clamp_factor_ppm' -v '6000000' # 600 % (same as 75% on hourly rate)
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.params.premium_vote_clamp_factor_ppm' -v '60000000' # 6000 % (some multiples of funding rate clamp factor)

	# Perpetuals.
	dasel put -t json -f "$GENESIS" '.app_state.perpetuals.perpetuals' -v "[]"
	dasel put -t json -f "$GENESIS" '.app_state.perpetuals.perpetuals.[]' -v "{}"
	dasel put -t string -f "$GENESIS" '.app_state.perpetuals.perpetuals.[0].ticker' -v 'BTC-USD'
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.perpetuals.[0].id' -v '0'
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.perpetuals.[0].market_id' -v '0'
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.perpetuals.[0].atomic_resolution' -v '-10'
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.perpetuals.[0].default_funding_ppm' -v '0'
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.perpetuals.[0].liquidity_tier' -v '0'

	dasel put -t json -f "$GENESIS" '.app_state.perpetuals.perpetuals.[]' -v "{}"
	dasel put -t string -f "$GENESIS" '.app_state.perpetuals.perpetuals.[1].ticker' -v 'ETH-USD'
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.perpetuals.[1].id' -v '1'
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.perpetuals.[1].market_id' -v '1'
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.perpetuals.[1].atomic_resolution' -v '-9'
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.perpetuals.[1].default_funding_ppm' -v '0'
	dasel put -t int -f "$GENESIS" '.app_state.perpetuals.perpetuals.[0].liquidity_tier' -v '0'

	# Update prices module.
	# ExchangeFeed: Binance
	dasel put -t json -f "$GENESIS" '.app_state.prices.exchange_feeds' -v "[]"
	dasel put -t json -f "$GENESIS" '.app_state.prices.exchange_feeds.[]' -v "{}"
	dasel put -t int -f "$GENESIS" '.app_state.prices.exchange_feeds.[0].id' -v '0'
	dasel put -t string -f "$GENESIS" '.app_state.prices.exchange_feeds.[0].name' -v 'Binance'
	dasel put -t string -f "$GENESIS" '.app_state.prices.exchange_feeds.[0].memo' -v 'Memo for Binance'
	# ExchangeFeed: BinanceUS
	dasel put -t json -f "$GENESIS" '.app_state.prices.exchange_feeds.[]' -v "{}"
	dasel put -t int -f "$GENESIS" '.app_state.prices.exchange_feeds.[1].id' -v '1'
	dasel put -t string -f "$GENESIS" '.app_state.prices.exchange_feeds.[1].name' -v 'BinanceUS'
	dasel put -t string -f "$GENESIS" '.app_state.prices.exchange_feeds.[1].memo' -v 'Memo for Binance US'
	# ExchangeFeed: Bitfinex
	dasel put -t json -f "$GENESIS" '.app_state.prices.exchange_feeds.[]' -v "{}"
	dasel put -t int -f "$GENESIS" '.app_state.prices.exchange_feeds.[2].id' -v '2'
	dasel put -t string -f "$GENESIS" '.app_state.prices.exchange_feeds.[2].name' -v 'Bitfinex'
	dasel put -t string -f "$GENESIS" '.app_state.prices.exchange_feeds.[2].memo' -v 'Memo for Bitfinex'
# ExchangeFeed: Kraken
	dasel put -t json -f "$GENESIS" '.app_state.prices.exchange_feeds.[]' -v "{}"
	dasel put -t int -f "$GENESIS" '.app_state.prices.exchange_feeds.[3].id' -v '3'
	dasel put -t string -f "$GENESIS" '.app_state.prices.exchange_feeds.[3].name' -v 'Kraken'
	dasel put -t string -f "$GENESIS" '.app_state.prices.exchange_feeds.[3].memo' -v 'Memo for Kraken'

	# Market: BTC-USD
	dasel put -t json -f "$GENESIS" '.app_state.prices.markets' -v "[]"
	dasel put -t json -f "$GENESIS" '.app_state.prices.markets.[]' -v "{}"
	dasel put -t string -f "$GENESIS" '.app_state.prices.markets.[0].pair' -v 'BTC-USD'
	dasel put -t int -f "$GENESIS" '.app_state.prices.markets.[0].id' -v '0'
	dasel put -t int -f "$GENESIS" '.app_state.prices.markets.[0].exponent' -v '-5'
	dasel put -t int -f "$GENESIS" '.app_state.prices.markets.[0].min_exchanges' -v '2'
	dasel put -t int -f "$GENESIS" '.app_state.prices.markets.[0].min_price_change_ppm' -v '1000' # 0.1%
	dasel put -t int -f "$GENESIS" '.app_state.prices.markets.[0].price' -v '2000000000'          # $20,000 = 1 BTC.
	dasel put -t json -f "$GENESIS" '.app_state.prices.markets.[0].exchanges' -v '[]'
	dasel put -t int -f "$GENESIS" '.app_state.prices.markets.[0].exchanges.[]' -v '0'
	dasel put -t int -f "$GENESIS" '.app_state.prices.markets.[0].exchanges.[]' -v '1'
	dasel put -t int -f "$GENESIS" '.app_state.prices.markets.[0].exchanges.[]' -v '2'
	dasel put -t int -f "$GENESIS" '.app_state.prices.markets.[0].exchanges.[]' -v '3'
	# Market: ETH-USD
	dasel put -t json -f "$GENESIS" '.app_state.prices.markets.[]' -v "{}"
	dasel put -t string -f "$GENESIS" '.app_state.prices.markets.[1].pair' -v 'ETH-USD'
	dasel put -t int -f "$GENESIS" '.app_state.prices.markets.[1].id' -v '1'
	dasel put -t int -f "$GENESIS" '.app_state.prices.markets.[1].exponent' -v '-6'
	dasel put -t int -f "$GENESIS" '.app_state.prices.markets.[1].min_exchanges' -v '2'
	dasel put -t int -f "$GENESIS" '.app_state.prices.markets.[1].min_price_change_ppm' -v '1000' # 0.1%
	dasel put -t int -f "$GENESIS" '.app_state.prices.markets.[1].price' -v '1500000000'          # $1,500 = 1 ETH.
	dasel put -t json -f "$GENESIS" '.app_state.prices.markets.[1].exchanges' -v '[]'
	dasel put -t int -f "$GENESIS" '.app_state.prices.markets.[1].exchanges.[]' -v '0'
	dasel put -t int -f "$GENESIS" '.app_state.prices.markets.[1].exchanges.[]' -v '1'
	dasel put -t int -f "$GENESIS" '.app_state.prices.markets.[1].exchanges.[]' -v '2'
	dasel put -t int -f "$GENESIS" '.app_state.prices.markets.[1].exchanges.[]' -v '3'

	# Market: LINK-USD
	dasel put -t json -f "$GENESIS" '.app_state.prices.markets.[]' -v "{}"
	dasel put -t string -f "$GENESIS" '.app_state.prices.markets.[2].pair' -v 'LINK-USD'
	dasel put -t int -f "$GENESIS" '.app_state.prices.markets.[2].id' -v '2'
	dasel put -t int -f "$GENESIS" '.app_state.prices.markets.[2].exponent' -v '-8'
	dasel put -t int -f "$GENESIS" '.app_state.prices.markets.[2].min_exchanges' -v '1'
	dasel put -t int -f "$GENESIS" '.app_state.prices.markets.[2].min_price_change_ppm' -v '1000' # 0.1%
	dasel put -t int -f "$GENESIS" '.app_state.prices.markets.[2].price' -v '1000000000'          # $10 = 1 LINK.
	dasel put -t json -f "$GENESIS" '.app_state.prices.markets.[2].exchanges' -v '[]'
	dasel put -t int -f "$GENESIS" '.app_state.prices.markets.[2].exchanges.[]' -v '0'
	dasel put -t int -f "$GENESIS" '.app_state.prices.markets.[2].exchanges.[]' -v '1'
	dasel put -t int -f "$GENESIS" '.app_state.prices.markets.[2].exchanges.[]' -v '3'

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
	dasel put -t json -f "$GENESIS" ".app_state.bank.balances.[]" -v "{}"
	dasel put -t string -f "$GENESIS" ".app_state.bank.balances.[0].address" -v "${SUBACCOUNTS_MODACC_ADDR}"
	dasel put -t json -f "$GENESIS" ".app_state.bank.balances.[0].coins.[]" -v "{}"
	dasel put -t string -f "$GENESIS" ".app_state.bank.balances.[0].coins.[0].denom" -v "$USDC_DENOM"
	# TODO(DEC-969): For testnet, ensure subaccounts module balance >= sum of subaccount quote balances.
	dasel put -t string -f "$GENESIS" ".app_state.bank.balances.[0].coins.[0].amount" -v "${total_accounts_quote_balance}"

	# Update clob module.
	dasel put -t json -f "$GENESIS" '.app_state.clob.clob_pairs.[]' -v '{}'
	dasel put -t int -f "$GENESIS" '.app_state.clob.clob_pairs.[0].id' -v '0'
	dasel put -t string -f "$GENESIS" '.app_state.clob.clob_pairs.[0].status' -v 'STATUS_ACTIVE'
	dasel put -t int -f "$GENESIS" '.app_state.clob.clob_pairs.[0].perpetual_clob_metadata.perpetual_id' -v '0'
	dasel put -t int -f "$GENESIS" '.app_state.clob.clob_pairs.[0].step_base_quantums' -v '1000000'
	dasel put -t int -f "$GENESIS" '.app_state.clob.clob_pairs.[0].subticks_per_tick' -v '10000'
	dasel put -t int -f "$GENESIS" '.app_state.clob.clob_pairs.[0].min_order_base_quantums' -v '1000000'
	dasel put -t int -f "$GENESIS" '.app_state.clob.clob_pairs.[0].quantum_conversion_exponent' -v '-8'
	dasel put -t int -f "$GENESIS" '.app_state.clob.clob_pairs.[0].maker_fee_ppm' -v '200'
	dasel put -t int -f "$GENESIS" '.app_state.clob.clob_pairs.[0].taker_fee_ppm' -v '500'

	dasel put -t json -f "$GENESIS" '.app_state.clob.clob_pairs.[]' -v '{}'
	dasel put -t int -f "$GENESIS" '.app_state.clob.clob_pairs.[1].id' -v '1'
	dasel put -t string -f "$GENESIS" '.app_state.clob.clob_pairs.[1].status' -v 'STATUS_ACTIVE'
	dasel put -t int -f "$GENESIS" '.app_state.clob.clob_pairs.[1].perpetual_clob_metadata.perpetual_id' -v '1'
	dasel put -t int -f "$GENESIS" '.app_state.clob.clob_pairs.[1].step_base_quantums' -v '1000000'
	dasel put -t int -f "$GENESIS" '.app_state.clob.clob_pairs.[1].subticks_per_tick' -v '100000'
	dasel put -t int -f "$GENESIS" '.app_state.clob.clob_pairs.[1].min_order_base_quantums' -v '1000000'
	dasel put -t int -f "$GENESIS" '.app_state.clob.clob_pairs.[1].quantum_conversion_exponent' -v '-9'
	dasel put -t int -f "$GENESIS" '.app_state.clob.clob_pairs.[1].maker_fee_ppm' -v '200'
	dasel put -t int -f "$GENESIS" '.app_state.clob.clob_pairs.[1].taker_fee_ppm' -v '500'

	dasel put -t int -f "$GENESIS" '.app_state.clob.liquidations_config.max_liquidation_fee_ppm' -v '5000'
	dasel put -t int -f "$GENESIS" '.app_state.clob.liquidations_config.fillable_price_config.bankruptcy_adjustment_ppm' -v '1000000'
	dasel put -t int -f "$GENESIS" '.app_state.clob.liquidations_config.fillable_price_config.spread_to_maintenance_margin_ratio_ppm' -v '100000'
	dasel put -t int -f "$GENESIS" '.app_state.clob.liquidations_config.position_block_limits.min_position_notional_liquidated' -v '1000'
	dasel put -t int -f "$GENESIS" '.app_state.clob.liquidations_config.position_block_limits.max_position_portion_liquidated_ppm' -v '1000000'
	dasel put -t int -f "$GENESIS" '.app_state.clob.liquidations_config.subaccount_block_limits.max_notional_liquidated' -v '100000000000000'
	dasel put -t int -f "$GENESIS" '.app_state.clob.liquidations_config.subaccount_block_limits.max_quantums_insurance_lost' -v '100000000000000'
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
