package price_fetcher

import (
	"github.com/dydxprotocol/v4-chain/protocol/daemons/pricefeed/client/types"
	"golang.org/x/exp/maps"
	"gopkg.in/typ.v4/lists"
	"sync"
)

// mutableState stores the mutable state of the price fetcher.
// These parameters are updated periodically by a go routine that polls the pricefeed server
// for the current exchange and market configs. By keeping these parameters in a separate
// struct, we can better ensure that the price fetcher is respecting the synchronization
// requirements of accessing this data in order to stay in a consistent state.
// All access is synchronized by a RW mutex.
type mutableState struct {
	// Market configuration for the price fetcher is synchronized as it is subject to change over time.
	sync.Mutex
	// Access to all following fields is protected.
	// mutableExchangeConfig contains a copy of the current MutableExchangeMarketConfig for the exchange.
	// It is updated periodically by the daemon's PriceFeedMutableMarketConfig when a change in the
	// exchange's configuration is reported by the pricefeed server.
	mutableExchangeConfig *types.MutableExchangeMarketConfig
	// marketExponenets maps market ids to exponents for all markets supported by the exchange.
	// It is updated every time the price fetcher's `UpdateMutableExchangeConfig` method is called.
	marketExponents map[types.MarketId]types.Exponent
	// marketIdsRing tracks the pointer to the next market to query for the upcoming task loop.
	// It is replaced every time the price fetcher's `UpdateMutableExchangeConfig` method is called.
	marketIdsRing *lists.Ring[types.MarketId]
}

// GetMarketIds returns the current set of markets the price fetcher queries for this exchange.
// This method is synchronized.
func (ms *mutableState) GetMarketIds() []types.MarketId {
	ms.Lock()
	defer ms.Unlock()

	return ms.mutableExchangeConfig.GetMarketIds()
}

// GetNextNMarkets returns the next n markets the price fetcher queries for this exchange,
// and advances the mutable state's ring of market ids so that the next call returns
// the next n markets. This method is synchronized.
func (ms *mutableState) GetNextNMarkets(n int) []types.MarketId {
	// A writer lock is used here because the marketIdsRing pointer is advanced below.
	ms.Lock()
	defer ms.Unlock()

	markets := make([]types.MarketId, 0, n)
	for i := 0; i < n; i++ {
		markets = append(markets, ms.marketIdsRing.Value)
		ms.marketIdsRing = ms.marketIdsRing.Next()
	}
	return markets
}

// GetMarketExponents returns a copy of the current set of market exponents for this exchange.
// This method is synchronized.
func (ms *mutableState) GetMarketExponents() map[types.MarketId]types.Exponent {
	ms.Lock()
	defer ms.Unlock()

	return maps.Clone(ms.marketExponents)
}

// GetMutableExchangeConfig returns a copy of the current MutableExchangeMarketConfig for the exchange.
// This method is synchronized.
func (ms *mutableState) GetMutableExchangeConfig() *types.MutableExchangeMarketConfig {
	ms.Lock()
	defer ms.Unlock()

	return ms.mutableExchangeConfig.Copy()
}

// Update updates all fields in the mutableState atomically. Updates are required
// to be atomic in order to keep the mutableState consistent. This method expects
// validation to occur in the PriceFetcher. This method is synchronized.
func (ms *mutableState) Update(
	config *types.MutableExchangeMarketConfig,
	marketExponents map[types.MarketId]types.Exponent,
	marketIdsRing *lists.Ring[types.MarketId],
) {
	ms.Lock()
	defer ms.Unlock()

	ms.mutableExchangeConfig = config
	ms.marketExponents = marketExponents
	ms.marketIdsRing = marketIdsRing
}

// getTaskLoopDefinition returns a snapshot of the current price fetcher mutable state, while
// advancing the price fetcher's marketIdsRing to the next set of markets to query for
// single-market exchanges.
// This method is used to prevent R/W collisions on the price fetcher's mutable state from
// putting the task loop into an inconsistent state.
func (ms *mutableState) getTaskLoopDefinition(
	isMultiMarket bool,
	maxQueries int,
) (
	definition *taskLoopDefinition,
) {
	ms.Lock()
	defer ms.Unlock()

	// Compute which markets to query for this task loop.
	// For single-market exchanges, we want to query the next set of markets in the ring.
	var marketIds []types.MarketId
	if isMultiMarket {
		marketIds = ms.mutableExchangeConfig.GetMarketIds()
	} else {
		marketIds = make([]types.MarketId, 0, maxQueries)
		for i := 0; i < maxQueries; i++ {
			marketIds = append(marketIds, ms.marketIdsRing.Value)
			ms.marketIdsRing = ms.marketIdsRing.Next()
		}
	}

	// Create a copy of all state to pass to the task loop.
	return &taskLoopDefinition{
		mutableExchangeConfig: ms.mutableExchangeConfig.Copy(),
		marketExponents:       maps.Clone(ms.marketExponents),
		marketIds:             marketIds,
	}
}
