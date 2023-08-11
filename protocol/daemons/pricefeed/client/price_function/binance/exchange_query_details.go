package binance

import (
	"github.com/dydxprotocol/v4/daemons/pricefeed/client/constants/exchange_common"
	"github.com/dydxprotocol/v4/daemons/pricefeed/client/types"
)

var (
	BinanceDetails = types.ExchangeQueryDetails{
		Exchange: exchange_common.EXCHANGE_FEED_BINANCE,
		Url:      "https://data.binance.com/api/v3/ticker/24hr?symbol=$",
		MarketSymbols: map[types.MarketId]string{
			exchange_common.MARKET_BTC_USD:  "BTCUSDT",
			exchange_common.MARKET_ETH_USD:  "ETHUSDT",
			exchange_common.MARKET_LINK_USD: "LINKUSDT",
		},
		PriceFunction: BinancePriceFunction,
	}

	BinanceUSDetails = types.ExchangeQueryDetails{
		Exchange: exchange_common.EXCHANGE_FEED_BINANCE_US,
		Url:      "https://api.binance.us/api/v3/ticker/24hr?symbol=$",
		MarketSymbols: map[types.MarketId]string{
			exchange_common.MARKET_BTC_USD:  "BTCUSD",
			exchange_common.MARKET_ETH_USD:  "ETHUSD",
			exchange_common.MARKET_LINK_USD: "LINKUSD",
		},
		PriceFunction: BinancePriceFunction,
	}
)
