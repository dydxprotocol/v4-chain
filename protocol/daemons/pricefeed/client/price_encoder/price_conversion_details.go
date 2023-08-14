package price_encoder

import "github.com/dydxprotocol/v4-chain/protocol/daemons/pricefeed/client/types"

// priceConversionDetails contains all information required to convert a ticker price from an exchange API response
// into a market price in the correct quote currency for the market (which, at this time, is uniformly USD.)
// If `adjustByMarketDetails` is nil, then the price is not adjusted by another market's price.
type priceConversionDetails struct {
	Invert                bool
	Exponent              types.Exponent
	AdjustByMarketDetails *adjustByMarketDetails
}
