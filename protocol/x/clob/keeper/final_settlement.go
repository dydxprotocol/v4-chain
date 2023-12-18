package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	indexerevents "github.com/dydxprotocol/v4-chain/protocol/indexer/events"
	"github.com/dydxprotocol/v4-chain/protocol/indexer/indexer_manager"
	indexershared "github.com/dydxprotocol/v4-chain/protocol/indexer/shared"
	"github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
)

// MustEnterFinalSettlement holds logic executed when a market transitions to FINAL_SETTLEMENT status.
// This function will forcefully cancel all stateful open orders for the clob pair.
func (k Keeper) mustEnterFinalSettlement(ctx sdk.Context, clobPairId types.ClobPairId) {
	// Forcefully cancel all stateful orders from state for this clob pair
	k.mustCancelStatefulOrdersForFinalSettlement(ctx, clobPairId)

	// Delete untriggered conditional orders for this clob pair from memory
	delete(k.UntriggeredConditionalOrders, clobPairId)
}

// mustCancelStatefulOrdersForFinalSettlement forcefully cancels all stateful orders
// for the provided ClobPair. These orders will be removed from the memclob in PrepareCheckState.
func (k Keeper) mustCancelStatefulOrdersForFinalSettlement(ctx sdk.Context, clobPairId types.ClobPairId) {
	statefulOrders := k.GetAllStatefulOrders(ctx)
	processProposerMatchesEvents := k.GetProcessProposerMatchesEvents(ctx)

	// This logic is executed in EndBlocker and should not panic. This would be unexpected,
	// but if it happens we would rather recover and continue if an order fails to be removed from state
	// rather than halt the chain.
	safelyRemoveStatefulOrder := func(ctx sdk.Context, orderId types.OrderId) {
		defer func() {
			if r := recover(); r != nil {
				k.Logger(ctx).Error(
					"mustCancelStatefulOrdersForFinalSettlement: Failed to remove stateful order with OrderId %+v: %v",
					orderId,
					r,
				)
			}
		}()
		k.MustRemoveStatefulOrder(ctx, orderId)
	}

	// TODO(CLOB-1053): Iterate over stateful orders for only specified clob pair
	for _, order := range statefulOrders {
		if order.GetClobPairId() == clobPairId {
			// Remove from state, recovering from panic if necessary
			safelyRemoveStatefulOrder(ctx, order.OrderId)

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
	}

	k.MustSetProcessProposerMatchesEvents(ctx, processProposerMatchesEvents)
}
