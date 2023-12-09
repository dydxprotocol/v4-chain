package constants_test

import (
	"testing"

	"github.com/dydxprotocol/v4-chain/protocol/daemons/pricefeed/client/constants"
	"github.com/dydxprotocol/v4-chain/protocol/daemons/pricefeed/client/constants/exchange_common"
	"github.com/dydxprotocol/v4-chain/protocol/daemons/pricefeed/client/types"
	"github.com/stretchr/testify/require"
)

func TestStaticExchangeQueryConfigCache(t *testing.T) {
	tests := map[string]struct {
		// parameters
		exchangeId types.ExchangeId

		// expectations
		expectedValue *types.ExchangeQueryConfig
		expectedFound bool
	}{
		"Get BINANCE exchangeDetails": {
			exchangeId: exchange_common.EXCHANGE_ID_BINANCE,
			expectedValue: &types.ExchangeQueryConfig{
				ExchangeId: exchange_common.EXCHANGE_ID_BINANCE,
				IntervalMs: 2_500,
				TimeoutMs:  3_000,
				MaxQueries: 1,
			},
			expectedFound: true,
		},
		"Get BINANCEUS exchangeDetails": {
			exchangeId: exchange_common.EXCHANGE_ID_BINANCE_US,
			expectedValue: &types.ExchangeQueryConfig{
				ExchangeId: exchange_common.EXCHANGE_ID_BINANCE_US,
				IntervalMs: 2_500,
				TimeoutMs:  3_000,
				MaxQueries: 1,
			},
			expectedFound: true,
		},
		"Get BITFINEX exchangeDetails": {
			exchangeId: exchange_common.EXCHANGE_ID_BITFINEX,
			expectedValue: &types.ExchangeQueryConfig{
				ExchangeId: exchange_common.EXCHANGE_ID_BITFINEX,
				IntervalMs: 2_500,
				TimeoutMs:  3_000,
				MaxQueries: 1,
			},
			expectedFound: true,
		},
		"Get Kraken exchangeDetails": {
			exchangeId: exchange_common.EXCHANGE_ID_KRAKEN,
			expectedValue: &types.ExchangeQueryConfig{
				ExchangeId: exchange_common.EXCHANGE_ID_KRAKEN,
				IntervalMs: 2_000,
				TimeoutMs:  3_000,
				MaxQueries: 1,
			},
			expectedFound: true,
		},
		"Get GATE exchangeDetails": {
			exchangeId: exchange_common.EXCHANGE_ID_GATE,
			expectedValue: &types.ExchangeQueryConfig{
				ExchangeId: exchange_common.EXCHANGE_ID_GATE,
				IntervalMs: 2_000,
				TimeoutMs:  3_000,
				MaxQueries: 1,
			},
			expectedFound: true,
		},
		"Get Bitstamp exchangeDetails": {
			exchangeId: exchange_common.EXCHANGE_ID_BITSTAMP,
			expectedValue: &types.ExchangeQueryConfig{
				ExchangeId: exchange_common.EXCHANGE_ID_BITSTAMP,
				IntervalMs: 2_000,
				TimeoutMs:  3_000,
				MaxQueries: 1,
			},
			expectedFound: true,
		},
		"Get Bybit exchangeDetails": {
			exchangeId: exchange_common.EXCHANGE_ID_BYBIT,
			expectedValue: &types.ExchangeQueryConfig{
				ExchangeId: exchange_common.EXCHANGE_ID_BYBIT,
				IntervalMs: 2_000,
				TimeoutMs:  3_000,
				MaxQueries: 1,
			},
			expectedFound: true,
		},
		"Get CryptoCom exchangeDetails": {
			exchangeId: exchange_common.EXCHANGE_ID_CRYPTO_COM,
			expectedValue: &types.ExchangeQueryConfig{
				ExchangeId: exchange_common.EXCHANGE_ID_CRYPTO_COM,
				IntervalMs: 2_000,
				TimeoutMs:  3_000,
				MaxQueries: 1,
			},
			expectedFound: true,
		},
		"Get Huobi exchangeDetails": {
			exchangeId: exchange_common.EXCHANGE_ID_HUOBI,
			expectedValue: &types.ExchangeQueryConfig{
				ExchangeId: exchange_common.EXCHANGE_ID_HUOBI,
				IntervalMs: 2_000,
				TimeoutMs:  3_000,
				MaxQueries: 1,
			},
			expectedFound: true,
		},
		"Get Kucoin exchangeDetails": {
			exchangeId: exchange_common.EXCHANGE_ID_KUCOIN,
			expectedValue: &types.ExchangeQueryConfig{
				ExchangeId: exchange_common.EXCHANGE_ID_KUCOIN,
				IntervalMs: 2_000,
				TimeoutMs:  3_000,
				MaxQueries: 1,
			},
			expectedFound: true,
		},
		"Get Okx exchangeDetails": {
			exchangeId: exchange_common.EXCHANGE_ID_OKX,
			expectedValue: &types.ExchangeQueryConfig{
				ExchangeId: exchange_common.EXCHANGE_ID_OKX,
				IntervalMs: 2_000,
				TimeoutMs:  3_000,
				MaxQueries: 1,
			},
			expectedFound: true,
		},
		"Get Mexc exchangeDetails": {
			exchangeId: exchange_common.EXCHANGE_ID_MEXC,
			expectedValue: &types.ExchangeQueryConfig{
				ExchangeId: exchange_common.EXCHANGE_ID_MEXC,
				IntervalMs: 2_000,
				TimeoutMs:  3_000,
				MaxQueries: 1,
			},
			expectedFound: true,
		},
		"Get CoinbasePro exchangeDetails": {
			exchangeId: exchange_common.EXCHANGE_ID_COINBASE_PRO,
			expectedValue: &types.ExchangeQueryConfig{
				ExchangeId: exchange_common.EXCHANGE_ID_COINBASE_PRO,
				IntervalMs: 2_000,
				TimeoutMs:  3_000,
				MaxQueries: 3,
			},
			expectedFound: true,
		},
		"Get unknown exchangeDetails": {
			exchangeId:    "unknown",
			expectedFound: false,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			value, ok := constants.StaticExchangeQueryConfig[tc.exchangeId]
			require.Equal(t, tc.expectedValue, value)
			require.Equal(t, ok, tc.expectedFound)
		})
	}
}

func TestStaticExchangeQueryConfigCacheLength(t *testing.T) {
	require.Len(t, constants.StaticExchangeQueryConfig, 15)
}
