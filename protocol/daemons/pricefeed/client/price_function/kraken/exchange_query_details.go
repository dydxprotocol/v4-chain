package kraken

import (
	"github.com/dydxprotocol/v4/daemons/pricefeed/client/constants/exchange_common"
	"github.com/dydxprotocol/v4/daemons/pricefeed/client/types"
)

var (
	KrakenDetails = types.ExchangeQueryDetails{
		Exchange: exchange_common.EXCHANGE_FEED_KRAKEN,
		Url:      "https://api.kraken.com/0/public/Ticker?pair=$",
		MarketSymbols: map[types.MarketId]string{
			exchange_common.MARKET_BTC_USD:  "XXBTZUSD",
			exchange_common.MARKET_ETH_USD:  "XETHZUSD",
			exchange_common.MARKET_LINK_USD: "LINKUSD",
		},
		PriceFunction: KrakenPriceFunction,
		IsMultiMarket: true,
	}
)
