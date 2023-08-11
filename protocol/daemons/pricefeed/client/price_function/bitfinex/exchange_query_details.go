package bitfinex

import (
	"github.com/dydxprotocol/v4/daemons/pricefeed/client/constants/exchange_common"
	"github.com/dydxprotocol/v4/daemons/pricefeed/client/types"
)

var (
	BitfinexDetails = types.ExchangeQueryDetails{
		Exchange: exchange_common.EXCHANGE_FEED_BITFINEX,
		Url:      "https://api.bitfinex.com/v2/ticker/$",
		MarketSymbols: map[uint32]string{
			exchange_common.MARKET_BTC_USD: "tBTCUSD",
			exchange_common.MARKET_ETH_USD: "tETHUSD",
		},
		PriceFunction: BitfinexPriceFunction,
	}
)
