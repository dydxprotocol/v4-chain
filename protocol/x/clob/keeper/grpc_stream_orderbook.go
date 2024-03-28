package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	v1types "github.com/dydxprotocol/v4-chain/protocol/indexer/protocol/v1/types"
	"github.com/dydxprotocol/v4-chain/protocol/lib"
	"github.com/dydxprotocol/v4-chain/protocol/streaming/grpc/client"
	"github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
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

	// Compare bids.
	bids := lib.GetSortedKeys[lib.Sortable[types.Subticks]](orderbook.Bids)

	logger.Info("Comparing bids", "bids", bids)
	if len(bids) != len(localOrderbook.Bids) {
		logger.Error(
			"Bids length mismatch",
			"expected", len(bids),
			"actual", len(localOrderbook.Bids),
		)
	}

	for _, bid := range bids {
		level := orderbook.Bids[bid]

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
			remainingAmount := localOrderbook.OrderRemainingAmount[order.OrderId]
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
	}

	// Compare asks.
	asks := lib.GetSortedKeys[lib.Sortable[types.Subticks]](orderbook.Asks)

	logger.Info("Comparing asks", "asks", asks)
	if len(asks) != len(localOrderbook.Asks) {
		logger.Error(
			"Asks length mismatch",
			"expected", len(asks),
			"actual", len(localOrderbook.Asks),
		)
	}

	for _, ask := range asks {
		level := orderbook.Asks[ask]

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
			remainingAmount := localOrderbook.OrderRemainingAmount[order.OrderId]
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
	}

	logger.Info("Orderbook comparison done!")
}
