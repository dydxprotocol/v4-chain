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
	invalidInnerMsgErr_Unsupported = fmt.Errorf("Invalid nested msg: unsupported msg type")
	invalidInnerMsgErr_AppInjected = fmt.Errorf("Invalid nested msg: app-injected msg type")
	invalidInnerMsgErr_Nested      = fmt.Errorf("Invalid nested msg: double-nested msg type")
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

// func TestValidateNestedMsg_IterateEachMsgSample(t *testing.T) {
// 	type addMsgType string
// 	const (
// 		unsupportedMsg addMsgType = "AddUnsupportedMsg"
// 		appInjectedMsg addMsgType = "AddAppInjectedMsg"
// 		nestedMsg      addMsgType = "AddNestedMsg"
// 	)

// 	tests := map[string]struct {
// 		innerMsgs  []sdk.Msg
// 		addMsgType addMsgType

// 		expectedErr error
// 	}{
// 		"Invalid: SingleMsg / AddUnsupportedMsg=true": {
// 			innerMsgs:   []sdk.Msg{},
// 			addMsgType:  unsupportedMsg,
// 			expectedErr: invalidInnerMsgErr_Unsupported,
// 		},
// 		"Invalid: MultiMsgs / AddUnsupportedMsg=true": {
// 			innerMsgs:   []sdk.Msg{&bank.MsgSend{}, &bank.MsgMultiSend{}},
// 			addMsgType:  unsupportedMsg,
// 			expectedErr: invalidInnerMsgErr_Unsupported,
// 		},
// 		"Invalid: SingleMsg / AddAppInjectedMsg=true": {
// 			innerMsgs:   []sdk.Msg{},
// 			addMsgType:  appInjectedMsg,
// 			expectedErr: invalidInnerMsgErr_AppInjected,
// 		},
// 		"Invalid: MultiMsgs / AddAppInjectedMsg=true": {
// 			innerMsgs:   []sdk.Msg{&bank.MsgSend{}, &bank.MsgMultiSend{}},
// 			addMsgType:  appInjectedMsg,
// 			expectedErr: invalidInnerMsgErr_AppInjected,
// 		},
// 		"Invalid: SingleMsg / AddNestedMsg=true": {
// 			innerMsgs:   []sdk.Msg{},
// 			addMsgType:  nestedMsg,
// 			expectedErr: invalidInnerMsgErr_Nested,
// 		},
// 		"Invalid: MultiMsgs / AddNestedMsg=true": {
// 			innerMsgs:   []sdk.Msg{&bank.MsgSend{}, &bank.MsgMultiSend{}},
// 			addMsgType:  nestedMsg,
// 			expectedErr: invalidInnerMsgErr_Nested,
// 		},
// 	}

// 	type testCase struct {
// 		name        string
// 		msgs        []sdk.Msg
// 		expectedErr error
// 	}
// 	allTestCases := make([]testCase, 0, len(tests))

// 	unsupportedCnt := 0
// 	appInjectedCnt := 0
// 	nestedCnt := 0
// 	for tcName, tc := range tests {
// 		var msgSampleTestCase map[string]sdk.Msg

// 		switch tc.addMsgType {
// 		case unsupportedMsg:
// 			unsupportedCnt++
// 			msgSampleTestCase = appmsgs.UnsupportedMsgSamples
// 		case appInjectedMsg:
// 			appInjectedCnt++
// 			msgSampleTestCase = appmsgs.AppInjectedMsgSamples
// 		case nestedMsg:
// 			nestedCnt++
// 			msgSampleTestCase = appmsgs.NestedMsgSamples
// 		default:
// 			panic(fmt.Errorf("unexpected addMsgType: %s", tc.addMsgType))
// 		}

// 		allSampleMsgs := testmsgs.GetNonNilSampleMsgs(msgSampleTestCase)
// 		for _, sampleMsg := range allSampleMsgs {
// 			testName := fmt.Sprintf("%s / %s", tcName, sampleMsg.Name)
// 			testMsgs := append(tc.innerMsgs, sampleMsg.Msg)
// 			require.True(t, len(testMsgs) > 0 && len(testMsgs) <= 3)
// 			allTestCases = append(allTestCases, testCase{testName, testMsgs, tc.expectedErr})
// 		}
// 	}

// 	expectedTotalCnt := 0
// 	expectedTotalCnt += unsupportedCnt * len(testmsgs.GetNonNilSampleMsgs(appmsgs.UnsupportedMsgSamples))
// 	expectedTotalCnt += appInjectedCnt * len(testmsgs.GetNonNilSampleMsgs(appmsgs.AppInjectedMsgSamples))
// 	expectedTotalCnt += nestedCnt * len(testmsgs.GetNonNilSampleMsgs(appmsgs.NestedMsgSamples))
// 	require.Len(t, allTestCases, expectedTotalCnt)

// 	for _, tc := range allTestCases {
// 		t.Run(tc.name, func(t *testing.T) {
// 			nestedMsg, err := gov.NewMsgSubmitProposal(tc.msgs, nil, "", "", "", "", false)
// 			require.NoError(t, err)
// 			result := ante.ValidateNestedMsg(nestedMsg)
// 			require.Equal(t, tc.expectedErr, result)
// 		})
// 	}
// }
