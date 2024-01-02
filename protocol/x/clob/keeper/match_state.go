package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	gogotypes "github.com/cosmos/gogoproto/types"
	"github.com/dydxprotocol/v4-chain/protocol/lib"
	"github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
)

// GetLastTradePriceForPerpetual gets the last trade price for a perpetual.
func (k Keeper) GetLastTradePriceForPerpetual(
	ctx sdk.Context,
	perpetualId uint32,
) (subticks types.Subticks, found bool) {
	store := k.GetLastTradePriceStore(ctx)

	b := store.Get(lib.Uint32ToKey(perpetualId))
	if b == nil {
		return 0, false
	}

	result := gogotypes.UInt64Value{Value: 0}
	k.cdc.MustUnmarshal(b, &result)
	return types.Subticks(result.Value), true
}

// SetLastTradePriceForPerpetual sets the last trade price for a perpetual.
func (k Keeper) SetLastTradePriceForPerpetual(
	ctx sdk.Context,
	perpetualId uint32,
	subticks types.Subticks,
) {
	store := k.GetLastTradePriceStore(ctx)

	value := gogotypes.UInt64Value{Value: subticks.ToUint64()}
	store.Set(
		lib.Uint32ToKey(perpetualId),
		k.cdc.MustMarshal(&value),
	)
}
