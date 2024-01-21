package keeper

import (
	"fmt"

	"cosmossdk.io/log"
	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4-chain/protocol/indexer/indexer_manager"
	"github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"
)

type (
	Keeper struct {
		cdc                 codec.BinaryCodec
		storeKey            storetypes.StoreKey
		assetsKeeper        types.AssetsKeeper
		bankKeeper          types.BankKeeper
		perpetualsKeeper    types.PerpetualsKeeper
		blocktimeKeeper     types.BlocktimeKeeper
		indexerEventManager indexer_manager.IndexerEventManager
	}
)

func NewKeeper(
	cdc codec.BinaryCodec,
	storeKey storetypes.StoreKey,
	assetsKeeper types.AssetsKeeper,
	bankKeeper types.BankKeeper,
	perpetualsKeeper types.PerpetualsKeeper,
	blocktimeKeeper types.BlocktimeKeeper,
	indexerEventManager indexer_manager.IndexerEventManager,
) *Keeper {
	return &Keeper{
		cdc:                 cdc,
		storeKey:            storeKey,
		assetsKeeper:        assetsKeeper,
		bankKeeper:          bankKeeper,
		perpetualsKeeper:    perpetualsKeeper,
		blocktimeKeeper:     blocktimeKeeper,
		indexerEventManager: indexerEventManager,
	}
}

func (k Keeper) GetIndexerEventManager() indexer_manager.IndexerEventManager {
	return k.indexerEventManager
}

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With(log.ModuleKey, fmt.Sprintf("x/%s", types.ModuleName))
}

func (k Keeper) InitializeForGenesis(ctx sdk.Context) {
}
