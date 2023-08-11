package constants

import (
	"github.com/dydxprotocol/v4/daemons/pricefeed/client/constants/exchange_common"
	"github.com/dydxprotocol/v4/daemons/pricefeed/client/types"
)

var (
	// StaticMarketExponent is the static mapping of `MarketId` to its price exponent.
	// TODO(DEC-663): replace all values with their real value from the prices module.
	StaticMarketPriceExponent = map[types.MarketId]types.Exponent{
		exchange_common.MARKET_BTC_USD:  -5,
		exchange_common.MARKET_ETH_USD:  -6,
		exchange_common.MARKET_LINK_USD: -8,
	}
)
