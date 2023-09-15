package keeper

import (
	"context"
	errorsmod "cosmossdk.io/errors"

	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"
	indexerevents "github.com/dydxprotocol/v4-chain/protocol/indexer/events"
	"github.com/dydxprotocol/v4-chain/protocol/indexer/indexer_manager"
	"github.com/dydxprotocol/v4-chain/protocol/lib"
	"github.com/dydxprotocol/v4-chain/protocol/lib/metrics"
	"github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
)

// PlaceOrder is the entry point for stateful `MsgPlaceOrder` messages executed in `runMsgs` during `DeliverTx`.
// This handler is only invoked for stateful orders due to the filtering logic in the mempool in our CometBFT fork.
// TODO (CLOB-646) - Support stateful order replacements.
func (k msgServer) PlaceOrder(goCtx context.Context, msg *types.MsgPlaceOrder) (*types.MsgPlaceOrderResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// 1. Ensure the order is not a Short-Term order.
	order := msg.GetOrder()
	order.MustBeStatefulOrder()

	// 2. Return an error if an associated cancellation or removal already exists in the current block.
	processProposerMatchesEvents := k.Keeper.GetProcessProposerMatchesEvents(ctx)
	cancelledOrderIds := lib.SliceToSet(processProposerMatchesEvents.PlacedStatefulCancellationOrderIds)
	if _, found := cancelledOrderIds[order.GetOrderId()]; found {
		return nil, errorsmod.Wrapf(
			types.ErrStatefulOrderPreviouslyCancelled,
			"PlaceOrder: order (%+v)",
			order,
		)
	}
	removedOrderIds := lib.SliceToSet(processProposerMatchesEvents.RemovedStatefulOrderIds)
	if _, found := removedOrderIds[order.GetOrderId()]; found {
		return nil, errorsmod.Wrapf(
			types.ErrStatefulOrderPreviouslyRemoved,
			"PlaceOrder: order (%+v)",
			order,
		)
	}

	// 3. Place the order on the ClobKeeper which is responsible for:
	//   - stateful order validation.
	//   - collateralization check.
	//   - writing the order to state and the memstore.
	if err := k.Keeper.PlaceStatefulOrder(ctx, msg); err != nil {
		return nil, err
	}

	// 4. Emit the new order placement indexer event.
	if order.IsConditionalOrder() {
		k.Keeper.GetIndexerEventManager().AddTxnEvent(
			ctx,
			indexerevents.SubtypeStatefulOrder,
			indexer_manager.GetB64EncodedEventMessage(
				indexerevents.NewConditionalOrderPlacementEvent(
					order,
				),
			),
		)
		processProposerMatchesEvents.PlacedConditionalOrderIds = append(
			processProposerMatchesEvents.PlacedConditionalOrderIds,
			order.OrderId,
		)
	} else {
		k.Keeper.GetIndexerEventManager().AddTxnEvent(
			ctx,
			indexerevents.SubtypeStatefulOrder,
			indexer_manager.GetB64EncodedEventMessage(
				indexerevents.NewLongTermOrderPlacementEvent(
					order,
				),
			),
		)
		processProposerMatchesEvents.PlacedLongTermOrderIds = append(
			processProposerMatchesEvents.PlacedLongTermOrderIds,
			order.OrderId,
		)
	}
	// 5. Add the newly-placed stateful order to `ProcessProposerMatchesEvents` for use in `PrepareCheckState`.
	k.Keeper.MustSetProcessProposerMatchesEvents(
		ctx,
		processProposerMatchesEvents,
	)

	telemetry.IncrCounterWithLabels(
		[]string{types.ModuleName, metrics.StatefulOrderMsgHandlerSuccess, metrics.Count},
		1,
		order.GetOrderLabels(),
	)

	return &types.MsgPlaceOrderResponse{}, nil
}
