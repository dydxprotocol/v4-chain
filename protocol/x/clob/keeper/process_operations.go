package keeper

import (
	"fmt"
	"math/big"
	"time"

	errorsmod "cosmossdk.io/errors"

	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"
	indexerevents "github.com/dydxprotocol/v4-chain/protocol/indexer/events"
	"github.com/dydxprotocol/v4-chain/protocol/indexer/indexer_manager"
	indexershared "github.com/dydxprotocol/v4-chain/protocol/indexer/shared"
	"github.com/dydxprotocol/v4-chain/protocol/lib"
	"github.com/dydxprotocol/v4-chain/protocol/lib/log"
	"github.com/dydxprotocol/v4-chain/protocol/lib/metrics"
	affiliatetypes "github.com/dydxprotocol/v4-chain/protocol/x/affiliates/types"
	"github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
	satypes "github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"
)

// fetchOrdersInvolvedInOpQueue fetches all OrderIds involved in an operations
// queue's matches + short term order placements and returns them as a set.
func fetchOrdersInvolvedInOpQueue(
	operations []types.InternalOperation,
) (orderIdSet map[types.OrderId]struct{}) {
	orderIdSet = make(map[types.OrderId]struct{})
	for _, operation := range operations {
		if shortTermOrderPlacement := operation.GetShortTermOrderPlacement(); shortTermOrderPlacement != nil {
			orderId := shortTermOrderPlacement.GetOrder().OrderId
			orderIdSet[orderId] = struct{}{}
		}
		if clobMatch := operation.GetMatch(); clobMatch != nil {
			orderIdSetForClobMatch := clobMatch.GetAllOrderIds()
			orderIdSet = lib.MergeMaps(orderIdSet, orderIdSetForClobMatch)
		}
	}
	return orderIdSet
}

// ProcessProposerOperations updates on-chain state given an []OperationRaw operations queue
// representing matches that occurred in the previous block. It performs validation on an operations
// queue. If all validation passes, the operations queue is written to state.
// The following operations are written to state:
// - Order Matches, Liquidation Matches, Deleveraging Matches
func (k Keeper) ProcessProposerOperations(
	ctx sdk.Context,
	rawOperations []types.OperationRaw,
) error {
	// This function should be only run in DeliverTx mode.
	lib.AssertDeliverTxMode(ctx)
	defer telemetry.ModuleMeasureSince(types.ModuleName, time.Now(), metrics.ProcessOperations)

	// Stateless validation of RawOperations and transforms them into InternalOperations to be used internally by memclob.
	operations, err := types.ValidateAndTransformRawOperations(ctx, rawOperations, k.txDecoder, k.antehandler)
	if err != nil {
		return errorsmod.Wrapf(types.ErrInvalidMsgProposedOperations, "Error: %+v", err)
	}

	log.DebugLog(ctx, "Processing operations queue",
		log.OperationsQueue, types.GetInternalOperationsQueueTextString(operations))

	// Write results of the operations queue to state. Performs stateful validation as well.
	if err := k.ProcessInternalOperations(ctx, operations); err != nil {
		return err
	}

	// Collect the list of order ids filled and set the field in the `ProcessProposerMatchesEvents` object.
	processProposerMatchesEvents := k.GenerateProcessProposerMatchesEvents(ctx, operations)

	// Remove fully filled orders from state.
	for _, orderId := range processProposerMatchesEvents.OrderIdsFilledInLastBlock {
		if orderId.IsShortTermOrder() {
			continue
		}

		orderPlacement, placementExists := k.GetLongTermOrderPlacement(ctx, orderId)
		if placementExists {
			fillAmountExists, orderStateFillAmount, _ := k.GetOrderFillAmount(ctx, orderId)
			if !fillAmountExists {
				panic("ProcessProposerOperations: Order fill amount does not exist in state")
			}
			if orderStateFillAmount > orderPlacement.Order.GetBaseQuantums() {
				panic("ProcessProposerOperations: Order fill amount exceeds order amount")
			}

			// If the order is fully filled, remove it from state.
			if orderStateFillAmount == orderPlacement.Order.GetBaseQuantums() {
				k.MustRemoveStatefulOrder(ctx, orderId)
				telemetry.IncrCounterWithLabels(
					[]string{types.ModuleName, metrics.ProcessOperations, metrics.StatefulOrderRemoved, metrics.Count},
					1,
					append(
						orderPlacement.Order.GetOrderLabels(),
						metrics.GetLabelForStringValue(metrics.RemovalReason, types.OrderRemoval_REMOVAL_REASON_FULLY_FILLED.String()),
					),
				)

				processProposerMatchesEvents.RemovedStatefulOrderIds = append(
					processProposerMatchesEvents.RemovedStatefulOrderIds,
					orderId,
				)
			}
		}
	}

	// Update the memstore with list of orderIds filled during this block.
	// During commit, all orders that have been fully filled during this block will be removed from the memclob.
	k.MustSetProcessProposerMatchesEvents(
		ctx,
		processProposerMatchesEvents,
	)

	// Emit stats about the proposed operations.
	operationsStats := types.StatMsgProposedOperations(rawOperations)
	operationsStats.EmitStats(metrics.DeliverTx)

	return nil
}

// ProcessInternalOperations takes in an InternalOperations slice and writes all relevant
// operations to state. This function assumes that the operations have passed all stateless validation.
// This function will perform stateful validation as it processes operations.
// The following operations modify state:
// - Order Matches, Liquidation Matches, Deleveraging Matches
// - Order Removals

// Function will panic if:
// - any orderId referenced in clobMatch cannot be found.
// - any orderId referenced in order removal operations cannot be found.
func (k Keeper) ProcessInternalOperations(
	ctx sdk.Context,
	operations []types.InternalOperation,
) error {
	// Collect all the short-term orders placed for subsequent lookups.
	// All short term orders in this map have passed validation.
	placedShortTermOrders := make(map[types.OrderId]types.Order, 0)

	var affiliateOverrides map[string]bool = nil
	var affiliateParameters affiliatetypes.AffiliateParameters
	// Write the matches to state if all stateful validation passes.
	for _, operation := range operations {
		if err := k.validateInternalOperationAgainstClobPairStatus(ctx, operation); err != nil {
			return err
		}

		switch castedOperation := operation.Operation.(type) {
		case *types.InternalOperation_Match:
			// check if affiliate whitelist map is nil and initialize it if it is.
			// This is done to avoid getting whitelist map on list of operations
			// where there are no matches.
			if affiliateOverrides == nil {
				var err error
				affiliateOverrides, err = k.affiliatesKeeper.GetAffiliateOverridesMap(ctx)
				if err != nil {
					return errorsmod.Wrapf(
						err,
						"ProcessInternalOperations: Failed to get affiliates whitelist map",
					)
				}
			}
			var err error
			affiliateParameters, err = k.affiliatesKeeper.GetAffiliateParameters(ctx)
			if err != nil {
				return errorsmod.Wrapf(
					err,
					"ProcessInternalOperations: Failed to get affiliates parameters",
				)
			}
			clobMatch := castedOperation.Match
			if err := k.PersistMatchToState(ctx, clobMatch, placedShortTermOrders,
				affiliateOverrides, affiliateParameters); err != nil {
				return errorsmod.Wrapf(
					err,
					"ProcessInternalOperations: Failed to process clobMatch: %+v",
					clobMatch,
				)
			}
		case *types.InternalOperation_ShortTermOrderPlacement:
			order := castedOperation.ShortTermOrderPlacement.GetOrder()
			if err := k.PerformStatefulOrderValidation(
				ctx,
				&order,
				lib.MustConvertIntegerToUint32(ctx.BlockHeight()),
				false,
			); err != nil {
				return err
			}
			placedShortTermOrders[order.GetOrderId()] = order
		case *types.InternalOperation_OrderRemoval:
			orderRemoval := castedOperation.OrderRemoval

			if err := k.PersistOrderRemovalToState(ctx, *orderRemoval); err != nil {
				return errorsmod.Wrapf(
					types.ErrInvalidOrderRemoval,
					"Order Removal (%+v) invalid. Error: %+v",
					*orderRemoval,
					err,
				)
			}
		case *types.InternalOperation_PreexistingStatefulOrder:
			// When we fetch operations to propose, preexisting stateful orders are not included
			// in the operations queue.
			panic(
				fmt.Sprintf(
					"ProcessInternalOperations: Preexisting Stateful Orders should not exist in operations queue: %+v",
					castedOperation.PreexistingStatefulOrder,
				),
			)
		default:
			panic(
				fmt.Sprintf(
					"ProcessInternalOperations: Unrecognized operation type for operation: %+v",
					operation.GetInternalOperationTextString(),
				),
			)
		}
	}
	return nil
}

// PersistMatchToState takes in an ClobMatch and writes the match to state. A map of orderId
// to Order is required to fetch the whole Order object for short term orders.
func (k Keeper) PersistMatchToState(
	ctx sdk.Context,
	clobMatch *types.ClobMatch,
	ordersMap map[types.OrderId]types.Order,
	affiliateOverrides map[string]bool,
	affiliateParameters affiliatetypes.AffiliateParameters,
) error {
	switch castedMatch := clobMatch.Match.(type) {
	case *types.ClobMatch_MatchOrders:
		if err := k.PersistMatchOrdersToState(ctx, castedMatch.MatchOrders, ordersMap,
			affiliateOverrides, affiliateParameters); err != nil {
			return err
		}
	case *types.ClobMatch_MatchPerpetualLiquidation:
		if err := k.PersistMatchLiquidationToState(
			ctx,
			castedMatch.MatchPerpetualLiquidation,
			ordersMap,
			affiliateOverrides,
			affiliateParameters,
		); err != nil {
			return err
		}
	case *types.ClobMatch_MatchPerpetualDeleveraging:
		if err := k.PersistMatchDeleveragingToState(
			ctx,
			castedMatch.MatchPerpetualDeleveraging,
		); err != nil {
			return err
		}
	default:
		panic(
			fmt.Sprintf(
				"PersistMatchToState: Unrecognized operation type for match: %+v",
				clobMatch,
			),
		)
	}
	return nil
}

// statUnverifiedOrderRemoval increments the unverified order removal counter
// and the base quantums counter for the order to be removed.
func (k Keeper) statUnverifiedOrderRemoval(
	ctx sdk.Context,
	orderRemoval types.OrderRemoval,
) {
	proposerConsAddress := sdk.ConsAddress(ctx.BlockHeader().ProposerAddress)
	telemetry.IncrCounterWithLabels(
		[]string{types.ModuleName, metrics.ProcessOperations, metrics.UnverifiedStatefulOrderRemoval, metrics.Count},
		1,
		[]metrics.Label{
			metrics.GetLabelForStringValue(metrics.RemovalReason, orderRemoval.GetRemovalReason().String()),
			metrics.GetLabelForStringValue(metrics.Proposer, proposerConsAddress.String()),
		},
	)
}

// PersistOrderRemovalToState takes in an OrderRemoval, statefully validates it according to
// RemovalReason, and writes the removal to state.
func (k Keeper) PersistOrderRemovalToState(
	ctx sdk.Context,
	orderRemoval types.OrderRemoval,
) error {
	lib.AssertDeliverTxMode(ctx)
	orderIdToRemove := orderRemoval.GetOrderId()
	orderIdToRemove.MustBeStatefulOrder()

	// Order removals are always for long-term orders which must exist or conditional orders
	// which must be triggered.
	orderToRemove, err := k.FetchOrderFromOrderId(ctx, orderIdToRemove, nil)
	if err != nil {
		return err
	}

	// Statefully validate that the removal reason is valid.
	switch removalReason := orderRemoval.RemovalReason; removalReason {
	case types.OrderRemoval_REMOVAL_REASON_UNDERCOLLATERALIZED:
		k.statUnverifiedOrderRemoval(ctx, orderRemoval)
		// TODO (CLOB-877) - These validations are commented out because margin requirements can be non-linear.
		// For the collateralization check, use the remaining amount of the order that is resting on the book.
		// remainingAmount, hasRemainingAmount := k.MemClob.GetOrderRemainingAmount(ctx, orderToRemove)
		// if !hasRemainingAmount {
		// 	return types.ErrOrderFullyFilled
		// }

		// pendingOpenOrder := types.PendingOpenOrder{
		// 	RemainingQuantums: remainingAmount,
		// 	IsBuy:             orderToRemove.IsBuy(),
		// 	Subticks:          orderToRemove.GetOrderSubticks(),
		// 	ClobPairId:        orderToRemove.GetClobPairId(),
		// }

		// // Temporarily construct the subaccountOpenOrders with a single PendingOpenOrder.
		// subaccountOpenOrders := map[satypes.SubaccountId][]types.PendingOpenOrder{
		// 	orderIdToRemove.SubaccountId: {
		// 		pendingOpenOrder,
		// 	},
		// }

		// // TODO(DEC-1896): AddOrderToOrderbookSubaccountUpdatesCheck should accept a single PendingOpenOrder as a
		// // parameter rather than the subaccountOpenOrders map.
		// _, successPerSubaccountUpdate := k.AddOrderToOrderbookSubaccountUpdatesCheck(
		// 	ctx,
		// 	orderToRemove.GetClobPairId(),
		// 	subaccountOpenOrders,
		// )
		// if successPerSubaccountUpdate[orderIdToRemove.SubaccountId].IsSuccess() {
		// 	return errorsmod.Wrapf(
		// 		types.ErrInvalidOrderRemoval,
		// 		"Order Removal (%+v) invalid. Order passes collateralization check.",
		// 		orderRemoval,
		// 	)
		// }
	case types.OrderRemoval_REMOVAL_REASON_POST_ONLY_WOULD_CROSS_MAKER_ORDER:
		// TODO(CLOB-877)
		k.statUnverifiedOrderRemoval(ctx, orderRemoval)

		// The order should be post-only
		if orderToRemove.TimeInForce != types.Order_TIME_IN_FORCE_POST_ONLY {
			return errorsmod.Wrap(
				types.ErrUnexpectedTimeInForce,
				"Order is not post-only.",
			)
		}
	case types.OrderRemoval_REMOVAL_REASON_INVALID_SELF_TRADE:
		// TODO(CLOB-877)
		k.statUnverifiedOrderRemoval(ctx, orderRemoval)
	case types.OrderRemoval_REMOVAL_REASON_CONDITIONAL_FOK_COULD_NOT_BE_FULLY_FILLED:
		// TODO(CLOB-877)
		k.statUnverifiedOrderRemoval(ctx, orderRemoval)

		// The order should be FOK
		if orderToRemove.TimeInForce != types.Order_TIME_IN_FORCE_FILL_OR_KILL {
			return errorsmod.Wrap(
				types.ErrUnexpectedTimeInForce,
				"Order is not fill-or-kill.",
			)
		}

		// The order should not be fully filled.
		_, hasRemainingAmount := k.MemClob.GetOrderRemainingAmount(ctx, orderToRemove)
		if !hasRemainingAmount {
			return errorsmod.Wrap(
				types.ErrOrderFullyFilled,
				"Fill-or-kill order is fully filled.",
			)
		}
	case types.OrderRemoval_REMOVAL_REASON_CONDITIONAL_IOC_WOULD_REST_ON_BOOK:
		// TODO(CLOB-877)
		k.statUnverifiedOrderRemoval(ctx, orderRemoval)

		// The order should be IOC.
		if orderToRemove.TimeInForce != types.Order_TIME_IN_FORCE_IOC {
			return errorsmod.Wrap(
				types.ErrUnexpectedTimeInForce,
				"Order is not immediate-or-cancel.",
			)
		}

		// The order should not be fully filled.
		_, hasRemainingAmount := k.MemClob.GetOrderRemainingAmount(ctx, orderToRemove)
		if !hasRemainingAmount {
			return errorsmod.Wrapf(
				types.ErrOrderFullyFilled,
				"Immediate-or-cancel order is fully filled.",
			)
		}
	case types.OrderRemoval_REMOVAL_REASON_FULLY_FILLED:
		// Order removal reason fully filled is only used within indexer services.
		// Fully filled orders are removed from the protocol after persisting the operations queue
		// to state, instead of through the operations queue.
		return errorsmod.Wrapf(
			types.ErrInvalidOrderRemovalReason,
			"Order removal reason fully filled should not be part of the operations queue.",
		)
	case types.OrderRemoval_REMOVAL_REASON_INVALID_REDUCE_ONLY:
		if !orderToRemove.IsReduceOnly() {
			return errorsmod.Wrapf(
				types.ErrInvalidOrderRemoval,
				"Order Removal (%+v) invalid. Order must be reduce only.",
				orderRemoval,
			)
		}

		currentPositionSize := k.GetStatePosition(
			ctx,
			orderIdToRemove.SubaccountId,
			orderToRemove.GetClobPairId(),
		)

		// If the position is not fully closed, check that the order fill would increase the position size or change side.
		if currentPositionSize.Sign() != 0 {
			orderQuantumsToFill := orderToRemove.GetBigQuantums()
			orderFillWouldIncreasePositionSize := orderQuantumsToFill.Sign() == currentPositionSize.Sign()

			newPositionSize := new(big.Int).Add(currentPositionSize, orderQuantumsToFill)
			orderChangedSide := currentPositionSize.Sign()*newPositionSize.Sign() == -1
			if !orderFillWouldIncreasePositionSize && !orderChangedSide {
				return errorsmod.Wrapf(
					types.ErrInvalidOrderRemoval,
					"Order Removal (%+v) invalid. Order fill must increase position size or change side.",
					orderRemoval,
				)
			}
		}
	case types.OrderRemoval_REMOVAL_REASON_VIOLATES_ISOLATED_SUBACCOUNT_CONSTRAINTS:
		// TODO(CLOB-877)
		k.statUnverifiedOrderRemoval(ctx, orderRemoval)
	default:
		return errorsmod.Wrapf(
			types.ErrInvalidOrderRemovalReason,
			"PersistOrderRemovalToState: Unrecognized order removal type",
		)
	}

	// Remove the stateful order from state.
	k.MustRemoveStatefulOrder(ctx, orderIdToRemove)

	// Emit an on-chain indexer event for Stateful Order Removal.
	k.GetIndexerEventManager().AddTxnEvent(
		ctx,
		indexerevents.SubtypeStatefulOrder,
		indexerevents.StatefulOrderEventVersion,
		indexer_manager.GetBytes(
			indexerevents.NewStatefulOrderRemovalEvent(
				orderIdToRemove,
				indexershared.ConvertOrderRemovalReasonToIndexerOrderRemovalReason(
					orderRemoval.RemovalReason,
				),
			),
		),
	)

	telemetry.IncrCounterWithLabels(
		[]string{types.ModuleName, metrics.ProcessOperations, metrics.StatefulOrderRemoved, metrics.Count},
		1,
		append(
			orderIdToRemove.GetOrderIdLabels(),
			metrics.GetLabelForStringValue(metrics.RemovalReason, orderRemoval.GetRemovalReason().String()),
		),
	)
	return nil
}

// PersistMatchOrdersToState writes a MatchOrders object to state and emits an onchain
// indexer event for the match.
func (k Keeper) PersistMatchOrdersToState(
	ctx sdk.Context,
	matchOrders *types.MatchOrders,
	ordersMap map[types.OrderId]types.Order,
	affiliateOverrides map[string]bool,
	affiliateParameters affiliatetypes.AffiliateParameters,
) error {
	takerOrderId := matchOrders.GetTakerOrderId()
	// Fetch the taker order from either short term orders or state
	takerOrder, err := k.FetchOrderFromOrderId(ctx, takerOrderId, ordersMap)
	if err != nil {
		return err
	}

	// Taker order cannot be post only.
	if takerOrder.GetTimeInForce() == types.Order_TIME_IN_FORCE_POST_ONLY {
		return errorsmod.Wrapf(
			types.ErrInvalidMatchOrder,
			"Taker order %+v cannot be post only.",
			takerOrder.GetOrderTextString(),
		)
	}

	if takerOrder.RequiresImmediateExecution() {
		_, fillAmount, _ := k.GetOrderFillAmount(ctx, takerOrder.OrderId)
		if fillAmount != 0 {
			return errorsmod.Wrapf(
				types.ErrImmediateExecutionOrderAlreadyFilled,
				"Order %s",
				takerOrder.GetOrderTextString(),
			)
		}
	}

	makerOrders := make([]types.Order, 0)
	makerFills := matchOrders.GetFills()
	for _, makerFill := range makerFills {
		// Fetch the maker order from either short term orders or state.
		makerOrder, err := k.FetchOrderFromOrderId(ctx, makerFill.MakerOrderId, ordersMap)
		if err != nil {
			return err
		}

		matchWithOrders := types.MatchWithOrders{
			TakerOrder: &takerOrder,
			MakerOrder: &makerOrder,
			FillAmount: satypes.BaseQuantums(makerFill.GetFillAmount()),
		}
		makerOrders = append(makerOrders, makerOrder)

		_, _, _, affiliateRevSharesQuoteQuantums, err := k.ProcessSingleMatch(
			ctx,
			&matchWithOrders,
			affiliateOverrides,
			affiliateParameters,
		)
		if err != nil {
			return err
		}

		// Send on-chain update for the match. The events are stored in a TransientStore which should be rolled-back
		// if the branched state is discarded, so batching is not necessary.

		makerExists, totalFilledMaker, _ := k.GetOrderFillAmount(ctx, matchWithOrders.MakerOrder.MustGetOrder().OrderId)
		takerExists, totalFilledTaker, _ := k.GetOrderFillAmount(ctx, matchWithOrders.TakerOrder.MustGetOrder().OrderId)
		if !makerExists {
			panic(
				fmt.Sprintf("PersistMatchOrdersToState: Order fill amount not found for maker order: %+v",
					matchWithOrders.MakerOrder.MustGetOrder().OrderId,
				),
			)
		}
		if !takerExists {
			panic(
				fmt.Sprintf("PersistMatchOrdersToState: Order fill amount not found for taker order: %+v",
					matchWithOrders.TakerOrder.MustGetOrder().OrderId,
				),
			)
		}
		// TODO: (anmol) update fill event to include builder codes [CT-1363]
		k.GetIndexerEventManager().AddTxnEvent(
			ctx,
			indexerevents.SubtypeOrderFill,
			indexerevents.OrderFillEventVersion,
			indexer_manager.GetBytes(
				indexerevents.NewOrderFillEvent(
					matchWithOrders.MakerOrder.MustGetOrder(),
					matchWithOrders.TakerOrder.MustGetOrder(),
					matchWithOrders.FillAmount,
					matchWithOrders.MakerFee,
					matchWithOrders.TakerFee,
					matchWithOrders.MakerBuilderFee,
					matchWithOrders.TakerBuilderFee,
					totalFilledMaker,
					totalFilledTaker,
					affiliateRevSharesQuoteQuantums,
					matchWithOrders.MakerOrderRouterFee,
					matchWithOrders.TakerOrderRouterFee,
				),
			),
		)
	}

	// if GRPC streaming is on, emit a generated clob match to stream.
	if streamingManager := k.GetFullNodeStreamingManager(); streamingManager.Enabled() {
		// Note: GenerateStreamOrderbookFill doesn't rely on MemClob state.
		streamOrderbookFill := k.MemClob.GenerateStreamOrderbookFill(
			ctx,
			types.ClobMatch{
				Match: &types.ClobMatch_MatchOrders{
					MatchOrders: matchOrders,
				},
			},
			&takerOrder,
			makerOrders,
		)

		k.GetFullNodeStreamingManager().SendOrderbookFillUpdate(
			streamOrderbookFill,
			ctx,
			k.PerpetualIdToClobPairId,
		)
	}

	return nil
}

// PersistMatchLiquidationToState writes a MatchPerpetualLiquidation event and updates the keeper transient store.
// It also performs stateful validation on the matchLiquidations object.
func (k Keeper) PersistMatchLiquidationToState(
	ctx sdk.Context,
	matchLiquidation *types.MatchPerpetualLiquidation,
	ordersMap map[types.OrderId]types.Order,
	affiliateOverrides map[string]bool,
	affiliateParameters affiliatetypes.AffiliateParameters,
) error {
	// If the subaccount is not liquidatable, do nothing.
	if err := k.EnsureIsLiquidatable(ctx, matchLiquidation.Liquidated); err != nil {
		return err
	}

	takerOrder, err := k.GetLiquidationOrderForPerpetual(
		ctx,
		matchLiquidation.Liquidated,
		matchLiquidation.PerpetualId,
	)
	if err != nil {
		return err
	}

	// Perform stateless validation on the liquidation order.
	if err := k.ValidateLiquidationOrderAgainstProposedLiquidation(ctx, takerOrder, matchLiquidation); err != nil {
		return err
	}

	makerOrders := make([]types.Order, 0)
	for _, fill := range matchLiquidation.GetFills() {
		// Fetch the maker order from either short term orders or state.
		makerOrder, err := k.FetchOrderFromOrderId(ctx, fill.MakerOrderId, ordersMap)
		if err != nil {
			return err
		}
		makerOrders = append(makerOrders, makerOrder)

		matchWithOrders := types.MatchWithOrders{
			MakerOrder: &makerOrder,
			TakerOrder: takerOrder,
			FillAmount: satypes.BaseQuantums(fill.FillAmount),
		}

		// Write the position updates and state fill amounts for this match.
		// Note stateless validation on the constructed `matchWithOrders` is performed within this function.
		_, _, _, affiliateRevSharesQuoteQuantums, err := k.ProcessSingleMatch(
			ctx,
			&matchWithOrders,
			affiliateOverrides,
			affiliateParameters,
		)
		if err != nil {
			return err
		}

		makerExists, totalFilledMaker, _ := k.GetOrderFillAmount(ctx, matchWithOrders.MakerOrder.MustGetOrder().OrderId)
		if !makerExists {
			panic(
				fmt.Sprintf("PersistMatchLiquidationToState: Order fill amount not found for maker order: %+v",
					matchWithOrders.MakerOrder.MustGetOrder().OrderId,
				),
			)
		}

		// Send on-chain update for the liquidation. The events are stored in a TransientStore which should be rolled-back
		// if the branched state is discarded, so batching is not necessary.
		// There is potentially a maker builder fee for liquidations, but no taker builder fee since the protocol is always
		// the taker in the case of liquidations.
		k.GetIndexerEventManager().AddTxnEvent(
			ctx,
			indexerevents.SubtypeOrderFill,
			indexerevents.OrderFillEventVersion,
			indexer_manager.GetBytes(
				indexerevents.NewLiquidationOrderFillEvent(
					matchWithOrders.MakerOrder.MustGetOrder(),
					matchWithOrders.TakerOrder,
					matchWithOrders.FillAmount,
					matchWithOrders.MakerFee,
					matchWithOrders.TakerFee,
					matchWithOrders.MakerBuilderFee,
					totalFilledMaker,
					affiliateRevSharesQuoteQuantums,
					matchWithOrders.MakerOrderRouterFee,
				),
			),
		)
	}

	// Update the keeper transient store if-and-only-if the liquidation is valid.
	k.MustUpdateSubaccountPerpetualLiquidated(
		ctx,
		matchLiquidation.Liquidated,
		matchLiquidation.PerpetualId,
	)

	// if GRPC streaming is on, emit a generated clob match to stream.
	if streamingManager := k.GetFullNodeStreamingManager(); streamingManager.Enabled() {
		streamOrderbookFill := k.MemClob.GenerateStreamOrderbookFill(
			ctx,
			types.ClobMatch{
				Match: &types.ClobMatch_MatchPerpetualLiquidation{
					MatchPerpetualLiquidation: matchLiquidation,
				},
			},
			takerOrder,
			makerOrders,
		)
		k.GetFullNodeStreamingManager().SendOrderbookFillUpdate(
			streamOrderbookFill,
			ctx,
			k.PerpetualIdToClobPairId,
		)
	}
	return nil
}

// PersistMatchDeleveragingToState writes a MatchPerpetualDeleveraging object to state.
// This function returns an error if:
// - CanDeleverageSubaccount returns false for both boolean return values, indicating the
// subaccount failed deleveraging validation.
// - The IsFinalSettlement flag on the operation does not match the expected value based on collateralization
// and market status.
// - OffsetSubaccountPerpetualPosition returns an error.
// - The generated fills do not match the fills in the Operations object.
// TODO(CLOB-654) Verify deleveraging is triggered by unmatched liquidation orders and for the correct amount.
func (k Keeper) PersistMatchDeleveragingToState(
	ctx sdk.Context,
	matchDeleveraging *types.MatchPerpetualDeleveraging,
) error {
	liquidatedSubaccountId := matchDeleveraging.GetLiquidated()
	perpetualId := matchDeleveraging.GetPerpetualId()

	// Validate that the provided subaccount can be deleveraged.
	shouldDeleverageAtBankruptcyPrice, shouldDeleverageAtOraclePrice, err := k.CanDeleverageSubaccount(
		ctx,
		liquidatedSubaccountId,
		perpetualId,
	)
	if err != nil {
		panic(
			fmt.Sprintf(
				"PersistMatchDeleveragingToState: Failed to determine if subaccount can be deleveraged. "+
					"SubaccountId %+v, error %+v",
				liquidatedSubaccountId,
				err,
			),
		)
	}

	if !shouldDeleverageAtBankruptcyPrice && !shouldDeleverageAtOraclePrice {
		// TODO(CLOB-853): Add more verbose error logging about why deleveraging failed validation.
		return errorsmod.Wrapf(
			types.ErrInvalidDeleveragedSubaccount,
			"Subaccount %+v failed deleveraging validation",
			liquidatedSubaccountId,
		)
	}

	if matchDeleveraging.IsFinalSettlement != shouldDeleverageAtOraclePrice {
		// Throw error if the isFinalSettlement flag does not match the expected value. This prevents misuse or lack
		// of use of the isFinalSettlement flag. The isFinalSettlement flag should be set to true if-and-only-if the
		// subaccount has non-negative TNC and the market is in final settlement. Otherwise, it must be false.
		return errorsmod.Wrapf(
			types.ErrDeleveragingIsFinalSettlementFlagMismatch,
			"MatchPerpetualDeleveraging %+v has isFinalSettlement flag (%v), expected (%v)",
			matchDeleveraging,
			matchDeleveraging.IsFinalSettlement,
			shouldDeleverageAtOraclePrice,
		)
	}

	liquidatedSubaccount := k.subaccountsKeeper.GetSubaccount(ctx, liquidatedSubaccountId)
	position, exists := liquidatedSubaccount.GetPerpetualPositionForId(perpetualId)
	if !exists {
		return errorsmod.Wrapf(
			types.ErrNoOpenPositionForPerpetual,
			"Subaccount %+v does not have an open position for perpetual %+v",
			liquidatedSubaccountId,
			perpetualId,
		)
	}
	deltaBaseQuantumsIsNegative := position.GetIsLong()

	// If there are zero-fill deleveraging operations, this is a sentinel value to indicate a subaccount could not be
	// liquidated or deleveraged and still has negative equity. Mark the current block number in state to indicate a
	// negative TNC subaccount was seen.
	if len(matchDeleveraging.GetFills()) == 0 {
		if !shouldDeleverageAtBankruptcyPrice {
			return errorsmod.Wrap(
				types.ErrZeroFillDeleveragingForNonNegativeTncSubaccount,
				fmt.Sprintf(
					"PersistMatchDeleveragingToState: zero-fill deleveraging operation included for subaccount %+v"+
						" and perpetual %d but subaccount isn't negative TNC",
					liquidatedSubaccountId,
					perpetualId,
				),
			)
		}

		metrics.IncrCountMetricWithLabels(
			types.ModuleName,
			metrics.SubaccountsNegativeTncSubaccountSeen,
			metrics.GetLabelForIntValue(metrics.PerpetualId, int(perpetualId)),
			metrics.GetLabelForBoolValue(metrics.IsLong, position.GetIsLong()),
			metrics.GetLabelForBoolValue(metrics.DeliverTx, true),
		)
		if err = k.subaccountsKeeper.SetNegativeTncSubaccountSeenAtBlock(
			ctx,
			perpetualId,
			lib.MustConvertIntegerToUint32(ctx.BlockHeight()),
		); err != nil {
			return err
		}
		return nil
	}

	for _, fill := range matchDeleveraging.GetFills() {
		deltaBaseQuantums := new(big.Int).SetUint64(fill.FillAmount)
		if deltaBaseQuantumsIsNegative {
			deltaBaseQuantums.Neg(deltaBaseQuantums)
		}

		deltaQuoteQuantums, err := k.getDeleveragingQuoteQuantumsDelta(
			ctx,
			perpetualId,
			liquidatedSubaccountId,
			deltaBaseQuantums,
			matchDeleveraging.IsFinalSettlement,
		)
		if err != nil {
			return err
		}

		if err := k.ProcessDeleveraging(
			ctx,
			liquidatedSubaccountId,
			fill.OffsettingSubaccountId,
			perpetualId,
			deltaBaseQuantums,
			deltaQuoteQuantums,
		); err != nil {
			return errorsmod.Wrapf(
				types.ErrInvalidDeleveragingFill,
				"Failed to process deleveraging fill: %+v. liquidatedSubaccountId: %+v, "+
					"perpetualId: %v, deltaBaseQuantums: %v, deltaQuoteQuantums: %v, error: %v",
				fill,
				liquidatedSubaccountId,
				perpetualId,
				deltaBaseQuantums,
				deltaQuoteQuantums,
				err,
			)
		}

		// Send on-chain update for the deleveraging. The events are stored in a TransientStore which should be rolled-back
		// if the branched state is discarded, so batching is not necessary.
		k.GetIndexerEventManager().AddTxnEvent(
			ctx,
			indexerevents.SubtypeDeleveraging,
			indexerevents.DeleveragingEventVersion,
			indexer_manager.GetBytes(
				indexerevents.NewDeleveragingEvent(
					liquidatedSubaccountId,
					fill.OffsettingSubaccountId,
					perpetualId,
					satypes.BaseQuantums(new(big.Int).Abs(deltaBaseQuantums).Uint64()),
					satypes.BaseQuantums(deltaQuoteQuantums.Uint64()),
					deltaBaseQuantums.Sign() > 0,
					matchDeleveraging.IsFinalSettlement,
				),
			),
		)
		// if GRPC streaming is on, emit a generated clob match to stream.
		if streamingManager := k.GetFullNodeStreamingManager(); streamingManager.Enabled() {
			streamOrderbookFill := types.StreamOrderbookFill{
				ClobMatch: &types.ClobMatch{
					Match: &types.ClobMatch_MatchPerpetualDeleveraging{
						MatchPerpetualDeleveraging: matchDeleveraging,
					},
				},
			}
			k.SendOrderbookFillUpdate(
				ctx,
				streamOrderbookFill,
			)
		}
	}

	return nil
}

// GenerateProcessProposerMatchesEvents generates a `ProcessProposerMatchesEvents` object from
// an operations queue.
// Currently, it sets the `OrderIdsFilledInLastBlock` field and the `BlockHeight` field.
// This function expects the proposed operations to be valid, and does not verify that the `GoodTilBlockTime`
// of order replacement and cancellation is greater than the `GoodTilBlockTime` of the existing order.
func (k Keeper) GenerateProcessProposerMatchesEvents(
	ctx sdk.Context,
	operations []types.InternalOperation,
) types.ProcessProposerMatchesEvents {
	// Seen set for filled order ids
	seenOrderIdsFilledInLastBlock := make(map[types.OrderId]struct{}, 0)
	seenOrderIdsRemovedInLastBlock := make(map[types.OrderId]struct{}, 0)

	// Collect all filled order ids in this block.
	for _, operation := range operations {
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
		} else if operationRemoval := operation.GetOrderRemoval(); operationRemoval != nil {
			// For order removal, add order id to `seenOrderIdsRemovedInLastBlock`
			orderId := operationRemoval.GetOrderId()
			seenOrderIdsRemovedInLastBlock[orderId] = struct{}{}
		}
	}

	// Sort for deterministic ordering when writing to memstore.
	filledOrderIds := lib.GetSortedKeys[types.SortedOrders](seenOrderIdsFilledInLastBlock)
	removedOrderIds := lib.GetSortedKeys[types.SortedOrders](seenOrderIdsRemovedInLastBlock)

	// ConditionalOrderIdsTriggeredInLastBlock to be populated in EndBlocker.
	// ExpiredOrderId to be populated in the EndBlocker.
	return types.ProcessProposerMatchesEvents{
		ExpiredStatefulOrderIds:                 []types.OrderId{},
		OrderIdsFilledInLastBlock:               filledOrderIds,
		RemovedStatefulOrderIds:                 removedOrderIds,
		ConditionalOrderIdsTriggeredInLastBlock: []types.OrderId{},
		BlockHeight:                             lib.MustConvertIntegerToUint32(ctx.BlockHeight()),
	}
}
