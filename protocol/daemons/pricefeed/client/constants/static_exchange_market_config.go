package constants

import (
	"fmt"
	"github.com/dydxprotocol/v4/daemons/pricefeed/client/constants/exchange_common"
	"github.com/dydxprotocol/v4/daemons/pricefeed/client/types"
	"sort"
	"strings"
)

var (
	// StaticExchangeMarketConfig maps exchange feed ids to exchange market config. This map is used to generate
	// the exchange config json used by the genesis state. See `GenerateExchangeConfigJson` below.
	StaticExchangeMarketConfig = map[types.ExchangeId]*types.MutableExchangeMarketConfig{
		exchange_common.EXCHANGE_ID_BINANCE: {
			Id: exchange_common.EXCHANGE_ID_BINANCE,
			// example `symbols` parameter: ["BTCUSDT","BNBUSDT"]
			MarketToTicker: map[types.MarketId]string{
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
		},
		exchange_common.EXCHANGE_ID_BINANCE_US: {
			Id: exchange_common.EXCHANGE_ID_BINANCE_US,
			// example `symbols` parameter: ["BTCUSD","BNBUSD"]
			MarketToTicker: map[types.MarketId]string{
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
		},
		exchange_common.EXCHANGE_ID_BITFINEX: {
			Id: exchange_common.EXCHANGE_ID_BITFINEX,
			MarketToTicker: map[uint32]string{
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
		},
		exchange_common.EXCHANGE_ID_KRAKEN: {
			Id: exchange_common.EXCHANGE_ID_KRAKEN,
			MarketToTicker: map[types.MarketId]string{
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
		},
		exchange_common.EXCHANGE_ID_GATE: {
			Id: exchange_common.EXCHANGE_ID_GATE,
			MarketToTicker: map[types.MarketId]string{
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
		},
		exchange_common.EXCHANGE_ID_BITSTAMP: {
			Id: exchange_common.EXCHANGE_ID_BITSTAMP,
			MarketToTicker: map[types.MarketId]string{
				exchange_common.MARKET_BTC_USD: "BTC/USD",
				exchange_common.MARKET_ETH_USD: "ETH/USD",
			},
		},
		exchange_common.EXCHANGE_ID_BYBIT: {
			Id: exchange_common.EXCHANGE_ID_BYBIT,
			MarketToTicker: map[types.MarketId]string{
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
		},
		exchange_common.EXCHANGE_ID_CRYPTO_COM: {
			Id: exchange_common.EXCHANGE_ID_CRYPTO_COM,
			MarketToTicker: map[types.MarketId]string{
				exchange_common.MARKET_BTC_USD:  "BTC_USD",
				exchange_common.MARKET_ETH_USD:  "ETH_USD",
				exchange_common.MARKET_LINK_USD: "LINK_USD",
			},
		},
		exchange_common.EXCHANGE_ID_HUOBI: {
			Id: exchange_common.EXCHANGE_ID_HUOBI,
			MarketToTicker: map[types.MarketId]string{
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
		},
		exchange_common.EXCHANGE_ID_KUCOIN: {
			Id: exchange_common.EXCHANGE_ID_KUCOIN,
			MarketToTicker: map[types.MarketId]string{
				exchange_common.MARKET_LINK_USD:  "LINK-USDT",
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
		},
		exchange_common.EXCHANGE_ID_OKX: {
			Id: exchange_common.EXCHANGE_ID_OKX,
			MarketToTicker: map[types.MarketId]string{
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
		},
		exchange_common.EXCHANGE_ID_MEXC: {
			Id: exchange_common.EXCHANGE_ID_MEXC,
			MarketToTicker: map[types.MarketId]string{
				exchange_common.MARKET_UNI_USD: "UNI_USDT",
				exchange_common.MARKET_XMR_USD: "XMR_USDT",
			},
		},
		exchange_common.EXCHANGE_ID_COINBASE_PRO: {
			Id: exchange_common.EXCHANGE_ID_COINBASE_PRO,
			MarketToTicker: map[types.MarketId]string{
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
		},
		exchange_common.EXCHANGE_ID_TEST_EXCHANGE: {
			Id: exchange_common.EXCHANGE_ID_TEST_EXCHANGE,
			MarketToTicker: map[types.MarketId]string{
				exchange_common.MARKET_BTC_USD:  "BTC-USD",
				exchange_common.MARKET_ETH_USD:  "ETH-USD",
				exchange_common.MARKET_LINK_USD: "LINK-USD",
			},
		},
	}
)

// GenerateExchangeConfigJson generates the exchange config json for each market based on the contents
// of an exchangeToExchangeConfig map. For the default exchange configs, pass in the
// `StaticExchangeMarketConfig` map above as the argument.
func GenerateExchangeConfigJson(
	exchangeToExchangeConfig map[types.ExchangeId]*types.MutableExchangeMarketConfig,
) (
	marketToExchangeConfigJson map[types.MarketId]string,
) {
	// marketToExchangeConfigJson maps markets to a map of exchange, exchange market config. This
	// is used to generate the exchange config json from the above map that is keyed by exchange id.
	// We keep this intermediate map so that we can sort the exchange configs for each market by exchange name
	// in order to make the output deterministic.
	marketToExchangeMarketConfigs := make(map[types.MarketId]map[string]string)

	// Generate the market-specific exchange config for each market, exchange.
	for id, exchangeConfig := range exchangeToExchangeConfig {
		// Skip config for the test exchange.
		if id == exchange_common.EXCHANGE_ID_TEST_EXCHANGE {
			continue
		}
		for marketId, ticker := range exchangeConfig.MarketToTicker {
			marketExchangeConfigs, ok := marketToExchangeMarketConfigs[marketId]
			if !ok {
				marketToExchangeMarketConfigs[marketId] = map[string]string{}
				marketExchangeConfigs = marketToExchangeMarketConfigs[marketId]
			}

			// Escape double quotes in the ticker.
			if strings.Contains(ticker, `"`) {
				ticker = strings.ReplaceAll(ticker, `"`, `\"`)
			}

			marketExchangeConfigs[id] = fmt.Sprintf(`{"exchangeName":"%v","ticker":"%v"}`, id, ticker)
		}
	}

	// Initialize the output map.
	marketToExchangeConfigJson = make(map[types.MarketId]string, len(marketToExchangeMarketConfigs))

	// Generate the output map of market to exchange config json.
	for marketId, exchangeToConfigs := range marketToExchangeMarketConfigs {
		// Sort output exchange configs by exchange name in order to make output deterministic.
		exchangeNames := make([]string, 0, len(exchangeToConfigs))

		// 1. Generate sorted list of exchange names.
		for name := range exchangeToConfigs {
			exchangeNames = append(exchangeNames, name)
		}
		sort.Strings(exchangeNames)

		// 2. Generate sorted list of exchange configs by exchange name.
		sortedExchangeConfigs := make([]string, 0, len(exchangeNames))
		for _, exchangeName := range exchangeNames {
			sortedExchangeConfigs = append(sortedExchangeConfigs, exchangeToConfigs[exchangeName])
		}
		// 3. Generate the output json for the market, sorted by exchange name.
		marketToExchangeConfigJson[marketId] = fmt.Sprintf(`{"exchanges":[%v]}`, strings.Join(sortedExchangeConfigs, ","))
	}
	return marketToExchangeConfigJson
}
