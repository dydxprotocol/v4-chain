package keeper

import (
	"fmt"
	"time"

	sdklog "cosmossdk.io/log"

	"github.com/cometbft/cometbft/libs/log"
	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	pricefeedtypes "github.com/dydxprotocol/v4-chain/protocol/daemons/server/types/pricefeed"
	"github.com/dydxprotocol/v4-chain/protocol/indexer/indexer_manager"
	"github.com/dydxprotocol/v4-chain/protocol/lib"
	libtime "github.com/dydxprotocol/v4-chain/protocol/lib/time"
	"github.com/dydxprotocol/v4-chain/protocol/x/prices/types"
)

type (
	Keeper struct {
		cdc                    codec.BinaryCodec
		storeKey               storetypes.StoreKey
		indexPriceCache        *pricefeedtypes.MarketToExchangePrices
		marketToSmoothedPrices types.MarketToSmoothedPrices
		timeProvider           libtime.TimeProvider
		indexerEventManager    indexer_manager.IndexerEventManager
		marketToCreatedAt      map[uint32]time.Time
		authorities            map[string]struct{}
	}
)

var _ types.PricesKeeper = &Keeper{}

func NewKeeper(
	cdc codec.BinaryCodec,
	storeKey storetypes.StoreKey,
	indexPriceCache *pricefeedtypes.MarketToExchangePrices,
	marketToSmoothedPrices types.MarketToSmoothedPrices,
	timeProvider libtime.TimeProvider,
	indexerEventManager indexer_manager.IndexerEventManager,
	authorities []string,
) *Keeper {
	return &Keeper{
		cdc:                    cdc,
		storeKey:               storeKey,
		indexPriceCache:        indexPriceCache,
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
	return ctx.Logger().With(sdklog.ModuleKey, fmt.Sprintf("x/%s", types.ModuleName))
}

func (k Keeper) HasAuthority(authority string) bool {
	_, ok := k.authorities[authority]
	return ok
}
