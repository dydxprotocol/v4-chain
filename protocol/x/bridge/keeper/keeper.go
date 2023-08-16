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
		authority string
		// The address capable of executing a MsgCompleteBridge message. Typically, this
		// should be the x/bridge module account.
		selfAuthority string
	}
)

func NewKeeper(
	cdc codec.BinaryCodec,
	storeKey storetypes.StoreKey,
	bridgeEventManager *bridgeserver.BridgeEventManager,
	bankKeeper types.BankKeeper,
	authority string,
) *Keeper {
	return &Keeper{
		cdc:                cdc,
		storeKey:           storeKey,
		bridgeEventManager: bridgeEventManager,
		bankKeeper:         bankKeeper,
		authority:          authority,
		selfAuthority:      authtypes.NewModuleAddress(types.ModuleName).String(),
	}
}

// GetAuthority returns the x/bridge module's authority for updating parameters.
func (k Keeper) GetAuthority() string {
	return k.authority
}

// GetSelfAuthority returns the x/bridge module's authority for completing bridges.
func (k Keeper) GetSelfAuthority() string {
	return k.selfAuthority
}

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With(sdklog.ModuleKey, fmt.Sprintf("x/%s", types.ModuleName))
}

func (k Keeper) InitializeForGenesis(ctx sdk.Context) {
}
