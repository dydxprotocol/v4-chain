package types_test

import (
	"errors"
	fmt "fmt"
	"testing"

	"github.com/dydxprotocol/v4/testutil/constants"
	"github.com/dydxprotocol/v4/x/clob/types"
	satypes "github.com/dydxprotocol/v4/x/subaccounts/types"
	"github.com/stretchr/testify/require"
)

// TestValidateBasic provides black box testing for the ValidateBasic method of the MsgProposedOperations msg.
func TestValidateBasic(t *testing.T) {
	tests := map[string]struct {
		operations    []types.Operation
		expectedError error
	}{
		// Tests for functionality not related to one message type
		"passes validation": {
			operations: []types.Operation{
				types.NewOrderPlacementOperation(constants.Order_Carl_Num0_Id0_Clob0_Buy1BTC_Price50000_GTB10),
				types.NewOrderPlacementOperation(constants.Order_Dave_Num0_Id0_Clob0_Sell1BTC_Price50000_GTB10),
				types.NewMatchOperation(&constants.Order_Carl_Num0_Id0_Clob0_Buy1BTC_Price50000_GTB10,
					[]types.MakerFill{
						{
							FillAmount:   100_000_000, // 1 BTC
							MakerOrderId: constants.Order_Dave_Num0_Id0_Clob0_Sell1BTC_Price50000_GTB10.GetOrderId(),
						},
					},
				),
				types.NewPreexistingStatefulOrderPlacementOperation(
					constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15,
				),
			},
			expectedError: nil,
		},
		"fails when TakerOrderHash invalid length": {
			operations: []types.Operation{
				types.NewOrderPlacementOperation(constants.Order_Carl_Num0_Id0_Clob0_Buy1BTC_Price50000_GTB10),
				types.NewOrderPlacementOperation(constants.Order_Dave_Num0_Id0_Clob0_Sell1BTC_Price50000_GTB10),
				{
					Operation: &types.Operation_Match{
						Match: &types.ClobMatch{
							Match: &types.ClobMatch_MatchOrders{
								MatchOrders: &types.MatchOrders{
									TakerOrderId:   constants.Order_Carl_Num0_Id0_Clob0_Buy1BTC_Price50000_GTB10.GetOrderId(),
									TakerOrderHash: make([]byte, 32+1),
									Fills: []types.MakerFill{
										{
											FillAmount:   100_000_000, // 1 BTC
											MakerOrderId: constants.Order_Dave_Num0_Id0_Clob0_Sell1BTC_Price50000_GTB10.GetOrderId(),
										},
									},
								},
							},
						},
					},
				},
			},
			expectedError: types.ErrInvalidMatchOrder,
		},
		"fails when TakerOrderHash not specified": {
			operations: []types.Operation{
				types.NewOrderPlacementOperation(constants.Order_Carl_Num0_Id0_Clob0_Buy1BTC_Price50000_GTB10),
				types.NewOrderPlacementOperation(constants.Order_Dave_Num0_Id0_Clob0_Sell1BTC_Price50000_GTB10),
				{
					Operation: &types.Operation_Match{
						Match: &types.ClobMatch{
							Match: &types.ClobMatch_MatchOrders{
								MatchOrders: &types.MatchOrders{
									TakerOrderId: constants.Order_Carl_Num0_Id0_Clob0_Buy1BTC_Price50000_GTB10.GetOrderId(),
									Fills: []types.MakerFill{
										{
											FillAmount:   100_000_000, // 1 BTC
											MakerOrderId: constants.Order_Dave_Num0_Id0_Clob0_Sell1BTC_Price50000_GTB10.GetOrderId(),
										},
									},
								},
							},
						},
					},
				},
			},
			expectedError: types.ErrInvalidMatchOrder,
		},

		"a cancelled pre existing order cannot be re-referenced by stateful order id": {
			operations: []types.Operation{
				types.NewPreexistingStatefulOrderPlacementOperation(
					constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15,
				),
				types.NewOrderCancellationOperation(&types.MsgCancelOrder{
					OrderId:      constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15.GetOrderId(),
					GoodTilOneof: &types.MsgCancelOrder_GoodTilBlockTime{GoodTilBlockTime: 20},
				}),
				types.NewMatchOperation(
					&constants.Order_Carl_Num0_Id0_Clob0_Buy1BTC_Price50000_GTB10,
					[]types.MakerFill{
						{
							FillAmount:   5,
							MakerOrderId: constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15.GetOrderId(),
						},
					},
				),
			},
			expectedError: types.ErrOrderPlacementNotInOperationsQueue,
		},
		"a pre existing stateful order can be re-referenced by stateful order id": {
			operations: []types.Operation{
				types.NewPreexistingStatefulOrderPlacementOperation(
					constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15,
				),
				types.NewOrderCancellationOperation(&types.MsgCancelOrder{
					OrderId:      constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15.GetOrderId(),
					GoodTilOneof: &types.MsgCancelOrder_GoodTilBlockTime{GoodTilBlockTime: 20},
				}),
				types.NewMatchOperation(
					&constants.Order_Carl_Num0_Id0_Clob0_Buy1BTC_Price50000_GTB10,
					[]types.MakerFill{
						{
							FillAmount:   5,
							MakerOrderId: constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15.GetOrderId(),
						},
					},
				),
			},
			expectedError: types.ErrOrderPlacementNotInOperationsQueue,
		},

		// tests for invalid subaccount id
		"Stateless order validation: Place Order has invalid SubaccountId": {
			operations: []types.Operation{
				types.NewOrderPlacementOperation(types.Order{
					OrderId:      constants.InvalidSubaccountIdOwner_OrderId,
					Side:         types.Order_SIDE_BUY,
					Quantums:     100_000_000,
					Subticks:     50_000_000_000,
					GoodTilOneof: &types.Order_GoodTilBlock{GoodTilBlock: 10},
				}),
			},
			expectedError: errors.New("invalid SubaccountId Owner address"),
		},
		"Stateless order validation: Cancel Order has invalid SubaccountId": {
			operations: []types.Operation{
				types.NewOrderCancellationOperation(&types.MsgCancelOrder{
					OrderId:      constants.InvalidSubaccountIdNumber_OrderId,
					GoodTilOneof: &types.MsgCancelOrder_GoodTilBlock{GoodTilBlock: 1},
				}),
			},
			expectedError: satypes.ErrInvalidSubaccountIdNumber,
		},
		"Stateless order validation: matchOrders Taker order id invalid SubaccountId": {
			operations: []types.Operation{
				types.NewOrderPlacementOperation(constants.Order_Dave_Num0_Id0_Clob0_Sell1BTC_Price50000),
				{
					Operation: &types.Operation_Match{
						Match: &types.ClobMatch{
							Match: &types.ClobMatch_MatchOrders{
								MatchOrders: &types.MatchOrders{
									TakerOrderHash: constants.OrderHash_Empty[:],
									TakerOrderId:   constants.InvalidSubaccountIdNumber_OrderId,
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
			operations: []types.Operation{
				types.NewOrderPlacementOperation(constants.Order_Dave_Num0_Id0_Clob0_Sell1BTC_Price50000),
				{
					Operation: &types.Operation_Match{
						Match: &types.ClobMatch{
							Match: &types.ClobMatch_MatchOrders{
								MatchOrders: &types.MatchOrders{
									TakerOrderHash: constants.OrderHash_Empty[:],
									TakerOrderId:   constants.OrderId_Alice_Num0_ClientId0_Clob0,
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
			operations: []types.Operation{
				types.NewOrderPlacementOperation(constants.Order_Dave_Num0_Id0_Clob0_Sell1BTC_Price50000),
				{
					Operation: &types.Operation_Match{
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

		// tests for order cancellation
		"cancel order validation: good till block cannot be 0": {
			operations: []types.Operation{
				types.NewOrderCancellationOperation(&types.MsgCancelOrder{
					OrderId:      constants.OrderId_Alice_Num0_ClientId0_Clob0,
					GoodTilOneof: &types.MsgCancelOrder_GoodTilBlock{GoodTilBlock: 0}, // invalid
				}),
			},
			expectedError: errors.New("cancellation goodTilBlock cannot be 0"),
		},
		"cancel order validation: short term order cancellations must reference an order in the same block": {
			operations: []types.Operation{
				types.NewOrderCancellationOperation(&types.MsgCancelOrder{
					OrderId:      constants.OrderId_Alice_Num0_ClientId0_Clob0,
					GoodTilOneof: &types.MsgCancelOrder_GoodTilBlock{GoodTilBlock: 20},
				}),
			},
			expectedError: types.ErrOrderPlacementNotInOperationsQueue,
		},

		// tests for Order Placement
		"Stateless place order validation: replacement order higher priority": {
			operations: []types.Operation{
				types.NewOrderPlacementOperation(constants.Order_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB20), // higher pri
				types.NewOrderPlacementOperation(constants.Order_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB15),
			},
			expectedError: errors.New("Replacement order is not higher priority"),
		},
		"Stateless place order validation: placeOrder has invalid side": {
			operations: []types.Operation{
				types.NewOrderPlacementOperation(types.Order{
					OrderId:      types.OrderId{SubaccountId: constants.Alice_Num0, ClientId: 0, ClobPairId: 0},
					Side:         types.Order_Side(uint32(999)),
					Quantums:     100_000_000,
					Subticks:     50_000_000_000,
					GoodTilOneof: &types.Order_GoodTilBlock{GoodTilBlock: 10},
				}),
				types.NewOrderPlacementOperation(constants.Order_Dave_Num0_Id0_Clob0_Sell1BTC_Price50000_GTB10),
			},
			expectedError: errors.New("invalid order side"),
		},
		"Stateless place order validation: placeOrder has unspecified side": {
			operations: []types.Operation{
				types.NewOrderPlacementOperation(types.Order{
					OrderId:      types.OrderId{SubaccountId: constants.Alice_Num0, ClientId: 0, ClobPairId: 0},
					Quantums:     100_000_000,
					Subticks:     50_000_000_000,
					GoodTilOneof: &types.Order_GoodTilBlock{GoodTilBlock: 10},
				}),
				types.NewOrderPlacementOperation(constants.Order_Dave_Num0_Id0_Clob0_Sell1BTC_Price50000_GTB10),
			},
			expectedError: errors.New("UNSPECIFIED is not a valid order side"),
		},
		"Stateless place order validation: no duplicate order placements": {
			operations: []types.Operation{
				types.NewOrderPlacementOperation(constants.Order_Dave_Num0_Id0_Clob0_Sell1BTC_Price50000_GTB10),
				types.NewOrderPlacementOperation(constants.Order_Dave_Num0_Id0_Clob0_Sell1BTC_Price50000_GTB10),
			},
			expectedError: errors.New("Duplicate Order"),
		},
		"Stateless place order validation: order quantums cannot be 0": {
			operations: []types.Operation{
				types.NewOrderPlacementOperation(types.Order{
					OrderId:      types.OrderId{SubaccountId: constants.Alice_Num0, ClientId: 0, ClobPairId: 0},
					Side:         types.Order_SIDE_BUY,
					Quantums:     0,
					Subticks:     50_000_000_000,
					GoodTilOneof: &types.Order_GoodTilBlock{GoodTilBlock: 10},
				}),
				types.NewOrderPlacementOperation(constants.Order_Dave_Num0_Id0_Clob0_Sell1BTC_Price50000_GTB10),
			},
			expectedError: errors.New("order size quantums cannot be 0"),
		},
		"Stateless place order validation: order goodTilBlock cannot be 0": {
			operations: []types.Operation{
				types.NewOrderPlacementOperation(types.Order{
					OrderId:      types.OrderId{SubaccountId: constants.Alice_Num0, ClientId: 0, ClobPairId: 0},
					Side:         types.Order_SIDE_BUY,
					Quantums:     10,
					Subticks:     50_000_000_000,
					GoodTilOneof: &types.Order_GoodTilBlock{GoodTilBlock: 0},
				}),
				types.NewOrderPlacementOperation(constants.Order_Dave_Num0_Id0_Clob0_Sell1BTC_Price50000_GTB10),
			},
			expectedError: errors.New("order goodTilBlock cannot be 0"),
		},
		"Stateless place order validation: order subticks cannot be 0": {
			operations: []types.Operation{
				types.NewOrderPlacementOperation(types.Order{
					OrderId:      types.OrderId{SubaccountId: constants.Alice_Num0, ClientId: 0, ClobPairId: 0},
					Side:         types.Order_SIDE_BUY,
					Quantums:     10,
					Subticks:     0,
					GoodTilOneof: &types.Order_GoodTilBlock{GoodTilBlock: 10},
				}),
				types.NewOrderPlacementOperation(constants.Order_Dave_Num0_Id0_Clob0_Sell1BTC_Price50000_GTB10),
			},
			expectedError: errors.New("order subticks cannot be 0"),
		},

		// tests for Match Orders
		"Stateless match order validation: fill amount is zero": {
			operations: []types.Operation{
				types.NewOrderPlacementOperation(constants.Order_Dave_Num0_Id0_Clob0_Sell1BTC_Price50000),
				types.NewOrderPlacementOperation(constants.Order_Carl_Num1_Id0_Clob0_Buy1BTC_Price50000),
				types.NewMatchOperation(
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
			operations: []types.Operation{
				types.NewOrderPlacementOperation(constants.Order_Dave_Num0_Id0_Clob0_Sell1BTC_Price50000),
				types.NewOrderPlacementOperation(constants.Order_Carl_Num1_Id0_Clob0_Buy1BTC_Price50000),
				types.NewMatchOperation(
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

		// tests for Pereptual Liquidations
		"Stateless liquidation validation: total size of liquidation order is zero": {
			operations: []types.Operation{
				types.NewOrderPlacementOperation(constants.Order_Dave_Num0_Id0_Clob0_Sell1BTC_Price50000),
				{
					Operation: &types.Operation_Match{
						Match: &types.ClobMatch{
							Match: &types.ClobMatch_MatchPerpetualLiquidation{
								MatchPerpetualLiquidation: &types.MatchPerpetualLiquidation{
									Liquidated:  constants.Carl_Num0,
									ClobPairId:  0,
									PerpetualId: 0,
									TotalSize:   0, // size is zero
									IsBuy:       true,
									Fills: []types.MakerFill{
										{
											MakerOrderId: constants.Order_Dave_Num0_Id0_Clob0_Sell1BTC_Price50000.GetOrderId(),
											FillAmount:   1,
										},
									},
								},
							},
						},
					},
				},
			},
			expectedError: errors.New("Liquidation match total size is zero"),
		},
		"Stateless liquidation validation: fails if total fill amount exceeds order size": {
			operations: []types.Operation{
				types.NewOrderPlacementOperation(types.Order{
					OrderId:      constants.OrderId_Alice_Num0_ClientId0_Clob0,
					Side:         types.Order_SIDE_SELL,
					Quantums:     150_000_000,
					Subticks:     50_000_000_000,
					GoodTilOneof: &types.Order_GoodTilBlock{GoodTilBlock: 11},
				}),
				types.NewMatchOperationFromPerpetualLiquidation(types.MatchPerpetualLiquidation{
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
			operations: []types.Operation{
				types.NewOrderPlacementOperation(constants.Order_Dave_Num0_Id0_Clob0_Sell1BTC_Price50000),
				types.NewMatchOperationFromPerpetualLiquidation(types.MatchPerpetualLiquidation{
					Liquidated:  constants.Carl_Num0,
					ClobPairId:  0,
					PerpetualId: 0,
					TotalSize:   1,
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
		// tests for Match Perpetual Deleveraging
		"Stateless match perpetual deleveraging validation: forwards errors from validate": {
			operations: []types.Operation{
				types.NewMatchOperationFromPerpetualDeleveragingLiquidation(
					types.MatchPerpetualDeleveraging{
						Liquidated:  constants.Alice_Num0,
						PerpetualId: constants.ClobPair_Eth.MustGetPerpetualId(),
						Fills:       []types.MatchPerpetualDeleveraging_Fill{},
					},
				),
			},
			expectedError: types.ErrEmptyDeleveragingFills,
		},

		// Tests for Pre Existing Stateful Validations
		"Pre Existing Stateful Order validation: duplicate Pre Existing Stateful Orders": {
			operations: []types.Operation{
				types.NewPreexistingStatefulOrderPlacementOperation(
					constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15,
				),
				types.NewPreexistingStatefulOrderPlacementOperation(
					constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15,
				),
			},
			expectedError: errors.New("Duplicate Pre Existing Order Operation"),
		},
		"Pre Existing Stateful Order validation: Non stateful order Id": {
			operations: []types.Operation{
				{
					Operation: &types.Operation_PreexistingStatefulOrder{
						PreexistingStatefulOrder: &constants.OrderId_Alice_Num0_ClientId0_Clob0,
					},
				},
			},
			expectedError: fmt.Errorf(
				"Invalid Preexisting Order Operation: OrderId %+v is not stateful.",
				constants.OrderId_Alice_Num0_ClientId0_Clob0,
			),
		},
		"Pre Existing Stateful Order validation: Fails Order Id Validation": {
			operations: []types.Operation{
				{
					Operation: &types.Operation_PreexistingStatefulOrder{
						PreexistingStatefulOrder: &constants.InvalidSubaccountIdOwner_OrderId,
					},
				},
			},
			expectedError: satypes.ErrInvalidSubaccountIdOwner,
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			msg := types.MsgProposedOperations{
				OperationsQueue: tc.operations,
			}
			err := msg.ValidateBasic()
			if tc.expectedError != nil {
				require.ErrorContains(t, err, tc.expectedError.Error())
				return
			}
			require.NoError(t, err)
		})
	}
}

func TestGetAddToOrderbookCollatCheckOrderHashesSet(t *testing.T) {
	tests := map[string]struct {
		addToOrderbookCollatCheckOrderHashes [][]byte

		expected map[types.OrderHash]bool
	}{
		"Empty list": {
			addToOrderbookCollatCheckOrderHashes: [][]byte{},

			expected: map[types.OrderHash]bool{},
		},
		"List with one zero element": {
			addToOrderbookCollatCheckOrderHashes: [][]byte{
				{
					0, 0, 0, 0, 0, 0, 0, 0,
					0, 0, 0, 0, 0, 0, 0, 0,
					0, 0, 0, 0, 0, 0, 0, 0,
				},
			},

			expected: map[types.OrderHash]bool{
				{
					0, 0, 0, 0, 0, 0, 0, 0,
					0, 0, 0, 0, 0, 0, 0, 0,
					0, 0, 0, 0, 0, 0, 0, 0,
				}: true,
			},
		},
		"List with one non-zero element": {
			addToOrderbookCollatCheckOrderHashes: [][]byte{
				{
					8, 0, 3, 0, 0, 0, 0, 0,
					0, 0, 0, 0, 0, 2, 0, 0,
					0, 0, 0, 0, 0, 0, 0, 0,
				},
			},

			expected: map[types.OrderHash]bool{
				{
					8, 0, 3, 0, 0, 0, 0, 0,
					0, 0, 0, 0, 0, 2, 0, 0,
					0, 0, 0, 0, 0, 0, 0, 0,
				}: true,
			},
		},
		"List with multiple non-zero elements": {
			addToOrderbookCollatCheckOrderHashes: [][]byte{
				{
					8, 0, 3, 0, 0, 0, 0, 0,
					0, 0, 0, 0, 0, 2, 0, 0,
					0, 0, 0, 0, 0, 0, 0, 0,
				},
				{
					7, 0, 3, 0, 0, 0, 0, 0,
					0, 0, 0, 0, 0, 2, 0, 0,
					0, 0, 0, 0, 0, 0, 0, 0,
				},
				{
					7, 6, 3, 0, 0, 0, 0, 0,
					0, 0, 0, 0, 0, 2, 0, 0,
					0, 0, 0, 0, 0, 0, 0, 2,
				},
				{
					1, 2, 3, 0, 0, 0, 0, 0,
					0, 0, 0, 0, 0, 2, 0, 0,
					0, 0, 0, 0, 0, 0, 0, 2,
				},
			},

			expected: map[types.OrderHash]bool{
				{
					8, 0, 3, 0, 0, 0, 0, 0,
					0, 0, 0, 0, 0, 2, 0, 0,
					0, 0, 0, 0, 0, 0, 0, 0,
				}: true,
				{
					7, 0, 3, 0, 0, 0, 0, 0,
					0, 0, 0, 0, 0, 2, 0, 0,
					0, 0, 0, 0, 0, 0, 0, 0,
				}: true,
				{
					7, 6, 3, 0, 0, 0, 0, 0,
					0, 0, 0, 0, 0, 2, 0, 0,
					0, 0, 0, 0, 0, 0, 0, 2,
				}: true,
				{
					1, 2, 3, 0, 0, 0, 0, 0,
					0, 0, 0, 0, 0, 2, 0, 0,
					0, 0, 0, 0, 0, 0, 0, 2,
				}: true,
			},
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			msg := types.MsgProposedOperations{
				AddToOrderbookCollatCheckOrderHashes: tc.addToOrderbookCollatCheckOrderHashes,
			}

			require.Equal(t, tc.expected, msg.GetAddToOrderbookCollatCheckOrderHashesSet())
		})
	}
}

func TestGetAddToOrderbookCollatCheckOrderHashesSet_PanicsOnDuplicate(t *testing.T) {
	hashes := [][]byte{
		{
			1, 2, 3, 0, 0, 0, 0, 0,
			0, 0, 0, 0, 0, 2, 0, 0,
			0, 0, 0, 0, 0, 0, 0, 2,
		},
		{
			1, 2, 3, 0, 0, 0, 0, 0,
			0, 0, 0, 0, 0, 2, 0, 0,
			0, 0, 0, 0, 0, 0, 0, 2,
		},
	}

	msg := types.MsgProposedOperations{AddToOrderbookCollatCheckOrderHashes: hashes}

	require.PanicsWithValue(
		t,
		fmt.Sprintf(
			"GetAddToOrderbookCollatCheckOrderHashesSet: duplicate order hash in AddToOrderbookCollatCheckOrderHashes: %+v",
			hashes,
		),
		func() { msg.GetAddToOrderbookCollatCheckOrderHashesSet() },
	)
}

func TestGetSigners(t *testing.T) {
	msg := types.MsgProposedOperations{}
	require.Empty(t, msg.GetSigners())
}
