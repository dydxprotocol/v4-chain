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
		// total shares.
		totalShares *big.Int
		// function input.
		vaultId          vaulttypes.VaultId
		sharesToWithdraw *big.Int
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
			sharesToWithdraw:     big.NewInt(10),
			totalShares:          big.NewInt(100),
			// no slippage when leverage is 0.
			expectedSlippagePpm: big.NewInt(0),
		},
		"Success: leverage 0.00003, skew 3, spread 0.005, withdraw 9_999_999 out of 10_000_000 shares": {
			skewFactorPpm:        3_000_000,
			spreadMinPpm:         5_000,
			spreadBufferPpm:      1_500,
			minPriceChangePpm:    1_500,
			assetQuoteQuantums:   big.NewInt(999_700_000),
			positionBaseQuantums: big.NewInt(10_000), // 1 ETH
			vaultId:              testVaultId,
			sharesToWithdraw:     big.NewInt(9_999_999),
			totalShares:          big.NewInt(10_000_000),
			// open_notional = 10_000 * 10^-9 * 3_000 * 10^6 = 30_000
			// leverage = 30_000 / (999_700_000 + 3_000_000) = 0.00003
			// posterior_leverage = 0.00003 * 10_000_000 / (10_000_000 - 9_999_999) = 300
			// integral
			// = skew_antiderivative(skew_factor, posterior_leverage) - skew_antiderivative(skew_factor, leverage)
			// = 3 * 300^2 + 3^2 * 300^3 / 3 - (3 * 1.5^2 + 3^2 * 1.5^3 / 3)
			// ~= 81_270_000
			// average_skew = 81_270_000 / (300 - 0.00003) = 270_900
			// slippage = spread * (1 + average_skew) * leverage
			// = 0.005 * (1 + 270_900) * 0.00003
			// = 0.04063515
			// slippage = min(0.04063515, leverage * imf)
			// = min(0.04063515, 0.00003 * 0.2) = 0.000006
			expectedSlippagePpm: big.NewInt(6),
		},
		"Success: leverage 0.000003, skew 3, spread 0.005, withdraw 9_999_999 out of 10_000_000 shares": {
			skewFactorPpm:        3_000_000,
			spreadMinPpm:         5_000,
			spreadBufferPpm:      1_500,
			minPriceChangePpm:    1_500,
			assetQuoteQuantums:   big.NewInt(999_997_000),
			positionBaseQuantums: big.NewInt(1_000), // 1 ETH
			vaultId:              testVaultId,
			sharesToWithdraw:     big.NewInt(9_999_999),
			totalShares:          big.NewInt(10_000_000),
			// open_notional = 1_000 * 10^-9 * 3_000 * 10^6 = 3_000
			// leverage = 3_000 / (999_997_000 + 3_000) = 0.000003
			// posterior_leverage = 0.000003 * 10_000_000 / (10_000_000 - 9_999_999) = 30
			// integral
			// = skew_antiderivative(skew_factor, posterior_leverage) - skew_antiderivative(skew_factor, leverage)
			// = 3 * 30^2 + 3^2 * 30^3 / 3 - (3 * 1.5^2 + 3^2 * 1.5^3 / 3)
			// ~= 83,700
			// average_skew = 83,700 / (30 - 0.000003) = 2,790.000279
			// slippage = spread * (1 + average_skew) * leverage
			// = 0.005 * (1 + 2,790.000279) * 0.000003
			// = 0.00004186500419
			// slippage = min(0.00004186500419, leverage * imf)
			// = min(0.00004186500419, 0.000003 * 0.2)
			// ~= 0.000001 (0.0000006 gets rounded up to 0.000001)
			expectedSlippagePpm: big.NewInt(1),
		},
		"Success: leverage 0.5, skew 2, spread 0.003, withdraw 10%": {
			skewFactorPpm:        2_000_000,
			spreadMinPpm:         2_000,
			spreadBufferPpm:      1_500,
			minPriceChangePpm:    1_500,
			assetQuoteQuantums:   big.NewInt(3_000_000_000), // 3,000 USDC
			positionBaseQuantums: big.NewInt(1_000_000_000), // 1 ETH
			vaultId:              testVaultId,
			sharesToWithdraw:     big.NewInt(100_000),
			totalShares:          big.NewInt(1_000_000),
			// open_notional = 1_000_000_000 * 10^-9 * 3_000 * 10^6 = 3_000_000_000
			// leverage = 3_000_000_000 / (3_000_000_000 + 3_000_000_000) = 0.5
			// posterior_leverage = 0.5 / (1 - 0.1) = 0.5555555556 ~= 0.555556
			// integral
			// = skew_antiderivative(skew_factor, posterior_leverage) - skew_antiderivative(skew_factor, leverage)
			// = 2 * 0.555556^2 + 2^2 * 0.555556^3 / 3 - (2 * 0.5^2 + 2^2 * 0.5^3 / 3)
			// = 0.845910 - 0.666667
			// = 0.179243
			// average_skew = 0.179243 / (0.555556 - 0.5) = 3.2263481892
			// slippage = spread * (1 + average_skew) * leverage
			// = 0.003 * (1 + 3.2263481892) * 0.5
			// ~= 0.00634
			// slippage = min(0.006339, leverage * imf)
			// = min(0.00634, 0.5 * 0.2) = 0.00634
			expectedSlippagePpm: big.NewInt(6_340),
		},
		"Success: leverage 1.5, skew 2, spread 0.003, withdraw 0.0001%": {
			skewFactorPpm:        2_000_000,
			spreadMinPpm:         2_000,
			spreadBufferPpm:      1_500,
			minPriceChangePpm:    1_500,
			assetQuoteQuantums:   big.NewInt(-1_000_000_000), // -1,000 USDC
			positionBaseQuantums: big.NewInt(1_000_000_000),  // 1 ETH
			vaultId:              testVaultId,
			sharesToWithdraw:     big.NewInt(1),
			totalShares:          big.NewInt(1_000_000),
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
			sharesToWithdraw:     big.NewInt(100_000),
			totalShares:          big.NewInt(1_000_000),
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
			sharesToWithdraw:     big.NewInt(500_000),
			totalShares:          big.NewInt(1_000_000),
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
		"Success: leverage -1.5, skew 3, spread 0.003, withdraw 50%, slippage is same as when leverage is 1.5": {
			skewFactorPpm:        3_000_000,
			spreadMinPpm:         2_000,
			spreadBufferPpm:      1_500,
			minPriceChangePpm:    1_500,
			assetQuoteQuantums:   big.NewInt(5_000_000_000),  // 5,000 USDC
			positionBaseQuantums: big.NewInt(-1_000_000_000), // -1 ETH
			vaultId:              testVaultId,
			sharesToWithdraw:     big.NewInt(1_111),
			totalShares:          big.NewInt(2_222),
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
			sharesToWithdraw:     big.NewInt(2_345_678),
			totalShares:          big.NewInt(4_691_356),
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
			sharesToWithdraw:     big.NewInt(2_345_678),
			totalShares:          big.NewInt(2_345_678),
			// open_notional = 1_000_000_000 * 10^-9 * 3_000 * 10^6 = 3_000_000_000
			// leverage = 3_000_000_000 / (-1_000_000_000 + 3_000_000_000) = 1.5
			// slippage = leverage * imf = 1.5 * 0.2 = 0.3
			expectedSlippagePpm: big.NewInt(300_000),
		},
		"Success: leverage 3, skew 2, spread 0.003, withdraw 1 out of 10 million shares": {
			skewFactorPpm:        2_000_000,
			spreadMinPpm:         2_000,
			spreadBufferPpm:      1_500,
			minPriceChangePpm:    1_500,
			assetQuoteQuantums:   big.NewInt(-2_000_000_000), // -2,000 USDC
			positionBaseQuantums: big.NewInt(1_000_000_000),  // 1 ETH
			vaultId:              testVaultId,
			sharesToWithdraw:     big.NewInt(1),
			totalShares:          big.NewInt(10_000_000),
			// open_notional = 1_000_000_000 * 10^-9 * 3_000 * 10^6 = 3_000_000_000
			// leverage = 3_000_000_000 / (-2_000_000_000 + 3_000_000_000) = 3
			// posterior_leverage = 3 / (1 - 0.0000001) = 3.0000003 ~= 3.000001 (rounds up to 6 decimals)
			// integral
			// = skew_antiderivative(skew_factor, posterior_leverage) - skew_antiderivative(skew_factor, leverage)
			// = 2 * 3.000001^2 + 2^2 * 3.000001^3 / 3 - (2 * 3^2 + 2^2 * 3^3 / 3)
			// ~= 54.000050 - 54
			// = 0.000050
			// average_skew = 0.000050 / (3.000001 - 3) = 50
			// slippage = spread * (1 + average_skew) * leverage
			// = 0.003 * (1 + 50) * 3
			// = 0.459
			// slippage = min(0.459, leverage * imf)
			// = min(0.459, 3 * 0.2) = 0.459
			expectedSlippagePpm: big.NewInt(459_000),
		},
		"Success: leverage 3, skew 2, spread 0.003, withdraw 1234 out of 12345 shares": {
			skewFactorPpm:        2_000_000,
			spreadMinPpm:         3_000,
			spreadBufferPpm:      1_500,
			minPriceChangePpm:    1_000,
			assetQuoteQuantums:   big.NewInt(-2_000_000_000), // -2,000 USDC
			positionBaseQuantums: big.NewInt(1_000_000_000),  // 1 ETH
			vaultId:              testVaultId,
			sharesToWithdraw:     big.NewInt(1_234),
			totalShares:          big.NewInt(12_345),
			// open_notional = 1_000_000_000 * 10^-9 * 3_000 * 10^6 = 3_000_000_000
			// leverage = 3_000_000_000 / (-2_000_000_000 + 3_000_000_000) = 3
			// posterior_leverage = 3 * 12345 / (12345 - 1234) ~= 3.333184
			// integral
			// = skew_antiderivative(skew_factor, posterior_leverage) - skew_antiderivative(skew_factor, leverage)
			// = 2 * 3.333184^2 + 2^2 * 3.333184^3 / 3 - (2 * 3^2 + 2^2 * 3^3 / 3)
			// = 71.596311 - 54
			// = 17.596311
			// average_skew = 17.596311 / (3.333184 - 3) ~= 52.812594
			// slippage = spread * (1 + average_skew) * leverage
			// = 0.003 * (1 + 52.812594) * 3
			// = 0.484313346 ~= 0.484314
			// slippage = min(0.484314, leverage * imf)
			// = min(0.484314, 3 * 0.2) = 0.484314
			expectedSlippagePpm: big.NewInt(484_314),
		},
		"Success: leverage 3, skew 2, spread 0.003, withdraw 50%": {
			skewFactorPpm:        2_000_000,
			spreadMinPpm:         2_000,
			spreadBufferPpm:      1_500,
			minPriceChangePpm:    1_500,
			assetQuoteQuantums:   big.NewInt(-2_000_000_000), // -2,000 USDC
			positionBaseQuantums: big.NewInt(1_000_000_000),  // 1 ETH
			vaultId:              testVaultId,
			sharesToWithdraw:     big.NewInt(222_222),
			totalShares:          big.NewInt(444_444),
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
			sharesToWithdraw:     big.NewInt(999_999),
			totalShares:          big.NewInt(1_000_000),
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
			sharesToWithdraw:     big.NewInt(10),
			totalShares:          big.NewInt(100),
			expectedErr:          vaulttypes.ErrVaultParamsNotFound.Error(),
		},
		"Error: negative shares to withdraw": {
			skewFactorPpm:        2_000_000,
			spreadMinPpm:         2_000,
			spreadBufferPpm:      1_500,
			minPriceChangePpm:    1_500,
			assetQuoteQuantums:   big.NewInt(-2_000_000_000), // -2,000 USDC
			positionBaseQuantums: big.NewInt(1_000_000_000),  // 1 ETH
			vaultId:              testVaultId,
			sharesToWithdraw:     big.NewInt(-1),
			totalShares:          big.NewInt(100),
			expectedErr:          vaulttypes.ErrInvalidSharesToWithdraw.Error(),
		},
		"Error: zero shares to withdraw": {
			skewFactorPpm:        2_000_000,
			spreadMinPpm:         2_000,
			spreadBufferPpm:      1_500,
			minPriceChangePpm:    1_500,
			assetQuoteQuantums:   big.NewInt(-2_000_000_000), // -2,000 USDC
			positionBaseQuantums: big.NewInt(1_000_000_000),  // 1 ETH
			vaultId:              testVaultId,
			sharesToWithdraw:     big.NewInt(0),
			totalShares:          big.NewInt(100),
			expectedErr:          vaulttypes.ErrInvalidSharesToWithdraw.Error(),
		},
		"Error: shares to withdraw greater than total shares": {
			skewFactorPpm:        2_000_000,
			spreadMinPpm:         2_000,
			spreadBufferPpm:      1_500,
			minPriceChangePpm:    1_500,
			assetQuoteQuantums:   big.NewInt(-2_000_000_000), // -2,000 USDC
			positionBaseQuantums: big.NewInt(1_000_000_000),  // 1 ETH
			vaultId:              testVaultId,
			sharesToWithdraw:     big.NewInt(1_000_001),
			totalShares:          big.NewInt(1_000_000),
			expectedErr:          vaulttypes.ErrInvalidSharesToWithdraw.Error(),
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

			slippage, err := k.GetVaultWithdrawalSlippagePpm(ctx, tc.vaultId, tc.sharesToWithdraw)

			if tc.expectedErr == "" {
				require.NoError(t, err)
				require.Equal(t, tc.expectedSlippagePpm, slippage)
			} else {
				require.ErrorContains(t, err, tc.expectedErr)
			}
		})
	}
}
