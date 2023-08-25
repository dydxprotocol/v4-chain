package keeper

import (
	"fmt"

	sdklog "cosmossdk.io/log"

	"github.com/cometbft/cometbft/libs/log"

	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	bridgeserver "github.com/dydxprotocol/v4-chain/protocol/daemons/server/types/bridge"
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

		// The address capable of executing MsgUpdateEventParams, MsgUpdateProposeParams, and
		// MsgUpdateSafetyParams messages. Typically, this should be the x/gov module account.
		govAuthority string
	}
)

func NewKeeper(
	cdc codec.BinaryCodec,
	storeKey storetypes.StoreKey,
	bridgeEventManager *bridgeserver.BridgeEventManager,
	bankKeeper types.BankKeeper,
	delayMsgKeeper delaymsgtypes.DelayMsgKeeper,
	govAuthority string,
) *Keeper {
	return &Keeper{
		cdc:                cdc,
		storeKey:           storeKey,
		bridgeEventManager: bridgeEventManager,
		bankKeeper:         bankKeeper,
		delayMsgKeeper:     delayMsgKeeper,
		govAuthority:       govAuthority,
	}
}

// GetGovAuthority returns the x/bridge module's authority for updating parameters.
func (k Keeper) GetGovAuthority() string {
	return k.govAuthority
}

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With(sdklog.ModuleKey, fmt.Sprintf("x/%s", types.ModuleName))
}

func (k Keeper) InitializeForGenesis(ctx sdk.Context) {
}
