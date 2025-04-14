package keeper

import (
	"encoding/binary"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
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
