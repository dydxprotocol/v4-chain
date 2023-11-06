package metrics

import (
	"fmt"
	"sync"

	"github.com/dydxprotocol/v4-chain/protocol/daemons/pricefeed/client/types"
)

var (
	// marketToPair maps live marketIds to marketPair strings and is used for labelling metrics.
	// This map is populated whenever markets are created or updated and access to this map is
	// synchronized by the below mutex.
	// The most correct approach here would be to get the marketPair from current chain state, so that
	// we do not accidentally capture outdated or incorrect labels for market updates from rejected blocks.
	// (A rejected block with a market creation should not present a problem as we will not log regular
	// metrics for this market id unless it is later re-assigned - and re-updated.) However, these labels
	// are used across the prices keeper, daemon, and daemon server code, so there is motive to avoid
	// a solution that requires propagating references to centralized state. Furthermore, we judge that
	// these market pairs are very unlikely to be updated, so this solution, while not perfect, is
	// acceptable for the use case of logging/metrics in order to manage code complexity.
	marketToPair = map[types.MarketId]string{}
	// lock synchronizes access to the marketToPair map.
	lock sync.Mutex
)

// SetMarketPairForTelemetry sets a market pair to an in-memory map of marketId to marketPair strings used
// for labelling metrics. This method is synchronized.
func SetMarketPairForTelemetry(marketId types.MarketId, marketPair string) {
	lock.Lock()
	defer lock.Unlock()
	marketToPair[marketId] = marketPair
}

// GetMarketPairForTelemetry returns the market pair string for a given marketId. If the marketId is not
// found in the map, returns the INVALID string. This method is synchronized.
func GetMarketPairForTelemetry(marketId types.MarketId) string {
	lock.Lock()
	defer lock.Unlock()

	marketPair, exists := marketToPair[marketId]
	if !exists {
		return fmt.Sprintf("invalid_id:%v", marketId)
	}

	return marketPair
}
