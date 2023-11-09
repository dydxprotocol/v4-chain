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
	GENESIS_EVENT_PARAMS = bridgetypes.EventParams{
		Denom:      "testdenom",
		EthChainId: 123,
		EthAddress: "0x0123",
	}
	GENESIS_PROPOSE_PARAMS = bridgetypes.ProposeParams{
		MaxBridgesPerBlock:           10,
		ProposeDelayDuration:         time.Minute,
		SkipRatePpm:                  800_000,
		SkipIfBlockDelayedByDuration: time.Minute,
	}
	GENESIS_SAFETY_PARAMS = bridgetypes.SafetyParams{
		IsDisabled:  false,
		DelayBlocks: 10,
	}
	INVALID_BRIDGE_AUTHORITY = constants.AliceAccAddress.String()
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
				Params: bridgetypes.EventParams{
					Denom:      "adv4tnt",
					EthChainId: 1,
					EthAddress: "0xabcd",
				},
			},
			expectedProposalStatus: govtypesv1.ProposalStatus_PROPOSAL_STATUS_PASSED,
		},
		"Failure: empty eth address": {
			msg: &bridgetypes.MsgUpdateEventParams{
				Authority: lib.GovModuleAddress.String(),
				Params: bridgetypes.EventParams{
					Denom:      "adv4tnt",
					EthChainId: 1,
					EthAddress: "",
				},
			},
			expectCheckTxFails: true,
		},
		"Failure: invalid authority": {
			msg: &bridgetypes.MsgUpdateEventParams{
				Authority: constants.BobAccAddress.String(),
				Params: bridgetypes.EventParams{
					Denom:      "adv4tnt",
					EthChainId: 1,
					EthAddress: "0xabcd",
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
				testapp.UpdateGenesisDocWithAppStateForModule(
					&genesis,
					func(genesisState *bridgetypes.GenesisState) {
						genesisState.EventParams = GENESIS_EVENT_PARAMS
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
				require.Equal(t, GENESIS_EVENT_PARAMS, tApp.App.BridgeKeeper.GetEventParams(ctx))
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
				Params: bridgetypes.ProposeParams{
					MaxBridgesPerBlock:           7,
					ProposeDelayDuration:         time.Second,
					SkipRatePpm:                  700_001,
					SkipIfBlockDelayedByDuration: time.Second,
				},
			},
			expectedProposalStatus: govtypesv1.ProposalStatus_PROPOSAL_STATUS_PASSED,
		},
		"Failure: negative propose delay duration": {
			msg: &bridgetypes.MsgUpdateProposeParams{
				Authority: lib.GovModuleAddress.String(),
				Params: bridgetypes.ProposeParams{
					MaxBridgesPerBlock:           7,
					ProposeDelayDuration:         -time.Second,
					SkipRatePpm:                  700_001,
					SkipIfBlockDelayedByDuration: time.Second,
				},
			},
			expectCheckTxFails: true,
		},
		"Failure: negative skip if block delayed by duration": {
			msg: &bridgetypes.MsgUpdateProposeParams{
				Authority: lib.GovModuleAddress.String(),
				Params: bridgetypes.ProposeParams{
					MaxBridgesPerBlock:           7,
					ProposeDelayDuration:         time.Second,
					SkipRatePpm:                  700_001,
					SkipIfBlockDelayedByDuration: -time.Second,
				},
			},
			expectCheckTxFails: true,
		},
		"Failure: skip rate ppm out of bounds": {
			msg: &bridgetypes.MsgUpdateProposeParams{
				Authority: lib.GovModuleAddress.String(),
				Params: bridgetypes.ProposeParams{
					MaxBridgesPerBlock:           7,
					ProposeDelayDuration:         time.Second,
					SkipRatePpm:                  1_000_001, // greater than 1 million.
					SkipIfBlockDelayedByDuration: time.Second,
				},
			},
			expectCheckTxFails: true,
		},
		"Failure: invalid authority": {
			msg: &bridgetypes.MsgUpdateProposeParams{
				Authority: INVALID_BRIDGE_AUTHORITY,
				Params: bridgetypes.ProposeParams{
					MaxBridgesPerBlock:           7,
					ProposeDelayDuration:         time.Second,
					SkipRatePpm:                  700_001,
					SkipIfBlockDelayedByDuration: time.Second,
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
				testapp.UpdateGenesisDocWithAppStateForModule(
					&genesis,
					func(genesisState *bridgetypes.GenesisState) {
						genesisState.ProposeParams = GENESIS_PROPOSE_PARAMS
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
				require.Equal(t, GENESIS_PROPOSE_PARAMS, tApp.App.BridgeKeeper.GetProposeParams(ctx))
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
				Params: bridgetypes.SafetyParams{
					IsDisabled:  true,
					DelayBlocks: 123,
				},
			},
			expectedProposalStatus: govtypesv1.ProposalStatus_PROPOSAL_STATUS_PASSED,
		},
		"Failure: invalid authority": {
			msg: &bridgetypes.MsgUpdateSafetyParams{
				Authority: INVALID_BRIDGE_AUTHORITY,
				Params: bridgetypes.SafetyParams{
					IsDisabled:  true,
					DelayBlocks: 123,
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
				testapp.UpdateGenesisDocWithAppStateForModule(
					&genesis,
					func(genesisState *bridgetypes.GenesisState) {
						genesisState.SafetyParams = GENESIS_SAFETY_PARAMS
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
				require.Equal(t, GENESIS_SAFETY_PARAMS, tApp.App.BridgeKeeper.GetSafetyParams(ctx))
			}
		})
	}
}
