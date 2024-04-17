package constants

import (
	"github.com/StreamFinance-Protocol/stream-chain/protocol/daemons/pricefeed/client/constants/exchange_common"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/daemons/pricefeed/client/price_function/binance"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/daemons/pricefeed/client/price_function/bitfinex"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/daemons/pricefeed/client/price_function/bitstamp"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/daemons/pricefeed/client/price_function/bybit"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/daemons/pricefeed/client/price_function/coinbase_pro"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/daemons/pricefeed/client/price_function/crypto_com"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/daemons/pricefeed/client/price_function/gate"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/daemons/pricefeed/client/price_function/huobi"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/daemons/pricefeed/client/price_function/kraken"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/daemons/pricefeed/client/price_function/kucoin"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/daemons/pricefeed/client/price_function/mexc"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/daemons/pricefeed/client/price_function/okx"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/daemons/pricefeed/client/price_function/test_fixed_price_exchange"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/daemons/pricefeed/client/price_function/test_volatile_exchange"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/daemons/pricefeed/client/price_function/testexchange"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/daemons/pricefeed/client/types"
)

var (
	// StaticExchangeDetails is the static mapping of `ExchangeId` to its `ExchangeQueryDetails`.
	StaticExchangeDetails = map[types.ExchangeId]types.ExchangeQueryDetails{
		exchange_common.EXCHANGE_ID_BINANCE:                   binance.BinanceDetails,
		exchange_common.EXCHANGE_ID_BINANCE_US:                binance.BinanceUSDetails,
		exchange_common.EXCHANGE_ID_BITFINEX:                  bitfinex.BitfinexDetails,
		exchange_common.EXCHANGE_ID_KRAKEN:                    kraken.KrakenDetails,
		exchange_common.EXCHANGE_ID_GATE:                      gate.GateDetails,
		exchange_common.EXCHANGE_ID_BITSTAMP:                  bitstamp.BitstampDetails,
		exchange_common.EXCHANGE_ID_BYBIT:                     bybit.BybitDetails,
		exchange_common.EXCHANGE_ID_CRYPTO_COM:                crypto_com.CryptoComDetails,
		exchange_common.EXCHANGE_ID_HUOBI:                     huobi.HuobiDetails,
		exchange_common.EXCHANGE_ID_KUCOIN:                    kucoin.KucoinDetails,
		exchange_common.EXCHANGE_ID_OKX:                       okx.OkxDetails,
		exchange_common.EXCHANGE_ID_MEXC:                      mexc.MexcDetails,
		exchange_common.EXCHANGE_ID_COINBASE_PRO:              coinbase_pro.CoinbaseProDetails,
		exchange_common.EXCHANGE_ID_TEST_EXCHANGE:             testexchange.TestExchangeDetails,
		exchange_common.EXCHANGE_ID_TEST_VOLATILE_EXCHANGE:    test_volatile_exchange.TestVolatileExchangeDetails,
		exchange_common.EXCHANGE_ID_TEST_FIXED_PRICE_EXCHANGE: test_fixed_price_exchange.TestFixedPriceExchangeDetails,
	}
)
