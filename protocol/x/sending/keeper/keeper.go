package keeper

import (
	"fmt"

	"github.com/cometbft/cometbft/libs/log"
	"github.com/dydxprotocol/v4-chain/protocol/indexer/indexer_manager"

	sdklog "cosmossdk.io/log"
	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4-chain/protocol/x/sending/types"
)

type (
	Keeper struct {
		cdc                 codec.BinaryCodec
		storeKey            storetypes.StoreKey
		accountKeeper       types.AccountKeeper
		subaccountsKeeper   types.SubaccountsKeeper
		indexerEventManager indexer_manager.IndexerEventManager
	}
)

func NewKeeper(
	cdc codec.BinaryCodec,
	storeKey storetypes.StoreKey,
	accountKeeper types.AccountKeeper,
	subaccountsKeeper types.SubaccountsKeeper,
	indexerEventManager indexer_manager.IndexerEventManager,
) *Keeper {
	return &Keeper{
		cdc:                 cdc,
		storeKey:            storeKey,
		accountKeeper:       accountKeeper,
		subaccountsKeeper:   subaccountsKeeper,
		indexerEventManager: indexerEventManager,
	}
}

func (k Keeper) GetIndexerEventManager() indexer_manager.IndexerEventManager {
	return k.indexerEventManager
}

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With(sdklog.ModuleKey, fmt.Sprintf("x/%s", types.ModuleName))
}

func (k Keeper) InitializeForGenesis(ctx sdk.Context) {
}
