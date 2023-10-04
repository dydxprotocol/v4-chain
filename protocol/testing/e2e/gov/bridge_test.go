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
	bridgetypes "github.com/dydxprotocol/v4-chain/protocol/x/bridge/types"
	"github.com/stretchr/testify/require"
)

func TestUpdateEventParams(t *testing.T) {
	genesisEventParams := bridgetypes.DefaultGenesis().EventParams

	tests := map[string]struct {
		msg                      *bridgetypes.MsgUpdateEventParams
		expectSubmitProposalFail bool
		expectedProposalStatus   govtypesv1.ProposalStatus
	}{
		"Success": {
			msg: &bridgetypes.MsgUpdateEventParams{
				Authority: authtypes.NewModuleAddress(govtypes.ModuleName).String(),
				Params: bridgetypes.EventParams{
					Denom:      genesisEventParams.Denom + "updated",
					EthChainId: genesisEventParams.EthChainId + 1,
					EthAddress: genesisEventParams.EthAddress,
				},
			},
			expectedProposalStatus: govtypesv1.ProposalStatus_PROPOSAL_STATUS_PASSED,
		},
		"Fail: invalid authority": {
			msg: &bridgetypes.MsgUpdateEventParams{
				Authority: constants.BobAccAddress.String(),
				Params: bridgetypes.EventParams{
					Denom:      genesisEventParams.Denom + "updated",
					EthChainId: genesisEventParams.EthChainId + 1,
					EthAddress: genesisEventParams.EthAddress,
				},
			},
			expectSubmitProposalFail: true,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tApp := testapp.NewTestAppBuilder().WithGenesisDocFn(func() (genesis types.GenesisDoc) {
				genesis = testapp.DefaultGenesis()
				testapp.UpdateGenesisDocWithAppStateForModule(
					&genesis,
					func(genesisState *govtypesv1.GenesisState) {
						genesisState.Params.VotingPeriod = &testapp.TestVotingPeriod
					},
				)
				return genesis
			}).WithTesting(t).Build()
			ctx := tApp.InitChain()

			// Submit and tally governance proposal that includes `MsgUpdateEventParams`.
			ctx = testapp.SubmitAndTallyProposal(
				t,
				ctx,
				&tApp,
				[]sdk.Msg{tc.msg},
				tc.expectSubmitProposalFail,
				tc.expectedProposalStatus,
			)

			// If proposal is supposed to pass, verify that event params are updated.
			if tc.expectedProposalStatus == govtypesv1.ProposalStatus_PROPOSAL_STATUS_PASSED {
				updatedEventParams := tApp.App.BridgeKeeper.GetEventParams(ctx)
				require.Equal(t, tc.msg.Params, updatedEventParams)
			}
		})
	}
}

func TestUpdateProposeParams(t *testing.T) {
	genesisProposeParams := bridgetypes.DefaultGenesis().ProposeParams

	tests := map[string]struct {
		msg                      *bridgetypes.MsgUpdateProposeParams
		expectSubmitProposalFail bool
		expectedProposalStatus   govtypesv1.ProposalStatus
	}{
		"Success": {
			msg: &bridgetypes.MsgUpdateProposeParams{
				Authority: authtypes.NewModuleAddress(govtypes.ModuleName).String(),
				Params: bridgetypes.ProposeParams{
					MaxBridgesPerBlock:           genesisProposeParams.MaxBridgesPerBlock + 1,
					ProposeDelayDuration:         genesisProposeParams.ProposeDelayDuration + 1,
					SkipRatePpm:                  genesisProposeParams.SkipRatePpm + 1,
					SkipIfBlockDelayedByDuration: genesisProposeParams.SkipIfBlockDelayedByDuration + 1,
				},
			},
			expectedProposalStatus: govtypesv1.ProposalStatus_PROPOSAL_STATUS_PASSED,
		},
		"Fail: invalid authority": {
			msg: &bridgetypes.MsgUpdateProposeParams{
				Authority: constants.AliceAccAddress.String(),
				Params: bridgetypes.ProposeParams{
					MaxBridgesPerBlock:           genesisProposeParams.MaxBridgesPerBlock + 1,
					ProposeDelayDuration:         genesisProposeParams.ProposeDelayDuration + 1,
					SkipRatePpm:                  genesisProposeParams.SkipRatePpm + 1,
					SkipIfBlockDelayedByDuration: genesisProposeParams.SkipIfBlockDelayedByDuration + 1,
				},
			},
			expectSubmitProposalFail: true,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tApp := testapp.NewTestAppBuilder().WithGenesisDocFn(func() (genesis types.GenesisDoc) {
				genesis = testapp.DefaultGenesis()
				testapp.UpdateGenesisDocWithAppStateForModule(
					&genesis,
					func(genesisState *govtypesv1.GenesisState) {
						genesisState.Params.VotingPeriod = &testapp.TestVotingPeriod
					},
				)
				return genesis
			}).WithTesting(t).Build()
			ctx := tApp.InitChain()

			// Submit and tally governance proposal that includes `MsgUpdateProposeParams`.
			ctx = testapp.SubmitAndTallyProposal(
				t,
				ctx,
				&tApp,
				[]sdk.Msg{tc.msg},
				tc.expectSubmitProposalFail,
				tc.expectedProposalStatus,
			)

			// If proposal is supposed to pass, verify that propose params are updated.
			if tc.expectedProposalStatus == govtypesv1.ProposalStatus_PROPOSAL_STATUS_PASSED {
				updatedProposeParams := tApp.App.BridgeKeeper.GetProposeParams(ctx)
				require.Equal(t, tc.msg.Params, updatedProposeParams)
			}
		})
	}
}

func TestUpdateSafetyParams(t *testing.T) {
	genesisSafetyParams := bridgetypes.DefaultGenesis().SafetyParams

	tests := map[string]struct {
		msg                      *bridgetypes.MsgUpdateSafetyParams
		expectSubmitProposalFail bool
		expectedProposalStatus   govtypesv1.ProposalStatus
	}{
		"Success": {
			msg: &bridgetypes.MsgUpdateSafetyParams{
				Authority: authtypes.NewModuleAddress(govtypes.ModuleName).String(),
				Params: bridgetypes.SafetyParams{
					IsDisabled:  !genesisSafetyParams.IsDisabled,
					DelayBlocks: genesisSafetyParams.DelayBlocks + 1,
				},
			},
			expectedProposalStatus: govtypesv1.ProposalStatus_PROPOSAL_STATUS_PASSED,
		},
		"Fail: invalid authority": {
			msg: &bridgetypes.MsgUpdateSafetyParams{
				Authority: constants.AliceAccAddress.String(),
				Params: bridgetypes.SafetyParams{
					IsDisabled:  !genesisSafetyParams.IsDisabled,
					DelayBlocks: genesisSafetyParams.DelayBlocks + 1,
				},
			},
			expectSubmitProposalFail: true,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tApp := testapp.NewTestAppBuilder().WithGenesisDocFn(func() (genesis types.GenesisDoc) {
				genesis = testapp.DefaultGenesis()
				testapp.UpdateGenesisDocWithAppStateForModule(
					&genesis,
					func(genesisState *govtypesv1.GenesisState) {
						genesisState.Params.VotingPeriod = &testapp.TestVotingPeriod
					},
				)
				return genesis
			}).WithTesting(t).Build()
			ctx := tApp.InitChain()

			// Submit and tally governance proposal that includes `MsgUpdateSafetyParams`.
			ctx = testapp.SubmitAndTallyProposal(
				t,
				ctx,
				&tApp,
				[]sdk.Msg{tc.msg},
				tc.expectSubmitProposalFail,
				tc.expectedProposalStatus,
			)

			// If proposal is supposed to pass, verify that safety params are updated.
			if tc.expectedProposalStatus == govtypesv1.ProposalStatus_PROPOSAL_STATUS_PASSED {
				updatedSafetyParams := tApp.App.BridgeKeeper.GetSafetyParams(ctx)
				require.Equal(t, tc.msg.Params, updatedSafetyParams)
			}
		})
	}
}
