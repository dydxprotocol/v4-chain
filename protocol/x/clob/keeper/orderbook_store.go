package keeper

import (
	"cosmossdk.io/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
)

func (k Keeper) GetOrderbookFromStore(
	ctx sdk.Context,
	clobPairId types.ClobPairId,
) (ret *types.Orderbook, exists bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte(types.OrderbookStoreKeyPrefix))

	val := types.OrderbookStore{}
	b := store.Get(clobPairKey(clobPairId))
	if b == nil {
		return ret, false
	}

	k.cdc.MustUnmarshal(b, &val)
	return types.NewOrderbookFromStore(val), true
}
