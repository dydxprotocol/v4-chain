package gov_test

import (
	"testing"
	"time"

	"github.com/cometbft/cometbft/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	govtypesv1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
	"github.com/dydxprotocol/v4-chain/protocol/lib"
	testapp "github.com/dydxprotocol/v4-chain/protocol/testutil/app"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	bridgetypes "github.com/dydxprotocol/v4-chain/protocol/x/bridge/types"
	"github.com/stretchr/testify/require"
)

var (
	// Genesis params.
	GenesisEventParams = bridgetypes.EventParams{
		Denom:      "testdenom",
		EthChainId: 123,
		EthAddress: "0x0123",
	}
	GenesisProposeParams = bridgetypes.ProposeParams{
		MaxBridgesPerBlock:           10,
		ProposeDelayDuration:         time.Minute,
		SkipRatePpm:                  800_000,
		SkipIfBlockDelayedByDuration: time.Minute,
	}
	GenesisSafetyParams = bridgetypes.SafetyParams{
		IsDisabled:  false,
		DelayBlocks: 10,
	}
	// Modified params.
	ModifiedEventParams = bridgetypes.EventParams{
		Denom:      "advtnt",
		EthChainId: 1,
		EthAddress: "0xabcd",
	}
	ModifiedProposeParams = bridgetypes.ProposeParams{
		MaxBridgesPerBlock:           7,
		ProposeDelayDuration:         time.Second,
		SkipRatePpm:                  700_007,
		SkipIfBlockDelayedByDuration: time.Second,
	}
	ModifiedSafetyParams = bridgetypes.SafetyParams{
		IsDisabled:  true,
		DelayBlocks: 5,
	}
	// Invalid authority address.
	InvalidBridgeAuthority = constants.AliceAccAddress.String()
)

func TestUpdateEventParams(t *testing.T) {
	tests := map[string]struct {
		msg                       *bridgetypes.MsgUpdateEventParams
		expectCheckTxFails        bool
		expectSubmitProposalFails bool
		expectedProposalStatus    govtypesv1.ProposalStatus
	}{
		"Success": {
			msg: &bridgetypes.MsgUpdateEventParams{
				Authority: lib.GovModuleAddress.String(),
				Params:    ModifiedEventParams,
			},
			expectedProposalStatus: govtypesv1.ProposalStatus_PROPOSAL_STATUS_PASSED,
		},
		"Failure: empty eth address": {
			msg: &bridgetypes.MsgUpdateEventParams{
				Authority: lib.GovModuleAddress.String(),
				Params: bridgetypes.EventParams{
					Denom:      ModifiedEventParams.Denom,
					EthChainId: ModifiedEventParams.EthChainId,
					EthAddress: "",
				},
			},
			expectCheckTxFails: true,
		},
		"Failure: invalid authority": {
			msg: &bridgetypes.MsgUpdateEventParams{
				Authority: constants.BobAccAddress.String(),
				Params:    ModifiedEventParams,
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
				testapp.UpdateGenesisDocWithAppStateForModule(
					&genesis,
					func(genesisState *bridgetypes.GenesisState) {
						genesisState.EventParams = GenesisEventParams
					},
				)
				return genesis
			}).Build()
			ctx := tApp.InitChain()

			// Submit and tally governance proposal that includes `MsgUpdateEventParams`.
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
				// If proposal is supposed to pass, verify that event params are updated.
				updatedEventParams := tApp.App.BridgeKeeper.GetEventParams(ctx)
				require.Equal(t, tc.msg.Params, updatedEventParams)
			} else {
				// Otherwise, verify that event params are unchanged.
				require.Equal(t, GenesisEventParams, tApp.App.BridgeKeeper.GetEventParams(ctx))
			}
		})
	}
}

func TestUpdateProposeParams(t *testing.T) {
	tests := map[string]struct {
		msg                       *bridgetypes.MsgUpdateProposeParams
		expectCheckTxFails        bool
		expectSubmitProposalFails bool
		expectedProposalStatus    govtypesv1.ProposalStatus
	}{
		"Success": {
			msg: &bridgetypes.MsgUpdateProposeParams{
				Authority: lib.GovModuleAddress.String(),
				Params:    ModifiedProposeParams,
			},
			expectedProposalStatus: govtypesv1.ProposalStatus_PROPOSAL_STATUS_PASSED,
		},
		"Failure: negative propose delay duration": {
			msg: &bridgetypes.MsgUpdateProposeParams{
				Authority: lib.GovModuleAddress.String(),
				Params: bridgetypes.ProposeParams{
					MaxBridgesPerBlock:           ModifiedProposeParams.MaxBridgesPerBlock,
					ProposeDelayDuration:         -time.Second,
					SkipRatePpm:                  ModifiedProposeParams.SkipRatePpm,
					SkipIfBlockDelayedByDuration: ModifiedProposeParams.SkipIfBlockDelayedByDuration,
				},
			},
			expectCheckTxFails: true,
		},
		"Failure: negative skip if block delayed by duration": {
			msg: &bridgetypes.MsgUpdateProposeParams{
				Authority: lib.GovModuleAddress.String(),
				Params: bridgetypes.ProposeParams{
					MaxBridgesPerBlock:           ModifiedProposeParams.MaxBridgesPerBlock,
					ProposeDelayDuration:         ModifiedProposeParams.ProposeDelayDuration,
					SkipRatePpm:                  ModifiedProposeParams.SkipRatePpm,
					SkipIfBlockDelayedByDuration: -time.Second,
				},
			},
			expectCheckTxFails: true,
		},
		"Failure: skip rate ppm out of bounds": {
			msg: &bridgetypes.MsgUpdateProposeParams{
				Authority: lib.GovModuleAddress.String(),
				Params: bridgetypes.ProposeParams{
					MaxBridgesPerBlock:           ModifiedProposeParams.MaxBridgesPerBlock,
					ProposeDelayDuration:         ModifiedProposeParams.ProposeDelayDuration,
					SkipRatePpm:                  1_000_001, // greater than 1 million.
					SkipIfBlockDelayedByDuration: ModifiedProposeParams.SkipIfBlockDelayedByDuration,
				},
			},
			expectCheckTxFails: true,
		},
		"Failure: invalid authority": {
			msg: &bridgetypes.MsgUpdateProposeParams{
				Authority: InvalidBridgeAuthority,
				Params:    ModifiedProposeParams,
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
				testapp.UpdateGenesisDocWithAppStateForModule(
					&genesis,
					func(genesisState *bridgetypes.GenesisState) {
						genesisState.ProposeParams = GenesisProposeParams
					},
				)
				return genesis
			}).Build()
			ctx := tApp.InitChain()

			// Submit and tally governance proposal that includes `MsgUpdateProposeParams`.
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
				// If proposal is supposed to pass, verify that propose params are updated.
				updatedProposeParams := tApp.App.BridgeKeeper.GetProposeParams(ctx)
				require.Equal(t, tc.msg.Params, updatedProposeParams)
			} else {
				// Otherwise, verify that propose params are unchanged.
				require.Equal(t, GenesisProposeParams, tApp.App.BridgeKeeper.GetProposeParams(ctx))
			}
		})
	}
}

func TestUpdateSafetyParams(t *testing.T) {
	tests := map[string]struct {
		msg                       *bridgetypes.MsgUpdateSafetyParams
		expectSubmitProposalFails bool
		expectedProposalStatus    govtypesv1.ProposalStatus
	}{
		"Success": {
			msg: &bridgetypes.MsgUpdateSafetyParams{
				Authority: lib.GovModuleAddress.String(),
				Params:    ModifiedSafetyParams,
			},
			expectedProposalStatus: govtypesv1.ProposalStatus_PROPOSAL_STATUS_PASSED,
		},
		"Failure: invalid authority": {
			msg: &bridgetypes.MsgUpdateSafetyParams{
				Authority: InvalidBridgeAuthority,
				Params:    ModifiedSafetyParams,
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
				testapp.UpdateGenesisDocWithAppStateForModule(
					&genesis,
					func(genesisState *bridgetypes.GenesisState) {
						genesisState.SafetyParams = GenesisSafetyParams
					},
				)
				return genesis
			}).Build()
			ctx := tApp.InitChain()

			// Submit and tally governance proposal that includes `MsgUpdateSafetyParams`.
			ctx = testapp.SubmitAndTallyProposal(
				t,
				ctx,
				tApp,
				[]sdk.Msg{tc.msg},
				testapp.TestSubmitProposalTxHeight,
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
				require.Equal(t, GenesisSafetyParams, tApp.App.BridgeKeeper.GetSafetyParams(ctx))
			}
		})
	}
}
