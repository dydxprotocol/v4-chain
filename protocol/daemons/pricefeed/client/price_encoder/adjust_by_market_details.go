package price_encoder

import "github.com/dydxprotocol/v4-chain/protocol/daemons/pricefeed/client/types"

// adjustByMarketDetails contains all information required to find and interpret the adjust-by market's price
// for the purposes of converting an exchange's raw API response price into a market price.
type adjustByMarketDetails struct {
	MarketId     types.MarketId
	Exponent     types.Exponent
	MinExchanges uint32
}
