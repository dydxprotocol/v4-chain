package exchange_common

import "github.com/dydxprotocol/v4/daemons/pricefeed/client/types"

// All ids must match with the prices module.
// TODO(DEC-663): replace all values with their real value from the prices module.
const (
	// EXCHANGE_FEED_BINANCE is the id for Binance exchange.
	EXCHANGE_FEED_BINANCE types.ExchangeFeedId = 0
	// EXCHANGE_FEED_BINANCE_US is the id for BinanceUS exchange.
	EXCHANGE_FEED_BINANCE_US types.ExchangeFeedId = 1
	// EXCHANGE_FEED_BITFINEX is the id for Bitfinex exchange.
	EXCHANGE_FEED_BITFINEX types.ExchangeFeedId = 2
	// EXCHANGE_FEED_KRAKEN is the id for Kraken exchange
	EXCHANGE_FEED_KRAKEN types.ExchangeFeedId = 3
)
