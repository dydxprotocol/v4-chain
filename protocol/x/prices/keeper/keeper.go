package keeper

import (
	"fmt"
	"time"

	"cosmossdk.io/log"
	storetypes "cosmossdk.io/store/types"

	pricefeedtypes "github.com/StreamFinance-Protocol/stream-chain/protocol/daemons/server/types/pricefeed"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/indexer/indexer_manager"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/lib"
	libtime "github.com/StreamFinance-Protocol/stream-chain/protocol/lib/time"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/x/prices/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type Keeper struct {
	cdc                    codec.BinaryCodec
	storeKey               storetypes.StoreKey
	marketToSmoothedPrices types.MarketToSmoothedPrices
	DaemonPriceCache       *pricefeedtypes.MarketToExchangePrices
	timeProvider           libtime.TimeProvider
	indexerEventManager    indexer_manager.IndexerEventManager
	marketToCreatedAt      map[uint32]time.Time
	authorities            map[string]struct{}
}

var _ types.PricesKeeper = &Keeper{}

func NewKeeper(
	cdc codec.BinaryCodec,
	storeKey storetypes.StoreKey,
	daemonPriceCache *pricefeedtypes.MarketToExchangePrices,
	marketToSmoothedPrices types.MarketToSmoothedPrices,
	timeProvider libtime.TimeProvider,
	indexerEventManager indexer_manager.IndexerEventManager,
	authorities []string,
) *Keeper {
	return &Keeper{
		cdc:                    cdc,
		storeKey:               storeKey,
		DaemonPriceCache:       daemonPriceCache,
		marketToSmoothedPrices: marketToSmoothedPrices,
		timeProvider:           timeProvider,
		indexerEventManager:    indexerEventManager,
		marketToCreatedAt:      map[uint32]time.Time{},
		authorities:            lib.UniqueSliceToSet(authorities),
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
