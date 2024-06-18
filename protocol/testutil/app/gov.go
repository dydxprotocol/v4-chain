package app

import (
	"bytes"
	"testing"
	"time"

	abcitypes "github.com/cometbft/cometbft/abci/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	govtypesv1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	"github.com/stretchr/testify/require"
)

const (
	TestMetadata               = "test metadata"
	TestTitle                  = "test title"
	TestSummary                = "test summary"
	TestSubmitProposalTxHeight = uint32(2)
	TestProposalTallyHeight    = uint32(10)
)

var (
	TestVotingPeriod = 1 * time.Minute
	TestDeposit      = sdk.Coins{sdk.NewInt64Coin(constants.TestNativeTokenDenom, 10_000_000)}
)

// SubmitAndTallyProposal simulates the following:
//   - A proposal with the given messages is submitted. Check proposal submission succeeds or fails as expected.
//   - If proposal successfully submitted:
//     -- All validators vote for the proposal.
//     -- The proposal is tallied after voting period ends.
func SubmitAndTallyProposal(
	t testing.TB,
	ctx sdk.Context,
	tApp *TestApp,
	messages []sdk.Msg,
	submitProposalTxHeight uint32,
	expectCheckTxFails bool,
	expectSubmitProposalFails bool,
	expectedProposalStatus govtypesv1.ProposalStatus,
) sdk.Context {
	// Create a MsgSubmitProposal
	msgSubmitProposal, err := govtypesv1.NewMsgSubmitProposal(
		messages,
		TestDeposit,
		constants.AliceAccAddress.String(),
		TestMetadata,
		TestTitle,
		TestSummary,
		false,
	)
	require.NoError(t, err)

	// Create a signed transaction with MsgSubmitProposal
	submitProposalCheckTx := MustMakeCheckTxWithPrivKeySupplier(
		ctx,
		tApp.App,
		MustMakeCheckTxOptions{
			AccAddressForSigning: constants.AliceAccAddress.String(),
			Gas:                  constants.TestGasLimit,
			FeeAmt:               constants.TestFeeCoins_5Cents,
		},
		constants.GetPrivateKeyFromAddress,
		msgSubmitProposal,
	)
	result := tApp.CheckTx(submitProposalCheckTx)
	if expectCheckTxFails {
		require.False(t, result.IsOK())
		// CheckTx failed. Return early.
		return ctx
	} else {
		require.True(t, result.IsOK())
	}

	if expectSubmitProposalFails {
		ctx = tApp.AdvanceToBlock(submitProposalTxHeight, AdvanceToBlockOptions{
			ValidateFinalizeBlock: func(
				context sdk.Context,
				request abcitypes.RequestFinalizeBlock,
				response abcitypes.ResponseFinalizeBlock,
			) (haltChain bool) {
				for i := range request.Txs {
					if bytes.Equal(request.Txs[i], submitProposalCheckTx.Tx) {
						require.True(t, response.TxResults[i].IsErr())
					} else {
						require.True(t, response.TxResults[i].IsOK())
					}
				}
				return false
			},
		})
		// Proposal submission failed. Return early.
		return ctx
	} else {
		ctx = tApp.AdvanceToBlock(submitProposalTxHeight, AdvanceToBlockOptions{})
	}

	proposalsIterator, err := tApp.App.GovKeeper.Proposals.Iterate(ctx, nil)
	require.NoError(t, err)
	proposals, err := proposalsIterator.Values()
	require.NoError(t, err)
	require.Len(t, proposals, 1)

	// Have all 4 validators vote for the proposal.
	for _, validator := range []sdk.AccAddress{
		constants.AliceAccAddress,
		constants.BobAccAddress,
		constants.CarlAccAddress,
		constants.DaveAccAddress,
	} {
		// Create a signed vote transaction
		msgVote := govtypesv1.NewMsgVote(
			validator,
			1, // proposal ID
			govtypesv1.VoteOption_VOTE_OPTION_YES,
			"", // metadata
		)
		voteCheckTx := MustMakeCheckTxWithPrivKeySupplier(
			ctx,
			tApp.App,
			MustMakeCheckTxOptions{
				AccAddressForSigning: validator.String(),
				Gas:                  constants.TestGasLimit,
				FeeAmt:               constants.TestFeeCoins_5Cents,
			},
			constants.GetPrivateKeyFromAddress,
			msgVote,
		)
		require.True(t, tApp.CheckTx(voteCheckTx).IsOK())
	}

	// Advance to the height right before voting period ends.
	ctx = tApp.AdvanceToBlock(TestProposalTallyHeight-1, AdvanceToBlockOptions{})
	// Advance height to TestProposalTallyHeight; advance block time to right after voting period.
	ctx = tApp.AdvanceToBlock(TestProposalTallyHeight, AdvanceToBlockOptions{
		// Ensure the voting period has passed.
		BlockTime: ctx.BlockTime().Add(TestVotingPeriod).Add(1 * time.Second),
	})

	proposalsIterator, err = tApp.App.GovKeeper.Proposals.Iterate(ctx, nil)
	require.NoError(t, err)
	proposals, err = proposalsIterator.Values()
	require.NoError(t, err)
	require.Len(t, proposals, 1)
	// Check that proposal was executed or failed as expected.
	require.Equal(t, expectedProposalStatus, proposals[0].Status)
	return ctx
}
