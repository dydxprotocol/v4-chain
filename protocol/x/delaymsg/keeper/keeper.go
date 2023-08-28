package keeper

import (
	sdklog "cosmossdk.io/log"
	"fmt"
	"github.com/cometbft/cometbft/libs/log"
	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4-chain/protocol/x/delaymsg/types"
)

type (
	Keeper struct {
		cdc      codec.BinaryCodec
		storeKey storetypes.StoreKey
		// authorities stores addresses capable of submitting a delayed message.
		authorities map[string]struct{}
		router      *baseapp.MsgServiceRouter
	}
)

// NewKeeper creates a new x/delaymsg keeper.
func NewKeeper(
	cdc codec.BinaryCodec,
	storeKey storetypes.StoreKey,
	router *baseapp.MsgServiceRouter,
	authorities []string,
) *Keeper {
	authoritiesMap := make(map[string]struct{}, len(authorities))
	for _, authority := range authorities {
		authoritiesMap[authority] = struct{}{}
	}
	return &Keeper{
		cdc:         cdc,
		storeKey:    storeKey,
		authorities: authoritiesMap,
		router:      router,
	}
}

// GetAuthorities returns the set of authorities permitted to sign delayed messages.
func (k Keeper) GetAuthorities() map[string]struct{} {
	authorities := make(map[string]struct{}, len(k.authorities))
	for authority := range k.authorities {
		authorities[authority] = struct{}{}
	}
	return authorities
}

// Router returns the x/delaymsg router.
func (k Keeper) Router() *baseapp.MsgServiceRouter {
	return k.router
}

// InitializeForGenesis initializes the x/delaymsg keeper for genesis.
func (k Keeper) InitializeForGenesis(ctx sdk.Context) {
	k.SetNumMessages(ctx, 0)
}

// Logger returns a module-specific logger for x/delaymsg.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With(sdklog.ModuleKey, fmt.Sprintf("x/%s", types.ModuleName))
}
