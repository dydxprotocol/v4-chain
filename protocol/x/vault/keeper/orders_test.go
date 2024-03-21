package keeper_test

import (
	"math"
	"testing"

	"github.com/cometbft/cometbft/types"
	testapp "github.com/dydxprotocol/v4-chain/protocol/testutil/app"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	clobtypes "github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
	perptypes "github.com/dydxprotocol/v4-chain/protocol/x/perpetuals/types"
	pricestypes "github.com/dydxprotocol/v4-chain/protocol/x/prices/types"
	satypes "github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"
	vaulttypes "github.com/dydxprotocol/v4-chain/protocol/x/vault/types"
	"github.com/stretchr/testify/require"
)

// TODO (TRA-118): store vault strategy constants in x/vault state.
const (
	numLayers                          = uint8(2)
	minBaseSpreadPpm                   = uint32(3_000) // 30 bps
	baseSpreadMinPriceChangePremiumPpm = uint32(1_500) // 15 bps
	orderExpirationSeconds             = uint32(5)     // 5 seconds
)

func TestGetVaultClobOrders(t *testing.T) {
	tests := map[string]struct {
		/* --- Setup --- */
		// Vault ID.
		vaultId vaulttypes.VaultId
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
			vaultId:     constants.Vault_Clob_0,
			clobPair:    constants.ClobPair_Btc,
			marketParam: constants.TestMarketParams[0],
			marketPrice: constants.TestMarketPrices[0],
			perpetual:   constants.BtcUsd_0DefaultFunding_10AtomicResolution,
			// To calculate order subticks:
			// 1. spreadPpm = max(minBaseSpreadPpm, baseSpreadMinPriceChangePremiumPpm + minPriceChangePpm)
			// 2. priceSubticks = marketPrice.Price * 10^(marketPrice.Exponent - quantumConversionExponent +
			//                    baseAtomicResolution - quoteAtomicResolution)
			// 3. askSubticks at layer i = priceSubticks * (1 + spread)^i
			//    bidSubticks at layer i = priceSubticks * (1 - spread)^i
			// 4. subticks needs to be a multiple of subtickPerTick (round up for asks, round down for bids)
			expectedOrderSubticks: []uint64{
				// spreadPpm = max(3_000, 1_500 + 50) = 3_000
				// priceSubticks = 5_000_000_000 * 10^(-5 - (-8) + (-10) - (-6)) = 5 * 10^8
				// a_1 = 5 * 10^8 * (1 + 0.003)^1 = 501_500_000
				501_500_000,
				// b_1 = 5 * 10^8 * (1 - 0.003)^1 = 498_500_000
				498_500_000,
				// a_2 = 5 * 10^8 * (1 + 0.003)^2 = 503_004_500
				503_004_500,
				// b_2 = 5 * 10^8 * (1 - 0.003)^2 = 497_004_500
				497_004_500,
			},
			expectedOrderQuantums: []uint64{ // TODO (TRA-144): Implement order size
				5,
				5,
				5,
				5,
			},
		},
		"Success - Get orders from Vault for Clob Pair 1": {
			vaultId:  constants.Vault_Clob_1,
			clobPair: constants.ClobPair_Eth,
			marketParam: pricestypes.MarketParam{
				Id:                 constants.TestMarketParams[1].Id,
				Pair:               constants.TestMarketParams[1].Pair,
				Exponent:           constants.TestMarketParams[1].Exponent,
				MinExchanges:       constants.TestMarketParams[1].MinExchanges,
				MinPriceChangePpm:  4_200, // Set a high min price change to test spread calculation.
				ExchangeConfigJson: constants.TestMarketParams[1].ExchangeConfigJson,
			},
			marketPrice: constants.TestMarketPrices[1],
			perpetual:   constants.EthUsd_20PercentInitial_10PercentMaintenance,
			expectedOrderSubticks: []uint64{
				// spreadPpm = max(3_000, 1_500 + 4_200) = 5_700
				// priceSubticks = 3_000_000_000 * 10^(-6 - (-9) + (-9) - (-6)) = 3 * 10^9
				// a_1 = 3 * 10^9 * (1 + 0.0057)^1 = 3_017_100_000
				3_017_100_000,
				// b_1 = 3 * 10^9 * (1 - 0.0057)^1 = 2_982_900_000
				2_982_900_000,
				// a_2 = 3 * 10^9 * (1 + 0.0057)^2 = 3_034_297_470
				// round up to nearest multiple of subticksPerTick=1000.
				3_034_298_000,
				// b_2 = 3 * 10^9 * (1 - 0.0057)^2 = 2_965_897_470
				// round down to nearest multiple of subticksPerTick=1000.
				2_965_897_000,
			},
			expectedOrderQuantums: []uint64{ // TODO (TRA-144): Implement order size
				1000,
				1000,
				1000,
				1000,
			},
		},
		"Error - Clob Pair doesn't exist": {
			vaultId:     constants.Vault_Clob_0,
			clobPair:    constants.ClobPair_Eth,
			marketParam: constants.TestMarketParams[1],
			marketPrice: constants.TestMarketPrices[1],
			perpetual:   constants.EthUsd_NoMarginRequirement,
			expectedErr: vaulttypes.ErrClobPairNotFound,
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
			buildVaultClobOrder := func(
				layer uint8,
				side clobtypes.Order_Side,
				quantums uint64,
				subticks uint64,
			) *clobtypes.Order {
				return &clobtypes.Order{
					OrderId: clobtypes.OrderId{
						SubaccountId: satypes.SubaccountId{
							Owner:  tc.vaultId.ToModuleAccountAddress(),
							Number: 0,
						},
						ClientId:   tApp.App.VaultKeeper.GetVaultClobOrderClientId(ctx, side, layer),
						OrderFlags: clobtypes.OrderIdFlags_LongTerm,
						ClobPairId: tc.vaultId.Number,
					},
					Side:     side,
					Quantums: quantums,
					Subticks: subticks,
					GoodTilOneof: &clobtypes.Order_GoodTilBlockTime{
						GoodTilBlockTime: uint32(ctx.BlockTime().Unix()) + orderExpirationSeconds,
					},
				}
			}
			expectedOrders := make([]*clobtypes.Order, 0)
			for i := uint8(0); i < numLayers; i++ {
				expectedOrders = append(
					expectedOrders,
					// ask.
					buildVaultClobOrder(
						i+1,
						clobtypes.Order_SIDE_SELL,
						tc.expectedOrderQuantums[2*i],
						tc.expectedOrderSubticks[2*i],
					),
					// bid.
					buildVaultClobOrder(
						i+1,
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
