package keeper

import (
	"fmt"

	"github.com/cometbft/cometbft/libs/log"
	"github.com/dydxprotocol/v4/indexer/indexer_manager"

	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4/x/sending/types"
)

type (
	Keeper struct {
		cdc                 codec.BinaryCodec
		storeKey            storetypes.StoreKey
		accountKeeper       types.AccountKeeper
		subaccountsKeeper   types.SubaccountsKeeper
		indexerEventManager *indexer_manager.IndexerEventManagerImpl
	}
)

func NewKeeper(
	cdc codec.BinaryCodec,
	storeKey storetypes.StoreKey,
	accountKeeper types.AccountKeeper,
	subaccountsKeeper types.SubaccountsKeeper,
	indexerEventManager *indexer_manager.IndexerEventManagerImpl,
) *Keeper {
	return &Keeper{
		cdc:                 cdc,
		storeKey:            storeKey,
		accountKeeper:       accountKeeper,
		subaccountsKeeper:   subaccountsKeeper,
		indexerEventManager: indexerEventManager,
	}
}

func (k Keeper) GetIndexerEventManager() *indexer_manager.IndexerEventManagerImpl {
	return k.indexerEventManager
}

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

func (k Keeper) InitializeForGenesis(ctx sdk.Context) {
}
