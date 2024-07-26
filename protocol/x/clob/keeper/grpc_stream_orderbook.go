package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	v1 "github.com/dydxprotocol/v4-chain/protocol/indexer/protocol/v1"
	v1types "github.com/dydxprotocol/v4-chain/protocol/indexer/protocol/v1/types"
	"github.com/dydxprotocol/v4-chain/protocol/lib"
	"github.com/dydxprotocol/v4-chain/protocol/streaming/grpc/client"
	"github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
	satypes "github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"
)

func (k Keeper) StreamOrderbookUpdates(
	req *types.StreamOrderbookUpdatesRequest,
	stream types.Query_StreamOrderbookUpdatesServer,
) error {
	err := k.GetGrpcStreamingManager().Subscribe(*req, stream)
	if err != nil {
		return err
	}

	// Keep this scope alive because once this scope exits - the stream is closed
	ctx := stream.Context()
	<-ctx.Done()
	return nil
}

// Compare the aggregated size for each price level.
func (k Keeper) CompareMemclobOrderbookWithLocalOrderbook(
	ctx sdk.Context,
	localOrderbook *client.LocalOrderbook,
	id types.ClobPairId,
) {
	localOrderbook.Lock()
	defer localOrderbook.Unlock()

	logger := k.Logger(ctx).With("module", "grpc-example-client").With("block", ctx.BlockHeight()).With("clob_pair_id", id)

	logger.Info("Comparing grpc orderbook with actual memclob orderbook!")

	orderbook := k.MemClob.GetOrderbook(ctx, id)
	orderbookBids := orderbook.GetBids()

	// Compare bids.
	bids := lib.GetSortedKeys[lib.Sortable[types.Subticks]](orderbookBids)

	logger.Info("Comparing bids", "bids", bids)
	if len(bids) != len(localOrderbook.Bids) {
		logger.Error(
			"Bids length mismatch",
			"expected", len(bids),
			"actual", len(localOrderbook.Bids),
		)
	}

	for _, bid := range bids {
		level := orderbookBids[bid]

		expectedAggregatedQuantity := uint64(0)
		expectedOrders := make([]types.Order, 0)
		expectedRemainingAmounts := make([]uint64, 0)
		for levelOrder := level.LevelOrders.Front; levelOrder != nil; levelOrder = levelOrder.Next {
			order := levelOrder.Value
			_, filledAmount, _ := k.GetOrderFillAmount(ctx, order.Order.OrderId)
			expectedAggregatedQuantity += (order.Order.Quantums - filledAmount.ToUint64())
			expectedOrders = append(expectedOrders, order.Order)
			expectedRemainingAmounts = append(expectedRemainingAmounts, order.Order.Quantums-filledAmount.ToUint64())
		}

		actualAggregatedQuantity := uint64(0)
		actualOrders := make([]v1types.IndexerOrder, 0)
		actualRemainingAmounts := make([]uint64, 0)
		for _, order := range localOrderbook.Bids[bid.ToUint64()] {
			remainingAmount := order.Quantums - localOrderbook.FillAmounts[order.OrderId]
			actualAggregatedQuantity += remainingAmount
			actualOrders = append(actualOrders, order)
			actualRemainingAmounts = append(actualRemainingAmounts, remainingAmount)
		}

		// Compare the aggregated quantity as a basic sanity check.
		if expectedAggregatedQuantity != actualAggregatedQuantity {
			logger.Error(
				"Aggregated quantity mismatch for bid level",
				"price", bid,
				"expected", expectedAggregatedQuantity,
				"actual", actualAggregatedQuantity,
				"expected_orders", expectedOrders,
				"actual_orders", actualOrders,
				"expected_remaining_amounts", expectedRemainingAmounts,
				"actual_remaining_amounts", actualRemainingAmounts,
			)
		}

		if len(expectedOrders) != len(actualOrders) {
			logger.Error(
				"Different number of orders at bid level",
				"price", bid,
				"expected", expectedOrders,
				"actual", actualOrders,
			)
		} else {
			for i, expected := range expectedOrders {
				if expected.OrderId.ClientId != actualOrders[i].OrderId.ClientId {
					logger.Error(
						"Different order at bid level",
						"price", bid,
						"expected", expected,
						"actual", actualOrders[i],
					)
				}
				if expectedRemainingAmounts[i] != actualRemainingAmounts[i] {
					logger.Error(
						"Different remaining amount at bid level",
						"price", bid,
						"expected", expectedRemainingAmounts[i],
						"actual", actualRemainingAmounts[i],
					)
				}
			}
		}
	}

	orderbookAsks := orderbook.GetAsks()

	// Compare asks.
	asks := lib.GetSortedKeys[lib.Sortable[types.Subticks]](orderbookAsks)

	logger.Info("Comparing asks", "asks", asks)
	if len(asks) != len(localOrderbook.Asks) {
		logger.Error(
			"Asks length mismatch",
			"expected", len(asks),
			"actual", len(localOrderbook.Asks),
		)
	}

	for _, ask := range asks {
		level := orderbookAsks[ask]

		expectedAggregatedQuantity := uint64(0)
		expectedOrders := make([]types.Order, 0)
		expectedRemainingAmounts := make([]uint64, 0)
		for levelOrder := level.LevelOrders.Front; levelOrder != nil; levelOrder = levelOrder.Next {
			order := levelOrder.Value
			_, filledAmount, _ := k.GetOrderFillAmount(ctx, order.Order.OrderId)
			expectedAggregatedQuantity += (order.Order.Quantums - filledAmount.ToUint64())
			expectedOrders = append(expectedOrders, order.Order)
			expectedRemainingAmounts = append(expectedRemainingAmounts, order.Order.Quantums-filledAmount.ToUint64())
		}

		actualAggregatedQuantity := uint64(0)
		actualOrders := make([]v1types.IndexerOrder, 0)
		actualRemainingAmounts := make([]uint64, 0)
		for _, order := range localOrderbook.Asks[ask.ToUint64()] {
			remainingAmount := order.Quantums - localOrderbook.FillAmounts[order.OrderId]
			actualAggregatedQuantity += remainingAmount
			actualOrders = append(actualOrders, order)
			actualRemainingAmounts = append(actualRemainingAmounts, remainingAmount)
		}

		// Compare the aggregated quantity as a basic sanity check.
		if expectedAggregatedQuantity != actualAggregatedQuantity {
			logger.Error(
				"Aggregated quantity mismatch for ask level",
				"price", ask,
				"expected", expectedAggregatedQuantity,
				"actual", actualAggregatedQuantity,
				"expected_orders", expectedOrders,
				"actual_orders", actualOrders,
				"expected_remaining_amounts", expectedRemainingAmounts,
				"actual_remaining_amounts", actualRemainingAmounts,
			)
		}

		if len(expectedOrders) != len(actualOrders) {
			logger.Error(
				"Different number of orders at ask level",
				"price", ask,
				"expected", expectedOrders,
				"actual", actualOrders,
			)
		} else {
			for i, expected := range expectedOrders {
				if expected.OrderId.ClientId != actualOrders[i].OrderId.ClientId {
					logger.Error(
						"Different order at ask level",
						"price", ask,
						"expected", expected,
						"actual", actualOrders[i],
					)
				}
				if expectedRemainingAmounts[i] != actualRemainingAmounts[i] {
					logger.Error(
						"Different remaining amount at ask level",
						"price", ask,
						"expected", expectedRemainingAmounts[i],
						"actual", actualRemainingAmounts[i],
					)
				}
			}
		}
	}

	// Compare Fills in State with fills on the locally constructed orderbook from
	// grpc stream.
	numFailed := 0
	numPassed := 0
	allFillStates := k.GetAllOrderFillStates(ctx)
	for _, fillState := range allFillStates {
		orderFillAmount := fillState.FillAmount
		orderId := fillState.OrderId
		// Skip check for non-relevant clob pair id
		if orderId.ClobPairId != uint32(id) {
			continue
		}

		indexerOrderId := v1.OrderIdToIndexerOrderId(orderId)
		localOrderbookFillAmount := localOrderbook.FillAmounts[indexerOrderId]

		if orderFillAmount != localOrderbookFillAmount {
			logger.Error(
				"Fill Amount Mismatch",
				"orderId", orderId.String(),
				"state_fill_amt", orderFillAmount,
				"local_fill_amt", localOrderbookFillAmount,
			)
			numFailed += 1
		} else {
			numPassed += 1
		}
	}

	// Check if the locally constructed orderbook has extraneous order ids in the fill amounts
	// when compared to state.
	numInOrderbookButNotState := 0
	for indexerOrderId, localFillAmount := range localOrderbook.FillAmounts {
		clobOrderId := types.OrderId{
			SubaccountId: satypes.SubaccountId{
				Owner:  indexerOrderId.SubaccountId.Owner,
				Number: indexerOrderId.SubaccountId.Number,
			},
			ClientId:   indexerOrderId.ClientId,
			OrderFlags: indexerOrderId.OrderFlags,
			ClobPairId: indexerOrderId.ClobPairId,
		}
		exists, _, _ := k.GetOrderFillAmount(ctx, clobOrderId)
		if !exists {
			numInOrderbookButNotState += 1
			logger.Error(
				"Fill Amount Mismatch: exists in local orderbook but not in state",
				"orderId", clobOrderId.String(),
				"local_fill_amt", localFillAmount,
			)
		}
	}

	ratio := float32(numFailed) * 100 / float32(numPassed+numFailed)
	if numFailed+numPassed > 0 {
		if numFailed > 0 {
			logger.Error(
				fmt.Sprintf("Fill amount comparison results: %.2f failed", ratio),
				"failed", numFailed,
				"passed", numPassed,
				"in_orderbook_not_state", numInOrderbookButNotState,
			)
		} else {
			logger.Info("Orderbook comparison done successfully!")
		}
	} else {
		logger.Info("No checks done on fill amounts.")
	}
}
