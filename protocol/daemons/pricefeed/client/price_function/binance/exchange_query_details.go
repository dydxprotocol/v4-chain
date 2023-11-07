package binance

import (
	"github.com/dydxprotocol/v4-chain/protocol/daemons/pricefeed/client/constants/exchange_common"
	"github.com/dydxprotocol/v4-chain/protocol/daemons/pricefeed/client/types"
)

var (
	BinanceDetails = types.ExchangeQueryDetails{
		Exchange:      exchange_common.EXCHANGE_ID_BINANCE,
		Url:           "https://data-api.binance.vision/api/v3/ticker/24hr",
		PriceFunction: BinancePriceFunction,
		IsMultiMarket: true,
	}

	BinanceUSDetails = types.ExchangeQueryDetails{
		Exchange:      exchange_common.EXCHANGE_ID_BINANCE_US,
		Url:           "https://api.binance.us/api/v3/ticker/24hr",
		PriceFunction: BinancePriceFunction,
		IsMultiMarket: true,
	}
)
