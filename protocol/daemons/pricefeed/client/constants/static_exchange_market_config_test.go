package constants

import (
	"os"
	"testing"

	"github.com/dydxprotocol/v4-chain/protocol/daemons/pricefeed/client/constants/exchange_common"
	"github.com/dydxprotocol/v4-chain/protocol/daemons/pricefeed/client/types"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/daemons/pricefeed"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/json"
	"github.com/stretchr/testify/require"
)

func TestStaticExchangeMarketConfigCacheLen(t *testing.T) {
	require.Len(t, StaticExchangeMarketConfig, 13)
}

func TestGenerateExchangeConfigJsonLength(t *testing.T) {
	configs := GenerateExchangeConfigJson(StaticExchangeMarketConfig)
	require.Len(t, configs, 35)
}

func TestGenerateExchangeConfigJson(t *testing.T) {
	tests := map[string]struct {
		id                             types.MarketId
		expectedExchangeConfigJsonFile string
	}{
		"BTC exchange config": {
			id:                             exchange_common.MARKET_BTC_USD,
			expectedExchangeConfigJsonFile: "btc_exchange_config.json",
		},
		"ETH exchange config": {
			id:                             exchange_common.MARKET_ETH_USD,
			expectedExchangeConfigJsonFile: "eth_exchange_config.json",
		},
		"LINK exchange config": {
			id:                             exchange_common.MARKET_LINK_USD,
			expectedExchangeConfigJsonFile: "link_exchange_config.json",
		},
		"MATIC exchange config": {
			id:                             exchange_common.MARKET_MATIC_USD,
			expectedExchangeConfigJsonFile: "matic_exchange_config.json",
		},
		"CRV exchange config": {
			id:                             exchange_common.MARKET_CRV_USD,
			expectedExchangeConfigJsonFile: "crv_exchange_config.json",
		},
		"SOL exchange config": {
			id:                             exchange_common.MARKET_SOL_USD,
			expectedExchangeConfigJsonFile: "sol_exchange_config.json",
		},
		"ADA exchange config": {
			id:                             exchange_common.MARKET_ADA_USD,
			expectedExchangeConfigJsonFile: "ada_exchange_config.json",
		},
		"AVAX exchange config": {
			id:                             exchange_common.MARKET_AVAX_USD,
			expectedExchangeConfigJsonFile: "avax_exchange_config.json",
		},
		"FIL exchange config": {
			id:                             exchange_common.MARKET_FIL_USD,
			expectedExchangeConfigJsonFile: "fil_exchange_config.json",
		},
		"LTC exchange config": {
			id:                             exchange_common.MARKET_LTC_USD,
			expectedExchangeConfigJsonFile: "ltc_exchange_config.json",
		},
		"DOGE exchange config": {
			id:                             exchange_common.MARKET_DOGE_USD,
			expectedExchangeConfigJsonFile: "doge_exchange_config.json",
		},
		"ATOM exchange config": {
			id:                             exchange_common.MARKET_ATOM_USD,
			expectedExchangeConfigJsonFile: "atom_exchange_config.json",
		},
		"DOT exchange config": {
			id:                             exchange_common.MARKET_DOT_USD,
			expectedExchangeConfigJsonFile: "dot_exchange_config.json",
		},
		"UNI exchange config": {
			id:                             exchange_common.MARKET_UNI_USD,
			expectedExchangeConfigJsonFile: "uni_exchange_config.json",
		},
		"BCH exchange config": {
			id:                             exchange_common.MARKET_BCH_USD,
			expectedExchangeConfigJsonFile: "bch_exchange_config.json",
		},
		"TRX exchange config": {
			id:                             exchange_common.MARKET_TRX_USD,
			expectedExchangeConfigJsonFile: "trx_exchange_config.json",
		},
		"NEAR exchange config": {
			id:                             exchange_common.MARKET_NEAR_USD,
			expectedExchangeConfigJsonFile: "near_exchange_config.json",
		},
		"MKR exchange config": {
			id:                             exchange_common.MARKET_MKR_USD,
			expectedExchangeConfigJsonFile: "mkr_exchange_config.json",
		},
		"XLM exchange config": {
			id:                             exchange_common.MARKET_XLM_USD,
			expectedExchangeConfigJsonFile: "xlm_exchange_config.json",
		},
		"ETC exchange config": {
			id:                             exchange_common.MARKET_ETC_USD,
			expectedExchangeConfigJsonFile: "etc_exchange_config.json",
		},
		"COMP exchange config": {
			id:                             exchange_common.MARKET_COMP_USD,
			expectedExchangeConfigJsonFile: "comp_exchange_config.json",
		},
		"WLD exchange config": {
			id:                             exchange_common.MARKET_WLD_USD,
			expectedExchangeConfigJsonFile: "wld_exchange_config.json",
		},
		"APE exchange config": {
			id:                             exchange_common.MARKET_APE_USD,
			expectedExchangeConfigJsonFile: "ape_exchange_config.json",
		},
		"APT exchange config": {
			id:                             exchange_common.MARKET_APT_USD,
			expectedExchangeConfigJsonFile: "apt_exchange_config.json",
		},
		"ARB exchange config": {
			id:                             exchange_common.MARKET_ARB_USD,
			expectedExchangeConfigJsonFile: "arb_exchange_config.json",
		},
		"BLUR exchange config": {
			id:                             exchange_common.MARKET_BLUR_USD,
			expectedExchangeConfigJsonFile: "blur_exchange_config.json",
		},
		"LDO exchange config": {
			id:                             exchange_common.MARKET_LDO_USD,
			expectedExchangeConfigJsonFile: "ldo_exchange_config.json",
		},
		"OP exchange config": {
			id:                             exchange_common.MARKET_OP_USD,
			expectedExchangeConfigJsonFile: "op_exchange_config.json",
		},
		"PEPE exchange config": {
			id:                             exchange_common.MARKET_PEPE_USD,
			expectedExchangeConfigJsonFile: "pepe_exchange_config.json",
		},
		"SEI exchange config": {
			id:                             exchange_common.MARKET_SEI_USD,
			expectedExchangeConfigJsonFile: "sei_exchange_config.json",
		},
		"SHIB exchange config": {
			id:                             exchange_common.MARKET_SHIB_USD,
			expectedExchangeConfigJsonFile: "shib_exchange_config.json",
		},
		"SUI exchange config": {
			id:                             exchange_common.MARKET_SUI_USD,
			expectedExchangeConfigJsonFile: "sui_exchange_config.json",
		},
		"XRP exchange config": {
			id:                             exchange_common.MARKET_XRP_USD,
			expectedExchangeConfigJsonFile: "xrp_exchange_config.json",
		},
		"TEST exchange config": {
			id:                             exchange_common.MARKET_TEST_USD,
			expectedExchangeConfigJsonFile: "test_exchange_config.json",
		},
		"USDT exchange config": {
			id:                             exchange_common.MARKET_USDT_USD,
			expectedExchangeConfigJsonFile: "usdt_exchange_config.json",
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			exchangeCount := 0
			for _, exchangeConfig := range StaticExchangeMarketConfig {
				if _, ok := exchangeConfig.MarketToMarketConfig[tc.id]; ok {
					exchangeCount++
				}
			}
			if tc.id != exchange_common.MARKET_TEST_USD {
				// Ok to drop this to 5 for some markets if needed.
				require.GreaterOrEqual(t, exchangeCount, MinimumRequiredExchangesPerMarket)
			}

			configs := GenerateExchangeConfigJson(StaticExchangeMarketConfig)

			// Uncomment to update test data
			f, err := os.OpenFile("testdata/"+tc.expectedExchangeConfigJsonFile, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
			require.NoError(t, err)
			defer f.Close()
			_, err = f.WriteString(configs[tc.id] + "\n") // Final newline added manually.
			require.NoError(t, err)

			actualExchangeConfigJson := json.CompactJsonString(t, configs[tc.id])
			expectedExchangeConfigJson := pricefeed.ReadJsonTestFile(t, tc.expectedExchangeConfigJsonFile)
			require.Equal(t, expectedExchangeConfigJson, actualExchangeConfigJson)
		})
	}
}
