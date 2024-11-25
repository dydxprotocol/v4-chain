package keeper

import (
	"fmt"

	"cosmossdk.io/log"
	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	pricefeedtypes "github.com/dydxprotocol/v4-chain/protocol/daemons/server/types/pricefeed"
	"github.com/dydxprotocol/v4-chain/protocol/indexer/indexer_manager"
	"github.com/dydxprotocol/v4-chain/protocol/lib"
	libtime "github.com/dydxprotocol/v4-chain/protocol/lib/time"
	streamingtypes "github.com/dydxprotocol/v4-chain/protocol/streaming/types"
	"github.com/dydxprotocol/v4-chain/protocol/x/prices/types"
)

type (
	Keeper struct {
		cdc                 codec.BinaryCodec
		storeKey            storetypes.StoreKey
		indexPriceCache     *pricefeedtypes.MarketToExchangePrices
		timeProvider        libtime.TimeProvider
		indexerEventManager indexer_manager.IndexerEventManager
		authorities         map[string]struct{}
		RevShareKeeper      types.RevShareKeeper
		MarketMapKeeper     types.MarketMapKeeper

		streamingManager streamingtypes.FullNodeStreamingManager
	}
)

var _ types.PricesKeeper = &Keeper{}

func NewKeeper(
	cdc codec.BinaryCodec,
	storeKey storetypes.StoreKey,
	indexPriceCache *pricefeedtypes.MarketToExchangePrices,
	timeProvider libtime.TimeProvider,
	indexerEventManager indexer_manager.IndexerEventManager,
	authorities []string,
	revShareKeeper types.RevShareKeeper,
	marketMapKeeper types.MarketMapKeeper,
	streamingManager streamingtypes.FullNodeStreamingManager,
) *Keeper {
	return &Keeper{
		cdc:                 cdc,
		storeKey:            storeKey,
		indexPriceCache:     indexPriceCache,
		timeProvider:        timeProvider,
		indexerEventManager: indexerEventManager,
		authorities:         lib.UniqueSliceToSet(authorities),
		RevShareKeeper:      revShareKeeper,
		MarketMapKeeper:     marketMapKeeper,
		streamingManager:    streamingManager,
	}
}

func (k Keeper) GetIndexerEventManager() indexer_manager.IndexerEventManager {
	return k.indexerEventManager
}

func (k Keeper) InitializeForGenesis(ctx sdk.Context) {
}

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With(log.ModuleKey, fmt.Sprintf("x/%s", types.ModuleName))
}

func (k Keeper) HasAuthority(authority string) bool {
	_, ok := k.authorities[authority]
	return ok
}

func (k Keeper) GetFullNodeStreamingManager() streamingtypes.FullNodeStreamingManager {
	return k.streamingManager
}
