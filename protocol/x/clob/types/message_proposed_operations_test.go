package types_test

import (
	"errors"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	clobtestutils "github.com/dydxprotocol/v4-chain/protocol/testutil/clob"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	testtx "github.com/dydxprotocol/v4-chain/protocol/testutil/tx"
	"github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
	satypes "github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"
	"github.com/stretchr/testify/require"
)

func TestValidateBasic(t *testing.T) {
	tests := map[string]struct {
		msg           types.MsgProposedOperations
		expectedError error
	}{
		"valid message": {
			msg: types.MsgProposedOperations{
				OperationsQueue: []types.OperationRaw{
					{
						Operation: &types.OperationRaw_ShortTermOrderPlacement{
							ShortTermOrderPlacement: []byte{1, 2, 3},
						},
					},
					{
						Operation: &types.OperationRaw_ShortTermOrderPlacement{
							ShortTermOrderPlacement: []byte{4, 5, 6},
						},
					},
					clobtestutils.NewMatchOperationRaw(
						&constants.Order_Alice_Num0_Id0_Clob0_Buy10_Price10_GTB16,
						[]types.MakerFill{
							{
								FillAmount:   5,
								MakerOrderId: constants.Order_Bob_Num0_Id8_Clob0_Sell20_Price10_GTB22.OrderId,
							},
						},
					),
					{
						Operation: &types.OperationRaw_OrderRemoval{
							OrderRemoval: &types.OrderRemoval{
								OrderId:       constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy100_Price10_GTBT15.OrderId,
								RemovalReason: types.OrderRemoval_REMOVAL_REASON_INVALID_SELF_TRADE,
							},
						},
					},
				},
			},
		},
		"short term order removal returns error": {
			msg: types.MsgProposedOperations{
				OperationsQueue: []types.OperationRaw{
					{
						Operation: &types.OperationRaw_OrderRemoval{
							OrderRemoval: &types.OrderRemoval{
								OrderId:       constants.Order_Alice_Num0_Id0_Clob0_Buy10_Price10_GTB16.OrderId,
								RemovalReason: types.OrderRemoval_REMOVAL_REASON_INVALID_SELF_TRADE,
							},
						},
					},
				},
			},
			expectedError: errors.New("order removal is not allowed for short-term orders"),
		},
		"unspecified removal reason returns error": {
			msg: types.MsgProposedOperations{
				OperationsQueue: []types.OperationRaw{
					{
						Operation: &types.OperationRaw_OrderRemoval{
							OrderRemoval: &types.OrderRemoval{
								OrderId:       constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy100_Price10_GTBT15.OrderId,
								RemovalReason: types.OrderRemoval_REMOVAL_REASON_UNSPECIFIED,
							},
						},
					},
				},
			},
			expectedError: errors.New("order removal reason must be specified"),
		},
		"reduce-only removal reason returns error": {
			msg: types.MsgProposedOperations{
				OperationsQueue: []types.OperationRaw{
					{
						Operation: &types.OperationRaw_OrderRemoval{
							OrderRemoval: &types.OrderRemoval{
								OrderId:       constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy100_Price10_GTBT15.OrderId,
								RemovalReason: types.OrderRemoval_REMOVAL_REASON_INVALID_REDUCE_ONLY,
							},
						},
					},
				},
			},
			expectedError: nil,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			err := tc.msg.ValidateBasic()
			if tc.expectedError != nil {
				require.ErrorContains(t, err, tc.expectedError.Error())
			} else {
				require.NoError(t, err)
			}
		})
	}
}

// TestValidateBasic provides black box testing for the ValidateBasic method of the MsgProposedOperations msg.
func TestValidateAndTransformRawOperations(t *testing.T) {
	tests := map[string]struct {
		operations    []types.OperationRaw
		expectedError error
	}{
		// Tests for functionality not related to one message type
		"passes validation": {
			operations: []types.OperationRaw{
				clobtestutils.NewShortTermOrderPlacementOperationRaw(
					constants.Order_Carl_Num0_Id0_Clob0_Buy1BTC_Price50000_GTB10,
				),
				clobtestutils.NewShortTermOrderPlacementOperationRaw(
					constants.Order_Dave_Num0_Id0_Clob0_Sell1BTC_Price50000_GTB10,
				),
				clobtestutils.NewMatchOperationRaw(
					&constants.Order_Carl_Num0_Id0_Clob0_Buy1BTC_Price50000_GTB10,
					[]types.MakerFill{
						{
							FillAmount:   100_000_000, // 1 BTC
							MakerOrderId: constants.Order_Dave_Num0_Id0_Clob0_Sell1BTC_Price50000_GTB10.GetOrderId(),
						},
					},
				),
				clobtestutils.NewOrderRemovalOperationRaw(
					constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy100_Price10_GTBT15.OrderId,
					types.OrderRemoval_REMOVAL_REASON_INVALID_SELF_TRADE,
				),
			},
			expectedError: nil,
		},

		// tests for invalid subaccount id
		"Stateless order validation: Place Order has invalid SubaccountId": {
			operations: []types.OperationRaw{
				clobtestutils.NewShortTermOrderPlacementOperationRaw(types.Order{
					OrderId:      constants.InvalidSubaccountIdOwner_OrderId,
					Side:         types.Order_SIDE_BUY,
					Quantums:     100_000_000,
					Subticks:     50_000_000_000,
					GoodTilOneof: &types.Order_GoodTilBlock{GoodTilBlock: 10},
				}),
			},
			expectedError: errors.New("invalid SubaccountId Owner address"),
		},
		"Stateless order validation: matchOrders Taker order id invalid SubaccountId": {
			operations: []types.OperationRaw{
				clobtestutils.NewShortTermOrderPlacementOperationRaw(constants.Order_Dave_Num0_Id0_Clob0_Sell1BTC_Price50000),
				{
					Operation: &types.OperationRaw_Match{
						Match: &types.ClobMatch{
							Match: &types.ClobMatch_MatchOrders{
								MatchOrders: &types.MatchOrders{
									TakerOrderId: constants.InvalidSubaccountIdNumber_OrderId,
									Fills: []types.MakerFill{
										{
											FillAmount:   1,
											MakerOrderId: constants.Order_Dave_Num0_Id0_Clob0_Sell1BTC_Price50000.GetOrderId(),
										},
									},
								},
							},
						},
					},
				},
			},
			expectedError: satypes.ErrInvalidSubaccountIdNumber,
		},
		"Stateless order validation: matchOrders Maker order id invalid SubaccountId": {
			operations: []types.OperationRaw{
				clobtestutils.NewShortTermOrderPlacementOperationRaw(constants.Order_Dave_Num0_Id0_Clob0_Sell1BTC_Price50000),
				{
					Operation: &types.OperationRaw_Match{
						Match: &types.ClobMatch{
							Match: &types.ClobMatch_MatchOrders{
								MatchOrders: &types.MatchOrders{
									TakerOrderId: constants.OrderId_Alice_Num0_ClientId0_Clob0,
									Fills: []types.MakerFill{
										{
											FillAmount:   1,
											MakerOrderId: constants.Order_Dave_Num0_Id0_Clob0_Sell1BTC_Price50000.GetOrderId(),
										},
										{
											FillAmount:   1,
											MakerOrderId: constants.InvalidSubaccountIdOwner_OrderId,
										},
									},
								},
							},
						},
					},
				},
			},
			expectedError: errors.New("invalid SubaccountId Owner address"),
		},
		"Stateless order validation: perpLiquidation fill maker order id invalid SubaccountId": {
			operations: []types.OperationRaw{
				clobtestutils.NewShortTermOrderPlacementOperationRaw(constants.Order_Dave_Num0_Id0_Clob0_Sell1BTC_Price50000),
				{
					Operation: &types.OperationRaw_Match{
						Match: &types.ClobMatch{
							Match: &types.ClobMatch_MatchPerpetualLiquidation{
								MatchPerpetualLiquidation: &types.MatchPerpetualLiquidation{
									Liquidated:  constants.Carl_Num0,
									ClobPairId:  0,
									PerpetualId: 0,
									TotalSize:   1,
									IsBuy:       true,
									Fills: []types.MakerFill{
										{
											MakerOrderId: constants.InvalidSubaccountIdOwner_OrderId,
											FillAmount:   1,
										},
									},
								},
							},
						},
					},
				},
			},
			expectedError: errors.New("invalid SubaccountId Owner address"),
		},

		// tests for Order Placement
		"Stateless place order validation: replacement order higher priority": {
			operations: []types.OperationRaw{
				clobtestutils.NewShortTermOrderPlacementOperationRaw(constants.Order_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB20),
				clobtestutils.NewShortTermOrderPlacementOperationRaw(constants.Order_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB15),
			},
			expectedError: errors.New("Replacement order is not higher priority"),
		},
		"Stateless place order validation: placeOrder has invalid side": {
			operations: []types.OperationRaw{
				clobtestutils.NewShortTermOrderPlacementOperationRaw(types.Order{
					OrderId:      types.OrderId{SubaccountId: constants.Alice_Num0, ClientId: 0, ClobPairId: 0},
					Side:         types.Order_Side(uint32(999)),
					Quantums:     100_000_000,
					Subticks:     50_000_000_000,
					GoodTilOneof: &types.Order_GoodTilBlock{GoodTilBlock: 10},
				}),
				clobtestutils.NewShortTermOrderPlacementOperationRaw(
					constants.Order_Dave_Num0_Id0_Clob0_Sell1BTC_Price50000_GTB10,
				),
			},
			expectedError: errors.New("invalid order side"),
		},
		"Stateless place order validation: placeOrder has unspecified side": {
			operations: []types.OperationRaw{
				clobtestutils.NewShortTermOrderPlacementOperationRaw(types.Order{
					OrderId:      types.OrderId{SubaccountId: constants.Alice_Num0, ClientId: 0, ClobPairId: 0},
					Quantums:     100_000_000,
					Subticks:     50_000_000_000,
					GoodTilOneof: &types.Order_GoodTilBlock{GoodTilBlock: 10},
				}),
				clobtestutils.NewShortTermOrderPlacementOperationRaw(
					constants.Order_Dave_Num0_Id0_Clob0_Sell1BTC_Price50000_GTB10,
				),
			},
			expectedError: errors.New("UNSPECIFIED is not a valid order side"),
		},
		"Stateless place order validation: no duplicate order placements": {
			operations: []types.OperationRaw{
				clobtestutils.NewShortTermOrderPlacementOperationRaw(
					constants.Order_Dave_Num0_Id0_Clob0_Sell1BTC_Price50000_GTB10,
				),
				clobtestutils.NewShortTermOrderPlacementOperationRaw(
					constants.Order_Dave_Num0_Id0_Clob0_Sell1BTC_Price50000_GTB10,
				),
			},
			expectedError: errors.New("Duplicate Order"),
		},
		"Stateless place order validation: order quantums cannot be 0": {
			operations: []types.OperationRaw{
				clobtestutils.NewShortTermOrderPlacementOperationRaw(types.Order{
					OrderId:      types.OrderId{SubaccountId: constants.Alice_Num0, ClientId: 0, ClobPairId: 0},
					Side:         types.Order_SIDE_BUY,
					Quantums:     0,
					Subticks:     50_000_000_000,
					GoodTilOneof: &types.Order_GoodTilBlock{GoodTilBlock: 10},
				}),
				clobtestutils.NewShortTermOrderPlacementOperationRaw(
					constants.Order_Dave_Num0_Id0_Clob0_Sell1BTC_Price50000_GTB10,
				),
			},
			expectedError: errors.New("order size quantums cannot be 0"),
		},
		"Stateless place order validation: order goodTilBlock cannot be 0": {
			operations: []types.OperationRaw{
				clobtestutils.NewShortTermOrderPlacementOperationRaw(types.Order{
					OrderId:      types.OrderId{SubaccountId: constants.Alice_Num0, ClientId: 0, ClobPairId: 0},
					Side:         types.Order_SIDE_BUY,
					Quantums:     10,
					Subticks:     50_000_000_000,
					GoodTilOneof: &types.Order_GoodTilBlock{GoodTilBlock: 0},
				}),
				clobtestutils.NewShortTermOrderPlacementOperationRaw(
					constants.Order_Dave_Num0_Id0_Clob0_Sell1BTC_Price50000_GTB10,
				),
			},
			expectedError: errors.New("order goodTilBlock cannot be 0"),
		},
		"Stateless place order validation: order subticks cannot be 0": {
			operations: []types.OperationRaw{
				clobtestutils.NewShortTermOrderPlacementOperationRaw(types.Order{
					OrderId:      types.OrderId{SubaccountId: constants.Alice_Num0, ClientId: 0, ClobPairId: 0},
					Side:         types.Order_SIDE_BUY,
					Quantums:     10,
					Subticks:     0,
					GoodTilOneof: &types.Order_GoodTilBlock{GoodTilBlock: 10},
				}),
				clobtestutils.NewShortTermOrderPlacementOperationRaw(
					constants.Order_Dave_Num0_Id0_Clob0_Sell1BTC_Price50000_GTB10,
				),
			},
			expectedError: errors.New("order subticks cannot be 0"),
		},

		// tests for Match Orders
		"Stateless match order validation: match contains no fills": {
			operations: []types.OperationRaw{
				clobtestutils.NewShortTermOrderPlacementOperationRaw(constants.Order_Dave_Num0_Id0_Clob0_Sell1BTC_Price50000),
				clobtestutils.NewShortTermOrderPlacementOperationRaw(constants.Order_Carl_Num1_Id0_Clob0_Buy1BTC_Price50000),
				clobtestutils.NewMatchOperationRaw(
					&constants.Order_Carl_Num1_Id0_Clob0_Buy1BTC_Price50000,
					[]types.MakerFill{},
				),
			},
			expectedError: types.ErrInvalidMatchOrder,
		},
		"Stateless match order validation: fill amount is zero": {
			operations: []types.OperationRaw{
				clobtestutils.NewShortTermOrderPlacementOperationRaw(constants.Order_Dave_Num0_Id0_Clob0_Sell1BTC_Price50000),
				clobtestutils.NewShortTermOrderPlacementOperationRaw(constants.Order_Carl_Num1_Id0_Clob0_Buy1BTC_Price50000),
				clobtestutils.NewMatchOperationRaw(
					&constants.Order_Carl_Num1_Id0_Clob0_Buy1BTC_Price50000,
					[]types.MakerFill{
						{
							FillAmount:   0, // zero
							MakerOrderId: constants.Order_Dave_Num0_Id0_Clob0_Sell1BTC_Price50000.GetOrderId(),
						},
					},
				),
			},
			expectedError: types.ErrFillAmountIsZero,
		},
		"Stateless match order validation: Duplicate Maker OrderId in Fill List": {
			operations: []types.OperationRaw{
				clobtestutils.NewShortTermOrderPlacementOperationRaw(constants.Order_Dave_Num0_Id0_Clob0_Sell1BTC_Price50000),
				clobtestutils.NewShortTermOrderPlacementOperationRaw(constants.Order_Carl_Num1_Id0_Clob0_Buy1BTC_Price50000),
				clobtestutils.NewMatchOperationRaw(
					&constants.Order_Carl_Num1_Id0_Clob0_Buy1BTC_Price50000,
					[]types.MakerFill{
						{
							FillAmount:   1,
							MakerOrderId: constants.Order_Dave_Num0_Id0_Clob0_Sell1BTC_Price50000.GetOrderId(),
						},
						{
							FillAmount:   1,
							MakerOrderId: constants.Order_Dave_Num0_Id0_Clob0_Sell1BTC_Price50000.GetOrderId(),
						},
					},
				),
			},
			expectedError: errors.New("duplicate Maker OrderId in a MatchOrder's fills"),
		},

		// tests for Perpetual Liquidations
		"Stateless liquidation validation: fails if total fill amount exceeds order size": {
			operations: []types.OperationRaw{
				clobtestutils.NewShortTermOrderPlacementOperationRaw(types.Order{
					OrderId:      constants.OrderId_Alice_Num0_ClientId0_Clob0,
					Side:         types.Order_SIDE_SELL,
					Quantums:     150_000_000,
					Subticks:     50_000_000_000,
					GoodTilOneof: &types.Order_GoodTilBlock{GoodTilBlock: 11},
				}),
				clobtestutils.NewMatchOperationRawFromPerpetualLiquidation(types.MatchPerpetualLiquidation{
					Liquidated:  constants.Carl_Num0,
					ClobPairId:  0,
					PerpetualId: 0,
					TotalSize:   100_000_000, // 1 BTCw
					IsBuy:       true,
					Fills: []types.MakerFill{
						{
							MakerOrderId: constants.OrderId_Alice_Num0_ClientId0_Clob0,
							FillAmount:   50_000_000, // .50 BTC does not exceed order quantums
						},
						{
							MakerOrderId: constants.OrderId_Alice_Num0_ClientId0_Clob0,
							FillAmount:   50_000_000, // .50 BTC does not exceed order quantums
						},
						{
							MakerOrderId: constants.OrderId_Alice_Num0_ClientId0_Clob0,
							FillAmount:   50_000_000, // another .50 BTC EXCEEDS liquidation order size
						},
					},
				}),
			},
			expectedError: types.ErrTotalFillAmountExceedsOrderSize,
		},
		"Stateless liquidation validation: fails when total fill amount is zero": {
			operations: []types.OperationRaw{
				clobtestutils.NewShortTermOrderPlacementOperationRaw(constants.Order_Dave_Num0_Id0_Clob0_Sell1BTC_Price50000),
				clobtestutils.NewMatchOperationRawFromPerpetualLiquidation(types.MatchPerpetualLiquidation{
					Liquidated:  constants.Carl_Num0,
					ClobPairId:  0,
					PerpetualId: 0,
					TotalSize:   0, // Total size is zero.
					IsBuy:       true,
					Fills: []types.MakerFill{
						{
							MakerOrderId: constants.Order_Dave_Num0_Id0_Clob0_Sell1BTC_Price50000.GetOrderId(),
							FillAmount:   100,
						},
					},
				}),
			},
			expectedError: types.ErrInvalidLiquidationOrderTotalSize,
		},
		"Stateless liquidation validation: fails when fill amount is zero": {
			operations: []types.OperationRaw{
				clobtestutils.NewShortTermOrderPlacementOperationRaw(constants.Order_Dave_Num0_Id0_Clob0_Sell1BTC_Price50000),
				clobtestutils.NewMatchOperationRawFromPerpetualLiquidation(types.MatchPerpetualLiquidation{
					Liquidated:  constants.Carl_Num0,
					ClobPairId:  0,
					PerpetualId: 0,
					TotalSize:   100,
					IsBuy:       true,
					Fills: []types.MakerFill{
						{
							MakerOrderId: constants.Order_Dave_Num0_Id0_Clob0_Sell1BTC_Price50000.GetOrderId(),
							FillAmount:   0, // fill amount is zero
						},
					},
				}),
			},
			expectedError: types.ErrFillAmountIsZero,
		},
		"Stateless liquidation validation: fails when liquidation match contains no fills": {
			operations: []types.OperationRaw{
				clobtestutils.NewShortTermOrderPlacementOperationRaw(constants.Order_Dave_Num0_Id0_Clob0_Sell1BTC_Price50000),
				clobtestutils.NewMatchOperationRawFromPerpetualLiquidation(types.MatchPerpetualLiquidation{
					Liquidated:  constants.Carl_Num0,
					ClobPairId:  0,
					PerpetualId: 0,
					TotalSize:   100,
					IsBuy:       true,
					Fills:       []types.MakerFill{},
				}),
			},
			expectedError: types.ErrInvalidMatchOrder,
		},

		// Tests for Match Perpetual Deleveraging
		"Stateless match perpetual deleveraging validation: forwards errors from validate": {
			operations: []types.OperationRaw{
				clobtestutils.NewMatchOperationRawFromPerpetualDeleveragingLiquidation(
					types.MatchPerpetualDeleveraging{
						Liquidated:  constants.Alice_Num0,
						PerpetualId: constants.ClobPair_Eth.MustGetPerpetualId(),
						Fills:       []types.MatchPerpetualDeleveraging_Fill{},
					},
				),
			},
			expectedError: nil,
		},

		// Tests for byte functionality
		"Short term order placement tx bytes fail to decode": {
			operations: []types.OperationRaw{
				{
					Operation: &types.OperationRaw_ShortTermOrderPlacement{
						ShortTermOrderPlacement: []byte("invalid"),
					},
				},
			},
			expectedError: errors.New("tx parse error"),
		},
		"Short term order placement tx bytes contains too many messages": {
			operations: []types.OperationRaw{
				{
					Operation: &types.OperationRaw_ShortTermOrderPlacement{
						ShortTermOrderPlacement: testtx.MustGetTxBytes(
							constants.Msg_PlaceOrder,
							constants.Msg_PlaceOrder,
						),
					},
				},
			},
			expectedError: errors.New("expected 1 msg, got 2"),
		},
		"Short term order placement tx bytes contains a cancel instead of placement": {
			operations: []types.OperationRaw{
				{
					Operation: &types.OperationRaw_ShortTermOrderPlacement{
						ShortTermOrderPlacement: testtx.MustGetTxBytes(
							constants.Msg_CancelOrder,
						),
					},
				},
			},
			expectedError: errors.New("expected MsgPlaceOrder, got *types.MsgCancelOrder"),
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			var ctx sdk.Context
			_, err := types.ValidateAndTransformRawOperations(
				ctx,
				tc.operations,
				constants.TestEncodingCfg.TxConfig.TxDecoder(),
				constants.EmptyAnteHandler,
			)
			if tc.expectedError != nil {
				require.ErrorContains(t, err, tc.expectedError.Error())
				return
			}
			require.NoError(t, err)
		})
	}
}
