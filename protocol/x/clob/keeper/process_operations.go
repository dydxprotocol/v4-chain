package keeper

import (
	"errors"
	"fmt"
	"time"

	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	indexerevents "github.com/dydxprotocol/v4/indexer/events"
	indexer_manager "github.com/dydxprotocol/v4/indexer/indexer_manager"
	"github.com/dydxprotocol/v4/lib"
	"github.com/dydxprotocol/v4/lib/metrics"
	"github.com/dydxprotocol/v4/x/clob/memclob"
	"github.com/dydxprotocol/v4/x/clob/types"
)

// ProcessProposerOperations updates on-chain state given an operations queue
// representing matches that occurred in the previous block.
// A new copy of the memclob will be generated to replay the operations on, which will
// be discarded at the end of DeliverTx.
func (k Keeper) ProcessProposerOperations(
	ctx sdk.Context,
	operations []types.Operation,
	addToOrderbookCollatCheckOrderHashesSet map[types.OrderHash]bool,
) error {
	// This function should be only run in DeliverTx mode.
	lib.AssertDeliverTxMode(ctx)
	defer telemetry.ModuleMeasureSince(types.ModuleName, time.Now(), metrics.ProcessOperations)

	// Perform stateful validation of Operations.
	if err := k.performOperationProcessing(
		ctx,
		operations,
		addToOrderbookCollatCheckOrderHashesSet,
	); err != nil {
		return err
	}

	return nil
}

// performOperationProcessing performs all stateful / some stateless validations on the
// provided list of operations. A list of user operations will be generated from the operations
// queue. User operations are order placements, cancellations, and liquidations.
// These user operations will be replayed on a new memclob based off of deliverState.
// After the operations are replayed on the memclob, the temporary memclob's operation queue
// will be compared with the msg's operation queue. If the operation queues match,
// that implies the underlying state is the same. If not, we will error out.
// The memclob is then discarded at the end and the state updates are then persisted to deliverState.
// Note that `addToOrderbookCollatCheckOrderHashesSet` is necessary to determine which orders we should
// perform the add-to-orderbook collateralization check on.
func (k Keeper) performOperationProcessing(
	ctx sdk.Context,
	operations []types.Operation,
	addToOrderbookCollatCheckOrderHashesSet map[types.OrderHash]bool,
) error {
	// Create a temporary memclob to replay operations on.
	replayMemclob := memclob.NewMemClobPriceTimePriority(false)
	replayMemclob.SetClobKeeper(k)

	// Initialize the temporary memclob with clobPairs.
	clobPairs := k.GetAllClobPair(ctx)
	for _, clobPair := range clobPairs {
		// Create the corresponding orderbook in the memclob.
		replayMemclob.CreateOrderbook(
			ctx,
			clobPair,
		)
	}

	// Perform initial stateless checks and generate a user operation queue with state information
	// necessary to place orders on temporary memclob.
	replayableOperations, err := k.getReplayableOperationsFromOpQueue(
		ctx,
		operations,
		addToOrderbookCollatCheckOrderHashesSet,
	)
	if err != nil {
		return err
	}

	// replay all user operations on the temporary memclob.
	for _, operation := range replayableOperations {
		err := operation.replay(ctx, k, replayMemclob)
		if err != nil {
			return err
		}
	}

	memclobOperations := replayMemclob.GetOperations(ctx)
	if len(operations) != len(memclobOperations) {
		minOpLength := lib.Min(len(operations), len(memclobOperations))

		for i := 0; i < minOpLength; i++ {
			if !types.IsEqual(&operations[i], &memclobOperations[i]) {
				return sdkerrors.Wrapf(
					types.ErrOperationsQueueValidationFailure,
					"Mismatched lengths. Proposed operations length: %d. Validation operations length: %d. "+
						"First mismatched index: %d. Proposed operation: %s."+
						"Local operations: %s \nProposer operations: %s",
					len(operations),
					len(memclobOperations),
					i,
					operations[i].GetOperationTextString(),
					types.GetOperationsQueueTextString(memclobOperations),
					types.GetOperationsQueueTextString(operations),
				)
			}
		}
		// TODO(DEC-1772) Make the errors more readable.
		return sdkerrors.Wrapf(
			types.ErrOperationsQueueValidationFailure,
			"Same operations but mismatched lengths. Proposed operations length: %d. Validation operations length: %d."+
				"Local operations: %s \nProposer operations: %s",
			len(operations),
			len(memclobOperations),
			types.GetOperationsQueueTextString(memclobOperations),
			types.GetOperationsQueueTextString(operations),
		)
	}
	for idx, originalOperation := range operations {
		memclobOperation := memclobOperations[idx]
		if !types.IsEqual(&originalOperation, &memclobOperation) {
			return sdkerrors.Wrapf(
				types.ErrOperationsQueueValidationFailure,
				"Mismatched Operation Queues. Length %d, Index %v, Proposed Operation %+v, Memclob Operation %+v. "+
					"Local operations: %s \nProposer operations: %s",
				len(operations),
				idx,
				operations[idx].GetOperationTextString(),
				memclobOperations[idx].GetOperationTextString(),
				types.GetOperationsQueueTextString(memclobOperations),
				types.GetOperationsQueueTextString(operations),
			)
		}
	}

	// Fetch process proposer matches events to write to state.
	processProposerMatchesEvents := k.GenerateProcessProposerMatchesEvents(ctx, operations)
	// Update the memstore with stateful orders placed in this block and list of order Ids filled during this block.
	// During commit, placed stateful orders will be updated in the memclob. All orders that have been fully filled
	// during this block will be removed from the memclob.
	k.MustSetProcessProposerMatchesEvents(
		ctx,
		processProposerMatchesEvents,
	)

	return nil
}

// GenerateProcessProposerMatchesEvents generates a `ProcessProposerMatchesEvents` object from
// an operations queue.
// This function expects the proposed operations to be valid, and does not verify that the `GoodTilBlockTime`
// of order replacement and cancellation is greater than the `GoodTilBlockTime` of the existing order.
func (k Keeper) GenerateProcessProposerMatchesEvents(
	ctx sdk.Context,
	operations []types.Operation,
) types.ProcessProposerMatchesEvents {
	// Set of stateful orders placed in this block, accounting for cancellations and replacements.
	finalPlacedStatefulOrders := make(map[types.OrderId]types.Order, 0)

	// Seen set for filled order ids
	seenOrderIdsFilledInLastBlock := make(map[types.OrderId]struct{}, 0)

	// Collect all stateful placed orders and filled order ids in this block.
	for _, operation := range operations {
		// Add stateful placed order to `placedStatefulOrders`.
		if orderPlacement := operation.GetOrderPlacement(); orderPlacement != nil {
			order := orderPlacement.Order
			if order.IsStatefulOrder() {
				finalPlacedStatefulOrders[order.OrderId] = order
			}
		}
		// Delete the associating stateful orders if it has been canceled in the same block.
		if orderCancellation := operation.GetOrderCancellation(); orderCancellation != nil {
			// Deleting a short-term order is a no-op.
			delete(finalPlacedStatefulOrders, orderCancellation.OrderId)
		}
		if operationMatch := operation.GetMatch(); operationMatch != nil {
			if matchOrders := operationMatch.GetMatchOrders(); matchOrders != nil {
				// For match order, add taker order id to `seenOrderIdsFilledInLastBlock`
				takerOrderId := matchOrders.GetTakerOrderId()
				seenOrderIdsFilledInLastBlock[takerOrderId] = struct{}{}
				// For each fill of a match order, add maker order id to `seenOrderIdsFilledInLastBlock`
				for _, fill := range matchOrders.GetFills() {
					makerOrderId := fill.GetMakerOrderId()
					seenOrderIdsFilledInLastBlock[makerOrderId] = struct{}{}
				}
			}
			// For each fill of a perpetual liquidation match, add maker order id to `seenOrderIdsFilledInLastBlock`
			if perpLiquidationMatch := operationMatch.GetMatchPerpetualLiquidation(); perpLiquidationMatch != nil {
				for _, fill := range perpLiquidationMatch.GetFills() {
					makerOrderId := fill.GetMakerOrderId()
					seenOrderIdsFilledInLastBlock[makerOrderId] = struct{}{}
				}
			}
		}
	}
	orderIds := lib.ConvertMapToSliceOfKeys(seenOrderIdsFilledInLastBlock)
	// Sort for deterministic ordering when writing to memstore.
	types.MustSortAndHaveNoDuplicates(orderIds)

	// Append stateful orders placed in this block to a slice
	// in the order they appear in `operations`.
	orderedPlacedStatefulOrders := make([]types.Order, 0, len(finalPlacedStatefulOrders))
	for _, operation := range operations {
		if orderPlacement := operation.GetOrderPlacement(); orderPlacement != nil {
			order := orderPlacement.Order

			// Short-term orders won't exist in `finalPlacedStatefulOrders` and won't
			// be added to `orderedPlacedStatefulOrders`.
			currentOrder, exists := finalPlacedStatefulOrders[order.OrderId]

			// Double check the order hash to correctly account for order replacements.
			if exists && currentOrder.GetOrderHash() == order.GetOrderHash() {
				orderedPlacedStatefulOrders = append(orderedPlacedStatefulOrders, order)
			}
		}
	}

	return types.ProcessProposerMatchesEvents{
		PlacedStatefulOrders:          orderedPlacedStatefulOrders,
		ExpiredStatefulOrderIds:       []types.OrderId{}, // ExpiredOrderId to be populated in the EndBlocker.
		OrdersIdsFilledInLastBlock:    orderIds,
		OperationsProposedInLastBlock: operations,
		BlockHeight:                   lib.MustConvertIntegerToUint32(ctx.BlockHeight()),
	}
}

// replayableOperation represents a user operation that can be played on the memclob.
// Each user operation is constructed from Msg operations, is standalone and self-encapsulates
// all the state data required to be replayed on a memclob.
//
//	An replayableOperation has a replay method that replays these user operations on the given memclob.
//
//	Types that are valid replayableOperations:
//	*replayableNoOpOperation
//	*replayableMatchPerpetualLiquidationOperation
//	*replayablePlaceOrderOperation
//	*replayableCancelOrderOperation
//	*replayablePreexistingStatefulOrder
type replayableOperation interface {
	replay(ctx sdk.Context, k Keeper, memclob *memclob.MemClobPriceTimePriority) error
}

// replayableNoOpOperation represents a no-op for the memclob.
// We emit these for OperationQueue operations that do not generate user events,
// such as the `MatchOrders` type.
type replayableNoOpOperation struct{}

// replay performs a NoOpOperation on the memclob.
func (op *replayableNoOpOperation) replay(
	ctx sdk.Context,
	k Keeper,
	memclob *memclob.MemClobPriceTimePriority,
) error {
	return nil
}

// replayablePerpetualLiquidationOperation contains all the stateful information
// required to place a perpetual liquidation on the memclob. It is constructed from the
// the `MatchPerpetualLiquidationNew` type in the operations queue.
type replayablePerpetualLiquidationOperation struct {
	liquidationOrder types.LiquidationOrder
}

// replay places the liquidation order on the memclob.
func (op *replayablePerpetualLiquidationOperation) replay(
	ctx sdk.Context,
	k Keeper,
	memclob *memclob.MemClobPriceTimePriority,
) error {
	_, _, err := k.DeliverTxPlacePerpetualLiquidation(ctx, op.liquidationOrder, memclob)
	return err
}

// replayablePlaceOrderOperation contains all the stateful information
// required to place an order on the memclob. It is constructed from the
// the `Operation_OrderPlacement` type in the Operations queue.
type replayablePlaceOrderOperation struct {
	msgPlaceOrder                    types.MsgPlaceOrder
	performAddToOrderbookCollatCheck bool
}

// replay places an order on the memclob.
func (op *replayablePlaceOrderOperation) replay(
	ctx sdk.Context,
	k Keeper,
	memclob *memclob.MemClobPriceTimePriority,
) error {
	// Branch the multistore before calling `DeliverTxPlaceOrder` so we can revert state updates
	// if placing the order returns an error.
	deliverTxPlaceOrderCtx, writeCache := ctx.CacheContext()

	// If the order was placed successfully and was a stateful order, send an order placement event to the Indexer.
	// Note that if there is an error and `writeCache` is not called below, then the event will never be included.
	placedOrder := op.msgPlaceOrder.Order
	if placedOrder.IsStatefulOrder() {
		k.GetIndexerEventManager().AddTxnEvent(
			deliverTxPlaceOrderCtx,
			indexerevents.SubtypeStatefulOrder,
			indexer_manager.GetB64EncodedEventMessage(
				indexerevents.NewStatefulOrderPlacementEvent(
					op.msgPlaceOrder.Order,
				),
			),
		)
	}

	_, _, err := k.DeliverTxPlaceOrder(
		deliverTxPlaceOrderCtx,
		&op.msgPlaceOrder,
		op.performAddToOrderbookCollatCheck,
		memclob,
	)

	// If the order was placed successfully, commit all updates to state.
	if err == nil {
		writeCache()
	}

	// If this was an expected error, return `nil` to indicate replaying the operation was successful.
	// TODO(DEC-1964) Remove this.
	if errors.Is(err, types.ErrFokOrderCouldNotBeFullyFilled) ||
		errors.Is(err, types.ErrPostOnlyWouldCrossMakerOrder) {
		return nil
	}
	return err
}

// replayablePreexistingStatefulOrder contains all the stateful information
// required to replay a stateful order to a memclob. It is constructed from the
// the Operation_PreexistingStatefulOrder type in the Operations queue.
type replayablePreexistingStatefulOrder struct {
	order types.Order
}

// replay fetches a stateful validation from keeper state and places it on
// the memclob. It returns an error when:
// - stateful order placement cannot be found
// - an error is found when placing order (i.e post-only and crosses the book)
func (op *replayablePreexistingStatefulOrder) replay(
	ctx sdk.Context,
	k Keeper,
	memclob *memclob.MemClobPriceTimePriority,
) error {
	// Use the height of the current block being processed.
	currBlockHeight := lib.MustConvertIntegerToUint32(ctx.BlockHeight())

	_, _, _, err := k.AddPreexistingStatefulOrder(ctx, &op.order, currBlockHeight, memclob)
	return err
}

// replayableCancelOrderOperation contains all the stateful information
// required to cancel an order on the memclob. It is constructed from the
// the Operation_OrderCancellation type in the Operations queue. The cancel
// order can be stateful or short-term.
type replayableCancelOrderOperation struct {
	msgCancelOrder types.MsgCancelOrder
}

// replay cancels an order on the memclob.
func (op *replayableCancelOrderOperation) replay(
	ctx sdk.Context,
	k Keeper,
	memclob *memclob.MemClobPriceTimePriority,
) error {
	err := k.DeliverTxCancelOrder(
		ctx,
		&op.msgCancelOrder,
		memclob,
	)

	if err == nil && op.msgCancelOrder.OrderId.IsStatefulOrder() {
		k.GetIndexerEventManager().AddTxnEvent(
			ctx,
			indexerevents.SubtypeStatefulOrder,
			indexer_manager.GetB64EncodedEventMessage(
				indexerevents.NewStatefulOrderCancelationEvent(
					op.msgCancelOrder.OrderId,
				),
			),
		)
	}

	return err
}

// getReplayableOperationsFromOpQueue takes in the operations queue and generates a list of replayable
// operations we will replay on the temporary memclob. Each replayable operation struct contains all the
// information necessary to perform a replayable operation on the memclob and a replay method to do so.
func (k Keeper) getReplayableOperationsFromOpQueue(
	ctx sdk.Context,
	operations []types.Operation,
	addToOrderbookCollatCheckOrderHashesSet map[types.OrderHash]bool,
) ([]replayableOperation, error) {
	replayableOperations := []replayableOperation{}

	for _, operation := range operations {
		var replayableOp replayableOperation
		var err error

		switch operation.Operation.(type) {
		case *types.Operation_Match:
			operationMatch := operation.GetMatch()
			replayableOp, err = k.processClobMatch(ctx, operationMatch)
		case *types.Operation_OrderPlacement:
			orderPlacement := operation.GetOrderPlacement()
			replayableOp, err = k.processOrderPlacement(
				ctx,
				orderPlacement,
				addToOrderbookCollatCheckOrderHashesSet[orderPlacement.Order.GetOrderHash()],
			)
		case *types.Operation_OrderCancellation:
			orderCancellation := operation.GetOrderCancellation()
			replayableOp, err = k.processOrderCancellation(ctx, orderCancellation)
		case *types.Operation_PreexistingStatefulOrder:
			statefulOrder := operation.GetPreexistingStatefulOrder()
			replayableOp, err = k.processPreexistingStatefulOrder(ctx, statefulOrder)
		default:
			err = fmt.Errorf("operation queue type not implemented yet for operation %+v", operation)
		}

		if err != nil {
			return nil, sdkerrors.Wrapf(
				err,
				"processing failed on operation %+v",
				operation,
			)
		}

		// Append the generated memclob user operation to the queue to be returned.
		replayableOperations = append(replayableOperations, replayableOp)
	}
	return replayableOperations, nil
}

// processPreexistingStatefulOrder takes in a `OrderId` object that represents a stateful order
// and generates a `replayablePreexistingStatefulOrder` object that can be replayed on a memclob.
func (k Keeper) processPreexistingStatefulOrder(
	ctx sdk.Context,
	orderId *types.OrderId,
) (replayableOperation, error) {
	placement, found := k.GetStatefulOrderPlacement(ctx, *orderId)
	if !found {
		return nil, fmt.Errorf("failed to find stateful order placement for order id %+v", orderId)
	}
	statefulOrderOperation := replayablePreexistingStatefulOrder{
		order: placement.Order,
	}
	return &statefulOrderOperation, nil
}

// processOrderCancellation takes in a `MsgCancelOrder` object and generates a
// `replayableCancelOrderOperation` that can be replayed on a memclob.
func (k Keeper) processOrderCancellation(
	ctx sdk.Context,
	cancelOrder *types.MsgCancelOrder,
) (replayableOperation, error) {
	cancelOrderOperation := replayableCancelOrderOperation{
		msgCancelOrder: *cancelOrder,
	}
	return &cancelOrderOperation, nil
}

// processOrderPlacement takes in a `MsgPlaceOrder` object and generates a
// `replayablePlaceOrderOperation` that can be replayed on a memclob.
func (k Keeper) processOrderPlacement(
	ctx sdk.Context,
	placeOrder *types.MsgPlaceOrder,
	performAddOrderbookCollatCheck bool,
) (replayableOperation, error) {
	placeOrderOperation := replayablePlaceOrderOperation{
		msgPlaceOrder:                    *placeOrder,
		performAddToOrderbookCollatCheck: performAddOrderbookCollatCheck,
	}
	return &placeOrderOperation, nil
}

// processClobMatch takes in a `ClobMatch“ object and forwards the object to
// the corresponding Processing function in order to generate a `replayableOperation`
// that can be replayed on a memclob.
func (k Keeper) processClobMatch(
	ctx sdk.Context,
	match *types.ClobMatch,
) (replayableOperation, error) {
	switch match.Match.(type) {
	// Replaying the taker order placement should also replay the match orders event.
	// We do not need to replay this. This operation is for verification purposes.
	case *types.ClobMatch_MatchOrders:
		return &replayableNoOpOperation{}, nil
	case *types.ClobMatch_MatchPerpetualLiquidation:
		perpetualLiquidationMatch := match.GetMatchPerpetualLiquidation()
		return k.processMatchPerpetualLiquidationNew(ctx, perpetualLiquidationMatch)
	// Replaying the liquidation order should also replay the deleveraging event.
	// We do not need to replay this. This operation is for verification purposes.
	case *types.ClobMatch_MatchPerpetualDeleveraging:
		return &replayableNoOpOperation{}, nil
	default:
		panic("PreprocessOperationMatch: Unsupported Clob Match type")
	}
}

// processMatchPerpetualLiquidationNew takes in a `MatchPerpetualLiquidation` object
// and generates a `replayablePerpetualLiquidationOperation` object that can be replayed
// on a memclob.
// TODO(DEC-1653) Remove New suffix after name clash is gone when old ProcessMatches code is deprecated.
func (k Keeper) processMatchPerpetualLiquidationNew(
	ctx sdk.Context,
	perpetualLiquidationMatch *types.MatchPerpetualLiquidation,
) (replayableOperation, error) {
	takerOrder, err := k.ConstructTakerOrderFromMatchPerpetualLiquidationNew(ctx, perpetualLiquidationMatch)
	if err != nil {
		return nil, err
	}

	memclobPerpLiquidation := replayablePerpetualLiquidationOperation{
		liquidationOrder: *takerOrder,
	}

	return &memclobPerpLiquidation, nil
}
