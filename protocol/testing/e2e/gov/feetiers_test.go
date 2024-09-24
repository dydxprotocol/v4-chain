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
	feetierstypes "github.com/dydxprotocol/v4-chain/protocol/x/feetiers/types"
	"github.com/stretchr/testify/require"
)

// This tests `MsgUpdatePerpetualFeeParams` in `x/feetiers`.
func TestUpdateFeeTiersModuleParams(t *testing.T) {
	testPerpetualFeeParams := feetierstypes.PerpetualFeeParams{
		Tiers: []*feetierstypes.PerpetualFeeTier{
			{
				Name:        "test_tier_0",
				MakerFeePpm: 11_000,
				TakerFeePpm: 22_000,
			},
			{
				Name:                           "test_tier_1",
				AbsoluteVolumeRequirement:      200_000,
				TotalVolumeShareRequirementPpm: 100_000,
				MakerVolumeShareRequirementPpm: 50_000,
				MakerFeePpm:                    1_000,
				TakerFeePpm:                    2_000,
			},
		},
	}

	tests := map[string]struct {
		msg                       *feetierstypes.MsgUpdatePerpetualFeeParams
		expectCheckTxFails        bool
		expectSubmitProposalFails bool
		expectedProposalStatus    govtypesv1.ProposalStatus
	}{
		"Success": {
			msg: &feetierstypes.MsgUpdatePerpetualFeeParams{
				Authority: authtypes.NewModuleAddress(govtypes.ModuleName).String(),
				Params:    testPerpetualFeeParams,
			},
			expectedProposalStatus: govtypesv1.ProposalStatus_PROPOSAL_STATUS_PASSED,
		},
		"Failure: no tiers": {
			msg: &feetierstypes.MsgUpdatePerpetualFeeParams{
				Authority: authtypes.NewModuleAddress(govtypes.ModuleName).String(),
				Params:    feetierstypes.PerpetualFeeParams{},
			},
			expectCheckTxFails: true,
		},
		"Failure: first tier has non-zero volume requirement": {
			msg: &feetierstypes.MsgUpdatePerpetualFeeParams{
				Authority: authtypes.NewModuleAddress(govtypes.ModuleName).String(),
				Params: feetierstypes.PerpetualFeeParams{
					Tiers: []*feetierstypes.PerpetualFeeTier{
						{
							Name:                      "test_tier_0",
							AbsoluteVolumeRequirement: 1,
							MakerFeePpm:               1_000,
							TakerFeePpm:               2_000,
						},
					},
				},
			},
			expectCheckTxFails: true,
		},
		"Failure: sum of lowest make fee and taker fee is negative": {
			msg: &feetierstypes.MsgUpdatePerpetualFeeParams{
				Authority: authtypes.NewModuleAddress(govtypes.ModuleName).String(),
				Params: feetierstypes.PerpetualFeeParams{
					Tiers: []*feetierstypes.PerpetualFeeTier{
						{
							Name:                      "test_tier_0",
							AbsoluteVolumeRequirement: 1,
							MakerFeePpm:               -1_000, // lowest maker fee.
							TakerFeePpm:               2_000,
						},
						{
							Name:                      "test_tier_1",
							AbsoluteVolumeRequirement: 1,
							MakerFeePpm:               -888,
							TakerFeePpm:               500, // lowest taker fee.
						},
					},
				},
			},
			expectCheckTxFails: true,
		},
		"Failure: invalid authority": {
			msg: &feetierstypes.MsgUpdatePerpetualFeeParams{
				Authority: authtypes.NewModuleAddress(feetierstypes.ModuleName).String(),
				Params:    testPerpetualFeeParams,
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
				// Initialize feetiers module with params that are different from the proposal.
				testapp.UpdateGenesisDocWithAppStateForModule(
					&genesis,
					func(genesisState *feetierstypes.GenesisState) {
						genesisState.Params = constants.PerpetualFeeParams
					},
				)
				return genesis
			}).Build()
			ctx := tApp.InitChain()
			initialPerpetualFeeParams := tApp.App.FeeTiersKeeper.GetPerpetualFeeParams(ctx)

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

			if tc.expectedProposalStatus == govtypesv1.ProposalStatus_PROPOSAL_STATUS_PASSED {
				// If proposal is supposed to pass, verify that perpetual fee params have been updated.
				require.Equal(t, tc.msg.Params, tApp.App.FeeTiersKeeper.GetPerpetualFeeParams(ctx))
			} else {
				// Otherwise, verify that perpetual fee params match the ones before proposal submission.
				require.Equal(t, initialPerpetualFeeParams, tApp.App.FeeTiersKeeper.GetPerpetualFeeParams(ctx))
			}
		})
	}
}
