package types_test

import (
	"testing"

	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	"github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
	"github.com/stretchr/testify/require"
)

func TestStatMsgProposedOperations(t *testing.T) {
	tests := map[string]struct {
		operations []types.OperationRaw
		expected   types.OperationsStats
	}{
		"Successfully increments each stat when seen in the list of operations": {
			operations: []types.OperationRaw{
				{
					Operation: &types.OperationRaw_Match{
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
					Operation: &types.OperationRaw_Match{
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
					Operation: &types.OperationRaw_Match{
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
					Operation: &types.OperationRaw_Match{
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
					Operation: &types.OperationRaw_OrderRemoval{
						OrderRemoval: &types.OrderRemoval{
							OrderId:       constants.LongTermOrderId_Alice_Num0_ClientId0_Clob0,
							RemovalReason: types.OrderRemoval_REMOVAL_REASON_INVALID_SELF_TRADE,
						},
					},
				},
				{
					Operation: &types.OperationRaw_OrderRemoval{
						OrderRemoval: &types.OrderRemoval{
							OrderId:       constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT5_SL_15.OrderId,
							RemovalReason: types.OrderRemoval_REMOVAL_REASON_UNDERCOLLATERALIZED,
						},
					},
				},
				{
					Operation: &types.OperationRaw_OrderRemoval{
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
				TakerOrdersCount:                       2,
				LiquidationOrdersCount:                 1,
				DeleveragingOperationsCount:            1,
				TotalFillsCount:                        3,
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
			require.Equal(t, tt.expected.TakerOrdersCount, stats.TakerOrdersCount)
			require.Equal(t, tt.expected.LiquidationOrdersCount, stats.LiquidationOrdersCount)
			require.Equal(t, tt.expected.DeleveragingOperationsCount, stats.DeleveragingOperationsCount)
			require.Equal(t, tt.expected.TotalFillsCount, stats.TotalFillsCount)
			require.Equal(t, tt.expected.LongTermOrderRemovalsCount, stats.LongTermOrderRemovalsCount)
			require.Equal(t, tt.expected.ConditionalOrderRemovalsCount, stats.ConditionalOrderRemovalsCount)
			require.Equal(t, tt.expected.UniqueSubaccountsLiquidated, stats.UniqueSubaccountsLiquidated)
			require.Equal(t, tt.expected.UniqueSubaccountsDeleveraged, stats.UniqueSubaccountsDeleveraged)
			require.Equal(t, tt.expected.UniqueSubaccountsOffsettingDeleveraged, stats.UniqueSubaccountsOffsettingDeleveraged)
		})
	}
}
