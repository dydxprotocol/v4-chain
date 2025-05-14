package memclob

import (
	"errors"
	"fmt"
	"math"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	mapset "github.com/deckarep/golang-set/v2"
	"github.com/dydxprotocol/v4-chain/protocol/indexer/msgsender"
	"github.com/dydxprotocol/v4-chain/protocol/indexer/off_chain_updates"
	ocutypes "github.com/dydxprotocol/v4-chain/protocol/indexer/off_chain_updates/types"
	indexershared "github.com/dydxprotocol/v4-chain/protocol/indexer/shared/types"
	"github.com/dydxprotocol/v4-chain/protocol/lib"
	testutil_memclob "github.com/dydxprotocol/v4-chain/protocol/testutil/memclob"
	"github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
	satypes "github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"
	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/require"
)

// expectedMatch is a testing utility struct used for verifying a match occurred between a maker and taker order, with
// a specific amount of quantums.
type expectedMatch struct {
	makerOrder      types.MatchableOrder
	takerOrder      types.MatchableOrder
	matchedQuantums satypes.BaseQuantums
}

// OrderWithRemainingSize is a testing utility struct used for storing an order along with the expected remaining size.
type OrderWithRemainingSize struct {
	Order         types.Order
	RemainingSize satypes.BaseQuantums
}

// createCollatCheckExpectationsFromPendingMatches is a testing utility function used for creating the expected
// collateralization check parameters.
// This function iterates through `expectedPendingMatches`, and constructs the expected collateralization checks from
// the matches. If any previous collateralization checks succeeded, it will add all matches for both maker and
// taker subaccounts that previously passed collateralization checks to the expected collateralization check.
// For handling the case where the taker order is placed on the book, this function additionally takes parameters for an
// order to add to the orderbook along with a size parameter, and if the size is greater than zero it will expect a
// collateralization check when adding the order to the orderbook.
// Note that this function also adds in matches currently within the match queue to the expected pending matches,
// for each subaccount, within each expected collateralization check.
func createCollatCheckExpectationsFromPendingMatches(
	ctx sdk.Context,
	t *testing.T,
	memclob *MemClobPriceTimePriority,
	expectedPendingMatches []expectedMatch,
	expectedCollatCheckFailures map[int]bool,
	order types.MatchableOrder,
	addToOrderbookSize satypes.BaseQuantums,
) (
	expectedCollatChecks map[int]map[satypes.SubaccountId][]types.PendingOpenOrder,
) {
	expectedCollatChecks = make(map[int]map[satypes.SubaccountId][]types.PendingOpenOrder)
	clobPairId := order.GetClobPairId()

	for i, expectedMatch := range expectedPendingMatches {
		// Ensure the maker and taker orders are on the same CLOB, opposite sides, and
		// not the same subaccount.
		require.Equal(
			t,
			expectedMatch.takerOrder.GetClobPairId(),
			expectedMatch.makerOrder.GetClobPairId(),
		)
		require.NotEqual(
			t,
			expectedMatch.takerOrder.IsBuy(),
			expectedMatch.makerOrder.IsBuy(),
		)
		require.NotEqual(
			t,
			expectedMatch.takerOrder.GetSubaccountId(),
			expectedMatch.makerOrder.GetSubaccountId(),
		)

		makerSubaccountId := expectedMatch.makerOrder.GetSubaccountId()
		takerSubaccountId := expectedMatch.takerOrder.GetSubaccountId()
		subticks := expectedMatch.makerOrder.GetOrderSubticks()
		expectedPendingMatchesForCollatCheck := map[satypes.SubaccountId][]types.PendingOpenOrder{
			makerSubaccountId: {
				{
					RemainingQuantums: expectedMatch.matchedQuantums,
					IsBuy:             expectedMatch.makerOrder.IsBuy(),
					Subticks:          subticks,
					ClobPairId:        clobPairId,
				},
			},
			takerSubaccountId: {
				{
					RemainingQuantums: expectedMatch.matchedQuantums,
					IsBuy:             expectedMatch.takerOrder.IsBuy(),
					Subticks:          subticks,
					ClobPairId:        clobPairId,
				},
			},
		}

		expectedCollatChecks[i] = expectedPendingMatchesForCollatCheck
	}

	return expectedCollatChecks
}

// createMatchExpectationsFromMatches is a testing utility function for calculating the expected values of memclob
// fields for tracking matches, based off a list of expected matches.
func createMatchExpectationsFromMatches(
	expectedMatches []expectedMatch,
) (
	expectedOrderIdToFilledAmount map[types.OrderId]satypes.BaseQuantums,
) {
	expectedOrderIdToFilledAmount = make(map[types.OrderId]satypes.BaseQuantums)

	for _, em := range expectedMatches {
		// Update the filled size of the maker order.
		expectedOrderIdToFilledAmount[em.makerOrder.MustGetOrder().OrderId] += em.matchedQuantums
		// If it's not a liquidation order, update the filled size of the taker order.
		if !em.takerOrder.IsLiquidation() {
			expectedOrderIdToFilledAmount[em.takerOrder.MustGetOrder().OrderId] += em.matchedQuantums
		}
	}

	return expectedOrderIdToFilledAmount
}

// createMatchExpectationsFromOperations is a testing utility function for calculating the expected
// values of memclob fields for tracking the operations queue, based off a list of expected operations.
func createMatchExpectationsFromOperations(
	expectedOperations []types.Operation,
) (
	expectedOrderIdToFilledAmount map[types.OrderId]satypes.BaseQuantums,
	orderHashToOperationsQueueOrder map[types.OrderHash]types.Order,
) {
	expectedOrderIdToFilledAmount = make(map[types.OrderId]satypes.BaseQuantums)
	orderHashToOperationsQueueOrder = make(map[types.OrderHash]types.Order)

	for _, op := range expectedOperations {
		switch operation := op.Operation.(type) {
		case *types.Operation_ShortTermOrderPlacement:
			order := operation.ShortTermOrderPlacement.Order
			orderHashToOperationsQueueOrder[order.GetOrderHash()] = order
		case *types.Operation_Match:
			switch match := operation.Match.Match.(type) {
			case *types.ClobMatch_MatchOrders:
				// For each fill, add the fill amount to the maker and taker order's filled amount.
				for _, fill := range match.MatchOrders.Fills {
					expectedOrderIdToFilledAmount[fill.MakerOrderId] += satypes.BaseQuantums(fill.FillAmount)
					expectedOrderIdToFilledAmount[match.MatchOrders.TakerOrderId] += satypes.BaseQuantums(fill.FillAmount)
				}
			case *types.ClobMatch_MatchPerpetualLiquidation:
				// For each fill, add the fill amount to the maker order's filled amount.
				// Note we skip the taker order because it's a liquidation order.
				for _, fill := range match.MatchPerpetualLiquidation.Fills {
					expectedOrderIdToFilledAmount[fill.MakerOrderId] += satypes.BaseQuantums(fill.FillAmount)
				}
			default:
				panic(
					fmt.Sprintf(
						"Unknown match operation type %+v",
						operation,
					),
				)
			}
		case *types.Operation_ShortTermOrderCancellation:
			// Do nothing since cancels are not tracked in any operations queue data structures.
		default:
			panic(
				fmt.Sprintf(
					"Unsupported operation type %+v",
					operation,
				),
			)
		}
	}

	return expectedOrderIdToFilledAmount, orderHashToOperationsQueueOrder
}

// assertMemclobHasMatches asserts that the memclob contains each passed in bid and ask order on the respective side,
// along with the state from these bid and ask orders (such as `BestBid`, `BestAsk`, and the total quantums at
// each level).
func assertMemclobHasMatches(
	t *testing.T,
	ctx sdk.Context,
	memclob *MemClobPriceTimePriority,
	expectedMatches []expectedMatch,
) {
	expectedOrderIdToFilledAmount := createMatchExpectationsFromMatches(expectedMatches)

	for orderId, amount := range expectedOrderIdToFilledAmount {
		require.Equal(
			t,
			amount,
			memclob.GetOrderFilledAmount(ctx, orderId),
		)
	}
}

// assertMemclobHasOperations asserts that the memclob contains each passed in bid and ask order on the respective side,
// along with the state from these bid and ask orders (such as `BestBid`, `BestAsk`, and the total quantums at
// each level).
func assertMemclobHasOperations(
	t *testing.T,
	ctx sdk.Context,
	memclob *MemClobPriceTimePriority,
	expectedOperations []types.Operation,
	expectedInternalOperations []types.InternalOperation,
) {
	expectedOrderIdToFilledAmount, _ := createMatchExpectationsFromOperations(expectedOperations)

	for orderId, amount := range expectedOrderIdToFilledAmount {
		require.Equal(
			t,
			amount,
			memclob.GetOrderFilledAmount(ctx, orderId),
		)
	}

	require.Equal(
		t,
		expectedInternalOperations,
		memclob.operationsToPropose.OperationsQueue,
	)
}

// AssertMemclobHasShortTermTxBytes asserts that the memclob contains a TX bytes entry for each
// Short-Term order on the orderbook and in the operations queue.
func AssertMemclobHasShortTermTxBytes(
	t *testing.T,
	ctx sdk.Context,
	memclob *MemClobPriceTimePriority,
	expectedInternalOperations []types.InternalOperation,
	expectedRemainingBids []OrderWithRemainingSize,
	expectedRemainingAsks []OrderWithRemainingSize,
) {
	expectedShortTermOrderHashes := make(map[types.OrderHash]bool)
	for _, operation := range expectedInternalOperations {
		if shortTermOrder := operation.GetShortTermOrderPlacement(); shortTermOrder != nil {
			expectedShortTermOrderHashes[shortTermOrder.Order.GetOrderHash()] = true
		}
	}

	for _, orders := range [][]OrderWithRemainingSize{expectedRemainingBids, expectedRemainingAsks} {
		for _, order := range orders {
			if order.Order.IsShortTermOrder() {
				expectedShortTermOrderHashes[order.Order.GetOrderHash()] = true
			}
		}
	}

	expectedShortTermOrdersWithTxBytes := lib.GetSortedKeys[types.SortedOrderHashes](
		expectedShortTermOrderHashes,
	)
	shortTermOrdersWithTxBytes := lib.GetSortedKeys[types.SortedOrderHashes](
		memclob.operationsToPropose.ShortTermOrderHashToTxBytes,
	)

	// TODO(CLOB-631): Add test coverage to verify values of `ShortTermOrderHashToTxBytes`.
	require.ElementsMatch(
		t,
		expectedShortTermOrdersWithTxBytes,
		shortTermOrdersWithTxBytes,
	)
}

// createOrderbookExpectationsFromOrders is a testing utility function for calculating the expected bookkeeping
// variables of a side of the orderbook, given an order CLOB pair ID, side, and orders. It also performs validation on
// the orders to verify they match the passed in parameters.
// It returns the expected bookkeeping variables for that orderbook and side.
func createOrderbookExpectationsFromOrders(
	t *testing.T,
	expectedClobPairId types.ClobPairId,
	expectedIsBuy bool,
	orders []OrderWithRemainingSize,
) (
	expectedBestOrder types.Subticks,
	numLevels int,
) {
	if expectedIsBuy {
		expectedBestOrder = 0
	} else {
		expectedBestOrder = math.MaxUint64
	}

	seenLevels := make(map[types.Subticks]bool)

	for _, order := range orders {
		// Verify all orders have the expected CLOB pair ID.
		if order.Order.GetClobPairId() != expectedClobPairId {
			require.Fail(
				t,
				fmt.Sprintf(
					"Bid with order ID %s has CLOB pair ID %d, expected %d",
					order.Order.OrderId.String(),
					order.Order.GetClobPairId(),
					expectedClobPairId,
				),
			)
		}

		// Verify all orders have the correct side.
		if order.Order.IsBuy() != expectedIsBuy {
			require.Fail(
				t,
				fmt.Sprintf(
					"Order with order ID %s has side %s, expected %s",
					order.Order.OrderId.String(),
					testutil_memclob.OrderSideHumanReadable(order.Order.IsBuy()),
					testutil_memclob.OrderSideHumanReadable(expectedIsBuy),
				),
			)
		}

		// Update the best seen order if necessary.
		subticks := order.Order.GetOrderSubticks()
		if expectedIsBuy && subticks > expectedBestOrder {
			expectedBestOrder = subticks
		} else if !expectedIsBuy && subticks < expectedBestOrder {
			expectedBestOrder = subticks
		}

		seenLevels[subticks] = true
	}

	return expectedBestOrder, len(seenLevels)
}

// AssertMemclobHasOrders asserts that the memclob contains each passed in bid and ask order on the
// respective side and orderbook. It also verifies the expected state of each orderbook that holds
// at least one order, specifically the `BestBid` and `BestAsk`.
func AssertMemclobHasOrders(
	t *testing.T,
	ctx sdk.Context,
	memclob *MemClobPriceTimePriority,
	expectedBids []OrderWithRemainingSize,
	expectedAsks []OrderWithRemainingSize,
) {
	allClobPairIds := make(map[types.ClobPairId]bool)
	clobPairToExpectedBids := make(map[types.ClobPairId][]OrderWithRemainingSize)
	for _, bid := range expectedBids {
		clobPairId := types.ClobPairId(bid.Order.OrderId.ClobPairId)
		if _, exists := clobPairToExpectedBids[clobPairId]; !exists {
			clobPairToExpectedBids[clobPairId] = make([]OrderWithRemainingSize, 0)
		}
		clobPairToExpectedBids[clobPairId] = append(clobPairToExpectedBids[clobPairId], bid)
		allClobPairIds[clobPairId] = true
	}
	clobPairToExpectedAsks := make(map[types.ClobPairId][]OrderWithRemainingSize)
	for _, ask := range expectedAsks {
		clobPairId := types.ClobPairId(ask.Order.OrderId.ClobPairId)
		if _, exists := clobPairToExpectedAsks[clobPairId]; !exists {
			clobPairToExpectedAsks[clobPairId] = make([]OrderWithRemainingSize, 0)
		}
		clobPairToExpectedAsks[clobPairId] = append(clobPairToExpectedAsks[clobPairId], ask)
		allClobPairIds[clobPairId] = true
	}

	// Get all orders and CLOB pair IDs from the memclob.
	unseenOrdersOnMemclob := make(map[types.OrderId]types.Order)
	for clobPairId, orderbook := range memclob.orderbooks {
		allClobPairIds[clobPairId] = true
		for _, levelOrder := range orderbook.Bids {
			curr := levelOrder.LevelOrders.Front
			for curr != nil {
				unseenOrdersOnMemclob[curr.Value.Order.OrderId] = curr.Value.Order
				curr = curr.Next
			}
		}
		for _, levelOrder := range orderbook.Asks {
			curr := levelOrder.LevelOrders.Front
			for curr != nil {
				unseenOrdersOnMemclob[curr.Value.Order.OrderId] = curr.Value.Order
				curr = curr.Next
			}
		}
	}

	// Perform assertions on each seen CLOB on the provided list of expected bids and asks.
	for clobPairId := range allClobPairIds {
		// Test setup.
		expectedBids := make([]OrderWithRemainingSize, 0)
		if _, exists := clobPairToExpectedBids[clobPairId]; exists {
			expectedBids = clobPairToExpectedBids[clobPairId]
		}
		expectedBestBid, numBidLevels := createOrderbookExpectationsFromOrders(
			t,
			clobPairId,
			true,
			expectedBids,
		)

		expectedAsks := make([]OrderWithRemainingSize, 0)
		if _, exists := clobPairToExpectedAsks[clobPairId]; exists {
			expectedAsks = clobPairToExpectedAsks[clobPairId]
		}
		expectedBestAsk, numAskLevels := createOrderbookExpectationsFromOrders(
			t,
			clobPairId,
			false,
			expectedAsks,
		)

		// Verify memclob state matches the expected state.
		orderbook, exists := memclob.orderbooks[clobPairId]
		require.True(t, exists)

		// Verify best ask and best bid.
		require.Equal(
			t,
			expectedBestAsk,
			orderbook.BestAsk,
			fmt.Sprintf("Best ask for CLOB pair ID %d is incorrect", clobPairId),
		)
		require.Equal(
			t,
			expectedBestBid,
			orderbook.BestBid,
			fmt.Sprintf("Best bid for CLOB pair ID %d is incorrect", clobPairId),
		)

		// Verify total number of levels.
		require.Len(t, orderbook.Bids, numBidLevels)
		require.Len(t, orderbook.Asks, numAskLevels)

		// Verify each order has the correct remaining size and an assigned nonce.
		expectedRestingOrders := append(expectedAsks, expectedBids...)
		require.Equal(t, len(expectedRestingOrders), int(orderbook.TotalOpenOrders))
		for _, order := range expectedRestingOrders {
			orderId := order.Order.OrderId

			// Verify the order has nonzero remaining size.
			quantums := order.RemainingSize
			if quantums == 0 {
				require.Fail(t, fmt.Sprintf("Bid with order ID %s has 0 remaining quantums", orderId.String()))
			}

			remainingAmount, hasRemainingAmount := memclob.GetOrderRemainingAmount(
				ctx,
				order.Order,
			)
			require.True(t, hasRemainingAmount)
			require.Equal(t, order.RemainingSize, remainingAmount)

			// Mark the order as seen.
			delete(unseenOrdersOnMemclob, orderId)
		}
	}

	// Verify we have seen all orders that are currently on the memclob.
	require.Empty(t, unseenOrdersOnMemclob)
}

// assertOrderbookStateExpectations enforces various expectations around the in-memory data structures for
// orders in the clob.
func assertOrderbookStateExpectations(
	t *testing.T,
	memclob *MemClobPriceTimePriority,
	order types.Order,
	expectedBestBid types.Subticks,
	expectedBestAsk types.Subticks,
	expectedTotalLevels int,
	expectLevelToExist bool,
	expectBlockExpirationsForOrdersToExist bool,
	expectSubaccountOpenClobOrdersForSideToExist bool,
	expectSubaccountOpenClobOrdersToExist bool,
) {
	orderbook := memclob.orderbooks[order.GetClobPairId()]
	// Expect that the `blockExpirationsForOrders` map exists for the `GoodTilBlock` associated with the
	// passed in order.
	if expectBlockExpirationsForOrdersToExist {
		require.NotEmpty(t, orderbook.blockExpirationsForOrders[order.GetGoodTilBlock()])
	} else {
		require.Empty(t, orderbook.blockExpirationsForOrders[order.GetGoodTilBlock()])
	}

	// Expect that the relevant `SubaccountOpenClobOrders` map exists for the `SubaccountId` and `ClobPairId`
	// associated with the passed in order.
	if expectSubaccountOpenClobOrdersToExist {
		require.NotEmpty(
			t,
			memclob.orderbooks[order.GetClobPairId()].SubaccountOpenClobOrders[order.OrderId.SubaccountId],
		)
	} else {
		require.Empty(
			t,
			memclob.orderbooks[order.GetClobPairId()].SubaccountOpenClobOrders[order.OrderId.SubaccountId],
		)
	}

	// Expect that the relevant `subaccountOpenClobOrders` map exists for the `SubaccountId`, `ClobPairId`, and `Side`
	// associated with the passed in order.
	if expectSubaccountOpenClobOrdersForSideToExist {
		require.NotEmpty(
			t,
			memclob.orderbooks[order.GetClobPairId()].
				SubaccountOpenClobOrders[order.OrderId.SubaccountId][order.Side],
		)
	} else {
		require.Empty(
			t,
			memclob.orderbooks[order.GetClobPairId()].
				SubaccountOpenClobOrders[order.OrderId.SubaccountId][order.Side],
		)
	}

	// Verify the `BestBid` for this orderbook associated with the passed in order matches expectedBestBid.
	require.Equal(t, expectedBestBid, orderbook.BestBid)
	// Verify the `BestBid` for this orderbook associated with the passed in order matches expectedBestAsk.
	require.Equal(t, expectedBestAsk, orderbook.BestAsk)

	isBuy := order.IsBuy()

	// Verify the price level's total number of quantums was updated properly.
	levels := orderbook.GetSide(isBuy)

	require.Len(t, levels, expectedTotalLevels)
	_, exists := levels[order.GetOrderSubticks()]
	require.Equal(t, expectLevelToExist, exists)
}

// createOrderbooks creates orderbooks up to the provided `maxId`.
func createOrderbooks(
	t *testing.T,
	ctx sdk.Context,
	memclob *MemClobPriceTimePriority,
	maxOrderbooks uint,
) {
	// Create all unique orderbooks.
	for clobPairId := uint32(0); clobPairId < uint32(maxOrderbooks); clobPairId++ {
		clobPair := types.ClobPair{
			Id:               clobPairId,
			SubticksPerTick:  5,
			StepBaseQuantums: 5,
			Metadata: &types.ClobPair_PerpetualClobMetadata{
				PerpetualClobMetadata: &types.PerpetualClobMetadata{
					// Set the `PerpetualId` field to be the same as the CLOB pair ID.
					PerpetualId: clobPairId,
				},
			},
		}
		memclob.CreateOrderbook(clobPair)
	}
}

// createAllOrderbooksForMatchableOrders creates relevant orderbooks in the memclob for each
// matchable order if the orderbook does not already exist. The only difference between each created
// CLOB pair is the `Id` field, all other fields are the same.
func createAllOrderbooksForMatchableOrders(
	t *testing.T,
	ctx sdk.Context,
	memclob *MemClobPriceTimePriority,
	orders []types.MatchableOrder,
) {
	// Create all unique orderbooks.
	createdOrderbooks := make(map[types.ClobPairId]bool)
	for _, order := range orders {
		// Note the for-loop is necessary due to the auto-incrementing ID after creating CLOB pairs.
		if _, exists := createdOrderbooks[order.GetClobPairId()]; !exists {
			clobPair := types.ClobPair{
				Id:               order.GetClobPairId().ToUint32(),
				SubticksPerTick:  5,
				StepBaseQuantums: 5,
				Metadata: &types.ClobPair_PerpetualClobMetadata{
					PerpetualClobMetadata: &types.PerpetualClobMetadata{
						PerpetualId: 0,
					},
				},
			}
			memclob.CreateOrderbook(clobPair)
			createdOrderbooks[order.GetClobPairId()] = true
		}
	}
}

// createAllOrderbooksForOrders creates relevant orderbooks in the memclob for each order if the orderbook
// does not already exist.
func createAllOrderbooksForOrders(
	t *testing.T,
	ctx sdk.Context,
	memclob *MemClobPriceTimePriority,
	orders []types.Order,
) {
	// Create all unique orderbooks.
	createdOrderbooks := make(map[types.ClobPairId]bool)
	for _, order := range orders {
		if _, exists := createdOrderbooks[order.GetClobPairId()]; !exists {
			clobPair := types.ClobPair{
				Id:               order.GetClobPairId().ToUint32(),
				SubticksPerTick:  5,
				StepBaseQuantums: 1,
				Metadata: &types.ClobPair_PerpetualClobMetadata{
					PerpetualClobMetadata: &types.PerpetualClobMetadata{
						PerpetualId: 0,
					},
				},
			}
			memclob.CreateOrderbook(clobPair)
			createdOrderbooks[order.GetClobPairId()] = true
		}
	}
}

// applyOperationsToMemclob applies all operations to the provided memclob.
// It currently only supports order placements and order cancellations. Note that the orderbook
// referenced by each operation should already exist.
// TODO(DEC-1712): Change how memclob test setup is performed.
func applyOperationsToMemclob(
	t *testing.T,
	ctx sdk.Context,
	memclob *MemClobPriceTimePriority,
	operations []types.Operation,
	memClobKeeper types.MemClobKeeper,
) {
	// Perform all existing operations on the memclob.
	for _, op := range operations {
		switch op.Operation.(type) {
		case *types.Operation_ShortTermOrderPlacement:
			orderPlacement := op.GetShortTermOrderPlacement()
			_, _, _, err := memclob.PlaceOrder(ctx, orderPlacement.Order)
			require.NoError(t, err)
		case *types.Operation_ShortTermOrderCancellation:
			orderCancellation := op.GetShortTermOrderCancellation()
			_, err := memclob.CancelOrder(ctx, orderCancellation)
			require.NoError(t, err)
		case *types.Operation_PreexistingStatefulOrder:
			preexistingStatefulOrderId := op.GetPreexistingStatefulOrder()
			orderPlacement, found := memClobKeeper.GetLongTermOrderPlacement(
				ctx,
				*preexistingStatefulOrderId,
			)
			require.True(t, found)
			_, _, _, err := memclob.PlaceOrder(ctx, orderPlacement.Order)
			require.NoError(t, err)
		default:
			panic(
				fmt.Sprintf(
					"applyOperationsToMemclob: cannot apply operation %+v to memclob",
					op,
				),
			)
		}
	}
}

// createAllMatchableOrders creates matchable orders in the memclob. Note that it expects the
// orderbook referenced by the order to already exist.
func createAllMatchableOrders(
	t *testing.T,
	ctx sdk.Context,
	memclob *MemClobPriceTimePriority,
	matchableOrders []types.MatchableOrder,
	fakeMemClobKeeper *testutil_memclob.FakeMemClobKeeper,
) {
	// Place all existing orders on the orderbook.
	for _, matchableOrder := range matchableOrders {
		// If the order is a liquidation order, place the liquidation.
		// Else, assume it's a regular order and place it.
		if matchableOrder.IsLiquidation() {
			clobPair := types.ClobPair{
				Id: matchableOrder.GetClobPairId().ToUint32(),
				Metadata: &types.ClobPair_PerpetualClobMetadata{
					PerpetualClobMetadata: &types.PerpetualClobMetadata{
						PerpetualId: matchableOrder.MustGetLiquidatedPerpetualId(),
					},
				},
			}
			liquidationOrder := types.NewLiquidationOrder(
				matchableOrder.GetSubaccountId(),
				clobPair,
				matchableOrder.IsBuy(),
				matchableOrder.GetBaseQuantums(),
				matchableOrder.GetOrderSubticks(),
			)

			_, _, _, err := memclob.PlacePerpetualLiquidation(
				ctx,
				*liquidationOrder,
			)
			require.NoError(t, err)

			fakeMemClobKeeper.CommitState()
		} else {
			order := matchableOrder.MustGetOrder()
			_, _, _, err := memclob.PlaceOrder(
				ctx,
				order,
			)
			require.NoError(t, err)

			fakeMemClobKeeper.CommitState()
		}
	}
}

// createAllOrders creates orders in the memclob. Note that it expects the orderbook referenced by the order to
// already exist.
func createAllOrders(
	t *testing.T,
	ctx sdk.Context,
	memclob *MemClobPriceTimePriority,
	orders []types.Order,
) {
	// Place all existing orders on the orderbook.
	for _, order := range orders {
		_, _, _, err := memclob.PlaceOrder(
			ctx,
			order,
		)
		require.NoError(t, err)
	}
}

// doesOrderExistOnSide is a testing helper used for checking whether `orderId` exists in `orderLevels`.
func doesOrderExistOnSide(
	t *testing.T,
	orderId types.OrderId,
	orderLevels map[types.Subticks]*types.Level,
) bool {
	foundOrderId := false
	for _, level := range orderLevels {
		// Starting from the first element in the list of orders in this level, check if `orderId` exists in the level.
		// Note that we can assume `level.LevelOrders.Front` exists, since levels should only exist if they have more than
		// zero orders.
		level.LevelOrders.Front.Each(func(levelOrder types.ClobOrder) {
			if !foundOrderId && levelOrder.Order.OrderId == orderId {
				foundOrderId = true
			}
		})

		if foundOrderId {
			break
		}
	}

	return foundOrderId
}

// requireOrderExistsInMemclob is a testing helper used for verifying that `order.OrderId` exists in all expected
// memclob data structures. As a caveat, this helper function will only verify the order is contained
// within the orderbook specified by `order.ClobPairId`. It will not check if the order exists somewhere where it
// _should not_ (i.e. it will not check all orderbooks in the memclob).
// Note that this function must be defined within this package, since it accesses internal fields of the memclob.
func requireOrderExistsInMemclob(
	t *testing.T,
	ctx sdk.Context,
	order types.Order,
	memclob *MemClobPriceTimePriority,
) {
	// Verify the order exists on the correct side of the orderbook.
	orderbook, exists := memclob.orderbooks[order.GetClobPairId()]
	require.True(t, exists)
	require.True(t, doesOrderExistOnSide(t, order.OrderId, orderbook.GetSide(order.IsBuy())))

	// Verify the order exists in `orderIdToLevelOrder`.
	require.Contains(t, orderbook.orderIdToLevelOrder, order.OrderId)
	require.Equal(t, order, orderbook.orderIdToLevelOrder[order.OrderId].Value.Order)

	// Verify the order was added to the subaccounts open orders.
	subaccountOpenOrders, err := memclob.GetSubaccountOrders(
		order.GetClobPairId(),
		order.OrderId.SubaccountId,
		order.Side,
	)
	require.NoError(t, err)
	require.Contains(t, subaccountOpenOrders, order)

	// If the order is a Short-Term order, verify the order was added to the block expiration map.
	// Else, verify it was not added to the block expiration map.
	if order.OrderId.IsShortTermOrder() {
		require.Contains(t, orderbook.blockExpirationsForOrders[order.GetGoodTilBlock()], order.OrderId)
	} else {
		require.NotContains(t, orderbook.blockExpirationsForOrders[order.GetGoodTilBlock()], order.OrderId)
	}

	// If this is a reduce-only order, verify the order exists in the open reduce-only orders for
	// this subaccount. Else, verify it is not present.
	if order.IsReduceOnly() {
		require.Contains(
			t,
			orderbook.SubaccountOpenReduceOnlyOrders[order.OrderId.SubaccountId],
			order.OrderId,
		)
	} else {
		require.NotContains(
			t,
			orderbook.SubaccountOpenReduceOnlyOrders[order.OrderId.SubaccountId],
			order.OrderId,
		)
	}
}

// requireOrderDoesNotExistInMemclob is a testing helper used for verifying that `order.OrderId` does not exist in any
// of the memclob's data structures. As a caveat, this helper function will only verify the order is not contained
// within the orderbook specified by `order.ClobPairId` (and it will not check all orderbooks in the memclob).
// Note that this function must be defined within this package, since it accesses internal fields of the memclob.
func requireOrderDoesNotExistInMemclob(
	t *testing.T,
	ctx sdk.Context,
	order types.Order,
	memclob *MemClobPriceTimePriority,
) {
	// Verify the order is not found.
	_, found := memclob.GetOrder(order.OrderId)
	require.False(t, found)

	// Verify the order was not added to the subaccounts open orders.

	subaccountOpenOrders, err := memclob.GetSubaccountOrders(
		order.GetClobPairId(),
		order.OrderId.SubaccountId,
		order.Side,
	)
	require.NoError(t, err)
	require.NotContains(t, subaccountOpenOrders, order)

	// Verify the order does not exist on either side of the orderbook.
	orderbook, exists := memclob.orderbooks[order.GetClobPairId()]
	require.True(t, exists)
	require.False(t, doesOrderExistOnSide(t, order.OrderId, orderbook.Bids))
	require.False(t, doesOrderExistOnSide(t, order.OrderId, orderbook.Asks))

	// Verify the order was not added to the block expiration map.
	if order.OrderId.IsShortTermOrder() {
		require.NotContains(t, orderbook.blockExpirationsForOrders[order.GetGoodTilBlock()], order.OrderId)
	}

	// Verify the order does not exist in the open reduce-only orders for this subaccount.
	require.NotContains(
		t,
		orderbook.SubaccountOpenReduceOnlyOrders[order.OrderId.SubaccountId],
		order.OrderId,
	)
}

func setUpMemclobAndOrderbook(
	t *testing.T,
	ctx sdk.Context,
	placedMatchableOrders []types.MatchableOrder,
	getStatePosition types.GetStatePositionFn,
	newOrder []types.MatchableOrder,
) (
	memclob *MemClobPriceTimePriority,
	fakeMemClobKeeper *testutil_memclob.FakeMemClobKeeper,
) {
	// Setup the memclob state.
	memClobKeeper := testutil_memclob.NewFakeMemClobKeeper().
		WithStatePositionFn(getStatePosition).
		WithCollatCheckFnForProcessSingleMatch()
	memclob = NewMemClobPriceTimePriority(true)
	memclob.SetClobKeeper(memClobKeeper)

	// Create all unique orderbooks.
	createAllOrderbooksForMatchableOrders(
		t,
		ctx,
		memclob,
		append(placedMatchableOrders, newOrder...),
	)

	// Create all orders.
	createAllMatchableOrders(
		t,
		ctx,
		memclob,
		placedMatchableOrders,
		memClobKeeper,
	)
	return memclob, memClobKeeper
}

// placeOrderTestSetup is a testing helper used for properly setting up memclob state and expectations for testing
// `PlaceOrder`.
// Note that this function expects that every `UpdateResult` specified in `collatCheckFailures` is not successful.
// If `UpdateResult.Success` is specified for any collateralization check and subaccount, this testing function will
// fail.
func placeOrderTestSetup(
	t *testing.T,
	ctx sdk.Context,
	placedMatchableOrders []types.MatchableOrder,
	newOrder types.MatchableOrder,
	expectedPendingMatches []expectedMatch,
	expectedOrderStatus types.OrderStatus,
	addOrderToOrderbookSize satypes.BaseQuantums,
	expectedErr error,
	collatCheckFailures map[int]map[satypes.SubaccountId]satypes.UpdateResult,
	getStatePosition types.GetStatePositionFn,
) (
	memclob *MemClobPriceTimePriority,
	fakeMemClobKeeper *testutil_memclob.FakeMemClobKeeper,
	expectedNumCollateralizationChecks int,
	numCollateralizationChecksCounter *int,
) {
	memclob, memClobKeeper := setUpMemclobAndOrderbook(
		t,
		ctx,
		placedMatchableOrders,
		getStatePosition,
		[]types.MatchableOrder{newOrder},
	)
	collatCheckFailuresSet := make(map[int]bool)
	for i := range collatCheckFailures {
		collatCheckFailuresSet[i] = true
	}

	expectedCollatChecks := createCollatCheckExpectationsFromPendingMatches(
		ctx,
		t,
		memclob,
		expectedPendingMatches,
		collatCheckFailuresSet,
		newOrder,
		addOrderToOrderbookSize,
	)

	numCollateralChecks := 0

	collatCheckFn := testutil_memclob.CreateCollatCheckFunction(
		t,
		&numCollateralChecks,
		expectedCollatChecks,
		collatCheckFailures,
	)

	memClobKeeper.WithCollatCheckFn(collatCheckFn)

	expectedNumCollateralizationChecks = len(expectedCollatChecks)

	return memclob, memClobKeeper, expectedNumCollateralizationChecks, &numCollateralChecks
}

// simplePlaceOrderTestSetup is a testing helper used for properly setting up memclob state and expectations for testing
// `PlaceOrder`.
func simplePlaceOrderTestSetup(
	t *testing.T,
	ctx sdk.Context,
	placedMatchableOrders []types.MatchableOrder,
	collateralizationChecks map[int]testutil_memclob.CollateralizationCheck,
	getStatePosition types.GetStatePositionFn,
	newOrders ...types.MatchableOrder,
) (
	memclob *MemClobPriceTimePriority,
	fakeMemClobKeeper *testutil_memclob.FakeMemClobKeeper,
	expectedNumCollateralizationChecks int,
	numCollateralizationChecksCounter *int,
) {
	// Setup the memclob state.
	memclob, memClobKeeper := setUpMemclobAndOrderbook(
		t,
		ctx,
		placedMatchableOrders,
		getStatePosition,
		newOrders,
	)

	numCollateralChecks := 0

	collatCheckFn := testutil_memclob.CreateSimpleCollatCheckFunction(
		t,
		&numCollateralChecks,
		collateralizationChecks,
	)

	memClobKeeper.WithCollatCheckFn(collatCheckFn)

	expectedNumCollateralizationChecks = len(collateralizationChecks)

	return memclob, memClobKeeper, expectedNumCollateralizationChecks, &numCollateralChecks
}

// memclobOperationsTestSetupWithCustomCollatCheck a testing helper used for properly setting up
// memclob state by applying multiple operations and using a custom collat check function.
func memclobOperationsTestSetupWithCustomCollatCheck(
	t *testing.T,
	ctx sdk.Context,
	placedOperations []types.Operation,
	collatCheckFn types.AddOrderToOrderbookCollateralizationCheckFn,
	getStatePosition types.GetStatePositionFn,
	preexistingStatefulOrders []types.LongTermOrderPlacement,
) (
	memclob *MemClobPriceTimePriority,
	fakeMemClobKeeper *testutil_memclob.FakeMemClobKeeper,
) {
	// Setup the memclob state.
	memClobKeeper := testutil_memclob.NewFakeMemClobKeeper().
		WithStatePositionFn(getStatePosition).
		WithCollatCheckFnForProcessSingleMatch()
	memclob = NewMemClobPriceTimePriority(true)
	memclob.SetClobKeeper(memClobKeeper)

	// Create all unique orderbooks.
	createOrderbooks(
		t,
		ctx,
		memclob,
		3,
	)

	// Create all pre-existing stateful orders, verify the specified transaction index
	// is correct, and then commit the state.
	for i, statefulOrderPlacement := range preexistingStatefulOrders {
		require.Equal(t, uint32(i), statefulOrderPlacement.PlacementIndex.TransactionIndex)
		memClobKeeper.SetLongTermOrderPlacement(
			ctx,
			statefulOrderPlacement.Order,
			statefulOrderPlacement.PlacementIndex.BlockHeight,
		)
		memClobKeeper.MustAddOrderToStatefulOrdersTimeSlice(
			ctx,
			statefulOrderPlacement.Order.MustGetUnixGoodTilBlockTime(),
			statefulOrderPlacement.Order.OrderId,
		)
	}
	memClobKeeper.CommitState()

	memClobKeeper.WithCollatCheckFn(collatCheckFn)

	// Create all orders.
	applyOperationsToMemclob(
		t,
		ctx,
		memclob,
		placedOperations,
		memClobKeeper,
	)

	return memclob, memClobKeeper
}

// memclobOperationsTestSetup is a testing helper used for properly setting up memclob state by
// applying multiple operations.
func memclobOperationsTestSetup(
	t *testing.T,
	ctx sdk.Context,
	placedOperations []types.Operation,
	collateralizationChecks map[int]testutil_memclob.CollateralizationCheck,
	getStatePosition types.GetStatePositionFn,
	preexistingStatefulOrders []types.LongTermOrderPlacement,
) (
	memclob *MemClobPriceTimePriority,
	fakeMemClobKeeper *testutil_memclob.FakeMemClobKeeper,
) {
	numCollateralChecks := 0

	collatCheckFn := testutil_memclob.CreateSimpleCollatCheckFunction(
		t,
		&numCollateralChecks,
		collateralizationChecks,
	)

	memclob, memClobKeeper := memclobOperationsTestSetupWithCustomCollatCheck(
		t,
		ctx,
		placedOperations,
		collatCheckFn,
		getStatePosition,
		preexistingStatefulOrders,
	)

	require.Equal(t, len(collateralizationChecks), numCollateralChecks)

	return memclob, memClobKeeper
}

// placeOrderAndVerifyExpectations is a testing helper used for calling `PlaceOrder` and asserting expectations about
// the return values from `PlaceOrder` and memclob state. Specifically, it verifies the following are as expected:
// - The return values of `PlaceOrder`.
// - The number of collateralization checks performed.
// - Verifies that the order exists or does not exist on the orderbook.
//   - If the order is supposed to exist, it verifies that it was added to the back of the price level.
//
// - It asserts that the memclob contains all expected orders.
// - It asserts that the memclob contains all expected matches.
func placeOrderAndVerifyExpectations(
	t *testing.T,
	ctx sdk.Context,
	memclob *MemClobPriceTimePriority,
	order types.Order,
	numCollateralChecks *int,
	expectedFilledSize satypes.BaseQuantums,
	expectedTotalFilledSize satypes.BaseQuantums,
	expectedOrderStatus types.OrderStatus,
	expectedErr error,
	expectedNumCollateralizationChecks int,
	expectedRemainingBids []OrderWithRemainingSize,
	expectedRemainingAsks []OrderWithRemainingSize,
	expectedMatches []expectedMatch,
	fakeMemClobKeeper *testutil_memclob.FakeMemClobKeeper,
) *types.OffchainUpdates {
	filledSize,
		orderStatus,
		offchainUpdates,
		err := memclob.PlaceOrder(ctx, order)

	if fakeMemClobKeeper != nil {
		if err == nil {
			fakeMemClobKeeper.CommitState()
		} else {
			fakeMemClobKeeper.ResetState()
		}
	}

	orderbook := memclob.orderbooks[order.GetClobPairId()]

	// Verify the return values are correct.
	require.ErrorIs(t, err, expectedErr)
	require.Equal(t, expectedFilledSize, filledSize)
	// If no error occurred, validate the order status.
	if expectedErr == nil {
		require.Equal(t, expectedOrderStatus, orderStatus)
	}

	// Verify the collateralization check occurred the expected number of times.
	require.Equal(t, expectedNumCollateralizationChecks, *numCollateralChecks)

	orderQuantums := order.GetBaseQuantums()
	orderShouldNotExist := expectedTotalFilledSize == orderQuantums ||
		!expectedOrderStatus.IsSuccess() ||
		expectedErr != nil
	isReplaceError := errors.Is(expectedErr, types.ErrInvalidReplacement) ||
		errors.Is(expectedErr, types.ErrOrderFullyFilled)

	if !isReplaceError {
		if orderShouldNotExist {
			// No order was placed on the orderbook.
			requireOrderDoesNotExistInMemclob(t, ctx, order, memclob)
		} else {
			// Verify the order was placed on the orderbook.
			requireOrderExistsInMemclob(t, ctx, order, memclob)

			// Verify the order was placed at the end of the price level.
			var level *types.Level
			if order.IsBuy() {
				level = orderbook.Bids[order.GetOrderSubticks()]
			} else {
				level = orderbook.Asks[order.GetOrderSubticks()]
			}

			require.NotNil(t, level)
			require.Equal(t, order.OrderId, level.LevelOrders.Back.Value.Order.OrderId)
		}
	}

	AssertMemclobHasOrders(
		t,
		ctx,
		memclob,
		expectedRemainingBids,
		expectedRemainingAsks,
	)

	assertMemclobHasMatches(
		t,
		ctx,
		memclob,
		expectedMatches,
	)

	return offchainUpdates
}

// placeOrderAndVerifyExpectationsOperations is a testing helper used for calling `PlaceOrder` and
// asserting expectations about the return values from `PlaceOrder` and memclob state.
// Specifically, it verifies the following are as expected:
// - The return values of `PlaceOrder`.
// - The number of collateralization checks performed.
// - It asserts that the memclob operations queue contains all expected operations.
// - Verifies that the order exists or does not exist on the orderbook.
//   - If the order is supposed to exist, it verifies that it was added to the back of the price level.
func placeOrderAndVerifyExpectationsOperations(
	t *testing.T,
	ctx sdk.Context,
	memclob *MemClobPriceTimePriority,
	order types.Order,
	numCollateralChecks *int,
	expectedFilledSize satypes.BaseQuantums,
	expectedTotalFilledSize satypes.BaseQuantums,
	expectedOrderStatus types.OrderStatus,
	expectedErr error,
	expectedNumCollateralizationChecks int,
	expectedRemainingBids []OrderWithRemainingSize,
	expectedRemainingAsks []OrderWithRemainingSize,
	expectedOperations []types.Operation,
	expectedInternalOperations []types.InternalOperation,
	fakeMemClobKeeper *testutil_memclob.FakeMemClobKeeper,
) *types.OffchainUpdates {
	filledSize,
		orderStatus,
		offchainUpdates,
		err := memclob.PlaceOrder(ctx, order)

	if fakeMemClobKeeper != nil {
		if err == nil {
			fakeMemClobKeeper.CommitState()
		} else {
			fakeMemClobKeeper.ResetState()
		}
	}

	orderbook := memclob.orderbooks[order.GetClobPairId()]

	// Verify the return values are correct.
	require.ErrorIs(t, err, expectedErr)
	require.Equal(t, expectedFilledSize, filledSize)
	// If no error occurred, validate the order status.
	if expectedErr == nil {
		require.Equal(t, expectedOrderStatus, orderStatus)
	}

	// Verify the collateralization check occurred the expected number of times.
	require.Equal(t, expectedNumCollateralizationChecks, *numCollateralChecks)

	orderQuantums := order.GetBaseQuantums()
	orderShouldNotExist := expectedTotalFilledSize == orderQuantums ||
		!expectedOrderStatus.IsSuccess() ||
		expectedErr != nil
	isReplaceError := errors.Is(expectedErr, types.ErrInvalidReplacement) ||
		errors.Is(expectedErr, types.ErrOrderFullyFilled)

	if !isReplaceError {
		if orderShouldNotExist {
			// No order was placed on the orderbook.
			requireOrderDoesNotExistInMemclob(t, ctx, order, memclob)
		} else {
			// Verify the order was placed on the orderbook.
			requireOrderExistsInMemclob(t, ctx, order, memclob)

			// Verify the order was placed at the end of the price level.
			var level *types.Level
			if order.IsBuy() {
				level = orderbook.Bids[order.GetOrderSubticks()]
			} else {
				level = orderbook.Asks[order.GetOrderSubticks()]
			}

			require.NotNil(t, level)
			require.Equal(t, order.OrderId, level.LevelOrders.Back.Value.Order.OrderId)
		}
	}

	AssertMemclobHasOrders(
		t,
		ctx,
		memclob,
		expectedRemainingBids,
		expectedRemainingAsks,
	)

	assertMemclobHasOperations(
		t,
		ctx,
		memclob,
		expectedOperations,
		expectedInternalOperations,
	)

	AssertMemclobHasShortTermTxBytes(
		t,
		ctx,
		memclob,
		expectedInternalOperations,
		expectedRemainingBids,
		expectedRemainingAsks,
	)

	return offchainUpdates
}

// assertPlaceOrderOffchainMessages checks that the expected offchain update messages are returned
// from a call to `PlaceOrder`.
// This includes:
//   - an OrderPlace message is sent if order placement does not result in an error
//   - an OrderUpdate message is sent for the placed order if non-zero matches occur
//   - an OrderUpdate message is sent for each maker order matched against the placed order
//   - an OrderRemove message is sent for each maker order that failed collateralization checks
//     when matching with the placed order
//   - an OrderRemove message is sent if order placement does not result in an error but has a
//     non-success status
func assertPlaceOrderOffchainMessages(
	t *testing.T,
	ctx sdk.Context,
	offchainUpdates *types.OffchainUpdates,
	order types.Order,
	placedMatchableOrders []types.MatchableOrder,
	collatCheckFailures map[int]map[satypes.SubaccountId]satypes.UpdateResult,
	expectedErr error,
	expectedTotalFilledSize satypes.BaseQuantums,
	expectedOrderStatus types.OrderStatus,
	expectedExistingMatches []expectedMatch,
	expectedNewMatches []expectedMatch,
	expectedCancelledReduceOnlyOrders []types.OrderId,
	expectedToReplaceOrder bool,
) {
	actualOffchainMessages := offchainUpdates.GetMessages()
	expectedOffchainMessages := []msgsender.Message{}
	seenCancelledReduceOnlyOrders := mapset.NewSet[types.OrderId]()

	// If there are no errors expected, an order place message should be sent.
	if expectedErr == nil || doesErrorProduceOffchainMessages(expectedErr) {
		if expectedToReplaceOrder {
			removeMessage := off_chain_updates.MustCreateOrderRemoveMessageWithReason(
				ctx,
				order.OrderId,
				indexershared.OrderRemovalReason_ORDER_REMOVAL_REASON_REPLACED,
				ocutypes.OrderRemoveV1_ORDER_REMOVAL_STATUS_BEST_EFFORT_CANCELED,
			)
			expectedOffchainMessages = append(
				expectedOffchainMessages,
				removeMessage,
			)
		}
		placeMessage := off_chain_updates.MustCreateOrderPlaceMessage(
			ctx,
			order,
		)
		expectedOffchainMessages = append(
			expectedOffchainMessages,
			placeMessage,
		)
		require.Equal(t, expectedOffchainMessages, actualOffchainMessages[:len(expectedOffchainMessages)])
	}

	// Reduce-only order removals are sent before updates if the maker orders are removed during
	// matching.
	// Reduce-only order removals are also sent after matching, for orders from the same subaccount
	// as the taker order.
	for _, orderId := range expectedCancelledReduceOnlyOrders {
		cancelMessage := off_chain_updates.MustCreateOrderRemoveMessageWithReason(
			ctx,
			orderId,
			indexershared.OrderRemovalReason_ORDER_REMOVAL_REASON_REDUCE_ONLY_RESIZE,
			ocutypes.OrderRemoveV1_ORDER_REMOVAL_STATUS_BEST_EFFORT_CANCELED,
		)
		// If the reduce-only order was seen before updates, add it to the set so we don't try to check
		// for it later.
		if cmp.Equal(actualOffchainMessages[len(expectedOffchainMessages)], cancelMessage) {
			seenCancelledReduceOnlyOrders.Add(orderId)
			expectedOffchainMessages = append(
				expectedOffchainMessages,
				cancelMessage,
			)
			require.Equal(t, expectedOffchainMessages, actualOffchainMessages[:len(expectedOffchainMessages)])
		}
	}

	// For each order that was already placed on the book, if the order will fail collateralization
	// checks during matching with the placed order, an order removal message should be sent.
	for i, matchableOrder := range placedMatchableOrders {
		matchOrder := matchableOrder.MustGetOrder()
		// If the OrderIds match, this is a replacement order. Messages for replacements are tested
		// separately by the expectedToReplaceOrder argument passed in.
		if matchOrder.OrderId == order.OrderId {
			continue
		}
		if subaccountFailures, exists := collatCheckFailures[i]; exists {
			if _, exists = subaccountFailures[matchOrder.OrderId.SubaccountId]; exists {
				updateMessage := off_chain_updates.MustCreateOrderRemoveMessageWithReason(
					ctx,
					matchOrder.OrderId,
					indexershared.OrderRemovalReason_ORDER_REMOVAL_REASON_UNDERCOLLATERALIZED,
					ocutypes.OrderRemoveV1_ORDER_REMOVAL_STATUS_BEST_EFFORT_CANCELED,
				)

				expectedOffchainMessages = append(
					expectedOffchainMessages,
					updateMessage,
				)
				require.Equal(t, expectedOffchainMessages, actualOffchainMessages[:len(expectedOffchainMessages)])
			}
		}
	}

	// Keep track of the existing amount filled for orders.
	existingFilledAmounts := make(map[types.OrderId]satypes.BaseQuantums)
	for _, existingMatch := range expectedExistingMatches {
		ordersInMatch := []types.MatchableOrder{existingMatch.makerOrder, existingMatch.takerOrder}
		for _, order := range ordersInMatch {
			orderId := order.MustGetOrder().OrderId
			existingFilledAmounts[orderId] = existingMatch.matchedQuantums
		}
	}
	// For each new match generated from placing the order, the maker order in the match should have
	// an update message sent for it.
	for _, newMatch := range expectedNewMatches {
		makerOrder := newMatch.makerOrder.MustGetOrder()
		makerOrderId := makerOrder.OrderId
		totalFilled := newMatch.matchedQuantums
		// If there is an existing match for the maker order, the update message should contain the
		// sum of the existing match's fill and the new match's fill
		if existingFilled, exists := existingFilledAmounts[makerOrderId]; exists {
			totalFilled += existingFilled
		}

		updateMessage := off_chain_updates.MustCreateOrderUpdateMessage(
			ctx,
			makerOrderId,
			totalFilled,
		)
		expectedOffchainMessages = append(
			expectedOffchainMessages,
			updateMessage,
		)
		require.Equal(t, expectedOffchainMessages, actualOffchainMessages[:len(expectedOffchainMessages)])
	}

	// Reduce-only order removals are sent before updates if the maker orders are removed during
	// matching.
	// Reduce-only order removals are also sent after matching, for orders from the same subaccount
	// as the taker order.
	for _, orderId := range expectedCancelledReduceOnlyOrders {
		// Any reduce-only order removals that happened during matching don't need to be checked here.
		if seenCancelledReduceOnlyOrders.Contains(orderId) {
			continue
		}
		cancelMessage := off_chain_updates.MustCreateOrderRemoveMessageWithReason(
			ctx,
			orderId,
			indexershared.OrderRemovalReason_ORDER_REMOVAL_REASON_REDUCE_ONLY_RESIZE,
			ocutypes.OrderRemoveV1_ORDER_REMOVAL_STATUS_BEST_EFFORT_CANCELED,
		)
		expectedOffchainMessages = append(
			expectedOffchainMessages,
			cancelMessage,
		)
		require.Equal(t, expectedOffchainMessages, actualOffchainMessages[:len(expectedOffchainMessages)])
	}

	// If there is no error and the order status is success, the placed order should have an update
	// message sent for it.
	if expectedErr == nil && expectedOrderStatus.IsSuccess() {
		updateMessage := off_chain_updates.MustCreateOrderUpdateMessage(
			ctx,
			order.OrderId,
			expectedTotalFilledSize,
		)
		expectedOffchainMessages = append(
			expectedOffchainMessages,
			updateMessage,
		)
		require.Equal(t, expectedOffchainMessages, actualOffchainMessages[:len(expectedOffchainMessages)])
	}

	// If the order status after placing is not a success/resize but no error was returned, or an error
	// was returned which should result in a removal, an order removal message should be sent for the placed order.
	if doesErrorProduceOffchainMessages(expectedErr) ||
		(expectedErr == nil && !expectedOrderStatus.IsSuccess() && expectedOrderStatus != types.ReduceOnlyResized) {
		updateMessage := off_chain_updates.MustCreateOrderRemoveMessage(
			ctx,
			order.OrderId,
			expectedOrderStatus,
			expectedErr,
			ocutypes.OrderRemoveV1_ORDER_REMOVAL_STATUS_BEST_EFFORT_CANCELED,
		)
		expectedOffchainMessages = append(
			expectedOffchainMessages,
			updateMessage,
		)
		require.Equal(t, expectedOffchainMessages, actualOffchainMessages[:len(expectedOffchainMessages)])
	}

	require.Equal(t, expectedOffchainMessages, actualOffchainMessages)
}

// assertPlacePerpetualLiquidationOffchainMessages checks that the correct off-chain update messages
// were generated by a liquidation order.
func assertPlacePerpetualLiquidationOffchainMessages(
	t *testing.T,
	ctx sdk.Context,
	offchainUpdates *types.OffchainUpdates,
	liquidationOrder types.LiquidationOrder,
	placedMatchableOrders []types.MatchableOrder,
	collatCheckFailures map[int]map[satypes.SubaccountId]satypes.UpdateResult,
	expectedMatches []expectedMatch,
) {
	expectedOffchainMessages := getExpectedPlacePerpetualLiquidationOffchainMessages(
		t,
		ctx,
		liquidationOrder,
		placedMatchableOrders,
		collatCheckFailures,
		expectedMatches,
	)
	require.ElementsMatch(t, expectedOffchainMessages, offchainUpdates.GetMessages())
}

// getExpectedPlacePerpetualLiquidationOffchainMessages gets the expected off-chain update messages
// generated by a liquidation order. The messages include:
//   - `OrderRemove` message generated for any maker orders that fail collateralization
//   - `OrderRemove` message generated for any maker orders from the same subaccount as the liquidation
//     order that cross the liquidation order. This is due to the subaccount being undercollateralized.
//   - `OrderUpdate` message generated for any maker orders that match with the liquidation order
func getExpectedPlacePerpetualLiquidationOffchainMessages(
	t *testing.T,
	ctx sdk.Context,
	liquidationOrder types.LiquidationOrder,
	placedMatchableOrders []types.MatchableOrder,
	collatCheckFailures map[int]map[satypes.SubaccountId]satypes.UpdateResult,
	expectedMatches []expectedMatch,
) []msgsender.Message {
	expectedOffchainMessages := []msgsender.Message(nil)

	orderIdx := 0
	for _, matchableOrder := range placedMatchableOrders {
		// Skip any existing matchable orders that are liquidations, these should already have had
		// messages generated for them.
		if matchableOrder.IsLiquidation() {
			continue
		}
		order := matchableOrder.MustGetOrder()
		// Any orders from the same subaccount that cross with the liquidation order should be removed
		// due to it being a self-trade.
		if order.OrderId.SubaccountId == liquidationOrder.GetSubaccountId() &&
			doesLiquidationCrossOrder(liquidationOrder, order) {
			updateMessage := off_chain_updates.MustCreateOrderRemoveMessageWithReason(
				ctx,
				order.OrderId,
				indexershared.OrderRemovalReason_ORDER_REMOVAL_REASON_SELF_TRADE_ERROR,
				ocutypes.OrderRemoveV1_ORDER_REMOVAL_STATUS_BEST_EFFORT_CANCELED,
			)

			expectedOffchainMessages = append(
				expectedOffchainMessages,
				updateMessage,
			)
			// If the order would fail collateralization checks, an order removal message should be sent for
			// for the order.
		} else if subaccountFailures, exists := collatCheckFailures[orderIdx]; exists {
			if _, exists = subaccountFailures[order.OrderId.SubaccountId]; exists {
				updateMessage := off_chain_updates.MustCreateOrderRemoveMessageWithReason(
					ctx,
					order.OrderId,
					indexershared.OrderRemovalReason_ORDER_REMOVAL_REASON_UNDERCOLLATERALIZED,
					ocutypes.OrderRemoveV1_ORDER_REMOVAL_STATUS_BEST_EFFORT_CANCELED,
				)

				expectedOffchainMessages = append(
					expectedOffchainMessages,
					updateMessage,
				)
				// Increment the index to look at the in the expected collateralization failtures map.
				// As orders from the same subaccount are expected to already be canceled, they are not
				// included in the map of expected collateralization check failures.
				orderIdx = orderIdx + 1
			}
		}
	}

	// Any matches with maker orders from the liquidation order should result in an update to the
	// total filled of the maker order.
	for _, match := range expectedMatches {
		makerOrder := match.makerOrder.MustGetOrder()
		makerOrderId := makerOrder.OrderId
		totalFilled := match.matchedQuantums

		updateMessage := off_chain_updates.MustCreateOrderUpdateMessage(
			ctx,
			makerOrderId,
			totalFilled,
		)
		expectedOffchainMessages = append(
			expectedOffchainMessages,
			updateMessage,
		)
	}

	return expectedOffchainMessages
}

// doesErrorProduceOffchainMessages returns true if the provided error is not nil and
// should return offchain add/removal messages
func doesErrorProduceOffchainMessages(err error) bool {
	return err != nil &&
		(errors.Is(err, types.ErrPostOnlyWouldCrossMakerOrder) ||
			errors.Is(err, types.ErrFokOrderCouldNotBeFullyFilled))
}

// doesLiquidationOverlapOrder checks if a liquidation order crosses another order.
// For BUY liquidation orders, this is true if the other order has a price equal to or lower than
// the price of the liquidation order.
// For SELL liquidation orders, this is true if the other order has a price equal to or greater than
// the price of the liquidation order.
func doesLiquidationCrossOrder(
	liquidationOrder types.LiquidationOrder,
	order types.Order,
) bool {
	if liquidationOrder.IsBuy() {
		return order.Subticks <= liquidationOrder.GetOrderSubticks().ToUint64()
	} else {
		return order.Subticks >= liquidationOrder.GetOrderSubticks().ToUint64()
	}
}

// placePerpetualLiquidationAndVerifyExpectations is a testing helper used for calling
// `PlacePerpetualLiquidation` and asserting expectations about memclob state. Specifically, it
// verifies the following are as expected:
// - The number of collateralization checks performed.
// - Verifies that the order does not exist on the orderbook.
// - It asserts that the memclob contains all expected orders.
// - It asserts that the memclob contains all expected matches.
func placePerpetualLiquidationAndVerifyExpectations(
	t *testing.T,
	ctx sdk.Context,
	memclob *MemClobPriceTimePriority,
	order types.LiquidationOrder,
	numCollateralChecks *int,
	expectedNumCollateralizationChecks int,
	expectedRemainingBids []OrderWithRemainingSize,
	expectedRemainingAsks []OrderWithRemainingSize,
	expectedMatches []expectedMatch,
) *types.OffchainUpdates {
	// Run the test case.
	_, _, offchainUpdates, err := memclob.PlacePerpetualLiquidation(
		ctx,
		order,
	)
	require.NoError(t, err)

	// Verify the collateralization check occurred the expected number of times.
	require.Equal(t, expectedNumCollateralizationChecks, *numCollateralChecks)

	AssertMemclobHasOrders(
		t,
		ctx,
		memclob,
		expectedRemainingBids,
		expectedRemainingAsks,
	)

	assertMemclobHasMatches(
		t,
		ctx,
		memclob,
		expectedMatches,
	)

	return offchainUpdates
}

// placePerpetualLiquidationAndVerifyExpectationsOperations is a testing helper used for calling
// `PlacePerpetualLiquidation` and asserting expectations about memclob state. Specifically, it
// verifies the following are as expected:
// - The number of collateralization checks performed.
// - Verifies that the order does not exist on the orderbook.
// - It asserts that the memclob contains all expected orders.
// - It asserts that the memclob operations queue contains all expected operations.
func placePerpetualLiquidationAndVerifyExpectationsOperations(
	t *testing.T,
	ctx sdk.Context,
	memclob *MemClobPriceTimePriority,
	order types.LiquidationOrder,
	numCollateralChecks *int,
	expectedNumCollateralizationChecks int,
	expectedRemainingBids []OrderWithRemainingSize,
	expectedRemainingAsks []OrderWithRemainingSize,
	expectedOperations []types.Operation,
	expectedInternalOperations []types.InternalOperation,
) *types.OffchainUpdates {
	// Run the test case.
	_, _, offchainUpdates, err := memclob.PlacePerpetualLiquidation(
		ctx,
		order,
	)
	require.NoError(t, err)

	// Verify the collateralization check occurred the expected number of times.
	require.Equal(t, expectedNumCollateralizationChecks, *numCollateralChecks)

	AssertMemclobHasOrders(
		t,
		ctx,
		memclob,
		expectedRemainingBids,
		expectedRemainingAsks,
	)

	assertMemclobHasOperations(
		t,
		ctx,
		memclob,
		expectedOperations,
		expectedInternalOperations,
	)

	AssertMemclobHasShortTermTxBytes(
		t,
		ctx,
		memclob,
		expectedInternalOperations,
		expectedRemainingBids,
		expectedRemainingAsks,
	)

	return offchainUpdates
}
