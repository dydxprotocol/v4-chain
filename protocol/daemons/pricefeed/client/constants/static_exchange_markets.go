package constants

import (
	"github.com/dydxprotocol/v4/daemons/pricefeed/client/constants/exchange_common"
	"github.com/dydxprotocol/v4/daemons/pricefeed/client/types"
)

var (
	// StaticExchangeMarkets maps exchanges to the markets they will be queried for.
	// NOTE: this mapping must be modified for the price-daemon client to query different
	// markets per exchange.
	StaticExchangeMarkets = map[types.ExchangeFeedId][]types.MarketId{
		exchange_common.EXCHANGE_FEED_BINANCE: {
			exchange_common.MARKET_BTC_USD,
			exchange_common.MARKET_ETH_USD,
			exchange_common.MARKET_LINK_USD,
		},
		exchange_common.EXCHANGE_FEED_BINANCE_US: {
			exchange_common.MARKET_BTC_USD,
			exchange_common.MARKET_ETH_USD,
			exchange_common.MARKET_LINK_USD,
		},
		exchange_common.EXCHANGE_FEED_BITFINEX: {
			exchange_common.MARKET_BTC_USD,
			exchange_common.MARKET_ETH_USD,
		},
		exchange_common.EXCHANGE_FEED_KRAKEN: {
			exchange_common.MARKET_BTC_USD,
			exchange_common.MARKET_ETH_USD,
			exchange_common.MARKET_LINK_USD,
		},
	}
)
