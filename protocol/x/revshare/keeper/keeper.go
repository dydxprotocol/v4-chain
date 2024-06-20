package keeper

import (
	"fmt"

	"cosmossdk.io/store/prefix"

	"cosmossdk.io/log"
	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4-chain/protocol/lib"
	"github.com/dydxprotocol/v4-chain/protocol/x/revshare/types"
)

type (
	Keeper struct {
		cdc         codec.BinaryCodec
		storeKey    storetypes.StoreKey
		authorities map[string]struct{}
	}
)

func NewKeeper(
	cdc codec.BinaryCodec,
	storeKey storetypes.StoreKey,
	authorities []string,
) *Keeper {
	return &Keeper{
		cdc:         cdc,
		storeKey:    storeKey,
		authorities: lib.UniqueSliceToSet(authorities),
	}
}

func (k Keeper) HasAuthority(authority string) bool {
	_, ok := k.authorities[authority]
	return ok
}

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With(log.ModuleKey, fmt.Sprintf("x/%s", types.ModuleName))
}

func (k Keeper) InitializeForGenesis(ctx sdk.Context) {}

// Function to serialize market mapper rev share params and store in the
// module store
func (k Keeper) SetMarketMapperRevShareDetails(
	ctx sdk.Context,
	marketId uint32,
	params types.MarketMapperRevShareDetails,
) (err error) {
	// Store the rev share details for provided market in module store
	store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte(types.MarketMapperRevSharePrefix))
	b := k.cdc.MustMarshal(&params)
	store.Set(lib.Uint32ToKey(marketId), b)
	return nil
}

// Function to retrieve marketmapper revshare params from module store
func (k Keeper) GetMarketMapperRevShareDetails(
	ctx sdk.Context,
	marketId uint32,
) (params types.MarketMapperRevShareDetails, err error) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte(types.MarketMapperRevSharePrefix))
	b := store.Get(lib.Uint32ToKey(marketId))
	if b == nil {
		return params, fmt.Errorf("MarketMapperRevShareDetails not found for marketId: %d", marketId)
	}
	k.cdc.MustUnmarshal(b, &params)
	return params, nil
}
