package keeper

import (
	"fmt"

	"github.com/cometbft/cometbft/libs/log"

	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	bridgeserver "github.com/dydxprotocol/v4/daemons/server/types/bridge"
	"github.com/dydxprotocol/v4/x/bridge/types"
)

type (
	Keeper struct {
		cdc                codec.BinaryCodec
		storeKey           storetypes.StoreKey
		bridgeEventManager *bridgeserver.BridgeEventManager
	}
)

func NewKeeper(
	cdc codec.BinaryCodec,
	storeKey storetypes.StoreKey,
	bridgeEventManager *bridgeserver.BridgeEventManager,
) *Keeper {
	return &Keeper{
		cdc:                cdc,
		storeKey:           storeKey,
		bridgeEventManager: bridgeEventManager,
	}
}

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

func (k Keeper) InitializeForGenesis(ctx sdk.Context) {
	k.SetNextAcknowledgedEventId(ctx, 0)
}
