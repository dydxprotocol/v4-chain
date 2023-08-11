package keeper

import (
	"fmt"

	"github.com/cometbft/cometbft/libs/log"
	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	pricefeedtypes "github.com/dydxprotocol/v4/daemons/server/types/pricefeed"
	"github.com/dydxprotocol/v4/indexer/indexer_manager"
	"github.com/dydxprotocol/v4/lib"
	"github.com/dydxprotocol/v4/x/prices/types"
)

type (
	Keeper struct {
		cdc                    codec.BinaryCodec
		storeKey               storetypes.StoreKey
		indexPriceCache        *pricefeedtypes.MarketToExchangePrices
		marketToSmoothedPrices types.MarketToSmoothedPrices
		timeProvider           lib.TimeProvider
		indexerEventManager    indexer_manager.IndexerEventManager
	}
)

var _ types.PricesKeeper = &Keeper{}

func NewKeeper(
	cdc codec.BinaryCodec,
	storeKey storetypes.StoreKey,
	indexPriceCache *pricefeedtypes.MarketToExchangePrices,
	marketToSmoothedPrices types.MarketToSmoothedPrices,
	timeProvider lib.TimeProvider,
	indexerEventManager indexer_manager.IndexerEventManager,
) *Keeper {
	return &Keeper{
		cdc:                    cdc,
		storeKey:               storeKey,
		indexPriceCache:        indexPriceCache,
		marketToSmoothedPrices: marketToSmoothedPrices,
		timeProvider:           timeProvider,
		indexerEventManager:    indexerEventManager,
	}
}

func (k Keeper) GetIndexerEventManager() indexer_manager.IndexerEventManager {
	return k.indexerEventManager
}

func (k Keeper) InitializeForGenesis(ctx sdk.Context) {
	k.setNumMarkets(ctx, uint32(0))
}

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}
