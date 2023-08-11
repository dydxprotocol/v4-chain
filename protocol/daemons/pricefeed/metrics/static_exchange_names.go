package metrics

import (
	"github.com/dydxprotocol/v4/daemons/pricefeed/client/constants/exchange_common"
	"github.com/dydxprotocol/v4/daemons/pricefeed/client/types"
)

var (
	// TODO(DEC-586): get information from protocol.
	// StaticExchangeNames is the mapping of exchangeFeedIds to their names.
	StaticExchangeNames = map[types.ExchangeFeedId]string{
		exchange_common.EXCHANGE_FEED_BINANCE:    exchange_common.EXCHANGE_NAME_BINANCE,
		exchange_common.EXCHANGE_FEED_BINANCE_US: exchange_common.EXCHANGE_NAME_BINANCEUS,
		exchange_common.EXCHANGE_FEED_BITFINEX:   exchange_common.EXCHANGE_NAME_BITFINEX,
		exchange_common.EXCHANGE_FEED_KRAKEN:     exchange_common.EXCHANGE_NAME_KRAKEN,
	}
)
