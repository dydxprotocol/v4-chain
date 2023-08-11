package kraken

import (
	"github.com/dydxprotocol/v4/daemons/pricefeed/client/constants/exchange_common"
	"github.com/dydxprotocol/v4/daemons/pricefeed/client/types"
)

var (
	KrakenDetails = types.ExchangeQueryDetails{
		Exchange:      exchange_common.EXCHANGE_ID_KRAKEN,
		Url:           "https://api.kraken.com/0/public/Ticker?pair=$",
		PriceFunction: KrakenPriceFunction,
		IsMultiMarket: true,
	}
)
