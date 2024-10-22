package exchange_config

import (
	"github.com/dydxprotocol/v4-chain/protocol/daemons/pricefeed/client/types"
)

var (
	// StaticMarketNames maps marketIds to their human-readable market names. This list is
	// used for generating market exchange config that is then read back into the daemon.
	StaticMarketNames = map[types.MarketId]string{
		MARKET_BTC_USD:  "BTC-USD",
		MARKET_ETH_USD:  "ETH-USD",
		MARKET_LINK_USD: "LINK-USD",
		MARKET_POL_USD:  "POL-USD",
		MARKET_CRV_USD:  "CRV-USD",
		MARKET_SOL_USD:  "SOL-USD",
		MARKET_ADA_USD:  "ADA-USD",
		MARKET_AVAX_USD: "AVAX-USD",
		MARKET_FIL_USD:  "FIL-USD",
		MARKET_LTC_USD:  "LTC-USD",
		MARKET_DOGE_USD: "DOGE-USD",
		MARKET_ATOM_USD: "ATOM-USD",
		MARKET_DOT_USD:  "DOT-USD",
		MARKET_UNI_USD:  "UNI-USD",
		MARKET_BCH_USD:  "BCH-USD",
		MARKET_TRX_USD:  "TRX-USD",
		MARKET_NEAR_USD: "NEAR-USD",
		MARKET_MKR_USD:  "MKR-USD",
		MARKET_XLM_USD:  "XLM-USD",
		MARKET_ETC_USD:  "ETC-USD",
		MARKET_COMP_USD: "COMP-USD",
		MARKET_WLD_USD:  "WLD-USD",
		MARKET_APE_USD:  "APE-USD",
		MARKET_APT_USD:  "APT-USD",
		MARKET_ARB_USD:  "ARB-USD",
		MARKET_BLUR_USD: "BLUR-USD",
		MARKET_LDO_USD:  "LDO-USD",
		MARKET_OP_USD:   "OP-USD",
		MARKET_PEPE_USD: "PEPE-USD",
		MARKET_SEI_USD:  "SEI-USD",
		MARKET_SHIB_USD: "SHIB-USD",
		MARKET_SUI_USD:  "SUI-USD",
		MARKET_XRP_USD:  "XRP-USD",
		MARKET_USDT_USD: "USDT-USD",
	}
)
