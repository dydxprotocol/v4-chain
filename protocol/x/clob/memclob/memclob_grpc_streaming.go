package memclob

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4-chain/protocol/indexer/off_chain_updates"
	ocutypes "github.com/dydxprotocol/v4-chain/protocol/indexer/off_chain_updates/types"
	indexersharedtypes "github.com/dydxprotocol/v4-chain/protocol/indexer/shared/types"
	"github.com/dydxprotocol/v4-chain/protocol/lib"
	"github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
)

// GetOffchainUpdatesForOrderbookSnapshot returns the offchain updates for the orderbook snapshot.
// This is used by the gRPC streaming server to send the orderbook snapshot to the client.
func (m *MemClobPriceTimePriority) GetOffchainUpdatesForOrderbookSnapshot(
	ctx sdk.Context,
	clobPairId types.ClobPairId,
) (offchainUpdates *types.OffchainUpdates) {
	offchainUpdates = types.NewOffchainUpdates()

	if orderbook, exists := m.openOrders.orderbooksMap[clobPairId]; exists {
		// Generate the offchain updates for buy orders.
		// Updates are sorted in descending order of price.
		buyPriceLevels := lib.GetSortedKeys[lib.Sortable[types.Subticks]](orderbook.Bids)
		for i := len(buyPriceLevels) - 1; i >= 0; i-- {
			subticks := buyPriceLevels[i]
			level := orderbook.Bids[subticks]

			// For each price level, generate offchain updates for each order in the level.
			level.LevelOrders.Front.Each(
				func(order types.ClobOrder) {
					offchainUpdates.Append(
						m.GetOrderbookUpdatesForOrderPlacement(ctx, order.Order),
					)
				},
			)
		}

		// Generate the offchain updates for sell orders.
		// Updates are sorted in ascending order of price.
		sellPriceLevels := lib.GetSortedKeys[lib.Sortable[types.Subticks]](orderbook.Asks)
		for i := 0; i < len(sellPriceLevels); i++ {
			subticks := sellPriceLevels[i]
			level := orderbook.Asks[subticks]

			// For each price level, generate offchain updates for each order in the level.
			level.LevelOrders.Front.Each(
				func(order types.ClobOrder) {
					offchainUpdates.Append(
						m.GetOrderbookUpdatesForOrderPlacement(ctx, order.Order),
					)
				},
			)
		}
	}

	return offchainUpdates
}

// GetOrderbookUpdatesForOrderPlacement returns a place order offchain message and
// a update order offchain message used to add an order for
// the orderbook grpc stream.
func (m *MemClobPriceTimePriority) GetOrderbookUpdatesForOrderPlacement(
	ctx sdk.Context,
	order types.Order,
) (offchainUpdates *types.OffchainUpdates) {
	offchainUpdates = types.NewOffchainUpdates()
	orderId := order.OrderId

	// Generate a order place message.
	if message, success := off_chain_updates.CreateOrderPlaceMessage(
		ctx,
		order,
	); success {
		offchainUpdates.AddPlaceMessage(orderId, message)
	}

	// Get the current fill amount of the order.
	fillAmount := m.GetOrderFilledAmount(ctx, orderId)

	// Generate an update message updating the total filled amount of order.
	if message, success := off_chain_updates.CreateOrderUpdateMessage(
		ctx,
		orderId,
		fillAmount,
	); success {
		offchainUpdates.AddUpdateMessage(orderId, message)
	}

	return offchainUpdates
}

// GetOrderbookUpdatesForOrderRemoval returns a remove order offchain message
// used to remove an order for the orderbook grpc stream.
func (m *MemClobPriceTimePriority) GetOrderbookUpdatesForOrderRemoval(
	ctx sdk.Context,
	orderId types.OrderId,
) (offchainUpdates *types.OffchainUpdates) {
	offchainUpdates = types.NewOffchainUpdates()
	if message, success := off_chain_updates.CreateOrderRemoveMessageWithReason(
		ctx,
		orderId,
		indexersharedtypes.OrderRemovalReason_ORDER_REMOVAL_REASON_UNSPECIFIED,
		ocutypes.OrderRemoveV1_ORDER_REMOVAL_STATUS_BEST_EFFORT_CANCELED,
	); success {
		offchainUpdates.AddRemoveMessage(orderId, message)
	}
	return offchainUpdates
}

// GetOrderbookUpdatesForOrderUpdate returns an update order offchain message
// used to update an order for the orderbook grpc stream.
func (m *MemClobPriceTimePriority) GetOrderbookUpdatesForOrderUpdate(
	ctx sdk.Context,
	orderId types.OrderId,
) (offchainUpdates *types.OffchainUpdates) {
	offchainUpdates = types.NewOffchainUpdates()

	// Get the current fill amount of the order.
	fillAmount := m.GetOrderFilledAmount(ctx, orderId)

	// Generate an update message updating the total filled amount of order.
	if message, success := off_chain_updates.CreateOrderUpdateMessage(
		ctx,
		orderId,
		fillAmount,
	); success {
		offchainUpdates.AddUpdateMessage(orderId, message)
	}
	return offchainUpdates
}
