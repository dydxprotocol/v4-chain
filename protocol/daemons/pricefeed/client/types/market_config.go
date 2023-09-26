package types

// MarketConfig specifies the exchange-specific market configuration used to resolve a market's price on
// a particular exchange.
type MarketConfig struct {
	// Ticker specifies the string to use to query the relevant ticker price for a market on this exchange.
	Ticker string

	// AdjustByMarket optionally identifies the appropriate market that should be used to adjust
	// the price of the market ticker to arrive at the final USD price of the market. This is used in
	// cases where we choose to query a price for a market on a particular exchange in a different quote
	// currency - perhaps because that market is more robust - and convert the price back to the quote
	// currency of the original market with another market's price.
	//
	// For example, for resolving the BTC-USD price on this exchange, we may use a "BTC-USDT" ticker with
	// an adjust-by market of USDT-USD, and compute BTC-USD as
	//
	// BTC-USD = BTC-USDT * USDT-USD
	//
	// If this field is nil, then the market has no adjust-by market.
	AdjustByMarket *MarketId

	// Invert specifies the inversion strategy to use when converting the price of the market's ticker to arrive
	// at the final USD price of the market. The application of the inversion strategy depends on whether an adjust-by
	// market is defined for this market.
	//
	// If an adjust-by market is defined for this market, then the inversion strategy is applied with respect
	// to the adjustment market. For example, say we use a "BTC-USDT" ticker for USDT-USD on this exchange, with
	// an adjust-by market of BTC-USD, and an inversion value of true. In this case, we are describing that
	// we will derive the BTC-USD price by multiplying the BTC-USD index price by the inverse of the BTC-USDT ticker
	// price:
	//
	// USDT-USD = BTC-USD / BTC-USDT
	//
	// If an adjust-by market is not defined for this market, then the inversion strategy is applied to the ticker
	// price itself. For example, for BTC, say we use "USD-BTC" as the BTC-USD ticker on this exchange with an
	// inversion value of true. In that case, we would derive the BTC-USD price by taking the inverse of the
	// USD-BTC price:
	//
	// BTC-USD = 1 / USD-BTC
	Invert bool
}

// Equal returns true if the two MarketConfigs are equal.
func (mc *MarketConfig) Equal(other MarketConfig) bool {
	return mc.Ticker == other.Ticker &&
		mc.Invert == other.Invert &&
		((mc.AdjustByMarket == nil && other.AdjustByMarket == nil) ||
			(mc.AdjustByMarket != nil && other.AdjustByMarket != nil &&
				*mc.AdjustByMarket == *other.AdjustByMarket))
}

// Copy returns a deep copy of the MarketConfig.
func (mc *MarketConfig) Copy() MarketConfig {
	var adjustByMarket *MarketId
	if mc.AdjustByMarket != nil {
		adjustByMarket = new(MarketId)
		*adjustByMarket = *mc.AdjustByMarket
	}
	return MarketConfig{
		Ticker:         mc.Ticker,
		AdjustByMarket: adjustByMarket,
		Invert:         mc.Invert,
	}
}
