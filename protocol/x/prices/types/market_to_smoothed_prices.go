package types

// MarketToSmoothedPrices is shorthand for the type that maps Market Ids to this validator's exponentiated local index
// price, to which we are also applying exponential smoothing on every block.
type MarketToSmoothedPrices map[uint32]uint64

// NewMarketToSmoothedPrices returns a new empty map of market ids to market prices.
func NewMarketToSmoothedPrices() MarketToSmoothedPrices {
	return make(MarketToSmoothedPrices, 0)
}
