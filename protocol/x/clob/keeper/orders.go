package keeper

import (
	errorsmod "cosmossdk.io/errors"
	"fmt"
	"math"
	"math/big"
	"time"

	"github.com/cometbft/cometbft/crypto/tmhash"
	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4-chain/protocol/indexer/msgsender"
	"github.com/dydxprotocol/v4-chain/protocol/indexer/off_chain_updates"
	indexershared "github.com/dydxprotocol/v4-chain/protocol/indexer/shared"
	"github.com/dydxprotocol/v4-chain/protocol/lib"
	"github.com/dydxprotocol/v4-chain/protocol/lib/metrics"
	"github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
	satypes "github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"
)

func (k Keeper) GetOperations(ctx sdk.Context) *types.MsgProposedOperations {
	operationsQueueRaw := k.MemClob.GetOperationsRaw(ctx)

	msgProposedOperations := &types.MsgProposedOperations{
		OperationsQueue: operationsQueueRaw,
	}

	if err := msgProposedOperations.ValidateBasic(); err != nil {
		operations, _ := k.MemClob.GetOperationsToReplay(ctx)
		panic(fmt.Sprintf("MsgProposedOperations failed validation: %s. Operations to replay: %+v", err, operations))
	}

	if _, err := types.ValidateAndTransformRawOperations(
		ctx,
		operationsQueueRaw,
		k.txDecoder,
		k.antehandler,
	); err != nil {
		operations, _ := k.MemClob.GetOperationsToReplay(ctx)
		panic(fmt.Sprintf("MsgProposedOperations failed stateful validation: %s. Operations to replay: %+v", err, operations))
	}

	return msgProposedOperations
}

// CancelShortTermOrder removes a Short-Term order by `OrderId` (if it exists) from all order-related data structures
// in the memclob. As well, CancelShortTermOrder adds (or updates) a cancel to the desired `goodTilBlock` in the
// memclob.
// If a cancel already exists for this order with a lower `goodTilBlock`, the cancel is updated to the
// new `goodTilBlock`. This method is meant to be used in the CheckTx flow. It uses the next block height.
//
// An error will be returned if any of the following conditions are true:
// - The cancel's `GoodTilblock` is less than or equal to the next block height.
// - The cancel's `GoodTilblock` is greater than the sum of the next block height and `ShortBlockWindow`.
// - The memclob itself returns an error.
//
// This method assumes the provided MsgCancelOrder has already passed ValidateBasic in CheckTx.
func (k Keeper) CancelShortTermOrder(
	ctx sdk.Context,
	msgCancelOrder *types.MsgCancelOrder,
) error {
	lib.AssertCheckTxMode(ctx)
	// Note that we add `+1` here to account for the fact that `ctx.BlockHeight()` is technically the
	// previously mined block, not the next block that will be proposed. This is due to the fact that
	// this function is only ever called during `CheckTx`.
	nextBlockHeight := lib.MustConvertIntegerToUint32(ctx.BlockHeight() + 1)

	defer telemetry.ModuleMeasureSince(types.ModuleName, time.Now(), metrics.CancelShortTermOrder, metrics.Latency)
	telemetry.IncrCounter(1, types.ModuleName, metrics.CancelShortTermOrder, metrics.Count)

	// Perform all stateful validation on the order. Order may be stateful or short term.
	if err := k.PerformOrderCancellationStatefulValidation(ctx, msgCancelOrder, nextBlockHeight); err != nil {
		return err
	}

	// Update in-memory orderbook to remove order.
	offchainUpdates, err := k.MemClob.CancelOrder(ctx, msgCancelOrder)

	if err == nil {
		k.sendOffchainMessagesWithTxHash(
			offchainUpdates,
			tmhash.Sum(ctx.TxBytes()),
			metrics.SendCancelOrderOffchainUpdates,
		)
	}
	return err
}

// PlaceShortTermOrder places an order on the corresponding orderbook, and performs matching if placing the
// order causes an overlap. This function will return the result of calling `PlaceOrder` on the
// keeper's memclob. This method is meant to be used in the CheckTx flow. It uses the next block height.
//
// An error will be returned if any of the following conditions are true:
//   - Standard stateful validation fails.
//   - The memclob itself returns an error.
//
// This method will panic if the provided order is not a Short-Term order.
func (k Keeper) PlaceShortTermOrder(
	ctx sdk.Context,
	msg *types.MsgPlaceOrder,
) (
	orderSizeOptimisticallyFilledFromMatchingQuantums satypes.BaseQuantums,
	orderStatus types.OrderStatus,
	err error,
) {
	msg.Order.OrderId.MustBeShortTermOrder()

	lib.AssertCheckTxMode(ctx)
	nextBlockHeight := lib.MustConvertIntegerToUint32(ctx.BlockHeight() + 1)

	return k.placeOrder(ctx, msg, nextBlockHeight, k.MemClob)
}

// CancelStatefulOrder performs stateful order cancellation validation and removes the stateful order
// from state and the memstore.
//
// The following conditions must be true otherwise an error will be returned:
//   - Stateful Order Cancellation cancels an existing stateful order.
//   - Stateful Order Cancellation GTBT is greater than or equal to than stateful order GTBT.
//   - Stateful Order Cancellation GTBT is greater than the block time of previous block.
//   - Stateful Order Cancellation GTBT is less than or equal to `StatefulOrderTimeWindow` away from block time of
//     previous block.
//
// Note that this method conditionally updates state depending on the context. This is needed
// to separate updating committed state during DeliverTx (the stateful order and the ToBeCommitted stateful order
// count) from uncommitted state that is modified during CheckTx.
func (k Keeper) CancelStatefulOrder(
	ctx sdk.Context,
	msg *types.MsgCancelOrder,
) error {
	// 1. If this is a Short-Term order, panic.
	msg.OrderId.MustBeStatefulOrder()

	// 2. Perform stateful validation on the order cancellation.
	err := k.PerformOrderCancellationStatefulValidation(
		ctx,
		msg,
		// Note that the blockHeight is not used during stateful order cancellation validation.
		0,
	)
	if err != nil {
		return err
	}

	// 3. Update uncommitted or committed state depending on whether we are in `checkTx` or `deliverTx`.
	if lib.IsDeliverTxMode(ctx) {
		// Remove the stateful order from state. Note that if the stateful order did not
		// exist in state, then it would have failed validation in the previous step.
		k.MustRemoveStatefulOrder(ctx, msg.OrderId)

		// Decrement the `to be committed` stateful order count.
		k.SetToBeCommittedStatefulOrderCount(
			ctx,
			msg.OrderId,
			k.GetToBeCommittedStatefulOrderCount(ctx, msg.OrderId)-1,
		)
	} else {
		// Write the stateful order cancellation to uncommitted state. PerformOrderCancellationStatefulValidation will
		// return an error if the order cancellation already exists which will prevent
		// MustAddUncommittedStatefulOrderCancellation from panicking.
		k.MustAddUncommittedStatefulOrderCancellation(ctx, msg)
		// TODO(DEC-1238): Support stateful order replacements by removing the uncommitted order placement.
		// This should allow a cycle of place + cancel + place + cancel + ... which we currently disallow during
		// `DeliverTx`.
	}

	return nil
}

// PlaceStatefulOrder performs order validation, equity tier limit check, a collateralization check and writes the
// order to state and the memstore. The order will not be placed on the orderbook.
//
// An error will be returned if any of the following conditions are true:
//   - Standard stateful validation fails.
//   - Equity tier limit exceeded.
//   - Collateralization check fails.
//
// Note that this method conditionally updates state depending on the context. This is needed
// to separate updating committed state during DeliverTx from uncommitted state that is modified during
// CheckTx.
//
// This method will panic if the provided order is not a Stateful order.
func (k Keeper) PlaceStatefulOrder(
	ctx sdk.Context,
	msg *types.MsgPlaceOrder,
) error {
	// 1. Ensure the order is not a Short-Term order.
	order := msg.Order
	order.OrderId.MustBeStatefulOrder()

	// 2. Perform stateful validation on the order.
	if err := k.PerformStatefulOrderValidation(
		ctx,
		&order,
		// Note that the blockHeight is not used during stateful order validation.
		0,
		false,
	); err != nil {
		return err
	}

	// 3. Check that adding the order would not exceed the equity tier for the account.
	if err := k.ValidateSubaccountEquityTierLimitForNewOrder(ctx, order); err != nil {
		return err
	}

	// 4. Perform a collateralization check for the full size of the order to mitigate spam.
	// TODO(CLOB-725): Consider using a pessimistic collateralization check.
	_, successPerSubaccountUpdate := k.AddOrderToOrderbookCollatCheck(
		ctx,
		order.GetClobPairId(),
		map[satypes.SubaccountId][]types.PendingOpenOrder{
			order.OrderId.SubaccountId: {
				{
					RemainingQuantums: order.GetBaseQuantums(),
					IsBuy:             order.IsBuy(),
					Subticks:          order.GetOrderSubticks(),
					ClobPairId:        order.GetClobPairId(),
				},
			},
		},
	)

	if !successPerSubaccountUpdate[order.OrderId.SubaccountId].IsSuccess() {
		return errorsmod.Wrapf(
			types.ErrStatefulOrderCollateralizationCheckFailed,
			"PlaceStatefulOrder: order (%+v), result (%s)",
			order,
			successPerSubaccountUpdate[order.OrderId.SubaccountId].String(),
		)
	}

	// 5. If we are in `deliverTx` then we write the order to committed state otherwise add the order to uncommitted
	// state.
	if lib.IsDeliverTxMode(ctx) {
		// Write the stateful order to state and the memstore.
		k.SetLongTermOrderPlacement(ctx, order, lib.MustConvertIntegerToUint32(ctx.BlockHeight()))
		k.MustAddOrderToStatefulOrdersTimeSlice(
			ctx,
			order.MustGetUnixGoodTilBlockTime(),
			order.GetOrderId(),
		)
	} else {
		// Write the stateful order to a transient store. PerformStatefulOrderValidation will ensure that the order does
		// not exist which will prevent MustAddUncommittedStatefulOrderPlacement from panicking.
		k.MustAddUncommittedStatefulOrderPlacement(ctx, msg)
		// TODO(DEC-1238): Support stateful order replacements by removing the uncommitted order cancellation.
		// This should allow a cycle of place + cancel + place + cancel + ... which we currently disallow during
		// `DeliverTx`.
	}

	return nil
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
	err = k.PerformStatefulOrderValidation(ctx, &order, nextBlockHeight, true)
	if err != nil {
		return 0, 0, nil, err
	}

	// Place the order on the memclob and return the result.
	orderSizeOptimisticallyFilledFromMatchingQuantums, orderStatus, offchainUpdates, err = k.MemClob.PlaceOrder(
		ctx,
		msg.Order,
	)

	return orderSizeOptimisticallyFilledFromMatchingQuantums, orderStatus, offchainUpdates, err
}

// placeOrder contains shared logic for order placement. It performs
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
	}()

	// Perform stateful validation.
	err = k.PerformStatefulOrderValidation(ctx, &order, blockHeight, true)
	if err != nil {
		return 0, 0, err
	}

	// Place the order on the memclob and return the result.
	orderSizeOptimisticallyFilledFromMatchingQuantums, orderStatus, offchainUpdates, err := memclob.PlaceOrder(
		ctx,
		msg.Order,
	)

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
	// the add-to-orderbook collateralization check here for long-term orders since it was performed in a prior block,
	// but for triggered conditional orders we have not yet performed the collaterlization check.
	return memclob.PlaceOrder(
		ctx,
		*order,
	)
}

// PlaceStatefulOrdersFromLastBlock validates and places stateful orders from the last block onto the memclob.
// Note that stateful orders could fail to be placed due to various reasons such as collateralization
// check failures, self-trade errors, etc. In these cases the `checkState` will not be written to.
// This function is used in:
// 1. `PrepareCheckState` to place newly placed long term orders from the last
// block from ProcessProposerMatchesEvents.PlacedStatefulOrderIds. This is step 3 in PrepareCheckState.
// 2. `PlaceConditionalOrdersTriggeredInLastBlock` to place conditional orders triggered in the last block
// from ProcessProposerMatchesEvents.ConditionalOrderIdsTriggeredInLastBlock. This is step 4 in PrepareCheckState.
func (k Keeper) PlaceStatefulOrdersFromLastBlock(
	ctx sdk.Context,
	placedStatefulOrderIds []types.OrderId,
	existingOffchainUpdates *types.OffchainUpdates,
) (
	offchainUpdates *types.OffchainUpdates,
) {
	lib.AssertCheckTxMode(ctx)

	for _, orderId := range placedStatefulOrderIds {
		orderId.MustBeStatefulOrder()

		orderPlacement, exists := k.GetLongTermOrderPlacement(ctx, orderId)
		if !exists {
			// Order does not exist in state and therefore should not be placed. This likely
			// indicates that the order was cancelled.
			continue
		}

		placeOrderCtx, writeCache := ctx.CacheContext()

		order := orderPlacement.GetOrder()
		// Validate and place order.
		_, orderStatus, placeOrderOffchainUpdates, err := k.AddPreexistingStatefulOrder(
			placeOrderCtx,
			&order,
			0,
			k.MemClob,
		)

		if err != nil {
			k.Logger(ctx).Debug(
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
					k.Logger(ctx),
					order.OrderId,
					orderStatus,
					err,
					off_chain_updates.OrderRemoveV1_ORDER_REMOVAL_STATUS_BEST_EFFORT_CANCELED,
					indexershared.OrderRemovalReason_ORDER_REMOVAL_REASON_INTERNAL_ERROR,
				); success {
					existingOffchainUpdates.AddRemoveMessage(order.OrderId, message)
				}
			}
		} else {
			writeCache()

			if k.indexerEventManager.Enabled() {
				existingOffchainUpdates.Append(placeOrderOffchainUpdates)
			}
		}
	}

	// Clear place messages as BEST_EFFORT_OPEN messages should not be
	// sent for stateful order placements.
	existingOffchainUpdates.CondenseMessagesForReplay()

	return existingOffchainUpdates
}

// PlaceConditionalOrdersTriggeredInLastBlock takes in a list of conditional order ids that were triggered
// in the last block, verifies they are conditional orders, verifies they are in triggered state, and places
// the orders on the memclob.
func (k Keeper) PlaceConditionalOrdersTriggeredInLastBlock(
	ctx sdk.Context,
	conditionalOrderIdsTriggeredInLastBlock []types.OrderId,
	existingOffchainUpdates *types.OffchainUpdates,
) (
	offchainUpdates *types.OffchainUpdates,
) {
	for _, orderId := range conditionalOrderIdsTriggeredInLastBlock {
		// Panic if the order is not in triggered state.
		if !k.IsConditionalOrderTriggered(ctx, orderId) {
			panic(
				fmt.Sprintf(
					"PlaceConditionalOrdersTriggeredInLastBlock: Order with OrderId %+v is not in triggered state",
					orderId,
				),
			)
		}
	}

	return k.PlaceStatefulOrdersFromLastBlock(ctx, conditionalOrderIdsTriggeredInLastBlock, existingOffchainUpdates)
}

// PerformOrderCancellationStatefulValidation performs stateful validation on an order cancellation.
// The order cancellation can be either stateful or short term. This validation performs state reads.
//
// This validation ensures:
//   - Stateful Order Cancellation for the order does not already exist in uncommitted state.
//   - Stateful Order Cancellation cancels an uncommitted or existing stateful order.
//   - Stateful Order Cancellation GTBT is greater than or equal to than stateful order GTBT.
//   - Stateful Order Cancellation GTBT is greater than the block time of previous block.
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
		previousBlockTime := k.blockTimeKeeper.GetPreviousBlockInfo(ctx).Timestamp

		// Return an error if `goodTilBlockTime` is less than previous block's blockTime
		if cancelGoodTilBlockTime <= lib.MustConvertIntegerToUint32(previousBlockTime.Unix()) {
			return types.ErrTimeExceedsGoodTilBlockTime
		}

		// Return an error if `goodTilBlockTime` is further into the future
		// than the previous block time plus `StatefulOrderTimeWindow`.
		endTime := previousBlockTime.Add(types.StatefulOrderTimeWindow)
		if cancelGoodTilBlockTime > lib.MustConvertIntegerToUint32(endTime.Unix()) {
			return errorsmod.Wrapf(
				types.ErrGoodTilBlockTimeExceedsStatefulOrderTimeWindow,
				"GoodTilBlockTime %v exceeds the previous blockTime plus StatefulOrderTimeWindow %v. MsgCancelOrder: %+v",
				cancelGoodTilBlockTime,
				endTime,
				msgCancelOrder,
			)
		}

		// Return an error if we are attempting to submit another cancellation when the mempool already has an
		// existing uncommitted cancellation for this order ID.
		existingCancellation, uncommittedCancelExists := k.GetUncommittedStatefulOrderCancellation(ctx, orderIdToCancel)
		if uncommittedCancelExists {
			return errorsmod.Wrapf(
				types.ErrStatefulOrderCancellationAlreadyExists,
				"An uncommitted stateful order cancellation with this OrderId already exists and stateful "+
					"order cancellation replacement is not supported. Existing order cancellation GoodTilBlockTime "+
					"(%v), New order cancellation GoodTilBlockTime (%v). Existing order cancellation: (%+v). New "+
					"order cancellation: (%+v).",
				existingCancellation.GetGoodTilBlockTime(),
				cancelGoodTilBlockTime,
				existingCancellation,
				msgCancelOrder,
			)
		}

		// Fetch the highest priority order we are trying to cancel from state.
		statefulOrderPlacement, orderToCancelExists := k.GetLongTermOrderPlacement(ctx, orderIdToCancel)

		// The order we are cancelling must exist in uncommitted or committed state.
		if !orderToCancelExists {
			statefulOrderPlacement, orderToCancelExists = k.GetUncommittedStatefulOrderPlacement(ctx, orderIdToCancel)

			if !orderToCancelExists {
				return errorsmod.Wrapf(
					types.ErrStatefulOrderDoesNotExist,
					"Order Id to cancel does not exist. OrderId : %+v",
					orderIdToCancel,
				)
			}
		}

		// Highest priority stateful matching order to cancel.
		existingStatefulOrder := statefulOrderPlacement.Order
		// Return an error if cancellation's GTBT is less than stateful order's GTBT.
		if cancelGoodTilBlockTime < existingStatefulOrder.GetGoodTilBlockTime() {
			return errorsmod.Wrapf(
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
//   - The `Subticks` of the order is a multiple of the ClobPair's `SubticksPerTick`.
//   - The `Quantums` of the order is a multiple of the ClobPair's `StepBaseQuantums`.
//
// This validation also ensures that the order is valid for the ClobPair's status.
//
// For short term orders it also ensures:
//   - The `GoodTilBlock` of the order is greater than the provided `blockHeight`.
//   - The `GoodTilBlock` of the order does not exceed the provided `blockHeight + ShortBlockWindow`.
//
// For stateful orders it also ensures:
//   - GTBT is greater than the block time of previous block.
//   - GTBT is less than or equal to `StatefulOrderTimeWindow` away from block time of
//     previous block.
//   - That there isn't an order cancellation in uncommitted state.
//   - That the order does not already exist in uncommitted state unless `isPreexistingStatefulOrder` is true.
//   - That the order does not already exist in committed state unless `isPreexistingStatefulOrder` is true.
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
		return errorsmod.Wrapf(
			types.ErrInvalidClob,
			"Clob %v is not a valid clob",
			order.GetClobPairId(),
		)
	}

	if order.Subticks%uint64(clobPair.SubticksPerTick) != 0 {
		return errorsmod.Wrapf(
			types.ErrInvalidPlaceOrder,
			"Order subticks %v must be a multiple of the ClobPair's SubticksPerTick %v",
			order.Subticks,
			clobPair.SubticksPerTick,
		)
	}

	if order.Quantums%clobPair.StepBaseQuantums != 0 {
		return errorsmod.Wrapf(
			types.ErrInvalidPlaceOrder,
			"Order Quantums %v must be a multiple of the ClobPair's StepBaseQuantums %v",
			order.Quantums,
			clobPair.StepBaseQuantums,
		)
	}

	// Validates the order against the ClobPair's status.
	if err := k.validateOrderAgainstClobPairStatus(ctx, order.MustGetOrder(), clobPair); err != nil {
		telemetry.IncrCounterWithLabels(
			[]string{types.ModuleName, metrics.ValidateOrder, metrics.OrderConflictsWithClobPairStatus, metrics.Count},
			1,
			append(
				order.GetOrderLabels(),
				metrics.GetLabelForBoolValue(metrics.CheckTx, ctx.IsCheckTx()),
				metrics.GetLabelForBoolValue(metrics.DeliverTx, lib.IsDeliverTxMode(ctx)),
			),
		)
		return err
	}

	if order.OrderId.IsShortTermOrder() {
		goodTilBlock := order.GetGoodTilBlock()

		// Return an error if `goodTilBlock` is in the past.
		if goodTilBlock < blockHeight {
			return errorsmod.Wrapf(
				types.ErrHeightExceedsGoodTilBlock,
				"GoodTilBlock %v is less than the current blockHeight %v",
				goodTilBlock,
				blockHeight,
			)
		}

		// Return an error if `goodTilBlock` is further into the future than `ShortBlockWindow`.
		if goodTilBlock > types.ShortBlockWindow+blockHeight {
			return errorsmod.Wrapf(
				types.ErrGoodTilBlockExceedsShortBlockWindow,
				"The GoodTilBlock %v exceeds the current blockHeight %v plus ShortBlockWindow %v",
				goodTilBlock,
				blockHeight,
				types.ShortBlockWindow,
			)
		}
	} else {
		goodTilBlockTimeUnix := order.GetGoodTilBlockTime()
		previousBlockTime := k.blockTimeKeeper.GetPreviousBlockInfo(ctx).Timestamp
		previousBlockTimeUnix := lib.MustConvertIntegerToUint32(previousBlockTime.Unix())

		// Return an error if `goodTilBlockTime` is less than or equal to the
		// block time of the previous block.
		if goodTilBlockTimeUnix <= previousBlockTimeUnix {
			return errorsmod.Wrapf(
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
			return errorsmod.Wrapf(
				types.ErrGoodTilBlockTimeExceedsStatefulOrderTimeWindow,
				"GoodTilBlockTime %v exceeds the previous blockTime plus StatefulOrderTimeWindow %v",
				goodTilBlockTimeUnix,
				endTimeUnix,
			)
		}

		// Check to see if we are aware of a cancellation that is part of the mempool and has yet to be included
		// in a block for the order in state.
		// TODO(DEC-1238): Support stateful order replacements.
		if uncommittedCancel, uncommittedCancelExists := k.GetUncommittedStatefulOrderCancellation(
			ctx, order.OrderId); uncommittedCancelExists {
			return errorsmod.Wrapf(
				types.ErrStatefulOrderPreviouslyCancelled,
				"An uncommitted stateful order cancellation with this OrderId already exists. "+
					"Existing order cancellation: (%+v). New order: (%+v).",
				uncommittedCancel,
				order,
			)
		}

		// If this is not pre-existing stateful order, then we expect it does not exist in uncommitted state.
		// TODO(DEC-1238): Support stateful order replacements.
		statefulOrderPlacement, found := k.GetUncommittedStatefulOrderPlacement(ctx, order.OrderId)
		if !isPreexistingStatefulOrder && found {
			existingOrder := statefulOrderPlacement.GetOrder()
			return errorsmod.Wrapf(
				types.ErrStatefulOrderAlreadyExists,
				"An uncommitted stateful order with this OrderId already exists and stateful order replacement is not supported. "+
					"Existing order GoodTilBlockTime (%v), New order GoodTilBlockTime (%v). "+
					"Existing order: (%+v). New order: (%+v).",
				existingOrder.GetGoodTilBlockTime(),
				goodTilBlockTimeUnix,
				existingOrder,
				order,
			)
		}

		// If the stateful order already exists in state, validate
		// that the new stateful order has a higher priority than the existing order.
		statefulOrderPlacement, found = k.GetLongTermOrderPlacement(ctx, order.OrderId)

		// If this is a pre-existing stateful order, then we expect it to exist in state.
		// Panic if the order is not in state, as this indicates an application error.
		if isPreexistingStatefulOrder && !found {
			panic(
				fmt.Sprintf(
					"PerformStatefulOrderValidation: Expected pre-existing stateful order to exist in state "+
						"order: (%+v).",
					order,
				),
			)
		}

		// If this is not pre-existing stateful order, then we expect it does not exist in state.
		// TODO(DEC-1238): Support stateful order replacements.
		if !isPreexistingStatefulOrder && found {
			existingOrder := statefulOrderPlacement.GetOrder()
			return errorsmod.Wrapf(
				types.ErrStatefulOrderAlreadyExists,
				"A stateful order with this OrderId already exists and stateful order replacement is not supported. "+
					"Existing order GoodTilBlockTime (%v), New order GoodTilBlockTime (%v). "+
					"Existing order: (%+v). New order: (%+v).",
				existingOrder.GetGoodTilBlockTime(),
				goodTilBlockTimeUnix,
				existingOrder,
				order,
			)
		}

		if order.IsConditionalOrder() {
			if order.ConditionalOrderTriggerSubticks%uint64(clobPair.SubticksPerTick) != 0 {
				return errorsmod.Wrapf(
					types.ErrInvalidPlaceOrder,
					"Conditional order trigger subticks %v must be a multiple of the ClobPair's SubticksPerTick %v",
					order.ConditionalOrderTriggerSubticks,
					clobPair.StepBaseQuantums,
				)
			}
		}
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
			return errorsmod.Wrapf(
				types.ErrReduceOnlyWouldIncreasePositionSize,
				"Reduce-only order failed validation while matching. Order: (%+v), position-size: (%+v)",
				order,
				currentPositionSize,
			)
		}
	} else {
		if currentPositionSize.Sign() != 1 {
			return errorsmod.Wrapf(
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
		return errorsmod.Wrapf(
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

	pendingUpdates := types.NewPendingUpdates()

	// Retrieve the associated `PerpetualId` for the `ClobPair`.
	oraclePriceSubticksRat := k.GetOraclePriceSubticksRat(ctx, clobPair)
	perpetualId := clobPair.MustGetPerpetualId()

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

		makerFeePpm := k.feeTiersKeeper.GetPerpetualFeePpm(ctx, subaccountId.Owner, false)
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
				k.Logger(ctx).Error(
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
					errorsmod.Wrapf(
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
				makerFeePpm,
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

// GetOraclePriceSubticksRat returns the oracle price in subticks for the given `ClobPair`.
func (k Keeper) GetOraclePriceSubticksRat(ctx sdk.Context, clobPair types.ClobPair) *big.Rat {
	// Retrieve the associated `PerpetualId` for the `ClobPair`.
	perpetualId := clobPair.MustGetPerpetualId()

	// Use the `PerpetualId` to retrieve the `Perpetual` and `Market` so we can determine the oracle price.
	perpetual, marketPrice, err := k.perpetualsKeeper.GetPerpetualAndMarketPrice(ctx, perpetualId)
	// If an error is returned, this implies stateful order validation was not performed properly, therefore panic.
	if err != nil {
		panic(errorsmod.Wrapf(err, "perpetual ID = (%d)", perpetualId))
	}

	// Get the oracle price for the market.
	oraclePriceSubticksRat := types.PriceToSubticks(
		marketPrice,
		clobPair,
		perpetual.Params.AtomicResolution,
		lib.QuoteCurrencyAtomicResolution,
	)
	if oraclePriceSubticksRat.Cmp(big.NewRat(0, 1)) == 0 {
		panic(
			errorsmod.Wrapf(
				types.ErrZeroPriceForOracle,
				"clob pair ID = (%d), perpetual ID = (%d), market ID = (%d)",
				clobPair.Id,
				perpetualId,
				marketPrice.Id,
			),
		)
	}
	return oraclePriceSubticksRat
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
			errorsmod.Wrap(
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

	// Get all placed stateful orders in state, ordered by time priority ascending order.
	// Place each order in the memclob, ignoring errors if they occur.
	statefulOrders := k.GetAllPlacedStatefulOrders(ctx)
	for _, statefulOrder := range statefulOrders {
		// First fork the multistore. If `PlaceOrder` fails, we don't want to write to state.
		placeOrderCtx, writeCache := ctx.CacheContext()

		// Place the order on the memclob and return the result.
		// Note that we skip stateful validation since these orders are already in state and don't
		// need to be statefully validated.
		orderSizeOptimisticallyFilledFromMatchingQuantums, _, offchainUpdates, err := k.MemClob.PlaceOrder(
			placeOrderCtx,
			statefulOrder,
		)

		// If the order was placed successfully, write to the underlying `checkState`.
		if err == nil {
			writeCache()
		}

		telemetry.IncrCounter(1, types.ModuleName, metrics.PlaceOrder, metrics.Hydrate, metrics.Count)
		if err != nil {
			telemetry.IncrCounter(1, types.ModuleName, metrics.PlaceOrder, metrics.Hydrate, metrics.Rejected)
		}

		if err != nil {
			// TODO(DEC-847): Revisit this error log once `MsgRemoveOrder` is implemented,
			// since it should potentially be a panic.
			k.Logger(ctx).Error(
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

// HydrateUntriggeredConditionalOrders inserts all untriggered conditional orders in state into the
// `UntriggeredConditionalOrders` data structure. Note that all untriggered conditional orders will
// be ordered by time priority. This function should only be called on application startup.
func (k Keeper) HydrateUntriggeredConditionalOrders(
	ctx sdk.Context,
) {
	defer telemetry.ModuleMeasureSince(
		types.ModuleName,
		time.Now(),
		metrics.ConditionalOrderUntriggered,
		metrics.Hydrate,
		metrics.Latency,
	)

	// Get all untriggered conditional orders in state, ordered by time priority ascending order,
	// and add them to the `UntriggeredConditionalOrders` data structure.
	// Place each order in the memclob, ignoring errors if they occur.
	untriggeredConditionalOrders := k.GetAllUntriggeredConditionalOrders(ctx)
	k.AddUntriggeredConditionalOrders(
		ctx,
		lib.MapSlice(
			untriggeredConditionalOrders,
			func(o types.Order) types.OrderId {
				return o.OrderId
			},
		),
		// Note both of these arguments are empty slices since the untriggered conditional orders
		// shouldn't be expired or canceled.
		map[types.OrderId]struct{}{},
		map[types.OrderId]struct{}{},
	)
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
