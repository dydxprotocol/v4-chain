package keeper

import (
	"fmt"

	sdklog "cosmossdk.io/log"

	"github.com/cometbft/cometbft/libs/log"

	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	bridgeserver "github.com/dydxprotocol/v4-chain/protocol/daemons/server/types/bridge"
	"github.com/dydxprotocol/v4-chain/protocol/x/bridge/types"
)

type (
	Keeper struct {
		cdc                codec.BinaryCodec
		storeKey           storetypes.StoreKey
		bridgeEventManager *bridgeserver.BridgeEventManager
		bankKeeper         types.BankKeeper

		// The address capable of executing MsgUpdateEventParams, MsgUpdateProposeParams, and
		// MsgUpdateSafetyParams messages. Typically, this should be the x/gov module account.
		govAuthority string
		// The address capable of executing a MsgCompleteBridge message. Typically, this
		// should be the x/bridge module account.
		bridgeAuthority string
	}
)

func NewKeeper(
	cdc codec.BinaryCodec,
	storeKey storetypes.StoreKey,
	bridgeEventManager *bridgeserver.BridgeEventManager,
	bankKeeper types.BankKeeper,
	govAuthority string,
) *Keeper {
	return &Keeper{
		cdc:                cdc,
		storeKey:           storeKey,
		bridgeEventManager: bridgeEventManager,
		bankKeeper:         bankKeeper,
		govAuthority:       govAuthority,
		bridgeAuthority:    authtypes.NewModuleAddress(types.ModuleName).String(),
	}
}

// GetGovAuthority returns the x/bridge module's authority for updating parameters.
func (k Keeper) GetGovAuthority() string {
	return k.govAuthority
}

// GetBridgeAuthority returns the x/bridge module's authority for completing bridges.
func (k Keeper) GetBridgeAuthority() string {
	return k.bridgeAuthority
}

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With(sdklog.ModuleKey, fmt.Sprintf("x/%s", types.ModuleName))
}

func (k Keeper) InitializeForGenesis(ctx sdk.Context) {
}
