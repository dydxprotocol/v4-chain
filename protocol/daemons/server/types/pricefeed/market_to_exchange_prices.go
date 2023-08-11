package types

import (
	"sync"
	"time"

	gometrics "github.com/armon/go-metrics"
	"github.com/cosmos/cosmos-sdk/telemetry"
	"github.com/dydxprotocol/v4/daemons/pricefeed"
	"github.com/dydxprotocol/v4/daemons/pricefeed/api"
	pricefeedmetrics "github.com/dydxprotocol/v4/daemons/pricefeed/metrics"
	"github.com/dydxprotocol/v4/lib"
	"github.com/dydxprotocol/v4/lib/metrics"
	"github.com/dydxprotocol/v4/x/prices/types"
)

// MarketToExchangePrices maintains price info for multiple markets. Each
// market can support prices from multiple exchange sources. Specifically,
// MarketToExchangePrices supports methods to update prices and to retrieve
// median prices. Methods are goroutine safe.
type MarketToExchangePrices struct {
	sync.RWMutex                                       // reader-writer lock
	marketToExchangePrices map[uint32]*ExchangeToPrice // {k: market id, v: exchange prices}
}

// NewMarketToExchangePrices creates a new MarketToExchangePrices.
func NewMarketToExchangePrices() *MarketToExchangePrices {
	return &MarketToExchangePrices{
		marketToExchangePrices: make(map[uint32]*ExchangeToPrice),
	}
}

// UpdatePrices updates market prices given a list of price updates. Prices are
// only updated if the timestamp on the updates are greater than the timestamp
// on existing prices.
func (mte *MarketToExchangePrices) UpdatePrices(
	updates []*api.MarketPriceUpdate) {
	mte.Lock()
	defer mte.Unlock()
	for _, marketPriceUpdate := range updates {
		marketId := marketPriceUpdate.MarketId
		exchangeToPrices, ok := mte.marketToExchangePrices[marketId]
		if !ok {
			exchangeToPrices = NewExchangeToPrice()
			mte.marketToExchangePrices[marketId] = exchangeToPrices
		}
		exchangeToPrices.UpdatePrices(marketPriceUpdate.ExchangePrices)
	}
}

// GetValidMedianPrices returns median prices for multiple markets.
// Specifically, it returns a map where the key is the market ID and the value
// is the median price for the market. It only returns "valid" prices where
// a price is valid iff
// 1) the last update time is within a predefined threshold away from the given
// read time.
// 2) the exchange where the price comes from is included in the given list of
// accepted exchanges.
// 3) the number of prices that meet 1) and 2) are greater than the minimum
// number of exchanges specified in the given input.
func (mte *MarketToExchangePrices) GetValidMedianPrices(
	markets []types.Market,
	readTime time.Time,
) map[uint32]uint64 {
	cutoffTime := readTime.Add(-pricefeed.MaxPriceAge)
	marketIdToValidExchangeFeedIds := getMarketIdToValidExchangeFeedIds(markets)
	marketIdToMedianPrice := make(map[uint32]uint64)

	mte.RLock()
	defer mte.RUnlock()
	for _, market := range markets {
		marketId := market.Id
		validExchangeFeedIds, ok := marketIdToValidExchangeFeedIds[marketId]
		if !ok || len(validExchangeFeedIds) == 0 {
			// No valid exchanges, skip this market.
			telemetry.IncrCounterWithLabels(
				[]string{
					metrics.PricefeedServer,
					metrics.NoValidExchanges,
					metrics.Count,
				},
				1,
				[]gometrics.Label{
					pricefeedmetrics.GetLabelForMarketId(marketId),
				},
			)
			continue
		}
		exchangeToPrice, ok := mte.marketToExchangePrices[marketId]
		if !ok {
			// No market price info yet, skip this market.
			telemetry.IncrCounterWithLabels(
				[]string{
					metrics.PricefeedServer,
					metrics.NoMarketPrice,
					metrics.Count,
				},
				1,
				[]gometrics.Label{
					pricefeedmetrics.GetLabelForMarketId(marketId),
				},
			)
			continue
		}

		// GetValidPrice filters prices based on valid exchanges and cutoff time.
		validPrices := exchangeToPrice.GetValidPrices(validExchangeFeedIds, cutoffTime)

		// The number of valid prices must be >= min number of exchanges.
		if len(validPrices) >= int(market.MinExchanges) {
			// Calculate the median. Returns an error if the input is empty.
			median, err := lib.MedianUint64(validPrices)
			if err != nil {
				telemetry.IncrCounterWithLabels(
					[]string{
						metrics.PricefeedServer,
						metrics.NoValidMedianPrice,
						metrics.Count,
					},
					1,
					[]gometrics.Label{
						pricefeedmetrics.GetLabelForMarketId(marketId),
					},
				)
				continue
			}
			marketIdToMedianPrice[marketId] = median
		}
	}

	return marketIdToMedianPrice
}

// getMarketIdToValidExchangeFeedIds returns a map {k: market id, v: map of valid exchange ids}
// given a list of markets.
func getMarketIdToValidExchangeFeedIds(
	markets []types.Market,
) map[uint32]map[uint32]bool {
	marketIdToValidExchangeFeedIds := make(map[uint32]map[uint32]bool, len(markets))
	for _, market := range markets {
		exchangeFeedIdsSet, ok := marketIdToValidExchangeFeedIds[market.Id]
		if !ok {
			exchangeFeedIdsSet = make(map[uint32]bool, len(market.Exchanges))
			marketIdToValidExchangeFeedIds[market.Id] = exchangeFeedIdsSet
		}
		for _, exchangeFeedId := range market.Exchanges {
			exchangeFeedIdsSet[exchangeFeedId] = true
		}
	}
	return marketIdToValidExchangeFeedIds
}
