package keeper

import (
	"context"
	"errors"

	errorsmod "cosmossdk.io/errors"
	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"
	indexerevents "github.com/dydxprotocol/v4-chain/protocol/indexer/events"
	"github.com/dydxprotocol/v4-chain/protocol/indexer/indexer_manager"
	indexershared "github.com/dydxprotocol/v4-chain/protocol/indexer/shared/types"
	"github.com/dydxprotocol/v4-chain/protocol/lib"
	"github.com/dydxprotocol/v4-chain/protocol/lib/log"
	"github.com/dydxprotocol/v4-chain/protocol/lib/metrics"
	"github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
)

// CancelOrder performs order cancellation functionality for stateful orders.
func (k msgServer) CancelOrder(
	goCtx context.Context,
	msg *types.MsgCancelOrder,
) (resp *types.MsgCancelOrderResponse, err error) {
	ctx := lib.UnwrapSDKContext(goCtx, types.ModuleName)

	if err := k.Keeper.HandleMsgCancelOrder(ctx, msg); err != nil {
		return nil, err
	}

	return &types.MsgCancelOrderResponse{}, nil
}

// HandleMsgCancelOrder handles a MsgCancelOrder by
// 1. persisting the cancellation on chain.
// 2. updating ProcessProposerMatchesEvents with the new stateful order cancellation.
// 3. adding order cancellation on-chain indexer event.
func (k Keeper) HandleMsgCancelOrder(
	ctx sdk.Context,
	msg *types.MsgCancelOrder,
) (err error) {
	lib.AssertDeliverTxMode(ctx)

	// Attach various logging tags relative to this request. These should be static with no changes.
	ctx = log.AddPersistentTagsToLogger(ctx,
		log.Module, log.Clob,
		log.ProposerConsAddress, sdk.ConsAddress(ctx.BlockHeader().ProposerAddress),
		log.Callback, lib.TxMode(ctx),
		log.BlockHeight, ctx.BlockHeight(),
		log.Handler, log.CancelOrder,
	)

	defer func() {
		metrics.IncrSuccessOrErrorCounter(
			err,
			types.ModuleName,
			metrics.CancelOrder,
			metrics.DeliverTx,
			msg.OrderId.GetOrderIdLabels()...,
		)
		if err != nil {
			// Gracefully handle the case where the order was already removed from state. This can happen if an Order
			// Removal Operation was included in the same block as the MsgCancelOrder. By the time we try to cancel
			// the order, it has already been removed from state due to errors encountered while matching.
			// TODO(CLOB-778): Prevent invalid MsgCancelOrder messages from being included in the block.
			if errors.Is(err, types.ErrStatefulOrderDoesNotExist) {
				processProposerMatchesEvents := k.GetProcessProposerMatchesEvents(ctx)
				removedOrderIds := lib.UniqueSliceToSet(processProposerMatchesEvents.RemovedStatefulOrderIds)
				if _, found := removedOrderIds[msg.GetOrderId()]; found {
					telemetry.IncrCounterWithLabels(
						[]string{
							types.ModuleName,
							metrics.StatefulCancellationMsgHandlerFailure,
							metrics.StatefulOrderAlreadyRemoved,
							metrics.Count,
						},
						1,
						msg.OrderId.GetOrderIdLabels(),
					)
					err = errorsmod.Wrapf(
						types.ErrStatefulOrderCancellationFailedForAlreadyRemovedOrder,
						"Error: %s",
						err.Error(),
					)
					log.InfoLog(ctx, "Cancel Order Expected Error", log.Error, err)
					return
				}
			}
			log.ErrorLogWithError(ctx, "Error cancelling order", err)
		}
	}()

	// 1. If this is a Short-Term order, panic.
	msg.OrderId.MustBeStatefulOrder()

	// 2. Cancel the order on the ClobKeeper which is responsible for:
	//   - stateful cancellation validation.
	//   - removing the order from state and the memstore.
	if err := k.CancelStatefulOrder(ctx, msg); err != nil {
		return err
	}

	// 3. Update memstore with the new stateful order cancellation.
	k.AddDeliveredCancelledOrderId(ctx, msg.OrderId)

	// 4. Add the relevant on-chain Indexer event for the cancellation.
	k.GetIndexerEventManager().AddTxnEvent(
		ctx,
		indexerevents.SubtypeStatefulOrder,
		indexerevents.StatefulOrderEventVersion,
		indexer_manager.GetBytes(
			indexerevents.NewStatefulOrderRemovalEvent(
				msg.OrderId,
				indexershared.OrderRemovalReason_ORDER_REMOVAL_REASON_USER_CANCELED,
			),
		),
	)

	return nil
}
