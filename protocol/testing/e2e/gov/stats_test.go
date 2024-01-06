package gov_test

import (
	"testing"
	"time"

	"github.com/dydxprotocol/v4-chain/protocol/lib"

	"github.com/cometbft/cometbft/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	govtypesv1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
	testapp "github.com/dydxprotocol/v4-chain/protocol/testutil/app"
	statstypes "github.com/dydxprotocol/v4-chain/protocol/x/stats/types"
	"github.com/stretchr/testify/require"
)

func TestUpdateParams(t *testing.T) {
	tests := map[string]struct {
		msg                      *statstypes.MsgUpdateParams
		expectCheckTxFails       bool
		expectedProposalStatus   govtypesv1.ProposalStatus
		expectSubmitProposalFail bool
	}{
		"Success": {
			msg: &statstypes.MsgUpdateParams{
				Authority: lib.GovModuleAddress.String(),
				Params: statstypes.Params{
					WindowDuration: time.Hour,
				},
			},
			expectedProposalStatus: govtypesv1.ProposalStatus_PROPOSAL_STATUS_PASSED,
		},
		"Fail: invalid authority": {
			msg: &statstypes.MsgUpdateParams{
				Authority: authtypes.NewModuleAddress(statstypes.ModuleName).String(),
				Params: statstypes.Params{
					WindowDuration: 100,
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
				// Initialize stats module with params that are different from the proposal.
				testapp.UpdateGenesisDocWithAppStateForModule(
					&genesis,
					func(genesisState *statstypes.GenesisState) {
						genesisState.Params = statstypes.Params{
							WindowDuration: tc.msg.Params.WindowDuration + time.Second,
						}
					},
				)
				return genesis
			}).Build()
			ctx := tApp.InitChain()
			initialParams := tApp.App.StatsKeeper.GetParams(ctx)

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

			// If governance proposal is supposed to fail submission, verify that stats params
			// match the ones before proposal submission.
			if tc.expectSubmitProposalFail {
				require.Equal(t, initialParams, tApp.App.StatsKeeper.GetParams(ctx))
			}

			// If proposal is supposed to pass, verify that stats params have been updated.
			if tc.expectedProposalStatus == govtypesv1.ProposalStatus_PROPOSAL_STATUS_PASSED {
				require.Equal(t, tc.msg.Params, tApp.App.StatsKeeper.GetParams(ctx))
			}
		})
	}
}
