package keeper

import (
	"fmt"
	"math/big"
	"time"

	gometrics "github.com/armon/go-metrics"
	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	indexerevents "github.com/dydxprotocol/v4-chain/protocol/indexer/events"
	"github.com/dydxprotocol/v4-chain/protocol/indexer/indexer_manager"
	indexershared "github.com/dydxprotocol/v4-chain/protocol/indexer/shared"
	"github.com/dydxprotocol/v4-chain/protocol/lib"
	"github.com/dydxprotocol/v4-chain/protocol/lib/metrics"
	"github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
	satypes "github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"
)

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

	// Validates RawOperations and transforms them into InternalOperations to be used internally by memclob.
	operations, err := types.ValidateAndTransformRawOperations(ctx, rawOperations, k.txDecoder, k.antehandler)
	if err != nil {
		return sdkerrors.Wrapf(types.ErrInvalidMsgProposedOperations, "Error: %+v", err)
	}

	ctx.Logger().Debug(
		"Processing operations queue",
		"operationsQueue",
		types.GetInternalOperationsQueueTextString(operations),
		"block",
		ctx.BlockHeight(),
	)

	// Write results of the operations queue to state.
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
// operations to state. This function assumes that the operations have passed all validation.
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
	placedShortTermOrders := make(map[types.OrderId]types.Order, 0)

	// Write the matches to state if all stateful validation passes.
	for _, operation := range operations {
		switch castedOperation := operation.Operation.(type) {
		case *types.InternalOperation_Match:
			clobMatch := castedOperation.Match
			if err := k.PersistMatchToState(ctx, clobMatch, placedShortTermOrders); err != nil {
				return sdkerrors.Wrapf(
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
			// Remove the stateful order from state.
			// TODO(CLOB-85): Perform additional validation on the order removal.
			orderRemoval := castedOperation.OrderRemoval
			k.MustRemoveStatefulOrder(ctx, orderRemoval.GetOrderId())

			// Emit an on-chain indexer event for Stateful Order Removal.
			k.GetIndexerEventManager().AddTxnEvent(
				ctx,
				indexerevents.SubtypeStatefulOrder,
				indexer_manager.GetB64EncodedEventMessage(
					indexerevents.NewStatefulOrderRemovalEvent(
						orderRemoval.OrderId,
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
					orderRemoval.OrderId.GetOrderIdLabels(),
					metrics.GetLabelForStringValue(metrics.RemovalReason, orderRemoval.GetRemovalReason().String()),
				),
			)
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
// Function will panic if any orderId referenced in clobMatch cannot be found.
func (k Keeper) PersistMatchToState(
	ctx sdk.Context,
	clobMatch *types.ClobMatch,
	ordersMap map[types.OrderId]types.Order,
) error {
	switch castedMatch := clobMatch.Match.(type) {
	case *types.ClobMatch_MatchOrders:
		if err := k.PersistMatchOrdersToState(ctx, castedMatch.MatchOrders, ordersMap); err != nil {
			return err
		}
	case *types.ClobMatch_MatchPerpetualLiquidation:
		if err := k.PersistMatchLiquidationToState(
			ctx,
			castedMatch.MatchPerpetualLiquidation,
			ordersMap,
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

// PersistMatchOrdersToState writes a MatchOrders object to state and emits an onchain
// indexer event for the match.
// Function will panic if any orderId referenced in matchOrders cannot be found.
func (k Keeper) PersistMatchOrdersToState(
	ctx sdk.Context,
	matchOrders *types.MatchOrders,
	ordersMap map[types.OrderId]types.Order,
) error {
	makerFills := matchOrders.GetFills()
	takerOrder := k.MustFetchOrderFromOrderId(ctx, matchOrders.GetTakerOrderId(), ordersMap)

	for _, makerFill := range makerFills {
		makerOrder := k.MustFetchOrderFromOrderId(ctx, makerFill.GetMakerOrderId(), ordersMap)

		matchWithOrders := types.MatchWithOrders{
			TakerOrder: &takerOrder,
			MakerOrder: &makerOrder,
			FillAmount: satypes.BaseQuantums(makerFill.GetFillAmount()),
		}

		_, _, _, _, err := k.ProcessSingleMatch(ctx, &matchWithOrders)
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
		k.GetIndexerEventManager().AddTxnEvent(
			ctx,
			indexerevents.SubtypeOrderFill,
			indexer_manager.GetB64EncodedEventMessage(
				indexerevents.NewOrderFillEvent(
					matchWithOrders.MakerOrder.MustGetOrder(),
					matchWithOrders.TakerOrder.MustGetOrder(),
					matchWithOrders.FillAmount,
					matchWithOrders.MakerFee,
					matchWithOrders.TakerFee,
					totalFilledMaker,
					totalFilledTaker,
				),
			),
		)
	}

	return nil
}

// PersistMatchLiquidationToState writes a MatchPerpetualLiquidation event and updates the keeper transient store.
// Function will panic if any orderId referenced in matchOrders cannot be matchLiquidation.
func (k Keeper) PersistMatchLiquidationToState(
	ctx sdk.Context,
	matchLiquidation *types.MatchPerpetualLiquidation,
	ordersMap map[types.OrderId]types.Order,
) error {
	isLiquidatable, err := k.IsLiquidatable(ctx, matchLiquidation.Liquidated)
	if err != nil {
		return err
	}
	if !isLiquidatable {
		return sdkerrors.Wrapf(
			types.ErrSubaccountNotLiquidatable,
			"PersistMatchLiquidationToState: Subaccount %s is not liquidatable",
			matchLiquidation.Liquidated,
		)
	}

	takerOrder, err := k.ConstructTakerOrderFromMatchPerpetualLiquidation(ctx, matchLiquidation)
	if err != nil {
		return err
	}

	for _, fill := range matchLiquidation.GetFills() {
		makerOrderId := fill.GetMakerOrderId()
		makerOrder := k.MustFetchOrderFromOrderId(ctx, makerOrderId, ordersMap)

		matchWithOrders := types.MatchWithOrders{
			MakerOrder: &makerOrder,
			TakerOrder: takerOrder,
			FillAmount: satypes.BaseQuantums(fill.FillAmount),
		}

		if err := matchWithOrders.Validate(); err != nil {
			return err
		}

		// Write the position updates and state fill amounts for this match.
		_, _, _, _, err := k.ProcessSingleMatch(
			ctx,
			&matchWithOrders,
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
		k.GetIndexerEventManager().AddTxnEvent(
			ctx,
			indexerevents.SubtypeOrderFill,
			indexer_manager.GetB64EncodedEventMessage(
				indexerevents.NewLiquidationOrderFillEvent(
					matchWithOrders.MakerOrder.MustGetOrder(),
					matchWithOrders.TakerOrder,
					matchWithOrders.FillAmount,
					matchWithOrders.MakerFee,
					matchWithOrders.TakerFee,
					totalFilledMaker,
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
	return nil
}

// PersistMatchDeleveragingToState writes a MatchPerpetualDeleveraging object to state.
// This function returns an error if:
// - OffsetSubaccountPerpetualPosition returns an error.
// - The generated fills do not match the fills in the Operations object.
// TODO(CLOB-654) Verify deleveraging is triggered by liquidation orders and for the correct amount.
func (k Keeper) PersistMatchDeleveragingToState(
	ctx sdk.Context,
	matchDeleveraging *types.MatchPerpetualDeleveraging,
) error {
	liquidatedSubaccountId := matchDeleveraging.GetLiquidated()

	isLiquidatable, err := k.IsLiquidatable(ctx, liquidatedSubaccountId)
	if err != nil {
		panic(
			fmt.Sprintf(
				"PersistMatchDeleveragingToState: Failed to determine if subaccount is liquidatable. "+
					"SubaccountId %v, error %s",
				liquidatedSubaccountId,
				err,
			),
		)
	}

	if !isLiquidatable {
		return sdkerrors.Wrapf(
			types.ErrDeleveragedSubaccountNotLiquidatable,
			"Subaccount %+v is not liquidatable",
			liquidatedSubaccountId,
		)
	}

	perpetualId := matchDeleveraging.GetPerpetualId()

	liquidatedSubaccount := k.subaccountsKeeper.GetSubaccount(ctx, liquidatedSubaccountId)
	position, exists := liquidatedSubaccount.GetPerpetualPositionForId(perpetualId)
	if !exists {
		return sdkerrors.Wrapf(
			types.ErrNoOpenPositionForPerpetual,
			"Subaccount %+v does not have an open position for perpetual %+v",
			liquidatedSubaccountId,
			perpetualId,
		)
	}
	deltaQuantumsIsNegative := position.GetIsLong()

	telemetry.IncrCounterWithLabels(
		[]string{types.ModuleName, metrics.Deleveraging, metrics.DeltaQuoteQuantums},
		1,
		[]gometrics.Label{
			metrics.GetLabelForBoolValue(metrics.Positive, !deltaQuantumsIsNegative),
		},
	)

	for _, fill := range matchDeleveraging.GetFills() {
		deltaQuantums := new(big.Int).SetUint64(fill.FillAmount)
		if deltaQuantumsIsNegative {
			deltaQuantums.Neg(deltaQuantums)
		}

		if err := k.ProcessDeleveraging(
			ctx,
			liquidatedSubaccountId,
			fill.OffsettingSubaccountId,
			perpetualId,
			deltaQuantums,
		); err != nil {
			return sdkerrors.Wrapf(
				types.ErrInvalidDeleveragingFill,
				"Failed to process deleveraging fill: %+v. liquidatedSubaccountId: %+v, "+
					"perpetualId: %v, deltaQuantums: %v, error: %v",
				fill,
				liquidatedSubaccountId,
				perpetualId,
				deltaQuantums,
				err,
			)
		}
	}

	return nil
}

// MustFetchOrderFromOrderId fetches an Order object given an orderId. If it is a short term order,
// `ordersMap` will be used to populate the order. If it is a stateful order, read from state.
// Note that this function is meant to be used for operation processing during DeliverTx and does not
// fetch untriggered conditional orders.
//
// Function will panic if for any reason, the order cannot be searched up.
func (k Keeper) MustFetchOrderFromOrderId(
	ctx sdk.Context,
	orderId types.OrderId,
	ordersMap map[types.OrderId]types.Order,
) types.Order {
	// In the case of short term orders, fetch from the orders map.
	// It should always exist in the orders map because short term order placement operations
	// should precede operations with order ids that reference them.
	if orderId.IsShortTermOrder() {
		order, exists := ordersMap[orderId]
		if !exists {
			panic(
				fmt.Sprintf("MustFetchOrderFromOrderId: failed fetching short term order for order id: %+v", orderId),
			)
		}
		return order
	}

	// For stateful orders, fetch from state.
	if orderId.IsLongTermOrder() {
		statefulOrderPlacement, found := k.GetLongTermOrderPlacement(ctx, orderId)
		if !found {
			panic(
				fmt.Sprintf("MustFetchOrderFromOrderId: failed fetching stateful term order for order id: %+v", orderId),
			)
		}
		return statefulOrderPlacement.Order
	} else if orderId.IsConditionalOrder() {
		conditionalOrderPlacement, found := k.GetTriggeredConditionalOrderPlacement(ctx, orderId)
		if !found {
			panic(
				fmt.Sprintf("MustFetchOrderFromOrderId: failed fetching triggered conditional order for order id: %+v", orderId),
			)
		}
		return conditionalOrderPlacement.Order
	}

	panic(
		fmt.Sprintf("MustFetchOrderFromOrderId: unknown order type. order id: %+v", orderId),
	)
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
	filledOrderIds := lib.ConvertMapToSliceOfKeys(seenOrderIdsFilledInLastBlock)
	removedOrderIds := lib.ConvertMapToSliceOfKeys(seenOrderIdsRemovedInLastBlock)
	// Sort for deterministic ordering when writing to memstore.
	types.MustSortAndHaveNoDuplicates(filledOrderIds)
	types.MustSortAndHaveNoDuplicates(removedOrderIds)

	// PlacedLongTermOrderIds to be populated in MsgHandler for MsgPlaceOrder.
	// PlacedConditionalOrderIds to be populated in MsgHandler for MsgPlaceOrder.
	// ConditionalOrderIdsTriggeredInLastBlock to be populated in EndBlocker.
	// ExpiredOrderId to be populated in the EndBlocker.
	// PlacedStatefulCancellation to be populated in MsgHandler for MsgCancelOrder.
	return types.ProcessProposerMatchesEvents{
		PlacedLongTermOrderIds:                  []types.OrderId{},
		ExpiredStatefulOrderIds:                 []types.OrderId{},
		OrderIdsFilledInLastBlock:               filledOrderIds,
		PlacedStatefulCancellationOrderIds:      []types.OrderId{},
		RemovedStatefulOrderIds:                 removedOrderIds,
		PlacedConditionalOrderIds:               []types.OrderId{},
		ConditionalOrderIdsTriggeredInLastBlock: []types.OrderId{},
		BlockHeight:                             lib.MustConvertIntegerToUint32(ctx.BlockHeight()),
	}
}
