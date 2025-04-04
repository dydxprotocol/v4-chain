package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
)

func (k Keeper) SetTWAPOrderPlacement(ctx sdk.Context,
	order types.Order,
	blockHeight uint32,
) {
	store := k.GetTWAPOrderPlacementStore(ctx)
	orderKey := order.OrderId.ToStateKey()

	total_legs := order.TwapParameters.Duration / order.TwapParameters.Interval

	twapOrderPlacement := types.TwapOrderPlacement{
		Order:             order,
		TotalLegs:         total_legs,
		RemainingLegs:     total_legs,
		RemainingQuantums: order.Quantums,
		BlockHeight:       blockHeight,
	}

	k.addSuborderToTriggerStore(ctx, twapOrderPlacement, 0, 0)

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
) (val types.TwapTriggerPlacement, found bool) {
	store := k.GetTWAPTriggerOrderPlacementStore(ctx)

	iterator := store.Iterator(nil, nil)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var triggerPlacement types.TwapTriggerPlacement
		k.cdc.MustUnmarshal(iterator.Value(), &triggerPlacement)

		if triggerPlacement.Order.OrderId.SubaccountId.Owner == orderId.SubaccountId.Owner &&
			triggerPlacement.Order.OrderId.SubaccountId.Number == orderId.SubaccountId.Number &&
			triggerPlacement.Order.OrderId.ClientId == orderId.ClientId &&
			triggerPlacement.Order.OrderId.ClobPairId == orderId.ClobPairId &&
			triggerPlacement.Order.OrderId.SequenceNumber == orderId.SequenceNumber {
			return triggerPlacement, true
		}
	}
	return types.TwapTriggerPlacement{}, false
}

// addSuborderToTriggerStore creates a TWAP suborder from a parent TWAP order and adds it to the trigger store.
// The suborder's size is calculated by dividing the parent order's quantums by the total number of legs.
// The suborder is marked with the TWAP suborder flag and given the specified sequence number.
func (k Keeper) addSuborderToTriggerStore(
	ctx sdk.Context,
	twapOrderPlacement types.TwapOrderPlacement,
	sequenceNumber uint32,
	triggerOffset int64,
) {
	triggerStore := k.GetTWAPTriggerOrderPlacementStore(ctx)

	if sequenceNumber > twapOrderPlacement.TotalLegs {
		// remove the parent twap order from the store
		store := k.GetTWAPOrderPlacementStore(ctx)
		orderKey := twapOrderPlacement.Order.OrderId.ToStateKey()
		store.Delete(orderKey)
		// TODO: (anmol) emit event?
		return
	}

	triggerTime := ctx.BlockTime().Unix() + triggerOffset
	suborder := twapOrderPlacement.Order

	// suborder quantums and subticks are set to 0 until triggered
	// and updated in the end blocker based off oracle price and
	// remaining quantums and legs
	suborder.Quantums = 0
	suborder.Subticks = 0

	// Set the order flag to indicate this is a TWAP suborder
	suborder.OrderId.OrderFlags = types.OrderIdFlags_TwapSuborder
	suborder.OrderId.SequenceNumber = sequenceNumber

	suborder.GoodTilOneof = &types.Order_GoodTilBlockTime{
		GoodTilBlockTime: uint32(triggerTime + 3),
	}

	// Create trigger placement
	triggerPlacement := types.TwapTriggerPlacement{
		Order:            suborder,
		TriggerBlockTime: uint64(triggerTime),
	}

	triggerPlacementBytes := k.cdc.MustMarshal(&triggerPlacement)
	triggerKey := types.GetTWAPTriggerKey(triggerTime, suborder.OrderId)
	triggerStore.Set(triggerKey, triggerPlacementBytes)
}
