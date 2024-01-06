package gov_test

import (
	"testing"

	sdkmath "cosmossdk.io/math"

	"github.com/dydxprotocol/v4-chain/protocol/lib"

	"github.com/cometbft/cometbft/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
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
		initialModuleBalance     int64
		expectCheckTxFails       bool
		expectedProposalStatus   govtypesv1.ProposalStatus
		expectSubmitProposalFail bool
	}{
		"Success: send from module to user account": {
			msg: &sendingtypes.MsgSendFromModuleToAccount{
				Authority:        lib.GovModuleAddress.String(),
				SenderModuleName: vesttypes.CommunityTreasuryAccountName,
				Recipient:        constants.AliceAccAddress.String(),
				Coin:             sdk.NewCoin("adv4tnt", sdkmath.NewInt(123)),
			},
			initialModuleBalance:   200,
			expectedProposalStatus: govtypesv1.ProposalStatus_PROPOSAL_STATUS_PASSED,
		},
		"Success: send from module to module account": {
			msg: &sendingtypes.MsgSendFromModuleToAccount{
				Authority:        lib.GovModuleAddress.String(),
				SenderModuleName: vesttypes.CommunityTreasuryAccountName,
				Recipient:        authtypes.NewModuleAddress(vesttypes.CommunityVesterAccountName).String(),
				Coin:             sdk.NewCoin("adv4tnt", sdkmath.NewInt(123)),
			},
			initialModuleBalance:   123,
			expectedProposalStatus: govtypesv1.ProposalStatus_PROPOSAL_STATUS_PASSED,
		},
		"Failure: insufficient balance": {
			msg: &sendingtypes.MsgSendFromModuleToAccount{
				Authority:        lib.GovModuleAddress.String(),
				SenderModuleName: vesttypes.CommunityTreasuryAccountName,
				Recipient:        authtypes.NewModuleAddress(vesttypes.CommunityVesterAccountName).String(),
				Coin:             sdk.NewCoin("adv4tnt", sdkmath.NewInt(124)),
			},
			initialModuleBalance:   123,
			expectedProposalStatus: govtypesv1.ProposalStatus_PROPOSAL_STATUS_FAILED,
		},
		"Failure: invalid authority": {
			msg: &sendingtypes.MsgSendFromModuleToAccount{
				Authority:        authtypes.NewModuleAddress(sendingtypes.ModuleName).String(),
				SenderModuleName: vesttypes.CommunityTreasuryAccountName,
				Recipient:        constants.AliceAccAddress.String(),
				Coin:             sdk.NewCoin("adv4tnt", sdkmath.NewInt(123)),
			},
			initialModuleBalance:     123,
			expectSubmitProposalFail: true,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			senderModuleAddress := authtypes.NewModuleAddress(tc.msg.SenderModuleName)
			tApp := testapp.NewTestAppBuilder(t).WithGenesisDocFn(func() (genesis types.GenesisDoc) {
				genesis = testapp.DefaultGenesis()
				testapp.UpdateGenesisDocWithAppStateForModule(
					&genesis,
					func(genesisState *govtypesv1.GenesisState) {
						genesisState.Params.VotingPeriod = &testapp.TestVotingPeriod
					},
				)
				// Initialize sender module with its initial balance.
				testapp.UpdateGenesisDocWithAppStateForModule(
					&genesis,
					func(genesisState *banktypes.GenesisState) {
						genesisState.Balances = append(genesisState.Balances, banktypes.Balance{
							Address: senderModuleAddress.String(),
							Coins: sdk.Coins{
								sdk.NewCoin(tc.msg.Coin.Denom, sdkmath.NewInt(tc.initialModuleBalance)),
							},
						})
					},
				)
				return genesis
			}).Build()
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
				tApp,
				[]sdk.Msg{tc.msg},
				testapp.TestSubmitProposalTxHeight,
				tc.expectCheckTxFails,
				tc.expectSubmitProposalFail,
				tc.expectedProposalStatus,
			)

			updatedModuleBalance := tApp.App.BankKeeper.GetBalance(ctx, senderModuleAddress, tc.msg.Coin.Denom)
			updatedRecipientBalance := tApp.App.BankKeeper.GetBalance(
				ctx,
				sdk.MustAccAddressFromBech32(tc.msg.Recipient),
				tc.msg.Coin.Denom,
			)

			// If governance proposal is supposed to fail submission or execution, verify that module and
			// recipient balances match the ones before proposal submission.
			if tc.expectSubmitProposalFail || tc.expectedProposalStatus == govtypesv1.ProposalStatus_PROPOSAL_STATUS_FAILED {
				require.Equal(t, initialModuleBalance, updatedModuleBalance)
				require.Equal(t, initialRecipientBalance, updatedRecipientBalance)
			}

			// If proposal is supposed to pass, verify that module and recipient balances are updated.
			if tc.expectedProposalStatus == govtypesv1.ProposalStatus_PROPOSAL_STATUS_PASSED {
				require.True(t, updatedModuleBalance.Equal(initialModuleBalance.Sub(tc.msg.Coin)))
				require.True(t, updatedRecipientBalance.Equal(initialRecipientBalance.Add(tc.msg.Coin)))
			}
		})
	}
}
