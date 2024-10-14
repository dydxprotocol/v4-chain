package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	indexerevents "github.com/dydxprotocol/v4-chain/protocol/indexer/events"
	"github.com/dydxprotocol/v4-chain/protocol/indexer/indexer_manager"
	indexershared "github.com/dydxprotocol/v4-chain/protocol/indexer/shared/types"
	"github.com/dydxprotocol/v4-chain/protocol/lib/log"
	"github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
)

// mustTransitionToFinalSettlement holds logic executed when a market transitions to FINAL_SETTLEMENT status.
// This function will forcefully cancel all stateful open orders for the clob pair.
func (k Keeper) mustTransitionToFinalSettlement(ctx sdk.Context, clobPairId types.ClobPairId) {
	// Forcefully cancel all stateful orders from state for this clob pair.
	k.mustCancelStatefulOrdersForFinalSettlement(ctx, clobPairId)
}

// mustCancelStatefulOrdersForFinalSettlement forcefully cancels all stateful orders
// for the provided ClobPair. These orders will be removed from the memclob in PrepareCheckState.
func (k Keeper) mustCancelStatefulOrdersForFinalSettlement(ctx sdk.Context, clobPairId types.ClobPairId) {
	statefulOrders := k.GetAllStatefulOrders(ctx)
	processProposerMatchesEvents := k.GetProcessProposerMatchesEvents(ctx)

	// This logic is executed in EndBlocker and should not panic. This would be unexpected,
	// but if it happens we would rather recover and continue if an order fails to be removed from state
	// rather than halt the chain.
	removeStatefulOrderWithoutPanicing := func(ctx sdk.Context, orderId types.OrderId) {
		defer func() {
			if err := recover(); err != nil {
				log.ErrorLog(
					ctx,
					"mustCancelStatefulOrdersForFinalSettlement: Failed to remove stateful order",
					"orderId",
					orderId,
					"error",
					err,
				)
			}
		}()
		k.MustRemoveStatefulOrder(ctx, orderId)
	}

	// TODO(CLOB-1053): Iterate over stateful orders for only specified clob pair
	for _, order := range statefulOrders {
		if order.GetClobPairId() != clobPairId {
			continue
		}

		// Remove from state, recovering from panic if necessary
		removeStatefulOrderWithoutPanicing(ctx, order.OrderId)

		// Append to RemovedStatefulOrderIds so this order gets removed
		// from the memclob in PrepareCheckState during the PurgeInvalidMemclobState step
		processProposerMatchesEvents.RemovedStatefulOrderIds = append(
			processProposerMatchesEvents.RemovedStatefulOrderIds,
			order.OrderId,
		)

		k.GetIndexerEventManager().AddTxnEvent(
			ctx,
			indexerevents.SubtypeStatefulOrder,
			indexerevents.StatefulOrderEventVersion,
			indexer_manager.GetBytes(
				indexerevents.NewStatefulOrderRemovalEvent(
					order.OrderId,
					indexershared.OrderRemovalReason_ORDER_REMOVAL_REASON_FINAL_SETTLEMENT,
				),
			),
		)
	}

	k.MustSetProcessProposerMatchesEvents(ctx, processProposerMatchesEvents)
}
