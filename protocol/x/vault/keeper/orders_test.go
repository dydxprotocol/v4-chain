package keeper_test

import (
	"math"
	"testing"

	"github.com/cometbft/cometbft/types"
	"github.com/dydxprotocol/v4-chain/protocol/lib"
	testapp "github.com/dydxprotocol/v4-chain/protocol/testutil/app"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	clobtypes "github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
	perptypes "github.com/dydxprotocol/v4-chain/protocol/x/perpetuals/types"
	pricestypes "github.com/dydxprotocol/v4-chain/protocol/x/prices/types"
	satypes "github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"
	vaulttypes "github.com/dydxprotocol/v4-chain/protocol/x/vault/types"
	"github.com/stretchr/testify/require"
)

func TestGetVaultClobOrders(t *testing.T) {
	// Set up test clob pairs, markets, perpetuals, and vault params.
	testClobPairs := []clobtypes.ClobPair{
		constants.ClobPair_Btc,
		constants.ClobPair_Eth,
	}
	testMarketParams := []pricestypes.MarketParam{
		constants.TestMarketParams[0],
		{
			Id:                 constants.TestMarketParams[1].Id,
			Pair:               constants.TestMarketParams[1].Pair,
			Exponent:           constants.TestMarketParams[1].Exponent,
			MinExchanges:       constants.TestMarketParams[1].MinExchanges,
			MinPriceChangePpm:  4_200, // Set a high min price change to test spread calculation.
			ExchangeConfigJson: constants.TestMarketParams[1].ExchangeConfigJson,
		},
	}
	testMarketPrices := []pricestypes.MarketPrice{
		constants.TestMarketPrices[0],
		constants.TestMarketPrices[1],
	}
	testPerps := []perptypes.Perpetual{
		constants.BtcUsd_0DefaultFunding_0AtomicResolution,
		constants.EthUsd_0DefaultFunding_9AtomicResolution,
	}
	// TODO (TRA-118): store vault strategy constants in x/vault state.
	minBaseSpreadPpm := uint32(3_000)                   // 30bps
	baseSpreadMinPriceChangePremiumPpm := uint32(1_500) // 15bps
	orderExpirationSeconds := uint32(5)                 // 5 seconds

	// Initialize tApp and ctx with above test parameters.
	tApp := testapp.NewTestAppBuilder(t).WithGenesisDocFn(func() (genesis types.GenesisDoc) {
		genesis = testapp.DefaultGenesis()
		// Initialize prices module with test markets.
		testapp.UpdateGenesisDocWithAppStateForModule(
			&genesis,
			func(genesisState *pricestypes.GenesisState) {
				genesisState.MarketParams = testMarketParams
				genesisState.MarketPrices = testMarketPrices
			},
		)
		// Initialize perpetuals module with test perpetuals.
		testapp.UpdateGenesisDocWithAppStateForModule(
			&genesis,
			func(genesisState *perptypes.GenesisState) {
				genesisState.LiquidityTiers = constants.LiquidityTiers
				genesisState.Perpetuals = testPerps
			},
		)
		// Initialize clob module with test clob pairs.
		testapp.UpdateGenesisDocWithAppStateForModule(
			&genesis,
			func(genesisState *clobtypes.GenesisState) {
				genesisState.ClobPairs = testClobPairs
			},
		)
		return genesis
	}).Build()
	ctx := tApp.InitChain()

	// Calculate subticks of BTC and ETH.
	btcSubticks := clobtypes.PriceToSubticks(
		testMarketPrices[0],
		testClobPairs[0],
		testPerps[0].Params.AtomicResolution,
		lib.QuoteCurrencyAtomicResolution,
	)
	ethSubticks := clobtypes.PriceToSubticks(
		testMarketPrices[1],
		testClobPairs[1],
		testPerps[1].Params.AtomicResolution,
		lib.QuoteCurrencyAtomicResolution,
	)
	// Calculate spreads of BTC and ETH.
	// - spread = max(minBaseSpreadPpm, baseSpreadMinPriceChangePremiumPpm + minPriceChangePpm)
	// - btcSpreadPpm = max(3_000, 1_500+50) = 3_000
	btcSpreadPpm := lib.Max(minBaseSpreadPpm, baseSpreadMinPriceChangePremiumPpm+testMarketParams[0].MinPriceChangePpm)
	// - ethSpreadPpm = max(3_000, 1_500+4_200) = 5_700
	ethSpreadPpm := lib.Max(minBaseSpreadPpm, baseSpreadMinPriceChangePremiumPpm+testMarketParams[1].MinPriceChangePpm)
	// Calculate order good-till-block-time.
	orderGtbt := &clobtypes.Order_GoodTilBlockTime{
		GoodTilBlockTime: uint32(ctx.BlockTime().Unix()) + orderExpirationSeconds,
	}

	tests := map[string]struct {
		/* --- Setup --- */
		// Vault ID.
		vaultId vaulttypes.VaultId

		/* --- Expectations --- */
		// Expected orders.
		expectedOrders []clobtypes.Order
	}{
		"Get orders from Vault for Clob Pair 0": {
			vaultId: constants.Vault_Clob_0,
			expectedOrders: []clobtypes.Order{
				// ask at layer 1.
				{
					OrderId: clobtypes.OrderId{
						SubaccountId: satypes.SubaccountId{
							Owner:  constants.Vault_Clob_0.ToModuleAccountAddress(),
							Number: 0,
						},
						ClientId: tApp.App.VaultKeeper.GetVaultClobOrderClientId(
							ctx,
							clobtypes.Order_SIDE_SELL,
							uint8(1),
						),
						OrderFlags: clobtypes.OrderIdFlags_LongTerm,
						ClobPairId: constants.Vault_Clob_0.Number,
					},
					Side:     clobtypes.Order_SIDE_SELL,
					Quantums: testClobPairs[0].StepBaseQuantums, // TODO (TRA-144): Implement order size
					Subticks: lib.BigRatRoundToNearestMultiple(
						// subticks of ask_1 = subticks * (1 + spread)
						lib.BigRatMulPpm(btcSubticks, lib.OneMillion+btcSpreadPpm),
						testClobPairs[0].SubticksPerTick,
						true, // round up for asks
					),
					GoodTilOneof: orderGtbt,
				},
				// bid at layer 1.
				{
					OrderId: clobtypes.OrderId{
						SubaccountId: satypes.SubaccountId{
							Owner:  constants.Vault_Clob_0.ToModuleAccountAddress(),
							Number: 0,
						},
						ClientId: tApp.App.VaultKeeper.GetVaultClobOrderClientId(
							ctx,
							clobtypes.Order_SIDE_BUY,
							uint8(1),
						),
						OrderFlags: clobtypes.OrderIdFlags_LongTerm,
						ClobPairId: constants.Vault_Clob_0.Number,
					},
					Side:     clobtypes.Order_SIDE_BUY,
					Quantums: testClobPairs[0].StepBaseQuantums, // TODO (TRA-144): Implement order size
					Subticks: lib.BigRatRoundToNearestMultiple(
						// subticks of bid_1 = subticks * (1 - spread)
						lib.BigRatMulPpm(btcSubticks, lib.OneMillion-btcSpreadPpm),
						testClobPairs[0].SubticksPerTick,
						false, // round down for bids
					),
					GoodTilOneof: orderGtbt,
				},
				// ask at layer 2.
				{
					OrderId: clobtypes.OrderId{
						SubaccountId: satypes.SubaccountId{
							Owner:  constants.Vault_Clob_0.ToModuleAccountAddress(),
							Number: 0,
						},
						ClientId: tApp.App.VaultKeeper.GetVaultClobOrderClientId(
							ctx,
							clobtypes.Order_SIDE_SELL,
							uint8(2),
						),
						OrderFlags: clobtypes.OrderIdFlags_LongTerm,
						ClobPairId: constants.Vault_Clob_0.Number,
					},
					Side:     clobtypes.Order_SIDE_SELL,
					Quantums: testClobPairs[0].StepBaseQuantums, // TODO (TRA-144): Implement order size
					Subticks: lib.BigRatRoundToNearestMultiple(
						// subticks of ask_2 = subticks * (1 + spread)^2
						lib.BigRatMulPpm(
							lib.BigRatMulPpm(btcSubticks, lib.OneMillion+btcSpreadPpm),
							lib.OneMillion+btcSpreadPpm,
						),
						testClobPairs[0].SubticksPerTick,
						true, // round up for asks
					),
					GoodTilOneof: orderGtbt,
				},
				// bid at layer 2.
				{
					OrderId: clobtypes.OrderId{
						SubaccountId: satypes.SubaccountId{
							Owner:  constants.Vault_Clob_0.ToModuleAccountAddress(),
							Number: 0,
						},
						ClientId: tApp.App.VaultKeeper.GetVaultClobOrderClientId(
							ctx,
							clobtypes.Order_SIDE_BUY,
							uint8(2),
						),
						OrderFlags: clobtypes.OrderIdFlags_LongTerm,
						ClobPairId: constants.Vault_Clob_0.Number,
					},
					Side:     clobtypes.Order_SIDE_BUY,
					Quantums: testClobPairs[0].StepBaseQuantums, // TODO (TRA-144): Implement order size
					Subticks: lib.BigRatRoundToNearestMultiple(
						// subticks of bid_2 = subticks * (1 - spread)^2
						lib.BigRatMulPpm(
							lib.BigRatMulPpm(btcSubticks, lib.OneMillion-btcSpreadPpm),
							lib.OneMillion-btcSpreadPpm,
						),
						testClobPairs[0].SubticksPerTick,
						false, // round down for bids
					),
					GoodTilOneof: orderGtbt,
				},
			},
		},
		"Get orders from Vault for Clob Pair 1": {
			vaultId: constants.Vault_Clob_1,
			expectedOrders: []clobtypes.Order{
				// ask at layer 1.
				{
					OrderId: clobtypes.OrderId{
						SubaccountId: satypes.SubaccountId{
							Owner:  constants.Vault_Clob_1.ToModuleAccountAddress(),
							Number: 0,
						},
						ClientId: tApp.App.VaultKeeper.GetVaultClobOrderClientId(
							ctx,
							clobtypes.Order_SIDE_SELL,
							uint8(1),
						),
						OrderFlags: clobtypes.OrderIdFlags_LongTerm,
						ClobPairId: constants.Vault_Clob_1.Number,
					},
					Side:     clobtypes.Order_SIDE_SELL,
					Quantums: testClobPairs[1].StepBaseQuantums, // TODO (TRA-144): Implement order size
					Subticks: lib.BigRatRoundToNearestMultiple(
						// subticks of ask_1 = subticks * (1 + spread)
						lib.BigRatMulPpm(ethSubticks, lib.OneMillion+ethSpreadPpm),
						testClobPairs[1].SubticksPerTick,
						true, // round up for asks
					),
					GoodTilOneof: orderGtbt,
				},
				// bid at layer 1.
				{
					OrderId: clobtypes.OrderId{
						SubaccountId: satypes.SubaccountId{
							Owner:  constants.Vault_Clob_1.ToModuleAccountAddress(),
							Number: 0,
						},
						ClientId: tApp.App.VaultKeeper.GetVaultClobOrderClientId(
							ctx,
							clobtypes.Order_SIDE_BUY,
							uint8(1),
						),
						OrderFlags: clobtypes.OrderIdFlags_LongTerm,
						ClobPairId: constants.Vault_Clob_1.Number,
					},
					Side:     clobtypes.Order_SIDE_BUY,
					Quantums: testClobPairs[1].StepBaseQuantums, // TODO (TRA-144): Implement order size
					Subticks: lib.BigRatRoundToNearestMultiple(
						// subticks of bid_1 = subticks * (1 - spread)
						lib.BigRatMulPpm(ethSubticks, lib.OneMillion-ethSpreadPpm),
						testClobPairs[1].SubticksPerTick,
						false, // round down for bids
					),
					GoodTilOneof: orderGtbt,
				},
				// ask at layer 2.
				{
					OrderId: clobtypes.OrderId{
						SubaccountId: satypes.SubaccountId{
							Owner:  constants.Vault_Clob_1.ToModuleAccountAddress(),
							Number: 0,
						},
						ClientId: tApp.App.VaultKeeper.GetVaultClobOrderClientId(
							ctx,
							clobtypes.Order_SIDE_SELL,
							uint8(2),
						),
						OrderFlags: clobtypes.OrderIdFlags_LongTerm,
						ClobPairId: constants.Vault_Clob_1.Number,
					},
					Side:     clobtypes.Order_SIDE_SELL,
					Quantums: testClobPairs[1].StepBaseQuantums, // TODO (TRA-144): Implement order size
					Subticks: lib.BigRatRoundToNearestMultiple(
						// subticks of ask_2 = subticks * (1 + spread)^2
						lib.BigRatMulPpm(
							lib.BigRatMulPpm(ethSubticks, lib.OneMillion+ethSpreadPpm),
							lib.OneMillion+ethSpreadPpm,
						),
						testClobPairs[1].SubticksPerTick,
						true, // round up for asks
					),
					GoodTilOneof: orderGtbt,
				},
				// bid at layer 2.
				{
					OrderId: clobtypes.OrderId{
						SubaccountId: satypes.SubaccountId{
							Owner:  constants.Vault_Clob_1.ToModuleAccountAddress(),
							Number: 0,
						},
						ClientId: tApp.App.VaultKeeper.GetVaultClobOrderClientId(
							ctx,
							clobtypes.Order_SIDE_BUY,
							uint8(2),
						),
						OrderFlags: clobtypes.OrderIdFlags_LongTerm,
						ClobPairId: constants.Vault_Clob_1.Number,
					},
					Side:     clobtypes.Order_SIDE_BUY,
					Quantums: testClobPairs[1].StepBaseQuantums, // TODO (TRA-144): Implement order size
					Subticks: lib.BigRatRoundToNearestMultiple(
						// subticks of bid_2 = subticks * (1 - spread)^2
						lib.BigRatMulPpm(
							lib.BigRatMulPpm(ethSubticks, lib.OneMillion-ethSpreadPpm),
							lib.OneMillion-ethSpreadPpm,
						),
						testClobPairs[1].SubticksPerTick,
						false, // round down for bids
					),
					GoodTilOneof: orderGtbt,
				},
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			vaultOrders := tApp.App.VaultKeeper.GetVaultClobOrders(
				ctx,
				tc.vaultId,
			)
			require.Equal(t, tc.expectedOrders, vaultOrders)
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
			side:             clobtypes.Order_SIDE_SELL, // 1<<31
			blockHeight:      -678987,                   // 1<<30
			layer:            202,                       // 202<<22
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
