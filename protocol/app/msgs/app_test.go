package msgs_test

import (
	"testing"

	abcitypes "github.com/cometbft/cometbft/abci/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	testapp "github.com/dydxprotocol/v4-chain/protocol/testutil/app"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	testmsgs "github.com/dydxprotocol/v4-chain/protocol/testutil/msgs"
	"github.com/stretchr/testify/require"
)

type testCase struct {
	name      string
	msg       sdk.Msg
	multiMsgs bool
}

var (
	appInjectedSingle = testCase{
		name:      "app-injected, multimsgs=false",
		msg:       constants.ValidMsgUpdateMarketPrices,
		multiMsgs: false,
	}
	appInjectedMulti = testCase{
		name:      "app-injected, multimsgs=true",
		msg:       constants.ValidMsgUpdateMarketPrices,
		multiMsgs: true,
	}
	internalSingle = testCase{
		name:      "internal, multimsgs=false",
		msg:       testmsgs.MsgSoftwareUpgrade,
		multiMsgs: false,
	}
	internalMulti = testCase{
		name:      "internal, multimsgs=true",
		msg:       testmsgs.MsgSoftwareUpgrade,
		multiMsgs: true,
	}
	nestedSingle = testCase{
		name:      "nested, multimsgs=false",
		msg:       testmsgs.MsgSubmitProposalWithAppInjectedInner,
		multiMsgs: false,
	}
	nestedMulti = testCase{
		name:      "nested, multimsgs=true",
		msg:       testmsgs.MsgSubmitProposalWithAppInjectedInner,
		multiMsgs: true,
	}
	unsupportedSingle = testCase{
		name:      "unsupported, multimsgs=false",
		msg:       testmsgs.GovBetaMsgSubmitProposal,
		multiMsgs: false,
	}
	unsupportedMulti = testCase{
		name:      "unsupported, multimsgs=true",
		msg:       testmsgs.GovBetaMsgSubmitProposal,
		multiMsgs: true,
	}
)

func TestDisallowMsgs_CheckTx_Fail(t *testing.T) {
	tests := []struct {
		name   string
		base   testCase
		signTx bool

		expectedErrorLog string
	}{
		// app-injected
		{
			name:   appInjectedSingle.name + ", signed=false",
			base:   appInjectedSingle,
			signTx: false,

			expectedErrorLog: "app-injected msg must only be included in DeliverTx: invalid request",
		},
		{
			name:   appInjectedSingle.name + ", signed=true",
			base:   appInjectedSingle,
			signTx: true,

			expectedErrorLog: "app-injected msg must only be included in DeliverTx: invalid request",
		},
		{
			name:   appInjectedMulti.name + ", signed=false",
			base:   appInjectedMulti,
			signTx: false,

			expectedErrorLog: "app-injected msg must be the only msg in a tx: invalid request",
		},
		{
			name:   appInjectedMulti.name + ", signed=true",
			base:   appInjectedMulti,
			signTx: true,

			expectedErrorLog: "app-injected msg must be the only msg in a tx: invalid request",
		},

		// internal
		{
			name:   internalSingle.name + ", signed=false",
			base:   internalSingle,
			signTx: false,

			expectedErrorLog: "internal msg cannot be submitted externally: invalid request",
		},
		{
			name:   internalSingle.name + ", signed=true",
			base:   internalSingle,
			signTx: true,

			expectedErrorLog: "internal msg cannot be submitted externally: invalid request",
		},
		{
			name:   internalMulti.name + ", signed=false",
			base:   internalMulti,
			signTx: false,

			expectedErrorLog: "internal msg cannot be submitted externally: invalid request",
		},
		{
			name:   internalMulti.name + ", signed=true",
			base:   internalMulti,
			signTx: true,

			expectedErrorLog: "internal msg cannot be submitted externally: invalid request",
		},

		// nested
		{
			name:   nestedSingle.name + ", signed=false",
			base:   nestedSingle,
			signTx: false,

			expectedErrorLog: "Invalid nested msg: app-injected msg type: invalid request",
		},
		{
			name:   nestedSingle.name + ", signed=true",
			base:   nestedSingle,
			signTx: true,

			expectedErrorLog: "Invalid nested msg: app-injected msg type: invalid request",
		},
		{
			name:   nestedMulti.name + ", signed=false",
			base:   nestedMulti,
			signTx: false,

			expectedErrorLog: "Invalid nested msg: app-injected msg type: invalid request",
		},
		{
			name:   nestedMulti.name + ", signed=true",
			base:   nestedMulti,
			signTx: true,

			expectedErrorLog: "Invalid nested msg: app-injected msg type: invalid request",
		},

		// unsupported
		{
			name:   unsupportedSingle.name + ", signed=false",
			base:   unsupportedSingle,
			signTx: false,

			expectedErrorLog: "unsupported msg: invalid request",
		},
		{
			name:   unsupportedSingle.name + ", signed=true",
			base:   unsupportedSingle,
			signTx: true,

			expectedErrorLog: "unsupported msg: invalid request",
		},
		{
			name:   unsupportedMulti.name + ", signed=false",
			base:   unsupportedMulti,
			signTx: false,

			expectedErrorLog: "unsupported msg: invalid request",
		},
		{
			name:   unsupportedMulti.name + ", signed=true",
			base:   unsupportedMulti,
			signTx: true,

			expectedErrorLog: "unsupported msg: invalid request",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Setup app.
			tApp := testapp.NewTestAppBuilder(t).Build()
			ctx := tApp.AdvanceToBlock(2, testapp.AdvanceToBlockOptions{})

			// Setup msgs.
			msgs := getMsgs(tc.base)

			// Setup reqCheckTx.
			var reqCheckTx abcitypes.RequestCheckTx
			if tc.signTx {
				reqCheckTx = testapp.MustMakeCheckTx(
					ctx,
					tApp.App,
					testapp.MustMakeCheckTxOptions{
						AccAddressForSigning: constants.Alice_Num0.Owner,
					},
					msgs...,
				)
			} else { // simply encode the tx without signing the msgs.
				txBuilder := tApp.App.TxConfig().NewTxBuilder()
				err := txBuilder.SetMsgs(msgs...)
				require.NoError(t, err)
				bytes, err := tApp.App.TxConfig().TxEncoder()(txBuilder.GetTx())
				require.NoError(t, err)

				reqCheckTx = abcitypes.RequestCheckTx{
					Tx:   bytes,
					Type: abcitypes.CheckTxType_New,
				}
			}

			// Run & Validate.
			resp := tApp.CheckTx(reqCheckTx)
			require.Conditionf(t, resp.IsErr, "Expected CheckTx to error. Response: %+v", resp)
			require.Equal(t, sdkerrors.ErrInvalidRequest.ABCICode(), resp.Code)
			require.Equal(t, tc.expectedErrorLog, resp.Log)
		})
	}
}

func TestDisallowMsgs_PrepareProposal_Filter(t *testing.T) {
	tests := []testCase{
		// app-injected
		appInjectedSingle,
		appInjectedMulti,

		// internal
		internalSingle,
		internalMulti,

		// nested
		nestedSingle,
		nestedMulti,

		// unsupported
		unsupportedSingle,
		unsupportedMulti,
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Setup app.
			tApp := testapp.NewTestAppBuilder(t).Build()
			tApp.AdvanceToBlock(2, testapp.AdvanceToBlockOptions{})

			// Setup msg tx.
			msgs := getMsgs(tc)
			txBuilder := tApp.App.TxConfig().NewTxBuilder()
			err := txBuilder.SetMsgs(msgs...)
			require.NoError(t, err)
			otherTxsBytes, err := tApp.App.TxConfig().TxEncoder()(txBuilder.GetTx())
			require.NoError(t, err)

			// Test that PrepareProposal filters out the disallow msgs and ProcessProposal accepts the result.
			tApp.AdvanceToBlock(
				3,
				testapp.AdvanceToBlockOptions{
					// 1. Override the PrepareProposal txs with test case ones.
					RequestPrepareProposalTxsOverride: [][]byte{otherTxsBytes},

					// 2. Validate that PrepareProposal would filter out the disallow msgs.
					ValidateRespPrepare: func(ctx sdk.Context, resp abcitypes.ResponsePrepareProposal) (haltChain bool) {
						proposalTxs := resp.GetTxs()
						require.Len(t, proposalTxs, 4)
						require.Equal(t, constants.ValidEmptyMsgProposedOperationsTxBytes, proposalTxs[0])
						require.Equal(t, constants.MsgAcknowledgeBridges_NoEvents_TxBytes, proposalTxs[1])
						require.Equal(t, constants.EmptyMsgAddPremiumVotesTxBytes, proposalTxs[2])
						require.Equal(t, constants.EmptyMsgUpdateMarketPricesTxBytes, proposalTxs[3])
						return false
					},

					// 3. Validate that the filtered PrepareProposal txs are accepted during ProcessProposal.
					ValidateRespProcess: func(ctx sdk.Context, resp abcitypes.ResponseProcessProposal) (haltChain bool) {
						require.Equal(t, abcitypes.ResponseProcessProposal_ACCEPT, resp.Status)
						return false
					},
				},
			)
		})
	}
}

func TestDisallowMsgs_ProcessProposal_Fail(t *testing.T) {
	tests := []testCase{
		// app-injected
		appInjectedSingle,
		appInjectedMulti,

		// internal
		internalSingle,
		internalMulti,

		// nested
		nestedSingle,
		nestedMulti,

		// unsupported
		unsupportedSingle,
		unsupportedMulti,
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Setup app.
			tApp := testapp.NewTestAppBuilder(t).Build()
			tApp.AdvanceToBlock(2, testapp.AdvanceToBlockOptions{})

			// Setup msgs.
			msgs := getMsgs(tc)
			txBuilder := tApp.App.TxConfig().NewTxBuilder()
			err := txBuilder.SetMsgs(msgs...)
			require.NoError(t, err)
			otherTxsBytes, err := tApp.App.TxConfig().TxEncoder()(txBuilder.GetTx())
			require.NoError(t, err)

			tApp.AdvanceToBlock(
				3,
				testapp.AdvanceToBlockOptions{
					// 1. Override the ProcessProposal txs with test case ones.
					RequestProcessProposalTxsOverride: getProposalTxsWithOtherTxs(otherTxsBytes),

					// 2. Test that ProcessProposal rejects when disallow msgs are in `OtherTxs`.
					ValidateRespProcess: func(ctx sdk.Context, resp abcitypes.ResponseProcessProposal) (haltChain bool) {
						require.Equal(t, abcitypes.ResponseProcessProposal_REJECT, resp.Status)
						return true // halt chain.
					},
				},
			)
		})
	}
}

func getMsgs(tc testCase) []sdk.Msg {
	if tc.multiMsgs { // append extra msg to the tx
		return []sdk.Msg{tc.msg, constants.Msg_Send}
	}
	return []sdk.Msg{tc.msg}
}

func getProposalTxsWithOtherTxs(otherTxsToAppend []byte) [][]byte {
	if otherTxsToAppend == nil {
		panic("otherTxsToAppend cannot be nil")
	}
	return [][]byte{
		constants.ValidEmptyMsgProposedOperationsTxBytes,
		otherTxsToAppend,
		constants.EmptyMsgAddPremiumVotesTxBytes,
		constants.EmptyMsgUpdateMarketPricesTxBytes,
	}
}
