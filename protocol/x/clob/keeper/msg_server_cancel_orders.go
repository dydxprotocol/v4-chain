package keeper

import (
	"context"
	"errors"

	errorsmod "cosmossdk.io/errors"

	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"
	indexerevents "github.com/dydxprotocol/v4-chain/protocol/indexer/events"
	"github.com/dydxprotocol/v4-chain/protocol/indexer/indexer_manager"
	indexershared "github.com/dydxprotocol/v4-chain/protocol/indexer/shared"
	"github.com/dydxprotocol/v4-chain/protocol/lib"
	errorlib "github.com/dydxprotocol/v4-chain/protocol/lib/error"
	"github.com/dydxprotocol/v4-chain/protocol/lib/metrics"
	"github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
)

// CancelOrder performs order cancellation functionality for stateful orders.
func (k msgServer) CancelOrder(
	goCtx context.Context,
	msg *types.MsgCancelOrder,
) (resp *types.MsgCancelOrderResponse, err error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

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
				processProposerMatchesEvents := k.Keeper.GetProcessProposerMatchesEvents(ctx)
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
					k.Keeper.Logger(ctx).Info(
						err.Error(),
						metrics.BlockHeight, ctx.BlockHeight(),
						metrics.Handler, "CancelOrder",
						metrics.Callback, metrics.DeliverTx,
						metrics.Msg, msg,
					)
					return
				}
			}
			errorlib.LogDeliverTxError(k.Keeper.Logger(ctx), err, ctx.BlockHeight(), "CancelOrder", msg)
		}
	}()

	// 1. If this is a Short-Term order, panic.
	msg.OrderId.MustBeStatefulOrder()

	// 2. Cancel the order on the ClobKeeper which is responsible for:
	//   - stateful cancellation validation.
	//   - removing the order from state and the memstore.
	if err := k.Keeper.CancelStatefulOrder(ctx, msg); err != nil {
		return nil, err
	}

	// 3. Update `ProcessProposerMatchesEvents` with the new stateful order cancellation.
	processProposerMatchesEvents := k.Keeper.GetProcessProposerMatchesEvents(ctx)

	processProposerMatchesEvents.PlacedStatefulCancellationOrderIds = append(
		processProposerMatchesEvents.PlacedStatefulCancellationOrderIds,
		msg.OrderId,
	)

	k.Keeper.MustSetProcessProposerMatchesEvents(ctx, processProposerMatchesEvents)

	// 4. Add the relevant on-chain Indexer event for the cancellation.
	k.Keeper.GetIndexerEventManager().AddTxnEvent(
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

	return &types.MsgCancelOrderResponse{}, nil
}
