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

var (
	TEST_START_TIME         = time.Date(2023, 10, 2, 0, 0, 0, 0, time.UTC)
	TEST_END_TIME           = time.Date(2024, 10, 1, 0, 0, 0, 0, time.UTC)
	TEST_VESTER_ACCOUNT_1   = "random_vester"
	TEST_VESTER_ACCOUNT_2   = "random_vester_2"
	TEST_GENESIS_VEST_ENTRY = vesttypes.VestEntry{
		VesterAccount:   "genesis_vester",
		TreasuryAccount: "genesis_treasury",
		Denom:           "genesis_denom",
		StartTime:       TEST_START_TIME.AddDate(-2022, -1, -1),
		EndTime:         TEST_END_TIME.AddDate(-2022, -1, -1),
	}
)

func TestSetVestEntry_Success(t *testing.T) {
	tests := map[string]struct {
		msgs                 []sdk.Msg
		genesisVestEntryKeys []string // keys of vest entries in genesis state.
	}{
		"Success: create a new vest entry": {
			msgs: []sdk.Msg{
				&vesttypes.MsgSetVestEntry{
					Authority: authtypes.NewModuleAddress(govtypes.ModuleName).String(),
					Entry: vesttypes.VestEntry{
						VesterAccount:   TEST_VESTER_ACCOUNT_1,
						TreasuryAccount: "random_treasury",
						Denom:           "adv4tnt",
						StartTime:       TEST_START_TIME,
						EndTime:         TEST_END_TIME,
					},
				},
			},
		},
		"Success: update an existing vest entry": {
			msgs: []sdk.Msg{
				&vesttypes.MsgSetVestEntry{
					Authority: authtypes.NewModuleAddress(govtypes.ModuleName).String(),
					Entry: vesttypes.VestEntry{
						VesterAccount:   TEST_VESTER_ACCOUNT_1,
						TreasuryAccount: "random_treasury",
						Denom:           "adv4tnt",
						StartTime:       TEST_START_TIME,
						EndTime:         TEST_END_TIME,
					},
				},
			},
			genesisVestEntryKeys: []string{TEST_VESTER_ACCOUNT_1},
		},
		"Success: create two new vest entries": {
			msgs: []sdk.Msg{
				&vesttypes.MsgSetVestEntry{
					Authority: authtypes.NewModuleAddress(govtypes.ModuleName).String(),
					Entry: vesttypes.VestEntry{
						VesterAccount:   TEST_VESTER_ACCOUNT_1,
						TreasuryAccount: "random_treasury",
						Denom:           "adv4tnt",
						StartTime:       TEST_START_TIME,
						EndTime:         TEST_END_TIME,
					},
				},
				&vesttypes.MsgSetVestEntry{
					Authority: authtypes.NewModuleAddress(govtypes.ModuleName).String(),
					Entry: vesttypes.VestEntry{
						VesterAccount:   TEST_VESTER_ACCOUNT_2,
						TreasuryAccount: "random_treasury_2",
						Denom:           "adv4tnt",
						StartTime:       TEST_START_TIME.Add(time.Hour),
						EndTime:         TEST_END_TIME.Add(time.Hour),
					},
				},
			},
		},
		"Success: create and then update a vest entry": {
			msgs: []sdk.Msg{
				&vesttypes.MsgSetVestEntry{
					Authority: authtypes.NewModuleAddress(govtypes.ModuleName).String(),
					Entry: vesttypes.VestEntry{
						VesterAccount:   TEST_VESTER_ACCOUNT_1,
						TreasuryAccount: "random_treasury",
						Denom:           "adv4tnt",
						StartTime:       TEST_START_TIME,
						EndTime:         TEST_END_TIME,
					},
				},
				&vesttypes.MsgSetVestEntry{
					Authority: authtypes.NewModuleAddress(govtypes.ModuleName).String(),
					Entry: vesttypes.VestEntry{
						VesterAccount:   TEST_VESTER_ACCOUNT_1,
						TreasuryAccount: "random_treasury_2",
						Denom:           "adv4tnt",
						StartTime:       TEST_START_TIME.Add(time.Hour),
						EndTime:         TEST_END_TIME.Add(time.Hour),
					},
				},
			},
		},
		"Success: update a vest entry twice": {
			msgs: []sdk.Msg{
				&vesttypes.MsgSetVestEntry{
					Authority: authtypes.NewModuleAddress(govtypes.ModuleName).String(),
					Entry: vesttypes.VestEntry{
						VesterAccount:   TEST_VESTER_ACCOUNT_1,
						TreasuryAccount: "random_treasury",
						Denom:           "adv4tnt",
						StartTime:       TEST_START_TIME,
						EndTime:         TEST_END_TIME,
					},
				},
				&vesttypes.MsgSetVestEntry{
					Authority: authtypes.NewModuleAddress(govtypes.ModuleName).String(),
					Entry: vesttypes.VestEntry{
						VesterAccount:   TEST_VESTER_ACCOUNT_1,
						TreasuryAccount: "random_treasury_2",
						Denom:           "adv4tnt",
						StartTime:       TEST_START_TIME.Add(time.Hour),
						EndTime:         TEST_END_TIME.Add(time.Hour),
					},
				},
			},
			genesisVestEntryKeys: []string{TEST_VESTER_ACCOUNT_1},
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
				// Set vest module genesis state with vest entries.
				testapp.UpdateGenesisDocWithAppStateForModule(
					&genesis,
					func(genesisState *vesttypes.GenesisState) {
						genesisState.VestEntries = make([]vesttypes.VestEntry, len(tc.genesisVestEntryKeys))
						for i, key := range tc.genesisVestEntryKeys {
							genesisState.VestEntries[i] = vesttypes.VestEntry{
								VesterAccount:   key,
								TreasuryAccount: TEST_GENESIS_VEST_ENTRY.TreasuryAccount,
								Denom:           TEST_GENESIS_VEST_ENTRY.Denom,
								StartTime:       TEST_GENESIS_VEST_ENTRY.StartTime,
								EndTime:         TEST_GENESIS_VEST_ENTRY.EndTime,
							}
						}
					},
				)
				return genesis
			}).Build()
			ctx := tApp.InitChain()

			// Submit and tally governance proposal that includes `MsgSetVestEntry`s.
			ctx = testapp.SubmitAndTallyProposal(
				t,
				ctx,
				tApp,
				tc.msgs,
				testapp.TestSubmitProposalTxHeight,
				false, // checkTx should not fail.
				false, // submitProposal should not fail.
				govtypesv1.ProposalStatus_PROPOSAL_STATUS_PASSED,
			)

			// Verify that expected vest entries are set in state.
			finalVestEntries := make(map[string]vesttypes.VestEntry)
			for _, msg := range tc.msgs {
				msgSetVestEntry, ok := msg.(*vesttypes.MsgSetVestEntry)
				require.True(t, ok)
				finalVestEntries[msgSetVestEntry.Entry.VesterAccount] = msgSetVestEntry.Entry
			}
			for vesterAccount, expectedVestEntry := range finalVestEntries {
				vestEntryInState, err := tApp.App.VestKeeper.GetVestEntry(ctx, vesterAccount)
				require.NoError(t, err)
				require.Equal(t, expectedVestEntry, vestEntryInState)
			}
		})
	}
}

func TestSetVestEntry_Failure(t *testing.T) {
	tests := map[string]struct {
		msgs                     []sdk.Msg
		expectCheckTxFails       bool
		expectSubmitProposalFail bool
	}{
		"Failure: vester account is empty": {
			msgs: []sdk.Msg{
				&vesttypes.MsgSetVestEntry{
					Authority: authtypes.NewModuleAddress(govtypes.ModuleName).String(),
					Entry: vesttypes.VestEntry{
						VesterAccount:   "",
						TreasuryAccount: "random_treasury",
						Denom:           "adv4tnt",
						StartTime:       TEST_START_TIME,
						EndTime:         TEST_END_TIME,
					},
				},
			},
			expectCheckTxFails: true,
		},
		"Failure: treasury account is empty": {
			msgs: []sdk.Msg{
				&vesttypes.MsgSetVestEntry{
					Authority: authtypes.NewModuleAddress(govtypes.ModuleName).String(),
					Entry: vesttypes.VestEntry{
						VesterAccount:   TEST_VESTER_ACCOUNT_1,
						TreasuryAccount: "",
						Denom:           "adv4tnt",
						StartTime:       TEST_START_TIME,
						EndTime:         TEST_END_TIME,
					},
				},
			},
			expectCheckTxFails: true,
		},
		"Failure: start time after end time": {
			msgs: []sdk.Msg{
				&vesttypes.MsgSetVestEntry{
					Authority: authtypes.NewModuleAddress(govtypes.ModuleName).String(),
					Entry: vesttypes.VestEntry{
						VesterAccount:   TEST_VESTER_ACCOUNT_1,
						TreasuryAccount: "",
						Denom:           "adv4tnt",
						StartTime:       TEST_END_TIME,
						EndTime:         TEST_START_TIME,
					},
				},
			},
			expectCheckTxFails: true,
		},
		"Failure: invalid authority": {
			msgs: []sdk.Msg{
				&vesttypes.MsgSetVestEntry{
					Authority: constants.BobAccAddress.String(),
					Entry: vesttypes.VestEntry{
						VesterAccount:   TEST_VESTER_ACCOUNT_1,
						TreasuryAccount: "random_treasury",
						Denom:           "adv4tnt",
						StartTime:       TEST_START_TIME,
						EndTime:         TEST_END_TIME,
					},
				},
			},
			expectSubmitProposalFail: true,
		},
		"Failure: failure of one message causes rollback of others": {
			msgs: []sdk.Msg{
				&vesttypes.MsgSetVestEntry{ // Valid message.
					Authority: authtypes.NewModuleAddress(govtypes.ModuleName).String(),
					Entry: vesttypes.VestEntry{
						VesterAccount:   TEST_VESTER_ACCOUNT_1,
						TreasuryAccount: "random_treasury",
						Denom:           "adv4tnt",
						StartTime:       TEST_START_TIME,
						EndTime:         TEST_END_TIME,
					},
				},
				&vesttypes.MsgSetVestEntry{ // Invalid message (due to invalid authority).
					Authority: constants.BobAccAddress.String(),
					Entry: vesttypes.VestEntry{
						VesterAccount:   TEST_VESTER_ACCOUNT_1,
						TreasuryAccount: "random_treasury",
						Denom:           "adv4tnt",
						StartTime:       TEST_START_TIME,
						EndTime:         TEST_END_TIME,
					},
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
				return genesis
			}).Build()
			ctx := tApp.InitChain()
			initialVestEntries := tApp.App.VestKeeper.GetAllVestEntries(ctx)

			// Submit and tally governance proposal that includes `MsgSetVestEntry`s.
			ctx = testapp.SubmitAndTallyProposal(
				t,
				ctx,
				tApp,
				tc.msgs,
				testapp.TestSubmitProposalTxHeight,
				tc.expectCheckTxFails,
				tc.expectSubmitProposalFail,
				govtypesv1.ProposalStatus_PROPOSAL_STATUS_FAILED,
			)

			// Verify that vest entries in state match the ones before proposal submission.
			require.Equal(t, initialVestEntries, tApp.App.VestKeeper.GetAllVestEntries(ctx))
		})
	}
}

func TestDeleteVestEntry_Success(t *testing.T) {
	tests := map[string]struct {
		msgs                 []sdk.Msg
		genesisVestEntryKeys []string // keys of vest entries in genesis state.
	}{
		"Success: delete one vest entry": {
			msgs: []sdk.Msg{
				&vesttypes.MsgDeleteVestEntry{
					Authority:     authtypes.NewModuleAddress(govtypes.ModuleName).String(),
					VesterAccount: TEST_VESTER_ACCOUNT_1,
				},
			},
			genesisVestEntryKeys: []string{TEST_VESTER_ACCOUNT_1},
		},
		"Success: delete two vest entries": {
			msgs: []sdk.Msg{
				&vesttypes.MsgDeleteVestEntry{
					Authority:     authtypes.NewModuleAddress(govtypes.ModuleName).String(),
					VesterAccount: TEST_VESTER_ACCOUNT_1,
				},
				&vesttypes.MsgDeleteVestEntry{
					Authority:     authtypes.NewModuleAddress(govtypes.ModuleName).String(),
					VesterAccount: TEST_VESTER_ACCOUNT_2,
				},
			},
			genesisVestEntryKeys: []string{TEST_VESTER_ACCOUNT_1, TEST_VESTER_ACCOUNT_2},
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
				// Set vest module genesis state with vest entries.
				testapp.UpdateGenesisDocWithAppStateForModule(
					&genesis,
					func(genesisState *vesttypes.GenesisState) {
						genesisState.VestEntries = make([]vesttypes.VestEntry, len(tc.genesisVestEntryKeys))
						for i, key := range tc.genesisVestEntryKeys {
							genesisState.VestEntries[i] = vesttypes.VestEntry{
								VesterAccount:   key,
								TreasuryAccount: TEST_GENESIS_VEST_ENTRY.TreasuryAccount,
								Denom:           TEST_GENESIS_VEST_ENTRY.Denom,
								StartTime:       TEST_GENESIS_VEST_ENTRY.StartTime,
								EndTime:         TEST_GENESIS_VEST_ENTRY.EndTime,
							}
						}
					},
				)
				return genesis
			}).Build()
			ctx := tApp.InitChain()

			// Submit and tally governance proposal that includes `MsgDeleteVestEntry`s.
			ctx = testapp.SubmitAndTallyProposal(
				t,
				ctx,
				tApp,
				tc.msgs,
				testapp.TestSubmitProposalTxHeight,
				false, // checkTx should not fail.
				false, // submitProposal should not fail.
				govtypesv1.ProposalStatus_PROPOSAL_STATUS_PASSED,
			)

			// Verify that vest entries have been deleted.
			for _, msg := range tc.msgs {
				vestEntry, err := tApp.App.VestKeeper.GetVestEntry(ctx, msg.(*vesttypes.MsgDeleteVestEntry).VesterAccount)
				require.Equal(t, vesttypes.VestEntry{}, vestEntry)
				require.Error(t, err)
			}
		})
	}
}

func TestDeleteVestEntry_Failure(t *testing.T) {
	tests := map[string]struct {
		msgs                     []sdk.Msg
		genesisVestEntryKeys     []string // key of vest entries in genesis state.
		expectSubmitProposalFail bool
	}{
		"Failure: vest entry does not exist": {
			msgs: []sdk.Msg{
				&vesttypes.MsgDeleteVestEntry{
					Authority:     authtypes.NewModuleAddress(govtypes.ModuleName).String(),
					VesterAccount: TEST_VESTER_ACCOUNT_1,
				},
			},
		},
		"Failure: delete the same vest entry twice": {
			msgs: []sdk.Msg{
				&vesttypes.MsgDeleteVestEntry{
					Authority:     authtypes.NewModuleAddress(govtypes.ModuleName).String(),
					VesterAccount: TEST_VESTER_ACCOUNT_1,
				},
				&vesttypes.MsgDeleteVestEntry{
					Authority:     authtypes.NewModuleAddress(govtypes.ModuleName).String(),
					VesterAccount: TEST_VESTER_ACCOUNT_1,
				},
			},
			genesisVestEntryKeys: []string{TEST_VESTER_ACCOUNT_1},
		},
		"Failure: second vest entry to delete does not exist": {
			msgs: []sdk.Msg{
				&vesttypes.MsgDeleteVestEntry{
					Authority:     authtypes.NewModuleAddress(govtypes.ModuleName).String(),
					VesterAccount: TEST_VESTER_ACCOUNT_1,
				},
				&vesttypes.MsgDeleteVestEntry{
					Authority:     authtypes.NewModuleAddress(govtypes.ModuleName).String(),
					VesterAccount: TEST_VESTER_ACCOUNT_2,
				},
			},
			genesisVestEntryKeys: []string{TEST_VESTER_ACCOUNT_1},
		},
		"Failure: invalid authority": {
			msgs: []sdk.Msg{
				&vesttypes.MsgDeleteVestEntry{
					Authority:     constants.BobAccAddress.String(),
					VesterAccount: TEST_VESTER_ACCOUNT_1,
				},
			},
			genesisVestEntryKeys:     []string{TEST_VESTER_ACCOUNT_1},
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
				// Set vest module genesis state with vest entries.
				testapp.UpdateGenesisDocWithAppStateForModule(
					&genesis,
					func(genesisState *vesttypes.GenesisState) {
						genesisState.VestEntries = make([]vesttypes.VestEntry, len(tc.genesisVestEntryKeys))
						for i, key := range tc.genesisVestEntryKeys {
							genesisState.VestEntries[i] = vesttypes.VestEntry{
								VesterAccount:   key,
								TreasuryAccount: TEST_GENESIS_VEST_ENTRY.TreasuryAccount,
								Denom:           TEST_GENESIS_VEST_ENTRY.Denom,
								StartTime:       TEST_GENESIS_VEST_ENTRY.StartTime,
								EndTime:         TEST_GENESIS_VEST_ENTRY.EndTime,
							}
						}
					},
				)
				return genesis
			}).Build()
			ctx := tApp.InitChain()
			initialVestEntries := tApp.App.VestKeeper.GetAllVestEntries(ctx)

			// Submit and tally governance proposal that includes `MsgDeleteVestEntry`s.
			ctx = testapp.SubmitAndTallyProposal(
				t,
				ctx,
				tApp,
				tc.msgs,
				testapp.TestSubmitProposalTxHeight,
				false, // checkTx should not fail.
				tc.expectSubmitProposalFail,
				govtypesv1.ProposalStatus_PROPOSAL_STATUS_FAILED,
			)

			// Verify that vest entries in state match vest entries before proposal submission.
			require.Equal(t, initialVestEntries, tApp.App.VestKeeper.GetAllVestEntries(ctx))
		})
	}
}
