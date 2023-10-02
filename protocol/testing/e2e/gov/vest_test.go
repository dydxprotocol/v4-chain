package gov_test

import (
	"testing"
	"time"

	"github.com/cometbft/cometbft/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	govtypesv1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
	testapp "github.com/dydxprotocol/v4-chain/protocol/testutil/app"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	vesttypes "github.com/dydxprotocol/v4-chain/protocol/x/vest/types"
	"github.com/stretchr/testify/require"
)

func TestSetVestEntry(t *testing.T) {
	tests := map[string]struct {
		msg                      *vesttypes.MsgSetVestEntry
		updateExistingVestEntry  bool // whether above msg should update an existing vest entry.
		expectSubmitProposalFail bool
		expectedProposalStatus   govtypesv1.ProposalStatus
	}{
		"Success: create a new vest entry": {
			msg: &vesttypes.MsgSetVestEntry{
				Authority: authtypes.NewModuleAddress(govtypes.ModuleName).String(),
				Entry: vesttypes.VestEntry{
					VesterAccount:   "random_vester",
					TreasuryAccount: "random_treasury",
					Denom:           "dv4tnt",
					StartTime:       time.Date(2023, 10, 2, 0, 0, 0, 0, time.UTC),
					EndTime:         time.Date(2024, 10, 1, 0, 0, 0, 0, time.UTC),
				},
			},
			expectedProposalStatus: govtypesv1.ProposalStatus_PROPOSAL_STATUS_PASSED,
		},
		"Success: update an existing vest entry": {
			msg: &vesttypes.MsgSetVestEntry{
				Authority: authtypes.NewModuleAddress(govtypes.ModuleName).String(),
				Entry: vesttypes.VestEntry{
					VesterAccount:   "random_vester",
					TreasuryAccount: "random_treasury",
					Denom:           "dv4tnt",
					StartTime:       time.Date(2023, 10, 2, 0, 0, 0, 0, time.UTC),
					EndTime:         time.Date(2024, 10, 1, 0, 0, 0, 0, time.UTC),
				},
			},
			updateExistingVestEntry: true,
			expectedProposalStatus:  govtypesv1.ProposalStatus_PROPOSAL_STATUS_PASSED,
		},
		"Fail: invalid authority": {
			msg: &vesttypes.MsgSetVestEntry{
				Authority: constants.BobAccAddress.String(),
				Entry: vesttypes.VestEntry{
					VesterAccount:   vesttypes.ModuleName,
					TreasuryAccount: "random_treasury",
					Denom:           "dv4tnt",
					StartTime:       time.Date(2023, 10, 2, 0, 0, 0, 0, time.UTC),
					EndTime:         time.Date(2024, 10, 1, 0, 0, 0, 0, time.UTC),
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
				// If updating an existing vest entry, initialize vest module with a vest entry that has the
				// same key. Otherwise, initialize vest module with no vest entry.
				genesisVestEntries := []vesttypes.VestEntry{}
				if tc.updateExistingVestEntry {
					genesisVestEntries = append(genesisVestEntries, vesttypes.VestEntry{
						VesterAccount:   tc.msg.Entry.VesterAccount,
						TreasuryAccount: tc.msg.Entry.TreasuryAccount + "_initial",
						Denom:           tc.msg.Entry.Denom + "_initial",
						StartTime:       tc.msg.Entry.StartTime.AddDate(0, 0, 1),
						EndTime:         tc.msg.Entry.EndTime.AddDate(0, 0, 1),
					})
				}
				testapp.UpdateGenesisDocWithAppStateForModule(
					&genesis,
					func(genesisState *vesttypes.GenesisState) {
						genesisState.VestEntries = genesisVestEntries
					},
				)
				return genesis
			}).WithTesting(t).Build()
			ctx := tApp.InitChain()
			initialVestEntries := tApp.App.VestKeeper.GetAllVestEntries(ctx)

			// Submit and tally governance proposal that includes `MsgSetVestEntry`.
			ctx = testapp.SubmitAndTallyProposal(
				t,
				ctx,
				&tApp,
				[]sdk.Msg{tc.msg},
				tc.expectSubmitProposalFail,
				tc.expectedProposalStatus,
			)

			// If governance proposal is supposed to fail submission, verify that vest entries in state
			// match vest entries before proposal submission.
			if tc.expectSubmitProposalFail {
				require.Equal(t, initialVestEntries, tApp.App.VestKeeper.GetAllVestEntries(ctx))
			}

			// If proposal is supposed to pass, verify that expected vest entry is set in state.
			if tc.expectedProposalStatus == govtypesv1.ProposalStatus_PROPOSAL_STATUS_PASSED {
				vestEntry, err := tApp.App.VestKeeper.GetVestEntry(ctx, tc.msg.Entry.VesterAccount)
				require.NoError(t, err)
				require.Equal(t, tc.msg.Entry, vestEntry)
			}
		})
	}
}

func TestDeleteVestEntry(t *testing.T) {
	tests := map[string]struct {
		msg                      *vesttypes.MsgDeleteVestEntry
		expectSubmitProposalFail bool
		expectedProposalStatus   govtypesv1.ProposalStatus
	}{
		"Success": {
			msg: &vesttypes.MsgDeleteVestEntry{
				Authority:     authtypes.NewModuleAddress(govtypes.ModuleName).String(),
				VesterAccount: "random_vester",
			},
			expectedProposalStatus: govtypesv1.ProposalStatus_PROPOSAL_STATUS_PASSED,
		},
		"Fail: invalid authority": {
			msg: &vesttypes.MsgDeleteVestEntry{
				Authority:     constants.BobAccAddress.String(),
				VesterAccount: "random_vester",
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
				// Initialize vest module with a vest entry that has the same key as the one to be deleted.
				testapp.UpdateGenesisDocWithAppStateForModule(
					&genesis,
					func(genesisState *vesttypes.GenesisState) {
						genesisState.VestEntries = []vesttypes.VestEntry{
							{
								VesterAccount:   tc.msg.VesterAccount,
								TreasuryAccount: "random_treasury",
								Denom:           "dv4tnt",
								StartTime:       time.Date(2023, 10, 2, 0, 0, 0, 0, time.UTC),
								EndTime:         time.Date(2024, 10, 1, 0, 0, 0, 0, time.UTC),
							},
						}
					},
				)
				return genesis
			}).WithTesting(t).Build()
			ctx := tApp.InitChain()
			initialVestEntries := tApp.App.VestKeeper.GetAllVestEntries(ctx)

			// Submit and tally governance proposal that includes `MsgDeleteVestEntry`.
			ctx = testapp.SubmitAndTallyProposal(
				t,
				ctx,
				&tApp,
				[]sdk.Msg{tc.msg},
				tc.expectSubmitProposalFail,
				tc.expectedProposalStatus,
			)

			// If governance proposal is supposed to fail submission, verify that vest entries in state
			// match vest entries before proposal submission.
			if tc.expectSubmitProposalFail {
				require.Equal(t, initialVestEntries, tApp.App.VestKeeper.GetAllVestEntries(ctx))
			}

			// If proposal is supposed to pass, verify that the vest entry has been deleted.
			if tc.expectedProposalStatus == govtypesv1.ProposalStatus_PROPOSAL_STATUS_PASSED {
				vestEntry, err := tApp.App.VestKeeper.GetVestEntry(ctx, tc.msg.VesterAccount)
				require.Equal(t, vesttypes.VestEntry{}, vestEntry)
				require.Error(t, err)
			}
		})
	}
}
