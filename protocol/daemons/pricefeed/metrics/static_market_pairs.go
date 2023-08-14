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
		exchange_common.MARKET_AAVE_USD:  "AAVEUSD",
		exchange_common.MARKET_LTC_USD:   "LTCUSD",
		exchange_common.MARKET_DOGE_USD:  "DOGEUSD",
		exchange_common.MARKET_ICP_USD:   "ICPUSD",
		exchange_common.MARKET_ATOM_USD:  "ATOMUSD",
		exchange_common.MARKET_DOT_USD:   "DOTUSD",
		exchange_common.MARKET_XTZ_USD:   "XTZUSD",
		exchange_common.MARKET_UNI_USD:   "UNIUSD",
		exchange_common.MARKET_BCH_USD:   "BCHUSD",
		exchange_common.MARKET_EOS_USD:   "EOSUSD",
		exchange_common.MARKET_TRX_USD:   "TRXUSD",
		exchange_common.MARKET_ALGO_USD:  "ALGOUSD",
		exchange_common.MARKET_NEAR_USD:  "NEARUSD",
		exchange_common.MARKET_SNX_USD:   "SNXUSD",
		exchange_common.MARKET_MKR_USD:   "MKRUSD",
		exchange_common.MARKET_SUSHI_USD: "SUSHIUSD",
		exchange_common.MARKET_XLM_USD:   "XLMUSD",
		exchange_common.MARKET_XMR_USD:   "XMRUSD",
		exchange_common.MARKET_ETC_USD:   "ETCUSD",
		exchange_common.MARKET_1INCH_USD: "1INCHUSD",
		exchange_common.MARKET_COMP_USD:  "COMPUSD",
		exchange_common.MARKET_ZEC_USD:   "ZECUSD",
		exchange_common.MARKET_ZRX_USD:   "ZRXUSD",
		exchange_common.MARKET_YFI_USD:   "YFIUSD",
		exchange_common.MARKET_USDT_USD:  "USDTUSD",
	}
)
