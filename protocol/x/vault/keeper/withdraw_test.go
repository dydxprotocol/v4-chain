package keeper_test

import (
	"math/big"
	"testing"

	"github.com/cometbft/cometbft/types"
	"github.com/dydxprotocol/v4-chain/protocol/dtypes"
	testapp "github.com/dydxprotocol/v4-chain/protocol/testutil/app"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	testutil "github.com/dydxprotocol/v4-chain/protocol/testutil/util"
	assettypes "github.com/dydxprotocol/v4-chain/protocol/x/assets/types"
	clobtypes "github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
	perptypes "github.com/dydxprotocol/v4-chain/protocol/x/perpetuals/types"
	pricestypes "github.com/dydxprotocol/v4-chain/protocol/x/prices/types"
	satypes "github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"
	vaulttypes "github.com/dydxprotocol/v4-chain/protocol/x/vault/types"
	"github.com/stretchr/testify/require"
)

func TestGetVaultWithdrawalSlippagePpm(t *testing.T) {
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
		assetQuoteQuantums   *big.Int
		positionBaseQuantums *big.Int
		// function input.
		vaultId              vaulttypes.VaultId
		withdrawalPortionPpm *big.Int
		/* --- Expectations --- */
		expectedSlippagePpm *big.Int
		expectedErr         string
	}{
		"Success: leverage 0, skew 2, spread 0.003, withdraw 10%": {
			skewFactorPpm:        2_000_000,
			spreadMinPpm:         2_000,
			spreadBufferPpm:      1_500,
			minPriceChangePpm:    1_500,
			assetQuoteQuantums:   big.NewInt(1_000_000_000), // 1,000 USDC
			positionBaseQuantums: big.NewInt(0),
			vaultId:              testVaultId,
			withdrawalPortionPpm: big.NewInt(100_000), // 10%
			// no slippage when leverage is 0.
			expectedSlippagePpm: big.NewInt(0),
		},
		"Success: leverage 1.5, skew 2, spread 0.003, withdraw the smallest portion 0.0001%": {
			skewFactorPpm:        2_000_000,
			spreadMinPpm:         2_000,
			spreadBufferPpm:      1_500,
			minPriceChangePpm:    1_500,
			assetQuoteQuantums:   big.NewInt(-1_000_000_000), // -1,000 USDC
			positionBaseQuantums: big.NewInt(1_000_000_000),  // 1 ETH
			vaultId:              testVaultId,
			withdrawalPortionPpm: big.NewInt(1), // 0.0001%
			// open_notional = 1_000_000_000 * 10^-9 * 3_000 * 10^6 = 3_000_000_000
			// leverage = 3_000_000_000 / (-1_000_000_000 + 3_000_000_000) = 1.5
			// posterior_leverage = 1.5 / (1 - 0.000001) = 1.5000015 ~= 1.500002
			// integral
			// = skew_antiderivative(skew_factor, posterior_leverage) - skew_antiderivative(skew_factor, leverage)
			// = 2 * 1.500002^2 + 2^2 * 1.500002^3 / 3 - (2 * 1.5^2 + 2^2 * 1.5^3 / 3)
			// = 9.000032 - 9
			// = 0.000032
			// average_skew = 0.000032 / (1.500002 - 1.5) = 16
			// slippage = spread * (1 + average_skew) * leverage
			// = 0.003 * (1 + 16) * 1.5
			// = 0.0765
			// slippage = min(0.0765, leverage * imf)
			// = min(0.0765, 1.5 * 0.2) = 0.0765
			expectedSlippagePpm: big.NewInt(76_500),
		},
		"Success: leverage 1.5, skew 3, spread 0.003, withdraw 10%": {
			skewFactorPpm:        3_000_000,
			spreadMinPpm:         2_000,
			spreadBufferPpm:      1_500,
			minPriceChangePpm:    1_500,
			assetQuoteQuantums:   big.NewInt(-1_000_000_000), // -1,000 USDC
			positionBaseQuantums: big.NewInt(1_000_000_000),  // 1 ETH
			vaultId:              testVaultId,
			withdrawalPortionPpm: big.NewInt(100_000), // 10%
			// open_notional = 1_000_000_000 * 10^-9 * 3_000 * 10^6 = 3_000_000_000
			// leverage = 3_000_000_000 / (-1_000_000_000 + 3_000_000_000) = 1.5
			// posterior_leverage = 1.5 / (1 - 0.1) ~= 1.666667
			// integral
			// = skew_antiderivative(skew_factor, posterior_leverage) - skew_antiderivative(skew_factor, leverage)
			// = 3 * 1.666667^2 + 3^2 * 1.666667^3 / 3 - (3 * 1.5^2 + 3^2 * 1.5^3 / 3)
			// = 22.222234 - 16.875
			// = 5.347234
			// average_skew = 5.347234 / (1.666667 - 1.5) ~= 32.083340
			// slippage = spread * (1 + average_skew) * leverage
			// = 0.003 * (1 + 32.083340) * 1.5
			// ~= 0.148876
			// slippage = min(0.148876, leverage * imf)
			// = min(0.148876, 1.5 * 0.2) = 0.148876
			expectedSlippagePpm: big.NewInt(148_876),
		},
		"Success: leverage 1.5, skew 3, spread 0.003, withdraw 50%": {
			skewFactorPpm:        3_000_000,
			spreadMinPpm:         2_000,
			spreadBufferPpm:      1_500,
			minPriceChangePpm:    1_500,
			assetQuoteQuantums:   big.NewInt(-1_000_000_000), // -1,000 USDC
			positionBaseQuantums: big.NewInt(1_000_000_000),  // 1 ETH
			vaultId:              testVaultId,
			withdrawalPortionPpm: big.NewInt(500_000), // 50%
			// open_notional = 1_000_000_000 * 10^-9 * 3_000 * 10^6 = 3_000_000_000
			// leverage = 3_000_000_000 / (-1_000_000_000 + 3_000_000_000) = 1.5
			// posterior_leverage = 1.5 / (1 - 0.5) = 3
			// integral
			// = skew_antiderivative(skew_factor, posterior_leverage) - skew_antiderivative(skew_factor, leverage)
			// = 3 * 3^2 + 3^2 * 3^3 / 3 - (3 * 1.5^2 + 3^2 * 1.5^3 / 3)
			// = 108 - 16.875
			// = 91.125
			// average_skew = 91.125 / (3 - 1.5) = 60.75
			// slippage = spread * (1 + average_skew) * leverage
			// = 0.003 * (1 + 60.75) * 1.5
			// = 0.277875
			// slippage = min(0.277875, leverage * imf)
			// = min(0.277875, 1.5 * 0.2) = 0.277875
			expectedSlippagePpm: big.NewInt(277_875),
		},
		"Success: leverage -1.5, skew 3, spread 0.003, withdraw 50%. slippage is same as when leverage is 1.5": {
			skewFactorPpm:        3_000_000,
			spreadMinPpm:         2_000,
			spreadBufferPpm:      1_500,
			minPriceChangePpm:    1_500,
			assetQuoteQuantums:   big.NewInt(5_000_000_000),  // 5,000 USDC
			positionBaseQuantums: big.NewInt(-1_000_000_000), // -1 ETH
			vaultId:              testVaultId,
			withdrawalPortionPpm: big.NewInt(500_000), // 50%
			// open_notional = -1_000_000_000 * 10^-9 * 3_000 * 10^6 = -3_000_000_000
			// |leverage| = |-3_000_000_000 / (5_000_000_000 + -3_000_000_000)| = |-1.5| = 1.5
			// posterior_leverage = 1.5 / (1 - 0.5) = 3
			// integral
			// = skew_antiderivative(skew_factor, posterior_leverage) - skew_antiderivative(skew_factor, leverage)
			// = 3 * 3^2 + 3^2 * 3^3 / 3 - (3 * 1.5^2 + 3^2 * 1.5^3 / 3)
			// = 108 - 16.875
			// = 91.125
			// average_skew = 91.125 / (3 - 1.5) = 60.75
			// slippage = spread * (1 + average_skew) * leverage
			// = 0.003 * (1 + 60.75) * 1.5
			// = 0.277875
			// slippage = min(0.277875, leverage * imf)
			// = min(0.277875, 1.5 * 0.2) = 0.277875
			expectedSlippagePpm: big.NewInt(277_875),
		},
		"Success: leverage 1.5, skew 3, spread 0.005, withdraw 50%": {
			skewFactorPpm:        3_000_000,
			spreadMinPpm:         5_000,
			spreadBufferPpm:      1_500,
			minPriceChangePpm:    1_500,
			assetQuoteQuantums:   big.NewInt(-1_000_000_000), // -1,000 USDC
			positionBaseQuantums: big.NewInt(1_000_000_000),  // 1 ETH
			vaultId:              testVaultId,
			withdrawalPortionPpm: big.NewInt(500_000), // 50%
			// open_notional = 1_000_000_000 * 10^-9 * 3_000 * 10^6 = 3_000_000_000
			// leverage = 3_000_000_000 / (-1_000_000_000 + 3_000_000_000) = 1.5
			// posterior_leverage = 1.5 / (1 - 0.5) = 3
			// integral
			// = skew_antiderivative(skew_factor, posterior_leverage) - skew_antiderivative(skew_factor, leverage)
			// = 3 * 3^2 + 3^2 * 3^3 / 3 - (3 * 1.5^2 + 3^2 * 1.5^3 / 3)
			// = 108 - 16.875
			// = 91.125
			// average_skew = 91.125 / (3 - 1.5) = 60.75
			// slippage = spread * (1 + average_skew) * leverage
			// = 0.005 * (1 + 60.75) * 1.5
			// = 0.463125
			// slippage = min(0.463125, leverage * imf)
			// = min(0.463125, 1.5 * 0.2) = 0.3
			expectedSlippagePpm: big.NewInt(300_000),
		},
		"Success: leverage 1.5, skew 3, spread 0.005, withdraw 100%": {
			skewFactorPpm:        3_000_000,
			spreadMinPpm:         5_000,
			spreadBufferPpm:      1_500,
			minPriceChangePpm:    1_500,
			assetQuoteQuantums:   big.NewInt(-1_000_000_000), // -1,000 USDC
			positionBaseQuantums: big.NewInt(1_000_000_000),  // 1 ETH
			vaultId:              testVaultId,
			withdrawalPortionPpm: big.NewInt(1_000_000), // 100%
			// open_notional = 1_000_000_000 * 10^-9 * 3_000 * 10^6 = 3_000_000_000
			// leverage = 3_000_000_000 / (-1_000_000_000 + 3_000_000_000) = 1.5
			// slippage = leverage * imf = 1.5 * 0.2 = 0.3
			expectedSlippagePpm: big.NewInt(300_000),
		},
		"Success: leverage 3, skew 2, spread 0.003, withdraw the smallest portion 0.0001%": {
			skewFactorPpm:        2_000_000,
			spreadMinPpm:         2_000,
			spreadBufferPpm:      1_500,
			minPriceChangePpm:    1_500,
			assetQuoteQuantums:   big.NewInt(-2_000_000_000), // -2,000 USDC
			positionBaseQuantums: big.NewInt(1_000_000_000),  // 1 ETH
			vaultId:              testVaultId,
			withdrawalPortionPpm: big.NewInt(1), // 0.0001%
			// open_notional = 1_000_000_000 * 10^-9 * 3_000 * 10^6 = 3_000_000_000
			// leverage = 3_000_000_000 / (-2_000_000_000 + 3_000_000_000) = 3
			// posterior_leverage = 3 / (1 - 0.000001) ~= 3.000004
			// integral
			// = skew_antiderivative(skew_factor, posterior_leverage) - skew_antiderivative(skew_factor, leverage)
			// = 2 * 3.000004^2 + 2^2 * 3.000004^3 / 3 - (2 * 3^2 + 2^2 * 3^3 / 3)
			// ~= 54.000194 - 54
			// = 0.000194
			// average_skew = 0.000194 / (3.000004 - 3) = 48.5
			// slippage = spread * (1 + average_skew) * leverage
			// = 0.003 * (1 + 48.5) * 3
			// = 0.4455
			// slippage = min(0.4455, leverage * imf)
			// = min(0.4455, 3 * 0.2) = 0.4455
			expectedSlippagePpm: big.NewInt(445_500),
		},
		"Success: leverage 3, skew 2, spread 0.003, withdraw 10%": {
			skewFactorPpm:        2_000_000,
			spreadMinPpm:         3_000,
			spreadBufferPpm:      1_500,
			minPriceChangePpm:    1_000,
			assetQuoteQuantums:   big.NewInt(-2_000_000_000), // -2,000 USDC
			positionBaseQuantums: big.NewInt(1_000_000_000),  // 1 ETH
			vaultId:              testVaultId,
			withdrawalPortionPpm: big.NewInt(100_000), // 10%
			// open_notional = 1_000_000_000 * 10^-9 * 3_000 * 10^6 = 3_000_000_000
			// leverage = 3_000_000_000 / (-2_000_000_000 + 3_000_000_000) = 3
			// posterior_leverage = 3 / (1 - 0.1) = 3.333333
			// integral
			// = skew_antiderivative(skew_factor, posterior_leverage) - skew_antiderivative(skew_factor, leverage)
			// = 2 * 3.333333^2 + 2^2 * 3.333333^3 / 3 - (2 * 3^2 + 2^2 * 3^3 / 3)
			// = 71.604919 - 54
			// = 17.604919
			// average_skew = 17.604919 / (3.333333 - 3) ~= 52.814810
			// slippage = spread * (1 + average_skew) * leverage
			// = 0.003 * (1 + 52.814810) * 3
			// = 0.48433329
			// round up to 0.484334
			// slippage = min(0.484334, leverage * imf)
			// = min(0.484334, 3 * 0.2) = 0.484334
			expectedSlippagePpm: big.NewInt(484_334),
		},
		"Success: leverage 3, skew 2, spread 0.003, withdraw 50%": {
			skewFactorPpm:        2_000_000,
			spreadMinPpm:         2_000,
			spreadBufferPpm:      1_500,
			minPriceChangePpm:    1_500,
			assetQuoteQuantums:   big.NewInt(-2_000_000_000), // -2,000 USDC
			positionBaseQuantums: big.NewInt(1_000_000_000),  // 1 ETH
			vaultId:              testVaultId,
			withdrawalPortionPpm: big.NewInt(500_000), // 50%
			// open_notional = 1_000_000_000 * 10^-9 * 3_000 * 10^6 = 3_000_000_000
			// leverage = 3_000_000_000 / (-2_000_000_000 + 3_000_000_000) = 3
			// posterior_leverage = 3 / (1 - 0.5) = 6
			// integral
			// = skew_antiderivative(skew_factor, posterior_leverage) - skew_antiderivative(skew_factor, leverage)
			// = 2 * 6^2 + 2^2 * 6^3 / 3 - (2 * 3^2 + 2^2 * 3^3 / 3)
			// = 360 - 54
			// = 306
			// average_skew = 306 / (6 - 3) = 102
			// slippage = spread * (1 + average_skew) * leverage
			// = 0.003 * (1 + 102) * 3
			// = 0.927
			// slippage = min(0.4455, leverage * imf)
			// = min(0.927, 3 * 0.2) = 0.6
			expectedSlippagePpm: big.NewInt(600_000),
		},
		"Success: leverage 3, skew 2, spread 0.003, withdraw 99.9999%": {
			skewFactorPpm:        2_000_000,
			spreadMinPpm:         2_000,
			spreadBufferPpm:      1_500,
			minPriceChangePpm:    1_500,
			assetQuoteQuantums:   big.NewInt(-2_000_000_000), // -2,000 USDC
			positionBaseQuantums: big.NewInt(1_000_000_000),  // 1 ETH
			vaultId:              testVaultId,
			withdrawalPortionPpm: big.NewInt(999_999), // 99.9999%
			// open_notional = 1_000_000_000 * 10^-9 * 3_000 * 10^6 = 3_000_000_000
			// leverage = 3_000_000_000 / (-2_000_000_000 + 3_000_000_000) = 3
			// posterior_leverage = 3 / (1 - 0.999999) = 3_000_000
			// integral
			// = skew_antiderivative(skew_factor, posterior_leverage) - skew_antiderivative(skew_factor, leverage)
			// = 2 * 3_000_000^2 + 2^2 * 3_000_000^3 / 3 - (2 * 3^2 + 2^2 * 3^3 / 3)
			// ~= 3.6 * 10^19 - 54
			// ~= 3.6e19
			// average_skew = 3.6e19 / (3_000_000 - 3) ~= 1.2e13
			// slippage = spread * (1 + average_skew) * leverage
			// = 0.003 * (1 + 1.2e13) * 3
			// ~= 108000000000
			// slippage = min(108000000000, leverage * imf)
			// = min(108000000000, 3 * 0.2) = 0.6
			expectedSlippagePpm: big.NewInt(600_000),
		},
		"Error: vault not found": {
			skewFactorPpm:        2_000_000,
			spreadMinPpm:         2_000,
			spreadBufferPpm:      1_500,
			minPriceChangePpm:    1_500,
			assetQuoteQuantums:   big.NewInt(-2_000_000_000), // -2,000 USDC
			positionBaseQuantums: big.NewInt(1_000_000_000),  // 1 ETH
			vaultId:              constants.Vault_Clob0,      // non-existent vault
			withdrawalPortionPpm: big.NewInt(500_000),        // 50%
			expectedErr:          vaulttypes.ErrVaultParamsNotFound.Error(),
		},
		"Error: negative withdrawal portion": {
			skewFactorPpm:        2_000_000,
			spreadMinPpm:         2_000,
			spreadBufferPpm:      1_500,
			minPriceChangePpm:    1_500,
			assetQuoteQuantums:   big.NewInt(-2_000_000_000), // -2,000 USDC
			positionBaseQuantums: big.NewInt(1_000_000_000),  // 1 ETH
			vaultId:              testVaultId,
			withdrawalPortionPpm: big.NewInt(-1), // negative
			expectedErr:          vaulttypes.ErrInvalidWithdrawalPortion.Error(),
		},
		"Error: zero withdrawal portion": {
			skewFactorPpm:        2_000_000,
			spreadMinPpm:         2_000,
			spreadBufferPpm:      1_500,
			minPriceChangePpm:    1_500,
			assetQuoteQuantums:   big.NewInt(-2_000_000_000), // -2,000 USDC
			positionBaseQuantums: big.NewInt(1_000_000_000),  // 1 ETH
			vaultId:              testVaultId,
			withdrawalPortionPpm: big.NewInt(0), // 0%
			expectedErr:          vaulttypes.ErrInvalidWithdrawalPortion.Error(),
		},
		"Error: withdrawal portion greater than 1": {
			skewFactorPpm:        2_000_000,
			spreadMinPpm:         2_000,
			spreadBufferPpm:      1_500,
			minPriceChangePpm:    1_500,
			assetQuoteQuantums:   big.NewInt(-2_000_000_000), // -2,000 USDC
			positionBaseQuantums: big.NewInt(1_000_000_000),  // 1 ETH
			vaultId:              testVaultId,
			withdrawalPortionPpm: big.NewInt(1_000_001), // 100.0001%
			expectedErr:          vaulttypes.ErrInvalidWithdrawalPortion.Error(),
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
				// Set up vault asset and perpetual positions.
				testapp.UpdateGenesisDocWithAppStateForModule(
					&genesis,
					func(genesisState *satypes.GenesisState) {
						genesisState.Subaccounts = []satypes.Subaccount{
							{
								Id: tc.vaultId.ToSubaccountId(),
								AssetPositions: []*satypes.AssetPosition{
									{
										AssetId:  assettypes.AssetUsdc.Id,
										Quantums: dtypes.NewIntFromBigInt(tc.assetQuoteQuantums),
									},
								},
								PerpetualPositions: []*satypes.PerpetualPosition{
									testutil.CreateSinglePerpetualPosition(
										testPerpetual.Params.Id,
										tc.positionBaseQuantums,
										big.NewInt(0),
										big.NewInt(0),
									),
								},
							},
						}
					},
				)
				return genesis
			}).Build()
			ctx := tApp.InitChain()
			k := tApp.App.VaultKeeper

			slippage, err := k.GetVaultWithdrawalSlippagePpm(ctx, tc.vaultId, tc.withdrawalPortionPpm)

			if tc.expectedErr == "" {
				require.NoError(t, err)
				require.Equal(t, tc.expectedSlippagePpm, slippage)
			} else {
				require.ErrorContains(t, err, tc.expectedErr)
			}
		})
	}
}
