package keeper

import (
	"context"

	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	indexerevents "github.com/dydxprotocol/v4/indexer/events"
	"github.com/dydxprotocol/v4/indexer/indexer_manager"
	"github.com/dydxprotocol/v4/lib"
	"github.com/dydxprotocol/v4/lib/metrics"
	"github.com/dydxprotocol/v4/x/clob/types"
)

// PlaceOrder is the entry point for stateful `MsgPlaceOrder` messages executed in `runMsgs` during `DeliverTx`.
// This handler is only invoked for stateful orders due to the filtering logic in the mempool in our CometBFT fork.
// TODO (CLOB-646) - Support stateful order replacements.
func (k msgServer) PlaceOrder(goCtx context.Context, msg *types.MsgPlaceOrder) (*types.MsgPlaceOrderResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// 1. Ensure the order is not a Short-Term order.
	order := msg.GetOrder()
	order.MustBeStatefulOrder()

	// 2. Return an error if an associated cancellation already exists in the current block.
	processProposerMatchesEvents := k.Keeper.GetProcessProposerMatchesEvents(ctx)
	cancelledOrderIds := lib.SliceToSet(processProposerMatchesEvents.PlacedStatefulCancellations)
	if _, found := cancelledOrderIds[order.GetOrderId()]; found {
		return nil, sdkerrors.Wrapf(
			types.ErrStatefulOrderPreviouslyCancelled,
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

	// 6. Emit the new stateful order placement indexer event.
	k.Keeper.GetIndexerEventManager().AddTxnEvent(
		ctx,
		indexerevents.SubtypeStatefulOrder,
		indexer_manager.GetB64EncodedEventMessage(
			indexerevents.NewStatefulOrderPlacementEvent(
				order,
			),
		),
	)

	// 5. Add the newly-placed stateful order to `ProcessProposerMatchesEvents` for use in `PrepareCheckState`.
	processProposerMatchesEvents.PlacedStatefulOrders = append(
		processProposerMatchesEvents.PlacedStatefulOrders,
		order,
	)
	k.Keeper.MustSetProcessProposerMatchesEvents(
		ctx,
		processProposerMatchesEvents,
	)

	telemetry.IncrCounter(1, types.ModuleName, metrics.StatefulOrderMsgHandlerSuccess, metrics.Count)

	return &types.MsgPlaceOrderResponse{}, nil
}
