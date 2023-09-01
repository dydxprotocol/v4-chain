package metrics

import (
	"github.com/dydxprotocol/v4-chain/protocol/daemons/pricefeed/client/constants/exchange_common"
	"github.com/dydxprotocol/v4-chain/protocol/daemons/pricefeed/client/types"
)

var (
	// TODO(DEC-586): get information from protocol.
	// StaticMarketPairs is the mapping of marketIds to their human readable market pairs.
	StaticMarketPairs = map[types.MarketId]string{
		exchange_common.MARKET_BTC_USD:   "BTCUSD",
		exchange_common.MARKET_ETH_USD:   "ETHUSD",
		exchange_common.MARKET_LINK_USD:  "LINKUSD",
		exchange_common.MARKET_MATIC_USD: "MATICUSD",
		exchange_common.MARKET_CRV_USD:   "CRVUSD",
		exchange_common.MARKET_SOL_USD:   "SOLUSD",
		exchange_common.MARKET_ADA_USD:   "ADAUSD",
		exchange_common.MARKET_AVAX_USD:  "AVAXUSD",
		exchange_common.MARKET_FIL_USD:   "FILUSD",
		exchange_common.MARKET_LTC_USD:   "LTCUSD",
		exchange_common.MARKET_DOGE_USD:  "DOGEUSD",
		exchange_common.MARKET_ATOM_USD:  "ATOMUSD",
		exchange_common.MARKET_DOT_USD:   "DOTUSD",
		exchange_common.MARKET_UNI_USD:   "UNIUSD",
		exchange_common.MARKET_BCH_USD:   "BCHUSD",
		exchange_common.MARKET_TRX_USD:   "TRXUSD",
		exchange_common.MARKET_NEAR_USD:  "NEARUSD",
		exchange_common.MARKET_MKR_USD:   "MKRUSD",
		exchange_common.MARKET_XLM_USD:   "XLMUSD",
		exchange_common.MARKET_ETC_USD:   "ETCUSD",
		exchange_common.MARKET_COMP_USD:  "COMPUSD",
		exchange_common.MARKET_USDT_USD:  "USDTUSD",
		exchange_common.MARKET_WLD_USD:   "WLDUSD",
		exchange_common.MARKET_APE_USD:   "APEUSD",
		exchange_common.MARKET_APT_USD:   "APTUSD",
		exchange_common.MARKET_ARB_USD:   "ARBUSD",
		exchange_common.MARKET_BLUR_USD:  "BLURUSD",
		exchange_common.MARKET_LDO_USD:   "LDOUSD",
		exchange_common.MARKET_OP_USD:    "OPUSD",
		exchange_common.MARKET_PEPE_USD:  "PEPEUSD",
		exchange_common.MARKET_SEI_USD:   "SEIUSD",
		exchange_common.MARKET_SHIB_USD:  "SHIBUSD",
		exchange_common.MARKET_SUI_USD:   "SUIUSD",
		exchange_common.MARKET_XRP_USD:   "XRPUSD",
		exchange_common.MARKET_TEST_USD:  "TESTUSD",
	}
)
