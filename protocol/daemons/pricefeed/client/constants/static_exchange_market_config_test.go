package constants

import (
	"testing"

	"github.com/dydxprotocol/v4-chain/protocol/daemons/pricefeed/client/constants/exchange_common"
	"github.com/dydxprotocol/v4-chain/protocol/daemons/pricefeed/client/types"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/daemons/pricefeed"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/json"
	"github.com/stretchr/testify/require"
)

func TestStaticExchangeMarketConfigCache(t *testing.T) {
	tests := map[string]struct {
		id             types.ExchangeId
		marketToConfig map[types.MarketId]types.MarketConfig
		expectedFound  bool
	}{
		"Binance": {
			id: "Binance",
			marketToConfig: map[types.MarketId]types.MarketConfig{
				exchange_common.MARKET_BTC_USD: {
					Ticker:         `"BTCUSDT"`,
					AdjustByMarket: newMarketIdWithValue(exchange_common.MARKET_USDT_USD),
				},
				exchange_common.MARKET_ETH_USD: {
					Ticker:         `"ETHUSDT"`,
					AdjustByMarket: newMarketIdWithValue(exchange_common.MARKET_USDT_USD),
				},
				exchange_common.MARKET_LINK_USD: {
					Ticker:         `"LINKUSDT"`,
					AdjustByMarket: newMarketIdWithValue(exchange_common.MARKET_USDT_USD),
				},
				exchange_common.MARKET_MATIC_USD: {
					Ticker:         `"MATICUSDT"`,
					AdjustByMarket: newMarketIdWithValue(exchange_common.MARKET_USDT_USD),
				},
				exchange_common.MARKET_CRV_USD: {
					Ticker:         `"CRVUSDT"`,
					AdjustByMarket: newMarketIdWithValue(exchange_common.MARKET_USDT_USD),
				},
				exchange_common.MARKET_SOL_USD: {
					Ticker:         `"SOLUSDT"`,
					AdjustByMarket: newMarketIdWithValue(exchange_common.MARKET_USDT_USD),
				},
				exchange_common.MARKET_ADA_USD: {
					Ticker:         `"ADAUSDT"`,
					AdjustByMarket: newMarketIdWithValue(exchange_common.MARKET_USDT_USD),
				},
				exchange_common.MARKET_AVAX_USD: {
					Ticker:         `"AVAXUSDT"`,
					AdjustByMarket: newMarketIdWithValue(exchange_common.MARKET_USDT_USD),
				},
				exchange_common.MARKET_FIL_USD: {
					Ticker:         `"FILUSDT"`,
					AdjustByMarket: newMarketIdWithValue(exchange_common.MARKET_USDT_USD),
				},
				exchange_common.MARKET_AAVE_USD: {
					Ticker:         `"AAVEUSDT"`,
					AdjustByMarket: newMarketIdWithValue(exchange_common.MARKET_USDT_USD),
				},
				exchange_common.MARKET_LTC_USD: {
					Ticker:         `"LTCUSDT"`,
					AdjustByMarket: newMarketIdWithValue(exchange_common.MARKET_USDT_USD),
				},
				exchange_common.MARKET_DOGE_USD: {
					Ticker:         `"DOGEUSDT"`,
					AdjustByMarket: newMarketIdWithValue(exchange_common.MARKET_USDT_USD),
				},
				exchange_common.MARKET_ICP_USD: {
					Ticker:         `"ICPUSDT"`,
					AdjustByMarket: newMarketIdWithValue(exchange_common.MARKET_USDT_USD),
				},
				exchange_common.MARKET_ATOM_USD: {
					Ticker:         `"ATOMUSDT"`,
					AdjustByMarket: newMarketIdWithValue(exchange_common.MARKET_USDT_USD),
				},
				exchange_common.MARKET_DOT_USD: {
					Ticker:         `"DOTUSDT"`,
					AdjustByMarket: newMarketIdWithValue(exchange_common.MARKET_USDT_USD),
				},
				exchange_common.MARKET_XTZ_USD: {
					Ticker:         `"XTZUSDT"`,
					AdjustByMarket: newMarketIdWithValue(exchange_common.MARKET_USDT_USD),
				},
				exchange_common.MARKET_UNI_USD: {
					Ticker:         `"UNIUSDT"`,
					AdjustByMarket: newMarketIdWithValue(exchange_common.MARKET_USDT_USD),
				},
				exchange_common.MARKET_BCH_USD: {
					Ticker:         `"BCHUSDT"`,
					AdjustByMarket: newMarketIdWithValue(exchange_common.MARKET_USDT_USD),
				},
				exchange_common.MARKET_EOS_USD: {
					Ticker:         `"EOSUSDT"`,
					AdjustByMarket: newMarketIdWithValue(exchange_common.MARKET_USDT_USD),
				},
				exchange_common.MARKET_TRX_USD: {
					Ticker:         `"TRXUSDT"`,
					AdjustByMarket: newMarketIdWithValue(exchange_common.MARKET_USDT_USD),
				},
				exchange_common.MARKET_ALGO_USD: {
					Ticker:         `"ALGOUSDT"`,
					AdjustByMarket: newMarketIdWithValue(exchange_common.MARKET_USDT_USD),
				},
				exchange_common.MARKET_NEAR_USD: {
					Ticker:         `"NEARUSDT"`,
					AdjustByMarket: newMarketIdWithValue(exchange_common.MARKET_USDT_USD),
				},
				exchange_common.MARKET_SNX_USD: {
					Ticker:         `"SNXUSDT"`,
					AdjustByMarket: newMarketIdWithValue(exchange_common.MARKET_USDT_USD),
				},
				exchange_common.MARKET_MKR_USD: {
					Ticker:         `"MKRUSDT"`,
					AdjustByMarket: newMarketIdWithValue(exchange_common.MARKET_USDT_USD),
				},
				exchange_common.MARKET_SUSHI_USD: {
					Ticker:         `"SUSHIUSDT"`,
					AdjustByMarket: newMarketIdWithValue(exchange_common.MARKET_USDT_USD),
				},
				exchange_common.MARKET_XLM_USD: {
					Ticker:         `"XLMUSDT"`,
					AdjustByMarket: newMarketIdWithValue(exchange_common.MARKET_USDT_USD),
				},
				exchange_common.MARKET_XMR_USD: {
					Ticker:         `"XMRUSDT"`,
					AdjustByMarket: newMarketIdWithValue(exchange_common.MARKET_USDT_USD),
				},
				exchange_common.MARKET_ETC_USD: {
					Ticker:         `"ETCUSDT"`,
					AdjustByMarket: newMarketIdWithValue(exchange_common.MARKET_USDT_USD),
				},
				exchange_common.MARKET_1INCH_USD: {
					Ticker:         `"1INCHUSDT"`,
					AdjustByMarket: newMarketIdWithValue(exchange_common.MARKET_USDT_USD),
				},
				exchange_common.MARKET_COMP_USD: {
					Ticker:         `"COMPUSDT"`,
					AdjustByMarket: newMarketIdWithValue(exchange_common.MARKET_USDT_USD),
				},
				exchange_common.MARKET_ZEC_USD: {
					Ticker:         `"ZECUSDT"`,
					AdjustByMarket: newMarketIdWithValue(exchange_common.MARKET_USDT_USD),
				},
				exchange_common.MARKET_ZRX_USD: {
					Ticker:         `"ZRXUSDT"`,
					AdjustByMarket: newMarketIdWithValue(exchange_common.MARKET_USDT_USD),
				},
				exchange_common.MARKET_YFI_USD: {
					Ticker:         `"YFIUSDT"`,
					AdjustByMarket: newMarketIdWithValue(exchange_common.MARKET_USDT_USD),
				},
				exchange_common.MARKET_USDT_USD: {
					Ticker:         `"BTCUSDT"`, // Adjusted with BTC index price.
					Invert:         true,
					AdjustByMarket: newMarketIdWithValue(exchange_common.MARKET_BTC_USD),
				},
			},
			expectedFound: true,
		},
		"BinanceUS": {
			id: "BinanceUS",
			marketToConfig: map[types.MarketId]types.MarketConfig{
				exchange_common.MARKET_BTC_USD: {
					Ticker: `"BTCUSD"`,
				},
				exchange_common.MARKET_ETH_USD: {
					Ticker: `"ETHUSD"`,
				},
				exchange_common.MARKET_LINK_USD: {
					Ticker: `"LINKUSD"`,
				},
				exchange_common.MARKET_MATIC_USD: {
					Ticker: `"MATICUSD"`,
				},
				exchange_common.MARKET_CRV_USD: {
					Ticker: `"CRVUSD"`,
				},
				exchange_common.MARKET_SOL_USD: {
					Ticker: `"SOLUSD"`,
				},
				exchange_common.MARKET_ADA_USD: {
					Ticker: `"ADAUSD"`,
				},
				exchange_common.MARKET_AVAX_USD: {
					Ticker: `"AVAXUSD"`,
				},
				exchange_common.MARKET_FIL_USD: {
					Ticker: `"FILUSD"`,
				},
				exchange_common.MARKET_AAVE_USD: {
					Ticker: `"AAVEUSD"`,
				},
				exchange_common.MARKET_LTC_USD: {
					Ticker: `"LTCUSD"`,
				},
				exchange_common.MARKET_DOGE_USD: {
					Ticker: `"DOGEUSD"`,
				},
				exchange_common.MARKET_ICP_USD: {
					Ticker: `"ICPUSD"`,
				},
				exchange_common.MARKET_ATOM_USD: {
					Ticker: `"ATOMUSD"`,
				},
				exchange_common.MARKET_DOT_USD: {
					Ticker: `"DOTUSD"`,
				},
				exchange_common.MARKET_XTZ_USD: {
					Ticker: `"XTZUSD"`,
				},
				exchange_common.MARKET_UNI_USD: {
					Ticker: `"UNIUSD"`,
				},
				exchange_common.MARKET_BCH_USD: {
					Ticker: `"BCHUSD"`,
				},
				exchange_common.MARKET_EOS_USD: {
					Ticker: `"EOSUSD"`,
				},
				exchange_common.MARKET_ALGO_USD: {
					Ticker: `"ALGOUSD"`,
				},
				exchange_common.MARKET_NEAR_USD: {
					Ticker: `"NEARUSD"`,
				},
				exchange_common.MARKET_SNX_USD: {
					Ticker: `"SNXUSD"`,
				},
				exchange_common.MARKET_MKR_USD: {
					Ticker: `"MKRUSD"`,
				},
				exchange_common.MARKET_SUSHI_USD: {
					Ticker: `"SUSHIUSD"`,
				},
				exchange_common.MARKET_XLM_USD: {
					Ticker: `"XLMUSD"`,
				},
				exchange_common.MARKET_ETC_USD: {
					Ticker: `"ETCUSD"`,
				},
				exchange_common.MARKET_1INCH_USD: {
					Ticker: `"1INCHUSD"`,
				},
				exchange_common.MARKET_COMP_USD: {
					Ticker: `"COMPUSD"`,
				},
				exchange_common.MARKET_ZEC_USD: {
					Ticker: `"ZECUSD"`,
				},
				exchange_common.MARKET_ZRX_USD: {
					Ticker: `"ZRXUSD"`,
				},
				exchange_common.MARKET_YFI_USD: {
					Ticker: `"YFIUSD"`,
				},
				exchange_common.MARKET_USDT_USD: {
					Ticker: `"USDTUSD"`,
				},
			},
			expectedFound: true,
		},
		"Bitfinex": {
			id: "Bitfinex",
			marketToConfig: map[types.MarketId]types.MarketConfig{
				exchange_common.MARKET_BTC_USD: {
					Ticker:         "tBTCUSD",
					AdjustByMarket: newMarketIdWithValue(exchange_common.MARKET_USDT_USD),
				},
				exchange_common.MARKET_ETH_USD: {
					Ticker:         "tETHUSD",
					AdjustByMarket: newMarketIdWithValue(exchange_common.MARKET_USDT_USD),
				},
				exchange_common.MARKET_SOL_USD: {
					Ticker:         "tSOLUSD",
					AdjustByMarket: newMarketIdWithValue(exchange_common.MARKET_USDT_USD),
				},
				exchange_common.MARKET_ADA_USD: {
					Ticker:         "tADAUSD",
					AdjustByMarket: newMarketIdWithValue(exchange_common.MARKET_USDT_USD),
				},
				exchange_common.MARKET_AVAX_USD: {
					Ticker:         "tAVAX:USD",
					AdjustByMarket: newMarketIdWithValue(exchange_common.MARKET_USDT_USD),
				},
				exchange_common.MARKET_DOT_USD: {
					Ticker:         "tDOTUSD",
					AdjustByMarket: newMarketIdWithValue(exchange_common.MARKET_USDT_USD),
				},
				exchange_common.MARKET_XTZ_USD: {
					Ticker:         "tXTZUSD",
					AdjustByMarket: newMarketIdWithValue(exchange_common.MARKET_USDT_USD),
				},
				exchange_common.MARKET_EOS_USD: {
					Ticker:         "tEOSUSD",
					AdjustByMarket: newMarketIdWithValue(exchange_common.MARKET_USDT_USD),
				},
				exchange_common.MARKET_TRX_USD: {
					Ticker:         "tTRXUSD",
					AdjustByMarket: newMarketIdWithValue(exchange_common.MARKET_USDT_USD),
				},
				exchange_common.MARKET_SNX_USD: {
					Ticker:         "tSNXUSD",
					AdjustByMarket: newMarketIdWithValue(exchange_common.MARKET_USDT_USD),
				},
				exchange_common.MARKET_MKR_USD: {
					Ticker:         "tMKRUSD",
					AdjustByMarket: newMarketIdWithValue(exchange_common.MARKET_USDT_USD),
				},
				exchange_common.MARKET_SUSHI_USD: {
					Ticker:         "tSUSHI:USD",
					AdjustByMarket: newMarketIdWithValue(exchange_common.MARKET_USDT_USD),
				},
				exchange_common.MARKET_XLM_USD: {
					Ticker:         "tXLMUSD",
					AdjustByMarket: newMarketIdWithValue(exchange_common.MARKET_USDT_USD),
				},
				exchange_common.MARKET_XMR_USD: {
					Ticker:         "tXMRUSD",
					AdjustByMarket: newMarketIdWithValue(exchange_common.MARKET_USDT_USD),
				},
				exchange_common.MARKET_ZEC_USD: {
					Ticker:         "tZECUSD",
					AdjustByMarket: newMarketIdWithValue(exchange_common.MARKET_USDT_USD),
				},
				exchange_common.MARKET_ZRX_USD: {
					Ticker:         "tZRXUSD",
					AdjustByMarket: newMarketIdWithValue(exchange_common.MARKET_USDT_USD),
				},
				exchange_common.MARKET_YFI_USD: {
					Ticker:         "tYFIUSD",
					AdjustByMarket: newMarketIdWithValue(exchange_common.MARKET_USDT_USD),
				},
				exchange_common.MARKET_USDT_USD: {
					Ticker:         "tUSTUSD",
					AdjustByMarket: newMarketIdWithValue(exchange_common.MARKET_USDT_USD),
				},
			},
			expectedFound: true,
		},
		"Kraken": {
			id: "Kraken",
			marketToConfig: map[types.MarketId]types.MarketConfig{
				exchange_common.MARKET_BTC_USD: {
					Ticker: "XXBTZUSD",
				},
				exchange_common.MARKET_ETH_USD: {
					Ticker: "XETHZUSD",
				},
				exchange_common.MARKET_LINK_USD: {
					Ticker: "LINKUSD",
				},
				exchange_common.MARKET_CRV_USD: {
					Ticker: "CRVUSD",
				},
				exchange_common.MARKET_SOL_USD: {
					Ticker: "SOLUSD",
				},
				exchange_common.MARKET_ADA_USD: {
					Ticker: "ADAUSD",
				},
				exchange_common.MARKET_FIL_USD: {
					Ticker: "FILUSD",
				},
				exchange_common.MARKET_AAVE_USD: {
					Ticker: "AAVEUSD",
				},
				exchange_common.MARKET_LTC_USD: {
					Ticker: "XLTCZUSD",
				},
				exchange_common.MARKET_ATOM_USD: {
					Ticker: "ATOMUSD",
				},
				exchange_common.MARKET_DOT_USD: {
					Ticker: "DOTUSD",
				},
				exchange_common.MARKET_XTZ_USD: {
					Ticker: "XTZUSD",
				},
				exchange_common.MARKET_UNI_USD: {
					Ticker: "UNIUSD",
				},
				exchange_common.MARKET_BCH_USD: {
					Ticker: "BCHUSD",
				},
				exchange_common.MARKET_EOS_USD: {
					Ticker: "EOSUSD",
				},
				exchange_common.MARKET_ALGO_USD: {
					Ticker: "ALGOUSD",
				},
				exchange_common.MARKET_SNX_USD: {
					Ticker: "SNXUSD",
				},
				exchange_common.MARKET_XLM_USD: {
					Ticker: "XXLMZUSD",
				},
				exchange_common.MARKET_XMR_USD: {
					Ticker: "XXMRZUSD",
				},
				exchange_common.MARKET_ETC_USD: {
					Ticker: "XETCZUSD",
				},
				exchange_common.MARKET_COMP_USD: {
					Ticker: "COMPUSD",
				},
				exchange_common.MARKET_ZEC_USD: {
					Ticker: "XZECZUSD",
				},
				exchange_common.MARKET_ZRX_USD: {
					Ticker: "ZRXUSD",
				},
				exchange_common.MARKET_YFI_USD: {
					Ticker: "YFIUSD",
				},
				exchange_common.MARKET_USDT_USD: {
					Ticker: "USDTZUSD",
				},
			},
			expectedFound: true,
		},
		"Gate": {
			id: "Gate",
			marketToConfig: map[types.MarketId]types.MarketConfig{
				exchange_common.MARKET_MATIC_USD: {
					Ticker:         "MATIC_USDT",
					AdjustByMarket: newMarketIdWithValue(exchange_common.MARKET_USDT_USD),
				},
				exchange_common.MARKET_CRV_USD: {
					Ticker:         "CRV_USDT",
					AdjustByMarket: newMarketIdWithValue(exchange_common.MARKET_USDT_USD),
				},
				exchange_common.MARKET_ADA_USD: {
					Ticker:         "ADA_USDT",
					AdjustByMarket: newMarketIdWithValue(exchange_common.MARKET_USDT_USD),
				},
				exchange_common.MARKET_AVAX_USD: {
					Ticker:         "AVAX_USDT",
					AdjustByMarket: newMarketIdWithValue(exchange_common.MARKET_USDT_USD),
				},
				exchange_common.MARKET_DOGE_USD: {
					Ticker:         "DOGE_USDT",
					AdjustByMarket: newMarketIdWithValue(exchange_common.MARKET_USDT_USD),
				},
				exchange_common.MARKET_ICP_USD: {
					Ticker:         "ICP_USDT",
					AdjustByMarket: newMarketIdWithValue(exchange_common.MARKET_USDT_USD),
				},
				exchange_common.MARKET_DOT_USD: {
					Ticker:         "DOT_USDT",
					AdjustByMarket: newMarketIdWithValue(exchange_common.MARKET_USDT_USD),
				},
				exchange_common.MARKET_XTZ_USD: {
					Ticker:         "XTZ_USDT",
					AdjustByMarket: newMarketIdWithValue(exchange_common.MARKET_USDT_USD),
				},
				exchange_common.MARKET_UNI_USD: {
					Ticker:         "UNI_USDT",
					AdjustByMarket: newMarketIdWithValue(exchange_common.MARKET_USDT_USD),
				},
				exchange_common.MARKET_BCH_USD: {
					Ticker:         "BCH_USDT",
					AdjustByMarket: newMarketIdWithValue(exchange_common.MARKET_USDT_USD),
				},
				exchange_common.MARKET_TRX_USD: {
					Ticker:         "TRX_USDT",
					AdjustByMarket: newMarketIdWithValue(exchange_common.MARKET_USDT_USD),
				},
				exchange_common.MARKET_NEAR_USD: {
					Ticker:         "NEAR_USDT",
					AdjustByMarket: newMarketIdWithValue(exchange_common.MARKET_USDT_USD),
				},
				exchange_common.MARKET_MKR_USD: {
					Ticker:         "MKR_USDT",
					AdjustByMarket: newMarketIdWithValue(exchange_common.MARKET_USDT_USD),
				},
				exchange_common.MARKET_SUSHI_USD: {
					Ticker:         "SUSHI_USDT",
					AdjustByMarket: newMarketIdWithValue(exchange_common.MARKET_USDT_USD),
				},
				exchange_common.MARKET_XLM_USD: {
					Ticker:         "XLM_USDT",
					AdjustByMarket: newMarketIdWithValue(exchange_common.MARKET_USDT_USD),
				},
				exchange_common.MARKET_XMR_USD: {
					Ticker:         "XMR_USDT",
					AdjustByMarket: newMarketIdWithValue(exchange_common.MARKET_USDT_USD),
				},
				exchange_common.MARKET_ETC_USD: {
					Ticker:         "ETC_USDT",
					AdjustByMarket: newMarketIdWithValue(exchange_common.MARKET_USDT_USD),
				},
				exchange_common.MARKET_1INCH_USD: {
					Ticker:         "1INCH_USDT",
					AdjustByMarket: newMarketIdWithValue(exchange_common.MARKET_USDT_USD),
				},
				exchange_common.MARKET_USDT_USD: {
					Ticker: "USDT_USD",
				},
			},
			expectedFound: true,
		},
		"Bitstamp": {
			id: "Bitstamp",
			marketToConfig: map[types.MarketId]types.MarketConfig{
				exchange_common.MARKET_BTC_USD: {
					Ticker: "BTC/USD",
				},
				exchange_common.MARKET_ETH_USD: {
					Ticker: "ETH/USD",
				},
				exchange_common.MARKET_USDT_USD: {
					Ticker: "USDT/USD",
				},
			},
			expectedFound: true,
		},
		"Bybit": {
			id: "Bybit",
			marketToConfig: map[types.MarketId]types.MarketConfig{
				exchange_common.MARKET_BTC_USD: {
					Ticker:         "BTCUSDT",
					AdjustByMarket: newMarketIdWithValue(exchange_common.MARKET_USDT_USD),
				},
				exchange_common.MARKET_ETH_USD: {
					Ticker:         "ETHUSDT",
					AdjustByMarket: newMarketIdWithValue(exchange_common.MARKET_USDT_USD),
				},
				exchange_common.MARKET_CRV_USD: {
					Ticker:         "CRVUSDT",
					AdjustByMarket: newMarketIdWithValue(exchange_common.MARKET_USDT_USD),
				},
				exchange_common.MARKET_LTC_USD: {
					Ticker:         "LTCUSDT",
					AdjustByMarket: newMarketIdWithValue(exchange_common.MARKET_USDT_USD),
				},
				exchange_common.MARKET_ATOM_USD: {
					Ticker:         "ATOMUSDT",
					AdjustByMarket: newMarketIdWithValue(exchange_common.MARKET_USDT_USD),
				},
				exchange_common.MARKET_UNI_USD: {
					Ticker:         "UNIUSDT",
					AdjustByMarket: newMarketIdWithValue(exchange_common.MARKET_USDT_USD),
				},
				exchange_common.MARKET_NEAR_USD: {
					Ticker:         "NEARUSDT",
					AdjustByMarket: newMarketIdWithValue(exchange_common.MARKET_USDT_USD),
				},
				exchange_common.MARKET_COMP_USD: {
					Ticker:         "COMPUSDT",
					AdjustByMarket: newMarketIdWithValue(exchange_common.MARKET_USDT_USD),
				},
				exchange_common.MARKET_YFI_USD: {
					Ticker:         "YFIUSDT",
					AdjustByMarket: newMarketIdWithValue(exchange_common.MARKET_USDT_USD),
				},
				exchange_common.MARKET_USDT_USD: {
					Ticker: "USDCUSDT",
					Invert: true,
				},
			},
			expectedFound: true,
		},
		"Crypto.com": {
			id: "CryptoCom",
			marketToConfig: map[types.MarketId]types.MarketConfig{
				exchange_common.MARKET_BTC_USD: {
					Ticker: "BTC_USD",
				},
				exchange_common.MARKET_ETH_USD: {
					Ticker: "ETH_USD",
				},
				exchange_common.MARKET_LINK_USD: {
					Ticker: "LINK_USD",
				},
				exchange_common.MARKET_USDT_USD: {
					Ticker: "USDT_USD",
				},
			},
			expectedFound: true,
		},
		"Huobi": {
			id: "Huobi",
			marketToConfig: map[types.MarketId]types.MarketConfig{
				exchange_common.MARKET_LINK_USD: {
					Ticker:         "linkusdt",
					AdjustByMarket: newMarketIdWithValue(exchange_common.MARKET_USDT_USD),
				},
				exchange_common.MARKET_MATIC_USD: {
					Ticker:         "maticusdt",
					AdjustByMarket: newMarketIdWithValue(exchange_common.MARKET_USDT_USD),
				},
				exchange_common.MARKET_CRV_USD: {
					Ticker:         "crvusdt",
					AdjustByMarket: newMarketIdWithValue(exchange_common.MARKET_USDT_USD),
				},
				exchange_common.MARKET_SOL_USD: {
					Ticker:         "solusdt",
					AdjustByMarket: newMarketIdWithValue(exchange_common.MARKET_USDT_USD),
				},
				exchange_common.MARKET_ADA_USD: {
					Ticker:         "adausdt",
					AdjustByMarket: newMarketIdWithValue(exchange_common.MARKET_USDT_USD),
				},
				exchange_common.MARKET_AVAX_USD: {
					Ticker:         "avaxusdt",
					AdjustByMarket: newMarketIdWithValue(exchange_common.MARKET_USDT_USD),
				},
				exchange_common.MARKET_FIL_USD: {
					Ticker:         "filusdt",
					AdjustByMarket: newMarketIdWithValue(exchange_common.MARKET_USDT_USD),
				},
				exchange_common.MARKET_AAVE_USD: {
					Ticker:         "aaveusdt",
					AdjustByMarket: newMarketIdWithValue(exchange_common.MARKET_USDT_USD),
				},
				exchange_common.MARKET_LTC_USD: {
					Ticker:         "ltcusdt",
					AdjustByMarket: newMarketIdWithValue(exchange_common.MARKET_USDT_USD),
				},
				exchange_common.MARKET_DOGE_USD: {
					Ticker:         "dogeusdt",
					AdjustByMarket: newMarketIdWithValue(exchange_common.MARKET_USDT_USD),
				},
				exchange_common.MARKET_ICP_USD: {
					Ticker:         "icpusdt",
					AdjustByMarket: newMarketIdWithValue(exchange_common.MARKET_USDT_USD),
				},
				exchange_common.MARKET_ATOM_USD: {
					Ticker:         "atomusdt",
					AdjustByMarket: newMarketIdWithValue(exchange_common.MARKET_USDT_USD),
				},
				exchange_common.MARKET_DOT_USD: {
					Ticker:         "dotusdt",
					AdjustByMarket: newMarketIdWithValue(exchange_common.MARKET_USDT_USD),
				},
				exchange_common.MARKET_XTZ_USD: {
					Ticker:         "xtzusdt",
					AdjustByMarket: newMarketIdWithValue(exchange_common.MARKET_USDT_USD),
				},
				exchange_common.MARKET_UNI_USD: {
					Ticker:         "uniusdt",
					AdjustByMarket: newMarketIdWithValue(exchange_common.MARKET_USDT_USD),
				},
				exchange_common.MARKET_BCH_USD: {
					Ticker:         "bchusdt",
					AdjustByMarket: newMarketIdWithValue(exchange_common.MARKET_USDT_USD),
				},
				exchange_common.MARKET_EOS_USD: {
					Ticker:         "eosusdt",
					AdjustByMarket: newMarketIdWithValue(exchange_common.MARKET_USDT_USD),
				},
				exchange_common.MARKET_TRX_USD: {
					Ticker:         "trxusdt",
					AdjustByMarket: newMarketIdWithValue(exchange_common.MARKET_USDT_USD),
				},
				exchange_common.MARKET_ALGO_USD: {
					Ticker:         "algousdt",
					AdjustByMarket: newMarketIdWithValue(exchange_common.MARKET_USDT_USD),
				},
				exchange_common.MARKET_NEAR_USD: {
					Ticker:         "nearusdt",
					AdjustByMarket: newMarketIdWithValue(exchange_common.MARKET_USDT_USD),
				},
				exchange_common.MARKET_SNX_USD: {
					Ticker:         "snxusdt",
					AdjustByMarket: newMarketIdWithValue(exchange_common.MARKET_USDT_USD),
				},
				exchange_common.MARKET_MKR_USD: {
					Ticker:         "mkrusdt",
					AdjustByMarket: newMarketIdWithValue(exchange_common.MARKET_USDT_USD),
				},
				exchange_common.MARKET_SUSHI_USD: {
					Ticker:         "sushiusdt",
					AdjustByMarket: newMarketIdWithValue(exchange_common.MARKET_USDT_USD),
				},
				exchange_common.MARKET_ETC_USD: {
					Ticker:         "etcusdt",
					AdjustByMarket: newMarketIdWithValue(exchange_common.MARKET_USDT_USD),
				},
				exchange_common.MARKET_1INCH_USD: {
					Ticker:         "1inchusdt",
					AdjustByMarket: newMarketIdWithValue(exchange_common.MARKET_USDT_USD),
				},
				exchange_common.MARKET_COMP_USD: {
					Ticker:         "compusdt",
					AdjustByMarket: newMarketIdWithValue(exchange_common.MARKET_USDT_USD),
				},
				exchange_common.MARKET_ZRX_USD: {
					Ticker:         "zrxusdt",
					AdjustByMarket: newMarketIdWithValue(exchange_common.MARKET_USDT_USD),
				},
				exchange_common.MARKET_YFI_USD: {
					Ticker:         "yfiusdt",
					AdjustByMarket: newMarketIdWithValue(exchange_common.MARKET_USDT_USD),
				},
				exchange_common.MARKET_USDT_USD: {
					Ticker:         "ethusdt", // Adjusted with ETH index price.
					AdjustByMarket: newMarketIdWithValue(exchange_common.MARKET_ETH_USD),
					Invert:         true,
				},
			},
			expectedFound: true,
		},
		"Kucoin": {
			id: "Kucoin",
			marketToConfig: map[types.MarketId]types.MarketConfig{
				exchange_common.MARKET_LINK_USD: {
					Ticker:         "LINK-USDT",
					AdjustByMarket: newMarketIdWithValue(exchange_common.MARKET_USDT_USD),
				},
				exchange_common.MARKET_MATIC_USD: {
					Ticker:         "MATIC-USDT",
					AdjustByMarket: newMarketIdWithValue(exchange_common.MARKET_USDT_USD),
				},
				exchange_common.MARKET_CRV_USD: {
					Ticker:         "CRV-USDT",
					AdjustByMarket: newMarketIdWithValue(exchange_common.MARKET_USDT_USD),
				},
				exchange_common.MARKET_SOL_USD: {
					Ticker:         "SOL-USDT",
					AdjustByMarket: newMarketIdWithValue(exchange_common.MARKET_USDT_USD),
				},
				exchange_common.MARKET_ADA_USD: {
					Ticker:         "ADA-USDT",
					AdjustByMarket: newMarketIdWithValue(exchange_common.MARKET_USDT_USD),
				},
				exchange_common.MARKET_AVAX_USD: {
					Ticker:         "AVAX-USDT",
					AdjustByMarket: newMarketIdWithValue(exchange_common.MARKET_USDT_USD),
				},
				exchange_common.MARKET_FIL_USD: {
					Ticker:         "FIL-USDT",
					AdjustByMarket: newMarketIdWithValue(exchange_common.MARKET_USDT_USD),
				},
				exchange_common.MARKET_AAVE_USD: {
					Ticker:         "AAVE-USDT",
					AdjustByMarket: newMarketIdWithValue(exchange_common.MARKET_USDT_USD),
				},
				exchange_common.MARKET_LTC_USD: {
					Ticker:         "LTC-USDT",
					AdjustByMarket: newMarketIdWithValue(exchange_common.MARKET_USDT_USD),
				},
				exchange_common.MARKET_DOGE_USD: {
					Ticker:         "DOGE-USDT",
					AdjustByMarket: newMarketIdWithValue(exchange_common.MARKET_USDT_USD),
				},
				exchange_common.MARKET_ICP_USD: {
					Ticker:         "ICP-USDT",
					AdjustByMarket: newMarketIdWithValue(exchange_common.MARKET_USDT_USD),
				},
				exchange_common.MARKET_ATOM_USD: {
					Ticker:         "ATOM-USDT",
					AdjustByMarket: newMarketIdWithValue(exchange_common.MARKET_USDT_USD),
				},
				exchange_common.MARKET_DOT_USD: {
					Ticker:         "DOT-USDT",
					AdjustByMarket: newMarketIdWithValue(exchange_common.MARKET_USDT_USD),
				},
				exchange_common.MARKET_ALGO_USD: {
					Ticker:         "ALGO-USDT",
					AdjustByMarket: newMarketIdWithValue(exchange_common.MARKET_USDT_USD),
				},
				exchange_common.MARKET_NEAR_USD: {
					Ticker:         "NEAR-USDT",
					AdjustByMarket: newMarketIdWithValue(exchange_common.MARKET_USDT_USD),
				},
				exchange_common.MARKET_SNX_USD: {
					Ticker:         "SNX-USDT",
					AdjustByMarket: newMarketIdWithValue(exchange_common.MARKET_USDT_USD),
				},
				exchange_common.MARKET_MKR_USD: {
					Ticker:         "MKR-USDT",
					AdjustByMarket: newMarketIdWithValue(exchange_common.MARKET_USDT_USD),
				},
				exchange_common.MARKET_XLM_USD: {
					Ticker:         "XLM-USDT",
					AdjustByMarket: newMarketIdWithValue(exchange_common.MARKET_USDT_USD),
				},
				exchange_common.MARKET_XMR_USD: {
					Ticker:         "XMR-USDT",
					AdjustByMarket: newMarketIdWithValue(exchange_common.MARKET_USDT_USD),
				},
				exchange_common.MARKET_1INCH_USD: {
					Ticker:         "1INCH-USDT",
					AdjustByMarket: newMarketIdWithValue(exchange_common.MARKET_USDT_USD),
				},
				exchange_common.MARKET_COMP_USD: {
					Ticker:         "COMP-USDT",
					AdjustByMarket: newMarketIdWithValue(exchange_common.MARKET_USDT_USD),
				},
				exchange_common.MARKET_ZEC_USD: {
					Ticker:         "ZEC-USDT",
					AdjustByMarket: newMarketIdWithValue(exchange_common.MARKET_USDT_USD),
				},
				exchange_common.MARKET_USDT_USD: {
					Ticker:         "BTC-USDT", // Adjusted with BTC index price.
					AdjustByMarket: newMarketIdWithValue(exchange_common.MARKET_BTC_USD),
					Invert:         true,
				},
			},
			expectedFound: true,
		},
		"Okx": {
			id: "Okx",
			marketToConfig: map[types.MarketId]types.MarketConfig{
				exchange_common.MARKET_BTC_USD: {
					Ticker:         "BTC-USDT",
					AdjustByMarket: newMarketIdWithValue(exchange_common.MARKET_USDT_USD),
				},
				exchange_common.MARKET_ETH_USD: {
					Ticker:         "ETH-USDT",
					AdjustByMarket: newMarketIdWithValue(exchange_common.MARKET_USDT_USD),
				},
				exchange_common.MARKET_LINK_USD: {
					Ticker:         "LINK-USDT",
					AdjustByMarket: newMarketIdWithValue(exchange_common.MARKET_USDT_USD),
				},
				exchange_common.MARKET_MATIC_USD: {
					Ticker:         "MATIC-USDT",
					AdjustByMarket: newMarketIdWithValue(exchange_common.MARKET_USDT_USD),
				},
				exchange_common.MARKET_CRV_USD: {
					Ticker:         "CRV-USDT",
					AdjustByMarket: newMarketIdWithValue(exchange_common.MARKET_USDT_USD),
				},
				exchange_common.MARKET_SOL_USD: {
					Ticker:         "SOL-USDT",
					AdjustByMarket: newMarketIdWithValue(exchange_common.MARKET_USDT_USD),
				},
				exchange_common.MARKET_AVAX_USD: {
					Ticker:         "AVAX-USDT",
					AdjustByMarket: newMarketIdWithValue(exchange_common.MARKET_USDT_USD),
				},
				exchange_common.MARKET_FIL_USD: {
					Ticker:         "FIL-USDT",
					AdjustByMarket: newMarketIdWithValue(exchange_common.MARKET_USDT_USD),
				},
				exchange_common.MARKET_AAVE_USD: {
					Ticker:         "AAVE-USDT",
					AdjustByMarket: newMarketIdWithValue(exchange_common.MARKET_USDT_USD),
				},
				exchange_common.MARKET_LTC_USD: {
					Ticker:         "LTC-USDT",
					AdjustByMarket: newMarketIdWithValue(exchange_common.MARKET_USDT_USD),
				},
				exchange_common.MARKET_DOGE_USD: {
					Ticker:         "DOGE-USDT",
					AdjustByMarket: newMarketIdWithValue(exchange_common.MARKET_USDT_USD),
				},
				exchange_common.MARKET_ICP_USD: {
					Ticker:         "ICP-USDT",
					AdjustByMarket: newMarketIdWithValue(exchange_common.MARKET_USDT_USD),
				},
				exchange_common.MARKET_ATOM_USD: {
					Ticker:         "ATOM-USDT",
					AdjustByMarket: newMarketIdWithValue(exchange_common.MARKET_USDT_USD),
				},
				exchange_common.MARKET_DOT_USD: {
					Ticker:         "DOT-USDT",
					AdjustByMarket: newMarketIdWithValue(exchange_common.MARKET_USDT_USD),
				},
				exchange_common.MARKET_XTZ_USD: {
					Ticker:         "XTZ-USDT",
					AdjustByMarket: newMarketIdWithValue(exchange_common.MARKET_USDT_USD),
				},
				exchange_common.MARKET_UNI_USD: {
					Ticker:         "UNI-USDT",
					AdjustByMarket: newMarketIdWithValue(exchange_common.MARKET_USDT_USD),
				},
				exchange_common.MARKET_BCH_USD: {
					Ticker:         "BCH-USDT",
					AdjustByMarket: newMarketIdWithValue(exchange_common.MARKET_USDT_USD),
				},
				exchange_common.MARKET_EOS_USD: {
					Ticker:         "EOS-USDT",
					AdjustByMarket: newMarketIdWithValue(exchange_common.MARKET_USDT_USD),
				},
				exchange_common.MARKET_TRX_USD: {
					Ticker:         "TRX-USDT",
					AdjustByMarket: newMarketIdWithValue(exchange_common.MARKET_USDT_USD),
				},
				exchange_common.MARKET_ALGO_USD: {
					Ticker:         "ALGO-USDT",
					AdjustByMarket: newMarketIdWithValue(exchange_common.MARKET_USDT_USD),
				},
				exchange_common.MARKET_NEAR_USD: {
					Ticker:         "NEAR-USDT",
					AdjustByMarket: newMarketIdWithValue(exchange_common.MARKET_USDT_USD),
				},
				exchange_common.MARKET_SNX_USD: {
					Ticker:         "SNX-USDT",
					AdjustByMarket: newMarketIdWithValue(exchange_common.MARKET_USDT_USD),
				},
				exchange_common.MARKET_MKR_USD: {
					Ticker:         "MKR-USDT",
					AdjustByMarket: newMarketIdWithValue(exchange_common.MARKET_USDT_USD),
				},
				exchange_common.MARKET_SUSHI_USD: {
					Ticker:         "SUSHI-USDT",
					AdjustByMarket: newMarketIdWithValue(exchange_common.MARKET_USDT_USD),
				},
				exchange_common.MARKET_XLM_USD: {
					Ticker:         "XLM-USDT",
					AdjustByMarket: newMarketIdWithValue(exchange_common.MARKET_USDT_USD),
				},
				exchange_common.MARKET_XMR_USD: {
					Ticker:         "XMR-USDT",
					AdjustByMarket: newMarketIdWithValue(exchange_common.MARKET_USDT_USD),
				},
				exchange_common.MARKET_ETC_USD: {
					Ticker:         "ETC-USDT",
					AdjustByMarket: newMarketIdWithValue(exchange_common.MARKET_USDT_USD),
				},
				exchange_common.MARKET_1INCH_USD: {
					Ticker:         "1INCH-USDT",
					AdjustByMarket: newMarketIdWithValue(exchange_common.MARKET_USDT_USD),
				},
				exchange_common.MARKET_COMP_USD: {
					Ticker:         "COMP-USDT",
					AdjustByMarket: newMarketIdWithValue(exchange_common.MARKET_USDT_USD),
				},
				exchange_common.MARKET_ZEC_USD: {
					Ticker:         "ZEC-USDT",
					AdjustByMarket: newMarketIdWithValue(exchange_common.MARKET_USDT_USD),
				},
				exchange_common.MARKET_ZRX_USD: {
					Ticker:         "ZRX-USDT",
					AdjustByMarket: newMarketIdWithValue(exchange_common.MARKET_USDT_USD),
				},
				exchange_common.MARKET_YFI_USD: {
					Ticker:         "YFI-USDT",
					AdjustByMarket: newMarketIdWithValue(exchange_common.MARKET_USDT_USD),
				},
				exchange_common.MARKET_USDT_USD: {
					Ticker:         "BTC-USDT", // Adjusted with BTC index price.
					AdjustByMarket: newMarketIdWithValue(exchange_common.MARKET_BTC_USD),
					Invert:         true,
				},
			},
			expectedFound: true,
		},
		"Mexc": {
			id: "Mexc",
			marketToConfig: map[types.MarketId]types.MarketConfig{
				exchange_common.MARKET_UNI_USD: {
					Ticker:         "UNI_USDT",
					AdjustByMarket: newMarketIdWithValue(exchange_common.MARKET_USDT_USD),
				},
				exchange_common.MARKET_XMR_USD: {
					Ticker:         "XMR_USDT",
					AdjustByMarket: newMarketIdWithValue(exchange_common.MARKET_USDT_USD),
				},
			},
			expectedFound: true,
		},
		"Coinbase Pro": {
			id: "CoinbasePro",
			marketToConfig: map[types.MarketId]types.MarketConfig{
				exchange_common.MARKET_BTC_USD: {
					Ticker: "BTC-USD",
				},
				exchange_common.MARKET_ETH_USD: {
					Ticker: "ETH-USD",
				},
				exchange_common.MARKET_LINK_USD: {
					Ticker: "LINK-USD",
				},
				exchange_common.MARKET_MATIC_USD: {
					Ticker: "MATIC-USD",
				},
				exchange_common.MARKET_CRV_USD: {
					Ticker: "CRV-USD",
				},
				exchange_common.MARKET_SOL_USD: {
					Ticker: "SOL-USD",
				},
				exchange_common.MARKET_ADA_USD: {
					Ticker: "ADA-USD",
				},
				exchange_common.MARKET_FIL_USD: {
					Ticker: "FIL-USD",
				},
				exchange_common.MARKET_AAVE_USD: {
					Ticker: "AAVE-USD",
				},
				exchange_common.MARKET_LTC_USD: {
					Ticker: "LTC-USD",
				},
				exchange_common.MARKET_ICP_USD: {
					Ticker: "ICP-USD",
				},
				exchange_common.MARKET_ATOM_USD: {
					Ticker: "ATOM-USD",
				},
				exchange_common.MARKET_XTZ_USD: {
					Ticker: "XTZ-USD",
				},
				exchange_common.MARKET_UNI_USD: {
					Ticker: "UNI-USD",
				},
				exchange_common.MARKET_BCH_USD: {
					Ticker: "BCH-USD",
				},
				exchange_common.MARKET_EOS_USD: {
					Ticker: "EOS-USD",
				},
				exchange_common.MARKET_ALGO_USD: {
					Ticker: "ALGO-USD",
				},
				exchange_common.MARKET_NEAR_USD: {
					Ticker: "NEAR-USD",
				},
				exchange_common.MARKET_SNX_USD: {
					Ticker: "SNX-USD",
				},
				exchange_common.MARKET_MKR_USD: {
					Ticker: "MKR-USD",
				},
				exchange_common.MARKET_SUSHI_USD: {
					Ticker: "SUSHI-USD",
				},
				exchange_common.MARKET_XLM_USD: {
					Ticker: "XLM-USD",
				},
				exchange_common.MARKET_ETC_USD: {
					Ticker: "ETC-USD",
				},
				exchange_common.MARKET_1INCH_USD: {
					Ticker: "1INCH-USD",
				},
				exchange_common.MARKET_COMP_USD: {
					Ticker: "COMP-USD",
				},
				exchange_common.MARKET_ZEC_USD: {
					Ticker: "ZEC-USD",
				},
				exchange_common.MARKET_ZRX_USD: {
					Ticker: "ZRX-USD",
				},
				exchange_common.MARKET_YFI_USD: {
					Ticker: "YFI-USD",
				},
				exchange_common.MARKET_USDT_USD: {
					Ticker: "USDT-USD",
				},
			},
			expectedFound: true,
		},
		"Test exchange": {
			id: "TestExchange",
			marketToConfig: map[types.MarketId]types.MarketConfig{
				exchange_common.MARKET_BTC_USD: {
					Ticker: "BTC-USD",
				},
				exchange_common.MARKET_ETH_USD: {
					Ticker: "ETH-USD",
				},
				exchange_common.MARKET_LINK_USD: {
					Ticker: "LINK-USD",
				},
				exchange_common.MARKET_USDT_USD: {
					Ticker: "USDT-USD",
				},
			},
			expectedFound: true,
		},
		"unknown": {
			id:            "unknown",
			expectedFound: false,
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			config, ok := StaticExchangeMarketConfig[tc.id]
			if tc.expectedFound {
				require.True(t, ok)
				require.Equal(t, tc.id, config.Id)
				require.Len(t, config.MarketToMarketConfig, len(tc.marketToConfig))
				for market, expectedConfig := range tc.marketToConfig {
					actualConfig, ok := config.MarketToMarketConfig[market]
					require.True(t, ok, "Market %v missing from exchange market config", market)
					require.Equal(t, expectedConfig, actualConfig)
				}
			} else {
				require.False(t, ok)
			}
		})
	}
}

func TestStaticExchangeMarketConfigCacheLen(t *testing.T) {
	require.Len(t, StaticExchangeMarketConfig, 14)
}

func TestGenerateExchangeConfigJsonLength(t *testing.T) {
	configs := GenerateExchangeConfigJson(StaticExchangeMarketConfig)
	require.Len(t, configs, 34)
}

func TestGenerateExchangeConfigJson(t *testing.T) {
	tests := map[string]struct {
		id                             types.MarketId
		expectedExchangeConfigJsonFile string
	}{
		"BTC exchange config": {
			id:                             exchange_common.MARKET_BTC_USD,
			expectedExchangeConfigJsonFile: "btc_exchange_config.json",
		},
		"ETH exchange config": {
			id:                             exchange_common.MARKET_ETH_USD,
			expectedExchangeConfigJsonFile: "eth_exchange_config.json",
		},
		"LINK exchange config": {
			id:                             exchange_common.MARKET_LINK_USD,
			expectedExchangeConfigJsonFile: "link_exchange_config.json",
		},
		"MATIC exchange config": {
			id:                             exchange_common.MARKET_MATIC_USD,
			expectedExchangeConfigJsonFile: "matic_exchange_config.json",
		},
		"CRV exchange config": {
			id:                             exchange_common.MARKET_CRV_USD,
			expectedExchangeConfigJsonFile: "crv_exchange_config.json",
		},
		"SOL exchange config": {
			id:                             exchange_common.MARKET_SOL_USD,
			expectedExchangeConfigJsonFile: "sol_exchange_config.json",
		},
		"ADA exchange config": {
			id:                             exchange_common.MARKET_ADA_USD,
			expectedExchangeConfigJsonFile: "ada_exchange_config.json",
		},
		"AVAX exchange config": {
			id:                             exchange_common.MARKET_AVAX_USD,
			expectedExchangeConfigJsonFile: "avax_exchange_config.json",
		},
		"FIL exchange config": {
			id:                             exchange_common.MARKET_FIL_USD,
			expectedExchangeConfigJsonFile: "fil_exchange_config.json",
		},
		"AAVE exchange config": {
			id:                             exchange_common.MARKET_AAVE_USD,
			expectedExchangeConfigJsonFile: "aave_exchange_config.json",
		},
		"LTC exchange config": {
			id:                             exchange_common.MARKET_LTC_USD,
			expectedExchangeConfigJsonFile: "ltc_exchange_config.json",
		},
		"DOGE exchange config": {
			id:                             exchange_common.MARKET_DOGE_USD,
			expectedExchangeConfigJsonFile: "doge_exchange_config.json",
		},
		"ICP exchange config": {
			id:                             exchange_common.MARKET_ICP_USD,
			expectedExchangeConfigJsonFile: "icp_exchange_config.json",
		},
		"ATOM exchange config": {
			id:                             exchange_common.MARKET_ATOM_USD,
			expectedExchangeConfigJsonFile: "atom_exchange_config.json",
		},
		"DOT exchange config": {
			id:                             exchange_common.MARKET_DOT_USD,
			expectedExchangeConfigJsonFile: "dot_exchange_config.json",
		},
		"XTZ exchange config": {
			id:                             exchange_common.MARKET_XTZ_USD,
			expectedExchangeConfigJsonFile: "xtz_exchange_config.json",
		},
		"UNI exchange config": {
			id:                             exchange_common.MARKET_UNI_USD,
			expectedExchangeConfigJsonFile: "uni_exchange_config.json",
		},
		"BCH exchange config": {
			id:                             exchange_common.MARKET_BCH_USD,
			expectedExchangeConfigJsonFile: "bch_exchange_config.json",
		},
		"EOS exchange config": {
			id:                             exchange_common.MARKET_EOS_USD,
			expectedExchangeConfigJsonFile: "eos_exchange_config.json",
		},
		"TRX exchange config": {
			id:                             exchange_common.MARKET_TRX_USD,
			expectedExchangeConfigJsonFile: "trx_exchange_config.json",
		},
		"ALGO exchange config": {
			id:                             exchange_common.MARKET_ALGO_USD,
			expectedExchangeConfigJsonFile: "algo_exchange_config.json",
		},
		"NEAR exchange config": {
			id:                             exchange_common.MARKET_NEAR_USD,
			expectedExchangeConfigJsonFile: "near_exchange_config.json",
		},
		"SNX exchange config": {
			id:                             exchange_common.MARKET_SNX_USD,
			expectedExchangeConfigJsonFile: "snx_exchange_config.json",
		},
		"MKR exchange config": {
			id:                             exchange_common.MARKET_MKR_USD,
			expectedExchangeConfigJsonFile: "mkr_exchange_config.json",
		},
		"SUSHI exchange config": {
			id:                             exchange_common.MARKET_SUSHI_USD,
			expectedExchangeConfigJsonFile: "sushi_exchange_config.json",
		},
		"XLM exchange config": {
			id:                             exchange_common.MARKET_XLM_USD,
			expectedExchangeConfigJsonFile: "xlm_exchange_config.json",
		},
		"XMR exchange config": {
			id:                             exchange_common.MARKET_XMR_USD,
			expectedExchangeConfigJsonFile: "xmr_exchange_config.json",
		},
		"ETC exchange config": {
			id:                             exchange_common.MARKET_ETC_USD,
			expectedExchangeConfigJsonFile: "etc_exchange_config.json",
		},
		"1INCH exchange config": {
			id:                             exchange_common.MARKET_1INCH_USD,
			expectedExchangeConfigJsonFile: "1inch_exchange_config.json",
		},
		"COMP exchange config": {
			id:                             exchange_common.MARKET_COMP_USD,
			expectedExchangeConfigJsonFile: "comp_exchange_config.json",
		},
		"ZEC exchange config": {
			id:                             exchange_common.MARKET_ZEC_USD,
			expectedExchangeConfigJsonFile: "zec_exchange_config.json",
		},
		"ZRX exchange config": {
			id:                             exchange_common.MARKET_ZRX_USD,
			expectedExchangeConfigJsonFile: "zrx_exchange_config.json",
		},
		"YFI exchange config": {
			id:                             exchange_common.MARKET_YFI_USD,
			expectedExchangeConfigJsonFile: "yfi_exchange_config.json",
		},
		"USDT exchange config": {
			id:                             exchange_common.MARKET_USDT_USD,
			expectedExchangeConfigJsonFile: "usdt_exchange_config.json",
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			configs := GenerateExchangeConfigJson(StaticExchangeMarketConfig)

			// Uncomment to update test data
			//f, err := os.OpenFile("testdata/"+tc.expectedExchangeConfigJsonFile, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
			//require.NoError(t, err)
			//defer f.Close()
			//_, err = f.WriteString(configs[tc.id] + "\n") // Final newline added manually.
			//require.NoError(t, err)

			actualExchangeConfigJson := json.CompactJsonString(t, configs[tc.id])
			expectedExchangeConfigJson := pricefeed.ReadJsonTestFile(t, tc.expectedExchangeConfigJsonFile)
			require.Equal(t, expectedExchangeConfigJson, actualExchangeConfigJson)
		})
	}
}
