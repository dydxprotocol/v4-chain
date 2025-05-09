package keeper

import (
	"context"
	"errors"

	errorsmod "cosmossdk.io/errors"

	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"
	indexerevents "github.com/dydxprotocol/v4-chain/protocol/indexer/events"
	"github.com/dydxprotocol/v4-chain/protocol/indexer/indexer_manager"
	"github.com/dydxprotocol/v4-chain/protocol/lib"
	"github.com/dydxprotocol/v4-chain/protocol/lib/log"
	"github.com/dydxprotocol/v4-chain/protocol/lib/metrics"
	"github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
)

// PlaceOrder is the entry point for stateful `MsgPlaceOrder` messages executed in `runMsgs` during `DeliverTx`.
// This handler is only invoked for stateful orders due to the filtering logic in the mempool in our CometBFT fork.
// TODO (CLOB-646) - Support stateful order replacements.
func (k msgServer) PlaceOrder(goCtx context.Context, msg *types.MsgPlaceOrder) (
	resp *types.MsgPlaceOrderResponse,
	err error,
) {
	ctx := lib.UnwrapSDKContext(goCtx, types.ModuleName)

	if err := k.Keeper.HandleMsgPlaceOrder(ctx, msg, false); err != nil {
		return nil, err
	}

	return &types.MsgPlaceOrderResponse{}, nil
}

// HandleMsgPlaceOrder handles a MsgPlaceOrder by
// 1. persisting the placement on chain.
// 2. updating ProcessProposerMatchesEvents with the new stateful order placement.
// 3. adding order placement on-chain indexer event.
// Various logs, metrics, and validations are skipped for orders internal to the protocol.
func (k Keeper) HandleMsgPlaceOrder(
	ctx sdk.Context,
	msg *types.MsgPlaceOrder,
	isInternalOrder bool,
) (err error) {
	lib.AssertDeliverTxMode(ctx)

	// Attach various logging tags relative to this request. These should be static with no changes.
	ctx = log.AddPersistentTagsToLogger(ctx,
		log.Module, log.Clob,
		log.ProposerConsAddress, sdk.ConsAddress(ctx.BlockHeader().ProposerAddress),
		log.Callback, lib.TxMode(ctx),
		log.BlockHeight, ctx.BlockHeight(),
		log.Handler, log.PlaceOrder,
	)

	if !isInternalOrder {
		defer func() {
			metrics.IncrSuccessOrErrorCounter(
				err,
				types.ModuleName,
				metrics.PlaceOrder,
				metrics.DeliverTx,
				msg.Order.GetOrderLabels()...,
			)
			if err != nil {
				if errors.Is(err, types.ErrStatefulOrderCollateralizationCheckFailed) {
					telemetry.IncrCounterWithLabels(
						[]string{
							types.ModuleName,
							metrics.PlaceOrder,
							metrics.CollateralizationCheckFailed,
						},
						1,
						msg.Order.GetOrderLabels(),
					)
					log.InfoLog(ctx, "Place Order Expected Error", log.Error, err)
					return
				}
				log.ErrorLogWithError(ctx, "Error placing order", err)
			}
		}()
	}

	// 1. Ensure the order is not a Short-Term order.
	order := msg.GetOrder()
	order.MustBeStatefulOrder()

	// 2. Return an error if an associated cancellation or removal already exists in the current block.
	processProposerMatchesEvents := k.GetProcessProposerMatchesEvents(ctx)
	cancelledOrderIds := lib.UniqueSliceToSet(k.GetDeliveredCancelledOrderIds(ctx))
	if _, found := cancelledOrderIds[order.GetOrderId()]; found {
		return errorsmod.Wrapf(
			types.ErrStatefulOrderPreviouslyCancelled,
			"PlaceOrder: order (%+v)",
			order,
		)
	}
	removedOrderIds := lib.UniqueSliceToSet(processProposerMatchesEvents.RemovedStatefulOrderIds)
	if _, found := removedOrderIds[order.GetOrderId()]; found {
		return errorsmod.Wrapf(
			types.ErrStatefulOrderPreviouslyRemoved,
			"PlaceOrder: order (%+v)",
			order,
		)
	}

	// 3. Place the order on the ClobKeeper which is responsible for:
	//   - stateful order validation.
	//   - collateralization check.
	//   - writing the order to state and the memstore.
	if err := k.PlaceStatefulOrder(ctx, msg, isInternalOrder); err != nil {
		return err
	}

	// 4. Emit the new order placement indexer event.
	if order.IsConditionalOrder() {
		k.GetIndexerEventManager().AddTxnEvent(
			ctx,
			indexerevents.SubtypeStatefulOrder,
			indexerevents.StatefulOrderEventVersion,
			indexer_manager.GetBytes(
				indexerevents.NewConditionalOrderPlacementEvent(
					order,
				),
			),
		)
		k.AddDeliveredConditionalOrderId(
			ctx,
			order.OrderId,
		)
	} else if order.IsTwapOrder() {
		k.GetIndexerEventManager().AddTxnEvent(
			ctx,
			indexerevents.SubtypeStatefulOrder,
			indexerevents.StatefulOrderEventVersion,
			indexer_manager.GetBytes(
				indexerevents.NewTwapOrderPlacementEvent(
					order,
				),
			),
		)
	} else {
		k.GetIndexerEventManager().AddTxnEvent(
			ctx,
			indexerevents.SubtypeStatefulOrder,
			indexerevents.StatefulOrderEventVersion,
			indexer_manager.GetBytes(
				indexerevents.NewLongTermOrderPlacementEvent(
					order,
				),
			),
		)
		k.AddDeliveredLongTermOrderId(
			ctx,
			order.OrderId,
		)
	}

	return nil
}
