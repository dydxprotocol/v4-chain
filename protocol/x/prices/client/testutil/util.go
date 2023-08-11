package testutil

import (
	"sort"
)

func GetTickersSortedByMarketId(marketIdToTicker map[uint32]string) []string {
	// Get all `marketId`s in `marketIdToTicker` as a sorted array.
	marketIds := make([]uint32, 0, len(marketIdToTicker))
	for marketId := range marketIdToTicker {
		marketIds = append(marketIds, marketId)
	}
	sort.Slice(marketIds, func(i, j int) bool {
		return marketIds[i] < marketIds[j]
	})

	// Get a list of tickers sorted by their corresponding `marketId`.
	tickers := make([]string, 0, len(marketIdToTicker))
	for _, marketId := range marketIds {
		tickers = append(tickers, marketIdToTicker[marketId])
	}

	return tickers
}
