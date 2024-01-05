package keeper

import (
	"fmt"

	"cosmossdk.io/log"
	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4-chain/protocol/lib"
	"github.com/dydxprotocol/v4-chain/protocol/x/delaymsg/types"
)

type (
	Keeper struct {
		cdc      codec.Codec
		storeKey storetypes.StoreKey
		// authorities stores addresses capable of submitting a delayed message.
		authorities map[string]struct{}
		router      *baseapp.MsgServiceRouter
	}
)

// NewKeeper creates a new x/delaymsg keeper.
func NewKeeper(
	cdc codec.Codec,
	storeKey storetypes.StoreKey,
	router *baseapp.MsgServiceRouter,
	authorities []string,
) *Keeper {
	return &Keeper{
		cdc:         cdc,
		storeKey:    storeKey,
		authorities: lib.UniqueSliceToSet(authorities),
		router:      router,
	}
}

func (k Keeper) HasAuthority(authority string) bool {
	_, ok := k.authorities[authority]
	return ok
}

// Router returns the x/delaymsg router.
func (k Keeper) Router() lib.MsgRouter {
	return k.router
}

// InitializeForGenesis initializes the x/delaymsg keeper for genesis.
func (k Keeper) InitializeForGenesis(ctx sdk.Context) {
	k.SetNextDelayedMessageId(ctx, 0)
}

// Logger returns a module-specific logger for x/delaymsg.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With(log.ModuleKey, fmt.Sprintf("x/%s", types.ModuleName))
}
