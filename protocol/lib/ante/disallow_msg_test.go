package ante_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"

	appmsgs "github.com/dydxprotocol/v4-chain/protocol/app/msgs"
	"github.com/dydxprotocol/v4-chain/protocol/lib/ante"
	testmsgs "github.com/dydxprotocol/v4-chain/protocol/testutil/msgs"
	"github.com/stretchr/testify/require"
)

func TestIsDisallowExternalSubmitMsg(t *testing.T) {
	// All disallow msgs should return true.
	for _, sampleMsg := range testmsgs.GetNonNilSampleMsgs(appmsgs.DisallowMsgs) {
		result := ante.IsDisallowExternalSubmitMsg(sampleMsg.Msg)
		if ante.IsNestedMsg(sampleMsg.Msg) {
			// nested msgs are allowed as long as the inner msgs are allowed.
			require.False(t, result, sampleMsg.Name)
		} else {
			require.True(t, result, sampleMsg.Name)
		}
	}

	// All allow msgs should return false.
	for _, sampleMsg := range testmsgs.GetNonNilSampleMsgs(appmsgs.AllowMsgs) {
		require.False(t, ante.IsDisallowExternalSubmitMsg(sampleMsg.Msg), sampleMsg.Name)
	}
}

func TestIsDisallowExternalSubmitMsg_InvalidInnerMsgs(t *testing.T) {
	containsInvalidInnerMsgs := []sdk.Msg{
		testmsgs.MsgSubmitProposalWithUnsupportedInner,
		testmsgs.MsgSubmitProposalWithAppInjectedInner,
		testmsgs.MsgSubmitProposalWithDoubleNestedInner,
	}

	for _, msg := range containsInvalidInnerMsgs {
		require.True(t, ante.IsDisallowExternalSubmitMsg(msg))
	}
}
