package keeper_test

import (
	"testing"

	vetesting "github.com/StreamFinance-Protocol/stream-chain/protocol/testutil/ve"

	"github.com/cometbft/cometbft/types"

	"github.com/StreamFinance-Protocol/stream-chain/protocol/app/ve"
	testapp "github.com/StreamFinance-Protocol/stream-chain/protocol/testutil/app"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/testutil/constants"
	assettypes "github.com/StreamFinance-Protocol/stream-chain/protocol/x/assets/types"
	clobtypes "github.com/StreamFinance-Protocol/stream-chain/protocol/x/clob/types"
	feetiertypes "github.com/StreamFinance-Protocol/stream-chain/protocol/x/feetiers/types"
	perptypes "github.com/StreamFinance-Protocol/stream-chain/protocol/x/perpetuals/types"
	prices "github.com/StreamFinance-Protocol/stream-chain/protocol/x/prices/types"
	satypes "github.com/StreamFinance-Protocol/stream-chain/protocol/x/subaccounts/types"
	"github.com/stretchr/testify/require"
)

func TestChangePriceVE_CauseNegativeTNC(t *testing.T) {
	tests := map[string]struct {
		// State.
		subaccounts                   []satypes.Subaccount
		marketIdToOraclePriceOverride map[uint32]uint64

		// Configuration.
		liquidityTiers []perptypes.LiquidityTier
		perpetuals     []perptypes.Perpetual
		clobPairs      []clobtypes.ClobPair

		// action
		priceUpdate map[uint32]ve.VEPricePair

		// Expectations.
		expectedNegativeTncSubaccountSeenAtBlock uint32
	}{
		`No price change`: {
			subaccounts: []satypes.Subaccount{
				constants.Carl_Num0_1BTC_Short_50000USD,
			},
			marketIdToOraclePriceOverride: map[uint32]uint64{
				constants.BtcUsd.MarketId: 4_000_000_000, // $40,000 / BTC
			},

			liquidityTiers: constants.LiquidityTiers,
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_20PercentInitial_10PercentMaintenance_OpenInterest1,
			},
			clobPairs:                                []clobtypes.ClobPair{constants.ClobPair_Btc},
			expectedNegativeTncSubaccountSeenAtBlock: 0,

			priceUpdate: map[uint32]ve.VEPricePair{},
		},
		`Price change causes Negative TNC`: {
			subaccounts: []satypes.Subaccount{
				constants.Carl_Num0_1BTC_Short_50000USD,
			},

			marketIdToOraclePriceOverride: map[uint32]uint64{
				constants.BtcUsd.MarketId: 4_000_000_000, // $40,000 / BTC
			},

			liquidityTiers: constants.LiquidityTiers,
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_20PercentInitial_10PercentMaintenance_OpenInterest1,
			},
			clobPairs: []clobtypes.ClobPair{constants.ClobPair_Btc},

			expectedNegativeTncSubaccountSeenAtBlock: 4,

			priceUpdate: map[uint32]ve.VEPricePair{
				0: {
					SpotPrice: 6_000_000_000,
					PnlPrice:  6_000_000_000,
				},
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tApp := testapp.NewTestAppBuilder(t).WithGenesisDocFn(func() (genesis types.GenesisDoc) {
				genesis = testapp.DefaultGenesis()
				testapp.UpdateGenesisDocWithAppStateForModule(
					&genesis,
					func(genesisState *assettypes.GenesisState) {
						genesisState.Assets = []assettypes.Asset{
							*constants.Usdc,
						}
					},
				)
				testapp.UpdateGenesisDocWithAppStateForModule(
					&genesis,
					func(genesisState *prices.GenesisState) {
						// Set oracle prices in the genesis.
						pricesGenesis := constants.TestPricesGenesisState

						// Make a copy of the MarketPrices slice to avoid modifying by reference.
						marketPricesCopy := make([]prices.MarketPrice, len(pricesGenesis.MarketPrices))
						copy(marketPricesCopy, pricesGenesis.MarketPrices)

						for marketId, oraclePrice := range tc.marketIdToOraclePriceOverride {

							exponent, exists := constants.TestMarketIdsToExponents[marketId]
							require.True(t, exists)

							marketPricesCopy[marketId] = prices.MarketPrice{
								Id:        marketId,
								SpotPrice: oraclePrice,
								PnlPrice:  oraclePrice,
								Exponent:  exponent,
							}
						}

						pricesGenesis.MarketPrices = marketPricesCopy
						*genesisState = pricesGenesis
					},
				)
				testapp.UpdateGenesisDocWithAppStateForModule(
					&genesis,
					func(genesisState *perptypes.GenesisState) {
						genesisState.Params = constants.PerpetualsGenesisParams
						genesisState.LiquidityTiers = tc.liquidityTiers
						genesisState.Perpetuals = tc.perpetuals
					},
				)
				testapp.UpdateGenesisDocWithAppStateForModule(
					&genesis,
					func(genesisState *satypes.GenesisState) {
						genesisState.Subaccounts = tc.subaccounts
					},
				)
				testapp.UpdateGenesisDocWithAppStateForModule(
					&genesis,
					func(genesisState *clobtypes.GenesisState) {
						genesisState.ClobPairs = tc.clobPairs
						genesisState.EquityTierLimitConfig = clobtypes.EquityTierLimitConfiguration{}
					},
				)
				testapp.UpdateGenesisDocWithAppStateForModule(
					&genesis,
					func(genesisState *feetiertypes.GenesisState) {
						genesisState.Params = constants.PerpetualFeeParamsNoFee
					},
				)
				return genesis
			}).Build()

			ctx := tApp.AdvanceToBlock(2, testapp.AdvanceToBlockOptions{})

			// Add the price update.
			deliverTxsOverride := make([][]byte, 0)
			_, extCommitBz, err := vetesting.GetInjectedExtendedCommitInfoForTestApp(
				&tApp.App.ConsumerKeeper,
				ctx,
				tc.priceUpdate,
				tApp.GetHeader().Height,
			)
			require.NoError(t, err)

			deliverTxsOverride = append([][]byte{extCommitBz}, deliverTxsOverride...)

			// Advance to the next block, updating the price.
			_ = tApp.AdvanceToBlock(3, testapp.AdvanceToBlockOptions{
				DeliverTxsOverride: deliverTxsOverride,
			})

			ctx = tApp.AdvanceToBlock(4, testapp.AdvanceToBlockOptions{})

			negativeTncSubaccountSeenAtBlock, _, err := tApp.App.SubaccountsKeeper.GetNegativeTncSubaccountSeenAtBlock(
				ctx,
				constants.BtcUsd_NoMarginRequirement.Params.Id,
			)
			require.NoError(t, err)
			require.Equal(t, tc.expectedNegativeTncSubaccountSeenAtBlock, negativeTncSubaccountSeenAtBlock)
		})
	}
}
