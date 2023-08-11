package constants

import (
	"github.com/dydxprotocol/v4/daemons/pricefeed/client/constants/exchange_common"
	"github.com/dydxprotocol/v4/daemons/pricefeed/client/price_function/binance"
	"github.com/dydxprotocol/v4/daemons/pricefeed/client/price_function/bitfinex"
	"github.com/dydxprotocol/v4/daemons/pricefeed/client/price_function/kraken"
	"github.com/dydxprotocol/v4/daemons/pricefeed/client/types"
)

var (
	// StaticExchangeDetails is the static mapping of `ExchangeId` to its `ExchangeQueryDetails`.
	StaticExchangeDetails = map[types.ExchangeFeedId]types.ExchangeQueryDetails{
		exchange_common.EXCHANGE_FEED_BINANCE:    binance.BinanceDetails,
		exchange_common.EXCHANGE_FEED_BINANCE_US: binance.BinanceUSDetails,
		exchange_common.EXCHANGE_FEED_BITFINEX:   bitfinex.BitfinexDetails,
		exchange_common.EXCHANGE_FEED_KRAKEN:     kraken.KrakenDetails,
	}
)
