package constants_test

import (
	"testing"

	"github.com/dydxprotocol/v4-chain/protocol/daemons/pricefeed/client/constants"
	"github.com/dydxprotocol/v4-chain/protocol/daemons/pricefeed/client/constants/exchange_common"
	"github.com/dydxprotocol/v4-chain/protocol/daemons/pricefeed/client/price_function/binance"
	"github.com/dydxprotocol/v4-chain/protocol/daemons/pricefeed/client/price_function/bitfinex"
	"github.com/dydxprotocol/v4-chain/protocol/daemons/pricefeed/client/price_function/bitstamp"
	"github.com/dydxprotocol/v4-chain/protocol/daemons/pricefeed/client/price_function/bybit"
	"github.com/dydxprotocol/v4-chain/protocol/daemons/pricefeed/client/price_function/coinbase_pro"
	"github.com/dydxprotocol/v4-chain/protocol/daemons/pricefeed/client/price_function/crypto_com"
	"github.com/dydxprotocol/v4-chain/protocol/daemons/pricefeed/client/price_function/gate"
	"github.com/dydxprotocol/v4-chain/protocol/daemons/pricefeed/client/price_function/huobi"
	"github.com/dydxprotocol/v4-chain/protocol/daemons/pricefeed/client/price_function/kraken"
	"github.com/dydxprotocol/v4-chain/protocol/daemons/pricefeed/client/price_function/kucoin"
	"github.com/dydxprotocol/v4-chain/protocol/daemons/pricefeed/client/price_function/mexc"
	"github.com/dydxprotocol/v4-chain/protocol/daemons/pricefeed/client/price_function/okx"
	"github.com/dydxprotocol/v4-chain/protocol/daemons/pricefeed/client/price_function/testexchange"
	"github.com/dydxprotocol/v4-chain/protocol/daemons/pricefeed/client/types"
	"github.com/stretchr/testify/require"
)

func TestStaticExchangeDetailsCache(t *testing.T) {
	tests := map[string]struct {
		// parameters
		exchangeId types.ExchangeId

		// expectations
		expectedValue types.ExchangeQueryDetails
		expectedFound bool
	}{
		"Get BINANCE exchangeDetails": {
			exchangeId:    exchange_common.EXCHANGE_ID_BINANCE,
			expectedValue: binance.BinanceDetails,
			expectedFound: true,
		},
		"Get BINANCEUS exchangeDetails": {
			exchangeId:    exchange_common.EXCHANGE_ID_BINANCE_US,
			expectedValue: binance.BinanceUSDetails,
			expectedFound: true,
		},
		"Get Bitfinex exchangeDetails": {
			exchangeId:    exchange_common.EXCHANGE_ID_BITFINEX,
			expectedValue: bitfinex.BitfinexDetails,
			expectedFound: true,
		},
		"Get Kraken exchangeDetails": {
			exchangeId:    exchange_common.EXCHANGE_ID_KRAKEN,
			expectedValue: kraken.KrakenDetails,
			expectedFound: true,
		},
		"Get Gate exchangeDetails": {
			exchangeId:    exchange_common.EXCHANGE_ID_GATE,
			expectedValue: gate.GateDetails,
			expectedFound: true,
		},
		"Get Bitstamp exchangeDetails": {
			exchangeId:    exchange_common.EXCHANGE_ID_BITSTAMP,
			expectedValue: bitstamp.BitstampDetails,
			expectedFound: true,
		},
		"Get Bybit exchangeDetails": {
			exchangeId:    exchange_common.EXCHANGE_ID_BYBIT,
			expectedValue: bybit.BybitDetails,
			expectedFound: true,
		},
		"Get CryptoCom exchangeDetails": {
			exchangeId:    exchange_common.EXCHANGE_ID_CRYPTO_COM,
			expectedValue: crypto_com.CryptoComDetails,
			expectedFound: true,
		},
		"Get Huobi exchangeDetails": {
			exchangeId:    exchange_common.EXCHANGE_ID_HUOBI,
			expectedValue: huobi.HuobiDetails,
			expectedFound: true,
		},
		"Get Kucoin exchangeDetails": {
			exchangeId:    exchange_common.EXCHANGE_ID_KUCOIN,
			expectedValue: kucoin.KucoinDetails,
			expectedFound: true,
		},
		"Get Okx exchangeDetails": {
			exchangeId:    exchange_common.EXCHANGE_ID_OKX,
			expectedValue: okx.OkxDetails,
			expectedFound: true,
		},
		"Get Mexc exchangeDetails": {
			exchangeId:    exchange_common.EXCHANGE_ID_MEXC,
			expectedValue: mexc.MexcDetails,
			expectedFound: true,
		},
		"Get CoinbasePro exchangeDetails": {
			exchangeId:    exchange_common.EXCHANGE_ID_COINBASE_PRO,
			expectedValue: coinbase_pro.CoinbaseProDetails,
			expectedFound: true,
		},
		"Get test exchange exchangeDetails": {
			exchangeId:    exchange_common.EXCHANGE_ID_TEST_EXCHANGE,
			expectedValue: testexchange.TestExchangeDetails,
			expectedFound: true,
		},
		"Get unknown exchangeDetails": {
			exchangeId:    "unknown",
			expectedFound: false,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			value, ok := constants.StaticExchangeDetails[tc.exchangeId]
			require.Equal(t, tc.expectedValue.Exchange, value.Exchange)
			require.Equal(t, tc.expectedValue.Url, value.Url)
			require.Equal(t, tc.expectedFound, ok)
		})
	}
}
