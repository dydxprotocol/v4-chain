package testutil

import (
	"sort"

	"github.com/dydxprotocol/v4-chain/protocol/daemons/pricefeed/client/types"
)

func GetTickersSortedByMarketId(marketToMarketConfig map[uint32]types.MarketConfig) []string {
	// Get all `marketId`s in `marketIdToTicker` as a sorted array.
	marketIds := make([]uint32, 0, len(marketToMarketConfig))
	for marketId := range marketToMarketConfig {
		marketIds = append(marketIds, marketId)
	}
	sort.Slice(marketIds, func(i, j int) bool {
		return marketIds[i] < marketIds[j]
	})

	// Get a list of tickers sorted by their corresponding `marketId`.
	tickers := make([]string, 0, len(marketToMarketConfig))
	for _, marketId := range marketIds {
		tickers = append(tickers, marketToMarketConfig[marketId].Ticker)
	}

	return tickers
}
