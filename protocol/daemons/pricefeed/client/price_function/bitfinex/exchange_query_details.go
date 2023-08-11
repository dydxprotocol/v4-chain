package bitfinex

import (
	"github.com/dydxprotocol/v4/daemons/pricefeed/client/constants/exchange_common"
	"github.com/dydxprotocol/v4/daemons/pricefeed/client/types"
)

var (
	BitfinexDetails = types.ExchangeQueryDetails{
		Exchange:      exchange_common.EXCHANGE_ID_BITFINEX,
		Url:           "https://api-pub.bitfinex.com/v2/tickers?symbols=$",
		PriceFunction: BitfinexPriceFunction,
		IsMultiMarket: true,
	}
)
