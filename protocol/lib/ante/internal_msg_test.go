package ante_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"

	appmsgs "github.com/dydxprotocol/v4-chain/protocol/app/msgs"
	"github.com/dydxprotocol/v4-chain/protocol/lib"
	"github.com/dydxprotocol/v4-chain/protocol/lib/ante"
	testmsgs "github.com/dydxprotocol/v4-chain/protocol/testutil/msgs"
	"github.com/stretchr/testify/require"
)

func TestIsInternalMsg_Empty(t *testing.T) {
	tests := map[string]struct {
		msg sdk.Msg
	}{
		"empty msg": {
			msg: nil,
		},
		"not internal msg": {
			msg: testMsg,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			require.False(t, ante.IsInternalMsg(tc.msg))
		})
	}
}

func TestIsInternalMsg_Invalid(t *testing.T) {
	allMsgsMinusInternal := lib.MergeAllMapsMustHaveDistinctKeys(appmsgs.AllowMsgs, appmsgs.DisallowMsgs)
	for key := range appmsgs.InternalMsgSamplesAll {
		delete(allMsgsMinusInternal, key)
	}
	allNonNilSampleMsgs := testmsgs.GetNonNilSampleMsgs(allMsgsMinusInternal)

	for _, sampleMsg := range allNonNilSampleMsgs {
		t.Run(sampleMsg.Name, func(t *testing.T) {
			require.False(t, ante.IsInternalMsg(sampleMsg.Msg))
		})
	}
}

func TestIsInternalMsg_Valid(t *testing.T) {
	sampleMsgs := testmsgs.GetNonNilSampleMsgs(appmsgs.InternalMsgSamplesAll)
	// +1 for "/cosmos.auth.v1beta1.MsgUpdateParams" not having a corresponding Response msg type.
	require.Len(t, sampleMsgs, len(appmsgs.InternalMsgSamplesAll)/2+1)
	for _, sampleMsg := range sampleMsgs {
		t.Run(sampleMsg.Name, func(t *testing.T) {
			require.True(t, ante.IsInternalMsg(sampleMsg.Msg))
		})
	}
}
