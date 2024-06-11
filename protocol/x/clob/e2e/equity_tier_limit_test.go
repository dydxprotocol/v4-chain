package clob_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/cometbft/cometbft/types"
	"github.com/dydxprotocol/v4-chain/protocol/dtypes"
	"github.com/dydxprotocol/v4-chain/protocol/lib"
	testapp "github.com/dydxprotocol/v4-chain/protocol/testutil/app"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	clobtypes "github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
	satypes "github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"
	"github.com/stretchr/testify/require"
)

func TestPlaceOrder_EquityTierLimit(t *testing.T) {
	tests := map[string]struct {
		allowedOrders                                 []clobtypes.Order
		limitedOrder                                  clobtypes.Order
		equityTierLimitConfiguration                  clobtypes.EquityTierLimitConfiguration
		cancellation                                  *clobtypes.MsgCancelOrder
		advanceBlock                                  bool
		expectError                                   bool
		crashingAppCheckTxNonDeterminsmChecksDisabled bool
	}{
		"Long-term order would exceed max open stateful orders across blocks": {
			allowedOrders: []clobtypes.Order{
				testapp.MustScaleOrder(
					constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15_StopLoss20,
					testapp.DefaultGenesis(),
				),
			},
			limitedOrder: testapp.MustScaleOrder(
				constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15,
				testapp.DefaultGenesis(),
			),
			equityTierLimitConfiguration: clobtypes.EquityTierLimitConfiguration{
				StatefulOrderEquityTiers: []clobtypes.EquityTierLimit{
					{
						UsdTncRequired: dtypes.NewInt(0),
						Limit:          0,
					},
					{
						UsdTncRequired: dtypes.NewInt(5_000_000_000), // $5,000
						Limit:          1,
					},
					{
						UsdTncRequired: dtypes.NewInt(70_000_000_000), // $70,000
						Limit:          100,
					},
				},
			},
			advanceBlock: true,
			expectError:  true,
		},
		"Long-term order would exceed max open stateful orders (due to untriggered conditional order) across blocks": {
			allowedOrders: []clobtypes.Order{
				testapp.MustScaleOrder(
					constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15_TakeProfit20,
					testapp.DefaultGenesis(),
				),
			},
			limitedOrder: testapp.MustScaleOrder(
				constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15,
				testapp.DefaultGenesis(),
			),
			equityTierLimitConfiguration: clobtypes.EquityTierLimitConfiguration{
				StatefulOrderEquityTiers: []clobtypes.EquityTierLimit{
					{
						UsdTncRequired: dtypes.NewInt(0),
						Limit:          0,
					},
					{
						UsdTncRequired: dtypes.NewInt(5_000_000_000), // $5,000
						Limit:          1,
					},
					{
						UsdTncRequired: dtypes.NewInt(70_000_000_000), // $70,000
						Limit:          100,
					},
				},
			},
			advanceBlock: true,
			expectError:  true,
		},
		"Conditional order would exceed max open stateful orders across blocks": {
			allowedOrders: []clobtypes.Order{
				testapp.MustScaleOrder(
					constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15,
					testapp.DefaultGenesis(),
				),
			},
			limitedOrder: testapp.MustScaleOrder(
				constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15_StopLoss20,
				testapp.DefaultGenesis(),
			),
			equityTierLimitConfiguration: clobtypes.EquityTierLimitConfiguration{
				StatefulOrderEquityTiers: []clobtypes.EquityTierLimit{
					{
						UsdTncRequired: dtypes.NewInt(0),
						Limit:          0,
					},
					{
						UsdTncRequired: dtypes.NewInt(5_000_000_000), // $5,000
						Limit:          1,
					},
					{
						UsdTncRequired: dtypes.NewInt(70_000_000_000), // $70,000
						Limit:          100,
					},
				},
			},
			advanceBlock: true,
			expectError:  true,
		},
		"Conditional IoC order would exceed max open stateful orders across blocks": {
			allowedOrders: []clobtypes.Order{
				testapp.MustScaleOrder(
					constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15,
					testapp.DefaultGenesis(),
				),
			},
			limitedOrder: testapp.MustScaleOrder(
				constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy5_Price50_GTBT10_StopLoss51_IOC,
				testapp.DefaultGenesis(),
			),
			equityTierLimitConfiguration: clobtypes.EquityTierLimitConfiguration{
				StatefulOrderEquityTiers: []clobtypes.EquityTierLimit{
					{
						UsdTncRequired: dtypes.NewInt(0),
						Limit:          0,
					},
					{
						UsdTncRequired: dtypes.NewInt(5_000_000_000), // $5,000
						Limit:          1,
					},
					{
						UsdTncRequired: dtypes.NewInt(70_000_000_000), // $70,000
						Limit:          100,
					},
				},
			},
			advanceBlock: true,
			expectError:  true,
		},
		"Order cancellation prevents exceeding max open stateful orders for long-term order in same block": {
			allowedOrders: []clobtypes.Order{
				testapp.MustScaleOrder(
					constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15_StopLoss20,
					testapp.DefaultGenesis(),
				),
			},
			limitedOrder: testapp.MustScaleOrder(
				constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15,
				testapp.DefaultGenesis(),
			),
			cancellation: clobtypes.NewMsgCancelOrderStateful(
				constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15_StopLoss20.OrderId,
				constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15_StopLoss20.GetGoodTilBlockTime(),
			),
			equityTierLimitConfiguration: clobtypes.EquityTierLimitConfiguration{
				StatefulOrderEquityTiers: []clobtypes.EquityTierLimit{
					{
						UsdTncRequired: dtypes.NewInt(0),
						Limit:          0,
					},
					{
						UsdTncRequired: dtypes.NewInt(5_000_000_000), // $5,000
						Limit:          1,
					},
					{
						UsdTncRequired: dtypes.NewInt(70_000_000_000), // $70,000
						Limit:          100,
					},
				},
			},
		},
		"Order cancellation of untriggered order prevents exceeding max open stateful orders for long-term order in " +
			"same block": {
			allowedOrders: []clobtypes.Order{
				testapp.MustScaleOrder(
					constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15_TakeProfit20,
					testapp.DefaultGenesis(),
				),
			},
			limitedOrder: testapp.MustScaleOrder(
				constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15,
				testapp.DefaultGenesis(),
			),
			cancellation: clobtypes.NewMsgCancelOrderStateful(
				constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15_TakeProfit20.OrderId,
				constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15_TakeProfit20.GetGoodTilBlockTime(),
			),
			equityTierLimitConfiguration: clobtypes.EquityTierLimitConfiguration{
				StatefulOrderEquityTiers: []clobtypes.EquityTierLimit{
					{
						UsdTncRequired: dtypes.NewInt(0),
						Limit:          0,
					},
					{
						UsdTncRequired: dtypes.NewInt(5_000_000_000), // $5,000
						Limit:          1,
					},
					{
						UsdTncRequired: dtypes.NewInt(70_000_000_000), // $70,000
						Limit:          100,
					},
				},
			},
		},
		"Order cancellation prevents exceeding max open stateful orders for conditional order in same block": {
			allowedOrders: []clobtypes.Order{
				testapp.MustScaleOrder(
					constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15,
					testapp.DefaultGenesis(),
				),
			},
			limitedOrder: testapp.MustScaleOrder(
				constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15_StopLoss20,
				testapp.DefaultGenesis(),
			),
			cancellation: clobtypes.NewMsgCancelOrderStateful(
				constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15.OrderId,
				constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15.GetGoodTilBlockTime(),
			),
			equityTierLimitConfiguration: clobtypes.EquityTierLimitConfiguration{
				StatefulOrderEquityTiers: []clobtypes.EquityTierLimit{
					{
						UsdTncRequired: dtypes.NewInt(0),
						Limit:          0,
					},
					{
						UsdTncRequired: dtypes.NewInt(5_000_000_000), // $5,000
						Limit:          1,
					},
					{
						UsdTncRequired: dtypes.NewInt(70_000_000_000), // $70,000
						Limit:          100,
					},
				},
			},
		},
		"Order cancellation prevents exceeding max open stateful orders for long-term order across blocks": {
			allowedOrders: []clobtypes.Order{
				testapp.MustScaleOrder(
					constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15_StopLoss20,
					testapp.DefaultGenesis(),
				),
			},
			limitedOrder: testapp.MustScaleOrder(
				constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15,
				testapp.DefaultGenesis(),
			),
			cancellation: clobtypes.NewMsgCancelOrderStateful(
				constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15_StopLoss20.OrderId,
				constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15_StopLoss20.GetGoodTilBlockTime(),
			),
			equityTierLimitConfiguration: clobtypes.EquityTierLimitConfiguration{
				StatefulOrderEquityTiers: []clobtypes.EquityTierLimit{
					{
						UsdTncRequired: dtypes.NewInt(0),
						Limit:          0,
					},
					{
						UsdTncRequired: dtypes.NewInt(5_000_000_000), // $5,000
						Limit:          1,
					},
					{
						UsdTncRequired: dtypes.NewInt(70_000_000_000), // $70,000
						Limit:          100,
					},
				},
			},
			advanceBlock: true,
		},
		"Order cancellation of untriggered order prevents exceeding max open stateful orders for long-term order " +
			"across blocks": {
			allowedOrders: []clobtypes.Order{
				testapp.MustScaleOrder(
					constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15_TakeProfit20,
					testapp.DefaultGenesis(),
				),
			},
			limitedOrder: testapp.MustScaleOrder(
				constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15,
				testapp.DefaultGenesis(),
			),
			cancellation: clobtypes.NewMsgCancelOrderStateful(
				constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15_TakeProfit20.OrderId,
				constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15_TakeProfit20.GetGoodTilBlockTime(),
			),
			equityTierLimitConfiguration: clobtypes.EquityTierLimitConfiguration{
				StatefulOrderEquityTiers: []clobtypes.EquityTierLimit{
					{
						UsdTncRequired: dtypes.NewInt(0),
						Limit:          0,
					},
					{
						UsdTncRequired: dtypes.NewInt(5_000_000_000), // $5,000
						Limit:          1,
					},
					{
						UsdTncRequired: dtypes.NewInt(70_000_000_000), // $70,000
						Limit:          100,
					},
				},
			},
			advanceBlock: true,
		},
		"Order cancellation prevents exceeding max open stateful orders for conditional order across blocks": {
			allowedOrders: []clobtypes.Order{
				testapp.MustScaleOrder(
					constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15,
					testapp.DefaultGenesis(),
				),
			},
			limitedOrder: testapp.MustScaleOrder(
				constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15_StopLoss20,
				testapp.DefaultGenesis(),
			),
			cancellation: clobtypes.NewMsgCancelOrderStateful(
				constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15.OrderId,
				constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15.GetGoodTilBlockTime(),
			),
			equityTierLimitConfiguration: clobtypes.EquityTierLimitConfiguration{
				StatefulOrderEquityTiers: []clobtypes.EquityTierLimit{
					{
						UsdTncRequired: dtypes.NewInt(0),
						Limit:          0,
					},
					{
						UsdTncRequired: dtypes.NewInt(5_000_000_000), // $5,000
						Limit:          1,
					},
					{
						UsdTncRequired: dtypes.NewInt(70_000_000_000), // $70,000
						Limit:          100,
					},
				},
			},
			advanceBlock: true,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tApp := testapp.NewTestAppBuilder(t).
				WithCrashingAppCheckTxNonDeterminismChecksEnabled(!tc.crashingAppCheckTxNonDeterminsmChecksDisabled).
				WithGenesisDocFn(func() types.GenesisDoc {
					genesis := testapp.DefaultGenesis()
					testapp.UpdateGenesisDocWithAppStateForModule(&genesis, func(state *satypes.GenesisState) {
						state.Subaccounts = []satypes.Subaccount{
							constants.Alice_Num0_10_000USD,
						}
					})
					testapp.UpdateGenesisDocWithAppStateForModule(&genesis, func(state *clobtypes.GenesisState) {
						state.EquityTierLimitConfig = tc.equityTierLimitConfiguration
						// Don't enforce the block rate limit.
						state.BlockRateLimitConfig = clobtypes.BlockRateLimitConfiguration{}
					})
					return genesis
				}).Build()

			ctx := tApp.InitChain()

			for _, allowedOrder := range tc.allowedOrders {
				for _, tx := range testapp.MustMakeCheckTxsWithClobMsg(ctx, tApp.App, *clobtypes.NewMsgPlaceOrder(allowedOrder)) {
					resp := tApp.CheckTx(tx)
					require.Conditionf(t, resp.IsOK, "Expected CheckTx to succeed. Response: %+v", resp)
				}
			}

			if tc.advanceBlock {
				ctx = tApp.AdvanceToBlock(2, testapp.AdvanceToBlockOptions{})
			}

			if tc.cancellation != nil {
				for _, tx := range testapp.MustMakeCheckTxsWithClobMsg(ctx, tApp.App, *tc.cancellation) {
					resp := tApp.CheckTx(tx)
					require.Conditionf(t, resp.IsOK, "Expected CheckTx to succeed. Response: %+v", resp)
				}
			}

			if tc.advanceBlock {
				ctx = tApp.AdvanceToBlock(3, testapp.AdvanceToBlockOptions{})
			}

			for _, tx := range testapp.MustMakeCheckTxsWithClobMsg(ctx, tApp.App, *clobtypes.NewMsgPlaceOrder(tc.limitedOrder)) {
				resp := tApp.CheckTx(tx)
				if tc.expectError {
					require.Conditionf(t, resp.IsErr, "Expected CheckTx to error. Response: %+v", resp)
					require.Contains(
						t,
						resp.Log,
						fmt.Sprintf(
							"Opening order would exceed equity tier limit of %d. Order count: %d,",
							len(tc.allowedOrders),
							len(tc.allowedOrders),
						),
					)
				} else {
					require.Conditionf(t, resp.IsOK, "Expected CheckTx to succeed. Response: %+v", resp)
				}
			}

			// Ensure that any successful transactions can be delivered.
			tApp.AdvanceToBlock(4, testapp.AdvanceToBlockOptions{})
		})
	}
}

func TestPlaceOrder_EquityTierLimit_OrderExpiry(t *testing.T) {
	tests := map[string]struct {
		firstOrder                                    clobtypes.Order
		secondOrder                                   clobtypes.Order
		equityTierLimitConfiguration                  clobtypes.EquityTierLimitConfiguration
		advanceToBlockAndTime                         uint32
		expectError                                   bool
		crashingAppCheckTxNonDeterminsmChecksDisabled bool
	}{
		"Stateful order has not expired": {
			firstOrder: testapp.MustScaleOrder(
				constants.LongTermOrder_Alice_Num0_Id0_Clob1_Buy5_Price10_GTBT5,
				testapp.DefaultGenesis(),
			),
			secondOrder: testapp.MustScaleOrder(
				constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT20,
				testapp.DefaultGenesis(),
			),
			equityTierLimitConfiguration: clobtypes.EquityTierLimitConfiguration{
				StatefulOrderEquityTiers: []clobtypes.EquityTierLimit{
					{
						UsdTncRequired: dtypes.NewInt(0),
						Limit:          0,
					},
					{
						UsdTncRequired: dtypes.NewInt(5_000_000_000), // $5,000
						Limit:          1,
					},
					{
						UsdTncRequired: dtypes.NewInt(70_000_000_000), // $70,000
						Limit:          100,
					},
				},
			},
			advanceToBlockAndTime: 4,
			expectError:           true,
		},
		"Stateful order has expired": {
			firstOrder: testapp.MustScaleOrder(
				constants.LongTermOrder_Alice_Num0_Id0_Clob1_Buy5_Price10_GTBT5,
				testapp.DefaultGenesis(),
			),
			secondOrder: testapp.MustScaleOrder(
				constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT20,
				testapp.DefaultGenesis(),
			),
			equityTierLimitConfiguration: clobtypes.EquityTierLimitConfiguration{
				StatefulOrderEquityTiers: []clobtypes.EquityTierLimit{
					{
						UsdTncRequired: dtypes.NewInt(0),
						Limit:          0,
					},
					{
						UsdTncRequired: dtypes.NewInt(5_000_000_000), // $5,000
						Limit:          1,
					},
					{
						UsdTncRequired: dtypes.NewInt(70_000_000_000), // $70,000
						Limit:          100,
					},
				},
			},
			advanceToBlockAndTime: 5,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tApp := testapp.NewTestAppBuilder(t).
				WithCrashingAppCheckTxNonDeterminismChecksEnabled(!tc.crashingAppCheckTxNonDeterminsmChecksDisabled).
				WithGenesisDocFn(func() types.GenesisDoc {
					genesis := testapp.DefaultGenesis()
					testapp.UpdateGenesisDocWithAppStateForModule(&genesis, func(state *satypes.GenesisState) {
						state.Subaccounts = []satypes.Subaccount{
							constants.Alice_Num0_10_000USD,
						}
					})
					testapp.UpdateGenesisDocWithAppStateForModule(&genesis, func(state *clobtypes.GenesisState) {
						state.EquityTierLimitConfig = tc.equityTierLimitConfiguration
					})
					return genesis
				}).Build()

			ctx := tApp.InitChain()

			for _, tx := range testapp.MustMakeCheckTxsWithClobMsg(ctx, tApp.App, *clobtypes.NewMsgPlaceOrder(tc.firstOrder)) {
				resp := tApp.CheckTx(tx)
				require.Conditionf(t, resp.IsOK, "Expected CheckTx to succeed. Response: %+v", resp)
			}

			ctx = tApp.AdvanceToBlock(
				tc.advanceToBlockAndTime,
				testapp.AdvanceToBlockOptions{BlockTime: time.Unix(int64(tc.advanceToBlockAndTime), 0)},
			)

			for _, tx := range testapp.MustMakeCheckTxsWithClobMsg(ctx, tApp.App, *clobtypes.NewMsgPlaceOrder(tc.secondOrder)) {
				resp := tApp.CheckTx(tx)
				if tc.expectError {
					require.Conditionf(t, resp.IsErr, "Expected CheckTx to error. Response: %+v", resp)
					require.Contains(t, resp.Log, "Opening order would exceed equity tier limit of 1. Order count: 1,")
				} else {
					require.Conditionf(t, resp.IsOK, "Expected CheckTx to succeed. Response: %+v", resp)
				}
			}

			// Ensure that any successful transactions can be delivered.
			tApp.AdvanceToBlock(lib.MustConvertIntegerToUint32(tApp.GetBlockHeight()+1), testapp.AdvanceToBlockOptions{})
		})
	}
}

func TestPlaceOrder_EquityTierLimit_OrderFill(t *testing.T) {
	tests := map[string]struct {
		makerOrder                                    clobtypes.Order
		takerOrder                                    clobtypes.Order
		extraOrder                                    clobtypes.Order
		equityTierLimitConfiguration                  clobtypes.EquityTierLimitConfiguration
		cancellation                                  *clobtypes.MsgCancelOrder
		advanceBlock                                  bool
		expectError                                   bool
		crashingAppCheckTxNonDeterminsmChecksDisabled bool
	}{
		"Order fully filled prevents exceeding max open stateful orders for conditional order across blocks": {
			makerOrder: testapp.MustScaleOrder(
				constants.LongTermOrder_Bob_Num0_Id0_Clob0_Sell5_Price5_GTBT10,
				testapp.DefaultGenesis(),
			),
			takerOrder: testapp.MustScaleOrder(
				constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15,
				testapp.DefaultGenesis(),
			),
			extraOrder: testapp.MustScaleOrder(
				constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15_StopLoss20,
				testapp.DefaultGenesis(),
			),
			equityTierLimitConfiguration: clobtypes.EquityTierLimitConfiguration{
				StatefulOrderEquityTiers: []clobtypes.EquityTierLimit{
					{
						UsdTncRequired: dtypes.NewInt(0),
						Limit:          0,
					},
					{
						UsdTncRequired: dtypes.NewInt(5_000_000_000), // $5,000
						Limit:          1,
					},
					{
						UsdTncRequired: dtypes.NewInt(70_000_000_000), // $70,000
						Limit:          100,
					},
				},
			},
			advanceBlock: true,
		},
		"Order fully filled prevents exceeding max open stateful orders for long-term order across blocks": {
			makerOrder: testapp.MustScaleOrder(
				constants.LongTermOrder_Bob_Num0_Id0_Clob0_Sell5_Price5_GTBT10,
				testapp.DefaultGenesis(),
			),
			takerOrder: testapp.MustScaleOrder(
				constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15_StopLoss20,
				testapp.DefaultGenesis(),
			),
			extraOrder: testapp.MustScaleOrder(
				constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15,
				testapp.DefaultGenesis(),
			),
			equityTierLimitConfiguration: clobtypes.EquityTierLimitConfiguration{
				StatefulOrderEquityTiers: []clobtypes.EquityTierLimit{
					{
						UsdTncRequired: dtypes.NewInt(0),
						Limit:          0,
					},
					{
						UsdTncRequired: dtypes.NewInt(5_000_000_000), // $5,000
						Limit:          1,
					},
					{
						UsdTncRequired: dtypes.NewInt(70_000_000_000), // $70,000
						Limit:          100,
					},
				},
			},
			advanceBlock: true,
		},
		"Order partially filled exceeds max open stateful orders for conditional order across blocks": {
			makerOrder: testapp.MustScaleOrder(
				constants.LongTermOrder_Bob_Num0_Id0_Clob0_Sell2_Price5_GTBT10,
				testapp.DefaultGenesis(),
			),
			takerOrder: testapp.MustScaleOrder(
				constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15,
				testapp.DefaultGenesis(),
			),
			extraOrder: testapp.MustScaleOrder(
				constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15_StopLoss20,
				testapp.DefaultGenesis(),
			),
			equityTierLimitConfiguration: clobtypes.EquityTierLimitConfiguration{
				StatefulOrderEquityTiers: []clobtypes.EquityTierLimit{
					{
						UsdTncRequired: dtypes.NewInt(0),
						Limit:          0,
					},
					{
						UsdTncRequired: dtypes.NewInt(5_000_000_000), // $5,000
						Limit:          1,
					},
					{
						UsdTncRequired: dtypes.NewInt(70_000_000_000), // $70,000
						Limit:          100,
					},
				},
			},
			advanceBlock: true,
			expectError:  true,
		},
		"Order partially filled exceeds max open stateful orders for long-term order across blocks": {
			makerOrder: testapp.MustScaleOrder(
				constants.LongTermOrder_Bob_Num0_Id0_Clob0_Sell2_Price5_GTBT10,
				testapp.DefaultGenesis(),
			),
			takerOrder: testapp.MustScaleOrder(
				constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15_StopLoss20,
				testapp.DefaultGenesis(),
			),
			extraOrder: testapp.MustScaleOrder(
				constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15,
				testapp.DefaultGenesis(),
			),
			equityTierLimitConfiguration: clobtypes.EquityTierLimitConfiguration{
				StatefulOrderEquityTiers: []clobtypes.EquityTierLimit{
					{
						UsdTncRequired: dtypes.NewInt(0),
						Limit:          0,
					},
					{
						UsdTncRequired: dtypes.NewInt(5_000_000_000), // $5,000
						Limit:          1,
					},
					{
						UsdTncRequired: dtypes.NewInt(70_000_000_000), // $70,000
						Limit:          100,
					},
				},
			},
			advanceBlock: true,
			expectError:  true,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tApp := testapp.NewTestAppBuilder(t).
				WithCrashingAppCheckTxNonDeterminismChecksEnabled(!tc.crashingAppCheckTxNonDeterminsmChecksDisabled).
				WithGenesisDocFn(func() types.GenesisDoc {
					genesis := testapp.DefaultGenesis()
					testapp.UpdateGenesisDocWithAppStateForModule(&genesis, func(state *satypes.GenesisState) {
						state.Subaccounts = []satypes.Subaccount{
							constants.Alice_Num0_10_000USD,
							constants.Bob_Num0_100_000USD,
						}
					})
					testapp.UpdateGenesisDocWithAppStateForModule(&genesis, func(state *clobtypes.GenesisState) {
						state.EquityTierLimitConfig = tc.equityTierLimitConfiguration
					})
					return genesis
				}).Build()

			ctx := tApp.InitChain()

			for _, tx := range testapp.MustMakeCheckTxsWithClobMsg(ctx, tApp.App, *clobtypes.NewMsgPlaceOrder(tc.makerOrder)) {
				resp := tApp.CheckTx(tx)
				require.Conditionf(t, resp.IsOK, "Expected CheckTx to succeed. Response: %+v", resp)
			}
			if tc.advanceBlock {
				ctx = tApp.AdvanceToBlock(2, testapp.AdvanceToBlockOptions{})
			}

			for _, tx := range testapp.MustMakeCheckTxsWithClobMsg(ctx, tApp.App, *clobtypes.NewMsgPlaceOrder(tc.takerOrder)) {
				resp := tApp.CheckTx(tx)
				require.Conditionf(t, resp.IsOK, "Expected CheckTx to succeed. Response: %+v", resp)
			}
			if tc.advanceBlock {
				ctx = tApp.AdvanceToBlock(4, testapp.AdvanceToBlockOptions{})
			}

			for _, tx := range testapp.MustMakeCheckTxsWithClobMsg(ctx, tApp.App, *clobtypes.NewMsgPlaceOrder(tc.extraOrder)) {
				resp := tApp.CheckTx(tx)
				if tc.expectError {
					require.Conditionf(t, resp.IsErr, "Expected CheckTx to error. Response: %+v", resp)
					require.Contains(t, resp.Log, "Opening order would exceed equity tier limit of 1. Order count: 1,")
				} else {
					require.Conditionf(t, resp.IsOK, "Expected CheckTx to succeed. Response: %+v", resp)
				}
			}

			// Ensure that any successful transactions can be delivered.
			tApp.AdvanceToBlock(5, testapp.AdvanceToBlockOptions{})
		})
	}
}
