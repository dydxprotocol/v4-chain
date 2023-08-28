package constants

import (
	"encoding/json"
	"fmt"
	"sort"

	"github.com/dydxprotocol/v4-chain/protocol/daemons/pricefeed/client/constants/exchange_common"
	"github.com/dydxprotocol/v4-chain/protocol/daemons/pricefeed/client/types"
)

var (
	// StaticExchangeMarketConfig maps exchange feed ids to exchange market config. This map is used to generate
	// the exchange config json used by the genesis state. See `GenerateExchangeConfigJson` below.
	StaticExchangeMarketConfig = map[types.ExchangeId]*types.MutableExchangeMarketConfig{
		exchange_common.EXCHANGE_ID_BINANCE: {
			Id: exchange_common.EXCHANGE_ID_BINANCE,
			// example `symbols` parameter: ["BTCUSDT","BNBUSDT"]
			MarketToMarketConfig: map[types.MarketId]types.MarketConfig{
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
				exchange_common.MARKET_LTC_USD: {
					Ticker:         `"LTCUSDT"`,
					AdjustByMarket: newMarketIdWithValue(exchange_common.MARKET_USDT_USD),
				},
				exchange_common.MARKET_DOGE_USD: {
					Ticker:         `"DOGEUSDT"`,
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
				exchange_common.MARKET_UNI_USD: {
					Ticker:         `"UNIUSDT"`,
					AdjustByMarket: newMarketIdWithValue(exchange_common.MARKET_USDT_USD),
				},
				exchange_common.MARKET_BCH_USD: {
					Ticker:         `"BCHUSDT"`,
					AdjustByMarket: newMarketIdWithValue(exchange_common.MARKET_USDT_USD),
				},
				exchange_common.MARKET_TRX_USD: {
					Ticker:         `"TRXUSDT"`,
					AdjustByMarket: newMarketIdWithValue(exchange_common.MARKET_USDT_USD),
				},
				exchange_common.MARKET_NEAR_USD: {
					Ticker:         `"NEARUSDT"`,
					AdjustByMarket: newMarketIdWithValue(exchange_common.MARKET_USDT_USD),
				},
				exchange_common.MARKET_MKR_USD: {
					Ticker:         `"MKRUSDT"`,
					AdjustByMarket: newMarketIdWithValue(exchange_common.MARKET_USDT_USD),
				},
				exchange_common.MARKET_XLM_USD: {
					Ticker:         `"XLMUSDT"`,
					AdjustByMarket: newMarketIdWithValue(exchange_common.MARKET_USDT_USD),
				},
				exchange_common.MARKET_ETC_USD: {
					Ticker:         `"ETCUSDT"`,
					AdjustByMarket: newMarketIdWithValue(exchange_common.MARKET_USDT_USD),
				},
				exchange_common.MARKET_COMP_USD: {
					Ticker:         `"COMPUSDT"`,
					AdjustByMarket: newMarketIdWithValue(exchange_common.MARKET_USDT_USD),
				},
				exchange_common.MARKET_APE_USD: {
					Ticker:         `"APEUSDT"`,
					AdjustByMarket: newMarketIdWithValue(exchange_common.MARKET_USDT_USD),
				},
				exchange_common.MARKET_APT_USD: {
					Ticker:         `"APTUSDT"`,
					AdjustByMarket: newMarketIdWithValue(exchange_common.MARKET_USDT_USD),
				},
				exchange_common.MARKET_ARB_USD: {
					Ticker:         `"ARBUSDT"`,
					AdjustByMarket: newMarketIdWithValue(exchange_common.MARKET_USDT_USD),
				},
				exchange_common.MARKET_LDO_USD: {
					Ticker:         `"LDOUSDT"`,
					AdjustByMarket: newMarketIdWithValue(exchange_common.MARKET_USDT_USD),
				},
				exchange_common.MARKET_OP_USD: {
					Ticker:         `"OPUSDT"`,
					AdjustByMarket: newMarketIdWithValue(exchange_common.MARKET_USDT_USD),
				},
				exchange_common.MARKET_PEPE_USD: {
					Ticker:         `"PEPEUSDT"`,
					AdjustByMarket: newMarketIdWithValue(exchange_common.MARKET_USDT_USD),
				},
				exchange_common.MARKET_SEI_USD: {
					Ticker:         `"SEIUSDT"`,
					AdjustByMarket: newMarketIdWithValue(exchange_common.MARKET_USDT_USD),
				},
				exchange_common.MARKET_SHIB_USD: {
					Ticker:         `"SHIBUSDT"`,
					AdjustByMarket: newMarketIdWithValue(exchange_common.MARKET_USDT_USD),
				},
				exchange_common.MARKET_SUI_USD: {
					Ticker:         `"SUIUSDT"`,
					AdjustByMarket: newMarketIdWithValue(exchange_common.MARKET_USDT_USD),
				},
				exchange_common.MARKET_WLD_USD: {
					Ticker:         `"WLDUSDT"`,
					AdjustByMarket: newMarketIdWithValue(exchange_common.MARKET_USDT_USD),
				},
				exchange_common.MARKET_XRP_USD: {
					Ticker:         `"XRPUSDT"`,
					AdjustByMarket: newMarketIdWithValue(exchange_common.MARKET_USDT_USD),
				},
				exchange_common.MARKET_USDT_USD: {
					Ticker:         `"BTCUSDT"`, // Adjusted with BTC index price.
					Invert:         true,
					AdjustByMarket: newMarketIdWithValue(exchange_common.MARKET_BTC_USD),
				},
			},
		},
		exchange_common.EXCHANGE_ID_BINANCE_US: {
			Id: exchange_common.EXCHANGE_ID_BINANCE_US,
			// example `symbols` parameter: ["BTCUSD","BNBUSD"]
			MarketToMarketConfig: map[types.MarketId]types.MarketConfig{
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
				exchange_common.MARKET_LTC_USD: {
					Ticker: `"LTCUSD"`,
				},
				exchange_common.MARKET_DOGE_USD: {
					Ticker: `"DOGEUSD"`,
				},
				exchange_common.MARKET_ATOM_USD: {
					Ticker: `"ATOMUSD"`,
				},
				exchange_common.MARKET_DOT_USD: {
					Ticker: `"DOTUSD"`,
				},
				exchange_common.MARKET_UNI_USD: {
					Ticker: `"UNIUSD"`,
				},
				exchange_common.MARKET_BCH_USD: {
					Ticker: `"BCHUSD"`,
				},
				exchange_common.MARKET_NEAR_USD: {
					Ticker: `"NEARUSD"`,
				},
				exchange_common.MARKET_MKR_USD: {
					Ticker: `"MKRUSD"`,
				},
				exchange_common.MARKET_XLM_USD: {
					Ticker: `"XLMUSD"`,
				},
				exchange_common.MARKET_ETC_USD: {
					Ticker: `"ETCUSD"`,
				},
				exchange_common.MARKET_COMP_USD: {
					Ticker: `"COMPUSD"`,
				},
				exchange_common.MARKET_APE_USD: {
					Ticker: `"APEUSD"`,
				},
				exchange_common.MARKET_APT_USD: {
					Ticker: `"APTUSD"`,
				},
				exchange_common.MARKET_LDO_USD: {
					Ticker: `"LDOUSD"`,
				},
				exchange_common.MARKET_SHIB_USD: {
					Ticker: `"SHIBUSD"`,
				},
				exchange_common.MARKET_XRP_USD: {
					Ticker: `"XRPUSD"`,
				},
				exchange_common.MARKET_USDT_USD: {
					Ticker: `"USDTUSD"`,
				},
			},
		},
		exchange_common.EXCHANGE_ID_BITFINEX: {
			Id: exchange_common.EXCHANGE_ID_BITFINEX,
			// Note: we treat all Bitfinex pairs as USDT.
			MarketToMarketConfig: map[types.MarketId]types.MarketConfig{
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
				exchange_common.MARKET_TRX_USD: {
					Ticker:         "tTRXUSD",
					AdjustByMarket: newMarketIdWithValue(exchange_common.MARKET_USDT_USD),
				},
				exchange_common.MARKET_MKR_USD: {
					Ticker:         "tMKRUSD",
					AdjustByMarket: newMarketIdWithValue(exchange_common.MARKET_USDT_USD),
				},
				exchange_common.MARKET_XLM_USD: {
					Ticker:         "tXLMUSD",
					AdjustByMarket: newMarketIdWithValue(exchange_common.MARKET_USDT_USD),
				},
				exchange_common.MARKET_APE_USD: {
					Ticker:         "tAPEUSD",
					AdjustByMarket: newMarketIdWithValue(exchange_common.MARKET_USDT_USD),
				},
				exchange_common.MARKET_APT_USD: {
					Ticker:         "tAPTUSD",
					AdjustByMarket: newMarketIdWithValue(exchange_common.MARKET_USDT_USD),
				},
				exchange_common.MARKET_ARB_USD: {
					Ticker:         "tARBUSD",
					AdjustByMarket: newMarketIdWithValue(exchange_common.MARKET_USDT_USD),
				},
				exchange_common.MARKET_BLUR_USD: {
					Ticker:         "tBLUR:USD",
					AdjustByMarket: newMarketIdWithValue(exchange_common.MARKET_USDT_USD),
				},
				exchange_common.MARKET_FIL_USD: {
					Ticker:         "tFILUSD",
					AdjustByMarket: newMarketIdWithValue(exchange_common.MARKET_USDT_USD),
				},
				exchange_common.MARKET_LDO_USD: {
					Ticker:         "LDOUSD",
					AdjustByMarket: newMarketIdWithValue(exchange_common.MARKET_USDT_USD),
				},
				exchange_common.MARKET_PEPE_USD: {
					Ticker:         "PEPE:USD",
					AdjustByMarket: newMarketIdWithValue(exchange_common.MARKET_USDT_USD),
				},
				exchange_common.MARKET_SHIB_USD: {
					Ticker:         "tSHIB:USD",
					AdjustByMarket: newMarketIdWithValue(exchange_common.MARKET_USDT_USD),
				},
				exchange_common.MARKET_SUI_USD: {
					Ticker:         "SUIUSD",
					AdjustByMarket: newMarketIdWithValue(exchange_common.MARKET_USDT_USD),
				},
				exchange_common.MARKET_XRP_USD: {
					Ticker:         "tXRPUSD",
					AdjustByMarket: newMarketIdWithValue(exchange_common.MARKET_USDT_USD),
				},
				exchange_common.MARKET_USDT_USD: {
					Ticker:         "tUSTUSD",
					AdjustByMarket: newMarketIdWithValue(exchange_common.MARKET_USDT_USD),
				},
			},
		},
		exchange_common.EXCHANGE_ID_KRAKEN: {
			Id: exchange_common.EXCHANGE_ID_KRAKEN,
			MarketToMarketConfig: map[types.MarketId]types.MarketConfig{
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
				exchange_common.MARKET_LTC_USD: {
					Ticker: "XLTCZUSD",
				},
				exchange_common.MARKET_ATOM_USD: {
					Ticker: "ATOMUSD",
				},
				exchange_common.MARKET_DOT_USD: {
					Ticker: "DOTUSD",
				},
				exchange_common.MARKET_UNI_USD: {
					Ticker: "UNIUSD",
				},
				exchange_common.MARKET_BCH_USD: {
					Ticker: "BCHUSD",
				},
				exchange_common.MARKET_XLM_USD: {
					Ticker: "XXLMZUSD",
				},
				exchange_common.MARKET_ETC_USD: {
					Ticker: "XETCZUSD",
				},
				exchange_common.MARKET_COMP_USD: {
					Ticker: "COMPUSD",
				},
				exchange_common.MARKET_APE_USD: {
					Ticker: "APEUSD",
				},
				exchange_common.MARKET_APT_USD: {
					Ticker: "APTUSD",
				},
				exchange_common.MARKET_ARB_USD: {
					Ticker: "ARBUSD",
				},
				exchange_common.MARKET_BLUR_USD: {
					Ticker: "BLURUSD",
				},
				exchange_common.MARKET_LDO_USD: {
					Ticker: "LDOUSD",
				},
				exchange_common.MARKET_PEPE_USD: {
					Ticker: "PEPEUSD",
				},
				exchange_common.MARKET_SHIB_USD: {
					Ticker: "SHIBUSD",
				},
				exchange_common.MARKET_SUI_USD: {
					Ticker: "SUIUSD",
				},
				exchange_common.MARKET_XRP_USD: {
					Ticker: "XXRPZUSD",
				},
				exchange_common.MARKET_USDT_USD: {
					Ticker: "USDTZUSD",
				},
			},
		},
		exchange_common.EXCHANGE_ID_GATE: {
			Id: exchange_common.EXCHANGE_ID_GATE,
			MarketToMarketConfig: map[types.MarketId]types.MarketConfig{
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
				exchange_common.MARKET_DOT_USD: {
					Ticker:         "DOT_USDT",
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
				exchange_common.MARKET_XLM_USD: {
					Ticker:         "XLM_USDT",
					AdjustByMarket: newMarketIdWithValue(exchange_common.MARKET_USDT_USD),
				},
				exchange_common.MARKET_ETC_USD: {
					Ticker:         "ETC_USDT",
					AdjustByMarket: newMarketIdWithValue(exchange_common.MARKET_USDT_USD),
				},
				exchange_common.MARKET_APE_USD: {
					Ticker:         "APE_USDT",
					AdjustByMarket: newMarketIdWithValue(exchange_common.MARKET_USDT_USD),
				},
				exchange_common.MARKET_APT_USD: {
					Ticker:         "APT_USDT",
					AdjustByMarket: newMarketIdWithValue(exchange_common.MARKET_USDT_USD),
				},
				exchange_common.MARKET_ARB_USD: {
					Ticker:         "ARB_USDT",
					AdjustByMarket: newMarketIdWithValue(exchange_common.MARKET_USDT_USD),
				},
				exchange_common.MARKET_BLUR_USD: {
					Ticker:         "BLUR_USDT",
					AdjustByMarket: newMarketIdWithValue(exchange_common.MARKET_USDT_USD),
				},
				exchange_common.MARKET_FIL_USD: {
					Ticker:         "FIL_USDT",
					AdjustByMarket: newMarketIdWithValue(exchange_common.MARKET_USDT_USD),
				},
				exchange_common.MARKET_LDO_USD: {
					Ticker:         "LDO_USDT",
					AdjustByMarket: newMarketIdWithValue(exchange_common.MARKET_USDT_USD),
				},
				exchange_common.MARKET_OP_USD: {
					Ticker:         "OP_USDT",
					AdjustByMarket: newMarketIdWithValue(exchange_common.MARKET_USDT_USD),
				},
				exchange_common.MARKET_PEPE_USD: {
					Ticker:         "PEPE_USDT",
					AdjustByMarket: newMarketIdWithValue(exchange_common.MARKET_USDT_USD),
				},
				exchange_common.MARKET_SEI_USD: {
					Ticker:         "SEI_USDT",
					AdjustByMarket: newMarketIdWithValue(exchange_common.MARKET_USDT_USD),
				},
				exchange_common.MARKET_SHIB_USD: {
					Ticker:         "SHIB_USDT",
					AdjustByMarket: newMarketIdWithValue(exchange_common.MARKET_USDT_USD),
				},
				exchange_common.MARKET_SUI_USD: {
					Ticker:         "SUI_USDT",
					AdjustByMarket: newMarketIdWithValue(exchange_common.MARKET_USDT_USD),
				},
				exchange_common.MARKET_WLD_USD: {
					Ticker:         "WLD_USDT",
					AdjustByMarket: newMarketIdWithValue(exchange_common.MARKET_USDT_USD),
				},
				exchange_common.MARKET_XRP_USD: {
					Ticker:         "XRP_USDT",
					AdjustByMarket: newMarketIdWithValue(exchange_common.MARKET_USDT_USD),
				},
				exchange_common.MARKET_USDT_USD: {
					Ticker: "USDT_USD",
				},
			},
		},
		exchange_common.EXCHANGE_ID_BITSTAMP: {
			Id: exchange_common.EXCHANGE_ID_BITSTAMP,
			MarketToMarketConfig: map[types.MarketId]types.MarketConfig{
				exchange_common.MARKET_BTC_USD: {
					Ticker: "BTC/USD",
				},
				exchange_common.MARKET_ETH_USD: {
					Ticker: "ETH/USD",
				},
				exchange_common.MARKET_APE_USD: {
					Ticker: "APE/USD",
				},
				exchange_common.MARKET_LDO_USD: {
					Ticker: "LDO/USD",
				},
				exchange_common.MARKET_SHIB_USD: {
					Ticker: "SHIB/USD",
				},
				exchange_common.MARKET_SUI_USD: {
					Ticker: "SUI/USD",
				},
				exchange_common.MARKET_XLM_USD: {
					Ticker: "XLM/USD",
				},
				exchange_common.MARKET_XRP_USD: {
					Ticker: "XRP/USD",
				},
				exchange_common.MARKET_USDT_USD: {
					Ticker: "USDT/USD",
				},
			},
		},
		exchange_common.EXCHANGE_ID_BYBIT: {
			Id: exchange_common.EXCHANGE_ID_BYBIT,
			MarketToMarketConfig: map[types.MarketId]types.MarketConfig{
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
				exchange_common.MARKET_APE_USD: {
					Ticker:         "APEUSDT",
					AdjustByMarket: newMarketIdWithValue(exchange_common.MARKET_USDT_USD),
				},
				exchange_common.MARKET_APT_USD: {
					Ticker:         "APTUSDT",
					AdjustByMarket: newMarketIdWithValue(exchange_common.MARKET_USDT_USD),
				},
				exchange_common.MARKET_ARB_USD: {
					Ticker:         "ARBUSDT",
					AdjustByMarket: newMarketIdWithValue(exchange_common.MARKET_USDT_USD),
				},
				exchange_common.MARKET_BLUR_USD: {
					Ticker:         "BLURUSDT",
					AdjustByMarket: newMarketIdWithValue(exchange_common.MARKET_USDT_USD),
				},
				exchange_common.MARKET_FIL_USD: {
					Ticker:         "FILUSDT",
					AdjustByMarket: newMarketIdWithValue(exchange_common.MARKET_USDT_USD),
				},
				exchange_common.MARKET_LDO_USD: {
					Ticker:         "LDOUSDT",
					AdjustByMarket: newMarketIdWithValue(exchange_common.MARKET_USDT_USD),
				},
				exchange_common.MARKET_OP_USD: {
					Ticker:         "OPUSDT",
					AdjustByMarket: newMarketIdWithValue(exchange_common.MARKET_USDT_USD),
				},
				exchange_common.MARKET_PEPE_USD: {
					Ticker:         "PEPEUSDT",
					AdjustByMarket: newMarketIdWithValue(exchange_common.MARKET_USDT_USD),
				},
				exchange_common.MARKET_SEI_USD: {
					Ticker:         "SEIUSDT",
					AdjustByMarket: newMarketIdWithValue(exchange_common.MARKET_USDT_USD),
				},
				exchange_common.MARKET_SHIB_USD: {
					Ticker:         "SHIBUSDT",
					AdjustByMarket: newMarketIdWithValue(exchange_common.MARKET_USDT_USD),
				},
				exchange_common.MARKET_SUI_USD: {
					Ticker:         "SUIUSDT",
					AdjustByMarket: newMarketIdWithValue(exchange_common.MARKET_USDT_USD),
				},
				exchange_common.MARKET_WLD_USD: {
					Ticker:         "WLDUSDT",
					AdjustByMarket: newMarketIdWithValue(exchange_common.MARKET_USDT_USD),
				},
				exchange_common.MARKET_XLM_USD: {
					Ticker:         "XLMUSDT",
					AdjustByMarket: newMarketIdWithValue(exchange_common.MARKET_USDT_USD),
				},
				exchange_common.MARKET_XRP_USD: {
					Ticker:         "XRPUSDT",
					AdjustByMarket: newMarketIdWithValue(exchange_common.MARKET_USDT_USD),
				},
				exchange_common.MARKET_USDT_USD: {
					Ticker: "USDCUSDT",
					Invert: true,
				},
			},
		},
		exchange_common.EXCHANGE_ID_CRYPTO_COM: {
			Id: exchange_common.EXCHANGE_ID_CRYPTO_COM,
			MarketToMarketConfig: map[types.MarketId]types.MarketConfig{
				exchange_common.MARKET_BTC_USD: {
					Ticker: "BTC_USD",
				},
				exchange_common.MARKET_ETH_USD: {
					Ticker: "ETH_USD",
				},
				exchange_common.MARKET_LINK_USD: {
					Ticker: "LINK_USD",
				},
				exchange_common.MARKET_SHIB_USD: {
					Ticker: "SHIB_USD",
				},
				exchange_common.MARKET_XRP_USD: {
					Ticker: "XRP_USD",
				},
				exchange_common.MARKET_XLM_USD: {
					Ticker: "XLM_USD",
				},
				exchange_common.MARKET_USDT_USD: {
					Ticker: "USDT_USD",
				},
			},
		},
		exchange_common.EXCHANGE_ID_HUOBI: {
			Id: exchange_common.EXCHANGE_ID_HUOBI,
			MarketToMarketConfig: map[types.MarketId]types.MarketConfig{
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
				exchange_common.MARKET_LTC_USD: {
					Ticker:         "ltcusdt",
					AdjustByMarket: newMarketIdWithValue(exchange_common.MARKET_USDT_USD),
				},
				exchange_common.MARKET_DOGE_USD: {
					Ticker:         "dogeusdt",
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
				exchange_common.MARKET_UNI_USD: {
					Ticker:         "uniusdt",
					AdjustByMarket: newMarketIdWithValue(exchange_common.MARKET_USDT_USD),
				},
				exchange_common.MARKET_BCH_USD: {
					Ticker:         "bchusdt",
					AdjustByMarket: newMarketIdWithValue(exchange_common.MARKET_USDT_USD),
				},
				exchange_common.MARKET_TRX_USD: {
					Ticker:         "trxusdt",
					AdjustByMarket: newMarketIdWithValue(exchange_common.MARKET_USDT_USD),
				},
				exchange_common.MARKET_NEAR_USD: {
					Ticker:         "nearusdt",
					AdjustByMarket: newMarketIdWithValue(exchange_common.MARKET_USDT_USD),
				},
				exchange_common.MARKET_MKR_USD: {
					Ticker:         "mkrusdt",
					AdjustByMarket: newMarketIdWithValue(exchange_common.MARKET_USDT_USD),
				},
				exchange_common.MARKET_ETC_USD: {
					Ticker:         "etcusdt",
					AdjustByMarket: newMarketIdWithValue(exchange_common.MARKET_USDT_USD),
				},
				exchange_common.MARKET_COMP_USD: {
					Ticker:         "compusdt",
					AdjustByMarket: newMarketIdWithValue(exchange_common.MARKET_USDT_USD),
				},
				exchange_common.MARKET_APE_USD: {
					Ticker:         "apeusdt",
					AdjustByMarket: newMarketIdWithValue(exchange_common.MARKET_USDT_USD),
				},
				exchange_common.MARKET_APT_USD: {
					Ticker:         "aptusdt",
					AdjustByMarket: newMarketIdWithValue(exchange_common.MARKET_USDT_USD),
				},
				exchange_common.MARKET_ARB_USD: {
					Ticker:         "arbusdt",
					AdjustByMarket: newMarketIdWithValue(exchange_common.MARKET_USDT_USD),
				},
				exchange_common.MARKET_BLUR_USD: {
					Ticker:         "blurusdt",
					AdjustByMarket: newMarketIdWithValue(exchange_common.MARKET_USDT_USD),
				},
				exchange_common.MARKET_LDO_USD: {
					Ticker:         "ldousdt",
					AdjustByMarket: newMarketIdWithValue(exchange_common.MARKET_USDT_USD),
				},
				exchange_common.MARKET_OP_USD: {
					Ticker:         "opusdt",
					AdjustByMarket: newMarketIdWithValue(exchange_common.MARKET_USDT_USD),
				},
				exchange_common.MARKET_PEPE_USD: {
					Ticker:         "pepeusdt",
					AdjustByMarket: newMarketIdWithValue(exchange_common.MARKET_USDT_USD),
				},
				exchange_common.MARKET_SEI_USD: {
					Ticker:         "seiusdt",
					AdjustByMarket: newMarketIdWithValue(exchange_common.MARKET_USDT_USD),
				},
				exchange_common.MARKET_SHIB_USD: {
					Ticker:         "shibusdt",
					AdjustByMarket: newMarketIdWithValue(exchange_common.MARKET_USDT_USD),
				},
				exchange_common.MARKET_SUI_USD: {
					Ticker:         "suiusdt",
					AdjustByMarket: newMarketIdWithValue(exchange_common.MARKET_USDT_USD),
				},
				exchange_common.MARKET_WLD_USD: {
					Ticker:         "wldusdt",
					AdjustByMarket: newMarketIdWithValue(exchange_common.MARKET_USDT_USD),
				},
				exchange_common.MARKET_XLM_USD: {
					Ticker:         "xlmusdt",
					AdjustByMarket: newMarketIdWithValue(exchange_common.MARKET_USDT_USD),
				},
				exchange_common.MARKET_XRP_USD: {
					Ticker:         "xrpusdt",
					AdjustByMarket: newMarketIdWithValue(exchange_common.MARKET_USDT_USD),
				},
				exchange_common.MARKET_USDT_USD: {
					Ticker:         "ethusdt", // Adjusted with ETH index price.
					AdjustByMarket: newMarketIdWithValue(exchange_common.MARKET_ETH_USD),
					Invert:         true,
				},
			},
		},
		exchange_common.EXCHANGE_ID_KUCOIN: {
			Id: exchange_common.EXCHANGE_ID_KUCOIN,
			MarketToMarketConfig: map[types.MarketId]types.MarketConfig{
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
				exchange_common.MARKET_LTC_USD: {
					Ticker:         "LTC-USDT",
					AdjustByMarket: newMarketIdWithValue(exchange_common.MARKET_USDT_USD),
				},
				exchange_common.MARKET_DOGE_USD: {
					Ticker:         "DOGE-USDT",
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
				exchange_common.MARKET_NEAR_USD: {
					Ticker:         "NEAR-USDT",
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
				exchange_common.MARKET_COMP_USD: {
					Ticker:         "COMP-USDT",
					AdjustByMarket: newMarketIdWithValue(exchange_common.MARKET_USDT_USD),
				},
				exchange_common.MARKET_APE_USD: {
					Ticker:         "APE-USDT",
					AdjustByMarket: newMarketIdWithValue(exchange_common.MARKET_USDT_USD),
				},
				exchange_common.MARKET_APT_USD: {
					Ticker:         "APT-USDT",
					AdjustByMarket: newMarketIdWithValue(exchange_common.MARKET_USDT_USD),
				},
				exchange_common.MARKET_ARB_USD: {
					Ticker:         "ARB-USDT",
					AdjustByMarket: newMarketIdWithValue(exchange_common.MARKET_USDT_USD),
				},
				exchange_common.MARKET_BLUR_USD: {
					Ticker:         "BLUR-USDT",
					AdjustByMarket: newMarketIdWithValue(exchange_common.MARKET_USDT_USD),
				},
				exchange_common.MARKET_LDO_USD: {
					Ticker:         "LDO-USDT",
					AdjustByMarket: newMarketIdWithValue(exchange_common.MARKET_USDT_USD),
				},
				exchange_common.MARKET_OP_USD: {
					Ticker:         "OP-USDT",
					AdjustByMarket: newMarketIdWithValue(exchange_common.MARKET_USDT_USD),
				},
				exchange_common.MARKET_PEPE_USD: {
					Ticker:         "PEPE-USDT",
					AdjustByMarket: newMarketIdWithValue(exchange_common.MARKET_USDT_USD),
				},
				exchange_common.MARKET_SEI_USD: {
					Ticker:         "SEI-USDT",
					AdjustByMarket: newMarketIdWithValue(exchange_common.MARKET_USDT_USD),
				},
				exchange_common.MARKET_SHIB_USD: {
					Ticker:         "SHIB-USDT",
					AdjustByMarket: newMarketIdWithValue(exchange_common.MARKET_USDT_USD),
				},
				exchange_common.MARKET_SUI_USD: {
					Ticker:         "SUI-USDT",
					AdjustByMarket: newMarketIdWithValue(exchange_common.MARKET_USDT_USD),
				},
				exchange_common.MARKET_WLD_USD: {
					Ticker:         "WLD-USDT",
					AdjustByMarket: newMarketIdWithValue(exchange_common.MARKET_USDT_USD),
				},
				exchange_common.MARKET_XRP_USD: {
					Ticker:         "XRP-USDT",
					AdjustByMarket: newMarketIdWithValue(exchange_common.MARKET_USDT_USD),
				},

				exchange_common.MARKET_USDT_USD: {
					Ticker:         "BTC-USDT", // Adjusted with BTC index price.
					AdjustByMarket: newMarketIdWithValue(exchange_common.MARKET_BTC_USD),
					Invert:         true,
				},
			},
		},
		exchange_common.EXCHANGE_ID_OKX: {
			Id: exchange_common.EXCHANGE_ID_OKX,
			MarketToMarketConfig: map[types.MarketId]types.MarketConfig{
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
				exchange_common.MARKET_LTC_USD: {
					Ticker:         "LTC-USDT",
					AdjustByMarket: newMarketIdWithValue(exchange_common.MARKET_USDT_USD),
				},
				exchange_common.MARKET_DOGE_USD: {
					Ticker:         "DOGE-USDT",
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
				exchange_common.MARKET_UNI_USD: {
					Ticker:         "UNI-USDT",
					AdjustByMarket: newMarketIdWithValue(exchange_common.MARKET_USDT_USD),
				},
				exchange_common.MARKET_BCH_USD: {
					Ticker:         "BCH-USDT",
					AdjustByMarket: newMarketIdWithValue(exchange_common.MARKET_USDT_USD),
				},
				exchange_common.MARKET_TRX_USD: {
					Ticker:         "TRX-USDT",
					AdjustByMarket: newMarketIdWithValue(exchange_common.MARKET_USDT_USD),
				},
				exchange_common.MARKET_NEAR_USD: {
					Ticker:         "NEAR-USDT",
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
				exchange_common.MARKET_ETC_USD: {
					Ticker:         "ETC-USDT",
					AdjustByMarket: newMarketIdWithValue(exchange_common.MARKET_USDT_USD),
				},
				exchange_common.MARKET_COMP_USD: {
					Ticker:         "COMP-USDT",
					AdjustByMarket: newMarketIdWithValue(exchange_common.MARKET_USDT_USD),
				},
				exchange_common.MARKET_APE_USD: {
					Ticker:         "APE-USDT",
					AdjustByMarket: newMarketIdWithValue(exchange_common.MARKET_USDT_USD),
				},
				exchange_common.MARKET_APT_USD: {
					Ticker:         "APT-USDT",
					AdjustByMarket: newMarketIdWithValue(exchange_common.MARKET_USDT_USD),
				},
				exchange_common.MARKET_ARB_USD: {
					Ticker:         "ARB-USDT",
					AdjustByMarket: newMarketIdWithValue(exchange_common.MARKET_USDT_USD),
				},
				exchange_common.MARKET_BLUR_USD: {
					Ticker:         "BLUR-USDT",
					AdjustByMarket: newMarketIdWithValue(exchange_common.MARKET_USDT_USD),
				},
				exchange_common.MARKET_LDO_USD: {
					Ticker:         "LDO-USDT",
					AdjustByMarket: newMarketIdWithValue(exchange_common.MARKET_USDT_USD),
				},
				exchange_common.MARKET_OP_USD: {
					Ticker:         "OP-USDT",
					AdjustByMarket: newMarketIdWithValue(exchange_common.MARKET_USDT_USD),
				},
				exchange_common.MARKET_PEPE_USD: {
					Ticker:         "PEPE-USDT",
					AdjustByMarket: newMarketIdWithValue(exchange_common.MARKET_USDT_USD),
				},
				exchange_common.MARKET_SHIB_USD: {
					Ticker:         "SHIB-USDT",
					AdjustByMarket: newMarketIdWithValue(exchange_common.MARKET_USDT_USD),
				},
				exchange_common.MARKET_SUI_USD: {
					Ticker:         "SUI-USDT",
					AdjustByMarket: newMarketIdWithValue(exchange_common.MARKET_USDT_USD),
				},
				exchange_common.MARKET_WLD_USD: {
					Ticker:         "WLD-USDT",
					AdjustByMarket: newMarketIdWithValue(exchange_common.MARKET_USDT_USD),
				},
				exchange_common.MARKET_XRP_USD: {
					Ticker:         "XRP-USDT",
					AdjustByMarket: newMarketIdWithValue(exchange_common.MARKET_USDT_USD),
				},
				exchange_common.MARKET_USDT_USD: {
					Ticker:         "BTC-USDT", // Adjusted with BTC index price.
					AdjustByMarket: newMarketIdWithValue(exchange_common.MARKET_BTC_USD),
					Invert:         true,
				},
			},
		},
		exchange_common.EXCHANGE_ID_MEXC: {
			Id: exchange_common.EXCHANGE_ID_MEXC,
			MarketToMarketConfig: map[types.MarketId]types.MarketConfig{
				exchange_common.MARKET_UNI_USD: {
					Ticker:         "UNI_USDT",
					AdjustByMarket: newMarketIdWithValue(exchange_common.MARKET_USDT_USD),
				},
				exchange_common.MARKET_APE_USD: {
					Ticker:         "APE_USDT",
					AdjustByMarket: newMarketIdWithValue(exchange_common.MARKET_USDT_USD),
				},
				exchange_common.MARKET_APT_USD: {
					Ticker:         "APT_USDT",
					AdjustByMarket: newMarketIdWithValue(exchange_common.MARKET_USDT_USD),
				},
				exchange_common.MARKET_ARB_USD: {
					Ticker:         "ARB_USDT",
					AdjustByMarket: newMarketIdWithValue(exchange_common.MARKET_USDT_USD),
				},
				exchange_common.MARKET_BLUR_USD: {
					Ticker:         "BLUR_USDT",
					AdjustByMarket: newMarketIdWithValue(exchange_common.MARKET_USDT_USD),
				},
				exchange_common.MARKET_FIL_USD: {
					Ticker:         "FIL_USDT",
					AdjustByMarket: newMarketIdWithValue(exchange_common.MARKET_USDT_USD),
				},
				exchange_common.MARKET_LDO_USD: {
					Ticker:         "LDO_USDT",
					AdjustByMarket: newMarketIdWithValue(exchange_common.MARKET_USDT_USD),
				},
				exchange_common.MARKET_OP_USD: {
					Ticker:         "OP_USDT",
					AdjustByMarket: newMarketIdWithValue(exchange_common.MARKET_USDT_USD),
				},
				exchange_common.MARKET_PEPE_USD: {
					Ticker:         "PEPE_USDT",
					AdjustByMarket: newMarketIdWithValue(exchange_common.MARKET_USDT_USD),
				},
				exchange_common.MARKET_SEI_USD: {
					Ticker:         "SEI_USDT",
					AdjustByMarket: newMarketIdWithValue(exchange_common.MARKET_USDT_USD),
				},
				exchange_common.MARKET_SHIB_USD: {
					Ticker:         "SHIB_USDT",
					AdjustByMarket: newMarketIdWithValue(exchange_common.MARKET_USDT_USD),
				},
				exchange_common.MARKET_SUI_USD: {
					Ticker:         "SUI_USDT",
					AdjustByMarket: newMarketIdWithValue(exchange_common.MARKET_USDT_USD),
				},
				exchange_common.MARKET_WLD_USD: {
					Ticker:         "WLD_USDT",
					AdjustByMarket: newMarketIdWithValue(exchange_common.MARKET_USDT_USD),
				},
				exchange_common.MARKET_XLM_USD: {
					Ticker:         "XLM_USDT",
					AdjustByMarket: newMarketIdWithValue(exchange_common.MARKET_USDT_USD),
				},
				exchange_common.MARKET_XRP_USD: {
					Ticker:         "XRP_USDT",
					AdjustByMarket: newMarketIdWithValue(exchange_common.MARKET_USDT_USD),
				},
			},
		},
		exchange_common.EXCHANGE_ID_COINBASE_PRO: {
			Id: exchange_common.EXCHANGE_ID_COINBASE_PRO,
			MarketToMarketConfig: map[types.MarketId]types.MarketConfig{
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
				exchange_common.MARKET_LTC_USD: {
					Ticker: "LTC-USD",
				},
				exchange_common.MARKET_ATOM_USD: {
					Ticker: "ATOM-USD",
				},
				exchange_common.MARKET_UNI_USD: {
					Ticker: "UNI-USD",
				},
				exchange_common.MARKET_BCH_USD: {
					Ticker: "BCH-USD",
				},
				exchange_common.MARKET_NEAR_USD: {
					Ticker: "NEAR-USD",
				},
				exchange_common.MARKET_MKR_USD: {
					Ticker: "MKR-USD",
				},
				exchange_common.MARKET_XLM_USD: {
					Ticker: "XLM-USD",
				},
				exchange_common.MARKET_ETC_USD: {
					Ticker: "ETC-USD",
				},
				exchange_common.MARKET_APE_USD: {
					Ticker: "APE-USD",
				},
				exchange_common.MARKET_APT_USD: {
					Ticker: "APT-USD",
				},
				exchange_common.MARKET_ARB_USD: {
					Ticker: "ARB-USD",
				},
				exchange_common.MARKET_BLUR_USD: {
					Ticker: "BLUR-USD",
				},
				exchange_common.MARKET_LDO_USD: {
					Ticker: "LDO-USD",
				},
				exchange_common.MARKET_OP_USD: {
					Ticker: "OP-USD",
				},
				exchange_common.MARKET_SEI_USD: {
					Ticker: "SEI-USD",
				},
				exchange_common.MARKET_SHIB_USD: {
					Ticker: "SHIB-USD",
				},
				exchange_common.MARKET_XRP_USD: {
					Ticker: "XRP-USD",
				},
				exchange_common.MARKET_USDT_USD: {
					Ticker: "USDT-USD",
				},
			},
		},
		exchange_common.EXCHANGE_ID_TEST_EXCHANGE: {
			Id: exchange_common.EXCHANGE_ID_TEST_EXCHANGE,
			MarketToMarketConfig: map[types.MarketId]types.MarketConfig{
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
		},
		exchange_common.EXCHANGE_ID_TEST_VOLATILE_EXCHANGE: {
			Id: exchange_common.EXCHANGE_ID_TEST_VOLATILE_EXCHANGE,
			MarketToMarketConfig: map[types.MarketId]types.MarketConfig{
				exchange_common.MARKET_TEST_USD: {
					Ticker: "TEST-USD",
				},
			},
		},
	}
)

// newMarketIdWithValue returns a pointer to a new market id set to the specified value. This helper method
// is used to initialize the `AdjustByMarket` field of the `MarketConfig` structs above.
func newMarketIdWithValue(id types.MarketId) *types.MarketId {
	ptr := new(types.MarketId)
	*ptr = id
	return ptr
}

// GenerateExchangeConfigJson generates human-readable exchange config json for each market based on the contents
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
	marketToExchangeMarketConfigs := make(map[types.MarketId]map[string]types.ExchangeMarketConfigJson)

	// Generate the market-specific exchange config for each market, exchange.
	for id, exchangeConfig := range exchangeToExchangeConfig {
		// Skip config for the test exchange.
		if id == exchange_common.EXCHANGE_ID_TEST_EXCHANGE {
			continue
		}
		for marketId, config := range exchangeConfig.MarketToMarketConfig {
			marketExchangeConfigs, ok := marketToExchangeMarketConfigs[marketId]
			if !ok {
				marketToExchangeMarketConfigs[marketId] = map[string]types.ExchangeMarketConfigJson{}
				marketExchangeConfigs = marketToExchangeMarketConfigs[marketId]
			}

			exchangeMarketConfigJson := types.ExchangeMarketConfigJson{
				ExchangeName: id,
				Ticker:       config.Ticker,
				Invert:       config.Invert,
			}

			// Convert adjust-by market id to name if specified.
			if config.AdjustByMarket != nil {
				adjustByMarketName, ok := StaticMarketNames[*config.AdjustByMarket]
				if !ok {
					panic(fmt.Sprintf("no name for adjust-by market %v", *config.AdjustByMarket))
				}
				exchangeMarketConfigJson.AdjustByMarket = adjustByMarketName
			}

			marketExchangeConfigs[id] = exchangeMarketConfigJson
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
		sortedExchangeConfigs := make([]types.ExchangeMarketConfigJson, 0, len(exchangeNames))
		for _, exchangeName := range exchangeNames {
			sortedExchangeConfigs = append(sortedExchangeConfigs, exchangeToConfigs[exchangeName])
		}
		exchangeConfigJson := types.ExchangeConfigJson{
			Exchanges: sortedExchangeConfigs,
		}

		// 3. Generate human-readable formatted output json for the market, sorted by exchange name.
		bytes, err := json.MarshalIndent(exchangeConfigJson, "", "  ")
		if err != nil {
			panic(err)
		}
		marketToExchangeConfigJson[marketId] = string(bytes)
	}
	return marketToExchangeConfigJson
}
