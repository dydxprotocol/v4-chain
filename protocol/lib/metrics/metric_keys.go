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
	SubaccountsNegativeTncSubaccountSeen               = "negative_tnc_subaccount_seen"
	GateWithdrawalsIfNegativeTncSubaccountSeen         = "gate_withdrawals_if_negative_tnc_subaccount_seen"
	ChainOutageSeen                                    = "chain_outage_seen"
	SubaccountCreatedCount                             = "subaccount_created_count"
	ClobRateLimitPlaceOrderCount                       = "clob_rate_limit_place_order_count"
	ClobRateLimitCancelOrderCount                      = "clob_rate_limit_cancel_order_count"
	ClobRateLimitBatchCancelCount                      = "clob_rate_limit_batch_cancel_count"
	StatsGetStakedBaseTokensCacheHit                   = "stats_get_staked_base_tokens_cache_hit"
	StatsGetStakedBaseTokensCacheMiss                  = "stats_get_staked_base_tokens_cache_miss"

	// Gauges
	InsuranceFundBalance                      = "insurance_fund_balance"
	ClobMev                                   = "clob_mev"
	ClobConditionalOrderTriggerPrice          = "clob_conditional_order_trigger_price"
	ClobConditionalOrderTriggered             = "clob_conditional_order_triggered"
	ClobSubaccountsRequiringDeleveragingCount = "clob_subaccounts_requiring_deleveraging_count"
	SendingProcessDepositToSubaccount         = "sending_process_deposit_to_subaccount"
	RateLimitInsufficientWithdrawalAmount     = "rate_limit_insufficient_withdrawal_amount"
	StatsGetStakedBaseTokensLatencyCacheHit   = "stats_get_staked_base_tokens_latency_cache_hit"
	StatsGetStakedBaseTokensLatencyCacheMiss  = "stats_get_staked_base_tokens_latency_cache_miss"

	// Samples
	ClobDeleverageSubaccountTotalQuoteQuantumsDistribution         = "clob_deleverage_subaccount_total_quote_quantums_distribution"
	DeleveragingPercentFilledDistribution                          = "deleveraging_percent_filled_distribution"
	ClobDeleveragingNumSubaccountsIteratedCount                    = "clob_deleveraging_num_subaccounts_iterated_count"
	ClobDeleveragingNonOverlappingBankrupcyPricesCount             = "clob_deleveraging_non_overlapping_bankruptcy_prices_count"
	ClobDeleveragingNoOpenPositionOnOppositeSideCount              = "clob_deleveraging_no_open_position_on_opposite_side_count"
	ClobDeleverageSubaccountFilledQuoteQuantums                    = "clob_deleverage_subaccount_filled_quote_quantums"
	ClobSubaccountsWithFinalSettlementPositionsCount               = "clob_subaccounts_with_final_settlement_positions_count"
	LiquidationsLiquidatableSubaccountIdsCount                     = "liquidations_liquidatable_subaccount_ids_count"
	LiquidationsPercentFilledDistribution                          = "liquidations_percent_filled_distribution"
	LiquidationsPlacePerpetualLiquidationQuoteQuantumsDistribution = "liquidations_place_perpetual_liquidation_quote_quantums_distribution"
	RateLimitWithdrawalAmount                                      = "rate_limit_withdrawal_amount"
	BlockTimeDistribution                                          = "block_time_dist"

	// Measure Since
	ClobOffsettingSubaccountPerpetualPosition         = "clob_offsetting_subaccount_perpetual_position"
	ClobMaybeTriggerConditionalOrders                 = "clob_maybe_trigger_conditional_orders"
	ClobNumUntriggeredOrders                          = "clob_num_untriggered_orders"
	DaemonGetPreviousBlockInfoLatency                 = "daemon_get_previous_block_info_latency"
	DaemonGetAllMarketPricesLatency                   = "daemon_get_all_market_prices_latency"
	DaemonGetMarketPricesPaginatedLatency             = "daemon_get_market_prices_paginated_latency"
	DaemonGetAllLiquidityTiersLatency                 = "daemon_get_all_liquidity_tiers_latency"
	DaemonGetLiquidityTiersPaginatedLatency           = "daemon_get_liquidity_tiers_paginated_latency"
	DaemonGetAllPerpetualsLatency                     = "daemon_get_all_perpetuals_latency"
	DaemonGetPerpetualsPaginatedLatency               = "daemon_get_perpetuals_paginated_latency"
	MevLatency                                        = "mev_latency"
	GateWithdrawalsIfNegativeTncSubaccountSeenLatency = "gate_withdrawals_if_negative_tnc_subaccount_seen_latency"

	// Full node grpc
	FullNodeGrpc                                  = "full_node_grpc"
	GrpcSendOrderbookUpdatesLatency               = "grpc_send_orderbook_updates_latency"
	GrpcSendOrderbookSnapshotLatency              = "grpc_send_orderbook_snapshot_latency"
	GrpcSendSubaccountUpdateCount                 = "grpc_send_subaccount_update_count"
	GrpcSendPriceUpdateCount                      = "grpc_send_price_update_count"
	GrpcSendOrderbookFillsLatency                 = "grpc_send_orderbook_fills_latency"
	GrpcAddUpdateToBufferCount                    = "grpc_add_update_to_buffer_count"
	GrpcAddToSubscriptionChannelCount             = "grpc_add_to_subscription_channel_count"
	GrpcSendResponseToSubscriberCount             = "grpc_send_response_to_subscriber_count"
	GrpcStreamSubscriberCount                     = "grpc_stream_subscriber_count"
	GrpcStreamNumUpdatesBuffered                  = "grpc_stream_num_updates_buffered"
	GrpcFlushUpdatesLatency                       = "grpc_flush_updates_latency"
	GrpcSubscriptionChannelLength                 = "grpc_subscription_channel_length"
	GrpcStagedAllFinalizeBlockUpdatesCount        = "grpc_staged_all_finalize_block_updates_count"
	GrpcStagedFillFinalizeBlockUpdatesCount       = "grpc_staged_finalize_block_fill_updates_count"
	GrpcStagedSubaccountFinalizeBlockUpdatesCount = "grpc_staged_finalize_block_subaccount_updates_count"
	SubscriptionId                                = "subscription_id"

	EndBlocker    = "end_blocker"
	EndBlockerLag = "end_blocker_lag"

	// Account plus
	AuthenticatorDecoratorAnteHandleLatency = "authenticator_decorator_ante_handle_latency"
	MissingRegisteredAuthenticator          = "missing_registered_authenticator"
	AuthenticatorTrackFailed                = "authenticator_track_failed"
)
