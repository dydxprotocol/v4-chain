package keeper

import (
	"math"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4-chain/protocol/lib"
	"github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
)

const (
	TWAP_SUBORDER_GOOD_TIL_BLOCK_TIME_OFFSET = 3

	TWAP_MAX_SUBORDER_CATCHUP_MULTIPLE = 3
)

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

	k.AddSuborderToTriggerStore(ctx, k.twapToSuborderId(order.OrderId), 0)

	twapOrderPlacementBytes := k.cdc.MustMarshal(&twapOrderPlacement)
	store.Set(orderKey, twapOrderPlacementBytes)
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
// formatted as [timestamp, orderId].
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
) {
	triggerStore := k.GetTWAPTriggerOrderPlacementStore(ctx)
	triggerTime := ctx.BlockTime().Unix() + triggerOffset

	triggerKey := types.GetTWAPTriggerKey(triggerTime, suborderId)

	// The value in the map is not used, so we can set it to an empty byte slice.
	triggerStore.Set(triggerKey, []byte{})
}

func (k Keeper) GenerateAndPlaceTriggeredTwapSuborders(ctx sdk.Context) {
	triggerStore := k.GetTWAPTriggerOrderPlacementStore(ctx)
	blockTime := ctx.BlockTime().Unix()
	iterator := triggerStore.Iterator(nil, nil)
	defer iterator.Close()

	var operationsToProcess []struct {
		keyToDelete        []byte
		suborderToPlace    types.Order
		twapOrderPlacement types.TwapOrderPlacement
	}

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
			panic("parent twap order not found") // TODO: (anmol) handle order cancellation
		}

		order := k.GenerateSuborder(ctx, orderId, twapOrderPlacement, blockTime)

		operationsToProcess = append(operationsToProcess, struct {
			keyToDelete        []byte
			suborderToPlace    types.Order
			twapOrderPlacement types.TwapOrderPlacement
		}{
			keyToDelete:        append([]byte{}, iterator.Key()...),
			suborderToPlace:    order,
			twapOrderPlacement: twapOrderPlacement,
		})
	}
	iterator.Close()

	for _, op := range operationsToProcess {
		// Delete from trigger store
		triggerStore.Delete(op.keyToDelete)

		// decrement remaining legs
		k.DecrementTwapOrderRemainingLegs(ctx, &op.twapOrderPlacement)

		if op.twapOrderPlacement.RemainingLegs == 0 {
			// remove the parent twap order from the store
			store := k.GetTWAPOrderPlacementStore(ctx)
			orderKey := op.twapOrderPlacement.Order.OrderId.ToStateKey()
			store.Delete(orderKey)
			// TODO: (anmol) handle missing parent order case
			// TODO: (anmol) emit event?
			continue
		}

		// add the next suborder to the trigger store
		k.AddSuborderToTriggerStore(
			ctx,
			op.suborderToPlace.OrderId,
			int64(op.twapOrderPlacement.Order.TwapParameters.Interval),
		)

		// place triggered suborder
		err := k.HandleMsgPlaceOrder(ctx, &types.MsgPlaceOrder{Order: op.suborderToPlace}, true)
		if err != nil {
			continue // TODO: (anmol) handle suborder placement failure
		}
	}
}

func (k Keeper) DecrementTwapOrderRemainingLegs(
	ctx sdk.Context,
	twapOrderPlacement *types.TwapOrderPlacement,
) {
	if twapOrderPlacement.RemainingLegs == 0 {
		return // TODO: (anmol) handle end of twap order case
	}

	twapOrderPlacement.RemainingLegs--

	// Store updated state
	store := k.GetTWAPOrderPlacementStore(ctx)
	orderKey := twapOrderPlacement.Order.OrderId.ToStateKey()
	twapOrderPlacementBytes := k.cdc.MustMarshal(twapOrderPlacement)
	store.Set(orderKey, twapOrderPlacementBytes)
}

func (k Keeper) UpdateTWAPOrderStateOnSuborderFill(
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

func (k Keeper) calculateSuborderQuantums(
	twapOrderPlacement types.TwapOrderPlacement,
	clobPair types.ClobPair,
) uint64 {
	originalQuantumsPerLeg := float64(twapOrderPlacement.Order.Quantums) / float64(twapOrderPlacement.Order.GetTotalLegsTWAPOrder())

	// Calculate the quantums for the suborder capping at 3x the original quantums per leg
	remainingPerLeg := float64(twapOrderPlacement.RemainingQuantums) / float64(twapOrderPlacement.RemainingLegs)
	suborderQuantums := lib.Min(remainingPerLeg, TWAP_MAX_SUBORDER_CATCHUP_MULTIPLE*originalQuantumsPerLeg)

	// Round down to nearest multiple of StepBaseQuantums
	quantumsByStepBaseQuantums := uint64(math.Floor(suborderQuantums / float64(clobPair.StepBaseQuantums)))
	suborderQuantumsRounded := quantumsByStepBaseQuantums * clobPair.StepBaseQuantums

	return suborderQuantumsRounded
}

func (k Keeper) GenerateSuborder(
	ctx sdk.Context,
	suborderId types.OrderId,
	twapOrderPlacement types.TwapOrderPlacement,
	blockTime int64,
) types.Order {
	parentOrder := twapOrderPlacement.Order
	order := types.Order{
		OrderId:    suborderId,
		Side:       twapOrderPlacement.Order.Side,
		ReduceOnly: twapOrderPlacement.Order.ReduceOnly,
	}

	priceTolerancePpm := int32(parentOrder.TwapParameters.PriceTolerance)
	if parentOrder.Side == types.Order_SIDE_SELL {
		// for sell orders, we want to adjust the price down
		priceTolerancePpm = -priceTolerancePpm
	}

	// calculate the suborder price with slippage adjustment
	clobPair := k.mustGetClobPair(ctx, parentOrder.GetClobPairId())
	order.Subticks = k.GetOraclePriceAdjustedByPercentageSubticks(ctx, clobPair, priceTolerancePpm)

	// calculate the suborder quantums based on remaining quantums and legs
	order.Quantums = k.calculateSuborderQuantums(twapOrderPlacement, clobPair)

	// set good til block time from current block time
	order.GoodTilOneof = &types.Order_GoodTilBlockTime{
		GoodTilBlockTime: uint32(blockTime + TWAP_SUBORDER_GOOD_TIL_BLOCK_TIME_OFFSET),
	}

	return order
}
