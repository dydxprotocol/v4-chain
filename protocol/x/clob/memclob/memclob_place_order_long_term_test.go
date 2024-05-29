package memclob

import (
	"testing"

	clobtest "github.com/dydxprotocol/v4-chain/protocol/testutil/clob"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	testutil_memclob "github.com/dydxprotocol/v4-chain/protocol/testutil/memclob"
	sdktest "github.com/dydxprotocol/v4-chain/protocol/testutil/sdk"
	"github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
	satypes "github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"
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
		expectedFilledSize         satypes.BaseQuantums
		expectedOrderStatus        types.OrderStatus
		expectedRemainingBids      []OrderWithRemainingSize
		expectedRemainingAsks      []OrderWithRemainingSize
		expectedOperations         []types.Operation
		expectedInternalOperations []types.InternalOperation
		expectedErr                error
	}{
		"Can place a valid Long-Term buy order on an empty orderbook": {
			placedMatchableOrders:  []types.MatchableOrder{},
			collateralizationCheck: map[int]testutil_memclob.CollateralizationCheck{},

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
				clobtest.NewOrderPlacementOperation(
					constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15,
				),
			},
			expectedInternalOperations: []types.InternalOperation{},
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
								Subticks:          15,
								ClobPairId:        0,
							},
						},
						constants.Bob_Num0: {
							{
								RemainingQuantums: 10,
								IsBuy:             true,
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
			},

			order: constants.LongTermOrder_Bob_Num0_Id0_Clob0_Buy25_Price30_GTBT10,

			expectedFilledSize:  10,
			expectedOrderStatus: types.Success,
			expectedOperations: []types.Operation{
				clobtest.NewOrderPlacementOperation(
					constants.Order_Alice_Num0_Id1_Clob0_Sell10_Price15_GTB15,
				),
				clobtest.NewOrderPlacementOperation(
					constants.LongTermOrder_Bob_Num0_Id0_Clob0_Buy25_Price30_GTBT10,
				),
				clobtest.NewMatchOperation(
					&constants.LongTermOrder_Bob_Num0_Id0_Clob0_Buy25_Price30_GTBT10,
					[]types.MakerFill{
						{
							MakerOrderId: constants.Order_Alice_Num0_Id1_Clob0_Sell10_Price15_GTB15.OrderId,
							FillAmount:   10,
						},
					},
				),
			},
			expectedInternalOperations: []types.InternalOperation{
				types.NewShortTermOrderPlacementInternalOperation(
					constants.Order_Alice_Num0_Id1_Clob0_Sell10_Price15_GTB15,
				),
				types.NewPreexistingStatefulOrderPlacementInternalOperation(
					constants.LongTermOrder_Bob_Num0_Id0_Clob0_Buy25_Price30_GTBT10,
				),
				types.NewMatchOrdersInternalOperation(
					constants.LongTermOrder_Bob_Num0_Id0_Clob0_Buy25_Price30_GTBT10,
					[]types.MakerFill{
						{
							MakerOrderId: constants.Order_Alice_Num0_Id1_Clob0_Sell10_Price15_GTB15.OrderId,
							FillAmount:   10,
						},
					},
				),
			},
			expectedRemainingBids: []OrderWithRemainingSize{
				{
					Order:         constants.LongTermOrder_Bob_Num0_Id0_Clob0_Buy25_Price30_GTBT10,
					RemainingSize: 15,
				},
			},
			expectedRemainingAsks: []OrderWithRemainingSize{},
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
								Subticks:          30,
								ClobPairId:        0,
							},
						},
						constants.Bob_Num0: {
							{
								RemainingQuantums: 25,
								IsBuy:             true,
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
								Subticks:          10,
								ClobPairId:        0,
							},
						},
						constants.Bob_Num0: {
							{
								RemainingQuantums: 40,
								IsBuy:             true,
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
				clobtest.NewOrderPlacementOperation(
					constants.LongTermOrder_Bob_Num0_Id0_Clob0_Buy25_Price30_GTBT10,
				),
				clobtest.NewOrderPlacementOperation(
					constants.LongTermOrder_Bob_Num0_Id1_Clob0_Buy45_Price10_GTBT10,
				),
				clobtest.NewOrderPlacementOperation(
					constants.LongTermOrder_Alice_Num0_Id2_Clob0_Sell65_Price10_GTBT25,
				),
				clobtest.NewMatchOperation(
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
			expectedInternalOperations: []types.InternalOperation{
				types.NewPreexistingStatefulOrderPlacementInternalOperation(
					constants.LongTermOrder_Bob_Num0_Id0_Clob0_Buy25_Price30_GTBT10,
				),
				types.NewPreexistingStatefulOrderPlacementInternalOperation(
					constants.LongTermOrder_Bob_Num0_Id1_Clob0_Buy45_Price10_GTBT10,
				),
				types.NewPreexistingStatefulOrderPlacementInternalOperation(
					constants.LongTermOrder_Alice_Num0_Id2_Clob0_Sell65_Price10_GTBT25,
				),
				types.NewMatchOrdersInternalOperation(
					constants.LongTermOrder_Alice_Num0_Id2_Clob0_Sell65_Price10_GTBT25,
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
			expectedRemainingBids: []OrderWithRemainingSize{
				{
					Order:         constants.LongTermOrder_Bob_Num0_Id1_Clob0_Buy45_Price10_GTBT10,
					RemainingSize: 5,
				},
			},
			expectedRemainingAsks: []OrderWithRemainingSize{},
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
								Subticks:          30,
								ClobPairId:        0,
							},
						},
						constants.Bob_Num0: {
							{
								RemainingQuantums: 5,
								IsBuy:             true,
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
				clobtest.NewOrderPlacementOperation(
					constants.LongTermOrder_Bob_Num0_Id0_Clob0_Buy25_Price30_GTBT10,
				),
				clobtest.NewOrderPlacementOperation(
					constants.LongTermOrder_Bob_Num0_Id1_Clob0_Buy45_Price10_GTBT10,
				),
				clobtest.NewOrderPlacementOperation(
					constants.Order_Alice_Num0_Id0_Clob0_Sell5_Price10_GTB20,
				),
				clobtest.NewMatchOperation(
					&constants.Order_Alice_Num0_Id0_Clob0_Sell5_Price10_GTB20,
					[]types.MakerFill{
						{
							MakerOrderId: constants.LongTermOrder_Bob_Num0_Id0_Clob0_Buy25_Price30_GTBT10.OrderId,
							FillAmount:   5,
						},
					},
				),
			},
			expectedInternalOperations: []types.InternalOperation{
				types.NewPreexistingStatefulOrderPlacementInternalOperation(
					constants.LongTermOrder_Bob_Num0_Id0_Clob0_Buy25_Price30_GTBT10,
				),
				types.NewShortTermOrderPlacementInternalOperation(
					constants.Order_Alice_Num0_Id0_Clob0_Sell5_Price10_GTB20,
				),
				types.NewMatchOrdersInternalOperation(
					constants.Order_Alice_Num0_Id0_Clob0_Sell5_Price10_GTB20,
					[]types.MakerFill{
						{
							MakerOrderId: constants.LongTermOrder_Bob_Num0_Id0_Clob0_Buy25_Price30_GTBT10.OrderId,
							FillAmount:   5,
						},
					},
				),
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
								Subticks:          30,
								ClobPairId:        0,
							},
						},
						constants.Bob_Num0: {
							{
								RemainingQuantums: 25,
								IsBuy:             true,
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
								Subticks:          10,
								ClobPairId:        0,
							},
						},
						constants.Bob_Num0: {
							{
								RemainingQuantums: 40,
								IsBuy:             true,
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
				clobtest.NewOrderPlacementOperation(
					constants.LongTermOrder_Bob_Num0_Id0_Clob0_Buy25_Price30_GTBT10,
				),
				clobtest.NewOrderPlacementOperation(
					constants.LongTermOrder_Bob_Num0_Id1_Clob0_Buy45_Price10_GTBT10,
				),
				clobtest.NewOrderPlacementOperation(
					constants.LongTermOrder_Alice_Num0_Id2_Clob0_Sell65_Price10_GTBT25,
				),
				clobtest.NewMatchOperation(
					&constants.LongTermOrder_Alice_Num0_Id2_Clob0_Sell65_Price10_GTBT25,
					[]types.MakerFill{
						{
							MakerOrderId: constants.LongTermOrder_Bob_Num0_Id0_Clob0_Buy25_Price30_GTBT10.OrderId,
							FillAmount:   25,
						},
					},
				),
			},
			expectedInternalOperations: []types.InternalOperation{
				types.NewOrderRemovalInternalOperation(
					constants.LongTermOrder_Bob_Num0_Id1_Clob0_Buy45_Price10_GTBT10.OrderId,
					types.OrderRemoval_REMOVAL_REASON_UNDERCOLLATERALIZED,
				),
				types.NewPreexistingStatefulOrderPlacementInternalOperation(
					constants.LongTermOrder_Bob_Num0_Id0_Clob0_Buy25_Price30_GTBT10,
				),
				types.NewPreexistingStatefulOrderPlacementInternalOperation(
					constants.LongTermOrder_Alice_Num0_Id2_Clob0_Sell65_Price10_GTBT25,
				),
				types.NewMatchOrdersInternalOperation(
					constants.LongTermOrder_Alice_Num0_Id2_Clob0_Sell65_Price10_GTBT25,
					[]types.MakerFill{
						{
							MakerOrderId: constants.LongTermOrder_Bob_Num0_Id0_Clob0_Buy25_Price30_GTBT10.OrderId,
							FillAmount:   25,
						},
					},
				),
				types.NewOrderRemovalInternalOperation(
					constants.LongTermOrder_Alice_Num0_Id2_Clob0_Sell65_Price10_GTBT25.OrderId,
					types.OrderRemoval_REMOVAL_REASON_UNDERCOLLATERALIZED,
				),
			},
			expectedRemainingBids: []OrderWithRemainingSize{},
			expectedRemainingAsks: []OrderWithRemainingSize{},
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
								Subticks:          30,
								ClobPairId:        0,
							},
						},
						constants.Bob_Num0: {
							{
								RemainingQuantums: 25,
								IsBuy:             true,
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
				clobtest.NewOrderPlacementOperation(
					constants.LongTermOrder_Bob_Num0_Id0_Clob0_Buy25_Price30_GTBT10,
				),
			},
			expectedInternalOperations: []types.InternalOperation{
				types.NewOrderRemovalInternalOperation(
					constants.LongTermOrder_Alice_Num0_Id2_Clob0_Sell65_Price10_GTBT25_PO.OrderId,
					types.OrderRemoval_REMOVAL_REASON_POST_ONLY_WOULD_CROSS_MAKER_ORDER,
				),
			},
			expectedRemainingBids: []OrderWithRemainingSize{
				{
					Order:         constants.LongTermOrder_Bob_Num0_Id0_Clob0_Buy25_Price30_GTBT10,
					RemainingSize: 25,
				},
			},
			expectedRemainingAsks: []OrderWithRemainingSize{},
			expectedErr:           types.ErrPostOnlyWouldCrossMakerOrder,
		},
		`A Long-term buy order can self-match against a Long-term sell order from the same subaccount,
			causing the maker order to be removed`: {
			placedMatchableOrders: []types.MatchableOrder{
				&constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15,
			},
			collateralizationCheck: map[int]testutil_memclob.CollateralizationCheck{},

			order: constants.LongTermOrder_Alice_Num0_Id2_Clob0_Sell65_Price10_GTBT25,

			expectedFilledSize:  0,
			expectedOrderStatus: types.Success,
			expectedOperations:  []types.Operation{},
			expectedInternalOperations: []types.InternalOperation{
				types.NewOrderRemovalInternalOperation(
					constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15.OrderId,
					types.OrderRemoval_REMOVAL_REASON_INVALID_SELF_TRADE,
				),
			},
			expectedRemainingBids: []OrderWithRemainingSize{},
			expectedRemainingAsks: []OrderWithRemainingSize{
				{
					Order:         constants.LongTermOrder_Alice_Num0_Id2_Clob0_Sell65_Price10_GTBT25,
					RemainingSize: 65,
				},
			},
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
				tc.expectedInternalOperations,
				fakeMemClobKeeper,
			)

			// TODO(DEC-1296): Verify the correct offchain update messages were returned for Long-Term orders.
		})
	}
}

func TestPlaceOrder_PreexistingStatefulOrder(t *testing.T) {
	// Setup memclob state and test expectations.
	ctx, _, _ := sdktest.NewSdkContextWithMultistore()
	ctx = ctx.WithIsCheckTx(true)
	longTermOrder := constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15
	collateralizationCheck := map[int]testutil_memclob.CollateralizationCheck{}
	memclob, fakeMemClobKeeper, expectedNumCollateralizationChecks, numCollateralChecks := simplePlaceOrderTestSetup(
		t,
		ctx,
		[]types.MatchableOrder{},
		collateralizationCheck,
		constants.GetStatePosition_ZeroPositionSize,
		&longTermOrder,
	)

	fakeMemClobKeeper.SetLongTermOrderPlacement(ctx, longTermOrder, uint32(5))

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
		fakeMemClobKeeper,
	)
}
