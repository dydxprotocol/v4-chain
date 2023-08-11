package metrics

import (
	"github.com/dydxprotocol/v4/daemons/pricefeed/client/constants/exchange_common"
	"github.com/dydxprotocol/v4/daemons/pricefeed/client/types"
)

var (
	// TODO(DEC-586): get information from protocol.
	// StaticMarketSymbols is the mapping of marketIds to their human readable symbols.
	StaticMarketSymbols = map[types.MarketId]string{
		exchange_common.MARKET_BTC_USD:  "BTCUSD",
		exchange_common.MARKET_ETH_USD:  "ETHUSD",
		exchange_common.MARKET_LINK_USD: "LINKUSD",
	}
)
