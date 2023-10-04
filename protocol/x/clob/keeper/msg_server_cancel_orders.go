package keeper

import (
	"context"

	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"
	indexerevents "github.com/dydxprotocol/v4-chain/protocol/indexer/events"
	"github.com/dydxprotocol/v4-chain/protocol/indexer/indexer_manager"
	indexershared "github.com/dydxprotocol/v4-chain/protocol/indexer/shared"
	errorlib "github.com/dydxprotocol/v4-chain/protocol/lib/error"
	"github.com/dydxprotocol/v4-chain/protocol/lib/metrics"
	"github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
)

// CancelOrder performs order cancellation functionality for stateful orders.
func (m msgServer) CancelOrder(
	goCtx context.Context,
	msg *types.MsgCancelOrder,
) (resp *types.MsgCancelOrderResponse, err error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	defer func() {
		if err != nil {
			errorlib.LogErrorWithBlockHeight(ctx.Logger(), err, ctx.BlockHeight(), metrics.DeliverTx)
		}
	}()

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

	processProposerMatchesEvents.PlacedStatefulCancellationOrderIds = append(
		processProposerMatchesEvents.PlacedStatefulCancellationOrderIds,
		msg.OrderId,
	)

	m.Keeper.MustSetProcessProposerMatchesEvents(ctx, processProposerMatchesEvents)

	// 4. Add the relevant on-chain Indexer event for the cancellation.
	m.Keeper.GetIndexerEventManager().AddTxnEvent(
		ctx,
		indexerevents.SubtypeStatefulOrder,
		indexer_manager.GetB64EncodedEventMessage(
			indexerevents.NewStatefulOrderRemovalEvent(
				msg.OrderId,
				indexershared.OrderRemovalReason_ORDER_REMOVAL_REASON_USER_CANCELED,
			),
		),
		indexerevents.StatefulOrderEventVersion,
		indexer_manager.GetBytes(
			indexerevents.NewStatefulOrderRemovalEvent(
				msg.OrderId,
				indexershared.OrderRemovalReason_ORDER_REMOVAL_REASON_USER_CANCELED,
			),
		),
	)

	telemetry.IncrCounterWithLabels(
		[]string{types.ModuleName, metrics.StatefulCancellationMsgHandlerSuccess, metrics.Count},
		1,
		msg.OrderId.GetOrderIdLabels(),
	)

	return &types.MsgCancelOrderResponse{}, nil
}
