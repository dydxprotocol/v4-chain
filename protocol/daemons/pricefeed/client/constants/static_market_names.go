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
		exchange_common.MARKET_AAVE_USD:  "AAVE-USD",
		exchange_common.MARKET_LTC_USD:   "LTC-USD",
		exchange_common.MARKET_DOGE_USD:  "DOGE-USD",
		exchange_common.MARKET_ICP_USD:   "ICP-USD",
		exchange_common.MARKET_ATOM_USD:  "ATOM-USD",
		exchange_common.MARKET_DOT_USD:   "DOT-USD",
		exchange_common.MARKET_XTZ_USD:   "XTZ-USD",
		exchange_common.MARKET_UNI_USD:   "UNI-USD",
		exchange_common.MARKET_BCH_USD:   "BCH-USD",
		exchange_common.MARKET_EOS_USD:   "EOS-USD",
		exchange_common.MARKET_TRX_USD:   "TRX-USD",
		exchange_common.MARKET_ALGO_USD:  "ALGO-USD",
		exchange_common.MARKET_NEAR_USD:  "NEAR-USD",
		exchange_common.MARKET_SNX_USD:   "SNX-USD",
		exchange_common.MARKET_MKR_USD:   "MKR-USD",
		exchange_common.MARKET_SUSHI_USD: "SUSHI-USD",
		exchange_common.MARKET_XLM_USD:   "XLM-USD",
		exchange_common.MARKET_XMR_USD:   "XMR-USD",
		exchange_common.MARKET_ETC_USD:   "ETC-USD",
		exchange_common.MARKET_1INCH_USD: "1INCH-USD",
		exchange_common.MARKET_COMP_USD:  "COMP-USD",
		exchange_common.MARKET_ZEC_USD:   "ZEC-USD",
		exchange_common.MARKET_ZRX_USD:   "ZRX-USD",
		exchange_common.MARKET_YFI_USD:   "YFI-USD",
		exchange_common.MARKET_USDT_USD:  "USDT-USD",
	}
)
