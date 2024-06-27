package ante_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	appmsgs "github.com/StreamFinance-Protocol/stream-chain/protocol/app/msgs"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/lib"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/lib/ante"
	testmsgs "github.com/StreamFinance-Protocol/stream-chain/protocol/testutil/msgs"
	perptypes "github.com/StreamFinance-Protocol/stream-chain/protocol/x/perpetuals/types"
)

func TestIsSingleAppInjectedMsg(t *testing.T) {
	tests := map[string]struct {
		msgs           []sdk.Msg
		expectedResult bool
	}{
		"empty msgs": {
			expectedResult: false,
		},
		"single msg: no app-injected msg": {
			msgs:           []sdk.Msg{testMsg},
			expectedResult: false,
		},
		"single msg: app-injected msg": {
			msgs: []sdk.Msg{
				&perptypes.MsgAddPremiumVotes{}, // app-injected.
			},
			expectedResult: true,
		},
		"mult msg: no app-injected msgs": {
			msgs:           []sdk.Msg{testMsg, testMsg},
			expectedResult: false,
		},
		"mult msg: all app-injected msgs": {
			msgs: []sdk.Msg{
				&perptypes.MsgAddPremiumVotes{}, // app-injected.
			},
			expectedResult: false,
		},
		"mult msg: mixed": {
			msgs: []sdk.Msg{
				testMsg,
				&perptypes.MsgAddPremiumVotes{}, // app-injected.
			},
			expectedResult: false,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			require.Equal(t, tc.expectedResult, ante.IsSingleAppInjectedMsg(tc.msgs))
		})
	}
}

func TestIsAppInjectedMsg_Empty(t *testing.T) {
	tests := map[string]struct {
		msg sdk.Msg
	}{
		"empty msg": {
			msg: nil,
		},
		"not app-injected msg": {
			msg: testMsg,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			require.False(t, ante.IsAppInjectedMsg(tc.msg))
		})
	}
}

func TestIsAppInjectedMsg_Invalid(t *testing.T) {
	allMsgsMinusAppInjected := lib.MergeAllMapsMustHaveDistinctKeys(appmsgs.AllowMsgs, appmsgs.DisallowMsgs)
	for key := range appmsgs.AppInjectedMsgSamples {
		delete(allMsgsMinusAppInjected, key)
	}
	allNonNilSampleMsgs := testmsgs.GetNonNilSampleMsgs(allMsgsMinusAppInjected)

	for _, sampleMsg := range allNonNilSampleMsgs {
		t.Run(sampleMsg.Name, func(t *testing.T) {
			require.False(t, ante.IsAppInjectedMsg(sampleMsg.Msg))
		})
	}
}

func TestIsAppInjectedMsg_Valid(t *testing.T) {
	appInjectedSampleMsgs := testmsgs.GetNonNilSampleMsgs(appmsgs.AppInjectedMsgSamples)
	require.Len(t, appInjectedSampleMsgs, len(appmsgs.AppInjectedMsgSamples)/2)
	for _, sampleMsg := range appInjectedSampleMsgs {
		t.Run(sampleMsg.Name, func(t *testing.T) {
			require.True(t, ante.IsAppInjectedMsg(sampleMsg.Msg))
		})
	}
}
