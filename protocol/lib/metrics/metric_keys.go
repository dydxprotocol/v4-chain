// nolint:lll
package metrics

// Metrics Keys Guidelines
// 1. Be wary of length
// 2. Prefix by module
// 3. Suffix keys with a unit of measurement
// 4. Delimit with '_'
// 5. Information such as callback type should be added as tags, not in key names.
// Example: clob_place_order_count, clob_msg_place_order_latency_ms, clob_operations_queue_length
// clob_expired_stateful_orders_count, clob_processed_orders_ms_total

// Clob Metrics Keys
const (
	// Stats
	ClobExpiredStatefulOrders                          = "clob_expired_stateful_order_removed"
	ClobPrepareCheckStateCannotDeleverageSubaccount    = "clob_prepare_check_state_cannot_deleverage_subaccount"
	ClobDeleverageSubaccountTotalQuoteQuantums         = "clob_deleverage_subaccount_total_quote_quantums"
	ClobDeleverageSubaccount                           = "clob_deleverage_subaccount"
	LiquidationsPlacePerpetualLiquidationQuoteQuantums = "liquidations_place_perpetual_liquidation_quote_quantums"
	LiquidationsLiquidationMatchNegativeTNC            = "liquidations_liquidation_match_negative_tnc"
	ClobMevErrorCount                                  = "clob_mev_error_count"

	// Gauges
	InsuranceFundBalance = "insurance_fund_balance"
	ClobMev              = "clob_mev"

	// Samples
	ClobDeleverageSubaccountTotalQuoteQuantumsDistribution         = "clob_deleverage_subaccount_total_quote_quantums_distribution"
	DeleveragingPercentFilledDistribution                          = "deleveraging_percent_filled_distribution"
	ClobDeleveragingNumSubaccountsIteratedCount                    = "clob_deleveraging_num_subaccounts_iterated_count"
	ClobDeleveragingNonOverlappingBankrupcyPricesCount             = "clob_deleveraging_non_overlapping_bankruptcy_prices_count"
	ClobDeleveragingNoOpenPositionOnOppositeSideCount              = "clob_deleveraging_no_open_position_on_opposite_side_count"
	ClobDeleverageSubaccountFilledQuoteQuantums                    = "clob_deleverage_subaccount_filled_quote_quantums"
	LiquidationsLiquidatableSubaccountIdsCount                     = "liquidations_liquidatable_subaccount_ids_count"
	LiquidationsPercentFilledDistribution                          = "liquidations_percent_filled_distribution"
	LiquidationsPlacePerpetualLiquidationQuoteQuantumsDistribution = "liquidations_place_perpetual_liquidation_quote_quantums_distribution"

	// Measure Since
	ClobOffsettingSubaccountPerpetualPosition = "clob_offsetting_subaccount_perpetual_position"
	DaemonGetPreviousBlockInfoLatency         = "daemon_get_previous_block_info_latency"
	DaemonGetAllMarketPricesLatency           = "daemon_get_all_market_prices_latency"
	DaemonGetMarketPricesPaginatedLatency     = "daemon_get_market_prices_paginated_latency"
	MevLatency                                = "mev_latency"
)
