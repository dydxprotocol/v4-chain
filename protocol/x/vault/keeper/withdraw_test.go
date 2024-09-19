package keeper_test

import (
	"math/big"
	"testing"

	"github.com/cometbft/cometbft/types"
	testapp "github.com/dydxprotocol/v4-chain/protocol/testutil/app"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	clobtypes "github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
	perptypes "github.com/dydxprotocol/v4-chain/protocol/x/perpetuals/types"
	pricestypes "github.com/dydxprotocol/v4-chain/protocol/x/prices/types"
	vaulttypes "github.com/dydxprotocol/v4-chain/protocol/x/vault/types"
	"github.com/stretchr/testify/require"
)

func TestGetVaultWithdrawalSlippage(t *testing.T) {
	testVaultId := constants.Vault_Clob1
	testClobPair := constants.ClobPair_Eth
	testPerpetual := constants.EthUsd_20PercentInitial_10PercentMaintenance
	testMarketParam := constants.TestMarketParams[1]
	testMarketPrice := constants.TestMarketPrices[1]

	tests := map[string]struct {
		/* --- Setup --- */
		// skew.
		skewFactorPpm uint32
		// spread.
		spreadMinPpm      uint32
		spreadBufferPpm   uint32
		minPriceChangePpm uint32
		// leverage.
		leverage *big.Rat
		// total shares.
		totalShares *big.Int
		// function input.
		vaultId          vaulttypes.VaultId
		sharesToWithdraw *big.Int
		/* --- Expectations --- */
		expectedSlippage *big.Rat
		expectedErr      string
	}{
		"Success: leverage 0, skew 2, spread 0.003, withdraw 10%": {
			skewFactorPpm:     2_000_000,
			spreadMinPpm:      2_000,
			spreadBufferPpm:   1_500,
			minPriceChangePpm: 1_500,
			leverage:          big.NewRat(0, 1),
			vaultId:           testVaultId,
			sharesToWithdraw:  big.NewInt(10),
			totalShares:       big.NewInt(100),
			// no slippage when leverage is 0.
			expectedSlippage: big.NewRat(0, 1),
		},
		"Success: leverage 0.00003, skew 3, spread 0.005, withdraw 9_999_999 out of 10_000_000 shares": {
			skewFactorPpm:     3_000_000,
			spreadMinPpm:      5_000,
			spreadBufferPpm:   1_500,
			minPriceChangePpm: 1_500,
			leverage:          big.NewRat(3, 100_000),
			vaultId:           testVaultId,
			sharesToWithdraw:  big.NewInt(9_999_999),
			totalShares:       big.NewInt(10_000_000),
			// posterior_leverage = 0.00003 * 10_000_000 / (10_000_000 - 9_999_999) = 300
			// integral
			// = skew_antiderivative(skew_factor, posterior_leverage) - skew_antiderivative(skew_factor, leverage)
			// = 3 * 300^2 + 3^2 * 300^3 / 3 - (3 * 1.5^2 + 3^2 * 1.5^3 / 3)
			// ~= 81_270_000
			// estimated_slippage
			// = spread * (leverage + integral * (total_shares - shares_to_withdraw) / shares_to_withdraw)
			// = 0.005 * (0.00003 + 81_270_000 * 1 / 9_999_999)
			// ~= 0.0406352
			// slippage = min(0.0406352, leverage * imf)
			// = min(0.0406352, 0.00003 * 0.2) = 0.000006
			expectedSlippage: big.NewRat(6, 1_000_000),
		},
		"Success: leverage 0.000003, skew 3, spread 0.005, withdraw 9_999_999 out of 10_000_000 shares": {
			skewFactorPpm:     3_000_000,
			spreadMinPpm:      5_000,
			spreadBufferPpm:   1_500,
			minPriceChangePpm: 1_500,
			leverage:          big.NewRat(3, 1_000_000),
			vaultId:           testVaultId,
			sharesToWithdraw:  big.NewInt(9_999_999),
			totalShares:       big.NewInt(10_000_000),
			// posterior_leverage = 0.000003 * 10_000_000 / (10_000_000 - 9_999_999) = 30
			// integral
			// = skew_antiderivative(skew_factor, posterior_leverage) - skew_antiderivative(skew_factor, leverage)
			// = 3 * 30^2 + 3^2 * 30^3 / 3 - (3 * 1.5^2 + 3^2 * 1.5^3 / 3)
			// ~= 83_700
			// estimated_slippage
			// = spread * (leverage + integral * (total_shares - shares_to_withdraw) / shares_to_withdraw)
			// = 0.005 * (0.000003 + 83_700 * 1 / 9_999_999)
			// ~= 0.000041865
			// slippage = min(0.000041865, leverage * imf)
			// = min(0.000041865, 0.000003 * 0.2)
			// = 0.0000006
			expectedSlippage: big.NewRat(3, 5_000_000),
		},
		"Success: leverage 0.5, skew 2, spread 0.003, withdraw 10%": {
			skewFactorPpm:     2_000_000,
			spreadMinPpm:      2_000,
			spreadBufferPpm:   1_500,
			minPriceChangePpm: 1_500,
			leverage:          big.NewRat(1, 2),
			vaultId:           testVaultId,
			sharesToWithdraw:  big.NewInt(100_000),
			totalShares:       big.NewInt(1_000_000),
			// posterior_leverage = 0.5 * 1_000_000 / (1_000_000 - 100_000) = 5 / 9
			// integral
			// = skew_antiderivative(skew_factor, posterior_leverage) - skew_antiderivative(skew_factor, leverage)
			// = 2 * (5/9)^2 + 2^2 * (5/9)^3 / 3 - (2 * 0.5^2 + 2^2 * 0.5^3 / 3)
			// = 392/2187
			// estimated_slippage
			// = spread * (leverage + integral * (total_shares - shares_to_withdraw) / shares_to_withdraw)
			// = 0.003 * (0.5 + 392/2187 * 900_000 / 100_000)
			// = 1027 / 162000
			// slippage = min(0.010565935, leverage * imf)
			// = min(1027 / 162000, 0.5 * 0.2) = 1027 / 162000
			expectedSlippage: big.NewRat(1_027, 162_000),
		},
		"Success: leverage 1.5, skew 2, spread 0.003, withdraw 0.0001%": {
			skewFactorPpm:     2_000_000,
			spreadMinPpm:      2_000,
			spreadBufferPpm:   1_500,
			minPriceChangePpm: 1_500,
			leverage:          big.NewRat(3, 2),
			vaultId:           testVaultId,
			sharesToWithdraw:  big.NewInt(1),
			totalShares:       big.NewInt(1_000_000),
			// posterior_leverage = 1.5 * 1_000_000 / (1_000_000 - 1) = 500_000/333_333
			// integral
			// = skew_antiderivative(skew_factor, posterior_leverage) - skew_antiderivative(skew_factor, leverage)
			// = 2 * (500_000/333_333)^2 + 2^2 * (500_000/333_333)^3 / 3 - (2 * 1.5^2 + 2^2 * 1.5^3 / 3)
			// = 2499997000001 / 111110777778111111
			// estimated_slippage
			// = spread * (leverage + integral * (total_shares - shares_to_withdraw) / shares_to_withdraw)
			// = 0.003 * (1.5 + 2499997000001 / 111110777778111111 * 999_999 / 1)
			// = 5333326666669 / 74073925926000
			// slippage = min(0.100499, leverage * imf)
			// = min(5333326666669 / 74073925926000, 1.5 * 0.2) = 5333326666669 / 74073925926000
			// ~= 0.072000054
			expectedSlippage: big.NewRat(5333326666669, 74073925926000),
		},
		"Success: leverage 1.5, skew 3, spread 0.003, withdraw 10%": {
			skewFactorPpm:     3_000_000,
			spreadMinPpm:      2_000,
			spreadBufferPpm:   1_500,
			minPriceChangePpm: 1_500,
			leverage:          big.NewRat(3, 2),
			vaultId:           testVaultId,
			sharesToWithdraw:  big.NewInt(100_000),
			totalShares:       big.NewInt(1_000_000),
			// posterior_leverage = 1.5 * 1_000_000 / (1_000_000 - 100_000) = 5/3
			// integral
			// = skew_antiderivative(skew_factor, posterior_leverage) - skew_antiderivative(skew_factor, leverage)
			// = 3 * (5/3)^2 + 3^2 * (5/3)^3 / 3 - (3 * 1.5^2 + 3^2 * 1.5^3 / 3)
			// = 385/72
			// estimated_slippage
			// = spread * (leverage + integral * (total_shares - shares_to_withdraw) / shares_to_withdraw)
			// = 0.003 * (1.5 + 385/72 * 900_000 / 100_000)
			// = 1191/8000
			// slippage = min(1191/8000, leverage * imf)
			// = min(1191/8000, 1.5 * 0.2)
			// ~= 0.148875
			expectedSlippage: big.NewRat(1_191, 8_000),
		},
		"Success: leverage 1.5, skew 3, spread 0.003, withdraw 50%": {
			skewFactorPpm:     3_000_000,
			spreadMinPpm:      2_000,
			spreadBufferPpm:   1_500,
			minPriceChangePpm: 1_500,
			leverage:          big.NewRat(3, 2),
			vaultId:           testVaultId,
			sharesToWithdraw:  big.NewInt(500_000),
			totalShares:       big.NewInt(1_000_000),
			// posterior_leverage = 1.5 * 1_000_000 / (1_000_000 - 500_000) = 3
			// integral
			// = skew_antiderivative(skew_factor, posterior_leverage) - skew_antiderivative(skew_factor, leverage)
			// = 3 * 3^2 + 3^2 * 3^3 / 3 - (3 * 1.5^2 + 3^2 * 1.5^3 / 3)
			// = 729/8
			// estimated_slippage
			// = spread * (leverage + integral * (total_shares - shares_to_withdraw) / shares_to_withdraw)
			// = 0.003 * (1.5 + 729/8 * (1_000_000 - 500_000) / 500_000)
			// = 2223/8000
			// slippage = min(441/800, leverage * imf)
			// = min(2223/8000, 1.5 * 0.2) = 2223/8000
			// = 0.277875
			expectedSlippage: big.NewRat(2_223, 8_000),
		},
		"Success: leverage -1.5, skew 3, spread 0.003, withdraw 50%, slippage is same as when leverage is 1.5": {
			skewFactorPpm:     3_000_000,
			spreadMinPpm:      2_000,
			spreadBufferPpm:   1_500,
			minPriceChangePpm: 1_500,
			leverage:          big.NewRat(-3, 2),
			vaultId:           testVaultId,
			sharesToWithdraw:  big.NewInt(1_111),
			totalShares:       big.NewInt(2_222),
			// |leverage| = |-1.5| = 1.5
			// posterior_leverage = 1.5 * 2_222 / (2_222 - 1_111) = 3
			// integral
			// = skew_antiderivative(skew_factor, posterior_leverage) - skew_antiderivative(skew_factor, leverage)
			// = 3 * 3^2 + 3^2 * 3^3 / 3 - (3 * 1.5^2 + 3^2 * 1.5^3 / 3)
			// = 729/8
			// estimated_slippage
			// = spread * (leverage + integral * (total_shares - shares_to_withdraw) / shares_to_withdraw)
			// = 0.003 * (1.5 + 729/8 * (2_222 - 1_111) / 1_111)
			// = 2223/8000
			// slippage = min(441/800, leverage * imf)
			// = min(2223/8000, 1.5 * 0.2) = 2223/8000
			// = 0.277875
			expectedSlippage: big.NewRat(2_223, 8_000),
		},
		"Success: leverage 1.5, skew 3, spread 0.005, withdraw 50%": {
			skewFactorPpm:     3_000_000,
			spreadMinPpm:      5_000,
			spreadBufferPpm:   1_500,
			minPriceChangePpm: 1_500,
			leverage:          big.NewRat(3, 2),
			vaultId:           testVaultId,
			sharesToWithdraw:  big.NewInt(2_345_678),
			totalShares:       big.NewInt(4_691_356),
			// posterior_leverage = 1.5 * 4_691_356 / (4_691_356 - 2_345_678) = 3
			// integral
			// = skew_antiderivative(skew_factor, posterior_leverage) - skew_antiderivative(skew_factor, leverage)
			// = 3 * 3^2 + 3^2 * 3^3 / 3 - (3 * 1.5^2 + 3^2 * 1.5^3 / 3)
			// = 729/8
			// estimated_slippage
			// = spread * (leverage + integral * (total_shares - shares_to_withdraw) / shares_to_withdraw)
			// = 0.005 * (1.5 + 729/8 * (4_691_356 - 2_345_678) / 2_345_678)
			// = 741/1600
			// slippage = min(741/1600, leverage * imf)
			// = min(741/1600, 1.5 * 0.2) = 0.3
			expectedSlippage: big.NewRat(3, 10),
		},
		"Success: leverage 1.5, skew 3, spread 0.005, withdraw 100%": {
			skewFactorPpm:     3_000_000,
			spreadMinPpm:      5_000,
			spreadBufferPpm:   1_500,
			minPriceChangePpm: 1_500,
			leverage:          big.NewRat(3, 2),
			vaultId:           testVaultId,
			sharesToWithdraw:  big.NewInt(2_345_678),
			totalShares:       big.NewInt(2_345_678),
			// slippage = leverage * imf = 1.5 * 0.2 = 0.3
			expectedSlippage: big.NewRat(3, 10),
		},
		"Success: leverage 3, skew 2, spread 0.003, withdraw 1 out of 10 million shares": {
			skewFactorPpm:     2_000_000,
			spreadMinPpm:      2_000,
			spreadBufferPpm:   1_500,
			minPriceChangePpm: 1_500,
			leverage:          big.NewRat(3, 1),
			vaultId:           testVaultId,
			sharesToWithdraw:  big.NewInt(1),
			totalShares:       big.NewInt(10_000_000),
			// posterior_leverage = 3 * 10_000_000 / (10_000_000 - 1) = 3.0000003
			// integral
			// = skew_antiderivative(skew_factor, posterior_leverage) - skew_antiderivative(skew_factor, leverage)
			// = 2 * 3.0000003^2 + 2^2 * 3.0000003^3 / 3 - (2 * 3^2 + 2^2 * 3^3 / 3)
			// = 144_000_013/10_000_000_000_000
			// estimated_slippage
			// = spread * (leverage + integral * (total_shares - shares_to_withdraw) / shares_to_withdraw)
			// = 0.003 * (3 + 144000013/10000000000000 * (10_000_000 - 1) / 1)
			// = 1633333146666673 / 3703702962963000
			// slippage = min(1633333146666673 / 3703702962963000, leverage * imf)
			// = min(1633333146666673 / 3703702962963000, 3 * 0.2) = 1633333146666673 / 3703702962963000
			// = 0.4410000378
			expectedSlippage: big.NewRat(1_633_333_146_666_673, 3_703_702_962_963_000),
		},
		"Success: leverage 3, skew 2, spread 0.003, withdraw 1234 out of 12345 shares": {
			skewFactorPpm:     2_000_000,
			spreadMinPpm:      3_000,
			spreadBufferPpm:   1_500,
			minPriceChangePpm: 1_000,
			leverage:          big.NewRat(3, 1),
			vaultId:           testVaultId,
			sharesToWithdraw:  big.NewInt(1_234),
			totalShares:       big.NewInt(12_345),
			// posterior_leverage = 3 * 12345 / (12345 - 1234) = 37035/11111
			// integral
			// = skew_antiderivative(skew_factor, posterior_leverage) - skew_antiderivative(skew_factor, leverage)
			// = 2 * (37035/11111)^2 + 2^2 * (37035/11111)^3 / 3 - (2 * 3^2 + 2^2 * 3^3 / 3)
			// ~= 17.59627186
			// estimated_slippage
			// = spread * (leverage + integral * (total_shares - shares_to_withdraw) / shares_to_withdraw)
			// = 0.003 * (3 + 17.59627186 * (12345 - 1234) / 1234)
			// = 59_790_561_381 / 123_454_321_000
			// slippage = min(59_790_561_381 / 123_454_321_000, leverage * imf)
			// = min(59_790_561_381 / 123_454_321_000, 3 * 0.2)
			// = 59_790_561_381 / 123_454_321_000
			// ~= 0.484313
			expectedSlippage: big.NewRat(59_790_561_381, 123_454_321_000),
		},
		"Success: leverage 3, skew 2, spread 0.003, withdraw 50%": {
			skewFactorPpm:     2_000_000,
			spreadMinPpm:      2_000,
			spreadBufferPpm:   1_500,
			minPriceChangePpm: 1_500,
			leverage:          big.NewRat(3, 1),
			vaultId:           testVaultId,
			sharesToWithdraw:  big.NewInt(222_222),
			totalShares:       big.NewInt(444_444),
			// posterior_leverage = 3 * 444_444 / (444_444 - 222_222) = 6
			// integral
			// = skew_antiderivative(skew_factor, posterior_leverage) - skew_antiderivative(skew_factor, leverage)
			// = 2 * 6^2 + 2^2 * 6^3 / 3 - (2 * 3^2 + 2^2 * 3^3 / 3)
			// = 360 - 54
			// = 306
			// estimated_slippage
			// = spread * (leverage + integral * (total_shares - shares_to_withdraw) / shares_to_withdraw)
			// = 0.003 * (3 + 306 * (444_444 - 222_222) / 222_222)
			// = 927/1000
			// slippage = min(927/1000, leverage * imf)
			// = min(927/1000, 3 * 0.2)
			// = 0.6
			expectedSlippage: big.NewRat(3, 5),
		},
		"Success: leverage 3, skew 2, spread 0.003, withdraw 99.9999%": {
			skewFactorPpm:     2_000_000,
			spreadMinPpm:      2_000,
			spreadBufferPpm:   1_500,
			minPriceChangePpm: 1_500,
			leverage:          big.NewRat(3, 1),
			vaultId:           testVaultId,
			sharesToWithdraw:  big.NewInt(999_999),
			totalShares:       big.NewInt(1_000_000),
			// posterior_leverage = 3 * 1_000_000 / (1_000_000 - 999_999) = 3_000_000
			// integral
			// = skew_antiderivative(skew_factor, posterior_leverage) - skew_antiderivative(skew_factor, leverage)
			// = 2 * 3_000_000^2 + 2^2 * 3_000_000^3 / 3 - (2 * 3^2 + 2^2 * 3^3 / 3)
			// = 36000018e12
			// estimated_slippage
			// = spread * (leverage + integral * (total_shares - shares_to_withdraw) / shares_to_withdraw)
			// = 0.003 * (3 + 36000018e12 * (1_000_000 - 999_999) / 999_999)
			// = large number
			// slippage = min(large number, leverage * imf)
			// = min(large number, 3 * 0.2)
			// = 0.6
			expectedSlippage: big.NewRat(3, 5),
		},
		"Error: vault not found": {
			skewFactorPpm:     2_000_000,
			spreadMinPpm:      2_000,
			spreadBufferPpm:   1_500,
			minPriceChangePpm: 1_500,
			leverage:          big.NewRat(1, 1),
			vaultId:           constants.Vault_Clob0, // non-existent vault
			sharesToWithdraw:  big.NewInt(10),
			totalShares:       big.NewInt(100),
			expectedErr:       vaulttypes.ErrVaultParamsNotFound.Error(),
		},
		"Error: negative shares to withdraw": {
			skewFactorPpm:     2_000_000,
			spreadMinPpm:      2_000,
			spreadBufferPpm:   1_500,
			minPriceChangePpm: 1_500,
			leverage:          big.NewRat(1, 1),
			vaultId:           testVaultId,
			sharesToWithdraw:  big.NewInt(-1),
			totalShares:       big.NewInt(100),
			expectedErr:       vaulttypes.ErrInvalidSharesToWithdraw.Error(),
		},
		"Error: zero shares to withdraw": {
			skewFactorPpm:     2_000_000,
			spreadMinPpm:      2_000,
			spreadBufferPpm:   1_500,
			minPriceChangePpm: 1_500,
			leverage:          big.NewRat(1, 1),
			vaultId:           testVaultId,
			sharesToWithdraw:  big.NewInt(0),
			totalShares:       big.NewInt(100),
			expectedErr:       vaulttypes.ErrInvalidSharesToWithdraw.Error(),
		},
		"Error: shares to withdraw greater than total shares": {
			skewFactorPpm:     2_000_000,
			spreadMinPpm:      2_000,
			spreadBufferPpm:   1_500,
			minPriceChangePpm: 1_500,
			leverage:          big.NewRat(-1, 1),
			vaultId:           testVaultId,
			sharesToWithdraw:  big.NewInt(1_000_001),
			totalShares:       big.NewInt(1_000_000),
			expectedErr:       vaulttypes.ErrInvalidSharesToWithdraw.Error(),
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tApp := testapp.NewTestAppBuilder(t).WithGenesisDocFn(func() (genesis types.GenesisDoc) {
				genesis = testapp.DefaultGenesis()
				// Set up vault's quoting params.
				testapp.UpdateGenesisDocWithAppStateForModule(
					&genesis,
					func(genesisState *vaulttypes.GenesisState) {
						quotingParams := vaulttypes.DefaultQuotingParams()
						quotingParams.SkewFactorPpm = tc.skewFactorPpm
						quotingParams.SpreadMinPpm = tc.spreadMinPpm
						quotingParams.SpreadBufferPpm = tc.spreadBufferPpm
						genesisState.Vaults = []vaulttypes.Vault{
							{
								VaultId: testVaultId,
								VaultParams: vaulttypes.VaultParams{
									Status:        vaulttypes.VaultStatus_VAULT_STATUS_QUOTING,
									QuotingParams: &quotingParams,
								},
							},
						}
						genesisState.TotalShares = vaulttypes.BigIntToNumShares(tc.totalShares)
					},
				)
				// Set up markets.
				testMarketParam.MinPriceChangePpm = tc.minPriceChangePpm
				testapp.UpdateGenesisDocWithAppStateForModule(
					&genesis,
					func(genesisState *pricestypes.GenesisState) {
						genesisState.MarketParams = []pricestypes.MarketParam{testMarketParam}
						genesisState.MarketPrices = []pricestypes.MarketPrice{testMarketPrice}
					},
				)
				// Set up perpetuals.
				testapp.UpdateGenesisDocWithAppStateForModule(
					&genesis,
					func(genesisState *perptypes.GenesisState) {
						genesisState.LiquidityTiers = constants.LiquidityTiers
						genesisState.Perpetuals = []perptypes.Perpetual{testPerpetual}
					},
				)
				// Set up clob pairs.
				testapp.UpdateGenesisDocWithAppStateForModule(
					&genesis,
					func(genesisState *clobtypes.GenesisState) {
						genesisState.ClobPairs = []clobtypes.ClobPair{testClobPair}
					},
				)
				return genesis
			}).Build()
			ctx := tApp.InitChain()
			k := tApp.App.VaultKeeper

			slippage, err := k.GetVaultWithdrawalSlippage(
				ctx,
				tc.vaultId,
				tc.sharesToWithdraw,
				tc.totalShares,
				tc.leverage,
				&testPerpetual,
				&testMarketParam,
			)

			if tc.expectedErr == "" {
				require.NoError(t, err)
				require.Equal(t, tc.expectedSlippage, slippage)
			} else {
				require.ErrorContains(t, err, tc.expectedErr)
			}
		})
	}
}
