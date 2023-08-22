package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4-chain/protocol/lib"
	"github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
	satypes "github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"
)

// To be committed stateful orders are ones that this validator is aware of during block processing (e.g. `DeliverTx`).
// See `uncommitted_stateful_order_state.go` for associated functions related to stateful orders
// that this validator is aware of that have yet to be included in a block (e.g. `CheckTx`).

// GetToBeCommittedStatefulOrderCount gets a count of how many stateful orders will be added for the associated
// subaccount during `DeliverTx`. This is represented by the number of stateful order `placements - removals`.
// Note that this value can be negative (for example if the stateful order is already on the book and the cancellation
// is to be committed).
// OrderId can be conditional or long term.
func (k Keeper) GetToBeCommittedStatefulOrderCount(
	ctx sdk.Context,
	orderId types.OrderId,
) int32 {
	// If this is a Short-Term order, panic.
	orderId.MustBeStatefulOrder()

	store := k.GetToBeCommittedStatefulOrderCountTransientStore(ctx)

	b := store.Get(satypes.SubaccountKey(orderId.SubaccountId))
	if b == nil {
		return 0
	}

	return lib.BytesToInt32(b)
}

// SetToBeCommittedStatefulOrderCount sets a count of how many stateful orders will be added for the associated
// subaccount during `DeliverTx`. This represents the number of stateful order `placements - cancellations`.
// Note that this value can be negative (for example if the stateful order is already on the book and the cancellation
// is to be committed).
// OrderId can be conditional or long term.
func (k Keeper) SetToBeCommittedStatefulOrderCount(
	ctx sdk.Context,
	orderId types.OrderId,
	count int32,
) {
	// If this is a Short-Term order, panic.
	orderId.MustBeStatefulOrder()

	store := k.GetToBeCommittedStatefulOrderCountTransientStore(ctx)
	store.Set(
		satypes.SubaccountKey(orderId.SubaccountId),
		lib.Int32ToBytes(count),
	)
}
