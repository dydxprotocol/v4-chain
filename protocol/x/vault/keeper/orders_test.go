package keeper_test

import (
	"math"
	"math/big"
	"testing"

	"github.com/cometbft/cometbft/types"
	"github.com/dydxprotocol/v4-chain/protocol/dtypes"
	testapp "github.com/dydxprotocol/v4-chain/protocol/testutil/app"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	assettypes "github.com/dydxprotocol/v4-chain/protocol/x/assets/types"
	clobtypes "github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
	perptypes "github.com/dydxprotocol/v4-chain/protocol/x/perpetuals/types"
	pricestypes "github.com/dydxprotocol/v4-chain/protocol/x/prices/types"
	satypes "github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"
	vaulttypes "github.com/dydxprotocol/v4-chain/protocol/x/vault/types"
	"github.com/stretchr/testify/require"
)

func TestRefreshAllVaultOrders(t *testing.T) {
	tests := map[string]struct {
		// Vault IDs.
		vaultIds []vaulttypes.VaultId
		// Total Shares of each vault ID above.
		totalShares []*big.Int
		// Asset quantums of each vault ID above.
		assetQuantums []*big.Int
		// Activation threshold (quote quantums) of vaults.
		activationThresholdQuoteQuantums *big.Int
	}{
		"Two Vaults, Both Positive Shares, Both above Activation Threshold": {
			vaultIds: []vaulttypes.VaultId{
				constants.Vault_Clob_0,
				constants.Vault_Clob_1,
			},
			totalShares: []*big.Int{
				big.NewInt(1_000),
				big.NewInt(200),
			},
			assetQuantums: []*big.Int{
				big.NewInt(1_000_000_000), // 1,000 USDC
				big.NewInt(1_000_000_001),
			},
			activationThresholdQuoteQuantums: big.NewInt(1_000_000_000),
		},
		"Two Vaults, One Positive Shares, One Zero Shares, Both above Activation Threshold": {
			vaultIds: []vaulttypes.VaultId{
				constants.Vault_Clob_0,
				constants.Vault_Clob_1,
			},
			totalShares: []*big.Int{
				big.NewInt(1_000),
				big.NewInt(0),
			},
			assetQuantums: []*big.Int{
				big.NewInt(1_000_000_000), // 1,000 USDC
				big.NewInt(1_000_000_001),
			},
			activationThresholdQuoteQuantums: big.NewInt(1_000_000_000),
		},
		"Two Vaults, Both Zero Shares, Both above Activation Threshold": {
			vaultIds: []vaulttypes.VaultId{
				constants.Vault_Clob_0,
				constants.Vault_Clob_1,
			},
			totalShares: []*big.Int{
				big.NewInt(0),
				big.NewInt(0),
			},
			assetQuantums: []*big.Int{
				big.NewInt(1_000_000_000), // 1,000 USDC
				big.NewInt(1_000_000_001),
			},
			activationThresholdQuoteQuantums: big.NewInt(1_000_000_000),
		},
		"Two Vaults, Both Positive Shares, Only One above Activation Threshold": {
			vaultIds: []vaulttypes.VaultId{
				constants.Vault_Clob_0,
				constants.Vault_Clob_1,
			},
			totalShares: []*big.Int{
				big.NewInt(1_000),
				big.NewInt(200),
			},
			assetQuantums: []*big.Int{
				big.NewInt(1_000_000_000),
				big.NewInt(999_999_999),
			},
			activationThresholdQuoteQuantums: big.NewInt(1_000_000_000),
		},
		"Two Vaults, Both Positive Shares, Both below Activation Threshold": {
			vaultIds: []vaulttypes.VaultId{
				constants.Vault_Clob_0,
				constants.Vault_Clob_1,
			},
			totalShares: []*big.Int{
				big.NewInt(1_000),
				big.NewInt(200),
			},
			assetQuantums: []*big.Int{
				big.NewInt(123_456_788),
				big.NewInt(123_456_787),
			},
			activationThresholdQuoteQuantums: big.NewInt(123_456_789),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// Initialize tApp and ctx (in deliverTx mode).
			tApp := testapp.NewTestAppBuilder(t).WithGenesisDocFn(func() (genesis types.GenesisDoc) {
				genesis = testapp.DefaultGenesis()
				// Initialize each vault with quote quantums to be able to place orders.
				testapp.UpdateGenesisDocWithAppStateForModule(
					&genesis,
					func(genesisState *satypes.GenesisState) {
						subaccounts := make([]satypes.Subaccount, len(tc.vaultIds))
						for i, vaultId := range tc.vaultIds {
							subaccounts[i] = satypes.Subaccount{
								Id: vaultId.ToSubaccountId(),
								AssetPositions: []*satypes.AssetPosition{
									{
										AssetId:  assettypes.AssetUsdc.Id,
										Quantums: dtypes.NewIntFromBigInt(tc.assetQuantums[i]),
									},
								},
							}
						}
						genesisState.Subaccounts = subaccounts
					},
				)
				testapp.UpdateGenesisDocWithAppStateForModule(
					&genesis,
					func(genesisState *vaulttypes.GenesisState) {
						vaultParams := genesisState.Params
						vaultParams.ActivationThresholdQuoteQuantums = dtypes.NewIntFromBigInt(
							tc.activationThresholdQuoteQuantums,
						)
						genesisState.Params = vaultParams
					},
				)
				return genesis
			}).Build()
			ctx := tApp.InitChain().WithIsCheckTx(false)

			// Set total shares for each vault ID.
			for i, vaultId := range tc.vaultIds {
				err := tApp.App.VaultKeeper.SetTotalShares(
					ctx,
					vaultId,
					vaulttypes.BigIntToNumShares(tc.totalShares[i]),
				)
				require.NoError(t, err)
			}

			// Check that there's no stateful orders yet.
			allStatefulOrders := tApp.App.ClobKeeper.GetAllStatefulOrders(ctx)
			require.Len(t, allStatefulOrders, 0)

			// Refresh all vault orders.
			tApp.App.VaultKeeper.RefreshAllVaultOrders(ctx)

			// Check orders are as expected, i.e. orders from last block have been
			// cancelled and orders from this block have been placed.
			numExpectedOrders := 0
			allExpectedOrderIds := make(map[clobtypes.OrderId]bool)
			for i, vaultId := range tc.vaultIds {
				if tc.totalShares[i].Sign() > 0 && tc.assetQuantums[i].Cmp(tc.activationThresholdQuoteQuantums) >= 0 {
					expectedOrders, err := tApp.App.VaultKeeper.GetVaultClobOrders(ctx, vaultId)
					require.NoError(t, err)
					numExpectedOrders += len(expectedOrders)
					for _, order := range expectedOrders {
						allExpectedOrderIds[order.OrderId] = true
					}
				}
			}
			allStatefulOrders = tApp.App.ClobKeeper.GetAllStatefulOrders(ctx)
			require.Len(t, allStatefulOrders, numExpectedOrders)
			for _, order := range allStatefulOrders {
				require.True(t, allExpectedOrderIds[order.OrderId])
			}
		})
	}
}

func TestRefreshVaultClobOrders(t *testing.T) {
	tests := map[string]struct {
		/* --- Setup --- */
		// Vault ID.
		vaultId vaulttypes.VaultId

		/* --- Expectations --- */
		expectedErr error
	}{
		"Success - Refresh Orders from Vault for Clob Pair 0": {
			vaultId: constants.Vault_Clob_0,
		},
		"Error - Refresh Orders from Vault for Clob Pair 4321 (non-existent clob pair)": {
			vaultId: vaulttypes.VaultId{
				Type:   vaulttypes.VaultType_VAULT_TYPE_CLOB,
				Number: 4321,
			},
			expectedErr: vaulttypes.ErrClobPairNotFound,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// Initialize tApp and ctx (in deliverTx mode).
			tApp := testapp.NewTestAppBuilder(t).WithGenesisDocFn(func() (genesis types.GenesisDoc) {
				genesis = testapp.DefaultGenesis()
				// Initialize vault with quote quantums to be able to place orders.
				testapp.UpdateGenesisDocWithAppStateForModule(
					&genesis,
					func(genesisState *satypes.GenesisState) {
						genesisState.Subaccounts = []satypes.Subaccount{
							{
								Id: tc.vaultId.ToSubaccountId(),
								AssetPositions: []*satypes.AssetPosition{
									{
										AssetId:  assettypes.AssetUsdc.Id,
										Quantums: dtypes.NewInt(1_000_000_000), // 1,000 USDC
									},
								},
							},
						}
					},
				)
				return genesis
			}).Build()
			ctx := tApp.InitChain().WithIsCheckTx(false)

			// Check that there's no stateful orders yet.
			allStatefulOrders := tApp.App.ClobKeeper.GetAllStatefulOrders(ctx)
			require.Len(t, allStatefulOrders, 0)

			// Refresh vault orders.
			err := tApp.App.VaultKeeper.RefreshVaultClobOrders(ctx, tc.vaultId)
			allStatefulOrders = tApp.App.ClobKeeper.GetAllStatefulOrders(ctx)
			if tc.expectedErr != nil {
				// Check that the error is as expected.
				require.ErrorContains(t, err, tc.expectedErr.Error())
				// Check that there's no stateful orders.
				require.Len(t, allStatefulOrders, 0)
				return
			} else {
				// Check that there's no error.
				require.NoError(t, err)
				// Check that the number of orders is as expected.
				params := tApp.App.VaultKeeper.GetParams(ctx)
				require.Len(t, allStatefulOrders, int(params.Layers*2))
				// Check that the orders are as expected.
				expectedOrders, err := tApp.App.VaultKeeper.GetVaultClobOrders(ctx, tc.vaultId)
				require.NoError(t, err)
				for i := uint32(0); i < params.Layers*2; i++ {
					require.Equal(t, *expectedOrders[i], allStatefulOrders[i])
				}
			}
		})
	}
}

func TestGetVaultClobOrders(t *testing.T) {
	tests := map[string]struct {
		/* --- Setup --- */
		// Vault params.
		vaultParams vaulttypes.Params
		// Vault ID.
		vaultId vaulttypes.VaultId
		// Vault asset.
		vaultAssetQuoteQuantums *big.Int
		// Vault inventory.
		vaultInventoryBaseQuantums *big.Int
		// Clob pair.
		clobPair clobtypes.ClobPair
		// Market param.
		marketParam pricestypes.MarketParam
		// Market price.
		marketPrice pricestypes.MarketPrice
		// Perpetual.
		perpetual perptypes.Perpetual

		/* --- Expectations --- */
		expectedOrderSubticks []uint64
		expectedOrderQuantums []uint64
		expectedErr           error
	}{
		"Success - Get orders from Vault for Clob Pair 0": {
			vaultParams: vaulttypes.Params{
				Layers:                           2,       // 2 layers
				SpreadMinPpm:                     3_123,   // 31.23 bps
				SpreadBufferPpm:                  1_500,   // 15 bps
				SkewFactorPpm:                    554_321, // 0.554321
				OrderSizePctPpm:                  100_000, // 10%
				OrderExpirationSeconds:           2,       // 2 seconds
				ActivationThresholdQuoteQuantums: dtypes.NewInt(1_000_000_000),
			},
			vaultId:                    constants.Vault_Clob_0,
			vaultAssetQuoteQuantums:    big.NewInt(1_000_000_000), // 1,000 USDC
			vaultInventoryBaseQuantums: big.NewInt(0),
			clobPair:                   constants.ClobPair_Btc,
			marketParam:                constants.TestMarketParams[0],
			marketPrice: pricestypes.MarketPrice{
				Id:       0,
				Exponent: -5,
				Price:    5_000_000, // $50
			},
			perpetual: constants.BtcUsd_0DefaultFunding_10AtomicResolution,
			// To calculate order subticks:
			// 1. spread = max(spread_min, spread_buffer + min_price_change)
			// 2. leverage = open_notional / equity
			// 3. leverage_i = leverage +/- i * order_size_pct (- for ask and + for bid)
			// 4. skew_i = -leverage_i * spread * skew_factor
			// 5. a_i = max(oracle_price * (1 + skew_i + spread * {i+1}), oracle_price)
			//    b_i = min(oracle_price * (1 + skew_i - spread * {i+1}), oracle_price)
			// 6. subticks needs to be a multiple of subticks_per_tick (round up for asks, round down for bids)
			// To calculate size of each order
			// 1. `order_size_pct_ppm * equity / oracle_price`.
			expectedOrderSubticks: []uint64{
				// spreadPpm = max(3_123, 1_500 + 50) = 3_123
				// spread = 0.003123
				// leverage = 0 / 1_000 = 0
				// oracleSubticks = 5_000_000_000 * 10^(-5 - (-8) + (-10) - (-6)) = 5e8
				// leverage_0 = leverage = 0
				// skew_0 = -0 * 3_123 * 0.554321 = 0
				// a_0 = 5e5 * (1 + 0 + 0.003123*1) = 501_561.5 = 501_565 (rounded up to 5)
				501_565,
				// b_0 = 5e5 * (1 + 0 - 0.003123*1) = 498_438.5 = 498435 (rounded down to 5)
				498_435,
				// leverage_1 = leverage - 0.1 = -0.1
				// skew_1 = 0.1 * 0.003123 * 0.554321 ~= 0.000173
				// a_1 = 5e5 * (1 + 0.000173 + 0.003123*2) = 503209.5 ~= 503_210 (rounded up to 5)
				503_210,
				// leverage_1 = leverage + 0.1 = 0.1
				// skew_1 = -0.1 * 0.003123 * 0.554321 = -0.000173
				// b_2 = 5e5 * (1 - 0.000173 - 0.003123*2) = 496790.5 ~= 496_790 (rounded down to 5)
				496_790,
			},
			// order_size = 10% * $1_000 / $50 = 2
			// order_size_base_quantums = 2 * 10^10 = 20_000_000_000
			expectedOrderQuantums: []uint64{
				20_000_000_000,
				20_000_000_000,
				20_000_000_000,
				20_000_000_000,
			},
		},
		"Success - Get orders from Vault for Clob Pair 1, bids bounded by oracle price.": {
			vaultParams: vaulttypes.Params{
				Layers:                           3,       // 3 layers
				SpreadMinPpm:                     3_000,   // 30 bps
				SpreadBufferPpm:                  8_500,   // 85 bps
				SkewFactorPpm:                    900_000, // 0.9
				OrderSizePctPpm:                  200_000, // 20%
				OrderExpirationSeconds:           4,       // 4 seconds
				ActivationThresholdQuoteQuantums: dtypes.NewInt(1_000_000_000),
			},
			vaultId:                    constants.Vault_Clob_1,
			vaultAssetQuoteQuantums:    big.NewInt(2_000_000_000), // 2,000 USDC
			vaultInventoryBaseQuantums: big.NewInt(-500_000_000),  // -0.5 ETH
			clobPair:                   constants.ClobPair_Eth,
			marketParam:                constants.TestMarketParams[1],
			marketPrice:                constants.TestMarketPrices[1],
			perpetual:                  constants.EthUsd_0DefaultFunding_9AtomicResolution,
			// To calculate order subticks:
			// 1. spread = max(spread_min, spread_buffer + min_price_change)
			// 2. leverage = open_notional / equity
			// 3. leverage_i = leverage +/- i * order_size_pct (- for ask and + for bid)
			// 4. skew_i = -leverage_i * spread * skew_factor
			// 5. a_i = max(oracle_price * (1 + skew_i + spread*{i+1}), oracle_price)
			//    b_i = min(oracle_price * (1 + skew_i - spread*{i+1}), oracle_price)
			// 6. subticks needs to be a multiple of subticks_per_tick (round up for asks, round down for bids)
			// To calculate size of each order
			// 1. `order_size_pct_ppm * equity / oracle_price`.
			expectedOrderSubticks: []uint64{
				// spreadPpm = max(3_000, 8_500 + 50) = 8_550
				// spread = 0.00855
				// open_notional = -500_000_000 * 10^-9 * 3_000 * 10^6 = -1_500_000_000
				// leverage = -1_500_000_000 / (2_000_000_000 - 1_500_000_000) = -3
				// oracleSubticks = 3_000_000_000 * 10^(-6 - (-9) + (-9) - (-6)) = 3e9
				// leverage_0 = leverage - 0 * 0.2 = -3
				// skew_0 = 3 * 0.00855 * 0.9
				// a_0 = 3e9 * (1 + skew_0 + 0.00855*1) = 3_094_905_000
				// a_0 = max(a_0, oracle_price) = 3_094_905_000
				3_094_905_000,
				// b_0 = 3e9 * (1 + skew_0 - 0.00855*1) = 3_043_605_000
				// b_0 = min(b_0, oracle_price) = 3e9 (bound)
				3_000_000_000,
				// leverage_1 = leverage - 1 * 0.2
				// skew_1 = -leverage_1 * 0.00855 * 0.9
				// a_1 = 3e9 * (1 + skew_1 + 0.00855*2) = 3_125_172_000
				// a_1 = max(a_1, oracle_price) = 3_125_172_000
				3_125_172_000,
				// leverage_1 = leverage + 1 * 0.2
				// skew_1 = -leverage_1 * 0.00855 * 0.9
				// b_1 = 3e9 * (1 + skew_1 - 0.00855*2) = 3_013_338_000
				// b_1 = min(b_1, oracle_price) = 3e9 (bound)
				3_000_000_000,
				// leverage_2 = leverage - 2 * 0.2
				// skew_2 = -leverage_2 * 0.00855 * 0.9
				// a_2 = 3e9 * (1 + skew_2 + 0.00855*3) = 3_155_439_000
				// a_2 = max(a_2, oracle_price) = 3_155_439_000
				3_155_439_000,
				// leverage_2 = leverage + 2 * 0.2
				// skew_2 = -leverage_2 * 0.00855 * 0.9
				// b_2 = 3e9 * (1 + skew_2 - 0.00855*3) = 2_983_071_000
				// b_2 = min(b_2, oracle_price) = 2_983_071_000
				2_983_071_000,
			},
			// order_size = 20% * 500 / 3000 ~= 0.0333333333
			// order_size_base_quantums = 0.0333333333e9 ~= 33_333_333.33
			// round down to nearest multiple of step_base_quantums=1_000.
			expectedOrderQuantums: []uint64{
				33_333_000,
				33_333_000,
				33_333_000,
				33_333_000,
				33_333_000,
				33_333_000,
			},
		},
		"Success - Get orders from Vault for Clob Pair 1, asks bounded by oracle price.": {
			vaultParams: vaulttypes.Params{
				Layers:                           2,         // 2 layers
				SpreadMinPpm:                     3_000,     // 30 bps
				SpreadBufferPpm:                  1_500,     // 15 bps
				SkewFactorPpm:                    500_000,   // 0.5
				OrderSizePctPpm:                  1_000_000, // 100%
				OrderExpirationSeconds:           4,         // 4 seconds
				ActivationThresholdQuoteQuantums: dtypes.NewInt(1_000_000_000),
			},
			vaultId:                    constants.Vault_Clob_1,
			vaultAssetQuoteQuantums:    big.NewInt(-2_000_000_000), // -2,000 USDC
			vaultInventoryBaseQuantums: big.NewInt(1_000_000_000),  // 1 ETH
			clobPair:                   constants.ClobPair_Eth,
			marketParam:                constants.TestMarketParams[1],
			marketPrice:                constants.TestMarketPrices[1],
			perpetual:                  constants.EthUsd_0DefaultFunding_9AtomicResolution,
			// To calculate order subticks:
			// 1. spread = max(spread_min, spread_buffer + min_price_change)
			// 2. leverage = open_notional / equity
			// 3. leverage_i = leverage +/- i * order_size_pct (- for ask and + for bid)
			// 4. skew_i = -leverage_i * spread * skew_factor
			// 5. a_i = max(oracle_price * (1 + skew_i + spread*{i+1}), oracle_price)
			//    b_i = min(oracle_price * (1 + skew_i - spread*{i+1}), oracle_price)
			// 6. subticks needs to be a multiple of subticks_per_tick (round up for asks, round down for bids)
			// To calculate size of each order
			// 1. `order_size_pct_ppm * equity / oracle_price`.
			expectedOrderSubticks: []uint64{
				// spreadPpm = max(3_000, 1_500 + 50) = 3_000
				// spread = 0.003
				// open_notional = 1_000_000_000 * 10^-9 * 3_000 * 10^6 = 3_000_000_000
				// leverage = 3_000_000_000 / (-2_000_000_000 + 3_000_000_000) = 3
				// oracleSubticks = 3_000_000_000 * 10^(-6 - (-9) + (-9) - (-6)) = 3e9
				// leverage_0 = leverage - 0 * 1 = 3
				// skew_0 = -3 * 0.003 * 0.5
				// a_0 = 3e9 * (1 + skew_0 + 0.003*1) = 2_995_500_000
				// a_0 = max(a_0, oracle_price) = 3e9 (bound)
				3_000_000_000,
				// b_0 = 3e9 * (1 + skew_0 - 0.003*1) = 2_977_500_000
				// b_0 = min(b_0, oracle_price) = 2_977_500_000
				2_977_500_000,
				// leverage_1 = leverage - 1 * 1 = 2
				// skew_1 = -2 * 0.003 * .5
				// a_1 = 3e9 * (1 + skew_1 + 0.003*2) = 3_009_000_000
				// a_1 = max(a_1, oracle_price) = 3_009_000_000
				3_009_000_000,
				// leverage_1 = leverage + 1 * 1 = 4
				// skew_1 = -4 * 0.003 * .5
				// b_1 = 3e9 * (1 + skew_1 - 0.003*2) = 2_964_000_000
				// b_1 = min(b_1, oracle_price) = 2_964_000_000
				2_964_000_000,
			},
			// order_size = 100% * 1000 / 3000 ~= 0.333333333
			// order_size_base_quantums = 0.333333333e9 ~= 333_333_333.33
			// round down to nearest multiple of step_base_quantums=1_000.
			expectedOrderQuantums: []uint64{
				333_333_000,
				333_333_000,
				333_333_000,
				333_333_000,
			},
		},
		"Success - Get orders from Vault for Clob Pair 1, No Orders due to Zero Order Size": {
			vaultParams: vaulttypes.Params{
				Layers:                           2,       // 2 layers
				SpreadMinPpm:                     3_000,   // 30 bps
				SpreadBufferPpm:                  1_500,   // 15 bps
				SkewFactorPpm:                    500_000, // 0.5
				OrderSizePctPpm:                  1_000,   // 0.1%
				OrderExpirationSeconds:           2,       // 2 seconds
				ActivationThresholdQuoteQuantums: dtypes.NewInt(1_000_000_000),
			},
			vaultId:                    constants.Vault_Clob_1,
			vaultAssetQuoteQuantums:    big.NewInt(1_000_000), // 1 USDC
			vaultInventoryBaseQuantums: big.NewInt(0),
			clobPair:                   constants.ClobPair_Eth,
			marketParam:                constants.TestMarketParams[1],
			marketPrice:                constants.TestMarketPrices[1],
			perpetual:                  constants.EthUsd_0DefaultFunding_9AtomicResolution,
			expectedOrderSubticks:      []uint64{},
			// order_size = 0.1% * 1 / 3_000 ~= 0.00000033333
			// order_size_base_quantums = 0.000033333e9 = 333
			// round down to nearest multiple of step_base_quantums=1_000.
			// order size is 0.
			expectedOrderQuantums: []uint64{},
		},
		"Error - Clob Pair doesn't exist": {
			vaultParams: vaulttypes.DefaultParams(),
			vaultId:     constants.Vault_Clob_0,
			clobPair:    constants.ClobPair_Eth,
			marketParam: constants.TestMarketParams[1],
			marketPrice: constants.TestMarketPrices[1],
			perpetual:   constants.EthUsd_NoMarginRequirement,
			expectedErr: vaulttypes.ErrClobPairNotFound,
		},
		"Error - Vault equity is zero": {
			vaultParams:                vaulttypes.DefaultParams(),
			vaultId:                    constants.Vault_Clob_0,
			vaultAssetQuoteQuantums:    big.NewInt(0),
			vaultInventoryBaseQuantums: big.NewInt(0),
			clobPair:                   constants.ClobPair_Btc,
			marketParam:                constants.TestMarketParams[0],
			marketPrice:                constants.TestMarketPrices[0],
			perpetual:                  constants.BtcUsd_0DefaultFunding_10AtomicResolution,
			expectedErr:                vaulttypes.ErrNonPositiveEquity,
		},
		"Error - Vault equity is negative": {
			vaultParams:                vaulttypes.DefaultParams(),
			vaultId:                    constants.Vault_Clob_0,
			vaultAssetQuoteQuantums:    big.NewInt(5_000_000), // 5 USDC
			vaultInventoryBaseQuantums: big.NewInt(-10_000_000),
			clobPair:                   constants.ClobPair_Btc,
			marketParam:                constants.TestMarketParams[0],
			marketPrice:                constants.TestMarketPrices[0],
			perpetual:                  constants.BtcUsd_0DefaultFunding_10AtomicResolution,
			expectedErr:                vaulttypes.ErrNonPositiveEquity,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// Initialize tApp and ctx.
			tApp := testapp.NewTestAppBuilder(t).WithGenesisDocFn(func() (genesis types.GenesisDoc) {
				genesis = testapp.DefaultGenesis()
				// Initialize prices module with test market param and market price.
				testapp.UpdateGenesisDocWithAppStateForModule(
					&genesis,
					func(genesisState *pricestypes.GenesisState) {
						genesisState.MarketParams = []pricestypes.MarketParam{tc.marketParam}
						genesisState.MarketPrices = []pricestypes.MarketPrice{tc.marketPrice}
					},
				)
				// Initialize perpetuals module with test perpetual.
				testapp.UpdateGenesisDocWithAppStateForModule(
					&genesis,
					func(genesisState *perptypes.GenesisState) {
						genesisState.LiquidityTiers = constants.LiquidityTiers
						genesisState.Perpetuals = []perptypes.Perpetual{tc.perpetual}
					},
				)
				// Initialize clob module with test clob pair.
				testapp.UpdateGenesisDocWithAppStateForModule(
					&genesis,
					func(genesisState *clobtypes.GenesisState) {
						genesisState.ClobPairs = []clobtypes.ClobPair{tc.clobPair}
					},
				)
				// Initialize vault module with test params.
				testapp.UpdateGenesisDocWithAppStateForModule(
					&genesis,
					func(genesisState *vaulttypes.GenesisState) {
						genesisState.Params = tc.vaultParams
					},
				)
				// Initialize subaccounts module with vault's equity and inventory.
				testapp.UpdateGenesisDocWithAppStateForModule(
					&genesis,
					func(genesisState *satypes.GenesisState) {
						assetPositions := []*satypes.AssetPosition{}
						if tc.vaultAssetQuoteQuantums != nil && tc.vaultAssetQuoteQuantums.Sign() != 0 {
							assetPositions = append(
								assetPositions,
								&satypes.AssetPosition{
									AssetId:  assettypes.AssetUsdc.Id,
									Quantums: dtypes.NewIntFromBigInt(tc.vaultAssetQuoteQuantums),
								},
							)
						}
						perpPositions := []*satypes.PerpetualPosition{}
						if tc.vaultInventoryBaseQuantums != nil && tc.vaultInventoryBaseQuantums.Sign() != 0 {
							perpPositions = append(
								perpPositions,
								&satypes.PerpetualPosition{
									PerpetualId: tc.perpetual.Params.Id,
									Quantums:    dtypes.NewIntFromBigInt(tc.vaultInventoryBaseQuantums),
								},
							)
						}
						genesisState.Subaccounts = []satypes.Subaccount{
							{
								Id:                 tc.vaultId.ToSubaccountId(),
								AssetPositions:     assetPositions,
								PerpetualPositions: perpPositions,
							},
						}
					},
				)
				return genesis
			}).Build()
			ctx := tApp.InitChain()

			// Get vault orders.
			orders, err := tApp.App.VaultKeeper.GetVaultClobOrders(ctx, tc.vaultId)
			if tc.expectedErr != nil {
				require.ErrorContains(t, err, tc.expectedErr.Error())
				return
			}
			require.NoError(t, err)

			// Get expected orders.
			params := tApp.App.VaultKeeper.GetParams(ctx)
			buildVaultClobOrder := func(
				layer uint8,
				side clobtypes.Order_Side,
				quantums uint64,
				subticks uint64,
			) *clobtypes.Order {
				return &clobtypes.Order{
					OrderId: clobtypes.OrderId{
						SubaccountId: *tc.vaultId.ToSubaccountId(),
						ClientId:     tApp.App.VaultKeeper.GetVaultClobOrderClientId(ctx, side, layer),
						OrderFlags:   clobtypes.OrderIdFlags_LongTerm,
						ClobPairId:   tc.vaultId.Number,
					},
					Side:     side,
					Quantums: quantums,
					Subticks: subticks,
					GoodTilOneof: &clobtypes.Order_GoodTilBlockTime{
						GoodTilBlockTime: uint32(ctx.BlockTime().Unix()) + params.OrderExpirationSeconds,
					},
				}
			}
			expectedOrders := make([]*clobtypes.Order, 0)
			for i := 0; i < len(tc.expectedOrderQuantums); i += 2 {
				expectedOrders = append(
					expectedOrders,
					// ask.
					buildVaultClobOrder(
						uint8(i/2),
						clobtypes.Order_SIDE_SELL,
						tc.expectedOrderQuantums[i],
						tc.expectedOrderSubticks[i],
					),
					// bid.
					buildVaultClobOrder(
						uint8(i/2),
						clobtypes.Order_SIDE_BUY,
						tc.expectedOrderQuantums[i+1],
						tc.expectedOrderSubticks[i+1],
					),
				)
			}

			// Compare expected orders with actual orders.
			require.Equal(
				t,
				expectedOrders,
				orders,
			)
		})
	}
}

func TestGetVaultClobOrderClientId(t *testing.T) {
	tests := map[string]struct {
		/* --- Setup --- */
		// side.
		side clobtypes.Order_Side
		// layer.
		layer uint8

		/* --- Expectations --- */
		// Expected client ID.
		expectedClientId uint32
	}{
		"Buy, Layer 1": {
			side:             clobtypes.Order_SIDE_BUY, // 0<<31
			layer:            1,                        // 1<<23
			expectedClientId: 0<<31 | 1<<23,
		},
		"Sell, Layer 2": {
			side:             clobtypes.Order_SIDE_SELL, // 1<<31
			layer:            2,                         // 2<<23
			expectedClientId: 1<<31 | 2<<23,
		},
		"Buy, Layer Max Uint8": {
			side:             clobtypes.Order_SIDE_BUY, // 0<<31
			layer:            math.MaxUint8,            // 255<<23
			expectedClientId: 0<<31 | 255<<23,
		},
		"Sell, Layer 0": {
			side:             clobtypes.Order_SIDE_SELL, // 1<<31
			layer:            0,                         // 0<<23
			expectedClientId: 1<<31 | 0<<23,
		},
		"Sell, Layer 202": {
			side:             clobtypes.Order_SIDE_SELL, // 1<<31
			layer:            202,                       // 202<<23
			expectedClientId: 1<<31 | 202<<23,
		},
		"Buy, Layer 157": {
			side:             clobtypes.Order_SIDE_SELL, // 1<<31
			layer:            157,                       // 157<<23
			expectedClientId: 1<<31 | 157<<23,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tApp := testapp.NewTestAppBuilder(t).Build()
			ctx := tApp.InitChain()

			clientId := tApp.App.VaultKeeper.GetVaultClobOrderClientId(
				ctx,
				tc.side,
				tc.layer,
			)
			require.Equal(t, tc.expectedClientId, clientId)
		})
	}
}
