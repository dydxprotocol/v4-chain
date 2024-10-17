package exchange_config

import "github.com/dydxprotocol/v4-chain/protocol/daemons/pricefeed/client/types"

// All market ids must match with the genesis state.
const (
	// MARKET_BTC_USD is the id for the BTC-USD market pair.
	MARKET_BTC_USD types.MarketId = 0
	// MARKET_ETH_USD is the id for the ETH-USD market pair.
	MARKET_ETH_USD types.MarketId = 1
	// MARKET_LINK_USD is the id for the LINK-USD market pair.
	MARKET_LINK_USD types.MarketId = 2
	// MARKET_POL_USD is the id for the POL-USD market pair.
	MARKET_POL_USD types.MarketId = 3
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
	// MARKET_LTC_USD is the id for the LTC-USD market pair.
	MARKET_LTC_USD types.MarketId = 9
	// MARKET_DOGE_USD is the id for the DOGE-USD market pair.
	MARKET_DOGE_USD types.MarketId = 10
	// MARKET_ATOM_USD is the id for the ATOM-USD market pair.
	MARKET_ATOM_USD types.MarketId = 11
	// MARKET_DOT_USD is the id for the DOT-USD market pair.
	MARKET_DOT_USD types.MarketId = 12
	// MARKET_UNI_USD is the id for the UNI-USD market pair.
	MARKET_UNI_USD types.MarketId = 13
	// MARKET_BCH_USD is the id for the BCH-USD market pair.
	MARKET_BCH_USD types.MarketId = 14
	// MARKET_TRX_USD is the id for the TRX-USD market pair.
	MARKET_TRX_USD types.MarketId = 15
	// MARKET_NEAR_USD is the id for the NEAR-USD market pair.
	MARKET_NEAR_USD types.MarketId = 16
	// MARKET_MKR_USD is the id for the MKR-USD market pair.
	MARKET_MKR_USD types.MarketId = 17
	// MARKET_XLM_USD is the id for the XLM-USD market pair.
	MARKET_XLM_USD types.MarketId = 18
	// MARKET_ETC_USD is the id for the ETC-USD market pair.
	MARKET_ETC_USD types.MarketId = 19
	// MARKET_COMP_USD is the id for the COMP-USD market pair.
	MARKET_COMP_USD types.MarketId = 20
	// MARKET_WLD_USD is the id for the WLD-USD market pair.
	MARKET_WLD_USD types.MarketId = 21
	// MARKET_APE_USD is the id for the APE-USD market pair.
	MARKET_APE_USD types.MarketId = 22
	// MARKET_APT_USD is the id for the APT-USD market pair.
	MARKET_APT_USD types.MarketId = 23
	// MARKET_ARB_USD is the id for the ARB-USD market pair.
	MARKET_ARB_USD types.MarketId = 24
	// MARKET_BLUR_USD is the id for the BLUR-USD market pair.
	MARKET_BLUR_USD types.MarketId = 25
	// MARKET_LDO_USD is the id for the LDO-USD market pair.
	MARKET_LDO_USD types.MarketId = 26
	// MARKET_OP_USD is the id for the OP-USD market pair.
	MARKET_OP_USD types.MarketId = 27
	// MARKET_PEPE_USD is the id for the PEPE-USD market pair.
	MARKET_PEPE_USD types.MarketId = 28
	// MARKET_SEI_USD is the id for the SEI-USD market pair.
	MARKET_SEI_USD types.MarketId = 29
	// MARKET_SHIB_USD is the id for the SHIB-USD market pair.
	MARKET_SHIB_USD types.MarketId = 30
	// MARKET_SUI_USD is the id for the SUI-USD market pair.
	MARKET_SUI_USD types.MarketId = 31
	// MARKET_XRP_USD is the id for the XRP-USD market pair.
	MARKET_XRP_USD types.MarketId = 32

	// Testing markets used in local, staging, dev
	// MARKET_TEST_USD is the id used for the TEST-USD market pair.
	MARKET_TEST_USD types.MarketId = 33

	// Arbitrary isolated markets
	MARKET_ISO2_USD types.MarketId = 999_998
	MARKET_ISO_USD  types.MarketId = 999_999

	// Non-trading markets.
	// MARKET_USDT_USD is the id for the USDT-USD market pair.
	MARKET_USDT_USD types.MarketId = 1_000_000
	// MARKET_DYDX_USD is the id for the DYDX-USD market pair.
	MARKET_DYDX_USD types.MarketId = 1_000_001
)
