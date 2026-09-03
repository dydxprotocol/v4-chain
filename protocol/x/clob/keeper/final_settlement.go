package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	indexerevents "github.com/dydxprotocol/v4-chain/protocol/indexer/events"
	"github.com/dydxprotocol/v4-chain/protocol/indexer/indexer_manager"
	indexershared "github.com/dydxprotocol/v4-chain/protocol/indexer/shared/types"
	"github.com/dydxprotocol/v4-chain/protocol/lib/log"
	"github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
)

// transitionToFinalSettlement holds logic executed when a market transitions to FINAL_SETTLEMENT status.
// This function will forcefully cancel all stateful open orders for the clob pair and sweep any isolated
// insurance fund for the market into the cross insurance fund. An error fails the enclosing
// UpdateClobPair atomically; recovery requires a new proposal, since x/delaymsg deletes a
// delayed MsgUpdateClobPair even when its execution fails.
func (k Keeper) transitionToFinalSettlement(ctx sdk.Context, clobPairId types.ClobPairId, perpetualId uint32) error {
	// Forcefully cancel all stateful orders from state for this clob pair.
	k.mustCancelStatefulOrdersForFinalSettlement(ctx, clobPairId)

	// Sweep any isolated insurance fund into the cross insurance fund (no-op for cross markets).
	// This is safe even with positions still open: only deleveraging can occur in final settlement,
	// which never touches the insurance fund. Deposits made to the isolated fund address after this
	// sweep are not re-swept.
	return k.subaccountsKeeper.TransferIsolatedInsuranceFundToCross(ctx, perpetualId)
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
