package metrics

import (
	"github.com/dydxprotocol/v4-chain/protocol/daemons/pricefeed/client/types"
	"sync"
)

var (
	marketToPair = map[types.MarketId]string{}
	lock         sync.RWMutex
)

// AddMarketPairForTelemetry adds a market pair to an in-memory map of marketId to marketPair strings used
// for labelling metrics.
func AddMarketPairForTelemetry(marketId types.MarketId, marketPair string) {
	lock.Lock()
	defer lock.Unlock()
	marketToPair[marketId] = marketPair
}

// GetMarketPairForTelemetry returns the market pair string for a given marketId. If the marketId is not
// found in the map, returns the INVALID string.
func GetMarketPairForTelemetry(marketId types.MarketId) string {
	lock.RLock()
	defer lock.RUnlock()

	marketPair, exists := marketToPair[marketId]
	if !exists {
		return INVALID
	}

	return marketPair
}
