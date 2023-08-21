package keeper

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4-chain/protocol/lib"
	"github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
	satypes "github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"
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

// GetUncommittedStatefulOrderCount gets a count of uncommitted stateful orders for the associated subaccount.
// This is represented by the number of stateful order `placements - cancellations`. Note that this value
// can be negative (for example if the stateful order is already on the book and the cancellation is uncommitted).
// OrderId can be conditional or long term.
func (k Keeper) GetUncommittedStatefulOrderCount(
	ctx sdk.Context,
	orderId types.OrderId,
) int32 {
	// If this is a Short-Term order, panic.
	orderId.MustBeStatefulOrder()

	store := k.GetUncommittedStatefulOrderCountTransientStore(ctx)

	b := store.Get(satypes.SubaccountKey(orderId.SubaccountId))
	if b == nil {
		return 0
	}

	return lib.BytesToInt32(b)
}

// SetUncommittedStatefulOrderCount sets a count of uncommitted stateful orders for the associated subaccount.
// This represents the number of stateful order `placements - cancellations`. Note that this value
// can be negative (for example if the stateful order is already on the book and the cancellation is uncommitted).
// OrderId can be conditional or long term.
func (k Keeper) SetUncommittedStatefulOrderCount(
	ctx sdk.Context,
	orderId types.OrderId,
	count int32,
) {
	// If this is a Short-Term order, panic.
	orderId.MustBeStatefulOrder()

	store := k.GetUncommittedStatefulOrderCountTransientStore(ctx)
	store.Set(
		satypes.SubaccountKey(orderId.SubaccountId),
		lib.Int32ToBytes(count),
	)
}

// MustAddUncommittedStatefulOrderPlacement adds a new order placemenet by `OrderId` to a transient store and
// increments the per subaccount uncommitted stateful order count.
//
// This method will panic if the order already exists or if the order count overflows a uint32.
func (k Keeper) MustAddUncommittedStatefulOrderPlacement(ctx sdk.Context, msg *types.MsgPlaceOrder) {
	lib.AssertCheckTxMode(ctx)

	orderId := msg.Order.OrderId
	// If this is a Short-Term order, panic.
	orderId.MustBeStatefulOrder()

	if _, exists := k.GetUncommittedStatefulOrderPlacement(ctx, orderId); exists {
		panic(fmt.Sprintf("MustAddUncommittedStatefulOrderPlacement: order %v already exists", orderId))
	}

	longTermOrderPlacement := types.LongTermOrderPlacement{
		Order: msg.Order,
	}

	orderIdBytes := types.OrderIdKey(orderId)
	b := k.cdc.MustMarshal(&longTermOrderPlacement)

	store := k.GetUncommittedStatefulOrderPlacementTransientStore(ctx)
	store.Set(orderIdBytes, b)

	k.SetUncommittedStatefulOrderCount(
		ctx,
		orderId,
		k.GetUncommittedStatefulOrderCount(ctx, orderId)+1,
	)
}

// MustAddUncommittedStatefulOrderCancellation adds a new order cancellation by `OrderId` to a transient store and
// decrements the per subaccount uncommitted stateful order count.
//
// This method will panic if the order cancellation already exists or if the order count underflows a uint32.
func (k Keeper) MustAddUncommittedStatefulOrderCancellation(ctx sdk.Context, msg *types.MsgCancelOrder) {
	lib.AssertCheckTxMode(ctx)

	orderId := msg.OrderId
	// If this is a Short-Term order, panic.
	orderId.MustBeStatefulOrder()

	if _, exists := k.GetUncommittedStatefulOrderCancellation(ctx, orderId); exists {
		panic(fmt.Sprintf("MustAddUncommittedStatefulOrderPlacement: order cancellation %v already exists", orderId))
	}

	orderIdBytes := types.OrderIdKey(orderId)
	b := k.cdc.MustMarshal(msg)

	store := k.GetUncommittedStatefulOrderCancellationTransientStore(ctx)
	store.Set(orderIdBytes, b)

	k.SetUncommittedStatefulOrderCount(
		ctx,
		orderId,
		k.GetUncommittedStatefulOrderCount(ctx, orderId)-1,
	)
}
