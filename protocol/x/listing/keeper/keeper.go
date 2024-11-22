package keeper

import (
	"fmt"

	"cosmossdk.io/log"
	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4-chain/protocol/indexer/indexer_manager"
	"github.com/dydxprotocol/v4-chain/protocol/lib"
	"github.com/dydxprotocol/v4-chain/protocol/x/listing/types"
)

type (
	Keeper struct {
		cdc                 codec.BinaryCodec
		storeKey            storetypes.StoreKey
		authorities         map[string]struct{}
		indexerEventManager indexer_manager.IndexerEventManager
		PricesKeeper        types.PricesKeeper
		ClobKeeper          types.ClobKeeper
		MarketMapKeeper     types.MarketMapKeeper
		PerpetualsKeeper    types.PerpetualsKeeper
		SubaccountsKeeper   types.SubaccountsKeeper
		VaultKeeper         types.VaultKeeper
	}
)

func NewKeeper(
	cdc codec.BinaryCodec,
	storeKey storetypes.StoreKey,
	authorities []string,
	indexerEventsManager indexer_manager.IndexerEventManager,
	pricesKeeper types.PricesKeeper,
	clobKeeper types.ClobKeeper,
	marketMapKeeper types.MarketMapKeeper,
	perpetualsKeeper types.PerpetualsKeeper,
	subaccountsKeeper types.SubaccountsKeeper,
	vaultKeeper types.VaultKeeper,
) *Keeper {
	return &Keeper{
		cdc:                 cdc,
		storeKey:            storeKey,
		authorities:         lib.UniqueSliceToSet(authorities),
		indexerEventManager: indexerEventsManager,
		PricesKeeper:        pricesKeeper,
		ClobKeeper:          clobKeeper,
		MarketMapKeeper:     marketMapKeeper,
		PerpetualsKeeper:    perpetualsKeeper,
		SubaccountsKeeper:   subaccountsKeeper,
		VaultKeeper:         vaultKeeper,
	}
}

func (k Keeper) HasAuthority(authority string) bool {
	_, ok := k.authorities[authority]
	return ok
}

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With(log.ModuleKey, fmt.Sprintf("x/%s", types.ModuleName))
}

func (k Keeper) InitializeForGenesis(ctx sdk.Context) {}

func (k Keeper) GetIndexerEventManager() indexer_manager.IndexerEventManager {
	return k.indexerEventManager
}
