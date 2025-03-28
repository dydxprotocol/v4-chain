package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
)

func (k Keeper) SetTWAPOrderPlacement(ctx sdk.Context, 
	order types.Order,
	blockHeight uint32,
) (err error) {
	store := k.GetTWAPOrderPlacementStore(ctx)
	orderKey := order.OrderId.ToStateKey()

	total_legs := order.TwapConfig.Duration / 30 // TODO: (anmol) add checks and also consider configurable time increments
	err = k.partitionAndSetTWAPTriggerOrders(ctx, order, total_legs)
	if err != nil {
		return err
	}

	// TODO: (anmol) this is assuming we fire off an order immediately
	// also need to consider the case where that initial order is not filled. probably handle initial order outside of this function
	// suborder_size := order.Quantums / total_legs
	// remaining_quantums := order.Quantums - suborder_size

	// TODO: (anmol) potentially add a suborder array here which gets modified as we fire off suborders - maybe status as well?
	twapOrderPlacement := types.TWAPOrderPlacement{
		Order: order,
		RemainingLegs: total_legs,
		RemainingQuantums: order.Quantums,
		BlockHeight: blockHeight,
	}

	twapOrderPlacementBytes := k.cdc.MustMarshal(&twapOrderPlacement)

	store.Set(orderKey, twapOrderPlacementBytes)
	return nil
}

// partitionTwapOrder splits a TWAP order into equal-sized suborders and stores them in the trigger store.
// Each suborder will be triggered at its designated block time.
func (k Keeper) partitionAndSetTWAPTriggerOrders(
	ctx sdk.Context,
	twapOrder types.Order,
	totalLegs uint32,
) error {
	triggerStore := k.GetTWAPTriggerOrderPlacementStore(ctx)

	// Create and store suborders in the trigger store
	for i := uint32(0); i < totalLegs; i++ {
		triggerTime := ctx.BlockTime().Unix() + (int64(i) * int64(twapOrder.TwapConfig.Interval))
		// Create a suborder with quantums split evenly across legs
		suborder := twapOrder
		suborder.Quantums = twapOrder.Quantums / uint64(totalLegs) // TODO: (anmol) what if not evenly divisible? front/backload load the remainder?
		suborder.TimeInForce = types.Order_TIME_IN_FORCE_IOC
		suborder.GoodTilOneof = &types.Order_GoodTilBlockTime{GoodTilBlockTime: uint32(triggerTime)} // TODO: (anmol) add some buffer? how does IOC work with this?
		// Set the order flag to indicate this is a TWAP suborder
		suborder.OrderId.OrderFlags = types.OrderIdFlags_TwapSuborder
		suborder.OrderId.SequenceNumber = uint64(i)
		// Create trigger placement
		triggerPlacement := types.TWAPTriggerPlacement{
			Order: suborder,
			TriggerBlockHeight: uint64(triggerTime),
		}

		// Marshal and store the trigger placement
		triggerPlacementBytes := k.cdc.MustMarshal(&triggerPlacement)
		triggerKey := types.GetTWAPTriggerKey(triggerTime, suborder.OrderId)
		triggerStore.Set(triggerKey, triggerPlacementBytes)
	}

	return nil
}
