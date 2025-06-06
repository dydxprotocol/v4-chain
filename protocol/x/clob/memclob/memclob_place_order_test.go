package memclob

import (
	"encoding/json"
	"math/big"
	"testing"

	"github.com/cosmos/cosmos-sdk/telemetry"
	"github.com/dydxprotocol/v4-chain/protocol/mocks"
	clobtest "github.com/dydxprotocol/v4-chain/protocol/testutil/clob"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	testutil_memclob "github.com/dydxprotocol/v4-chain/protocol/testutil/memclob"
	sdktest "github.com/dydxprotocol/v4-chain/protocol/testutil/sdk"
	"github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
	satypes "github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestPlaceOrder_AddOrderToOrderbook(t *testing.T) {
	ctx, _, _ := sdktest.NewSdkContextWithMultistore()
	ctx = ctx.WithIsCheckTx(true)
	tests := map[string]struct {
		// State.
		existingOrders         []types.MatchableOrder
		canceledOrderGTB       uint32
		collateralizationCheck satypes.UpdateResult

		// Parameters.
		order types.Order

		// Expectations.
		expectedOrderStatus    types.OrderStatus
		expectedErr            error
		expectedToReplaceOrder bool
	}{
		"Can place a valid buy order on an empty orderbook": {
			existingOrders:         []types.MatchableOrder{},
			collateralizationCheck: satypes.Success,

			order: constants.Order_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB15,

			expectedOrderStatus:    types.Success,
			expectedToReplaceOrder: false,
		},
		"Can place a valid sell order on an empty orderbook": {
			existingOrders:         []types.MatchableOrder{},
			collateralizationCheck: satypes.Success,

			order: constants.Order_Alice_Num0_Id1_Clob0_Sell10_Price15_GTB15,

			expectedOrderStatus:    types.Success,
			expectedToReplaceOrder: false,
		},
		"Can place a new buy order on an orderbook with bids, and best bid is updated": {
			existingOrders: []types.MatchableOrder{
				&constants.Order_Bob_Num0_Id3_Clob1_Buy10_Price10_GTB20,
			},
			collateralizationCheck: satypes.Success,

			order: constants.Order_Bob_Num0_Id4_Clob1_Buy20_Price35_GTB22,

			expectedOrderStatus:    types.Success,
			expectedToReplaceOrder: false,
		},
		"Can place a new sell order on an orderbook with asks, and best ask is updated": {
			existingOrders: []types.MatchableOrder{
				&constants.Order_Alice_Num0_Id5_Clob1_Sell25_Price15_GTB20,
			},
			collateralizationCheck: satypes.Success,

			order: constants.Order_Alice_Num0_Id3_Clob1_Sell5_Price10_GTB15,

			expectedOrderStatus:    types.Success,
			expectedToReplaceOrder: false,
		},
		`Can place a new sell order on an orderbook with asks at same price level, and best ask is not updated but total
				level quantums is updated`: {
			existingOrders: []types.MatchableOrder{
				&constants.Order_Alice_Num1_Id1_Clob1_Sell10_Price15_GTB20,
			},
			collateralizationCheck: satypes.Success,

			order: constants.Order_Bob_Num0_Id0_Clob1_Sell10_Price15_GTB20,

			expectedOrderStatus:    types.Success,
			expectedToReplaceOrder: false,
		},
		`Can place a new buy order on an orderbook with bids at same price level, and best bid is not updated but total
					level quantums is updated`: {
			existingOrders: []types.MatchableOrder{
				&constants.Order_Alice_Num1_Id2_Clob1_Buy67_Price5_GTB20,
			},
			collateralizationCheck: satypes.Success,

			order: constants.Order_Alice_Num1_Id3_Clob1_Buy7_Price5,

			expectedOrderStatus:    types.Success,
			expectedToReplaceOrder: false,
		},
		"Can place a new sell order on an orderbook with asks at a better price level, and best ask is not updated": {
			existingOrders: []types.MatchableOrder{
				&constants.Order_Bob_Num0_Id0_Clob1_Sell10_Price15_GTB20,
			},
			collateralizationCheck: satypes.Success,

			order: constants.Order_Bob_Num0_Id1_Clob1_Sell11_Price16_GTB20,

			expectedOrderStatus:    types.Success,
			expectedToReplaceOrder: false,
		},
		"Can place a new buy order on an orderbook with bids at a better price level, and best bid is not updated": {
			existingOrders: []types.MatchableOrder{
				&constants.Order_Bob_Num0_Id3_Clob1_Buy10_Price10_GTB20,
				&constants.Order_Bob_Num0_Id4_Clob1_Buy20_Price35_GTB22,
			},
			collateralizationCheck: satypes.Success,

			order: constants.Order_Alice_Num0_Id4_Clob1_Buy25_Price5_GTB20,

			expectedOrderStatus:    types.Success,
			expectedToReplaceOrder: false,
		},
		"Can place a new buy order on an orderbook with multiple bids and asks at the same price level": {
			existingOrders: []types.MatchableOrder{
				&constants.Order_Alice_Num1_Id2_Clob1_Buy67_Price5_GTB20,
				&constants.Order_Alice_Num1_Id3_Clob1_Buy7_Price5,
				&constants.Order_Bob_Num0_Id0_Clob1_Sell10_Price15_GTB20,
				&constants.Order_Alice_Num0_Id5_Clob1_Sell25_Price15_GTB20,
			},
			collateralizationCheck: satypes.Success,

			order: constants.Order_Alice_Num0_Id4_Clob1_Buy25_Price5_GTB20,

			expectedOrderStatus:    types.Success,
			expectedToReplaceOrder: false,
		},
		"Can place a new sell order on an orderbook with multiple bids and asks at different price levels": {
			existingOrders: []types.MatchableOrder{
				&constants.Order_Alice_Num1_Id2_Clob1_Buy67_Price5_GTB20,
				&constants.Order_Bob_Num0_Id3_Clob1_Buy10_Price10_GTB20,
				&constants.Order_Bob_Num0_Id0_Clob1_Sell10_Price15_GTB20,
				&constants.Order_Bob_Num0_Id1_Clob1_Sell11_Price16_GTB20,
			},
			collateralizationCheck: satypes.Success,

			order: constants.Order_Alice_Num0_Id5_Clob1_Sell25_Price15_GTB20,

			expectedOrderStatus:    types.Success,
			expectedToReplaceOrder: false,
		},
		"Placing a canceled order fails": {
			existingOrders: []types.MatchableOrder{
				&constants.Order_Alice_Num1_Id2_Clob1_Buy67_Price5_GTB20,
				&constants.Order_Bob_Num0_Id3_Clob1_Buy10_Price10_GTB20,
				&constants.Order_Bob_Num0_Id0_Clob1_Sell10_Price15_GTB20,
				&constants.Order_Alice_Num0_Id5_Clob1_Sell25_Price15_GTB20,
			},
			canceledOrderGTB: constants.Order_Bob_Num0_Id1_Clob1_Sell11_Price16_GTB20.GetGoodTilBlock(),

			order: constants.Order_Bob_Num0_Id1_Clob1_Sell11_Price16_GTB20,

			expectedErr:            types.ErrOrderIsCanceled,
			expectedToReplaceOrder: false,
		},
		"Placing a stale canceled order succeeds": {
			existingOrders: []types.MatchableOrder{
				&constants.Order_Alice_Num1_Id2_Clob1_Buy67_Price5_GTB20,
				&constants.Order_Bob_Num0_Id3_Clob1_Buy10_Price10_GTB20,
				&constants.Order_Bob_Num0_Id0_Clob1_Sell10_Price15_GTB20,
				&constants.Order_Alice_Num0_Id5_Clob1_Sell25_Price15_GTB20,
			},
			canceledOrderGTB: constants.Order_Bob_Num0_Id1_Clob1_Sell11_Price16_GTB20.GetGoodTilBlock() - 1,

			order: constants.Order_Bob_Num0_Id1_Clob1_Sell11_Price16_GTB20,

			expectedOrderStatus:    types.Success,
			expectedToReplaceOrder: false,
		},
		"Replacing an order fails if GoodTilBlock is lower than existing order": {
			existingOrders: []types.MatchableOrder{
				&constants.Order_Bob_Num0_Id1_Clob1_Sell11_Price16_GTB20,
			},
			order:                  constants.Order_Bob_Num0_Id1_Clob1_Sell11_Price16_GTB18,
			expectedErr:            types.ErrInvalidReplacement,
			expectedToReplaceOrder: false,
		},
		"Replacing an order fails if the existing order has the same GoodTilBlock and hash": {
			existingOrders: []types.MatchableOrder{
				&constants.Order_Bob_Num0_Id1_Clob1_Sell11_Price16_GTB20,
			},
			order:                  constants.Order_Bob_Num0_Id1_Clob1_Sell11_Price16_GTB20,
			expectedErr:            types.ErrInvalidReplacement,
			expectedToReplaceOrder: false,
		},
		"Replacing an order succeeds if GoodTilBlock is greater than existing order": {
			existingOrders: []types.MatchableOrder{
				&constants.Order_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB15,
			},

			order: constants.Order_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB20,

			collateralizationCheck: satypes.Success,
			expectedOrderStatus:    types.Success,
			expectedToReplaceOrder: true,
		},
		"Replacing an order fails if GoodTilBlock is greater than existing order and changes sides": {
			existingOrders: []types.MatchableOrder{
				&constants.Order_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB15,
			},

			order: constants.Order_Alice_Num0_Id0_Clob0_Sell5_Price10_GTB20,

			expectedErr:            types.ErrInvalidReplacement,
			expectedToReplaceOrder: false,
		},
		"Replacing an order succeeds if GoodTilBlock is greater than existing order and changes price": {
			existingOrders: []types.MatchableOrder{
				&constants.Order_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB15,
			},

			order: constants.Order_Alice_Num0_Id0_Clob0_Buy5_Price5_GTB20,

			collateralizationCheck: satypes.Success,
			expectedOrderStatus:    types.Success,
			expectedToReplaceOrder: true,
		},
		"Replacing an order succeeds if GoodTilBlock is greater than existing order and changes size": {
			existingOrders: []types.MatchableOrder{
				&constants.Order_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB15,
			},

			order: constants.Order_Alice_Num0_Id0_Clob0_Buy6_Price10_GTB20,

			collateralizationCheck: satypes.Success,
			expectedOrderStatus:    types.Success,
			expectedToReplaceOrder: true,
		},
		"Replacing an order fails if OrderHash is less than existing order but GoodTilBlock is the same": {
			existingOrders: []types.MatchableOrder{
				&constants.Order_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB20,
			},

			order: constants.Order_Alice_Num0_Id0_Clob0_Buy6_Price10_GTB20,

			expectedErr:            types.ErrInvalidReplacement,
			expectedToReplaceOrder: false,
		},
		"Replacing an order succeeds if OrderHash is greater than existing order but GoodTilBlock is the same": {
			existingOrders: []types.MatchableOrder{
				&constants.Order_Alice_Num0_Id0_Clob0_Buy6_Price10_GTB20,
			},

			order: constants.Order_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB20,

			collateralizationCheck: satypes.Success,
			expectedOrderStatus:    types.Success,
			expectedToReplaceOrder: true,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// Setup memclob state and test expectations.

			collatCheckFailures := make(map[int]map[satypes.SubaccountId]satypes.UpdateResult)
			// Only include the subaccount within `collatCheckFailures` if it's not successful.
			if tc.collateralizationCheck != satypes.Success {
				collatCheckFailures[0] = map[satypes.SubaccountId]satypes.UpdateResult{
					tc.order.OrderId.SubaccountId: tc.collateralizationCheck,
				}
			}
			addOrderToOrderbookSize := satypes.BaseQuantums(0)
			// If we don't expect an error, then `PlaceOrder` should attempt to place the order on the orderbook.
			if tc.expectedErr == nil {
				addOrderToOrderbookSize = tc.order.GetBaseQuantums()
			}
			memclob, fakeMemClobKeeper, expectedNumCollateralizationChecks, numCollateralChecks := placeOrderTestSetup(
				t,
				ctx,
				tc.existingOrders,
				&tc.order,
				[]expectedMatch{},
				tc.expectedOrderStatus,
				addOrderToOrderbookSize,
				tc.expectedErr,
				collatCheckFailures,
				constants.GetStatePosition_ZeroPositionSize,
			)

			// Mark the current order as canceled if necessary.
			if tc.canceledOrderGTB != 0 {
				orderbook := memclob.mustGetOrderbook(tc.order.GetClobPairId())
				orderbook.addShortTermCancel(tc.order.OrderId, tc.canceledOrderGTB)
			}

			// TODO(DEC-1640): Explicitly specify expected remaining orders on the book in test case.
			expectedRemainingBids := make([]OrderWithRemainingSize, 0)
			expectedRemainingAsks := make([]OrderWithRemainingSize, 0)

			ordersOnBook := tc.existingOrders
			// If we expect the order to have been successfully placed on the book, add it to the existing orders.
			shouldExpectOrderToBePlaced := tc.expectedOrderStatus == types.Success && tc.expectedErr == nil
			if shouldExpectOrderToBePlaced {
				order := tc.order
				ordersOnBook = append(ordersOnBook, &order)
			}

			for _, matchableOrder := range ordersOnBook {
				// Note we assume these are regular orders since liquidation orders cannot rest on
				// the book.

				// If this is an order replacement and it was successful, we assert that the old order being replaced
				// is no longer on the book.
				matchableOrderOrder := matchableOrder.MustGetOrder()
				if tc.expectedToReplaceOrder && matchableOrderOrder.OrderId == tc.order.OrderId &&
					tc.order.MustCmpReplacementOrder(&matchableOrderOrder) > 0 {
					continue
				}

				order := matchableOrder.MustGetOrder()
				if order.IsBuy() {
					expectedRemainingBids = append(expectedRemainingBids, OrderWithRemainingSize{
						Order:         order,
						RemainingSize: order.GetBaseQuantums(),
					})
				} else {
					expectedRemainingAsks = append(expectedRemainingAsks, OrderWithRemainingSize{
						Order:         order,
						RemainingSize: order.GetBaseQuantums(),
					})
				}
			}

			// Run the test case and verify expectations.
			offchainUpdates := placeOrderAndVerifyExpectationsOperations(
				t,
				ctx,
				memclob,
				tc.order,
				numCollateralChecks,
				0, // expectedFilledSize is 0 since no matches are expected.
				0, // expectedTotalFilledSize is 0 since no matches are expected.
				tc.expectedOrderStatus,
				tc.expectedErr,
				expectedNumCollateralizationChecks,
				expectedRemainingBids,
				expectedRemainingAsks,
				[]types.Operation{},         // expectedOperations is empty since no matches are expected.
				[]types.InternalOperation{}, // expectedInternalOperations is empty since no matches are expected.
				fakeMemClobKeeper,
			)

			// Verify the correct offchain update messages were returned.
			assertPlaceOrderOffchainMessages(
				t,
				ctx,
				offchainUpdates,
				tc.order,
				tc.existingOrders,
				collatCheckFailures,
				tc.expectedErr,
				0,
				tc.expectedOrderStatus,
				[]expectedMatch{},
				[]expectedMatch{},
				[]types.OrderId{},
				tc.expectedToReplaceOrder,
			)
		})
	}
}

func TestPlaceOrder_MatchOrders(t *testing.T) {
	ctx, _, _ := sdktest.NewSdkContextWithMultistore()
	ctx = ctx.WithIsCheckTx(true)
	tests := map[string]struct {
		// State.
		placedMatchableOrders []types.MatchableOrder

		// Parameters.
		order types.Order

		// Expectations.
		expectedFilledSize         satypes.BaseQuantums
		expectedOrderStatus        types.OrderStatus
		expectedCollatCheck        []expectedMatch
		expectedRemainingBids      []OrderWithRemainingSize
		expectedRemainingAsks      []OrderWithRemainingSize
		expectedMatches            []expectedMatch
		expectedOperations         []types.Operation
		expectedInternalOperations []types.InternalOperation
	}{
		`Matches a buy order when it overlaps the orderbook, and with no orders on the other side it places the remaining
					size on the orderbook`: {
			placedMatchableOrders: []types.MatchableOrder{
				&constants.Order_Alice_Num0_Id2_Clob1_Sell5_Price10_GTB15,
			},

			order: constants.Order_Bob_Num0_Id4_Clob1_Buy20_Price35_GTB22,

			expectedFilledSize:  5,
			expectedOrderStatus: types.Success,
			expectedCollatCheck: []expectedMatch{
				{
					makerOrder:      &constants.Order_Alice_Num0_Id2_Clob1_Sell5_Price10_GTB15,
					takerOrder:      &constants.Order_Bob_Num0_Id4_Clob1_Buy20_Price35_GTB22,
					matchedQuantums: 5,
				},
			},
			expectedRemainingBids: []OrderWithRemainingSize{
				{
					Order:         constants.Order_Bob_Num0_Id4_Clob1_Buy20_Price35_GTB22,
					RemainingSize: 15,
				},
			},
			expectedRemainingAsks: []OrderWithRemainingSize{},
			expectedMatches: []expectedMatch{
				{
					makerOrder:      &constants.Order_Alice_Num0_Id2_Clob1_Sell5_Price10_GTB15,
					takerOrder:      &constants.Order_Bob_Num0_Id4_Clob1_Buy20_Price35_GTB22,
					matchedQuantums: 5,
				},
			},
			expectedOperations: []types.Operation{
				clobtest.NewOrderPlacementOperation(
					constants.Order_Alice_Num0_Id2_Clob1_Sell5_Price10_GTB15,
				),
				clobtest.NewOrderPlacementOperation(
					constants.Order_Bob_Num0_Id4_Clob1_Buy20_Price35_GTB22,
				),
				clobtest.NewMatchOperation(
					&constants.Order_Bob_Num0_Id4_Clob1_Buy20_Price35_GTB22,
					[]types.MakerFill{
						{
							MakerOrderId: constants.Order_Alice_Num0_Id2_Clob1_Sell5_Price10_GTB15.OrderId,
							FillAmount:   5,
						},
					},
				),
			},
			expectedInternalOperations: []types.InternalOperation{
				types.NewShortTermOrderPlacementInternalOperation(
					constants.Order_Alice_Num0_Id2_Clob1_Sell5_Price10_GTB15,
				),
				types.NewShortTermOrderPlacementInternalOperation(
					constants.Order_Bob_Num0_Id4_Clob1_Buy20_Price35_GTB22,
				),
				types.NewMatchOrdersInternalOperation(
					constants.Order_Bob_Num0_Id4_Clob1_Buy20_Price35_GTB22,
					[]types.MakerFill{
						{
							MakerOrderId: constants.Order_Alice_Num0_Id2_Clob1_Sell5_Price10_GTB15.OrderId,
							FillAmount:   5,
						},
					},
				),
			},
		},
		`Matches a sell order when it overlaps the orderbook, and with no orders on the other side it places the remaining
					size on the orderbook`: {
			placedMatchableOrders: []types.MatchableOrder{
				&constants.Order_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB15,
			},

			order: constants.Order_Bob_Num0_Id8_Clob0_Sell20_Price10_GTB22,

			expectedFilledSize:    5,
			expectedOrderStatus:   types.Success,
			expectedRemainingBids: []OrderWithRemainingSize{},
			expectedRemainingAsks: []OrderWithRemainingSize{
				{
					Order:         constants.Order_Bob_Num0_Id8_Clob0_Sell20_Price10_GTB22,
					RemainingSize: 15,
				},
			},
			expectedCollatCheck: []expectedMatch{
				{
					makerOrder:      &constants.Order_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB15,
					takerOrder:      &constants.Order_Bob_Num0_Id8_Clob0_Sell20_Price10_GTB22,
					matchedQuantums: 5,
				},
			},
			expectedMatches: []expectedMatch{
				{
					makerOrder:      &constants.Order_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB15,
					takerOrder:      &constants.Order_Bob_Num0_Id8_Clob0_Sell20_Price10_GTB22,
					matchedQuantums: 5,
				},
			},
			expectedOperations: []types.Operation{
				clobtest.NewOrderPlacementOperation(
					constants.Order_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB15,
				),
				clobtest.NewOrderPlacementOperation(
					constants.Order_Bob_Num0_Id8_Clob0_Sell20_Price10_GTB22,
				),
				clobtest.NewMatchOperation(
					&constants.Order_Bob_Num0_Id8_Clob0_Sell20_Price10_GTB22,
					[]types.MakerFill{
						{
							MakerOrderId: constants.Order_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB15.OrderId,
							FillAmount:   5,
						},
					},
				),
			},
			expectedInternalOperations: []types.InternalOperation{
				types.NewShortTermOrderPlacementInternalOperation(
					constants.Order_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB15,
				),
				types.NewShortTermOrderPlacementInternalOperation(
					constants.Order_Bob_Num0_Id8_Clob0_Sell20_Price10_GTB22,
				),
				types.NewMatchOrdersInternalOperation(
					constants.Order_Bob_Num0_Id8_Clob0_Sell20_Price10_GTB22,
					[]types.MakerFill{
						{
							MakerOrderId: constants.Order_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB15.OrderId,
							FillAmount:   5,
						},
					},
				),
			},
		},
		`Matches a buy order at max subticks when it overlaps the orderbook, and with no orders on the other side it places
					the remaining size on the orderbook`: {
			placedMatchableOrders: []types.MatchableOrder{
				&constants.Order_Alice_Num0_Id8_Clob1_Sell25_PriceMax_GTB20,
			},

			order: constants.Order_Alice_Num1_Id7_Clob1_Buy35_PriceMax_GTB30,

			expectedFilledSize:  25,
			expectedOrderStatus: types.Success,
			expectedRemainingBids: []OrderWithRemainingSize{
				{
					Order:         constants.Order_Alice_Num1_Id7_Clob1_Buy35_PriceMax_GTB30,
					RemainingSize: 10,
				},
			},
			expectedRemainingAsks: []OrderWithRemainingSize{},
			expectedCollatCheck: []expectedMatch{
				{
					makerOrder:      &constants.Order_Alice_Num0_Id8_Clob1_Sell25_PriceMax_GTB20,
					takerOrder:      &constants.Order_Alice_Num1_Id7_Clob1_Buy35_PriceMax_GTB30,
					matchedQuantums: 25,
				},
			},
			expectedMatches: []expectedMatch{
				{
					makerOrder:      &constants.Order_Alice_Num0_Id8_Clob1_Sell25_PriceMax_GTB20,
					takerOrder:      &constants.Order_Alice_Num1_Id7_Clob1_Buy35_PriceMax_GTB30,
					matchedQuantums: 25,
				},
			},
			expectedOperations: []types.Operation{
				clobtest.NewOrderPlacementOperation(
					constants.Order_Alice_Num0_Id8_Clob1_Sell25_PriceMax_GTB20,
				),
				clobtest.NewOrderPlacementOperation(
					constants.Order_Alice_Num1_Id7_Clob1_Buy35_PriceMax_GTB30,
				),
				clobtest.NewMatchOperation(
					&constants.Order_Alice_Num1_Id7_Clob1_Buy35_PriceMax_GTB30,
					[]types.MakerFill{
						{
							MakerOrderId: constants.Order_Alice_Num0_Id8_Clob1_Sell25_PriceMax_GTB20.OrderId,
							FillAmount:   25,
						},
					},
				),
			},
			expectedInternalOperations: []types.InternalOperation{
				types.NewShortTermOrderPlacementInternalOperation(
					constants.Order_Alice_Num0_Id8_Clob1_Sell25_PriceMax_GTB20,
				),
				types.NewShortTermOrderPlacementInternalOperation(
					constants.Order_Alice_Num1_Id7_Clob1_Buy35_PriceMax_GTB30,
				),
				types.NewMatchOrdersInternalOperation(
					constants.Order_Alice_Num1_Id7_Clob1_Buy35_PriceMax_GTB30,
					[]types.MakerFill{
						{
							MakerOrderId: constants.Order_Alice_Num0_Id8_Clob1_Sell25_PriceMax_GTB20.OrderId,
							FillAmount:   25,
						},
					},
				),
			},
		},
		`Matches a sell order when it overlaps the orderbook, and consumes multiple buy orders on the other side
					from the same subaccount until the sell order is fully matched`: {
			placedMatchableOrders: []types.MatchableOrder{
				&constants.Order_Bob_Num0_Id4_Clob1_Buy20_Price35_GTB22,
				&constants.Order_Bob_Num0_Id11_Clob1_Buy5_Price40_GTB32,
			},

			order: constants.Order_Alice_Num1_Id1_Clob1_Sell10_Price15_GTB20,

			expectedFilledSize:  10,
			expectedOrderStatus: types.Success,
			expectedRemainingBids: []OrderWithRemainingSize{
				{
					Order:         constants.Order_Bob_Num0_Id4_Clob1_Buy20_Price35_GTB22,
					RemainingSize: 15,
				},
			},
			expectedRemainingAsks: []OrderWithRemainingSize{},
			expectedCollatCheck: []expectedMatch{
				{
					makerOrder:      &constants.Order_Bob_Num0_Id11_Clob1_Buy5_Price40_GTB32,
					takerOrder:      &constants.Order_Alice_Num1_Id1_Clob1_Sell10_Price15_GTB20,
					matchedQuantums: 5,
				},
				{
					makerOrder:      &constants.Order_Bob_Num0_Id4_Clob1_Buy20_Price35_GTB22,
					takerOrder:      &constants.Order_Alice_Num1_Id1_Clob1_Sell10_Price15_GTB20,
					matchedQuantums: 5,
				},
			},
			expectedMatches: []expectedMatch{
				{
					makerOrder:      &constants.Order_Bob_Num0_Id11_Clob1_Buy5_Price40_GTB32,
					takerOrder:      &constants.Order_Alice_Num1_Id1_Clob1_Sell10_Price15_GTB20,
					matchedQuantums: 5,
				},
				{
					makerOrder:      &constants.Order_Bob_Num0_Id4_Clob1_Buy20_Price35_GTB22,
					takerOrder:      &constants.Order_Alice_Num1_Id1_Clob1_Sell10_Price15_GTB20,
					matchedQuantums: 5,
				},
			},
			expectedOperations: []types.Operation{
				clobtest.NewOrderPlacementOperation(
					constants.Order_Bob_Num0_Id4_Clob1_Buy20_Price35_GTB22,
				),
				clobtest.NewOrderPlacementOperation(
					constants.Order_Bob_Num0_Id11_Clob1_Buy5_Price40_GTB32,
				),
				clobtest.NewOrderPlacementOperation(
					constants.Order_Alice_Num1_Id1_Clob1_Sell10_Price15_GTB20,
				),
				clobtest.NewMatchOperation(
					&constants.Order_Alice_Num1_Id1_Clob1_Sell10_Price15_GTB20,
					[]types.MakerFill{
						{
							MakerOrderId: constants.Order_Bob_Num0_Id11_Clob1_Buy5_Price40_GTB32.OrderId,
							FillAmount:   5,
						},
						{
							MakerOrderId: constants.Order_Bob_Num0_Id4_Clob1_Buy20_Price35_GTB22.OrderId,
							FillAmount:   5,
						},
					},
				),
			},
			expectedInternalOperations: []types.InternalOperation{
				types.NewShortTermOrderPlacementInternalOperation(
					constants.Order_Bob_Num0_Id11_Clob1_Buy5_Price40_GTB32,
				),
				types.NewShortTermOrderPlacementInternalOperation(
					constants.Order_Bob_Num0_Id4_Clob1_Buy20_Price35_GTB22,
				),
				types.NewShortTermOrderPlacementInternalOperation(
					constants.Order_Alice_Num1_Id1_Clob1_Sell10_Price15_GTB20,
				),
				types.NewMatchOrdersInternalOperation(
					constants.Order_Alice_Num1_Id1_Clob1_Sell10_Price15_GTB20,
					[]types.MakerFill{
						{
							MakerOrderId: constants.Order_Bob_Num0_Id11_Clob1_Buy5_Price40_GTB32.OrderId,
							FillAmount:   5,
						},
						{
							MakerOrderId: constants.Order_Bob_Num0_Id4_Clob1_Buy20_Price35_GTB22.OrderId,
							FillAmount:   5,
						},
					},
				),
			},
		},
		"Buy order is fully matched by sell order, and orderbook is empty": {
			placedMatchableOrders: []types.MatchableOrder{
				&constants.Order_Alice_Num0_Id1_Clob0_Sell10_Price15_GTB15,
			},

			order: constants.Order_Alice_Num1_Clob0_Id4_Buy10_Price45_GTB20,

			expectedFilledSize:    10,
			expectedOrderStatus:   types.Success,
			expectedRemainingBids: []OrderWithRemainingSize{},
			expectedRemainingAsks: []OrderWithRemainingSize{},
			expectedCollatCheck: []expectedMatch{
				{
					makerOrder:      &constants.Order_Alice_Num0_Id1_Clob0_Sell10_Price15_GTB15,
					takerOrder:      &constants.Order_Alice_Num1_Clob0_Id4_Buy10_Price45_GTB20,
					matchedQuantums: 10,
				},
			},
			expectedMatches: []expectedMatch{
				{
					makerOrder:      &constants.Order_Alice_Num0_Id1_Clob0_Sell10_Price15_GTB15,
					takerOrder:      &constants.Order_Alice_Num1_Clob0_Id4_Buy10_Price45_GTB20,
					matchedQuantums: 10,
				},
			},
			expectedOperations: []types.Operation{
				clobtest.NewOrderPlacementOperation(
					constants.Order_Alice_Num0_Id1_Clob0_Sell10_Price15_GTB15,
				),
				clobtest.NewOrderPlacementOperation(
					constants.Order_Alice_Num1_Clob0_Id4_Buy10_Price45_GTB20,
				),
				clobtest.NewMatchOperation(
					&constants.Order_Alice_Num1_Clob0_Id4_Buy10_Price45_GTB20,
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
				types.NewShortTermOrderPlacementInternalOperation(
					constants.Order_Alice_Num1_Clob0_Id4_Buy10_Price45_GTB20,
				),
				types.NewMatchOrdersInternalOperation(
					constants.Order_Alice_Num1_Clob0_Id4_Buy10_Price45_GTB20,
					[]types.MakerFill{
						{
							MakerOrderId: constants.Order_Alice_Num0_Id1_Clob0_Sell10_Price15_GTB15.OrderId,
							FillAmount:   10,
						},
					},
				),
			},
		},
		`Continues matching if two orders from the same subaccount overlap, and cancels any maker orders
			that would cause a self-trade`: {
			placedMatchableOrders: []types.MatchableOrder{
				&constants.Order_Bob_Num0_Id0_Clob1_Sell10_Price15_GTB20,
				&constants.Order_Alice_Num1_Id1_Clob1_Sell10_Price15_GTB20,
			},

			order: constants.Order_Alice_Num1_Id7_Clob1_Buy35_PriceMax_GTB30,

			expectedFilledSize:  10,
			expectedOrderStatus: types.Success,
			expectedRemainingBids: []OrderWithRemainingSize{
				{
					Order:         constants.Order_Alice_Num1_Id7_Clob1_Buy35_PriceMax_GTB30,
					RemainingSize: 25,
				},
			},
			expectedRemainingAsks: []OrderWithRemainingSize{},
			expectedCollatCheck: []expectedMatch{
				{
					makerOrder:      &constants.Order_Bob_Num0_Id0_Clob1_Sell10_Price15_GTB20,
					takerOrder:      &constants.Order_Alice_Num1_Id7_Clob1_Buy35_PriceMax_GTB30,
					matchedQuantums: 10,
				},
			},
			expectedMatches: []expectedMatch{
				{
					makerOrder:      &constants.Order_Bob_Num0_Id0_Clob1_Sell10_Price15_GTB20,
					takerOrder:      &constants.Order_Alice_Num1_Id7_Clob1_Buy35_PriceMax_GTB30,
					matchedQuantums: 10,
				},
			},
			expectedOperations: []types.Operation{
				clobtest.NewOrderPlacementOperation(
					constants.Order_Bob_Num0_Id0_Clob1_Sell10_Price15_GTB20,
				),
				clobtest.NewOrderPlacementOperation(
					constants.Order_Alice_Num1_Id7_Clob1_Buy35_PriceMax_GTB30,
				),
				clobtest.NewMatchOperation(
					&constants.Order_Alice_Num1_Id7_Clob1_Buy35_PriceMax_GTB30,
					[]types.MakerFill{
						{
							MakerOrderId: constants.Order_Bob_Num0_Id0_Clob1_Sell10_Price15_GTB20.OrderId,
							FillAmount:   10,
						},
					},
				),
			},
			expectedInternalOperations: []types.InternalOperation{
				types.NewShortTermOrderPlacementInternalOperation(
					constants.Order_Bob_Num0_Id0_Clob1_Sell10_Price15_GTB20,
				),
				types.NewShortTermOrderPlacementInternalOperation(
					constants.Order_Alice_Num1_Id7_Clob1_Buy35_PriceMax_GTB30,
				),
				types.NewMatchOrdersInternalOperation(
					constants.Order_Alice_Num1_Id7_Clob1_Buy35_PriceMax_GTB30,
					[]types.MakerFill{
						{
							MakerOrderId: constants.Order_Bob_Num0_Id0_Clob1_Sell10_Price15_GTB20.OrderId,
							FillAmount:   10,
						},
					},
				),
			},
		},
		"Buy order fully matches multiple sell orders and remaining size is added to the orderbook after it uncrosses": {
			placedMatchableOrders: []types.MatchableOrder{
				&constants.Order_Alice_Num1_Id5_Clob1_Sell50_Price40_GTB20,
				&constants.Order_Alice_Num0_Id3_Clob1_Sell5_Price10_GTB15,
				&constants.Order_Alice_Num0_Id2_Clob1_Sell5_Price10_GTB15,
			},

			order: constants.Order_Bob_Num0_Id4_Clob1_Buy20_Price35_GTB22,

			expectedFilledSize:  10,
			expectedOrderStatus: types.Success,
			expectedRemainingBids: []OrderWithRemainingSize{
				{
					Order:         constants.Order_Bob_Num0_Id4_Clob1_Buy20_Price35_GTB22,
					RemainingSize: 10,
				},
			},
			expectedRemainingAsks: []OrderWithRemainingSize{
				{
					Order:         constants.Order_Alice_Num1_Id5_Clob1_Sell50_Price40_GTB20,
					RemainingSize: 50,
				},
			},
			expectedCollatCheck: []expectedMatch{
				{
					makerOrder:      &constants.Order_Alice_Num0_Id3_Clob1_Sell5_Price10_GTB15,
					takerOrder:      &constants.Order_Bob_Num0_Id4_Clob1_Buy20_Price35_GTB22,
					matchedQuantums: 5,
				},
				{
					makerOrder:      &constants.Order_Alice_Num0_Id2_Clob1_Sell5_Price10_GTB15,
					takerOrder:      &constants.Order_Bob_Num0_Id4_Clob1_Buy20_Price35_GTB22,
					matchedQuantums: 5,
				},
			},
			expectedMatches: []expectedMatch{
				{
					makerOrder:      &constants.Order_Alice_Num0_Id3_Clob1_Sell5_Price10_GTB15,
					takerOrder:      &constants.Order_Bob_Num0_Id4_Clob1_Buy20_Price35_GTB22,
					matchedQuantums: 5,
				},
				{
					makerOrder:      &constants.Order_Alice_Num0_Id2_Clob1_Sell5_Price10_GTB15,
					takerOrder:      &constants.Order_Bob_Num0_Id4_Clob1_Buy20_Price35_GTB22,
					matchedQuantums: 5,
				},
			},
			expectedOperations: []types.Operation{
				clobtest.NewOrderPlacementOperation(
					constants.Order_Alice_Num0_Id3_Clob1_Sell5_Price10_GTB15,
				),
				clobtest.NewOrderPlacementOperation(
					constants.Order_Alice_Num0_Id2_Clob1_Sell5_Price10_GTB15,
				),
				clobtest.NewOrderPlacementOperation(
					constants.Order_Bob_Num0_Id4_Clob1_Buy20_Price35_GTB22,
				),
				clobtest.NewMatchOperation(
					&constants.Order_Bob_Num0_Id4_Clob1_Buy20_Price35_GTB22,
					[]types.MakerFill{
						{
							MakerOrderId: constants.Order_Alice_Num0_Id3_Clob1_Sell5_Price10_GTB15.OrderId,
							FillAmount:   5,
						},
						{
							MakerOrderId: constants.Order_Alice_Num0_Id2_Clob1_Sell5_Price10_GTB15.OrderId,
							FillAmount:   5,
						},
					},
				),
			},
			expectedInternalOperations: []types.InternalOperation{
				types.NewShortTermOrderPlacementInternalOperation(
					constants.Order_Alice_Num0_Id3_Clob1_Sell5_Price10_GTB15,
				),
				types.NewShortTermOrderPlacementInternalOperation(
					constants.Order_Alice_Num0_Id2_Clob1_Sell5_Price10_GTB15,
				),
				types.NewShortTermOrderPlacementInternalOperation(
					constants.Order_Bob_Num0_Id4_Clob1_Buy20_Price35_GTB22,
				),
				types.NewMatchOrdersInternalOperation(
					constants.Order_Bob_Num0_Id4_Clob1_Buy20_Price35_GTB22,
					[]types.MakerFill{
						{
							MakerOrderId: constants.Order_Alice_Num0_Id3_Clob1_Sell5_Price10_GTB15.OrderId,
							FillAmount:   5,
						},
						{
							MakerOrderId: constants.Order_Alice_Num0_Id2_Clob1_Sell5_Price10_GTB15.OrderId,
							FillAmount:   5,
						},
					},
				),
			},
		},
		"Buy order matches multiple sell orders, before partially matching a sell order": {
			placedMatchableOrders: []types.MatchableOrder{
				&constants.Order_Alice_Num1_Id6_Clob1_Sell15_Price22_GTB30,
				&constants.Order_Alice_Num0_Id3_Clob1_Sell5_Price10_GTB15,
				&constants.Order_Alice_Num0_Id2_Clob1_Sell5_Price10_GTB15,
				&constants.Order_Alice_Num1_Id2_Clob1_Buy67_Price5_GTB20,
				&constants.Order_Alice_Num0_Id4_Clob1_Buy25_Price5_GTB20,
				&constants.Order_Alice_Num1_Id5_Clob1_Sell50_Price40_GTB20,
			},

			order: constants.Order_Bob_Num0_Id4_Clob1_Buy20_Price35_GTB22,

			expectedFilledSize:  20,
			expectedOrderStatus: types.Success,
			expectedRemainingBids: []OrderWithRemainingSize{
				{
					Order:         constants.Order_Alice_Num1_Id2_Clob1_Buy67_Price5_GTB20,
					RemainingSize: 67,
				},
				{
					Order:         constants.Order_Alice_Num0_Id4_Clob1_Buy25_Price5_GTB20,
					RemainingSize: 25,
				},
			},
			expectedRemainingAsks: []OrderWithRemainingSize{
				{
					Order:         constants.Order_Alice_Num1_Id6_Clob1_Sell15_Price22_GTB30,
					RemainingSize: 5,
				},
				{
					Order:         constants.Order_Alice_Num1_Id5_Clob1_Sell50_Price40_GTB20,
					RemainingSize: 50,
				},
			},
			expectedCollatCheck: []expectedMatch{
				{
					makerOrder:      &constants.Order_Alice_Num0_Id3_Clob1_Sell5_Price10_GTB15,
					takerOrder:      &constants.Order_Bob_Num0_Id4_Clob1_Buy20_Price35_GTB22,
					matchedQuantums: 5,
				},
				{
					makerOrder:      &constants.Order_Alice_Num0_Id2_Clob1_Sell5_Price10_GTB15,
					takerOrder:      &constants.Order_Bob_Num0_Id4_Clob1_Buy20_Price35_GTB22,
					matchedQuantums: 5,
				},
				{
					makerOrder:      &constants.Order_Alice_Num1_Id6_Clob1_Sell15_Price22_GTB30,
					takerOrder:      &constants.Order_Bob_Num0_Id4_Clob1_Buy20_Price35_GTB22,
					matchedQuantums: 10,
				},
			},
			expectedMatches: []expectedMatch{
				{
					makerOrder:      &constants.Order_Alice_Num0_Id3_Clob1_Sell5_Price10_GTB15,
					takerOrder:      &constants.Order_Bob_Num0_Id4_Clob1_Buy20_Price35_GTB22,
					matchedQuantums: 5,
				},
				{
					makerOrder:      &constants.Order_Alice_Num0_Id2_Clob1_Sell5_Price10_GTB15,
					takerOrder:      &constants.Order_Bob_Num0_Id4_Clob1_Buy20_Price35_GTB22,
					matchedQuantums: 5,
				},
				{
					makerOrder:      &constants.Order_Alice_Num1_Id6_Clob1_Sell15_Price22_GTB30,
					takerOrder:      &constants.Order_Bob_Num0_Id4_Clob1_Buy20_Price35_GTB22,
					matchedQuantums: 10,
				},
			},
			expectedOperations: []types.Operation{
				clobtest.NewOrderPlacementOperation(
					constants.Order_Alice_Num1_Id6_Clob1_Sell15_Price22_GTB30,
				),
				clobtest.NewOrderPlacementOperation(
					constants.Order_Alice_Num0_Id3_Clob1_Sell5_Price10_GTB15,
				),
				clobtest.NewOrderPlacementOperation(
					constants.Order_Alice_Num0_Id2_Clob1_Sell5_Price10_GTB15,
				),
				clobtest.NewOrderPlacementOperation(
					constants.Order_Bob_Num0_Id4_Clob1_Buy20_Price35_GTB22,
				),
				clobtest.NewMatchOperation(
					&constants.Order_Bob_Num0_Id4_Clob1_Buy20_Price35_GTB22,
					[]types.MakerFill{
						{
							MakerOrderId: constants.Order_Alice_Num0_Id3_Clob1_Sell5_Price10_GTB15.OrderId,
							FillAmount:   5,
						},
						{
							MakerOrderId: constants.Order_Alice_Num0_Id2_Clob1_Sell5_Price10_GTB15.OrderId,
							FillAmount:   5,
						},
						{
							MakerOrderId: constants.Order_Alice_Num1_Id6_Clob1_Sell15_Price22_GTB30.OrderId,
							FillAmount:   10,
						},
					},
				),
			},
			expectedInternalOperations: []types.InternalOperation{
				types.NewShortTermOrderPlacementInternalOperation(
					constants.Order_Alice_Num0_Id3_Clob1_Sell5_Price10_GTB15,
				),
				types.NewShortTermOrderPlacementInternalOperation(
					constants.Order_Alice_Num0_Id2_Clob1_Sell5_Price10_GTB15,
				),
				types.NewShortTermOrderPlacementInternalOperation(
					constants.Order_Alice_Num1_Id6_Clob1_Sell15_Price22_GTB30,
				),
				types.NewShortTermOrderPlacementInternalOperation(
					constants.Order_Bob_Num0_Id4_Clob1_Buy20_Price35_GTB22,
				),
				types.NewMatchOrdersInternalOperation(
					constants.Order_Bob_Num0_Id4_Clob1_Buy20_Price35_GTB22,
					[]types.MakerFill{
						{
							MakerOrderId: constants.Order_Alice_Num0_Id3_Clob1_Sell5_Price10_GTB15.OrderId,
							FillAmount:   5,
						},
						{
							MakerOrderId: constants.Order_Alice_Num0_Id2_Clob1_Sell5_Price10_GTB15.OrderId,
							FillAmount:   5,
						},
						{
							MakerOrderId: constants.Order_Alice_Num1_Id6_Clob1_Sell15_Price22_GTB30.OrderId,
							FillAmount:   10,
						},
					},
				),
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// Setup memclob state and test expectations.
			addOrderToOrderbookSize := satypes.BaseQuantums(0)
			if tc.expectedOrderStatus.IsSuccess() {
				addOrderToOrderbookSize = tc.order.GetBaseQuantums() - tc.expectedFilledSize
			}
			memclob, fakeMemClobKeeper, expectedNumCollateralizationChecks, numCollateralChecks := placeOrderTestSetup(
				t,
				ctx,
				tc.placedMatchableOrders,
				&tc.order,
				tc.expectedCollatCheck,
				tc.expectedOrderStatus,
				addOrderToOrderbookSize,
				nil,
				map[int]map[satypes.SubaccountId]satypes.UpdateResult{},
				constants.GetStatePosition_ZeroPositionSize,
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
				nil,
				expectedNumCollateralizationChecks,
				tc.expectedRemainingBids,
				tc.expectedRemainingAsks,
				tc.expectedOperations,
				tc.expectedInternalOperations,
				fakeMemClobKeeper,
			)

			// Verify the correct offchain update messages were returned.
			// TODO(DEC-1588): Update the indexer tests to properly handle self-trades.
			// assertPlaceOrderOffchainMessages(
			// 	t,
			// 	offchainUpdates,
			// 	tc.order,
			// 	tc.placedMatchableOrders,
			// 	map[int]map[satypes.SubaccountId]satypes.UpdateResult{},
			// 	nil,
			// 	tc.expectedFilledSize,
			// 	tc.expectedOrderStatus,
			// 	[]expectedMatch{},
			// 	tc.expectedMatches,
			// 	[]types.OrderId{},
			// )
		})
	}
}

// TestPlaceOrder_MatchOrders_PreexistingMatches is different from TestPlaceOrder_MatchOrders because there
// exist matches in the match queue before the test case is ran.
func TestPlaceOrder_MatchOrders_PreexistingMatches(t *testing.T) {
	ctx, _, _ := sdktest.NewSdkContextWithMultistore()
	ctx = ctx.WithIsCheckTx(true)
	tests := map[string]struct {
		// State.
		placedMatchableOrders []types.MatchableOrder

		// Parameters.
		order types.Order

		// Expectations.
		expectedFilledSize         satypes.BaseQuantums
		expectedTotalFilledSize    satypes.BaseQuantums
		expectedOrderStatus        types.OrderStatus
		expectedCollatCheck        []expectedMatch
		expectedRemainingBids      []OrderWithRemainingSize
		expectedRemainingAsks      []OrderWithRemainingSize
		expectedExistingMatches    []expectedMatch
		expectedNewMatches         []expectedMatch
		expectedMatches            []expectedMatch
		expectedOperations         []types.Operation
		expectedInternalOperations []types.InternalOperation
		expectedErr                error
		expectedToReplaceOrder     bool
	}{
		"A partially matched sell order is fully matched by a buy order, and the buy order is also fully matched": {
			placedMatchableOrders: []types.MatchableOrder{
				// Match #1: This order is partially matched before the test case as a maker order with the below order.
				&constants.Order_Bob_Num0_Id8_Clob0_Sell20_Price10_GTB22,
				// Match #1: This order is fully matched before the test case as a taker order with the above order.
				&constants.Order_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB15,
				&constants.Order_Alice_Num1_Id9_Clob0_Sell10_Price10_GTB31,
			},

			order: constants.Order_Alice_Num1_Id8_Clob0_Buy15_Price25_GTB31,

			expectedFilledSize:      15,
			expectedTotalFilledSize: 15,
			expectedOrderStatus:     types.Success,
			expectedRemainingBids:   []OrderWithRemainingSize{},
			expectedRemainingAsks: []OrderWithRemainingSize{
				{
					Order:         constants.Order_Alice_Num1_Id9_Clob0_Sell10_Price10_GTB31,
					RemainingSize: 10,
				},
			},
			expectedCollatCheck: []expectedMatch{
				{
					makerOrder:      &constants.Order_Bob_Num0_Id8_Clob0_Sell20_Price10_GTB22,
					takerOrder:      &constants.Order_Alice_Num1_Id8_Clob0_Buy15_Price25_GTB31,
					matchedQuantums: 15,
				},
			},
			expectedExistingMatches: []expectedMatch{
				// Match #1: This match is generated before the test case.
				{
					makerOrder:      &constants.Order_Bob_Num0_Id8_Clob0_Sell20_Price10_GTB22,
					takerOrder:      &constants.Order_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB15,
					matchedQuantums: 5,
				},
			},
			expectedNewMatches: []expectedMatch{
				{
					makerOrder:      &constants.Order_Bob_Num0_Id8_Clob0_Sell20_Price10_GTB22,
					takerOrder:      &constants.Order_Alice_Num1_Id8_Clob0_Buy15_Price25_GTB31,
					matchedQuantums: 15,
				},
			},
			expectedOperations: []types.Operation{
				clobtest.NewOrderPlacementOperation(
					constants.Order_Bob_Num0_Id8_Clob0_Sell20_Price10_GTB22,
				),
				clobtest.NewOrderPlacementOperation(
					constants.Order_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB15,
				),
				clobtest.NewMatchOperation(
					&constants.Order_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB15,
					[]types.MakerFill{
						{
							MakerOrderId: constants.Order_Bob_Num0_Id8_Clob0_Sell20_Price10_GTB22.OrderId,
							FillAmount:   5,
						},
					},
				),
				clobtest.NewOrderPlacementOperation(
					constants.Order_Alice_Num1_Id8_Clob0_Buy15_Price25_GTB31,
				),
				clobtest.NewMatchOperation(
					&constants.Order_Alice_Num1_Id8_Clob0_Buy15_Price25_GTB31,
					[]types.MakerFill{
						{
							MakerOrderId: constants.Order_Bob_Num0_Id8_Clob0_Sell20_Price10_GTB22.OrderId,
							FillAmount:   15,
						},
					},
				),
			},
			expectedInternalOperations: []types.InternalOperation{
				types.NewShortTermOrderPlacementInternalOperation(
					constants.Order_Bob_Num0_Id8_Clob0_Sell20_Price10_GTB22,
				),
				types.NewShortTermOrderPlacementInternalOperation(
					constants.Order_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB15,
				),
				types.NewMatchOrdersInternalOperation(
					constants.Order_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB15,
					[]types.MakerFill{
						{
							MakerOrderId: constants.Order_Bob_Num0_Id8_Clob0_Sell20_Price10_GTB22.OrderId,
							FillAmount:   5,
						},
					},
				),
				types.NewShortTermOrderPlacementInternalOperation(
					constants.Order_Alice_Num1_Id8_Clob0_Buy15_Price25_GTB31,
				),
				types.NewMatchOrdersInternalOperation(
					constants.Order_Alice_Num1_Id8_Clob0_Buy15_Price25_GTB31,
					[]types.MakerFill{
						{
							MakerOrderId: constants.Order_Bob_Num0_Id8_Clob0_Sell20_Price10_GTB22.OrderId,
							FillAmount:   15,
						},
					},
				),
			},
		},
		`Taker has multiple previous matches in the match queue and can submit another taker order and collateralization
		checks include the current matches in the match queue in all collateralization checks`: {
			placedMatchableOrders: []types.MatchableOrder{
				// Match #1: This order is fully matched before the test case as a maker order with the below order.
				&constants.Order_Alice_Num1_Id11_Clob1_Buy10_Price45_GTB20,
				// Match #1: This order is fully matched before the test case as a taker order with the above order.
				&constants.Order_Bob_Num0_Id0_Clob1_Sell10_Price15_GTB20,
				// Match #2: This order is fully matched before the test case as a maker order with the below order.
				&constants.Order_Bob_Num0_Id11_Clob1_Buy5_Price40_GTB32,
				// Match #2: This order is fully matched before the test case as a taker order with the above order.
				&constants.Order_Alice_Num1_Id6_Clob1_Sell15_Price22_GTB30,
			},

			order: constants.Order_Bob_Num0_Id4_Clob1_Buy20_Price35_GTB22,

			expectedFilledSize:      10,
			expectedTotalFilledSize: 10,
			expectedOrderStatus:     types.Success,
			expectedRemainingBids: []OrderWithRemainingSize{
				{
					Order:         constants.Order_Bob_Num0_Id4_Clob1_Buy20_Price35_GTB22,
					RemainingSize: 10,
				},
			},
			expectedRemainingAsks: []OrderWithRemainingSize{},
			expectedCollatCheck: []expectedMatch{
				{
					makerOrder:      &constants.Order_Alice_Num1_Id6_Clob1_Sell15_Price22_GTB30,
					takerOrder:      &constants.Order_Bob_Num0_Id4_Clob1_Buy20_Price35_GTB22,
					matchedQuantums: 10,
				},
			},
			expectedExistingMatches: []expectedMatch{
				// Match #1: This match is generated before the test case.
				{
					makerOrder:      &constants.Order_Alice_Num1_Id11_Clob1_Buy10_Price45_GTB20,
					takerOrder:      &constants.Order_Bob_Num0_Id0_Clob1_Sell10_Price15_GTB20,
					matchedQuantums: 10,
				},
				// Match #2: This match is generated before the test case.
				{
					makerOrder:      &constants.Order_Bob_Num0_Id11_Clob1_Buy5_Price40_GTB32,
					takerOrder:      &constants.Order_Alice_Num1_Id6_Clob1_Sell15_Price22_GTB30,
					matchedQuantums: 5,
				},
			},
			expectedNewMatches: []expectedMatch{
				{
					makerOrder:      &constants.Order_Alice_Num1_Id6_Clob1_Sell15_Price22_GTB30,
					takerOrder:      &constants.Order_Bob_Num0_Id4_Clob1_Buy20_Price35_GTB22,
					matchedQuantums: 10,
				},
			},
			expectedOperations: []types.Operation{
				clobtest.NewOrderPlacementOperation(
					constants.Order_Alice_Num1_Id11_Clob1_Buy10_Price45_GTB20,
				),
				clobtest.NewOrderPlacementOperation(
					constants.Order_Bob_Num0_Id0_Clob1_Sell10_Price15_GTB20,
				),
				clobtest.NewMatchOperation(
					&constants.Order_Bob_Num0_Id0_Clob1_Sell10_Price15_GTB20,
					[]types.MakerFill{
						{
							MakerOrderId: constants.Order_Alice_Num1_Id11_Clob1_Buy10_Price45_GTB20.OrderId,
							FillAmount:   10,
						},
					},
				),
				clobtest.NewOrderPlacementOperation(
					constants.Order_Bob_Num0_Id11_Clob1_Buy5_Price40_GTB32,
				),
				clobtest.NewOrderPlacementOperation(
					constants.Order_Alice_Num1_Id6_Clob1_Sell15_Price22_GTB30,
				),
				clobtest.NewMatchOperation(
					&constants.Order_Alice_Num1_Id6_Clob1_Sell15_Price22_GTB30,
					[]types.MakerFill{
						{
							MakerOrderId: constants.Order_Bob_Num0_Id11_Clob1_Buy5_Price40_GTB32.OrderId,
							FillAmount:   5,
						},
					},
				),
				clobtest.NewOrderPlacementOperation(
					constants.Order_Bob_Num0_Id4_Clob1_Buy20_Price35_GTB22,
				),
				clobtest.NewMatchOperation(
					&constants.Order_Bob_Num0_Id4_Clob1_Buy20_Price35_GTB22,
					[]types.MakerFill{
						{
							MakerOrderId: constants.Order_Alice_Num1_Id6_Clob1_Sell15_Price22_GTB30.OrderId,
							FillAmount:   10,
						},
					},
				),
			},
			expectedInternalOperations: []types.InternalOperation{
				types.NewShortTermOrderPlacementInternalOperation(
					constants.Order_Alice_Num1_Id11_Clob1_Buy10_Price45_GTB20,
				),
				types.NewShortTermOrderPlacementInternalOperation(
					constants.Order_Bob_Num0_Id0_Clob1_Sell10_Price15_GTB20,
				),
				types.NewMatchOrdersInternalOperation(
					constants.Order_Bob_Num0_Id0_Clob1_Sell10_Price15_GTB20,
					[]types.MakerFill{
						{
							MakerOrderId: constants.Order_Alice_Num1_Id11_Clob1_Buy10_Price45_GTB20.OrderId,
							FillAmount:   10,
						},
					},
				),
				types.NewShortTermOrderPlacementInternalOperation(
					constants.Order_Bob_Num0_Id11_Clob1_Buy5_Price40_GTB32,
				),
				types.NewShortTermOrderPlacementInternalOperation(
					constants.Order_Alice_Num1_Id6_Clob1_Sell15_Price22_GTB30,
				),
				types.NewMatchOrdersInternalOperation(
					constants.Order_Alice_Num1_Id6_Clob1_Sell15_Price22_GTB30,
					[]types.MakerFill{
						{
							MakerOrderId: constants.Order_Bob_Num0_Id11_Clob1_Buy5_Price40_GTB32.OrderId,
							FillAmount:   5,
						},
					},
				),
				types.NewShortTermOrderPlacementInternalOperation(
					constants.Order_Bob_Num0_Id4_Clob1_Buy20_Price35_GTB22,
				),
				types.NewMatchOrdersInternalOperation(
					constants.Order_Bob_Num0_Id4_Clob1_Buy20_Price35_GTB22,
					[]types.MakerFill{
						{
							MakerOrderId: constants.Order_Alice_Num1_Id6_Clob1_Sell15_Price22_GTB30.OrderId,
							FillAmount:   10,
						},
					},
				),
			},
		},
		`Taker has no previous matches in the match queue and is fully matched, but maker has previous matches`: {
			placedMatchableOrders: []types.MatchableOrder{
				// Match #1: This order is fully matched before the test case as a maker order with the below order.
				&constants.Order_Alice_Num1_Id11_Clob1_Buy10_Price45_GTB20,
				// Match #1: This order is fully matched before the test case as a taker order with the above order.
				&constants.Order_Bob_Num0_Id0_Clob1_Sell10_Price15_GTB20,
				&constants.Order_Bob_Num0_Id11_Clob1_Buy5_Price40_GTB32,
			},

			order: constants.Order_Alice_Num0_Id2_Clob1_Sell5_Price10_GTB15,

			expectedFilledSize:      5,
			expectedTotalFilledSize: 5,
			expectedOrderStatus:     types.Success,
			expectedRemainingBids:   []OrderWithRemainingSize{},
			expectedRemainingAsks:   []OrderWithRemainingSize{},
			expectedCollatCheck: []expectedMatch{
				{
					makerOrder:      &constants.Order_Bob_Num0_Id11_Clob1_Buy5_Price40_GTB32,
					takerOrder:      &constants.Order_Alice_Num0_Id2_Clob1_Sell5_Price10_GTB15,
					matchedQuantums: 5,
				},
			},
			expectedExistingMatches: []expectedMatch{
				// Match #1: This match is generated before the test case.
				{
					makerOrder:      &constants.Order_Alice_Num1_Id11_Clob1_Buy10_Price45_GTB20,
					takerOrder:      &constants.Order_Bob_Num0_Id0_Clob1_Sell10_Price15_GTB20,
					matchedQuantums: 10,
				},
			},
			expectedNewMatches: []expectedMatch{
				{
					makerOrder:      &constants.Order_Bob_Num0_Id11_Clob1_Buy5_Price40_GTB32,
					takerOrder:      &constants.Order_Alice_Num0_Id2_Clob1_Sell5_Price10_GTB15,
					matchedQuantums: 5,
				},
			},
			expectedOperations: []types.Operation{
				clobtest.NewOrderPlacementOperation(
					constants.Order_Alice_Num1_Id11_Clob1_Buy10_Price45_GTB20,
				),
				clobtest.NewOrderPlacementOperation(
					constants.Order_Bob_Num0_Id0_Clob1_Sell10_Price15_GTB20,
				),
				clobtest.NewMatchOperation(
					&constants.Order_Bob_Num0_Id0_Clob1_Sell10_Price15_GTB20,
					[]types.MakerFill{
						{
							MakerOrderId: constants.Order_Alice_Num1_Id11_Clob1_Buy10_Price45_GTB20.OrderId,
							FillAmount:   10,
						},
					},
				),
				clobtest.NewOrderPlacementOperation(
					constants.Order_Bob_Num0_Id11_Clob1_Buy5_Price40_GTB32,
				),
				clobtest.NewOrderPlacementOperation(
					constants.Order_Alice_Num0_Id2_Clob1_Sell5_Price10_GTB15,
				),
				clobtest.NewMatchOperation(
					&constants.Order_Alice_Num0_Id2_Clob1_Sell5_Price10_GTB15,
					[]types.MakerFill{
						{
							MakerOrderId: constants.Order_Bob_Num0_Id11_Clob1_Buy5_Price40_GTB32.OrderId,
							FillAmount:   5,
						},
					},
				),
			},
			expectedInternalOperations: []types.InternalOperation{
				types.NewShortTermOrderPlacementInternalOperation(
					constants.Order_Alice_Num1_Id11_Clob1_Buy10_Price45_GTB20,
				),
				types.NewShortTermOrderPlacementInternalOperation(
					constants.Order_Bob_Num0_Id0_Clob1_Sell10_Price15_GTB20,
				),
				types.NewMatchOrdersInternalOperation(
					constants.Order_Bob_Num0_Id0_Clob1_Sell10_Price15_GTB20,
					[]types.MakerFill{
						{
							MakerOrderId: constants.Order_Alice_Num1_Id11_Clob1_Buy10_Price45_GTB20.OrderId,
							FillAmount:   10,
						},
					},
				),
				types.NewShortTermOrderPlacementInternalOperation(
					constants.Order_Bob_Num0_Id11_Clob1_Buy5_Price40_GTB32,
				),
				types.NewShortTermOrderPlacementInternalOperation(
					constants.Order_Alice_Num0_Id2_Clob1_Sell5_Price10_GTB15,
				),
				types.NewMatchOrdersInternalOperation(
					constants.Order_Alice_Num0_Id2_Clob1_Sell5_Price10_GTB15,
					[]types.MakerFill{
						{
							MakerOrderId: constants.Order_Bob_Num0_Id11_Clob1_Buy5_Price40_GTB32.OrderId,
							FillAmount:   5,
						},
					},
				),
			},
		},
		`Taker replaces fully matched order which has a match in the match queue with a larger order that
			should fully match`: {
			placedMatchableOrders: []types.MatchableOrder{
				// Match #1: This order is fully matched before the test case as a maker order with the below order.
				&constants.Order_Alice_Num1_Id10_Clob0_Buy5_Price30_GTB32,
				// Match #1: This order is partially matched before the test case as a taker order with the above order.
				&constants.Order_Alice_Num0_Id1_Clob0_Sell10_Price15_GTB15,
			},

			// This order replaces the existing fully-matched order with an order that increases the size.
			// This order should fully match and no order should remain on the book.
			order: constants.Order_Alice_Num1_Id10_Clob0_Buy10_Price30_GTB33,

			expectedFilledSize:      5,
			expectedTotalFilledSize: 10,
			expectedOrderStatus:     types.Success,
			expectedRemainingBids:   []OrderWithRemainingSize{},
			expectedRemainingAsks:   []OrderWithRemainingSize{},
			expectedCollatCheck: []expectedMatch{
				{
					makerOrder:      &constants.Order_Alice_Num0_Id1_Clob0_Sell10_Price15_GTB15,
					takerOrder:      &constants.Order_Alice_Num1_Id10_Clob0_Buy10_Price30_GTB33,
					matchedQuantums: 5,
				},
			},
			expectedExistingMatches: []expectedMatch{
				// Match #1: This match is generated before the test case.
				{
					makerOrder:      &constants.Order_Alice_Num1_Id10_Clob0_Buy5_Price30_GTB32,
					takerOrder:      &constants.Order_Alice_Num0_Id1_Clob0_Sell10_Price15_GTB15,
					matchedQuantums: 5,
				},
			},
			expectedNewMatches: []expectedMatch{
				// Match #2: This match is generated based on the `order` in the test case.
				{
					makerOrder:      &constants.Order_Alice_Num0_Id1_Clob0_Sell10_Price15_GTB15,
					takerOrder:      &constants.Order_Alice_Num1_Id10_Clob0_Buy10_Price30_GTB33,
					matchedQuantums: 5,
				},
			},
			expectedOperations: []types.Operation{
				clobtest.NewOrderPlacementOperation(
					constants.Order_Alice_Num1_Id10_Clob0_Buy5_Price30_GTB32,
				),
				clobtest.NewOrderPlacementOperation(
					constants.Order_Alice_Num0_Id1_Clob0_Sell10_Price15_GTB15,
				),
				clobtest.NewMatchOperation(
					&constants.Order_Alice_Num0_Id1_Clob0_Sell10_Price15_GTB15,
					[]types.MakerFill{
						{
							MakerOrderId: constants.Order_Alice_Num1_Id10_Clob0_Buy5_Price30_GTB32.OrderId,
							FillAmount:   5,
						},
					},
				),
				clobtest.NewOrderPlacementOperation(
					constants.Order_Alice_Num1_Id10_Clob0_Buy10_Price30_GTB33,
				),
				clobtest.NewMatchOperation(
					&constants.Order_Alice_Num1_Id10_Clob0_Buy10_Price30_GTB33,
					[]types.MakerFill{
						{
							MakerOrderId: constants.Order_Alice_Num0_Id1_Clob0_Sell10_Price15_GTB15.OrderId,
							FillAmount:   5,
						},
					},
				),
			},
			expectedInternalOperations: []types.InternalOperation{
				types.NewShortTermOrderPlacementInternalOperation(
					constants.Order_Alice_Num1_Id10_Clob0_Buy5_Price30_GTB32,
				),
				types.NewShortTermOrderPlacementInternalOperation(
					constants.Order_Alice_Num0_Id1_Clob0_Sell10_Price15_GTB15,
				),
				types.NewMatchOrdersInternalOperation(
					constants.Order_Alice_Num0_Id1_Clob0_Sell10_Price15_GTB15,
					[]types.MakerFill{
						{
							MakerOrderId: constants.Order_Alice_Num1_Id10_Clob0_Buy5_Price30_GTB32.OrderId,
							FillAmount:   5,
						},
					},
				),
				types.NewShortTermOrderPlacementInternalOperation(
					constants.Order_Alice_Num1_Id10_Clob0_Buy10_Price30_GTB33,
				),
				types.NewMatchOrdersInternalOperation(
					constants.Order_Alice_Num1_Id10_Clob0_Buy10_Price30_GTB33,
					[]types.MakerFill{
						{
							MakerOrderId: constants.Order_Alice_Num0_Id1_Clob0_Sell10_Price15_GTB15.OrderId,
							FillAmount:   5,
						},
					},
				),
			},
		},
		`Taker replaces fully filled matched order which has matches in the match queue with a larger order that
		is added to the book.`: {
			placedMatchableOrders: []types.MatchableOrder{
				// Match #1: This order is fully matched before the test case as a maker order with the below order.
				&constants.Order_Alice_Num1_Id10_Clob0_Buy5_Price30_GTB32,
				// Match #1: This order is partially matched before the test case as a taker order with the above order.
				&constants.Order_Alice_Num0_Id1_Clob0_Sell10_Price15_GTB15,
			},

			// This order replaces the existing fully-matched order with an order that increases the size.
			// This order should fully fill the crossing order and have its remaining size placed on the book.
			order: constants.Order_Alice_Num1_Id10_Clob0_Buy15_Price30_GTB33,

			expectedFilledSize:      5,
			expectedTotalFilledSize: 10,
			expectedOrderStatus:     types.Success,
			expectedRemainingBids: []OrderWithRemainingSize{
				{
					Order:         constants.Order_Alice_Num1_Id10_Clob0_Buy15_Price30_GTB33,
					RemainingSize: 5,
				},
			},
			expectedRemainingAsks: []OrderWithRemainingSize{},
			expectedCollatCheck: []expectedMatch{
				{
					makerOrder:      &constants.Order_Alice_Num0_Id1_Clob0_Sell10_Price15_GTB15,
					takerOrder:      &constants.Order_Alice_Num1_Id10_Clob0_Buy15_Price30_GTB33,
					matchedQuantums: 5,
				},
			},
			expectedExistingMatches: []expectedMatch{
				// Match #1: This match is generated before the test case.
				{
					makerOrder:      &constants.Order_Alice_Num1_Id10_Clob0_Buy5_Price30_GTB32,
					takerOrder:      &constants.Order_Alice_Num0_Id1_Clob0_Sell10_Price15_GTB15,
					matchedQuantums: 5,
				},
			},
			expectedNewMatches: []expectedMatch{
				// Match #2: This match is generated based on the `order` in the test case.
				{
					makerOrder:      &constants.Order_Alice_Num0_Id1_Clob0_Sell10_Price15_GTB15,
					takerOrder:      &constants.Order_Alice_Num1_Id10_Clob0_Buy15_Price30_GTB33,
					matchedQuantums: 5,
				},
			},
			expectedOperations: []types.Operation{
				clobtest.NewOrderPlacementOperation(
					constants.Order_Alice_Num1_Id10_Clob0_Buy5_Price30_GTB32,
				),
				clobtest.NewOrderPlacementOperation(
					constants.Order_Alice_Num0_Id1_Clob0_Sell10_Price15_GTB15,
				),
				clobtest.NewMatchOperation(
					&constants.Order_Alice_Num0_Id1_Clob0_Sell10_Price15_GTB15,
					[]types.MakerFill{
						{
							MakerOrderId: constants.Order_Alice_Num1_Id10_Clob0_Buy5_Price30_GTB32.OrderId,
							FillAmount:   5,
						},
					},
				),
				clobtest.NewOrderPlacementOperation(
					constants.Order_Alice_Num1_Id10_Clob0_Buy15_Price30_GTB33,
				),
				clobtest.NewMatchOperation(
					&constants.Order_Alice_Num1_Id10_Clob0_Buy15_Price30_GTB33,
					[]types.MakerFill{
						{
							MakerOrderId: constants.Order_Alice_Num0_Id1_Clob0_Sell10_Price15_GTB15.OrderId,
							FillAmount:   5,
						},
					},
				),
			},
			expectedInternalOperations: []types.InternalOperation{
				types.NewShortTermOrderPlacementInternalOperation(
					constants.Order_Alice_Num1_Id10_Clob0_Buy5_Price30_GTB32,
				),
				types.NewShortTermOrderPlacementInternalOperation(
					constants.Order_Alice_Num0_Id1_Clob0_Sell10_Price15_GTB15,
				),
				types.NewMatchOrdersInternalOperation(
					constants.Order_Alice_Num0_Id1_Clob0_Sell10_Price15_GTB15,
					[]types.MakerFill{
						{
							MakerOrderId: constants.Order_Alice_Num1_Id10_Clob0_Buy5_Price30_GTB32.OrderId,
							FillAmount:   5,
						},
					},
				),
				types.NewShortTermOrderPlacementInternalOperation(
					constants.Order_Alice_Num1_Id10_Clob0_Buy15_Price30_GTB33,
				),
				types.NewMatchOrdersInternalOperation(
					constants.Order_Alice_Num1_Id10_Clob0_Buy15_Price30_GTB33,
					[]types.MakerFill{
						{
							MakerOrderId: constants.Order_Alice_Num0_Id1_Clob0_Sell10_Price15_GTB15.OrderId,
							FillAmount:   5,
						},
					},
				),
			},
		},
		"Taker replaces partially matched order with smaller order. Smaller order is added to the book.": {
			placedMatchableOrders: []types.MatchableOrder{
				// Match #1: This order is fully matched before the test case as a maker order with the below order.
				&constants.Order_Alice_Num1_Id10_Clob0_Buy15_Price30_GTB33,
				// Match #1: This order is partially matched before the test case as a taker order with the above order.
				&constants.Order_Alice_Num0_Id1_Clob0_Sell5_Price15_GTB15,
			},

			// This order replaces the existing partially-matched order with an order that decreases the size.
			order: constants.Order_Alice_Num1_Id10_Clob0_Buy10_Price30_GTB34,

			expectedFilledSize:      0,
			expectedTotalFilledSize: 5,
			expectedOrderStatus:     types.Success,
			expectedRemainingBids: []OrderWithRemainingSize{
				{
					Order:         constants.Order_Alice_Num1_Id10_Clob0_Buy10_Price30_GTB34,
					RemainingSize: 5,
				},
			},
			expectedRemainingAsks: []OrderWithRemainingSize{},
			expectedCollatCheck:   []expectedMatch{},
			expectedExistingMatches: []expectedMatch{
				// Match #1: This match is generated before the test case.
				{
					makerOrder:      &constants.Order_Alice_Num1_Id10_Clob0_Buy15_Price30_GTB33,
					takerOrder:      &constants.Order_Alice_Num0_Id1_Clob0_Sell5_Price15_GTB15,
					matchedQuantums: 5,
				},
			},
			expectedNewMatches: []expectedMatch{},
			expectedOperations: []types.Operation{
				clobtest.NewOrderPlacementOperation(
					constants.Order_Alice_Num1_Id10_Clob0_Buy15_Price30_GTB33,
				),
				clobtest.NewOrderPlacementOperation(
					constants.Order_Alice_Num0_Id1_Clob0_Sell5_Price15_GTB15,
				),
				clobtest.NewMatchOperation(
					&constants.Order_Alice_Num0_Id1_Clob0_Sell5_Price15_GTB15,
					[]types.MakerFill{
						{
							MakerOrderId: constants.Order_Alice_Num1_Id10_Clob0_Buy15_Price30_GTB33.OrderId,
							FillAmount:   5,
						},
					},
				),
				clobtest.NewOrderPlacementOperation(
					constants.Order_Alice_Num1_Id10_Clob0_Buy10_Price30_GTB34,
				),
			},
			expectedInternalOperations: []types.InternalOperation{
				types.NewShortTermOrderPlacementInternalOperation(
					constants.Order_Alice_Num1_Id10_Clob0_Buy15_Price30_GTB33,
				),
				types.NewShortTermOrderPlacementInternalOperation(
					constants.Order_Alice_Num0_Id1_Clob0_Sell5_Price15_GTB15,
				),
				types.NewMatchOrdersInternalOperation(
					constants.Order_Alice_Num0_Id1_Clob0_Sell5_Price15_GTB15,
					[]types.MakerFill{
						{
							MakerOrderId: constants.Order_Alice_Num1_Id10_Clob0_Buy15_Price30_GTB33.OrderId,
							FillAmount:   5,
						},
					},
				),
			},
			expectedToReplaceOrder: true,
		},
		"Error: Taker replaces fully matched order with order which has smaller GTB": {
			placedMatchableOrders: []types.MatchableOrder{
				// Match #1: This order is fully matched before the test case as a maker order with the below order.
				&constants.Order_Alice_Num1_Id10_Clob0_Buy5_Price30_GTB32,
				// Match #1: This order is partially matched before the test case as a taker order with the above order.
				&constants.Order_Alice_Num0_Id1_Clob0_Sell5_Price15_GTB15,
			},

			// This order replaces the existing fully-matched order with an order that is the same size, however
			// the replacement has a smaller GTB. The existing order is _not_ on the book, but only in the match queue.
			order: constants.Order_Alice_Num1_Id10_Clob0_Buy5_Price30_GTB31,

			expectedFilledSize:      0,
			expectedTotalFilledSize: 5,
			expectedOrderStatus:     types.InternalError,
			expectedRemainingBids:   []OrderWithRemainingSize{},
			expectedRemainingAsks:   []OrderWithRemainingSize{},
			expectedCollatCheck:     []expectedMatch{},
			expectedExistingMatches: []expectedMatch{
				// Match #1: This match is generated before the test case.
				{
					makerOrder:      &constants.Order_Alice_Num1_Id10_Clob0_Buy5_Price30_GTB32,
					takerOrder:      &constants.Order_Alice_Num0_Id1_Clob0_Sell5_Price15_GTB15,
					matchedQuantums: 5,
				},
			},
			expectedNewMatches: []expectedMatch{},
			expectedOperations: []types.Operation{
				clobtest.NewOrderPlacementOperation(
					constants.Order_Alice_Num1_Id10_Clob0_Buy5_Price30_GTB32,
				),
				clobtest.NewOrderPlacementOperation(
					constants.Order_Alice_Num0_Id1_Clob0_Sell5_Price15_GTB15,
				),
				clobtest.NewMatchOperation(
					&constants.Order_Alice_Num0_Id1_Clob0_Sell5_Price15_GTB15,
					[]types.MakerFill{
						{
							MakerOrderId: constants.Order_Alice_Num1_Id10_Clob0_Buy5_Price30_GTB32.OrderId,
							FillAmount:   5,
						},
					},
				),
			},
			expectedInternalOperations: []types.InternalOperation{
				types.NewShortTermOrderPlacementInternalOperation(
					constants.Order_Alice_Num1_Id10_Clob0_Buy5_Price30_GTB32,
				),
				types.NewShortTermOrderPlacementInternalOperation(
					constants.Order_Alice_Num0_Id1_Clob0_Sell5_Price15_GTB15,
				),
				types.NewMatchOrdersInternalOperation(
					constants.Order_Alice_Num0_Id1_Clob0_Sell5_Price15_GTB15,
					[]types.MakerFill{
						{
							MakerOrderId: constants.Order_Alice_Num1_Id10_Clob0_Buy5_Price30_GTB32.OrderId,
							FillAmount:   5,
						},
					},
				),
			},
			expectedErr: types.ErrInvalidReplacement,
		},
		"Error: Taker replaces fully matched order with an order of the same size": {
			placedMatchableOrders: []types.MatchableOrder{
				// Match #1: This order is fully matched before the test case as a maker order with the below order.
				&constants.Order_Alice_Num1_Id10_Clob0_Buy5_Price30_GTB31,
				// Match #1: This order is partially matched before the test case as a taker order with the above order.
				&constants.Order_Alice_Num0_Id1_Clob0_Sell5_Price15_GTB15,
			},

			// This order replaces the existing fully-matched order with an order that is the same size.
			// The order is therefore already fully filled and an error is returned.
			order: constants.Order_Alice_Num1_Id10_Clob0_Buy5_Price30_GTB32,

			expectedFilledSize:      0,
			expectedTotalFilledSize: 5,
			expectedOrderStatus:     types.InternalError,
			expectedRemainingBids:   []OrderWithRemainingSize{},
			expectedRemainingAsks:   []OrderWithRemainingSize{},
			expectedCollatCheck:     []expectedMatch{},
			expectedExistingMatches: []expectedMatch{
				// Match #1: This match is generated before the test case.
				{
					makerOrder:      &constants.Order_Alice_Num1_Id10_Clob0_Buy5_Price30_GTB31,
					takerOrder:      &constants.Order_Alice_Num0_Id1_Clob0_Sell5_Price15_GTB15,
					matchedQuantums: 5,
				},
			},
			expectedNewMatches: []expectedMatch{},
			expectedOperations: []types.Operation{
				clobtest.NewOrderPlacementOperation(
					constants.Order_Alice_Num1_Id10_Clob0_Buy5_Price30_GTB31,
				),
				clobtest.NewOrderPlacementOperation(
					constants.Order_Alice_Num0_Id1_Clob0_Sell5_Price15_GTB15,
				),
				clobtest.NewMatchOperation(
					&constants.Order_Alice_Num0_Id1_Clob0_Sell5_Price15_GTB15,
					[]types.MakerFill{
						{
							MakerOrderId: constants.Order_Alice_Num1_Id10_Clob0_Buy5_Price30_GTB31.OrderId,
							FillAmount:   5,
						},
					},
				),
			},
			expectedInternalOperations: []types.InternalOperation{
				types.NewShortTermOrderPlacementInternalOperation(
					constants.Order_Alice_Num1_Id10_Clob0_Buy5_Price30_GTB31,
				),
				types.NewShortTermOrderPlacementInternalOperation(
					constants.Order_Alice_Num0_Id1_Clob0_Sell5_Price15_GTB15,
				),
				types.NewMatchOrdersInternalOperation(
					constants.Order_Alice_Num0_Id1_Clob0_Sell5_Price15_GTB15,
					[]types.MakerFill{
						{
							MakerOrderId: constants.Order_Alice_Num1_Id10_Clob0_Buy5_Price30_GTB31.OrderId,
							FillAmount:   5,
						},
					},
				),
			},
			expectedErr: types.ErrOrderFullyFilled,
		},
		`Error: Taker replaces partially matched order with an order of larger size but remaining fillable amount is less than
		MinOrderBaseQuantums`: {
			placedMatchableOrders: []types.MatchableOrder{
				// Match #1: This order is partially matched before the test case as a maker order with the below order.
				&constants.Order_Alice_Num1_Id10_Clob0_Buy6_Price30_GTB32,
				// Match #1: This order is fully matched before the test case as a taker order with the above order.
				&constants.Order_Alice_Num0_Id1_Clob0_Sell5_Price15_GTB15,
			},

			// This order replaces the existing partially-matched order with an order that is is larger by 1 base quantum.
			// The CLOB has a `MinOrderBaseQuantums` of 5 so therefore the placement could only possibly result in a fill
			// of size 2, which is lower than the `MinOrderBaseQuantums` of the orderbook, therefore an error is returned.
			order: constants.Order_Alice_Num1_Id10_Clob0_Buy7_Price30_GTB33,

			expectedFilledSize:      0,
			expectedTotalFilledSize: 5,
			expectedOrderStatus:     types.InternalError,
			expectedRemainingBids: []OrderWithRemainingSize{
				{
					Order:         constants.Order_Alice_Num1_Id10_Clob0_Buy6_Price30_GTB32,
					RemainingSize: 1,
				},
			},
			expectedRemainingAsks: []OrderWithRemainingSize{},
			expectedCollatCheck:   []expectedMatch{},
			expectedExistingMatches: []expectedMatch{
				// Match #1: This match is generated before the test case.
				{
					makerOrder:      &constants.Order_Alice_Num1_Id10_Clob0_Buy6_Price30_GTB32,
					takerOrder:      &constants.Order_Alice_Num0_Id1_Clob0_Sell5_Price15_GTB15,
					matchedQuantums: 5,
				},
			},
			expectedNewMatches: []expectedMatch{},
			expectedOperations: []types.Operation{
				clobtest.NewOrderPlacementOperation(
					constants.Order_Alice_Num1_Id10_Clob0_Buy6_Price30_GTB32,
				),
				clobtest.NewOrderPlacementOperation(
					constants.Order_Alice_Num0_Id1_Clob0_Sell5_Price15_GTB15,
				),
				clobtest.NewMatchOperation(
					&constants.Order_Alice_Num0_Id1_Clob0_Sell5_Price15_GTB15,
					[]types.MakerFill{
						{
							MakerOrderId: constants.Order_Alice_Num1_Id10_Clob0_Buy6_Price30_GTB32.OrderId,
							FillAmount:   5,
						},
					},
				),
			},
			expectedInternalOperations: []types.InternalOperation{
				types.NewShortTermOrderPlacementInternalOperation(
					constants.Order_Alice_Num1_Id10_Clob0_Buy6_Price30_GTB32,
				),
				types.NewShortTermOrderPlacementInternalOperation(
					constants.Order_Alice_Num0_Id1_Clob0_Sell5_Price15_GTB15,
				),
				types.NewMatchOrdersInternalOperation(
					constants.Order_Alice_Num0_Id1_Clob0_Sell5_Price15_GTB15,
					[]types.MakerFill{
						{
							MakerOrderId: constants.Order_Alice_Num1_Id10_Clob0_Buy6_Price30_GTB32.OrderId,
							FillAmount:   5,
						},
					},
				),
			},
			expectedErr: types.ErrOrderFullyFilled,
		},
		"Error: Taker replaces fully matched order with an order of smaller size and order is now fully filled": {
			placedMatchableOrders: []types.MatchableOrder{
				// Match #1: This order is fully matched before the test case as a maker order with the below order.
				&constants.Order_Alice_Num1_Id10_Clob0_Buy15_Price30_GTB33,
				// Match #1: This order is partially matched before the test case as a taker order with the above order.
				&constants.Order_Alice_Num0_Id1_Clob0_Sell5_Price15_GTB15,
			},

			// This order replaces the existing partially-matched order with an order that is smaller and equal to the total
			// fill size of the order. This means that the replacement would leave the order fully filled and therefore an
			// error is returned.
			order: constants.Order_Alice_Num1_Id10_Clob0_Buy5_Price30_GTB34,

			expectedFilledSize:      0,
			expectedTotalFilledSize: 5,
			expectedOrderStatus:     types.InternalError,
			expectedRemainingBids: []OrderWithRemainingSize{
				{
					Order:         constants.Order_Alice_Num1_Id10_Clob0_Buy15_Price30_GTB33,
					RemainingSize: 10,
				},
			},
			expectedRemainingAsks: []OrderWithRemainingSize{},
			expectedCollatCheck:   []expectedMatch{},
			expectedExistingMatches: []expectedMatch{
				// Match #1: This match is generated before the test case.
				{
					makerOrder:      &constants.Order_Alice_Num1_Id10_Clob0_Buy15_Price30_GTB33,
					takerOrder:      &constants.Order_Alice_Num0_Id1_Clob0_Sell5_Price15_GTB15,
					matchedQuantums: 5,
				},
			},
			expectedNewMatches: []expectedMatch{},
			expectedOperations: []types.Operation{
				clobtest.NewOrderPlacementOperation(
					constants.Order_Alice_Num1_Id10_Clob0_Buy15_Price30_GTB33,
				),
				clobtest.NewOrderPlacementOperation(
					constants.Order_Alice_Num0_Id1_Clob0_Sell5_Price15_GTB15,
				),
				clobtest.NewMatchOperation(
					&constants.Order_Alice_Num0_Id1_Clob0_Sell5_Price15_GTB15,
					[]types.MakerFill{
						{
							MakerOrderId: constants.Order_Alice_Num1_Id10_Clob0_Buy15_Price30_GTB33.OrderId,
							FillAmount:   5,
						},
					},
				),
			},
			expectedInternalOperations: []types.InternalOperation{
				types.NewShortTermOrderPlacementInternalOperation(
					constants.Order_Alice_Num1_Id10_Clob0_Buy15_Price30_GTB33,
				),
				types.NewShortTermOrderPlacementInternalOperation(
					constants.Order_Alice_Num0_Id1_Clob0_Sell5_Price15_GTB15,
				),
				types.NewMatchOrdersInternalOperation(
					constants.Order_Alice_Num0_Id1_Clob0_Sell5_Price15_GTB15,
					[]types.MakerFill{
						{
							MakerOrderId: constants.Order_Alice_Num1_Id10_Clob0_Buy15_Price30_GTB33.OrderId,
							FillAmount:   5,
						},
					},
				),
			},
			expectedErr: types.ErrOrderFullyFilled,
		},
		"Error: Taker is IOC replacement for partially filled IOC order": {
			placedMatchableOrders: []types.MatchableOrder{
				&constants.Order_Bob_Num0_Id11_Clob1_Buy5_Price40_GTB20,
				&constants.Order_Alice_Num1_Id1_Clob1_Sell10_Price15_GTB20_IOC,
			},

			order: constants.Order_Alice_Num1_Id1_Clob1_Sell10_Price15_GTB21_IOC,

			expectedTotalFilledSize: 5,
			expectedOrderStatus:     types.InternalError,
			expectedExistingMatches: []expectedMatch{
				{
					makerOrder:      &constants.Order_Bob_Num0_Id11_Clob1_Buy5_Price40_GTB20,
					takerOrder:      &constants.Order_Alice_Num1_Id1_Clob1_Sell10_Price15_GTB20_IOC,
					matchedQuantums: 5,
				},
			},
			expectedOperations: []types.Operation{
				clobtest.NewOrderPlacementOperation(
					constants.Order_Bob_Num0_Id11_Clob1_Buy5_Price40_GTB20,
				),
				clobtest.NewOrderPlacementOperation(
					constants.Order_Alice_Num1_Id1_Clob1_Sell10_Price15_GTB20_IOC,
				),
				clobtest.NewMatchOperation(
					&constants.Order_Alice_Num1_Id1_Clob1_Sell10_Price15_GTB20_IOC,
					[]types.MakerFill{
						{
							MakerOrderId: constants.Order_Bob_Num0_Id11_Clob1_Buy5_Price40_GTB20.OrderId,
							FillAmount:   5,
						},
					},
				),
			},
			expectedInternalOperations: []types.InternalOperation{
				types.NewShortTermOrderPlacementInternalOperation(
					constants.Order_Bob_Num0_Id11_Clob1_Buy5_Price40_GTB20,
				),
				types.NewShortTermOrderPlacementInternalOperation(
					constants.Order_Alice_Num1_Id1_Clob1_Sell10_Price15_GTB20_IOC,
				),
				types.NewMatchOrdersInternalOperation(
					constants.Order_Alice_Num1_Id1_Clob1_Sell10_Price15_GTB20_IOC,
					[]types.MakerFill{
						{
							MakerOrderId: constants.Order_Bob_Num0_Id11_Clob1_Buy5_Price40_GTB20.OrderId,
							FillAmount:   5,
						},
					},
				),
			},

			expectedErr: types.ErrImmediateExecutionOrderAlreadyFilled,
		},
		"Error: Taker is IOC replacement for partially filled order": {
			placedMatchableOrders: []types.MatchableOrder{
				&constants.Order_Bob_Num0_Id11_Clob1_Buy5_Price40_GTB20,
				&constants.Order_Alice_Num1_Id1_Clob1_Sell10_Price15_GTB20,
			},

			order: constants.Order_Alice_Num1_Id1_Clob1_Sell10_Price15_GTB21_IOC,

			expectedTotalFilledSize: 5,
			expectedOrderStatus:     types.InternalError,
			expectedExistingMatches: []expectedMatch{
				{
					makerOrder:      &constants.Order_Bob_Num0_Id11_Clob1_Buy5_Price40_GTB20,
					takerOrder:      &constants.Order_Alice_Num1_Id1_Clob1_Sell10_Price15_GTB20,
					matchedQuantums: 5,
				},
			},
			expectedRemainingAsks: []OrderWithRemainingSize{
				{
					Order:         constants.Order_Alice_Num1_Id1_Clob1_Sell10_Price15_GTB20,
					RemainingSize: 5,
				},
			},
			expectedOperations: []types.Operation{
				clobtest.NewOrderPlacementOperation(
					constants.Order_Bob_Num0_Id11_Clob1_Buy5_Price40_GTB20,
				),
				clobtest.NewOrderPlacementOperation(
					constants.Order_Alice_Num1_Id1_Clob1_Sell10_Price15_GTB20,
				),
				clobtest.NewMatchOperation(
					&constants.Order_Alice_Num1_Id1_Clob1_Sell10_Price15_GTB20,
					[]types.MakerFill{
						{
							MakerOrderId: constants.Order_Bob_Num0_Id11_Clob1_Buy5_Price40_GTB20.OrderId,
							FillAmount:   5,
						},
					},
				),
			},
			expectedInternalOperations: []types.InternalOperation{
				types.NewShortTermOrderPlacementInternalOperation(
					constants.Order_Bob_Num0_Id11_Clob1_Buy5_Price40_GTB20,
				),
				types.NewShortTermOrderPlacementInternalOperation(
					constants.Order_Alice_Num1_Id1_Clob1_Sell10_Price15_GTB20,
				),
				types.NewMatchOrdersInternalOperation(
					constants.Order_Alice_Num1_Id1_Clob1_Sell10_Price15_GTB20,
					[]types.MakerFill{
						{
							MakerOrderId: constants.Order_Bob_Num0_Id11_Clob1_Buy5_Price40_GTB20.OrderId,
							FillAmount:   5,
						},
					},
				),
			},

			expectedErr: types.ErrInvalidReplacement,
		},
		"IOC Taker replaces unfilled non IOC order": {
			placedMatchableOrders: []types.MatchableOrder{
				&constants.Order_Alice_Num1_Id1_Clob1_Sell10_Price15_GTB20,
			},
			expectedRemainingAsks: []OrderWithRemainingSize{
				{
					Order:         constants.Order_Alice_Num1_Id1_Clob1_Sell10_Price15_GTB20,
					RemainingSize: 10,
				},
			},

			order:                      constants.Order_Alice_Num1_Id1_Clob1_Sell10_Price15_GTB21_IOC,
			expectedInternalOperations: []types.InternalOperation{},

			expectedErr: types.ErrInvalidReplacement,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// Setup memclob state and test expectations.
			addOrderToOrderbookSize := satypes.BaseQuantums(0)
			if tc.expectedOrderStatus.IsSuccess() {
				addOrderToOrderbookSize = tc.order.GetBaseQuantums() - tc.expectedTotalFilledSize
			}

			memclob, _, expectedNumCollateralizationChecks, numCollateralChecks := placeOrderTestSetup(
				t,
				ctx,
				tc.placedMatchableOrders,
				&tc.order,
				tc.expectedCollatCheck,
				tc.expectedOrderStatus,
				addOrderToOrderbookSize,
				tc.expectedErr,
				map[int]map[satypes.SubaccountId]satypes.UpdateResult{},
				constants.GetStatePosition_ZeroPositionSize,
			)

			// Run the test case and verify expectations.
			offchainUpdates := placeOrderAndVerifyExpectationsOperations(
				t,
				ctx,
				memclob,
				tc.order,
				numCollateralChecks,
				tc.expectedFilledSize,
				tc.expectedTotalFilledSize,
				tc.expectedOrderStatus,
				tc.expectedErr,
				expectedNumCollateralizationChecks,
				tc.expectedRemainingBids,
				tc.expectedRemainingAsks,
				tc.expectedOperations,
				tc.expectedInternalOperations,
				nil,
			)

			// Verify the correct offchain update messages were returned.
			assertPlaceOrderOffchainMessages(
				t,
				ctx,
				offchainUpdates,
				tc.order,
				tc.placedMatchableOrders,
				map[int]map[satypes.SubaccountId]satypes.UpdateResult{},
				tc.expectedErr,
				tc.expectedTotalFilledSize,
				tc.expectedOrderStatus,
				tc.expectedExistingMatches,
				tc.expectedNewMatches,
				[]types.OrderId{},
				tc.expectedToReplaceOrder,
			)
		})
	}
}

func TestPlaceOrder_MatchOrders_CollatCheckFailure(t *testing.T) {
	ctx, _, _ := sdktest.NewSdkContextWithMultistore()
	ctx = ctx.WithIsCheckTx(true)
	tests := map[string]struct {
		// State.
		placedMatchableOrders          []types.MatchableOrder
		collateralizationCheckFailures map[int]map[satypes.SubaccountId]satypes.UpdateResult

		// Parameters.
		order types.Order

		// Expectations.
		expectedFilledSize         satypes.BaseQuantums
		expectedOrderStatus        types.OrderStatus
		expectedCollatCheck        []expectedMatch
		expectedRemainingBids      []OrderWithRemainingSize
		expectedRemainingAsks      []OrderWithRemainingSize
		expectedMatches            []expectedMatch
		expectedOperations         []types.Operation
		expectedInternalOperations []types.InternalOperation
	}{
		`When all maker orders that cross the taker order fail collateralization checks, they are removed and the taker
		order is placed on the book`: {
			placedMatchableOrders: []types.MatchableOrder{
				&constants.Order_Alice_Num0_Id2_Clob1_Sell5_Price10_GTB15,
				&constants.Order_Alice_Num0_Id3_Clob1_Sell5_Price10_GTB15,
				&constants.Order_Alice_Num1_Id1_Clob1_Sell10_Price15_GTB20,
			},
			collateralizationCheckFailures: map[int]map[satypes.SubaccountId]satypes.UpdateResult{
				0: {
					constants.Alice_Num0: satypes.NewlyUndercollateralized,
				},
				1: {
					constants.Alice_Num0: satypes.StillUndercollateralized,
				},
				2: {
					constants.Alice_Num1: satypes.UpdateCausedError,
				},
			},

			order: constants.Order_Bob_Num0_Id4_Clob1_Buy20_Price35_GTB22,

			expectedFilledSize:  0,
			expectedOrderStatus: types.Success,
			expectedRemainingBids: []OrderWithRemainingSize{
				{
					Order:         constants.Order_Bob_Num0_Id4_Clob1_Buy20_Price35_GTB22,
					RemainingSize: 20,
				},
			},
			expectedRemainingAsks: []OrderWithRemainingSize{},
			expectedCollatCheck: []expectedMatch{
				{
					makerOrder:      &constants.Order_Alice_Num0_Id2_Clob1_Sell5_Price10_GTB15,
					takerOrder:      &constants.Order_Bob_Num0_Id4_Clob1_Buy20_Price35_GTB22,
					matchedQuantums: 5,
				},
				{
					makerOrder:      &constants.Order_Alice_Num0_Id3_Clob1_Sell5_Price10_GTB15,
					takerOrder:      &constants.Order_Bob_Num0_Id4_Clob1_Buy20_Price35_GTB22,
					matchedQuantums: 5,
				},
				{
					makerOrder:      &constants.Order_Alice_Num1_Id1_Clob1_Sell10_Price15_GTB20,
					takerOrder:      &constants.Order_Bob_Num0_Id4_Clob1_Buy20_Price35_GTB22,
					matchedQuantums: 10,
				},
			},
			expectedMatches:            []expectedMatch{},
			expectedOperations:         []types.Operation{},
			expectedInternalOperations: []types.InternalOperation{},
		},
		`Matching stops if taker order fails collateralization, and no orders are removed from the book`: {
			placedMatchableOrders: []types.MatchableOrder{
				&constants.Order_Alice_Num0_Id2_Clob1_Sell5_Price10_GTB15,
				&constants.Order_Alice_Num0_Id3_Clob1_Sell5_Price10_GTB15,
				&constants.Order_Alice_Num1_Id1_Clob1_Sell10_Price15_GTB20,
			},
			collateralizationCheckFailures: map[int]map[satypes.SubaccountId]satypes.UpdateResult{
				0: {
					constants.Bob_Num0: satypes.StillUndercollateralized,
				},
			},

			order: constants.Order_Bob_Num0_Id4_Clob1_Buy20_Price35_GTB22,

			expectedFilledSize:    0,
			expectedOrderStatus:   types.Undercollateralized,
			expectedRemainingBids: []OrderWithRemainingSize{},
			expectedRemainingAsks: []OrderWithRemainingSize{
				{
					Order:         constants.Order_Alice_Num0_Id2_Clob1_Sell5_Price10_GTB15,
					RemainingSize: 5,
				},
				{
					Order:         constants.Order_Alice_Num0_Id3_Clob1_Sell5_Price10_GTB15,
					RemainingSize: 5,
				},
				{
					Order:         constants.Order_Alice_Num1_Id1_Clob1_Sell10_Price15_GTB20,
					RemainingSize: 10,
				},
			},
			expectedCollatCheck: []expectedMatch{
				{
					makerOrder:      &constants.Order_Alice_Num0_Id2_Clob1_Sell5_Price10_GTB15,
					takerOrder:      &constants.Order_Bob_Num0_Id4_Clob1_Buy20_Price35_GTB22,
					matchedQuantums: 5,
				},
			},
			expectedMatches:            []expectedMatch{},
			expectedOperations:         []types.Operation{},
			expectedInternalOperations: []types.InternalOperation{},
		},
		`Matching stops if taker order is partially filled then fails collateralization, and all filled maker orders are
		removed from the book`: {
			placedMatchableOrders: []types.MatchableOrder{
				&constants.Order_Alice_Num1_Clob0_Id4_Buy10_Price45_GTB20,
				&constants.Order_Alice_Num1_Id10_Clob0_Buy5_Price30_GTB32,
				&constants.Order_Alice_Num1_Id8_Clob0_Buy15_Price25_GTB31,
				&constants.Order_Bob_Num0_Id9_Clob0_Sell20_Price1000,
			},
			collateralizationCheckFailures: map[int]map[satypes.SubaccountId]satypes.UpdateResult{
				2: {
					constants.Alice_Num0: satypes.NewlyUndercollateralized,
				},
			},

			order: constants.Order_Alice_Num0_Id7_Clob0_Sell25_Price15_GTB20,

			expectedFilledSize:  15,
			expectedOrderStatus: types.Undercollateralized,
			expectedRemainingBids: []OrderWithRemainingSize{
				{
					Order:         constants.Order_Alice_Num1_Id8_Clob0_Buy15_Price25_GTB31,
					RemainingSize: 15,
				},
			},
			expectedRemainingAsks: []OrderWithRemainingSize{
				{
					Order:         constants.Order_Bob_Num0_Id9_Clob0_Sell20_Price1000,
					RemainingSize: 20,
				},
			},
			expectedCollatCheck: []expectedMatch{
				{
					makerOrder:      &constants.Order_Alice_Num1_Clob0_Id4_Buy10_Price45_GTB20,
					takerOrder:      &constants.Order_Alice_Num0_Id7_Clob0_Sell25_Price15_GTB20,
					matchedQuantums: 10,
				},
				{
					makerOrder:      &constants.Order_Alice_Num1_Id10_Clob0_Buy5_Price30_GTB32,
					takerOrder:      &constants.Order_Alice_Num0_Id7_Clob0_Sell25_Price15_GTB20,
					matchedQuantums: 5,
				},
				{
					makerOrder:      &constants.Order_Alice_Num1_Id8_Clob0_Buy15_Price25_GTB31,
					takerOrder:      &constants.Order_Alice_Num0_Id7_Clob0_Sell25_Price15_GTB20,
					matchedQuantums: 10,
				},
			},
			expectedMatches: []expectedMatch{
				{
					makerOrder:      &constants.Order_Alice_Num1_Clob0_Id4_Buy10_Price45_GTB20,
					takerOrder:      &constants.Order_Alice_Num0_Id7_Clob0_Sell25_Price15_GTB20,
					matchedQuantums: 10,
				},
				{
					makerOrder:      &constants.Order_Alice_Num1_Id10_Clob0_Buy5_Price30_GTB32,
					takerOrder:      &constants.Order_Alice_Num0_Id7_Clob0_Sell25_Price15_GTB20,
					matchedQuantums: 5,
				},
			},
			expectedOperations: []types.Operation{
				clobtest.NewOrderPlacementOperation(
					constants.Order_Alice_Num1_Clob0_Id4_Buy10_Price45_GTB20,
				),
				clobtest.NewOrderPlacementOperation(
					constants.Order_Alice_Num1_Id10_Clob0_Buy5_Price30_GTB32,
				),
				clobtest.NewOrderPlacementOperation(
					constants.Order_Alice_Num1_Id8_Clob0_Buy15_Price25_GTB31,
				),
				clobtest.NewOrderPlacementOperation(
					constants.Order_Alice_Num0_Id7_Clob0_Sell25_Price15_GTB20,
				),
				clobtest.NewMatchOperation(
					&constants.Order_Alice_Num0_Id7_Clob0_Sell25_Price15_GTB20,
					[]types.MakerFill{
						{
							MakerOrderId: constants.Order_Alice_Num1_Clob0_Id4_Buy10_Price45_GTB20.OrderId,
							FillAmount:   10,
						},
						{
							MakerOrderId: constants.Order_Alice_Num1_Id10_Clob0_Buy5_Price30_GTB32.OrderId,
							FillAmount:   5,
						},
					},
				),
			},
			expectedInternalOperations: []types.InternalOperation{
				types.NewShortTermOrderPlacementInternalOperation(
					constants.Order_Alice_Num1_Clob0_Id4_Buy10_Price45_GTB20,
				),
				types.NewShortTermOrderPlacementInternalOperation(
					constants.Order_Alice_Num1_Id10_Clob0_Buy5_Price30_GTB32,
				),
				types.NewShortTermOrderPlacementInternalOperation(
					constants.Order_Alice_Num0_Id7_Clob0_Sell25_Price15_GTB20,
				),
				types.NewMatchOrdersInternalOperation(
					constants.Order_Alice_Num0_Id7_Clob0_Sell25_Price15_GTB20,
					[]types.MakerFill{
						{
							MakerOrderId: constants.Order_Alice_Num1_Clob0_Id4_Buy10_Price45_GTB20.OrderId,
							FillAmount:   10,
						},
						{
							MakerOrderId: constants.Order_Alice_Num1_Id10_Clob0_Buy5_Price30_GTB32.OrderId,
							FillAmount:   5,
						},
					},
				),
			},
		},
		`Matching stops if taker order fails collateralization, all partial fills are added to the match queue, and all
		maker orders that failed collateralization are removed from the book`: {
			placedMatchableOrders: []types.MatchableOrder{
				&constants.Order_Alice_Num0_Id2_Clob1_Sell5_Price10_GTB15,
				&constants.Order_Alice_Num0_Id3_Clob1_Sell5_Price10_GTB15,
				&constants.Order_Alice_Num1_Id1_Clob1_Sell10_Price15_GTB20,
			},
			collateralizationCheckFailures: map[int]map[satypes.SubaccountId]satypes.UpdateResult{
				1: {
					constants.Alice_Num0: satypes.NewlyUndercollateralized,
				},
				2: {
					constants.Bob_Num0: satypes.UpdateCausedError,
				},
			},

			order: constants.Order_Bob_Num0_Id4_Clob1_Buy20_Price35_GTB22,

			expectedFilledSize:    5,
			expectedOrderStatus:   types.InternalError,
			expectedRemainingBids: []OrderWithRemainingSize{},
			expectedRemainingAsks: []OrderWithRemainingSize{
				{
					Order:         constants.Order_Alice_Num1_Id1_Clob1_Sell10_Price15_GTB20,
					RemainingSize: 10,
				},
			},
			expectedCollatCheck: []expectedMatch{
				{
					makerOrder:      &constants.Order_Alice_Num0_Id2_Clob1_Sell5_Price10_GTB15,
					takerOrder:      &constants.Order_Bob_Num0_Id4_Clob1_Buy20_Price35_GTB22,
					matchedQuantums: 5,
				},
				{
					makerOrder:      &constants.Order_Alice_Num0_Id3_Clob1_Sell5_Price10_GTB15,
					takerOrder:      &constants.Order_Bob_Num0_Id4_Clob1_Buy20_Price35_GTB22,
					matchedQuantums: 5,
				},
				{
					makerOrder:      &constants.Order_Alice_Num1_Id1_Clob1_Sell10_Price15_GTB20,
					takerOrder:      &constants.Order_Bob_Num0_Id4_Clob1_Buy20_Price35_GTB22,
					matchedQuantums: 10,
				},
			},
			expectedMatches: []expectedMatch{
				{
					makerOrder:      &constants.Order_Alice_Num0_Id2_Clob1_Sell5_Price10_GTB15,
					takerOrder:      &constants.Order_Bob_Num0_Id4_Clob1_Buy20_Price35_GTB22,
					matchedQuantums: 5,
				},
			},
			expectedOperations: []types.Operation{
				clobtest.NewOrderPlacementOperation(
					constants.Order_Alice_Num0_Id2_Clob1_Sell5_Price10_GTB15,
				),
				clobtest.NewOrderPlacementOperation(
					constants.Order_Alice_Num1_Id1_Clob1_Sell10_Price15_GTB20,
				),
				clobtest.NewOrderPlacementOperation(
					constants.Order_Bob_Num0_Id4_Clob1_Buy20_Price35_GTB22,
				),
				clobtest.NewMatchOperation(
					&constants.Order_Bob_Num0_Id4_Clob1_Buy20_Price35_GTB22,
					[]types.MakerFill{
						{
							MakerOrderId: constants.Order_Alice_Num0_Id2_Clob1_Sell5_Price10_GTB15.OrderId,
							FillAmount:   5,
						},
					},
				),
			},
			expectedInternalOperations: []types.InternalOperation{
				types.NewShortTermOrderPlacementInternalOperation(
					constants.Order_Alice_Num0_Id2_Clob1_Sell5_Price10_GTB15,
				),
				types.NewShortTermOrderPlacementInternalOperation(
					constants.Order_Bob_Num0_Id4_Clob1_Buy20_Price35_GTB22,
				),
				types.NewMatchOrdersInternalOperation(
					constants.Order_Bob_Num0_Id4_Clob1_Buy20_Price35_GTB22,
					[]types.MakerFill{
						{
							MakerOrderId: constants.Order_Alice_Num0_Id2_Clob1_Sell5_Price10_GTB15.OrderId,
							FillAmount:   5,
						},
					},
				),
			},
		},
		`Matching stops if taker and maker order fail collateralization, and maker order is removed from the book`: {
			placedMatchableOrders: []types.MatchableOrder{
				&constants.Order_Alice_Num0_Id2_Clob1_Sell5_Price10_GTB15,
				&constants.Order_Alice_Num0_Id3_Clob1_Sell5_Price10_GTB15,
				&constants.Order_Alice_Num1_Id2_Clob1_Buy67_Price5_GTB20,
			},
			collateralizationCheckFailures: map[int]map[satypes.SubaccountId]satypes.UpdateResult{
				0: {
					constants.Alice_Num0: satypes.NewlyUndercollateralized,
					constants.Bob_Num0:   satypes.StillUndercollateralized,
				},
			},

			order: constants.Order_Bob_Num0_Id4_Clob1_Buy20_Price35_GTB22,

			expectedFilledSize:  0,
			expectedOrderStatus: types.Undercollateralized,
			expectedRemainingBids: []OrderWithRemainingSize{
				{
					Order:         constants.Order_Alice_Num1_Id2_Clob1_Buy67_Price5_GTB20,
					RemainingSize: 67,
				},
			},
			expectedRemainingAsks: []OrderWithRemainingSize{
				{
					Order:         constants.Order_Alice_Num0_Id3_Clob1_Sell5_Price10_GTB15,
					RemainingSize: 5,
				},
			},
			expectedCollatCheck: []expectedMatch{
				{
					makerOrder:      &constants.Order_Alice_Num0_Id2_Clob1_Sell5_Price10_GTB15,
					takerOrder:      &constants.Order_Bob_Num0_Id4_Clob1_Buy20_Price35_GTB22,
					matchedQuantums: 5,
				},
			},
			expectedMatches:            []expectedMatch{},
			expectedOperations:         []types.Operation{},
			expectedInternalOperations: []types.InternalOperation{},
		},
		`Matching stops, taker order is added to the book, and taker causes a
			partially filled maker order to fail collateralization checks`: {
			placedMatchableOrders: []types.MatchableOrder{
				&constants.Order_Alice_Num1_Id10_Clob0_Buy5_Price30_GTB32,
				&constants.Order_Alice_Num0_Id1_Clob0_Sell10_Price15_GTB15,
			},
			collateralizationCheckFailures: map[int]map[satypes.SubaccountId]satypes.UpdateResult{
				0: {
					constants.Alice_Num0: satypes.NewlyUndercollateralized,
				},
			},

			order: constants.Order_Alice_Num1_Clob0_Id4_Buy10_Price45_GTB20,

			expectedFilledSize:  0,
			expectedOrderStatus: types.Success,
			expectedRemainingBids: []OrderWithRemainingSize{
				{
					Order:         constants.Order_Alice_Num1_Clob0_Id4_Buy10_Price45_GTB20,
					RemainingSize: 10,
				},
			},
			expectedRemainingAsks: []OrderWithRemainingSize{},
			expectedCollatCheck: []expectedMatch{
				{
					makerOrder:      &constants.Order_Alice_Num0_Id1_Clob0_Sell10_Price15_GTB15,
					takerOrder:      &constants.Order_Alice_Num1_Clob0_Id4_Buy10_Price45_GTB20,
					matchedQuantums: 5,
				},
			},
			expectedMatches: []expectedMatch{
				{
					makerOrder:      &constants.Order_Alice_Num1_Id10_Clob0_Buy5_Price30_GTB32,
					takerOrder:      &constants.Order_Alice_Num0_Id1_Clob0_Sell10_Price15_GTB15,
					matchedQuantums: 5,
				},
			},
			expectedOperations: []types.Operation{
				clobtest.NewOrderPlacementOperation(
					constants.Order_Alice_Num1_Id10_Clob0_Buy5_Price30_GTB32,
				),
				clobtest.NewOrderPlacementOperation(
					constants.Order_Alice_Num0_Id1_Clob0_Sell10_Price15_GTB15,
				),
				clobtest.NewMatchOperation(
					&constants.Order_Alice_Num0_Id1_Clob0_Sell10_Price15_GTB15,
					[]types.MakerFill{
						{
							MakerOrderId: constants.Order_Alice_Num1_Id10_Clob0_Buy5_Price30_GTB32.OrderId,
							FillAmount:   5,
						},
					},
				),
				// Note that this order does not match with any orders afterwards.
				clobtest.NewOrderPlacementOperation(
					constants.Order_Alice_Num1_Clob0_Id4_Buy10_Price45_GTB20,
				),
			},
			expectedInternalOperations: []types.InternalOperation{
				types.NewShortTermOrderPlacementInternalOperation(
					constants.Order_Alice_Num1_Id10_Clob0_Buy5_Price30_GTB32,
				),
				types.NewShortTermOrderPlacementInternalOperation(
					constants.Order_Alice_Num0_Id1_Clob0_Sell10_Price15_GTB15,
				),
				types.NewMatchOrdersInternalOperation(
					constants.Order_Alice_Num0_Id1_Clob0_Sell10_Price15_GTB15,
					[]types.MakerFill{
						{
							MakerOrderId: constants.Order_Alice_Num1_Id10_Clob0_Buy5_Price30_GTB32.OrderId,
							FillAmount:   5,
						},
					},
				),
			},
		},
		`Placing an order without a builder code parameter matches successfully with an existing sell order`: {
			placedMatchableOrders: []types.MatchableOrder{
				&constants.Order_Bob_Num0_Id8_Clob0_Sell20_Price10_GTB22,
			},

			order:                          constants.Order_Alice_Num0_Id0_Clob0_Buy10_Price10_GTB20,
			expectedFilledSize:             10,
			expectedOrderStatus:            types.Success,
			collateralizationCheckFailures: map[int]map[satypes.SubaccountId]satypes.UpdateResult{},
			expectedRemainingBids:          []OrderWithRemainingSize{},
			expectedRemainingAsks: []OrderWithRemainingSize{
				{
					Order:         constants.Order_Bob_Num0_Id8_Clob0_Sell20_Price10_GTB22,
					RemainingSize: 10,
				},
			},
			expectedCollatCheck: []expectedMatch{
				{
					makerOrder:      &constants.Order_Bob_Num0_Id8_Clob0_Sell20_Price10_GTB22,
					takerOrder:      &constants.Order_Alice_Num0_Id0_Clob0_Buy10_Price10_GTB20,
					matchedQuantums: 10,
				},
			},
			expectedOperations: []types.Operation{
				clobtest.NewOrderPlacementOperation(constants.Order_Bob_Num0_Id8_Clob0_Sell20_Price10_GTB22),
				clobtest.NewOrderPlacementOperation(constants.Order_Alice_Num0_Id0_Clob0_Buy10_Price10_GTB20),
				clobtest.NewMatchOperation(
					&constants.Order_Alice_Num0_Id0_Clob0_Buy10_Price10_GTB20,
					[]types.MakerFill{
						{
							MakerOrderId: constants.Order_Bob_Num0_Id8_Clob0_Sell20_Price10_GTB22.OrderId,
							FillAmount:   10,
						},
					},
				),
			},
			expectedInternalOperations: []types.InternalOperation{
				types.NewShortTermOrderPlacementInternalOperation(constants.Order_Bob_Num0_Id8_Clob0_Sell20_Price10_GTB22),
				types.NewShortTermOrderPlacementInternalOperation(constants.Order_Alice_Num0_Id0_Clob0_Buy10_Price10_GTB20),
				types.NewMatchOrdersInternalOperation(
					constants.Order_Alice_Num0_Id0_Clob0_Buy10_Price10_GTB20,
					[]types.MakerFill{
						{
							MakerOrderId: constants.Order_Bob_Num0_Id8_Clob0_Sell20_Price10_GTB22.OrderId,
							FillAmount:   10,
						},
					},
				),
			},
		},
		`Placing an order with a builder code parameter matches does not match an existing sell order
			because it will fail collateralization checks due to builder fees`: {
			placedMatchableOrders: []types.MatchableOrder{
				&constants.LongTermOrder_Bob_Num0_Id0_Clob0_Sell10_Price10_GTBT10_PO,
			},

			order:               constants.Order_Alice_Num0_Id0_Clob0_Buy10_Price10_GTB20_BuilderCode,
			expectedFilledSize:  0,
			expectedOrderStatus: types.Undercollateralized,
			collateralizationCheckFailures: map[int]map[satypes.SubaccountId]satypes.UpdateResult{
				0: {
					constants.Alice_Num0: satypes.NewlyUndercollateralized,
				},
			},
			expectedRemainingBids: []OrderWithRemainingSize{},
			expectedRemainingAsks: []OrderWithRemainingSize{
				{
					Order:         constants.LongTermOrder_Bob_Num0_Id0_Clob0_Sell10_Price10_GTBT10_PO,
					RemainingSize: 10,
				},
			},
			expectedCollatCheck: []expectedMatch{
				{
					makerOrder:      &constants.LongTermOrder_Bob_Num0_Id0_Clob0_Sell10_Price10_GTBT10_PO,
					takerOrder:      &constants.Order_Alice_Num0_Id0_Clob0_Buy10_Price10_GTB20_BuilderCode,
					matchedQuantums: 10,
				},
			},
			expectedOperations: []types.Operation{
				clobtest.NewOrderPlacementOperation(constants.LongTermOrder_Bob_Num0_Id0_Clob0_Sell10_Price10_GTBT10_PO),
				clobtest.NewOrderPlacementOperation(constants.Order_Alice_Num0_Id0_Clob0_Buy10_Price10_GTB20_BuilderCode),
			},
			expectedInternalOperations: []types.InternalOperation{},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// Setup memclob state and test expectations.
			addOrderToOrderbookSize := satypes.BaseQuantums(0)
			if tc.expectedOrderStatus.IsSuccess() {
				addOrderToOrderbookSize = tc.order.GetBaseQuantums() - tc.expectedFilledSize
			}
			memclob, _, expectedNumCollateralizationChecks, numCollateralChecks := placeOrderTestSetup(
				t,
				ctx,
				tc.placedMatchableOrders,
				&tc.order,
				tc.expectedCollatCheck,
				tc.expectedOrderStatus,
				addOrderToOrderbookSize,
				nil,
				tc.collateralizationCheckFailures,
				constants.GetStatePosition_ZeroPositionSize,
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
				nil,
				expectedNumCollateralizationChecks,
				tc.expectedRemainingBids,
				tc.expectedRemainingAsks,
				tc.expectedOperations,
				tc.expectedInternalOperations,
				nil,
			)

			// Verify the correct offchain update messages were returned.
			// TODO(DEC-1587): Update the indexer tests to perform assertions on the expected operations queue.
			// assertPlaceOrderOffchainMessages(
			// 	t,
			// 	offchainUpdates,
			// 	tc.order,
			// 	tc.placedMatchableOrders,
			// 	tc.collateralizationCheckFailures,
			// 	nil,
			// 	tc.expectedFilledSize,
			// 	tc.expectedOrderStatus,
			// 	[]expectedMatch{},
			// 	tc.expectedMatches,
			// 	[]types.OrderId{},
			// )
		})
	}
}

func TestAddOrderToOrderbook_PanicsOnInvalidSide(t *testing.T) {
	ctx, _, _ := sdktest.NewSdkContextWithMultistore()
	ctx = ctx.WithIsCheckTx(true)
	memclob := NewMemClobPriceTimePriority(false)

	require.Panics(t, func() {
		memclob.mustAddOrderToOrderbook(
			ctx,
			types.Order{},
			false,
		)
	})
}

func TestAddOrderToOrderbook_ErrorPlaceNewFullyFilledOrder(t *testing.T) {
	ctx, _, _ := sdktest.NewSdkContextWithMultistore()
	ctx = ctx.WithIsCheckTx(true)

	memClobKeeper := mocks.MemClobKeeper{}
	memclob := NewMemClobPriceTimePriority(false)
	memclob.SetClobKeeper(&memClobKeeper)
	memclob.CreateOrderbook(constants.ClobPair_Btc)

	memClobKeeper.On("GetStatePosition", mock.Anything, mock.Anything, mock.Anything).
		Return(big.NewInt(0))
	memClobKeeper.On("ValidateSubaccountEquityTierLimitForNewOrder", mock.Anything, mock.Anything).
		Return(nil)

	order := constants.Order_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB15
	orderId := order.OrderId
	quantums := order.GetBaseQuantums()

	// Set state filled amount as though we learned about this fill
	// from a block, but had not yet learned about the order.
	memClobKeeper.On("GetOrderFillAmount", mock.Anything, orderId).
		Return(true, quantums, uint32(0))

	// Place an order which was already fully-filled in a previous block as though
	// we are only now learning of the order itself via p2p.
	_, _, _, err := memclob.PlaceOrder(ctx, order)

	// Should fail as the order has already been fully filled.
	require.ErrorIs(t, err, types.ErrOrderFullyFilled)
}

func TestAddOrderToOrderbook_PanicsIfFullyFilled(t *testing.T) {
	ctx, _, _ := sdktest.NewSdkContextWithMultistore()
	ctx = ctx.WithIsCheckTx(true)
	memClobKeeper := mocks.MemClobKeeper{}
	memclob := NewMemClobPriceTimePriority(false)
	memclob.SetClobKeeper(&memClobKeeper)
	order := constants.Order_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB15
	orderId := order.OrderId
	quantums := order.GetBaseQuantums()

	memClobKeeper.On("GetStatePosition", mock.Anything, mock.Anything, mock.Anything).
		Return(big.NewInt(0))

	// Fully-filled.
	memClobKeeper.On("GetOrderFillAmount", mock.Anything, orderId).
		Return(true, quantums, uint32(0))

	require.Panics(t, func() {
		memclob.mustAddOrderToOrderbook(ctx, order, false)
	})
}

func TestUpdateOrderbookStateWithMatchedMakerOrder_PanicsOnInvalidFillAmount(t *testing.T) {
	ctx, _, _ := sdktest.NewSdkContextWithMultistore()
	ctx = ctx.WithIsCheckTx(true)
	memclob := NewMemClobPriceTimePriority(false)

	require.Panics(t, func() {
		memclob.mustUpdateOrderbookStateWithMatchedMakerOrder(
			ctx,
			types.Order{Quantums: 1},
		)
	})
}

func TestPlaceOrder_PostOnly(t *testing.T) {
	ctx, _, _ := sdktest.NewSdkContextWithMultistore()
	ctx = ctx.WithIsCheckTx(true)
	tests := map[string]struct {
		// State.
		placedMatchableOrders          []types.MatchableOrder
		collateralizationCheckFailures map[int]map[satypes.SubaccountId]satypes.UpdateResult
		// Parameters.
		order types.Order

		// Expectations.
		expectedErr                error
		expectedOrderStatus        types.OrderStatus
		expectedCollatCheck        []expectedMatch
		expectedRemainingBids      []OrderWithRemainingSize
		expectedRemainingAsks      []OrderWithRemainingSize
		expectedExistingMatches    []expectedMatch
		expectedOperations         []types.Operation
		expectedInternalOperations []types.InternalOperation
	}{
		`Can place a valid post-only order on an empty book`: {
			placedMatchableOrders:          []types.MatchableOrder{},
			collateralizationCheckFailures: map[int]map[satypes.SubaccountId]satypes.UpdateResult{},

			order: constants.Order_Bob_Num0_Id4_Clob1_Buy20_Price35_GTB22_PO,

			expectedOrderStatus: types.Success,
			expectedRemainingBids: []OrderWithRemainingSize{
				{
					Order:         constants.Order_Bob_Num0_Id4_Clob1_Buy20_Price35_GTB22_PO,
					RemainingSize: 20,
				},
			},
			expectedRemainingAsks:      []OrderWithRemainingSize{},
			expectedExistingMatches:    []expectedMatch{},
			expectedOperations:         []types.Operation{},
			expectedInternalOperations: []types.InternalOperation{},
		},
		`A fully matched post-only sell order cannot be placed`: {
			placedMatchableOrders: []types.MatchableOrder{
				&constants.Order_Bob_Num0_Id4_Clob1_Buy20_Price35_GTB22,
				&constants.Order_Bob_Num0_Id11_Clob1_Buy5_Price40_GTB32,
			},

			order: constants.Order_Alice_Num1_Id1_Clob1_Sell10_Price15_GTB20_PO,

			expectedErr: types.ErrPostOnlyWouldCrossMakerOrder,
			expectedRemainingBids: []OrderWithRemainingSize{
				{
					Order:         constants.Order_Bob_Num0_Id4_Clob1_Buy20_Price35_GTB22,
					RemainingSize: 20,
				},
				{
					Order:         constants.Order_Bob_Num0_Id11_Clob1_Buy5_Price40_GTB32,
					RemainingSize: 5,
				},
			},
			expectedRemainingAsks: []OrderWithRemainingSize{},
			// Second order is not collat check'd since the first order generates a valid
			// match, so the matching loop ends.
			expectedCollatCheck: []expectedMatch{
				{
					makerOrder:      &constants.Order_Bob_Num0_Id11_Clob1_Buy5_Price40_GTB32,
					takerOrder:      &constants.Order_Alice_Num1_Id1_Clob1_Sell10_Price15_GTB20_PO,
					matchedQuantums: 5,
				},
			},
			expectedExistingMatches:    []expectedMatch{},
			expectedOperations:         []types.Operation{},
			expectedInternalOperations: []types.InternalOperation{},
		},
		`Placing a post-only order which matches a partially-filled order on the books
					which subsequently fails collateralization
					causes the PO order to be added to the book`: {
			placedMatchableOrders: []types.MatchableOrder{
				// Match #1: This order is partially matched before the test case as a maker order with the below order.
				&constants.Order_Bob_Num0_Id8_Clob0_Sell20_Price10_GTB22,
				// Match #1: This order is fully matched before the test case as a taker order with the above order.
				&constants.Order_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB15,
			},

			order: constants.Order_Alice_Num0_Id1_Clob0_Buy15_Price10_GTB18_PO,

			collateralizationCheckFailures: map[int]map[satypes.SubaccountId]satypes.UpdateResult{
				0: {
					constants.Bob_Num0: satypes.NewlyUndercollateralized,
				},
			},
			expectedExistingMatches: []expectedMatch{
				// Match #1: This match is generated before the test case.
				{
					makerOrder:      &constants.Order_Bob_Num0_Id8_Clob0_Sell20_Price10_GTB22,
					takerOrder:      &constants.Order_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB15,
					matchedQuantums: 5,
				},
			},
			expectedRemainingBids: []OrderWithRemainingSize{
				{
					Order:         constants.Order_Alice_Num0_Id1_Clob0_Buy15_Price10_GTB18_PO,
					RemainingSize: 15,
				},
			},
			expectedRemainingAsks: []OrderWithRemainingSize{},
			expectedCollatCheck: []expectedMatch{
				{
					makerOrder:      &constants.Order_Bob_Num0_Id8_Clob0_Sell20_Price10_GTB22,
					takerOrder:      &constants.Order_Alice_Num0_Id1_Clob0_Buy15_Price10_GTB18_PO,
					matchedQuantums: 15,
				},
			},
			expectedOperations: []types.Operation{
				clobtest.NewOrderPlacementOperation(constants.Order_Bob_Num0_Id8_Clob0_Sell20_Price10_GTB22),
				clobtest.NewOrderPlacementOperation(constants.Order_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB15),
				clobtest.NewMatchOperation(
					&constants.Order_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB15,
					[]types.MakerFill{
						{
							MakerOrderId: constants.Order_Bob_Num0_Id8_Clob0_Sell20_Price10_GTB22.OrderId,
							FillAmount:   5,
						},
					},
				),
				clobtest.NewOrderPlacementOperation(constants.Order_Alice_Num0_Id1_Clob0_Buy15_Price10_GTB18_PO),
			},
			expectedInternalOperations: []types.InternalOperation{
				types.NewShortTermOrderPlacementInternalOperation(constants.Order_Bob_Num0_Id8_Clob0_Sell20_Price10_GTB22),
				types.NewShortTermOrderPlacementInternalOperation(constants.Order_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB15),
				types.NewMatchOrdersInternalOperation(
					constants.Order_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB15,
					[]types.MakerFill{
						{
							MakerOrderId: constants.Order_Bob_Num0_Id8_Clob0_Sell20_Price10_GTB22.OrderId,
							FillAmount:   5,
						},
					},
				),
			},
		},
		`Placing a post-only order which causes a partially-filled maker order to fail collateralization`: {
			placedMatchableOrders: []types.MatchableOrder{
				// Match #1: This order is partially matched before the test case as a maker order with the below order.
				&constants.Order_Bob_Num0_Id8_Clob0_Sell20_Price10_GTB22,
				// Match #1: This order is fully matched before the test case as a taker order with the above order.
				&constants.Order_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB15,
				// Will cause the PO order to cross the book.
				&constants.Order_Carl_Num0_Id2_Clob0_Sell5_Price10_GTB15,
			},

			order: constants.Order_Alice_Num0_Id1_Clob0_Buy15_Price10_GTB18_PO,

			expectedErr: types.ErrPostOnlyWouldCrossMakerOrder,
			collateralizationCheckFailures: map[int]map[satypes.SubaccountId]satypes.UpdateResult{
				0: {
					constants.Bob_Num0: satypes.NewlyUndercollateralized,
				},
			},
			expectedExistingMatches: []expectedMatch{
				// Match #1: This match is generated before the test case.
				{
					makerOrder:      &constants.Order_Bob_Num0_Id8_Clob0_Sell20_Price10_GTB22,
					takerOrder:      &constants.Order_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB15,
					matchedQuantums: 5,
				},
			},
			expectedRemainingBids: []OrderWithRemainingSize{}, // PO order crossed and was canceled.
			expectedRemainingAsks: []OrderWithRemainingSize{
				{
					Order:         constants.Order_Carl_Num0_Id2_Clob0_Sell5_Price10_GTB15,
					RemainingSize: 5,
				},
			},
			expectedCollatCheck: []expectedMatch{
				{
					makerOrder:      &constants.Order_Bob_Num0_Id8_Clob0_Sell20_Price10_GTB22,
					takerOrder:      &constants.Order_Alice_Num0_Id1_Clob0_Buy15_Price10_GTB18_PO,
					matchedQuantums: 15,
				},
				{
					makerOrder:      &constants.Order_Carl_Num0_Id2_Clob0_Sell5_Price10_GTB15,
					takerOrder:      &constants.Order_Alice_Num0_Id1_Clob0_Buy15_Price10_GTB18_PO,
					matchedQuantums: 5,
				},
			},
			expectedOperations: []types.Operation{
				clobtest.NewOrderPlacementOperation(constants.Order_Bob_Num0_Id8_Clob0_Sell20_Price10_GTB22),
				clobtest.NewOrderPlacementOperation(constants.Order_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB15),
				clobtest.NewMatchOperation(
					&constants.Order_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB15,
					[]types.MakerFill{
						{
							MakerOrderId: constants.Order_Bob_Num0_Id8_Clob0_Sell20_Price10_GTB22.OrderId,
							FillAmount:   5,
						},
					},
				),
				clobtest.NewOrderPlacementOperation(constants.Order_Carl_Num0_Id2_Clob0_Sell5_Price10_GTB15),
				clobtest.NewOrderPlacementOperation(constants.Order_Alice_Num0_Id1_Clob0_Buy15_Price10_GTB18_PO),
			},
			expectedInternalOperations: []types.InternalOperation{
				types.NewShortTermOrderPlacementInternalOperation(constants.Order_Bob_Num0_Id8_Clob0_Sell20_Price10_GTB22),
				types.NewShortTermOrderPlacementInternalOperation(constants.Order_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB15),
				types.NewMatchOrdersInternalOperation(
					constants.Order_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB15,
					[]types.MakerFill{
						{
							MakerOrderId: constants.Order_Bob_Num0_Id8_Clob0_Sell20_Price10_GTB22.OrderId,
							FillAmount:   5,
						},
					},
				),
			},
		},
		`A partially matched post-only sell order cannot be placed`: {
			placedMatchableOrders: []types.MatchableOrder{
				&constants.Order_Bob_Num0_Id11_Clob1_Buy5_Price40_GTB32,
			},

			order: constants.Order_Alice_Num1_Id1_Clob1_Sell10_Price15_GTB20_PO,

			expectedErr: types.ErrPostOnlyWouldCrossMakerOrder,
			expectedRemainingBids: []OrderWithRemainingSize{
				{
					Order:         constants.Order_Bob_Num0_Id11_Clob1_Buy5_Price40_GTB32,
					RemainingSize: 5,
				},
			},
			expectedRemainingAsks: []OrderWithRemainingSize{},
			expectedCollatCheck: []expectedMatch{
				{
					makerOrder:      &constants.Order_Bob_Num0_Id11_Clob1_Buy5_Price40_GTB32,
					takerOrder:      &constants.Order_Alice_Num1_Id1_Clob1_Sell10_Price15_GTB20_PO,
					matchedQuantums: 5,
				},
			},
			expectedExistingMatches:    []expectedMatch{},
			expectedOperations:         []types.Operation{},
			expectedInternalOperations: []types.InternalOperation{},
		},
		`A post-only buy order can be placed if all crossing maker orders fail collateralization
					checks, and all crossing maker orders are removed`: {
			placedMatchableOrders: []types.MatchableOrder{
				&constants.Order_Alice_Num0_Id2_Clob1_Sell5_Price10_GTB15,
				&constants.Order_Alice_Num0_Id3_Clob1_Sell5_Price10_GTB15,
				&constants.Order_Alice_Num1_Id1_Clob1_Sell10_Price15_GTB20,
			},
			collateralizationCheckFailures: map[int]map[satypes.SubaccountId]satypes.UpdateResult{
				0: {
					constants.Alice_Num0: satypes.NewlyUndercollateralized,
				},
				1: {
					constants.Alice_Num0: satypes.StillUndercollateralized,
				},
				2: {
					constants.Alice_Num1: satypes.UpdateCausedError,
				},
			},

			order: constants.Order_Bob_Num0_Id4_Clob1_Buy20_Price35_GTB22_PO,

			expectedOrderStatus: types.Success,
			expectedRemainingBids: []OrderWithRemainingSize{
				{
					Order:         constants.Order_Bob_Num0_Id4_Clob1_Buy20_Price35_GTB22_PO,
					RemainingSize: 20,
				},
			},
			expectedRemainingAsks: []OrderWithRemainingSize{},
			expectedCollatCheck: []expectedMatch{
				{
					makerOrder:      &constants.Order_Alice_Num0_Id2_Clob1_Sell5_Price10_GTB15,
					takerOrder:      &constants.Order_Bob_Num0_Id4_Clob1_Buy20_Price35_GTB22_PO,
					matchedQuantums: 5,
				},
				{
					makerOrder:      &constants.Order_Alice_Num0_Id3_Clob1_Sell5_Price10_GTB15,
					takerOrder:      &constants.Order_Bob_Num0_Id4_Clob1_Buy20_Price35_GTB22_PO,
					matchedQuantums: 5,
				},
				{
					makerOrder:      &constants.Order_Alice_Num1_Id1_Clob1_Sell10_Price15_GTB20,
					takerOrder:      &constants.Order_Bob_Num0_Id4_Clob1_Buy20_Price35_GTB22_PO,
					matchedQuantums: 10,
				},
			},
			expectedExistingMatches:    []expectedMatch{},
			expectedOperations:         []types.Operation{},
			expectedInternalOperations: []types.InternalOperation{},
		},
		`A partially matched post-only buy order cannot be placed, and all crossing maker orders that
					failed collateralization checks are removed`: {
			placedMatchableOrders: []types.MatchableOrder{
				&constants.Order_Alice_Num0_Id2_Clob1_Sell5_Price10_GTB15,
				&constants.Order_Alice_Num0_Id3_Clob1_Sell5_Price10_GTB15,
				&constants.Order_Alice_Num1_Id1_Clob1_Sell10_Price15_GTB20,
			},
			collateralizationCheckFailures: map[int]map[satypes.SubaccountId]satypes.UpdateResult{
				0: {
					constants.Alice_Num0: satypes.NewlyUndercollateralized,
				},
				1: {
					constants.Alice_Num0: satypes.StillUndercollateralized,
				},
			},

			order: constants.Order_Bob_Num0_Id4_Clob1_Buy20_Price35_GTB22_PO,

			expectedErr:           types.ErrPostOnlyWouldCrossMakerOrder,
			expectedRemainingBids: []OrderWithRemainingSize{},
			expectedRemainingAsks: []OrderWithRemainingSize{
				{
					Order:         constants.Order_Alice_Num1_Id1_Clob1_Sell10_Price15_GTB20,
					RemainingSize: 10,
				},
			},
			expectedCollatCheck: []expectedMatch{
				{
					makerOrder:      &constants.Order_Alice_Num0_Id2_Clob1_Sell5_Price10_GTB15,
					takerOrder:      &constants.Order_Bob_Num0_Id4_Clob1_Buy20_Price35_GTB22_PO,
					matchedQuantums: 5,
				},
				{
					makerOrder:      &constants.Order_Alice_Num0_Id3_Clob1_Sell5_Price10_GTB15,
					takerOrder:      &constants.Order_Bob_Num0_Id4_Clob1_Buy20_Price35_GTB22_PO,
					matchedQuantums: 5,
				},
				{
					makerOrder:      &constants.Order_Alice_Num1_Id1_Clob1_Sell10_Price15_GTB20,
					takerOrder:      &constants.Order_Bob_Num0_Id4_Clob1_Buy20_Price35_GTB22_PO,
					matchedQuantums: 10,
				},
			},
			expectedExistingMatches:    []expectedMatch{},
			expectedOperations:         []types.Operation{},
			expectedInternalOperations: []types.InternalOperation{},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// Setup memclob state and test expectations.
			addOrderToOrderbookSize := tc.order.GetBaseQuantums()
			if !tc.expectedOrderStatus.IsSuccess() || tc.expectedErr != nil {
				addOrderToOrderbookSize = 0
			}
			expectedFilledSize := satypes.BaseQuantums(0)
			order := tc.order
			memclob, fakeMemClobKeeper, expectedNumCollateralizationChecks, numCollateralChecks := placeOrderTestSetup(
				t,
				ctx,
				tc.placedMatchableOrders,
				&order,
				tc.expectedCollatCheck,
				tc.expectedOrderStatus,
				addOrderToOrderbookSize,
				tc.expectedErr,
				tc.collateralizationCheckFailures,
				constants.GetStatePosition_ZeroPositionSize,
			)

			// Run the test case and verify expectations.
			offchainUpdates := placeOrderAndVerifyExpectationsOperations(
				t,
				ctx,
				memclob,
				order,
				numCollateralChecks,
				expectedFilledSize,
				expectedFilledSize,
				tc.expectedOrderStatus,
				tc.expectedErr,
				expectedNumCollateralizationChecks,
				tc.expectedRemainingBids,
				tc.expectedRemainingAsks,
				tc.expectedOperations,
				tc.expectedInternalOperations,
				fakeMemClobKeeper,
			)

			// Verify the correct offchain update messages were returned.
			assertPlaceOrderOffchainMessages(
				t,
				ctx,
				offchainUpdates,
				order,
				tc.placedMatchableOrders,
				tc.collateralizationCheckFailures,
				tc.expectedErr,
				expectedFilledSize,
				tc.expectedOrderStatus,
				tc.expectedExistingMatches,
				[]expectedMatch{},
				[]types.OrderId{},
				false,
			)
		})
	}
}

func TestPlaceOrder_MatchOrders_WithBuilderCode(t *testing.T) {
	ctx, _, _ := sdktest.NewSdkContextWithMultistore()
	ctx = ctx.WithIsCheckTx(true)
	tests := map[string]struct {
		// State.
		placedMatchableOrders []types.MatchableOrder

		// Parameters.
		order       types.Order
		builderCode *types.BuilderCodeParameters

		// Expectations.
		expectedFilledSize             satypes.BaseQuantums
		expectedOrderStatus            types.OrderStatus
		collateralizationCheckFailures map[int]map[satypes.SubaccountId]satypes.UpdateResult
		expectedErr                    error
		expectedCollatCheck            []expectedMatch
		expectedRemainingBids          []OrderWithRemainingSize
		expectedRemainingAsks          []OrderWithRemainingSize
		expectedOperations             []types.Operation
		expectedInternalOperations     []types.InternalOperation
	}{
		"Order succeeds without builder code": {
			placedMatchableOrders: []types.MatchableOrder{
				&constants.Order_Alice_Num0_Id2_Clob1_Sell5_Price10_GTB15,
			},

			order:       constants.Order_Bob_Num0_Id4_Clob1_Buy20_Price35_GTB22,
			builderCode: nil,

			expectedFilledSize:  5,
			expectedOrderStatus: types.Success,
			expectedRemainingBids: []OrderWithRemainingSize{
				{
					Order:         constants.Order_Bob_Num0_Id4_Clob1_Buy20_Price35_GTB22,
					RemainingSize: 15,
				},
			},
			expectedRemainingAsks: []OrderWithRemainingSize{},
			expectedCollatCheck: []expectedMatch{
				{
					makerOrder:      &constants.Order_Alice_Num0_Id2_Clob1_Sell5_Price10_GTB15,
					takerOrder:      &constants.Order_Bob_Num0_Id4_Clob1_Buy20_Price35_GTB22,
					matchedQuantums: 5,
				},
			},
			expectedOperations: []types.Operation{
				clobtest.NewOrderPlacementOperation(
					constants.Order_Alice_Num0_Id2_Clob1_Sell5_Price10_GTB15,
				),
				clobtest.NewOrderPlacementOperation(
					constants.Order_Bob_Num0_Id4_Clob1_Buy20_Price35_GTB22,
				),
				clobtest.NewMatchOperation(
					&constants.Order_Bob_Num0_Id4_Clob1_Buy20_Price35_GTB22,
					[]types.MakerFill{
						{
							MakerOrderId: constants.Order_Alice_Num0_Id2_Clob1_Sell5_Price10_GTB15.OrderId,
							FillAmount:   5,
						},
					},
				),
			},
			expectedInternalOperations: []types.InternalOperation{
				types.NewShortTermOrderPlacementInternalOperation(
					constants.Order_Alice_Num0_Id2_Clob1_Sell5_Price10_GTB15,
				),
				types.NewShortTermOrderPlacementInternalOperation(
					constants.Order_Bob_Num0_Id4_Clob1_Buy20_Price35_GTB22,
				),
				types.NewMatchOrdersInternalOperation(
					constants.Order_Bob_Num0_Id4_Clob1_Buy20_Price35_GTB22,
					[]types.MakerFill{
						{
							MakerOrderId: constants.Order_Alice_Num0_Id2_Clob1_Sell5_Price10_GTB15.OrderId,
							FillAmount:   5,
						},
					},
				),
			},
		},
		"Order fails with builder code": {
			placedMatchableOrders: []types.MatchableOrder{
				&constants.Order_Alice_Num0_Id2_Clob1_Sell5_Price10_GTB15,
			},

			order: constants.Order_Bob_Num0_Id4_Clob1_Buy20_Price35_GTB22,
			builderCode: &types.BuilderCodeParameters{
				BuilderAddress: constants.Alice_Num0.Owner,
				FeePpm:         10000,
			},

			expectedFilledSize: 0,
			collateralizationCheckFailures: map[int]map[satypes.SubaccountId]satypes.UpdateResult{
				0: {
					constants.Bob_Num0: satypes.NewlyUndercollateralized,
				},
			},
			expectedOrderStatus:   types.Undercollateralized,
			expectedRemainingBids: []OrderWithRemainingSize{},
			expectedRemainingAsks: []OrderWithRemainingSize{
				{
					Order:         constants.Order_Alice_Num0_Id2_Clob1_Sell5_Price10_GTB15,
					RemainingSize: 5,
				},
			},
			expectedCollatCheck: []expectedMatch{
				{
					makerOrder:      &constants.Order_Alice_Num0_Id2_Clob1_Sell5_Price10_GTB15,
					takerOrder:      &constants.Order_Bob_Num0_Id4_Clob1_Buy20_Price35_GTB22,
					matchedQuantums: 5,
				},
			},
			expectedOperations:         []types.Operation{},
			expectedInternalOperations: []types.InternalOperation{},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			order := tc.order
			if tc.builderCode != nil {
				order.BuilderCodeParameters = tc.builderCode
			}

			addOrderToOrderbookSize := satypes.BaseQuantums(0)
			if tc.expectedOrderStatus.IsSuccess() {
				addOrderToOrderbookSize = order.GetBaseQuantums() - tc.expectedFilledSize
			}

			memclob, fakeMemClobKeeper, expectedNumCollateralizationChecks, numCollateralChecks := placeOrderTestSetup(
				t,
				ctx,
				tc.placedMatchableOrders,
				&order,
				tc.expectedCollatCheck,
				tc.expectedOrderStatus,
				addOrderToOrderbookSize,
				nil,
				tc.collateralizationCheckFailures,
				constants.GetStatePosition_ZeroPositionSize,
			)

			placeOrderAndVerifyExpectationsOperations(
				t,
				ctx,
				memclob,
				order,
				numCollateralChecks,
				tc.expectedFilledSize,
				tc.expectedFilledSize,
				tc.expectedOrderStatus,
				nil,
				expectedNumCollateralizationChecks,
				tc.expectedRemainingBids,
				tc.expectedRemainingAsks,
				tc.expectedOperations,
				tc.expectedInternalOperations,
				fakeMemClobKeeper,
			)
		})
	}
}

func TestPlaceOrder_ImmediateOrCancel(t *testing.T) {
	ctx, _, _ := sdktest.NewSdkContextWithMultistore()
	ctx = ctx.WithIsCheckTx(true)
	tests := map[string]struct {
		// State.
		placedMatchableOrders          []types.MatchableOrder
		collateralizationCheckFailures map[int]map[satypes.SubaccountId]satypes.UpdateResult

		// Parameters.
		order types.Order

		// Expectations.
		expectedFilledSize    satypes.BaseQuantums
		expectedOrderStatus   types.OrderStatus
		expectedCollatCheck   []expectedMatch
		expectedRemainingBids []OrderWithRemainingSize
		expectedRemainingAsks []OrderWithRemainingSize
	}{
		`Can place an IOC order on an empty book and it's canceled`: {
			placedMatchableOrders:          []types.MatchableOrder{},
			collateralizationCheckFailures: map[int]map[satypes.SubaccountId]satypes.UpdateResult{},

			order: constants.Order_Alice_Num1_Id1_Clob1_Sell10_Price15_GTB20_IOC,

			expectedFilledSize:    0,
			expectedOrderStatus:   types.ImmediateOrCancelWouldRestOnBook,
			expectedRemainingBids: []OrderWithRemainingSize{},
			expectedRemainingAsks: []OrderWithRemainingSize{},
		},
		`An IOC order can be fully matched`: {
			placedMatchableOrders: []types.MatchableOrder{
				&constants.Order_Bob_Num0_Id4_Clob1_Buy20_Price35_GTB22,
				&constants.Order_Bob_Num0_Id11_Clob1_Buy5_Price40_GTB32,
			},

			order: constants.Order_Alice_Num1_Id1_Clob1_Sell10_Price15_GTB20_IOC,

			expectedFilledSize:  10,
			expectedOrderStatus: types.Success,
			expectedRemainingBids: []OrderWithRemainingSize{
				{
					Order:         constants.Order_Bob_Num0_Id4_Clob1_Buy20_Price35_GTB22,
					RemainingSize: 15,
				},
			},
			expectedRemainingAsks: []OrderWithRemainingSize{},
			expectedCollatCheck: []expectedMatch{
				{
					makerOrder:      &constants.Order_Bob_Num0_Id11_Clob1_Buy5_Price40_GTB32,
					takerOrder:      &constants.Order_Alice_Num1_Id1_Clob1_Sell10_Price15_GTB20_IOC,
					matchedQuantums: 5,
				},
				{
					makerOrder:      &constants.Order_Bob_Num0_Id4_Clob1_Buy20_Price35_GTB22,
					takerOrder:      &constants.Order_Alice_Num1_Id1_Clob1_Sell10_Price15_GTB20_IOC,
					matchedQuantums: 5,
				},
			},
		},
		`An IOC order can be partially matched and the remaining unmatched size is canceled`: {
			placedMatchableOrders: []types.MatchableOrder{
				&constants.Order_Alice_Num0_Id4_Clob1_Buy25_Price5_GTB20,
				&constants.Order_Bob_Num0_Id11_Clob1_Buy5_Price40_GTB32,
			},

			order: constants.Order_Alice_Num1_Id1_Clob1_Sell10_Price15_GTB20_IOC,

			expectedFilledSize:  5,
			expectedOrderStatus: types.ImmediateOrCancelWouldRestOnBook,
			expectedRemainingBids: []OrderWithRemainingSize{
				{
					Order:         constants.Order_Alice_Num0_Id4_Clob1_Buy25_Price5_GTB20,
					RemainingSize: 25,
				},
			},
			expectedRemainingAsks: []OrderWithRemainingSize{},
			expectedCollatCheck: []expectedMatch{
				{
					makerOrder:      &constants.Order_Bob_Num0_Id11_Clob1_Buy5_Price40_GTB32,
					takerOrder:      &constants.Order_Alice_Num1_Id1_Clob1_Sell10_Price15_GTB20_IOC,
					matchedQuantums: 5,
				},
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// Setup memclob state and test expectations.
			addOrderToOrderbookSize := satypes.BaseQuantums(0)
			order := tc.order
			require.Equal(t, types.Order_TIME_IN_FORCE_IOC, order.TimeInForce)
			memclob, _, expectedNumCollateralizationChecks, numCollateralChecks := placeOrderTestSetup(
				t,
				ctx,
				tc.placedMatchableOrders,
				&order,
				tc.expectedCollatCheck,
				tc.expectedOrderStatus,
				addOrderToOrderbookSize,
				nil,
				tc.collateralizationCheckFailures,
				constants.GetStatePosition_ZeroPositionSize,
			)

			// Run the test case and verify expectations.
			offchainUpdates := placeOrderAndVerifyExpectations(
				t,
				ctx,
				memclob,
				order,
				numCollateralChecks,
				tc.expectedFilledSize,
				tc.expectedFilledSize,
				tc.expectedOrderStatus,
				nil,
				expectedNumCollateralizationChecks,
				tc.expectedRemainingBids,
				tc.expectedRemainingAsks,
				tc.expectedCollatCheck,
				nil,
			)

			// Verify the correct offchain update messages were returned.
			assertPlaceOrderOffchainMessages(
				t,
				ctx,
				offchainUpdates,
				order,
				tc.placedMatchableOrders,
				tc.collateralizationCheckFailures,
				nil,
				tc.expectedFilledSize,
				tc.expectedOrderStatus,
				[]expectedMatch{},
				tc.expectedCollatCheck,
				[]types.OrderId{},
				false,
			)
		})
	}
}

func TestPlaceOrder_Telemetry(t *testing.T) {
	m, err := telemetry.New(telemetry.Config{
		Enabled:        true,
		EnableHostname: false,
		ServiceName:    "test",
	})
	require.NoError(t, err)
	require.NotNil(t, m)

	ctx, _, _ := sdktest.NewSdkContextWithMultistore()
	ctx = ctx.WithIsCheckTx(true)

	// Setup the memclob state.
	memClobKeeper := testutil_memclob.NewFakeMemClobKeeper()
	memclob := NewMemClobPriceTimePriority(false)
	memclob.SetClobKeeper(memClobKeeper)

	clobPairId := uint32(0)

	orders := make([]types.Order, 0, 5)
	for i := 0; i < 5; i++ {
		order := types.Order{
			OrderId: types.OrderId{
				SubaccountId: constants.Alice_Num0,
				ClientId:     uint32(i),
				ClobPairId:   clobPairId,
			},
			Side:         types.Order_SIDE_BUY,
			Quantums:     100,
			Subticks:     5,
			GoodTilOneof: &types.Order_GoodTilBlock{GoodTilBlock: 5},
		}
		orders = append(orders, order)
	}

	// Create the orderbook.
	memclob.CreateOrderbook(constants.ClobPair_Btc)

	// Create all orders.
	createAllOrders(
		t,
		ctx,
		memclob,
		orders,
	)

	gr, err := m.Gather(telemetry.FormatText)
	require.NoError(t, err)
	require.Equal(t, "application/json", gr.ContentType)

	jsonMetrics := make(map[string]interface{})
	require.NoError(t, json.Unmarshal(gr.Metrics, &jsonMetrics))

	counters := jsonMetrics["Counters"].([]any)
	require.Condition(t, func() bool {
		for _, counter := range counters {
			if counter.(map[string]any)["Name"].(string) == "test.clob.place_order.added_to_orderbook" &&
				counter.(map[string]any)["Count"].(float64) == 5.0 {
				return true
			}
		}
		return false
	})

	samples := jsonMetrics["Samples"].([]interface{})
	require.Condition(t, func() bool {
		for _, sample := range samples {
			if sample.(map[string]any)["Name"].(string) == "test.place_order.memclob.added_to_orderbook.latency" &&
				sample.(map[string]any)["Count"].(float64) == 5.0 {
				return true
			}
		}
		return false
	})
}

func TestPlaceOrder_GenerateOffchainUpdatesFalse_NoMessagesSent(t *testing.T) {
	ctx, _, _ := sdktest.NewSdkContextWithMultistore()
	ctx = ctx.WithIsCheckTx(true)
	// Setup the memclob state.
	memClobKeeper := testutil_memclob.NewFakeMemClobKeeper()
	memclob := NewMemClobPriceTimePriority(false)
	memclob.SetClobKeeper(memClobKeeper)

	order := constants.Order_Bob_Num0_Id13_Clob0_Sell35_Price35_GTB30

	// Create the orderbook.
	memclob.CreateOrderbook(constants.ClobPair_Btc)

	// Place a new order.
	_, _, offchainUpdates, err := memclob.PlaceOrder(ctx, order)
	require.NoError(t, err)
	require.Empty(t, offchainUpdates.GetMessages())
}

// TestPlaceOrder_DuplicateOrder tests that placing the same order twice returns an ErrInvalidReplacement
// error. Adding this test because we depend on this being the case in PrepareCheckState. There are certain
// situations in which PrepareCheckState may attempt to place the same order twice, and we want to make sure
// that the second call to PlaceOrder will return this error instead of undergoing any placement/matching logic.
func TestPlaceOrder_DuplicateOrder(t *testing.T) {
	ctx, _, _ := sdktest.NewSdkContextWithMultistore()
	ctx = ctx.WithIsCheckTx(true)

	memClobKeeper := testutil_memclob.NewFakeMemClobKeeper()
	memclob := NewMemClobPriceTimePriority(false)
	memclob.SetClobKeeper(memClobKeeper)

	memclob.CreateOrderbook(constants.ClobPair_Btc)

	order := constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy100_Price10_GTBT15
	_, _, _, err := memclob.PlaceOrder(ctx, order)
	require.NoError(t, err)
	_, _, _, err = memclob.PlaceOrder(ctx, order)
	require.Error(t, err, types.ErrInvalidReplacement)
}
