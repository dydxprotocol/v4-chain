package keeper

import (
	"context"

	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"
	indexerevents "github.com/dydxprotocol/v4/indexer/events"
	"github.com/dydxprotocol/v4/indexer/indexer_manager"
	"github.com/dydxprotocol/v4/lib/metrics"
	"github.com/dydxprotocol/v4/x/clob/types"
)

// CancelOrder performs order cancellation functionality for stateful orders.
func (m msgServer) CancelOrder(
	goCtx context.Context,
	msg *types.MsgCancelOrder,
) (*types.MsgCancelOrderResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// 1. If this is a Short-Term order, panic.
	msg.OrderId.MustBeStatefulOrder()

	// 2. Cancel the order on the ClobKeeper which is responsible for:
	//   - stateful cancellation validation.
	//   - removing the order from state and the memstore.
	if err := m.Keeper.CancelStatefulOrder(ctx, msg); err != nil {
		return nil, err
	}

	// 3. Update `ProcessProposerMatchesEvents` with the new stateful order cancellation.
	processProposerMatchesEvents := m.Keeper.GetProcessProposerMatchesEvents(ctx)

	processProposerMatchesEvents.PlacedStatefulCancellations = append(
		processProposerMatchesEvents.PlacedStatefulCancellations,
		msg.OrderId,
	)

	m.Keeper.MustSetProcessProposerMatchesEvents(ctx, processProposerMatchesEvents)

	// 4. Add the relevant on-chain Indexer event for the cancellation.
	m.Keeper.GetIndexerEventManager().AddTxnEvent(
		ctx,
		indexerevents.SubtypeStatefulOrder,
		indexer_manager.GetB64EncodedEventMessage(
			indexerevents.NewStatefulOrderCancelationEvent(
				msg.OrderId,
			),
		),
	)

	telemetry.IncrCounter(1, types.ModuleName, metrics.StatefulCancellationMsgHandlerSuccess, metrics.Count)

	return &types.MsgCancelOrderResponse{}, nil
}
