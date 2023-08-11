package msgs_test

import (
	"testing"

	abcitypes "github.com/cometbft/cometbft/abci/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	testapp "github.com/dydxprotocol/v4/testutil/app"
	"github.com/dydxprotocol/v4/testutil/constants"
	testmsgs "github.com/dydxprotocol/v4/testutil/msgs"
	"github.com/stretchr/testify/require"
)

func TestDisallowMsgs_CheckTx_Fail(t *testing.T) {
	tests := map[string]struct {
		msg       sdk.Msg
		signTx    bool
		multiMsgs bool

		expectedErrorLog string
	}{
		// app-injected
		"app-injected, signed=false, multimsgs=false": {
			msg:       constants.ValidMsgUpdateMarketPrices,
			signTx:    false,
			multiMsgs: false,

			expectedErrorLog: "app-injected msg must only be included in DeliverTx: invalid request",
		},
		"app-injected, signed=true,  multimsgs=false": {
			msg:       constants.ValidMsgUpdateMarketPrices,
			signTx:    true,
			multiMsgs: false,

			expectedErrorLog: "app-injected msg must only be included in DeliverTx: invalid request",
		},
		"app-injected, signed=false, multimsgs=true": {
			msg:       constants.ValidMsgUpdateMarketPrices,
			signTx:    false,
			multiMsgs: true,

			expectedErrorLog: "app-injected msg must be the only msg in a tx: invalid request",
		},
		"app-injected, signed=true,  multimsgs=true": {
			msg:       constants.ValidMsgUpdateMarketPrices,
			signTx:    true,
			multiMsgs: true,

			expectedErrorLog: "app-injected msg must be the only msg in a tx: invalid request",
		},

		// internal
		"internal, signed=false, multimsgs=false": {
			msg:       testmsgs.MsgSoftwareUpgrade,
			signTx:    false,
			multiMsgs: false,

			expectedErrorLog: "internal msg cannot be submitted externally: invalid request",
		},
		"internal, signed=true,  multimsgs=false": {
			msg:       testmsgs.MsgSoftwareUpgrade,
			signTx:    true,
			multiMsgs: false,

			expectedErrorLog: "internal msg cannot be submitted externally: invalid request",
		},
		"internal, signed=false, multimsgs=true": {
			msg:       testmsgs.MsgSoftwareUpgrade,
			signTx:    false,
			multiMsgs: true,

			expectedErrorLog: "internal msg cannot be submitted externally: invalid request",
		},
		"internal, signed=true,  multimsgs=true": {
			msg:       testmsgs.MsgSoftwareUpgrade,
			signTx:    true,
			multiMsgs: true,

			expectedErrorLog: "internal msg cannot be submitted externally: invalid request",
		},

		// nested
		"nested, signed=false, multimsgs=false": {
			msg:       testmsgs.MsgSubmitProposalWithAppInjectedInner,
			signTx:    false,
			multiMsgs: false,

			expectedErrorLog: "Invalid nested msg: app-injected msg type: invalid request",
		},
		"nested, signed=true,  multimsgs=false": {
			msg:       testmsgs.MsgSubmitProposalWithAppInjectedInner,
			signTx:    false,
			multiMsgs: true,

			expectedErrorLog: "Invalid nested msg: app-injected msg type: invalid request",
		},
		"nested, signed=false, multimsgs=true": {
			msg:       testmsgs.MsgSubmitProposalWithAppInjectedInner,
			signTx:    true,
			multiMsgs: false,

			expectedErrorLog: "Invalid nested msg: app-injected msg type: invalid request",
		},
		"nested, signed=true,  multimsgs=true": {
			msg:       testmsgs.MsgSubmitProposalWithAppInjectedInner,
			signTx:    true,
			multiMsgs: true,

			expectedErrorLog: "Invalid nested msg: app-injected msg type: invalid request",
		},

		// unsupported
		"unsupported, signed=false, multimsgs=false": {
			msg:       testmsgs.GovBetaMsgSubmitProposal,
			signTx:    false,
			multiMsgs: false,

			expectedErrorLog: "unsupported msg: invalid request",
		},
		"unsupported, signed=true,  multimsgs=false": {
			msg:       testmsgs.GovBetaMsgSubmitProposal,
			signTx:    true,
			multiMsgs: false,

			expectedErrorLog: "unsupported msg: invalid request",
		},
		"unsupported, signed=false, multimsgs=true": {
			msg:       testmsgs.GovBetaMsgSubmitProposal,
			signTx:    false,
			multiMsgs: true,

			expectedErrorLog: "unsupported msg: invalid request",
		},
		"unsupported, signed=true,  multimsgs=true": {
			msg:       testmsgs.GovBetaMsgSubmitProposal,
			signTx:    true,
			multiMsgs: true,

			expectedErrorLog: "unsupported msg: invalid request",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// Setup app.
			tApp := testapp.NewTestAppBuilder().WithTesting(t).Build()
			// Note that due to TODO(DEC-1248) the minimum block height is 2.
			ctx := tApp.AdvanceToBlock(2)

			// Setup msgs.
			msgs := make([]sdk.Msg, 0)
			msgs = append(msgs, tc.msg)
			if tc.multiMsgs { // append extra msg to the tx
				msgs = append(msgs, constants.Msg_Send)
			}

			// Setup reqCheckTx.
			var reqCheckTx abcitypes.RequestCheckTx
			if tc.signTx {
				reqCheckTx = testapp.MustMakeCheckTx(
					ctx,
					tApp.App,
					testapp.MustMakeCheckTxOptions{
						AccAddressForSigning: string(constants.Alice_Num0.Owner),
					},
					msgs...,
				)
			} else { // simply encoded the tx without signing the msgs.
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
			result := tApp.CheckTx(reqCheckTx)
			require.False(t, result.IsOK(), "expected CheckTx to fail")
			require.Equal(t, sdkerrors.ErrInvalidRequest.ABCICode(), result.Code)
			require.Equal(t, tc.expectedErrorLog, result.Log)
		})
	}
}
