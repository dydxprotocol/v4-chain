package constants

import (
	"github.com/dydxprotocol/v4-chain/protocol/daemons/pricefeed/client/constants/exchange_common"
	"github.com/dydxprotocol/v4-chain/protocol/daemons/pricefeed/client/types"
)

var (
	// StaticMarketNames maps marketIds to their human-readable market names. This list is
	// used for generating market exchange config that is then read back into the daemon.
	// Please do not use this mapping to determine the market name of given id, or vice versa,
	// in code that executes when the protocol is running.
	StaticMarketNames = map[types.MarketId]string{
		exchange_common.MARKET_BTC_USD:   "BTC-USD",
		exchange_common.MARKET_ETH_USD:   "ETH-USD",
		exchange_common.MARKET_LINK_USD:  "LINK-USD",
		exchange_common.MARKET_MATIC_USD: "MATIC-USD",
		exchange_common.MARKET_CRV_USD:   "CRV-USD",
		exchange_common.MARKET_SOL_USD:   "SOL-USD",
		exchange_common.MARKET_ADA_USD:   "ADA-USD",
		exchange_common.MARKET_AVAX_USD:  "AVAX-USD",
		exchange_common.MARKET_FIL_USD:   "FIL-USD",
		exchange_common.MARKET_LTC_USD:   "LTC-USD",
		exchange_common.MARKET_DOGE_USD:  "DOGE-USD",
		exchange_common.MARKET_ATOM_USD:  "ATOM-USD",
		exchange_common.MARKET_DOT_USD:   "DOT-USD",
		exchange_common.MARKET_UNI_USD:   "UNI-USD",
		exchange_common.MARKET_BCH_USD:   "BCH-USD",
		exchange_common.MARKET_TRX_USD:   "TRX-USD",
		exchange_common.MARKET_NEAR_USD:  "NEAR-USD",
		exchange_common.MARKET_MKR_USD:   "MKR-USD",
		exchange_common.MARKET_XLM_USD:   "XLM-USD",
		exchange_common.MARKET_ETC_USD:   "ETC-USD",
		exchange_common.MARKET_COMP_USD:  "COMP-USD",
		exchange_common.MARKET_WLD_USD:   "WLD-USD",
		exchange_common.MARKET_APE_USD:   "APE-USD",
		exchange_common.MARKET_APT_USD:   "APT-USD",
		exchange_common.MARKET_ARB_USD:   "ARB-USD",
		exchange_common.MARKET_BLUR_USD:  "BLUR-USD",
		exchange_common.MARKET_LDO_USD:   "LDO-USD",
		exchange_common.MARKET_OP_USD:    "OP-USD",
		exchange_common.MARKET_PEPE_USD:  "PEPE-USD",
		exchange_common.MARKET_SEI_USD:   "SEI-USD",
		exchange_common.MARKET_SHIB_USD:  "SHIB-USD",
		exchange_common.MARKET_SUI_USD:   "SUI-USD",
		exchange_common.MARKET_XRP_USD:   "XRP-USD",
		exchange_common.MARKET_USDT_USD:  "USDT-USD",
	}
)
