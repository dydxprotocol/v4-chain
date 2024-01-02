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

func TestIsUnsupportedMsg_Empty(t *testing.T) {
	tests := map[string]struct {
		msg sdk.Msg
	}{
		"empty msg": {
			msg: nil,
		},
		"not unsupported msg": {
			msg: testMsg,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			require.False(t, ante.IsUnsupportedMsg(tc.msg))
		})
	}
}

func TestIsUnsupportedMsg_Invalid(t *testing.T) {
	allMsgsMinusUnsupported := lib.MergeAllMapsMustHaveDistinctKeys(appmsgs.AllowMsgs, appmsgs.DisallowMsgs)
	for key := range appmsgs.UnsupportedMsgSamples {
		delete(allMsgsMinusUnsupported, key)
	}
	allNonNilSampleMsgs := testmsgs.GetNonNilSampleMsgs(allMsgsMinusUnsupported)

	for _, sampleMsg := range allNonNilSampleMsgs {
		t.Run(sampleMsg.Name, func(t *testing.T) {
			require.False(t, ante.IsUnsupportedMsg(sampleMsg.Msg))
		})
	}
}

func TestIsUnsupportedMsg_Valid(t *testing.T) {
	sampleMsgs := testmsgs.GetNonNilSampleMsgs(appmsgs.UnsupportedMsgSamples)

	for _, sampleMsg := range sampleMsgs {
		t.Run(sampleMsg.Name, func(t *testing.T) {
			require.True(t, ante.IsUnsupportedMsg(sampleMsg.Msg))
		})
	}
}
