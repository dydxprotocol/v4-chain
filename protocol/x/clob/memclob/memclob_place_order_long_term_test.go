package memclob

import (
	"testing"

	"github.com/dydxprotocol/v4/testutil/constants"
	testutil_memclob "github.com/dydxprotocol/v4/testutil/memclob"
	sdktest "github.com/dydxprotocol/v4/testutil/sdk"
	"github.com/dydxprotocol/v4/x/clob/types"
	satypes "github.com/dydxprotocol/v4/x/subaccounts/types"
)

func TestPlaceOrder_LongTerm(t *testing.T) {
	ctx, _, _ := sdktest.NewSdkContextWithMultistore()
	ctx = ctx.WithIsCheckTx(true)
	tests := map[string]struct {
		// State.
		placedMatchableOrders  []types.MatchableOrder
		collateralizationCheck map[int]testutil_memclob.CollateralizationCheck

		// Parameters.
		order types.Order

		// Expectations.
		expectedFilledSize            satypes.BaseQuantums
		expectedOrderStatus           types.OrderStatus
		expectedRemainingBids         []OrderWithRemainingSize
		expectedRemainingAsks         []OrderWithRemainingSize
		expectedOperations            []types.Operation
		expectedOperationToNonce      map[types.Operation]types.Nonce
		expectedPendingStatefulOrders []types.Order
		expectedErr                   error
	}{
		"Can place a valid Long-Term buy order on an empty orderbook": {
			placedMatchableOrders: []types.MatchableOrder{},
			collateralizationCheck: map[int]testutil_memclob.CollateralizationCheck{
				0: {
					CollatCheck: map[satypes.SubaccountId][]types.PendingOpenOrder{
						constants.Alice_Num0: {
							{
								RemainingQuantums: constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15.GetBaseQuantums(),
								IsBuy:             constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15.IsBuy(),
								IsTaker:           false,
								Subticks:          constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15.GetOrderSubticks(),
								ClobPairId:        constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15.GetClobPairId(),
							},
						},
					},
					Result: map[satypes.SubaccountId]satypes.UpdateResult{
						constants.Alice_Num0: satypes.Success,
					},
				},
			},

			order: constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15,

			expectedFilledSize:  0,
			expectedOrderStatus: types.Success,
			expectedRemainingBids: []OrderWithRemainingSize{
				{
					Order:         constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15,
					RemainingSize: 5,
				},
			},
			expectedRemainingAsks: []OrderWithRemainingSize{},
			expectedOperations: []types.Operation{
				types.NewOrderPlacementOperation(
					constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15,
				),
			},
			expectedOperationToNonce: map[types.Operation]types.Nonce{
				types.NewOrderPlacementOperation(constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15): 0,
			},
			expectedPendingStatefulOrders: []types.Order{
				constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15,
			},
		},
		`Matches a buy order when it overlaps the orderbook, and with no orders on the other side it places the remaining
		size on the orderbook`: {
			placedMatchableOrders: []types.MatchableOrder{
				&constants.Order_Alice_Num0_Id1_Clob0_Sell10_Price15_GTB15,
			},
			collateralizationCheck: map[int]testutil_memclob.CollateralizationCheck{
				0: {
					CollatCheck: map[satypes.SubaccountId][]types.PendingOpenOrder{
						constants.Alice_Num0: {
							{
								RemainingQuantums: 10,
								IsBuy:             false,
								IsTaker:           false,
								Subticks:          15,
								ClobPairId:        0,
							},
						},
						constants.Bob_Num0: {
							{
								RemainingQuantums: 10,
								IsBuy:             true,
								IsTaker:           true,
								Subticks:          15,
								ClobPairId:        0,
							},
						},
					},
					Result: map[satypes.SubaccountId]satypes.UpdateResult{
						constants.Alice_Num0: satypes.Success,
						constants.Bob_Num0:   satypes.Success,
					},
				},
				1: {
					CollatCheck: map[satypes.SubaccountId][]types.PendingOpenOrder{
						constants.Bob_Num0: {
							{
								RemainingQuantums: 15,
								IsBuy:             true,
								IsTaker:           false,
								Subticks:          30,
								ClobPairId:        0,
							},
						},
					},
					Result: map[satypes.SubaccountId]satypes.UpdateResult{
						constants.Bob_Num0: satypes.Success,
					},
				},
			},

			order: constants.LongTermOrder_Bob_Num0_Id0_Clob0_Buy25_Price30_GTBT10,

			expectedFilledSize:  10,
			expectedOrderStatus: types.Success,
			expectedOperations: []types.Operation{
				types.NewOrderPlacementOperation(
					constants.Order_Alice_Num0_Id1_Clob0_Sell10_Price15_GTB15,
				),
				types.NewOrderPlacementOperation(
					constants.LongTermOrder_Bob_Num0_Id0_Clob0_Buy25_Price30_GTBT10,
				),
				types.NewMatchOperation(
					&constants.LongTermOrder_Bob_Num0_Id0_Clob0_Buy25_Price30_GTBT10,
					[]types.MakerFill{
						{
							MakerOrderId: constants.Order_Alice_Num0_Id1_Clob0_Sell10_Price15_GTB15.OrderId,
							FillAmount:   10,
						},
					},
				),
			},
			expectedOperationToNonce: map[types.Operation]types.Nonce{
				types.NewOrderPlacementOperation(constants.Order_Alice_Num0_Id1_Clob0_Sell10_Price15_GTB15):       0,
				types.NewOrderPlacementOperation(constants.LongTermOrder_Bob_Num0_Id0_Clob0_Buy25_Price30_GTBT10): 1,
				types.NewMatchOperation(
					&constants.LongTermOrder_Bob_Num0_Id0_Clob0_Buy25_Price30_GTBT10,
					[]types.MakerFill{
						{
							MakerOrderId: constants.Order_Alice_Num0_Id1_Clob0_Sell10_Price15_GTB15.OrderId,
							FillAmount:   10,
						},
					},
				): 2,
			},
			expectedRemainingBids: []OrderWithRemainingSize{
				{
					Order:         constants.LongTermOrder_Bob_Num0_Id0_Clob0_Buy25_Price30_GTBT10,
					RemainingSize: 15,
				},
			},
			expectedRemainingAsks: []OrderWithRemainingSize{},
			expectedPendingStatefulOrders: []types.Order{
				constants.LongTermOrder_Bob_Num0_Id0_Clob0_Buy25_Price30_GTBT10,
			},
		},
		`Fully matches a Long-Term sell order with other Long-Term buy orders when it overlaps the
		orderbook`: {
			placedMatchableOrders: []types.MatchableOrder{
				&constants.LongTermOrder_Bob_Num0_Id0_Clob0_Buy25_Price30_GTBT10,
				&constants.LongTermOrder_Bob_Num0_Id1_Clob0_Buy45_Price10_GTBT10,
			},
			collateralizationCheck: map[int]testutil_memclob.CollateralizationCheck{
				0: {
					CollatCheck: map[satypes.SubaccountId][]types.PendingOpenOrder{
						constants.Alice_Num0: {
							{
								RemainingQuantums: 25,
								IsBuy:             false,
								IsTaker:           true,
								Subticks:          30,
								ClobPairId:        0,
							},
						},
						constants.Bob_Num0: {
							{
								RemainingQuantums: 25,
								IsBuy:             true,
								IsTaker:           false,
								Subticks:          30,
								ClobPairId:        0,
							},
						},
					},
					Result: map[satypes.SubaccountId]satypes.UpdateResult{
						constants.Alice_Num0: satypes.Success,
						constants.Bob_Num0:   satypes.Success,
					},
				},
				1: {
					CollatCheck: map[satypes.SubaccountId][]types.PendingOpenOrder{
						constants.Alice_Num0: {
							{
								RemainingQuantums: 40,
								IsBuy:             false,
								IsTaker:           true,
								Subticks:          10,
								ClobPairId:        0,
							},
						},
						constants.Bob_Num0: {
							{
								RemainingQuantums: 40,
								IsBuy:             true,
								IsTaker:           false,
								Subticks:          10,
								ClobPairId:        0,
							},
						},
					},
					Result: map[satypes.SubaccountId]satypes.UpdateResult{
						constants.Alice_Num0: satypes.Success,
						constants.Bob_Num0:   satypes.Success,
					},
				},
			},

			order: constants.LongTermOrder_Alice_Num0_Id2_Clob0_Sell65_Price10_GTBT25,

			expectedFilledSize:  65,
			expectedOrderStatus: types.Success,
			expectedOperations: []types.Operation{
				types.NewOrderPlacementOperation(
					constants.LongTermOrder_Bob_Num0_Id0_Clob0_Buy25_Price30_GTBT10,
				),
				types.NewOrderPlacementOperation(
					constants.LongTermOrder_Bob_Num0_Id1_Clob0_Buy45_Price10_GTBT10,
				),
				types.NewOrderPlacementOperation(
					constants.LongTermOrder_Alice_Num0_Id2_Clob0_Sell65_Price10_GTBT25,
				),
				types.NewMatchOperation(
					&constants.LongTermOrder_Alice_Num0_Id2_Clob0_Sell65_Price10_GTBT25,
					[]types.MakerFill{
						{
							MakerOrderId: constants.LongTermOrder_Bob_Num0_Id0_Clob0_Buy25_Price30_GTBT10.OrderId,
							FillAmount:   25,
						},
						{
							MakerOrderId: constants.LongTermOrder_Bob_Num0_Id1_Clob0_Buy45_Price10_GTBT10.OrderId,
							FillAmount:   40,
						},
					},
				),
			},
			expectedOperationToNonce: map[types.Operation]types.Nonce{
				types.NewOrderPlacementOperation(constants.LongTermOrder_Bob_Num0_Id0_Clob0_Buy25_Price30_GTBT10):    0,
				types.NewOrderPlacementOperation(constants.LongTermOrder_Bob_Num0_Id1_Clob0_Buy45_Price10_GTBT10):    1,
				types.NewOrderPlacementOperation(constants.LongTermOrder_Alice_Num0_Id2_Clob0_Sell65_Price10_GTBT25): 2,
				types.NewMatchOperation(
					&constants.LongTermOrder_Alice_Num0_Id2_Clob0_Sell65_Price10_GTBT25,
					[]types.MakerFill{
						{
							MakerOrderId: constants.LongTermOrder_Bob_Num0_Id0_Clob0_Buy25_Price30_GTBT10.OrderId,
							FillAmount:   25,
						},
						{
							MakerOrderId: constants.LongTermOrder_Bob_Num0_Id1_Clob0_Buy45_Price10_GTBT10.OrderId,
							FillAmount:   40,
						},
					},
				): 3,
			},
			expectedRemainingBids: []OrderWithRemainingSize{
				{
					Order:         constants.LongTermOrder_Bob_Num0_Id1_Clob0_Buy45_Price10_GTBT10,
					RemainingSize: 5,
				},
			},
			expectedRemainingAsks: []OrderWithRemainingSize{},
			expectedPendingStatefulOrders: []types.Order{
				constants.LongTermOrder_Bob_Num0_Id0_Clob0_Buy25_Price30_GTBT10,
				constants.LongTermOrder_Bob_Num0_Id1_Clob0_Buy45_Price10_GTBT10,
				constants.LongTermOrder_Alice_Num0_Id2_Clob0_Sell65_Price10_GTBT25,
			},
		},
		`Short-Term taker order can fully match with Long-Term maker order`: {
			placedMatchableOrders: []types.MatchableOrder{
				&constants.LongTermOrder_Bob_Num0_Id0_Clob0_Buy25_Price30_GTBT10,
				&constants.LongTermOrder_Bob_Num0_Id1_Clob0_Buy45_Price10_GTBT10,
			},
			collateralizationCheck: map[int]testutil_memclob.CollateralizationCheck{
				0: {
					CollatCheck: map[satypes.SubaccountId][]types.PendingOpenOrder{
						constants.Alice_Num0: {
							{
								RemainingQuantums: 5,
								IsBuy:             false,
								IsTaker:           true,
								Subticks:          30,
								ClobPairId:        0,
							},
						},
						constants.Bob_Num0: {
							{
								RemainingQuantums: 5,
								IsBuy:             true,
								IsTaker:           false,
								Subticks:          30,
								ClobPairId:        0,
							},
						},
					},
					Result: map[satypes.SubaccountId]satypes.UpdateResult{
						constants.Alice_Num0: satypes.Success,
						constants.Bob_Num0:   satypes.Success,
					},
				},
			},

			order: constants.Order_Alice_Num0_Id0_Clob0_Sell5_Price10_GTB20,

			expectedFilledSize:  5,
			expectedOrderStatus: types.Success,
			expectedOperations: []types.Operation{
				types.NewOrderPlacementOperation(
					constants.LongTermOrder_Bob_Num0_Id0_Clob0_Buy25_Price30_GTBT10,
				),
				types.NewOrderPlacementOperation(
					constants.LongTermOrder_Bob_Num0_Id1_Clob0_Buy45_Price10_GTBT10,
				),
				types.NewOrderPlacementOperation(
					constants.Order_Alice_Num0_Id0_Clob0_Sell5_Price10_GTB20,
				),
				types.NewMatchOperation(
					&constants.Order_Alice_Num0_Id0_Clob0_Sell5_Price10_GTB20,
					[]types.MakerFill{
						{
							MakerOrderId: constants.LongTermOrder_Bob_Num0_Id0_Clob0_Buy25_Price30_GTBT10.OrderId,
							FillAmount:   5,
						},
					},
				),
			},
			expectedOperationToNonce: map[types.Operation]types.Nonce{
				types.NewOrderPlacementOperation(constants.LongTermOrder_Bob_Num0_Id0_Clob0_Buy25_Price30_GTBT10): 0,
				types.NewOrderPlacementOperation(constants.LongTermOrder_Bob_Num0_Id1_Clob0_Buy45_Price10_GTBT10): 1,
				types.NewOrderPlacementOperation(constants.Order_Alice_Num0_Id0_Clob0_Sell5_Price10_GTB20):        2,
				types.NewMatchOperation(
					&constants.Order_Alice_Num0_Id0_Clob0_Sell5_Price10_GTB20,
					[]types.MakerFill{
						{
							MakerOrderId: constants.LongTermOrder_Bob_Num0_Id0_Clob0_Buy25_Price30_GTBT10.OrderId,
							FillAmount:   5,
						},
					},
				): 3,
			},
			expectedRemainingBids: []OrderWithRemainingSize{
				{
					Order:         constants.LongTermOrder_Bob_Num0_Id0_Clob0_Buy25_Price30_GTBT10,
					RemainingSize: 20,
				},
				{
					Order:         constants.LongTermOrder_Bob_Num0_Id1_Clob0_Buy45_Price10_GTBT10,
					RemainingSize: 45,
				},
			},
			expectedRemainingAsks: []OrderWithRemainingSize{},
			expectedPendingStatefulOrders: []types.Order{
				constants.LongTermOrder_Bob_Num0_Id0_Clob0_Buy25_Price30_GTBT10,
				constants.LongTermOrder_Bob_Num0_Id1_Clob0_Buy45_Price10_GTBT10,
			},
		},
		`A Long-Term sell order can partially match with a Long-Term buy order, fail collateralization
			checks while matching, and all existing matches are considered valid`: {
			placedMatchableOrders: []types.MatchableOrder{
				&constants.LongTermOrder_Bob_Num0_Id0_Clob0_Buy25_Price30_GTBT10,
				&constants.LongTermOrder_Bob_Num0_Id1_Clob0_Buy45_Price10_GTBT10,
			},
			collateralizationCheck: map[int]testutil_memclob.CollateralizationCheck{
				0: {
					CollatCheck: map[satypes.SubaccountId][]types.PendingOpenOrder{
						constants.Alice_Num0: {
							{
								RemainingQuantums: 25,
								IsBuy:             false,
								IsTaker:           true,
								Subticks:          30,
								ClobPairId:        0,
							},
						},
						constants.Bob_Num0: {
							{
								RemainingQuantums: 25,
								IsBuy:             true,
								IsTaker:           false,
								Subticks:          30,
								ClobPairId:        0,
							},
						},
					},
					Result: map[satypes.SubaccountId]satypes.UpdateResult{
						constants.Alice_Num0: satypes.Success,
						constants.Bob_Num0:   satypes.Success,
					},
				},
				1: {
					CollatCheck: map[satypes.SubaccountId][]types.PendingOpenOrder{
						constants.Alice_Num0: {
							{
								RemainingQuantums: 40,
								IsBuy:             false,
								IsTaker:           true,
								Subticks:          10,
								ClobPairId:        0,
							},
						},
						constants.Bob_Num0: {
							{
								RemainingQuantums: 40,
								IsBuy:             true,
								IsTaker:           false,
								Subticks:          10,
								ClobPairId:        0,
							},
						},
					},
					Result: map[satypes.SubaccountId]satypes.UpdateResult{
						constants.Alice_Num0: satypes.StillUndercollateralized,
						constants.Bob_Num0:   satypes.NewlyUndercollateralized,
					},
				},
			},

			order: constants.LongTermOrder_Alice_Num0_Id2_Clob0_Sell65_Price10_GTBT25,

			expectedFilledSize:  25,
			expectedOrderStatus: types.Undercollateralized,
			expectedOperations: []types.Operation{
				types.NewOrderPlacementOperation(
					constants.LongTermOrder_Bob_Num0_Id0_Clob0_Buy25_Price30_GTBT10,
				),
				types.NewOrderPlacementOperation(
					constants.LongTermOrder_Bob_Num0_Id1_Clob0_Buy45_Price10_GTBT10,
				),
				types.NewOrderPlacementOperation(
					constants.LongTermOrder_Alice_Num0_Id2_Clob0_Sell65_Price10_GTBT25,
				),
				types.NewMatchOperation(
					&constants.LongTermOrder_Alice_Num0_Id2_Clob0_Sell65_Price10_GTBT25,
					[]types.MakerFill{
						{
							MakerOrderId: constants.LongTermOrder_Bob_Num0_Id0_Clob0_Buy25_Price30_GTBT10.OrderId,
							FillAmount:   25,
						},
					},
				),
			},
			expectedOperationToNonce: map[types.Operation]types.Nonce{
				types.NewOrderPlacementOperation(constants.LongTermOrder_Bob_Num0_Id0_Clob0_Buy25_Price30_GTBT10):    0,
				types.NewOrderPlacementOperation(constants.LongTermOrder_Bob_Num0_Id1_Clob0_Buy45_Price10_GTBT10):    1,
				types.NewOrderPlacementOperation(constants.LongTermOrder_Alice_Num0_Id2_Clob0_Sell65_Price10_GTBT25): 2,
				types.NewMatchOperation(
					&constants.LongTermOrder_Alice_Num0_Id2_Clob0_Sell65_Price10_GTBT25,
					[]types.MakerFill{
						{
							MakerOrderId: constants.LongTermOrder_Bob_Num0_Id0_Clob0_Buy25_Price30_GTBT10.OrderId,
							FillAmount:   25,
						},
					},
				): 3,
			},
			expectedRemainingBids: []OrderWithRemainingSize{},
			expectedRemainingAsks: []OrderWithRemainingSize{},
			expectedPendingStatefulOrders: []types.Order{
				constants.LongTermOrder_Bob_Num0_Id0_Clob0_Buy25_Price30_GTBT10,
				constants.LongTermOrder_Bob_Num0_Id1_Clob0_Buy45_Price10_GTBT10,
				constants.LongTermOrder_Alice_Num0_Id2_Clob0_Sell65_Price10_GTBT25,
			},
		},
		`A Long-Term sell order can partially match with a Long-Term buy order, fail collateralization
			checks when adding to orderbook, and all existing matches are considered valid`: {
			placedMatchableOrders: []types.MatchableOrder{
				&constants.LongTermOrder_Bob_Num0_Id0_Clob0_Buy25_Price30_GTBT10,
			},
			collateralizationCheck: map[int]testutil_memclob.CollateralizationCheck{
				0: {
					CollatCheck: map[satypes.SubaccountId][]types.PendingOpenOrder{
						constants.Alice_Num0: {
							{
								RemainingQuantums: 25,
								IsBuy:             false,
								IsTaker:           true,
								Subticks:          30,
								ClobPairId:        0,
							},
						},
						constants.Bob_Num0: {
							{
								RemainingQuantums: 25,
								IsBuy:             true,
								IsTaker:           false,
								Subticks:          30,
								ClobPairId:        0,
							},
						},
					},
					Result: map[satypes.SubaccountId]satypes.UpdateResult{
						constants.Alice_Num0: satypes.Success,
						constants.Bob_Num0:   satypes.Success,
					},
				},
				1: {
					CollatCheck: map[satypes.SubaccountId][]types.PendingOpenOrder{
						constants.Alice_Num0: {
							{
								RemainingQuantums: 40,
								IsBuy:             false,
								IsTaker:           false,
								Subticks:          10,
								ClobPairId:        0,
							},
						},
					},
					Result: map[satypes.SubaccountId]satypes.UpdateResult{
						constants.Alice_Num0: satypes.NewlyUndercollateralized,
					},
				},
			},

			order: constants.LongTermOrder_Alice_Num0_Id2_Clob0_Sell65_Price10_GTBT25,

			expectedFilledSize:  25,
			expectedOrderStatus: types.Undercollateralized,
			expectedOperations: []types.Operation{
				types.NewOrderPlacementOperation(
					constants.LongTermOrder_Bob_Num0_Id0_Clob0_Buy25_Price30_GTBT10,
				),
				types.NewOrderPlacementOperation(
					constants.LongTermOrder_Alice_Num0_Id2_Clob0_Sell65_Price10_GTBT25,
				),
				types.NewMatchOperation(
					&constants.LongTermOrder_Alice_Num0_Id2_Clob0_Sell65_Price10_GTBT25,
					[]types.MakerFill{
						{
							MakerOrderId: constants.LongTermOrder_Bob_Num0_Id0_Clob0_Buy25_Price30_GTBT10.OrderId,
							FillAmount:   25,
						},
					},
				),
			},
			expectedOperationToNonce: map[types.Operation]types.Nonce{
				types.NewOrderPlacementOperation(constants.LongTermOrder_Bob_Num0_Id0_Clob0_Buy25_Price30_GTBT10):    0,
				types.NewOrderPlacementOperation(constants.LongTermOrder_Alice_Num0_Id2_Clob0_Sell65_Price10_GTBT25): 1,
				types.NewMatchOperation(
					&constants.LongTermOrder_Alice_Num0_Id2_Clob0_Sell65_Price10_GTBT25,
					[]types.MakerFill{
						{
							MakerOrderId: constants.LongTermOrder_Bob_Num0_Id0_Clob0_Buy25_Price30_GTBT10.OrderId,
							FillAmount:   25,
						},
					},
				): 2,
			},
			expectedRemainingBids: []OrderWithRemainingSize{},
			expectedRemainingAsks: []OrderWithRemainingSize{},
			expectedPendingStatefulOrders: []types.Order{
				constants.LongTermOrder_Bob_Num0_Id0_Clob0_Buy25_Price30_GTBT10,
				constants.LongTermOrder_Alice_Num0_Id2_Clob0_Sell65_Price10_GTBT25,
			},
		},
		`A Long-Term post-only sell order can partially match with a Long-Term buy order,
				all existing matches are reverted and it's not added to pendingStatefulOrders`: {
			placedMatchableOrders: []types.MatchableOrder{
				&constants.LongTermOrder_Bob_Num0_Id0_Clob0_Buy25_Price30_GTBT10,
			},
			collateralizationCheck: map[int]testutil_memclob.CollateralizationCheck{
				0: {
					CollatCheck: map[satypes.SubaccountId][]types.PendingOpenOrder{
						constants.Alice_Num0: {
							{
								RemainingQuantums: 25,
								IsBuy:             false,
								IsTaker:           true,
								Subticks:          30,
								ClobPairId:        0,
							},
						},
						constants.Bob_Num0: {
							{
								RemainingQuantums: 25,
								IsBuy:             true,
								IsTaker:           false,
								Subticks:          30,
								ClobPairId:        0,
							},
						},
					},
					Result: map[satypes.SubaccountId]satypes.UpdateResult{
						constants.Alice_Num0: satypes.Success,
						constants.Bob_Num0:   satypes.Success,
					},
				},
			},

			order: constants.LongTermOrder_Alice_Num0_Id2_Clob0_Sell65_Price10_GTBT25_PO,

			expectedFilledSize:  0,
			expectedOrderStatus: types.Success,
			expectedOperations: []types.Operation{
				types.NewOrderPlacementOperation(
					constants.LongTermOrder_Bob_Num0_Id0_Clob0_Buy25_Price30_GTBT10,
				),
			},
			expectedOperationToNonce: map[types.Operation]types.Nonce{
				types.NewOrderPlacementOperation(constants.LongTermOrder_Bob_Num0_Id0_Clob0_Buy25_Price30_GTBT10): 0,
				// Post-only order not added to OperationsToPropose.
			},
			expectedRemainingBids: []OrderWithRemainingSize{
				{
					Order:         constants.LongTermOrder_Bob_Num0_Id0_Clob0_Buy25_Price30_GTBT10,
					RemainingSize: 25,
				},
			},
			expectedRemainingAsks: []OrderWithRemainingSize{},
			expectedPendingStatefulOrders: []types.Order{
				constants.LongTermOrder_Bob_Num0_Id0_Clob0_Buy25_Price30_GTBT10,
			},
			expectedErr: types.ErrPostOnlyWouldCrossMakerOrder,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// Setup memclob state and test expectations.
			memclob, fakeMemClobKeeper, expectedNumCollateralizationChecks, numCollateralChecks := simplePlaceOrderTestSetup(
				t,
				ctx,
				tc.placedMatchableOrders,
				tc.collateralizationCheck,
				constants.GetStatePosition_ZeroPositionSize,
				&tc.order,
			)

			// Run the test case and verify expectations.
			placeOrderAndVerifyExpectationsOperations(
				t,
				ctx,
				memclob,
				tc.order,
				numCollateralChecks,
				tc.expectedFilledSize,
				tc.expectedFilledSize,
				tc.expectedOrderStatus,
				tc.expectedErr,
				expectedNumCollateralizationChecks,
				tc.expectedRemainingBids,
				tc.expectedRemainingAsks,
				tc.expectedOperations,
				tc.expectedOperationToNonce,
				tc.expectedPendingStatefulOrders,
				fakeMemClobKeeper,
			)

			// TODO(DEC-1296): Verify the correct offchain update messages were returned for Long-Term orders.
		})
	}
}

func TestPlaceOrder_PreexistingStatefulOrder(t *testing.T) {
	// Setup memclob state and test expectations.
	ctx, _, _ := sdktest.NewSdkContextWithMultistore()
	longTermOrder := constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15
	collateralizationCheck := map[int]testutil_memclob.CollateralizationCheck{
		0: {
			CollatCheck: map[satypes.SubaccountId][]types.PendingOpenOrder{
				constants.Alice_Num0: {
					{
						RemainingQuantums: 5,
						IsBuy:             true,
						IsTaker:           false,
						Subticks:          10,
						ClobPairId:        0,
					},
				},
			},
			Result: map[satypes.SubaccountId]satypes.UpdateResult{
				constants.Alice_Num0: satypes.Success,
			},
		},
	}
	memclob, fakeMemClobKeeper, expectedNumCollateralizationChecks, numCollateralChecks := simplePlaceOrderTestSetup(
		t,
		ctx,
		[]types.MatchableOrder{},
		collateralizationCheck,
		constants.GetStatePosition_ZeroPositionSize,
		&longTermOrder,
	)

	fakeMemClobKeeper.SetStatefulOrderPlacement(ctx, longTermOrder, uint32(5))

	// Run the test case and verify expectations.
	placeOrderAndVerifyExpectations(
		t,
		ctx,
		memclob,
		longTermOrder,
		numCollateralChecks,
		0,
		0,
		types.Success,
		nil,
		expectedNumCollateralizationChecks,
		[]OrderWithRemainingSize{
			{
				Order:         longTermOrder,
				RemainingSize: 5,
			},
		},
		[]OrderWithRemainingSize{},
		[]expectedMatch{},
		[]types.Order{}, // Note the Long-Term order should not be added as a pending stateful order.
		fakeMemClobKeeper,
	)
}
