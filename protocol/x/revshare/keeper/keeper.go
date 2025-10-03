package keeper

import (
	"fmt"

	"cosmossdk.io/log"
	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4-chain/protocol/lib"
	affiliateskeeper "github.com/dydxprotocol/v4-chain/protocol/x/affiliates/keeper"
	feetierskeeper "github.com/dydxprotocol/v4-chain/protocol/x/feetiers/keeper"
	"github.com/dydxprotocol/v4-chain/protocol/x/revshare/types"
	statsKeeper "github.com/dydxprotocol/v4-chain/protocol/x/stats/keeper"
)

type (
	Keeper struct {
		cdc              codec.BinaryCodec
		storeKey         storetypes.StoreKey
		authorities      map[string]struct{}
		affiliatesKeeper affiliateskeeper.Keeper
		feetiersKeeper   feetierskeeper.Keeper
		statsKeeper      statsKeeper.Keeper
	}
)

func NewKeeper(
	cdc codec.BinaryCodec,
	storeKey storetypes.StoreKey,
	authorities []string,
	affiliatesKeeper affiliateskeeper.Keeper,
	feetiersKeeper feetierskeeper.Keeper,
	statsKeeper statsKeeper.Keeper,
) *Keeper {
	return &Keeper{
		cdc:              cdc,
		storeKey:         storeKey,
		authorities:      lib.UniqueSliceToSet(authorities),
		affiliatesKeeper: affiliatesKeeper,
		feetiersKeeper:   feetiersKeeper,
		statsKeeper:      statsKeeper,
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
