package keeper

import (
	"encoding/binary"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4-chain/protocol/lib"
	"github.com/dydxprotocol/v4-chain/protocol/x/clob/types"

	"time"
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
) (o types.OrderId, t uint64, found bool) {
	store := k.GetTWAPTriggerOrderPlacementStore(ctx)

	iterator := store.Iterator(nil, nil)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var suborderId types.OrderId
		k.cdc.MustUnmarshal(iterator.Key()[8:], &suborderId)

		timestamp := binary.BigEndian.Uint64(iterator.Key()[0:8])
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

func (k Keeper) GenerateAndPlaceTriggeredTwapSuborders(ctx sdk.Context, block_time time.Time) {
	triggerStore := k.GetTWAPTriggerOrderPlacementStore(ctx)

	iterator := triggerStore.Iterator(nil, nil)

	var operationsToProcess []struct {
		keyToDelete        []byte
		orderToPlace       types.Order
		twapOrderPlacement types.TwapOrderPlacement
	}

	for ; iterator.Valid(); iterator.Next() {
		var triggerPlacement types.TwapTriggerPlacement
		k.cdc.MustUnmarshal(iterator.Value(), &triggerPlacement)

		if triggerPlacement.TriggerBlockTime > uint64(block_time.Unix()) {
			break // all remaining suborders are in the future
		}

		order := triggerPlacement.Order

		parentOrderId := types.OrderId{
			SubaccountId:   order.OrderId.SubaccountId,
			ClientId:       order.OrderId.ClientId,
			OrderFlags:     types.OrderIdFlags_Twap, // Set directly to TWAP
			ClobPairId:     order.OrderId.ClobPairId,
			SequenceNumber: 0, // Set directly to 0
		}

		twapOrderPlacement, found := k.GetTwapOrderPlacement(ctx, parentOrderId)
		if !found {
			panic("parent twap order not found") // TODO: (anmol) handle order cancellation
		}

		clobPair := k.mustGetClobPair(ctx, order.GetClobPairId())
		slippage_adjustment := order.TwapParameters.SlippagePercent
		if order.Side == types.Order_SIDE_SELL {
			slippage_adjustment = -slippage_adjustment
		}
		// calculate the suborder price with slippage adjustment
		order.Subticks = k.GetOraclePriceAdjustedByPercentageSubticks(ctx, clobPair, float64(slippage_adjustment)/10000.0)

		// calculate the suborder quantums based on remaining quantums and legs
		order.Quantums = k.calculateSuborderQuantums(twapOrderPlacement)

		operationsToProcess = append(operationsToProcess, struct {
			keyToDelete        []byte
			orderToPlace       types.Order
			twapOrderPlacement types.TwapOrderPlacement
		}{
			keyToDelete:        append([]byte{}, iterator.Key()...),
			orderToPlace:       order,
			twapOrderPlacement: twapOrderPlacement,
		})
	}
	iterator.Close()

	for _, op := range operationsToProcess {
		// Delete from trigger store
		triggerStore.Delete(op.keyToDelete)

		// decrement remaining legs
		k.DecrementTwapOrderRemainingLegs(ctx, &op.twapOrderPlacement)

		// add the next suborder to the trigger store
		k.addSuborderToTriggerStore(
			ctx,
			op.twapOrderPlacement,
			op.orderToPlace.OrderId.SequenceNumber+1,
			int64(op.twapOrderPlacement.Order.TwapParameters.Interval),
		)

		// place triggered suborder
		err := k.HandleMsgPlaceOrder(ctx, &types.MsgPlaceOrder{Order: op.orderToPlace}, true)
		if err != nil {
			continue // TODO: (anmol) handle suborder placement failure
		}
	}
}

func (k Keeper) calculateSuborderQuantums(
	twapOrderPlacement types.TwapOrderPlacement,
) uint64 {
	// TODO: (anmol) ensure rounding is correct. same as subticks (check factor) and min
	// total order size > min_size * legs
	originalQuantumsPerLeg := twapOrderPlacement.Order.Quantums / uint64(twapOrderPlacement.TotalLegs)

	// Calculate the quantums for the suborder capping at 3x the original quantums per leg
	remainingPerLeg := twapOrderPlacement.RemainingQuantums / uint64(twapOrderPlacement.RemainingLegs)
	return lib.Min(remainingPerLeg, 3*originalQuantumsPerLeg)
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
