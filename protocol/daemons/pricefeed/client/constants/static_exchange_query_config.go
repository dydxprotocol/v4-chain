package constants

import (
	"github.com/dydxprotocol/v4-chain/protocol/daemons/pricefeed/client/constants/exchange_common"
	"github.com/dydxprotocol/v4-chain/protocol/daemons/pricefeed/client/types"
)

const (
	defaultIntervalMs = 2_000
	// Binance / BinanceUS has a limit of 1_200 request weight per minute. At 40 request weight per
	// iteration, we can query once every 2 seconds, but increase to 2.5s to allow for some jitter.
	binanceIntervalMs = 2_500
	// Bitfinex has a rate limit of 30 requests per minute, so we query every 2.5 seconds
	// to allow for some jitter, as the 2-second interval occasionally produces 429 responses.
	bitfinexIntervalMs           = 2_500
	defaultTimeoutMs             = 3_000
	defaultMaxQueries            = 3
	defaultMultiMarketMaxQueries = 1
)

var (
	StaticExchangeQueryConfig = map[types.ExchangeId]*types.ExchangeQueryConfig{
		// See above for rate limiting information of Binance.
		// https://binance-docs.github.io/apidocs/spot/en/#24hr-ticker-price-change-statistics
		exchange_common.EXCHANGE_ID_BINANCE: {
			ExchangeId: exchange_common.EXCHANGE_ID_BINANCE,
			IntervalMs: binanceIntervalMs,
			TimeoutMs:  defaultTimeoutMs,
			MaxQueries: defaultMultiMarketMaxQueries,
		},
		// See above for rate limiting information of BinanceUS.
		exchange_common.EXCHANGE_ID_BINANCE_US: {
			ExchangeId: exchange_common.EXCHANGE_ID_BINANCE_US,
			IntervalMs: binanceIntervalMs,
			TimeoutMs:  defaultTimeoutMs,
			MaxQueries: defaultMultiMarketMaxQueries,
		},
		// Bitfinex has a limit of 30 requests per minute.
		// https://docs.bitfinex.com/reference/rest-public-tickers
		exchange_common.EXCHANGE_ID_BITFINEX: {
			ExchangeId: exchange_common.EXCHANGE_ID_BITFINEX,
			IntervalMs: bitfinexIntervalMs,
			TimeoutMs:  defaultTimeoutMs,
			MaxQueries: defaultMultiMarketMaxQueries,
		},
		exchange_common.EXCHANGE_ID_KRAKEN: {
			ExchangeId: exchange_common.EXCHANGE_ID_KRAKEN,
			IntervalMs: defaultIntervalMs,
			TimeoutMs:  defaultTimeoutMs,
			MaxQueries: defaultMultiMarketMaxQueries,
		},
		// Gate has a limit of 900 requests/second
		// https://www.gate.io/docs/developers/apiv4/en/#frequency-limit-rule
		exchange_common.EXCHANGE_ID_GATE: {
			ExchangeId: exchange_common.EXCHANGE_ID_GATE,
			IntervalMs: defaultIntervalMs,
			TimeoutMs:  defaultTimeoutMs,
			MaxQueries: defaultMultiMarketMaxQueries,
		},
		// Bitstamp has a limit of 8000 requests per 10 minutes.
		// https://www.bitstamp.net/api/#request-limits
		exchange_common.EXCHANGE_ID_BITSTAMP: {
			ExchangeId: exchange_common.EXCHANGE_ID_BITSTAMP,
			IntervalMs: defaultIntervalMs,
			TimeoutMs:  defaultTimeoutMs,
			MaxQueries: defaultMultiMarketMaxQueries,
		},
		// Bybit has a limit of 120 requests per second for 5 consecutive seconds.
		// https://bybit-exchange.github.io/docs/v5/rate-limit
		exchange_common.EXCHANGE_ID_BYBIT: {
			ExchangeId: exchange_common.EXCHANGE_ID_BYBIT,
			IntervalMs: defaultIntervalMs,
			TimeoutMs:  defaultTimeoutMs,
			MaxQueries: defaultMultiMarketMaxQueries,
		},
		// Crypto.com has a limit of 100 requests per second.
		// https://exchange-docs.crypto.com/derivatives/index.html#rate-limits
		exchange_common.EXCHANGE_ID_CRYPTO_COM: {
			ExchangeId: exchange_common.EXCHANGE_ID_CRYPTO_COM,
			IntervalMs: defaultIntervalMs,
			TimeoutMs:  defaultTimeoutMs,
			MaxQueries: defaultMultiMarketMaxQueries,
		},
		// Huobi has a limit of 100 requests per second.
		// https://huobiapi.github.io/docs/spot/v1/en/#api-access
		exchange_common.EXCHANGE_ID_HUOBI: {
			ExchangeId: exchange_common.EXCHANGE_ID_HUOBI,
			IntervalMs: defaultIntervalMs,
			TimeoutMs:  defaultTimeoutMs,
			MaxQueries: defaultMultiMarketMaxQueries,
		},
		// Kucoin has a limit of 500 requests per 10 seconds.
		// https://docs.kucoin.com/#request-rate-limit
		exchange_common.EXCHANGE_ID_KUCOIN: {
			ExchangeId: exchange_common.EXCHANGE_ID_KUCOIN,
			IntervalMs: defaultIntervalMs,
			TimeoutMs:  defaultTimeoutMs,
			MaxQueries: defaultMultiMarketMaxQueries,
		},
		// Okx has a limit of 20 requests per 2 seconds.
		// https://www.okx.com/docs-v5/en/#rest-api-market-data-get-tickers
		exchange_common.EXCHANGE_ID_OKX: {
			ExchangeId: exchange_common.EXCHANGE_ID_OKX,
			IntervalMs: defaultIntervalMs,
			TimeoutMs:  defaultTimeoutMs,
			MaxQueries: defaultMultiMarketMaxQueries,
		},
		// Mexc has a limit of 20 requests per second.
		// https://mxcdevelop.github.io/apidocs/spot_v2_en/#rate-limit
		exchange_common.EXCHANGE_ID_MEXC: {
			ExchangeId: exchange_common.EXCHANGE_ID_MEXC,
			IntervalMs: defaultIntervalMs,
			TimeoutMs:  defaultTimeoutMs,
			MaxQueries: defaultMultiMarketMaxQueries,
		},
		// CoinbasePro has a limit of 10 requests per second.
		// https://docs.cloud.coinbase.com/exchange/docs/rest-rate-limits
		exchange_common.EXCHANGE_ID_COINBASE_PRO: {
			ExchangeId: exchange_common.EXCHANGE_ID_COINBASE_PRO,
			IntervalMs: defaultIntervalMs,
			TimeoutMs:  defaultTimeoutMs,
			MaxQueries: defaultMaxQueries,
		},
		exchange_common.EXCHANGE_ID_TEST_VOLATILE_EXCHANGE: {
			ExchangeId: exchange_common.EXCHANGE_ID_TEST_VOLATILE_EXCHANGE,
			IntervalMs: defaultIntervalMs,
			TimeoutMs:  defaultTimeoutMs,
			MaxQueries: defaultMaxQueries,
		},
		exchange_common.EXCHANGE_ID_TEST_FIXED_PRICE_EXCHANGE: {
			ExchangeId: exchange_common.EXCHANGE_ID_TEST_FIXED_PRICE_EXCHANGE,
			IntervalMs: defaultIntervalMs,
			TimeoutMs:  defaultTimeoutMs,
			MaxQueries: defaultMaxQueries,
		},
	}
)
