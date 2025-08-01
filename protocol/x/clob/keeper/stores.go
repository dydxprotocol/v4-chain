package keeper

import (
	"fmt"

	storetypes "cosmossdk.io/store/types"

	"github.com/dydxprotocol/v4-chain/protocol/lib"

	"cosmossdk.io/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
)

// GetLongTermOrderPlacementStore fetches a state store used for creating,
// reading, updating, and deleting a stateful order placement from state.
func (k Keeper) GetLongTermOrderPlacementStore(ctx sdk.Context) prefix.Store {
	return prefix.NewStore(
		ctx.KVStore(k.storeKey),
		[]byte(types.LongTermOrderPlacementKeyPrefix),
	)
}

// GetUntriggeredConditionalOrderPlacementStore fetches a state store used for creating,
// reading, updating, and deleting untriggered conditional order placement from state.
func (k Keeper) GetUntriggeredConditionalOrderPlacementStore(ctx sdk.Context) prefix.Store {
	return prefix.NewStore(
		ctx.KVStore(k.storeKey),
		[]byte(types.UntriggeredConditionalOrderKeyPrefix),
	)
}

// GetUncommittedStatefulOrderPlacementTransientStore fetches a state store used for creating,
// reading, updating, and deleting a stateful order placement from transient state.
func (k Keeper) GetUncommittedStatefulOrderPlacementTransientStore(ctx sdk.Context) prefix.Store {
	return prefix.NewStore(
		ctx.KVStore(k.transientStoreKey),
		[]byte(types.UncommittedStatefulOrderPlacementKeyPrefix),
	)
}

// GetUncommittedStatefulOrderCancellationTransientStore fetches a state store used for creating,
// reading, updating, and deleting a stateful order cancellation from transient state.
func (k Keeper) GetUncommittedStatefulOrderCancellationTransientStore(ctx sdk.Context) prefix.Store {
	return prefix.NewStore(
		ctx.KVStore(k.transientStoreKey),
		[]byte(types.UncommittedStatefulOrderCancellationKeyPrefix),
	)
}

// GetUncommittedStatefulOrderCountTransientStore fetches a state store used for creating,
// reading, updating, and deleting a stateful order count from transient state. This represents
// the number of uncommitted `order placements - order cancellations` during `CheckTx`.
func (k Keeper) GetUncommittedStatefulOrderCountTransientStore(ctx sdk.Context) prefix.Store {
	lib.AssertCheckTxMode(ctx)
	return prefix.NewStore(
		ctx.KVStore(k.transientStoreKey),
		[]byte(types.UncommittedStatefulOrderCountPrefix),
	)
}

// GetStatefulOrderCountMemStore fetches a state store used for creating,
// reading, updating, and deleting a stateful order count from stores. This represents
// the number of long term and triggered conditional orders.
func (k Keeper) GetStatefulOrderCountMemStore(ctx sdk.Context) prefix.Store {
	return prefix.NewStore(
		ctx.KVStore(k.memKey),
		[]byte(types.StatefulOrderCountPrefix),
	)
}

// GetTriggeredConditionalOrderPlacementStore fetches a state store used for creating,
// reading, updating, and deleting a stateful order placement from state.
func (k Keeper) GetTriggeredConditionalOrderPlacementStore(ctx sdk.Context) prefix.Store {
	return prefix.NewStore(
		ctx.KVStore(k.storeKey),
		[]byte(types.TriggeredConditionalOrderKeyPrefix),
	)
}

func (k Keeper) GetTWAPOrderPlacementStore(ctx sdk.Context) prefix.Store {
	return prefix.NewStore(
		ctx.KVStore(k.storeKey),
		[]byte(types.TWAPOrderKeyPrefix),
	)
}

func (k Keeper) GetTWAPTriggerOrderPlacementStore(ctx sdk.Context) prefix.Store {
	return prefix.NewStore(
		ctx.KVStore(k.storeKey),
		[]byte(types.TWAPTriggerOrderKeyPrefix),
	)
}

// getTransientStore fetches a transient store used for reading and
// updating the transient store.
func (k Keeper) getTransientStore(ctx sdk.Context) storetypes.KVStore {
	return ctx.KVStore(k.transientStoreKey)
}

// fetchStateStoresForOrder fetches the state store and memstore for a given orderId. If the orderId is
// for a long term order, the long term order placement store will be returned. If it is conditional, the
// IsConditionalOrderTriggered function will be used to determine which conditional order placement
// state store is returned.
// Currently, this function supports conditional orders and long term orders.
// If the given order id is conditional, it will return the Untriggered conditional order state store.
func (k Keeper) fetchStateStoresForOrder(
	ctx sdk.Context,
	orderId types.OrderId,
) prefix.Store {
	orderId.MustBeStatefulOrder()

	if orderId.IsConditionalOrder() {
		triggered := k.IsConditionalOrderTriggered(ctx, orderId)
		if triggered {
			return k.GetTriggeredConditionalOrderPlacementStore(ctx)
		}
		return k.GetUntriggeredConditionalOrderPlacementStore(ctx)
	} else if orderId.IsLongTermOrder() || orderId.IsTwapSuborder() {
		return k.GetLongTermOrderPlacementStore(ctx)
	}
	panic(
		fmt.Sprintf(
			"FetchStateStoresForOrder: orderId (%+v) not supported",
			orderId,
		),
	)
}

// GetMinTradePriceStore fetches a transient store used for reading and updating the
// min trade prices for perpetuals.
func (k Keeper) GetMinTradePriceStore(ctx sdk.Context) prefix.Store {
	return prefix.NewStore(
		ctx.KVStore(k.transientStoreKey),
		[]byte(types.MinTradePricePrefix),
	)
}

// GetMaxTradePriceStore fetches a transient store used for reading and updating the
// max trade prices for perpetuals.
func (k Keeper) GetMaxTradePriceStore(ctx sdk.Context) prefix.Store {
	return prefix.NewStore(
		ctx.KVStore(k.transientStoreKey),
		[]byte(types.MaxTradePricePrefix),
	)
}
