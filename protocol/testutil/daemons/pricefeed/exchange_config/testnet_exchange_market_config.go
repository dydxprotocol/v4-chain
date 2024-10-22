package exchange_config

import (
	"github.com/dydxprotocol/v4-chain/protocol/daemons/pricefeed/client/constants/exchange_common"
	"github.com/dydxprotocol/v4-chain/protocol/daemons/pricefeed/client/types"
)

var (
	// TestnetExchangeMarketConfig maps exchange feed ids to exchange market config. This map is used to generate
	// the exchange config json used to construct the genesis file for various testnet deploys defined in the testing
	// package - namely, localnet, dev, and staging. Note that public testnet is not affected by this map.
	TestnetExchangeMarketConfig = map[types.ExchangeId]*types.MutableExchangeMarketConfig{
		exchange_common.EXCHANGE_ID_BINANCE: {
			Id: exchange_common.EXCHANGE_ID_BINANCE,
			// example `symbols` parameter: ["BTCUSDT","BNBUSDT"]
			MarketToMarketConfig: map[types.MarketId]types.MarketConfig{
				MARKET_DYDX_USD: {
					Ticker:         "DYDXUSDT",
					AdjustByMarket: newMarketIdWithValue(MARKET_USDT_USD),
				},
				MARKET_BTC_USD: {
					Ticker:         "BTCUSDT",
					AdjustByMarket: newMarketIdWithValue(MARKET_USDT_USD),
				},
				MARKET_ETH_USD: {
					Ticker:         "ETHUSDT",
					AdjustByMarket: newMarketIdWithValue(MARKET_USDT_USD),
				},
				MARKET_LINK_USD: {
					Ticker:         "LINKUSDT",
					AdjustByMarket: newMarketIdWithValue(MARKET_USDT_USD),
				},
				MARKET_POL_USD: {
					Ticker:         "POLUSDT",
					AdjustByMarket: newMarketIdWithValue(MARKET_USDT_USD),
				},
				MARKET_CRV_USD: {
					Ticker:         "CRVUSDT",
					AdjustByMarket: newMarketIdWithValue(MARKET_USDT_USD),
				},
				MARKET_SOL_USD: {
					Ticker:         "SOLUSDT",
					AdjustByMarket: newMarketIdWithValue(MARKET_USDT_USD),
				},
				MARKET_ADA_USD: {
					Ticker:         "ADAUSDT",
					AdjustByMarket: newMarketIdWithValue(MARKET_USDT_USD),
				},
				MARKET_AVAX_USD: {
					Ticker:         "AVAXUSDT",
					AdjustByMarket: newMarketIdWithValue(MARKET_USDT_USD),
				},
				MARKET_FIL_USD: {
					Ticker:         "FILUSDT",
					AdjustByMarket: newMarketIdWithValue(MARKET_USDT_USD),
				},
				MARKET_LTC_USD: {
					Ticker:         "LTCUSDT",
					AdjustByMarket: newMarketIdWithValue(MARKET_USDT_USD),
				},
				MARKET_DOGE_USD: {
					Ticker:         "DOGEUSDT",
					AdjustByMarket: newMarketIdWithValue(MARKET_USDT_USD),
				},
				MARKET_ATOM_USD: {
					Ticker:         "ATOMUSDT",
					AdjustByMarket: newMarketIdWithValue(MARKET_USDT_USD),
				},
				MARKET_DOT_USD: {
					Ticker:         "DOTUSDT",
					AdjustByMarket: newMarketIdWithValue(MARKET_USDT_USD),
				},
				MARKET_UNI_USD: {
					Ticker:         "UNIUSDT",
					AdjustByMarket: newMarketIdWithValue(MARKET_USDT_USD),
				},
				MARKET_BCH_USD: {
					Ticker:         "BCHUSDT",
					AdjustByMarket: newMarketIdWithValue(MARKET_USDT_USD),
				},
				MARKET_TRX_USD: {
					Ticker:         "TRXUSDT",
					AdjustByMarket: newMarketIdWithValue(MARKET_USDT_USD),
				},
				MARKET_NEAR_USD: {
					Ticker:         "NEARUSDT",
					AdjustByMarket: newMarketIdWithValue(MARKET_USDT_USD),
				},
				MARKET_MKR_USD: {
					Ticker:         "MKRUSDT",
					AdjustByMarket: newMarketIdWithValue(MARKET_USDT_USD),
				},
				MARKET_XLM_USD: {
					Ticker:         "XLMUSDT",
					AdjustByMarket: newMarketIdWithValue(MARKET_USDT_USD),
				},
				MARKET_ETC_USD: {
					Ticker:         "ETCUSDT",
					AdjustByMarket: newMarketIdWithValue(MARKET_USDT_USD),
				},
				MARKET_COMP_USD: {
					Ticker:         "COMPUSDT",
					AdjustByMarket: newMarketIdWithValue(MARKET_USDT_USD),
				},
				MARKET_APE_USD: {
					Ticker:         "APEUSDT",
					AdjustByMarket: newMarketIdWithValue(MARKET_USDT_USD),
				},
				MARKET_APT_USD: {
					Ticker:         "APTUSDT",
					AdjustByMarket: newMarketIdWithValue(MARKET_USDT_USD),
				},
				MARKET_ARB_USD: {
					Ticker:         "ARBUSDT",
					AdjustByMarket: newMarketIdWithValue(MARKET_USDT_USD),
				},
				MARKET_LDO_USD: {
					Ticker:         "LDOUSDT",
					AdjustByMarket: newMarketIdWithValue(MARKET_USDT_USD),
				},
				MARKET_OP_USD: {
					Ticker:         "OPUSDT",
					AdjustByMarket: newMarketIdWithValue(MARKET_USDT_USD),
				},
				MARKET_PEPE_USD: {
					Ticker:         "PEPEUSDT",
					AdjustByMarket: newMarketIdWithValue(MARKET_USDT_USD),
				},
				MARKET_SEI_USD: {
					Ticker:         "SEIUSDT",
					AdjustByMarket: newMarketIdWithValue(MARKET_USDT_USD),
				},
				MARKET_SHIB_USD: {
					Ticker:         "SHIBUSDT",
					AdjustByMarket: newMarketIdWithValue(MARKET_USDT_USD),
				},
				MARKET_SUI_USD: {
					Ticker:         "SUIUSDT",
					AdjustByMarket: newMarketIdWithValue(MARKET_USDT_USD),
				},
				MARKET_WLD_USD: {
					Ticker:         "WLDUSDT",
					AdjustByMarket: newMarketIdWithValue(MARKET_USDT_USD),
				},
				MARKET_XRP_USD: {
					Ticker:         "XRPUSDT",
					AdjustByMarket: newMarketIdWithValue(MARKET_USDT_USD),
				},
				MARKET_USDT_USD: {
					Ticker: "USDCUSDT",
					Invert: true,
				},
			},
		},
		exchange_common.EXCHANGE_ID_KRAKEN: {
			Id: exchange_common.EXCHANGE_ID_KRAKEN,
			MarketToMarketConfig: map[types.MarketId]types.MarketConfig{
				MARKET_BTC_USD: {
					Ticker: "XXBTZUSD",
				},
				MARKET_ETH_USD: {
					Ticker: "XETHZUSD",
				},
				MARKET_LINK_USD: {
					Ticker: "LINKUSD",
				},
				MARKET_SOL_USD: {
					Ticker: "SOLUSD",
				},
				MARKET_ADA_USD: {
					Ticker: "ADAUSD",
				},
				MARKET_FIL_USD: {
					Ticker: "FILUSD",
				},
				MARKET_LTC_USD: {
					Ticker: "XLTCZUSD",
				},
				MARKET_ATOM_USD: {
					Ticker: "ATOMUSD",
				},
				MARKET_DOT_USD: {
					Ticker: "DOTUSD",
				},
				MARKET_BCH_USD: {
					Ticker: "BCHUSD",
				},
				MARKET_XLM_USD: {
					Ticker: "XXLMZUSD",
				},
				MARKET_APE_USD: {
					Ticker: "APEUSD",
				},
				MARKET_BLUR_USD: {
					Ticker: "BLURUSD",
				},
				MARKET_PEPE_USD: {
					Ticker: "PEPEUSD",
				},
				MARKET_SHIB_USD: {
					Ticker: "SHIBUSD",
				},
				MARKET_XRP_USD: {
					Ticker: "XXRPZUSD",
				},
				MARKET_UNI_USD: {
					Ticker: "UNIUSD",
				},
				MARKET_CRV_USD: {
					Ticker: "CRVUSD",
				},
				MARKET_COMP_USD: {
					Ticker: "COMPUSD",
				},
				MARKET_AVAX_USD: {
					Ticker: "AVAXUSD",
				},
				MARKET_DOGE_USD: {
					Ticker: "XDGUSD",
				},
				MARKET_LDO_USD: {
					Ticker: "LDOUSD",
				},
				MARKET_USDT_USD: {
					Ticker: "USDTZUSD",
				},
				MARKET_MKR_USD: {
					Ticker: "MKRUSD",
				},
				MARKET_TRX_USD: {
					Ticker: "TRXUSD",
				},
			},
		},
		exchange_common.EXCHANGE_ID_GATE: {
			Id: exchange_common.EXCHANGE_ID_GATE,
			MarketToMarketConfig: map[types.MarketId]types.MarketConfig{
				MARKET_DYDX_USD: {
					Ticker:         "DYDX_USDT",
					AdjustByMarket: newMarketIdWithValue(MARKET_USDT_USD),
				},
				MARKET_CRV_USD: {
					Ticker:         "CRV_USDT",
					AdjustByMarket: newMarketIdWithValue(MARKET_USDT_USD),
				},
				MARKET_ADA_USD: {
					Ticker:         "ADA_USDT",
					AdjustByMarket: newMarketIdWithValue(MARKET_USDT_USD),
				},
				MARKET_DOGE_USD: {
					Ticker:         "DOGE_USDT",
					AdjustByMarket: newMarketIdWithValue(MARKET_USDT_USD),
				},
				MARKET_DOT_USD: {
					Ticker:         "DOT_USDT",
					AdjustByMarket: newMarketIdWithValue(MARKET_USDT_USD),
				},
				MARKET_UNI_USD: {
					Ticker:         "UNI_USDT",
					AdjustByMarket: newMarketIdWithValue(MARKET_USDT_USD),
				},
				MARKET_BCH_USD: {
					Ticker:         "BCH_USDT",
					AdjustByMarket: newMarketIdWithValue(MARKET_USDT_USD),
				},
				MARKET_TRX_USD: {
					Ticker:         "TRX_USDT",
					AdjustByMarket: newMarketIdWithValue(MARKET_USDT_USD),
				},
				MARKET_NEAR_USD: {
					Ticker:         "NEAR_USDT",
					AdjustByMarket: newMarketIdWithValue(MARKET_USDT_USD),
				},
				MARKET_ETC_USD: {
					Ticker:         "ETC_USDT",
					AdjustByMarket: newMarketIdWithValue(MARKET_USDT_USD),
				},
				MARKET_APE_USD: {
					Ticker:         "APE_USDT",
					AdjustByMarket: newMarketIdWithValue(MARKET_USDT_USD),
				},
				MARKET_APT_USD: {
					Ticker:         "APT_USDT",
					AdjustByMarket: newMarketIdWithValue(MARKET_USDT_USD),
				},
				MARKET_ARB_USD: {
					Ticker:         "ARB_USDT",
					AdjustByMarket: newMarketIdWithValue(MARKET_USDT_USD),
				},
				MARKET_BLUR_USD: {
					Ticker:         "BLUR_USDT",
					AdjustByMarket: newMarketIdWithValue(MARKET_USDT_USD),
				},
				MARKET_FIL_USD: {
					Ticker:         "FIL_USDT",
					AdjustByMarket: newMarketIdWithValue(MARKET_USDT_USD),
				},
				MARKET_OP_USD: {
					Ticker:         "OP_USDT",
					AdjustByMarket: newMarketIdWithValue(MARKET_USDT_USD),
				},
				MARKET_PEPE_USD: {
					Ticker:         "PEPE_USDT",
					AdjustByMarket: newMarketIdWithValue(MARKET_USDT_USD),
				},
				MARKET_SEI_USD: {
					Ticker:         "SEI_USDT",
					AdjustByMarket: newMarketIdWithValue(MARKET_USDT_USD),
				},
				MARKET_SHIB_USD: {
					Ticker:         "SHIB_USDT",
					AdjustByMarket: newMarketIdWithValue(MARKET_USDT_USD),
				},
				MARKET_SUI_USD: {
					Ticker:         "SUI_USDT",
					AdjustByMarket: newMarketIdWithValue(MARKET_USDT_USD),
				},
				MARKET_WLD_USD: {
					Ticker:         "WLD_USDT",
					AdjustByMarket: newMarketIdWithValue(MARKET_USDT_USD),
				},
				MARKET_XRP_USD: {
					Ticker:         "XRP_USDT",
					AdjustByMarket: newMarketIdWithValue(MARKET_USDT_USD),
				},
				MARKET_AVAX_USD: {
					Ticker:         "AVAX_USDT",
					AdjustByMarket: newMarketIdWithValue(MARKET_USDT_USD),
				},
				MARKET_ATOM_USD: {
					Ticker:         "ATOM_USDT",
					AdjustByMarket: newMarketIdWithValue(MARKET_USDT_USD),
				},
				MARKET_COMP_USD: {
					Ticker:         "COMP_USDT",
					AdjustByMarket: newMarketIdWithValue(MARKET_USDT_USD),
				},
			},
		},
		exchange_common.EXCHANGE_ID_BITSTAMP: {
			Id:                   exchange_common.EXCHANGE_ID_BITSTAMP,
			MarketToMarketConfig: map[types.MarketId]types.MarketConfig{},
		},
		exchange_common.EXCHANGE_ID_BYBIT: {
			Id: exchange_common.EXCHANGE_ID_BYBIT,
			MarketToMarketConfig: map[types.MarketId]types.MarketConfig{
				MARKET_DYDX_USD: {
					Ticker:         "DYDXUSDT",
					AdjustByMarket: newMarketIdWithValue(MARKET_USDT_USD),
				},
				MARKET_BTC_USD: {
					Ticker:         "BTCUSDT",
					AdjustByMarket: newMarketIdWithValue(MARKET_USDT_USD),
				},
				MARKET_ETH_USD: {
					Ticker:         "ETHUSDT",
					AdjustByMarket: newMarketIdWithValue(MARKET_USDT_USD),
				},
				MARKET_XRP_USD: {
					Ticker:         "XRPUSDT",
					AdjustByMarket: newMarketIdWithValue(MARKET_USDT_USD),
				},
				MARKET_LTC_USD: {
					Ticker:         "LTCUSDT",
					AdjustByMarket: newMarketIdWithValue(MARKET_USDT_USD),
				},
				MARKET_ATOM_USD: {
					Ticker:         "ATOMUSDT",
					AdjustByMarket: newMarketIdWithValue(MARKET_USDT_USD),
				},
				MARKET_UNI_USD: {
					Ticker:         "UNIUSDT",
					AdjustByMarket: newMarketIdWithValue(MARKET_USDT_USD),
				},
				MARKET_WLD_USD: {
					Ticker:         "WLDUSDT",
					AdjustByMarket: newMarketIdWithValue(MARKET_USDT_USD),
				},
				MARKET_APT_USD: {
					Ticker:         "APTUSDT",
					AdjustByMarket: newMarketIdWithValue(MARKET_USDT_USD),
				},
				MARKET_SOL_USD: {
					Ticker:         "SOLUSDT",
					AdjustByMarket: newMarketIdWithValue(MARKET_USDT_USD),
				},
				MARKET_DOGE_USD: {
					Ticker:         "DOGEUSDT",
					AdjustByMarket: newMarketIdWithValue(MARKET_USDT_USD),
				},
				MARKET_ADA_USD: {
					Ticker:         "ADAUSDT",
					AdjustByMarket: newMarketIdWithValue(MARKET_USDT_USD),
				},
				MARKET_XLM_USD: {
					Ticker:         "XLMUSDT",
					AdjustByMarket: newMarketIdWithValue(MARKET_USDT_USD),
				},
				MARKET_SHIB_USD: {
					Ticker:         "SHIBUSDT",
					AdjustByMarket: newMarketIdWithValue(MARKET_USDT_USD),
				},
				MARKET_LINK_USD: {
					Ticker:         "LINKUSDT",
					AdjustByMarket: newMarketIdWithValue(MARKET_USDT_USD),
				},
				MARKET_ARB_USD: {
					Ticker:         "ARBUSDT",
					AdjustByMarket: newMarketIdWithValue(MARKET_USDT_USD),
				},
				MARKET_SUI_USD: {
					Ticker:         "SUIUSDT",
					AdjustByMarket: newMarketIdWithValue(MARKET_USDT_USD),
				},
				MARKET_TRX_USD: {
					Ticker:         "TRXUSDT",
					AdjustByMarket: newMarketIdWithValue(MARKET_USDT_USD),
				},
				MARKET_SEI_USD: {
					Ticker:         "SEIUSDT",
					AdjustByMarket: newMarketIdWithValue(MARKET_USDT_USD),
				},
				MARKET_PEPE_USD: {
					Ticker:         "PEPEUSDT",
					AdjustByMarket: newMarketIdWithValue(MARKET_USDT_USD),
				},
				MARKET_AVAX_USD: {
					Ticker:         "AVAXUSDT",
					AdjustByMarket: newMarketIdWithValue(MARKET_USDT_USD),
				},
				MARKET_BCH_USD: {
					Ticker:         "BCHUSDT",
					AdjustByMarket: newMarketIdWithValue(MARKET_USDT_USD),
				},
				MARKET_DOT_USD: {
					Ticker:         "DOTUSDT",
					AdjustByMarket: newMarketIdWithValue(MARKET_USDT_USD),
				},
				MARKET_POL_USD: {
					Ticker:         "POLUSDT",
					AdjustByMarket: newMarketIdWithValue(MARKET_USDT_USD),
				},
				MARKET_USDT_USD: {
					Ticker: "USDCUSDT",
					Invert: true,
				},
			},
		},
		exchange_common.EXCHANGE_ID_CRYPTO_COM: {
			Id: exchange_common.EXCHANGE_ID_CRYPTO_COM,
			MarketToMarketConfig: map[types.MarketId]types.MarketConfig{
				MARKET_POL_USD: {
					Ticker: "POL_USD",
				},
			},
		},
		exchange_common.EXCHANGE_ID_HUOBI: {
			Id: exchange_common.EXCHANGE_ID_HUOBI,
			MarketToMarketConfig: map[types.MarketId]types.MarketConfig{
				MARKET_SOL_USD: {
					Ticker:         "solusdt",
					AdjustByMarket: newMarketIdWithValue(MARKET_USDT_USD),
				},
				MARKET_ADA_USD: {
					Ticker:         "adausdt",
					AdjustByMarket: newMarketIdWithValue(MARKET_USDT_USD),
				},
				MARKET_FIL_USD: {
					Ticker:         "filusdt",
					AdjustByMarket: newMarketIdWithValue(MARKET_USDT_USD),
				},
				MARKET_LTC_USD: {
					Ticker:         "ltcusdt",
					AdjustByMarket: newMarketIdWithValue(MARKET_USDT_USD),
				},
				MARKET_DOGE_USD: {
					Ticker:         "dogeusdt",
					AdjustByMarket: newMarketIdWithValue(MARKET_USDT_USD),
				},
				MARKET_BCH_USD: {
					Ticker:         "bchusdt",
					AdjustByMarket: newMarketIdWithValue(MARKET_USDT_USD),
				},
				MARKET_TRX_USD: {
					Ticker:         "trxusdt",
					AdjustByMarket: newMarketIdWithValue(MARKET_USDT_USD),
				},
				MARKET_APT_USD: {
					Ticker:         "aptusdt",
					AdjustByMarket: newMarketIdWithValue(MARKET_USDT_USD),
				},
				MARKET_ARB_USD: {
					Ticker:         "arbusdt",
					AdjustByMarket: newMarketIdWithValue(MARKET_USDT_USD),
				},
				MARKET_SEI_USD: {
					Ticker:         "seiusdt",
					AdjustByMarket: newMarketIdWithValue(MARKET_USDT_USD),
				},
				MARKET_SUI_USD: {
					Ticker:         "suiusdt",
					AdjustByMarket: newMarketIdWithValue(MARKET_USDT_USD),
				},
				MARKET_WLD_USD: {
					Ticker:         "wldusdt",
					AdjustByMarket: newMarketIdWithValue(MARKET_USDT_USD),
				},
				MARKET_XRP_USD: {
					Ticker:         "xrpusdt",
					AdjustByMarket: newMarketIdWithValue(MARKET_USDT_USD),
				},
				MARKET_AVAX_USD: {
					Ticker:         "avaxusdt",
					AdjustByMarket: newMarketIdWithValue(MARKET_USDT_USD),
				},
				MARKET_NEAR_USD: {
					Ticker:         "nearusdt",
					AdjustByMarket: newMarketIdWithValue(MARKET_USDT_USD),
				},
				MARKET_ETC_USD: {
					Ticker:         "etcusdt",
					AdjustByMarket: newMarketIdWithValue(MARKET_USDT_USD),
				},
				MARKET_BTC_USD: {
					Ticker:         "btcusdt",
					AdjustByMarket: newMarketIdWithValue(MARKET_USDT_USD),
				},
				MARKET_ETH_USD: {
					Ticker:         "ethusdt",
					AdjustByMarket: newMarketIdWithValue(MARKET_USDT_USD),
				},
				MARKET_USDT_USD: {
					Ticker:         "ethusdt", // Adjusted with ETH index price.
					AdjustByMarket: newMarketIdWithValue(MARKET_ETH_USD),
					Invert:         true,
				},
			},
		},
		exchange_common.EXCHANGE_ID_KUCOIN: {
			Id: exchange_common.EXCHANGE_ID_KUCOIN,
			MarketToMarketConfig: map[types.MarketId]types.MarketConfig{
				MARKET_DYDX_USD: {
					Ticker:         "DYDX-USDT",
					AdjustByMarket: newMarketIdWithValue(MARKET_USDT_USD),
				},
				MARKET_LINK_USD: {
					Ticker:         "LINK-USDT",
					AdjustByMarket: newMarketIdWithValue(MARKET_USDT_USD),
				},
				MARKET_CRV_USD: {
					Ticker:         "CRV-USDT",
					AdjustByMarket: newMarketIdWithValue(MARKET_USDT_USD),
				},
				MARKET_SOL_USD: {
					Ticker:         "SOL-USDT",
					AdjustByMarket: newMarketIdWithValue(MARKET_USDT_USD),
				},
				MARKET_ADA_USD: {
					Ticker:         "ADA-USDT",
					AdjustByMarket: newMarketIdWithValue(MARKET_USDT_USD),
				},
				MARKET_AVAX_USD: {
					Ticker:         "AVAX-USDT",
					AdjustByMarket: newMarketIdWithValue(MARKET_USDT_USD),
				},
				MARKET_LTC_USD: {
					Ticker:         "LTC-USDT",
					AdjustByMarket: newMarketIdWithValue(MARKET_USDT_USD),
				},
				MARKET_DOGE_USD: {
					Ticker:         "DOGE-USDT",
					AdjustByMarket: newMarketIdWithValue(MARKET_USDT_USD),
				},
				MARKET_ATOM_USD: {
					Ticker:         "ATOM-USDT",
					AdjustByMarket: newMarketIdWithValue(MARKET_USDT_USD),
				},
				MARKET_DOT_USD: {
					Ticker:         "DOT-USDT",
					AdjustByMarket: newMarketIdWithValue(MARKET_USDT_USD),
				},
				MARKET_XLM_USD: {
					Ticker:         "XLM-USDT",
					AdjustByMarket: newMarketIdWithValue(MARKET_USDT_USD),
				},
				MARKET_BCH_USD: {
					Ticker:         "BCH-USDT",
					AdjustByMarket: newMarketIdWithValue(MARKET_USDT_USD),
				},
				MARKET_TRX_USD: {
					Ticker:         "TRX-USDT",
					AdjustByMarket: newMarketIdWithValue(MARKET_USDT_USD),
				},
				MARKET_ARB_USD: {
					Ticker:         "ARB-USDT",
					AdjustByMarket: newMarketIdWithValue(MARKET_USDT_USD),
				},
				MARKET_BLUR_USD: {
					Ticker:         "BLUR-USDT",
					AdjustByMarket: newMarketIdWithValue(MARKET_USDT_USD),
				},
				MARKET_LDO_USD: {
					Ticker:         "LDO-USDT",
					AdjustByMarket: newMarketIdWithValue(MARKET_USDT_USD),
				},
				MARKET_OP_USD: {
					Ticker:         "OP-USDT",
					AdjustByMarket: newMarketIdWithValue(MARKET_USDT_USD),
				},
				MARKET_PEPE_USD: {
					Ticker:         "PEPE-USDT",
					AdjustByMarket: newMarketIdWithValue(MARKET_USDT_USD),
				},
				MARKET_SEI_USD: {
					Ticker:         "SEI-USDT",
					AdjustByMarket: newMarketIdWithValue(MARKET_USDT_USD),
				},
				MARKET_SHIB_USD: {
					Ticker:         "SHIB-USDT",
					AdjustByMarket: newMarketIdWithValue(MARKET_USDT_USD),
				},
				MARKET_SUI_USD: {
					Ticker:         "SUI-USDT",
					AdjustByMarket: newMarketIdWithValue(MARKET_USDT_USD),
				},
				MARKET_WLD_USD: {
					Ticker:         "WLD-USDT",
					AdjustByMarket: newMarketIdWithValue(MARKET_USDT_USD),
				},
				MARKET_XRP_USD: {
					Ticker:         "XRP-USDT",
					AdjustByMarket: newMarketIdWithValue(MARKET_USDT_USD),
				},
				MARKET_MKR_USD: {
					Ticker:         "MKR-USDT",
					AdjustByMarket: newMarketIdWithValue(MARKET_USDT_USD),
				},
				MARKET_NEAR_USD: {
					Ticker:         "NEAR-USDT",
					AdjustByMarket: newMarketIdWithValue(MARKET_USDT_USD),
				},
				MARKET_APE_USD: {
					Ticker:         "APE-USDT",
					AdjustByMarket: newMarketIdWithValue(MARKET_USDT_USD),
				},
				MARKET_APT_USD: {
					Ticker:         "APT-USDT",
					AdjustByMarket: newMarketIdWithValue(MARKET_USDT_USD),
				},
				MARKET_USDT_USD: {
					Ticker:         "BTC-USDT", // Adjusted with BTC index price.
					AdjustByMarket: newMarketIdWithValue(MARKET_BTC_USD),
					Invert:         true,
				},
				MARKET_BTC_USD: {
					Ticker:         "BTC-USDT",
					AdjustByMarket: newMarketIdWithValue(MARKET_USDT_USD),
				},
				MARKET_ETC_USD: {
					Ticker:         "ETC-USDT",
					AdjustByMarket: newMarketIdWithValue(MARKET_USDT_USD),
				},
				MARKET_ETH_USD: {
					Ticker:         "ETH-USDT",
					AdjustByMarket: newMarketIdWithValue(MARKET_USDT_USD),
				},
				MARKET_UNI_USD: {
					Ticker:         "UNI-USDT",
					AdjustByMarket: newMarketIdWithValue(MARKET_USDT_USD),
				},
			},
		},
		exchange_common.EXCHANGE_ID_OKX: {
			Id: exchange_common.EXCHANGE_ID_OKX,
			MarketToMarketConfig: map[types.MarketId]types.MarketConfig{
				MARKET_DYDX_USD: {
					Ticker:         "DYDX-USDT",
					AdjustByMarket: newMarketIdWithValue(MARKET_USDT_USD),
				},
				MARKET_BTC_USD: {
					Ticker:         "BTC-USDT",
					AdjustByMarket: newMarketIdWithValue(MARKET_USDT_USD),
				},
				MARKET_ETH_USD: {
					Ticker:         "ETH-USDT",
					AdjustByMarket: newMarketIdWithValue(MARKET_USDT_USD),
				},
				MARKET_LINK_USD: {
					Ticker:         "LINK-USDT",
					AdjustByMarket: newMarketIdWithValue(MARKET_USDT_USD),
				},
				MARKET_POL_USD: {
					Ticker:         "POL-USDT",
					AdjustByMarket: newMarketIdWithValue(MARKET_USDT_USD),
				},
				MARKET_CRV_USD: {
					Ticker:         "CRV-USDT",
					AdjustByMarket: newMarketIdWithValue(MARKET_USDT_USD),
				},
				MARKET_SOL_USD: {
					Ticker:         "SOL-USDT",
					AdjustByMarket: newMarketIdWithValue(MARKET_USDT_USD),
				},
				MARKET_AVAX_USD: {
					Ticker:         "AVAX-USDT",
					AdjustByMarket: newMarketIdWithValue(MARKET_USDT_USD),
				},
				MARKET_FIL_USD: {
					Ticker:         "FIL-USDT",
					AdjustByMarket: newMarketIdWithValue(MARKET_USDT_USD),
				},
				MARKET_LTC_USD: {
					Ticker:         "LTC-USDT",
					AdjustByMarket: newMarketIdWithValue(MARKET_USDT_USD),
				},
				MARKET_DOGE_USD: {
					Ticker:         "DOGE-USDT",
					AdjustByMarket: newMarketIdWithValue(MARKET_USDT_USD),
				},
				MARKET_DOT_USD: {
					Ticker:         "DOT-USDT",
					AdjustByMarket: newMarketIdWithValue(MARKET_USDT_USD),
				},
				MARKET_UNI_USD: {
					Ticker:         "UNI-USDT",
					AdjustByMarket: newMarketIdWithValue(MARKET_USDT_USD),
				},
				MARKET_BCH_USD: {
					Ticker:         "BCH-USDT",
					AdjustByMarket: newMarketIdWithValue(MARKET_USDT_USD),
				},
				MARKET_TRX_USD: {
					Ticker:         "TRX-USDT",
					AdjustByMarket: newMarketIdWithValue(MARKET_USDT_USD),
				},
				MARKET_ETC_USD: {
					Ticker:         "ETC-USDT",
					AdjustByMarket: newMarketIdWithValue(MARKET_USDT_USD),
				},
				MARKET_APE_USD: {
					Ticker:         "APE-USDT",
					AdjustByMarket: newMarketIdWithValue(MARKET_USDT_USD),
				},
				MARKET_ARB_USD: {
					Ticker:         "ARB-USDT",
					AdjustByMarket: newMarketIdWithValue(MARKET_USDT_USD),
				},
				MARKET_BLUR_USD: {
					Ticker:         "BLUR-USDT",
					AdjustByMarket: newMarketIdWithValue(MARKET_USDT_USD),
				},
				MARKET_OP_USD: {
					Ticker:         "OP-USDT",
					AdjustByMarket: newMarketIdWithValue(MARKET_USDT_USD),
				},
				MARKET_PEPE_USD: {
					Ticker:         "PEPE-USDT",
					AdjustByMarket: newMarketIdWithValue(MARKET_USDT_USD),
				},
				MARKET_SHIB_USD: {
					Ticker:         "SHIB-USDT",
					AdjustByMarket: newMarketIdWithValue(MARKET_USDT_USD),
				},
				MARKET_SUI_USD: {
					Ticker:         "SUI-USDT",
					AdjustByMarket: newMarketIdWithValue(MARKET_USDT_USD),
				},
				MARKET_WLD_USD: {
					Ticker:         "WLD-USDT",
					AdjustByMarket: newMarketIdWithValue(MARKET_USDT_USD),
				},
				MARKET_XRP_USD: {
					Ticker:         "XRP-USDT",
					AdjustByMarket: newMarketIdWithValue(MARKET_USDT_USD),
				},
				MARKET_COMP_USD: {
					Ticker:         "COMP-USDT",
					AdjustByMarket: newMarketIdWithValue(MARKET_USDT_USD),
				},
				MARKET_MKR_USD: {
					Ticker:         "MKR-USDT",
					AdjustByMarket: newMarketIdWithValue(MARKET_USDT_USD),
				},
				MARKET_APT_USD: {
					Ticker:         "APT-USDT",
					AdjustByMarket: newMarketIdWithValue(MARKET_USDT_USD),
				},
				MARKET_ATOM_USD: {
					Ticker:         "ATOM-USDT",
					AdjustByMarket: newMarketIdWithValue(MARKET_USDT_USD),
				},
				MARKET_ADA_USD: {
					Ticker:         "ADA-USDT",
					AdjustByMarket: newMarketIdWithValue(MARKET_USDT_USD),
				},
				MARKET_LDO_USD: {
					Ticker:         "LDO-USDT",
					AdjustByMarket: newMarketIdWithValue(MARKET_USDT_USD),
				},
				MARKET_USDT_USD: {
					Ticker: "USDC-USDT",
					Invert: true,
				},
				MARKET_XLM_USD: {
					Ticker:         "XLM-USDT",
					AdjustByMarket: newMarketIdWithValue(MARKET_USDT_USD),
				},
				MARKET_NEAR_USD: {
					Ticker:         "NEAR-USDT",
					AdjustByMarket: newMarketIdWithValue(MARKET_USDT_USD),
				},
			},
		},
		exchange_common.EXCHANGE_ID_COINBASE_PRO: {
			Id: exchange_common.EXCHANGE_ID_COINBASE_PRO,
			MarketToMarketConfig: map[types.MarketId]types.MarketConfig{
				MARKET_BTC_USD: {
					Ticker: "BTC-USD",
				},
				MARKET_ETH_USD: {
					Ticker: "ETH-USD",
				},
				MARKET_LINK_USD: {
					Ticker: "LINK-USD",
				},
				MARKET_POL_USD: {
					Ticker: "POL-USD",
				},
				MARKET_CRV_USD: {
					Ticker: "CRV-USD",
				},
				MARKET_SOL_USD: {
					Ticker: "SOL-USD",
				},
				MARKET_ADA_USD: {
					Ticker: "ADA-USD",
				},
				MARKET_FIL_USD: {
					Ticker: "FIL-USD",
				},
				MARKET_LTC_USD: {
					Ticker: "LTC-USD",
				},
				MARKET_ATOM_USD: {
					Ticker: "ATOM-USD",
				},
				MARKET_UNI_USD: {
					Ticker: "UNI-USD",
				},
				MARKET_BCH_USD: {
					Ticker: "BCH-USD",
				},
				MARKET_NEAR_USD: {
					Ticker: "NEAR-USD",
				},
				MARKET_MKR_USD: {
					Ticker: "MKR-USD",
				},
				MARKET_XLM_USD: {
					Ticker: "XLM-USD",
				},
				MARKET_ETC_USD: {
					Ticker: "ETC-USD",
				},
				MARKET_COMP_USD: {
					Ticker: "COMP-USD",
				},
				MARKET_APE_USD: {
					Ticker: "APE-USD",
				},
				MARKET_APT_USD: {
					Ticker: "APT-USD",
				},
				MARKET_ARB_USD: {
					Ticker: "ARB-USD",
				},
				MARKET_BLUR_USD: {
					Ticker: "BLUR-USD",
				},
				MARKET_LDO_USD: {
					Ticker: "LDO-USD",
				},
				MARKET_OP_USD: {
					Ticker: "OP-USD",
				},
				MARKET_SEI_USD: {
					Ticker: "SEI-USD",
				},
				MARKET_SHIB_USD: {
					Ticker: "SHIB-USD",
				},
				MARKET_SUI_USD: {
					Ticker: "SUI-USD",
				},
				MARKET_XRP_USD: {
					Ticker: "XRP-USD",
				},
				MARKET_AVAX_USD: {
					Ticker: "AVAX-USD",
				},
				MARKET_DOGE_USD: {
					Ticker: "DOGE-USD",
				},
				MARKET_DOT_USD: {
					Ticker: "DOT-USD",
				},
				MARKET_USDT_USD: {
					Ticker: "USDT-USD",
				},
			},
		},
		exchange_common.EXCHANGE_ID_TEST_EXCHANGE: {
			Id: exchange_common.EXCHANGE_ID_TEST_EXCHANGE,
			MarketToMarketConfig: map[types.MarketId]types.MarketConfig{
				MARKET_BTC_USD: {
					Ticker: "BTC-USD",
				},
				MARKET_ETH_USD: {
					Ticker: "ETH-USD",
				},
				MARKET_LINK_USD: {
					Ticker: "LINK-USD",
				},
				MARKET_USDT_USD: {
					Ticker: "USDT-USD",
				},
			},
		},
		exchange_common.EXCHANGE_ID_TEST_VOLATILE_EXCHANGE: {
			Id: exchange_common.EXCHANGE_ID_TEST_VOLATILE_EXCHANGE,
			MarketToMarketConfig: map[types.MarketId]types.MarketConfig{
				MARKET_TEST_USD: {
					Ticker: "TEST-USD",
				},
			},
		},
		exchange_common.EXCHANGE_ID_TEST_FIXED_PRICE_EXCHANGE: {
			Id: exchange_common.EXCHANGE_ID_TEST_FIXED_PRICE_EXCHANGE,
			MarketToMarketConfig: map[types.MarketId]types.MarketConfig{
				MARKET_BTC_USD: {
					Ticker: "BTC-USD",
				},
				MARKET_ETH_USD: {
					Ticker: "ETH-USD",
				},
				MARKET_SOL_USD: {
					Ticker: "SOL-USD",
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
