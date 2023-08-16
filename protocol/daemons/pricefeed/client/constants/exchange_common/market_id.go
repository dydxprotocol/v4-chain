package exchange_common

import "github.com/dydxprotocol/v4-chain/protocol/daemons/pricefeed/client/types"

// All market ids must match with the genesis state.
// TODO(CORE-296): Remove static daemon market config.
const (
	// MARKET_BTC_USD is the id for the BTC-USD market pair.
	MARKET_BTC_USD types.MarketId = 0
	// MARKET_ETH_USD is the id for the ETH-USD market pair.
	MARKET_ETH_USD types.MarketId = 1
	// MARKET_LINK_USD is the id for the LINK-USD market pair.
	MARKET_LINK_USD types.MarketId = 2
	// MARKET_MATIC_USD is the id for the MATIC-USD market pair.
	MARKET_MATIC_USD types.MarketId = 3
	// MARKET_CRV_USD is the id for the CRV-USD market pair.
	MARKET_CRV_USD types.MarketId = 4
	// MARKET_SOL_USD is the id for the SOL-USD market pair.
	MARKET_SOL_USD types.MarketId = 5
	// MARKET_ADA_USD is the id for the ADA-USD market pair.
	MARKET_ADA_USD types.MarketId = 6
	// MARKET_AVAX_USD is the id for the AVAX-USD market pair.
	MARKET_AVAX_USD types.MarketId = 7
	// MARKET_FIL_USD is the id for the FIL-USD market pair.
	MARKET_FIL_USD types.MarketId = 8
	// MARKET_AAVE_USD is the id for the AAVE-USD market pair.
	MARKET_AAVE_USD types.MarketId = 9
	// MARKET_LTC_USD is the id for the LTC-USD market pair.
	MARKET_LTC_USD types.MarketId = 10
	// MARKET_DOGE_USD is the id for the DOGE-USD market pair.
	MARKET_DOGE_USD types.MarketId = 11
	// MARKET_ICP_USD is the id for the ICP-USD market pair.
	MARKET_ICP_USD types.MarketId = 12
	// MARKET_ATOM_USD is the id for the ATOM-USD market pair.
	MARKET_ATOM_USD types.MarketId = 13
	// MARKET_DOT_USD is the id for the DOT-USD market pair.
	MARKET_DOT_USD types.MarketId = 14
	// MARKET_XTZ_USD is the id for the XTZ-USD market pair.
	MARKET_XTZ_USD types.MarketId = 15
	// MARKET_UNI_USD is the id for the UNI-USD market pair.
	MARKET_UNI_USD types.MarketId = 16
	// MARKET_BCH_USD is the id for the BCH-USD market pair.
	MARKET_BCH_USD types.MarketId = 17
	// MARKET_EOS_USD is the id for the EOS-USD market pair.
	MARKET_EOS_USD types.MarketId = 18
	// MARKET_TRX_USD is the id for the TRX-USD market pair.
	MARKET_TRX_USD types.MarketId = 19
	// MARKET_ALGO_USD is the id for the ALGO-USD market pair.
	MARKET_ALGO_USD types.MarketId = 20
	// MARKET_NEAR_USD is the id for the NEAR-USD market pair.
	MARKET_NEAR_USD types.MarketId = 21
	// MARKET_SNX_USD is the id for the SNX-USD market pair.
	MARKET_SNX_USD types.MarketId = 22
	// MARKET_MKR_USD is the id for the MKR-USD market pair.
	MARKET_MKR_USD types.MarketId = 23
	// MARKET_SUSHI_USD is the id for the SUSHI-USD market pair.
	MARKET_SUSHI_USD types.MarketId = 24
	// MARKET_XLM_USD is the id for the XLM-USD market pair.
	MARKET_XLM_USD types.MarketId = 25
	// MARKET_XMR_USD is the id for the XMR-USD market pair.
	MARKET_XMR_USD types.MarketId = 26
	// MARKET_ETC_USD is the id for the ETC-USD market pair.
	MARKET_ETC_USD types.MarketId = 27
	// MARKET_1INCH_USD is the id for the 1INCH-USD market pair.
	MARKET_1INCH_USD types.MarketId = 28
	// MARKET_COMP_USD is the id for the COMP-USD market pair.
	MARKET_COMP_USD types.MarketId = 29
	// MARKET_ZEC_USD is the id for the ZEC-USD market pair.
	MARKET_ZEC_USD types.MarketId = 30
	// MARKET_ZRX_USD is the id for the ZRX-USD market pair.
	MARKET_ZRX_USD types.MarketId = 31
	// MARKET_YFI_USD is the id for the YFI-USD market pair.
	MARKET_YFI_USD types.MarketId = 32

	// Non-trading adjust-by markets.
	// MARKET_USDT_USD is the id for the USDT-USD market pair.
	MARKET_USDT_USD types.MarketId = 33

	// Testing markets.
	// MARKET_TEST_USD is the id used for the TEST-USD market pair.
	MARKET_TEST_USD types.MarketId = 34
)
