package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4/x/clob/types"
)

// GetUncommittedStatefulOrderPlacement gets a stateful order and the placement information from uncommitted state.
// OrderId can be conditional or long term.
// Returns false if no stateful order exists in uncommitted state with `orderId`.
func (k Keeper) GetUncommittedStatefulOrderPlacement(
	ctx sdk.Context,
	orderId types.OrderId,
) (val types.LongTermOrderPlacement, found bool) {
	// If this is a Short-Term order, panic.
	orderId.MustBeStatefulOrder()

	store := k.GetUncommittedStatefulOrderPlacementTransientStore(ctx)

	b := store.Get(types.OrderIdKey(orderId))
	if b == nil {
		return val, false
	}

	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

// GetUncommittedStatefulOrderCancellation gets a stateful order cancellation from uncommitted state.
// OrderId can be conditional or long term.
// Returns false if no stateful order cancellation exists in uncommitted state with `orderId`.
func (k Keeper) GetUncommittedStatefulOrderCancellation(
	ctx sdk.Context,
	orderId types.OrderId,
) (val types.MsgCancelOrder, found bool) {
	// If this is a Short-Term order, panic.
	orderId.MustBeStatefulOrder()

	store := k.GetUncommittedStatefulOrderCancellationTransientStore(ctx)

	b := store.Get(types.OrderIdKey(orderId))
	if b == nil {
		return val, false
	}

	k.cdc.MustUnmarshal(b, &val)
	return val, true
}
