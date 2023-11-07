package gov_test

import (
	"testing"

	"github.com/cometbft/cometbft/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	govtypesv1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
	testapp "github.com/dydxprotocol/v4-chain/protocol/testutil/app"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	perptypes "github.com/dydxprotocol/v4-chain/protocol/x/perpetuals/types"
	"github.com/stretchr/testify/require"
)

// This tests `MsgUpdateParams` in `x/perpetuals`.
func TestUpdatePerpetualsModuleParams(t *testing.T) {
	tests := map[string]struct {
		msg                      *perptypes.MsgUpdateParams
		expectedProposalStatus   govtypesv1.ProposalStatus
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
		"Fail: invalid authority": {
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
				false,
				tc.expectSubmitProposalFail,
				tc.expectedProposalStatus,
			)

			// If governance proposal is supposed to fail submission, verify that perpetuals module
			// params match the ones before proposal submission.
			if tc.expectSubmitProposalFail {
				require.Equal(t, initialParams, tApp.App.PerpetualsKeeper.GetParams(ctx))
			}

			// If proposal is supposed to pass, verify that perpetuals params have been updated.
			if tc.expectedProposalStatus == govtypesv1.ProposalStatus_PROPOSAL_STATUS_PASSED {
				require.Equal(t, tc.msg.Params, tApp.App.PerpetualsKeeper.GetParams(ctx))
			}
		})
	}
}

// This tests `MsgUpdatePerpetualParams` in `x/perpetuals`.
func TestUpdatePerpetualsParams(t *testing.T) {
	tests := map[string]struct {
		msg                      *perptypes.MsgUpdatePerpetualParams
		expectedProposalStatus   govtypesv1.ProposalStatus
		expectSubmitProposalFail bool
	}{
		"Success": {
			msg: &perptypes.MsgUpdatePerpetualParams{
				Authority: authtypes.NewModuleAddress(govtypes.ModuleName).String(),
				PerpetualParams: perptypes.PerpetualParams{
					Id:                5000,
					Ticker:            "BTC-DV4TNT",
					MarketId:          4,
					AtomicResolution:  -7,
					DefaultFundingPpm: 500,
					LiquidityTier:     0,
				},
			},
			expectedProposalStatus: govtypesv1.ProposalStatus_PROPOSAL_STATUS_PASSED,
		},
		"Fail: invalid authority": {
			msg: &perptypes.MsgUpdatePerpetualParams{
				Authority: authtypes.NewModuleAddress(perptypes.ModuleName).String(),
				PerpetualParams: perptypes.PerpetualParams{
					Id:                5000,
					Ticker:            "BTC-DV4TNT",
					MarketId:          4,
					AtomicResolution:  -7,
					DefaultFundingPpm: 500,
					LiquidityTier:     0,
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
				// Initialize perpetuals module with
				// - a perpetual whose params are different from the proposal.
				// - liquidity tiers.
				testapp.UpdateGenesisDocWithAppStateForModule(
					&genesis,
					func(genesisState *perptypes.GenesisState) {
						genesisState.Perpetuals = append(genesisState.Perpetuals, perptypes.Perpetual{
							Params: perptypes.PerpetualParams{
								Id:                tc.msg.PerpetualParams.Id, // same ID as the proposal.
								Ticker:            tc.msg.PerpetualParams.Ticker + "_initial",
								MarketId:          tc.msg.PerpetualParams.MarketId + 1,
								AtomicResolution:  tc.msg.PerpetualParams.AtomicResolution,
								DefaultFundingPpm: tc.msg.PerpetualParams.DefaultFundingPpm + 234,
								LiquidityTier:     tc.msg.PerpetualParams.LiquidityTier + 2,
							},
						})
						genesisState.LiquidityTiers = constants.LiquidityTiers
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
				false,
				tc.expectSubmitProposalFail,
				tc.expectedProposalStatus,
			)

			updatedPerpetual, err := tApp.App.PerpetualsKeeper.GetPerpetual(ctx, tc.msg.PerpetualParams.Id)
			require.NoError(t, err)

			// If governance proposal is supposed to fail submission, verify that the perpetual's
			// params match the ones before proposal submission.
			if tc.expectSubmitProposalFail {
				require.Equal(t, initialPerpetual.Params, updatedPerpetual.Params)
			}

			// If proposal is supposed to pass, verify that the perpetual's params have been updated.
			if tc.expectedProposalStatus == govtypesv1.ProposalStatus_PROPOSAL_STATUS_PASSED {
				require.Equal(t, tc.msg.PerpetualParams, updatedPerpetual.Params)
			}
		})
	}
}

// This tests `MsgSetLiquidityTier` in `x/perpetuals`.
func TestSetLiquidityTier(t *testing.T) {
	tests := map[string]struct {
		msg                      *perptypes.MsgSetLiquidityTier
		updateExistingLt         bool // whether above msg updates an existing liquidity tier.
		expectedProposalStatus   govtypesv1.ProposalStatus
		expectSubmitProposalFail bool
	}{
		"Success: create a new liquidity tier": {
			msg: &perptypes.MsgSetLiquidityTier{
				Authority: authtypes.NewModuleAddress(govtypes.ModuleName).String(),
				LiquidityTier: perptypes.LiquidityTier{
					Id:                     5678,
					Name:                   "Test Tier",
					InitialMarginPpm:       765_432,
					MaintenanceFractionPpm: 345_678,
					BasePositionNotional:   123_456,
					ImpactNotional:         654_321,
				},
			},
			expectedProposalStatus: govtypesv1.ProposalStatus_PROPOSAL_STATUS_PASSED,
		},
		"Success: update an existing liquidity tier": {
			msg: &perptypes.MsgSetLiquidityTier{
				Authority: authtypes.NewModuleAddress(govtypes.ModuleName).String(),
				LiquidityTier: perptypes.LiquidityTier{
					Id:                     5678,
					Name:                   "Test Tier",
					InitialMarginPpm:       765_432,
					MaintenanceFractionPpm: 345_678,
					BasePositionNotional:   123_456,
					ImpactNotional:         654_321,
				},
			},
			updateExistingLt:       true,
			expectedProposalStatus: govtypesv1.ProposalStatus_PROPOSAL_STATUS_PASSED,
		},
		"Fail: invalid authority": {
			msg: &perptypes.MsgSetLiquidityTier{
				Authority: authtypes.NewModuleAddress(perptypes.ModuleName).String(),
				LiquidityTier: perptypes.LiquidityTier{
					Id:                     5678,
					Name:                   "Test Tier",
					InitialMarginPpm:       765_432,
					MaintenanceFractionPpm: 345_678,
					BasePositionNotional:   123_456,
					ImpactNotional:         654_321,
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
				if tc.updateExistingLt {
					// Initialize perpetuals module with a liquidity tier with different params from the proposal.
					testapp.UpdateGenesisDocWithAppStateForModule(
						&genesis,
						func(genesisState *perptypes.GenesisState) {
							genesisState.LiquidityTiers = append(genesisState.LiquidityTiers, perptypes.LiquidityTier{
								Id:                     tc.msg.LiquidityTier.Id,
								Name:                   tc.msg.LiquidityTier.Name + "_initial",
								InitialMarginPpm:       tc.msg.LiquidityTier.InitialMarginPpm + 1,
								MaintenanceFractionPpm: tc.msg.LiquidityTier.MaintenanceFractionPpm + 2,
								BasePositionNotional:   tc.msg.LiquidityTier.BasePositionNotional + 3,
								ImpactNotional:         tc.msg.LiquidityTier.ImpactNotional + 4,
							})
						},
					)
				}
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
				false,
				tc.expectSubmitProposalFail,
				tc.expectedProposalStatus,
			)

			// If governance proposal is supposed to fail submission, verify that liquidity tiers
			// match the ones before proposal submission.
			if tc.expectSubmitProposalFail {
				require.Equal(t, initialLts, tApp.App.PerpetualsKeeper.GetAllLiquidityTiers(ctx))
			}

			// If proposal is supposed to pass, verify that the liquidity tier has been createupdated.
			if tc.expectedProposalStatus == govtypesv1.ProposalStatus_PROPOSAL_STATUS_PASSED {
				updatedLt, err := tApp.App.PerpetualsKeeper.GetLiquidityTier(ctx, tc.msg.LiquidityTier.Id)
				require.NoError(t, err)
				require.Equal(t, tc.msg.LiquidityTier, updatedLt)
			}
		})
	}
}
