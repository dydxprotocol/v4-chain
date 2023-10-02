package gov_test

import (
	"testing"

	"github.com/cometbft/cometbft/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	govtypesv1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
	testapp "github.com/dydxprotocol/v4-chain/protocol/testutil/app"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	sendingtypes "github.com/dydxprotocol/v4-chain/protocol/x/sending/types"
	vesttypes "github.com/dydxprotocol/v4-chain/protocol/x/vest/types"
	"github.com/stretchr/testify/require"
)

func TestSendFromModuleToAccount(t *testing.T) {
	tests := map[string]struct {
		msg                      *sendingtypes.MsgSendFromModuleToAccount
		expectedProposalStatus   govtypesv1.ProposalStatus
		expectSubmitProposalFail bool
	}{
		"Success": {
			msg: &sendingtypes.MsgSendFromModuleToAccount{
				Authority:        authtypes.NewModuleAddress(govtypes.ModuleName).String(),
				SenderModuleName: vesttypes.CommunityTreasuryAccountName,
				Recipient:        constants.AliceAccAddress.String(),
				Coin:             sdk.NewCoin("dv4tnt", sdk.NewInt(123)),
			},
			expectedProposalStatus: govtypesv1.ProposalStatus_PROPOSAL_STATUS_PASSED,
		},
		"Fail: invalid authority": {
			msg: &sendingtypes.MsgSendFromModuleToAccount{
				Authority:        authtypes.NewModuleAddress(sendingtypes.ModuleName).String(),
				SenderModuleName: banktypes.ModuleName,
				Recipient:        constants.AliceAccAddress.String(),
				Coin:             sdk.NewCoin("dv4tnt", sdk.NewInt(123)),
			},
			expectSubmitProposalFail: true,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			senderModuleAddress := authtypes.NewModuleAddress(tc.msg.SenderModuleName)
			tApp := testapp.NewTestAppBuilder().WithGenesisDocFn(func() (genesis types.GenesisDoc) {
				genesis = testapp.DefaultGenesis()
				testapp.UpdateGenesisDocWithAppStateForModule(
					&genesis,
					func(genesisState *govtypesv1.GenesisState) {
						genesisState.Params.VotingPeriod = &testapp.TestVotingPeriod
					},
				)
				// Initialize sender module with enough balance to send.
				testapp.UpdateGenesisDocWithAppStateForModule(
					&genesis,
					func(genesisState *banktypes.GenesisState) {
						genesisState.Balances = append(genesisState.Balances, banktypes.Balance{
							Address: senderModuleAddress.String(),
							Coins: sdk.Coins{
								tc.msg.Coin.AddAmount(sdk.NewInt(567)),
							},
						})
					},
				)
				return genesis
			}).WithTesting(t).Build()
			ctx := tApp.InitChain()
			initialModuleBalance := tApp.App.BankKeeper.GetBalance(ctx, senderModuleAddress, tc.msg.Coin.Denom)
			initialRecipientBalance := tApp.App.BankKeeper.GetBalance(
				ctx,
				sdk.MustAccAddressFromBech32(tc.msg.Recipient),
				tc.msg.Coin.Denom,
			)

			// Submit and tally governance proposal that includes `MsgUpdateParams`.
			ctx = testapp.SubmitAndTallyProposal(
				t,
				ctx,
				&tApp,
				[]sdk.Msg{tc.msg},
				tc.expectSubmitProposalFail,
				tc.expectedProposalStatus,
			)

			updatedModuleBalance := tApp.App.BankKeeper.GetBalance(ctx, senderModuleAddress, tc.msg.Coin.Denom)
			updatedRecipientBalance := tApp.App.BankKeeper.GetBalance(
				ctx,
				sdk.MustAccAddressFromBech32(tc.msg.Recipient),
				tc.msg.Coin.Denom,
			)

			// If governance proposal is supposed to fail submission, verify that module and recipient
			// balances match the ones before proposal submission.
			if tc.expectSubmitProposalFail {
				require.Equal(t, initialModuleBalance, updatedModuleBalance)
				require.Equal(t, initialRecipientBalance, updatedRecipientBalance)
			}

			// If proposal is supposed to pass, verify that module and recipient balances are updated.
			if tc.expectedProposalStatus == govtypesv1.ProposalStatus_PROPOSAL_STATUS_PASSED {
				require.Equal(t, initialModuleBalance.Sub(tc.msg.Coin), updatedModuleBalance)
				require.Equal(t, initialRecipientBalance.Add(tc.msg.Coin), updatedRecipientBalance)
			}
		})
	}
}
