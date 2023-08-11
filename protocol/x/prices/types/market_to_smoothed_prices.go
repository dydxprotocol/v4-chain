package types

import (
	"gopkg.in/typ.v4/lists"
)

const (
	// SmoothedPriceTrackingBlockHistoryLength is the number of blocks we track smoothed prices for to determine if the
	// next smoothed price should be used in the price update proposal.
	// TODO(CORE-396): Put this on-chain and configure it in the application via chain state.
	SmoothedPriceTrackingBlockHistoryLength = uint32(5)
)

// MarketToSmoothedPrices tracks current and historical exponentially smoothed prices for each market.
type MarketToSmoothedPrices interface {
	GetSmoothedPrice(marketId uint32) (price uint64, ok bool)
	GetSmoothedPricesForTest() map[uint32]uint64
	GetHistoricalSmoothedPrices(marketId uint32) []uint64
	PushSmoothedPrice(marketId uint32, price uint64)
}

type MarketToSmoothedPricesImpl struct {
	historyLength          uint32
	marketToSmoothedPrices map[uint32]*lists.Ring[uint64]
}

// GetSmoothedPrice returns the smoothed price for the given market.
func (m *MarketToSmoothedPricesImpl) GetSmoothedPrice(marketId uint32) (
	price uint64,
	ok bool,
) {
	smoothedPrices, ok := m.marketToSmoothedPrices[marketId]
	if !ok {
		return 0, false
	}

	price = smoothedPrices.Value

	// If the price is zero, we do not have a valid price for this market.
	if price == 0 {
		return 0, false
	}
	return price, true
}

// GetSmoothedPricesForTest returns a map of market ids to smoothed prices. This is primarily here for testing.
func (m *MarketToSmoothedPricesImpl) GetSmoothedPricesForTest() map[uint32]uint64 {
	smoothedPrices := make(map[uint32]uint64)
	for marketId := range m.marketToSmoothedPrices {
		smoothedPrice, exists := m.GetSmoothedPrice(marketId)
		if exists {
			smoothedPrices[marketId] = smoothedPrice
		}
	}
	return smoothedPrices
}

// GetHistoricalSmoothedPrices returns up to the last `SmoothedPriceTrackingBlockHistoryLength` smoothed prices for the
// given market. The returned slice is ordered from newest to oldest, and the first entry in the slice will be the
// most recent valid smoothed price.
func (m *MarketToSmoothedPricesImpl) GetHistoricalSmoothedPrices(marketId uint32) []uint64 {
	smoothedPrices, ok := m.marketToSmoothedPrices[marketId]
	if !ok {
		return []uint64{}
	}

	prices := make([]uint64, 0, smoothedPrices.Len())
	for i := 0; i < smoothedPrices.Len(); i++ {
		price := smoothedPrices.Value
		if price != 0 {
			prices = append(prices, price)
		}
		// Walk backwards in time
		smoothedPrices = smoothedPrices.Prev()
	}
	return prices
}

// PushSmoothedPrice sets the smoothed price for the given market.
func (m *MarketToSmoothedPricesImpl) PushSmoothedPrice(id uint32, price uint64) {
	smoothedPrices, ok := m.marketToSmoothedPrices[id]
	if !ok {
		smoothedPrices = lists.NewRing[uint64](int(m.historyLength))
	}
	smoothedPrices = smoothedPrices.Next()
	smoothedPrices.Value = price
	m.marketToSmoothedPrices[id] = smoothedPrices
}

// NewMarketToSmoothedPrices returns a new `MarketToSmoothedPrices` that tracks the previous `historyLength` prices per
// market. The default value to use for the protocol is `SmoothedPriceTrackingBlockHistoryLength`.
func NewMarketToSmoothedPrices(historyLength uint32) MarketToSmoothedPrices {
	return &MarketToSmoothedPricesImpl{
		historyLength:          historyLength,
		marketToSmoothedPrices: make(map[uint32]*lists.Ring[uint64]),
	}
}
