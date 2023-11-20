package metrics

// 1. do not make it too long
// 2. prefix/suffix rules
const (
	// Clob metric keys
	// Currently Sorted by type of metric.
	// TODO make sure the metric type is clear.

	// stat
	ClobExpiredStatefulOrders                       = "clob_expired_stateful_order_removed"
	ClobPrepareCheckStateCannotDeleverageSubaccount = "clob_prepare_check_state_cannot_deleverage_subaccount"
	ClobDeleverageSubaccountTotalQuoteQuantums      = "clob_deleverage_subaccount_total_quote_quantums"
	ClobDeleverageSubaccount                        = "clob_deleverage_subaccount"

	// gauge
	InsuranceFundBalance = "insurance_fund_balance"

	// sample
	ClobDeleverageSubaccountTotalQuoteQuantumsDistribution = "clob_deleverage_subaccount_total_quote_quantums_distribution"
	DeleveragingPercentFilledDistribution                  = "deleveraging_percent_filled_distribution"
	ClobDeleveragingNumSubaccountsIteratedCount            = "clob_deleveraging_num_subaccounts_iterated_count"
	ClobDeleveragingNonOverlappingBankrupcyPricesCount     = "clob_deleveraging_non_overlapping_bankruptcy_prices_count"
	ClobDeleveragingNoOpenPositionOnOppositeSideCount      = "clob_deleveraging_no_open_position_on_opposite_side_count"
	ClobDeleverageSubaccountFilledQuoteQuantums            = "clob_deleverage_subaccount_filled_quote_quantums"

	// Measure since
	ClobOffsettingSubaccountPerpetualPosition = "clob_offsetting_subaccount_perpetual_position"
)
