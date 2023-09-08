package metrics

import (
	"github.com/dydxprotocol/v4-chain/protocol/daemons/pricefeed/client/types"
	"sync"
)

var (
	// marketToPair is a map of marketId to marketPair strings used for labelling metrics.
	// This map is populated whenever markets are created or updated and access to this map is
	// synchronized by the below mutex.
	marketToPair = map[types.MarketId]string{}
	// lock syncronizes access to the marketToPair map.
	lock sync.RWMutex
)

// AddMarketPairForTelemetry adds a market pair to an in-memory map of marketId to marketPair strings used
// for labelling metrics. This method is synchronized.
func AddMarketPairForTelemetry(marketId types.MarketId, marketPair string) {
	lock.Lock()
	defer lock.Unlock()
	marketToPair[marketId] = marketPair
}

// GetMarketPairForTelemetry returns the market pair string for a given marketId. If the marketId is not
// found in the map, returns the INVALID string. This method is synchronized.
func GetMarketPairForTelemetry(marketId types.MarketId) string {
	lock.RLock()
	defer lock.RUnlock()

	marketPair, exists := marketToPair[marketId]
	if !exists {
		return INVALID
	}

	return marketPair
}
