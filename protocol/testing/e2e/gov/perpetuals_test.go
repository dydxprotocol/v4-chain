package gov_test

import (
	"fmt"
	"testing"

	"github.com/cometbft/cometbft/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	govtypesv1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
	slinkytypes "github.com/dydxprotocol/slinky/pkg/types"
	marketmaptypes "github.com/dydxprotocol/slinky/x/marketmap/types"
	testapp "github.com/dydxprotocol/v4-chain/protocol/testutil/app"
	lttest "github.com/dydxprotocol/v4-chain/protocol/testutil/liquidity_tier"
	pricestest "github.com/dydxprotocol/v4-chain/protocol/testutil/prices"
	clobtypes "github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
	perptypes "github.com/dydxprotocol/v4-chain/protocol/x/perpetuals/types"
	pricestypes "github.com/dydxprotocol/v4-chain/protocol/x/prices/types"
	"github.com/stretchr/testify/require"
)

var (
	TEST_PERPETUAL_PARAMS = perptypes.PerpetualParams{
		Id:                0,
		Ticker:            "BTC-ADV4TNT",
		MarketId:          123,
		AtomicResolution:  -8,
		DefaultFundingPpm: 545,
		LiquidityTier:     1,
		MarketType:        perptypes.PerpetualMarketType_PERPETUAL_MARKET_TYPE_CROSS,
	}
	TEST_LIQUIDITY_TIER = perptypes.LiquidityTier{
		Id:                     1,
		Name:                   "Test Tier",
		InitialMarginPpm:       765_432,
		MaintenanceFractionPpm: 345_678,
		ImpactNotional:         654_321,
	}
)

// This tests `MsgUpdateParams` in `x/perpetuals`.
func TestUpdatePerpetualsModuleParams(t *testing.T) {
	tests := map[string]struct {
		msg                      *perptypes.MsgUpdateParams
		expectedProposalStatus   govtypesv1.ProposalStatus
		expectCheckTxFails       bool
		expectSubmitProposalFail bool
	}{
		"Success": {
			msg: &perptypes.MsgUpdateParams{
				Authority: authtypes.NewModuleAddress(govtypes.ModuleName).String(),
				Params: perptypes.Params{
					FundingRateClampFactorPpm: 123_456,
					PremiumVoteClampFactorPpm: 123_456_789,
					MinNumVotesPerSample:      15,
				},
			},
			expectedProposalStatus: govtypesv1.ProposalStatus_PROPOSAL_STATUS_PASSED,
		},
		"Failure: zero funding rate clamp factor": {
			msg: &perptypes.MsgUpdateParams{
				Authority: authtypes.NewModuleAddress(perptypes.ModuleName).String(),
				Params: perptypes.Params{
					FundingRateClampFactorPpm: 0,
					PremiumVoteClampFactorPpm: 100,
					MinNumVotesPerSample:      15,
				},
			},
			expectCheckTxFails: true,
		},
		"Failure: zero premium vote clamp factor": {
			msg: &perptypes.MsgUpdateParams{
				Authority: authtypes.NewModuleAddress(perptypes.ModuleName).String(),
				Params: perptypes.Params{
					FundingRateClampFactorPpm: 100,
					PremiumVoteClampFactorPpm: 0,
					MinNumVotesPerSample:      15,
				},
			},
			expectCheckTxFails: true,
		},
		"Failure: zero min number of votes per sample": {
			msg: &perptypes.MsgUpdateParams{
				Authority: authtypes.NewModuleAddress(perptypes.ModuleName).String(),
				Params: perptypes.Params{
					FundingRateClampFactorPpm: 100,
					PremiumVoteClampFactorPpm: 100,
					MinNumVotesPerSample:      0,
				},
			},
			expectCheckTxFails: true,
		},
		"Failure: invalid authority": {
			msg: &perptypes.MsgUpdateParams{
				Authority: authtypes.NewModuleAddress(perptypes.ModuleName).String(),
				Params: perptypes.Params{
					FundingRateClampFactorPpm: 100,
					PremiumVoteClampFactorPpm: 100,
					MinNumVotesPerSample:      15,
				},
			},
			expectSubmitProposalFail: true,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tApp := testapp.NewTestAppBuilder(t).WithGenesisDocFn(func() (genesis types.GenesisDoc) {
				genesis = testapp.DefaultGenesis()
				testapp.UpdateGenesisDocWithAppStateForModule(
					&genesis,
					func(genesisState *govtypesv1.GenesisState) {
						genesisState.Params.VotingPeriod = &testapp.TestVotingPeriod
					},
				)
				// Initialize perpetuals module with params that are different from the proposal.
				testapp.UpdateGenesisDocWithAppStateForModule(
					&genesis,
					func(genesisState *perptypes.GenesisState) {
						genesisState.Params = perptypes.Params{
							FundingRateClampFactorPpm: tc.msg.Params.FundingRateClampFactorPpm + 1,
							PremiumVoteClampFactorPpm: tc.msg.Params.PremiumVoteClampFactorPpm + 2,
							MinNumVotesPerSample:      tc.msg.Params.MinNumVotesPerSample + 3,
						}
					},
				)
				return genesis
			}).Build()
			ctx := tApp.InitChain()
			initialParams := tApp.App.PerpetualsKeeper.GetParams(ctx)

			// Submit and tally governance proposal that includes `MsgUpdateParams`.
			ctx = testapp.SubmitAndTallyProposal(
				t,
				ctx,
				tApp,
				[]sdk.Msg{tc.msg},
				testapp.TestSubmitProposalTxHeight,
				tc.expectCheckTxFails,
				tc.expectSubmitProposalFail,
				tc.expectedProposalStatus,
			)

			if tc.expectedProposalStatus == govtypesv1.ProposalStatus_PROPOSAL_STATUS_PASSED {
				// If proposal is supposed to pass, verify that perpetuals params have been updated.
				require.Equal(t, tc.msg.Params, tApp.App.PerpetualsKeeper.GetParams(ctx))
			} else {
				// Otherwise, verify that perpetuals module params match the ones before proposal submission.
				require.Equal(t, initialParams, tApp.App.PerpetualsKeeper.GetParams(ctx))
			}
		})
	}
}

// This tests `MsgUpdatePerpetualParams` in `x/perpetuals`.
func TestUpdatePerpetualsParams(t *testing.T) {
	tests := map[string]struct {
		msg                       *perptypes.MsgUpdatePerpetualParams
		genesisLiquidityTierIds   []uint32
		genesisMarketIds          []uint32
		expectedProposalStatus    govtypesv1.ProposalStatus
		expectCheckTxFails        bool
		expectSubmitProposalFails bool
	}{
		"Success": {
			msg: &perptypes.MsgUpdatePerpetualParams{
				Authority: authtypes.NewModuleAddress(govtypes.ModuleName).String(),
				PerpetualParams: perptypes.PerpetualParams{
					Id:                TEST_PERPETUAL_PARAMS.Id,
					Ticker:            "BTC-DV4TNT",
					MarketId:          4,
					AtomicResolution:  TEST_PERPETUAL_PARAMS.AtomicResolution,
					DefaultFundingPpm: 500,
					LiquidityTier:     123,
					MarketType:        perptypes.PerpetualMarketType_PERPETUAL_MARKET_TYPE_CROSS,
				},
			},
			genesisLiquidityTierIds: []uint32{123},
			genesisMarketIds:        []uint32{4},
			expectedProposalStatus:  govtypesv1.ProposalStatus_PROPOSAL_STATUS_PASSED,
		},
		"Failure: empty ticker": {
			msg: &perptypes.MsgUpdatePerpetualParams{
				Authority: authtypes.NewModuleAddress(govtypes.ModuleName).String(),
				PerpetualParams: perptypes.PerpetualParams{
					Id:                TEST_PERPETUAL_PARAMS.Id,
					Ticker:            "",
					MarketId:          4,
					AtomicResolution:  TEST_PERPETUAL_PARAMS.AtomicResolution,
					DefaultFundingPpm: 500,
					LiquidityTier:     123,
					MarketType:        perptypes.PerpetualMarketType_PERPETUAL_MARKET_TYPE_CROSS,
				},
			},
			genesisLiquidityTierIds: []uint32{123},
			genesisMarketIds:        []uint32{4},
			expectCheckTxFails:      true,
		},
		"Failure: default funding ppm exists maximum": {
			msg: &perptypes.MsgUpdatePerpetualParams{
				Authority: authtypes.NewModuleAddress(govtypes.ModuleName).String(),
				PerpetualParams: perptypes.PerpetualParams{
					Id:                TEST_PERPETUAL_PARAMS.Id,
					Ticker:            "BTC-DV4TNT",
					MarketId:          4,
					AtomicResolution:  TEST_PERPETUAL_PARAMS.AtomicResolution,
					DefaultFundingPpm: 1_000_001,
					LiquidityTier:     123,
					MarketType:        perptypes.PerpetualMarketType_PERPETUAL_MARKET_TYPE_CROSS,
				},
			},
			genesisLiquidityTierIds: []uint32{123},
			genesisMarketIds:        []uint32{4},
			expectCheckTxFails:      true,
		},
		"Failure: invalid authority": {
			msg: &perptypes.MsgUpdatePerpetualParams{
				Authority: authtypes.NewModuleAddress(perptypes.ModuleName).String(),
				PerpetualParams: perptypes.PerpetualParams{
					Id:                TEST_PERPETUAL_PARAMS.Id,
					Ticker:            "BTC-DV4TNT",
					MarketId:          4,
					AtomicResolution:  TEST_PERPETUAL_PARAMS.AtomicResolution,
					DefaultFundingPpm: 500,
					LiquidityTier:     123,
					MarketType:        perptypes.PerpetualMarketType_PERPETUAL_MARKET_TYPE_CROSS,
				},
			},
			genesisLiquidityTierIds:   []uint32{123},
			genesisMarketIds:          []uint32{4},
			expectSubmitProposalFails: true,
		},
		"Failure: liquidity tier does not exist": {
			msg: &perptypes.MsgUpdatePerpetualParams{
				Authority: authtypes.NewModuleAddress(govtypes.ModuleName).String(),
				PerpetualParams: perptypes.PerpetualParams{
					Id:                TEST_PERPETUAL_PARAMS.Id,
					Ticker:            "BTC-DV4TNT",
					MarketId:          4,
					AtomicResolution:  TEST_PERPETUAL_PARAMS.AtomicResolution,
					DefaultFundingPpm: 500,
					LiquidityTier:     124, // liquidity tier 124 does not exist.
					MarketType:        perptypes.PerpetualMarketType_PERPETUAL_MARKET_TYPE_CROSS,
				},
			},
			genesisLiquidityTierIds: []uint32{123},
			genesisMarketIds:        []uint32{4},
			expectedProposalStatus:  govtypesv1.ProposalStatus_PROPOSAL_STATUS_FAILED,
		},
		"Failure: market id does not exist": {
			msg: &perptypes.MsgUpdatePerpetualParams{
				Authority: authtypes.NewModuleAddress(govtypes.ModuleName).String(),
				PerpetualParams: perptypes.PerpetualParams{
					Id:                TEST_PERPETUAL_PARAMS.Id,
					Ticker:            "BTC-DV4TNT",
					MarketId:          5, // market id 5 does not exist.
					AtomicResolution:  TEST_PERPETUAL_PARAMS.AtomicResolution,
					DefaultFundingPpm: 500,
					LiquidityTier:     0,
					MarketType:        perptypes.PerpetualMarketType_PERPETUAL_MARKET_TYPE_CROSS,
				},
			},
			genesisMarketIds:       []uint32{4},
			expectedProposalStatus: govtypesv1.ProposalStatus_PROPOSAL_STATUS_FAILED,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tApp := testapp.NewTestAppBuilder(t).WithGenesisDocFn(func() (genesis types.GenesisDoc) {
				genesis = testapp.DefaultGenesis()
				testapp.UpdateGenesisDocWithAppStateForModule(
					&genesis,
					func(genesisState *govtypesv1.GenesisState) {
						genesisState.Params.VotingPeriod = &testapp.TestVotingPeriod
					},
				)
				// Initialize marketmap module with genesis markets.
				testapp.UpdateGenesisDocWithAppStateForModule(
					&genesis,
					func(genesisState *marketmaptypes.GenesisState) {
						markets := make(map[string]marketmaptypes.Market)
						marketIds := append(tc.genesisMarketIds, TEST_PERPETUAL_PARAMS.MarketId)
						for i := range marketIds {
							ticker := marketmaptypes.Ticker{
								CurrencyPair:     slinkytypes.CurrencyPair{Base: fmt.Sprintf("%d", i), Quote: fmt.Sprintf("%d", i)},
								Decimals:         8,
								MinProviderCount: 3,
								Enabled:          true,
								Metadata_JSON:    "",
							}
							markets[fmt.Sprintf("%d/%d", i, i)] = marketmaptypes.Market{
								Ticker: ticker,
								ProviderConfigs: []marketmaptypes.ProviderConfig{
									{Name: "binance_ws", OffChainTicker: "test"},
									{Name: "bybit_ws", OffChainTicker: "test"},
									{Name: "coinbase_ws", OffChainTicker: "test"},
								},
							}
						}
						marketMap := marketmaptypes.MarketMap{
							Markets: markets,
						}
						genesisState.MarketMap = marketMap
					},
				)
				// Initialize prices module with genesis markets.
				testapp.UpdateGenesisDocWithAppStateForModule(
					&genesis,
					func(genesisState *pricestypes.GenesisState) {
						marketParams := make([]pricestypes.MarketParam, len(tc.genesisMarketIds)+1)
						marketPrices := make([]pricestypes.MarketPrice, len(tc.genesisMarketIds)+1)
						for i, marketId := range append(tc.genesisMarketIds, TEST_PERPETUAL_PARAMS.MarketId) {
							marketParamPrice := pricestest.GenerateMarketParamPrice(
								pricestest.WithId(marketId),
								pricestest.WithPair(fmt.Sprintf("%d-%d", i, i)),
							)
							marketParams[i] = marketParamPrice.Param
							marketPrices[i] = marketParamPrice.Price
						}
						genesisState.MarketParams = marketParams
						genesisState.MarketPrices = marketPrices
					},
				)
				// Initialize perpetuals module with genesis perpetual and liquidity tiers.
				testapp.UpdateGenesisDocWithAppStateForModule(
					&genesis,
					func(genesisState *perptypes.GenesisState) {
						genesisState.Perpetuals = []perptypes.Perpetual{
							{
								Params: TEST_PERPETUAL_PARAMS,
							},
						}
						liquidityTiers := make([]perptypes.LiquidityTier, len(tc.genesisLiquidityTierIds)+1)
						for i, ltId := range append(tc.genesisLiquidityTierIds, TEST_PERPETUAL_PARAMS.LiquidityTier) {
							liquidityTiers[i] = *lttest.GenerateLiquidityTier(lttest.WithId(ltId))
						}
						genesisState.LiquidityTiers = liquidityTiers
					},
				)
				// Initialize clob module with no clob pairs.
				testapp.UpdateGenesisDocWithAppStateForModule(
					&genesis,
					func(genesisState *clobtypes.GenesisState) {
						genesisState.ClobPairs = []clobtypes.ClobPair{}
					},
				)
				return genesis
			}).Build()
			ctx := tApp.InitChain()
			initialPerpetual, err := tApp.App.PerpetualsKeeper.GetPerpetual(ctx, tc.msg.PerpetualParams.Id)
			require.NoError(t, err)

			// Submit and tally governance proposal that includes `MsgUpdateParams`.
			ctx = testapp.SubmitAndTallyProposal(
				t,
				ctx,
				tApp,
				[]sdk.Msg{tc.msg},
				testapp.TestSubmitProposalTxHeight,
				tc.expectCheckTxFails,
				tc.expectSubmitProposalFails,
				tc.expectedProposalStatus,
			)

			updatedPerpetual, err := tApp.App.PerpetualsKeeper.GetPerpetual(ctx, tc.msg.PerpetualParams.Id)
			require.NoError(t, err)

			if tc.expectedProposalStatus == govtypesv1.ProposalStatus_PROPOSAL_STATUS_PASSED {
				// If proposal is supposed to pass, verify that the perpetual's params have been updated.
				// All params except for MarketType should be updated.
				require.Equal(t, tc.msg.PerpetualParams.Ticker, updatedPerpetual.Params.Ticker)
				require.Equal(t, tc.msg.PerpetualParams.MarketId, updatedPerpetual.Params.MarketId)
				require.Equal(t, tc.msg.PerpetualParams.DefaultFundingPpm, updatedPerpetual.Params.DefaultFundingPpm)
				require.Equal(t, tc.msg.PerpetualParams.LiquidityTier, updatedPerpetual.Params.LiquidityTier)
				require.Equal(t, tc.msg.PerpetualParams.AtomicResolution, updatedPerpetual.Params.AtomicResolution)
				require.Equal(t, initialPerpetual.Params.MarketType, updatedPerpetual.Params.MarketType)
			} else {
				// If proposal is not supposed to succeed, verify that the perpetual's
				// params match the ones before proposal submission.
				require.Equal(t, initialPerpetual.Params, updatedPerpetual.Params)
			}
		})
	}
}

// This tests `MsgSetLiquidityTier` in `x/perpetuals`.
func TestSetLiquidityTier(t *testing.T) {
	tests := map[string]struct {
		msg                       *perptypes.MsgSetLiquidityTier
		genesisLiquidityTierIds   []uint32
		expectedProposalStatus    govtypesv1.ProposalStatus
		expectCheckTxFails        bool
		expectSubmitProposalFails bool
	}{
		"Success: create a new liquidity tier": {
			msg: &perptypes.MsgSetLiquidityTier{
				Authority: authtypes.NewModuleAddress(govtypes.ModuleName).String(),
				LiquidityTier: perptypes.LiquidityTier{
					Id:                     5678,
					Name:                   TEST_LIQUIDITY_TIER.Name,
					InitialMarginPpm:       TEST_LIQUIDITY_TIER.InitialMarginPpm,
					MaintenanceFractionPpm: TEST_LIQUIDITY_TIER.MaintenanceFractionPpm,
					ImpactNotional:         TEST_LIQUIDITY_TIER.ImpactNotional,
				},
			},
			expectedProposalStatus: govtypesv1.ProposalStatus_PROPOSAL_STATUS_PASSED,
		},
		"Success: update an existing liquidity tier": {
			msg: &perptypes.MsgSetLiquidityTier{
				Authority: authtypes.NewModuleAddress(govtypes.ModuleName).String(),
				LiquidityTier: perptypes.LiquidityTier{
					Id:                     5678,
					Name:                   TEST_LIQUIDITY_TIER.Name,
					InitialMarginPpm:       TEST_LIQUIDITY_TIER.InitialMarginPpm,
					MaintenanceFractionPpm: TEST_LIQUIDITY_TIER.MaintenanceFractionPpm,
					ImpactNotional:         TEST_LIQUIDITY_TIER.ImpactNotional,
				},
			},
			genesisLiquidityTierIds: []uint32{5678},
			expectedProposalStatus:  govtypesv1.ProposalStatus_PROPOSAL_STATUS_PASSED,
		},
		"Failure: initial margin ppm exceeds maximum": {
			msg: &perptypes.MsgSetLiquidityTier{
				Authority: authtypes.NewModuleAddress(perptypes.ModuleName).String(),
				LiquidityTier: perptypes.LiquidityTier{
					Id:                     5678,
					Name:                   TEST_LIQUIDITY_TIER.Name,
					InitialMarginPpm:       1_000_001,
					MaintenanceFractionPpm: TEST_LIQUIDITY_TIER.MaintenanceFractionPpm,
					ImpactNotional:         TEST_LIQUIDITY_TIER.ImpactNotional,
				},
			},
			expectCheckTxFails: true,
		},
		"Failure: maintenance fraction ppm exceeds maximum": {
			msg: &perptypes.MsgSetLiquidityTier{
				Authority: authtypes.NewModuleAddress(perptypes.ModuleName).String(),
				LiquidityTier: perptypes.LiquidityTier{
					Id:                     5678,
					Name:                   TEST_LIQUIDITY_TIER.Name,
					InitialMarginPpm:       TEST_LIQUIDITY_TIER.InitialMarginPpm,
					MaintenanceFractionPpm: 1_000_001,
					ImpactNotional:         TEST_LIQUIDITY_TIER.ImpactNotional,
				},
			},
			expectCheckTxFails: true,
		},
		"Failure: impact notional is 0": {
			msg: &perptypes.MsgSetLiquidityTier{
				Authority: authtypes.NewModuleAddress(perptypes.ModuleName).String(),
				LiquidityTier: perptypes.LiquidityTier{
					Id:                     5678,
					Name:                   TEST_LIQUIDITY_TIER.Name,
					InitialMarginPpm:       TEST_LIQUIDITY_TIER.InitialMarginPpm,
					MaintenanceFractionPpm: TEST_LIQUIDITY_TIER.MaintenanceFractionPpm,
					ImpactNotional:         0,
				},
			},
			expectCheckTxFails: true,
		},
		"Failure: invalid authority": {
			msg: &perptypes.MsgSetLiquidityTier{
				Authority: authtypes.NewModuleAddress(perptypes.ModuleName).String(),
				LiquidityTier: perptypes.LiquidityTier{
					Id:                     5678,
					Name:                   "Test Tier",
					InitialMarginPpm:       765_432,
					MaintenanceFractionPpm: 345_678,
					ImpactNotional:         654_321,
				},
			},
			expectSubmitProposalFails: true,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tApp := testapp.NewTestAppBuilder(t).WithGenesisDocFn(func() (genesis types.GenesisDoc) {
				genesis = testapp.DefaultGenesis()
				testapp.UpdateGenesisDocWithAppStateForModule(
					&genesis,
					func(genesisState *govtypesv1.GenesisState) {
						genesisState.Params.VotingPeriod = &testapp.TestVotingPeriod
					},
				)
				// Initialize perpetuals module with genesis liquidity tiers and no perpetuals.
				testapp.UpdateGenesisDocWithAppStateForModule(
					&genesis,
					func(genesisState *perptypes.GenesisState) {
						liquidityTiers := make([]perptypes.LiquidityTier, len(tc.genesisLiquidityTierIds))
						for i, ltId := range tc.genesisLiquidityTierIds {
							liquidityTiers[i] = *lttest.GenerateLiquidityTier(lttest.WithId(ltId))
						}
						genesisState.LiquidityTiers = liquidityTiers
						genesisState.Perpetuals = []perptypes.Perpetual{}
					},
				)
				// Initialize clob module with no clob pairs.
				testapp.UpdateGenesisDocWithAppStateForModule(
					&genesis,
					func(genesisState *clobtypes.GenesisState) {
						genesisState.ClobPairs = []clobtypes.ClobPair{}
					},
				)
				return genesis
			}).Build()
			ctx := tApp.InitChain()
			initialLts := tApp.App.PerpetualsKeeper.GetAllLiquidityTiers(ctx)

			// Submit and tally governance proposal that includes `MsgSetLiquidityTier`.
			ctx = testapp.SubmitAndTallyProposal(
				t,
				ctx,
				tApp,
				[]sdk.Msg{tc.msg},
				testapp.TestSubmitProposalTxHeight,
				tc.expectCheckTxFails,
				tc.expectSubmitProposalFails,
				tc.expectedProposalStatus,
			)

			if tc.expectedProposalStatus == govtypesv1.ProposalStatus_PROPOSAL_STATUS_PASSED {
				// If proposal is supposed to pass, verify that the liquidity tier has been createupdated.
				updatedLt, err := tApp.App.PerpetualsKeeper.GetLiquidityTier(ctx, tc.msg.LiquidityTier.Id)
				require.NoError(t, err)
				require.Equal(t, tc.msg.LiquidityTier, updatedLt)
			} else {
				// Otherwise, verify that liquidity tiers match the ones before proposal submission.
				require.Equal(t, initialLts, tApp.App.PerpetualsKeeper.GetAllLiquidityTiers(ctx))
			}
		})
	}
}
