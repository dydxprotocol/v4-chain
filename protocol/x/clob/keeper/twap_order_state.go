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

	total_legs := order.TwapConfig.Duration / order.TwapConfig.Interval
	k.addInitialSuborderToTriggerStore(ctx, order, total_legs)

	// TODO: (anmol) this is assuming we fire off an order immediately
	// also need to consider the case where that initial order is not filled.
	// probably handle initial order outside of this function
	// suborder_size := order.Quantums / total_legs
	// remaining_quantums := order.Quantums - suborder_size

	// TODO: (anmol) potentially add a suborder array here which gets
	// modified as we fire off suborders - maybe status as well?
	twapOrderPlacement := types.TwapOrderPlacement{
		Order:             order,
		RemainingLegs:     total_legs,
		RemainingQuantums: order.Quantums,
		BlockHeight:       blockHeight,
	}

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

// GetTwapTriggerPlacements gets all TWAP trigger placements for a given orderId.
// Returns empty slice if no trigger placements exist.
func (k Keeper) GetTwapTriggerPlacements(
	ctx sdk.Context,
	orderId types.OrderId,
) (val []*types.TwapTriggerPlacement, found bool) {
	store := k.GetTWAPTriggerOrderPlacementStore(ctx)
	var triggerPlacements []*types.TwapTriggerPlacement

	iterator := store.Iterator(nil, nil)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var triggerPlacement types.TwapTriggerPlacement
		k.cdc.MustUnmarshal(iterator.Value(), &triggerPlacement)

		if triggerPlacement.Order.OrderId.SubaccountId.Owner == orderId.SubaccountId.Owner &&
			triggerPlacement.Order.OrderId.SubaccountId.Number == orderId.SubaccountId.Number &&
			triggerPlacement.Order.OrderId.ClientId == orderId.ClientId &&
			triggerPlacement.Order.OrderId.ClobPairId == orderId.ClobPairId {
			triggerPlacements = append(triggerPlacements, &triggerPlacement)
		}
	}
	return triggerPlacements, len(triggerPlacements) > 0
}

// partitionTwapOrder splits a TWAP order into equal-sized suborders and stores them in the
// trigger store. Each suborder will be triggered at its designated block time.
func (k Keeper) addInitialSuborderToTriggerStore(
	ctx sdk.Context,
	twapOrder types.Order,
	totalLegs uint32,
) {
	triggerStore := k.GetTWAPTriggerOrderPlacementStore(ctx)

	// Create and store single suborder in the trigger store
	triggerTime := ctx.BlockTime().Unix()
	// Create a suborder with correct quantums
	suborder := twapOrder
	// TODO: (anmol) what if not evenly divisible? front/backload load the remainder?
	suborder.Quantums = twapOrder.Quantums / uint64(totalLegs)
	

	// suborder.TimeInForce = types.Order_TIME_IN_FORCE_IOC // how long should it stay resting?
	// suborder.GoodTilOneof = &types.Order_GoodTilBlockTime{GoodTilBlockTime: uint32(triggerTime)} 
	// TODO: (anmol) add some buffer? how does IOC work with this?

	// Set the order flag to indicate this is a TWAP suborder
	suborder.OrderId.OrderFlags = types.OrderIdFlags_TwapSuborder
	suborder.OrderId.SequenceNumber = 0
	// Create trigger placement
	triggerPlacement := types.TwapTriggerPlacement{
		Order:            suborder,
		TriggerBlockTime: uint64(triggerTime),
	}

	// Marshal and store the trigger placement
	triggerPlacementBytes := k.cdc.MustMarshal(&triggerPlacement)
	triggerKey := types.GetTWAPTriggerKey(triggerTime, suborder.OrderId)
	triggerStore.Set(triggerKey, triggerPlacementBytes)
}
