package gov_test

import (
	"testing"

	"github.com/cometbft/cometbft/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	govtypesv1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
	"github.com/dydxprotocol/v4-chain/protocol/lib"
	testapp "github.com/dydxprotocol/v4-chain/protocol/testutil/app"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	bridgetypes "github.com/dydxprotocol/v4-chain/protocol/x/bridge/types"
	"github.com/stretchr/testify/require"
)

func TestUpdateEventParams(t *testing.T) {
	genesisEventParams := bridgetypes.DefaultGenesis().EventParams

	tests := map[string]struct {
		msg                       *bridgetypes.MsgUpdateEventParams
		expectCheckTxFails        bool
		expectSubmitProposalFails bool
		expectedProposalStatus    govtypesv1.ProposalStatus
	}{
		"Success": {
			msg: &bridgetypes.MsgUpdateEventParams{
				Authority: lib.GovModuleAddress.String(),
				Params: bridgetypes.EventParams{
					Denom:      genesisEventParams.Denom + "updated",
					EthChainId: genesisEventParams.EthChainId + 1,
					EthAddress: genesisEventParams.EthAddress,
				},
			},
			expectedProposalStatus: govtypesv1.ProposalStatus_PROPOSAL_STATUS_PASSED,
		},
		"Failure: empty eth address": {
			msg: &bridgetypes.MsgUpdateEventParams{
				Authority: authtypes.NewModuleAddress(govtypes.ModuleName).String(),
				Params: bridgetypes.EventParams{
					Denom:      genesisEventParams.Denom,
					EthChainId: genesisEventParams.EthChainId,
					EthAddress: "",
				},
			},
			expectCheckTxFails: true,
		},
		"Failure: invalid authority": {
			msg: &bridgetypes.MsgUpdateEventParams{
				Authority: constants.BobAccAddress.String(),
				Params: bridgetypes.EventParams{
					Denom:      genesisEventParams.Denom + "updated",
					EthChainId: genesisEventParams.EthChainId + 1,
					EthAddress: genesisEventParams.EthAddress,
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
				return genesis
			}).Build()
			ctx := tApp.InitChain()
			initialEventParams := tApp.App.BridgeKeeper.GetEventParams(ctx)

			// Submit and tally governance proposal that includes `MsgUpdateEventParams`.
			ctx = testapp.SubmitAndTallyProposal(
				t,
				ctx,
				tApp,
				[]sdk.Msg{tc.msg},
				tc.expectCheckTxFails,
				tc.expectSubmitProposalFails,
				tc.expectedProposalStatus,
			)

			if tc.expectedProposalStatus == govtypesv1.ProposalStatus_PROPOSAL_STATUS_PASSED {
				// If proposal is supposed to pass, verify that event params are updated.
				updatedEventParams := tApp.App.BridgeKeeper.GetEventParams(ctx)
				require.Equal(t, tc.msg.Params, updatedEventParams)
			} else {
				// Otherwise, verify that event params are unchanged.
				require.Equal(t, initialEventParams, tApp.App.BridgeKeeper.GetEventParams(ctx))
			}
		})
	}
}

func TestUpdateProposeParams(t *testing.T) {
	genesisProposeParams := bridgetypes.DefaultGenesis().ProposeParams

	tests := map[string]struct {
		msg                       *bridgetypes.MsgUpdateProposeParams
		expectCheckTxFails        bool
		expectSubmitProposalFails bool
		expectedProposalStatus    govtypesv1.ProposalStatus
	}{
		"Success": {
			msg: &bridgetypes.MsgUpdateProposeParams{
				Authority: lib.GovModuleAddress.String(),
				Params: bridgetypes.ProposeParams{
					MaxBridgesPerBlock:           genesisProposeParams.MaxBridgesPerBlock + 1,
					ProposeDelayDuration:         genesisProposeParams.ProposeDelayDuration + 1,
					SkipRatePpm:                  genesisProposeParams.SkipRatePpm + 1,
					SkipIfBlockDelayedByDuration: genesisProposeParams.SkipIfBlockDelayedByDuration + 1,
				},
			},
			expectedProposalStatus: govtypesv1.ProposalStatus_PROPOSAL_STATUS_PASSED,
		},
		"Failure: negative propose delay duration": {
			msg: &bridgetypes.MsgUpdateProposeParams{
				Authority: authtypes.NewModuleAddress(govtypes.ModuleName).String(),
				Params: bridgetypes.ProposeParams{
					MaxBridgesPerBlock:           genesisProposeParams.MaxBridgesPerBlock,
					ProposeDelayDuration:         -genesisProposeParams.ProposeDelayDuration,
					SkipRatePpm:                  genesisProposeParams.SkipRatePpm,
					SkipIfBlockDelayedByDuration: genesisProposeParams.SkipIfBlockDelayedByDuration,
				},
			},
			expectCheckTxFails: true,
		},
		"Failure: negative skip if block delayed by duration": {
			msg: &bridgetypes.MsgUpdateProposeParams{
				Authority: authtypes.NewModuleAddress(govtypes.ModuleName).String(),
				Params: bridgetypes.ProposeParams{
					MaxBridgesPerBlock:           genesisProposeParams.MaxBridgesPerBlock,
					ProposeDelayDuration:         genesisProposeParams.ProposeDelayDuration,
					SkipRatePpm:                  genesisProposeParams.SkipRatePpm,
					SkipIfBlockDelayedByDuration: -genesisProposeParams.SkipIfBlockDelayedByDuration,
				},
			},
			expectCheckTxFails: true,
		},
		"Failure: skip rate ppm out of bounds": {
			msg: &bridgetypes.MsgUpdateProposeParams{
				Authority: authtypes.NewModuleAddress(govtypes.ModuleName).String(),
				Params: bridgetypes.ProposeParams{
					MaxBridgesPerBlock:           genesisProposeParams.MaxBridgesPerBlock,
					ProposeDelayDuration:         genesisProposeParams.ProposeDelayDuration,
					SkipRatePpm:                  1_000_001, // greater than 1 million.
					SkipIfBlockDelayedByDuration: genesisProposeParams.SkipIfBlockDelayedByDuration,
				},
			},
			expectCheckTxFails: true,
		},
		"Failure: invalid authority": {
			msg: &bridgetypes.MsgUpdateProposeParams{
				Authority: constants.AliceAccAddress.String(),
				Params: bridgetypes.ProposeParams{
					MaxBridgesPerBlock:           genesisProposeParams.MaxBridgesPerBlock + 1,
					ProposeDelayDuration:         genesisProposeParams.ProposeDelayDuration + 1,
					SkipRatePpm:                  genesisProposeParams.SkipRatePpm + 1,
					SkipIfBlockDelayedByDuration: genesisProposeParams.SkipIfBlockDelayedByDuration + 1,
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
				return genesis
			}).Build()
			ctx := tApp.InitChain()
			initialProposeParams := tApp.App.BridgeKeeper.GetProposeParams(ctx)

			// Submit and tally governance proposal that includes `MsgUpdateProposeParams`.
			ctx = testapp.SubmitAndTallyProposal(
				t,
				ctx,
				tApp,
				[]sdk.Msg{tc.msg},
				tc.expectCheckTxFails,
				tc.expectSubmitProposalFails,
				tc.expectedProposalStatus,
			)

			if tc.expectedProposalStatus == govtypesv1.ProposalStatus_PROPOSAL_STATUS_PASSED {
				// If proposal is supposed to pass, verify that propose params are updated.
				updatedProposeParams := tApp.App.BridgeKeeper.GetProposeParams(ctx)
				require.Equal(t, tc.msg.Params, updatedProposeParams)
			} else {
				// Otherwise, verify that propose params are unchanged.
				require.Equal(t, initialProposeParams, tApp.App.BridgeKeeper.GetProposeParams(ctx))
			}
		})
	}
}

func TestUpdateSafetyParams(t *testing.T) {
	genesisSafetyParams := bridgetypes.DefaultGenesis().SafetyParams

	tests := map[string]struct {
		msg                       *bridgetypes.MsgUpdateSafetyParams
		expectSubmitProposalFails bool
		expectedProposalStatus    govtypesv1.ProposalStatus
	}{
		"Success": {
			msg: &bridgetypes.MsgUpdateSafetyParams{
				Authority: lib.GovModuleAddress.String(),
				Params: bridgetypes.SafetyParams{
					IsDisabled:  !genesisSafetyParams.IsDisabled,
					DelayBlocks: genesisSafetyParams.DelayBlocks + 1,
				},
			},
			expectedProposalStatus: govtypesv1.ProposalStatus_PROPOSAL_STATUS_PASSED,
		},
		"Failure: invalid authority": {
			msg: &bridgetypes.MsgUpdateSafetyParams{
				Authority: constants.AliceAccAddress.String(),
				Params: bridgetypes.SafetyParams{
					IsDisabled:  !genesisSafetyParams.IsDisabled,
					DelayBlocks: genesisSafetyParams.DelayBlocks + 1,
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
				return genesis
			}).Build()
			ctx := tApp.InitChain()
			initialSafetyParams := tApp.App.BridgeKeeper.GetSafetyParams(ctx)

			// Submit and tally governance proposal that includes `MsgUpdateSafetyParams`.
			ctx = testapp.SubmitAndTallyProposal(
				t,
				ctx,
				tApp,
				[]sdk.Msg{tc.msg},
				false,
				tc.expectSubmitProposalFails,
				tc.expectedProposalStatus,
			)

			if tc.expectedProposalStatus == govtypesv1.ProposalStatus_PROPOSAL_STATUS_PASSED {
				// If proposal is supposed to pass, verify that safety params are updated.
				updatedSafetyParams := tApp.App.BridgeKeeper.GetSafetyParams(ctx)
				require.Equal(t, tc.msg.Params, updatedSafetyParams)
			} else {
				// Otherwise, verify that safety params are unchanged.
				require.Equal(t, initialSafetyParams, tApp.App.BridgeKeeper.GetSafetyParams(ctx))
			}
		})
	}
}
