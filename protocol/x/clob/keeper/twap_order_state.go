package keeper

import (
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

	total_legs := order.TwapConfig.Duration / order.TwapConfig.Interval

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

	// suborder quantums set to 0 until triggered and updated based on remaining quantums and legs
	suborder.Quantums = 0

	// Set the order flag to indicate this is a TWAP suborder
	suborder.OrderId.OrderFlags = types.OrderIdFlags_TwapSuborder
	suborder.OrderId.SequenceNumber = sequenceNumber

	goodTilBlockTimeOffset := twapOrderPlacement.Order.TwapConfig.GoodTillBlockTimeOffset
	suborder.GoodTilOneof = &types.Order_GoodTilBlockTime{
		GoodTilBlockTime: uint32(triggerTime + int64(goodTilBlockTimeOffset)),
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
		slippage_adjustment := order.TwapConfig.SlippagePercent
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
			int64(op.twapOrderPlacement.Order.TwapConfig.Interval),
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
