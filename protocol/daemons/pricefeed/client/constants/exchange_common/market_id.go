package exchange_common

import "github.com/dydxprotocol/v4/daemons/pricefeed/client/types"

// All ids must match with the prices module.
// TODO(DEC-663): replace all values with their real value from the prices module.
const (
	// MARKET_BTC_USD is the id for the BTC-USD market pair.
	MARKET_BTC_USD types.MarketId = 0
	// MARKET_ETH_USD is the id for the ETH-USD market pair.
	MARKET_ETH_USD types.MarketId = 1
	// MARKET_LINK_USD is the id for the LINK-USD market pair.
	MARKET_LINK_USD types.MarketId = 2
)
