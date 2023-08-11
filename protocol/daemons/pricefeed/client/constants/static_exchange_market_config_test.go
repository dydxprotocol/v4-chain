package constants

import (
	"github.com/dydxprotocol/v4/daemons/pricefeed/client/constants/exchange_common"
	"github.com/dydxprotocol/v4/daemons/pricefeed/client/types"
	"github.com/dydxprotocol/v4/testutil/daemons/pricefeed"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestStaticExchangeMarketConfigCache(t *testing.T) {
	tests := map[string]struct {
		id             types.ExchangeId
		marketToTicker map[types.MarketId]string
		expectedFound  bool
	}{
		"Binance": {
			id: "Binance",
			marketToTicker: map[types.MarketId]string{
				exchange_common.MARKET_BTC_USD:   `"BTCUSDT"`,
				exchange_common.MARKET_ETH_USD:   `"ETHUSDT"`,
				exchange_common.MARKET_LINK_USD:  `"LINKUSDT"`,
				exchange_common.MARKET_MATIC_USD: `"MATICUSDT"`,
				exchange_common.MARKET_CRV_USD:   `"CRVUSDT"`,
				exchange_common.MARKET_SOL_USD:   `"SOLUSDT"`,
				exchange_common.MARKET_ADA_USD:   `"ADAUSDT"`,
				exchange_common.MARKET_AVAX_USD:  `"AVAXUSDT"`,
				exchange_common.MARKET_FIL_USD:   `"FILUSDT"`,
				exchange_common.MARKET_AAVE_USD:  `"AAVEUSDT"`,
				exchange_common.MARKET_LTC_USD:   `"LTCUSDT"`,
				exchange_common.MARKET_DOGE_USD:  `"DOGEUSDT"`,
				exchange_common.MARKET_ICP_USD:   `"ICPUSDT"`,
				exchange_common.MARKET_ATOM_USD:  `"ATOMUSDT"`,
				exchange_common.MARKET_DOT_USD:   `"DOTUSDT"`,
				exchange_common.MARKET_XTZ_USD:   `"XTZUSDT"`,
				exchange_common.MARKET_UNI_USD:   `"UNIUSDT"`,
				exchange_common.MARKET_BCH_USD:   `"BCHUSDT"`,
				exchange_common.MARKET_EOS_USD:   `"EOSUSDT"`,
				exchange_common.MARKET_TRX_USD:   `"TRXUSDT"`,
				exchange_common.MARKET_ALGO_USD:  `"ALGOUSDT"`,
				exchange_common.MARKET_NEAR_USD:  `"NEARUSDT"`,
				exchange_common.MARKET_SNX_USD:   `"SNXUSDT"`,
				exchange_common.MARKET_MKR_USD:   `"MKRUSDT"`,
				exchange_common.MARKET_SUSHI_USD: `"SUSHIUSDT"`,
				exchange_common.MARKET_XLM_USD:   `"XLMUSDT"`,
				exchange_common.MARKET_XMR_USD:   `"XMRUSDT"`,
				exchange_common.MARKET_ETC_USD:   `"ETCUSDT"`,
				exchange_common.MARKET_1INCH_USD: `"1INCHUSDT"`,
				exchange_common.MARKET_COMP_USD:  `"COMPUSDT"`,
				exchange_common.MARKET_ZEC_USD:   `"ZECUSDT"`,
				exchange_common.MARKET_ZRX_USD:   `"ZRXUSDT"`,
				exchange_common.MARKET_YFI_USD:   `"YFIUSDT"`,
			},
			expectedFound: true,
		},
		"BinanceUS": {
			id: "BinanceUS",
			marketToTicker: map[types.MarketId]string{
				exchange_common.MARKET_BTC_USD:   `"BTCUSD"`,
				exchange_common.MARKET_ETH_USD:   `"ETHUSD"`,
				exchange_common.MARKET_LINK_USD:  `"LINKUSD"`,
				exchange_common.MARKET_MATIC_USD: `"MATICUSD"`,
				exchange_common.MARKET_CRV_USD:   `"CRVUSD"`,
				exchange_common.MARKET_SOL_USD:   `"SOLUSD"`,
				exchange_common.MARKET_ADA_USD:   `"ADAUSD"`,
				exchange_common.MARKET_AVAX_USD:  `"AVAXUSD"`,
				exchange_common.MARKET_FIL_USD:   `"FILUSD"`,
				exchange_common.MARKET_AAVE_USD:  `"AAVEUSD"`,
				exchange_common.MARKET_LTC_USD:   `"LTCUSD"`,
				exchange_common.MARKET_DOGE_USD:  `"DOGEUSD"`,
				exchange_common.MARKET_ICP_USD:   `"ICPUSD"`,
				exchange_common.MARKET_ATOM_USD:  `"ATOMUSD"`,
				exchange_common.MARKET_DOT_USD:   `"DOTUSD"`,
				exchange_common.MARKET_XTZ_USD:   `"XTZUSD"`,
				exchange_common.MARKET_UNI_USD:   `"UNIUSD"`,
				exchange_common.MARKET_BCH_USD:   `"BCHUSD"`,
				exchange_common.MARKET_EOS_USD:   `"EOSUSD"`,
				exchange_common.MARKET_ALGO_USD:  `"ALGOUSD"`,
				exchange_common.MARKET_NEAR_USD:  `"NEARUSD"`,
				exchange_common.MARKET_SNX_USD:   `"SNXUSD"`,
				exchange_common.MARKET_MKR_USD:   `"MKRUSD"`,
				exchange_common.MARKET_SUSHI_USD: `"SUSHIUSD"`,
				exchange_common.MARKET_XLM_USD:   `"XLMUSD"`,
				exchange_common.MARKET_ETC_USD:   `"ETCUSD"`,
				exchange_common.MARKET_1INCH_USD: `"1INCHUSD"`,
				exchange_common.MARKET_COMP_USD:  `"COMPUSD"`,
				exchange_common.MARKET_ZEC_USD:   `"ZECUSD"`,
				exchange_common.MARKET_ZRX_USD:   `"ZRXUSD"`,
				exchange_common.MARKET_YFI_USD:   `"YFIUSD"`,
			},
			expectedFound: true,
		},
		"Bitfinex": {
			id: "Bitfinex",
			marketToTicker: map[types.MarketId]string{
				exchange_common.MARKET_BTC_USD:   "tBTCUSD",
				exchange_common.MARKET_ETH_USD:   "tETHUSD",
				exchange_common.MARKET_SOL_USD:   "tSOLUSD",
				exchange_common.MARKET_ADA_USD:   "tADAUSD",
				exchange_common.MARKET_AVAX_USD:  "tAVAX:USD",
				exchange_common.MARKET_DOT_USD:   "tDOTUSD",
				exchange_common.MARKET_XTZ_USD:   "tXTZUSD",
				exchange_common.MARKET_EOS_USD:   "tEOSUSD",
				exchange_common.MARKET_TRX_USD:   "tTRXUSD",
				exchange_common.MARKET_SNX_USD:   "tSNXUSD",
				exchange_common.MARKET_MKR_USD:   "tMKRUSD",
				exchange_common.MARKET_SUSHI_USD: "tSUSHI:USD",
				exchange_common.MARKET_XLM_USD:   "tXLMUSD",
				exchange_common.MARKET_XMR_USD:   "tXMRUSD",
				exchange_common.MARKET_ZEC_USD:   "tZECUSD",
				exchange_common.MARKET_ZRX_USD:   "tZRXUSD",
				exchange_common.MARKET_YFI_USD:   "tYFIUSD",
			},
			expectedFound: true,
		},
		"Kraken": {
			id: "Kraken",
			marketToTicker: map[types.MarketId]string{
				exchange_common.MARKET_BTC_USD:  "XXBTZUSD",
				exchange_common.MARKET_ETH_USD:  "XETHZUSD",
				exchange_common.MARKET_LINK_USD: "LINKUSD",
				exchange_common.MARKET_CRV_USD:  "CRVUSD",
				exchange_common.MARKET_SOL_USD:  "SOLUSD",
				exchange_common.MARKET_ADA_USD:  "ADAUSD",
				exchange_common.MARKET_FIL_USD:  "FILUSD",
				exchange_common.MARKET_AAVE_USD: "AAVEUSD",
				exchange_common.MARKET_LTC_USD:  "XLTCZUSD",
				exchange_common.MARKET_ATOM_USD: "ATOMUSD",
				exchange_common.MARKET_DOT_USD:  "DOTUSD",
				exchange_common.MARKET_XTZ_USD:  "XTZUSD",
				exchange_common.MARKET_UNI_USD:  "UNIUSD",
				exchange_common.MARKET_BCH_USD:  "BCHUSD",
				exchange_common.MARKET_EOS_USD:  "EOSUSD",
				exchange_common.MARKET_ALGO_USD: "ALGOUSD",
				exchange_common.MARKET_SNX_USD:  "SNXUSD",
				exchange_common.MARKET_XLM_USD:  "XXLMZUSD",
				exchange_common.MARKET_XMR_USD:  "XXMRZUSD",
				exchange_common.MARKET_ETC_USD:  "XETCZUSD",
				exchange_common.MARKET_COMP_USD: "COMPUSD",
				exchange_common.MARKET_ZEC_USD:  "XZECZUSD",
				exchange_common.MARKET_ZRX_USD:  "ZRXUSD",
				exchange_common.MARKET_YFI_USD:  "YFIUSD",
			},
			expectedFound: true,
		},
		"Gate": {
			id: "Gate",
			marketToTicker: map[types.MarketId]string{
				exchange_common.MARKET_MATIC_USD: "MATIC_USDT",
				exchange_common.MARKET_CRV_USD:   "CRV_USDT",
				exchange_common.MARKET_ADA_USD:   "ADA_USDT",
				exchange_common.MARKET_AVAX_USD:  "AVAX_USDT",
				exchange_common.MARKET_DOGE_USD:  "DOGE_USDT",
				exchange_common.MARKET_ICP_USD:   "ICP_USDT",
				exchange_common.MARKET_DOT_USD:   "DOT_USDT",
				exchange_common.MARKET_XTZ_USD:   "XTZ_USDT",
				exchange_common.MARKET_UNI_USD:   "UNI_USDT",
				exchange_common.MARKET_BCH_USD:   "BCH_USDT",
				exchange_common.MARKET_TRX_USD:   "TRX_USDT",
				exchange_common.MARKET_NEAR_USD:  "NEAR_USDT",
				exchange_common.MARKET_MKR_USD:   "MKR_USDT",
				exchange_common.MARKET_SUSHI_USD: "SUSHI_USDT",
				exchange_common.MARKET_XLM_USD:   "XLM_USDT",
				exchange_common.MARKET_XMR_USD:   "XMR_USDT",
				exchange_common.MARKET_ETC_USD:   "ETC_USDT",
				exchange_common.MARKET_1INCH_USD: "1INCH_USDT",
			},
			expectedFound: true,
		},
		"Bitstamp": {
			id: "Bitstamp",
			marketToTicker: map[types.MarketId]string{
				exchange_common.MARKET_BTC_USD: "BTC/USD",
				exchange_common.MARKET_ETH_USD: "ETH/USD",
			},
			expectedFound: true,
		},
		"Bybit": {
			id: "Bybit",
			marketToTicker: map[types.MarketId]string{
				exchange_common.MARKET_BTC_USD:  "BTCUSDT",
				exchange_common.MARKET_ETH_USD:  "ETHUSDT",
				exchange_common.MARKET_CRV_USD:  "CRVUSDT",
				exchange_common.MARKET_LTC_USD:  "LTCUSDT",
				exchange_common.MARKET_ATOM_USD: "ATOMUSDT",
				exchange_common.MARKET_UNI_USD:  "UNIUSDT",
				exchange_common.MARKET_NEAR_USD: "NEARUSDT",
				exchange_common.MARKET_COMP_USD: "COMPUSDT",
				exchange_common.MARKET_YFI_USD:  "YFIUSDT",
			},
			expectedFound: true,
		},
		"Crypto.com": {
			id: "CryptoCom",
			marketToTicker: map[types.MarketId]string{
				exchange_common.MARKET_BTC_USD:  "BTC_USD",
				exchange_common.MARKET_ETH_USD:  "ETH_USD",
				exchange_common.MARKET_LINK_USD: "LINK_USD",
			},
			expectedFound: true,
		},
		"Huobi": {
			id: "Huobi",
			marketToTicker: map[types.MarketId]string{
				exchange_common.MARKET_LINK_USD:  "linkusdt",
				exchange_common.MARKET_MATIC_USD: "maticusdt",
				exchange_common.MARKET_CRV_USD:   "crvusdt",
				exchange_common.MARKET_SOL_USD:   "solusdt",
				exchange_common.MARKET_ADA_USD:   "adausdt",
				exchange_common.MARKET_AVAX_USD:  "avaxusdt",
				exchange_common.MARKET_FIL_USD:   "filusdt",
				exchange_common.MARKET_AAVE_USD:  "aaveusdt",
				exchange_common.MARKET_LTC_USD:   "ltcusdt",
				exchange_common.MARKET_DOGE_USD:  "dogeusdt",
				exchange_common.MARKET_ICP_USD:   "icpusdt",
				exchange_common.MARKET_ATOM_USD:  "atomusdt",
				exchange_common.MARKET_DOT_USD:   "dotusdt",
				exchange_common.MARKET_XTZ_USD:   "xtzusdt",
				exchange_common.MARKET_UNI_USD:   "uniusdt",
				exchange_common.MARKET_BCH_USD:   "bchusdt",
				exchange_common.MARKET_EOS_USD:   "eosusdt",
				exchange_common.MARKET_TRX_USD:   "trxusdt",
				exchange_common.MARKET_ALGO_USD:  "algousdt",
				exchange_common.MARKET_NEAR_USD:  "nearusdt",
				exchange_common.MARKET_SNX_USD:   "snxusdt",
				exchange_common.MARKET_MKR_USD:   "mkrusdt",
				exchange_common.MARKET_SUSHI_USD: "sushiusdt",
				exchange_common.MARKET_ETC_USD:   "etcusdt",
				exchange_common.MARKET_1INCH_USD: "1inchusdt",
				exchange_common.MARKET_COMP_USD:  "compusdt",
				exchange_common.MARKET_ZRX_USD:   "zrxusdt",
				exchange_common.MARKET_YFI_USD:   "yfiusdt",
			},
			expectedFound: true,
		},
		"Kucoin": {
			id: "Kucoin",
			marketToTicker: map[types.MarketId]string{
				2:                                "LINK-USDT",
				exchange_common.MARKET_MATIC_USD: "MATIC-USDT",
				exchange_common.MARKET_CRV_USD:   "CRV-USDT",
				exchange_common.MARKET_SOL_USD:   "SOL-USDT",
				exchange_common.MARKET_ADA_USD:   "ADA-USDT",
				exchange_common.MARKET_AVAX_USD:  "AVAX-USDT",
				exchange_common.MARKET_FIL_USD:   "FIL-USDT",
				exchange_common.MARKET_AAVE_USD:  "AAVE-USDT",
				exchange_common.MARKET_LTC_USD:   "LTC-USDT",
				exchange_common.MARKET_DOGE_USD:  "DOGE-USDT",
				exchange_common.MARKET_ICP_USD:   "ICP-USDT",
				exchange_common.MARKET_ATOM_USD:  "ATOM-USDT",
				exchange_common.MARKET_DOT_USD:   "DOT-USDT",
				exchange_common.MARKET_ALGO_USD:  "ALGO-USDT",
				exchange_common.MARKET_NEAR_USD:  "NEAR-USDT",
				exchange_common.MARKET_SNX_USD:   "SNX-USDT",
				exchange_common.MARKET_MKR_USD:   "MKR-USDT",
				exchange_common.MARKET_XLM_USD:   "XLM-USDT",
				exchange_common.MARKET_XMR_USD:   "XMR-USDT",
				exchange_common.MARKET_1INCH_USD: "1INCH-USDT",
				exchange_common.MARKET_COMP_USD:  "COMP-USDT",
				exchange_common.MARKET_ZEC_USD:   "ZEC-USDT",
			},
			expectedFound: true,
		},
		"Okx": {
			id: "Okx",
			marketToTicker: map[types.MarketId]string{
				exchange_common.MARKET_BTC_USD:   "BTC-USDT",
				exchange_common.MARKET_ETH_USD:   "ETH-USDT",
				exchange_common.MARKET_LINK_USD:  "LINK-USDT",
				exchange_common.MARKET_MATIC_USD: "MATIC-USDT",
				exchange_common.MARKET_CRV_USD:   "CRV-USDT",
				exchange_common.MARKET_SOL_USD:   "SOL-USDT",
				exchange_common.MARKET_AVAX_USD:  "AVAX-USDT",
				exchange_common.MARKET_FIL_USD:   "FIL-USDT",
				exchange_common.MARKET_AAVE_USD:  "AAVE-USDT",
				exchange_common.MARKET_LTC_USD:   "LTC-USDT",
				exchange_common.MARKET_DOGE_USD:  "DOGE-USDT",
				exchange_common.MARKET_ICP_USD:   "ICP-USDT",
				exchange_common.MARKET_ATOM_USD:  "ATOM-USDT",
				exchange_common.MARKET_DOT_USD:   "DOT-USDT",
				exchange_common.MARKET_XTZ_USD:   "XTZ-USDT",
				exchange_common.MARKET_UNI_USD:   "UNI-USDT",
				exchange_common.MARKET_BCH_USD:   "BCH-USDT",
				exchange_common.MARKET_EOS_USD:   "EOS-USDT",
				exchange_common.MARKET_TRX_USD:   "TRX-USDT",
				exchange_common.MARKET_ALGO_USD:  "ALGO-USDT",
				exchange_common.MARKET_NEAR_USD:  "NEAR-USDT",
				exchange_common.MARKET_SNX_USD:   "SNX-USDT",
				exchange_common.MARKET_MKR_USD:   "MKR-USDT",
				exchange_common.MARKET_SUSHI_USD: "SUSHI-USDT",
				exchange_common.MARKET_XLM_USD:   "XLM-USDT",
				exchange_common.MARKET_XMR_USD:   "XMR-USDT",
				exchange_common.MARKET_ETC_USD:   "ETC-USDT",
				exchange_common.MARKET_1INCH_USD: "1INCH-USDT",
				exchange_common.MARKET_COMP_USD:  "COMP-USDT",
				exchange_common.MARKET_ZEC_USD:   "ZEC-USDT",
				exchange_common.MARKET_ZRX_USD:   "ZRX-USDT",
				exchange_common.MARKET_YFI_USD:   "YFI-USDT",
			},
			expectedFound: true,
		},
		"Mexc": {
			id: "Mexc",
			marketToTicker: map[types.MarketId]string{
				exchange_common.MARKET_UNI_USD: "UNI_USDT",
				exchange_common.MARKET_XMR_USD: "XMR_USDT",
			},
			expectedFound: true,
		},
		"Coinbase Pro": {
			id: "CoinbasePro",
			marketToTicker: map[types.MarketId]string{
				exchange_common.MARKET_BTC_USD:   "BTC-USD",
				exchange_common.MARKET_ETH_USD:   "ETH-USD",
				exchange_common.MARKET_LINK_USD:  "LINK-USD",
				exchange_common.MARKET_MATIC_USD: "MATIC-USD",
				exchange_common.MARKET_CRV_USD:   "CRV-USD",
				exchange_common.MARKET_SOL_USD:   "SOL-USD",
				exchange_common.MARKET_ADA_USD:   "ADA-USD",
				exchange_common.MARKET_FIL_USD:   "FIL-USD",
				exchange_common.MARKET_AAVE_USD:  "AAVE-USD",
				exchange_common.MARKET_LTC_USD:   "LTC-USD",
				exchange_common.MARKET_ICP_USD:   "ICP-USD",
				exchange_common.MARKET_ATOM_USD:  "ATOM-USD",
				exchange_common.MARKET_XTZ_USD:   "XTZ-USD",
				exchange_common.MARKET_UNI_USD:   "UNI-USD",
				exchange_common.MARKET_BCH_USD:   "BCH-USD",
				exchange_common.MARKET_EOS_USD:   "EOS-USD",
				exchange_common.MARKET_ALGO_USD:  "ALGO-USD",
				exchange_common.MARKET_NEAR_USD:  "NEAR-USD",
				exchange_common.MARKET_SNX_USD:   "SNX-USD",
				exchange_common.MARKET_MKR_USD:   "MKR-USD",
				exchange_common.MARKET_SUSHI_USD: "SUSHI-USD",
				exchange_common.MARKET_XLM_USD:   "XLM-USD",
				exchange_common.MARKET_ETC_USD:   "ETC-USD",
				exchange_common.MARKET_1INCH_USD: "1INCH-USD",
				exchange_common.MARKET_COMP_USD:  "COMP-USD",
				exchange_common.MARKET_ZEC_USD:   "ZEC-USD",
				exchange_common.MARKET_ZRX_USD:   "ZRX-USD",
				exchange_common.MARKET_YFI_USD:   "YFI-USD",
			},
			expectedFound: true,
		},
		"Test exchange": {
			id: "TestExchange",
			marketToTicker: map[types.MarketId]string{
				exchange_common.MARKET_BTC_USD:  "BTC-USD",
				exchange_common.MARKET_ETH_USD:  "ETH-USD",
				exchange_common.MARKET_LINK_USD: "LINK-USD",
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
				require.Equal(t, tc.marketToTicker, config.MarketToTicker)
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
	require.Len(t, configs, 33)
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
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			configs := GenerateExchangeConfigJson(StaticExchangeMarketConfig)
			actualExchangeConfigJson, ok := configs[tc.id]
			expectedExchangeConfigJson := pricefeed.ReadJsonTestFile(t, tc.expectedExchangeConfigJsonFile)
			require.True(t, ok)
			require.Equal(t, expectedExchangeConfigJson, actualExchangeConfigJson)
		})
	}
}
