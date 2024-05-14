package keeper

import (
	"fmt"
	"sync/atomic"
	"time"

	"cosmossdk.io/log"
	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	pricefeedtypes "github.com/dydxprotocol/v4-chain/protocol/daemons/server/types/pricefeed"
	"github.com/dydxprotocol/v4-chain/protocol/indexer/indexer_manager"
	"github.com/dydxprotocol/v4-chain/protocol/lib"
	libtime "github.com/dydxprotocol/v4-chain/protocol/lib/time"
	"github.com/dydxprotocol/v4-chain/protocol/x/prices/types"
)

type (
	Keeper struct {
		cdc                            codec.BinaryCodec
		storeKey                       storetypes.StoreKey
		indexPriceCache                *pricefeedtypes.MarketToExchangePrices
		timeProvider                   libtime.TimeProvider
		indexerEventManager            indexer_manager.IndexerEventManager
		marketToCreatedAt              map[uint32]time.Time
		authorities                    map[string]struct{}
		currencyPairIDCache            *CurrencyPairIDCache
		currencyPairIdCacheInitialized *atomic.Bool
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
) *Keeper {
	return &Keeper{
		cdc:                            cdc,
		storeKey:                       storeKey,
		indexPriceCache:                indexPriceCache,
		timeProvider:                   timeProvider,
		indexerEventManager:            indexerEventManager,
		marketToCreatedAt:              map[uint32]time.Time{},
		authorities:                    lib.UniqueSliceToSet(authorities),
		currencyPairIDCache:            NewCurrencyPairIDCache(),
		currencyPairIdCacheInitialized: &atomic.Bool{}, // Initialized to false
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

func (k Keeper) InitializeCurrencyPairIdCache(ctx sdk.Context) {
	alreadyInitialized := k.currencyPairIdCacheInitialized.Swap(true)
	if alreadyInitialized {
		return
	}

	// Load the currency pair IDs for the markets from the x/prices state.
	k.LoadCurrencyPairIDCache(ctx)
}

// func (k Keeper) DebugUpdatePrice(ctx sdk.Context) {
// 	price := uint64(2000000000)
// 	k.UpdateMarketPrices(
// 		ctx,
// 		[]*types.MsgUpdateMarketPrices_MarketPrice{
// 			{
// 				MarketId: 33,
// 				// Price:    1000000,
// 				// Price:    90000000,
// 				Price: price,
// 			},
// 		},
// 	)

// 	fmt.Printf("!!!! Updated price for market 33, %v\n", price)
// }
