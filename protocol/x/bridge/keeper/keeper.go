package keeper

import (
	"fmt"

	"cosmossdk.io/log"
	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	bridgeserver "github.com/dydxprotocol/v4-chain/protocol/daemons/server/types/bridge"
	"github.com/dydxprotocol/v4-chain/protocol/lib"
	"github.com/dydxprotocol/v4-chain/protocol/x/bridge/types"
	delaymsgtypes "github.com/dydxprotocol/v4-chain/protocol/x/delaymsg/types"
)

type (
	Keeper struct {
		cdc                codec.BinaryCodec
		storeKey           storetypes.StoreKey
		bridgeEventManager *bridgeserver.BridgeEventManager
		bankKeeper         types.BankKeeper
		delayMsgKeeper     delaymsgtypes.DelayMsgKeeper

		// authorities stores addresses capable of sending a bridge message.
		authorities map[string]struct{}
	}
)

func NewKeeper(
	cdc codec.BinaryCodec,
	storeKey storetypes.StoreKey,
	bridgeEventManager *bridgeserver.BridgeEventManager,
	bankKeeper types.BankKeeper,
	delayMsgKeeper delaymsgtypes.DelayMsgKeeper,
	authorities []string,
) *Keeper {
	return &Keeper{
		cdc:                cdc,
		storeKey:           storeKey,
		bridgeEventManager: bridgeEventManager,
		bankKeeper:         bankKeeper,
		delayMsgKeeper:     delayMsgKeeper,
		authorities:        lib.UniqueSliceToSet(authorities),
	}
}

// HasAuthority returns whether `authority` exists in `k.authorities`.
func (k Keeper) HasAuthority(authority string) bool {
	_, ok := k.authorities[authority]
	return ok
}

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With(log.ModuleKey, fmt.Sprintf("x/%s", types.ModuleName))
}

func (k Keeper) InitializeForGenesis(ctx sdk.Context) {
}
