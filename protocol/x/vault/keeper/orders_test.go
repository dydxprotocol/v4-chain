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
	}{
		"Two Vaults, Both Positive Shares": {
			vaultIds: []vaulttypes.VaultId{
				constants.Vault_Clob_0,
				constants.Vault_Clob_1,
			},
			totalShares: []*big.Int{
				big.NewInt(1_000),
				big.NewInt(200),
			},
		},
		"Two Vaults, One Positive Shares, One Zero Shares": {
			vaultIds: []vaulttypes.VaultId{
				constants.Vault_Clob_0,
				constants.Vault_Clob_1,
			},
			totalShares: []*big.Int{
				big.NewInt(1_000),
				big.NewInt(0),
			},
		},
		"Two Vaults, Both Zero Shares": {
			vaultIds: []vaulttypes.VaultId{
				constants.Vault_Clob_0,
				constants.Vault_Clob_1,
			},
			totalShares: []*big.Int{
				big.NewInt(0),
				big.NewInt(0),
			},
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
										Quantums: dtypes.NewInt(1_000_000_000), // 1,000 USDC
									},
								},
							}
						}
						genesisState.Subaccounts = subaccounts
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

			// Simulate vault orders placed in last block.
			numPreviousOrders := 0
			for i, vaultId := range tc.vaultIds {
				if tc.totalShares[i].Sign() > 0 {
					orders, err := tApp.App.VaultKeeper.GetVaultClobOrders(
						ctx.WithBlockHeight(ctx.BlockHeight()-1),
						vaultId,
					)
					require.NoError(t, err)
					for _, order := range orders {
						err := tApp.App.ClobKeeper.HandleMsgPlaceOrder(
							ctx,
							clobtypes.NewMsgPlaceOrder(*order),
						)
						require.NoError(t, err)
					}
					numPreviousOrders += len(orders)
				}
			}
			require.Len(t, tApp.App.ClobKeeper.GetAllStatefulOrders(ctx), numPreviousOrders)

			// Refresh all vault orders.
			tApp.App.VaultKeeper.RefreshAllVaultOrders(ctx)

			// Check orders are as expected, i.e. orders from last block have been
			// cancelled and orders from this block have been placed.
			numExpectedOrders := 0
			allExpectedOrderIds := make(map[clobtypes.OrderId]bool)
			for i, vaultId := range tc.vaultIds {
				if tc.totalShares[i].Sign() > 0 {
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
		vaultAssetQuantums *big.Int
		// Vault inventory.
		vaultInventory *big.Int
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
				Layers:                 2,       // 2 layers
				SpreadMinPpm:           3_000,   // 30 bps
				SpreadBufferPpm:        1_500,   // 15 bps
				SkewFactorPpm:          500_000, // 0.5
				OrderSizePpm:           100_000, // 10%
				OrderExpirationSeconds: 2,       // 2 seconds
			},
			vaultId:            constants.Vault_Clob_0,
			vaultAssetQuantums: big.NewInt(1_000_000_000), // 1,000 USDC
			vaultInventory:     big.NewInt(0),
			clobPair:           constants.ClobPair_Btc,
			marketParam:        constants.TestMarketParams[0],
			marketPrice:        constants.TestMarketPrices[0],
			perpetual:          constants.BtcUsd_0DefaultFunding_10AtomicResolution,
			// To calculate order subticks:
			// 1. spread = max(spread_min, spread_buffer + min_price_change)
			// 2. leverage = open_notional / equity
			// 3. leverage_i = leverage +/- i * order_size_pct (- for ask and + for bid)
			// 4. skew_i = -leverage_i * spread * skew_factor
			// 5. a_i = oracle_price * (1 + skew_i) * (1 + spread)^{i+1}
			//    b_i = oracle_price * (1 + skew_i) / (1 + spread)^{i+1}
			// 6. subticks needs to be a multiple of subticks_per_tick (round up for asks, round down for bids)
			// To calculate size of each order
			// 1. `order_size_ppm * equity / oracle_price`.
			expectedOrderSubticks: []uint64{
				// spreadPpm = max(3_000, 1_500 + 50) = 3_000
				// spread = 0.003
				// leverage = 0 / 1_000 = 0
				// oracleSubticks = 5_000_000_000 * 10^(-5 - (-8) + (-10) - (-6)) = 5 * 10^8
				// leverage_0 = leverage = 0
				// skew_0 = -0 * 3_000 * 0.5 = 0
				// a_0 = 5 * 10^8 * (1 + 0) * (1 + 0.003)^1 = 501_500_000
				501_500_000,
				// b_0 = 5 * 10^8 * (1 + 0) / (1 + 0.003)^1 = 498_504_486
				// round down to nearest multiple of subticks_per_tick=5.
				498_504_485,
				// leverage_1 = leverage - 0.1 = -0.1
				// skew_1 = 0.1 * 0.003 * 0.5 = 0.00015
				// a_1 = 5 * 10^8 * (1 + 0.00015) * (1 + 0.003)^2 = 503_079_950.675
				// round up to nearest multiple of subticks_per_tick=5.
				503_079_955,
				// leverage_1 = leverage + 0.1 = 0.1
				// skew_1 = -0.1 * 0.003 * 0.5 = -0.00015
				// b_2 = 5 * 10^8 * (1 - 0.00015) / (1 + 0.003)^2 ~= 496_938_894.184
				// round down to nearest multiple of subticks_per_tick=5.
				496_938_890,
			},
			// order_size = 10% * 1_000 / 50_000 = 0.002
			// order_size_base_quantums = 0.002 * 10^10 = 20_000_000
			expectedOrderQuantums: []uint64{
				20_000_000,
				20_000_000,
				20_000_000,
				20_000_000,
			},
		},
		"Success - Get orders from Vault for Clob Pair 1": {
			vaultParams: vaulttypes.Params{
				Layers:                 3,       // 3 layers
				SpreadMinPpm:           3_000,   // 30 bps
				SpreadBufferPpm:        4_500,   // 15 bps
				SkewFactorPpm:          700_000, // 0.7
				OrderSizePpm:           50_000,  // 5%
				OrderExpirationSeconds: 4,       // 4 seconds
			},
			vaultId:            constants.Vault_Clob_1,
			vaultAssetQuantums: big.NewInt(2_000_000_000), // 2,000 USDC
			vaultInventory:     big.NewInt(-1_000),
			clobPair:           constants.ClobPair_Eth,
			marketParam:        constants.TestMarketParams[1],
			marketPrice:        constants.TestMarketPrices[1],
			perpetual:          constants.EthUsd_0DefaultFunding_9AtomicResolution,
			// To calculate order subticks:
			// 1. spread = max(spread_min, spread_buffer + min_price_change)
			// 2. leverage = open_notional / equity
			// 3. leverage_i = leverage +/- i * order_size_pct (- for ask and + for bid)
			// 4. skew_i = -leverage_i * spread * skew_factor
			// 5. a_i = oracle_price * (1 + skew_i) * (1 + spread)^{i+1}
			//    b_i = oracle_price * (1 + skew_i) / (1 + spread)^{i+1}
			// 6. subticks needs to be a multiple of subticks_per_tick (round up for asks, round down for bids)
			// To calculate size of each order
			// 1. `order_size_ppm * equity / oracle_price`.
			expectedOrderSubticks: []uint64{
				// spreadPpm = max(3_000, 4_500 + 50) = 4_550
				// spread = 0.00455
				// open_notional = -1_000 * 10^-9 * 3_000 * 10^6 = -3_000
				// leverage = -3_000 / (2_000_000_000 - 3_000) = -3/1_999_997
				// oracleSubticks = 3_000_000_000 * 10^(-6 - (-9) + (-9) - (-6)) = 3 * 10^9
				// leverage_0 = leverage - 0 * 0.05 = -3/1_999_997
				// skew_0 = 3 / 1_999_997 * 0.00455 * 0.7
				// a_0 = 3 * 10^9 * (1 + skew_0) * (1 + 0.00455)^1 = 3_013_650_014.4
				// round up to nearest multiple of subticks_per_tick=1_000.
				3_013_651_000,
				// b_0 = 3 * 10^9 * (1 + skew_0) / (1 + 0.00455)^1 = 2_986_411_840.46
				// round down to nearest multiple of subticks_per_tick=1_000.
				2_986_411_000,
				// leverage_1 = leverage - 1 * 0.05
				// skew_1 = -leverage_1 * 0.00455 * 0.7
				// a_1 = 3 * 10^9 * (1 + skew_1) * (1 + 0.00455)^2 ~= 3_027_844_229.378
				// round up to nearest multiple of subticks_per_tick=1_000.
				3_027_845_000,
				// leverage_1 = leverage + 1 * 0.05
				// skew_1 = -leverage_1 * 0.00455 * 0.7
				// b_1 = 3 * 10^9 * (1 + skew_1) / (1 + 0.00455)^2 ~= 2_972_411_780.773
				// round down to nearest multiple of subticks_per_tick=5.
				2_972_411_000,
				// leverage_2 = leverage - 2 * 0.05
				// skew_2 = -leverage_2 * 0.00455 * 0.7
				// a_2 = 3 * 10^9 * (1 + skew_2) * (1 + 0.00455)^3 ~= 3_042_105_221.627
				// round up to nearest multiple of subticks_per_tick=1_000.
				3_042_106_000,
				// leverage_2 = leverage + 2 * 0.05
				// skew_2 = -leverage_2 * 0.00455 * 0.7
				// b_2 = 3 * 10^9 * (1 + skew_2) / (1 + 0.00455)^3 ~= 2958477277.194
				// round down to nearest multiple of subticks_per_tick=1_000.
				2_958_477_000,
			},
			// order_size = 5% * 1999.997 / 3000 ~= 0.03333328333333334
			// order_size_base_quantums = 0.03333328333333334 * 10^9 = 333_33_283
			// round down to nearest multiple of step_base_quantums=1_000.
			expectedOrderQuantums: []uint64{
				333_33_000,
				333_33_000,
				333_33_000,
				333_33_000,
				333_33_000,
				333_33_000,
			},
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
			vaultParams:        vaulttypes.DefaultParams(),
			vaultId:            constants.Vault_Clob_0,
			vaultAssetQuantums: big.NewInt(0),
			vaultInventory:     big.NewInt(0),
			clobPair:           constants.ClobPair_Btc,
			marketParam:        constants.TestMarketParams[0],
			marketPrice:        constants.TestMarketPrices[0],
			perpetual:          constants.BtcUsd_0DefaultFunding_10AtomicResolution,
			expectedErr:        vaulttypes.ErrNonPositiveEquity,
		},
		"Error - Vault equity is negative": {
			vaultParams:        vaulttypes.DefaultParams(),
			vaultId:            constants.Vault_Clob_0,
			vaultAssetQuantums: big.NewInt(5_000_000), // 5 USDC
			vaultInventory:     big.NewInt(-10_000_000),
			clobPair:           constants.ClobPair_Btc,
			marketParam:        constants.TestMarketParams[0],
			marketPrice:        constants.TestMarketPrices[0],
			perpetual:          constants.BtcUsd_0DefaultFunding_10AtomicResolution,
			expectedErr:        vaulttypes.ErrNonPositiveEquity,
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
						if tc.vaultAssetQuantums != nil && tc.vaultAssetQuantums.Sign() != 0 {
							assetPositions = append(
								assetPositions,
								&satypes.AssetPosition{
									AssetId:  assettypes.AssetUsdc.Id,
									Quantums: dtypes.NewIntFromBigInt(tc.vaultAssetQuantums),
								},
							)
						}
						perpPositions := []*satypes.PerpetualPosition{}
						if tc.vaultInventory != nil && tc.vaultInventory.Sign() != 0 {
							perpPositions = append(
								perpPositions,
								&satypes.PerpetualPosition{
									PerpetualId: tc.perpetual.Params.Id,
									Quantums:    dtypes.NewIntFromBigInt(tc.vaultInventory),
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
			for i := uint32(0); i < params.Layers; i++ {
				expectedOrders = append(
					expectedOrders,
					// ask.
					buildVaultClobOrder(
						uint8(i),
						clobtypes.Order_SIDE_SELL,
						tc.expectedOrderQuantums[2*i],
						tc.expectedOrderSubticks[2*i],
					),
					// bid.
					buildVaultClobOrder(
						uint8(i),
						clobtypes.Order_SIDE_BUY,
						tc.expectedOrderQuantums[2*i+1],
						tc.expectedOrderSubticks[2*i+1],
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
		// block height.
		blockHeight int64
		// layer.
		layer uint8

		/* --- Expectations --- */
		// Expected client ID.
		expectedClientId uint32
	}{
		"Buy, Block Height Odd, Layer 1": {
			side:             clobtypes.Order_SIDE_BUY, // 0<<31
			blockHeight:      1,                        // 1<<30
			layer:            1,                        // 1<<22
			expectedClientId: 0<<31 | 1<<30 | 1<<22,
		},
		"Buy, Block Height Even, Layer 1": {
			side:             clobtypes.Order_SIDE_BUY, // 0<<31
			blockHeight:      2,                        // 0<<30
			layer:            1,                        // 1<<22
			expectedClientId: 0<<31 | 0<<30 | 1<<22,
		},
		"Sell, Block Height Odd, Layer 2": {
			side:             clobtypes.Order_SIDE_SELL, // 1<<31
			blockHeight:      1,                         // 1<<30
			layer:            2,                         // 2<<22
			expectedClientId: 1<<31 | 1<<30 | 2<<22,
		},
		"Sell, Block Height Even, Layer 2": {
			side:             clobtypes.Order_SIDE_SELL, // 1<<31
			blockHeight:      2,                         // 0<<30
			layer:            2,                         // 2<<22
			expectedClientId: 1<<31 | 0<<30 | 2<<22,
		},
		"Buy, Block Height Even, Layer Max Uint8": {
			side:             clobtypes.Order_SIDE_BUY, // 0<<31
			blockHeight:      123456,                   // 0<<30
			layer:            math.MaxUint8,            // 255<<22
			expectedClientId: 0<<31 | 0<<30 | 255<<22,
		},
		"Sell, Block Height Odd, Layer 0": {
			side:             clobtypes.Order_SIDE_SELL, // 1<<31
			blockHeight:      12345654321,               // 1<<30
			layer:            0,                         // 0<<22
			expectedClientId: 1<<31 | 1<<30 | 0<<22,
		},
		"Sell, Block Height Odd (negative), Layer 202": {
			side: clobtypes.Order_SIDE_SELL, // 1<<31
			// Negative block height shouldn't happen but blockHeight
			// is represented as int64.
			blockHeight:      -678987, // 1<<30
			layer:            202,     // 202<<22
			expectedClientId: 1<<31 | 1<<30 | 202<<22,
		},
		"Buy, Block Height Even (zero), Layer 157": {
			side:             clobtypes.Order_SIDE_SELL, // 1<<31
			blockHeight:      0,                         // 0<<30
			layer:            157,                       // 157<<22
			expectedClientId: 1<<31 | 0<<30 | 157<<22,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tApp := testapp.NewTestAppBuilder(t).Build()
			ctx := tApp.InitChain()

			clientId := tApp.App.VaultKeeper.GetVaultClobOrderClientId(
				ctx.WithBlockHeight(tc.blockHeight),
				tc.side,
				tc.layer,
			)
			require.Equal(t, tc.expectedClientId, clientId)
		})
	}
}
