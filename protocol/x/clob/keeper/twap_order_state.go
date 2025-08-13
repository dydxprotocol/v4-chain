package keeper

import (
	"errors"
	"fmt"
	"math/big"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	indexerevents "github.com/dydxprotocol/v4-chain/protocol/indexer/events"
	"github.com/dydxprotocol/v4-chain/protocol/indexer/indexer_manager"
	indexershared "github.com/dydxprotocol/v4-chain/protocol/indexer/shared/types"
	"github.com/dydxprotocol/v4-chain/protocol/lib"
	"github.com/dydxprotocol/v4-chain/protocol/lib/abci"
	"github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
)

// TWAP_SUBORDER_GOOD_TIL_BLOCK_TIME_OFFSET is the offset in seconds added to the
// current block time to set the good til block time for a suborder.
const TWAP_SUBORDER_GOOD_TIL_BLOCK_TIME_OFFSET = 3

// TWAP_MAX_SUBORDER_CATCHUP_MULTIPLE is the maximum multiple of the original
// quantums per leg that a suborder can be.
var TWAP_MAX_SUBORDER_CATCHUP_MULTIPLE = big.NewInt(3)

type twapOperationType int

const (
	parentTwapCompleted twapOperationType = iota
	parentTwapCancelled
	createSuborder
)

type twapOrderOperation struct {
	operationType      twapOperationType
	keyToDelete        []byte
	suborderToPlace    *types.Order
	twapOrderPlacement *types.TwapOrderPlacement
}

func (k Keeper) SetTWAPOrderPlacement(ctx sdk.Context,
	order types.Order,
	blockHeight uint32,
) {
	store := k.GetTWAPOrderPlacementStore(ctx)
	orderKey := order.OrderId.ToStateKey()

	total_legs := order.GetTotalLegsTWAPOrder()

	twapOrderPlacement := types.TwapOrderPlacement{
		Order:             order,
		RemainingLegs:     total_legs,
		RemainingQuantums: order.Quantums,
	}

	k.CheckAndIncrementStatefulOrderCount(ctx, order.OrderId)

	twapOrderPlacementBytes := k.cdc.MustMarshal(&twapOrderPlacement)
	store.Set(orderKey, twapOrderPlacementBytes)

	k.AddSuborderToTriggerStore(ctx, k.twapToSuborderId(order.OrderId), 0)
}

// GetTwapOrderPlacement gets a TWAP order placement from the store.
// Returns false if no TWAP order exists in store with `orderId`.
func (k Keeper) GetTwapOrderPlacement(
	ctx sdk.Context,
	orderId types.OrderId,
) (val types.TwapOrderPlacement, found bool) {
	store := k.GetTWAPOrderPlacementStore(ctx)

	b := store.Get(orderId.ToStateKey())
	if b == nil {
		return val, false
	}

	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

// GetTwapTriggerPlacement gets a TWAP trigger placement for a given orderId.
// Returns false if no trigger placement exists in store with `orderId`.
// This iterates over the entire store because the keys in the store are
// formatted as [timestamp, orderId]. This is primarily used for testing.
func (k Keeper) GetTwapTriggerPlacement(
	ctx sdk.Context,
	orderId types.OrderId,
) (o types.OrderId, t int64, found bool) {
	store := k.GetTWAPTriggerOrderPlacementStore(ctx)

	iterator := store.Iterator(nil, nil)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var suborderId types.OrderId
		k.cdc.MustUnmarshal(iterator.Key()[8:], &suborderId)

		timestamp := types.TimeFromTriggerKey(iterator.Key())
		if suborderId == orderId {
			return suborderId, timestamp, true
		}
	}
	return types.OrderId{}, 0, false
}

func (k Keeper) twapToSuborderId(twapOrderId types.OrderId) types.OrderId {
	return types.OrderId{
		SubaccountId: twapOrderId.SubaccountId,
		ClientId:     twapOrderId.ClientId,
		OrderFlags:   types.OrderIdFlags_TwapSuborder,
		ClobPairId:   twapOrderId.ClobPairId,
	}
}

// AddSuborderToTriggerStore adds a suborder to the trigger store with the
// binary encoded [timestamp][suborderId] key.
func (k Keeper) AddSuborderToTriggerStore(
	ctx sdk.Context,
	suborderId types.OrderId,
	triggerOffset int64,
) []byte {
	triggerStore := k.GetTWAPTriggerOrderPlacementStore(ctx)
	triggerTime := ctx.BlockTime().Unix() + triggerOffset

	triggerKey := types.GetTWAPTriggerKey(triggerTime, suborderId)

	// The value in the map is not used, so we can set it to an empty byte slice.
	triggerStore.Set(triggerKey, []byte{})
	return triggerKey
}

func (k Keeper) DeleteSuborderFromTriggerStore(
	ctx sdk.Context,
	triggerKey []byte,
) {
	triggerStore := k.GetTWAPTriggerOrderPlacementStore(ctx)
	triggerStore.Delete(triggerKey)
}

// GenerateAndPlaceTriggeredTwapSuborders will iterate over the twap trigger
// store and generate the suborders that need to be placed on the orderbook
// based on the timestamps their triggers were set for.
// It also has the responsibility of checking if the parent twap was
// cancelled, in which case no subsequent suborders should be placed, and
// also the case that the parent twap order was completed (no remaining legs),
// in which case the parent twap order should be removed from the store.
func (k Keeper) GenerateAndPlaceTriggeredTwapSuborders(ctx sdk.Context) {
	triggerStore := k.GetTWAPTriggerOrderPlacementStore(ctx)
	blockTime := ctx.BlockTime().Unix()
	var operationsToProcess []twapOrderOperation

	iterator := triggerStore.Iterator(nil, nil)
	for ; iterator.Valid(); iterator.Next() {
		var orderId types.OrderId
		k.cdc.MustUnmarshal(iterator.Key()[8:], &orderId)

		triggerTime := types.TimeFromTriggerKey(iterator.Key())

		if triggerTime > blockTime {
			break // all remaining suborders are in the future
		}

		parentOrderId := types.OrderId{
			SubaccountId: orderId.SubaccountId,
			ClientId:     orderId.ClientId,
			OrderFlags:   types.OrderIdFlags_Twap, // Set directly to TWAP
			ClobPairId:   orderId.ClobPairId,
		}

		twapOrderPlacement, found := k.GetTwapOrderPlacement(ctx, parentOrderId)
		if !found {
			// If parent TWAP was cancelled/not found, do not place any pending suborders.
			operationsToProcess = append(operationsToProcess, twapOrderOperation{
				operationType: parentTwapCancelled,
				keyToDelete:   append([]byte{}, iterator.Key()...),
			})
			continue
		}

		operationType := createSuborder
		order, isGenerated := k.GenerateSuborder(ctx, orderId, twapOrderPlacement, blockTime)
		if !isGenerated {
			operationType = parentTwapCompleted
		}
		operationsToProcess = append(operationsToProcess, twapOrderOperation{
			operationType:      operationType,
			keyToDelete:        append([]byte{}, iterator.Key()...),
			suborderToPlace:    order,
			twapOrderPlacement: &twapOrderPlacement,
		})
	}
	iterator.Close()

	for _, op := range operationsToProcess {
		// Delete from trigger store
		triggerStore.Delete(op.keyToDelete)

		switch op.operationType {
		case parentTwapCancelled:
			// no-op after trigger key has been deleted
		case parentTwapCompleted:
			k.DeleteTWAPOrderPlacement(ctx, op.twapOrderPlacement.Order.GetOrderId())
		case createSuborder:
			// decrement remaining legs
			k.DecrementTwapOrderRemainingLegs(ctx, *op.twapOrderPlacement)
			// add the next suborder to the trigger store
			triggerKey := k.AddSuborderToTriggerStore(
				ctx,
				op.suborderToPlace.OrderId,
				int64(op.twapOrderPlacement.Order.TwapParameters.Interval),
			)

			// place triggered suborder
			err := k.safeHandleMsgPlaceOrder(ctx, &types.MsgPlaceOrder{Order: *op.suborderToPlace}, true)
			if err != nil {
				k.DeleteTWAPOrderPlacement(ctx, op.twapOrderPlacement.Order.GetOrderId())
				k.DeleteSuborderFromTriggerStore(ctx, triggerKey)

				k.GetIndexerEventManager().AddTxnEvent(
					ctx,
					indexerevents.SubtypeStatefulOrder,
					indexerevents.StatefulOrderEventVersion,
					indexer_manager.GetBytes(
						indexerevents.NewStatefulOrderRemovalEvent(
							op.twapOrderPlacement.Order.GetOrderId(),
							getTwapOrderRemovalReason(err),
						),
					),
				)
			}
		default:
			k.Logger(ctx).Error(
				"unsupported twap operation type can not be processed",
				"operationType", op.operationType,
			)
		}
	}
}

func (k Keeper) DeleteTWAPOrderPlacement(
	ctx sdk.Context,
	orderId types.OrderId,
) {
	// Decrement the stateful order count for the TWAP order.
	k.CheckAndDecrementStatefulOrderCount(ctx, orderId)

	store := k.GetTWAPOrderPlacementStore(ctx)
	orderKey := orderId.ToStateKey()
	store.Delete(orderKey)
}

func (k Keeper) DecrementTwapOrderRemainingLegs(
	ctx sdk.Context,
	twapOrderPlacement types.TwapOrderPlacement,
) {
	store := k.GetTWAPOrderPlacementStore(ctx)
	orderKey := twapOrderPlacement.Order.OrderId.ToStateKey()
	if twapOrderPlacement.IsCompleted() {
		k.Logger(ctx).Error(
			"twap order has already been completed",
			"orderId", twapOrderPlacement.Order.OrderId,
		)
		return
	}

	twapOrderPlacement.RemainingLegs--
	twapOrderPlacementBytes := k.cdc.MustMarshal(&twapOrderPlacement)
	store.Set(orderKey, twapOrderPlacementBytes)
}

// UpdateTWAPOrderRemainingQuantityOnFill updates the remaining quantity of the
// parent twap order after a suborder has been filled.
func (k Keeper) UpdateTWAPOrderRemainingQuantityOnFill(
	ctx sdk.Context,
	parentOrderId types.OrderId,
	filledQuantums uint64,
) error {
	// Get the TWAP order placement
	twapOrderPlacement, found := k.GetTwapOrderPlacement(ctx, parentOrderId)
	if !found {
		return errorsmod.Wrapf(
			types.ErrInvalidTwapOrderPlacement,
			"TWAP order %+v does not exist",
			parentOrderId,
		)
	}

	// Update remaining quantums
	twapOrderPlacement.RemainingQuantums -= filledQuantums

	// Store updated state
	store := k.GetTWAPOrderPlacementStore(ctx)
	orderKey := parentOrderId.ToStateKey()
	twapOrderPlacementBytes := k.cdc.MustMarshal(&twapOrderPlacement)
	store.Set(orderKey, twapOrderPlacementBytes)

	return nil
}

// GetSubticksForTWAPSuborder returns the oracle price in subticks
// adjusted by a given directional price tolerance in ppm, rounded to the nearest multiple
// of SubticksPerTick. A positive price tolerance increases the price, while a negative price
// tolerance decreases it. If a TWAP order has a subticks value, it will always be used
// instead of the adjusted oracle price.
//
// For example:
//   - price tolerance = 500_000 means 50% higher than oracle price
//   - price tolerance = -500_000 means 50% lower than oracle price
func (k Keeper) calculateSuborderSubticks(
	ctx sdk.Context,
	clobPair types.ClobPair,
	twapOrderPlacement types.TwapOrderPlacement,
) uint64 {
	if twapOrderPlacement.Order.Subticks != 0 {
		return twapOrderPlacement.Order.Subticks
	}

	oraclePriceSubticksRat := k.GetOraclePriceSubticksRat(ctx, clobPair)

	twapOrder := twapOrderPlacement.Order
	priceTolerancePpm := int32(twapOrder.TwapParameters.PriceTolerance)
	if twapOrder.Side == types.Order_SIDE_SELL {
		// for sell orders, we want to adjust the price down
		priceTolerancePpm = -priceTolerancePpm
	}
	adjustment := int32(1_000_000) + priceTolerancePpm

	adjustedPrice := lib.BigRatMulPpm(oraclePriceSubticksRat, uint32(adjustment))
	// Round to the nearest multiple of SubticksPerTick
	roundedSubticks := lib.BigRatRoundToMultiple(
		adjustedPrice,
		new(big.Int).SetUint64(uint64(clobPair.SubticksPerTick)),
		priceTolerancePpm >= 0, // round up for positive adjustments, down for negative
	)

	return roundedSubticks.Uint64()
}

func (k Keeper) calculateSuborderQuantums(
	twapOrderPlacement types.TwapOrderPlacement,
	clobPair types.ClobPair,
) uint64 {
	totalLegs := twapOrderPlacement.Order.GetTotalLegsTWAPOrder()
	originalQuantums := twapOrderPlacement.Order.Quantums
	originalQuantumsPerLeg := lib.BigDivCeil(lib.BigU(originalQuantums), lib.BigU(totalLegs))

	// Calculate the quantums for the suborder capping at 3x the original quantums per leg
	remainingQuantums := twapOrderPlacement.RemainingQuantums
	remainingLegs := twapOrderPlacement.RemainingLegs
	remainingQuantumsPerLeg := lib.BigDivCeil(lib.BigU(remainingQuantums), lib.BigU(remainingLegs))

	maxSuborderSize := new(big.Int).Mul(
		originalQuantumsPerLeg,
		TWAP_MAX_SUBORDER_CATCHUP_MULTIPLE,
	)

	suborderQuantums := lib.BigMin(
		remainingQuantumsPerLeg,
		maxSuborderSize,
	)

	// Round down to nearest multiple of StepBaseQuantums
	quantumsByStepBaseQuantums := lib.BigDivFloor(suborderQuantums, lib.BigU(clobPair.StepBaseQuantums))
	if quantumsByStepBaseQuantums.Uint64() == 0 {
		return 0
	}

	suborderQuantumsRounded := quantumsByStepBaseQuantums.Mul(
		quantumsByStepBaseQuantums,
		lib.BigU(clobPair.StepBaseQuantums),
	)
	return suborderQuantumsRounded.Uint64()
}

func getTwapOrderRemovalReason(err error) indexershared.OrderRemovalReason {
	reason := indexershared.OrderRemovalReason_ORDER_REMOVAL_REASON_UNSPECIFIED
	if errors.Is(err, types.ErrStatefulOrderCollateralizationCheckFailed) {
		reason = indexershared.OrderRemovalReason_ORDER_REMOVAL_REASON_UNDERCOLLATERALIZED
	} else if errors.Is(err, types.ErrWouldViolateIsolatedSubaccountConstraints) {
		reason = indexershared.OrderRemovalReason_ORDER_REMOVAL_REASON_VIOLATES_ISOLATED_SUBACCOUNT_CONSTRAINTS
	}

	return reason
}

// GenerateSuborder generates a suborder when it has been triggered via the
// trigger store. The suborderId is given  by the store, and this method
// generates the remaining required fields. Configured price tolerance is
// applied to the oracle price to get the suborder price. Quantity is determined
// as a function of the remaining quantums and legs. Good til block time is set
// to the current block time plus the protocol defined offset.
func (k Keeper) GenerateSuborder(
	ctx sdk.Context,
	suborderId types.OrderId,
	twapOrderPlacement types.TwapOrderPlacement,
	blockTime int64,
) (*types.Order, bool) {
	if twapOrderPlacement.IsCompleted() {
		return nil, false
	}

	parentOrder := twapOrderPlacement.Order
	order := types.Order{
		OrderId:        suborderId,
		Side:           twapOrderPlacement.Order.Side,
		ReduceOnly:     twapOrderPlacement.Order.ReduceOnly,
		ClientMetadata: twapOrderPlacement.Order.ClientMetadata,
	}

	// calculate the suborder price with slippage adjustment
	clobPair := k.mustGetClobPair(ctx, parentOrder.GetClobPairId())

	// set the subticks based on the adjusted price and the limit price (if configured)
	// by the parent twap order
	order.Subticks = k.calculateSuborderSubticks(ctx, clobPair, twapOrderPlacement)

	// calculate the suborder quantums based on remaining quantums and legs
	order.Quantums = k.calculateSuborderQuantums(twapOrderPlacement, clobPair)

	// set good til block time from current block time
	order.GoodTilOneof = &types.Order_GoodTilBlockTime{
		GoodTilBlockTime: uint32(blockTime + TWAP_SUBORDER_GOOD_TIL_BLOCK_TIME_OFFSET),
	}

	return &order, true
}

// safeHandleMsgPlaceOrder safely calls HandleMsgPlaceOrder with panic recovery.
// This is used in end blockers where panics should be caught and logged rather than
// causing the entire block to fail.
func (k Keeper) safeHandleMsgPlaceOrder(
	ctx sdk.Context,
	msg *types.MsgPlaceOrder,
	isStateful bool,
) (err error) {
	if err = abci.RunCached(ctx, func(ctx sdk.Context) error {
		return k.HandleMsgPlaceOrder(ctx, msg, isStateful)
	}); err != nil {
		k.Logger(ctx).Error(
			"failed to handle TWAP suborder placement via HandleMsgPlaceOrder (panic recovered or error)",
			"cause", err,
			"orderId", msg.GetOrder().OrderId,
			"isStateful", isStateful,
			"stack", fmt.Sprintf("%+v", err),
		)

		return err
	}

	return nil
}
