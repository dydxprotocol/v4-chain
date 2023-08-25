package types_test

import (
	"testing"

	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	"github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
	"github.com/stretchr/testify/require"
)

func TestStatMsgProposedOperations(t *testing.T) {
	tests := map[string]struct {
		operations []types.InternalOperation
		expected   types.OperationsStats
	}{
		"Successfully increments each stat when seen in the list of operations": {
			operations: []types.InternalOperation{
				{
					Operation: &types.InternalOperation_ShortTermOrderPlacement{
						ShortTermOrderPlacement: types.NewMsgPlaceOrder(constants.Order_Alice_Num0_Id0_Clob0_Buy10_Price10_GTB16),
					},
				},
				{
					Operation: &types.InternalOperation_Match{
						Match: &types.ClobMatch{
							Match: &types.ClobMatch_MatchOrders{
								MatchOrders: &types.MatchOrders{
									TakerOrderId: constants.Order_Alice_Num0_Id0_Clob0_Buy10_Price10_GTB16.GetOrderId(),
									Fills: []types.MakerFill{
										{
											FillAmount:   5,
											MakerOrderId: constants.ConditionalOrder_Alice_Num1_Id0_Clob0_Sell5_Price10_GTBT15_StopLoss15.OrderId,
										},
									},
								},
							},
						},
					},
				},
				{
					Operation: &types.InternalOperation_Match{
						Match: &types.ClobMatch{
							Match: &types.ClobMatch_MatchOrders{
								MatchOrders: &types.MatchOrders{
									TakerOrderId: constants.LongTermOrderId_Alice_Num0_ClientId0_Clob0,
									Fills: []types.MakerFill{
										{
											FillAmount:   5,
											MakerOrderId: constants.ConditionalOrder_Bob_Num0_Id0_Clob0_Sell1BTC_Price50000_GTBT10_SL_49999.OrderId,
										},
									},
								},
							},
						},
					},
				},
				{
					Operation: &types.InternalOperation_Match{
						Match: &types.ClobMatch{
							Match: &types.ClobMatch_MatchPerpetualLiquidation{
								MatchPerpetualLiquidation: &types.MatchPerpetualLiquidation{
									Liquidated:  constants.LiquidationOrder_Bob_Num0_Clob0_Buy25_Price30_BTC.GetSubaccountId(),
									ClobPairId:  constants.LiquidationOrder_Bob_Num0_Clob0_Buy25_Price30_BTC.GetClobPairId().ToUint32(),
									PerpetualId: constants.LiquidationOrder_Bob_Num0_Clob0_Buy25_Price30_BTC.MustGetLiquidatedPerpetualId(),
									TotalSize:   constants.LiquidationOrder_Bob_Num0_Clob0_Buy25_Price30_BTC.GetBaseQuantums().ToUint64(),
									IsBuy:       constants.LiquidationOrder_Bob_Num0_Clob0_Buy25_Price30_BTC.IsBuy(),
									Fills: []types.MakerFill{
										{
											FillAmount:   5,
											MakerOrderId: constants.ConditionalOrder_Alice_Num1_Id0_Clob0_Sell5_Price10_GTBT15_StopLoss15.OrderId,
										},
									},
								},
							},
						},
					},
				},
				{
					Operation: &types.InternalOperation_Match{
						Match: &types.ClobMatch{
							Match: &types.ClobMatch_MatchPerpetualDeleveraging{
								MatchPerpetualDeleveraging: &types.MatchPerpetualDeleveraging{
									Liquidated:  constants.LiquidationOrder_Bob_Num0_Clob0_Buy100_Price20_BTC.GetSubaccountId(),
									PerpetualId: constants.ClobPair_Btc.GetPerpetualClobMetadata().PerpetualId,
									Fills: []types.MatchPerpetualDeleveraging_Fill{
										{
											FillAmount:             5,
											OffsettingSubaccountId: constants.LiquidationOrder_Dave_Num0_Clob0_Sell100BTC_Price98.GetSubaccountId(),
										},
									},
								},
							},
						},
					},
				},
				{
					Operation: &types.InternalOperation_OrderRemoval{
						OrderRemoval: &types.OrderRemoval{
							OrderId:       constants.LongTermOrderId_Alice_Num0_ClientId0_Clob0,
							RemovalReason: types.OrderRemoval_REMOVAL_REASON_INVALID_SELF_TRADE,
						},
					},
				},
				{
					Operation: &types.InternalOperation_OrderRemoval{
						OrderRemoval: &types.OrderRemoval{
							OrderId:       constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT5_SL_15.OrderId,
							RemovalReason: types.OrderRemoval_REMOVAL_REASON_UNDERCOLLATERALIZED,
						},
					},
				},
				{
					Operation: &types.InternalOperation_OrderRemoval{
						OrderRemoval: &types.OrderRemoval{
							OrderId:       constants.LongTermOrder_Bob_Num0_Id0_Clob0_Sell5_Price5_GTBT10.OrderId,
							RemovalReason: types.OrderRemoval_REMOVAL_REASON_UNDERCOLLATERALIZED,
						},
					},
				},
			},
			expected: types.OperationsStats{
				UniqueSubaccountsLiquidated:            1,
				MatchedShortTermOrdersCount:            1,
				MatchedLongTermOrdersCount:             1,
				MatchedConditionalOrdersCount:          2,
				RegularMatchesCount:                    2,
				LiquidationMatchesCount:                1,
				DeleveragingMatchesCount:               1,
				LongTermOrderRemovalsCount:             2,
				ConditionalOrderRemovalsCount:          1,
				UniqueSubaccountsDeleveraged:           1,
				UniqueSubaccountsOffsettingDeleveraged: 1,
			},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			// Run the test case.
			stats := types.StatMsgProposedOperations(tt.operations)

			// Verify the results.
			require.Equal(t, tt.expected.MatchedShortTermOrdersCount, stats.MatchedShortTermOrdersCount)
			require.Equal(t, tt.expected.MatchedLongTermOrdersCount, stats.MatchedLongTermOrdersCount)
			require.Equal(t, tt.expected.MatchedConditionalOrdersCount, stats.MatchedConditionalOrdersCount)
			require.Equal(t, tt.expected.RegularMatchesCount, stats.RegularMatchesCount)
			require.Equal(t, tt.expected.LiquidationMatchesCount, stats.LiquidationMatchesCount)
			require.Equal(t, tt.expected.DeleveragingMatchesCount, stats.DeleveragingMatchesCount)
			require.Equal(t, tt.expected.LongTermOrderRemovalsCount, stats.LongTermOrderRemovalsCount)
			require.Equal(t, tt.expected.ConditionalOrderRemovalsCount, stats.ConditionalOrderRemovalsCount)
			require.Equal(t, tt.expected.UniqueSubaccountsLiquidated, stats.UniqueSubaccountsLiquidated)
			require.Equal(t, tt.expected.UniqueSubaccountsDeleveraged, stats.UniqueSubaccountsDeleveraged)
			require.Equal(t, tt.expected.UniqueSubaccountsOffsettingDeleveraged, stats.UniqueSubaccountsOffsettingDeleveraged)
		})
	}
}
