package keeper

import (
	"fmt"

	"cosmossdk.io/log"
	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4-chain/protocol/lib"
	"github.com/dydxprotocol/v4-chain/protocol/x/govplus/types"
)

type (
	Keeper struct {
		cdc           codec.BinaryCodec
		stakingKeeper types.StakingKeeper
		storeKey      storetypes.StoreKey
		authorities   map[string]struct{}
	}
)

func NewKeeper(
	cdc codec.BinaryCodec,
	stakingKeeper types.StakingKeeper,
	storeKey storetypes.StoreKey,
	authorities []string,
) *Keeper {
	return &Keeper{
		cdc:           cdc,
		stakingKeeper: stakingKeeper,
		storeKey:      storeKey,
		authorities:   lib.UniqueSliceToSet(authorities),
	}
}

func (k Keeper) HasAuthority(authority string) bool {
	_, ok := k.authorities[authority]
	return ok
}

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With(log.ModuleKey, fmt.Sprintf("x/%s", types.ModuleName))
}

func (k Keeper) InitializeForGenesis(ctx sdk.Context) {
}
