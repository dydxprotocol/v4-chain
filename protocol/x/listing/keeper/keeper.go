package keeper

import (
	"fmt"

	"cosmossdk.io/log"
	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4-chain/protocol/lib"
	"github.com/dydxprotocol/v4-chain/protocol/x/listing/types"
)

type (
	Keeper struct {
		cdc              codec.BinaryCodec
		storeKey         storetypes.StoreKey
		authorities      map[string]struct{}
		PricesKeeper     types.PricesKeeper
		ClobKeeper       types.ClobKeeper
		MarketMapKeeper  types.MarketMapKeeper
		PerpetualsKeeper types.PerpetualsKeeper
	}
)

func NewKeeper(
	cdc codec.BinaryCodec,
	storeKey storetypes.StoreKey,
	authorities []string,
	pricesKeeper types.PricesKeeper,
	clobKeeper types.ClobKeeper,
	marketMapKeeper types.MarketMapKeeper,
	perpetualsKeeper types.PerpetualsKeeper,
) *Keeper {
	return &Keeper{
		cdc:              cdc,
		storeKey:         storeKey,
		authorities:      lib.UniqueSliceToSet(authorities),
		PricesKeeper:     pricesKeeper,
		ClobKeeper:       clobKeeper,
		MarketMapKeeper:  marketMapKeeper,
		PerpetualsKeeper: perpetualsKeeper,
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
