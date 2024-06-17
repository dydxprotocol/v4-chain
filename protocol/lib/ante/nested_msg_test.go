package ante_test

import (
	"fmt"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	bank "github.com/cosmos/cosmos-sdk/x/bank/types"

	appmsgs "github.com/StreamFinance-Protocol/stream-chain/protocol/app/msgs"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/lib"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/lib/ante"
	testmsgs "github.com/StreamFinance-Protocol/stream-chain/protocol/testutil/msgs"

	"github.com/stretchr/testify/require"
)

var (
	invalidInnerMsgErr_AppInjected = fmt.Errorf("Invalid nested msg: app-injected msg type")
	invalidInnerMsgErr_Dydx        = fmt.Errorf("Invalid nested msg for MsgExec: dydx msg type")
)

func TestIsNestedMsg_Empty(t *testing.T) {
	tests := map[string]struct {
		msg sdk.Msg
	}{
		"empty msg": {
			msg: nil,
		},
		"not nested msg": {
			msg: testMsg,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			require.False(t, ante.IsNestedMsg(tc.msg))
		})
	}
}

func TestIsNestedMsg_Invalid(t *testing.T) {
	allMsgsMinusNested := lib.MergeAllMapsMustHaveDistinctKeys(appmsgs.AllowMsgs, appmsgs.DisallowMsgs)
	for key := range appmsgs.NestedMsgSamples {
		delete(allMsgsMinusNested, key)
	}
	allNonNilSampleMsgs := testmsgs.GetNonNilSampleMsgs(allMsgsMinusNested)

	for _, sampleMsg := range allNonNilSampleMsgs {
		t.Run(sampleMsg.Name, func(t *testing.T) {
			require.False(t, ante.IsNestedMsg(sampleMsg.Msg))
		})
	}
}

func TestIsNestedMsg_Valid(t *testing.T) {
	sampleMsgs := testmsgs.GetNonNilSampleMsgs(appmsgs.NestedMsgSamples)
	for _, sampleMsg := range sampleMsgs {
		t.Run(sampleMsg.Name, func(t *testing.T) {
			require.True(t, ante.IsNestedMsg(sampleMsg.Msg))
		})
	}
}

func TestIsDydxMsg_Invalid(t *testing.T) {
	allDydxMsgs := lib.MergeAllMapsMustHaveDistinctKeys(
		appmsgs.AppInjectedMsgSamples,
		appmsgs.NormalMsgsDydxCustom,
		appmsgs.InternalMsgSamplesDydxCustom,
	)
	allMsgsMinusDydx := lib.MergeAllMapsMustHaveDistinctKeys(appmsgs.AllowMsgs, appmsgs.DisallowMsgs)
	for key := range allDydxMsgs {
		delete(allMsgsMinusDydx, key)
	}
	allNonNilSampleMsgs := testmsgs.GetNonNilSampleMsgs(allMsgsMinusDydx)

	for _, sampleMsg := range allNonNilSampleMsgs {
		t.Run(sampleMsg.Name, func(t *testing.T) {
			require.False(t, ante.IsDydxMsg(sampleMsg.Msg))
		})
	}
}

func TestIsDydxMsg_Valid(t *testing.T) {
	allDydxMsgs := lib.MergeAllMapsMustHaveDistinctKeys(
		appmsgs.AppInjectedMsgSamples,
		appmsgs.NormalMsgsDydxCustom,
		appmsgs.InternalMsgSamplesDydxCustom,
	)
	allNonNilSampleMsgs := testmsgs.GetNonNilSampleMsgs(allDydxMsgs)

	for _, sampleMsg := range allNonNilSampleMsgs {
		t.Run(sampleMsg.Name, func(t *testing.T) {
			require.True(t, ante.IsDydxMsg(sampleMsg.Msg))
		})
	}
}

func TestValidateNestedMsg(t *testing.T) {
	tests := map[string]struct {
		msg         sdk.Msg
		expectedErr error
	}{
		"Invalid: not a nested msg": {
			msg:         &bank.MsgSend{},
			expectedErr: fmt.Errorf("not a nested msg"),
		},

		"Invalid MsgExec: app-injected inner msg": {
			msg:         &testmsgs.MsgExecWithAppInjectedInner,
			expectedErr: invalidInnerMsgErr_AppInjected,
		},

		"Invalid MsgExec: dydx custom msg": {
			msg:         &testmsgs.MsgExecWithDydxMessage,
			expectedErr: invalidInnerMsgErr_Dydx,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			err := ante.ValidateNestedMsg(tc.msg)
			require.Equal(t, tc.expectedErr, err)
		})
	}
}
