package okx

import (
	"github.com/dydxprotocol/v4-chain/protocol/daemons/pricefeed/client/constants/exchange_common"
	"github.com/dydxprotocol/v4-chain/protocol/daemons/pricefeed/client/types"
)

var (
	OkxDetails = types.ExchangeQueryDetails{
		Exchange:      exchange_common.EXCHANGE_ID_OKX,
		Url:           "https://www.okx.com/api/v5/market/tickers?instType=SPOT",
		PriceFunction: OkxPriceFunction,
		IsMultiMarket: true,
	}
)
