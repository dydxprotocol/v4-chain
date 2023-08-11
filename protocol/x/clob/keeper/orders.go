package keeper

import (
	"fmt"
	"math"
	"math/big"
	"time"

	"github.com/cometbft/cometbft/crypto/tmhash"
	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/dydxprotocol/v4/indexer/msgsender"
	"github.com/dydxprotocol/v4/indexer/off_chain_updates"
	"github.com/dydxprotocol/v4/lib"
	"github.com/dydxprotocol/v4/lib/metrics"
	"github.com/dydxprotocol/v4/x/clob/types"
	satypes "github.com/dydxprotocol/v4/x/subaccounts/types"
)

func (k Keeper) GetOperations(ctx sdk.Context) *types.MsgProposedOperations {
	operationsQueue := k.MemClob.GetOperations(ctx)
	ordersWithAddToOrderbookCollatCheck := k.MemClob.GetOrdersWithAddToOrderbookCollatCheck(ctx)

	addToOrderbookCollatCheckOrderHashes := make([][]byte, len(ordersWithAddToOrderbookCollatCheck))
	for i, b := range ordersWithAddToOrderbookCollatCheck {
		copy := b
		addToOrderbookCollatCheckOrderHashes[i] = copy[:]
	}

	return &types.MsgProposedOperations{
		OperationsQueue:                      operationsQueue,
		AddToOrderbookCollatCheckOrderHashes: addToOrderbookCollatCheckOrderHashes,
	}
}

// DeliverTxCancelOrder removes an order by `OrderId` (if it exists) from all order-related data structures
// in the specified memclob. As well, DeliverTxCancelOrder adds (or updates) a cancel to the desired
// `goodTilBlock` in the memclob. If a cancel already exists for this order with a lower `goodTilBlock`,
// the cancel is updated to the new `goodTilBlock`. This method is meant to be used in the DeliverTx consensus flow,
// where we apply the cancel order on a new, specified memclob.
//
// An error will be returned if any of the following conditions are true:
// - The cancel's `GoodTilblock` is less than or equal to the prev block height.
// - The cancel's `GoodTilblock` is greater than the sum of the prev block height and `ShortBlockWindow`.
// - The memclob itself returns an error.
func (k Keeper) DeliverTxCancelOrder(
	ctx sdk.Context,
	msgCancelOrder *types.MsgCancelOrder,
	memclob types.MemClob,
) error {
	lib.AssertDeliverTxMode(ctx)
	// Use the height of the previously validated block.
	currBlockHeight := lib.MustConvertIntegerToUint32(ctx.BlockHeight())
	return k.cancelOrder(ctx, msgCancelOrder, currBlockHeight, memclob)
}

// CheckTxCancelOrder removes an order by `OrderId` (if it exists) from all order-related data structures
// in the memclob. As well, CheckTxCancelOrder adds (or updates) a cancel to the desired `goodTilBlock` in the memclob.
// If a cancel already exists for this order with a lower `goodTilBlock`, the cancel is updated to the
// new `goodTilBlock`. This method is meant to be used in the CheckTx flow. It uses the next block height and the keeper
// memclob.
//
// An error will be returned if any of the following conditions are true:
// - The cancel's `GoodTilblock` is less than or equal to the next block height.
// - The cancel's `GoodTilblock` is greater than the sum of the next block height and `ShortBlockWindow`.
// - The memclob itself returns an error.
func (k Keeper) CheckTxCancelOrder(
	ctx sdk.Context,
	msgCancelOrder *types.MsgCancelOrder,
) error {
	lib.AssertCheckTxMode(ctx)
	// Note that we add `+1` here to account for the fact that `ctx.BlockHeight()` is technically the
	// previously mined block, not the next block that will be proposed. This is due to the fact that
	// this function is only ever called during `CheckTx`.
	nextBlockHeight := lib.MustConvertIntegerToUint32(ctx.BlockHeight() + 1)
	return k.cancelOrder(ctx, msgCancelOrder, nextBlockHeight, k.MemClob)
}

// cancelOrder contains shared logic for `CheckTxCancelOrder` and `DeliverTxCancelOrder`.
// cancelOrder removes an order by `OrderId` (if it exists) from all order-related data structures
// in the memclob. As well, cancelOrder adds (or updates) a cancel to the desired `goodTilBlock` in the memclob.
// If a cancel already exists for this order with a lower `goodTilBlock`, the cancel is updated to the
// new `goodTilBlock`.
//
// An error will be returned if any of the following conditions are true:
// - The cancel's `GoodTilblock` is less than or equal to the block height.
// - The cancel's `GoodTilblock` is greater than the sum of the block height and `ShortBlockWindow`.
// - The memclob itself returns an error.
func (k Keeper) cancelOrder(
	ctx sdk.Context,
	msgCancelOrder *types.MsgCancelOrder,
	blockHeight uint32,
	memclob types.MemClob,
) error {
	defer telemetry.ModuleMeasureSince(types.ModuleName, time.Now(), metrics.CancelOrder, metrics.Latency)
	telemetry.IncrCounter(1, types.ModuleName, metrics.CancelOrder, metrics.Count)
	orderIdToCancel := msgCancelOrder.GetOrderId()

	if err := msgCancelOrder.ValidateBasic(); err != nil {
		return err
	}

	// Perform all stateful validation on the order. Order may be stateful or short term.
	if err := k.PerformOrderCancellationStatefulValidation(ctx, msgCancelOrder, blockHeight); err != nil {
		return err
	}

	if orderIdToCancel.IsStatefulOrder() {
		// Fetch the order from state. Return an error if it doesn't exist.
		statefulOrderPlacement, found := k.GetStatefulOrderPlacement(ctx, orderIdToCancel)
		if !found {
			return sdkerrors.Wrapf(
				types.ErrStatefulOrderDoesNotExist,
				"CancelOrder: order id to cancel %+v does not exist in state",
				orderIdToCancel,
			)
		}
		statefulOrderToBeCancelled := statefulOrderPlacement.GetOrder()
		// Remove the order from state, time slice, and order fill amounts.
		k.MustRemoveStatefulOrder(ctx, statefulOrderToBeCancelled.MustGetUnixGoodTilBlockTime(), orderIdToCancel)
	}

	// Update in-memory orderbook to remove order.
	offchainUpdates, err := memclob.CancelOrder(ctx, msgCancelOrder)

	if err == nil {
		k.sendOffchainMessagesWithTxHash(
			offchainUpdates,
			tmhash.Sum(ctx.TxBytes()),
			metrics.SendCancelOrderOffchainUpdates,
		)
	}
	return err
}

// DeliverTxPlaceOrder places an order on the corresponding orderbook, and performs matching if placing the
// order causes an overlap. This function will return the result of calling `PlaceOrder` on the
// specified memclob. This method is meant to be used in the DeliverTx consensus flow,
// where we apply the place order on a new, specified memclob.
//
// An error will be returned if any of the following conditions are true:
// - Standard stateful validation fails.
// - The memclob itself returns an error.
func (k Keeper) DeliverTxPlaceOrder(
	ctx sdk.Context,
	msg *types.MsgPlaceOrder,
	performAddToOrderbookCollatCheck bool,
	memclob types.MemClob,
) (
	orderSizeOptimisticallyFilledFromMatchingQuantums satypes.BaseQuantums,
	orderStatus types.OrderStatus,
	err error,
) {
	lib.AssertDeliverTxMode(ctx)
	// Use the height of the current block being processed.
	currBlockHeight := lib.MustConvertIntegerToUint32(ctx.BlockHeight())
	return k.placeOrder(ctx, msg, currBlockHeight, performAddToOrderbookCollatCheck, memclob)
}

// CheckTxPlaceOrder places an order on the corresponding orderbook, and performs matching if placing the
// order causes an overlap. This function will return the result of calling `PlaceOrder` on the
// keeper's memclob. This method is meant to be used in the CheckTx flow. It uses the next block height and the keeper
// memclob.
//
// An error will be returned if any of the following conditions are true:
// - Standard stateful validation fails.
// - The memclob itself returns an error.
func (k Keeper) CheckTxPlaceOrder(
	ctx sdk.Context,
	msg *types.MsgPlaceOrder,
) (
	orderSizeOptimisticallyFilledFromMatchingQuantums satypes.BaseQuantums,
	orderStatus types.OrderStatus,
	err error,
) {
	lib.AssertCheckTxMode(ctx)
	nextBlockHeight := lib.MustConvertIntegerToUint32(ctx.BlockHeight() + 1)

	return k.placeOrder(ctx, msg, nextBlockHeight, true, k.MemClob)
}

// ReplayPlaceOrder returns the result of calling `PlaceOrder` on the memclob.
// This method does not forward events directly to indexer, but instead returns
// them in the form of `OffchainUpdates`. This method is meant to be used in the
// `ReplayOperations` flow, where we replay Short-Term and newly-played stateful
// orders back onto the memclob.
//
// An error will be returned if any of the following conditions are true:
// - Standard stateful validation fails.
// - The memclob itself returns an error.
func (k Keeper) ReplayPlaceOrder(
	ctx sdk.Context,
	msg *types.MsgPlaceOrder,
) (
	orderSizeOptimisticallyFilledFromMatchingQuantums satypes.BaseQuantums,
	orderStatus types.OrderStatus,
	offchainUpdates *types.OffchainUpdates,
	err error,
) {
	order := msg.GetOrder()

	// Use the height of the next block. Check if this order would be valid if it were included
	// in the next block height, not in the block that was already committed.
	nextBlockHeight := lib.MustConvertIntegerToUint32(ctx.BlockHeight() + 1)

	// Perform stateful validation.
	err = k.PerformStatefulOrderValidation(ctx, &order, nextBlockHeight, false)
	if err != nil {
		return 0, 0, nil, err
	}

	// Place the order on the memclob and return the result.
	orderSizeOptimisticallyFilledFromMatchingQuantums, orderStatus, offchainUpdates, err = k.MemClob.PlaceOrder(
		ctx,
		msg.Order,
		true,
	)

	// If there wasn't an error from placing the order and it's a stateful order, write the order
	// placement to state with the next available transaction index.
	if err == nil && order.IsStatefulOrder() {
		k.SetStatefulOrderPlacement(ctx, order, nextBlockHeight)

		// TODO(DEC-1238): Ensure in the case of stateful order replacements, that the old time slice
		// entry is deleted properly.
		k.MustAddOrderToStatefulOrdersTimeSlice(
			ctx,
			order.MustGetUnixGoodTilBlockTime(),
			order.GetOrderId(),
		)
	}

	return orderSizeOptimisticallyFilledFromMatchingQuantums, orderStatus, offchainUpdates, err
}

// placeOrder contains shared logic for `CheckTxPlaceOrder` and `DeliverTxPlaceOrder`. It performs
// matching if placing the order causes an overlap. This function will return the result of
// calling `PlaceOrder` on the specified memclob.
//
// An error will be returned if any of the following conditions are true:
// - Standard stateful validation fails.
// - The memclob itself returns an error.
func (k Keeper) placeOrder(
	ctx sdk.Context,
	msg *types.MsgPlaceOrder,
	blockHeight uint32,
	performAddToOrderbookCollatCheck bool,
	memclob types.MemClob,
) (
	orderSizeOptimisticallyFilledFromMatchingQuantums satypes.BaseQuantums,
	orderStatus types.OrderStatus,
	err error,
) {
	order := msg.GetOrder()
	orderLabels := order.GetOrderLabels()
	defer telemetry.ModuleMeasureSince(types.ModuleName, time.Now(), metrics.PlaceOrder, metrics.Latency)
	defer func() {
		telemetry.IncrCounterWithLabels(
			[]string{types.ModuleName, metrics.PlaceOrder, metrics.Count},
			1,
			orderLabels,
		)
		if err != nil {
			telemetry.IncrCounterWithLabels(
				[]string{types.ModuleName, metrics.PlaceOrder, metrics.Rejected},
				1,
				orderLabels,
			)
		}

		// Add seen place order to in-memory map for metrics purposes. Only add if we passed all
		// validation and the memclob attempted to place it.
		if err == nil {
			k.AddSeenPlaceOrder(ctx, *msg)
		}
	}()

	// Perform stateful validation.
	err = k.PerformStatefulOrderValidation(ctx, &order, blockHeight, false)
	if err != nil {
		return 0, 0, err
	}

	// Place the order on the memclob and return the result.
	orderSizeOptimisticallyFilledFromMatchingQuantums, orderStatus, offchainUpdates, err := memclob.PlaceOrder(
		ctx,
		msg.Order,
		performAddToOrderbookCollatCheck,
	)

	// If there wasn't an error from placing the order and it's a stateful order, write the order
	// placement to state with the next available transaction index.
	if err == nil && order.IsStatefulOrder() {
		k.SetStatefulOrderPlacement(ctx, order, blockHeight)

		// TODO(DEC-1238): Ensure in the case of stateful order replacements, that the old time slice
		// entry is deleted properly.
		k.MustAddOrderToStatefulOrdersTimeSlice(
			ctx,
			order.MustGetUnixGoodTilBlockTime(),
			order.GetOrderId(),
		)
	}

	// Send off-chain updates generated from placing the order. `SendOffchainData` enqueues the
	// the messages to be sent in a channel and should be non-blocking.
	// Off-chain update messages should be only be returned if the `IndexerMessageSender`
	// is enabled (`msgSender.Enabled()` returns true).
	k.sendOffchainMessagesWithTxHash(
		offchainUpdates,
		tmhash.Sum(ctx.TxBytes()),
		metrics.SendPlaceOrderOffchainUpdates,
	)

	if orderSizeOptimisticallyFilledFromMatchingQuantums > 0 {
		telemetry.IncrCounterWithLabels(
			[]string{types.ModuleName, metrics.PlaceOrder, metrics.Matched},
			1,
			orderLabels,
		)
	}

	return orderSizeOptimisticallyFilledFromMatchingQuantums, orderStatus, err
}

// AddPreexistingStatefulOrder performs stateful validation on an order and adds it to the specified memclob.
// This function does not add the order into state, since it is assumed to be preexisting. Function panics
// if the specified order is not stateful.
func (k Keeper) AddPreexistingStatefulOrder(
	ctx sdk.Context,
	order *types.Order,
	blockHeight uint32,
	memclob types.MemClob,
) (
	orderSizeOptimisticallyFilledFromMatchingQuantums satypes.BaseQuantums,
	orderStatus types.OrderStatus,
	offchainUpdates *types.OffchainUpdates,
	err error,
) {
	order.MustBeStatefulOrder()
	// Perform stateful validation without checking existing order in state.
	err = k.PerformStatefulOrderValidation(ctx, order, blockHeight, true)
	if err != nil {
		return 0, 0, nil, err
	}

	// Place the order on the memclob and return the result. Note that we shouldn't perform
	// the add-to-orderbook collateralization check here since it was performed in a prior block.
	return memclob.PlaceOrder(
		ctx,
		*order,
		false,
	)
}

// PlaceStatefulOrdersFromLastBlock runs stateful validation for the provided placed stateful orders
// included in the last block. For each valid order, the keeper's memclob will be updated.
// Note that stateful orders could fail to be placed due to various reasons such as collateralization
// check failures, self-trade errors, etc. In these cases the `checkState` will not be written to.
// Also note this has to be done for all stateful order placements included in the last block to ensure each
// operation is inserted correctly into `operationsHashToNonce` as a pre-existing stateful order placement.
// This step is done in `PrepareCheckState`.
func (k Keeper) PlaceStatefulOrdersFromLastBlock(
	ctx sdk.Context,
	placedStatefulOrders []types.Order,
	existingOffchainUpdates *types.OffchainUpdates,
) *types.OffchainUpdates {
	lib.AssertCheckTxMode(ctx)
	// Use the height of the next block. Check if this order would be valid if it were included
	// in the next block height, not in the block that was already committed.
	currBlockHeight := lib.MustConvertIntegerToUint32(ctx.BlockHeight() + 1)
	for _, order := range placedStatefulOrders {
		// Panic if called with a non-stateful order.
		order.MustBeStatefulOrder()

		placeOrderCtx, writeCache := ctx.CacheContext()

		// Validate and place order.
		_, orderStatus, placeOrderOffchainUpdates, err := k.AddPreexistingStatefulOrder(
			placeOrderCtx,
			&order,
			currBlockHeight,
			k.MemClob,
		)

		if err != nil {
			ctx.Logger().Debug(
				fmt.Sprintf(
					"MustPlaceStatefulOrdersFromLastBlock: PlaceOrder() returned an error %+v for order %+v",
					err,
					order,
				),
			)

			// Note: Currently, the error returned from placing the order determines whether an order
			// removal message is sent to the Indexer. This may change later on to be a check on whether
			// the order has an existing nonce.
			if k.indexerEventManager.Enabled() && off_chain_updates.ShouldSendOrderRemovalOnReplay(err) {
				// If the stateful order is dropped while adding it to the book, return an off-chain order remove
				// message for the order. It's possible that this validator already knows about this order, in which
				// case an `ErrInvalidReplacement` error would be returned here.

				// It's possible that this is a new stateful order that this validator has never learned about,
				// but the validator failed to place on it on the book, even though it does exist in state.
				// In this case, Indexer could be learning of this order for the first
				// time with this removal.
				if message, success := off_chain_updates.CreateOrderRemoveMessageWithDefaultReason(
					ctx.Logger(),
					order.OrderId,
					orderStatus,
					err,
					off_chain_updates.OrderRemove_ORDER_REMOVAL_STATUS_BEST_EFFORT_CANCELED,
					off_chain_updates.OrderRemove_ORDER_REMOVAL_REASON_INTERNAL_ERROR,
				); success {
					existingOffchainUpdates.AddRemoveMessage(order.OrderId, message)
				}
			}
		} else {
			writeCache()

			if k.indexerEventManager.Enabled() {
				existingOffchainUpdates.BulkUpdate(placeOrderOffchainUpdates)
			}
		}
	}

	return existingOffchainUpdates
}

// PerformOrderCancellationStatefulValidation performs stateful validation on an order cancellation.
// The order cancellation can be either stateful or short term. This validation performs state reads.
//
// This validation ensures:
//   - Stateful Order Cancellation cancels an existing stateful order.
//   - Stateful Order Cancellation GTBT is greater than or equal to than stateful order GTBT.
//   - Stateful Order Cancellation GTBT is greater than than block time of previous block.
//   - Stateful Order Cancellation GTBT is less than or equal to `StatefulOrderTimeWindow` away from block time of
//     previous block.
//   - Short term Order Cancellation GTB must be greater than or equal to blockHeight
//   - Short term Order Cancellation GTB is less than or equal to ShortBlockWindow block hight in the future.
func (k Keeper) PerformOrderCancellationStatefulValidation(
	ctx sdk.Context,
	msgCancelOrder *types.MsgCancelOrder,
	blockHeight uint32,
) error {
	orderIdToCancel := msgCancelOrder.GetOrderId()
	if orderIdToCancel.IsStatefulOrder() {
		cancelGoodTilBlockTime := msgCancelOrder.GetGoodTilBlockTime()
		previousBlockTime := k.MustGetBlockTimeForLastCommittedBlock(ctx)

		// Return an error if `goodTilBlockTime` is less than previous block's blockTime
		if cancelGoodTilBlockTime <= lib.MustConvertIntegerToUint32(previousBlockTime.Unix()) {
			return types.ErrTimeExceedsGoodTilBlockTime
		}

		// Return an error if `goodTilBlockTime` is further into the future
		// than the previous block time plus `StatefulOrderTimeWindow`.
		endTime := previousBlockTime.Add(types.StatefulOrderTimeWindow)
		if cancelGoodTilBlockTime > lib.MustConvertIntegerToUint32(endTime.Unix()) {
			return sdkerrors.Wrapf(
				types.ErrGoodTilBlockTimeExceedsStatefulOrderTimeWindow,
				"GoodTilBlockTime %v exceeds the previous blockTime plus StatefulOrderTimeWindow %v. MsgCancelOrder: %+v",
				cancelGoodTilBlockTime,
				endTime,
				msgCancelOrder,
			)
		}

		// Fetch the highest priority order we are trying to cancel from state.
		statefulOrderPlacement, orderToCancelExists := k.GetStatefulOrderPlacement(ctx, orderIdToCancel)

		// The order we are cancelling must exist in state.
		if !orderToCancelExists {
			return sdkerrors.Wrapf(
				types.ErrStatefulOrderDoesNotExist,
				"Order Id to cancel does not exist. OrderId : %+v",
				orderIdToCancel,
			)
		}

		// Highest priority stateful matching order to cancel.
		existingStatefulOrder := statefulOrderPlacement.Order
		// Return an error if cancellation's GTBT is less than stateful order's GTBT.
		if cancelGoodTilBlockTime < existingStatefulOrder.GetGoodTilBlockTime() {
			return sdkerrors.Wrapf(
				types.ErrInvalidStatefulOrderCancellation,
				"cancellation goodTilBlockTime less than stateful order goodTilBlockTime."+
					" cancellation %+v, order %+v",
				msgCancelOrder,
				statefulOrderPlacement,
			)
		}
	} else {
		goodTilBlock := msgCancelOrder.GetGoodTilBlock()
		// Return an error if `goodTilBlock` is in the past.
		if goodTilBlock < blockHeight {
			return types.ErrHeightExceedsGoodTilBlock
		}

		// Return an error if `goodTilBlock` is further into the future than `ShortBlockWindow`.
		if goodTilBlock > types.ShortBlockWindow+blockHeight {
			return types.ErrGoodTilBlockExceedsShortBlockWindow
		}
	}
	return nil
}

// PerformStatefulOrderValidation performs stateful validation on an order. This validation performs
// state reads.
//
// This validation ensures:
//   - The `ClobPairId` on the order is for a valid CLOB.
//   - The `GoodTilBlock` of the order does not exceed the provided `blockHeight`.
//   - The `GoodTilBlock` of the order does not exceed the provided `blockHeight + ShortBlockWindow`.
//   - The `Subticks` of the order is a multiple of the ClobPair's `SubticksPerTick`.
//   - The `Quantums` of the order is greater than the ClobPair's `MinOrderBaseQuantums`.
//   - The `Quantums` of the order is a multiple of the ClobPair's `StepBaseQuantums`.
func (k Keeper) PerformStatefulOrderValidation(
	ctx sdk.Context,
	order *types.Order,
	blockHeight uint32,
	isPreexistingStatefulOrder bool,
) error {
	defer telemetry.ModuleMeasureSince(
		types.ModuleName,
		time.Now(),
		metrics.PlaceOrder,
		metrics.ValidateOrder,
		metrics.Latency,
	)
	clobPair, found := k.GetClobPair(ctx, order.GetClobPairId())
	if !found {
		return sdkerrors.Wrapf(
			types.ErrInvalidClob,
			"Clob %v is not a valid clob",
			order.GetClobPairId(),
		)
	}

	if order.Subticks%uint64(clobPair.SubticksPerTick) != 0 {
		return sdkerrors.Wrapf(
			types.ErrInvalidPlaceOrder,
			"Order subticks %v must be a multiple of the ClobPair's SubticksPerTick %v",
			order.Subticks,
			clobPair.SubticksPerTick,
		)
	}

	if order.Quantums < clobPair.MinOrderBaseQuantums {
		return sdkerrors.Wrapf(
			types.ErrInvalidPlaceOrder,
			"Order Quantums %v must be greater than the ClobPair's MinOrderBaseQuantums %v",
			order.Quantums,
			clobPair.MinOrderBaseQuantums,
		)
	}

	if order.Quantums%clobPair.StepBaseQuantums != 0 {
		return sdkerrors.Wrapf(
			types.ErrInvalidPlaceOrder,
			"Order Quantums %v must be a multiple of the ClobPair's StepBaseQuantums %v",
			order.Quantums,
			clobPair.StepBaseQuantums,
		)
	}

	if order.OrderId.IsShortTermOrder() {
		goodTilBlock := order.GetGoodTilBlock()

		// Return an error if `goodTilBlock` is in the past.
		if goodTilBlock < blockHeight {
			return sdkerrors.Wrapf(
				types.ErrHeightExceedsGoodTilBlock,
				"GoodTilBlock %v is less than the current blockHeight %v",
				goodTilBlock,
				blockHeight,
			)
		}

		// Return an error if `goodTilBlock` is further into the future than `ShortBlockWindow`.
		if goodTilBlock > types.ShortBlockWindow+blockHeight {
			return sdkerrors.Wrapf(
				types.ErrGoodTilBlockExceedsShortBlockWindow,
				"The GoodTilBlock %v exceeds the current blockHeight %v plus ShortBlockWindow %v",
				goodTilBlock,
				blockHeight,
				types.ShortBlockWindow,
			)
		}
	} else {
		goodTilBlockTimeUnix := order.GetGoodTilBlockTime()
		previousBlockTime := k.MustGetBlockTimeForLastCommittedBlock(ctx)
		previousBlockTimeUnix := lib.MustConvertIntegerToUint32(previousBlockTime.Unix())

		// Return an error if `goodTilBlockTime` is less than or equal to the
		// block time of the previous block.
		if goodTilBlockTimeUnix <= previousBlockTimeUnix {
			return sdkerrors.Wrapf(
				types.ErrTimeExceedsGoodTilBlockTime,
				"GoodTilBlockTime %v is less than the previous blockTime %v",
				goodTilBlockTimeUnix,
				previousBlockTimeUnix,
			)
		}

		// Return an error if `goodTilBlockTime` is further into the future
		// than the previous block time plus `StatefulOrderTimeWindow`.
		endTimeUnix := lib.MustConvertIntegerToUint32(
			previousBlockTime.Add(types.StatefulOrderTimeWindow).Unix(),
		)
		if goodTilBlockTimeUnix > endTimeUnix {
			return sdkerrors.Wrapf(
				types.ErrGoodTilBlockTimeExceedsStatefulOrderTimeWindow,
				"GoodTilBlockTime %v exceeds the previous blockTime plus StatefulOrderTimeWindow %v",
				goodTilBlockTimeUnix,
				endTimeUnix,
			)
		}

		// If the stateful order already exists in state, validate
		// that the new stateful order has a higher priority than the existing order.
		// TODO(DEC-1238): Support stateful order replacements by accounting for pendingStatefulOrders.
		statefulOrderPlacement, found := k.GetStatefulOrderPlacement(ctx, order.OrderId)
		if !isPreexistingStatefulOrder {
			if found {
				existingOrder := statefulOrderPlacement.GetOrder()
				if existingOrder.MustCmpReplacementOrder(order) >= 0 {
					return sdkerrors.Wrapf(
						types.ErrStatefulOrderAlreadyExists,
						"Existing order GoodTilBlockTime (%v), new order GoodTilBlockTime (%v). "+
							"Existing order: (%+v). New order: (%+v).",
						existingOrder.GetGoodTilBlockTime(),
						goodTilBlockTimeUnix,
						existingOrder,
						order,
					)
				}
			}
		}
		// TODO(CLOB-249): Add Sanity check that preexisting stateful order exists in state.
	}

	return nil
}

// MustValidateReduceOnlyOrder makes sure the given reduce-only
// order is valid with respect to the current position size.
// Specifically, this function validates:
//   - The reduce-only order is on the opposite side of the existing position.
//   - The reduce-only order does not change the subaccount's position side.
func (k Keeper) MustValidateReduceOnlyOrder(
	ctx sdk.Context,
	order types.MatchableOrder,
	matchedAmount uint64,
) error {
	if !order.IsReduceOnly() {
		panic("Order is not reduce-only.")
	}

	// Get the current position size from state.
	currentPositionSize := k.GetStatePosition(
		ctx,
		order.GetSubaccountId(),
		order.GetClobPairId(),
	)

	// Validate that the reduce-only order is on the opposite side of the existing position.
	if order.IsBuy() {
		if currentPositionSize.Sign() != -1 {
			return sdkerrors.Wrapf(
				types.ErrReduceOnlyWouldIncreasePositionSize,
				"Reduce-only order failed validation while matching. Order: (%+v), position-size: (%+v)",
				order,
				currentPositionSize,
			)
		}
	} else {
		if currentPositionSize.Sign() != 1 {
			return sdkerrors.Wrapf(
				types.ErrReduceOnlyWouldIncreasePositionSize,
				"Reduce-only order failed validation while matching. Order: (%+v), position-size: (%+v)",
				order,
				currentPositionSize,
			)
		}
	}

	// Validate that the reduce-only order does not change the subaccount's position side.
	bigMatchedAmount := new(big.Int).SetUint64(matchedAmount)
	if bigMatchedAmount.CmpAbs(currentPositionSize) == 1 {
		return sdkerrors.Wrapf(
			types.ErrReduceOnlyWouldChangePositionSide,
			"Current position size: %v, reduce-only order fill amount: %v",
			currentPositionSize,
			bigMatchedAmount,
		)
	}
	return nil
}

// AddOrderToOrderbookCollatCheck performs collateralization checks for orders to determine whether or not they may
// be added to the orderbook.
func (k Keeper) AddOrderToOrderbookCollatCheck(
	ctx sdk.Context,
	clobPairId types.ClobPairId,
	// TODO(DEC-1713): Convert this to 2 parameters: SubaccountId and a slice of PendingOpenOrders.
	subaccountOpenOrders map[satypes.SubaccountId][]types.PendingOpenOrder,
) (
	success bool,
	successPerUpdate map[satypes.SubaccountId]satypes.UpdateResult,
) {
	defer telemetry.ModuleMeasureSince(
		types.ModuleName,
		time.Now(),
		metrics.CollateralizationCheck,
		metrics.Latency,
	)

	telemetry.SetGauge(
		float32(len(subaccountOpenOrders)),
		types.ModuleName,
		metrics.CollateralizationCheckSubaccounts,
		metrics.Count,
	)

	clobPair, found := k.GetClobPair(ctx, clobPairId)
	if !found {
		panic(types.ErrInvalidClob)
	}

	makerFee := clobPair.GetFeePpm(false)

	pendingUpdates := types.NewPendingUpdates()

	// Retrieve the associated `PerpetualId` for the `ClobPair`.
	perpetualId, err := clobPair.GetPerpetualId()
	// If an error is returned, this implies stateful order validation was not performed properly, therefore panic.
	if err != nil {
		panic(sdkerrors.Wrapf(err, "clob pair ID = (%d)", clobPairId))
	}

	// Use the `PerpetualId` to retrieve the `Perpetual` and `Market` so we can determine the oracle price.
	perpetual, market, err := k.perpetualsKeeper.GetPerpetualAndMarket(ctx, perpetualId)
	// If an error is returned, this implies stateful order validation was not performed properly, therefore panic.
	if err != nil {
		panic(sdkerrors.Wrapf(err, "perpetual ID = (%d)", perpetualId))
	}

	// Get the oracle price for the market for our open orders.
	oraclePriceSubticksRat := types.PriceToSubticks(
		market,
		clobPair,
		perpetual.AtomicResolution,
		lib.QuoteCurrencyAtomicResolution,
	)
	if oraclePriceSubticksRat.Cmp(big.NewRat(0, 1)) == 0 {
		panic(
			sdkerrors.Wrapf(
				types.ErrZeroPriceForOracle,
				"clob pair ID = (%d), perpetual ID = (%d), market ID = (%d)",
				clobPairId,
				perpetualId,
				market.Id,
			),
		)
	}
	// TODO(DEC-1713): Complete as many calculations from getPessimisticCollateralCheckPrice as possible here
	// so we aren't recalculating the same thing within the loop.

	iterateOverOpenOrdersStart := time.Now()
	for subaccountId, openOrders := range subaccountOpenOrders {
		telemetry.SetGauge(
			float32(len(openOrders)),
			types.ModuleName,
			metrics.SubaccountPendingMatches,
			metrics.Count,
		)
		// For each subaccount ID, create the update from all of its existing open orders for the clob and side.
		for _, openOrder := range openOrders {
			if openOrder.ClobPairId != clobPairId {
				panic(fmt.Sprintf("Order `ClobPairId` must equal `clobPairId` for order %+v", openOrder))
			}

			// We don't want to allow users to place orders to improve their collateralization, so we choose between the
			// order price (user-input) or the oracle price (sane default) and select the price that results in the
			// most pessimistic collateral-check outcome.
			collatCheckPriceSubticks, err := getPessimisticCollateralCheckPrice(
				oraclePriceSubticksRat,
				openOrder.Subticks,
				openOrder.IsBuy,
			)
			if satypes.ErrIntegerOverflow.Is(err) {
				// TODO(DEC-1701): Determine best action to take if the oracle price overflows max uint64
				ctx.Logger().Error(
					fmt.Sprintf(
						"Integer overflow: oracle price (subticks) exceeded uint64 max. "+
							"perpetual ID = (%d), oracle price = (%+v), is buy = (%t)",
						perpetualId,
						oraclePriceSubticksRat,
						openOrder.IsBuy,
					),
				)
			} else if err != nil {
				panic(
					sdkerrors.Wrapf(
						err,
						"perpetual id = (%d), oracle price = (%+v), is buy = (%t)",
						perpetualId,
						oraclePriceSubticksRat,
						openOrder.IsBuy,
					),
				)
			}

			bigFillQuoteQuantums, err := getFillQuoteQuantums(
				clobPair,
				collatCheckPriceSubticks,
				openOrder.RemainingQuantums,
			)

			// If an error is returned, this implies stateful order validation was not performed properly, therefore panic.
			if err != nil {
				panic(err)
			}

			bigFillAmount := openOrder.RemainingQuantums.ToBigInt()
			addPerpetualFillAmountStart := time.Now()
			pendingUpdates.AddPerpetualFill(
				subaccountId,
				perpetualId,
				openOrder.IsBuy,
				makerFee,
				bigFillAmount,
				bigFillQuoteQuantums,
			)
			telemetry.ModuleMeasureSince(
				types.ModuleName,
				addPerpetualFillAmountStart,
				metrics.AddPerpetualFillAmount,
				metrics.Latency,
			)
		}
	}
	telemetry.ModuleMeasureSince(
		types.ModuleName,
		iterateOverOpenOrdersStart,
		metrics.IterateOverPendingMatches,
		metrics.Latency,
	)

	covertToUpdatesStart := time.Now()
	updates := pendingUpdates.ConvertToUpdates()
	telemetry.ModuleMeasureSince(
		types.ModuleName,
		covertToUpdatesStart,
		metrics.ConvertToUpdates,
		metrics.Latency,
	)

	success, successPerSubaccountUpdate, err := k.subaccountsKeeper.CanUpdateSubaccounts(ctx, updates)
	// TODO(DEC-191): Remove the error case from `CanUpdateSubaccounts`, which can only occur on overflow and specifying
	// duplicate accounts.
	if err != nil {
		panic(err)
	}

	result := make(map[satypes.SubaccountId]satypes.UpdateResult, len(updates))
	for i, update := range updates {
		result[update.SubaccountId] = successPerSubaccountUpdate[i]
	}

	return success, result
}

// GetStatePosition returns the current size of a subaccount's position for the specified `clobPairId`.
func (k Keeper) GetStatePosition(ctx sdk.Context, subaccountId satypes.SubaccountId, clobPairId types.ClobPairId,
) (
	positionSizeQuantums *big.Int,
) {
	// Get the CLOB pair, and panic if it does not exist.
	clobPair, found := k.GetClobPair(ctx, clobPairId)
	if !found {
		panic(fmt.Sprintf("GetStatePosition: CLOB pair %d not found", clobPairId))
	}

	// Get the perpetual ID for this CLOB pair, and panic if it is not a perpetual CLOB.
	perpetualId, err := clobPair.GetPerpetualId()
	if err != nil {
		panic(
			sdkerrors.Wrap(
				err,
				"GetStatePosition: Reduce-only orders for assets not implemented",
			),
		)
	}

	// Get the position size corresponding to `perpetualId` held by this subaccount, negative
	// if short and positive if long. If the subaccount does not have an open position
	// corresponding to `perpetualId`, a position size of zero is returned.
	subaccount := k.subaccountsKeeper.GetSubaccount(ctx, subaccountId)
	position, _ := subaccount.GetPerpetualPositionForId(perpetualId)
	return position.GetBigQuantums()
}

// InitStatefulOrdersInMemClob places all stateful orders in state on the memclob, placed in ascending
// order by time priority.
// This is called during app initialization in `app.go`, before any ABCI calls are received
// and after all MemClob orderbooks are instantiated.
func (k Keeper) InitStatefulOrdersInMemClob(
	ctx sdk.Context,
) {
	defer telemetry.ModuleMeasureSince(
		types.ModuleName,
		time.Now(),
		metrics.PlaceOrder,
		metrics.Hydrate,
		metrics.Latency,
	)

	// Get all stateful orders in state, ordered by time priority ascending order.
	// Place each order in the memclob, ignoring errors if they occur.
	statefulOrders := k.GetAllStatefulOrders(ctx)
	for _, statefulOrder := range statefulOrders {
		// First fork the multistore. If `PlaceOrder` fails, we don't want to write to state.
		placeOrderCtx, writeCache := ctx.CacheContext()

		// Place the order on the memclob and return the result.
		// Note that we skip stateful validation since these orders are already in state and don't
		// need to be statefully validated.
		orderSizeOptimisticallyFilledFromMatchingQuantums, _, offchainUpdates, err := k.MemClob.PlaceOrder(
			placeOrderCtx,
			statefulOrder,
			false,
		)

		// If the order was placed successfully, write to the underlying `checkState`.
		if err == nil {
			writeCache()
		}

		telemetry.IncrCounter(1, types.ModuleName, metrics.PlaceOrder, metrics.Hydrate, metrics.Count)
		if err != nil {
			telemetry.IncrCounter(1, types.ModuleName, metrics.PlaceOrder, metrics.Hydrate, metrics.Rejected)
		}

		// Add seen place order to in-memory map for metrics purposes. Only add if we passed all
		// validation and the memclob attempted to place it.
		if err == nil {
			k.AddSeenPlaceOrder(ctx, *types.NewMsgPlaceOrder(statefulOrder))
		} else {
			// TODO(DEC-847): Revisit this error log once `MsgRemoveOrder` is implemented,
			// since it should potentially be a panic.
			ctx.Logger().Error(
				"InitStatefulOrdersInMemClob: PlaceOrder() returned an error",
				"error",
				err,
			)
		}

		// Send off-chain updates generated from placing the order. `SendOffchainData` enqueues the
		// the messages to be sent in a channel and should be non-blocking.
		// Off-chain update messages should be only be returned if the `IndexerMessageSender`
		// is enabled (`msgSender.Enabled()` returns true).
		k.SendOffchainMessages(offchainUpdates, nil, metrics.SendPlaceOrderOffchainUpdates)

		if orderSizeOptimisticallyFilledFromMatchingQuantums > 0 {
			telemetry.IncrCounter(1, types.ModuleName, metrics.PlaceOrder, metrics.Hydrate, metrics.Matched)
		}
	}
}

// sendOffchainMessagesWithTxHash sends all the `Message` in the offchainUpdates passed in along with
// an additional header for the transaction hash passed in.
func (k Keeper) sendOffchainMessagesWithTxHash(
	offchainUpdates *types.OffchainUpdates,
	txHash []byte,
	metric string,
) {
	k.SendOffchainMessages(
		offchainUpdates,
		[]msgsender.MessageHeader{
			{
				Key:   msgsender.TransactionHashHeaderKey,
				Value: txHash,
			},
		},
		metric,
	)
}

// SendOffchainMessages sends all the `Message` in the offchainUpdates passed in along with
// any additional headers passed in. No headers will be added if a `nil` or empty list of additional
// headers is passed in.
func (k Keeper) SendOffchainMessages(
	offchainUpdates *types.OffchainUpdates,
	additionalHeaders []msgsender.MessageHeader,
	metric string,
) {
	defer telemetry.ModuleMeasureSince(
		types.ModuleName,
		time.Now(),
		metric,
		metrics.Latency,
	)
	for _, update := range offchainUpdates.GetMessages() {
		for _, header := range additionalHeaders {
			update = update.AddHeader(header)
		}
		k.GetIndexerEventManager().SendOffchainData(update)
	}
}

// getPessimisticCollateralCheckPrice returns the price in subticks we should use for collateralization checks.
// It pessimistically rounds oraclePriceSubticksRat (up for buys, down for sells) and then pessimistically
// chooses the subticks value to return: the highest for buys, the lowest for sells.
//
// Returns an error if:
// - the oraclePriceSubticksRat, after being rounded pessimistically, overflows a uint64
func getPessimisticCollateralCheckPrice(
	oraclePriceSubticksRat *big.Rat,
	makerOrderSubticks types.Subticks,
	isBuy bool,
) (price types.Subticks, err error) {
	// TODO(DEC-1713): Move this rounding to before the PendingOpenOrders loop.
	// The oracle price is guaranteed to be >= 0. Since we are using this value for a collateralization check,
	// we want to round pessimistically (up for buys, down for sells).
	oraclePriceSubticksInt := lib.BigRatRound(oraclePriceSubticksRat, isBuy)

	// TODO(DEC-1701): Determine best action to take if the oracle price overflows max uint64
	var oraclePriceSubticks types.Subticks
	if oraclePriceSubticksInt.IsUint64() {
		oraclePriceSubticks = types.Subticks(lib.Max(1, oraclePriceSubticksInt.Uint64()))
	} else {
		// Clamping the oracle price here should be fine because there are 2 outcomes:
		// 1. This is a sell order, in which case we choose the lowest value, so we choose the maker sell price which
		// would be identical to the old logic.
		// 2. This is a buy order, and we select uint64 max as the price, which will fail the collateral check in any
		// real-world example.
		oraclePriceSubticks = types.Subticks(math.MaxUint64)
		err = satypes.ErrIntegerOverflow
	}

	if isBuy {
		return lib.Max(makerOrderSubticks, oraclePriceSubticks), err
	}
	return lib.Min(makerOrderSubticks, oraclePriceSubticks), err
}

// getFillQuoteQuantums returns the total fillAmount price in quote quantums based on the maker subticks.
// This value is always positive.
//
// Returns an error if:
// - The Order is not for a `PerpetualClob`.
// - The underlying `Price` does not exist.
func getFillQuoteQuantums(
	clobPair types.ClobPair,
	makerSubticks types.Subticks,
	fillAmount satypes.BaseQuantums,
) (*big.Int, error) {
	defer telemetry.ModuleMeasureSince(
		types.ModuleName,
		time.Now(),
		metrics.GetFillQuoteQuantums,
		metrics.Latency,
	)

	if perpetualClobMetadata := clobPair.GetPerpetualClobMetadata(); perpetualClobMetadata == nil {
		return nil, types.ErrAssetOrdersNotImplemented
	}

	quantumConversionExponent := clobPair.QuantumConversionExponent

	quoteQuantums := types.FillAmountToQuoteQuantums(
		makerSubticks,
		fillAmount,
		quantumConversionExponent,
	)

	return quoteQuantums, nil
}
