package constants

import (
	"github.com/dydxprotocol/v4/daemons/pricefeed/client/constants/exchange_common"
	"github.com/dydxprotocol/v4/daemons/pricefeed/client/types"
)

const (
	defaultIntervalMs            = 2_000
	defaultTimeoutMs             = 3_000
	defaultMaxQueries            = 3
	defaultMultiMarketMaxQueries = 1
	bitfinexNumSupportedMarkets  = 2
)

var (
	StaticExchangeStartupConfig = map[types.ExchangeFeedId]*types.ExchangeStartupConfig{
		exchange_common.EXCHANGE_FEED_BINANCE: {
			ExchangeFeedId: exchange_common.EXCHANGE_FEED_BINANCE,
			IntervalMs:     defaultIntervalMs,
			TimeoutMs:      defaultTimeoutMs,
			MaxQueries:     defaultMaxQueries,
		},
		exchange_common.EXCHANGE_FEED_BINANCE_US: {
			ExchangeFeedId: exchange_common.EXCHANGE_FEED_BINANCE_US,
			IntervalMs:     defaultIntervalMs,
			TimeoutMs:      defaultTimeoutMs,
			MaxQueries:     defaultMaxQueries,
		},
		exchange_common.EXCHANGE_FEED_BITFINEX: {
			ExchangeFeedId: exchange_common.EXCHANGE_FEED_BITFINEX,
			IntervalMs:     defaultIntervalMs,
			TimeoutMs:      defaultTimeoutMs,
			MaxQueries:     bitfinexNumSupportedMarkets,
		},
		exchange_common.EXCHANGE_FEED_KRAKEN: {
			ExchangeFeedId: exchange_common.EXCHANGE_FEED_KRAKEN,
			IntervalMs:     defaultIntervalMs,
			TimeoutMs:      defaultTimeoutMs,
			MaxQueries:     defaultMultiMarketMaxQueries,
		},
	}
)
